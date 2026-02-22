---
id: zc-x6xe
status: closed
deps: [zc-06pa]
links: []
created: 2026-02-21T16:17:24Z
type: task
priority: 2
assignee: Jasmin Le Roux
---
# go: pagination helpers (CRM, Projects, WorkDrive)

Port pagination.py to Go. Three styles: 1) CRM - page_token based pagination for >2000 records, 2) Projects - page/per_page with has_next_page (string or bool), 3) WorkDrive - offset/limit with meta.has_next. Each returns []json.RawMessage.

