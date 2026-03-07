# Panelist Context Package

## Project Overview

zoho-cli is a Go CLI tool wrapping Zoho REST APIs. Each product module lives in `internal/<module>/<module>.go`, exports `Commands() *cli.Command`, and integration tests go in `internal/<module>_integration_test.go`.

## Tech Stack
- Go 1.25+, urfave/cli v3, stdlib net/http
- No comments in code
- JSON output to stdout, errors to stderr
- Exit codes: 0=success, 1=general error, 2=auth error, 3=not found, 4=validation error

## Module Pattern

Every module follows this exact pattern (from expense.go):

```go
package mymodule

import (
    "context"
    "encoding/json"
    "fmt"
    "os"

    "github.com/omin8tor/zoho-cli/internal"
    "github.com/omin8tor/zoho-cli/internal/auth"
    zohttp "github.com/omin8tor/zoho-cli/internal/http"
    "github.com/omin8tor/zoho-cli/internal/output"
    "github.com/urfave/cli/v3"
)

func getClient() (*zohttp.Client, error) {
    config, err := auth.ResolveAuth()
    if err != nil {
        return nil, err
    }
    return zohttp.NewClient(config)
}

// For products requiring org ID:
func resolveOrgID(cmd *cli.Command) (string, error) {
    org := cmd.String("org")
    if org == "" {
        org = os.Getenv("ZOHO_MYMODULE_ORG_ID")
    }
    if org == "" {
        return "", internal.NewValidationError("--org flag or ZOHO_MYMODULE_ORG_ID env var required")
    }
    return org, nil
}

func Commands() *cli.Command {
    return &cli.Command{
        Name:  "mymodule",
        Usage: "Zoho MyModule operations",
        Commands: []*cli.Command{
            resourceOneCmd(),
            resourceTwoCmd(),
        },
    }
}

func resourceOneCmd() *cli.Command {
    return &cli.Command{
        Name:  "resource-one",
        Usage: "Resource One operations",
        Commands: []*cli.Command{
            {
                Name:  "list",
                Usage: "List resource ones",
                Action: func(_ context.Context, cmd *cli.Command) error {
                    c, err := getClient()
                    if err != nil {
                        return err
                    }
                    raw, err := c.Request("GET", c.MyModuleBase+"/resource-ones", nil)
                    if err != nil {
                        return err
                    }
                    return output.JSONRaw(raw)
                },
            },
            {
                Name:      "get",
                Usage:     "Get a resource one",
                ArgsUsage: "<id>",
                Action: func(_ context.Context, cmd *cli.Command) error {
                    if cmd.Args().Len() < 1 {
                        return internal.NewValidationError("resource one ID required")
                    }
                    c, err := getClient()
                    if err != nil {
                        return err
                    }
                    raw, err := c.Request("GET", c.MyModuleBase+"/resource-ones/"+cmd.Args().First(), nil)
                    if err != nil {
                        return err
                    }
                    return output.JSONRaw(raw)
                },
            },
            {
                Name:  "create",
                Usage: "Create a resource one",
                Flags: []cli.Flag{
                    &cli.StringFlag{Name: "json", Required: true, Usage: "JSON body"},
                },
                Action: func(_ context.Context, cmd *cli.Command) error {
                    c, err := getClient()
                    if err != nil {
                        return err
                    }
                    var body any
                    if err := json.Unmarshal([]byte(cmd.String("json")), &body); err != nil {
                        return internal.NewValidationError(fmt.Sprintf("invalid JSON: %v", err))
                    }
                    raw, err := c.Request("POST", c.MyModuleBase+"/resource-ones", &zohttp.RequestOpts{JSON: body})
                    if err != nil {
                        return err
                    }
                    return output.JSONRaw(raw)
                },
            },
        },
    }
}
```

## Client Base Fields

The HTTP client (`internal/http/client.go`) has these base URL fields for new modules:

| Field | Base URL | Products Using It |
|-------|----------|-------------------|
| `AnalyticsBase` | `{APIURL}/analytics/v2` | Analytics |
| `AssistBase` | `{APIURL}/assist/v1` | Assist |
| `BackstageBase` | `{APIURL}/backstage/v1` | Backstage |
| `BiginBase` | `{APIURL}/bigin/v2` | Bigin |
| `BillingBase` | `{APIURL}/billing/v1` | Billing/Subscriptions |
| `BookingsBase` | `{APIURL}/bookings/v1` | Bookings |
| `CampaignsBase` | `{APIURL}/campaigns/v1` | Campaigns |
| `CreatorBase` | `{APIURL}/creator/v2.1` | Creator |
| `InventoryBase` | `{APIURL}/inventory/v1` | Inventory |
| `InvoiceBase` | `{APIURL}/invoice/v3` | Invoice |
| `LearnBase` | `{APIURL}/learn/v1` | Learn |
| `MarketingAutoBase` | `{APIURL}/marketingautomation/v1` | Marketing Automation |
| `MeetingBase` | `{APIURL}/meeting/v1` | Meeting |
| `PageSenseBase` | `{APIURL}/pagesense/v1` | PageSense |
| `PeopleBase` | `{PeopleURL}/people/api` | People |
| `RecruitBase` | `{APIURL}/recruit/v2` | Recruit |
| `SalesIQBase` | `{APIURL}/salesiq/v2` | SalesIQ |
| `ShowtimeBase` | `{APIURL}/showtime/v1` | Showtime |
| `SprintsBase` | `{SprintsURL}/zsapi` | Sprints |
| `VaultBase` | `{APIURL}/vault/v1` | Vault |
| `VoiceBase` | `{APIURL}/voice/v1` | Voice |

Where `{APIURL}` = `https://www.zohoapis.{tld}` (generic Zoho API domain).

**IMPORTANT**: The base URL versions above are best-guesses. You MUST verify the correct API version and path via Context7 docs. For example, Inventory might be v1 or v3 depending on the endpoint.

## Integration Test Pattern

Tests go in `internal/<module>_integration_test.go`. They:
- Have `//go:build integration` tag
- Use `package internal_test`
- Shell out to the `./zoho` binary via helper functions
- Test real CRUD flows against Zoho APIs
- Use cleanup trackers to delete test data

Available helpers from `helpers_test.go`:
- `zoho(t, args...)` — run CLI, fail on non-zero exit
- `zohoMayFail(t, args...)` — run CLI, return (stdout, error)
- `zohoIgnoreError(t, args...)` — run CLI, log but don't fail on error
- `runZoho(t, args...)` — run CLI, return Result{Stdout, Stderr, ExitCode}
- `parseJSON(t, s)` — parse JSON object
- `parseJSONArray(t, s)` — parse JSON array
- `toJSON(t, v)` — marshal to JSON string
- `testName(t)` — generate unique test name like "ZOHOTEST_1234_abc"
- `randomSuffix()` — 8-char hex suffix
- `requireID(t, id, msg)` — skip if id is empty
- `assertEqual(t, got, want)` — compare stringified values
- `assertContains(t, s, substr)` — check string contains substring
- `assertExitCode(t, r, code)` — check exit code
- `assertNonEmpty(t, arr, msg)` — check array not empty
- `newCleanup(t)` — create cleanup tracker
- `truncate(s, n)` — truncate string for error messages
- `testPrefix` = "ZOHOTEST"

Test file structure:
```go
//go:build integration

package internal_test

import (
    "fmt"
    "os"
    "strings"
    "testing"
)

// Cleanup tracker methods (if module creates resources)
func (c *testCleanup) trackMyResource(id string) {
    c.add("delete my resource "+id, func() {
        zohoIgnoreError(c.t, "mymodule", "resources", "delete", id)
    })
}

// Helper to get required env var
func requireMyModuleOrgID(t *testing.T) string {
    t.Helper()
    id := os.Getenv("ZOHO_MYMODULE_ORG_ID")
    if id == "" {
        t.Skip("skipping: ZOHO_MYMODULE_ORG_ID not set")
    }
    return id
}

func TestMyModuleResources(t *testing.T) {
    t.Parallel()
    cleanup := newCleanup(t)

    var resourceID string

    t.Run("list", func(t *testing.T) {
        out := zoho(t, "mymodule", "resources", "list")
        m := parseJSON(t, out)
        // verify response structure
    })

    t.Run("create", func(t *testing.T) {
        name := fmt.Sprintf("%s Res %s", testPrefix, randomSuffix())
        out := zoho(t, "mymodule", "resources", "create",
            "--json", toJSON(t, map[string]any{"name": name}))
        m := parseJSON(t, out)
        resourceID = fmt.Sprintf("%v", m["id"])
        cleanup.trackMyResource(resourceID)
    })

    t.Run("get", func(t *testing.T) {
        requireID(t, resourceID, "create must have succeeded")
        out := zoho(t, "mymodule", "resources", "get", resourceID)
        m := parseJSON(t, out)
        // verify fields
    })
}

func TestMyModuleErrors(t *testing.T) {
    t.Parallel()
    // Test validation errors, missing flags, bad IDs
    t.Run("missing-required", func(t *testing.T) {
        r := runZoho(t, "mymodule", "resources", "get")
        assertExitCode(t, r, 4)
    })
}
```

## Context7 Library IDs

Use these to look up API docs:
- Bigin: `/websites/bigin_developer_apis_v2`
- Recruit: `/websites/zoho_recruit_developer-guide_apiv2`
- People: `/websites/zoho_people_api`
- Inventory: `/websites/zoho_inventory_api_v1`
- Invoice: `/websites/zoho_invoice_api_v3`
- Billing: `/websites/zoho_billing_api_v1_introduction`
- Analytics: (search for "zoho analytics api")
- Campaigns: (search for "zoho campaigns api")
- Creator: `/websites/zoho_creator_help_api_v2_1`
- Sprints: `/websites/sprints_zoho_apidoc`
- SalesIQ: `/websites/zoho_salesiq_help`
- Learn, Meeting, Bookings, Assist, PageSense, Showtime, Backstage, Vault, Voice, Marketing Automation: Search Context7

## Key Conventions
- No comments in Go code
- Use `--json` flag for complex POST/PUT bodies (not individual flags)
- Use `output.JSONRaw(raw)` for all responses
- Use `internal.NewValidationError()` for validation errors
- Pass through raw Zoho API responses — no data transformation
- For products needing org ID: use `--org` flag with `ZOHO_<PRODUCT>_ORG_ID` env var fallback
- Org ID shared across Books/Expense/Inventory/Invoice/Billing = `916897749` (set as `ZOHO_BOOKS_ORG_ID`)

## Environment Variables Available in Tests
```
ZOHO_DC=com
ZOHO_TEAM_ID=otynj0f288fbba7f04c6cb7dcb94c4f6dac2d
ZOHO_PORTAL_ID=916898406
ZOHO_EXPENSE_ORG_ID=916897749
ZOHO_DESK_ORG_ID=916897754
ZOHO_BOOKS_ORG_ID=916897749
```

## Products That May Have Limited/No Public REST API
Vault, Voice, Backstage, Showtime, Meeting, PageSense, Assist, Bookings, Learn — these may have very thin or no documented REST API. If you can't find API docs via Context7, report back with what you found. Don't make up endpoints.
