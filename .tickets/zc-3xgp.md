---
id: zc-3xgp
status: closed
deps: [zc-06pa, zc-x6xe, zc-1g8v]
links: []
created: 2026-02-21T16:18:14Z
type: task
priority: 2
assignee: Jasmin Le Roux
---
# go: Drive files (list, get, search, rename, copy, move, trash, delete, restore, trash-list, versions)

Port 11 WorkDrive file commands. All use JSON:API content-type (application/vnd.api+json). list paginates via paginate_workdrive. rename/move/trash/delete/restore use PATCH with status codes (51=trash, 61=delete, 1=restore). copy uses reversed semantics: POST to /files/{destination}/copy with source in body. search at /teams/{team}/records. versions at /files/{id}/versions.

