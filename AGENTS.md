# Agent Instructions

## Project: zoho-cli

CLI for Zoho REST APIs (CRM, Projects, WorkDrive, Writer, Cliq).

## Tech Stack

- Go 1.25+, urfave/cli v3 (CLI framework), stdlib net/http (HTTP)
- Task runner: mise
- Issue tracker: tk

## Commands

```bash
mise run build         # Build Go binary
mise run test          # Run Go tests
mise run test:all      # Run all Go tests including integration
mise run lint          # Run Go linter (go vet)
mise run fmt           # Format Go code (gofmt)
mise run typecheck     # Go type check (go build)
./zoho --help          # Run CLI
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

- `cmd/zoho/main.go` - Entry point
- `internal/errors.go` - Error types, exit codes, `RequireFlag` helper
- `internal/flags.go` - `MergeJSON`, `MergeJSONForm` helpers
- `internal/auth/` - OAuth flows, token management, config resolution
- `internal/http/` - HTTP client with auto-refresh, `GetClient()`, DC maps
- `internal/dc/` - Datacenter URL config (9 DCs)
- `internal/output/` - JSON output, --help-all schema display
- `internal/crm/` - CRM subcommands (29 commands)
- `internal/projects/` - Projects subcommands (split into ~10 files by resource)
- `internal/books/` - Books subcommands (split into ~16 files by resource)
- `internal/drive/` - WorkDrive subcommands (26 commands)
- `internal/writer/` - Writer subcommands (7 commands)
- `internal/cliq/` - Cliq subcommands (12 commands)

### Reference implementations
- `~/Projects/work/rhi/ai_agent/rhi-agent/src/zoho/` (original endpoints)

## Environment Variables

- `ZOHO_CLIENT_ID`, `ZOHO_CLIENT_SECRET`, `ZOHO_REFRESH_TOKEN`, `ZOHO_DC` - Auth (handled in internal/auth/config.go)
- `ZOHO_PORTAL_ID` - Default for `--portal` flag (Projects commands)
- `ZOHO_TEAM_ID` - Default for `--team` flag (WorkDrive commands)
- `ZOHO_SPRINTS_TEAM_ID` - Default for `--team` flag (Sprints commands)
- `ZOHO_BOOKS_ORG_ID` - Default for `--org` flag (Books, Billing, Invoice, Inventory)
- `ZOHO_EXPENSE_ORG_ID` - Default for `--org` flag (Expense commands)
- `ZOHO_DESK_ORG_ID` - Default for `--org` flag (Desk commands)
- `ZOHO_MAIL_ORG_ID` - Default for `--org` flag (Mail org commands)
- `ZOHO_MAIL_ACCOUNT_ID` - Default for `--account` flag (Mail account commands)
- `ZOHO_CREATOR_OWNER` - Default for `--owner` flag (Creator commands)
- `ZOHO_CREATOR_APP` - Default for `--app` flag (Creator commands)

Flag passed on CLI always overrides the env var. If neither is set, commands fail with a clear error.
Env vars are wired via `Sources: cli.EnvVars(...)` on flag definitions so they appear in `--help`.

## API Documentation (Context7)

- **Always use Context7** to look up Zoho API docs before writing tests or fixing CLI bugs
- Context7 provides up-to-date API documentation with code examples
- Cross-reference CLI endpoint URLs, method names, and parameters against Context7 docs
- The CLI was auto-generated and contains many errors — Context7 is the source of truth

## Conventions

- No comments in code unless asked
- JSON output to stdout by default, errors to stderr
- Exit codes: 0=success, 1=general error, 2=auth error, 3=not found, 4=validation error
- Typed envelope structs for API responses, raw map[string]any for record data
- Pass through raw Zoho API responses (thin wrapper, no data transformation)
- --help-all shows jq-friendly output schemas per command

### Flag patterns (do NOT "fix" these — the distinction is intentional)
- **Query params**: `if v := cmd.String("x"); v != "" { params["x"] = v }` — don't send empty params
- **Body fields**: `if cmd.IsSet("x") { body["x"] = cmd.String("x") }` — allow intentional empty strings
- **Required body fields**: `body["x"] = cmd.String("x")` — no guard needed
- **Typed body fields**: `FloatFlag` + `cmd.Float()`, `IntFlag` + `cmd.Int()`, `BoolFlag` + `cmd.Bool()`
- **JSON array fields**: `StringFlag` + `json.Unmarshal([]byte(cmd.String("x")), &v)` in action
- **--json escape hatch**: `internal.MergeJSON(cmd, body)` — flags always win over --json

## Pagination Policy

All list commands support `--all` and `--limit N` for auto-pagination.

### Flags
- `--all` — fetch all pages automatically
- `--limit N` — fetch up to N total records, paginating as needed
- Default (neither flag) — single request, returns raw Zoho envelope
- `--all`/`--limit` mode returns a flat JSON array of items via `output.JSON()`

### Architecture (`internal/pagination/pagination.go`)
One generic `Paginate(PaginationConfig)` function with composable `SetPage`/`HasMore` callback pairs per product pattern:

| Pattern | Products | SetPage | HasMore | ItemsKey |
|---------|----------|---------|---------|----------|
| A: page/per_page + page_context | Books, Billing, Invoice, Inventory, Expense | `PagePerPage(200)` | `HasMoreBooks` | varies per resource |
| B: from/limit offset | Desk | `FromLimit(100)` | `HasMoreByCount` | `"data"` |
| C: page + page_token | CRM, Bigin, Recruit | `SetPageCRM` | `HasMoreCRM` | `"data"` |
| D: page + page_info | Projects | `PagePerPage(100)` | `HasMoreProjects` | varies per resource |
| E: page[offset]/page[limit] | WorkDrive | `PageOffsetLimit(50)` | `HasMoreWorkDrive` | `"data"` |
| F: index/range | Sprints | `IndexRange(100)` | `HasMoreByCount` | varies |
| G: JSON data param | Sign | `SignPageContext(100)` | `HasMoreSign` | varies |
| H: sIndex/limit | People | `SIndexLimit(200)` | `HasMoreByCount` | `"response.result"` |

### Adding pagination to a new list command
```go
if cmd.Bool("all") || cmd.IsSet("limit") {
    items, err := pagination.Paginate(pagination.PaginationConfig{
        Client:   c,
        URL:      url,
        Opts:     &zohttp.RequestOpts{Params: params},
        ItemsKey: "items_key_here",
        PageSize: 200,
        Limit:    cmd.Int("limit"),
        SetPage:  pagination.PagePerPage(200),
        HasMore:  pagination.HasMoreBooks,
    })
    if err != nil { return err }
    return output.JSON(items)
}
// default single-page fallback
raw, err := c.Request("GET", url, &zohttp.RequestOpts{Params: params})
if err != nil { return err }
return output.JSONRaw(raw)
```

## Key Zoho API Quirks (absorb internally)

- CRM v8 requires `fields` param on list/related/notes/attachments endpoints
- CRM v8 search-global uses `searchword` param (not `word`)
- CRM v8 tags add/remove use JSON body (not query params)
- CRM pagination uses page_token for >2000 records
- Projects pagination: has_next_page can be string "true" or bool true
- WorkDrive uses JSON:API content-type (application/vnd.api+json)
- WorkDrive copy has reversed semantics: POST to destination with source in body
- WorkDrive file status codes: 1=active, 51=trash, 61=delete
- Writer R3002: empty documents cannot be exported
- Download endpoint: use workdrive.zoho.com/api/v1/download/{id} (not download.zoho.com)
- Zoho rate-limits: 10 access token refreshes per refresh_token per 10 minutes
- Go net/http needs explicit Accept: */* header (WorkDrive returns 415 without it)
