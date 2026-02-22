---
id: zc-w5p4
status: closed
deps: [zc-sgty]
links: []
created: 2026-02-21T16:18:09Z
type: task
priority: 2
assignee: Jasmin Le Roux
---
# go: Projects tasks (list, my, get, create, update, delete, subtasks, add-subtask)

Port 8 task commands. list/get/create/update/delete at /projects/{pid}/tasks/. my lists tasks across all projects at /portal/{portal}/tasks. subtasks lists at /tasks/{tid}/subtasks. add-subtask POSTs to /tasks/{tid}/subtasks. All support --status and --priority filters. Typed envelope with tasks key.

