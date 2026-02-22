---
id: zc-wj3j
status: closed
deps: [zc-x6xe, zc-1g8v]
links: []
created: 2026-02-21T16:18:28Z
type: task
priority: 2
assignee: Jasmin Le Roux
---
# go: unit tests for core infrastructure

Write Go unit tests for: DC maps (all 9 DCs resolve correctly), auth config resolution (env vars, config file, missing), token caching (file read/write, expiry), pagination helpers (CRM page_token, Projects has_next_page string/bool, WorkDrive offset), output formatting (JSON to stdout). Aim for >80% coverage on internal/ packages.

