---
id: zc-dgak
status: closed
deps: [zc-06pa, zc-1g8v]
links: []
created: 2026-02-21T16:17:33Z
type: task
priority: 2
assignee: Jasmin Le Roux
---
# go: CRM notes (list, add, update, delete)

Port 4 CRM notes commands. list requires fields param (default: id,Note_Title,Note_Content,Created_Time,Modified_Time). add takes --content and --title, wraps as {data: [{Note_Content, Note_Title}]}. update/delete operate on note ID at /Notes/{id}. Typed envelope: {data: []Note}.

