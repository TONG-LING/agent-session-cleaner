<div align="center">
  <h1>agent-session-cleaner</h1>
  <p><strong>Browse, resume, and clean up your Codex and Pi sessions from one terminal.</strong></p>
  <p>
    <a href="https://github.com/haowang02/agent-session-cleaner/releases/latest"><img src="https://img.shields.io/github/v/release/haowang02/agent-session-cleaner?label=release" alt="Latest release"></a>
    <a href="https://github.com/haowang02/agent-session-cleaner/actions/workflows/ci.yml"><img src="https://github.com/haowang02/agent-session-cleaner/actions/workflows/ci.yml/badge.svg" alt="CI status"></a>
    <img src="https://img.shields.io/badge/platforms-macOS%20%7C%20Linux%20%7C%20Windows-blue" alt="Platforms: macOS, Linux, and Windows">
    <a href="./LICENSE"><img src="https://img.shields.io/badge/license-MIT-green" alt="MIT License"></a>
  </p>
  <p><strong>English</strong> · <a href="./README.zh-CN.md">简体中文</a></p>
</div>

![Interface preview](./assets/screenshots/example.png)

## Install

Install or update on macOS and Linux:

```bash
curl -LsSf https://raw.githubusercontent.com/haowang02/agent-session-cleaner/main/install.sh | sh
```

Install or update on Windows from PowerShell:

```powershell
powershell -ExecutionPolicy Bypass -c "irm https://raw.githubusercontent.com/haowang02/agent-session-cleaner/main/install.ps1 | iex"
```

The Windows installer adds `%LOCALAPPDATA%\Programs\agent-session-cleaner\bin` to your user `PATH`. Open a new terminal after the first install.

## Usage

```bash
asc
# The full command name agent-session-cleaner is also available.
```

Run without arguments to choose an agent, or name one directly:

```bash
asc codex
asc pi
```

By default, the app uses `$CODEX_HOME` or `~/.codex` for Codex and `$PI_CODING_AGENT_DIR` or `~/.pi/agent` for Pi. On Windows, `~` is your user profile directory. To override a location explicitly:

```bash
asc --codex-home /path/to/codex
asc --pi-home /path/to/pi/agent
```

Session discovery and previews read the stored data directly. Codex changes are delegated to the Codex CLI; if it is unavailable, Codex opens in browse-only mode. Pi deletions operate directly on its session files.

### Language

The interface follows your locale (`LC_ALL`, `LC_MESSAGES`, `LANGUAGE`, or `LANG`). Chinese locales use Simplified Chinese; all other locales use English. To override detection:

```bash
ASC_LANG=en asc
ASC_LANG=zh-CN asc
```

In PowerShell, use `$env:ASC_LANG = "zh-CN"` before running `asc`.

## Keyboard shortcuts

The footer shows the shortcuts available for the current agent and installation. Press `h` for the complete in-app reference.

| Key | Action | Supported by |
|---|---|---|
| `↑` `↓` / `j` `k` | Move to the previous / next session | All agents |
| `g` / `G` | Jump to the top / bottom of the list | All agents |
| `Tab` | Switch between the session list and conversation | All agents |
| `␣` | Select or deselect the current session | All agents |
| `/` | Search titles, working directories, session IDs, and clients | All agents |
| `?` | Search backward | All agents |
| `n` / `N` | Jump to the next / previous match | All agents |
| `c` | Copy the command that resumes the current session | All agents |
| `y` | Copy the current session's working directory | All agents |
| `d` | Delete the current session or selected sessions | All agents |
| `a` | Archive the current session or selected sessions | Codex only |
| `u` | Unarchive the current session or selected sessions | Codex only |
| `D` | Delete all archived sessions | Codex only |
| `O` | Delete all orphaned sub-agent sessions | All agents |
| `r` | Refresh the session list | All agents |
| `h` | Show keyboard shortcuts | All agents |
| `!` | Toggle danger mode; individual deletions skip confirmation | All agents |
| `f` | Show every session, including the short ones | All agents |
| `Esc` | Exit multi-select, turn off danger mode, or clear the search | All agents |
| `q` | Quit | All agents |

Press Space or double-click to select or deselect a session. Selecting a session includes all of its descendant sub-agent sessions. Actions supported by the current agent apply to every selected session; press `Esc` to clear the selection.

### Which sessions are listed

A session holding fewer than three messages is left out of the list: a test run, or a conversation that was opened and abandoned, is not something to come back to. The banner counts what was left out, and `f` puts it back on screen.

Sessions that are held back are held back from the sweeps as well, so nothing is deleted without having been on screen first. A session whose conversation could not be read is listed as well, since not knowing what is in it is not the same as knowing it is empty.

## Resume a session

To resume the current session, press `c` to copy the command that continues it, then paste it into a shell:

```bash
codex resume <session-id>
pi --session <session-id>
```

The session's own ID is already filled in, so the copied command runs as it stands.

## Deletion and archiving

> [!WARNING]
> This tool does not create backups. Deleted sessions cannot be recovered.

- Deleting or archiving a session also includes its descendant sub-agent sessions.
- Codex archive, unarchive, and delete operations are delegated to the Codex CLI.
- Pi deletion removes the session file directly.
- In danger mode (`!`), individual deletions skip confirmation. Bulk deletion still asks for confirmation.

## Acknowledgements

- [LINUX DO](https://linux.do/) — a community for builders and curious minds
