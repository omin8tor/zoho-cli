# zoho-cli

CLI for Zoho's REST APIs. Covers 20 products — CRM, Projects, WorkDrive, Books, Mail, Desk, and 14 more — with 1600+ commands in a single binary. JSON to stdout.

No other tool does this. Zoho's own CLIs (ZET, Catalyst) are for building extensions, not talking to the API. Their Python SDK covers CRM only and requires MySQL for token storage.

## Getting started

You'll need two things: the binary and a Zoho API app. Takes about five minutes.

### Install

From source (Go 1.22+):

```bash
go install github.com/omin8tor/zoho-cli/cmd/zoho@latest
```

Or grab a prebuilt binary from [releases](https://github.com/omin8tor/zoho-cli/releases) — unpack it, put it on your PATH.

### Set up a Zoho API app

Go to [api-console.zoho.com](https://api-console.zoho.com/) and create a "Self Client" app. You'll get a client ID and secret.

### Authenticate

**If you're a human at a terminal**, use device flow:

```bash
zoho auth login --client-id YOUR_ID --client-secret YOUR_SECRET
```

This opens your browser, you approve the scopes, and you're done. Tokens are saved to `~/.config/zoho-cli/config.json` and auto-refresh.

**If you're a script or agent**, use env vars:

```bash
export ZOHO_CLIENT_ID=1000.ABC123
export ZOHO_CLIENT_SECRET=xyz789
export ZOHO_REFRESH_TOKEN=1000.refresh_token_here
export ZOHO_DC=com
```

(You can also use `zoho auth self-client` to exchange a code from the API Console.)

Check that it works:

```bash
zoho auth status
```

## Tutorial: your first queries

Everything below assumes you've authenticated. Output is always JSON to stdout, errors to stderr. Pipe into `jq` for filtering.

### CRM: look up contacts

List contacts, pulling specific fields:

```bash
zoho crm records list Contacts --fields "Full_Name,Email,Phone"
```

Get one contact by ID:

```bash
zoho crm records get Contacts 5551234000000012345 --fields "Full_Name,Email"
```

Search by criteria — find closed-won deals:

```bash
zoho crm records search Deals --criteria "(Stage:equals:Closed Won)" --fields "Deal_Name,Amount"
```

Create a lead:

```bash
zoho crm records create Leads --json '{"Last_Name":"Ochoa","Company":"Acme"}'
```

CRM v8 requires the `--fields` param on most read endpoints. If you forget it, the API returns empty records. That's Zoho, not us.

### CRM: search across everything

```bash
zoho crm search-global --searchword "Ochoa" --fields "Full_Name,Email"
```

### CRM: COQL queries

If criteria-based search isn't flexible enough, use COQL (Zoho's SQL-like query language):

```bash
zoho crm coql --query "SELECT Full_Name, Email FROM Contacts WHERE Email LIKE '%@acme.com' LIMIT 10"
```

This needs the `ZohoCRM.coql.READ` scope — it's separate from the general CRM scopes.

### Projects: find your tasks

You need a portal ID. List your portals first:

```bash
zoho projects core list
```

Grab the portal ID from the output, then:

```bash
zoho projects tasks my --portal 12345
```

List all tasks in a project:

```bash
zoho projects tasks list --portal 12345 --project 67890
```

Filter open tasks with jq:

```bash
zoho projects tasks my --portal 12345 | jq '[.[] | select(.status.name == "Open")]'
```

### WorkDrive: navigate and download

Find your team:

```bash
zoho drive teams me
```

List top-level folders:

```bash
zoho drive folders list --team TEAM_ID
```

List files in a folder:

```bash
zoho drive files list --folder FOLDER_ID
```

Download a file:

```bash
zoho drive download FILE_ID --output ./report.pdf
```

Upload a file:

```bash
zoho drive upload ./quarterly.xlsx --folder FOLDER_ID
```

### Books: invoices and contacts

Every Books command needs `--org` (or set `ZOHO_BOOKS_ORG_ID`):

```bash
zoho books invoices list --org 12345
zoho books contacts list --org 12345
zoho books invoices create --org 12345 --json '{"customer_id":"460000000026049","line_items":[{"item_id":"460000000026065","quantity":1}]}'
```

### Desk: support tickets

```bash
zoho desk tickets list --org 12345
zoho desk tickets get --org 12345 --id 98765
zoho desk contacts list --org 12345
zoho desk search --org 12345 --module tickets --query "printer issue"
```

### Mail: messages and folders

```bash
zoho mail accounts list
zoho mail folders list --account ACCOUNT_ID
zoho mail messages list --account ACCOUNT_ID --folder FOLDER_ID
zoho mail messages get --account ACCOUNT_ID --folder FOLDER_ID --message MESSAGE_ID
```

### People: employee records

```bash
zoho people forms list
zoho people records list --form employee
zoho people attendance list --sdate 2025-01-01 --edate 2025-01-31
```

### Sprints: agile project management

```bash
zoho sprints teams list
zoho sprints projects list --team TEAM_ID
zoho sprints items list --team TEAM_ID --project PROJECT_ID --sprint SPRINT_ID
```

### Expense: track spending

```bash
zoho expense expenses list --org 12345
zoho expense reports list --org 12345
zoho expense categories list --org 12345
```

### Sign: document signing

```bash
zoho sign requests list
zoho sign templates list
zoho sign requests get --id REQUEST_ID
```

### Writer: work with documents

Get document details:

```bash
zoho writer details DOC_ID
```

Merge data into a template and export as PDF:

```bash
zoho writer merge DOC_ID --json '{"name":"Alice","date":"2025-01-15"}' --format pdf --output ./letter.pdf
```

Get doc IDs from WorkDrive — Writer has no "list documents" endpoint.

### Cliq: send messages

List channels:

```bash
zoho cliq channels list
```

Send a DM:

```bash
zoho cliq buddies message someone@company.com --text "quarterly report is ready"
```

## Pagination

All list commands support auto-pagination:

```bash
zoho crm records list Contacts --fields "Email" --all              # fetch everything
zoho books invoices list --org 12345 --limit 50                    # fetch up to 50
zoho desk tickets list --org 12345 --all | jq length               # count all tickets
```

Without `--all` or `--limit`, you get a single page with the raw Zoho envelope. With either flag, you get a flat JSON array.

## Piping and composition

The whole point is composability. Everything is JSON, so chain with `jq`, `xargs`, whatever:

```bash
# Get all emails from contacts
zoho crm records list Contacts --fields "Email" --all | jq -r '.[].Email'

# Download every file in a folder
zoho drive files list --folder FOLDER_ID | jq -r '.[].id' | xargs -I{} zoho drive download {} --output ./downloads/

# Find overdue tasks
zoho projects tasks my --portal 12345 | jq '[.[] | select(.end_date < "2025-01-01")]'

# List all unpaid invoices
zoho books invoices list --org 12345 --all | jq '[.[] | select(.status == "unpaid")]'
```

## Default IDs via environment

If you use the same portal or team repeatedly, set these env vars to skip the flags:

```bash
export ZOHO_PORTAL_ID=12345           # default for --portal (Projects)
export ZOHO_TEAM_ID=abc123            # default for --team (WorkDrive)
export ZOHO_SPRINTS_TEAM_ID=xyz789    # default for --team (Sprints)
export ZOHO_BOOKS_ORG_ID=12345        # default for --org (Books)
export ZOHO_DESK_ORG_ID=12345         # default for --org (Desk)
export ZOHO_EXPENSE_ORG_ID=12345      # default for --org (Expense)
export ZOHO_MAIL_ACCOUNT_ID=12345     # default for --account (Mail)
export ZOHO_MAIL_ORG_ID=12345         # default for --org (Mail admin)
export ZOHO_CREATOR_OWNER=owner       # default for --owner (Creator)
export ZOHO_CREATOR_APP=appname       # default for --app (Creator)
```

The flag always overrides the env var. If neither is set, the command fails with a clear error.

## Data centers

Zoho runs in 9 data centers. Set via `ZOHO_DC` env var or `--dc` flag on auth commands:

`com` (US, default) · `eu` · `in` · `com.au` · `jp` · `ca` · `sa` · `uk` · `com.cn`

## Supported products

| Product | Status | Commands |
|---------|--------|----------|
| **Books** | Supported | 515 |
| **Projects** | Supported | 232 |
| **Mail** | Supported | 159 |
| **Inventory** | Supported | 108 |
| **Invoice** | Supported | 103 |
| **Billing** | Supported | 88 |
| **Sheet** | Supported | 85 |
| **Expense** | Supported | 81 |
| **Bigin** | Supported | 34 |
| **CRM** | Supported | 29 |
| **Desk** | Supported | 28 |
| **Sign** | Supported | 24 |
| **Drive (WorkDrive)** | Supported | 24 |
| **Sprints** | Supported | 24 |
| **Recruit** | Supported | 23 |
| **Creator** | Supported | 15 |
| **People** | Supported | 18 |
| **Cliq** | Supported | 13 |
| **Writer** | Supported | 8 |
| **Auth** | Supported | 5 |
| Analytics | Planned | — |
| Campaigns | Planned | — |
| SalesIQ | Planned | — |
| Meeting | Planned | — |
| Bookings | Planned | — |
| Voice | Planned | — |
| Vault | Planned | — |
| Marketing Automation | Planned | — |
| PageSense | Planned | — |
| Assist | Planned | — |
| Learn | Planned | — |
| Showtime | Planned | — |
| Backstage | Planned | — |
| Calendar | Planned | — |
| Show | Planned | — |
| Forms | Planned | — |
| Social | Planned | — |
| Survey | Planned | — |
| Connect | Planned | — |
| Flow | Planned | — |
| BugTracker | Planned | — |
| Commerce | Planned | — |
| FSM | Planned | — |
| Directory | Planned | — |
| Shifts | Planned | — |
| Contracts | Planned | — |
| Practice | Planned | — |
| Checkout | Planned | — |
| Lens | Planned | — |
| ZeptoMail | Planned | — |
| Notebook | Planned | — |
| TeamInbox | Planned | — |
| Office Integrator | Planned | — |
| ToDo | Planned | — |
| PDF Editor | Planned | — |
| IoT | Planned | — |
| DataPrep | Planned | — |
| Apptics | Planned | — |
| Catalyst | Planned | — |
| Webinar | Planned | — |
| LandingPage | Planned | — |
| CommunitySpaces | Planned | — |
| Thrive | Planned | — |
| Sites | Planned | — |
| RouteIQ | Planned | — |
| Workerly | Planned | — |
| Solo | Planned | — |
| Procurement | Planned | — |

Want a product prioritized? [Open an issue](https://github.com/omin8tor/zoho-cli/issues).

Run `zoho --help-all` for the full command reference with every flag.

## Agent Skill

This repo ships as an [Agent Skill](https://agentskills.io/) so LLM agents can discover and use it. The skill definition is in [`SKILL.md`](./SKILL.md) with detailed references in [`references/`](./references/).

## Development

```bash
go build -o zoho ./cmd/zoho/
go test ./...
go vet ./...
```

Or with [mise](https://mise.jdx.dev/):

```bash
mise run build    # build binary
mise run test     # unit tests
mise run lint     # go vet
```

## License

GPL-3.0
