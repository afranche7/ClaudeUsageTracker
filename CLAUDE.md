# CLAUDE.md — AI Assistant Guide for ClaudeUsageTracker

This file provides context and conventions for AI assistants (Claude Code and others) working in this repository.

---

## Project Overview

**ClaudeUsageTracker** is a lightweight desktop widget for Windows 11 that displays Claude Code session, daily, and weekly usage at a glance. It reads Claude Code's local JSONL session files to show token consumption, costs, and usage patterns — enabling developers to monitor their Claude API spend in real time.

### Core Goals

- Parse Claude Code's local session data (`~/.claude/projects/*.jsonl`)
- Calculate costs per model with cache token support
- Display session, daily, and weekly usage with progress bars against plan limits
- Run as an always-on-top transparent widget with minimal resource usage (~8MB binary, ~5-10MB RAM)

---

## Repository Structure

```
ClaudeUsageTracker/
├── cmd/widget/               # Application entry point
│   ├── main.go               # HTTP server, browser mode, CLI flags
│   ├── native_windows.go     # Win32 Edge app-mode launcher
│   ├── native_other.go       # No-op stub for non-Windows
│   ├── tray_windows.go       # System tray (Windows only)
│   └── tray_other.go         # No-op stub for non-Windows
├── internal/
│   ├── tracker/              # Usage data parsing and cost calculation
│   │   ├── config.go         # Model pricing, plan limits, paths
│   │   ├── parser.go         # JSONL session file parser
│   │   └── parser_test.go    # Tests for parser and cost calculation
│   └── ui/                   # Widget UI rendering
│       ├── template.go       # HTML/CSS template (Claude color palette)
│       ├── format.go         # Token/cost/duration formatters + data builder
│       └── format_test.go    # Tests for formatters and HTML rendering
├── assets/                   # Static assets (icons, etc.)
├── build.bat                 # Windows build script
├── build.sh                  # Linux/macOS build script
├── Makefile                  # Cross-platform build targets
├── go.mod                    # Go module definition
├── go.sum                    # Dependency checksums
├── README.md                 # User-facing documentation
└── CLAUDE.md                 # This file
```

---

## Tech Stack

| Layer | Choice |
|---|---|
| Language | Go 1.22+ |
| UI | HTML/CSS rendered via embedded Go templates |
| Widget runtime | HTTP server + browser, or Edge app-mode (Windows) |
| System tray | `github.com/getlantern/systray` (Windows only) |
| Testing | `go test` (stdlib) |
| Build | `go build` with `-ldflags="-s -w"` for small binaries |

---

## Development Setup

### Prerequisites

- Go >= 1.22
- Claude Code installed (creates `~/.claude/` directory with session data)

### Getting Started

```bash
# Clone the repo
git clone https://github.com/afranche7/ClaudeUsageTracker.git
cd ClaudeUsageTracker

# Build
make build              # Linux/macOS
# or
.\build.bat             # Windows

# Run
./claude-usage-tracker --browser          # Browser mode (any OS)
.\claude-usage-tracker.exe                # Native mode (Windows)
```

### Environment Variables

| Variable | Required | Default | Description |
|---|---|---|---|
| `CLAUDE_PLAN` | No | `max_5x` | Plan tier: `pro`, `max_5x`, `max_20x` |
| `CLAUDE_HOME` | No | `~/.claude` | Path to Claude Code config directory |

### CLI Flags

| Flag | Default | Description |
|---|---|---|
| `--browser` | `false` | Force browser mode instead of native widget |
| `--port` | `17429` | HTTP port for browser mode |
| `--interval` | `30` | Refresh interval in seconds |

---

## Git Workflow

### Branches

| Branch pattern | Purpose |
|---|---|
| `main` | Stable, production-ready code |
| `feat/<description>` | New features |
| `fix/<description>` | Bug fixes |
| `chore/<description>` | Non-functional changes (deps, config) |
| `claude/<description>` | AI-assisted development branches |

### Commit Style

Follow [Conventional Commits](https://www.conventionalcommits.org/):

```
feat: add daily cost summary command
fix: correct token count for streaming responses
chore: update anthropic sdk to v0.30.0
docs: update setup instructions in README
test: add unit tests for cost calculator
```

- Keep commits small and focused
- Reference issue numbers where applicable: `feat: add export (#42)`

---

## Testing

```bash
# Run all tests
go test ./...

# Run with verbose output
go test -v ./...

# Run a specific package
go test ./internal/tracker/...
go test ./internal/ui/...
```

- Tests live alongside source files (`*_test.go`)
- Unit tests for cost calculation, formatting, timestamp parsing, and HTML rendering
- No external API calls in tests — all data is parsed from local files

---

## Key Conventions for AI Assistants

### Do

- Read relevant source files before modifying them
- Follow Go conventions: `gofmt`, exported names for public API, unexported for internal
- Keep changes focused — one concern per PR
- Add tests for new logic
- Update this `CLAUDE.md` when the project structure or conventions change

### Do Not

- Do not commit `.env` files or API keys
- Do not make real Anthropic API calls — all data comes from local JSONL files
- Do not add speculative abstractions or features not in scope
- Do not install new dependencies without justification

### Architecture Notes

- **`internal/tracker/`** — Reads `~/.claude/projects/**/*.jsonl` files, extracts token usage from JSON entries, calculates costs using the pricing table in `config.go`
- **`internal/ui/`** — Go `html/template` renders a self-contained HTML/CSS widget; `format.go` provides helpers for tokens (e.g. "1.5M"), costs, and durations
- **`cmd/widget/`** — Entry point with two modes: browser (HTTP server at localhost) and native Windows (Edge `--app` mode with Win32 always-on-top/transparency)
- **Platform files** use Go build tags (`//go:build windows` / `//go:build !windows`) so the project cross-compiles cleanly with `CGO_ENABLED=0`

### Model Pricing

Pricing is defined in `internal/tracker/config.go`. When models or prices change, update the `PricingTable` map. Cache read tokens are priced at 10% of input; cache write tokens at 125% of input.

---

## Useful Commands

```bash
make build            # Build for current platform
make build-windows    # Cross-compile for Windows (CGO_ENABLED=0)
make run              # Build and run in browser mode
make test             # Run all tests
make clean            # Remove build artifacts
```

---

## Security

- **Never commit secrets** — use environment variables
- The widget only reads local files from `~/.claude/` — no network calls are made
- The HTTP server binds to `127.0.0.1` only (not exposed to network)

---

*Last updated: 2026-03-28*
