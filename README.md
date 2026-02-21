# zoho-cli

CLI for Zoho's REST APIs. Covers CRM, Projects, WorkDrive, Writer, and Cliq.

There is no other tool that does this. Zoho's own CLIs (ZET, Catalyst) are for building extensions, not for talking to the API. Their Python SDK covers CRM only and requires MySQL for token storage. We checked.

## Install

**From source** (requires Go 1.22+):

```bash
go install github.com/omin8tor/zoho-cli/cmd/zoho@latest
```

**From release** (no Go required):

Download the binary for your platform from [releases](https://github.com/omin8tor/zoho-cli/releases), unpack, put it on your PATH.

## Auth

You need a Zoho API Console app. Go to [api-console.zoho.com](https://api-console.zoho.com/), create a "Self Client" app, and get a client ID and secret.

**Device flow** (interactive, for humans):

```bash
zoho auth login --client-id YOUR_ID --client-secret YOUR_SECRET
```

This opens a browser, you approve, done. Tokens are saved to `~/.config/zoho-cli/config.json`.

**Self-client code** (non-interactive, for CI/agents):

```bash
zoho auth self-client --code GENERATED_CODE --client-id YOUR_ID --client-secret YOUR_SECRET
```

**Environment variables** (simplest for automation):

```bash
export ZOHO_CLIENT_ID=...
export ZOHO_CLIENT_SECRET=...
export ZOHO_REFRESH_TOKEN=...
export ZOHO_DC=com  # or eu, in, com.au, jp, ca, sa, uk
```

The CLI checks env vars first, then config file.

## Usage

```bash
# CRM
zoho crm records list Contacts --fields "Full_Name,Email"
zoho crm records get Contacts 12345 --fields "Full_Name,Email,Phone"
zoho crm records search Deals --criteria "(Stage:equals:Closed Won)"
zoho crm records create Leads --json '{"Last_Name":"Smith","Company":"Acme"}'

# Projects
zoho projects core list --portal 12345
zoho projects tasks list --portal 12345 --project 67890
zoho projects tasks my --portal 12345

# WorkDrive
zoho drive teams me
zoho drive folders list --team TEAM_ID
zoho drive files list --folder FOLDER_ID
zoho drive download FILE_ID --output ./local-copy.pdf

# Writer
zoho writer details DOC_ID
zoho writer merge DOC_ID --json '{"name":"Alice"}' --format pdf --output ./merged.pdf

# Cliq
zoho cliq channels list
zoho cliq buddies message someone@company.com --text "hey"
```

Output is JSON to stdout. Errors go to stderr. Pipe into `jq` for filtering:

```bash
zoho crm records list Contacts --fields "Full_Name,Email" | jq '.[].Email'
zoho projects tasks my --portal 12345 | jq '[.[] | select(.status.name == "Open")]'
```

Run `zoho --help-all` to see every command and its flags.

## What's here

118 commands across 6 groups:

| Group | Commands | Covers |
|-------|----------|--------|
| auth | 5 | login, self-client, status, refresh, logout |
| crm | 29 | records CRUD, search, notes, related, tags, attachments, COQL, users |
| projects | 39 | projects, tasks, issues, comments, tasklists, timelogs, milestones, dependencies |
| drive | 26 | files, folders, download, upload, sharing, teams |
| writer | 7 | create, details, merge, read, download |
| cliq | 12 | channels, chats, buddies, messages, users |

Single binary, no runtime dependencies. Builds for Linux, macOS, and Windows on amd64 and arm64.

## Data centers

Zoho runs in 9 data centers. Set via `--dc` flag on auth commands or `ZOHO_DC` env var:

`com` (US), `eu`, `in`, `com.au`, `jp`, `ca`, `sa`, `uk`, `com.cn`

Default is `com`.

## Development

```bash
go build -o zoho ./cmd/zoho/    # build
go test ./...                    # unit tests
go vet ./...                     # lint
```

Or with [mise](https://mise.jdx.dev/):

```bash
mise run build
mise run test
mise run lint
```

## License

GPL-3.0
