---
id: zc-c09i
status: closed
deps: [zc-aum3, zc-bhmw]
links: []
created: 2026-02-21T16:17:21Z
type: task
priority: 2
assignee: Jasmin Le Roux
---
# go: auth config resolution (env -> config -> error)

Port auth/config.py to Go. Resolve auth from: 1) env vars (ZOHO_CLIENT_ID, ZOHO_CLIENT_SECRET, ZOHO_REFRESH_TOKEN, ZOHO_DC), 2) config file (~/.config/zoho-cli/config.toml), 3) error if neither. AuthConfig struct with client_id, client_secret, refresh_token, dc fields.

