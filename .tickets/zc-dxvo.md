---
id: zc-dxvo
status: closed
deps: [zc-06pa, zc-1g8v]
links: []
created: 2026-02-21T16:18:17Z
type: task
priority: 2
assignee: Jasmin Le Roux
---
# go: Drive share + teams

Port 9 WorkDrive commands. share permissions GET /files/{id}/permissions. add POST /permissions with role_id mapping (viewer=7, commenter=6, editor=5, organizer=4). revoke DELETE /permissions/{id}. link create POST /links (known Zoho 500 - handle gracefully). links list GET /files/{id}/links. unlink DELETE /links/{id}. teams me GET /users/me. teams list GET /users/me/teams. members GET /teams/{id}/members.

