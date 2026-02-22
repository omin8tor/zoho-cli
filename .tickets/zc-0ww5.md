---
id: zc-0ww5
status: closed
deps: [zc-06pa, zc-1g8v]
links: []
created: 2026-02-21T16:17:35Z
type: task
priority: 2
assignee: Jasmin Le Roux
---
# go: CRM attachments (list, upload, download, delete)

Port 4 CRM attachment commands. list requires fields param (default: id,File_Name,Size,Created_Time). upload sends multipart file to /{module}/{id}/Attachments. download uses RequestRaw, writes to --output or stdout. delete DELETEs /{module}/{id}/Attachments/{att_id}.

