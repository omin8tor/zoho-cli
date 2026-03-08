package books

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/omin8tor/zoho-cli/internal"
	zohttp "github.com/omin8tor/zoho-cli/internal/http"
	"github.com/omin8tor/zoho-cli/internal/output"
	"github.com/omin8tor/zoho-cli/internal/pagination"
	"github.com/urfave/cli/v3"
)

func expensesCmd() *cli.Command {
	return &cli.Command{
		Name:  "expenses",
		Usage: "Expense operations",
		Commands: []*cli.Command{
			{
				Name:  "create",
				Usage: "Create an expense",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "account_id", Required: true, Usage: "Account ID"},
					&cli.StringFlag{Name: "date", Required: true, Usage: "Expense date (YYYY-MM-DD)"},
					&cli.FloatFlag{Name: "amount", Required: true, Usage: "Expense amount"},
					&cli.StringFlag{Name: "customer_id", Usage: "Customer ID"},
					&cli.StringFlag{Name: "vendor_id", Usage: "Vendor ID"},
					&cli.StringFlag{Name: "json", Usage: "Additional fields as JSON"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					body := map[string]any{}
					body["account_id"] = cmd.String("account_id")
					body["date"] = cmd.String("date")
					body["amount"] = cmd.Float("amount")
					if cmd.IsSet("customer_id") {
						body["customer_id"] = cmd.String("customer_id")
					}
					if cmd.IsSet("vendor_id") {
						body["vendor_id"] = cmd.String("vendor_id")
					}
					if err := internal.MergeJSON(cmd, body); err != nil {
						return err
					}
					raw, err := c.Request("POST", c.BooksBase+"/expenses", &zohttp.RequestOpts{
						Params: orgParams(orgID),
						JSON:   body,
					})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:  "update-by-custom-field",
				Usage: "Update an expense by custom field",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "account_id", Usage: "Account ID"},
					&cli.StringFlag{Name: "date", Usage: "Expense date (YYYY-MM-DD)"},
					&cli.FloatFlag{Name: "amount", Usage: "Expense amount"},
					&cli.StringFlag{Name: "customer_id", Usage: "Customer ID"},
					&cli.StringFlag{Name: "vendor_id", Usage: "Vendor ID"},
					&cli.StringFlag{Name: "json", Usage: "Additional fields as JSON"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					body := map[string]any{}
					if cmd.IsSet("account_id") {
						body["account_id"] = cmd.String("account_id")
					}
					if cmd.IsSet("date") {
						body["date"] = cmd.String("date")
					}
					if cmd.IsSet("amount") {
						body["amount"] = cmd.Float("amount")
					}
					if cmd.IsSet("customer_id") {
						body["customer_id"] = cmd.String("customer_id")
					}
					if cmd.IsSet("vendor_id") {
						body["vendor_id"] = cmd.String("vendor_id")
					}
					if err := internal.MergeJSON(cmd, body); err != nil {
						return err
					}
					raw, err := c.Request("PUT", c.BooksBase+"/expenses", &zohttp.RequestOpts{
						Params: orgParams(orgID),
						JSON:   body,
					})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:  "list",
				Usage: "List expenses",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "status", Usage: "Filter"},
					&cli.StringFlag{Name: "customer-id", Usage: "Filter"},
					&cli.StringFlag{Name: "vendor-id", Usage: "Filter"},
					&cli.BoolFlag{Name: "all", Usage: "Fetch all records"},
					&cli.IntFlag{Name: "limit", Usage: "Max total records to fetch"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					params := orgParams(orgID)
					if v := cmd.String("status"); v != "" {
						params["status"] = v
					}
					if v := cmd.String("customer-id"); v != "" {
						params["customer_id"] = v
					}
					if v := cmd.String("vendor-id"); v != "" {
						params["vendor_id"] = v
					}

					if cmd.Bool("all") || cmd.IsSet("limit") {

						items, err := pagination.Paginate(pagination.PaginationConfig{

							Client: c,

							URL: c.BooksBase + "/expenses",

							Opts: &zohttp.RequestOpts{Params: params},

							ItemsKey: "expenses",

							PageSize: 200,

							Limit: int(cmd.Int("limit")),

							SetPage: pagination.PagePerPage(200),

							HasMore: pagination.HasMoreBooks,
						})

						if err != nil {

							return err

						}

						return output.JSON(items)

					}
					raw, err := c.Request("GET", c.BooksBase+"/expenses", &zohttp.RequestOpts{Params: params})
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
					&cli.StringFlag{Name: "account_id", Usage: "Account ID"},
					&cli.StringFlag{Name: "date", Usage: "Expense date (YYYY-MM-DD)"},
					&cli.FloatFlag{Name: "amount", Usage: "Expense amount"},
					&cli.StringFlag{Name: "customer_id", Usage: "Customer ID"},
					&cli.StringFlag{Name: "vendor_id", Usage: "Vendor ID"},
					&cli.StringFlag{Name: "json", Usage: "Additional fields as JSON"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					body := map[string]any{}
					if cmd.IsSet("account_id") {
						body["account_id"] = cmd.String("account_id")
					}
					if cmd.IsSet("date") {
						body["date"] = cmd.String("date")
					}
					if cmd.IsSet("amount") {
						body["amount"] = cmd.Float("amount")
					}
					if cmd.IsSet("customer_id") {
						body["customer_id"] = cmd.String("customer_id")
					}
					if cmd.IsSet("vendor_id") {
						body["vendor_id"] = cmd.String("vendor_id")
					}
					if err := internal.MergeJSON(cmd, body); err != nil {
						return err
					}
					raw, err := c.Request("PUT", c.BooksBase+"/expenses/"+cmd.Args().First(), &zohttp.RequestOpts{
						Params: orgParams(orgID),
						JSON:   body,
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
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("GET", c.BooksBase+"/expenses/"+cmd.Args().First(), &zohttp.RequestOpts{Params: orgParams(orgID)})
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
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("DELETE", c.BooksBase+"/expenses/"+cmd.Args().First(), &zohttp.RequestOpts{Params: orgParams(orgID)})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "list-history",
				Usage:     "List history of an expense",
				ArgsUsage: "<expense-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("GET", c.BooksBase+"/expenses/"+cmd.Args().First()+"/comments", &zohttp.RequestOpts{Params: orgParams(orgID)})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:  "create-employee",
				Usage: "Create an employee",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "employee_name", Required: true, Usage: "Employee name"},
					&cli.StringFlag{Name: "email", Usage: "Employee email"},
					&cli.StringFlag{Name: "json", Usage: "Additional fields as JSON"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					body := map[string]any{}
					body["employee_name"] = cmd.String("employee_name")
					if cmd.IsSet("email") {
						body["email"] = cmd.String("email")
					}
					if err := internal.MergeJSON(cmd, body); err != nil {
						return err
					}
					raw, err := c.Request("POST", c.BooksBase+"/expenses"+"/employees", &zohttp.RequestOpts{
						Params: orgParams(orgID),
						JSON:   body,
					})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:  "list-employees",
				Usage: "List employees",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("GET", c.BooksBase+"/expenses"+"/employees", &zohttp.RequestOpts{Params: orgParams(orgID)})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "get-employee",
				Usage:     "Get an employee",
				ArgsUsage: "<employee-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("GET", c.BooksBase+"/expenses"+"/employees/"+cmd.Args().First(), &zohttp.RequestOpts{Params: orgParams(orgID)})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "delete-employee",
				Usage:     "Delete an employee",
				ArgsUsage: "<employee-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("DELETE", c.BooksBase+"/expenses"+"/employees/"+cmd.Args().First(), &zohttp.RequestOpts{Params: orgParams(orgID)})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "add-receipt",
				Usage:     "Add receipt to an expense",
				ArgsUsage: "<expense-id>",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "file", Required: true, Usage: "File path to upload"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					data, err := os.ReadFile(cmd.String("file"))
					if err != nil {
						return fmt.Errorf("reading file: %w", err)
					}
					raw, err := c.Request("POST", c.BooksBase+"/expenses/"+cmd.Args().First()+"/receipt", &zohttp.RequestOpts{
						Params: orgParams(orgID),
						Files:  map[string]zohttp.FileUpload{"receipt": {Filename: filepath.Base(cmd.String("file")), Data: data}},
					})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "get-receipt",
				Usage:     "Get receipt of an expense",
				ArgsUsage: "<expense-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("GET", c.BooksBase+"/expenses/"+cmd.Args().First()+"/receipt", &zohttp.RequestOpts{Params: orgParams(orgID)})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "delete-receipt",
				Usage:     "Delete receipt of an expense",
				ArgsUsage: "<expense-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("DELETE", c.BooksBase+"/expenses/"+cmd.Args().First()+"/receipt", &zohttp.RequestOpts{Params: orgParams(orgID)})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "add-attachment",
				Usage:     "Add attachment to an expense",
				ArgsUsage: "<expense-id>",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "file", Required: true, Usage: "File path to upload"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					data, err := os.ReadFile(cmd.String("file"))
					if err != nil {
						return fmt.Errorf("reading file: %w", err)
					}
					raw, err := c.Request("POST", c.BooksBase+"/expenses/"+cmd.Args().First()+"/attachment", &zohttp.RequestOpts{
						Params: orgParams(orgID),
						Files:  map[string]zohttp.FileUpload{"attachment": {Filename: filepath.Base(cmd.String("file")), Data: data}},
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

func recurringExpensesCmd() *cli.Command {
	return &cli.Command{
		Name:  "recurring-expenses",
		Usage: "Recurring expense operations",
		Commands: []*cli.Command{
			{
				Name:  "create",
				Usage: "Create a recurring expense",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "account_id", Required: true, Usage: "Account ID"},
					&cli.FloatFlag{Name: "amount", Required: true, Usage: "Amount"},
					&cli.StringFlag{Name: "start_date", Usage: "Start date (YYYY-MM-DD)"},
					&cli.IntFlag{Name: "repeat_every", Usage: "Repeat interval"},
					&cli.StringFlag{Name: "json", Usage: "Additional fields as JSON"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					body := map[string]any{}
					body["account_id"] = cmd.String("account_id")
					body["amount"] = cmd.Float("amount")
					if cmd.IsSet("start_date") {
						body["start_date"] = cmd.String("start_date")
					}
					if cmd.IsSet("repeat_every") {
						body["repeat_every"] = cmd.Int("repeat_every")
					}
					if err := internal.MergeJSON(cmd, body); err != nil {
						return err
					}
					raw, err := c.Request("POST", c.BooksBase+"/recurringexpenses", &zohttp.RequestOpts{
						Params: orgParams(orgID),
						JSON:   body,
					})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:  "update-by-custom-field",
				Usage: "Update a recurring expense by custom field",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "account_id", Usage: "Account ID"},
					&cli.FloatFlag{Name: "amount", Usage: "Amount"},
					&cli.StringFlag{Name: "start_date", Usage: "Start date (YYYY-MM-DD)"},
					&cli.IntFlag{Name: "repeat_every", Usage: "Repeat interval"},
					&cli.StringFlag{Name: "json", Usage: "Additional fields as JSON"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					body := map[string]any{}
					if cmd.IsSet("account_id") {
						body["account_id"] = cmd.String("account_id")
					}
					if cmd.IsSet("amount") {
						body["amount"] = cmd.Float("amount")
					}
					if cmd.IsSet("start_date") {
						body["start_date"] = cmd.String("start_date")
					}
					if cmd.IsSet("repeat_every") {
						body["repeat_every"] = cmd.Int("repeat_every")
					}
					if err := internal.MergeJSON(cmd, body); err != nil {
						return err
					}
					raw, err := c.Request("PUT", c.BooksBase+"/recurringexpenses", &zohttp.RequestOpts{
						Params: orgParams(orgID),
						JSON:   body,
					})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:  "list",
				Usage: "List recurring expenses",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "status", Usage: "Filter"},
					&cli.BoolFlag{Name: "all", Usage: "Fetch all records"},
					&cli.IntFlag{Name: "limit", Usage: "Max total records to fetch"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					params := orgParams(orgID)
					if v := cmd.String("status"); v != "" {
						params["status"] = v
					}

					if cmd.Bool("all") || cmd.IsSet("limit") {

						items, err := pagination.Paginate(pagination.PaginationConfig{

							Client: c,

							URL: c.BooksBase + "/recurringexpenses",

							Opts: &zohttp.RequestOpts{Params: params},

							ItemsKey: "recurring_expenses",

							PageSize: 200,

							Limit: int(cmd.Int("limit")),

							SetPage: pagination.PagePerPage(200),

							HasMore: pagination.HasMoreBooks,
						})

						if err != nil {

							return err

						}

						return output.JSON(items)

					}
					raw, err := c.Request("GET", c.BooksBase+"/recurringexpenses", &zohttp.RequestOpts{Params: params})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "update",
				Usage:     "Update a recurring expense",
				ArgsUsage: "<recurringexpense-id>",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "account_id", Usage: "Account ID"},
					&cli.FloatFlag{Name: "amount", Usage: "Amount"},
					&cli.StringFlag{Name: "start_date", Usage: "Start date (YYYY-MM-DD)"},
					&cli.IntFlag{Name: "repeat_every", Usage: "Repeat interval"},
					&cli.StringFlag{Name: "json", Usage: "Additional fields as JSON"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					body := map[string]any{}
					if cmd.IsSet("account_id") {
						body["account_id"] = cmd.String("account_id")
					}
					if cmd.IsSet("amount") {
						body["amount"] = cmd.Float("amount")
					}
					if cmd.IsSet("start_date") {
						body["start_date"] = cmd.String("start_date")
					}
					if cmd.IsSet("repeat_every") {
						body["repeat_every"] = cmd.Int("repeat_every")
					}
					if err := internal.MergeJSON(cmd, body); err != nil {
						return err
					}
					raw, err := c.Request("PUT", c.BooksBase+"/recurringexpenses/"+cmd.Args().First(), &zohttp.RequestOpts{
						Params: orgParams(orgID),
						JSON:   body,
					})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "get",
				Usage:     "Get a recurring expense",
				ArgsUsage: "<recurringexpense-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("GET", c.BooksBase+"/recurringexpenses/"+cmd.Args().First(), &zohttp.RequestOpts{Params: orgParams(orgID)})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "delete",
				Usage:     "Delete a recurring expense",
				ArgsUsage: "<recurringexpense-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("DELETE", c.BooksBase+"/recurringexpenses/"+cmd.Args().First(), &zohttp.RequestOpts{Params: orgParams(orgID)})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "stop",
				Usage:     "Stop a recurring expense",
				ArgsUsage: "<recurringexpense-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("POST", c.BooksBase+"/recurringexpenses/"+cmd.Args().First()+"/status/stop", &zohttp.RequestOpts{Params: orgParams(orgID)})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "resume",
				Usage:     "Resume a recurring expense",
				ArgsUsage: "<recurringexpense-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("POST", c.BooksBase+"/recurringexpenses/"+cmd.Args().First()+"/status/resume", &zohttp.RequestOpts{Params: orgParams(orgID)})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "list-child-expenses",
				Usage:     "List child expenses of a recurring expense",
				ArgsUsage: "<recurringexpense-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("GET", c.BooksBase+"/recurringexpenses/"+cmd.Args().First()+"/expenses", &zohttp.RequestOpts{Params: orgParams(orgID)})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "list-history",
				Usage:     "List history of a recurring expense",
				ArgsUsage: "<recurringexpense-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("GET", c.BooksBase+"/recurringexpenses/"+cmd.Args().First()+"/comments", &zohttp.RequestOpts{Params: orgParams(orgID)})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
		},
	}
}
