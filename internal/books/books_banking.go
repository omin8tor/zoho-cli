package books

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/omin8tor/zoho-cli/internal"
	zohttp "github.com/omin8tor/zoho-cli/internal/http"
	"github.com/omin8tor/zoho-cli/internal/output"
	"github.com/omin8tor/zoho-cli/internal/pagination"
	"github.com/urfave/cli/v3"
)

func bankAccountsCmd() *cli.Command {
	return &cli.Command{
		Name:  "bank-accounts",
		Usage: "Bank account operations",
		Commands: []*cli.Command{
			{
				Name:  "create",
				Usage: "Create a bank account",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "account_name", Required: true, Usage: "Account name"},
					&cli.StringFlag{Name: "account_type", Required: true, Usage: "Account type"},
					&cli.StringFlag{Name: "account_number", Usage: "Account number"},
					&cli.StringFlag{Name: "currency_id", Usage: "Currency ID"},
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
					body["account_name"] = cmd.String("account_name")
					body["account_type"] = cmd.String("account_type")
					if cmd.IsSet("account_number") {
						body["account_number"] = cmd.String("account_number")
					}
					if cmd.IsSet("currency_id") {
						body["currency_id"] = cmd.String("currency_id")
					}
					if err := internal.MergeJSON(cmd, body); err != nil {
						return err
					}
					raw, err := c.Request("POST", c.BooksBase+"/bankaccounts", &zohttp.RequestOpts{
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
				Usage: "List bank accounts",
				Flags: []cli.Flag{
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

					if cmd.Bool("all") || cmd.IsSet("limit") {

						items, err := pagination.Paginate(pagination.PaginationConfig{

							Client: c,

							URL: c.BooksBase + "/bankaccounts",

							Opts: &zohttp.RequestOpts{Params: params},

							ItemsKey: "bankaccounts",

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
					raw, err := c.Request("GET", c.BooksBase+"/bankaccounts", &zohttp.RequestOpts{Params: params})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "update",
				Usage:     "Update a bank account",
				ArgsUsage: "<account-id>",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "account_name", Usage: "Account name"},
					&cli.StringFlag{Name: "account_type", Usage: "Account type"},
					&cli.StringFlag{Name: "account_number", Usage: "Account number"},
					&cli.StringFlag{Name: "currency_id", Usage: "Currency ID"},
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
					if cmd.IsSet("account_name") {
						body["account_name"] = cmd.String("account_name")
					}
					if cmd.IsSet("account_type") {
						body["account_type"] = cmd.String("account_type")
					}
					if cmd.IsSet("account_number") {
						body["account_number"] = cmd.String("account_number")
					}
					if cmd.IsSet("currency_id") {
						body["currency_id"] = cmd.String("currency_id")
					}
					if err := internal.MergeJSON(cmd, body); err != nil {
						return err
					}
					raw, err := c.Request("PUT", c.BooksBase+"/bankaccounts/"+cmd.Args().First(), &zohttp.RequestOpts{
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
				Usage:     "Get a bank account",
				ArgsUsage: "<account-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("GET", c.BooksBase+"/bankaccounts/"+cmd.Args().First(), &zohttp.RequestOpts{Params: orgParams(orgID)})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "delete",
				Usage:     "Delete a bank account",
				ArgsUsage: "<account-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("DELETE", c.BooksBase+"/bankaccounts/"+cmd.Args().First(), &zohttp.RequestOpts{Params: orgParams(orgID)})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "deactivate",
				Usage:     "Deactivate a bank account",
				ArgsUsage: "<account-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("POST", c.BooksBase+"/bankaccounts/"+cmd.Args().First()+"/inactive", &zohttp.RequestOpts{Params: orgParams(orgID)})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "activate",
				Usage:     "Activate a bank account",
				ArgsUsage: "<account-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("POST", c.BooksBase+"/bankaccounts/"+cmd.Args().First()+"/active", &zohttp.RequestOpts{Params: orgParams(orgID)})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "import-statement",
				Usage:     "Import a bank statement",
				ArgsUsage: "<account-id>",
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
					raw, err := c.Request("POST", c.BooksBase+"/bankaccounts/"+cmd.Args().First()+"/statement", &zohttp.RequestOpts{
						Params: orgParams(orgID),
						Files:  map[string]zohttp.FileUpload{"statement": {Filename: filepath.Base(cmd.String("file")), Data: data}},
					})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "get-last-statement",
				Usage:     "Get last imported statement",
				ArgsUsage: "<account-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("GET", c.BooksBase+"/bankaccounts/"+cmd.Args().First()+"/statement", &zohttp.RequestOpts{Params: orgParams(orgID)})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "delete-last-statement",
				Usage:     "Delete last imported statement",
				ArgsUsage: "<account-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("DELETE", c.BooksBase+"/bankaccounts/"+cmd.Args().First()+"/statement", &zohttp.RequestOpts{Params: orgParams(orgID)})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
		},
	}
}

func bankTransactionsCmd() *cli.Command {
	return &cli.Command{
		Name:  "bank-transactions",
		Usage: "Bank transaction operations",
		Commands: []*cli.Command{
			{
				Name:  "create",
				Usage: "Create a bank transaction",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "transaction_type", Required: true, Usage: "Transaction type"},
					&cli.FloatFlag{Name: "amount", Required: true, Usage: "Amount"},
					&cli.StringFlag{Name: "date", Required: true, Usage: "Date (YYYY-MM-DD)"},
					&cli.StringFlag{Name: "from_account_id", Usage: "From account ID"},
					&cli.StringFlag{Name: "to_account_id", Usage: "To account ID"},
					&cli.StringFlag{Name: "payment_mode", Usage: "Payment mode"},
					&cli.FloatFlag{Name: "exchange_rate", Usage: "Exchange rate"},
					&cli.StringFlag{Name: "customer_id", Usage: "Customer ID"},
					&cli.StringFlag{Name: "vendor_id", Usage: "Vendor ID"},
					&cli.StringFlag{Name: "reference_number", Usage: "Reference number"},
					&cli.StringFlag{Name: "description", Usage: "Description"},
					&cli.StringFlag{Name: "currency_id", Usage: "Currency ID"},
					&cli.StringFlag{Name: "tax_id", Usage: "Tax ID"},
					&cli.StringFlag{Name: "paid_through_account_id", Usage: "Paid through account ID"},
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
					body["transaction_type"] = cmd.String("transaction_type")
					body["amount"] = cmd.Float("amount")
					body["date"] = cmd.String("date")
					if cmd.IsSet("from_account_id") {
						body["from_account_id"] = cmd.String("from_account_id")
					}
					if cmd.IsSet("to_account_id") {
						body["to_account_id"] = cmd.String("to_account_id")
					}
					if cmd.IsSet("payment_mode") {
						body["payment_mode"] = cmd.String("payment_mode")
					}
					if cmd.IsSet("exchange_rate") {
						body["exchange_rate"] = cmd.Float("exchange_rate")
					}
					if cmd.IsSet("customer_id") {
						body["customer_id"] = cmd.String("customer_id")
					}
					if cmd.IsSet("vendor_id") {
						body["vendor_id"] = cmd.String("vendor_id")
					}
					if cmd.IsSet("reference_number") {
						body["reference_number"] = cmd.String("reference_number")
					}
					if cmd.IsSet("description") {
						body["description"] = cmd.String("description")
					}
					if cmd.IsSet("currency_id") {
						body["currency_id"] = cmd.String("currency_id")
					}
					if cmd.IsSet("tax_id") {
						body["tax_id"] = cmd.String("tax_id")
					}
					if cmd.IsSet("paid_through_account_id") {
						body["paid_through_account_id"] = cmd.String("paid_through_account_id")
					}
					if err := internal.MergeJSON(cmd, body); err != nil {
						return err
					}
					raw, err := c.Request("POST", c.BooksBase+"/banktransactions", &zohttp.RequestOpts{
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
				Usage: "List bank transactions",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "account-id", Usage: "Filter"},
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
					if v := cmd.String("account-id"); v != "" {
						params["account_id"] = v
					}
					if v := cmd.String("status"); v != "" {
						params["status"] = v
					}

					if cmd.Bool("all") || cmd.IsSet("limit") {

						items, err := pagination.Paginate(pagination.PaginationConfig{

							Client: c,

							URL: c.BooksBase + "/banktransactions",

							Opts: &zohttp.RequestOpts{Params: params},

							ItemsKey: "banktransactions",

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
					raw, err := c.Request("GET", c.BooksBase+"/banktransactions", &zohttp.RequestOpts{Params: params})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "update",
				Usage:     "Update a bank transaction",
				ArgsUsage: "<transaction-id>",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "transaction_type", Usage: "Transaction type"},
					&cli.FloatFlag{Name: "amount", Usage: "Amount"},
					&cli.StringFlag{Name: "date", Usage: "Date (YYYY-MM-DD)"},
					&cli.StringFlag{Name: "from_account_id", Usage: "From account ID"},
					&cli.StringFlag{Name: "to_account_id", Usage: "To account ID"},
					&cli.StringFlag{Name: "payment_mode", Usage: "Payment mode"},
					&cli.FloatFlag{Name: "exchange_rate", Usage: "Exchange rate"},
					&cli.StringFlag{Name: "customer_id", Usage: "Customer ID"},
					&cli.StringFlag{Name: "vendor_id", Usage: "Vendor ID"},
					&cli.StringFlag{Name: "reference_number", Usage: "Reference number"},
					&cli.StringFlag{Name: "description", Usage: "Description"},
					&cli.StringFlag{Name: "currency_id", Usage: "Currency ID"},
					&cli.StringFlag{Name: "tax_id", Usage: "Tax ID"},
					&cli.StringFlag{Name: "paid_through_account_id", Usage: "Paid through account ID"},
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
					if cmd.IsSet("transaction_type") {
						body["transaction_type"] = cmd.String("transaction_type")
					}
					if cmd.IsSet("amount") {
						body["amount"] = cmd.Float("amount")
					}
					if cmd.IsSet("date") {
						body["date"] = cmd.String("date")
					}
					if cmd.IsSet("from_account_id") {
						body["from_account_id"] = cmd.String("from_account_id")
					}
					if cmd.IsSet("to_account_id") {
						body["to_account_id"] = cmd.String("to_account_id")
					}
					if cmd.IsSet("payment_mode") {
						body["payment_mode"] = cmd.String("payment_mode")
					}
					if cmd.IsSet("exchange_rate") {
						body["exchange_rate"] = cmd.Float("exchange_rate")
					}
					if cmd.IsSet("customer_id") {
						body["customer_id"] = cmd.String("customer_id")
					}
					if cmd.IsSet("vendor_id") {
						body["vendor_id"] = cmd.String("vendor_id")
					}
					if cmd.IsSet("reference_number") {
						body["reference_number"] = cmd.String("reference_number")
					}
					if cmd.IsSet("description") {
						body["description"] = cmd.String("description")
					}
					if cmd.IsSet("currency_id") {
						body["currency_id"] = cmd.String("currency_id")
					}
					if cmd.IsSet("tax_id") {
						body["tax_id"] = cmd.String("tax_id")
					}
					if cmd.IsSet("paid_through_account_id") {
						body["paid_through_account_id"] = cmd.String("paid_through_account_id")
					}
					if err := internal.MergeJSON(cmd, body); err != nil {
						return err
					}
					raw, err := c.Request("PUT", c.BooksBase+"/banktransactions/"+cmd.Args().First(), &zohttp.RequestOpts{
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
				Usage:     "Get a bank transaction",
				ArgsUsage: "<transaction-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("GET", c.BooksBase+"/banktransactions/"+cmd.Args().First(), &zohttp.RequestOpts{Params: orgParams(orgID)})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "delete",
				Usage:     "Delete a bank transaction",
				ArgsUsage: "<transaction-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("DELETE", c.BooksBase+"/banktransactions/"+cmd.Args().First(), &zohttp.RequestOpts{Params: orgParams(orgID)})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "match",
				Usage:     "Match a bank transaction",
				ArgsUsage: "<transaction-id>",
				Flags: []cli.Flag{
					&cli.FloatFlag{Name: "amount", Usage: "Amount"},
					&cli.StringFlag{Name: "date", Usage: "Date (YYYY-MM-DD)"},
					&cli.StringFlag{Name: "description", Usage: "Description"},
					&cli.StringFlag{Name: "reference_number", Usage: "Reference number"},
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
					if cmd.IsSet("amount") {
						body["amount"] = cmd.Float("amount")
					}
					if cmd.IsSet("date") {
						body["date"] = cmd.String("date")
					}
					if cmd.IsSet("description") {
						body["description"] = cmd.String("description")
					}
					if cmd.IsSet("reference_number") {
						body["reference_number"] = cmd.String("reference_number")
					}
					if err := internal.MergeJSON(cmd, body); err != nil {
						return err
					}
					raw, err := c.Request("POST", c.BooksBase+"/banktransactions/"+cmd.Args().First()+"/match", &zohttp.RequestOpts{
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
				Name:      "get-matching",
				Usage:     "Get matching transactions",
				ArgsUsage: "<transaction-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("GET", c.BooksBase+"/banktransactions/"+cmd.Args().First()+"/match", &zohttp.RequestOpts{Params: orgParams(orgID)})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "unmatch",
				Usage:     "Unmatch a bank transaction",
				ArgsUsage: "<transaction-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("POST", c.BooksBase+"/banktransactions/"+cmd.Args().First()+"/unmatch", &zohttp.RequestOpts{Params: orgParams(orgID)})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "exclude",
				Usage:     "Exclude a bank transaction",
				ArgsUsage: "<transaction-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("POST", c.BooksBase+"/banktransactions/"+cmd.Args().First()+"/exclude", &zohttp.RequestOpts{Params: orgParams(orgID)})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "restore",
				Usage:     "Restore a bank transaction",
				ArgsUsage: "<transaction-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("POST", c.BooksBase+"/banktransactions/"+cmd.Args().First()+"/restore", &zohttp.RequestOpts{Params: orgParams(orgID)})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "categorize",
				Usage:     "Categorize a bank transaction",
				ArgsUsage: "<transaction-id>",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "transaction_type", Usage: "Transaction type"},
					&cli.StringFlag{Name: "account_id", Usage: "Account ID"},
					&cli.FloatFlag{Name: "amount", Usage: "Amount"},
					&cli.StringFlag{Name: "date", Usage: "Date (YYYY-MM-DD)"},
					&cli.StringFlag{Name: "payment_mode", Usage: "Payment mode"},
					&cli.StringFlag{Name: "reference_number", Usage: "Reference number"},
					&cli.StringFlag{Name: "description", Usage: "Description"},
					&cli.StringFlag{Name: "customer_id", Usage: "Customer ID"},
					&cli.StringFlag{Name: "vendor_id", Usage: "Vendor ID"},
					&cli.StringFlag{Name: "tax_id", Usage: "Tax ID"},
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
					if cmd.IsSet("transaction_type") {
						body["transaction_type"] = cmd.String("transaction_type")
					}
					if cmd.IsSet("account_id") {
						body["account_id"] = cmd.String("account_id")
					}
					if cmd.IsSet("amount") {
						body["amount"] = cmd.Float("amount")
					}
					if cmd.IsSet("date") {
						body["date"] = cmd.String("date")
					}
					if cmd.IsSet("payment_mode") {
						body["payment_mode"] = cmd.String("payment_mode")
					}
					if cmd.IsSet("reference_number") {
						body["reference_number"] = cmd.String("reference_number")
					}
					if cmd.IsSet("description") {
						body["description"] = cmd.String("description")
					}
					if cmd.IsSet("customer_id") {
						body["customer_id"] = cmd.String("customer_id")
					}
					if cmd.IsSet("vendor_id") {
						body["vendor_id"] = cmd.String("vendor_id")
					}
					if cmd.IsSet("tax_id") {
						body["tax_id"] = cmd.String("tax_id")
					}
					if err := internal.MergeJSON(cmd, body); err != nil {
						return err
					}
					raw, err := c.Request("POST", c.BooksBase+"/banktransactions/"+cmd.Args().First()+"/categorize", &zohttp.RequestOpts{
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
				Name:      "categorize-as-expense",
				Usage:     "Categorize as expense",
				ArgsUsage: "<transaction-id>",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "account_id", Usage: "Account ID"},
					&cli.FloatFlag{Name: "amount", Usage: "Amount"},
					&cli.StringFlag{Name: "date", Usage: "Date (YYYY-MM-DD)"},
					&cli.StringFlag{Name: "payment_mode", Usage: "Payment mode"},
					&cli.StringFlag{Name: "reference_number", Usage: "Reference number"},
					&cli.StringFlag{Name: "description", Usage: "Description"},
					&cli.StringFlag{Name: "vendor_id", Usage: "Vendor ID"},
					&cli.StringFlag{Name: "tax_id", Usage: "Tax ID"},
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
					if cmd.IsSet("date") {
						body["date"] = cmd.String("date")
					}
					if cmd.IsSet("payment_mode") {
						body["payment_mode"] = cmd.String("payment_mode")
					}
					if cmd.IsSet("reference_number") {
						body["reference_number"] = cmd.String("reference_number")
					}
					if cmd.IsSet("description") {
						body["description"] = cmd.String("description")
					}
					if cmd.IsSet("vendor_id") {
						body["vendor_id"] = cmd.String("vendor_id")
					}
					if cmd.IsSet("tax_id") {
						body["tax_id"] = cmd.String("tax_id")
					}
					if err := internal.MergeJSON(cmd, body); err != nil {
						return err
					}
					raw, err := c.Request("POST", c.BooksBase+"/banktransactions/"+cmd.Args().First()+"/categorize/expense", &zohttp.RequestOpts{
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
				Name:      "uncategorize",
				Usage:     "Uncategorize a bank transaction",
				ArgsUsage: "<transaction-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("POST", c.BooksBase+"/banktransactions/"+cmd.Args().First()+"/uncategorize", &zohttp.RequestOpts{Params: orgParams(orgID)})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "categorize-vendor-payment",
				Usage:     "Categorize as vendor payment",
				ArgsUsage: "<transaction-id>",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "vendor_id", Usage: "Vendor ID"},
					&cli.FloatFlag{Name: "amount", Usage: "Amount"},
					&cli.StringFlag{Name: "date", Usage: "Date (YYYY-MM-DD)"},
					&cli.StringFlag{Name: "payment_mode", Usage: "Payment mode"},
					&cli.StringFlag{Name: "reference_number", Usage: "Reference number"},
					&cli.StringFlag{Name: "description", Usage: "Description"},
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
					if cmd.IsSet("vendor_id") {
						body["vendor_id"] = cmd.String("vendor_id")
					}
					if cmd.IsSet("amount") {
						body["amount"] = cmd.Float("amount")
					}
					if cmd.IsSet("date") {
						body["date"] = cmd.String("date")
					}
					if cmd.IsSet("payment_mode") {
						body["payment_mode"] = cmd.String("payment_mode")
					}
					if cmd.IsSet("reference_number") {
						body["reference_number"] = cmd.String("reference_number")
					}
					if cmd.IsSet("description") {
						body["description"] = cmd.String("description")
					}
					if err := internal.MergeJSON(cmd, body); err != nil {
						return err
					}
					raw, err := c.Request("POST", c.BooksBase+"/banktransactions/"+cmd.Args().First()+"/categorize/vendorpayment", &zohttp.RequestOpts{
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
				Name:      "categorize-customer-payment",
				Usage:     "Categorize as customer payment",
				ArgsUsage: "<transaction-id>",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "customer_id", Usage: "Customer ID"},
					&cli.FloatFlag{Name: "amount", Usage: "Amount"},
					&cli.StringFlag{Name: "date", Usage: "Date (YYYY-MM-DD)"},
					&cli.StringFlag{Name: "payment_mode", Usage: "Payment mode"},
					&cli.StringFlag{Name: "reference_number", Usage: "Reference number"},
					&cli.StringFlag{Name: "description", Usage: "Description"},
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
					if cmd.IsSet("customer_id") {
						body["customer_id"] = cmd.String("customer_id")
					}
					if cmd.IsSet("amount") {
						body["amount"] = cmd.Float("amount")
					}
					if cmd.IsSet("date") {
						body["date"] = cmd.String("date")
					}
					if cmd.IsSet("payment_mode") {
						body["payment_mode"] = cmd.String("payment_mode")
					}
					if cmd.IsSet("reference_number") {
						body["reference_number"] = cmd.String("reference_number")
					}
					if cmd.IsSet("description") {
						body["description"] = cmd.String("description")
					}
					if err := internal.MergeJSON(cmd, body); err != nil {
						return err
					}
					raw, err := c.Request("POST", c.BooksBase+"/banktransactions/"+cmd.Args().First()+"/categorize/customerpayment", &zohttp.RequestOpts{
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
				Name:      "categorize-credit-note-refund",
				Usage:     "Categorize as credit note refund",
				ArgsUsage: "<transaction-id>",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "customer_id", Usage: "Customer ID"},
					&cli.StringFlag{Name: "date", Usage: "Date (YYYY-MM-DD)"},
					&cli.StringFlag{Name: "refund_mode", Usage: "Refund mode"},
					&cli.FloatFlag{Name: "amount", Usage: "Amount"},
					&cli.StringFlag{Name: "from_account_id", Usage: "From account ID"},
					&cli.StringFlag{Name: "reference_number", Usage: "Reference number"},
					&cli.StringFlag{Name: "description", Usage: "Description"},
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
					if cmd.IsSet("customer_id") {
						body["customer_id"] = cmd.String("customer_id")
					}
					if cmd.IsSet("date") {
						body["date"] = cmd.String("date")
					}
					if cmd.IsSet("refund_mode") {
						body["refund_mode"] = cmd.String("refund_mode")
					}
					if cmd.IsSet("amount") {
						body["amount"] = cmd.Float("amount")
					}
					if cmd.IsSet("from_account_id") {
						body["from_account_id"] = cmd.String("from_account_id")
					}
					if cmd.IsSet("reference_number") {
						body["reference_number"] = cmd.String("reference_number")
					}
					if cmd.IsSet("description") {
						body["description"] = cmd.String("description")
					}
					if err := internal.MergeJSON(cmd, body); err != nil {
						return err
					}
					raw, err := c.Request("POST", c.BooksBase+"/banktransactions/"+cmd.Args().First()+"/categorize/creditnoterefund", &zohttp.RequestOpts{
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
				Name:      "categorize-vendor-credit-refund",
				Usage:     "Categorize as vendor credit refund",
				ArgsUsage: "<transaction-id>",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "vendor_id", Usage: "Vendor ID"},
					&cli.StringFlag{Name: "date", Usage: "Date (YYYY-MM-DD)"},
					&cli.StringFlag{Name: "refund_mode", Usage: "Refund mode"},
					&cli.FloatFlag{Name: "amount", Usage: "Amount"},
					&cli.StringFlag{Name: "from_account_id", Usage: "From account ID"},
					&cli.StringFlag{Name: "reference_number", Usage: "Reference number"},
					&cli.StringFlag{Name: "description", Usage: "Description"},
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
					if cmd.IsSet("vendor_id") {
						body["vendor_id"] = cmd.String("vendor_id")
					}
					if cmd.IsSet("date") {
						body["date"] = cmd.String("date")
					}
					if cmd.IsSet("refund_mode") {
						body["refund_mode"] = cmd.String("refund_mode")
					}
					if cmd.IsSet("amount") {
						body["amount"] = cmd.Float("amount")
					}
					if cmd.IsSet("from_account_id") {
						body["from_account_id"] = cmd.String("from_account_id")
					}
					if cmd.IsSet("reference_number") {
						body["reference_number"] = cmd.String("reference_number")
					}
					if cmd.IsSet("description") {
						body["description"] = cmd.String("description")
					}
					if err := internal.MergeJSON(cmd, body); err != nil {
						return err
					}
					raw, err := c.Request("POST", c.BooksBase+"/banktransactions/"+cmd.Args().First()+"/categorize/vendorcreditrefund", &zohttp.RequestOpts{
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
				Name:      "categorize-customer-payment-refund",
				Usage:     "Categorize as customer payment refund",
				ArgsUsage: "<transaction-id>",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "customer_id", Usage: "Customer ID"},
					&cli.StringFlag{Name: "date", Usage: "Date (YYYY-MM-DD)"},
					&cli.StringFlag{Name: "refund_mode", Usage: "Refund mode"},
					&cli.FloatFlag{Name: "amount", Usage: "Amount"},
					&cli.StringFlag{Name: "from_account_id", Usage: "From account ID"},
					&cli.StringFlag{Name: "reference_number", Usage: "Reference number"},
					&cli.StringFlag{Name: "description", Usage: "Description"},
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
					if cmd.IsSet("customer_id") {
						body["customer_id"] = cmd.String("customer_id")
					}
					if cmd.IsSet("date") {
						body["date"] = cmd.String("date")
					}
					if cmd.IsSet("refund_mode") {
						body["refund_mode"] = cmd.String("refund_mode")
					}
					if cmd.IsSet("amount") {
						body["amount"] = cmd.Float("amount")
					}
					if cmd.IsSet("from_account_id") {
						body["from_account_id"] = cmd.String("from_account_id")
					}
					if cmd.IsSet("reference_number") {
						body["reference_number"] = cmd.String("reference_number")
					}
					if cmd.IsSet("description") {
						body["description"] = cmd.String("description")
					}
					if err := internal.MergeJSON(cmd, body); err != nil {
						return err
					}
					raw, err := c.Request("POST", c.BooksBase+"/banktransactions/"+cmd.Args().First()+"/categorize/customerpaymentrefund", &zohttp.RequestOpts{
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
				Name:      "categorize-vendor-payment-refund",
				Usage:     "Categorize as vendor payment refund",
				ArgsUsage: "<transaction-id>",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "vendor_id", Usage: "Vendor ID"},
					&cli.StringFlag{Name: "date", Usage: "Date (YYYY-MM-DD)"},
					&cli.StringFlag{Name: "refund_mode", Usage: "Refund mode"},
					&cli.FloatFlag{Name: "amount", Usage: "Amount"},
					&cli.StringFlag{Name: "from_account_id", Usage: "From account ID"},
					&cli.StringFlag{Name: "reference_number", Usage: "Reference number"},
					&cli.StringFlag{Name: "description", Usage: "Description"},
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
					if cmd.IsSet("vendor_id") {
						body["vendor_id"] = cmd.String("vendor_id")
					}
					if cmd.IsSet("date") {
						body["date"] = cmd.String("date")
					}
					if cmd.IsSet("refund_mode") {
						body["refund_mode"] = cmd.String("refund_mode")
					}
					if cmd.IsSet("amount") {
						body["amount"] = cmd.Float("amount")
					}
					if cmd.IsSet("from_account_id") {
						body["from_account_id"] = cmd.String("from_account_id")
					}
					if cmd.IsSet("reference_number") {
						body["reference_number"] = cmd.String("reference_number")
					}
					if cmd.IsSet("description") {
						body["description"] = cmd.String("description")
					}
					if err := internal.MergeJSON(cmd, body); err != nil {
						return err
					}
					raw, err := c.Request("POST", c.BooksBase+"/banktransactions/"+cmd.Args().First()+"/categorize/vendorpaymentrefund", &zohttp.RequestOpts{
						Params: orgParams(orgID),
						JSON:   body,
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

func bankRulesCmd() *cli.Command {
	return &cli.Command{
		Name:  "bank-rules",
		Usage: "Bank rule operations",
		Commands: []*cli.Command{
			{
				Name:  "create",
				Usage: "Create a bank rule",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "rule_name", Required: true, Usage: "Rule name"},
					&cli.StringFlag{Name: "transaction_type", Usage: "Transaction type"},
					&cli.StringFlag{Name: "criteria", Usage: "Rule criteria"},
					&cli.StringFlag{Name: "apply_to", Usage: "Apply to"},
					&cli.StringFlag{Name: "record_as", Usage: "Record as"},
					&cli.StringFlag{Name: "account_id", Usage: "Account ID"},
					&cli.StringFlag{Name: "customer_id", Usage: "Customer ID"},
					&cli.StringFlag{Name: "vendor_id", Usage: "Vendor ID"},
					&cli.StringFlag{Name: "tax_id", Usage: "Tax ID"},
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
					body["rule_name"] = cmd.String("rule_name")
					if cmd.IsSet("transaction_type") {
						body["transaction_type"] = cmd.String("transaction_type")
					}
					if cmd.IsSet("criteria") {
						var parsed any
						if err := json.Unmarshal([]byte(cmd.String("criteria")), &parsed); err != nil {
							return err
						}
						body["criteria"] = parsed
					}
					if cmd.IsSet("apply_to") {
						body["apply_to"] = cmd.String("apply_to")
					}
					if cmd.IsSet("record_as") {
						body["record_as"] = cmd.String("record_as")
					}
					if cmd.IsSet("account_id") {
						body["account_id"] = cmd.String("account_id")
					}
					if cmd.IsSet("customer_id") {
						body["customer_id"] = cmd.String("customer_id")
					}
					if cmd.IsSet("vendor_id") {
						body["vendor_id"] = cmd.String("vendor_id")
					}
					if cmd.IsSet("tax_id") {
						body["tax_id"] = cmd.String("tax_id")
					}
					if err := internal.MergeJSON(cmd, body); err != nil {
						return err
					}
					raw, err := c.Request("POST", c.BooksBase+"/bankrules", &zohttp.RequestOpts{
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
				Usage: "List bank rules",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "account-id", Usage: "Filter"},
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
					if v := cmd.String("account-id"); v != "" {
						params["account_id"] = v
					}
					raw, err := c.Request("GET", c.BooksBase+"/bankrules", &zohttp.RequestOpts{Params: params})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "update",
				Usage:     "Update a bank rule",
				ArgsUsage: "<rule-id>",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "rule_name", Usage: "Rule name"},
					&cli.StringFlag{Name: "transaction_type", Usage: "Transaction type"},
					&cli.StringFlag{Name: "criteria", Usage: "Rule criteria"},
					&cli.StringFlag{Name: "apply_to", Usage: "Apply to"},
					&cli.StringFlag{Name: "record_as", Usage: "Record as"},
					&cli.StringFlag{Name: "account_id", Usage: "Account ID"},
					&cli.StringFlag{Name: "customer_id", Usage: "Customer ID"},
					&cli.StringFlag{Name: "vendor_id", Usage: "Vendor ID"},
					&cli.StringFlag{Name: "tax_id", Usage: "Tax ID"},
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
					if cmd.IsSet("rule_name") {
						body["rule_name"] = cmd.String("rule_name")
					}
					if cmd.IsSet("transaction_type") {
						body["transaction_type"] = cmd.String("transaction_type")
					}
					if cmd.IsSet("criteria") {
						var parsed any
						if err := json.Unmarshal([]byte(cmd.String("criteria")), &parsed); err != nil {
							return err
						}
						body["criteria"] = parsed
					}
					if cmd.IsSet("apply_to") {
						body["apply_to"] = cmd.String("apply_to")
					}
					if cmd.IsSet("record_as") {
						body["record_as"] = cmd.String("record_as")
					}
					if cmd.IsSet("account_id") {
						body["account_id"] = cmd.String("account_id")
					}
					if cmd.IsSet("customer_id") {
						body["customer_id"] = cmd.String("customer_id")
					}
					if cmd.IsSet("vendor_id") {
						body["vendor_id"] = cmd.String("vendor_id")
					}
					if cmd.IsSet("tax_id") {
						body["tax_id"] = cmd.String("tax_id")
					}
					if err := internal.MergeJSON(cmd, body); err != nil {
						return err
					}
					raw, err := c.Request("PUT", c.BooksBase+"/bankrules/"+cmd.Args().First(), &zohttp.RequestOpts{
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
				Usage:     "Get a bank rule",
				ArgsUsage: "<rule-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("GET", c.BooksBase+"/bankrules/"+cmd.Args().First(), &zohttp.RequestOpts{Params: orgParams(orgID)})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "delete",
				Usage:     "Delete a bank rule",
				ArgsUsage: "<rule-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("DELETE", c.BooksBase+"/bankrules/"+cmd.Args().First(), &zohttp.RequestOpts{Params: orgParams(orgID)})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
		},
	}
}
