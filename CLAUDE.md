# CLAUDE.md — AI Assistant Guide for ClaudeUsageTracker

This file provides context and conventions for AI assistants (Claude Code and others) working in this repository.

---

## Project Overview

**ClaudeUsageTracker** is a tool for monitoring and analyzing usage of the Anthropic Claude API. It tracks token consumption, request counts, associated costs, and usage patterns over time — enabling developers and teams to gain visibility into their Claude API spend and behavior.

### Core Goals

- Aggregate usage data from the Anthropic API (tokens in/out, model used, cost)
- Store and persist usage history locally or in a lightweight database
- Provide summaries and reports (daily, weekly, by model, by project)
- Support CLI and/or web dashboard interfaces

---

## Repository Structure (Planned)

```
ClaudeUsageTracker/
├── src/                  # Application source code
│   ├── api/              # Anthropic API client and data fetching
│   ├── db/               # Database models and migrations
│   ├── cli/              # CLI entry points and commands
│   ├── dashboard/        # Web dashboard (if applicable)
│   └── utils/            # Shared utilities
├── tests/                # Test files mirroring src/ structure
├── scripts/              # Utility and setup scripts
├── .env.example          # Environment variable template
├── package.json          # (or pyproject.toml / go.mod depending on stack)
├── README.md
└── CLAUDE.md             # This file
```

> **Note:** This structure is a starting point. Update this section as the actual structure evolves.

---

## Tech Stack

> To be finalized. Common choices for this type of project:

| Layer | Options |
|---|---|
| Language | TypeScript (Node.js), Python, or Go |
| Database | SQLite (local), PostgreSQL (hosted) |
| CLI framework | `commander` (Node), `click` (Python), `cobra` (Go) |
| HTTP client | `axios` / `node-fetch` (Node), `httpx` (Python) |
| Testing | `jest` (Node), `pytest` (Python) |
| ORM / query | `drizzle-orm` / `prisma` (Node), `SQLAlchemy` (Python) |

When the stack is chosen, update this section with specific versions.

---

## Development Setup

### Prerequisites

- Node.js >= 20 (or Python >= 3.11 / Go >= 1.22 — update when stack is finalized)
- An Anthropic API key with access to usage data

### Getting Started

```bash
# Clone the repo
git clone <repo-url>
cd ClaudeUsageTracker

# Install dependencies (update command for chosen stack)
npm install          # Node.js
# pip install -e .   # Python
# go mod tidy        # Go

# Configure environment
cp .env.example .env
# Edit .env and set ANTHROPIC_API_KEY=sk-ant-...

# Run the app
npm start
```

### Environment Variables

| Variable | Required | Description |
|---|---|---|
| `ANTHROPIC_API_KEY` | Yes | Your Anthropic API key |
| `DATABASE_URL` | No | DB connection string (defaults to local SQLite) |
| `LOG_LEVEL` | No | `debug` \| `info` \| `warn` \| `error` (default: `info`) |

---

## Git Workflow

### Branches

| Branch pattern | Purpose |
|---|---|
| `main` | Stable, production-ready code |
| `develop` | Integration branch for features |
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

### Pull Requests

- Target `main` (or `develop` if used)
- Include a summary of changes and testing steps
- All CI checks must pass before merging

---

## Testing

```bash
# Run all tests
npm test

# Run tests in watch mode
npm run test:watch

# Run with coverage
npm run test:coverage
```

- Tests live in `tests/` and mirror the `src/` directory structure
- Unit tests for pure logic; integration tests for DB and API interactions
- Mock the Anthropic API in tests — never make real API calls in tests

---

## Key Conventions for AI Assistants

### Do

- Read relevant source files before modifying them
- Follow the existing code style (indentation, naming, module structure)
- Keep changes focused — one concern per PR
- Add tests for new logic
- Update this `CLAUDE.md` when the project structure or conventions change significantly

### Do Not

- Do not commit `.env` files or API keys
- Do not make real Anthropic API calls in tests
- Do not add speculative abstractions or features not in scope
- Do not modify `main` directly — always branch and PR
- Do not install new dependencies without justification

### Anthropic API Notes

- The usage data endpoint is `/v1/usage` (or embedded in response objects)
- Costs vary by model — keep a pricing table in config, not hardcoded
- Rate limits apply; implement backoff for batch data fetches
- Token counts come from response `usage` fields: `input_tokens`, `output_tokens`

---

## Useful Commands (update as project grows)

```bash
npm run build       # Compile TypeScript
npm run lint        # Lint source files
npm run format      # Auto-format with prettier
npm run migrate     # Run database migrations
npm run seed        # Seed database with sample data
```

---

## Security

- **Never commit secrets** — use `.env` (gitignored) or a secrets manager
- Validate all external input before storing or processing
- Keep dependencies up to date; run `npm audit` regularly

---

## Updating This File

Keep `CLAUDE.md` current as the project evolves. In particular, update:

- The **Repository Structure** section as directories are added
- The **Tech Stack** section once the stack is finalized
- The **Development Setup** section when install/run steps change
- The **Useful Commands** section as new scripts are added

---

*Last updated: 2026-03-28*
