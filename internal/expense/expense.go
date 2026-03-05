package expense

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

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

func resolveOrgID(cmd *cli.Command) (string, error) {
	org := cmd.String("org")
	if org == "" {
		org = os.Getenv("ZOHO_EXPENSE_ORG_ID")
	}
	if org == "" {
		return "", internal.NewValidationError("--org flag or ZOHO_EXPENSE_ORG_ID env var required")
	}
	return org, nil
}

func orgHeaders(orgID string) map[string]string {
	return map[string]string{
		"X-com-zoho-expense-organizationid": orgID,
	}
}

func v3Base(c *zohttp.Client) string {
	return strings.Replace(c.ExpenseBase, "/v1", "/v3", 1)
}

func v3OrgParams(orgID string) map[string]string {
	return map[string]string{
		"organization_id": orgID,
	}
}

func Commands() *cli.Command {
	return &cli.Command{
		Name:  "expense",
		Usage: "Zoho Expense operations",
		Flags: []cli.Flag{
			&cli.StringFlag{Name: "org", Usage: "Organization ID (or set ZOHO_EXPENSE_ORG_ID)"},
		},
		Commands: []*cli.Command{
			organizationsCmd(),
			expensesCmd(),
			reportsCmd(),
			categoriesCmd(),
			usersCmd(),
			customersCmd(),
			projectsCmd(),
			tripsCmd(),
			currenciesCmd(),
			taxesCmd(),
			receiptsCmd(),
			tagsCmd(),
		},
	}
}

func organizationsCmd() *cli.Command {
	return &cli.Command{
		Name:  "organizations",
		Usage: "Organization operations",
		Commands: []*cli.Command{
			{
				Name:  "list",
				Usage: "List organizations",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					raw, err := c.Request("GET", c.ExpenseBase+"/organizations", nil)
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "get",
				Usage:     "Get an organization",
				ArgsUsage: "<org-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					raw, err := c.Request("GET", c.ExpenseBase+"/organizations/"+cmd.Args().First(), nil)
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:  "create",
				Usage: "Create an organization",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "json", Required: true, Usage: "JSON body"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					var body any
					json.Unmarshal([]byte(cmd.String("json")), &body)
					raw, err := c.Request("POST", c.ExpenseBase+"/organizations", &zohttp.RequestOpts{
						JSON: body,
					})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "update",
				Usage:     "Update an organization",
				ArgsUsage: "<org-id>",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "json", Required: true, Usage: "JSON body"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					var body any
					json.Unmarshal([]byte(cmd.String("json")), &body)
					raw, err := c.Request("PUT", c.ExpenseBase+"/organizations/"+cmd.Args().First(), &zohttp.RequestOpts{
						JSON: body,
					})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
		},
	}
}

func expensesCmd() *cli.Command {
	return &cli.Command{
		Name:  "expenses",
		Usage: "Expense operations",
		Commands: []*cli.Command{
			{
				Name:  "list",
				Usage: "List expenses",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "status", Usage: "Filter by status"},
					&cli.StringFlag{Name: "date-start", Usage: "Start date"},
					&cli.StringFlag{Name: "date-end", Usage: "End date"},
					&cli.StringFlag{Name: "user-id", Usage: "Filter by user ID"},
					&cli.StringFlag{Name: "category-id", Usage: "Filter by category ID"},
					&cli.StringFlag{Name: "merchant-id", Usage: "Filter by merchant ID"},
					&cli.StringFlag{Name: "customer-id", Usage: "Filter by customer ID"},
					&cli.StringFlag{Name: "project-id", Usage: "Filter by project ID"},
					&cli.StringFlag{Name: "page", Usage: "Page number"},
					&cli.StringFlag{Name: "per-page", Usage: "Results per page"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					params := map[string]string{}
					if v := cmd.String("status"); v != "" {
						params["status"] = v
					}
					if v := cmd.String("date-start"); v != "" {
						params["date_start"] = v
					}
					if v := cmd.String("date-end"); v != "" {
						params["date_end"] = v
					}
					if v := cmd.String("user-id"); v != "" {
						params["user_id"] = v
					}
					if v := cmd.String("category-id"); v != "" {
						params["category_id"] = v
					}
					if v := cmd.String("merchant-id"); v != "" {
						params["merchant_id"] = v
					}
					if v := cmd.String("customer-id"); v != "" {
						params["customer_id"] = v
					}
					if v := cmd.String("project-id"); v != "" {
						params["project_id"] = v
					}
					if v := cmd.String("page"); v != "" {
						params["page"] = v
					}
					if v := cmd.String("per-page"); v != "" {
						params["per_page"] = v
					}
					raw, err := c.Request("GET", c.ExpenseBase+"/reports/expensedetails", &zohttp.RequestOpts{
						Params:  params,
						Headers: orgHeaders(orgID),
					})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:  "create",
				Usage: "Create an expense",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "json", Required: true, Usage: "JSON body"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					var body any
					json.Unmarshal([]byte(cmd.String("json")), &body)
					raw, err := c.Request("POST", c.ExpenseBase+"/expenses", &zohttp.RequestOpts{
						JSON:    body,
						Headers: orgHeaders(orgID),
					})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "update",
				Usage:     "Update an expense",
				ArgsUsage: "<expense-id>",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "json", Required: true, Usage: "JSON body"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					var body any
					json.Unmarshal([]byte(cmd.String("json")), &body)
					raw, err := c.Request("PUT", c.ExpenseBase+"/expenses/"+cmd.Args().First(), &zohttp.RequestOpts{
						JSON:    body,
						Headers: orgHeaders(orgID),
					})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "merge",
				Usage:     "Merge a duplicate expense",
				ArgsUsage: "<expense-id>",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "duplicate-id", Required: true, Usage: "Duplicate expense ID to merge"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("POST", c.ExpenseBase+"/expenses/"+cmd.Args().First()+"/merge", &zohttp.RequestOpts{
						Params:  map[string]string{"duplicate_expense_id": cmd.String("duplicate-id")},
						Headers: orgHeaders(orgID),
					})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
		},
	}
}

func reportsCmd() *cli.Command {
	return &cli.Command{
		Name:  "reports",
		Usage: "Expense report operations",
		Commands: []*cli.Command{
			{
				Name:  "list",
				Usage: "List expense reports",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "filter-by", Usage: "Filter criteria"},
					&cli.StringFlag{Name: "page", Usage: "Page number"},
					&cli.StringFlag{Name: "per-page", Usage: "Results per page"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					params := map[string]string{}
					if v := cmd.String("filter-by"); v != "" {
						params["filter_by"] = v
					}
					if v := cmd.String("page"); v != "" {
						params["page"] = v
					}
					if v := cmd.String("per-page"); v != "" {
						params["per_page"] = v
					}
					raw, err := c.Request("GET", c.ExpenseBase+"/expensereports", &zohttp.RequestOpts{
						Params:  params,
						Headers: orgHeaders(orgID),
					})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "get",
				Usage:     "Get an expense report",
				ArgsUsage: "<report-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("GET", c.ExpenseBase+"/expensereports/"+cmd.Args().First(), &zohttp.RequestOpts{
						Headers: orgHeaders(orgID),
					})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:  "create",
				Usage: "Create an expense report",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "json", Required: true, Usage: "JSON body"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					var body any
					json.Unmarshal([]byte(cmd.String("json")), &body)
					raw, err := c.Request("POST", c.ExpenseBase+"/expensereports", &zohttp.RequestOpts{
						JSON:    body,
						Headers: orgHeaders(orgID),
					})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "update",
				Usage:     "Update an expense report",
				ArgsUsage: "<report-id>",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "json", Required: true, Usage: "JSON body"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					var body any
					json.Unmarshal([]byte(cmd.String("json")), &body)
					raw, err := c.Request("PUT", c.ExpenseBase+"/expensereports/"+cmd.Args().First(), &zohttp.RequestOpts{
						JSON:    body,
						Headers: orgHeaders(orgID),
					})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "approve",
				Usage:     "Approve an expense report",
				ArgsUsage: "<report-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("POST", c.ExpenseBase+"/expensereports/"+cmd.Args().First()+"/approve", &zohttp.RequestOpts{
						Headers: orgHeaders(orgID),
					})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "reject",
				Usage:     "Reject an expense report",
				ArgsUsage: "<report-id>",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "json", Usage: "JSON body with rejection comments"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					opts := &zohttp.RequestOpts{
						Headers: orgHeaders(orgID),
					}
					if j := cmd.String("json"); j != "" {
						var body any
						json.Unmarshal([]byte(j), &body)
						opts.JSON = body
					}
					raw, err := c.Request("POST", c.ExpenseBase+"/expensereports/"+cmd.Args().First()+"/reject", opts)
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "reimburse",
				Usage:     "Reimburse an expense report",
				ArgsUsage: "<report-id>",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "json", Usage: "JSON body"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					opts := &zohttp.RequestOpts{
						Headers: orgHeaders(orgID),
					}
					if j := cmd.String("json"); j != "" {
						var body any
						json.Unmarshal([]byte(j), &body)
						opts.JSON = body
					}
					raw, err := c.Request("POST", c.ExpenseBase+"/expensereports/"+cmd.Args().First()+"/reimburse", opts)
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "approval-history",
				Usage:     "Get approval history of an expense report",
				ArgsUsage: "<report-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("GET", c.ExpenseBase+"/expensereports/"+cmd.Args().First()+"/approvalhistory", &zohttp.RequestOpts{
						Headers: orgHeaders(orgID),
					})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
		},
	}
}

func categoriesCmd() *cli.Command {
	return &cli.Command{
		Name:  "categories",
		Usage: "Expense category operations",
		Commands: []*cli.Command{
			{
				Name:  "list",
				Usage: "List expense categories",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "filter-by", Usage: "Filter criteria"},
					&cli.StringFlag{Name: "page", Usage: "Page number"},
					&cli.StringFlag{Name: "per-page", Usage: "Results per page"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					params := map[string]string{}
					if v := cmd.String("filter-by"); v != "" {
						params["filter_by"] = v
					}
					if v := cmd.String("page"); v != "" {
						params["page"] = v
					}
					if v := cmd.String("per-page"); v != "" {
						params["per_page"] = v
					}
					raw, err := c.Request("GET", c.ExpenseBase+"/expensecategories", &zohttp.RequestOpts{
						Params:  params,
						Headers: orgHeaders(orgID),
					})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "get",
				Usage:     "Get an expense category",
				ArgsUsage: "<category-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("GET", c.ExpenseBase+"/expensecategories/"+cmd.Args().First(), &zohttp.RequestOpts{
						Headers: orgHeaders(orgID),
					})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:  "create",
				Usage: "Create an expense category",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "json", Required: true, Usage: "JSON body"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					var body any
					json.Unmarshal([]byte(cmd.String("json")), &body)
					raw, err := c.Request("POST", c.ExpenseBase+"/expensecategories", &zohttp.RequestOpts{
						JSON:    body,
						Headers: orgHeaders(orgID),
					})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "update",
				Usage:     "Update an expense category",
				ArgsUsage: "<category-id>",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "json", Required: true, Usage: "JSON body"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					var body any
					json.Unmarshal([]byte(cmd.String("json")), &body)
					raw, err := c.Request("PUT", c.ExpenseBase+"/expensecategories/"+cmd.Args().First(), &zohttp.RequestOpts{
						JSON:    body,
						Headers: orgHeaders(orgID),
					})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "delete",
				Usage:     "Delete an expense category",
				ArgsUsage: "<category-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("DELETE", c.ExpenseBase+"/expensecategories/"+cmd.Args().First(), &zohttp.RequestOpts{
						Headers: orgHeaders(orgID),
					})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "enable",
				Usage:     "Enable an expense category",
				ArgsUsage: "<category-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("POST", c.ExpenseBase+"/expensecategories/"+cmd.Args().First()+"/show", &zohttp.RequestOpts{
						Headers: orgHeaders(orgID),
					})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "disable",
				Usage:     "Disable an expense category",
				ArgsUsage: "<category-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("POST", c.ExpenseBase+"/expensecategories/"+cmd.Args().First()+"/hide", &zohttp.RequestOpts{
						Headers: orgHeaders(orgID),
					})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
		},
	}
}

func usersCmd() *cli.Command {
	return &cli.Command{
		Name:  "users",
		Usage: "User operations",
		Commands: []*cli.Command{
			{
				Name:  "list",
				Usage: "List users",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "page", Usage: "Page number"},
					&cli.StringFlag{Name: "per-page", Usage: "Results per page"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					params := map[string]string{}
					if v := cmd.String("page"); v != "" {
						params["page"] = v
					}
					if v := cmd.String("per-page"); v != "" {
						params["per_page"] = v
					}
					raw, err := c.Request("GET", c.ExpenseBase+"/users", &zohttp.RequestOpts{
						Params:  params,
						Headers: orgHeaders(orgID),
					})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "get",
				Usage:     "Get a user",
				ArgsUsage: "<user-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("GET", c.ExpenseBase+"/users/"+cmd.Args().First(), &zohttp.RequestOpts{
						Headers: orgHeaders(orgID),
					})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:  "create",
				Usage: "Create a user",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "json", Required: true, Usage: "JSON body"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					var body any
					json.Unmarshal([]byte(cmd.String("json")), &body)
					raw, err := c.Request("POST", c.ExpenseBase+"/users", &zohttp.RequestOpts{
						JSON:    body,
						Headers: orgHeaders(orgID),
					})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "update",
				Usage:     "Update a user",
				ArgsUsage: "<user-id>",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "json", Required: true, Usage: "JSON body"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					var body any
					json.Unmarshal([]byte(cmd.String("json")), &body)
					raw, err := c.Request("PUT", c.ExpenseBase+"/users/"+cmd.Args().First(), &zohttp.RequestOpts{
						JSON:    body,
						Headers: orgHeaders(orgID),
					})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "delete",
				Usage:     "Delete a user",
				ArgsUsage: "<user-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("DELETE", c.ExpenseBase+"/users/"+cmd.Args().First(), &zohttp.RequestOpts{
						Headers: orgHeaders(orgID),
					})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "activate",
				Usage:     "Activate a user",
				ArgsUsage: "<user-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("POST", c.ExpenseBase+"/users/"+cmd.Args().First()+"/active", &zohttp.RequestOpts{
						Headers: orgHeaders(orgID),
					})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "deactivate",
				Usage:     "Deactivate a user",
				ArgsUsage: "<user-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("POST", c.ExpenseBase+"/users/"+cmd.Args().First()+"/inactive", &zohttp.RequestOpts{
						Headers: orgHeaders(orgID),
					})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "assign-role",
				Usage:     "Assign a role to a user",
				ArgsUsage: "<user-id> <role-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("POST", c.ExpenseBase+"/users/"+cmd.Args().First()+"/role/"+cmd.Args().Get(1), &zohttp.RequestOpts{
						Headers: orgHeaders(orgID),
					})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
		},
	}
}

func customersCmd() *cli.Command {
	return &cli.Command{
		Name:  "customers",
		Usage: "Customer (contact) operations",
		Commands: []*cli.Command{
			{
				Name:  "list",
				Usage: "List customers",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "page", Usage: "Page number"},
					&cli.StringFlag{Name: "per-page", Usage: "Results per page"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					params := map[string]string{}
					if v := cmd.String("page"); v != "" {
						params["page"] = v
					}
					if v := cmd.String("per-page"); v != "" {
						params["per_page"] = v
					}
					raw, err := c.Request("GET", c.ExpenseBase+"/contacts", &zohttp.RequestOpts{
						Params:  params,
						Headers: orgHeaders(orgID),
					})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "get",
				Usage:     "Get a customer",
				ArgsUsage: "<customer-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("GET", c.ExpenseBase+"/contacts/"+cmd.Args().First(), &zohttp.RequestOpts{
						Headers: orgHeaders(orgID),
					})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:  "create",
				Usage: "Create a customer",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "json", Required: true, Usage: "JSON body"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					var body any
					json.Unmarshal([]byte(cmd.String("json")), &body)
					raw, err := c.Request("POST", c.ExpenseBase+"/contacts", &zohttp.RequestOpts{
						JSON:    body,
						Headers: orgHeaders(orgID),
					})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "update",
				Usage:     "Update a customer",
				ArgsUsage: "<customer-id>",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "json", Required: true, Usage: "JSON body"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					var body any
					json.Unmarshal([]byte(cmd.String("json")), &body)
					raw, err := c.Request("PUT", c.ExpenseBase+"/contacts/"+cmd.Args().First(), &zohttp.RequestOpts{
						JSON:    body,
						Headers: orgHeaders(orgID),
					})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "delete",
				Usage:     "Delete a customer",
				ArgsUsage: "<customer-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("DELETE", c.ExpenseBase+"/contacts/"+cmd.Args().First(), &zohttp.RequestOpts{
						Headers: orgHeaders(orgID),
					})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
		},
	}
}

func projectsCmd() *cli.Command {
	return &cli.Command{
		Name:  "projects",
		Usage: "Expense project operations",
		Commands: []*cli.Command{
			{
				Name:  "list",
				Usage: "List projects",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "page", Usage: "Page number"},
					&cli.StringFlag{Name: "per-page", Usage: "Results per page"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					params := map[string]string{}
					if v := cmd.String("page"); v != "" {
						params["page"] = v
					}
					if v := cmd.String("per-page"); v != "" {
						params["per_page"] = v
					}
					raw, err := c.Request("GET", c.ExpenseBase+"/projects", &zohttp.RequestOpts{
						Params:  params,
						Headers: orgHeaders(orgID),
					})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "get",
				Usage:     "Get a project",
				ArgsUsage: "<project-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("GET", c.ExpenseBase+"/projects/"+cmd.Args().First(), &zohttp.RequestOpts{
						Headers: orgHeaders(orgID),
					})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:  "create",
				Usage: "Create a project",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "json", Required: true, Usage: "JSON body"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					var body any
					json.Unmarshal([]byte(cmd.String("json")), &body)
					raw, err := c.Request("POST", c.ExpenseBase+"/projects", &zohttp.RequestOpts{
						JSON:    body,
						Headers: orgHeaders(orgID),
					})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "update",
				Usage:     "Update a project",
				ArgsUsage: "<project-id>",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "json", Required: true, Usage: "JSON body"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					var body any
					json.Unmarshal([]byte(cmd.String("json")), &body)
					raw, err := c.Request("PUT", c.ExpenseBase+"/projects/"+cmd.Args().First(), &zohttp.RequestOpts{
						JSON:    body,
						Headers: orgHeaders(orgID),
					})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "delete",
				Usage:     "Delete a project",
				ArgsUsage: "<project-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("DELETE", c.ExpenseBase+"/projects/"+cmd.Args().First(), &zohttp.RequestOpts{
						Headers: orgHeaders(orgID),
					})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "activate",
				Usage:     "Activate a project",
				ArgsUsage: "<project-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("POST", c.ExpenseBase+"/projects/"+cmd.Args().First()+"/active", &zohttp.RequestOpts{
						Headers: orgHeaders(orgID),
					})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "deactivate",
				Usage:     "Deactivate a project",
				ArgsUsage: "<project-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("POST", c.ExpenseBase+"/projects/"+cmd.Args().First()+"/inactive", &zohttp.RequestOpts{
						Headers: orgHeaders(orgID),
					})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
		},
	}
}

func tripsCmd() *cli.Command {
	return &cli.Command{
		Name:  "trips",
		Usage: "Trip operations",
		Commands: []*cli.Command{
			{
				Name:  "list",
				Usage: "List trips",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "page", Usage: "Page number"},
					&cli.StringFlag{Name: "per-page", Usage: "Results per page"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					params := map[string]string{}
					if v := cmd.String("page"); v != "" {
						params["page"] = v
					}
					if v := cmd.String("per-page"); v != "" {
						params["per_page"] = v
					}
					raw, err := c.Request("GET", c.ExpenseBase+"/trips", &zohttp.RequestOpts{
						Params:  params,
						Headers: orgHeaders(orgID),
					})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "get",
				Usage:     "Get a trip",
				ArgsUsage: "<trip-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("GET", c.ExpenseBase+"/trips/"+cmd.Args().First(), &zohttp.RequestOpts{
						Headers: orgHeaders(orgID),
					})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:  "create",
				Usage: "Create a trip",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "json", Required: true, Usage: "JSON body"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					var body any
					json.Unmarshal([]byte(cmd.String("json")), &body)
					raw, err := c.Request("POST", c.ExpenseBase+"/trips", &zohttp.RequestOpts{
						JSON:    body,
						Headers: orgHeaders(orgID),
					})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "update",
				Usage:     "Update a trip",
				ArgsUsage: "<trip-id>",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "json", Required: true, Usage: "JSON body"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					var body any
					json.Unmarshal([]byte(cmd.String("json")), &body)
					raw, err := c.Request("PUT", c.ExpenseBase+"/trips/"+cmd.Args().First(), &zohttp.RequestOpts{
						JSON:    body,
						Headers: orgHeaders(orgID),
					})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "delete",
				Usage:     "Delete a trip",
				ArgsUsage: "<trip-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("DELETE", c.ExpenseBase+"/trips/"+cmd.Args().First(), &zohttp.RequestOpts{
						Headers: orgHeaders(orgID),
					})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "approve",
				Usage:     "Approve a trip",
				ArgsUsage: "<trip-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("POST", c.ExpenseBase+"/trips/"+cmd.Args().First()+"/approve", &zohttp.RequestOpts{
						Headers: orgHeaders(orgID),
					})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "reject",
				Usage:     "Reject a trip",
				ArgsUsage: "<trip-id>",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "json", Usage: "JSON body with rejection comments"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					opts := &zohttp.RequestOpts{
						Headers: orgHeaders(orgID),
					}
					if j := cmd.String("json"); j != "" {
						var body any
						json.Unmarshal([]byte(j), &body)
						opts.JSON = body
					}
					raw, err := c.Request("POST", c.ExpenseBase+"/trips/"+cmd.Args().First()+"/reject", opts)
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "cancel",
				Usage:     "Cancel a trip",
				ArgsUsage: "<trip-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("POST", c.ExpenseBase+"/trips/"+cmd.Args().First()+"/cancel", &zohttp.RequestOpts{
						Headers: orgHeaders(orgID),
					})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "close",
				Usage:     "Close a trip",
				ArgsUsage: "<trip-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("POST", c.ExpenseBase+"/trips/"+cmd.Args().First()+"/close", &zohttp.RequestOpts{
						Headers: orgHeaders(orgID),
					})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
		},
	}
}

func currenciesCmd() *cli.Command {
	return &cli.Command{
		Name:  "currencies",
		Usage: "Currency settings operations",
		Commands: []*cli.Command{
			{
				Name:  "list",
				Usage: "List currencies",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("GET", c.ExpenseBase+"/settings/currencies", &zohttp.RequestOpts{
						Headers: orgHeaders(orgID),
					})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "get",
				Usage:     "Get a currency",
				ArgsUsage: "<currency-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("GET", c.ExpenseBase+"/settings/currencies/"+cmd.Args().First(), &zohttp.RequestOpts{
						Headers: orgHeaders(orgID),
					})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:  "create",
				Usage: "Create a currency",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "json", Required: true, Usage: "JSON body"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					var body any
					json.Unmarshal([]byte(cmd.String("json")), &body)
					raw, err := c.Request("POST", c.ExpenseBase+"/settings/currencies", &zohttp.RequestOpts{
						JSON:    body,
						Headers: orgHeaders(orgID),
					})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "update",
				Usage:     "Update a currency",
				ArgsUsage: "<currency-id>",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "json", Required: true, Usage: "JSON body"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					var body any
					json.Unmarshal([]byte(cmd.String("json")), &body)
					raw, err := c.Request("PUT", c.ExpenseBase+"/settings/currencies/"+cmd.Args().First(), &zohttp.RequestOpts{
						JSON:    body,
						Headers: orgHeaders(orgID),
					})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "delete",
				Usage:     "Delete a currency",
				ArgsUsage: "<currency-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("DELETE", c.ExpenseBase+"/settings/currencies/"+cmd.Args().First(), &zohttp.RequestOpts{
						Headers: orgHeaders(orgID),
					})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
		},
	}
}

func taxesCmd() *cli.Command {
	return &cli.Command{
		Name:  "taxes",
		Usage: "Tax settings operations",
		Commands: []*cli.Command{
			{
				Name:  "list",
				Usage: "List taxes",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("GET", c.ExpenseBase+"/settings/taxes", &zohttp.RequestOpts{
						Headers: orgHeaders(orgID),
					})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "get",
				Usage:     "Get a tax",
				ArgsUsage: "<tax-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("GET", c.ExpenseBase+"/settings/taxes/"+cmd.Args().First(), &zohttp.RequestOpts{
						Headers: orgHeaders(orgID),
					})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:  "create",
				Usage: "Create a tax",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "json", Required: true, Usage: "JSON body"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					var body any
					json.Unmarshal([]byte(cmd.String("json")), &body)
					raw, err := c.Request("POST", c.ExpenseBase+"/settings/taxes", &zohttp.RequestOpts{
						JSON:    body,
						Headers: orgHeaders(orgID),
					})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "update",
				Usage:     "Update a tax",
				ArgsUsage: "<tax-id>",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "json", Required: true, Usage: "JSON body"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					var body any
					json.Unmarshal([]byte(cmd.String("json")), &body)
					raw, err := c.Request("PUT", c.ExpenseBase+"/settings/taxes/"+cmd.Args().First(), &zohttp.RequestOpts{
						JSON:    body,
						Headers: orgHeaders(orgID),
					})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "delete",
				Usage:     "Delete a tax",
				ArgsUsage: "<tax-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("DELETE", c.ExpenseBase+"/settings/taxes/"+cmd.Args().First(), &zohttp.RequestOpts{
						Headers: orgHeaders(orgID),
					})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "get-group",
				Usage:     "Get a tax group",
				ArgsUsage: "<group-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("GET", c.ExpenseBase+"/settings/taxgroups/"+cmd.Args().First(), &zohttp.RequestOpts{
						Headers: orgHeaders(orgID),
					})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
		},
	}
}

func receiptsCmd() *cli.Command {
	return &cli.Command{
		Name:  "receipts",
		Usage: "Receipt operations",
		Commands: []*cli.Command{
			{
				Name:      "upload",
				Usage:     "Upload a receipt for an expense",
				ArgsUsage: "<expense-id>",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "file", Required: true, Usage: "Path to receipt file"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					filePath := cmd.String("file")
					data, err := os.ReadFile(filePath)
					if err != nil {
						return fmt.Errorf("failed to read file: %w", err)
					}
					name := filepath.Base(filePath)
					raw, err := c.Request("POST", c.ExpenseBase+"/expenses", &zohttp.RequestOpts{
						Files:   map[string]zohttp.FileUpload{"receipt": {Filename: name, Data: data}},
						Form:    map[string]string{"expense_id": cmd.Args().First()},
						Headers: orgHeaders(orgID),
					})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
		},
	}
}

func tagsCmd() *cli.Command {
	return &cli.Command{
		Name:  "tags",
		Usage: "Reporting tag operations (V3 API)",
		Commands: []*cli.Command{
			{
				Name:  "list",
				Usage: "List reporting tags",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("GET", v3Base(c)+"/reportingtags", &zohttp.RequestOpts{
						Params: v3OrgParams(orgID),
					})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:  "create",
				Usage: "Create a reporting tag",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "json", Required: true, Usage: "JSON body"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					var body any
					json.Unmarshal([]byte(cmd.String("json")), &body)
					raw, err := c.Request("POST", v3Base(c)+"/reportingtags", &zohttp.RequestOpts{
						JSON:   body,
						Params: v3OrgParams(orgID),
					})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "get",
				Usage:     "Get a reporting tag",
				ArgsUsage: "<tag-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("GET", v3Base(c)+"/reportingtags/"+cmd.Args().First(), &zohttp.RequestOpts{
						Params: v3OrgParams(orgID),
					})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "update",
				Usage:     "Update a reporting tag",
				ArgsUsage: "<tag-id>",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "json", Required: true, Usage: "JSON body"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					var body any
					json.Unmarshal([]byte(cmd.String("json")), &body)
					raw, err := c.Request("PUT", v3Base(c)+"/reportingtags/"+cmd.Args().First(), &zohttp.RequestOpts{
						JSON:   body,
						Params: v3OrgParams(orgID),
					})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "delete",
				Usage:     "Delete a reporting tag",
				ArgsUsage: "<tag-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("DELETE", v3Base(c)+"/reportingtags/"+cmd.Args().First(), &zohttp.RequestOpts{
						Params: v3OrgParams(orgID),
					})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "update-options",
				Usage:     "Update reporting tag options",
				ArgsUsage: "<tag-id>",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "json", Required: true, Usage: "JSON body"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					var body any
					json.Unmarshal([]byte(cmd.String("json")), &body)
					raw, err := c.Request("PUT", v3Base(c)+"/reportingtags/"+cmd.Args().First()+"/options", &zohttp.RequestOpts{
						JSON:   body,
						Params: v3OrgParams(orgID),
					})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "update-criteria",
				Usage:     "Update reporting tag criteria",
				ArgsUsage: "<tag-id>",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "json", Required: true, Usage: "JSON body"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					var body any
					json.Unmarshal([]byte(cmd.String("json")), &body)
					raw, err := c.Request("PUT", v3Base(c)+"/reportingtags/"+cmd.Args().First()+"/criteria", &zohttp.RequestOpts{
						JSON:   body,
						Params: v3OrgParams(orgID),
					})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "activate",
				Usage:     "Activate a reporting tag",
				ArgsUsage: "<tag-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("POST", v3Base(c)+"/reportingtags/"+cmd.Args().First()+"/active", &zohttp.RequestOpts{
						Params: v3OrgParams(orgID),
					})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "deactivate",
				Usage:     "Deactivate a reporting tag",
				ArgsUsage: "<tag-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("POST", v3Base(c)+"/reportingtags/"+cmd.Args().First()+"/inactive", &zohttp.RequestOpts{
						Params: v3OrgParams(orgID),
					})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "activate-option",
				Usage:     "Activate a reporting tag option",
				ArgsUsage: "<tag-id> <option-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("POST", v3Base(c)+"/reportingtags/"+cmd.Args().First()+"/option/"+cmd.Args().Get(1)+"/active", &zohttp.RequestOpts{
						Params: v3OrgParams(orgID),
					})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "deactivate-option",
				Usage:     "Deactivate a reporting tag option",
				ArgsUsage: "<tag-id> <option-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("POST", v3Base(c)+"/reportingtags/"+cmd.Args().First()+"/option/"+cmd.Args().Get(1)+"/inactive", &zohttp.RequestOpts{
						Params: v3OrgParams(orgID),
					})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:  "list-options",
				Usage: "List all reporting tag options across tags",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("GET", v3Base(c)+"/reportingtags/options", &zohttp.RequestOpts{
						Params: v3OrgParams(orgID),
					})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "list-tag-options",
				Usage:     "List options for a specific reporting tag",
				ArgsUsage: "<tag-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("GET", v3Base(c)+"/reportingtags/"+cmd.Args().First()+"/options/all", &zohttp.RequestOpts{
						Params: v3OrgParams(orgID),
					})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:  "reorder",
				Usage: "Reorder reporting tags",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "json", Required: true, Usage: "JSON body"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					var body any
					json.Unmarshal([]byte(cmd.String("json")), &body)
					raw, err := c.Request("PUT", v3Base(c)+"/reportingtags/reorder", &zohttp.RequestOpts{
						JSON:   body,
						Params: v3OrgParams(orgID),
					})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
		},
	}
}
