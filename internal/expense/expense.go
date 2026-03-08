package expense

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/omin8tor/zoho-cli/internal"
	zohttp "github.com/omin8tor/zoho-cli/internal/http"
	"github.com/omin8tor/zoho-cli/internal/output"
	"github.com/omin8tor/zoho-cli/internal/pagination"
	"github.com/urfave/cli/v3"
)

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
			&cli.StringFlag{Name: "org", Sources: cli.EnvVars("ZOHO_EXPENSE_ORG_ID"), Usage: "Organization ID (or set ZOHO_EXPENSE_ORG_ID)"},
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
					c, err := zohttp.GetClient()
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
					c, err := zohttp.GetClient()
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
					&cli.StringFlag{Name: "name", Required: true, Usage: "Organization name"},
					&cli.StringFlag{Name: "fiscal_start_month", Usage: "Fiscal year start month"},
					&cli.StringFlag{Name: "time_zone", Usage: "Organization time zone"},
					&cli.StringFlag{Name: "json", Usage: "Additional fields as JSON"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					body := map[string]any{"name": cmd.String("name")}
					if cmd.IsSet("fiscal_start_month") {
						body["fiscal_start_month"] = cmd.String("fiscal_start_month")
					}
					if cmd.IsSet("time_zone") {
						body["time_zone"] = cmd.String("time_zone")
					}
					if err := internal.MergeJSON(cmd, body); err != nil {
						return err
					}
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
					&cli.StringFlag{Name: "name", Usage: "Organization name"},
					&cli.StringFlag{Name: "fiscal_start_month", Usage: "Fiscal year start month"},
					&cli.StringFlag{Name: "time_zone", Usage: "Organization time zone"},
					&cli.StringFlag{Name: "json", Usage: "Additional fields as JSON"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					body := map[string]any{}
					if cmd.IsSet("name") {
						body["name"] = cmd.String("name")
					}
					if cmd.IsSet("fiscal_start_month") {
						body["fiscal_start_month"] = cmd.String("fiscal_start_month")
					}
					if cmd.IsSet("time_zone") {
						body["time_zone"] = cmd.String("time_zone")
					}
					if err := internal.MergeJSON(cmd, body); err != nil {
						return err
					}
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
					&cli.BoolFlag{Name: "all", Usage: "Fetch all records"},
					&cli.IntFlag{Name: "limit", Usage: "Max total records to fetch"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := internal.RequireFlag(cmd, "org", "ZOHO_EXPENSE_ORG_ID")
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
					if cmd.Bool("all") || cmd.IsSet("limit") {
						items, err := pagination.Paginate(pagination.PaginationConfig{
							Client:   c,
							URL:      c.ExpenseBase + "/reports/expensedetails",
							Opts:     &zohttp.RequestOpts{Params: params, Headers: orgHeaders(orgID)},
							ItemsKey: "expenses",
							PageSize: 200,
							Limit:    int(cmd.Int("limit")),
							SetPage:  pagination.PagePerPage(200),
							HasMore:  pagination.HasMoreBooks,
						})
						if err != nil {
							return err
						}
						return output.JSON(items)
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
				Name:      "get",
				Usage:     "Get an expense",
				ArgsUsage: "<expense-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := internal.RequireFlag(cmd, "org", "ZOHO_EXPENSE_ORG_ID")
					if err != nil {
						return err
					}
					raw, err := c.Request("GET", c.ExpenseBase+"/expenses/"+cmd.Args().First(), &zohttp.RequestOpts{
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
					&cli.StringFlag{Name: "category_id", Required: true, Usage: "Expense category ID"},
					&cli.StringFlag{Name: "date", Required: true, Usage: "Expense date (YYYY-MM-DD)"},
					&cli.FloatFlag{Name: "amount", Required: true, Usage: "Line item amount"},
					&cli.StringFlag{Name: "currency_id", Usage: "Currency ID"},
					&cli.StringFlag{Name: "paid_through_account_id", Usage: "Paid through account ID"},
					&cli.StringFlag{Name: "description", Usage: "Expense description"},
					&cli.StringFlag{Name: "customer_id", Usage: "Customer ID"},
					&cli.StringFlag{Name: "project_id", Usage: "Project ID"},
					&cli.StringFlag{Name: "json", Usage: "Additional fields as JSON"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := internal.RequireFlag(cmd, "org", "ZOHO_EXPENSE_ORG_ID")
					if err != nil {
						return err
					}
					body := map[string]any{"date": cmd.String("date")}
					if cmd.IsSet("currency_id") {
						body["currency_id"] = cmd.String("currency_id")
					}
					if cmd.IsSet("paid_through_account_id") {
						body["paid_through_account_id"] = cmd.String("paid_through_account_id")
					}
					if cmd.IsSet("description") {
						body["description"] = cmd.String("description")
					}
					if cmd.IsSet("customer_id") {
						body["customer_id"] = cmd.String("customer_id")
					}
					if cmd.IsSet("project_id") {
						body["project_id"] = cmd.String("project_id")
					}
					if err := internal.MergeJSON(cmd, body); err != nil {
						return err
					}
					categoryID := cmd.String("category_id")
					if categoryID == "" {
						return internal.NewValidationError("--category_id and --amount are required")
					}
					amountVal := cmd.Float("amount")
					if _, exists := body["line_items"]; !exists {
						body["line_items"] = []map[string]any{{
							"category_id": categoryID,
							"amount":      amountVal,
						}}
					}
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
					&cli.StringFlag{Name: "category_id", Usage: "Expense category ID"},
					&cli.StringFlag{Name: "date", Usage: "Expense date (YYYY-MM-DD)"},
					&cli.FloatFlag{Name: "amount", Usage: "Line item amount"},
					&cli.StringFlag{Name: "currency_id", Usage: "Currency ID"},
					&cli.StringFlag{Name: "paid_through_account_id", Usage: "Paid through account ID"},
					&cli.StringFlag{Name: "description", Usage: "Expense description"},
					&cli.StringFlag{Name: "customer_id", Usage: "Customer ID"},
					&cli.StringFlag{Name: "project_id", Usage: "Project ID"},
					&cli.StringFlag{Name: "json", Usage: "Additional fields as JSON"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := internal.RequireFlag(cmd, "org", "ZOHO_EXPENSE_ORG_ID")
					if err != nil {
						return err
					}
					body := map[string]any{}
					if cmd.IsSet("date") {
						body["date"] = cmd.String("date")
					}
					if cmd.IsSet("currency_id") {
						body["currency_id"] = cmd.String("currency_id")
					}
					if cmd.IsSet("paid_through_account_id") {
						body["paid_through_account_id"] = cmd.String("paid_through_account_id")
					}
					if cmd.IsSet("description") {
						body["description"] = cmd.String("description")
					}
					if cmd.IsSet("customer_id") {
						body["customer_id"] = cmd.String("customer_id")
					}
					if cmd.IsSet("project_id") {
						body["project_id"] = cmd.String("project_id")
					}
					if err := internal.MergeJSON(cmd, body); err != nil {
						return err
					}
					categoryID := cmd.String("category_id")
					if cmd.IsSet("category_id") != cmd.IsSet("amount") {
						return internal.NewValidationError("--category_id and --amount must be provided together")
					}
					if cmd.IsSet("category_id") && cmd.IsSet("amount") {
						amountVal := cmd.Float("amount")
						if _, exists := body["line_items"]; !exists {
							body["line_items"] = []map[string]any{{
								"category_id": categoryID,
								"amount":      amountVal,
							}}
						}
					}
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
				Name:      "delete",
				Usage:     "Delete an expense",
				ArgsUsage: "<expense-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := internal.RequireFlag(cmd, "org", "ZOHO_EXPENSE_ORG_ID")
					if err != nil {
						return err
					}
					raw, err := c.Request("DELETE", c.ExpenseBase+"/expenses/"+cmd.Args().First(), &zohttp.RequestOpts{
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
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := internal.RequireFlag(cmd, "org", "ZOHO_EXPENSE_ORG_ID")
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
					&cli.BoolFlag{Name: "all", Usage: "Fetch all records"},
					&cli.IntFlag{Name: "limit", Usage: "Max total records to fetch"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := internal.RequireFlag(cmd, "org", "ZOHO_EXPENSE_ORG_ID")
					if err != nil {
						return err
					}
					params := map[string]string{}
					if v := cmd.String("filter-by"); v != "" {
						params["filter_by"] = v
					}
					if cmd.Bool("all") || cmd.IsSet("limit") {
						items, err := pagination.Paginate(pagination.PaginationConfig{
							Client:   c,
							URL:      c.ExpenseBase + "/expensereports",
							Opts:     &zohttp.RequestOpts{Params: params, Headers: orgHeaders(orgID)},
							ItemsKey: "expense_reports",
							PageSize: 200,
							Limit:    int(cmd.Int("limit")),
							SetPage:  pagination.PagePerPage(200),
							HasMore:  pagination.HasMoreBooks,
						})
						if err != nil {
							return err
						}
						return output.JSON(items)
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
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := internal.RequireFlag(cmd, "org", "ZOHO_EXPENSE_ORG_ID")
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
					&cli.StringFlag{Name: "report_name", Required: true, Usage: "Report name"},
					&cli.StringFlag{Name: "description", Usage: "Report description"},
					&cli.StringFlag{Name: "json", Usage: "Additional fields as JSON"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := internal.RequireFlag(cmd, "org", "ZOHO_EXPENSE_ORG_ID")
					if err != nil {
						return err
					}
					body := map[string]any{"report_name": cmd.String("report_name")}
					if cmd.IsSet("description") {
						body["description"] = cmd.String("description")
					}
					if err := internal.MergeJSON(cmd, body); err != nil {
						return err
					}
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
					&cli.StringFlag{Name: "report_name", Usage: "Report name"},
					&cli.StringFlag{Name: "description", Usage: "Report description"},
					&cli.StringFlag{Name: "json", Usage: "Additional fields as JSON"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := internal.RequireFlag(cmd, "org", "ZOHO_EXPENSE_ORG_ID")
					if err != nil {
						return err
					}
					body := map[string]any{}
					if cmd.IsSet("report_name") {
						body["report_name"] = cmd.String("report_name")
					}
					if cmd.IsSet("description") {
						body["description"] = cmd.String("description")
					}
					if err := internal.MergeJSON(cmd, body); err != nil {
						return err
					}
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
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := internal.RequireFlag(cmd, "org", "ZOHO_EXPENSE_ORG_ID")
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
					&cli.StringFlag{Name: "comments", Usage: "Rejection comments"},
					&cli.StringFlag{Name: "json", Usage: "Additional fields as JSON"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := internal.RequireFlag(cmd, "org", "ZOHO_EXPENSE_ORG_ID")
					if err != nil {
						return err
					}
					body := map[string]any{}
					if cmd.IsSet("comments") {
						body["comments"] = cmd.String("comments")
					}
					if err := internal.MergeJSON(cmd, body); err != nil {
						return err
					}
					opts := &zohttp.RequestOpts{Headers: orgHeaders(orgID)}
					if len(body) > 0 {
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
					&cli.StringFlag{Name: "payment_mode", Usage: "Payment mode"},
					&cli.StringFlag{Name: "reference_number", Usage: "Payment reference number"},
					&cli.StringFlag{Name: "json", Usage: "Additional fields as JSON"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := internal.RequireFlag(cmd, "org", "ZOHO_EXPENSE_ORG_ID")
					if err != nil {
						return err
					}
					body := map[string]any{}
					if cmd.IsSet("payment_mode") {
						body["payment_mode"] = cmd.String("payment_mode")
					}
					if cmd.IsSet("reference_number") {
						body["reference_number"] = cmd.String("reference_number")
					}
					if err := internal.MergeJSON(cmd, body); err != nil {
						return err
					}
					opts := &zohttp.RequestOpts{Headers: orgHeaders(orgID)}
					if len(body) > 0 {
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
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := internal.RequireFlag(cmd, "org", "ZOHO_EXPENSE_ORG_ID")
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
			{
				Name:      "delete",
				Usage:     "Delete an expense report",
				ArgsUsage: "<report-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := internal.RequireFlag(cmd, "org", "ZOHO_EXPENSE_ORG_ID")
					if err != nil {
						return err
					}
					raw, err := c.Request("DELETE", c.ExpenseBase+"/expensereports/"+cmd.Args().First(), &zohttp.RequestOpts{
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
					&cli.BoolFlag{Name: "all", Usage: "Fetch all records"},
					&cli.IntFlag{Name: "limit", Usage: "Max total records to fetch"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := internal.RequireFlag(cmd, "org", "ZOHO_EXPENSE_ORG_ID")
					if err != nil {
						return err
					}
					params := map[string]string{}
					if v := cmd.String("filter-by"); v != "" {
						params["filter_by"] = v
					}
					if cmd.Bool("all") || cmd.IsSet("limit") {
						items, err := pagination.Paginate(pagination.PaginationConfig{
							Client:   c,
							URL:      c.ExpenseBase + "/expensecategories",
							Opts:     &zohttp.RequestOpts{Params: params, Headers: orgHeaders(orgID)},
							ItemsKey: "categories",
							PageSize: 200,
							Limit:    int(cmd.Int("limit")),
							SetPage:  pagination.PagePerPage(200),
							HasMore:  pagination.HasMoreBooks,
						})
						if err != nil {
							return err
						}
						return output.JSON(items)
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
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := internal.RequireFlag(cmd, "org", "ZOHO_EXPENSE_ORG_ID")
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
					&cli.StringFlag{Name: "category_name", Required: true, Usage: "Category name"},
					&cli.StringFlag{Name: "description", Usage: "Category description"},
					&cli.StringFlag{Name: "json", Usage: "Additional fields as JSON"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := internal.RequireFlag(cmd, "org", "ZOHO_EXPENSE_ORG_ID")
					if err != nil {
						return err
					}
					body := map[string]any{"category_name": cmd.String("category_name")}
					if cmd.IsSet("description") {
						body["description"] = cmd.String("description")
					}
					if err := internal.MergeJSON(cmd, body); err != nil {
						return err
					}
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
					&cli.StringFlag{Name: "category_name", Usage: "Category name"},
					&cli.StringFlag{Name: "description", Usage: "Category description"},
					&cli.StringFlag{Name: "json", Usage: "Additional fields as JSON"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := internal.RequireFlag(cmd, "org", "ZOHO_EXPENSE_ORG_ID")
					if err != nil {
						return err
					}
					body := map[string]any{}
					if cmd.IsSet("category_name") {
						body["category_name"] = cmd.String("category_name")
					}
					if cmd.IsSet("description") {
						body["description"] = cmd.String("description")
					}
					if err := internal.MergeJSON(cmd, body); err != nil {
						return err
					}
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
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := internal.RequireFlag(cmd, "org", "ZOHO_EXPENSE_ORG_ID")
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
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := internal.RequireFlag(cmd, "org", "ZOHO_EXPENSE_ORG_ID")
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
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := internal.RequireFlag(cmd, "org", "ZOHO_EXPENSE_ORG_ID")
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
					&cli.BoolFlag{Name: "all", Usage: "Fetch all records"},
					&cli.IntFlag{Name: "limit", Usage: "Max total records to fetch"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := internal.RequireFlag(cmd, "org", "ZOHO_EXPENSE_ORG_ID")
					if err != nil {
						return err
					}
					params := map[string]string{}
					if cmd.Bool("all") || cmd.IsSet("limit") {
						items, err := pagination.Paginate(pagination.PaginationConfig{
							Client:   c,
							URL:      c.ExpenseBase + "/users",
							Opts:     &zohttp.RequestOpts{Params: params, Headers: orgHeaders(orgID)},
							ItemsKey: "users",
							PageSize: 200,
							Limit:    int(cmd.Int("limit")),
							SetPage:  pagination.PagePerPage(200),
							HasMore:  pagination.HasMoreBooks,
						})
						if err != nil {
							return err
						}
						return output.JSON(items)
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
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := internal.RequireFlag(cmd, "org", "ZOHO_EXPENSE_ORG_ID")
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
					&cli.StringFlag{Name: "name", Required: true, Usage: "User name"},
					&cli.StringFlag{Name: "email", Required: true, Usage: "User email"},
					&cli.StringFlag{Name: "role_id", Usage: "Role ID"},
					&cli.StringFlag{Name: "json", Usage: "Additional fields as JSON"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := internal.RequireFlag(cmd, "org", "ZOHO_EXPENSE_ORG_ID")
					if err != nil {
						return err
					}
					body := map[string]any{"name": cmd.String("name"), "email": cmd.String("email")}
					if cmd.IsSet("role_id") {
						body["role_id"] = cmd.String("role_id")
					}
					if err := internal.MergeJSON(cmd, body); err != nil {
						return err
					}
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
					&cli.StringFlag{Name: "name", Usage: "User name"},
					&cli.StringFlag{Name: "email", Usage: "User email"},
					&cli.StringFlag{Name: "role_id", Usage: "Role ID"},
					&cli.StringFlag{Name: "json", Usage: "Additional fields as JSON"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := internal.RequireFlag(cmd, "org", "ZOHO_EXPENSE_ORG_ID")
					if err != nil {
						return err
					}
					body := map[string]any{}
					if cmd.IsSet("name") {
						body["name"] = cmd.String("name")
					}
					if cmd.IsSet("email") {
						body["email"] = cmd.String("email")
					}
					if cmd.IsSet("role_id") {
						body["role_id"] = cmd.String("role_id")
					}
					if err := internal.MergeJSON(cmd, body); err != nil {
						return err
					}
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
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := internal.RequireFlag(cmd, "org", "ZOHO_EXPENSE_ORG_ID")
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
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := internal.RequireFlag(cmd, "org", "ZOHO_EXPENSE_ORG_ID")
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
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := internal.RequireFlag(cmd, "org", "ZOHO_EXPENSE_ORG_ID")
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
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := internal.RequireFlag(cmd, "org", "ZOHO_EXPENSE_ORG_ID")
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
					&cli.BoolFlag{Name: "all", Usage: "Fetch all records"},
					&cli.IntFlag{Name: "limit", Usage: "Max total records to fetch"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := internal.RequireFlag(cmd, "org", "ZOHO_EXPENSE_ORG_ID")
					if err != nil {
						return err
					}
					params := map[string]string{}
					if cmd.Bool("all") || cmd.IsSet("limit") {
						items, err := pagination.Paginate(pagination.PaginationConfig{
							Client:   c,
							URL:      c.ExpenseBase + "/contacts",
							Opts:     &zohttp.RequestOpts{Params: params, Headers: orgHeaders(orgID)},
							ItemsKey: "customers",
							PageSize: 200,
							Limit:    int(cmd.Int("limit")),
							SetPage:  pagination.PagePerPage(200),
							HasMore:  pagination.HasMoreBooks,
						})
						if err != nil {
							return err
						}
						return output.JSON(items)
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
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := internal.RequireFlag(cmd, "org", "ZOHO_EXPENSE_ORG_ID")
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
					&cli.StringFlag{Name: "contact_name", Required: true, Usage: "Customer name"},
					&cli.StringFlag{Name: "company_name", Usage: "Company name"},
					&cli.StringFlag{Name: "email", Usage: "Customer email"},
					&cli.StringFlag{Name: "json", Usage: "Additional fields as JSON"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := internal.RequireFlag(cmd, "org", "ZOHO_EXPENSE_ORG_ID")
					if err != nil {
						return err
					}
					body := map[string]any{"contact_name": cmd.String("contact_name")}
					if cmd.IsSet("company_name") {
						body["company_name"] = cmd.String("company_name")
					}
					if cmd.IsSet("email") {
						body["email"] = cmd.String("email")
					}
					if err := internal.MergeJSON(cmd, body); err != nil {
						return err
					}
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
					&cli.StringFlag{Name: "contact_name", Usage: "Customer name"},
					&cli.StringFlag{Name: "company_name", Usage: "Company name"},
					&cli.StringFlag{Name: "email", Usage: "Customer email"},
					&cli.StringFlag{Name: "json", Usage: "Additional fields as JSON"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := internal.RequireFlag(cmd, "org", "ZOHO_EXPENSE_ORG_ID")
					if err != nil {
						return err
					}
					body := map[string]any{}
					if cmd.IsSet("contact_name") {
						body["contact_name"] = cmd.String("contact_name")
					}
					if cmd.IsSet("company_name") {
						body["company_name"] = cmd.String("company_name")
					}
					if cmd.IsSet("email") {
						body["email"] = cmd.String("email")
					}
					if err := internal.MergeJSON(cmd, body); err != nil {
						return err
					}
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
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := internal.RequireFlag(cmd, "org", "ZOHO_EXPENSE_ORG_ID")
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
					&cli.BoolFlag{Name: "all", Usage: "Fetch all records"},
					&cli.IntFlag{Name: "limit", Usage: "Max total records to fetch"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := internal.RequireFlag(cmd, "org", "ZOHO_EXPENSE_ORG_ID")
					if err != nil {
						return err
					}
					params := map[string]string{}
					if cmd.Bool("all") || cmd.IsSet("limit") {
						items, err := pagination.Paginate(pagination.PaginationConfig{
							Client:   c,
							URL:      c.ExpenseBase + "/projects",
							Opts:     &zohttp.RequestOpts{Params: params, Headers: orgHeaders(orgID)},
							ItemsKey: "projects",
							PageSize: 200,
							Limit:    int(cmd.Int("limit")),
							SetPage:  pagination.PagePerPage(200),
							HasMore:  pagination.HasMoreBooks,
						})
						if err != nil {
							return err
						}
						return output.JSON(items)
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
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := internal.RequireFlag(cmd, "org", "ZOHO_EXPENSE_ORG_ID")
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
					&cli.StringFlag{Name: "project_name", Required: true, Usage: "Project name"},
					&cli.StringFlag{Name: "customer_id", Usage: "Customer ID"},
					&cli.StringFlag{Name: "description", Usage: "Project description"},
					&cli.StringFlag{Name: "json", Usage: "Additional fields as JSON"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := internal.RequireFlag(cmd, "org", "ZOHO_EXPENSE_ORG_ID")
					if err != nil {
						return err
					}
					body := map[string]any{"project_name": cmd.String("project_name")}
					if cmd.IsSet("customer_id") {
						body["customer_id"] = cmd.String("customer_id")
					}
					if cmd.IsSet("description") {
						body["description"] = cmd.String("description")
					}
					if err := internal.MergeJSON(cmd, body); err != nil {
						return err
					}
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
					&cli.StringFlag{Name: "project_name", Usage: "Project name"},
					&cli.StringFlag{Name: "customer_id", Usage: "Customer ID"},
					&cli.StringFlag{Name: "description", Usage: "Project description"},
					&cli.StringFlag{Name: "json", Usage: "Additional fields as JSON"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := internal.RequireFlag(cmd, "org", "ZOHO_EXPENSE_ORG_ID")
					if err != nil {
						return err
					}
					body := map[string]any{}
					if cmd.IsSet("project_name") {
						body["project_name"] = cmd.String("project_name")
					}
					if cmd.IsSet("customer_id") {
						body["customer_id"] = cmd.String("customer_id")
					}
					if cmd.IsSet("description") {
						body["description"] = cmd.String("description")
					}
					if err := internal.MergeJSON(cmd, body); err != nil {
						return err
					}
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
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := internal.RequireFlag(cmd, "org", "ZOHO_EXPENSE_ORG_ID")
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
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := internal.RequireFlag(cmd, "org", "ZOHO_EXPENSE_ORG_ID")
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
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := internal.RequireFlag(cmd, "org", "ZOHO_EXPENSE_ORG_ID")
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
					&cli.BoolFlag{Name: "all", Usage: "Fetch all records"},
					&cli.IntFlag{Name: "limit", Usage: "Max total records to fetch"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := internal.RequireFlag(cmd, "org", "ZOHO_EXPENSE_ORG_ID")
					if err != nil {
						return err
					}
					params := map[string]string{}
					if cmd.Bool("all") || cmd.IsSet("limit") {
						items, err := pagination.Paginate(pagination.PaginationConfig{
							Client:   c,
							URL:      c.ExpenseBase + "/trips",
							Opts:     &zohttp.RequestOpts{Params: params, Headers: orgHeaders(orgID)},
							ItemsKey: "trips",
							PageSize: 200,
							Limit:    int(cmd.Int("limit")),
							SetPage:  pagination.PagePerPage(200),
							HasMore:  pagination.HasMoreBooks,
						})
						if err != nil {
							return err
						}
						return output.JSON(items)
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
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := internal.RequireFlag(cmd, "org", "ZOHO_EXPENSE_ORG_ID")
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
					&cli.StringFlag{Name: "trip_number", Required: true, Usage: "Trip number"},
					&cli.BoolFlag{Name: "is_international", Required: true, Usage: "Whether the trip is international (true/false)"},
					&cli.StringFlag{Name: "business_purpose", Required: true, Usage: "Business purpose"},
					&cli.StringFlag{Name: "destination_country", Required: true, Usage: "Destination country"},
					&cli.StringFlag{Name: "json", Usage: "Additional fields as JSON"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := internal.RequireFlag(cmd, "org", "ZOHO_EXPENSE_ORG_ID")
					if err != nil {
						return err
					}
					body := map[string]any{
						"trip_number":         cmd.String("trip_number"),
						"is_international":    cmd.Bool("is_international"),
						"business_purpose":    cmd.String("business_purpose"),
						"destination_country": cmd.String("destination_country"),
					}
					if err := internal.MergeJSON(cmd, body); err != nil {
						return err
					}
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
					&cli.StringFlag{Name: "trip_number", Usage: "Trip number"},
					&cli.BoolFlag{Name: "is_international", Usage: "Whether the trip is international (true/false)"},
					&cli.StringFlag{Name: "business_purpose", Usage: "Business purpose"},
					&cli.StringFlag{Name: "destination_country", Usage: "Destination country"},
					&cli.StringFlag{Name: "json", Usage: "Additional fields as JSON"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := internal.RequireFlag(cmd, "org", "ZOHO_EXPENSE_ORG_ID")
					if err != nil {
						return err
					}
					body := map[string]any{}
					if cmd.IsSet("trip_number") {
						body["trip_number"] = cmd.String("trip_number")
					}
					if cmd.IsSet("is_international") {
						body["is_international"] = cmd.Bool("is_international")
					}
					if cmd.IsSet("business_purpose") {
						body["business_purpose"] = cmd.String("business_purpose")
					}
					if cmd.IsSet("destination_country") {
						body["destination_country"] = cmd.String("destination_country")
					}
					if err := internal.MergeJSON(cmd, body); err != nil {
						return err
					}
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
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := internal.RequireFlag(cmd, "org", "ZOHO_EXPENSE_ORG_ID")
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
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := internal.RequireFlag(cmd, "org", "ZOHO_EXPENSE_ORG_ID")
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
					&cli.StringFlag{Name: "comments", Usage: "Rejection comments"},
					&cli.StringFlag{Name: "json", Usage: "Additional fields as JSON"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := internal.RequireFlag(cmd, "org", "ZOHO_EXPENSE_ORG_ID")
					if err != nil {
						return err
					}
					body := map[string]any{}
					if cmd.IsSet("comments") {
						body["comments"] = cmd.String("comments")
					}
					if err := internal.MergeJSON(cmd, body); err != nil {
						return err
					}
					opts := &zohttp.RequestOpts{Headers: orgHeaders(orgID)}
					if len(body) > 0 {
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
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := internal.RequireFlag(cmd, "org", "ZOHO_EXPENSE_ORG_ID")
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
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := internal.RequireFlag(cmd, "org", "ZOHO_EXPENSE_ORG_ID")
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
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := internal.RequireFlag(cmd, "org", "ZOHO_EXPENSE_ORG_ID")
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
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := internal.RequireFlag(cmd, "org", "ZOHO_EXPENSE_ORG_ID")
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
					&cli.StringFlag{Name: "currency_code", Required: true, Usage: "Currency code"},
					&cli.StringFlag{Name: "currency_symbol", Usage: "Currency symbol"},
					&cli.StringFlag{Name: "currency_format", Usage: "Currency format"},
					&cli.IntFlag{Name: "price_precision", Usage: "Decimal precision"},
					&cli.StringFlag{Name: "json", Usage: "Additional fields as JSON"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := internal.RequireFlag(cmd, "org", "ZOHO_EXPENSE_ORG_ID")
					if err != nil {
						return err
					}
					body := map[string]any{"currency_code": cmd.String("currency_code")}
					if cmd.IsSet("currency_symbol") {
						body["currency_symbol"] = cmd.String("currency_symbol")
					}
					if cmd.IsSet("currency_format") {
						body["currency_format"] = cmd.String("currency_format")
					}
					if cmd.IsSet("price_precision") {
						body["price_precision"] = cmd.Int("price_precision")
					}
					if err := internal.MergeJSON(cmd, body); err != nil {
						return err
					}
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
					&cli.StringFlag{Name: "currency_code", Usage: "Currency code"},
					&cli.StringFlag{Name: "currency_symbol", Usage: "Currency symbol"},
					&cli.StringFlag{Name: "currency_format", Usage: "Currency format"},
					&cli.IntFlag{Name: "price_precision", Usage: "Decimal precision"},
					&cli.StringFlag{Name: "json", Usage: "Additional fields as JSON"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := internal.RequireFlag(cmd, "org", "ZOHO_EXPENSE_ORG_ID")
					if err != nil {
						return err
					}
					body := map[string]any{}
					if cmd.IsSet("currency_code") {
						body["currency_code"] = cmd.String("currency_code")
					}
					if cmd.IsSet("currency_symbol") {
						body["currency_symbol"] = cmd.String("currency_symbol")
					}
					if cmd.IsSet("currency_format") {
						body["currency_format"] = cmd.String("currency_format")
					}
					if cmd.IsSet("price_precision") {
						body["price_precision"] = cmd.Int("price_precision")
					}
					if err := internal.MergeJSON(cmd, body); err != nil {
						return err
					}
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
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := internal.RequireFlag(cmd, "org", "ZOHO_EXPENSE_ORG_ID")
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
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := internal.RequireFlag(cmd, "org", "ZOHO_EXPENSE_ORG_ID")
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
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := internal.RequireFlag(cmd, "org", "ZOHO_EXPENSE_ORG_ID")
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
					&cli.StringFlag{Name: "tax_name", Required: true, Usage: "Tax name"},
					&cli.FloatFlag{Name: "tax_percentage", Required: true, Usage: "Tax percentage"},
					&cli.StringFlag{Name: "tax_type", Usage: "Tax type"},
					&cli.StringFlag{Name: "json", Usage: "Additional fields as JSON"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := internal.RequireFlag(cmd, "org", "ZOHO_EXPENSE_ORG_ID")
					if err != nil {
						return err
					}
					body := map[string]any{"tax_name": cmd.String("tax_name"), "tax_percentage": cmd.Float("tax_percentage")}
					if cmd.IsSet("tax_type") {
						body["tax_type"] = cmd.String("tax_type")
					}
					if err := internal.MergeJSON(cmd, body); err != nil {
						return err
					}
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
					&cli.StringFlag{Name: "tax_name", Usage: "Tax name"},
					&cli.FloatFlag{Name: "tax_percentage", Usage: "Tax percentage"},
					&cli.StringFlag{Name: "tax_type", Usage: "Tax type"},
					&cli.StringFlag{Name: "json", Usage: "Additional fields as JSON"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := internal.RequireFlag(cmd, "org", "ZOHO_EXPENSE_ORG_ID")
					if err != nil {
						return err
					}
					body := map[string]any{}
					if cmd.IsSet("tax_name") {
						body["tax_name"] = cmd.String("tax_name")
					}
					if cmd.IsSet("tax_percentage") {
						body["tax_percentage"] = cmd.Float("tax_percentage")
					}
					if cmd.IsSet("tax_type") {
						body["tax_type"] = cmd.String("tax_type")
					}
					if err := internal.MergeJSON(cmd, body); err != nil {
						return err
					}
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
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := internal.RequireFlag(cmd, "org", "ZOHO_EXPENSE_ORG_ID")
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
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := internal.RequireFlag(cmd, "org", "ZOHO_EXPENSE_ORG_ID")
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
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := internal.RequireFlag(cmd, "org", "ZOHO_EXPENSE_ORG_ID")
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
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := internal.RequireFlag(cmd, "org", "ZOHO_EXPENSE_ORG_ID")
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
					&cli.StringFlag{Name: "tag_name", Required: true, Usage: "Tag name"},
					&cli.StringFlag{Name: "json", Usage: "Additional fields as JSON"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := internal.RequireFlag(cmd, "org", "ZOHO_EXPENSE_ORG_ID")
					if err != nil {
						return err
					}
					body := map[string]any{"tag_name": cmd.String("tag_name")}
					if err := internal.MergeJSON(cmd, body); err != nil {
						return err
					}
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
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := internal.RequireFlag(cmd, "org", "ZOHO_EXPENSE_ORG_ID")
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
					&cli.StringFlag{Name: "tag_name", Usage: "Tag name"},
					&cli.StringFlag{Name: "json", Usage: "Additional fields as JSON"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := internal.RequireFlag(cmd, "org", "ZOHO_EXPENSE_ORG_ID")
					if err != nil {
						return err
					}
					body := map[string]any{}
					if cmd.IsSet("tag_name") {
						body["tag_name"] = cmd.String("tag_name")
					}
					if err := internal.MergeJSON(cmd, body); err != nil {
						return err
					}
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
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := internal.RequireFlag(cmd, "org", "ZOHO_EXPENSE_ORG_ID")
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
					&cli.StringFlag{Name: "json", Usage: "Additional fields as JSON"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := internal.RequireFlag(cmd, "org", "ZOHO_EXPENSE_ORG_ID")
					if err != nil {
						return err
					}
					body := map[string]any{}
					if err := internal.MergeJSON(cmd, body); err != nil {
						return err
					}
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
					&cli.StringFlag{Name: "json", Usage: "Additional fields as JSON"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := internal.RequireFlag(cmd, "org", "ZOHO_EXPENSE_ORG_ID")
					if err != nil {
						return err
					}
					body := map[string]any{}
					if err := internal.MergeJSON(cmd, body); err != nil {
						return err
					}
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
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := internal.RequireFlag(cmd, "org", "ZOHO_EXPENSE_ORG_ID")
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
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := internal.RequireFlag(cmd, "org", "ZOHO_EXPENSE_ORG_ID")
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
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := internal.RequireFlag(cmd, "org", "ZOHO_EXPENSE_ORG_ID")
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
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := internal.RequireFlag(cmd, "org", "ZOHO_EXPENSE_ORG_ID")
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
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := internal.RequireFlag(cmd, "org", "ZOHO_EXPENSE_ORG_ID")
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
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := internal.RequireFlag(cmd, "org", "ZOHO_EXPENSE_ORG_ID")
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
					&cli.StringFlag{Name: "json", Usage: "Additional fields as JSON"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := internal.RequireFlag(cmd, "org", "ZOHO_EXPENSE_ORG_ID")
					if err != nil {
						return err
					}
					body := map[string]any{}
					if err := internal.MergeJSON(cmd, body); err != nil {
						return err
					}
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
