---
name: zoho-cli
description: Interact with Zoho REST APIs (20 products: CRM, Projects, WorkDrive, Books, Mail, Desk, Expense, Inventory, Invoice, Billing, Sheet, Sign, People, Sprints, Creator, Bigin, Recruit, Cliq, Writer) via CLI. 1600+ commands for records, tasks, files, invoices, tickets, email, and more. Use when the user needs to query, create, update, or manage data in any Zoho product, or when automating Zoho workflows.
license: GPL-3.0
compatibility: Requires the zoho-cli binary (Go). Needs network access to Zoho APIs. Needs ZOHO_CLIENT_ID, ZOHO_CLIENT_SECRET, and ZOHO_REFRESH_TOKEN env vars (or interactive auth).
metadata:
  author: omin8tor
  version: "1.0"
---

# zoho-cli

CLI for Zoho REST APIs. Single binary, JSON to stdout, 1600+ commands across 20 products.

## Install

```bash
go install github.com/omin8tor/zoho-cli/cmd/zoho@latest
```

Or download a prebuilt binary from GitHub releases.

## Authentication

Set these env vars:

```bash
export ZOHO_CLIENT_ID=1000.ABC123
export ZOHO_CLIENT_SECRET=xyz789
export ZOHO_REFRESH_TOKEN=1000.refresh_token_here
export ZOHO_DC=com
```

Tokens auto-refresh. The CLI caches access tokens at `~/.config/zoho-cli/cache/` to avoid hitting Zoho's 10-refreshes-per-10-minutes rate limit.

To get a refresh token, create a "Self Client" app at https://api-console.zoho.com/, generate a code with the needed scopes, then:

```bash
zoho auth self-client --code CODE --client-id ID --client-secret SECRET
```

Verify auth works:

```bash
zoho auth status
```

## How it works

Every command outputs JSON to stdout. Errors go to stderr. Exit codes: 0=success, 1=error, 2=auth error, 3=not found, 4=validation error.

The CLI is a thin wrapper — it passes through raw Zoho API responses without transformation. What Zoho returns is what you get.

All list commands support `--all` (fetch every page) and `--limit N` (fetch up to N records). Without these flags, you get a single page with the raw envelope.

Pipe into `jq` for filtering:

```bash
zoho crm records list Contacts --fields "Full_Name,Email" --all | jq '.[].Email'
```

## Quick reference by product

### CRM

```bash
zoho crm records list Contacts --fields "Full_Name,Email,Phone"
zoho crm records get Contacts RECORD_ID --fields "Full_Name,Email"
zoho crm records search Deals --criteria "(Stage:equals:Closed Won)" --fields "Deal_Name,Amount"
zoho crm records create Leads --json '{"Last_Name":"Smith","Company":"Acme"}'
zoho crm records update Leads RECORD_ID --json '{"Phone":"555-1234"}'
zoho crm records delete Leads RECORD_ID
zoho crm coql --query "SELECT Full_Name, Email FROM Contacts WHERE Email LIKE '%@acme.com' LIMIT 10"
zoho crm search-global "searchterm"
```

CRM v8 requires `--fields` on read endpoints. Without it, records come back empty.

COQL needs the `ZohoCRM.coql.READ` scope (separate from general CRM scopes).

### Projects

```bash
zoho projects core list --portal PORTAL_ID
zoho projects tasks my --portal PORTAL_ID
zoho projects tasks list --portal PORTAL_ID --project PROJECT_ID
zoho projects tasks create --portal PORTAL_ID --project PROJECT_ID --name "Task name"
zoho projects issues list --portal PORTAL_ID --project PROJECT_ID
```

Every Projects command needs `--portal`. You can set `ZOHO_PORTAL_ID` env var instead of passing it every time. The flag overrides the env var.

### WorkDrive

```bash
zoho drive teams me
zoho drive folders list --team TEAM_ID
zoho drive files list --folder FOLDER_ID
zoho drive files search --query "keyword" --team TEAM_ID
zoho drive download FILE_ID --output ./file.pdf
zoho drive upload ./file.xlsx --folder FOLDER_ID
zoho drive share add FILE_ID --email user@company.com --role editor
```

Navigate top-down: teams -> folders -> files. Set `ZOHO_TEAM_ID` env var to avoid passing `--team` every time.

### Books

```bash
zoho books invoices list --org ORG_ID
zoho books contacts list --org ORG_ID
zoho books invoices create --org ORG_ID --json '{"customer_id":"460000000026049","line_items":[{"item_id":"460000000026065","quantity":1}]}'
zoho books expenses list --org ORG_ID
zoho books items list --org ORG_ID
```

Every Books command needs `--org`. Set `ZOHO_BOOKS_ORG_ID` env var to skip the flag.

### Mail

```bash
zoho mail accounts list
zoho mail folders list --account ACCOUNT_ID
zoho mail messages list --account ACCOUNT_ID --folder FOLDER_ID
zoho mail messages get --account ACCOUNT_ID --folder FOLDER_ID --message MESSAGE_ID
```

Set `ZOHO_MAIL_ACCOUNT_ID` for the account default.

### Desk

```bash
zoho desk tickets list --org ORG_ID
zoho desk tickets get --org ORG_ID --id TICKET_ID
zoho desk contacts list --org ORG_ID
zoho desk search --org ORG_ID --module tickets --query "printer issue"
```

### Expense

```bash
zoho expense expenses list --org ORG_ID
zoho expense reports list --org ORG_ID
zoho expense categories list --org ORG_ID
```

### People

```bash
zoho people forms list
zoho people records list --form employee
zoho people attendance list --sdate 2025-01-01 --edate 2025-01-31
```

### Sprints

```bash
zoho sprints teams list
zoho sprints projects list --team TEAM_ID
zoho sprints items list --team TEAM_ID --project PROJECT_ID --sprint SPRINT_ID
```

### Sign

```bash
zoho sign requests list
zoho sign templates list
zoho sign requests get --id REQUEST_ID
```

### Writer

```bash
zoho writer details DOC_ID
zoho writer fields DOC_ID
zoho writer merge DOC_ID --json '{"name":"Alice"}' --format pdf --output ./out.pdf
zoho writer read DOC_ID
zoho writer download DOC_ID --format pdf --output ./doc.pdf
```

Writer has no "list documents" endpoint. Get doc IDs from WorkDrive.

### Cliq

```bash
zoho cliq channels list
zoho cliq channels message CHANNEL_NAME --text "message here"
zoho cliq buddies message user@company.com --text "hello"
zoho cliq messages list CHAT_ID
```

### Also supported

**Bigin** (34 commands), **Billing** (88), **Creator** (15), **Inventory** (108), **Invoice** (103), **Recruit** (23), **Sheet** (85). Run `zoho <product> --help` for details.

## Coverage

Supported: **Books** (515), **Projects** (232), **Mail** (159), **Inventory** (108), **Invoice** (103), **Billing** (88), **Sheet** (85), **Expense** (81), **Bigin** (34), **CRM** (29), **Desk** (28), **Sign** (24), **Drive** (24), **Sprints** (24), **Recruit** (23), **People** (18), **Creator** (15), **Cliq** (13), **Writer** (8), **Auth** (5).

Not yet supported: Analytics, Campaigns, SalesIQ, Meeting, Bookings, Voice, Vault, Marketing Automation, PageSense, Assist, Learn, Showtime, Backstage, and others. If the user asks about an unsupported product, tell them zoho-cli doesn't cover it yet and suggest they open an issue at https://github.com/omin8tor/zoho-cli/issues.

## Data centers

Set `ZOHO_DC` env var: `com` (US, default), `eu`, `in`, `com.au`, `jp`, `ca`, `sa`, `uk`, `com.cn`.

## Detailed references

- [Command reference](references/commands.md) — every command, flag, and usage pattern
- [API quirks](references/api-quirks.md) — Zoho-specific gotchas the CLI handles (or that you need to know about)
