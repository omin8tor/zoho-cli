---
id: zc-kkbv
status: closed
deps: [zc-c09i]
links: []
created: 2026-02-21T16:17:22Z
type: task
priority: 2
assignee: Jasmin Le Roux
---
# go: token refresh with file-based caching

Port auth/token.py to Go. Refresh access token via Zoho OAuth endpoint. Cache tokens at ~/.config/zoho-cli/cache/{hash}.json. Rate-limit awareness: Zoho limits 10 refreshes per refresh_token per 10 minutes. Return cached token if still valid.

