---
id: zc-shph
status: closed
deps: [zc-06pa, zc-1g8v]
links: []
created: 2026-02-21T16:18:20Z
type: task
priority: 2
assignee: Jasmin Le Roux
---
# go: Cliq channels (list, get, create, message, members)

Port 5 Cliq channel commands. list GET /api/v2/channels. get uses /channelsbyname/{name}. create POST /channels with name+description. message POST /channelsbyname/{name}/message with text and optional bot. members GET /channelsbyname/{name}/members.

