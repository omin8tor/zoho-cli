---
id: zc-hnje
status: closed
deps: [zc-06pa, zc-x6xe, zc-1g8v]
links: []
created: 2026-02-21T16:17:32Z
type: task
priority: 2
assignee: Jasmin Le Roux
---
# go: CRM records (list, get, create, update, delete, search, upsert, bulk-delete)

Port 8 CRM record commands. list requires fields param (default: id,Created_Time,Modified_Time). get returns data[0]. create/update wrap in {data: [record]}. search supports --word/--email/--phone/--criteria. upsert supports --duplicate-check. bulk-delete takes comma-separated IDs. Typed envelope: {data: []map[string]any, info: PageInfo}.

