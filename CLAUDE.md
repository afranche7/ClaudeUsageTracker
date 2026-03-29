# CLAUDE.md — AI Assistant Guide for ClaudeUsageTracker

This file provides context and conventions for AI assistants working in this repository.

---

## Project Overview

**ClaudeUsageTracker** is a lightweight transparent Windows 11 desktop widget that displays Claude Code token usage and estimated costs in real time. It reads Claude Code's local JSONL session files — no API calls required.

### What it shows
- **Session**: token count + estimated cost for the current active Claude Code session
- **Weekly**: token count + estimated cost across all sessions in the past 7 days

---

## Repository Structure

```
ClaudeUsageTracker/
├── main.py             # Entry point — run this to launch the widget
├── pyproject.toml      # Project metadata and dependencies (managed by uv)
├── src/
│   ├── __init__.py
│   ├── widget.py       # PyQt6 transparent window, drag, paint, timer
│   ├── tracker.py      # JSONL parser and session/weekly aggregation
│   └── pricing.py      # Model pricing table and cost calculator
├── README.md
└── CLAUDE.md           # This file
```

---

## Tech Stack

| Layer | Choice |
|---|---|
| Language | Python 3.11+ |
| UI framework | PyQt6 6.7+ |
| Package manager | uv |
| Data source | `~/.claude/projects/**/*.jsonl` (Claude Code local session files) |
| Testing | pytest |

No database, no API calls, no external services.

---

## Development Setup

### Prerequisites

- Python >= 3.11
- [uv](https://docs.astral.sh/uv/) package manager
- Claude Code installed (provides the JSONL data files this widget reads)

### Getting Started

```bash
# Clone the repo
git clone <repo-url>
cd ClaudeUsageTracker

# Install dependencies
uv sync

# Run the widget
uv run python main.py
```

No `.env` file or API keys needed — all data is read locally.

---

## How Data is Read

Claude Code writes one JSONL file per session under `~/.claude/projects/<project-hash>/`. Each line is a JSON object. Assistant response lines include a `usage` field:

```json
{
  "message": {
    "role": "assistant",
    "model": "claude-sonnet-4-6",
    "usage": {
      "input_tokens": 123,
      "output_tokens": 456,
      "cache_creation_input_tokens": 789,
      "cache_read_input_tokens": 0
    }
  },
  "sessionId": "abc-123",
  "timestamp": "2026-03-28T15:00:00.000Z"
}
```

`src/tracker.py` globs all `*.jsonl` files, filters for assistant messages with usage data, and aggregates by session ID or by timestamp window.

Current session is identified by reading `~/.claude/sessions/*.json` and finding the most recently started entry.

---

## Pricing Table

Costs are estimated in `src/pricing.py` using a static pricing dict keyed by model tier (`"sonnet"`, `"opus"`, `"haiku"`). Update this file when Anthropic adjusts rates — do not hardcode prices elsewhere.

Current rates (USD per 1M tokens):

| Model | Input | Output | Cache write | Cache read |
|---|---|---|---|---|
| Opus 4 | $15.00 | $75.00 | $18.75 | $1.50 |
| Sonnet 4 | $3.00 | $15.00 | $3.75 | $0.30 |
| Haiku 4 | $0.80 | $4.00 | $1.00 | $0.08 |

---

## Git Workflow

### Branches

| Pattern | Purpose |
|---|---|
| `main` | Stable, production-ready code |
| `feat/<description>` | New features |
| `fix/<description>` | Bug fixes |
| `chore/<description>` | Non-functional changes |
| `claude/<description>` | AI-assisted development branches |

### Commit Style

Follow [Conventional Commits](https://www.conventionalcommits.org/):

```
feat: add daily cost summary view
fix: correct cache read token cost calculation
chore: update PyQt6 to 6.8.0
docs: update pricing table in CLAUDE.md
test: add unit tests for format_tokens helper
```

---

## Testing

```bash
# Run all tests
pytest

# Run with coverage
pytest --cov=src
```

- Tests live in `tests/` mirroring `src/`
- Unit test `src/pricing.py` and `src/tracker.py` with fixture JSONL files
- Do not test `src/widget.py` with a real Qt display — use headless mocks
- Never read from a real `~/.claude/` directory in tests — use tmp fixtures

---

## Key Conventions for AI Assistants

### Do

- Read source files before modifying them
- Keep the pricing table in `src/pricing.py` — not scattered through the codebase
- Update this `CLAUDE.md` when structure or conventions change
- Keep the widget lightweight: no extra dependencies, no threads, no disk writes

### Do Not

- Do not add a database — the widget reads files directly, that's intentional
- Do not make Anthropic API calls — all data is local
- Do not modify `main` directly — always branch and PR
- Do not add dependencies beyond `PyQt6` without strong justification

---

## Useful Commands

```bash
uv run python main.py   # Launch the widget
uv sync                 # Install dependencies
uv run pytest            # Run tests
```

---

## Security

- No secrets or API keys are used or needed
- The widget only reads files from `~/.claude/` — it never writes to them
- No network calls are made at runtime

---

*Last updated: 2026-03-28*
