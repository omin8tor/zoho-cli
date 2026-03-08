package books

import (
	"context"
	"encoding/json"

	"github.com/omin8tor/zoho-cli/internal"
	zohttp "github.com/omin8tor/zoho-cli/internal/http"
	"github.com/omin8tor/zoho-cli/internal/output"
	"github.com/omin8tor/zoho-cli/internal/pagination"
	"github.com/urfave/cli/v3"
)

func customerPaymentsCmd() *cli.Command {
	return &cli.Command{
		Name:  "customer-payments",
		Usage: "Customer payment operations",
		Commands: []*cli.Command{
			{
				Name:  "create",
				Usage: "Create a customer payment",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "customer_id", Required: true, Usage: "Customer ID"},
					&cli.StringFlag{Name: "payment_mode", Required: true, Usage: "Payment mode"},
					&cli.FloatFlag{Name: "amount", Required: true, Usage: "Payment amount"},
					&cli.StringFlag{Name: "date", Required: true, Usage: "Payment date (YYYY-MM-DD)"},
					&cli.StringFlag{Name: "reference_number", Usage: "Reference number"},
					&cli.StringFlag{Name: "description", Usage: "Payment description"},
					&cli.StringFlag{Name: "account_id", Usage: "Deposit account ID"},
					&cli.StringFlag{Name: "json", Usage: "Additional fields as JSON"},
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					body := map[string]any{}
					body["customer_id"] = cmd.String("customer_id")
					body["payment_mode"] = cmd.String("payment_mode")
					body["amount"] = cmd.Float("amount")
					body["date"] = cmd.String("date")
					if cmd.IsSet("reference_number") {
						body["reference_number"] = cmd.String("reference_number")
					}
					if cmd.IsSet("description") {
						body["description"] = cmd.String("description")
					}
					if cmd.IsSet("account_id") {
						body["account_id"] = cmd.String("account_id")
					}
					if err := internal.MergeJSON(cmd, body); err != nil {
						return err
					}
					raw, err := c.Request(ctx, "POST", c.BooksBase+"/customerpayments", &zohttp.RequestOpts{
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
				Usage: "Update a customer payment by custom field",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "customer_id", Usage: "Customer ID"},
					&cli.StringFlag{Name: "payment_mode", Usage: "Payment mode"},
					&cli.FloatFlag{Name: "amount", Usage: "Payment amount"},
					&cli.StringFlag{Name: "date", Usage: "Payment date (YYYY-MM-DD)"},
					&cli.StringFlag{Name: "reference_number", Usage: "Reference number"},
					&cli.StringFlag{Name: "description", Usage: "Payment description"},
					&cli.StringFlag{Name: "account_id", Usage: "Deposit account ID"},
					&cli.StringFlag{Name: "json", Usage: "Additional fields as JSON"},
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
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
					if cmd.IsSet("payment_mode") {
						body["payment_mode"] = cmd.String("payment_mode")
					}
					if cmd.IsSet("amount") {
						body["amount"] = cmd.Float("amount")
					}
					if cmd.IsSet("date") {
						body["date"] = cmd.String("date")
					}
					if cmd.IsSet("reference_number") {
						body["reference_number"] = cmd.String("reference_number")
					}
					if cmd.IsSet("description") {
						body["description"] = cmd.String("description")
					}
					if cmd.IsSet("account_id") {
						body["account_id"] = cmd.String("account_id")
					}
					if err := internal.MergeJSON(cmd, body); err != nil {
						return err
					}
					raw, err := c.Request(ctx, "PUT", c.BooksBase+"/customerpayments", &zohttp.RequestOpts{
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
				Usage: "List customer payments",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "customer-id", Usage: "Filter"},
					&cli.BoolFlag{Name: "all", Usage: "Fetch all records"},
					&cli.IntFlag{Name: "limit", Usage: "Max total records to fetch"},
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					params := orgParams(orgID)
					if v := cmd.String("customer-id"); v != "" {
						params["customer_id"] = v
					}

					if cmd.Bool("all") || cmd.IsSet("limit") {

						items, err := pagination.Paginate(ctx, pagination.PaginationConfig{

							Client: c,

							URL: c.BooksBase + "/customerpayments",

							Opts: &zohttp.RequestOpts{Params: params},

							ItemsKey: "customerpayments",

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
					raw, err := c.Request(ctx, "GET", c.BooksBase+"/customerpayments", &zohttp.RequestOpts{Params: params})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "update",
				Usage:     "Update a customer payment",
				ArgsUsage: "<payment-id>",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "customer_id", Usage: "Customer ID"},
					&cli.StringFlag{Name: "payment_mode", Usage: "Payment mode"},
					&cli.FloatFlag{Name: "amount", Usage: "Payment amount"},
					&cli.StringFlag{Name: "date", Usage: "Payment date (YYYY-MM-DD)"},
					&cli.StringFlag{Name: "reference_number", Usage: "Reference number"},
					&cli.StringFlag{Name: "description", Usage: "Payment description"},
					&cli.StringFlag{Name: "account_id", Usage: "Deposit account ID"},
					&cli.StringFlag{Name: "json", Usage: "Additional fields as JSON"},
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
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
					if cmd.IsSet("payment_mode") {
						body["payment_mode"] = cmd.String("payment_mode")
					}
					if cmd.IsSet("amount") {
						body["amount"] = cmd.Float("amount")
					}
					if cmd.IsSet("date") {
						body["date"] = cmd.String("date")
					}
					if cmd.IsSet("reference_number") {
						body["reference_number"] = cmd.String("reference_number")
					}
					if cmd.IsSet("description") {
						body["description"] = cmd.String("description")
					}
					if cmd.IsSet("account_id") {
						body["account_id"] = cmd.String("account_id")
					}
					if err := internal.MergeJSON(cmd, body); err != nil {
						return err
					}
					raw, err := c.Request(ctx, "PUT", c.BooksBase+"/customerpayments/"+cmd.Args().First(), &zohttp.RequestOpts{
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
				Usage:     "Get a customer payment",
				ArgsUsage: "<payment-id>",
				Action: func(ctx context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request(ctx, "GET", c.BooksBase+"/customerpayments/"+cmd.Args().First(), &zohttp.RequestOpts{Params: orgParams(orgID)})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "delete",
				Usage:     "Delete a customer payment",
				ArgsUsage: "<payment-id>",
				Action: func(ctx context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request(ctx, "DELETE", c.BooksBase+"/customerpayments/"+cmd.Args().First(), &zohttp.RequestOpts{Params: orgParams(orgID)})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "refund-excess",
				Usage:     "Refund excess of a customer payment",
				ArgsUsage: "<payment-id>",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "date", Required: true, Usage: "Date (YYYY-MM-DD)"},
					&cli.StringFlag{Name: "refund_mode", Usage: "Refund mode"},
					&cli.StringFlag{Name: "from_account_id", Usage: "From account ID"},
					&cli.FloatFlag{Name: "amount", Required: true, Usage: "Amount"},
					&cli.StringFlag{Name: "reference_number", Usage: "Reference number"},
					&cli.StringFlag{Name: "description", Usage: "Description"},
					&cli.StringFlag{Name: "json", Usage: "Additional fields as JSON"},
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					body := map[string]any{}
					body["date"] = cmd.String("date")
					if cmd.IsSet("refund_mode") {
						body["refund_mode"] = cmd.String("refund_mode")
					}
					if cmd.IsSet("from_account_id") {
						body["from_account_id"] = cmd.String("from_account_id")
					}
					body["amount"] = cmd.Float("amount")
					if cmd.IsSet("reference_number") {
						body["reference_number"] = cmd.String("reference_number")
					}
					if cmd.IsSet("description") {
						body["description"] = cmd.String("description")
					}
					if err := internal.MergeJSON(cmd, body); err != nil {
						return err
					}
					raw, err := c.Request(ctx, "POST", c.BooksBase+"/customerpayments/"+cmd.Args().First()+"/refunds", &zohttp.RequestOpts{
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
				Name:      "list-refunds",
				Usage:     "List refunds of a customer payment",
				ArgsUsage: "<payment-id>",
				Action: func(ctx context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request(ctx, "GET", c.BooksBase+"/customerpayments/"+cmd.Args().First()+"/refunds", &zohttp.RequestOpts{Params: orgParams(orgID)})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "update-custom-fields",
				Usage:     "Update custom fields of a customer payment",
				ArgsUsage: "<payment-id>",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "custom_fields", Usage: "Custom fields JSON"},
					&cli.StringFlag{Name: "json", Usage: "Additional fields as JSON"},
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					body := map[string]any{}
					if cmd.IsSet("custom_fields") {
						var parsed any
						if err := json.Unmarshal([]byte(cmd.String("custom_fields")), &parsed); err != nil {
							return err
						}
						body["custom_fields"] = parsed
					}
					if err := internal.MergeJSON(cmd, body); err != nil {
						return err
					}
					raw, err := c.Request(ctx, "PUT", c.BooksBase+"/customerpayments/"+cmd.Args().First()+"/customfields", &zohttp.RequestOpts{
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
				Name:      "update-refund",
				Usage:     "Update a refund of a customer payment",
				ArgsUsage: "<payment-id> <refund-id>",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "date", Usage: "Refund date (YYYY-MM-DD)"},
					&cli.StringFlag{Name: "refund_mode", Usage: "Refund mode"},
					&cli.StringFlag{Name: "from_account_id", Usage: "Account ID for refund"},
					&cli.FloatFlag{Name: "amount", Usage: "Refund amount"},
					&cli.StringFlag{Name: "reference_number", Usage: "Reference number"},
					&cli.StringFlag{Name: "description", Usage: "Refund description"},
					&cli.StringFlag{Name: "json", Usage: "Additional fields as JSON"},
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					body := map[string]any{}
					if cmd.IsSet("date") {
						body["date"] = cmd.String("date")
					}
					if cmd.IsSet("refund_mode") {
						body["refund_mode"] = cmd.String("refund_mode")
					}
					if cmd.IsSet("from_account_id") {
						body["from_account_id"] = cmd.String("from_account_id")
					}
					if cmd.IsSet("amount") {
						body["amount"] = cmd.Float("amount")
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
					raw, err := c.Request(ctx, "PUT", c.BooksBase+"/customerpayments/"+cmd.Args().First()+"/refunds/"+cmd.Args().Get(1), &zohttp.RequestOpts{
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
				Name:      "get-refund",
				Usage:     "Get a refund of a customer payment",
				ArgsUsage: "<payment-id> <refund-id>",
				Action: func(ctx context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request(ctx, "GET", c.BooksBase+"/customerpayments/"+cmd.Args().First()+"/refunds/"+cmd.Args().Get(1), &zohttp.RequestOpts{Params: orgParams(orgID)})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "delete-refund",
				Usage:     "Delete a refund of a customer payment",
				ArgsUsage: "<payment-id> <refund-id>",
				Action: func(ctx context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request(ctx, "DELETE", c.BooksBase+"/customerpayments/"+cmd.Args().First()+"/refunds/"+cmd.Args().Get(1), &zohttp.RequestOpts{Params: orgParams(orgID)})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
		},
	}
}

func vendorPaymentsCmd() *cli.Command {
	return &cli.Command{
		Name:  "vendor-payments",
		Usage: "Vendor payment operations",
		Commands: []*cli.Command{
			{
				Name:  "create",
				Usage: "Create a vendor payment",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "vendor_id", Required: true, Usage: "Vendor ID"},
					&cli.StringFlag{Name: "payment_mode", Required: true, Usage: "Payment mode"},
					&cli.FloatFlag{Name: "amount", Required: true, Usage: "Payment amount"},
					&cli.StringFlag{Name: "date", Required: true, Usage: "Payment date (YYYY-MM-DD)"},
					&cli.StringFlag{Name: "reference_number", Usage: "Reference number"},
					&cli.StringFlag{Name: "description", Usage: "Payment description"},
					&cli.StringFlag{Name: "account_id", Usage: "Payment account ID"},
					&cli.StringFlag{Name: "json", Usage: "Additional fields as JSON"},
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					body := map[string]any{}
					body["vendor_id"] = cmd.String("vendor_id")
					body["payment_mode"] = cmd.String("payment_mode")
					body["amount"] = cmd.Float("amount")
					body["date"] = cmd.String("date")
					if cmd.IsSet("reference_number") {
						body["reference_number"] = cmd.String("reference_number")
					}
					if cmd.IsSet("description") {
						body["description"] = cmd.String("description")
					}
					if cmd.IsSet("account_id") {
						body["account_id"] = cmd.String("account_id")
					}
					if err := internal.MergeJSON(cmd, body); err != nil {
						return err
					}
					raw, err := c.Request(ctx, "POST", c.BooksBase+"/vendorpayments", &zohttp.RequestOpts{
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
				Usage: "Update a vendor payment by custom field",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "vendor_id", Usage: "Vendor ID"},
					&cli.StringFlag{Name: "payment_mode", Usage: "Payment mode"},
					&cli.FloatFlag{Name: "amount", Usage: "Payment amount"},
					&cli.StringFlag{Name: "date", Usage: "Payment date (YYYY-MM-DD)"},
					&cli.StringFlag{Name: "reference_number", Usage: "Reference number"},
					&cli.StringFlag{Name: "description", Usage: "Payment description"},
					&cli.StringFlag{Name: "account_id", Usage: "Payment account ID"},
					&cli.StringFlag{Name: "json", Usage: "Additional fields as JSON"},
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
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
					if cmd.IsSet("payment_mode") {
						body["payment_mode"] = cmd.String("payment_mode")
					}
					if cmd.IsSet("amount") {
						body["amount"] = cmd.Float("amount")
					}
					if cmd.IsSet("date") {
						body["date"] = cmd.String("date")
					}
					if cmd.IsSet("reference_number") {
						body["reference_number"] = cmd.String("reference_number")
					}
					if cmd.IsSet("description") {
						body["description"] = cmd.String("description")
					}
					if cmd.IsSet("account_id") {
						body["account_id"] = cmd.String("account_id")
					}
					if err := internal.MergeJSON(cmd, body); err != nil {
						return err
					}
					raw, err := c.Request(ctx, "PUT", c.BooksBase+"/vendorpayments", &zohttp.RequestOpts{
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
				Usage: "List vendor payments",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "vendor-id", Usage: "Filter"},
					&cli.BoolFlag{Name: "all", Usage: "Fetch all records"},
					&cli.IntFlag{Name: "limit", Usage: "Max total records to fetch"},
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					params := orgParams(orgID)
					if v := cmd.String("vendor-id"); v != "" {
						params["vendor_id"] = v
					}

					if cmd.Bool("all") || cmd.IsSet("limit") {

						items, err := pagination.Paginate(ctx, pagination.PaginationConfig{

							Client: c,

							URL: c.BooksBase + "/vendorpayments",

							Opts: &zohttp.RequestOpts{Params: params},

							ItemsKey: "vendorpayments",

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
					raw, err := c.Request(ctx, "GET", c.BooksBase+"/vendorpayments", &zohttp.RequestOpts{Params: params})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "update",
				Usage:     "Update a vendor payment",
				ArgsUsage: "<payment-id>",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "vendor_id", Usage: "Vendor ID"},
					&cli.StringFlag{Name: "payment_mode", Usage: "Payment mode"},
					&cli.FloatFlag{Name: "amount", Usage: "Payment amount"},
					&cli.StringFlag{Name: "date", Usage: "Payment date (YYYY-MM-DD)"},
					&cli.StringFlag{Name: "reference_number", Usage: "Reference number"},
					&cli.StringFlag{Name: "description", Usage: "Payment description"},
					&cli.StringFlag{Name: "account_id", Usage: "Payment account ID"},
					&cli.StringFlag{Name: "json", Usage: "Additional fields as JSON"},
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
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
					if cmd.IsSet("payment_mode") {
						body["payment_mode"] = cmd.String("payment_mode")
					}
					if cmd.IsSet("amount") {
						body["amount"] = cmd.Float("amount")
					}
					if cmd.IsSet("date") {
						body["date"] = cmd.String("date")
					}
					if cmd.IsSet("reference_number") {
						body["reference_number"] = cmd.String("reference_number")
					}
					if cmd.IsSet("description") {
						body["description"] = cmd.String("description")
					}
					if cmd.IsSet("account_id") {
						body["account_id"] = cmd.String("account_id")
					}
					if err := internal.MergeJSON(cmd, body); err != nil {
						return err
					}
					raw, err := c.Request(ctx, "PUT", c.BooksBase+"/vendorpayments/"+cmd.Args().First(), &zohttp.RequestOpts{
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
				Usage:     "Get a vendor payment",
				ArgsUsage: "<payment-id>",
				Action: func(ctx context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request(ctx, "GET", c.BooksBase+"/vendorpayments/"+cmd.Args().First(), &zohttp.RequestOpts{Params: orgParams(orgID)})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "delete",
				Usage:     "Delete a vendor payment",
				ArgsUsage: "<payment-id>",
				Action: func(ctx context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request(ctx, "DELETE", c.BooksBase+"/vendorpayments/"+cmd.Args().First(), &zohttp.RequestOpts{Params: orgParams(orgID)})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "refund-excess",
				Usage:     "Refund excess of a vendor payment",
				ArgsUsage: "<payment-id>",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "date", Required: true, Usage: "Date (YYYY-MM-DD)"},
					&cli.StringFlag{Name: "refund_mode", Usage: "Refund mode"},
					&cli.StringFlag{Name: "from_account_id", Usage: "From account ID"},
					&cli.FloatFlag{Name: "amount", Required: true, Usage: "Amount"},
					&cli.StringFlag{Name: "reference_number", Usage: "Reference number"},
					&cli.StringFlag{Name: "description", Usage: "Description"},
					&cli.StringFlag{Name: "json", Usage: "Additional fields as JSON"},
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					body := map[string]any{}
					body["date"] = cmd.String("date")
					if cmd.IsSet("refund_mode") {
						body["refund_mode"] = cmd.String("refund_mode")
					}
					if cmd.IsSet("from_account_id") {
						body["from_account_id"] = cmd.String("from_account_id")
					}
					body["amount"] = cmd.Float("amount")
					if cmd.IsSet("reference_number") {
						body["reference_number"] = cmd.String("reference_number")
					}
					if cmd.IsSet("description") {
						body["description"] = cmd.String("description")
					}
					if err := internal.MergeJSON(cmd, body); err != nil {
						return err
					}
					raw, err := c.Request(ctx, "POST", c.BooksBase+"/vendorpayments/"+cmd.Args().First()+"/refunds", &zohttp.RequestOpts{
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
				Name:      "list-refunds",
				Usage:     "List refunds of a vendor payment",
				ArgsUsage: "<payment-id>",
				Action: func(ctx context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request(ctx, "GET", c.BooksBase+"/vendorpayments/"+cmd.Args().First()+"/refunds", &zohttp.RequestOpts{Params: orgParams(orgID)})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "update-refund",
				Usage:     "Update a refund of a vendor payment",
				ArgsUsage: "<payment-id> <refund-id>",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "date", Usage: "Refund date (YYYY-MM-DD)"},
					&cli.StringFlag{Name: "refund_mode", Usage: "Refund mode"},
					&cli.StringFlag{Name: "from_account_id", Usage: "Account ID for refund"},
					&cli.FloatFlag{Name: "amount", Usage: "Refund amount"},
					&cli.StringFlag{Name: "reference_number", Usage: "Reference number"},
					&cli.StringFlag{Name: "description", Usage: "Refund description"},
					&cli.StringFlag{Name: "json", Usage: "Additional fields as JSON"},
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					body := map[string]any{}
					if cmd.IsSet("date") {
						body["date"] = cmd.String("date")
					}
					if cmd.IsSet("refund_mode") {
						body["refund_mode"] = cmd.String("refund_mode")
					}
					if cmd.IsSet("from_account_id") {
						body["from_account_id"] = cmd.String("from_account_id")
					}
					if cmd.IsSet("amount") {
						body["amount"] = cmd.Float("amount")
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
					raw, err := c.Request(ctx, "PUT", c.BooksBase+"/vendorpayments/"+cmd.Args().First()+"/refunds/"+cmd.Args().Get(1), &zohttp.RequestOpts{
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
				Name:      "get-refund",
				Usage:     "Get a refund of a vendor payment",
				ArgsUsage: "<payment-id> <refund-id>",
				Action: func(ctx context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request(ctx, "GET", c.BooksBase+"/vendorpayments/"+cmd.Args().First()+"/refunds/"+cmd.Args().Get(1), &zohttp.RequestOpts{Params: orgParams(orgID)})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "delete-refund",
				Usage:     "Delete a refund of a vendor payment",
				ArgsUsage: "<payment-id> <refund-id>",
				Action: func(ctx context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request(ctx, "DELETE", c.BooksBase+"/vendorpayments/"+cmd.Args().First()+"/refunds/"+cmd.Args().Get(1), &zohttp.RequestOpts{Params: orgParams(orgID)})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "email",
				Usage:     "Email a vendor payment",
				ArgsUsage: "<payment-id>",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "to_mail_ids", Required: true, Usage: "Recipient email addresses (comma-separated)"},
					&cli.StringFlag{Name: "cc_mail_ids", Usage: "CC email addresses (comma-separated)"},
					&cli.StringFlag{Name: "subject", Required: true, Usage: "Email subject"},
					&cli.StringFlag{Name: "body", Required: true, Usage: "Email body"},
					&cli.StringFlag{Name: "json", Usage: "Additional fields as JSON"},
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					body := map[string]any{}
					body["to_mail_ids"] = cmd.String("to_mail_ids")
					if cmd.IsSet("cc_mail_ids") {
						body["cc_mail_ids"] = cmd.String("cc_mail_ids")
					}
					body["subject"] = cmd.String("subject")
					body["body"] = cmd.String("body")
					if err := internal.MergeJSON(cmd, body); err != nil {
						return err
					}
					raw, err := c.Request(ctx, "POST", c.BooksBase+"/vendorpayments/"+cmd.Args().First()+"/email", &zohttp.RequestOpts{
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
				Name:      "get-email-content",
				Usage:     "Get email content of a vendor payment",
				ArgsUsage: "<payment-id>",
				Action: func(ctx context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request(ctx, "GET", c.BooksBase+"/vendorpayments/"+cmd.Args().First()+"/email", &zohttp.RequestOpts{Params: orgParams(orgID)})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
		},
	}
}
