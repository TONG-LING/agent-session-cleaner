package tui

import (
	"context"
	"sync"

	tea "charm.land/bubbletea/v2"

	"github.com/haowang02/agent-session-cleaner/internal/agent"
	"github.com/haowang02/agent-session-cleaner/internal/i18n"
	"github.com/haowang02/agent-session-cleaner/internal/session"
)

// listedMessages is how many messages a session needs before the listing shows
// it. A session holding a single exchange or none at all is a test run or a
// conversation that never started, and on a large store those outnumber the
// conversations worth returning to.
const listedMessages = 3

// countWorkers is how many conversations are read at once. Sessions are read
// independently, and eight of them measured a third of the time one at a time
// took over a store of five hundred, with nothing to gain past that.
const countWorkers = 8

// countMessages reads every session to learn how much was actually said in it,
// which is what the listing filters on. A session that could not be read is
// left out of the counts rather than counted as empty, so the listing keeps
// showing it: hiding a row would hide the reason to look at it.
func countMessages(ctx context.Context, target agent.Agent, found []session.Session) map[cacheKey]int {
	var (
		mu      sync.Mutex
		counts  = make(map[cacheKey]int, len(found))
		limit   = make(chan struct{}, countWorkers)
		running sync.WaitGroup
	)
	for _, s := range found {
		running.Add(1)
		go func(s session.Session) {
			defer running.Done()
			limit <- struct{}{}
			messages, err := target.Messages(ctx, s)
			<-limit
			if err != nil {
				return
			}
			mu.Lock()
			counts[keyOf(s)] = len(messages)
			mu.Unlock()
		}(s)
	}
	running.Wait()
	return counts
}

// listed reports whether a session has earned a row in the listing. A session
// whose conversation was never counted stays: not knowing what is in it is not
// the same as knowing it is empty.
func (m *Model) listed(s session.Session) bool {
	if m.showAll {
		return true
	}
	count, counted := m.counts[keyOf(s)]
	return !counted || count >= listedMessages
}

// visibleRows is the listing the filter leaves, in the order the tree gave it.
func (m *Model) visibleRows() []session.Row {
	rows := m.forest.Rows()
	kept := make([]session.Row, 0, len(rows))
	for _, row := range rows {
		if m.listed(row.Session) {
			kept = append(kept, row)
		}
	}
	return kept
}

// hiddenCount is how many sessions the filter is keeping out of the listing.
func (m *Model) hiddenCount() int { return m.forest.Len() - len(m.rows) }

// visibleOrphans are the stranded sub-agent sessions the listing is showing.
// A sweep works on what is on screen, so sessions the filter is holding back
// are never swept away unseen.
func (m *Model) visibleOrphans() []session.Session {
	return keepIf(m.forest.Orphans(), m.listed)
}

// toggleFilter shows every session, or goes back to leaving out the short ones.
func (m *Model) toggleFilter() tea.Cmd {
	m.showAll = !m.showAll
	m.rows = m.visibleRows()
	m.cursor = max(0, min(m.cursor, len(m.rows)-1))
	m.clampView()

	if m.showAll {
		m.say(i18n.FilterOff)
	} else {
		m.say(i18n.FilterOn, i18n.Args{"count": listedMessages})
	}
	return m.showCurrent()
}
