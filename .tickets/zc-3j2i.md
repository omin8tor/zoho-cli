---
id: zc-3j2i
status: closed
deps: [zc-sgty]
links: []
created: 2026-02-21T16:18:11Z
type: task
priority: 2
assignee: Jasmin Le Roux
---
# go: Projects issues (list, get, create, update, delete, defaults) + issue-comments

Port 8 issue commands. CRUD at /projects/{pid}/issues/. defaults at /issues/defaultfields. issue-comments list/add at /issues/{iid}/comments. Issue field is 'name' not 'title'. Typed envelope with issues key.

