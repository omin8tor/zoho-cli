# Agent Instructions

## Project: zoho-cli

CLI for Zoho REST APIs (CRM, Projects, WorkDrive, Writer, Cliq).

## Tech Stack

- Python 3.14, Cappa (CLI framework), httpx (HTTP), ruff (lint/fmt), mypy (typecheck), pytest (test)
- Package manager: uv
- Task runner: mise
- Issue tracker: tk

## Commands

```bash
mise run test          # Run tests (excluding slow)
mise run test:all      # Run all tests including slow
mise run lint          # Run linter (ruff check src/)
mise run fmt           # Format code (ruff format src/)
mise run typecheck     # Run type checker (mypy src/)
uv run zoho --help     # Run CLI
```

## Quality Gates (run before commit)

```bash
mise run lint && mise run typecheck && mise run test
```

## Issue Tracking (tk)

```bash
tk ready               # Find available work
tk show <id>           # View issue details
tk start <id>          # Claim work (set in_progress)
tk close <id>          # Complete work
tk ls                  # List all open issues
tk blocked             # Show blocked issues
```

## Architecture

- `src/zoho_cli/` - Main package
  - `main.py` - Entry point, root Cappa command
  - `auth/` - OAuth flows, token management, config resolution
  - `http/` - HTTP client with auto-refresh, DC maps
  - `output/` - JSON and table output formatting
  - `pagination.py` - Unified pagination (CRM, Projects, WorkDrive styles)
  - `crm/`, `projects/`, `drive/`, `writer/`, `cliq/` - Product subcommands
- `tests/` - pytest tests
- Reference implementations: `~/Projects/work/rhi/ai_agent/rhi-agent/src/zoho/` (port from here)

## Conventions

- No comments in code unless asked
- JSON output to stdout by default, errors to stderr
- Exit codes: 0=success, 1=general error, 2=auth error, 3=not found, 4=validation error
- Sync httpx (no asyncio)
- Cappa Dep injection for auth/HTTP client into commands
- ruff line-length=100
