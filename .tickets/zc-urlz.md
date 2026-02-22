---
id: zc-urlz
status: closed
deps: [zc-sgty]
links: []
created: 2026-02-21T16:18:12Z
type: task
priority: 2
assignee: Jasmin Le Roux
---
# go: Projects task-comments, tasklists, timelogs, users, milestones, dependencies

Port 16 misc Projects commands. Task comments CRUD at /tasks/{tid}/comments (field is 'comment' not 'content'). Tasklists CRUD at /tasklists/. Timelogs list uses module JSON param, add POSTs to /log (known Zoho 500 - handle gracefully). Users at /users (name in full_name/display_name). Milestones CRUD with --start/--end dates. Dependencies add/remove at /tasks/{tid}/dependencies.

