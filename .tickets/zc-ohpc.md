---
id: zc-ohpc
status: closed
deps: [zc-06pa, zc-x6xe, zc-1g8v]
links: []
created: 2026-02-21T16:17:34Z
type: task
priority: 2
assignee: Jasmin Le Roux
---
# go: CRM related, users, owner, coql, search-global

Port 5 CRM misc commands. related list requires fields param (default: id,Created_Time,Modified_Time). users list passes type=AllUsers. owner change POSTs {owner: {id}, notify: bool}. coql POSTs {select_query: string} - requires ZohoCRM.coql.READ scope. search-global uses searchword param (not word - v8 change). Typed envelopes per command.

