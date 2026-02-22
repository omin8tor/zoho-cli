---
id: zc-10zs
status: closed
deps: [zc-wj3j]
links: []
created: 2026-02-21T16:18:34Z
type: task
priority: 2
assignee: Jasmin Le Roux
---
# go: GitHub Actions CI (lint, test, build)

Add .github/workflows/ci.yml with: Go setup, golangci-lint, go test ./..., go build. Run on push to main and PRs. Matrix: Go 1.24, linux/macos. Cache Go modules. Integration tests as separate workflow with secrets.

