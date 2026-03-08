package invoice

import (
	"context"

	"github.com/omin8tor/zoho-cli/internal"
	zohttp "github.com/omin8tor/zoho-cli/internal/http"
	"github.com/omin8tor/zoho-cli/internal/output"
	"github.com/urfave/cli/v3"
)

func orgParams(orgID string) map[string]string {
	return map[string]string{"organization_id": orgID}
}

func Commands() *cli.Command {
	return &cli.Command{
		Name:  "invoice",
		Usage: "Zoho Invoice operations",
		Flags: []cli.Flag{
			&cli.StringFlag{Name: "org", Sources: cli.EnvVars("ZOHO_BOOKS_ORG_ID"), Usage: "Organization ID (or set ZOHO_BOOKS_ORG_ID)"},
		},
		Commands: []*cli.Command{
			organizationsCmd(),
			contactsCmd(),
			estimatesCmd(),
			invoicesCmd(),
			recurringInvoicesCmd(),
			creditNotesCmd(),
			customerPaymentsCmd(),
			expensesCmd(),
			recurringExpensesCmd(),
			retainerInvoicesCmd(),
			projectsCmd(),
			itemsCmd(),
			currenciesCmd(),
			taxesCmd(),
			usersCmd(),
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
					raw, err := c.Request("GET", c.InvoiceBase+"/organizations", nil)
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
					if cmd.Args().Len() < 1 {
						return internal.NewValidationError("organization ID required")
					}
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					raw, err := c.Request("GET", c.InvoiceBase+"/organizations/"+cmd.Args().First(), nil)
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
		},
	}
}

func contactsCmd() *cli.Command {
	return &cli.Command{
		Name:  "contacts",
		Usage: "Contact operations",
		Commands: []*cli.Command{
			{
				Name:  "list",
				Usage: "List contacts",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := internal.RequireFlag(cmd, "org", "ZOHO_BOOKS_ORG_ID")
					if err != nil {
						return err
					}
					raw, err := c.Request("GET", c.InvoiceBase+"/contacts", &zohttp.RequestOpts{Params: orgParams(orgID)})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:  "create",
				Usage: "Create a contact",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "contact_name", Required: true, Usage: "Name of the contact"},
					&cli.StringFlag{Name: "company_name", Usage: "Company name"},
					&cli.StringFlag{Name: "contact_type", Usage: "Contact type (customer or vendor)"},
					&cli.StringFlag{Name: "website", Usage: "Website URL"},
					&cli.StringFlag{Name: "language_code", Usage: "Language code (e.g. en)"},
					&cli.StringFlag{Name: "notes", Usage: "Notes for the contact"},
					&cli.StringFlag{Name: "json", Usage: "Additional fields as JSON"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := internal.RequireFlag(cmd, "org", "ZOHO_BOOKS_ORG_ID")
					if err != nil {
						return err
					}
					body := map[string]any{}
					body["contact_name"] = cmd.String("contact_name")
					if cmd.IsSet("company_name") {
						body["company_name"] = cmd.String("company_name")
					}
					if cmd.IsSet("contact_type") {
						body["contact_type"] = cmd.String("contact_type")
					}
					if cmd.IsSet("website") {
						body["website"] = cmd.String("website")
					}
					if cmd.IsSet("language_code") {
						body["language_code"] = cmd.String("language_code")
					}
					if cmd.IsSet("notes") {
						body["notes"] = cmd.String("notes")
					}
					if err := internal.MergeJSON(cmd, body); err != nil {
						return err
					}
					raw, err := c.Request("POST", c.InvoiceBase+"/contacts", &zohttp.RequestOpts{
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
				Action: func(_ context.Context, cmd *cli.Command) error {
					if cmd.Args().Len() < 1 {
						return internal.NewValidationError("contact ID required")
					}
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := internal.RequireFlag(cmd, "org", "ZOHO_BOOKS_ORG_ID")
					if err != nil {
						return err
					}
					raw, err := c.Request("GET", c.InvoiceBase+"/contacts/"+cmd.Args().First(), &zohttp.RequestOpts{Params: orgParams(orgID)})
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
					&cli.StringFlag{Name: "contact_name", Usage: "Name of the contact"},
					&cli.StringFlag{Name: "company_name", Usage: "Company name"},
					&cli.StringFlag{Name: "contact_type", Usage: "Contact type (customer or vendor)"},
					&cli.StringFlag{Name: "website", Usage: "Website URL"},
					&cli.StringFlag{Name: "language_code", Usage: "Language code (e.g. en)"},
					&cli.StringFlag{Name: "notes", Usage: "Notes for the contact"},
					&cli.StringFlag{Name: "json", Usage: "Additional fields as JSON"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					if cmd.Args().Len() < 1 {
						return internal.NewValidationError("contact ID required")
					}
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := internal.RequireFlag(cmd, "org", "ZOHO_BOOKS_ORG_ID")
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
					if cmd.IsSet("website") {
						body["website"] = cmd.String("website")
					}
					if cmd.IsSet("language_code") {
						body["language_code"] = cmd.String("language_code")
					}
					if cmd.IsSet("notes") {
						body["notes"] = cmd.String("notes")
					}
					if err := internal.MergeJSON(cmd, body); err != nil {
						return err
					}
					raw, err := c.Request("PUT", c.InvoiceBase+"/contacts/"+cmd.Args().First(), &zohttp.RequestOpts{
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
				Usage:     "Delete a contact",
				ArgsUsage: "<contact-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					if cmd.Args().Len() < 1 {
						return internal.NewValidationError("contact ID required")
					}
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := internal.RequireFlag(cmd, "org", "ZOHO_BOOKS_ORG_ID")
					if err != nil {
						return err
					}
					raw, err := c.Request("DELETE", c.InvoiceBase+"/contacts/"+cmd.Args().First(), &zohttp.RequestOpts{Params: orgParams(orgID)})
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
				Action: func(_ context.Context, cmd *cli.Command) error {
					if cmd.Args().Len() < 1 {
						return internal.NewValidationError("contact ID required")
					}
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := internal.RequireFlag(cmd, "org", "ZOHO_BOOKS_ORG_ID")
					if err != nil {
						return err
					}
					raw, err := c.Request("POST", c.InvoiceBase+"/contacts/"+cmd.Args().First()+"/active", &zohttp.RequestOpts{Params: orgParams(orgID)})
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
				Action: func(_ context.Context, cmd *cli.Command) error {
					if cmd.Args().Len() < 1 {
						return internal.NewValidationError("contact ID required")
					}
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := internal.RequireFlag(cmd, "org", "ZOHO_BOOKS_ORG_ID")
					if err != nil {
						return err
					}
					raw, err := c.Request("POST", c.InvoiceBase+"/contacts/"+cmd.Args().First()+"/inactive", &zohttp.RequestOpts{Params: orgParams(orgID)})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
		},
	}
}

func estimatesCmd() *cli.Command {
	return &cli.Command{
		Name:  "estimates",
		Usage: "Estimate operations",
		Commands: []*cli.Command{
			{
				Name:  "list",
				Usage: "List estimates",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := internal.RequireFlag(cmd, "org", "ZOHO_BOOKS_ORG_ID")
					if err != nil {
						return err
					}
					raw, err := c.Request("GET", c.InvoiceBase+"/estimates", &zohttp.RequestOpts{Params: orgParams(orgID)})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:  "create",
				Usage: "Create an estimate",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "customer_id", Required: true, Usage: "Customer ID"},
					&cli.StringFlag{Name: "estimate_number", Usage: "Estimate number"},
					&cli.StringFlag{Name: "reference_number", Usage: "Reference number"},
					&cli.StringFlag{Name: "date", Usage: "Estimate date (YYYY-MM-DD)"},
					&cli.StringFlag{Name: "expiry_date", Usage: "Expiry date (YYYY-MM-DD)"},
					&cli.StringFlag{Name: "notes", Usage: "Notes"},
					&cli.StringFlag{Name: "terms", Usage: "Terms and conditions"},
					&cli.StringFlag{Name: "json", Usage: "Additional fields as JSON"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := internal.RequireFlag(cmd, "org", "ZOHO_BOOKS_ORG_ID")
					if err != nil {
						return err
					}
					body := map[string]any{}
					body["customer_id"] = cmd.String("customer_id")
					if cmd.IsSet("estimate_number") {
						body["estimate_number"] = cmd.String("estimate_number")
					}
					if cmd.IsSet("reference_number") {
						body["reference_number"] = cmd.String("reference_number")
					}
					if cmd.IsSet("date") {
						body["date"] = cmd.String("date")
					}
					if cmd.IsSet("expiry_date") {
						body["expiry_date"] = cmd.String("expiry_date")
					}
					if cmd.IsSet("notes") {
						body["notes"] = cmd.String("notes")
					}
					if cmd.IsSet("terms") {
						body["terms"] = cmd.String("terms")
					}
					if err := internal.MergeJSON(cmd, body); err != nil {
						return err
					}
					raw, err := c.Request("POST", c.InvoiceBase+"/estimates", &zohttp.RequestOpts{
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
				Action: func(_ context.Context, cmd *cli.Command) error {
					if cmd.Args().Len() < 1 {
						return internal.NewValidationError("estimate ID required")
					}
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := internal.RequireFlag(cmd, "org", "ZOHO_BOOKS_ORG_ID")
					if err != nil {
						return err
					}
					raw, err := c.Request("GET", c.InvoiceBase+"/estimates/"+cmd.Args().First(), &zohttp.RequestOpts{Params: orgParams(orgID)})
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
					&cli.StringFlag{Name: "estimate_number", Usage: "Estimate number"},
					&cli.StringFlag{Name: "reference_number", Usage: "Reference number"},
					&cli.StringFlag{Name: "date", Usage: "Estimate date (YYYY-MM-DD)"},
					&cli.StringFlag{Name: "expiry_date", Usage: "Expiry date (YYYY-MM-DD)"},
					&cli.StringFlag{Name: "notes", Usage: "Notes"},
					&cli.StringFlag{Name: "terms", Usage: "Terms and conditions"},
					&cli.StringFlag{Name: "json", Usage: "Additional fields as JSON"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					if cmd.Args().Len() < 1 {
						return internal.NewValidationError("estimate ID required")
					}
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := internal.RequireFlag(cmd, "org", "ZOHO_BOOKS_ORG_ID")
					if err != nil {
						return err
					}
					body := map[string]any{}
					if cmd.IsSet("customer_id") {
						body["customer_id"] = cmd.String("customer_id")
					}
					if cmd.IsSet("estimate_number") {
						body["estimate_number"] = cmd.String("estimate_number")
					}
					if cmd.IsSet("reference_number") {
						body["reference_number"] = cmd.String("reference_number")
					}
					if cmd.IsSet("date") {
						body["date"] = cmd.String("date")
					}
					if cmd.IsSet("expiry_date") {
						body["expiry_date"] = cmd.String("expiry_date")
					}
					if cmd.IsSet("notes") {
						body["notes"] = cmd.String("notes")
					}
					if cmd.IsSet("terms") {
						body["terms"] = cmd.String("terms")
					}
					if err := internal.MergeJSON(cmd, body); err != nil {
						return err
					}
					raw, err := c.Request("PUT", c.InvoiceBase+"/estimates/"+cmd.Args().First(), &zohttp.RequestOpts{
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
				Usage:     "Delete an estimate",
				ArgsUsage: "<estimate-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					if cmd.Args().Len() < 1 {
						return internal.NewValidationError("estimate ID required")
					}
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := internal.RequireFlag(cmd, "org", "ZOHO_BOOKS_ORG_ID")
					if err != nil {
						return err
					}
					raw, err := c.Request("DELETE", c.InvoiceBase+"/estimates/"+cmd.Args().First(), &zohttp.RequestOpts{Params: orgParams(orgID)})
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
				Action: func(_ context.Context, cmd *cli.Command) error {
					if cmd.Args().Len() < 1 {
						return internal.NewValidationError("estimate ID required")
					}
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := internal.RequireFlag(cmd, "org", "ZOHO_BOOKS_ORG_ID")
					if err != nil {
						return err
					}
					raw, err := c.Request("POST", c.InvoiceBase+"/estimates/"+cmd.Args().First()+"/status/sent", &zohttp.RequestOpts{Params: orgParams(orgID)})
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
				Action: func(_ context.Context, cmd *cli.Command) error {
					if cmd.Args().Len() < 1 {
						return internal.NewValidationError("estimate ID required")
					}
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := internal.RequireFlag(cmd, "org", "ZOHO_BOOKS_ORG_ID")
					if err != nil {
						return err
					}
					raw, err := c.Request("POST", c.InvoiceBase+"/estimates/"+cmd.Args().First()+"/status/accepted", &zohttp.RequestOpts{Params: orgParams(orgID)})
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
				Action: func(_ context.Context, cmd *cli.Command) error {
					if cmd.Args().Len() < 1 {
						return internal.NewValidationError("estimate ID required")
					}
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := internal.RequireFlag(cmd, "org", "ZOHO_BOOKS_ORG_ID")
					if err != nil {
						return err
					}
					raw, err := c.Request("POST", c.InvoiceBase+"/estimates/"+cmd.Args().First()+"/status/declined", &zohttp.RequestOpts{Params: orgParams(orgID)})
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
				Action: func(_ context.Context, cmd *cli.Command) error {
					if cmd.Args().Len() < 1 {
						return internal.NewValidationError("estimate ID required")
					}
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := internal.RequireFlag(cmd, "org", "ZOHO_BOOKS_ORG_ID")
					if err != nil {
						return err
					}
					raw, err := c.Request("GET", c.InvoiceBase+"/estimates/"+cmd.Args().First()+"/comments", &zohttp.RequestOpts{Params: orgParams(orgID)})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
		},
	}
}

func invoicesCmd() *cli.Command {
	return &cli.Command{
		Name:  "invoices",
		Usage: "Invoice operations",
		Commands: []*cli.Command{
			{
				Name:  "list",
				Usage: "List invoices",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := internal.RequireFlag(cmd, "org", "ZOHO_BOOKS_ORG_ID")
					if err != nil {
						return err
					}
					raw, err := c.Request("GET", c.InvoiceBase+"/invoices", &zohttp.RequestOpts{Params: orgParams(orgID)})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:  "create",
				Usage: "Create an invoice",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "customer_id", Required: true, Usage: "Customer ID"},
					&cli.StringFlag{Name: "invoice_number", Usage: "Invoice number"},
					&cli.StringFlag{Name: "reference_number", Usage: "Reference number"},
					&cli.StringFlag{Name: "date", Usage: "Invoice date (YYYY-MM-DD)"},
					&cli.StringFlag{Name: "due_date", Usage: "Due date (YYYY-MM-DD)"},
					&cli.IntFlag{Name: "payment_terms", Usage: "Payment terms in days"},
					&cli.StringFlag{Name: "payment_terms_label", Usage: "Payment terms label (e.g. Net 30)"},
					&cli.StringFlag{Name: "notes", Usage: "Notes"},
					&cli.StringFlag{Name: "terms", Usage: "Terms and conditions"},
					&cli.StringFlag{Name: "salesperson_name", Usage: "Salesperson name"},
					&cli.StringFlag{Name: "json", Usage: "Additional fields as JSON"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := internal.RequireFlag(cmd, "org", "ZOHO_BOOKS_ORG_ID")
					if err != nil {
						return err
					}
					body := map[string]any{}
					body["customer_id"] = cmd.String("customer_id")
					if cmd.IsSet("invoice_number") {
						body["invoice_number"] = cmd.String("invoice_number")
					}
					if cmd.IsSet("reference_number") {
						body["reference_number"] = cmd.String("reference_number")
					}
					if cmd.IsSet("date") {
						body["date"] = cmd.String("date")
					}
					if cmd.IsSet("due_date") {
						body["due_date"] = cmd.String("due_date")
					}
					if cmd.IsSet("payment_terms") {
						body["payment_terms"] = cmd.Int("payment_terms")
					}
					if cmd.IsSet("payment_terms_label") {
						body["payment_terms_label"] = cmd.String("payment_terms_label")
					}
					if cmd.IsSet("notes") {
						body["notes"] = cmd.String("notes")
					}
					if cmd.IsSet("terms") {
						body["terms"] = cmd.String("terms")
					}
					if cmd.IsSet("salesperson_name") {
						body["salesperson_name"] = cmd.String("salesperson_name")
					}
					if err := internal.MergeJSON(cmd, body); err != nil {
						return err
					}
					raw, err := c.Request("POST", c.InvoiceBase+"/invoices", &zohttp.RequestOpts{
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
				Usage:     "Get an invoice",
				ArgsUsage: "<invoice-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					if cmd.Args().Len() < 1 {
						return internal.NewValidationError("invoice ID required")
					}
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := internal.RequireFlag(cmd, "org", "ZOHO_BOOKS_ORG_ID")
					if err != nil {
						return err
					}
					raw, err := c.Request("GET", c.InvoiceBase+"/invoices/"+cmd.Args().First(), &zohttp.RequestOpts{Params: orgParams(orgID)})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "update",
				Usage:     "Update an invoice",
				ArgsUsage: "<invoice-id>",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "customer_id", Usage: "Customer ID"},
					&cli.StringFlag{Name: "invoice_number", Usage: "Invoice number"},
					&cli.StringFlag{Name: "reference_number", Usage: "Reference number"},
					&cli.StringFlag{Name: "date", Usage: "Invoice date (YYYY-MM-DD)"},
					&cli.StringFlag{Name: "due_date", Usage: "Due date (YYYY-MM-DD)"},
					&cli.IntFlag{Name: "payment_terms", Usage: "Payment terms in days"},
					&cli.StringFlag{Name: "payment_terms_label", Usage: "Payment terms label (e.g. Net 30)"},
					&cli.StringFlag{Name: "notes", Usage: "Notes"},
					&cli.StringFlag{Name: "terms", Usage: "Terms and conditions"},
					&cli.StringFlag{Name: "salesperson_name", Usage: "Salesperson name"},
					&cli.StringFlag{Name: "json", Usage: "Additional fields as JSON"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					if cmd.Args().Len() < 1 {
						return internal.NewValidationError("invoice ID required")
					}
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := internal.RequireFlag(cmd, "org", "ZOHO_BOOKS_ORG_ID")
					if err != nil {
						return err
					}
					body := map[string]any{}
					if cmd.IsSet("customer_id") {
						body["customer_id"] = cmd.String("customer_id")
					}
					if cmd.IsSet("invoice_number") {
						body["invoice_number"] = cmd.String("invoice_number")
					}
					if cmd.IsSet("reference_number") {
						body["reference_number"] = cmd.String("reference_number")
					}
					if cmd.IsSet("date") {
						body["date"] = cmd.String("date")
					}
					if cmd.IsSet("due_date") {
						body["due_date"] = cmd.String("due_date")
					}
					if cmd.IsSet("payment_terms") {
						body["payment_terms"] = cmd.Int("payment_terms")
					}
					if cmd.IsSet("payment_terms_label") {
						body["payment_terms_label"] = cmd.String("payment_terms_label")
					}
					if cmd.IsSet("notes") {
						body["notes"] = cmd.String("notes")
					}
					if cmd.IsSet("terms") {
						body["terms"] = cmd.String("terms")
					}
					if cmd.IsSet("salesperson_name") {
						body["salesperson_name"] = cmd.String("salesperson_name")
					}
					if err := internal.MergeJSON(cmd, body); err != nil {
						return err
					}
					raw, err := c.Request("PUT", c.InvoiceBase+"/invoices/"+cmd.Args().First(), &zohttp.RequestOpts{
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
				Usage:     "Delete an invoice",
				ArgsUsage: "<invoice-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					if cmd.Args().Len() < 1 {
						return internal.NewValidationError("invoice ID required")
					}
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := internal.RequireFlag(cmd, "org", "ZOHO_BOOKS_ORG_ID")
					if err != nil {
						return err
					}
					raw, err := c.Request("DELETE", c.InvoiceBase+"/invoices/"+cmd.Args().First(), &zohttp.RequestOpts{Params: orgParams(orgID)})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "mark-sent",
				Usage:     "Mark an invoice as sent",
				ArgsUsage: "<invoice-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					if cmd.Args().Len() < 1 {
						return internal.NewValidationError("invoice ID required")
					}
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := internal.RequireFlag(cmd, "org", "ZOHO_BOOKS_ORG_ID")
					if err != nil {
						return err
					}
					raw, err := c.Request("POST", c.InvoiceBase+"/invoices/"+cmd.Args().First()+"/status/sent", &zohttp.RequestOpts{Params: orgParams(orgID)})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "void",
				Usage:     "Void an invoice",
				ArgsUsage: "<invoice-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					if cmd.Args().Len() < 1 {
						return internal.NewValidationError("invoice ID required")
					}
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := internal.RequireFlag(cmd, "org", "ZOHO_BOOKS_ORG_ID")
					if err != nil {
						return err
					}
					raw, err := c.Request("POST", c.InvoiceBase+"/invoices/"+cmd.Args().First()+"/status/void", &zohttp.RequestOpts{Params: orgParams(orgID)})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "mark-draft",
				Usage:     "Mark an invoice as draft",
				ArgsUsage: "<invoice-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					if cmd.Args().Len() < 1 {
						return internal.NewValidationError("invoice ID required")
					}
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := internal.RequireFlag(cmd, "org", "ZOHO_BOOKS_ORG_ID")
					if err != nil {
						return err
					}
					raw, err := c.Request("POST", c.InvoiceBase+"/invoices/"+cmd.Args().First()+"/status/draft", &zohttp.RequestOpts{Params: orgParams(orgID)})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "list-payments",
				Usage:     "List payments of an invoice",
				ArgsUsage: "<invoice-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					if cmd.Args().Len() < 1 {
						return internal.NewValidationError("invoice ID required")
					}
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := internal.RequireFlag(cmd, "org", "ZOHO_BOOKS_ORG_ID")
					if err != nil {
						return err
					}
					raw, err := c.Request("GET", c.InvoiceBase+"/invoices/"+cmd.Args().First()+"/payments", &zohttp.RequestOpts{Params: orgParams(orgID)})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "list-comments",
				Usage:     "List comments of an invoice",
				ArgsUsage: "<invoice-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					if cmd.Args().Len() < 1 {
						return internal.NewValidationError("invoice ID required")
					}
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := internal.RequireFlag(cmd, "org", "ZOHO_BOOKS_ORG_ID")
					if err != nil {
						return err
					}
					raw, err := c.Request("GET", c.InvoiceBase+"/invoices/"+cmd.Args().First()+"/comments", &zohttp.RequestOpts{Params: orgParams(orgID)})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "write-off",
				Usage:     "Write off an invoice",
				ArgsUsage: "<invoice-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					if cmd.Args().Len() < 1 {
						return internal.NewValidationError("invoice ID required")
					}
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := internal.RequireFlag(cmd, "org", "ZOHO_BOOKS_ORG_ID")
					if err != nil {
						return err
					}
					raw, err := c.Request("POST", c.InvoiceBase+"/invoices/"+cmd.Args().First()+"/writeoff", &zohttp.RequestOpts{Params: orgParams(orgID)})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "cancel-write-off",
				Usage:     "Cancel write-off of an invoice",
				ArgsUsage: "<invoice-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					if cmd.Args().Len() < 1 {
						return internal.NewValidationError("invoice ID required")
					}
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := internal.RequireFlag(cmd, "org", "ZOHO_BOOKS_ORG_ID")
					if err != nil {
						return err
					}
					raw, err := c.Request("POST", c.InvoiceBase+"/invoices/"+cmd.Args().First()+"/writeoff/cancel", &zohttp.RequestOpts{Params: orgParams(orgID)})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
		},
	}
}

func recurringInvoicesCmd() *cli.Command {
	return &cli.Command{
		Name:  "recurring-invoices",
		Usage: "Recurring invoice operations",
		Commands: []*cli.Command{
			{
				Name:  "list",
				Usage: "List recurring invoices",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := internal.RequireFlag(cmd, "org", "ZOHO_BOOKS_ORG_ID")
					if err != nil {
						return err
					}
					raw, err := c.Request("GET", c.InvoiceBase+"/recurringinvoices", &zohttp.RequestOpts{Params: orgParams(orgID)})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:  "create",
				Usage: "Create a recurring invoice",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "recurrence_name", Required: true, Usage: "Name for the recurring profile"},
					&cli.StringFlag{Name: "customer_id", Required: true, Usage: "Customer ID"},
					&cli.StringFlag{Name: "recurrence_frequency", Required: true, Usage: "Frequency (e.g. weekly, monthly, yearly)"},
					&cli.IntFlag{Name: "repeat_every", Usage: "Repeat interval"},
					&cli.StringFlag{Name: "start_date", Usage: "Start date (YYYY-MM-DD)"},
					&cli.StringFlag{Name: "end_date", Usage: "End date (YYYY-MM-DD)"},
					&cli.StringFlag{Name: "json", Usage: "Additional fields as JSON"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := internal.RequireFlag(cmd, "org", "ZOHO_BOOKS_ORG_ID")
					if err != nil {
						return err
					}
					body := map[string]any{}
					body["recurrence_name"] = cmd.String("recurrence_name")
					body["customer_id"] = cmd.String("customer_id")
					body["recurrence_frequency"] = cmd.String("recurrence_frequency")
					if cmd.IsSet("repeat_every") {
						body["repeat_every"] = cmd.Int("repeat_every")
					}
					if cmd.IsSet("start_date") {
						body["start_date"] = cmd.String("start_date")
					}
					if cmd.IsSet("end_date") {
						body["end_date"] = cmd.String("end_date")
					}
					if err := internal.MergeJSON(cmd, body); err != nil {
						return err
					}
					raw, err := c.Request("POST", c.InvoiceBase+"/recurringinvoices", &zohttp.RequestOpts{
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
				Usage:     "Get a recurring invoice",
				ArgsUsage: "<recurring-invoice-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					if cmd.Args().Len() < 1 {
						return internal.NewValidationError("recurring invoice ID required")
					}
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := internal.RequireFlag(cmd, "org", "ZOHO_BOOKS_ORG_ID")
					if err != nil {
						return err
					}
					raw, err := c.Request("GET", c.InvoiceBase+"/recurringinvoices/"+cmd.Args().First(), &zohttp.RequestOpts{Params: orgParams(orgID)})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "update",
				Usage:     "Update a recurring invoice",
				ArgsUsage: "<recurring-invoice-id>",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "recurrence_name", Usage: "Name for the recurring profile"},
					&cli.StringFlag{Name: "customer_id", Usage: "Customer ID"},
					&cli.StringFlag{Name: "recurrence_frequency", Usage: "Frequency (e.g. weekly, monthly, yearly)"},
					&cli.IntFlag{Name: "repeat_every", Usage: "Repeat interval"},
					&cli.StringFlag{Name: "start_date", Usage: "Start date (YYYY-MM-DD)"},
					&cli.StringFlag{Name: "end_date", Usage: "End date (YYYY-MM-DD)"},
					&cli.StringFlag{Name: "json", Usage: "Additional fields as JSON"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					if cmd.Args().Len() < 1 {
						return internal.NewValidationError("recurring invoice ID required")
					}
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := internal.RequireFlag(cmd, "org", "ZOHO_BOOKS_ORG_ID")
					if err != nil {
						return err
					}
					body := map[string]any{}
					if cmd.IsSet("recurrence_name") {
						body["recurrence_name"] = cmd.String("recurrence_name")
					}
					if cmd.IsSet("customer_id") {
						body["customer_id"] = cmd.String("customer_id")
					}
					if cmd.IsSet("recurrence_frequency") {
						body["recurrence_frequency"] = cmd.String("recurrence_frequency")
					}
					if cmd.IsSet("repeat_every") {
						body["repeat_every"] = cmd.Int("repeat_every")
					}
					if cmd.IsSet("start_date") {
						body["start_date"] = cmd.String("start_date")
					}
					if cmd.IsSet("end_date") {
						body["end_date"] = cmd.String("end_date")
					}
					if err := internal.MergeJSON(cmd, body); err != nil {
						return err
					}
					raw, err := c.Request("PUT", c.InvoiceBase+"/recurringinvoices/"+cmd.Args().First(), &zohttp.RequestOpts{
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
				Usage:     "Delete a recurring invoice",
				ArgsUsage: "<recurring-invoice-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					if cmd.Args().Len() < 1 {
						return internal.NewValidationError("recurring invoice ID required")
					}
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := internal.RequireFlag(cmd, "org", "ZOHO_BOOKS_ORG_ID")
					if err != nil {
						return err
					}
					raw, err := c.Request("DELETE", c.InvoiceBase+"/recurringinvoices/"+cmd.Args().First(), &zohttp.RequestOpts{Params: orgParams(orgID)})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "stop",
				Usage:     "Stop a recurring invoice",
				ArgsUsage: "<recurring-invoice-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					if cmd.Args().Len() < 1 {
						return internal.NewValidationError("recurring invoice ID required")
					}
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := internal.RequireFlag(cmd, "org", "ZOHO_BOOKS_ORG_ID")
					if err != nil {
						return err
					}
					raw, err := c.Request("POST", c.InvoiceBase+"/recurringinvoices/"+cmd.Args().First()+"/status/stop", &zohttp.RequestOpts{Params: orgParams(orgID)})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "resume",
				Usage:     "Resume a recurring invoice",
				ArgsUsage: "<recurring-invoice-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					if cmd.Args().Len() < 1 {
						return internal.NewValidationError("recurring invoice ID required")
					}
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := internal.RequireFlag(cmd, "org", "ZOHO_BOOKS_ORG_ID")
					if err != nil {
						return err
					}
					raw, err := c.Request("POST", c.InvoiceBase+"/recurringinvoices/"+cmd.Args().First()+"/status/resume", &zohttp.RequestOpts{Params: orgParams(orgID)})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
		},
	}
}

func creditNotesCmd() *cli.Command {
	return &cli.Command{
		Name:  "credit-notes",
		Usage: "Credit note operations",
		Commands: []*cli.Command{
			{
				Name:  "list",
				Usage: "List credit notes",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := internal.RequireFlag(cmd, "org", "ZOHO_BOOKS_ORG_ID")
					if err != nil {
						return err
					}
					raw, err := c.Request("GET", c.InvoiceBase+"/creditnotes", &zohttp.RequestOpts{Params: orgParams(orgID)})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:  "create",
				Usage: "Create a credit note",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "customer_id", Required: true, Usage: "Customer ID"},
					&cli.StringFlag{Name: "creditnote_number", Usage: "Credit note number"},
					&cli.StringFlag{Name: "reference_number", Usage: "Reference number"},
					&cli.StringFlag{Name: "date", Usage: "Date (YYYY-MM-DD)"},
					&cli.StringFlag{Name: "notes", Usage: "Notes"},
					&cli.StringFlag{Name: "terms", Usage: "Terms and conditions"},
					&cli.StringFlag{Name: "json", Usage: "Additional fields as JSON"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := internal.RequireFlag(cmd, "org", "ZOHO_BOOKS_ORG_ID")
					if err != nil {
						return err
					}
					body := map[string]any{}
					body["customer_id"] = cmd.String("customer_id")
					if cmd.IsSet("creditnote_number") {
						body["creditnote_number"] = cmd.String("creditnote_number")
					}
					if cmd.IsSet("reference_number") {
						body["reference_number"] = cmd.String("reference_number")
					}
					if cmd.IsSet("date") {
						body["date"] = cmd.String("date")
					}
					if cmd.IsSet("notes") {
						body["notes"] = cmd.String("notes")
					}
					if cmd.IsSet("terms") {
						body["terms"] = cmd.String("terms")
					}
					if err := internal.MergeJSON(cmd, body); err != nil {
						return err
					}
					raw, err := c.Request("POST", c.InvoiceBase+"/creditnotes", &zohttp.RequestOpts{
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
				Usage:     "Get a credit note",
				ArgsUsage: "<creditnote-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					if cmd.Args().Len() < 1 {
						return internal.NewValidationError("credit note ID required")
					}
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := internal.RequireFlag(cmd, "org", "ZOHO_BOOKS_ORG_ID")
					if err != nil {
						return err
					}
					raw, err := c.Request("GET", c.InvoiceBase+"/creditnotes/"+cmd.Args().First(), &zohttp.RequestOpts{Params: orgParams(orgID)})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "update",
				Usage:     "Update a credit note",
				ArgsUsage: "<creditnote-id>",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "customer_id", Usage: "Customer ID"},
					&cli.StringFlag{Name: "creditnote_number", Usage: "Credit note number"},
					&cli.StringFlag{Name: "reference_number", Usage: "Reference number"},
					&cli.StringFlag{Name: "date", Usage: "Date (YYYY-MM-DD)"},
					&cli.StringFlag{Name: "notes", Usage: "Notes"},
					&cli.StringFlag{Name: "terms", Usage: "Terms and conditions"},
					&cli.StringFlag{Name: "json", Usage: "Additional fields as JSON"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					if cmd.Args().Len() < 1 {
						return internal.NewValidationError("credit note ID required")
					}
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := internal.RequireFlag(cmd, "org", "ZOHO_BOOKS_ORG_ID")
					if err != nil {
						return err
					}
					body := map[string]any{}
					if cmd.IsSet("customer_id") {
						body["customer_id"] = cmd.String("customer_id")
					}
					if cmd.IsSet("creditnote_number") {
						body["creditnote_number"] = cmd.String("creditnote_number")
					}
					if cmd.IsSet("reference_number") {
						body["reference_number"] = cmd.String("reference_number")
					}
					if cmd.IsSet("date") {
						body["date"] = cmd.String("date")
					}
					if cmd.IsSet("notes") {
						body["notes"] = cmd.String("notes")
					}
					if cmd.IsSet("terms") {
						body["terms"] = cmd.String("terms")
					}
					if err := internal.MergeJSON(cmd, body); err != nil {
						return err
					}
					raw, err := c.Request("PUT", c.InvoiceBase+"/creditnotes/"+cmd.Args().First(), &zohttp.RequestOpts{
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
				Usage:     "Delete a credit note",
				ArgsUsage: "<creditnote-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					if cmd.Args().Len() < 1 {
						return internal.NewValidationError("credit note ID required")
					}
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := internal.RequireFlag(cmd, "org", "ZOHO_BOOKS_ORG_ID")
					if err != nil {
						return err
					}
					raw, err := c.Request("DELETE", c.InvoiceBase+"/creditnotes/"+cmd.Args().First(), &zohttp.RequestOpts{Params: orgParams(orgID)})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "void",
				Usage:     "Void a credit note",
				ArgsUsage: "<creditnote-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					if cmd.Args().Len() < 1 {
						return internal.NewValidationError("credit note ID required")
					}
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := internal.RequireFlag(cmd, "org", "ZOHO_BOOKS_ORG_ID")
					if err != nil {
						return err
					}
					raw, err := c.Request("POST", c.InvoiceBase+"/creditnotes/"+cmd.Args().First()+"/status/void", &zohttp.RequestOpts{Params: orgParams(orgID)})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "convert-to-draft",
				Usage:     "Convert a credit note to draft",
				ArgsUsage: "<creditnote-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					if cmd.Args().Len() < 1 {
						return internal.NewValidationError("credit note ID required")
					}
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := internal.RequireFlag(cmd, "org", "ZOHO_BOOKS_ORG_ID")
					if err != nil {
						return err
					}
					raw, err := c.Request("POST", c.InvoiceBase+"/creditnotes/"+cmd.Args().First()+"/status/draft", &zohttp.RequestOpts{Params: orgParams(orgID)})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "convert-to-open",
				Usage:     "Convert a credit note to open",
				ArgsUsage: "<creditnote-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					if cmd.Args().Len() < 1 {
						return internal.NewValidationError("credit note ID required")
					}
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := internal.RequireFlag(cmd, "org", "ZOHO_BOOKS_ORG_ID")
					if err != nil {
						return err
					}
					raw, err := c.Request("POST", c.InvoiceBase+"/creditnotes/"+cmd.Args().First()+"/status/open", &zohttp.RequestOpts{Params: orgParams(orgID)})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
		},
	}
}

func customerPaymentsCmd() *cli.Command {
	return &cli.Command{
		Name:  "customer-payments",
		Usage: "Customer payment operations",
		Commands: []*cli.Command{
			{
				Name:  "list",
				Usage: "List customer payments",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := internal.RequireFlag(cmd, "org", "ZOHO_BOOKS_ORG_ID")
					if err != nil {
						return err
					}
					raw, err := c.Request("GET", c.InvoiceBase+"/customerpayments", &zohttp.RequestOpts{Params: orgParams(orgID)})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:  "create",
				Usage: "Create a customer payment",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "customer_id", Required: true, Usage: "Customer ID"},
					&cli.StringFlag{Name: "payment_mode", Required: true, Usage: "Payment mode (e.g. cash, bank)"},
					&cli.FloatFlag{Name: "amount", Required: true, Usage: "Payment amount"},
					&cli.StringFlag{Name: "date", Required: true, Usage: "Payment date (YYYY-MM-DD)"},
					&cli.StringFlag{Name: "reference_number", Usage: "Reference number"},
					&cli.StringFlag{Name: "description", Usage: "Description"},
					&cli.StringFlag{Name: "json", Usage: "Additional fields as JSON"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := internal.RequireFlag(cmd, "org", "ZOHO_BOOKS_ORG_ID")
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
					if err := internal.MergeJSON(cmd, body); err != nil {
						return err
					}
					raw, err := c.Request("POST", c.InvoiceBase+"/customerpayments", &zohttp.RequestOpts{
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
				Action: func(_ context.Context, cmd *cli.Command) error {
					if cmd.Args().Len() < 1 {
						return internal.NewValidationError("payment ID required")
					}
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := internal.RequireFlag(cmd, "org", "ZOHO_BOOKS_ORG_ID")
					if err != nil {
						return err
					}
					raw, err := c.Request("GET", c.InvoiceBase+"/customerpayments/"+cmd.Args().First(), &zohttp.RequestOpts{Params: orgParams(orgID)})
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
					&cli.StringFlag{Name: "payment_mode", Usage: "Payment mode (e.g. cash, bank)"},
					&cli.FloatFlag{Name: "amount", Usage: "Payment amount"},
					&cli.StringFlag{Name: "date", Usage: "Payment date (YYYY-MM-DD)"},
					&cli.StringFlag{Name: "reference_number", Usage: "Reference number"},
					&cli.StringFlag{Name: "description", Usage: "Description"},
					&cli.StringFlag{Name: "json", Usage: "Additional fields as JSON"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					if cmd.Args().Len() < 1 {
						return internal.NewValidationError("payment ID required")
					}
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := internal.RequireFlag(cmd, "org", "ZOHO_BOOKS_ORG_ID")
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
					if err := internal.MergeJSON(cmd, body); err != nil {
						return err
					}
					raw, err := c.Request("PUT", c.InvoiceBase+"/customerpayments/"+cmd.Args().First(), &zohttp.RequestOpts{
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
				Usage:     "Delete a customer payment",
				ArgsUsage: "<payment-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					if cmd.Args().Len() < 1 {
						return internal.NewValidationError("payment ID required")
					}
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := internal.RequireFlag(cmd, "org", "ZOHO_BOOKS_ORG_ID")
					if err != nil {
						return err
					}
					raw, err := c.Request("DELETE", c.InvoiceBase+"/customerpayments/"+cmd.Args().First(), &zohttp.RequestOpts{Params: orgParams(orgID)})
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
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := internal.RequireFlag(cmd, "org", "ZOHO_BOOKS_ORG_ID")
					if err != nil {
						return err
					}
					raw, err := c.Request("GET", c.InvoiceBase+"/expenses", &zohttp.RequestOpts{Params: orgParams(orgID)})
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
					&cli.StringFlag{Name: "account_id", Required: true, Usage: "Expense account ID"},
					&cli.StringFlag{Name: "date", Required: true, Usage: "Expense date (YYYY-MM-DD)"},
					&cli.FloatFlag{Name: "amount", Required: true, Usage: "Expense amount"},
					&cli.StringFlag{Name: "customer_id", Usage: "Customer ID"},
					&cli.StringFlag{Name: "description", Usage: "Description"},
					&cli.StringFlag{Name: "reference_number", Usage: "Reference number"},
					&cli.StringFlag{Name: "json", Usage: "Additional fields as JSON"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := internal.RequireFlag(cmd, "org", "ZOHO_BOOKS_ORG_ID")
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
					if cmd.IsSet("description") {
						body["description"] = cmd.String("description")
					}
					if cmd.IsSet("reference_number") {
						body["reference_number"] = cmd.String("reference_number")
					}
					if err := internal.MergeJSON(cmd, body); err != nil {
						return err
					}
					raw, err := c.Request("POST", c.InvoiceBase+"/expenses", &zohttp.RequestOpts{
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
					if cmd.Args().Len() < 1 {
						return internal.NewValidationError("expense ID required")
					}
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := internal.RequireFlag(cmd, "org", "ZOHO_BOOKS_ORG_ID")
					if err != nil {
						return err
					}
					raw, err := c.Request("GET", c.InvoiceBase+"/expenses/"+cmd.Args().First(), &zohttp.RequestOpts{Params: orgParams(orgID)})
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
					&cli.StringFlag{Name: "account_id", Usage: "Expense account ID"},
					&cli.StringFlag{Name: "date", Usage: "Expense date (YYYY-MM-DD)"},
					&cli.FloatFlag{Name: "amount", Usage: "Expense amount"},
					&cli.StringFlag{Name: "customer_id", Usage: "Customer ID"},
					&cli.StringFlag{Name: "description", Usage: "Description"},
					&cli.StringFlag{Name: "reference_number", Usage: "Reference number"},
					&cli.StringFlag{Name: "json", Usage: "Additional fields as JSON"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					if cmd.Args().Len() < 1 {
						return internal.NewValidationError("expense ID required")
					}
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := internal.RequireFlag(cmd, "org", "ZOHO_BOOKS_ORG_ID")
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
					if cmd.IsSet("description") {
						body["description"] = cmd.String("description")
					}
					if cmd.IsSet("reference_number") {
						body["reference_number"] = cmd.String("reference_number")
					}
					if err := internal.MergeJSON(cmd, body); err != nil {
						return err
					}
					raw, err := c.Request("PUT", c.InvoiceBase+"/expenses/"+cmd.Args().First(), &zohttp.RequestOpts{
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
				Usage:     "Delete an expense",
				ArgsUsage: "<expense-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					if cmd.Args().Len() < 1 {
						return internal.NewValidationError("expense ID required")
					}
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := internal.RequireFlag(cmd, "org", "ZOHO_BOOKS_ORG_ID")
					if err != nil {
						return err
					}
					raw, err := c.Request("DELETE", c.InvoiceBase+"/expenses/"+cmd.Args().First(), &zohttp.RequestOpts{Params: orgParams(orgID)})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "list-history",
				Usage:     "List expense history and comments",
				ArgsUsage: "<expense-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					if cmd.Args().Len() < 1 {
						return internal.NewValidationError("expense ID required")
					}
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := internal.RequireFlag(cmd, "org", "ZOHO_BOOKS_ORG_ID")
					if err != nil {
						return err
					}
					raw, err := c.Request("GET", c.InvoiceBase+"/expenses/"+cmd.Args().First()+"/history", &zohttp.RequestOpts{Params: orgParams(orgID)})
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
				Name:  "list",
				Usage: "List recurring expenses",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := internal.RequireFlag(cmd, "org", "ZOHO_BOOKS_ORG_ID")
					if err != nil {
						return err
					}
					raw, err := c.Request("GET", c.InvoiceBase+"/recurringexpenses", &zohttp.RequestOpts{Params: orgParams(orgID)})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:  "create",
				Usage: "Create a recurring expense",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "account_id", Required: true, Usage: "Expense account ID"},
					&cli.StringFlag{Name: "recurrence_name", Required: true, Usage: "Name for the recurring expense"},
					&cli.StringFlag{Name: "start_date", Required: true, Usage: "Start date (YYYY-MM-DD)"},
					&cli.StringFlag{Name: "recurrence_frequency", Required: true, Usage: "Frequency (e.g. weekly, monthly, yearly)"},
					&cli.IntFlag{Name: "repeat_every", Required: true, Usage: "Repeat interval"},
					&cli.FloatFlag{Name: "amount", Required: true, Usage: "Expense amount"},
					&cli.StringFlag{Name: "end_date", Usage: "End date (YYYY-MM-DD)"},
					&cli.StringFlag{Name: "customer_id", Usage: "Customer ID"},
					&cli.StringFlag{Name: "json", Usage: "Additional fields as JSON"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := internal.RequireFlag(cmd, "org", "ZOHO_BOOKS_ORG_ID")
					if err != nil {
						return err
					}
					body := map[string]any{}
					body["account_id"] = cmd.String("account_id")
					body["recurrence_name"] = cmd.String("recurrence_name")
					body["start_date"] = cmd.String("start_date")
					body["recurrence_frequency"] = cmd.String("recurrence_frequency")
					body["repeat_every"] = cmd.Int("repeat_every")
					body["amount"] = cmd.Float("amount")
					if cmd.IsSet("end_date") {
						body["end_date"] = cmd.String("end_date")
					}
					if cmd.IsSet("customer_id") {
						body["customer_id"] = cmd.String("customer_id")
					}
					if err := internal.MergeJSON(cmd, body); err != nil {
						return err
					}
					raw, err := c.Request("POST", c.InvoiceBase+"/recurringexpenses", &zohttp.RequestOpts{
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
				ArgsUsage: "<recurring-expense-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					if cmd.Args().Len() < 1 {
						return internal.NewValidationError("recurring expense ID required")
					}
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := internal.RequireFlag(cmd, "org", "ZOHO_BOOKS_ORG_ID")
					if err != nil {
						return err
					}
					raw, err := c.Request("GET", c.InvoiceBase+"/recurringexpenses/"+cmd.Args().First(), &zohttp.RequestOpts{Params: orgParams(orgID)})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "update",
				Usage:     "Update a recurring expense",
				ArgsUsage: "<recurring-expense-id>",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "account_id", Usage: "Expense account ID"},
					&cli.StringFlag{Name: "recurrence_name", Usage: "Name for the recurring expense"},
					&cli.StringFlag{Name: "start_date", Usage: "Start date (YYYY-MM-DD)"},
					&cli.StringFlag{Name: "recurrence_frequency", Usage: "Frequency (e.g. weekly, monthly, yearly)"},
					&cli.IntFlag{Name: "repeat_every", Usage: "Repeat interval"},
					&cli.FloatFlag{Name: "amount", Usage: "Expense amount"},
					&cli.StringFlag{Name: "end_date", Usage: "End date (YYYY-MM-DD)"},
					&cli.StringFlag{Name: "customer_id", Usage: "Customer ID"},
					&cli.StringFlag{Name: "json", Usage: "Additional fields as JSON"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					if cmd.Args().Len() < 1 {
						return internal.NewValidationError("recurring expense ID required")
					}
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := internal.RequireFlag(cmd, "org", "ZOHO_BOOKS_ORG_ID")
					if err != nil {
						return err
					}
					body := map[string]any{}
					if cmd.IsSet("account_id") {
						body["account_id"] = cmd.String("account_id")
					}
					if cmd.IsSet("recurrence_name") {
						body["recurrence_name"] = cmd.String("recurrence_name")
					}
					if cmd.IsSet("start_date") {
						body["start_date"] = cmd.String("start_date")
					}
					if cmd.IsSet("recurrence_frequency") {
						body["recurrence_frequency"] = cmd.String("recurrence_frequency")
					}
					if cmd.IsSet("repeat_every") {
						body["repeat_every"] = cmd.Int("repeat_every")
					}
					if cmd.IsSet("amount") {
						body["amount"] = cmd.Float("amount")
					}
					if cmd.IsSet("end_date") {
						body["end_date"] = cmd.String("end_date")
					}
					if cmd.IsSet("customer_id") {
						body["customer_id"] = cmd.String("customer_id")
					}
					if err := internal.MergeJSON(cmd, body); err != nil {
						return err
					}
					raw, err := c.Request("PUT", c.InvoiceBase+"/recurringexpenses/"+cmd.Args().First(), &zohttp.RequestOpts{
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
				Usage:     "Delete a recurring expense",
				ArgsUsage: "<recurring-expense-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					if cmd.Args().Len() < 1 {
						return internal.NewValidationError("recurring expense ID required")
					}
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := internal.RequireFlag(cmd, "org", "ZOHO_BOOKS_ORG_ID")
					if err != nil {
						return err
					}
					raw, err := c.Request("DELETE", c.InvoiceBase+"/recurringexpenses/"+cmd.Args().First(), &zohttp.RequestOpts{Params: orgParams(orgID)})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "stop",
				Usage:     "Stop a recurring expense",
				ArgsUsage: "<recurring-expense-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					if cmd.Args().Len() < 1 {
						return internal.NewValidationError("recurring expense ID required")
					}
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := internal.RequireFlag(cmd, "org", "ZOHO_BOOKS_ORG_ID")
					if err != nil {
						return err
					}
					raw, err := c.Request("POST", c.InvoiceBase+"/recurringexpenses/"+cmd.Args().First()+"/status/stop", &zohttp.RequestOpts{Params: orgParams(orgID)})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "resume",
				Usage:     "Resume a recurring expense",
				ArgsUsage: "<recurring-expense-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					if cmd.Args().Len() < 1 {
						return internal.NewValidationError("recurring expense ID required")
					}
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := internal.RequireFlag(cmd, "org", "ZOHO_BOOKS_ORG_ID")
					if err != nil {
						return err
					}
					raw, err := c.Request("POST", c.InvoiceBase+"/recurringexpenses/"+cmd.Args().First()+"/status/resume", &zohttp.RequestOpts{Params: orgParams(orgID)})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "list-child-expenses",
				Usage:     "List child expenses of a recurring expense",
				ArgsUsage: "<recurring-expense-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					if cmd.Args().Len() < 1 {
						return internal.NewValidationError("recurring expense ID required")
					}
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := internal.RequireFlag(cmd, "org", "ZOHO_BOOKS_ORG_ID")
					if err != nil {
						return err
					}
					raw, err := c.Request("GET", c.InvoiceBase+"/recurringexpenses/"+cmd.Args().First()+"/expenses", &zohttp.RequestOpts{Params: orgParams(orgID)})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "list-history",
				Usage:     "List recurring expense history",
				ArgsUsage: "<recurring-expense-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					if cmd.Args().Len() < 1 {
						return internal.NewValidationError("recurring expense ID required")
					}
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := internal.RequireFlag(cmd, "org", "ZOHO_BOOKS_ORG_ID")
					if err != nil {
						return err
					}
					raw, err := c.Request("GET", c.InvoiceBase+"/recurringexpenses/"+cmd.Args().First()+"/comments", &zohttp.RequestOpts{Params: orgParams(orgID)})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
		},
	}
}

func retainerInvoicesCmd() *cli.Command {
	return &cli.Command{
		Name:  "retainer-invoices",
		Usage: "Retainer invoice operations",
		Commands: []*cli.Command{
			{
				Name:  "list",
				Usage: "List retainer invoices",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := internal.RequireFlag(cmd, "org", "ZOHO_BOOKS_ORG_ID")
					if err != nil {
						return err
					}
					raw, err := c.Request("GET", c.InvoiceBase+"/retainerinvoices", &zohttp.RequestOpts{Params: orgParams(orgID)})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:  "create",
				Usage: "Create a retainer invoice",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "customer_id", Required: true, Usage: "Customer ID"},
					&cli.StringFlag{Name: "reference_number", Usage: "Reference number"},
					&cli.StringFlag{Name: "date", Usage: "Date (YYYY-MM-DD)"},
					&cli.StringFlag{Name: "notes", Usage: "Notes"},
					&cli.StringFlag{Name: "terms", Usage: "Terms and conditions"},
					&cli.StringFlag{Name: "json", Usage: "Additional fields as JSON"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := internal.RequireFlag(cmd, "org", "ZOHO_BOOKS_ORG_ID")
					if err != nil {
						return err
					}
					body := map[string]any{}
					body["customer_id"] = cmd.String("customer_id")
					if cmd.IsSet("reference_number") {
						body["reference_number"] = cmd.String("reference_number")
					}
					if cmd.IsSet("date") {
						body["date"] = cmd.String("date")
					}
					if cmd.IsSet("notes") {
						body["notes"] = cmd.String("notes")
					}
					if cmd.IsSet("terms") {
						body["terms"] = cmd.String("terms")
					}
					if err := internal.MergeJSON(cmd, body); err != nil {
						return err
					}
					raw, err := c.Request("POST", c.InvoiceBase+"/retainerinvoices", &zohttp.RequestOpts{
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
				Usage:     "Get a retainer invoice",
				ArgsUsage: "<retainer-invoice-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					if cmd.Args().Len() < 1 {
						return internal.NewValidationError("retainer invoice ID required")
					}
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := internal.RequireFlag(cmd, "org", "ZOHO_BOOKS_ORG_ID")
					if err != nil {
						return err
					}
					raw, err := c.Request("GET", c.InvoiceBase+"/retainerinvoices/"+cmd.Args().First(), &zohttp.RequestOpts{Params: orgParams(orgID)})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "update",
				Usage:     "Update a retainer invoice",
				ArgsUsage: "<retainer-invoice-id>",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "customer_id", Usage: "Customer ID"},
					&cli.StringFlag{Name: "reference_number", Usage: "Reference number"},
					&cli.StringFlag{Name: "date", Usage: "Date (YYYY-MM-DD)"},
					&cli.StringFlag{Name: "notes", Usage: "Notes"},
					&cli.StringFlag{Name: "terms", Usage: "Terms and conditions"},
					&cli.StringFlag{Name: "json", Usage: "Additional fields as JSON"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					if cmd.Args().Len() < 1 {
						return internal.NewValidationError("retainer invoice ID required")
					}
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := internal.RequireFlag(cmd, "org", "ZOHO_BOOKS_ORG_ID")
					if err != nil {
						return err
					}
					body := map[string]any{}
					if cmd.IsSet("customer_id") {
						body["customer_id"] = cmd.String("customer_id")
					}
					if cmd.IsSet("reference_number") {
						body["reference_number"] = cmd.String("reference_number")
					}
					if cmd.IsSet("date") {
						body["date"] = cmd.String("date")
					}
					if cmd.IsSet("notes") {
						body["notes"] = cmd.String("notes")
					}
					if cmd.IsSet("terms") {
						body["terms"] = cmd.String("terms")
					}
					if err := internal.MergeJSON(cmd, body); err != nil {
						return err
					}
					raw, err := c.Request("PUT", c.InvoiceBase+"/retainerinvoices/"+cmd.Args().First(), &zohttp.RequestOpts{
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
				Usage:     "Delete a retainer invoice",
				ArgsUsage: "<retainer-invoice-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					if cmd.Args().Len() < 1 {
						return internal.NewValidationError("retainer invoice ID required")
					}
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := internal.RequireFlag(cmd, "org", "ZOHO_BOOKS_ORG_ID")
					if err != nil {
						return err
					}
					raw, err := c.Request("DELETE", c.InvoiceBase+"/retainerinvoices/"+cmd.Args().First(), &zohttp.RequestOpts{Params: orgParams(orgID)})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "mark-sent",
				Usage:     "Mark a retainer invoice as sent",
				ArgsUsage: "<retainer-invoice-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					if cmd.Args().Len() < 1 {
						return internal.NewValidationError("retainer invoice ID required")
					}
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := internal.RequireFlag(cmd, "org", "ZOHO_BOOKS_ORG_ID")
					if err != nil {
						return err
					}
					raw, err := c.Request("POST", c.InvoiceBase+"/retainerinvoices/"+cmd.Args().First()+"/status/sent", &zohttp.RequestOpts{Params: orgParams(orgID)})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "void",
				Usage:     "Void a retainer invoice",
				ArgsUsage: "<retainer-invoice-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					if cmd.Args().Len() < 1 {
						return internal.NewValidationError("retainer invoice ID required")
					}
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := internal.RequireFlag(cmd, "org", "ZOHO_BOOKS_ORG_ID")
					if err != nil {
						return err
					}
					raw, err := c.Request("POST", c.InvoiceBase+"/retainerinvoices/"+cmd.Args().First()+"/status/void", &zohttp.RequestOpts{Params: orgParams(orgID)})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "mark-draft",
				Usage:     "Mark a retainer invoice as draft",
				ArgsUsage: "<retainer-invoice-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					if cmd.Args().Len() < 1 {
						return internal.NewValidationError("retainer invoice ID required")
					}
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := internal.RequireFlag(cmd, "org", "ZOHO_BOOKS_ORG_ID")
					if err != nil {
						return err
					}
					raw, err := c.Request("POST", c.InvoiceBase+"/retainerinvoices/"+cmd.Args().First()+"/status/draft", &zohttp.RequestOpts{Params: orgParams(orgID)})
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
		Usage: "Project operations",
		Commands: []*cli.Command{
			{
				Name:  "list",
				Usage: "List projects",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := internal.RequireFlag(cmd, "org", "ZOHO_BOOKS_ORG_ID")
					if err != nil {
						return err
					}
					raw, err := c.Request("GET", c.InvoiceBase+"/projects", &zohttp.RequestOpts{Params: orgParams(orgID)})
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
					&cli.StringFlag{Name: "customer_id", Required: true, Usage: "Customer ID"},
					&cli.StringFlag{Name: "billing_type", Required: true, Usage: "Billing type (fixed_cost_for_project, based_on_project_hours, based_on_staff_hours, based_on_task_hours)"},
					&cli.StringFlag{Name: "description", Usage: "Project description"},
					&cli.FloatFlag{Name: "rate", Usage: "Hourly rate"},
					&cli.StringFlag{Name: "budget_type", Usage: "Budget type"},
					&cli.IntFlag{Name: "budget_hours", Usage: "Budget hours"},
					&cli.StringFlag{Name: "json", Usage: "Additional fields as JSON"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := internal.RequireFlag(cmd, "org", "ZOHO_BOOKS_ORG_ID")
					if err != nil {
						return err
					}
					body := map[string]any{}
					body["project_name"] = cmd.String("project_name")
					body["customer_id"] = cmd.String("customer_id")
					body["billing_type"] = cmd.String("billing_type")
					if cmd.IsSet("description") {
						body["description"] = cmd.String("description")
					}
					if cmd.IsSet("rate") {
						body["rate"] = cmd.Float("rate")
					}
					if cmd.IsSet("budget_type") {
						body["budget_type"] = cmd.String("budget_type")
					}
					if cmd.IsSet("budget_hours") {
						body["budget_hours"] = cmd.Int("budget_hours")
					}
					if err := internal.MergeJSON(cmd, body); err != nil {
						return err
					}
					raw, err := c.Request("POST", c.InvoiceBase+"/projects", &zohttp.RequestOpts{
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
				Usage:     "Get a project",
				ArgsUsage: "<project-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					if cmd.Args().Len() < 1 {
						return internal.NewValidationError("project ID required")
					}
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := internal.RequireFlag(cmd, "org", "ZOHO_BOOKS_ORG_ID")
					if err != nil {
						return err
					}
					raw, err := c.Request("GET", c.InvoiceBase+"/projects/"+cmd.Args().First(), &zohttp.RequestOpts{Params: orgParams(orgID)})
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
					&cli.StringFlag{Name: "billing_type", Usage: "Billing type"},
					&cli.StringFlag{Name: "description", Usage: "Project description"},
					&cli.FloatFlag{Name: "rate", Usage: "Hourly rate"},
					&cli.StringFlag{Name: "budget_type", Usage: "Budget type"},
					&cli.IntFlag{Name: "budget_hours", Usage: "Budget hours"},
					&cli.StringFlag{Name: "json", Usage: "Additional fields as JSON"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					if cmd.Args().Len() < 1 {
						return internal.NewValidationError("project ID required")
					}
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := internal.RequireFlag(cmd, "org", "ZOHO_BOOKS_ORG_ID")
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
					if cmd.IsSet("billing_type") {
						body["billing_type"] = cmd.String("billing_type")
					}
					if cmd.IsSet("description") {
						body["description"] = cmd.String("description")
					}
					if cmd.IsSet("rate") {
						body["rate"] = cmd.Float("rate")
					}
					if cmd.IsSet("budget_type") {
						body["budget_type"] = cmd.String("budget_type")
					}
					if cmd.IsSet("budget_hours") {
						body["budget_hours"] = cmd.Int("budget_hours")
					}
					if err := internal.MergeJSON(cmd, body); err != nil {
						return err
					}
					raw, err := c.Request("PUT", c.InvoiceBase+"/projects/"+cmd.Args().First(), &zohttp.RequestOpts{
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
				Usage:     "Delete a project",
				ArgsUsage: "<project-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					if cmd.Args().Len() < 1 {
						return internal.NewValidationError("project ID required")
					}
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := internal.RequireFlag(cmd, "org", "ZOHO_BOOKS_ORG_ID")
					if err != nil {
						return err
					}
					raw, err := c.Request("DELETE", c.InvoiceBase+"/projects/"+cmd.Args().First(), &zohttp.RequestOpts{Params: orgParams(orgID)})
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
					if cmd.Args().Len() < 1 {
						return internal.NewValidationError("project ID required")
					}
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := internal.RequireFlag(cmd, "org", "ZOHO_BOOKS_ORG_ID")
					if err != nil {
						return err
					}
					raw, err := c.Request("POST", c.InvoiceBase+"/projects/"+cmd.Args().First()+"/activate", &zohttp.RequestOpts{Params: orgParams(orgID)})
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
					if cmd.Args().Len() < 1 {
						return internal.NewValidationError("project ID required")
					}
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := internal.RequireFlag(cmd, "org", "ZOHO_BOOKS_ORG_ID")
					if err != nil {
						return err
					}
					raw, err := c.Request("POST", c.InvoiceBase+"/projects/"+cmd.Args().First()+"/deactivate", &zohttp.RequestOpts{Params: orgParams(orgID)})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "clone",
				Usage:     "Clone a project",
				ArgsUsage: "<project-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					if cmd.Args().Len() < 1 {
						return internal.NewValidationError("project ID required")
					}
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := internal.RequireFlag(cmd, "org", "ZOHO_BOOKS_ORG_ID")
					if err != nil {
						return err
					}
					raw, err := c.Request("POST", c.InvoiceBase+"/projects/"+cmd.Args().First()+"/clone", &zohttp.RequestOpts{Params: orgParams(orgID)})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "list-comments",
				Usage:     "List comments of a project",
				ArgsUsage: "<project-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					if cmd.Args().Len() < 1 {
						return internal.NewValidationError("project ID required")
					}
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := internal.RequireFlag(cmd, "org", "ZOHO_BOOKS_ORG_ID")
					if err != nil {
						return err
					}
					raw, err := c.Request("GET", c.InvoiceBase+"/projects/"+cmd.Args().First()+"/comments", &zohttp.RequestOpts{Params: orgParams(orgID)})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "list-invoices",
				Usage:     "List invoices of a project",
				ArgsUsage: "<project-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					if cmd.Args().Len() < 1 {
						return internal.NewValidationError("project ID required")
					}
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := internal.RequireFlag(cmd, "org", "ZOHO_BOOKS_ORG_ID")
					if err != nil {
						return err
					}
					raw, err := c.Request("GET", c.InvoiceBase+"/projects/"+cmd.Args().First()+"/invoices", &zohttp.RequestOpts{Params: orgParams(orgID)})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
		},
	}
}

func itemsCmd() *cli.Command {
	return &cli.Command{
		Name:  "items",
		Usage: "Item operations",
		Commands: []*cli.Command{
			{
				Name:  "list",
				Usage: "List items",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := internal.RequireFlag(cmd, "org", "ZOHO_BOOKS_ORG_ID")
					if err != nil {
						return err
					}
					raw, err := c.Request("GET", c.InvoiceBase+"/items", &zohttp.RequestOpts{Params: orgParams(orgID)})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:  "create",
				Usage: "Create an item",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "name", Required: true, Usage: "Item name"},
					&cli.FloatFlag{Name: "rate", Required: true, Usage: "Item rate"},
					&cli.StringFlag{Name: "description", Usage: "Item description"},
					&cli.StringFlag{Name: "unit", Usage: "Unit of measurement"},
					&cli.StringFlag{Name: "sku", Usage: "SKU"},
					&cli.StringFlag{Name: "product_type", Usage: "Product type (goods or service)"},
					&cli.StringFlag{Name: "json", Usage: "Additional fields as JSON"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := internal.RequireFlag(cmd, "org", "ZOHO_BOOKS_ORG_ID")
					if err != nil {
						return err
					}
					body := map[string]any{}
					body["name"] = cmd.String("name")
					body["rate"] = cmd.Float("rate")
					if cmd.IsSet("description") {
						body["description"] = cmd.String("description")
					}
					if cmd.IsSet("unit") {
						body["unit"] = cmd.String("unit")
					}
					if cmd.IsSet("sku") {
						body["sku"] = cmd.String("sku")
					}
					if cmd.IsSet("product_type") {
						body["product_type"] = cmd.String("product_type")
					}
					if err := internal.MergeJSON(cmd, body); err != nil {
						return err
					}
					raw, err := c.Request("POST", c.InvoiceBase+"/items", &zohttp.RequestOpts{
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
				Usage:     "Get an item",
				ArgsUsage: "<item-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					if cmd.Args().Len() < 1 {
						return internal.NewValidationError("item ID required")
					}
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := internal.RequireFlag(cmd, "org", "ZOHO_BOOKS_ORG_ID")
					if err != nil {
						return err
					}
					raw, err := c.Request("GET", c.InvoiceBase+"/items/"+cmd.Args().First(), &zohttp.RequestOpts{Params: orgParams(orgID)})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "update",
				Usage:     "Update an item",
				ArgsUsage: "<item-id>",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "name", Usage: "Item name"},
					&cli.FloatFlag{Name: "rate", Usage: "Item rate"},
					&cli.StringFlag{Name: "description", Usage: "Item description"},
					&cli.StringFlag{Name: "unit", Usage: "Unit of measurement"},
					&cli.StringFlag{Name: "sku", Usage: "SKU"},
					&cli.StringFlag{Name: "product_type", Usage: "Product type (goods or service)"},
					&cli.StringFlag{Name: "json", Usage: "Additional fields as JSON"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					if cmd.Args().Len() < 1 {
						return internal.NewValidationError("item ID required")
					}
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := internal.RequireFlag(cmd, "org", "ZOHO_BOOKS_ORG_ID")
					if err != nil {
						return err
					}
					body := map[string]any{}
					if cmd.IsSet("name") {
						body["name"] = cmd.String("name")
					}
					if cmd.IsSet("rate") {
						body["rate"] = cmd.Float("rate")
					}
					if cmd.IsSet("description") {
						body["description"] = cmd.String("description")
					}
					if cmd.IsSet("unit") {
						body["unit"] = cmd.String("unit")
					}
					if cmd.IsSet("sku") {
						body["sku"] = cmd.String("sku")
					}
					if cmd.IsSet("product_type") {
						body["product_type"] = cmd.String("product_type")
					}
					if err := internal.MergeJSON(cmd, body); err != nil {
						return err
					}
					raw, err := c.Request("PUT", c.InvoiceBase+"/items/"+cmd.Args().First(), &zohttp.RequestOpts{
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
				Usage:     "Delete an item",
				ArgsUsage: "<item-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					if cmd.Args().Len() < 1 {
						return internal.NewValidationError("item ID required")
					}
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := internal.RequireFlag(cmd, "org", "ZOHO_BOOKS_ORG_ID")
					if err != nil {
						return err
					}
					raw, err := c.Request("DELETE", c.InvoiceBase+"/items/"+cmd.Args().First(), &zohttp.RequestOpts{Params: orgParams(orgID)})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "mark-active",
				Usage:     "Mark an item as active",
				ArgsUsage: "<item-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					if cmd.Args().Len() < 1 {
						return internal.NewValidationError("item ID required")
					}
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := internal.RequireFlag(cmd, "org", "ZOHO_BOOKS_ORG_ID")
					if err != nil {
						return err
					}
					raw, err := c.Request("POST", c.InvoiceBase+"/items/"+cmd.Args().First()+"/active", &zohttp.RequestOpts{Params: orgParams(orgID)})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "mark-inactive",
				Usage:     "Mark an item as inactive",
				ArgsUsage: "<item-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					if cmd.Args().Len() < 1 {
						return internal.NewValidationError("item ID required")
					}
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := internal.RequireFlag(cmd, "org", "ZOHO_BOOKS_ORG_ID")
					if err != nil {
						return err
					}
					raw, err := c.Request("POST", c.InvoiceBase+"/items/"+cmd.Args().First()+"/inactive", &zohttp.RequestOpts{Params: orgParams(orgID)})
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
		Usage: "Currency operations",
		Commands: []*cli.Command{
			{
				Name:  "list",
				Usage: "List currencies",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := internal.RequireFlag(cmd, "org", "ZOHO_BOOKS_ORG_ID")
					if err != nil {
						return err
					}
					raw, err := c.Request("GET", c.InvoiceBase+"/settings/currencies", &zohttp.RequestOpts{Params: orgParams(orgID)})
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
					&cli.StringFlag{Name: "currency_code", Required: true, Usage: "Currency code (e.g. USD, EUR)"},
					&cli.StringFlag{Name: "currency_symbol", Required: true, Usage: "Currency symbol (e.g. $)"},
					&cli.StringFlag{Name: "currency_format", Usage: "Currency format"},
					&cli.IntFlag{Name: "price_precision", Usage: "Decimal precision"},
					&cli.StringFlag{Name: "json", Usage: "Additional fields as JSON"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := internal.RequireFlag(cmd, "org", "ZOHO_BOOKS_ORG_ID")
					if err != nil {
						return err
					}
					body := map[string]any{}
					body["currency_code"] = cmd.String("currency_code")
					body["currency_symbol"] = cmd.String("currency_symbol")
					if cmd.IsSet("currency_format") {
						body["currency_format"] = cmd.String("currency_format")
					}
					if cmd.IsSet("price_precision") {
						body["price_precision"] = cmd.Int("price_precision")
					}
					if err := internal.MergeJSON(cmd, body); err != nil {
						return err
					}
					raw, err := c.Request("POST", c.InvoiceBase+"/settings/currencies", &zohttp.RequestOpts{
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
				Usage:     "Get a currency",
				ArgsUsage: "<currency-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					if cmd.Args().Len() < 1 {
						return internal.NewValidationError("currency ID required")
					}
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := internal.RequireFlag(cmd, "org", "ZOHO_BOOKS_ORG_ID")
					if err != nil {
						return err
					}
					raw, err := c.Request("GET", c.InvoiceBase+"/settings/currencies/"+cmd.Args().First(), &zohttp.RequestOpts{Params: orgParams(orgID)})
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
					&cli.StringFlag{Name: "currency_code", Usage: "Currency code (e.g. USD, EUR)"},
					&cli.StringFlag{Name: "currency_symbol", Usage: "Currency symbol (e.g. $)"},
					&cli.StringFlag{Name: "currency_format", Usage: "Currency format"},
					&cli.IntFlag{Name: "price_precision", Usage: "Decimal precision"},
					&cli.StringFlag{Name: "json", Usage: "Additional fields as JSON"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					if cmd.Args().Len() < 1 {
						return internal.NewValidationError("currency ID required")
					}
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := internal.RequireFlag(cmd, "org", "ZOHO_BOOKS_ORG_ID")
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
					raw, err := c.Request("PUT", c.InvoiceBase+"/settings/currencies/"+cmd.Args().First(), &zohttp.RequestOpts{
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
				Usage:     "Delete a currency",
				ArgsUsage: "<currency-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					if cmd.Args().Len() < 1 {
						return internal.NewValidationError("currency ID required")
					}
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := internal.RequireFlag(cmd, "org", "ZOHO_BOOKS_ORG_ID")
					if err != nil {
						return err
					}
					raw, err := c.Request("DELETE", c.InvoiceBase+"/settings/currencies/"+cmd.Args().First(), &zohttp.RequestOpts{Params: orgParams(orgID)})
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
		Usage: "Tax operations",
		Commands: []*cli.Command{
			{
				Name:  "list",
				Usage: "List taxes",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := internal.RequireFlag(cmd, "org", "ZOHO_BOOKS_ORG_ID")
					if err != nil {
						return err
					}
					raw, err := c.Request("GET", c.InvoiceBase+"/settings/taxes", &zohttp.RequestOpts{Params: orgParams(orgID)})
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
					&cli.StringFlag{Name: "tax_type", Required: true, Usage: "Tax type (tax or compound_tax)"},
					&cli.StringFlag{Name: "json", Usage: "Additional fields as JSON"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := internal.RequireFlag(cmd, "org", "ZOHO_BOOKS_ORG_ID")
					if err != nil {
						return err
					}
					body := map[string]any{}
					body["tax_name"] = cmd.String("tax_name")
					body["tax_percentage"] = cmd.Float("tax_percentage")
					body["tax_type"] = cmd.String("tax_type")
					if err := internal.MergeJSON(cmd, body); err != nil {
						return err
					}
					raw, err := c.Request("POST", c.InvoiceBase+"/settings/taxes", &zohttp.RequestOpts{
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
				Usage:     "Get a tax",
				ArgsUsage: "<tax-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					if cmd.Args().Len() < 1 {
						return internal.NewValidationError("tax ID required")
					}
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := internal.RequireFlag(cmd, "org", "ZOHO_BOOKS_ORG_ID")
					if err != nil {
						return err
					}
					raw, err := c.Request("GET", c.InvoiceBase+"/settings/taxes/"+cmd.Args().First(), &zohttp.RequestOpts{Params: orgParams(orgID)})
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
					&cli.StringFlag{Name: "tax_type", Usage: "Tax type (tax or compound_tax)"},
					&cli.StringFlag{Name: "json", Usage: "Additional fields as JSON"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					if cmd.Args().Len() < 1 {
						return internal.NewValidationError("tax ID required")
					}
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := internal.RequireFlag(cmd, "org", "ZOHO_BOOKS_ORG_ID")
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
					raw, err := c.Request("PUT", c.InvoiceBase+"/settings/taxes/"+cmd.Args().First(), &zohttp.RequestOpts{
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
				Usage:     "Delete a tax",
				ArgsUsage: "<tax-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					if cmd.Args().Len() < 1 {
						return internal.NewValidationError("tax ID required")
					}
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := internal.RequireFlag(cmd, "org", "ZOHO_BOOKS_ORG_ID")
					if err != nil {
						return err
					}
					raw, err := c.Request("DELETE", c.InvoiceBase+"/settings/taxes/"+cmd.Args().First(), &zohttp.RequestOpts{Params: orgParams(orgID)})
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
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := internal.RequireFlag(cmd, "org", "ZOHO_BOOKS_ORG_ID")
					if err != nil {
						return err
					}
					raw, err := c.Request("GET", c.InvoiceBase+"/users", &zohttp.RequestOpts{Params: orgParams(orgID)})
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
					if cmd.Args().Len() < 1 {
						return internal.NewValidationError("user ID required")
					}
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := internal.RequireFlag(cmd, "org", "ZOHO_BOOKS_ORG_ID")
					if err != nil {
						return err
					}
					raw, err := c.Request("GET", c.InvoiceBase+"/users/"+cmd.Args().First(), &zohttp.RequestOpts{Params: orgParams(orgID)})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:  "get-current",
				Usage: "Get current user",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := internal.RequireFlag(cmd, "org", "ZOHO_BOOKS_ORG_ID")
					if err != nil {
						return err
					}
					raw, err := c.Request("GET", c.InvoiceBase+"/users/me", &zohttp.RequestOpts{Params: orgParams(orgID)})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
		},
	}
}
