# zoho-cli: Comprehensive Project Prompt

## What This Is

Build `zoho-cli`, an open-source command-line interface for Zoho's REST APIs. It covers CRM, Projects, WorkDrive, Writer, and Cliq — the five products that have no existing CLI tooling. This CLI will be used both by humans directly and by LLM agents operating inside code execution sandboxes.

There is **no existing Zoho CLI** that covers multiple products. Research (Feb 2026) found:
- One abandoned Java CLI for Zoho Projects only (0 stars, 6 commits)
- Zoho's own CLIs (ZET, Catalyst CLI) are for building extensions/serverless functions, NOT for API operations
- Zoho's official Python SDKs exist only for CRM, are poorly maintained, require MySQL for token persistence, and have incomplete type hints
- No SDK exists for WorkDrive, Writer, or Cliq
- Multiple MCP servers exist but are product-specific, immature, and don't solve the CLI use case
- The `zoho-crm-connector` PyPI package (community) is the only clean alternative, but covers CRM only

## Why This Exists

### The LLM Agent Problem

We operate a multi-agent LLM system (MIAIE) that integrates with Zoho via 50+ individual tool functions. This approach has critical failure modes:

**Tool spam**: Models call the same tool dozens of times in parallel instead of composing operations. We observed gpt-5.2-codex calling `list_project_timelogs` 40+ times simultaneously in a single turn — once per project — because it has no concept of iteration.

**No composability**: Native LLM tools are atomic. There's no way to pipe, filter, loop, or aggregate. "Find all overdue tasks across all projects" requires the LLM to manually orchestrate N+1 tool calls instead of writing `for project in $(zoho projects list --format ids); do zoho projects tasks list --project $project --filter overdue; done`.

**Schema quirks**: LLMs struggle with freeform `dict` parameters. We had to convert all `values: dict` tool params to `values: str` (JSON strings) because `pydantic-ai` generates `{"type": "object", "properties": {}, "additionalProperties": true}` which causes gpt-5.2-codex and gemini-2.5-flash to send empty objects `{}` every time. Only `type: string` works across all model families.

**Quirk multiplication**: Each of Zoho's 50+ endpoints has quirks (wrong field names, broken pagination, undocumented required params). Teaching these to an LLM via tool descriptions is fragile. A CLI absorbs the quirks once, permanently.

### The CLI Solution

Replace all native LLM tools with a CLI installed in the agent's Docker sandbox. The LLM writes shell scripts and Python that call the CLI, enabling:
- `for` loops, pipes, `grep`, `jq`, Python scripts
- Pagination handled transparently by the CLI
- All API quirks absorbed by the CLI implementation
- The LLM uses its strongest capability (writing code) instead of its weakest (calling 50 tools correctly)
- Composability: `zoho crm records search Deals --criteria "(Stage:equals:Closed Won)" | jq '.[].Deal_Name'`
- Auditable: you can see exactly what the agent ran and replay it

### The FOSS Opportunity

There is literally nothing in this space. A well-built Zoho CLI would be the first of its kind. The Zoho ecosystem has 90M+ users across 55+ apps. A CLI that "just works" across CRM, Projects, WorkDrive, Writer, and Cliq would fill a massive gap.

## Technical Decisions

### CLI Framework: TBD (Cappa vs Typer)

Two strong candidates. Final choice to be made during initial project setup.

**Cappa** (188 stars, 95 releases, Apache-2.0):
- Clap-inspired (Rust), dataclass/pydantic-native
- `Dep` injection system like FastAPI's `Depends` — ideal for injecting auth clients
- `parse` returns typed dataclass instances, not raw `Namespace`
- Methods as subcommands: class fields = arguments, methods = subcommands
- Clean testing: no `CliRunner`, just call the function
- Supports both `parse` (return parsed object) and `invoke` (call function) modes
- From the docs: "Cappa inverts the relationship to invoked functions. You do not wrap the functions, so they end up just being plain functions you can call in any scenario."

**Typer** (16k+ stars):
- Battle-tested, massive ecosystem, every LLM knows it
- Built on Click, inherits Click's drawbacks (wraps functions, changes behavior, complex testing via `CliRunner`)
- More verbose for deeply nested subcommands
- From Cappa's comparison: "Click imperitively accumulates the CLI's shape based on where/how you combine groups and commands. Intuiting what resultant CLI look like from the series of click commands strung together is not obvious."

The `--help-all` requirement (recursively show all subcommands) is not built into either framework but is straightforward to implement as a custom flag that walks the command tree.

### Package Name: `zoho-cli`

Published to PyPI as `zoho-cli`. The command name is `zoho`. Installation via `pip install zoho-cli` or `uv pip install zoho-cli` or `uvx zoho-cli`.

### Language & Tooling

- Python 3.12+
- `uv` for package management (consistent with our existing stack)
- `ruff` for linting and formatting (line-length=100)
- `httpx` for HTTP client (async internally, sync CLI interface)
- `pytest` for testing
- Repository: `~/Projects/work/rhi/ai_agent/zoho-cli` (separate repo from rhi-agent)

### Output Format

**JSON by default**, with `--format table` for human-readable output.

Rationale: The primary consumer is LLM agents piping to `jq` and Python. JSON is universally parseable. Table mode is a nice-to-have for interactive human use.

All commands output to stdout. Errors and progress messages go to stderr. This enables clean piping: `zoho crm records list Leads | jq '.[] | .Email'`.

### Pagination

**First page by default**. Options for more:
- `--all` — auto-paginate, fetch everything
- `--page N` — fetch specific page
- `--page 3-15` — fetch page range

Rationale: Default to first page to avoid accidentally fetching 10,000 records. `--all` is explicit opt-in. Page ranges enable targeted fetching without downloading everything.

The CLI handles all pagination mechanics internally (different for each Zoho product — see Zoho API Quirks section). The user never thinks about `page[offset]`, `page`, `per_page`, `index`, `range`, `has_next_page`, or any of that.

### Authentication

#### Token Interchangeability (Critical Insight)

All Zoho OAuth grant types produce **identical tokens**. A refresh token obtained via device flow, authorization code flow, self-client flow, or mobile PKCE flow all:
- Use the same `POST {accounts-server}/oauth/v2/token` endpoint to refresh
- Produce the same `access_token` format
- Use the same `Authorization: Zoho-oauthtoken {token}` header
- Have the same 1-hour access token lifetime
- Are subject to the same rate limits (10 access tokens per refresh token per 10 minutes, max 20 refresh tokens per client per user)

This means the CLI doesn't care how a token was obtained. It just needs a refresh token, client ID, client secret, and data center.

From the Zoho docs on token refresh:
```
POST {accounts-server-url}/oauth/v2/token
  ?grant_type=refresh_token
  &client_id=...
  &client_secret=...
  &refresh_token=...

Response:
{
  "access_token": "1004.ce70f...",
  "api_domain": "https://www.zohoapis.com",
  "token_type": "Bearer",
  "expires_in": 3600
}
```

#### Auth Methods (All Produce the Same Tokens)

**1. Device Flow (CLI-native, primary method)**

From Zoho docs (`/oauth/v3/device/code`):
```
zoho auth login
  1. POST accounts.zoho.com/oauth/v3/device/code
     params: client_id, grant_type=device_request, scope, access_type=offline
  2. Response: user_code, device_code, verification_url, interval, expires_in
  3. Display to user: "Go to {verification_url} and enter code: {user_code}"
     (or provide verification_uri_complete for direct link)
  4. Poll POST accounts.zoho.com/oauth/v3/device/token
     params: client_id, client_secret, grant_type=device_token, code={device_code}
     Poll every {interval}ms until: authorization_pending → success/denied/expired
  5. On success: store refresh_token + access_token + api_domain
```

This is the `gh auth login` pattern. Works in SSH, headless, Docker. No browser redirect server needed.

Polling responses include:
- `slow_down` — polling too fast
- `authorization_pending` — user hasn't approved yet
- `other_dc` — user is on a different DC, redirect polling there
- `access_denied` — user denied
- `expired` — 5 minute timeout

**2. Self-Client (quick setup, power users)**

Generate a code in Zoho API Console → exchange for tokens:
```
zoho auth self-client --code "1000.abc..." --server accounts.zoho.com
```

**3. Environment Variables (agent sandbox / CI)**

```
ZOHO_CLIENT_ID=...
ZOHO_CLIENT_SECRET=...
ZOHO_REFRESH_TOKEN=...
ZOHO_DC=com              # data center (com, eu, in, com.au, jp, ca, sa, uk, com.cn)
ZOHO_ACCOUNTS_URL=https://accounts.zoho.com  # optional, derived from DC if not set
```

The refresh token injected here can come from ANY source — the rhi-agent's webhook OAuth flow stores per-user tokens in SQLite; we just read the refresh_token and inject it as an env var into the sandbox. Same token, different transport.

**4. Config File (previously obtained by any method)**

```
~/.config/zoho-cli/config.toml
  [auth]
  client_id = "..."
  client_secret = "..."

~/.config/zoho-cli/tokens.json
  {
    "refresh_token": "1000.7ed4f...",
    "access_token": "1004.ce70f...",
    "access_token_expires_at": "2026-02-21T16:00:00Z",
    "dc": "com",
    "accounts_url": "https://accounts.zoho.com",
    "api_domain": "https://www.zohoapis.com"
  }
```

#### Auth Resolution Order

```
1. Environment variables (ZOHO_REFRESH_TOKEN etc) — highest priority
2. Config file (~/.config/zoho-cli/)
3. No auth found → "Run `zoho auth login` to authenticate"
```

No named profiles. Single user per installation. Agent use is via env vars per-task.

#### Client ID/Secret Ownership

Users must create their own Zoho OAuth client at https://api-console.zoho.com. The CLI does NOT ship a bundled client ID. Rationale:
- Shared client IDs share rate limits across all users
- Zoho enforces per-client token limits (max 20 refresh tokens per client per user)
- Users may want different scopes for different use cases
- This is what `gh` does (GitHub provides a default app, but Zoho's rate limits make sharing impractical)

#### Scopes

Default to requesting ALL scopes for all supported products unless the user specifies a subset. The full scope string:

```
ZohoCliq.Webhooks.CREATE,ZohoCliq.Channels.ALL,ZohoCliq.Messages.ALL,ZohoCliq.Chats.ALL,
ZohoCliq.Users.ALL,ZohoCliq.Bots.ALL,
WorkDrive.workspace.ALL,WorkDrive.files.ALL,WorkDrive.files.sharing.ALL,WorkDrive.links.ALL,
WorkDrive.team.ALL,WorkDrive.teamfolders.ALL,
ZohoSearch.securesearch.READ,
ZohoWriter.documentEditor.ALL,ZohoPC.files.ALL,
ZohoProjects.portals.ALL,ZohoProjects.projects.ALL,ZohoProjects.tasks.ALL,
ZohoProjects.tasklists.ALL,ZohoProjects.timesheets.ALL,ZohoProjects.bugs.ALL,
ZohoProjects.events.ALL,ZohoProjects.forums.ALL,ZohoProjects.milestones.ALL,
ZohoProjects.documents.ALL,ZohoProjects.users.ALL,
ZohoCRM.modules.ALL,ZohoCRM.settings.ALL,ZohoCRM.users.ALL,ZohoCRM.org.ALL,
ZohoCRM.change_owner.CREATE
```

Users can restrict with `--scopes "ZohoCRM.modules.ALL,ZohoCRM.settings.READ"` during `zoho auth login`.

Future: `--scope-profile readonly` or `--scope-profile crm-only` for predefined subsets.

#### Data Centers

Zoho operates 9 data centers. Each has different base URLs for EVERY service:

| DC | Accounts URL | CRM | Projects | WorkDrive | Writer |
|----|-------------|-----|----------|-----------|--------|
| com | accounts.zoho.com | zohoapis.com | projectsapi.zoho.com | workdrive.zoho.com | www.zohoapis.com/writer |
| eu | accounts.zoho.eu | zohoapis.eu | projectsapi.zoho.eu | workdrive.zoho.eu | www.zohoapis.eu/writer |
| in | accounts.zoho.in | zohoapis.in | projectsapi.zoho.in | workdrive.zoho.in | www.zohoapis.in/writer |
| com.au | accounts.zoho.com.au | zohoapis.com.au | projectsapi.zoho.com.au | workdrive.zoho.com.au | www.zohoapis.com.au/writer |
| jp | accounts.zoho.jp | zohoapis.jp | projectsapi.zoho.jp | workdrive.zoho.jp | www.zohoapis.jp/writer |
| ca | accounts.zohocloud.ca | zohoapis.ca | projectsapi.zohocloud.ca | workdrive.zohocloud.ca | www.zohoapis.ca/writer |
| sa | accounts.zoho.sa | zohoapis.sa | projectsapi.zoho.sa | workdrive.zoho.sa | www.zohoapis.sa/writer |
| uk | accounts.zoho.uk | zohoapis.uk | projectsapi.zoho.uk | workdrive.zoho.uk | www.zohoapis.uk/writer |
| com.cn | accounts.zoho.com.cn | zohoapis.com.cn | projectsapi.zoho.com.cn | workdrive.zoho.com.cn | www.zohoapis.com.cn/writer |

Note the irregularities:
- Canada uses `zohocloud.ca` for accounts/projects/workdrive, but `zohoapis.ca` for CRM/writer
- Writer base URL is under `www.zohoapis.{dc}/writer`, NOT a separate domain
- Projects uses `projectsapi.zoho.{dc}`, NOT `zohoapis.{dc}`
- Download service: `download.zoho.{dc}` (or `download.zohocloud.ca` for Canada)

The device flow's polling response includes `other_dc` if the user's DC differs from the app's DC — the CLI must handle this by redirecting subsequent polling to the correct DC.

#### Token Caching

The CLI caches the access token in the config file (or in memory for env var mode). On each API call:
1. Check if cached access token is still valid (expires_in check with 5-minute buffer)
2. If valid, use it
3. If expired or missing, refresh using the refresh token
4. Cache the new access token

This is identical to the pattern in rhi-agent's `TokenManager` (proactive refresh with 5-minute buffer, reactive refresh on 401).

## Command Structure

```
zoho
├── auth
│   ├── login                           # Device flow OAuth
│   ├── self-client --code CODE         # Self-client token exchange
│   ├── status                          # Show current auth status
│   ├── refresh                         # Force token refresh
│   └── logout                          # Clear stored tokens
├── crm
│   ├── modules
│   │   ├── list                        # List available CRM modules
│   │   └── fields <module>             # List fields for a module
│   ├── records
│   │   ├── list <module>               # List records (with --fields, --page, --all)
│   │   ├── get <module> <id>           # Get single record (with --fields)
│   │   ├── create <module> --json '{}'  # Create record
│   │   ├── update <module> <id> --json '{}'  # Update record
│   │   ├── delete <module> <id>        # Delete record
│   │   └── search <module>             # Search (--word, --email, --phone, --criteria)
│   ├── notes
│   │   ├── list <module> <id>          # List notes on a record
│   │   ├── add <module> <id>           # Add note (--title, --content)
│   │   ├── update <note_id>            # Update note
│   │   └── delete <note_id>            # Delete note
│   ├── related
│   │   └── list <module> <id> <related_list>  # List related records
│   ├── users
│   │   └── list                        # List CRM users
│   └── owner
│       └── change <module> <id> --owner <user_id>  # Change record owner
├── projects
│   ├── list                            # List all projects
│   ├── get <project_id>                # Get single project
│   ├── tasks
│   │   ├── list --project <id>         # List tasks in project
│   │   ├── my                          # List my tasks across all projects
│   │   ├── get --project <id> <task_id>
│   │   ├── create --project <id> --name "..."
│   │   └── update --project <id> <task_id> --json '{}'
│   ├── issues
│   │   ├── list --project <id>
│   │   ├── create --project <id> --name "..."
│   │   └── update --project <id> <issue_id> --json '{}'
│   ├── comments
│   │   ├── list --project <id> --task <task_id>
│   │   └── add --project <id> --task <task_id> --comment "..."
│   ├── tasklists
│   │   └── list --project <id>
│   ├── timelogs
│   │   ├── list --project <id>
│   │   └── add --project <id> --task <id> --hours "1:30" --date 2026-02-21
│   ├── users
│   │   └── list --project <id>
│   └── search --query "..."
├── drive
│   ├── files
│   │   ├── list --folder <id>          # List folder contents
│   │   ├── get <file_id>               # Get file info
│   │   ├── search --query "..."        # Search files
│   │   ├── rename <file_id> --name "..."
│   │   ├── copy <file_id> --to <folder_id>
│   │   ├── trash <file_id>             # Move to trash
│   │   └── delete <file_id>            # Permanent delete
│   ├── folders
│   │   ├── list                        # List team folders
│   │   ├── create --parent <id> --name "..."
│   │   └── breadcrumb <folder_id>      # Show folder path
│   ├── download <file_id> [--output path]
│   ├── upload --folder <id> <file_path>
│   ├── share
│   │   ├── permissions <file_id>       # List permissions
│   │   ├── add <file_id> --email "..." --role viewer|editor
│   │   ├── revoke <file_id> --email "..."
│   │   └── link <file_id>              # Create/get share link
│   └── info <file_id>                  # Alias for files get
├── writer
│   ├── create --folder <id> --name "..."  # Create new Writer doc
│   ├── read <doc_id>                   # Read document content as text
│   └── download <doc_id> --format txt|docx|pdf|html
└── cliq
    ├── channels
    │   ├── list
    │   └── message <channel> --text "..."
    └── chats
        └── message <chat_id> --text "..."
```

This maps 1:1 to our existing 50+ tools:

**CRM (16 tools → 16 commands)**:
list_crm_modules, list_module_fields, list_records, list_all_records, search_records, get_record, create_record, update_record, delete_record, get_related_records, list_record_notes, add_record_note, update_note, delete_note, list_crm_users, change_record_owner

**Projects (17 tools → 17 commands)**:
list_projects, get_project, list_tasks, list_my_tasks, get_task, create_task, update_task, list_task_comments, add_task_comment, list_tasklists, list_project_timelogs, add_timelog, list_issues, create_issue, update_issue, search_projects, list_project_users

**Files (14 tools → 14 commands)**:
search_files, list_team_folders, list_folder_contents, get_file_info, create_file, manage_file, bulk_manage_files, upload_file, read_document, download_file, list_file_permissions, share_file, revoke_file_access, create_share_link

**Cliq (4 commands)**: New — channel/chat messaging

## Zoho API Quirks (CRITICAL — Hard-Won Knowledge)

These quirks were discovered through extensive live API testing (47 integration tests). They are NOT documented in Zoho's official docs or OpenAPI specs. The CLI MUST handle all of them internally so users never encounter them.

### Projects V3 API

- **Base URL**: `https://projectsapi.zoho.{dc}/api/v3` — NOT under `zohoapis.{dc}`
- **ALL endpoints reject trailing slashes** with `URL_RULE_NOT_CONFIGURED`. Use `/projects` not `/projects/`. This applies to EVERY V3 endpoint.
- **Portals/Projects endpoints return bare arrays** — no wrapper key, no pagination metadata. Just `[{...}, {...}]`.
- **Tasks**: wrapper key `"tasks"`, pagination via `page_info.has_next_page` (boolean)
- **Issues**: path is `/issues` (NOT `/bugs`), wrapper key is `"issues"` (NOT `"bugs"`). The V2 path is deprecated.
- **Users**: `page_info.has_next_page` is a STRING `"true"`/`"false"`, not a boolean
- **Timelogs list**: requires `module` query param (JSON object) + `view_type=projectspan`. Path is `/timelogs` (NOT `/logs/`)
- **Timelog add**: body requires `module` as nested object `{"type": "task", "id": "..."}`, NOT flat `"module": "task"`. Also requires `log_name` field. Note: The V3 `POST /log` endpoint returns 500 Internal Server Error regardless of payload — this is a confirmed Zoho bug.
- **Comments**: field name is `comment` (NOT `content`)
- **Issues create/update**: field name is `name` (NOT `title`)
- **My Tasks**: use `/portal/{id}/tasks` (portal-level), NOT `/mytasks` (doesn't exist in V3)
- **Search**: requires `module` and `status` params (use `"all"` for both)
- **Date formats vary**: Tasks=ISO 8601, Phases=MM/DD/YYYY, Timelogs=YYYY-MM-DD

### CRM v8

- **Base URL**: `https://zohoapis.{dc}/crm/v8`
- **`fields` param is effectively required** for `list_records` and `list_all_records` — omitting it returns incomplete or empty results in v8
- **Modules endpoint**: returns `{"modules": [...]}` — NOT `{"data": [...]}`
- **Users endpoint**: returns `{"users": [...]}` — NOT `{"data": [...]}`
- **Create/Update**: use `{"data": [{...}]}` array wrapper, even for single records
- **Search**: requires at least one of `word`, `email`, `phone`, or `criteria`
- **Criteria syntax**: `(Field_Name:operator:value)` with `and`/`or` connectors, max 15 conditions
- **Related records**: `GET /{module}/{record_id}/{related_list}` **requires `fields` param** — 400 without it
- **Notes**: also require `fields` param despite docs marking it optional
- **Change owner**: `POST /{module}/{record_id}/actions/change_owner` with `{"owner": {"id": "..."}}`
- **Pagination**: `page` + `per_page`, max 200 per page, `info.more_records` (bool)

### WorkDrive

- **Base URL**: `https://workdrive.zoho.{dc}/api/v1` — uses JSON:API format
- **Search**: `GET /teams/{team_id}/records` with `search[all]`, `search[name]`, or `search[content]` params (NOT `search[all|name|content]` — that's doc shorthand for three separate params)
- **Search requires `ZohoSearch.securesearch.READ` scope** — without it you get 500 `Invalid OAuth scope`
- **Copy semantics are reversed**: `POST /files/{DESTINATION_FOLDER_ID}/copy` with body `{"data": {"attributes": {"resource_id": "SOURCE_FILE_ID"}}}`. URL = destination, body = source.
- **Upload response uses non-standard fields**: `Permalink` (capital P), `FileName`, `resource_id`, `parent_id`. The `id` field is null — use `attributes.resource_id`.
- **All write operations require `Content-Type: application/vnd.api+json`** with body `{"data": {"type": "files", "attributes": {...}}}`
- **Trash/restore/delete**: status codes in attributes: trash=51, restore=1, permanent delete=61
- **Pagination**: `page[offset]` + `page[limit]` with `meta.has_next`
- **`create_share_link`**: the `request_user_data` attribute is effectively required — omitting it causes a 500 Servlet exception
- **No "read file as text" endpoint** — use Writer API download instead

### Writer

- **Base URL**: `https://www.zohoapis.{dc}/writer/api/v1` — NOT `writer.zoho.{dc}`
- **Document IDs = WorkDrive resource IDs** — same file, two APIs
- **Download as text**: `GET /download/{document_id}?format=txt`
- **Request bodies use `formdata`** (multipart/form-data), NOT JSON
- **Response format is plain JSON**, NOT JSON:API
- **`service_type` for creating Writer docs is `zw`**, NOT `zohowriter`. Sheet=`zohosheet`, Show=`zohoshow`.
- **Empty docs can't be downloaded** — Writer API returns error R3002 for freshly created empty documents

### Cliq

- **Message posting requires `ZohoCliq.Webhooks.CREATE` scope** — NOT `ZohoCliq.Messages.ALL`
- **`response_url`** is single-use, returns 204
- **Bot identity**: `{"bot": {"name": "..."}}` in POST body for channel messages
- **DM responses must use `response_url`** — `POST /chats/{id}/message` silently fails with user OAuth tokens

### Cross-Product

- **Token refresh endpoint**: always `POST {accounts-url}/oauth/v2/token` regardless of grant type used to obtain the token
- **Access tokens expire in 1 hour** (3600 seconds)
- **Rate limit**: max 10 access token refreshes per refresh token per 10 minutes
- **Max 20 active refresh tokens per client per user** — oldest invalidated when exceeded
- **Device code validity**: 5 minutes
- **401 handling**: always retry once after refreshing the token. If `scope_invalid` or `scope_mismatch` in the 401 body, don't retry — the user needs to re-authorize with correct scopes.

## Project Structure

```
zoho-cli/
├── src/zoho_cli/
│   ├── __init__.py
│   ├── main.py                 # Entry point, root command group
│   ├── auth/
│   │   ├── __init__.py
│   │   ├── commands.py         # zoho auth login/logout/status/self-client/refresh
│   │   ├── device_flow.py      # Device authorization flow implementation
│   │   ├── token.py            # Token refresh, caching, validation
│   │   └── config.py           # Config file + env var resolution, auth hierarchy
│   ├── http/
│   │   ├── __init__.py
│   │   ├── client.py           # Shared httpx client, auto-refresh on 401, retry
│   │   └── dc.py               # DC→URL maps (single source of truth)
│   ├── output/
│   │   ├── __init__.py
│   │   ├── json.py             # JSON output (default)
│   │   └── table.py            # Table output (--format table)
│   ├── pagination.py           # Unified pagination (handles all 4 Zoho pagination styles)
│   ├── crm/
│   │   ├── __init__.py
│   │   ├── app.py              # CRM sub-app command group
│   │   ├── records.py          # CRUD, search
│   │   ├── modules.py          # List modules, list fields
│   │   ├── notes.py            # Record notes
│   │   ├── related.py          # Related records
│   │   ├── users.py            # CRM users
│   │   └── owner.py            # Change owner
│   ├── projects/
│   │   ├── __init__.py
│   │   ├── app.py              # Projects sub-app command group
│   │   ├── projects.py         # List/get projects
│   │   ├── tasks.py            # Task CRUD, my tasks
│   │   ├── issues.py           # Issue CRUD
│   │   ├── comments.py         # Task comments
│   │   ├── tasklists.py        # Tasklists
│   │   ├── timelogs.py         # Timelog list/add
│   │   ├── users.py            # Project users
│   │   └── search.py           # Project search
│   ├── drive/
│   │   ├── __init__.py
│   │   ├── app.py              # Drive sub-app command group
│   │   ├── files.py            # File operations (list, get, search, rename, copy, trash, delete)
│   │   ├── folders.py          # Folder operations (list, create, breadcrumb)
│   │   ├── transfer.py         # Upload/download
│   │   └── sharing.py          # Permissions, share, revoke, share links
│   ├── writer/
│   │   ├── __init__.py
│   │   ├── app.py              # Writer sub-app
│   │   └── documents.py        # Create, read, download
│   └── cliq/
│       ├── __init__.py
│       ├── app.py              # Cliq sub-app
│       └── messaging.py        # Channel/chat messaging
├── tests/
│   ├── conftest.py
│   ├── test_auth.py
│   ├── test_crm.py
│   ├── test_projects.py
│   ├── test_drive.py
│   ├── test_writer.py
│   └── test_cliq.py
├── pyproject.toml
├── README.md
├── LICENSE                     # MIT or Apache-2.0
└── AGENTS.md                   # AI agent instructions
```

## Reference Material

The following reference material exists in the rhi-agent repo (`~/Projects/work/rhi/ai_agent/rhi-agent/`) and should be consulted when implementing endpoints:

### OpenAPI Specs & API Docs

| Path | Product | Format | Notes |
|------|---------|--------|-------|
| `ref_repos/crm-oas/v8.0/` | CRM v8 | 106 JSON files (per-endpoint) | Comprehensive but contains inaccuracies |
| `ref_repos/cliq-oas/` | Cliq | 22 YAML files (per-resource) | bots, channels, chats, messages, users, etc. |
| `ref_repos/openapi_specs/workdrive.json` | WorkDrive | Single JSON file | Full OAS |
| `ref_repos/openapi_specs/writer-api-collection.json` | Writer | Postman collection | Not OAS format but complete |
| (no spec) | Projects V3 | N/A | Only live-tested; scraped docs cached locally |

**WARNING**: The OpenAPI specs contain inaccuracies. Use them as a catalog of "what endpoints and params exist" but ALWAYS cross-reference with the battle-tested tool implementations (below) for correct URLs, field names, and quirk handling.

### Battle-Tested Tool Implementations (Source of Truth)

| Path | Product | Tools | Notes |
|------|---------|-------|-------|
| `src/agents/zoho/tools/crm.py` | CRM | 16 tools | All v8 quirks handled |
| `src/agents/zoho/tools/projects.py` | Projects | 17 tools | All V3 quirks handled |
| `src/agents/zoho/tools/files.py` | Files | 14 tools | WorkDrive + Writer |

These implementations have been tested against live Zoho APIs with 47 integration tests (44 passing, 3 xfail for Zoho-side bugs). They contain the correct endpoint URLs, request/response formats, pagination handling, and all quirk workarounds. **Port the logic from these files, not from the OAS specs.**

### Infrastructure Reference

| Path | What | Port? |
|------|------|-------|
| `src/zoho/dc.py` | DC→URL maps for all 9 data centers | Yes — copy directly |
| `src/zoho/auth/oauth.py` | OAuth (authorization code + refresh) | Yes — add device flow |
| `src/zoho/api/session.py` | HTTP client with auto-refresh + 401 retry | Yes — adapt for CLI |
| `src/zoho/api/pagination.py` | Pagination for all Zoho styles | Yes — adapt for CLI |
| `tests/test_integration.py` | 47 live API integration tests | Yes — port as CLI tests |

## Development Approach

### Phase 1: Infrastructure

Build the foundation that every endpoint needs:
1. Project scaffolding (pyproject.toml, src layout, git init)
2. CLI framework setup (Cappa or Typer — decide and wire up root command + `--help`)
3. DC maps (port from `src/zoho/dc.py`)
4. Auth config resolution (env vars → config file → error)
5. Token refresh (port from `src/zoho/auth/oauth.py`)
6. Device flow (`POST /oauth/v3/device/code` → poll → store)
7. Self-client flow (exchange code → store)
8. HTTP client with auto-refresh on 401 (port from `src/zoho/api/session.py`)
9. Pagination helpers (port from `src/zoho/api/pagination.py`)
10. Output formatting (JSON default, `--format table`)
11. Error handling (clean stderr messages, exit codes)
12. `zoho auth login`, `zoho auth status`, `zoho auth logout`, `zoho auth refresh`

Test: `zoho auth login` works end-to-end with a real Zoho account.

### Phase 2: Endpoints (Task Swarm)

Each endpoint is an independent, testable unit of work. The approach:

```
while endpoint := get_next_task():
    implement_endpoint(endpoint)      # Port logic from rhi-agent tools
    write_test(endpoint)              # Unit test with mocked HTTP
    run_test(endpoint)                # Verify it passes
    run_integration_test(endpoint)    # Test against live API if possible
    fix_if_needed(endpoint)           # Iterate until green
    run_lint()                        # ruff check
    mark_complete(endpoint)
```

Use `bd` (or equivalent task tracker) to manage the swarm. Every endpoint gets its own task. Tasks are independent and can be worked in parallel by multiple agents, but within a single session we work them sequentially: pick one, implement it fully (with tests), verify, close, move to next.

#### CRM Endpoints (16 tasks)

| # | Command | Source | Priority |
|---|---------|--------|----------|
| 1 | `crm modules list` | `list_crm_modules` in crm.py | High |
| 2 | `crm modules fields <module>` | `list_module_fields` in crm.py | High |
| 3 | `crm records list <module>` | `list_records` in crm.py | High |
| 4 | `crm records get <module> <id>` | `get_record` in crm.py | High |
| 5 | `crm records create <module>` | `create_record` in crm.py | High |
| 6 | `crm records update <module> <id>` | `update_record` in crm.py | High |
| 7 | `crm records delete <module> <id>` | `delete_record` in crm.py | Medium |
| 8 | `crm records search <module>` | `search_records` in crm.py | High |
| 9 | `crm notes list <module> <id>` | `list_record_notes` in crm.py | Medium |
| 10 | `crm notes add <module> <id>` | `add_record_note` in crm.py | Medium |
| 11 | `crm notes update <note_id>` | `update_note` in crm.py | Low |
| 12 | `crm notes delete <note_id>` | `delete_note` in crm.py | Low |
| 13 | `crm related list <module> <id> <rel>` | `get_related_records` in crm.py | Medium |
| 14 | `crm users list` | `list_crm_users` in crm.py | Medium |
| 15 | `crm owner change <module> <id>` | `change_record_owner` in crm.py | Low |
| 16 | `crm records list-all <module>` | `list_all_records` in crm.py | Medium |

#### Projects Endpoints (17 tasks)

| # | Command | Source | Priority |
|---|---------|--------|----------|
| 17 | `projects list` | `list_projects` in projects.py | High |
| 18 | `projects get <id>` | `get_project` in projects.py | High |
| 19 | `projects tasks list` | `list_tasks` in projects.py | High |
| 20 | `projects tasks my` | `list_my_tasks` in projects.py | Medium |
| 21 | `projects tasks get` | `get_task` in projects.py | High |
| 22 | `projects tasks create` | `create_task` in projects.py | High |
| 23 | `projects tasks update` | `update_task` in projects.py | High |
| 24 | `projects comments list` | `list_task_comments` in projects.py | Medium |
| 25 | `projects comments add` | `add_task_comment` in projects.py | Medium |
| 26 | `projects tasklists list` | `list_tasklists` in projects.py | Medium |
| 27 | `projects timelogs list` | `list_project_timelogs` in projects.py | Medium |
| 28 | `projects timelogs add` | `add_timelog` in projects.py | Low (V3 bug) |
| 29 | `projects issues list` | `list_issues` in projects.py | High |
| 30 | `projects issues create` | `create_issue` in projects.py | High |
| 31 | `projects issues update` | `update_issue` in projects.py | Medium |
| 32 | `projects search` | `search_projects` in projects.py | Medium |
| 33 | `projects users list` | `list_project_users` in projects.py | Medium |

#### Drive Endpoints (14 tasks)

| # | Command | Source | Priority |
|---|---------|--------|----------|
| 34 | `drive folders list` | `list_team_folders` in files.py | High |
| 35 | `drive files list` | `list_folder_contents` in files.py | High |
| 36 | `drive files get <id>` | `get_file_info` in files.py | High |
| 37 | `drive files search` | `search_files` in files.py | High |
| 38 | `drive folders create` | `create_file` (folder mode) in files.py | High |
| 39 | `drive files rename` | `manage_file` (rename) in files.py | Medium |
| 40 | `drive files copy` | `manage_file` (copy) in files.py | Medium |
| 41 | `drive files trash` | `manage_file` (trash) in files.py | Medium |
| 42 | `drive files delete` | `manage_file` (delete) in files.py | Medium |
| 43 | `drive upload` | `upload_file` in files.py | High |
| 44 | `drive download` | `download_file` in files.py | High |
| 45 | `drive share permissions` | `list_file_permissions` in files.py | Medium |
| 46 | `drive share add/revoke` | `share_file`/`revoke_file_access` in files.py | Medium |
| 47 | `drive share link` | `create_share_link` in files.py | Medium |

#### Writer Endpoints (3 tasks)

| # | Command | Source | Priority |
|---|---------|--------|----------|
| 48 | `writer create` | `create_file` (writer mode) in files.py | High |
| 49 | `writer read <doc_id>` | `read_document` in files.py | High |
| 50 | `writer download <doc_id>` | `download_file` (writer format) in files.py | Medium |

#### Cliq Endpoints (4 tasks)

| # | Command | Source | Priority |
|---|---------|--------|----------|
| 51 | `cliq channels list` | New (from OAS) | Medium |
| 52 | `cliq channels message` | New (from OAS) | Medium |
| 53 | `cliq chats message` | New (from OAS) | Medium |
| 54 | `cliq channels message` (with bot identity) | New | Low |

**Total: 54 endpoint tasks + ~12 infrastructure tasks = ~66 tasks**

### Phase 3: Polish

- `--help-all` recursive help
- README with examples
- PyPI publishing setup
- Integration test suite (port from rhi-agent's test_integration.py)

## Integration with rhi-agent (Agent Sandbox Use)

When used inside rhi-agent's Docker sandbox:

1. **Install in sandbox image**: Add `uv pip install zoho-cli` to `docker/sandbox/Dockerfile`

2. **Inject auth per-task**: In `get_sandbox()` (`src/deps.py`), inject per-user credentials:
   ```python
   env_vars["ZOHO_CLIENT_ID"] = settings.zoho_client_id
   env_vars["ZOHO_CLIENT_SECRET"] = settings.zoho_client_secret
   env_vars["ZOHO_REFRESH_TOKEN"] = user_refresh_token  # from DB, per-task
   env_vars["ZOHO_DC"] = user_dc  # from token record
   ```

3. **Collapse subagents**: Replace Files, Projects, CRM subagents with a single "code" subagent that has the CLI. The orchestrator delegates to code, code writes scripts using `zoho` CLI.

4. **Code subagent instructions**: Updated to mention `zoho` CLI availability and common patterns (piping to jq, for loops over projects, etc.)

## Error Handling

- All API errors produce clean, actionable messages on stderr
- Include HTTP status code, Zoho error code (if present), and human-readable message
- On 401: attempt one token refresh and retry. If retry fails, suggest `zoho auth login`
- On scope errors (`scope_invalid`, `scope_mismatch`): tell user which scope is missing
- On rate limit (429 or throttle response): show retry-after duration
- Exit codes: 0=success, 1=general error, 2=auth error, 3=not found, 4=validation error

## Testing Strategy

- **Unit tests**: mock HTTP responses, test CLI argument parsing, test output formatting
- **Integration tests**: live API tests against real Zoho account (marked slow, require env vars)
- **Port existing 47 integration tests** from rhi-agent as starting point
- **Every endpoint gets a test before it's marked complete**
- **Lint on every task close**: `ruff check src/` must pass

Commands:
```bash
uv run pytest tests/                    # All unit tests
uv run pytest tests/ -m slow            # Integration tests only
uv run ruff check src/                  # Lint
uv run ruff format src/                 # Format
```

## Non-Goals (For Now)

- No Zoho Books, Desk, Recruit, or other products (can add later)
- No interactive/TUI mode (just CLI)
- No config wizard (just `zoho auth login` + config file)
- No built-in rate limiter beyond Zoho's (user responsibility)
- No webhook server (that stays in rhi-agent)
- No multi-user/profile support
