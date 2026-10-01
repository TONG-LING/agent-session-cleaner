package tui

import (
	"errors"
	"strings"
	"testing"

	"github.com/haowang02/agent-session-cleaner/internal/session"
)

// talk gives a session a conversation of the given length.
func talk(f *fake, id string, messages int) {
	said := make([]session.Message, messages)
	for i := range said {
		said[i] = session.Message{Role: session.User, Text: "something"}
	}
	f.talk[id] = said
}

// A session holding a single exchange is the test run or the abandoned start
// that buries the conversations worth returning to.
func TestSessionsTooShortToReturnToAreNotListed(t *testing.T) {
	t.Parallel()

	f := newFake(
		session.Session{ID: "test", Title: "quick test"},
		session.Session{ID: "real", Title: "a long conversation"},
	)
	talk(f, "test", 2)
	talk(f, "real", 3)

	m := startFiltered(t, f)
	contains(t, m, "a long conversation")
	omits(t, m, "quick test")
	contains(t, m, "1 hidden")
}

func TestShowAllIsOneKeyAway(t *testing.T) {
	t.Parallel()

	f := newFake(session.Session{ID: "test", Title: "quick test"})
	talk(f, "test", 1)

	m := startFiltered(t, f)
	omits(t, m, "quick test")

	press(t, m, "f")
	contains(t, m, "quick test")
	contains(t, m, "○ Every session")
	contains(t, m, "Showing every session, including the short ones.")

	press(t, m, "f")
	omits(t, m, "quick test")
	contains(t, m, "1 hidden")
	contains(t, m, "Listing only sessions with at least 3 messages.")
}

// Hiding a session that could not be read would hide the reason it could not
// be read.
func TestASessionThatCannotBeReadIsStillListed(t *testing.T) {
	t.Parallel()

	f := newFake(session.Session{ID: "broken", Title: "unreadable"})
	f.talkErr["broken"] = errors.New("the file could not be read")

	m := startFiltered(t, f)
	contains(t, m, "unreadable")
	omits(t, m, "1 hidden")
}

// What is held back from the listing is held back from the sweeps too: nothing
// is destroyed without having been on screen first.
func TestOnlyTheListedOrphansAreSwept(t *testing.T) {
	t.Parallel()

	// Both lost the conversation that spawned them, and only the longer one
	// has enough of a conversation of its own to be listed.
	f := newFake(
		session.Session{ID: "short", Title: "short orphan", Parent: "missing"},
		session.Session{ID: "long", Title: "long orphan", Parent: "missing"},
	)
	talk(f, "short", 1)
	talk(f, "long", 4)

	m := startFiltered(t, f)
	omits(t, m, "short orphan")

	press(t, m, "O", "y")
	if got := strings.Join(f.calls(), ","); got != "long" {
		t.Errorf("swept %q, want only the orphan the listing showed", got)
	}
}

// A key that cannot do anything is not mentioned, and nothing is being held
// back when everything is already listed.
func TestNoFilterIsOfferedWhenNothingIsHidden(t *testing.T) {
	t.Parallel()

	f := newFake(session.Session{ID: "real", Title: "a long conversation"})
	talk(f, "real", 3)

	m := startFiltered(t, f)
	contains(t, m, "a long conversation")
	omits(t, m, "f Show all sessions")
}
