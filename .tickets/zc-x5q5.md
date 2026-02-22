---
id: zc-x5q5
status: closed
deps: [zc-shph]
links: []
created: 2026-02-21T16:18:21Z
type: task
priority: 2
assignee: Jasmin Le Roux
---
# go: Cliq chats, buddies, messages, users

Port 6 Cliq commands. chats message POST /chats/{id}/message. buddies message POST /buddies/{email}/message. messages list GET /chats/{id}/messages with limit param. messages edit PUT and delete DELETE at /chats/{id}/messages/{mid}. users list GET /api/v2/users. users get GET /api/v2/users/{id}.

