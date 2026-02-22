---
id: zc-ncho
status: closed
deps: [zc-1g8v]
links: []
created: 2026-02-21T16:18:31Z
type: task
priority: 2
assignee: Jasmin Le Roux
---
# go: --help-all with jq-friendly output schemas

Implement --help-all flag that recursively prints all commands with their arguments, flags, and output schema. Output schema should show JSON structure with types, enough for users to write jq expressions. Example: 'zoho crm records list -> {"data": [{"id": "string", ...}], "info": {"per_page": int, "more_records": bool, "page_token": "string"}}'. Register schema metadata on each command.

