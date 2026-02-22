---
id: zc-aum3
status: closed
deps: [zc-0nge]
links: []
created: 2026-02-21T16:17:18Z
type: task
priority: 2
assignee: Jasmin Le Roux
---
# go: error types and exit codes

Port errors.py to Go. Define error types: ZohoCliError, AuthError, NotFoundError, ValidationError, ZohoAPIError. Exit codes: 0=success, 1=general, 2=auth, 3=not found, 4=validation. Stderr helper function.

