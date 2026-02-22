---
id: zc-0apo
status: closed
deps: [zc-06pa, zc-1g8v]
links: []
created: 2026-02-21T16:18:19Z
type: task
priority: 2
assignee: Jasmin Le Roux
---
# go: Writer commands (create, details, fields, merge, delete, read, download)

Port 7 Writer commands. create uses WorkDrive /files endpoint with service_type=zw. details/fields/delete at /writer/api/v1/documents/{id}. merge POSTs merge_data to /documents/{id}/merge, supports pdf/docx/inline output. read/download use /writer/api/v1/download/{id} with format param. Handle R3002 error (empty document cannot be exported).

