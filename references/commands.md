# Command reference

118 commands. Run `zoho --help-all` for the native output with all flags.

## Contents

- [auth](#auth) (5 commands)
- [crm](#crm) (29 commands)
- [projects](#projects) (39 commands)
- [drive](#drive) (26 commands)
- [writer](#writer) (7 commands)
- [cliq](#cliq) (12 commands)

---

## auth

### zoho auth login

Authenticate via device flow (interactive, opens browser).

```
zoho auth login --client-id ID --client-secret SECRET [--dc com] [--scopes "scope1,scope2"]
```

### zoho auth self-client

Exchange a self-client code from Zoho API Console.

```
zoho auth self-client --code CODE --client-id ID --client-secret SECRET [--dc com] [--server URL]
```

### zoho auth status

Show current auth status (token validity, data center, scopes).

```
zoho auth status
```

### zoho auth refresh

Force-refresh the access token.

```
zoho auth refresh
```

### zoho auth logout

Clear stored tokens.

```
zoho auth logout
```

---

## crm

### zoho crm modules list

List CRM modules.

```
zoho crm modules list [--include-hidden]
```

### zoho crm modules fields

List fields for a module.

```
zoho crm modules fields <module>
```

### zoho crm modules related-lists

List related lists for a module.

```
zoho crm modules related-lists <module>
```

### zoho crm modules layouts

List layouts for a module.

```
zoho crm modules layouts <module>
```

### zoho crm modules custom-views

List custom views for a module.

```
zoho crm modules custom-views <module>
```

### zoho crm records list

List records. Requires `--fields`.

```
zoho crm records list <module> --fields "Field1,Field2" [--sort-by Field] [--sort-order asc|desc] [--page N] [--per-page N] [--all]
```

### zoho crm records get

Get one record by ID. Requires `--fields`.

```
zoho crm records get <module> <record-id> --fields "Field1,Field2"
```

### zoho crm records create

Create a record.

```
zoho crm records create <module> --json '{"Field":"Value"}' [--trigger approval,workflow,blueprint]
```

### zoho crm records update

Update a record.

```
zoho crm records update <module> <record-id> --json '{"Field":"NewValue"}' [--trigger ...]
```

### zoho crm records delete

Delete a record.

```
zoho crm records delete <module> <record-id>
```

### zoho crm records search

Search records by keyword, email, phone, or criteria.

```
zoho crm records search <module> [--word text] [--email addr] [--phone num] [--criteria "(Field:op:Value)"] [--fields "..."] [--page N] [--per-page N]
```

### zoho crm records upsert

Insert or update based on duplicate check fields.

```
zoho crm records upsert <module> --json '{"..."}' --duplicate-check "Email" [--trigger ...]
```

### zoho crm records bulk-delete

Delete multiple records.

```
zoho crm records bulk-delete <module> "id1,id2,id3"
```

### zoho crm notes list

List notes on a record. Requires `--fields`.

```
zoho crm notes list <module> <record-id> [--fields "..."] [--page N] [--per-page N]
```

### zoho crm notes add

Add a note.

```
zoho crm notes add <module> <record-id> --content "text" [--title "title"]
```

### zoho crm notes update

```
zoho crm notes update <note-id> [--title "..."] [--content "..."]
```

### zoho crm notes delete

```
zoho crm notes delete <note-id>
```

### zoho crm related list

List related records.

```
zoho crm related list <module> <record-id> <related-list> [--fields "..."] [--page N] [--per-page N]
```

### zoho crm users list

```
zoho crm users list [--page N] [--per-page N]
```

### zoho crm owner change

Change record owner.

```
zoho crm owner change <module> <record-id> --owner USER_ID [--notify]
```

### zoho crm coql

Run a COQL query. Needs `ZohoCRM.coql.READ` scope.

```
zoho crm coql --query "SELECT Field FROM Module WHERE condition LIMIT N"
```

### zoho crm search-global

Search across all CRM modules.

```
zoho crm search-global <word> [--page N] [--per-page N]
```

### zoho crm attachments list

```
zoho crm attachments list <module> <record-id> [--fields "..."] [--page N] [--per-page N]
```

### zoho crm attachments upload

```
zoho crm attachments upload <module> <record-id> <file-path>
```

### zoho crm attachments download

```
zoho crm attachments download <module> <record-id> <attachment-id> [--output path]
```

### zoho crm attachments delete

```
zoho crm attachments delete <module> <record-id> <attachment-id>
```

### zoho crm convert

Convert a lead to contact/account/deal.

```
zoho crm convert <record-id> [--json '{"overwrite":true,"notify_lead_owner":true}']
```

### zoho crm tags add

```
zoho crm tags add <module> --ids "id1,id2" --tags "tag1,tag2"
```

### zoho crm tags remove

```
zoho crm tags remove <module> --ids "id1,id2" --tags "tag1,tag2"
```

---

## projects

All Projects commands need a portal ID. Pass `--portal PORTAL_ID` or set `ZOHO_PORTAL_ID` env var. The flag overrides the env var.

### zoho projects core list

```
zoho projects core list --portal ID
```

### zoho projects core get

```
zoho projects core get <project-id> --portal ID
```

### zoho projects core search

```
zoho projects core search --portal ID --query "text"
```

### zoho projects core create

```
zoho projects core create --portal ID --name "Name" [--json '{"..."}']
```

### zoho projects core update

```
zoho projects core update <project-id> --portal ID --json '{"..."}' 
```

### zoho projects tasks list

```
zoho projects tasks list --portal ID --project PID [--status open|closed|"in progress"] [--priority none|low|medium|high]
```

### zoho projects tasks my

List your tasks across all projects.

```
zoho projects tasks my --portal ID [--status ...] [--priority ...]
```

### zoho projects tasks get

```
zoho projects tasks get <task-id> --portal ID --project PID
```

### zoho projects tasks create

```
zoho projects tasks create --portal ID --project PID --name "Task" [--json '{"..."}']
```

### zoho projects tasks update

```
zoho projects tasks update <task-id> --portal ID --project PID --json '{"..."}' 
```

### zoho projects tasks delete

```
zoho projects tasks delete <task-id> --portal ID --project PID
```

### zoho projects tasks subtasks

```
zoho projects tasks subtasks <task-id> --portal ID --project PID
```

### zoho projects tasks add-subtask

```
zoho projects tasks add-subtask --portal ID --project PID --parent TASK_ID --name "Subtask" [--json '{"..."}']
```

### zoho projects issues list

```
zoho projects issues list --portal ID --project PID
```

### zoho projects issues get

```
zoho projects issues get <issue-id> --portal ID --project PID
```

### zoho projects issues create

```
zoho projects issues create --portal ID --project PID --name "Issue" [--json '{"..."}']
```

### zoho projects issues update

```
zoho projects issues update <issue-id> --portal ID --project PID --json '{"..."}' 
```

### zoho projects issues delete

```
zoho projects issues delete <issue-id> --portal ID --project PID
```

### zoho projects issues defaults

Get default statuses, severities, priorities for issues.

```
zoho projects issues defaults --portal ID --project PID
```

### zoho projects issue-comments list

```
zoho projects issue-comments list --portal ID --project PID --issue ISSUE_ID
```

### zoho projects issue-comments add

```
zoho projects issue-comments add --portal ID --project PID --issue ISSUE_ID --comment "text"
```

### zoho projects comments list

Task comments.

```
zoho projects comments list --portal ID --project PID --task TASK_ID
```

### zoho projects comments add

```
zoho projects comments add --portal ID --project PID --task TASK_ID --comment "text"
```

### zoho projects comments update

```
zoho projects comments update <comment-id> --portal ID --project PID --task TASK_ID --comment "text"
```

### zoho projects comments delete

```
zoho projects comments delete <comment-id> --portal ID --project PID --task TASK_ID
```

### zoho projects tasklists list

```
zoho projects tasklists list --portal ID --project PID
```

### zoho projects tasklists create

```
zoho projects tasklists create --portal ID --project PID --name "Name" [--json '{"..."}']
```

### zoho projects tasklists update

```
zoho projects tasklists update <tasklist-id> --portal ID --project PID --json '{"..."}' 
```

### zoho projects tasklists delete

```
zoho projects tasklists delete <tasklist-id> --portal ID --project PID
```

### zoho projects timelogs list

```
zoho projects timelogs list --portal ID --project PID [--module task|issue|general]
```

### zoho projects timelogs add

```
zoho projects timelogs add --portal ID --project PID --date 2025-01-15 --hours 2 [--task TASK_ID] [--bill-status "Billable"] [--notes "text"]
```

### zoho projects users list

```
zoho projects users list --portal ID --project PID
```

### zoho projects milestones list

```
zoho projects milestones list --portal ID --project PID
```

### zoho projects milestones get

```
zoho projects milestones get <milestone-id> --portal ID --project PID
```

### zoho projects milestones create

```
zoho projects milestones create --portal ID --project PID --name "Name" --start 2025-01-01 --end 2025-03-31 [--json '{"..."}']
```

### zoho projects milestones update

```
zoho projects milestones update <milestone-id> --portal ID --project PID --json '{"..."}' 
```

### zoho projects milestones delete

```
zoho projects milestones delete <milestone-id> --portal ID --project PID
```

### zoho projects dependencies add

```
zoho projects dependencies add <task-id> --portal ID --project PID --depends-on OTHER_TASK_ID [--type FS|SS|FF|SF]
```

### zoho projects dependencies remove

```
zoho projects dependencies remove <task-id> <dependency-id> --portal ID --project PID
```

---

## drive

### zoho drive teams me

Get current user info.

```
zoho drive teams me
```

### zoho drive teams list

```
zoho drive teams list
```

### zoho drive teams members

```
zoho drive teams members <team-id>
```

### zoho drive folders list

List top-level team folders.

```
zoho drive folders list --team TEAM_ID    # or set ZOHO_TEAM_ID env var
```

### zoho drive folders create

```
zoho drive folders create --name "Name" --parent FOLDER_ID [--type folder|zohowriter|zohosheet|zohoshow]
```

### zoho drive folders breadcrumb

Show the full path to a folder.

```
zoho drive folders breadcrumb <folder-id>
```

### zoho drive files list

List contents of a folder.

```
zoho drive files list --folder FOLDER_ID [--type file|folder|image]
```

### zoho drive files get

```
zoho drive files get <file-id>
```

### zoho drive files search

```
zoho drive files search --query "keyword" --team TEAM_ID [--mode all|name|content] [--type ...]    # or set ZOHO_TEAM_ID env var
```

### zoho drive files rename

```
zoho drive files rename <file-id> --name "New Name"
```

### zoho drive files copy

```
zoho drive files copy <file-id> --to DESTINATION_FOLDER_ID
```

### zoho drive files move

```
zoho drive files move <file-id> --to DESTINATION_FOLDER_ID
```

### zoho drive files trash

```
zoho drive files trash <file-id>
```

### zoho drive files delete

Permanently delete.

```
zoho drive files delete <file-id>
```

### zoho drive files restore

Restore from trash.

```
zoho drive files restore <file-id>
```

### zoho drive files trash-list

List trashed files in a team folder.

```
zoho drive files trash-list --team-folder TEAMFOLDER_ID
```

### zoho drive files versions

```
zoho drive files versions <file-id>
```

### zoho drive download

```
zoho drive download <file-id> [--output path] [--format native|txt|html|pdf|docx]
```

### zoho drive upload

```
zoho drive upload <file-path> --folder FOLDER_ID [--override]
```

### zoho drive share permissions

```
zoho drive share permissions <file-id>
```

### zoho drive share add

```
zoho drive share add <file-id> --email user@co.com [--role viewer|commenter|editor|organizer]
```

### zoho drive share revoke

```
zoho drive share revoke <permission-id>
```

### zoho drive share links

List share links for a file.

```
zoho drive share links <file-id>
```

### zoho drive share link

Create or get a share link.

```
zoho drive share link <file-id> [--role viewer|commenter|editor] [--allow-download] [--name "Link Name"] [--expiration 2025-12-31] [--password secret]
```

### zoho drive share unlink

```
zoho drive share unlink <link-id>
```

---

## writer

### zoho writer create

```
zoho writer create --name "Doc Name" [--folder FOLDER_ID] [--type writer|sheet|show]
```

### zoho writer details

```
zoho writer details <doc-id>
```

### zoho writer fields

List merge fields in a template.

```
zoho writer fields <doc-id>
```

### zoho writer merge

Merge data into a template and export.

```
zoho writer merge <doc-id> --json '{"field":"value"}' [--format pdf|docx|inline] [--output path]
```

### zoho writer delete

```
zoho writer delete <doc-id>
```

### zoho writer read

Read document content as text.

```
zoho writer read <doc-id> [--format txt|html]
```

### zoho writer download

```
zoho writer download <doc-id> [--format txt|html|pdf|docx|odt|rtf|epub] [--output path]
```

---

## cliq

### zoho cliq channels list

```
zoho cliq channels list
```

### zoho cliq channels get

```
zoho cliq channels get <channel-name>
```

### zoho cliq channels create

```
zoho cliq channels create --name "channel-name" [--description "..."]
```

### zoho cliq channels message

```
zoho cliq channels message <channel-name> --text "message" [--bot BOT_NAME]
```

### zoho cliq channels members

```
zoho cliq channels members <channel-name>
```

### zoho cliq chats message

Send a message to a chat by ID.

```
zoho cliq chats message <chat-id> --text "message"
```

### zoho cliq buddies message

Send a DM by email.

```
zoho cliq buddies message <email> --text "message"
```

### zoho cliq messages list

```
zoho cliq messages list <chat-id> [--limit N]
```

### zoho cliq messages edit

```
zoho cliq messages edit <chat-id> <message-id> --text "new text"
```

### zoho cliq messages delete

```
zoho cliq messages delete <chat-id> <message-id>
```

### zoho cliq users list

```
zoho cliq users list
```

### zoho cliq users get

```
zoho cliq users get <user-id>
```
