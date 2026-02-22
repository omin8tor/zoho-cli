---
id: zc-x8rt
status: closed
deps: [zc-hnje, zc-w5p4, zc-3xgp, zc-0apo, zc-shph]
links: []
created: 2026-02-21T16:18:30Z
type: task
priority: 2
assignee: Jasmin Le Roux
---
# go: integration tests against real Zoho API

Port Python integration test suite to Go. Test all 96 commands against real Zoho API using .env.test credentials. CRUD lifecycle tests for CRM records, Projects tasks/issues, Drive files/folders, Writer documents. Use testing subtests. Mark as integration build tag. Known Zoho-side failures: drive share link (500), projects timelogs add (500), coql (scope mismatch).

