# Claude Usage Tracker Widget

A lightweight desktop widget for Windows 11 that displays your Claude Code session and weekly usage at a glance. Built in Go for minimal resource consumption (~5-10MB RAM).

![Claude Colors](https://img.shields.io/badge/theme-Claude%20Palette-D4A27C)

## Features

- **Session stats** — tokens (in/out/cache), cost, request count, duration
- **Daily & weekly usage** — aggregated cost with progress bars against your plan limit
- **Transparent widget** — always-on-top, draggable, blends with your desktop
- **Auto-refresh** — updates every 30 seconds (configurable)
- **System tray** — minimize to tray, right-click to quit
- **Claude color palette** — warm tan/beige accent on dark background
- **Two modes** — native WebView2 widget (Windows) or browser-based (cross-platform)

## Requirements

- Go 1.22+ (to build from source)
- Windows 11 (for native widget mode — WebView2 is built-in)
- Claude Code installed (`~/.claude/` directory with session data)

## Quick Start

### Windows (Native Widget)

```powershell
# Clone and build
git clone https://github.com/afranche7/claudeusagetracker.git
cd claudeusagetracker
.\build.bat

# Run
.\claude-usage-tracker.exe
```

### Any Platform (Browser Mode)

```bash
# Build
make build
# or: go build -o claude-usage-tracker ./cmd/widget

# Run (opens in browser)
./claude-usage-tracker --browser
```

## Configuration

### Environment Variables

| Variable | Default | Description |
|---|---|---|
| `CLAUDE_PLAN` | `max_5x` | Your plan: `pro`, `max_5x`, `max_20x` |
| `CLAUDE_HOME` | `~/.claude` | Path to Claude Code config directory |

### Command-Line Flags

```
--browser       Force browser mode instead of native widget
--port 17429    HTTP port for browser mode (default: 17429)
--interval 30   Refresh interval in seconds (default: 30)
```

## How It Works

The widget reads Claude Code's local JSONL session files from `~/.claude/projects/` to extract:

- Token usage (input, output, cache read/write)
- Model information for cost calculation
- Session timestamps for duration and activity detection

All data stays local — no API calls are made.

## Architecture

```
cmd/widget/           → Entry point + platform-specific widget code
internal/tracker/     → JSONL parser, cost calculator, plan config
internal/ui/          → HTML/CSS template + data formatting
```

**Resource usage:** ~5-10MB RAM, negligible CPU (wakes every 30s to scan files).

## Plan Limits

The progress bars show usage as a percentage of your plan budget:

| Plan | Daily | Weekly |
|---|---|---|
| Pro | $20 | $100 |
| Max 5x | $100 | $500 |
| Max 20x | $400 | $2,000 |

Set your plan: `set CLAUDE_PLAN=max_5x` (Windows) or `export CLAUDE_PLAN=max_5x` (Unix).
