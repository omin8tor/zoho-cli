---
id: zc-1g8v
status: closed
deps: [zc-0nge]
links: []
created: 2026-02-21T16:17:26Z
type: task
priority: 2
assignee: Jasmin Le Roux
---
# go: JSON output and jq-friendly schema display

Port output/ to Go. JSON output to stdout (pretty-printed). Errors to stderr. Implement --help-all flag showing recursive help for all commands with jq-friendly output schema per command (e.g. 'Output: {"data": [{"id": "string", ...}], "info": {"more_records": bool}}').

