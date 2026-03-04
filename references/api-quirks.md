# Zoho API quirks

Things the CLI handles internally, plus things you need to know.

## CRM

- **`--fields` is required on v8 read endpoints.** List, get, related, notes, and attachments all need it. Without `--fields`, the API returns records with no data. This isn't optional — it's how CRM v8 works.
- **`search-global` uses `searchword` as the API param**, not `word`. The CLI handles this.
- **Tags add/remove use JSON body**, not query params. The CLI handles this.
- **Pagination past 2000 records uses `page_token`**, not page numbers. The `--all` flag handles this automatically.
- **COQL needs its own scope**: `ZohoCRM.coql.READ`. If you didn't include it when generating your refresh token, COQL queries will 403.

## Projects

- **Every command needs a portal ID.** Pass `--portal` or set `ZOHO_PORTAL_ID` env var. Get your portal ID from `zoho projects core list`.
- **`has_next_page` can be string `"true"` or boolean `true`** depending on the endpoint. The CLI normalizes this.
- **Pagination is index-based** (range=0-99, range=100-199). The CLI handles this when you use `--all`.

## WorkDrive

- **Uses JSON:API content type** (`application/vnd.api+json`). The CLI handles headers.
- **Go's net/http doesn't send `Accept: */*` by default.** WorkDrive returns 415 without it. The CLI adds this header to every request.
- **Copy has reversed semantics**: You POST to the destination folder with the source file ID in the body. The `--to` flag on `files copy` handles this correctly.
- **File status codes**: 1=active, 51=trash, 61=permanently deleted.
- **Download endpoint**: Uses `workdrive.zoho.com/api/v1/download/{id}`, not `download.zoho.com`.
- **`trash-list`**: The correct endpoint is `teamfolders/{id}/trashedfiles`, NOT `files/{id}/files?filter[status]=51`. The `filter[status]` param is not in the API spec.
- **`teams members`**: The correct endpoint is `teams/{id}/users`, NOT `teams/{id}/members`. The `/members` path returns "URL Rule is not configured".
- **Permanent delete (status=61)** returns 401 Unauthorized on some accounts. Use trash (status=51) instead.
- **`files get` on trashed copies** returns 401 Unauthorized. Trashed copies cannot be read back via `files get`.
- **Search indexing delay**: Newly uploaded files may take 30-60 seconds to appear in search results.
- **Folder names get timestamp suffix**: Created folders have ` DD-MM-YYYY HH:MM:SS:mmm` appended to the name.
- **`share link` (create)** returns 500 on some accounts — likely a plan/feature limitation, not a code bug.
- **`upload-url` (remotefile endpoint)** does not exist — removed from CLI. The `/files/{id}/remotefile` path was never a real API endpoint (returns F6016).

## Writer

- **No "list documents" endpoint.** You can't list Writer docs through the Writer API. Get document IDs from WorkDrive instead.
- **Empty documents can't be exported.** Writer returns error R3002. The CLI surfaces this clearly.

## Token management

- **Zoho rate-limits to 10 access token refreshes per refresh_token per 10 minutes.** The CLI caches tokens at `~/.config/zoho-cli/cache/{hash}.json` to avoid this. If you're running multiple parallel instances, they share the cache.
- **Token auto-refresh**: On 401, the CLI refreshes the access token once and retries. If that fails, exit code 2.

## Data centers

9 DCs, each with its own API base URLs. Set via `ZOHO_DC` env var or `--dc` on auth commands.

| Code | Region |
|------|--------|
| com | US (default) |
| eu | Europe |
| in | India |
| com.au | Australia |
| jp | Japan |
| ca | Canada |
| sa | Saudi Arabia |
| uk | United Kingdom |
| com.cn | China |
