---
id: zc-42io
status: closed
deps: [zc-3xgp]
links: []
created: 2026-02-21T16:18:15Z
type: task
priority: 2
assignee: Jasmin Le Roux
---
# go: Drive folders, download, upload, upload-url

Port 6 WorkDrive commands. folders list at /teams/{team}/teamfolders. create at /files with JSON:API body, supports folder/zohowriter/zohosheet/zohoshow types. breadcrumb at /files/{id}/breadcrumbs. download at /api/v1/download/{id} (not download.zoho.com). upload multipart to /upload with parent_id. upload-url POSTs remote URL to /files/{folder}/remotefile.

