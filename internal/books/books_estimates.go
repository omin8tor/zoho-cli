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

func estimatesCmd() *cli.Command {
	return &cli.Command{
		Name:  "estimates",
		Usage: "Estimate operations",
		Commands: []*cli.Command{
			{
				Name:  "create",
				Usage: "Create an estimate",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "customer_id", Required: true, Usage: "Customer ID"},
					&cli.StringFlag{Name: "date", Usage: "Estimate date (YYYY-MM-DD)"},
					&cli.StringFlag{Name: "expiry_date", Usage: "Expiry date (YYYY-MM-DD)"},
					&cli.StringFlag{Name: "reference_number", Usage: "Reference number"},
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
					if cmd.IsSet("date") {
						body["date"] = cmd.String("date")
					}
					if cmd.IsSet("expiry_date") {
						body["expiry_date"] = cmd.String("expiry_date")
					}
					if cmd.IsSet("reference_number") {
						body["reference_number"] = cmd.String("reference_number")
					}
					if err := internal.MergeJSON(cmd, body); err != nil {
						return err
					}
					raw, err := c.Request(ctx, "POST", c.BooksBase+"/estimates", &zohttp.RequestOpts{
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
				Usage: "Update an estimate by custom field",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "customer_id", Usage: "Customer ID"},
					&cli.StringFlag{Name: "date", Usage: "Estimate date (YYYY-MM-DD)"},
					&cli.StringFlag{Name: "expiry_date", Usage: "Expiry date (YYYY-MM-DD)"},
					&cli.StringFlag{Name: "reference_number", Usage: "Reference number"},
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
					if cmd.IsSet("date") {
						body["date"] = cmd.String("date")
					}
					if cmd.IsSet("expiry_date") {
						body["expiry_date"] = cmd.String("expiry_date")
					}
					if cmd.IsSet("reference_number") {
						body["reference_number"] = cmd.String("reference_number")
					}
					if err := internal.MergeJSON(cmd, body); err != nil {
						return err
					}
					raw, err := c.Request(ctx, "PUT", c.BooksBase+"/estimates", &zohttp.RequestOpts{
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
				Usage: "List estimates",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "status", Usage: "Filter"},
					&cli.StringFlag{Name: "customer-id", Usage: "Filter"},
					&cli.StringFlag{Name: "sort-column", Usage: "Filter"},
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
					if v := cmd.String("status"); v != "" {
						params["status"] = v
					}
					if v := cmd.String("customer-id"); v != "" {
						params["customer_id"] = v
					}
					if v := cmd.String("sort-column"); v != "" {
						params["sort_column"] = v
					}

					if cmd.Bool("all") || cmd.IsSet("limit") {

						items, err := pagination.Paginate(ctx, pagination.PaginationConfig{

							Client: c,

							URL: c.BooksBase + "/estimates",

							Opts: &zohttp.RequestOpts{Params: params},

							ItemsKey: "estimates",

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
					raw, err := c.Request(ctx, "GET", c.BooksBase+"/estimates", &zohttp.RequestOpts{Params: params})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "update",
				Usage:     "Update an estimate",
				ArgsUsage: "<estimate-id>",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "customer_id", Usage: "Customer ID"},
					&cli.StringFlag{Name: "date", Usage: "Estimate date (YYYY-MM-DD)"},
					&cli.StringFlag{Name: "expiry_date", Usage: "Expiry date (YYYY-MM-DD)"},
					&cli.StringFlag{Name: "reference_number", Usage: "Reference number"},
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
					if cmd.IsSet("date") {
						body["date"] = cmd.String("date")
					}
					if cmd.IsSet("expiry_date") {
						body["expiry_date"] = cmd.String("expiry_date")
					}
					if cmd.IsSet("reference_number") {
						body["reference_number"] = cmd.String("reference_number")
					}
					if err := internal.MergeJSON(cmd, body); err != nil {
						return err
					}
					raw, err := c.Request(ctx, "PUT", c.BooksBase+"/estimates/"+cmd.Args().First(), &zohttp.RequestOpts{
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
				Usage:     "Get an estimate",
				ArgsUsage: "<estimate-id>",
				Action: func(ctx context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request(ctx, "GET", c.BooksBase+"/estimates/"+cmd.Args().First(), &zohttp.RequestOpts{Params: orgParams(orgID)})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "delete",
				Usage:     "Delete an estimate",
				ArgsUsage: "<estimate-id>",
				Action: func(ctx context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request(ctx, "DELETE", c.BooksBase+"/estimates/"+cmd.Args().First(), &zohttp.RequestOpts{Params: orgParams(orgID)})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "update-custom-fields",
				Usage:     "Update custom fields of an estimate",
				ArgsUsage: "<estimate-id>",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "customfield_id", Usage: "Custom field ID"},
					&cli.StringFlag{Name: "value", Usage: "Custom field value"},
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
					if cmd.IsSet("customfield_id") {
						body["customfield_id"] = cmd.String("customfield_id")
					}
					if cmd.IsSet("value") {
						body["value"] = cmd.String("value")
					}
					if err := internal.MergeJSON(cmd, body); err != nil {
						return err
					}
					raw, err := c.Request(ctx, "PUT", c.BooksBase+"/estimates/"+cmd.Args().First()+"/customfields", &zohttp.RequestOpts{
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
				Name:      "mark-sent",
				Usage:     "Mark an estimate as sent",
				ArgsUsage: "<estimate-id>",
				Action: func(ctx context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request(ctx, "POST", c.BooksBase+"/estimates/"+cmd.Args().First()+"/status/sent", &zohttp.RequestOpts{Params: orgParams(orgID)})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "mark-accepted",
				Usage:     "Mark an estimate as accepted",
				ArgsUsage: "<estimate-id>",
				Action: func(ctx context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request(ctx, "POST", c.BooksBase+"/estimates/"+cmd.Args().First()+"/status/accepted", &zohttp.RequestOpts{Params: orgParams(orgID)})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "mark-declined",
				Usage:     "Mark an estimate as declined",
				ArgsUsage: "<estimate-id>",
				Action: func(ctx context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request(ctx, "POST", c.BooksBase+"/estimates/"+cmd.Args().First()+"/status/declined", &zohttp.RequestOpts{Params: orgParams(orgID)})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "submit-for-approval",
				Usage:     "Submit an estimate for approval",
				ArgsUsage: "<estimate-id>",
				Action: func(ctx context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request(ctx, "POST", c.BooksBase+"/estimates/"+cmd.Args().First()+"/submit", &zohttp.RequestOpts{Params: orgParams(orgID)})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "approve",
				Usage:     "Approve an estimate",
				ArgsUsage: "<estimate-id>",
				Action: func(ctx context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request(ctx, "POST", c.BooksBase+"/estimates/"+cmd.Args().First()+"/approve", &zohttp.RequestOpts{Params: orgParams(orgID)})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "email",
				Usage:     "Email an estimate",
				ArgsUsage: "<estimate-id>",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "to_mail_ids", Required: true, Usage: "Recipient email addresses (comma-separated)"},
					&cli.StringFlag{Name: "subject", Required: true, Usage: "Email subject"},
					&cli.StringFlag{Name: "body", Required: true, Usage: "Email body"},
					&cli.StringFlag{Name: "cc_mail_ids", Usage: "CC email addresses (comma-separated)"},
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
					body["subject"] = cmd.String("subject")
					body["body"] = cmd.String("body")
					if cmd.IsSet("cc_mail_ids") {
						body["cc_mail_ids"] = cmd.String("cc_mail_ids")
					}
					if err := internal.MergeJSON(cmd, body); err != nil {
						return err
					}
					raw, err := c.Request(ctx, "POST", c.BooksBase+"/estimates/"+cmd.Args().First()+"/email", &zohttp.RequestOpts{
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
				Usage:     "Get email content of an estimate",
				ArgsUsage: "<estimate-id>",
				Action: func(ctx context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request(ctx, "GET", c.BooksBase+"/estimates/"+cmd.Args().First()+"/email", &zohttp.RequestOpts{Params: orgParams(orgID)})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:  "email-multiple",
				Usage: "Email multiple estimates",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "estimate_ids", Required: true, Usage: "Estimate IDs (comma-separated)"},
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
					var parsed0 any
					if err := json.Unmarshal([]byte(cmd.String("estimate_ids")), &parsed0); err != nil {
						return err
					}
					body["estimate_ids"] = parsed0
					if err := internal.MergeJSON(cmd, body); err != nil {
						return err
					}
					raw, err := c.Request(ctx, "POST", c.BooksBase+"/estimates"+"/email", &zohttp.RequestOpts{
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
				Name:      "update-billing-address",
				Usage:     "Update billing address of an estimate",
				ArgsUsage: "<estimate-id>",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "address", Usage: "Street address"},
					&cli.StringFlag{Name: "city", Usage: "City"},
					&cli.StringFlag{Name: "state", Usage: "State"},
					&cli.StringFlag{Name: "zip", Usage: "Zip/postal code"},
					&cli.StringFlag{Name: "country", Usage: "Country"},
					&cli.StringFlag{Name: "fax", Usage: "Fax number"},
					&cli.StringFlag{Name: "attention", Usage: "Attention"},
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
					if cmd.IsSet("address") {
						body["address"] = cmd.String("address")
					}
					if cmd.IsSet("city") {
						body["city"] = cmd.String("city")
					}
					if cmd.IsSet("state") {
						body["state"] = cmd.String("state")
					}
					if cmd.IsSet("zip") {
						body["zip"] = cmd.String("zip")
					}
					if cmd.IsSet("country") {
						body["country"] = cmd.String("country")
					}
					if cmd.IsSet("fax") {
						body["fax"] = cmd.String("fax")
					}
					if cmd.IsSet("attention") {
						body["attention"] = cmd.String("attention")
					}
					if err := internal.MergeJSON(cmd, body); err != nil {
						return err
					}
					raw, err := c.Request(ctx, "PUT", c.BooksBase+"/estimates/"+cmd.Args().First()+"/address/billing", &zohttp.RequestOpts{
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
				Name:      "update-shipping-address",
				Usage:     "Update shipping address of an estimate",
				ArgsUsage: "<estimate-id>",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "address", Usage: "Street address"},
					&cli.StringFlag{Name: "city", Usage: "City"},
					&cli.StringFlag{Name: "state", Usage: "State"},
					&cli.StringFlag{Name: "zip", Usage: "Zip/postal code"},
					&cli.StringFlag{Name: "country", Usage: "Country"},
					&cli.StringFlag{Name: "fax", Usage: "Fax number"},
					&cli.StringFlag{Name: "attention", Usage: "Attention"},
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
					if cmd.IsSet("address") {
						body["address"] = cmd.String("address")
					}
					if cmd.IsSet("city") {
						body["city"] = cmd.String("city")
					}
					if cmd.IsSet("state") {
						body["state"] = cmd.String("state")
					}
					if cmd.IsSet("zip") {
						body["zip"] = cmd.String("zip")
					}
					if cmd.IsSet("country") {
						body["country"] = cmd.String("country")
					}
					if cmd.IsSet("fax") {
						body["fax"] = cmd.String("fax")
					}
					if cmd.IsSet("attention") {
						body["attention"] = cmd.String("attention")
					}
					if err := internal.MergeJSON(cmd, body); err != nil {
						return err
					}
					raw, err := c.Request(ctx, "PUT", c.BooksBase+"/estimates/"+cmd.Args().First()+"/address/shipping", &zohttp.RequestOpts{
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
				Name:  "list-templates",
				Usage: "List estimate templates",
				Action: func(ctx context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request(ctx, "GET", c.BooksBase+"/estimates"+"/templates", &zohttp.RequestOpts{Params: orgParams(orgID)})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "update-template",
				Usage:     "Update template of an estimate",
				ArgsUsage: "<estimate-id> <template-id>",
				Action: func(ctx context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request(ctx, "PUT", c.BooksBase+"/estimates/"+cmd.Args().First()+"/templates/"+cmd.Args().Get(1), &zohttp.RequestOpts{Params: orgParams(orgID)})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "add-comment",
				Usage:     "Add a comment to an estimate",
				ArgsUsage: "<estimate-id>",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "description", Required: true, Usage: "Comment text"},
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
					body["description"] = cmd.String("description")
					if err := internal.MergeJSON(cmd, body); err != nil {
						return err
					}
					raw, err := c.Request(ctx, "POST", c.BooksBase+"/estimates/"+cmd.Args().First()+"/comments", &zohttp.RequestOpts{
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
				Name:      "list-comments",
				Usage:     "List comments of an estimate",
				ArgsUsage: "<estimate-id>",
				Action: func(ctx context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request(ctx, "GET", c.BooksBase+"/estimates/"+cmd.Args().First()+"/comments", &zohttp.RequestOpts{Params: orgParams(orgID)})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "update-comment",
				Usage:     "Update a comment on an estimate",
				ArgsUsage: "<estimate-id> <comment-id>",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "description", Required: true, Usage: "Comment text"},
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
					body["description"] = cmd.String("description")
					if err := internal.MergeJSON(cmd, body); err != nil {
						return err
					}
					raw, err := c.Request(ctx, "PUT", c.BooksBase+"/estimates/"+cmd.Args().First()+"/comments/"+cmd.Args().Get(1), &zohttp.RequestOpts{
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
				Name:      "delete-comment",
				Usage:     "Delete a comment on an estimate",
				ArgsUsage: "<estimate-id> <comment-id>",
				Action: func(ctx context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request(ctx, "DELETE", c.BooksBase+"/estimates/"+cmd.Args().First()+"/comments/"+cmd.Args().Get(1), &zohttp.RequestOpts{Params: orgParams(orgID)})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
		},
	}
}
