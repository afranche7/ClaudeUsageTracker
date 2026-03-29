# Claude Usage Tracker

A lightweight transparent Windows 11 desktop widget that shows your Claude Code session and weekly token usage + costs in real time.

## Screenshot

```
┌─────────────────────────────────┐
│  ◆  Claude Usage          [×]   │
├────────────────┬────────────────┤
│   SESSION      │    WEEKLY      │
│                │                │
│   2.3k         │   45.1k        │
│   tokens       │   tokens       │
│   $0.041       │   $0.612       │
│                │                │
│   updated 12s ago               │
└─────────────────────────────────┘
```

## Requirements

- Python 3.11+
- Windows 11 (for Acrylic transparency; works on Windows 10 and macOS/Linux with solid background fallback)
- Claude Code installed and in use (data sourced from `~/.claude/projects/`)

## Setup

```bash
# Install dependencies
pip install -r requirements.txt

# Run the widget
python main.py
```

To auto-start with Windows, add a shortcut to `python main.py` in your Startup folder (`Win+R` → `shell:startup`).

## Usage

| Action | Result |
|---|---|
| Drag anywhere | Move widget |
| Click `×` | Close |
| Right-click | Context menu |
| Right-click → Refresh now | Force immediate data reload |
| Right-click → Reset session | Clear current session display |
| Right-click → Exit | Quit |

The widget refreshes automatically every 30 seconds.

## Data Source

Reads JSONL session logs from `~/.claude/projects/**/*.jsonl` — the same files Claude Code writes locally. No Anthropic API calls are made; all data is local.

Token costs are estimated using the pricing table in `src/pricing.py`. Update that file if Anthropic adjusts rates.

## Project Structure

```
ClaudeUsageTracker/
├── main.py             # Entry point
├── requirements.txt
├── src/
│   ├── widget.py       # PyQt6 transparent window
│   ├── tracker.py      # JSONL parsing and aggregation
│   └── pricing.py      # Model pricing table
└── README.md
```
