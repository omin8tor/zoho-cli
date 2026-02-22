---
id: zc-sgty
status: closed
deps: [zc-06pa, zc-x6xe, zc-1g8v]
links: []
created: 2026-02-21T16:18:08Z
type: task
priority: 2
assignee: Jasmin Le Roux
---
# go: Projects core (list, get, search, create, update)

Port 5 Projects commands. All use portal ID. list auto-paginates via paginate_projects. get returns single project. search uses search_term+module=all+status=all params. create/update take --name and --json for extra fields. URL pattern: /api/v3/portal/{portal}/projects/.

