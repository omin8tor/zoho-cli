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

func contactsCmd() *cli.Command {
	return &cli.Command{
		Name:  "contacts",
		Usage: "Contact operations",
		Commands: []*cli.Command{
			{
				Name:  "create",
				Usage: "Create a contact",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "contact_name", Required: true, Usage: "Contact name"},
					&cli.StringFlag{Name: "contact_type", Required: true, Usage: "Contact type (customer/vendor)"},
					&cli.StringFlag{Name: "company_name", Usage: "Company name"},
					&cli.StringFlag{Name: "currency_id", Usage: "Currency ID"},
					&cli.IntFlag{Name: "payment_terms", Usage: "Payment terms"},
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
					body["contact_name"] = cmd.String("contact_name")
					body["contact_type"] = cmd.String("contact_type")
					if cmd.IsSet("company_name") {
						body["company_name"] = cmd.String("company_name")
					}
					if cmd.IsSet("currency_id") {
						body["currency_id"] = cmd.String("currency_id")
					}
					if cmd.IsSet("payment_terms") {
						body["payment_terms"] = cmd.Int("payment_terms")
					}
					if err := internal.MergeJSON(cmd, body); err != nil {
						return err
					}
					raw, err := c.Request(ctx, "POST", c.BooksBase+"/contacts", &zohttp.RequestOpts{
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
				Usage: "Update a contact by custom field",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "contact_name", Usage: "Contact name"},
					&cli.StringFlag{Name: "company_name", Usage: "Company name"},
					&cli.StringFlag{Name: "contact_type", Usage: "Contact type (customer/vendor)"},
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
					if cmd.IsSet("contact_name") {
						body["contact_name"] = cmd.String("contact_name")
					}
					if cmd.IsSet("company_name") {
						body["company_name"] = cmd.String("company_name")
					}
					if cmd.IsSet("contact_type") {
						body["contact_type"] = cmd.String("contact_type")
					}
					if err := internal.MergeJSON(cmd, body); err != nil {
						return err
					}
					raw, err := c.Request(ctx, "PUT", c.BooksBase+"/contacts", &zohttp.RequestOpts{
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
				Usage: "List contacts",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "status", Usage: "Filter by status"},
					&cli.StringFlag{Name: "sort-column", Usage: "Sort column"},
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
					if v := cmd.String("sort-column"); v != "" {
						params["sort_column"] = v
					}

					if cmd.Bool("all") || cmd.IsSet("limit") {

						items, err := pagination.Paginate(ctx, pagination.PaginationConfig{

							Client: c,

							URL: c.BooksBase + "/contacts",

							Opts: &zohttp.RequestOpts{Params: params},

							ItemsKey: "contacts",

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
					raw, err := c.Request(ctx, "GET", c.BooksBase+"/contacts", &zohttp.RequestOpts{Params: params})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "update",
				Usage:     "Update a contact",
				ArgsUsage: "<contact-id>",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "contact_name", Usage: "Contact name"},
					&cli.StringFlag{Name: "contact_type", Usage: "Contact type (customer/vendor)"},
					&cli.StringFlag{Name: "company_name", Usage: "Company name"},
					&cli.StringFlag{Name: "currency_id", Usage: "Currency ID"},
					&cli.IntFlag{Name: "payment_terms", Usage: "Payment terms"},
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
					if cmd.IsSet("contact_name") {
						body["contact_name"] = cmd.String("contact_name")
					}
					if cmd.IsSet("contact_type") {
						body["contact_type"] = cmd.String("contact_type")
					}
					if cmd.IsSet("company_name") {
						body["company_name"] = cmd.String("company_name")
					}
					if cmd.IsSet("currency_id") {
						body["currency_id"] = cmd.String("currency_id")
					}
					if cmd.IsSet("payment_terms") {
						body["payment_terms"] = cmd.Int("payment_terms")
					}
					if err := internal.MergeJSON(cmd, body); err != nil {
						return err
					}
					raw, err := c.Request(ctx, "PUT", c.BooksBase+"/contacts/"+cmd.Args().First(), &zohttp.RequestOpts{
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
				Usage:     "Get a contact",
				ArgsUsage: "<contact-id>",
				Action: func(ctx context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request(ctx, "GET", c.BooksBase+"/contacts/"+cmd.Args().First(), &zohttp.RequestOpts{Params: orgParams(orgID)})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "delete",
				Usage:     "Delete a contact",
				ArgsUsage: "<contact-id>",
				Action: func(ctx context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request(ctx, "DELETE", c.BooksBase+"/contacts/"+cmd.Args().First(), &zohttp.RequestOpts{Params: orgParams(orgID)})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "mark-active",
				Usage:     "Mark a contact as active",
				ArgsUsage: "<contact-id>",
				Action: func(ctx context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request(ctx, "POST", c.BooksBase+"/contacts/"+cmd.Args().First()+"/active", &zohttp.RequestOpts{Params: orgParams(orgID)})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "mark-inactive",
				Usage:     "Mark a contact as inactive",
				ArgsUsage: "<contact-id>",
				Action: func(ctx context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request(ctx, "POST", c.BooksBase+"/contacts/"+cmd.Args().First()+"/inactive", &zohttp.RequestOpts{Params: orgParams(orgID)})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "enable-portal",
				Usage:     "Enable portal for a contact",
				ArgsUsage: "<contact-id>",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "contact_persons", Usage: "Contact person IDs (comma-separated or JSON array)"},
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
					if cmd.IsSet("contact_persons") {
						var parsed any
						if err := json.Unmarshal([]byte(cmd.String("contact_persons")), &parsed); err != nil {
							return err
						}
						body["contact_persons"] = parsed
					}
					if err := internal.MergeJSON(cmd, body); err != nil {
						return err
					}
					raw, err := c.Request(ctx, "POST", c.BooksBase+"/contacts/"+cmd.Args().First()+"/portal/enable", &zohttp.RequestOpts{
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
				Name:      "enable-payment-reminders",
				Usage:     "Enable payment reminders for a contact",
				ArgsUsage: "<contact-id>",
				Action: func(ctx context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request(ctx, "POST", c.BooksBase+"/contacts/"+cmd.Args().First()+"/paymentreminder/enable", &zohttp.RequestOpts{Params: orgParams(orgID)})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "disable-payment-reminders",
				Usage:     "Disable payment reminders for a contact",
				ArgsUsage: "<contact-id>",
				Action: func(ctx context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request(ctx, "POST", c.BooksBase+"/contacts/"+cmd.Args().First()+"/paymentreminder/disable", &zohttp.RequestOpts{Params: orgParams(orgID)})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "email-statement",
				Usage:     "Email statement to a contact",
				ArgsUsage: "<contact-id>",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "to_mail_ids", Required: true, Usage: "Recipient email addresses (comma-separated)"},
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
					body["subject"] = cmd.String("subject")
					body["body"] = cmd.String("body")
					if err := internal.MergeJSON(cmd, body); err != nil {
						return err
					}
					raw, err := c.Request(ctx, "POST", c.BooksBase+"/contacts/"+cmd.Args().First()+"/statements/email", &zohttp.RequestOpts{
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
				Name:      "get-statement-mail-content",
				Usage:     "Get statement email content for a contact",
				ArgsUsage: "<contact-id>",
				Action: func(ctx context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request(ctx, "GET", c.BooksBase+"/contacts/"+cmd.Args().First()+"/statements/email", &zohttp.RequestOpts{Params: orgParams(orgID)})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "email",
				Usage:     "Email a contact",
				ArgsUsage: "<contact-id>",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "to_mail_ids", Required: true, Usage: "Recipient email addresses (comma-separated)"},
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
					body["subject"] = cmd.String("subject")
					body["body"] = cmd.String("body")
					if err := internal.MergeJSON(cmd, body); err != nil {
						return err
					}
					raw, err := c.Request(ctx, "POST", c.BooksBase+"/contacts/"+cmd.Args().First()+"/email", &zohttp.RequestOpts{
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
				Usage:     "List comments for a contact",
				ArgsUsage: "<contact-id>",
				Action: func(ctx context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request(ctx, "GET", c.BooksBase+"/contacts/"+cmd.Args().First()+"/comments", &zohttp.RequestOpts{Params: orgParams(orgID)})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "add-address",
				Usage:     "Add address to a contact",
				ArgsUsage: "<contact-id>",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "attention", Usage: "Attention"},
					&cli.StringFlag{Name: "address", Usage: "Street address"},
					&cli.StringFlag{Name: "street2", Usage: "Street address line 2"},
					&cli.StringFlag{Name: "city", Usage: "City"},
					&cli.StringFlag{Name: "state", Usage: "State"},
					&cli.StringFlag{Name: "zip", Usage: "Zip/postal code"},
					&cli.StringFlag{Name: "country", Usage: "Country"},
					&cli.StringFlag{Name: "fax", Usage: "Fax number"},
					&cli.StringFlag{Name: "phone", Usage: "Phone number"},
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
					if cmd.IsSet("fax") {
						body["fax"] = cmd.String("fax")
					}
					if cmd.IsSet("phone") {
						body["phone"] = cmd.String("phone")
					}
					if err := internal.MergeJSON(cmd, body); err != nil {
						return err
					}
					raw, err := c.Request(ctx, "POST", c.BooksBase+"/contacts/"+cmd.Args().First()+"/address", &zohttp.RequestOpts{
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
				Name:      "get-addresses",
				Usage:     "Get addresses for a contact",
				ArgsUsage: "<contact-id>",
				Action: func(ctx context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request(ctx, "GET", c.BooksBase+"/contacts/"+cmd.Args().First()+"/address", &zohttp.RequestOpts{Params: orgParams(orgID)})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "edit-address",
				Usage:     "Edit an address for a contact",
				ArgsUsage: "<contact-id> <address-id>",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "attention", Usage: "Attention"},
					&cli.StringFlag{Name: "address", Usage: "Street address"},
					&cli.StringFlag{Name: "street2", Usage: "Street address line 2"},
					&cli.StringFlag{Name: "city", Usage: "City"},
					&cli.StringFlag{Name: "state", Usage: "State"},
					&cli.StringFlag{Name: "zip", Usage: "Zip/postal code"},
					&cli.StringFlag{Name: "country", Usage: "Country"},
					&cli.StringFlag{Name: "fax", Usage: "Fax number"},
					&cli.StringFlag{Name: "phone", Usage: "Phone number"},
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
					if cmd.IsSet("fax") {
						body["fax"] = cmd.String("fax")
					}
					if cmd.IsSet("phone") {
						body["phone"] = cmd.String("phone")
					}
					if err := internal.MergeJSON(cmd, body); err != nil {
						return err
					}
					raw, err := c.Request(ctx, "PUT", c.BooksBase+"/contacts/"+cmd.Args().First()+"/address/"+cmd.Args().Get(1), &zohttp.RequestOpts{
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
				Name:      "delete-address",
				Usage:     "Delete an address for a contact",
				ArgsUsage: "<contact-id> <address-id>",
				Action: func(ctx context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request(ctx, "DELETE", c.BooksBase+"/contacts/"+cmd.Args().First()+"/address/"+cmd.Args().Get(1), &zohttp.RequestOpts{Params: orgParams(orgID)})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "list-refunds",
				Usage:     "List refunds for a contact",
				ArgsUsage: "<contact-id>",
				Action: func(ctx context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request(ctx, "GET", c.BooksBase+"/contacts/"+cmd.Args().First()+"/refunds", &zohttp.RequestOpts{Params: orgParams(orgID)})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "track-1099",
				Usage:     "Track 1099 for a contact",
				ArgsUsage: "<contact-id>",
				Action: func(ctx context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request(ctx, "POST", c.BooksBase+"/contacts/"+cmd.Args().First()+"/track1099", &zohttp.RequestOpts{Params: orgParams(orgID)})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "untrack-1099",
				Usage:     "Untrack 1099 for a contact",
				ArgsUsage: "<contact-id>",
				Action: func(ctx context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request(ctx, "POST", c.BooksBase+"/contacts/"+cmd.Args().First()+"/untrack1099", &zohttp.RequestOpts{Params: orgParams(orgID)})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "get-unused-retainer-payments",
				Usage:     "Get unused retainer payments for a contact",
				ArgsUsage: "<contact-id>",
				Action: func(ctx context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request(ctx, "GET", c.BooksBase+"/contacts/"+cmd.Args().First()+"/retainerpayments", &zohttp.RequestOpts{Params: orgParams(orgID)})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
		},
	}
}

func contactPersonsCmd() *cli.Command {
	return &cli.Command{
		Name:  "contact-persons",
		Usage: "Contact person operations",
		Commands: []*cli.Command{
			{
				Name:  "create",
				Usage: "Create a contact person",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "contact-id", Required: true, Usage: "Parent contact ID"},
					&cli.StringFlag{Name: "first_name", Required: true, Usage: "First name"},
					&cli.StringFlag{Name: "last_name", Required: true, Usage: "Last name"},
					&cli.StringFlag{Name: "email", Required: true, Usage: "Email address"},
					&cli.StringFlag{Name: "salutation", Usage: "Salutation (Mr, Mrs, Ms, etc.)"},
					&cli.StringFlag{Name: "phone", Usage: "Phone number"},
					&cli.StringFlag{Name: "mobile", Usage: "Mobile number"},
					&cli.StringFlag{Name: "designation", Usage: "Designation"},
					&cli.StringFlag{Name: "department", Usage: "Department"},
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
					body["first_name"] = cmd.String("first_name")
					body["last_name"] = cmd.String("last_name")
					body["email"] = cmd.String("email")
					if cmd.IsSet("salutation") {
						body["salutation"] = cmd.String("salutation")
					}
					if cmd.IsSet("phone") {
						body["phone"] = cmd.String("phone")
					}
					if cmd.IsSet("mobile") {
						body["mobile"] = cmd.String("mobile")
					}
					if cmd.IsSet("designation") {
						body["designation"] = cmd.String("designation")
					}
					if cmd.IsSet("department") {
						body["department"] = cmd.String("department")
					}
					if err := internal.MergeJSON(cmd, body); err != nil {
						return err
					}
					raw, err := c.Request(ctx, "POST", c.BooksBase+"/contacts/"+cmd.String("contact-id")+"/contactpersons", &zohttp.RequestOpts{
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
				Name:      "update",
				Usage:     "Update a contact person",
				ArgsUsage: "<person-id>",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "contact-id", Required: true, Usage: "Parent contact ID"},
					&cli.StringFlag{Name: "first_name", Usage: "First name"},
					&cli.StringFlag{Name: "last_name", Usage: "Last name"},
					&cli.StringFlag{Name: "email", Usage: "Email address"},
					&cli.StringFlag{Name: "salutation", Usage: "Salutation (Mr, Mrs, Ms, etc.)"},
					&cli.StringFlag{Name: "phone", Usage: "Phone number"},
					&cli.StringFlag{Name: "mobile", Usage: "Mobile number"},
					&cli.StringFlag{Name: "designation", Usage: "Designation"},
					&cli.StringFlag{Name: "department", Usage: "Department"},
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
					if cmd.IsSet("first_name") {
						body["first_name"] = cmd.String("first_name")
					}
					if cmd.IsSet("last_name") {
						body["last_name"] = cmd.String("last_name")
					}
					if cmd.IsSet("email") {
						body["email"] = cmd.String("email")
					}
					if cmd.IsSet("salutation") {
						body["salutation"] = cmd.String("salutation")
					}
					if cmd.IsSet("phone") {
						body["phone"] = cmd.String("phone")
					}
					if cmd.IsSet("mobile") {
						body["mobile"] = cmd.String("mobile")
					}
					if cmd.IsSet("designation") {
						body["designation"] = cmd.String("designation")
					}
					if cmd.IsSet("department") {
						body["department"] = cmd.String("department")
					}
					if err := internal.MergeJSON(cmd, body); err != nil {
						return err
					}
					raw, err := c.Request(ctx, "PUT", c.BooksBase+"/contacts/"+cmd.String("contact-id")+"/contactpersons/"+cmd.Args().First(), &zohttp.RequestOpts{
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
				Name:      "delete",
				Usage:     "Delete a contact person",
				ArgsUsage: "<person-id>",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "contact-id", Required: true, Usage: "Parent contact ID"},
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
					raw, err := c.Request(ctx, "DELETE", c.BooksBase+"/contacts/"+cmd.String("contact-id")+"/contactpersons/"+cmd.Args().First(), &zohttp.RequestOpts{Params: orgParams(orgID)})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:  "list",
				Usage: "List contact persons",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "contact-id", Required: true, Usage: "Parent contact ID"},
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
					raw, err := c.Request(ctx, "GET", c.BooksBase+"/contacts/"+cmd.String("contact-id")+"/contactpersons", &zohttp.RequestOpts{Params: orgParams(orgID)})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "get",
				Usage:     "Get a contact person",
				ArgsUsage: "<person-id>",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "contact-id", Required: true, Usage: "Parent contact ID"},
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
					raw, err := c.Request(ctx, "GET", c.BooksBase+"/contacts/"+cmd.String("contact-id")+"/contactpersons/"+cmd.Args().First(), &zohttp.RequestOpts{Params: orgParams(orgID)})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "mark-primary",
				Usage:     "Mark a contact person as primary",
				ArgsUsage: "<person-id>",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "contact-id", Required: true, Usage: "Parent contact ID"},
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
					raw, err := c.Request(ctx, "POST", c.BooksBase+"/contacts/"+cmd.String("contact-id")+"/contactpersons/"+cmd.Args().First()+"/primary", &zohttp.RequestOpts{Params: orgParams(orgID)})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
		},
	}
}
