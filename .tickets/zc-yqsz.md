---
id: zc-yqsz
status: closed
deps: [zc-06pa, zc-c09i]
links: []
created: 2026-02-21T16:17:28Z
type: task
priority: 2
assignee: Jasmin Le Roux
---
# go: auth login device flow

Port auth/device_flow.py to Go. Implement 'zoho auth login' using Zoho's device flow OAuth: POST to /oauth/v3/device/code, poll /oauth/v3/token until user authorizes. Display user code and verification URL. Save refresh token to config.

