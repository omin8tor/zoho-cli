---
id: zc-06pa
status: closed
deps: [zc-kkbv]
links: []
created: 2026-02-21T16:17:23Z
type: task
priority: 2
assignee: Jasmin Le Roux
---
# go: HTTP client with auto-refresh on 401

Port http/client.py to Go using stdlib net/http. ZohoClient struct with Request() and RequestRaw() methods. Auto-refresh access token on 401 response. Support GET/POST/PUT/PATCH/DELETE. Handle JSON body, form data, file uploads, query params. Scope mismatch detection on 401.

