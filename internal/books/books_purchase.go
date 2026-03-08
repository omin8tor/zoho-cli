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

func purchaseOrdersCmd() *cli.Command {
	return &cli.Command{
		Name:  "purchase-orders",
		Usage: "Purchase order operations",
		Commands: []*cli.Command{
			{
				Name:  "create",
				Usage: "Create a purchase order",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "vendor_id", Required: true, Usage: "Vendor ID"},
					&cli.StringFlag{Name: "date", Usage: "Purchase order date (YYYY-MM-DD)"},
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
					body["vendor_id"] = cmd.String("vendor_id")
					if cmd.IsSet("date") {
						body["date"] = cmd.String("date")
					}
					if cmd.IsSet("reference_number") {
						body["reference_number"] = cmd.String("reference_number")
					}
					if err := internal.MergeJSON(cmd, body); err != nil {
						return err
					}
					raw, err := c.Request("POST", c.BooksBase+"/purchaseorders", &zohttp.RequestOpts{
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
				Usage: "Update a purchase order by custom field",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "vendor_id", Usage: "Vendor ID"},
					&cli.StringFlag{Name: "date", Usage: "Purchase order date (YYYY-MM-DD)"},
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
					if cmd.IsSet("vendor_id") {
						body["vendor_id"] = cmd.String("vendor_id")
					}
					if cmd.IsSet("date") {
						body["date"] = cmd.String("date")
					}
					if cmd.IsSet("reference_number") {
						body["reference_number"] = cmd.String("reference_number")
					}
					if err := internal.MergeJSON(cmd, body); err != nil {
						return err
					}
					raw, err := c.Request("PUT", c.BooksBase+"/purchaseorders", &zohttp.RequestOpts{
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
				Usage: "List purchase orders",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "status", Usage: "Filter"},
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
					if v := cmd.String("vendor-id"); v != "" {
						params["vendor_id"] = v
					}

					if cmd.Bool("all") || cmd.IsSet("limit") {

						items, err := pagination.Paginate(pagination.PaginationConfig{

							Client: c,

							URL: c.BooksBase + "/purchaseorders",

							Opts: &zohttp.RequestOpts{Params: params},

							ItemsKey: "purchaseorders",

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
					raw, err := c.Request("GET", c.BooksBase+"/purchaseorders", &zohttp.RequestOpts{Params: params})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "update",
				Usage:     "Update a purchase order",
				ArgsUsage: "<purchaseorder-id>",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "vendor_id", Usage: "Vendor ID"},
					&cli.StringFlag{Name: "date", Usage: "Purchase order date (YYYY-MM-DD)"},
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
					if cmd.IsSet("vendor_id") {
						body["vendor_id"] = cmd.String("vendor_id")
					}
					if cmd.IsSet("date") {
						body["date"] = cmd.String("date")
					}
					if cmd.IsSet("reference_number") {
						body["reference_number"] = cmd.String("reference_number")
					}
					if err := internal.MergeJSON(cmd, body); err != nil {
						return err
					}
					raw, err := c.Request("PUT", c.BooksBase+"/purchaseorders/"+cmd.Args().First(), &zohttp.RequestOpts{
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
				Usage:     "Get a purchase order",
				ArgsUsage: "<purchaseorder-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("GET", c.BooksBase+"/purchaseorders/"+cmd.Args().First(), &zohttp.RequestOpts{Params: orgParams(orgID)})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "delete",
				Usage:     "Delete a purchase order",
				ArgsUsage: "<purchaseorder-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("DELETE", c.BooksBase+"/purchaseorders/"+cmd.Args().First(), &zohttp.RequestOpts{Params: orgParams(orgID)})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "update-custom-fields",
				Usage:     "Update custom fields of a purchase order",
				ArgsUsage: "<purchaseorder-id>",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "custom_fields", Usage: "Custom fields JSON"},
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
					raw, err := c.Request("PUT", c.BooksBase+"/purchaseorders/"+cmd.Args().First()+"/customfields", &zohttp.RequestOpts{
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
				Name:      "mark-open",
				Usage:     "Mark a purchase order as open",
				ArgsUsage: "<purchaseorder-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("POST", c.BooksBase+"/purchaseorders/"+cmd.Args().First()+"/status/open", &zohttp.RequestOpts{Params: orgParams(orgID)})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "mark-billed",
				Usage:     "Mark a purchase order as billed",
				ArgsUsage: "<purchaseorder-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("POST", c.BooksBase+"/purchaseorders/"+cmd.Args().First()+"/status/billed", &zohttp.RequestOpts{Params: orgParams(orgID)})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "cancel",
				Usage:     "Cancel a purchase order",
				ArgsUsage: "<purchaseorder-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("POST", c.BooksBase+"/purchaseorders/"+cmd.Args().First()+"/status/cancelled", &zohttp.RequestOpts{Params: orgParams(orgID)})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "submit-for-approval",
				Usage:     "Submit a purchase order for approval",
				ArgsUsage: "<purchaseorder-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("POST", c.BooksBase+"/purchaseorders/"+cmd.Args().First()+"/submit", &zohttp.RequestOpts{Params: orgParams(orgID)})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "approve",
				Usage:     "Approve a purchase order",
				ArgsUsage: "<purchaseorder-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("POST", c.BooksBase+"/purchaseorders/"+cmd.Args().First()+"/approve", &zohttp.RequestOpts{Params: orgParams(orgID)})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "email",
				Usage:     "Email a purchase order",
				ArgsUsage: "<purchaseorder-id>",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "to_mail_ids", Required: true, Usage: "Recipient email addresses (comma-separated)"},
					&cli.StringFlag{Name: "cc_mail_ids", Usage: "CC email addresses (comma-separated)"},
					&cli.StringFlag{Name: "subject", Required: true, Usage: "Email subject"},
					&cli.StringFlag{Name: "body", Required: true, Usage: "Email body"},
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
					body["to_mail_ids"] = cmd.String("to_mail_ids")
					if cmd.IsSet("cc_mail_ids") {
						body["cc_mail_ids"] = cmd.String("cc_mail_ids")
					}
					body["subject"] = cmd.String("subject")
					body["body"] = cmd.String("body")
					if err := internal.MergeJSON(cmd, body); err != nil {
						return err
					}
					raw, err := c.Request("POST", c.BooksBase+"/purchaseorders/"+cmd.Args().First()+"/email", &zohttp.RequestOpts{
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
				Usage:     "Get email content of a purchase order",
				ArgsUsage: "<purchaseorder-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("GET", c.BooksBase+"/purchaseorders/"+cmd.Args().First()+"/email", &zohttp.RequestOpts{Params: orgParams(orgID)})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "update-billing-address",
				Usage:     "Update billing address of a purchase order",
				ArgsUsage: "<purchaseorder-id>",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "attention", Usage: "Attention"},
					&cli.StringFlag{Name: "address", Usage: "Street address"},
					&cli.StringFlag{Name: "street2", Usage: "Street 2"},
					&cli.StringFlag{Name: "city", Usage: "City"},
					&cli.StringFlag{Name: "state", Usage: "State"},
					&cli.StringFlag{Name: "zip", Usage: "ZIP code"},
					&cli.StringFlag{Name: "country", Usage: "Country"},
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
					if cmd.IsSet("attention") {
						body["attention"] = cmd.String("attention")
					}
					if cmd.IsSet("address") {
						body["address"] = cmd.String("address")
					}
					if cmd.IsSet("street2") {
						body["street2"] = cmd.String("street2")
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
					if err := internal.MergeJSON(cmd, body); err != nil {
						return err
					}
					raw, err := c.Request("PUT", c.BooksBase+"/purchaseorders/"+cmd.Args().First()+"/address/billing", &zohttp.RequestOpts{
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
				Usage: "List purchase order templates",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("GET", c.BooksBase+"/purchaseorders"+"/templates", &zohttp.RequestOpts{Params: orgParams(orgID)})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "update-template",
				Usage:     "Update template of a purchase order",
				ArgsUsage: "<purchaseorder-id> <template-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("PUT", c.BooksBase+"/purchaseorders/"+cmd.Args().First()+"/templates/"+cmd.Args().Get(1), &zohttp.RequestOpts{Params: orgParams(orgID)})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "add-attachment",
				Usage:     "Add attachment to a purchase order",
				ArgsUsage: "<purchaseorder-id>",
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
					raw, err := c.Request("POST", c.BooksBase+"/purchaseorders/"+cmd.Args().First()+"/attachment", &zohttp.RequestOpts{
						Params: orgParams(orgID),
						Files:  map[string]zohttp.FileUpload{"attachment": {Filename: filepath.Base(cmd.String("file")), Data: data}},
					})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "update-attachment-preference",
				Usage:     "Update attachment preference of a purchase order",
				ArgsUsage: "<purchaseorder-id>",
				Flags: []cli.Flag{
					&cli.BoolFlag{Name: "can_send_in_mail", Usage: "Include attachment in email (true/false)"},
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
					if cmd.IsSet("can_send_in_mail") {
						body["can_send_in_mail"] = cmd.Bool("can_send_in_mail")
					}
					if err := internal.MergeJSON(cmd, body); err != nil {
						return err
					}
					raw, err := c.Request("PUT", c.BooksBase+"/purchaseorders/"+cmd.Args().First()+"/attachment", &zohttp.RequestOpts{
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
				Name:      "get-attachment",
				Usage:     "Get attachment of a purchase order",
				ArgsUsage: "<purchaseorder-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("GET", c.BooksBase+"/purchaseorders/"+cmd.Args().First()+"/attachment", &zohttp.RequestOpts{Params: orgParams(orgID)})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "delete-attachment",
				Usage:     "Delete attachment of a purchase order",
				ArgsUsage: "<purchaseorder-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("DELETE", c.BooksBase+"/purchaseorders/"+cmd.Args().First()+"/attachment", &zohttp.RequestOpts{Params: orgParams(orgID)})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "add-comment",
				Usage:     "Add a comment to a purchase order",
				ArgsUsage: "<purchaseorder-id>",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "description", Required: true, Usage: "Description"},
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
					body["description"] = cmd.String("description")
					if err := internal.MergeJSON(cmd, body); err != nil {
						return err
					}
					raw, err := c.Request("POST", c.BooksBase+"/purchaseorders/"+cmd.Args().First()+"/comments", &zohttp.RequestOpts{
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
				Usage:     "List comments of a purchase order",
				ArgsUsage: "<purchaseorder-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("GET", c.BooksBase+"/purchaseorders/"+cmd.Args().First()+"/comments", &zohttp.RequestOpts{Params: orgParams(orgID)})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "update-comment",
				Usage:     "Update a comment on a purchase order",
				ArgsUsage: "<purchaseorder-id> <comment-id>",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "description", Required: true, Usage: "Comment text"},
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
					body["description"] = cmd.String("description")
					if err := internal.MergeJSON(cmd, body); err != nil {
						return err
					}
					raw, err := c.Request("PUT", c.BooksBase+"/purchaseorders/"+cmd.Args().First()+"/comments/"+cmd.Args().Get(1), &zohttp.RequestOpts{
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
				Usage:     "Delete a comment on a purchase order",
				ArgsUsage: "<purchaseorder-id> <comment-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("DELETE", c.BooksBase+"/purchaseorders/"+cmd.Args().First()+"/comments/"+cmd.Args().Get(1), &zohttp.RequestOpts{Params: orgParams(orgID)})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "reject",
				Usage:     "Reject a purchase order",
				ArgsUsage: "<purchaseorder-id>",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "reason", Required: true, Usage: "Reason"},
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
					body["reason"] = cmd.String("reason")
					if err := internal.MergeJSON(cmd, body); err != nil {
						return err
					}
					opts := &zohttp.RequestOpts{Params: orgParams(orgID), JSON: body}
					raw, err := c.Request("POST", c.BooksBase+"/purchaseorders/"+cmd.Args().First()+"/status/rejected", opts)
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
		},
	}
}

func billsCmd() *cli.Command {
	return &cli.Command{
		Name:  "bills",
		Usage: "Bill operations",
		Commands: []*cli.Command{
			{
				Name:  "create",
				Usage: "Create a bill",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "vendor_id", Required: true, Usage: "Vendor ID"},
					&cli.StringFlag{Name: "date", Required: true, Usage: "Bill date (YYYY-MM-DD)"},
					&cli.StringFlag{Name: "due_date", Usage: "Due date (YYYY-MM-DD)"},
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
					body["vendor_id"] = cmd.String("vendor_id")
					body["date"] = cmd.String("date")
					if cmd.IsSet("due_date") {
						body["due_date"] = cmd.String("due_date")
					}
					if cmd.IsSet("reference_number") {
						body["reference_number"] = cmd.String("reference_number")
					}
					if err := internal.MergeJSON(cmd, body); err != nil {
						return err
					}
					raw, err := c.Request("POST", c.BooksBase+"/bills", &zohttp.RequestOpts{
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
				Usage: "Update a bill by custom field",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "vendor_id", Usage: "Vendor ID"},
					&cli.StringFlag{Name: "date", Usage: "Bill date (YYYY-MM-DD)"},
					&cli.StringFlag{Name: "due_date", Usage: "Due date (YYYY-MM-DD)"},
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
					if cmd.IsSet("vendor_id") {
						body["vendor_id"] = cmd.String("vendor_id")
					}
					if cmd.IsSet("date") {
						body["date"] = cmd.String("date")
					}
					if cmd.IsSet("due_date") {
						body["due_date"] = cmd.String("due_date")
					}
					if cmd.IsSet("reference_number") {
						body["reference_number"] = cmd.String("reference_number")
					}
					if err := internal.MergeJSON(cmd, body); err != nil {
						return err
					}
					raw, err := c.Request("PUT", c.BooksBase+"/bills", &zohttp.RequestOpts{
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
				Usage: "List bills",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "status", Usage: "Filter"},
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
					if v := cmd.String("vendor-id"); v != "" {
						params["vendor_id"] = v
					}

					if cmd.Bool("all") || cmd.IsSet("limit") {

						items, err := pagination.Paginate(pagination.PaginationConfig{

							Client: c,

							URL: c.BooksBase + "/bills",

							Opts: &zohttp.RequestOpts{Params: params},

							ItemsKey: "bills",

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
					raw, err := c.Request("GET", c.BooksBase+"/bills", &zohttp.RequestOpts{Params: params})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "update",
				Usage:     "Update a bill",
				ArgsUsage: "<bill-id>",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "vendor_id", Usage: "Vendor ID"},
					&cli.StringFlag{Name: "date", Usage: "Bill date (YYYY-MM-DD)"},
					&cli.StringFlag{Name: "due_date", Usage: "Due date (YYYY-MM-DD)"},
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
					if cmd.IsSet("vendor_id") {
						body["vendor_id"] = cmd.String("vendor_id")
					}
					if cmd.IsSet("date") {
						body["date"] = cmd.String("date")
					}
					if cmd.IsSet("due_date") {
						body["due_date"] = cmd.String("due_date")
					}
					if cmd.IsSet("reference_number") {
						body["reference_number"] = cmd.String("reference_number")
					}
					if err := internal.MergeJSON(cmd, body); err != nil {
						return err
					}
					raw, err := c.Request("PUT", c.BooksBase+"/bills/"+cmd.Args().First(), &zohttp.RequestOpts{
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
				Usage:     "Get a bill",
				ArgsUsage: "<bill-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("GET", c.BooksBase+"/bills/"+cmd.Args().First(), &zohttp.RequestOpts{Params: orgParams(orgID)})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "delete",
				Usage:     "Delete a bill",
				ArgsUsage: "<bill-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("DELETE", c.BooksBase+"/bills/"+cmd.Args().First(), &zohttp.RequestOpts{Params: orgParams(orgID)})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "update-custom-fields",
				Usage:     "Update custom fields of a bill",
				ArgsUsage: "<bill-id>",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "custom_fields", Usage: "Custom fields JSON"},
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
					raw, err := c.Request("PUT", c.BooksBase+"/bills/"+cmd.Args().First()+"/customfields", &zohttp.RequestOpts{
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
				Name:      "void",
				Usage:     "Void a bill",
				ArgsUsage: "<bill-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("POST", c.BooksBase+"/bills/"+cmd.Args().First()+"/status/void", &zohttp.RequestOpts{Params: orgParams(orgID)})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "mark-open",
				Usage:     "Mark a bill as open",
				ArgsUsage: "<bill-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("POST", c.BooksBase+"/bills/"+cmd.Args().First()+"/status/open", &zohttp.RequestOpts{Params: orgParams(orgID)})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "submit-for-approval",
				Usage:     "Submit a bill for approval",
				ArgsUsage: "<bill-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("POST", c.BooksBase+"/bills/"+cmd.Args().First()+"/submit", &zohttp.RequestOpts{Params: orgParams(orgID)})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "approve",
				Usage:     "Approve a bill",
				ArgsUsage: "<bill-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("POST", c.BooksBase+"/bills/"+cmd.Args().First()+"/approve", &zohttp.RequestOpts{Params: orgParams(orgID)})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "update-billing-address",
				Usage:     "Update billing address of a bill",
				ArgsUsage: "<bill-id>",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "attention", Usage: "Attention"},
					&cli.StringFlag{Name: "address", Usage: "Street address"},
					&cli.StringFlag{Name: "street2", Usage: "Street 2"},
					&cli.StringFlag{Name: "city", Usage: "City"},
					&cli.StringFlag{Name: "state", Usage: "State"},
					&cli.StringFlag{Name: "zip", Usage: "ZIP code"},
					&cli.StringFlag{Name: "country", Usage: "Country"},
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
					if cmd.IsSet("attention") {
						body["attention"] = cmd.String("attention")
					}
					if cmd.IsSet("address") {
						body["address"] = cmd.String("address")
					}
					if cmd.IsSet("street2") {
						body["street2"] = cmd.String("street2")
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
					if err := internal.MergeJSON(cmd, body); err != nil {
						return err
					}
					raw, err := c.Request("PUT", c.BooksBase+"/bills/"+cmd.Args().First()+"/address/billing", &zohttp.RequestOpts{
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
				Name:      "list-payments",
				Usage:     "List payments of a bill",
				ArgsUsage: "<bill-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("GET", c.BooksBase+"/bills/"+cmd.Args().First()+"/payments", &zohttp.RequestOpts{Params: orgParams(orgID)})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "apply-credits",
				Usage:     "Apply credits to a bill",
				ArgsUsage: "<bill-id>",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "bill_id", Usage: "Bill ID"},
					&cli.FloatFlag{Name: "amount_applied", Usage: "Amount to apply"},
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
					if cmd.IsSet("bill_id") {
						body["bill_id"] = cmd.String("bill_id")
					}
					if cmd.IsSet("amount_applied") {
						body["amount_applied"] = cmd.Float("amount_applied")
					}
					if err := internal.MergeJSON(cmd, body); err != nil {
						return err
					}
					raw, err := c.Request("POST", c.BooksBase+"/bills/"+cmd.Args().First()+"/credits", &zohttp.RequestOpts{
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
				Name:      "delete-payment",
				Usage:     "Delete a payment from a bill",
				ArgsUsage: "<bill-id> <payment-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("DELETE", c.BooksBase+"/bills/"+cmd.Args().First()+"/payments/"+cmd.Args().Get(1), &zohttp.RequestOpts{Params: orgParams(orgID)})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "add-attachment",
				Usage:     "Add attachment to a bill",
				ArgsUsage: "<bill-id>",
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
					raw, err := c.Request("POST", c.BooksBase+"/bills/"+cmd.Args().First()+"/attachment", &zohttp.RequestOpts{
						Params: orgParams(orgID),
						Files:  map[string]zohttp.FileUpload{"attachment": {Filename: filepath.Base(cmd.String("file")), Data: data}},
					})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "get-attachment",
				Usage:     "Get attachment of a bill",
				ArgsUsage: "<bill-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("GET", c.BooksBase+"/bills/"+cmd.Args().First()+"/attachment", &zohttp.RequestOpts{Params: orgParams(orgID)})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "delete-attachment",
				Usage:     "Delete attachment of a bill",
				ArgsUsage: "<bill-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("DELETE", c.BooksBase+"/bills/"+cmd.Args().First()+"/attachment", &zohttp.RequestOpts{Params: orgParams(orgID)})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "add-comment",
				Usage:     "Add a comment to a bill",
				ArgsUsage: "<bill-id>",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "description", Required: true, Usage: "Description"},
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
					body["description"] = cmd.String("description")
					if err := internal.MergeJSON(cmd, body); err != nil {
						return err
					}
					raw, err := c.Request("POST", c.BooksBase+"/bills/"+cmd.Args().First()+"/comments", &zohttp.RequestOpts{
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
				Usage:     "List comments of a bill",
				ArgsUsage: "<bill-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("GET", c.BooksBase+"/bills/"+cmd.Args().First()+"/comments", &zohttp.RequestOpts{Params: orgParams(orgID)})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "delete-comment",
				Usage:     "Delete a comment on a bill",
				ArgsUsage: "<bill-id> <comment-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("DELETE", c.BooksBase+"/bills/"+cmd.Args().First()+"/comments/"+cmd.Args().Get(1), &zohttp.RequestOpts{Params: orgParams(orgID)})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
		},
	}
}

func recurringBillsCmd() *cli.Command {
	return &cli.Command{
		Name:  "recurring-bills",
		Usage: "Recurring bill operations",
		Commands: []*cli.Command{
			{
				Name:  "create",
				Usage: "Create a recurring bill",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "vendor_id", Required: true, Usage: "Vendor ID"},
					&cli.StringFlag{Name: "date", Usage: "Start date (YYYY-MM-DD)"},
					&cli.IntFlag{Name: "repeat_every", Usage: "Repeat interval"},
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
					body["vendor_id"] = cmd.String("vendor_id")
					if cmd.IsSet("date") {
						body["date"] = cmd.String("date")
					}
					if cmd.IsSet("repeat_every") {
						body["repeat_every"] = cmd.Int("repeat_every")
					}
					if cmd.IsSet("reference_number") {
						body["reference_number"] = cmd.String("reference_number")
					}
					if err := internal.MergeJSON(cmd, body); err != nil {
						return err
					}
					raw, err := c.Request("POST", c.BooksBase+"/recurringbills", &zohttp.RequestOpts{
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
				Usage: "Update a recurring bill by custom field",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "vendor_id", Usage: "Vendor ID"},
					&cli.StringFlag{Name: "date", Usage: "Start date (YYYY-MM-DD)"},
					&cli.IntFlag{Name: "repeat_every", Usage: "Repeat interval"},
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
					if cmd.IsSet("vendor_id") {
						body["vendor_id"] = cmd.String("vendor_id")
					}
					if cmd.IsSet("date") {
						body["date"] = cmd.String("date")
					}
					if cmd.IsSet("repeat_every") {
						body["repeat_every"] = cmd.Int("repeat_every")
					}
					if cmd.IsSet("reference_number") {
						body["reference_number"] = cmd.String("reference_number")
					}
					if err := internal.MergeJSON(cmd, body); err != nil {
						return err
					}
					raw, err := c.Request("PUT", c.BooksBase+"/recurringbills", &zohttp.RequestOpts{
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
				Usage: "List recurring bills",
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

							URL: c.BooksBase + "/recurringbills",

							Opts: &zohttp.RequestOpts{Params: params},

							ItemsKey: "recurring_bills",

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
					raw, err := c.Request("GET", c.BooksBase+"/recurringbills", &zohttp.RequestOpts{Params: params})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "update",
				Usage:     "Update a recurring bill",
				ArgsUsage: "<recurringbill-id>",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "vendor_id", Usage: "Vendor ID"},
					&cli.StringFlag{Name: "date", Usage: "Start date (YYYY-MM-DD)"},
					&cli.IntFlag{Name: "repeat_every", Usage: "Repeat interval"},
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
					if cmd.IsSet("vendor_id") {
						body["vendor_id"] = cmd.String("vendor_id")
					}
					if cmd.IsSet("date") {
						body["date"] = cmd.String("date")
					}
					if cmd.IsSet("repeat_every") {
						body["repeat_every"] = cmd.Int("repeat_every")
					}
					if cmd.IsSet("reference_number") {
						body["reference_number"] = cmd.String("reference_number")
					}
					if err := internal.MergeJSON(cmd, body); err != nil {
						return err
					}
					raw, err := c.Request("PUT", c.BooksBase+"/recurringbills/"+cmd.Args().First(), &zohttp.RequestOpts{
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
				Usage:     "Get a recurring bill",
				ArgsUsage: "<recurringbill-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("GET", c.BooksBase+"/recurringbills/"+cmd.Args().First(), &zohttp.RequestOpts{Params: orgParams(orgID)})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "delete",
				Usage:     "Delete a recurring bill",
				ArgsUsage: "<recurringbill-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("DELETE", c.BooksBase+"/recurringbills/"+cmd.Args().First(), &zohttp.RequestOpts{Params: orgParams(orgID)})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "stop",
				Usage:     "Stop a recurring bill",
				ArgsUsage: "<recurringbill-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("POST", c.BooksBase+"/recurringbills/"+cmd.Args().First()+"/status/stop", &zohttp.RequestOpts{Params: orgParams(orgID)})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "resume",
				Usage:     "Resume a recurring bill",
				ArgsUsage: "<recurringbill-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("POST", c.BooksBase+"/recurringbills/"+cmd.Args().First()+"/status/resume", &zohttp.RequestOpts{Params: orgParams(orgID)})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "list-history",
				Usage:     "List history of a recurring bill",
				ArgsUsage: "<recurringbill-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("GET", c.BooksBase+"/recurringbills/"+cmd.Args().First()+"/comments", &zohttp.RequestOpts{Params: orgParams(orgID)})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
		},
	}
}

func vendorCreditsCmd() *cli.Command {
	return &cli.Command{
		Name:  "vendor-credits",
		Usage: "Vendor credit operations",
		Commands: []*cli.Command{
			{
				Name:  "create",
				Usage: "Create a vendor credit",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "vendor_id", Required: true, Usage: "Vendor ID"},
					&cli.StringFlag{Name: "date", Required: true, Usage: "Vendor credit date (YYYY-MM-DD)"},
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
					body["vendor_id"] = cmd.String("vendor_id")
					body["date"] = cmd.String("date")
					if cmd.IsSet("reference_number") {
						body["reference_number"] = cmd.String("reference_number")
					}
					if err := internal.MergeJSON(cmd, body); err != nil {
						return err
					}
					raw, err := c.Request("POST", c.BooksBase+"/vendorcredits", &zohttp.RequestOpts{
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
				Usage: "List vendor credits",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "status", Usage: "Filter"},
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
					if v := cmd.String("vendor-id"); v != "" {
						params["vendor_id"] = v
					}

					if cmd.Bool("all") || cmd.IsSet("limit") {

						items, err := pagination.Paginate(pagination.PaginationConfig{

							Client: c,

							URL: c.BooksBase + "/vendorcredits",

							Opts: &zohttp.RequestOpts{Params: params},

							ItemsKey: "vendorcredits",

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
					raw, err := c.Request("GET", c.BooksBase+"/vendorcredits", &zohttp.RequestOpts{Params: params})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "update",
				Usage:     "Update a vendor credit",
				ArgsUsage: "<vendorcredit-id>",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "vendor_id", Usage: "Vendor ID"},
					&cli.StringFlag{Name: "date", Usage: "Vendor credit date (YYYY-MM-DD)"},
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
					if cmd.IsSet("vendor_id") {
						body["vendor_id"] = cmd.String("vendor_id")
					}
					if cmd.IsSet("date") {
						body["date"] = cmd.String("date")
					}
					if cmd.IsSet("reference_number") {
						body["reference_number"] = cmd.String("reference_number")
					}
					if err := internal.MergeJSON(cmd, body); err != nil {
						return err
					}
					raw, err := c.Request("PUT", c.BooksBase+"/vendorcredits/"+cmd.Args().First(), &zohttp.RequestOpts{
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
				Usage:     "Get a vendor credit",
				ArgsUsage: "<vendorcredit-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("GET", c.BooksBase+"/vendorcredits/"+cmd.Args().First(), &zohttp.RequestOpts{Params: orgParams(orgID)})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "delete",
				Usage:     "Delete a vendor credit",
				ArgsUsage: "<vendorcredit-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("DELETE", c.BooksBase+"/vendorcredits/"+cmd.Args().First(), &zohttp.RequestOpts{Params: orgParams(orgID)})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "convert-to-open",
				Usage:     "Convert a vendor credit to open",
				ArgsUsage: "<vendorcredit-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("POST", c.BooksBase+"/vendorcredits/"+cmd.Args().First()+"/status/open", &zohttp.RequestOpts{Params: orgParams(orgID)})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "void",
				Usage:     "Void a vendor credit",
				ArgsUsage: "<vendorcredit-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("POST", c.BooksBase+"/vendorcredits/"+cmd.Args().First()+"/status/void", &zohttp.RequestOpts{Params: orgParams(orgID)})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "submit-for-approval",
				Usage:     "Submit a vendor credit for approval",
				ArgsUsage: "<vendorcredit-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("POST", c.BooksBase+"/vendorcredits/"+cmd.Args().First()+"/submit", &zohttp.RequestOpts{Params: orgParams(orgID)})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "approve",
				Usage:     "Approve a vendor credit",
				ArgsUsage: "<vendorcredit-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("POST", c.BooksBase+"/vendorcredits/"+cmd.Args().First()+"/approve", &zohttp.RequestOpts{Params: orgParams(orgID)})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "apply-credits-to-bill",
				Usage:     "Apply vendor credit to bills",
				ArgsUsage: "<vendorcredit-id>",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "bill_id", Usage: "Bill ID"},
					&cli.FloatFlag{Name: "amount_applied", Usage: "Amount to apply"},
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
					if cmd.IsSet("bill_id") {
						body["bill_id"] = cmd.String("bill_id")
					}
					if cmd.IsSet("amount_applied") {
						body["amount_applied"] = cmd.Float("amount_applied")
					}
					if err := internal.MergeJSON(cmd, body); err != nil {
						return err
					}
					raw, err := c.Request("POST", c.BooksBase+"/vendorcredits/"+cmd.Args().First()+"/bills", &zohttp.RequestOpts{
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
				Name:      "list-bills-credited",
				Usage:     "List bills credited",
				ArgsUsage: "<vendorcredit-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("GET", c.BooksBase+"/vendorcredits/"+cmd.Args().First()+"/bills", &zohttp.RequestOpts{Params: orgParams(orgID)})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "delete-bills-credited",
				Usage:     "Delete bills credited",
				ArgsUsage: "<vendorcredit-id> <credited-bill-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("DELETE", c.BooksBase+"/vendorcredits/"+cmd.Args().First()+"/bills/"+cmd.Args().Get(1), &zohttp.RequestOpts{Params: orgParams(orgID)})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "refund",
				Usage:     "Refund a vendor credit",
				ArgsUsage: "<vendorcredit-id>",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "date", Required: true, Usage: "Date (YYYY-MM-DD)"},
					&cli.StringFlag{Name: "refund_mode", Usage: "Refund mode"},
					&cli.StringFlag{Name: "from_account_id", Usage: "From account ID"},
					&cli.FloatFlag{Name: "amount", Required: true, Usage: "Amount"},
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
					raw, err := c.Request("POST", c.BooksBase+"/vendorcredits/"+cmd.Args().First()+"/refunds", &zohttp.RequestOpts{
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
				Usage:     "List refunds of a vendor credit",
				ArgsUsage: "<vendorcredit-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("GET", c.BooksBase+"/vendorcredits/"+cmd.Args().First()+"/refunds", &zohttp.RequestOpts{Params: orgParams(orgID)})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "update-refund",
				Usage:     "Update a refund of a vendor credit",
				ArgsUsage: "<vendorcredit-id> <refund-id>",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "date", Usage: "Refund date (YYYY-MM-DD)"},
					&cli.StringFlag{Name: "refund_mode", Usage: "Refund mode"},
					&cli.StringFlag{Name: "from_account_id", Usage: "Account ID for refund"},
					&cli.FloatFlag{Name: "amount", Usage: "Refund amount"},
					&cli.StringFlag{Name: "reference_number", Usage: "Reference number"},
					&cli.StringFlag{Name: "description", Usage: "Refund description"},
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
					raw, err := c.Request("PUT", c.BooksBase+"/vendorcredits/"+cmd.Args().First()+"/refunds/"+cmd.Args().Get(1), &zohttp.RequestOpts{
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
				Usage:     "Get a refund of a vendor credit",
				ArgsUsage: "<vendorcredit-id> <refund-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("GET", c.BooksBase+"/vendorcredits/"+cmd.Args().First()+"/refunds/"+cmd.Args().Get(1), &zohttp.RequestOpts{Params: orgParams(orgID)})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "delete-refund",
				Usage:     "Delete a refund of a vendor credit",
				ArgsUsage: "<vendorcredit-id> <refund-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("DELETE", c.BooksBase+"/vendorcredits/"+cmd.Args().First()+"/refunds/"+cmd.Args().Get(1), &zohttp.RequestOpts{Params: orgParams(orgID)})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:  "list-all-refunds",
				Usage: "List all vendor credit refunds",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("GET", c.BooksBase+"/vendorcredits"+"/refunds", &zohttp.RequestOpts{Params: orgParams(orgID)})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "add-comment",
				Usage:     "Add a comment to a vendor credit",
				ArgsUsage: "<vendorcredit-id>",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "description", Required: true, Usage: "Description"},
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
					body["description"] = cmd.String("description")
					if err := internal.MergeJSON(cmd, body); err != nil {
						return err
					}
					raw, err := c.Request("POST", c.BooksBase+"/vendorcredits/"+cmd.Args().First()+"/comments", &zohttp.RequestOpts{
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
				Usage:     "List comments of a vendor credit",
				ArgsUsage: "<vendorcredit-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("GET", c.BooksBase+"/vendorcredits/"+cmd.Args().First()+"/comments", &zohttp.RequestOpts{Params: orgParams(orgID)})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "delete-comment",
				Usage:     "Delete a comment on a vendor credit",
				ArgsUsage: "<vendorcredit-id> <comment-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("DELETE", c.BooksBase+"/vendorcredits/"+cmd.Args().First()+"/comments/"+cmd.Args().Get(1), &zohttp.RequestOpts{Params: orgParams(orgID)})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
		},
	}
}
