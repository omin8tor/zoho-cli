package books

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

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
		org = os.Getenv("ZOHO_BOOKS_ORG_ID")
	}
	if org == "" {
		return "", internal.NewValidationError("--org flag or ZOHO_BOOKS_ORG_ID env var required")
	}
	return org, nil
}

func orgParams(orgID string) map[string]string {
	return map[string]string{"organization_id": orgID}
}

func mergeParams(orgID string, extra map[string]string) map[string]string {
	params := orgParams(orgID)
	for k, v := range extra {
		params[k] = v
	}
	return params
}

func Commands() *cli.Command {
	return &cli.Command{
		Name:  "books",
		Usage: "Zoho Books operations",
		Flags: []cli.Flag{
			&cli.StringFlag{Name: "org", Usage: "Organization ID (or set ZOHO_BOOKS_ORG_ID)"},
		},
		Commands: []*cli.Command{
			organizationsCmd(),
			contactsCmd(),
			contactPersonsCmd(),
			estimatesCmd(),
			salesOrdersCmd(),
			salesReceiptsCmd(),
			invoicesCmd(),
			recurringInvoicesCmd(),
			creditNotesCmd(),
			customerDebitNotesCmd(),
			customerPaymentsCmd(),
			expensesCmd(),
			recurringExpensesCmd(),
			retainerInvoicesCmd(),
			purchaseOrdersCmd(),
			billsCmd(),
			recurringBillsCmd(),
			vendorCreditsCmd(),
			vendorPaymentsCmd(),
			customModulesCmd(),
			bankAccountsCmd(),
			bankTransactionsCmd(),
			bankRulesCmd(),
			chartOfAccountsCmd(),
			journalsCmd(),
			fixedAssetsCmd(),
			baseCurrencyAdjCmd(),
			booksProjectsCmd(),
			tasksCmd(),
			timeEntriesCmd(),
			usersCmd(),
			itemsCmd(),
			locationsCmd(),
			currenciesCmd(),
			taxesCmd(),
			openingBalanceCmd(),
			crmIntegrationCmd(),
			reportingTagsCmd(),
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
					raw, err := c.Request("GET", c.BooksBase+"/organizations", nil)
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
					&cli.StringFlag{Name: "currency_code", Required: true, Usage: "Currency code (e.g. USD)"},
					&cli.StringFlag{Name: "time_zone", Required: true, Usage: "Time zone (e.g. PST)"},
					&cli.StringFlag{Name: "fiscal_year_start_month", Usage: "Fiscal year start month (e.g. january)"},
					&cli.StringFlag{Name: "date_format", Usage: "Date format (e.g. dd MMM yyyy)"},
					&cli.StringFlag{Name: "language_code", Usage: "Language code (e.g. en)"},
					&cli.StringFlag{Name: "industry_type", Usage: "Industry type"},
					&cli.StringFlag{Name: "portal_name", Usage: "Customer portal name"},
					&cli.StringFlag{Name: "json", Usage: "Additional fields as JSON"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					body := map[string]any{}
					body["name"] = cmd.String("name")
					body["currency_code"] = cmd.String("currency_code")
					body["time_zone"] = cmd.String("time_zone")
					if cmd.IsSet("fiscal_year_start_month") {
						body["fiscal_year_start_month"] = cmd.String("fiscal_year_start_month")
					}
					if cmd.IsSet("date_format") {
						body["date_format"] = cmd.String("date_format")
					}
					if cmd.IsSet("language_code") {
						body["language_code"] = cmd.String("language_code")
					}
					if cmd.IsSet("industry_type") {
						body["industry_type"] = cmd.String("industry_type")
					}
					if cmd.IsSet("portal_name") {
						body["portal_name"] = cmd.String("portal_name")
					}
					if err := internal.MergeJSON(cmd, body); err != nil {
						return err
					}
					raw, err := c.Request("POST", c.BooksBase+"/organizations", &zohttp.RequestOpts{JSON: body})
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
					raw, err := c.Request("GET", c.BooksBase+"/organizations/"+cmd.Args().First(), nil)
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
					&cli.StringFlag{Name: "currency_code", Usage: "Currency code (e.g. USD)"},
					&cli.StringFlag{Name: "time_zone", Usage: "Time zone (e.g. PST)"},
					&cli.StringFlag{Name: "fiscal_year_start_month", Usage: "Fiscal year start month (e.g. january)"},
					&cli.StringFlag{Name: "date_format", Usage: "Date format (e.g. dd MMM yyyy)"},
					&cli.StringFlag{Name: "language_code", Usage: "Language code (e.g. en)"},
					&cli.StringFlag{Name: "industry_type", Usage: "Industry type"},
					&cli.StringFlag{Name: "portal_name", Usage: "Customer portal name"},
					&cli.StringFlag{Name: "json", Usage: "Additional fields as JSON"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					body := map[string]any{}
					if cmd.IsSet("name") {
						body["name"] = cmd.String("name")
					}
					if cmd.IsSet("currency_code") {
						body["currency_code"] = cmd.String("currency_code")
					}
					if cmd.IsSet("time_zone") {
						body["time_zone"] = cmd.String("time_zone")
					}
					if cmd.IsSet("fiscal_year_start_month") {
						body["fiscal_year_start_month"] = cmd.String("fiscal_year_start_month")
					}
					if cmd.IsSet("date_format") {
						body["date_format"] = cmd.String("date_format")
					}
					if cmd.IsSet("language_code") {
						body["language_code"] = cmd.String("language_code")
					}
					if cmd.IsSet("industry_type") {
						body["industry_type"] = cmd.String("industry_type")
					}
					if cmd.IsSet("portal_name") {
						body["portal_name"] = cmd.String("portal_name")
					}
					if err := internal.MergeJSON(cmd, body); err != nil {
						return err
					}
					raw, err := c.Request("PUT", c.BooksBase+"/organizations/"+cmd.Args().First(), &zohttp.RequestOpts{JSON: body})
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
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
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
					raw, err := c.Request("POST", c.BooksBase+"/contacts", &zohttp.RequestOpts{
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
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
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
					raw, err := c.Request("PUT", c.BooksBase+"/contacts", &zohttp.RequestOpts{
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
					params := orgParams(orgID)
					if v := cmd.String("status"); v != "" {
						params["status"] = v
					}
					if v := cmd.String("sort-column"); v != "" {
						params["sort_column"] = v
					}
					if v := cmd.String("page"); v != "" {
						params["page"] = v
					}
					if v := cmd.String("per-page"); v != "" {
						params["per_page"] = v
					}
					raw, err := c.Request("GET", c.BooksBase+"/contacts", &zohttp.RequestOpts{Params: params})
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
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
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
					raw, err := c.Request("PUT", c.BooksBase+"/contacts/"+cmd.Args().First(), &zohttp.RequestOpts{
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
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("GET", c.BooksBase+"/contacts/"+cmd.Args().First(), &zohttp.RequestOpts{Params: orgParams(orgID)})
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
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("DELETE", c.BooksBase+"/contacts/"+cmd.Args().First(), &zohttp.RequestOpts{Params: orgParams(orgID)})
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
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("POST", c.BooksBase+"/contacts/"+cmd.Args().First()+"/active", &zohttp.RequestOpts{Params: orgParams(orgID)})
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
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("POST", c.BooksBase+"/contacts/"+cmd.Args().First()+"/inactive", &zohttp.RequestOpts{Params: orgParams(orgID)})
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
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
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
					raw, err := c.Request("POST", c.BooksBase+"/contacts/"+cmd.Args().First()+"/portal/enable", &zohttp.RequestOpts{
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
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("POST", c.BooksBase+"/contacts/"+cmd.Args().First()+"/paymentreminder/enable", &zohttp.RequestOpts{Params: orgParams(orgID)})
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
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("POST", c.BooksBase+"/contacts/"+cmd.Args().First()+"/paymentreminder/disable", &zohttp.RequestOpts{Params: orgParams(orgID)})
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
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
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
					raw, err := c.Request("POST", c.BooksBase+"/contacts/"+cmd.Args().First()+"/statements/email", &zohttp.RequestOpts{
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
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("GET", c.BooksBase+"/contacts/"+cmd.Args().First()+"/statements/email", &zohttp.RequestOpts{Params: orgParams(orgID)})
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
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
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
					raw, err := c.Request("POST", c.BooksBase+"/contacts/"+cmd.Args().First()+"/email", &zohttp.RequestOpts{
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
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("GET", c.BooksBase+"/contacts/"+cmd.Args().First()+"/comments", &zohttp.RequestOpts{Params: orgParams(orgID)})
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
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
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
					raw, err := c.Request("POST", c.BooksBase+"/contacts/"+cmd.Args().First()+"/address", &zohttp.RequestOpts{
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
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("GET", c.BooksBase+"/contacts/"+cmd.Args().First()+"/address", &zohttp.RequestOpts{Params: orgParams(orgID)})
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
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
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
					raw, err := c.Request("PUT", c.BooksBase+"/contacts/"+cmd.Args().First()+"/address/"+cmd.Args().Get(1), &zohttp.RequestOpts{
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
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("DELETE", c.BooksBase+"/contacts/"+cmd.Args().First()+"/address/"+cmd.Args().Get(1), &zohttp.RequestOpts{Params: orgParams(orgID)})
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
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("GET", c.BooksBase+"/contacts/"+cmd.Args().First()+"/refunds", &zohttp.RequestOpts{Params: orgParams(orgID)})
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
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("POST", c.BooksBase+"/contacts/"+cmd.Args().First()+"/track1099", &zohttp.RequestOpts{Params: orgParams(orgID)})
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
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("POST", c.BooksBase+"/contacts/"+cmd.Args().First()+"/untrack1099", &zohttp.RequestOpts{Params: orgParams(orgID)})
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
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("GET", c.BooksBase+"/contacts/"+cmd.Args().First()+"/retainerpayments", &zohttp.RequestOpts{Params: orgParams(orgID)})
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
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
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
					raw, err := c.Request("POST", c.BooksBase+"/contacts/"+cmd.String("contact-id")+"/contactpersons", &zohttp.RequestOpts{
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
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
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
					raw, err := c.Request("PUT", c.BooksBase+"/contacts/"+cmd.String("contact-id")+"/contactpersons/"+cmd.Args().First(), &zohttp.RequestOpts{
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
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("DELETE", c.BooksBase+"/contacts/"+cmd.String("contact-id")+"/contactpersons/"+cmd.Args().First(), &zohttp.RequestOpts{Params: orgParams(orgID)})
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
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("GET", c.BooksBase+"/contacts/"+cmd.String("contact-id")+"/contactpersons", &zohttp.RequestOpts{Params: orgParams(orgID)})
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
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("GET", c.BooksBase+"/contacts/"+cmd.String("contact-id")+"/contactpersons/"+cmd.Args().First(), &zohttp.RequestOpts{Params: orgParams(orgID)})
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
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("POST", c.BooksBase+"/contacts/"+cmd.String("contact-id")+"/contactpersons/"+cmd.Args().First()+"/primary", &zohttp.RequestOpts{Params: orgParams(orgID)})
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
				Name:  "create",
				Usage: "Create an estimate",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "customer_id", Required: true, Usage: "Customer ID"},
					&cli.StringFlag{Name: "date", Usage: "Estimate date (YYYY-MM-DD)"},
					&cli.StringFlag{Name: "expiry_date", Usage: "Expiry date (YYYY-MM-DD)"},
					&cli.StringFlag{Name: "reference_number", Usage: "Reference number"},
					&cli.StringFlag{Name: "json", Usage: "Additional fields as JSON"},
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
					raw, err := c.Request("POST", c.BooksBase+"/estimates", &zohttp.RequestOpts{
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
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
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
					raw, err := c.Request("PUT", c.BooksBase+"/estimates", &zohttp.RequestOpts{
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
					&cli.StringFlag{Name: "page", Usage: "Filter"},
					&cli.StringFlag{Name: "per-page", Usage: "Filter"},
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
					if v := cmd.String("page"); v != "" {
						params["page"] = v
					}
					if v := cmd.String("per-page"); v != "" {
						params["per_page"] = v
					}
					raw, err := c.Request("GET", c.BooksBase+"/estimates", &zohttp.RequestOpts{Params: params})
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
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
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
					raw, err := c.Request("PUT", c.BooksBase+"/estimates/"+cmd.Args().First(), &zohttp.RequestOpts{
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
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("GET", c.BooksBase+"/estimates/"+cmd.Args().First(), &zohttp.RequestOpts{Params: orgParams(orgID)})
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
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("DELETE", c.BooksBase+"/estimates/"+cmd.Args().First(), &zohttp.RequestOpts{Params: orgParams(orgID)})
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
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
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
					raw, err := c.Request("PUT", c.BooksBase+"/estimates/"+cmd.Args().First()+"/customfields", &zohttp.RequestOpts{
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
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("POST", c.BooksBase+"/estimates/"+cmd.Args().First()+"/status/sent", &zohttp.RequestOpts{Params: orgParams(orgID)})
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
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("POST", c.BooksBase+"/estimates/"+cmd.Args().First()+"/status/accepted", &zohttp.RequestOpts{Params: orgParams(orgID)})
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
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("POST", c.BooksBase+"/estimates/"+cmd.Args().First()+"/status/declined", &zohttp.RequestOpts{Params: orgParams(orgID)})
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
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("POST", c.BooksBase+"/estimates/"+cmd.Args().First()+"/submit", &zohttp.RequestOpts{Params: orgParams(orgID)})
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
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("POST", c.BooksBase+"/estimates/"+cmd.Args().First()+"/approve", &zohttp.RequestOpts{Params: orgParams(orgID)})
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
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
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
					raw, err := c.Request("POST", c.BooksBase+"/estimates/"+cmd.Args().First()+"/email", &zohttp.RequestOpts{
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
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("GET", c.BooksBase+"/estimates/"+cmd.Args().First()+"/email", &zohttp.RequestOpts{Params: orgParams(orgID)})
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
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
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
					raw, err := c.Request("POST", c.BooksBase+"/estimates"+"/email", &zohttp.RequestOpts{
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
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
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
					raw, err := c.Request("PUT", c.BooksBase+"/estimates/"+cmd.Args().First()+"/address/billing", &zohttp.RequestOpts{
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
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
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
					raw, err := c.Request("PUT", c.BooksBase+"/estimates/"+cmd.Args().First()+"/address/shipping", &zohttp.RequestOpts{
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
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("GET", c.BooksBase+"/estimates"+"/templates", &zohttp.RequestOpts{Params: orgParams(orgID)})
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
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("PUT", c.BooksBase+"/estimates/"+cmd.Args().First()+"/templates/"+cmd.Args().Get(1), &zohttp.RequestOpts{Params: orgParams(orgID)})
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
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
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
					raw, err := c.Request("POST", c.BooksBase+"/estimates/"+cmd.Args().First()+"/comments", &zohttp.RequestOpts{
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
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("GET", c.BooksBase+"/estimates/"+cmd.Args().First()+"/comments", &zohttp.RequestOpts{Params: orgParams(orgID)})
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
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
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
					raw, err := c.Request("PUT", c.BooksBase+"/estimates/"+cmd.Args().First()+"/comments/"+cmd.Args().Get(1), &zohttp.RequestOpts{
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
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("DELETE", c.BooksBase+"/estimates/"+cmd.Args().First()+"/comments/"+cmd.Args().Get(1), &zohttp.RequestOpts{Params: orgParams(orgID)})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
		},
	}
}

func salesOrdersCmd() *cli.Command {
	return &cli.Command{
		Name:  "sales-orders",
		Usage: "Sales order operations",
		Commands: []*cli.Command{
			{
				Name:  "create",
				Usage: "Create a sales order",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "customer_id", Required: true, Usage: "Customer ID"},
					&cli.StringFlag{Name: "date", Usage: "Sales order date (YYYY-MM-DD)"},
					&cli.StringFlag{Name: "reference_number", Usage: "Reference number"},
					&cli.StringFlag{Name: "json", Usage: "Additional fields as JSON"},
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
					body := map[string]any{}
					body["customer_id"] = cmd.String("customer_id")
					if cmd.IsSet("date") {
						body["date"] = cmd.String("date")
					}
					if cmd.IsSet("reference_number") {
						body["reference_number"] = cmd.String("reference_number")
					}
					if err := internal.MergeJSON(cmd, body); err != nil {
						return err
					}
					raw, err := c.Request("POST", c.BooksBase+"/salesorders", &zohttp.RequestOpts{
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
				Usage: "Update a sales order by custom field",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "customer_id", Usage: "Customer ID"},
					&cli.StringFlag{Name: "date", Usage: "Sales order date (YYYY-MM-DD)"},
					&cli.StringFlag{Name: "reference_number", Usage: "Reference number"},
					&cli.StringFlag{Name: "json", Usage: "Additional fields as JSON"},
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
					body := map[string]any{}
					if cmd.IsSet("customer_id") {
						body["customer_id"] = cmd.String("customer_id")
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
					raw, err := c.Request("PUT", c.BooksBase+"/salesorders", &zohttp.RequestOpts{
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
				Usage: "List sales orders",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "status", Usage: "Filter"},
					&cli.StringFlag{Name: "customer-id", Usage: "Filter"},
					&cli.StringFlag{Name: "sort-column", Usage: "Filter"},
					&cli.StringFlag{Name: "page", Usage: "Filter"},
					&cli.StringFlag{Name: "per-page", Usage: "Filter"},
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
					if v := cmd.String("page"); v != "" {
						params["page"] = v
					}
					if v := cmd.String("per-page"); v != "" {
						params["per_page"] = v
					}
					raw, err := c.Request("GET", c.BooksBase+"/salesorders", &zohttp.RequestOpts{Params: params})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "update",
				Usage:     "Update a sales order",
				ArgsUsage: "<salesorder-id>",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "customer_id", Usage: "Customer ID"},
					&cli.StringFlag{Name: "date", Usage: "Sales order date (YYYY-MM-DD)"},
					&cli.StringFlag{Name: "reference_number", Usage: "Reference number"},
					&cli.StringFlag{Name: "json", Usage: "Additional fields as JSON"},
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
					body := map[string]any{}
					if cmd.IsSet("customer_id") {
						body["customer_id"] = cmd.String("customer_id")
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
					raw, err := c.Request("PUT", c.BooksBase+"/salesorders/"+cmd.Args().First(), &zohttp.RequestOpts{
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
				Usage:     "Get a sales order",
				ArgsUsage: "<salesorder-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("GET", c.BooksBase+"/salesorders/"+cmd.Args().First(), &zohttp.RequestOpts{Params: orgParams(orgID)})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "delete",
				Usage:     "Delete a sales order",
				ArgsUsage: "<salesorder-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("DELETE", c.BooksBase+"/salesorders/"+cmd.Args().First(), &zohttp.RequestOpts{Params: orgParams(orgID)})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "update-custom-fields",
				Usage:     "Update custom fields of a sales order",
				ArgsUsage: "<salesorder-id>",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "customfield_id", Usage: "Custom field ID"},
					&cli.StringFlag{Name: "value", Usage: "Custom field value"},
					&cli.StringFlag{Name: "json", Usage: "Additional fields as JSON"},
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
					raw, err := c.Request("PUT", c.BooksBase+"/salesorders/"+cmd.Args().First()+"/customfields", &zohttp.RequestOpts{
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
				Usage:     "Mark a sales order as open",
				ArgsUsage: "<salesorder-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("POST", c.BooksBase+"/salesorders/"+cmd.Args().First()+"/status/open", &zohttp.RequestOpts{Params: orgParams(orgID)})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "mark-void",
				Usage:     "Mark a sales order as void",
				ArgsUsage: "<salesorder-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("POST", c.BooksBase+"/salesorders/"+cmd.Args().First()+"/status/void", &zohttp.RequestOpts{Params: orgParams(orgID)})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "update-sub-status",
				Usage:     "Update sub-status of a sales order",
				ArgsUsage: "<salesorder-id> <sub-status-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("POST", c.BooksBase+"/salesorders/"+cmd.Args().First()+"/substatus/"+cmd.Args().Get(1), &zohttp.RequestOpts{Params: orgParams(orgID)})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "email",
				Usage:     "Email a sales order",
				ArgsUsage: "<salesorder-id>",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "to_mail_ids", Required: true, Usage: "Recipient email addresses (comma-separated)"},
					&cli.StringFlag{Name: "subject", Required: true, Usage: "Email subject"},
					&cli.StringFlag{Name: "body", Required: true, Usage: "Email body"},
					&cli.StringFlag{Name: "cc_mail_ids", Usage: "CC email addresses (comma-separated)"},
					&cli.StringFlag{Name: "json", Usage: "Additional fields as JSON"},
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
					raw, err := c.Request("POST", c.BooksBase+"/salesorders/"+cmd.Args().First()+"/email", &zohttp.RequestOpts{
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
				Usage:     "Get email content of a sales order",
				ArgsUsage: "<salesorder-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("GET", c.BooksBase+"/salesorders/"+cmd.Args().First()+"/email", &zohttp.RequestOpts{Params: orgParams(orgID)})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "submit-for-approval",
				Usage:     "Submit a sales order for approval",
				ArgsUsage: "<salesorder-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("POST", c.BooksBase+"/salesorders/"+cmd.Args().First()+"/submit", &zohttp.RequestOpts{Params: orgParams(orgID)})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "approve",
				Usage:     "Approve a sales order",
				ArgsUsage: "<salesorder-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("POST", c.BooksBase+"/salesorders/"+cmd.Args().First()+"/approve", &zohttp.RequestOpts{Params: orgParams(orgID)})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "update-billing-address",
				Usage:     "Update billing address of a sales order",
				ArgsUsage: "<salesorder-id>",
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
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
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
					raw, err := c.Request("PUT", c.BooksBase+"/salesorders/"+cmd.Args().First()+"/address/billing", &zohttp.RequestOpts{
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
				Usage:     "Update shipping address of a sales order",
				ArgsUsage: "<salesorder-id>",
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
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
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
					raw, err := c.Request("PUT", c.BooksBase+"/salesorders/"+cmd.Args().First()+"/address/shipping", &zohttp.RequestOpts{
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
				Usage: "List sales order templates",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("GET", c.BooksBase+"/salesorders"+"/templates", &zohttp.RequestOpts{Params: orgParams(orgID)})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "update-template",
				Usage:     "Update template of a sales order",
				ArgsUsage: "<salesorder-id> <template-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("PUT", c.BooksBase+"/salesorders/"+cmd.Args().First()+"/templates/"+cmd.Args().Get(1), &zohttp.RequestOpts{Params: orgParams(orgID)})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "add-attachment",
				Usage:     "Add attachment to a sales order",
				ArgsUsage: "<salesorder-id>",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "file", Required: true, Usage: "File path to upload"},
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
					data, err := os.ReadFile(cmd.String("file"))
					if err != nil {
						return fmt.Errorf("reading file: %w", err)
					}
					raw, err := c.Request("POST", c.BooksBase+"/salesorders/"+cmd.Args().First()+"/attachment", &zohttp.RequestOpts{
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
				Usage:     "Update attachment preference of a sales order",
				ArgsUsage: "<salesorder-id>",
				Flags: []cli.Flag{
					&cli.BoolFlag{Name: "can_send_in_mail", Usage: "Attach to email (true/false)"},
					&cli.StringFlag{Name: "json", Usage: "Additional fields as JSON"},
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
					body := map[string]any{}
					if cmd.IsSet("can_send_in_mail") {
						body["can_send_in_mail"] = cmd.Bool("can_send_in_mail")
					}
					if err := internal.MergeJSON(cmd, body); err != nil {
						return err
					}
					raw, err := c.Request("PUT", c.BooksBase+"/salesorders/"+cmd.Args().First()+"/attachment", &zohttp.RequestOpts{
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
				Usage:     "Get attachment of a sales order",
				ArgsUsage: "<salesorder-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("GET", c.BooksBase+"/salesorders/"+cmd.Args().First()+"/attachment", &zohttp.RequestOpts{Params: orgParams(orgID)})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "delete-attachment",
				Usage:     "Delete attachment of a sales order",
				ArgsUsage: "<salesorder-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("DELETE", c.BooksBase+"/salesorders/"+cmd.Args().First()+"/attachment", &zohttp.RequestOpts{Params: orgParams(orgID)})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "add-comment",
				Usage:     "Add a comment to a sales order",
				ArgsUsage: "<salesorder-id>",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "description", Required: true, Usage: "Comment text"},
					&cli.StringFlag{Name: "json", Usage: "Additional fields as JSON"},
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
					body := map[string]any{}
					body["description"] = cmd.String("description")
					if err := internal.MergeJSON(cmd, body); err != nil {
						return err
					}
					raw, err := c.Request("POST", c.BooksBase+"/salesorders/"+cmd.Args().First()+"/comments", &zohttp.RequestOpts{
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
				Usage:     "List comments of a sales order",
				ArgsUsage: "<salesorder-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("GET", c.BooksBase+"/salesorders/"+cmd.Args().First()+"/comments", &zohttp.RequestOpts{Params: orgParams(orgID)})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "update-comment",
				Usage:     "Update a comment on a sales order",
				ArgsUsage: "<salesorder-id> <comment-id>",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "description", Required: true, Usage: "Comment text"},
					&cli.StringFlag{Name: "json", Usage: "Additional fields as JSON"},
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
					body := map[string]any{}
					body["description"] = cmd.String("description")
					if err := internal.MergeJSON(cmd, body); err != nil {
						return err
					}
					raw, err := c.Request("PUT", c.BooksBase+"/salesorders/"+cmd.Args().First()+"/comments/"+cmd.Args().Get(1), &zohttp.RequestOpts{
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
				Usage:     "Delete a comment on a sales order",
				ArgsUsage: "<salesorder-id> <comment-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("DELETE", c.BooksBase+"/salesorders/"+cmd.Args().First()+"/comments/"+cmd.Args().Get(1), &zohttp.RequestOpts{Params: orgParams(orgID)})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
		},
	}
}

func salesReceiptsCmd() *cli.Command {
	return &cli.Command{
		Name:  "sales-receipts",
		Usage: "Sales receipt operations",
		Commands: []*cli.Command{
			{
				Name:  "create",
				Usage: "Create a sales receipt",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "customer_id", Required: true, Usage: "Customer ID"},
					&cli.StringFlag{Name: "payment_mode", Usage: "Payment mode"},
					&cli.StringFlag{Name: "date", Usage: "Receipt date (YYYY-MM-DD)"},
					&cli.StringFlag{Name: "receipt_number", Usage: "Receipt number"},
					&cli.StringFlag{Name: "json", Usage: "Additional fields as JSON"},
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
					body := map[string]any{}
					body["customer_id"] = cmd.String("customer_id")
					if cmd.IsSet("payment_mode") {
						body["payment_mode"] = cmd.String("payment_mode")
					}
					if cmd.IsSet("date") {
						body["date"] = cmd.String("date")
					}
					if cmd.IsSet("receipt_number") {
						body["receipt_number"] = cmd.String("receipt_number")
					}
					if err := internal.MergeJSON(cmd, body); err != nil {
						return err
					}
					raw, err := c.Request("POST", c.BooksBase+"/salesreceipts", &zohttp.RequestOpts{
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
				Usage: "List sales receipts",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "page", Usage: "Filter"},
					&cli.StringFlag{Name: "per-page", Usage: "Filter"},
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
					params := orgParams(orgID)
					if v := cmd.String("page"); v != "" {
						params["page"] = v
					}
					if v := cmd.String("per-page"); v != "" {
						params["per_page"] = v
					}
					raw, err := c.Request("GET", c.BooksBase+"/salesreceipts", &zohttp.RequestOpts{Params: params})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "update",
				Usage:     "Update a sales receipt",
				ArgsUsage: "<salesreceipt-id>",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "customer_id", Usage: "Customer ID"},
					&cli.StringFlag{Name: "payment_mode", Usage: "Payment mode"},
					&cli.StringFlag{Name: "date", Usage: "Receipt date (YYYY-MM-DD)"},
					&cli.StringFlag{Name: "receipt_number", Usage: "Receipt number"},
					&cli.StringFlag{Name: "json", Usage: "Additional fields as JSON"},
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
					body := map[string]any{}
					if cmd.IsSet("customer_id") {
						body["customer_id"] = cmd.String("customer_id")
					}
					if cmd.IsSet("payment_mode") {
						body["payment_mode"] = cmd.String("payment_mode")
					}
					if cmd.IsSet("date") {
						body["date"] = cmd.String("date")
					}
					if cmd.IsSet("receipt_number") {
						body["receipt_number"] = cmd.String("receipt_number")
					}
					if err := internal.MergeJSON(cmd, body); err != nil {
						return err
					}
					raw, err := c.Request("PUT", c.BooksBase+"/salesreceipts/"+cmd.Args().First(), &zohttp.RequestOpts{
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
				Usage:     "Get a sales receipt",
				ArgsUsage: "<salesreceipt-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("GET", c.BooksBase+"/salesreceipts/"+cmd.Args().First(), &zohttp.RequestOpts{Params: orgParams(orgID)})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "delete",
				Usage:     "Delete a sales receipt",
				ArgsUsage: "<salesreceipt-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("DELETE", c.BooksBase+"/salesreceipts/"+cmd.Args().First(), &zohttp.RequestOpts{Params: orgParams(orgID)})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "email",
				Usage:     "Email a sales receipt",
				ArgsUsage: "<salesreceipt-id>",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "to_mail_ids", Required: true, Usage: "Recipient email addresses (comma-separated)"},
					&cli.StringFlag{Name: "subject", Required: true, Usage: "Email subject"},
					&cli.StringFlag{Name: "body", Required: true, Usage: "Email body"},
					&cli.StringFlag{Name: "json", Usage: "Additional fields as JSON"},
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
					body := map[string]any{}
					body["to_mail_ids"] = cmd.String("to_mail_ids")
					body["subject"] = cmd.String("subject")
					body["body"] = cmd.String("body")
					if err := internal.MergeJSON(cmd, body); err != nil {
						return err
					}
					raw, err := c.Request("POST", c.BooksBase+"/salesreceipts/"+cmd.Args().First()+"/email", &zohttp.RequestOpts{
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

func invoicesCmd() *cli.Command {
	return &cli.Command{
		Name:  "invoices",
		Usage: "Invoice operations",
		Commands: []*cli.Command{
			{
				Name:  "create",
				Usage: "Create an invoice",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "customer_id", Required: true, Usage: "Customer ID"},
					&cli.StringFlag{Name: "date", Usage: "Invoice date (YYYY-MM-DD)"},
					&cli.StringFlag{Name: "due_date", Usage: "Due date (YYYY-MM-DD)"},
					&cli.StringFlag{Name: "reference_number", Usage: "Reference number"},
					&cli.StringFlag{Name: "json", Usage: "Additional fields as JSON"},
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
					body := map[string]any{}
					body["customer_id"] = cmd.String("customer_id")
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
					raw, err := c.Request("POST", c.BooksBase+"/invoices", &zohttp.RequestOpts{
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
				Usage: "Update an invoice by custom field",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "customer_id", Usage: "Customer ID"},
					&cli.StringFlag{Name: "date", Usage: "Invoice date (YYYY-MM-DD)"},
					&cli.StringFlag{Name: "due_date", Usage: "Due date (YYYY-MM-DD)"},
					&cli.StringFlag{Name: "reference_number", Usage: "Reference number"},
					&cli.StringFlag{Name: "json", Usage: "Additional fields as JSON"},
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
					body := map[string]any{}
					if cmd.IsSet("customer_id") {
						body["customer_id"] = cmd.String("customer_id")
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
					raw, err := c.Request("PUT", c.BooksBase+"/invoices", &zohttp.RequestOpts{
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
				Usage: "List invoices",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "status", Usage: "Filter"},
					&cli.StringFlag{Name: "customer-id", Usage: "Filter"},
					&cli.StringFlag{Name: "sort-column", Usage: "Filter"},
					&cli.StringFlag{Name: "date", Usage: "Filter"},
					&cli.StringFlag{Name: "due-date", Usage: "Filter"},
					&cli.StringFlag{Name: "page", Usage: "Filter"},
					&cli.StringFlag{Name: "per-page", Usage: "Filter"},
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
					if v := cmd.String("date"); v != "" {
						params["date"] = v
					}
					if v := cmd.String("due-date"); v != "" {
						params["due_date"] = v
					}
					if v := cmd.String("page"); v != "" {
						params["page"] = v
					}
					if v := cmd.String("per-page"); v != "" {
						params["per_page"] = v
					}
					raw, err := c.Request("GET", c.BooksBase+"/invoices", &zohttp.RequestOpts{Params: params})
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
					&cli.StringFlag{Name: "date", Usage: "Invoice date (YYYY-MM-DD)"},
					&cli.StringFlag{Name: "due_date", Usage: "Due date (YYYY-MM-DD)"},
					&cli.StringFlag{Name: "reference_number", Usage: "Reference number"},
					&cli.StringFlag{Name: "json", Usage: "Additional fields as JSON"},
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
					body := map[string]any{}
					if cmd.IsSet("customer_id") {
						body["customer_id"] = cmd.String("customer_id")
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
					raw, err := c.Request("PUT", c.BooksBase+"/invoices/"+cmd.Args().First(), &zohttp.RequestOpts{
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
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("GET", c.BooksBase+"/invoices/"+cmd.Args().First(), &zohttp.RequestOpts{Params: orgParams(orgID)})
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
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("DELETE", c.BooksBase+"/invoices/"+cmd.Args().First(), &zohttp.RequestOpts{Params: orgParams(orgID)})
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
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("POST", c.BooksBase+"/invoices/"+cmd.Args().First()+"/status/sent", &zohttp.RequestOpts{Params: orgParams(orgID)})
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
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("POST", c.BooksBase+"/invoices/"+cmd.Args().First()+"/status/void", &zohttp.RequestOpts{Params: orgParams(orgID)})
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
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("POST", c.BooksBase+"/invoices/"+cmd.Args().First()+"/status/draft", &zohttp.RequestOpts{Params: orgParams(orgID)})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:  "email-multiple",
				Usage: "Email multiple invoices",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "invoice_ids", Required: true, Usage: "Invoice IDs (comma-separated)"},
					&cli.StringFlag{Name: "json", Usage: "Additional fields as JSON"},
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
					body := map[string]any{}
					var parsed0 any
					if err := json.Unmarshal([]byte(cmd.String("invoice_ids")), &parsed0); err != nil {
						return err
					}
					body["invoice_ids"] = parsed0
					if err := internal.MergeJSON(cmd, body); err != nil {
						return err
					}
					raw, err := c.Request("POST", c.BooksBase+"/invoices"+"/email", &zohttp.RequestOpts{
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
				Name:  "create-instant",
				Usage: "Create an instant invoice",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "customer_id", Required: true, Usage: "Customer ID"},
					&cli.StringFlag{Name: "date", Usage: "Invoice date (YYYY-MM-DD)"},
					&cli.StringFlag{Name: "due_date", Usage: "Due date (YYYY-MM-DD)"},
					&cli.StringFlag{Name: "json", Usage: "Additional fields as JSON"},
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
					body := map[string]any{}
					body["customer_id"] = cmd.String("customer_id")
					if cmd.IsSet("date") {
						body["date"] = cmd.String("date")
					}
					if cmd.IsSet("due_date") {
						body["due_date"] = cmd.String("due_date")
					}
					if err := internal.MergeJSON(cmd, body); err != nil {
						return err
					}
					raw, err := c.Request("POST", c.BooksBase+"/invoices"+"/instant", &zohttp.RequestOpts{
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
				Name:      "associate-salesorder",
				Usage:     "Associate a sales order with an invoice",
				ArgsUsage: "<invoice-id>",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "salesorder_id", Required: true, Usage: "Sales order ID to associate"},
					&cli.StringFlag{Name: "json", Usage: "Additional fields as JSON"},
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
					body := map[string]any{}
					body["salesorder_id"] = cmd.String("salesorder_id")
					if err := internal.MergeJSON(cmd, body); err != nil {
						return err
					}
					raw, err := c.Request("POST", c.BooksBase+"/invoices/"+cmd.Args().First()+"/salesorders", &zohttp.RequestOpts{
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
				Name:      "submit-for-approval",
				Usage:     "Submit an invoice for approval",
				ArgsUsage: "<invoice-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("POST", c.BooksBase+"/invoices/"+cmd.Args().First()+"/submit", &zohttp.RequestOpts{Params: orgParams(orgID)})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "approve",
				Usage:     "Approve an invoice",
				ArgsUsage: "<invoice-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("POST", c.BooksBase+"/invoices/"+cmd.Args().First()+"/approve", &zohttp.RequestOpts{Params: orgParams(orgID)})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "email",
				Usage:     "Email an invoice",
				ArgsUsage: "<invoice-id>",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "to_mail_ids", Required: true, Usage: "Recipient email addresses (comma-separated)"},
					&cli.StringFlag{Name: "subject", Required: true, Usage: "Email subject"},
					&cli.StringFlag{Name: "body", Required: true, Usage: "Email body"},
					&cli.StringFlag{Name: "cc_mail_ids", Usage: "CC email addresses (comma-separated)"},
					&cli.StringFlag{Name: "json", Usage: "Additional fields as JSON"},
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
					raw, err := c.Request("POST", c.BooksBase+"/invoices/"+cmd.Args().First()+"/email", &zohttp.RequestOpts{
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
				Usage:     "Get email content of an invoice",
				ArgsUsage: "<invoice-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("GET", c.BooksBase+"/invoices/"+cmd.Args().First()+"/email", &zohttp.RequestOpts{Params: orgParams(orgID)})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "remind-customer",
				Usage:     "Send payment reminder for an invoice",
				ArgsUsage: "<invoice-id>",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "to_mail_ids", Usage: "Recipient email addresses (comma-separated)"},
					&cli.StringFlag{Name: "subject", Usage: "Email subject"},
					&cli.StringFlag{Name: "body", Usage: "Email body"},
					&cli.StringFlag{Name: "cc_mail_ids", Usage: "CC email addresses (comma-separated)"},
					&cli.StringFlag{Name: "json", Usage: "Additional fields as JSON"},
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
					body := map[string]any{}
					if cmd.IsSet("to_mail_ids") {
						body["to_mail_ids"] = cmd.String("to_mail_ids")
					}
					if cmd.IsSet("subject") {
						body["subject"] = cmd.String("subject")
					}
					if cmd.IsSet("body") {
						body["body"] = cmd.String("body")
					}
					if cmd.IsSet("cc_mail_ids") {
						body["cc_mail_ids"] = cmd.String("cc_mail_ids")
					}
					if err := internal.MergeJSON(cmd, body); err != nil {
						return err
					}
					opts := &zohttp.RequestOpts{Params: orgParams(orgID), JSON: body}
					raw, err := c.Request("POST", c.BooksBase+"/invoices/"+cmd.Args().First()+"/paymentreminder", opts)
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "get-payment-reminder-content",
				Usage:     "Get payment reminder content of an invoice",
				ArgsUsage: "<invoice-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("GET", c.BooksBase+"/invoices/"+cmd.Args().First()+"/paymentreminder", &zohttp.RequestOpts{Params: orgParams(orgID)})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:  "bulk-reminder",
				Usage: "Send bulk payment reminders",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "invoice_ids", Required: true, Usage: "Invoice IDs (comma-separated)"},
					&cli.StringFlag{Name: "json", Usage: "Additional fields as JSON"},
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
					body := map[string]any{}
					var parsed0 any
					if err := json.Unmarshal([]byte(cmd.String("invoice_ids")), &parsed0); err != nil {
						return err
					}
					body["invoice_ids"] = parsed0
					if err := internal.MergeJSON(cmd, body); err != nil {
						return err
					}
					raw, err := c.Request("POST", c.BooksBase+"/invoices"+"/paymentreminder", &zohttp.RequestOpts{
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
				Name:      "disable-payment-reminder",
				Usage:     "Disable payment reminder for an invoice",
				ArgsUsage: "<invoice-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("POST", c.BooksBase+"/invoices/"+cmd.Args().First()+"/paymentreminder/disable", &zohttp.RequestOpts{Params: orgParams(orgID)})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "enable-payment-reminder",
				Usage:     "Enable payment reminder for an invoice",
				ArgsUsage: "<invoice-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("POST", c.BooksBase+"/invoices/"+cmd.Args().First()+"/paymentreminder/enable", &zohttp.RequestOpts{Params: orgParams(orgID)})
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
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("POST", c.BooksBase+"/invoices/"+cmd.Args().First()+"/writeoff", &zohttp.RequestOpts{Params: orgParams(orgID)})
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
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("POST", c.BooksBase+"/invoices/"+cmd.Args().First()+"/writeoff/cancel", &zohttp.RequestOpts{Params: orgParams(orgID)})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "update-billing-address",
				Usage:     "Update billing address of an invoice",
				ArgsUsage: "<invoice-id>",
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
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
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
					raw, err := c.Request("PUT", c.BooksBase+"/invoices/"+cmd.Args().First()+"/address/billing", &zohttp.RequestOpts{
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
				Usage:     "Update shipping address of an invoice",
				ArgsUsage: "<invoice-id>",
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
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
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
					raw, err := c.Request("PUT", c.BooksBase+"/invoices/"+cmd.Args().First()+"/address/shipping", &zohttp.RequestOpts{
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
				Usage: "List invoice templates",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("GET", c.BooksBase+"/invoices"+"/templates", &zohttp.RequestOpts{Params: orgParams(orgID)})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "update-template",
				Usage:     "Update template of an invoice",
				ArgsUsage: "<invoice-id> <template-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("PUT", c.BooksBase+"/invoices/"+cmd.Args().First()+"/templates/"+cmd.Args().Get(1), &zohttp.RequestOpts{Params: orgParams(orgID)})
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
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("GET", c.BooksBase+"/invoices/"+cmd.Args().First()+"/payments", &zohttp.RequestOpts{Params: orgParams(orgID)})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "list-credits-applied",
				Usage:     "List credits applied to an invoice",
				ArgsUsage: "<invoice-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("GET", c.BooksBase+"/invoices/"+cmd.Args().First()+"/creditsapplied", &zohttp.RequestOpts{Params: orgParams(orgID)})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "apply-credits",
				Usage:     "Apply credits to an invoice",
				ArgsUsage: "<invoice-id>",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "invoice_id", Usage: "Invoice ID"},
					&cli.FloatFlag{Name: "amount_applied", Usage: "Amount to apply"},
					&cli.StringFlag{Name: "json", Usage: "Additional fields as JSON"},
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
					body := map[string]any{}
					if cmd.IsSet("invoice_id") {
						body["invoice_id"] = cmd.String("invoice_id")
					}
					if cmd.IsSet("amount_applied") {
						body["amount_applied"] = cmd.Float("amount_applied")
					}
					if err := internal.MergeJSON(cmd, body); err != nil {
						return err
					}
					raw, err := c.Request("POST", c.BooksBase+"/invoices/"+cmd.Args().First()+"/credits", &zohttp.RequestOpts{
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
				Usage:     "Delete a payment from an invoice",
				ArgsUsage: "<invoice-id> <payment-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("DELETE", c.BooksBase+"/invoices/"+cmd.Args().First()+"/payments/"+cmd.Args().Get(1), &zohttp.RequestOpts{Params: orgParams(orgID)})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "delete-applied-credit",
				Usage:     "Delete an applied credit from an invoice",
				ArgsUsage: "<invoice-id> <credit-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("DELETE", c.BooksBase+"/invoices/"+cmd.Args().First()+"/creditsapplied/"+cmd.Args().Get(1), &zohttp.RequestOpts{Params: orgParams(orgID)})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "add-attachment",
				Usage:     "Add attachment to an invoice",
				ArgsUsage: "<invoice-id>",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "file", Required: true, Usage: "File path to upload"},
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
					data, err := os.ReadFile(cmd.String("file"))
					if err != nil {
						return fmt.Errorf("reading file: %w", err)
					}
					raw, err := c.Request("POST", c.BooksBase+"/invoices/"+cmd.Args().First()+"/attachment", &zohttp.RequestOpts{
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
				Usage:     "Update attachment preference of an invoice",
				ArgsUsage: "<invoice-id>",
				Flags: []cli.Flag{
					&cli.BoolFlag{Name: "can_send_in_mail", Usage: "Attach to email (true/false)"},
					&cli.StringFlag{Name: "json", Usage: "Additional fields as JSON"},
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
					body := map[string]any{}
					if cmd.IsSet("can_send_in_mail") {
						body["can_send_in_mail"] = cmd.Bool("can_send_in_mail")
					}
					if err := internal.MergeJSON(cmd, body); err != nil {
						return err
					}
					raw, err := c.Request("PUT", c.BooksBase+"/invoices/"+cmd.Args().First()+"/attachment", &zohttp.RequestOpts{
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
				Usage:     "Get attachment of an invoice",
				ArgsUsage: "<invoice-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("GET", c.BooksBase+"/invoices/"+cmd.Args().First()+"/attachment", &zohttp.RequestOpts{Params: orgParams(orgID)})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "delete-attachment",
				Usage:     "Delete attachment of an invoice",
				ArgsUsage: "<invoice-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("DELETE", c.BooksBase+"/invoices/"+cmd.Args().First()+"/attachment", &zohttp.RequestOpts{Params: orgParams(orgID)})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "retrieve-document",
				Usage:     "Retrieve documents of an invoice",
				ArgsUsage: "<invoice-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("GET", c.BooksBase+"/invoices/"+cmd.Args().First()+"/documents", &zohttp.RequestOpts{Params: orgParams(orgID)})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "delete-invoice-attachment",
				Usage:     "Delete an invoice document attachment",
				ArgsUsage: "<invoice-id> <document-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("DELETE", c.BooksBase+"/invoices/"+cmd.Args().First()+"/documents/"+cmd.Args().Get(1), &zohttp.RequestOpts{Params: orgParams(orgID)})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "delete-expense-receipt",
				Usage:     "Delete an expense receipt from an invoice",
				ArgsUsage: "<invoice-id> <expense-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("DELETE", c.BooksBase+"/invoices/"+cmd.Args().First()+"/receipt/"+cmd.Args().Get(1), &zohttp.RequestOpts{Params: orgParams(orgID)})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "update-custom-fields",
				Usage:     "Update custom fields of an invoice",
				ArgsUsage: "<invoice-id>",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "customfield_id", Usage: "Custom field ID"},
					&cli.StringFlag{Name: "value", Usage: "Custom field value"},
					&cli.StringFlag{Name: "json", Usage: "Additional fields as JSON"},
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
					raw, err := c.Request("PUT", c.BooksBase+"/invoices/"+cmd.Args().First()+"/customfields", &zohttp.RequestOpts{
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
				Name:      "add-comment",
				Usage:     "Add a comment to an invoice",
				ArgsUsage: "<invoice-id>",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "description", Required: true, Usage: "Comment text"},
					&cli.StringFlag{Name: "json", Usage: "Additional fields as JSON"},
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
					body := map[string]any{}
					body["description"] = cmd.String("description")
					if err := internal.MergeJSON(cmd, body); err != nil {
						return err
					}
					raw, err := c.Request("POST", c.BooksBase+"/invoices/"+cmd.Args().First()+"/comments", &zohttp.RequestOpts{
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
				Usage:     "List comments of an invoice",
				ArgsUsage: "<invoice-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("GET", c.BooksBase+"/invoices/"+cmd.Args().First()+"/comments", &zohttp.RequestOpts{Params: orgParams(orgID)})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "update-comment",
				Usage:     "Update a comment on an invoice",
				ArgsUsage: "<invoice-id> <comment-id>",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "description", Required: true, Usage: "Comment text"},
					&cli.StringFlag{Name: "json", Usage: "Additional fields as JSON"},
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
					body := map[string]any{}
					body["description"] = cmd.String("description")
					if err := internal.MergeJSON(cmd, body); err != nil {
						return err
					}
					raw, err := c.Request("PUT", c.BooksBase+"/invoices/"+cmd.Args().First()+"/comments/"+cmd.Args().Get(1), &zohttp.RequestOpts{
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
				Usage:     "Delete a comment on an invoice",
				ArgsUsage: "<invoice-id> <comment-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("DELETE", c.BooksBase+"/invoices/"+cmd.Args().First()+"/comments/"+cmd.Args().Get(1), &zohttp.RequestOpts{Params: orgParams(orgID)})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "generate-payment-link",
				Usage:     "Generate payment link for an invoice",
				ArgsUsage: "<invoice-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("POST", c.BooksBase+"/invoices/"+cmd.Args().First()+"/paymentlink", &zohttp.RequestOpts{Params: orgParams(orgID)})
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
				Name:  "create",
				Usage: "Create a recurring invoice",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "customer_id", Required: true, Usage: "Customer ID"},
					&cli.StringFlag{Name: "recurrence_name", Usage: "Recurrence name"},
					&cli.StringFlag{Name: "start_date", Usage: "Start date (YYYY-MM-DD)"},
					&cli.StringFlag{Name: "end_date", Usage: "End date (YYYY-MM-DD)"},
					&cli.StringFlag{Name: "recurrence_frequency", Usage: "Recurrence frequency"},
					&cli.IntFlag{Name: "repeat_every", Usage: "Repeat interval"},
					&cli.StringFlag{Name: "json", Usage: "Additional fields as JSON"},
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
					body := map[string]any{}
					body["customer_id"] = cmd.String("customer_id")
					if cmd.IsSet("recurrence_name") {
						body["recurrence_name"] = cmd.String("recurrence_name")
					}
					if cmd.IsSet("start_date") {
						body["start_date"] = cmd.String("start_date")
					}
					if cmd.IsSet("end_date") {
						body["end_date"] = cmd.String("end_date")
					}
					if cmd.IsSet("recurrence_frequency") {
						body["recurrence_frequency"] = cmd.String("recurrence_frequency")
					}
					if cmd.IsSet("repeat_every") {
						body["repeat_every"] = cmd.Int("repeat_every")
					}
					if err := internal.MergeJSON(cmd, body); err != nil {
						return err
					}
					raw, err := c.Request("POST", c.BooksBase+"/recurringinvoices", &zohttp.RequestOpts{
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
				Usage: "Update a recurring invoice by custom field",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "customer_id", Usage: "Customer ID"},
					&cli.StringFlag{Name: "recurrence_name", Usage: "Recurrence name"},
					&cli.StringFlag{Name: "start_date", Usage: "Start date (YYYY-MM-DD)"},
					&cli.StringFlag{Name: "end_date", Usage: "End date (YYYY-MM-DD)"},
					&cli.StringFlag{Name: "recurrence_frequency", Usage: "Recurrence frequency"},
					&cli.IntFlag{Name: "repeat_every", Usage: "Repeat interval"},
					&cli.StringFlag{Name: "json", Usage: "Additional fields as JSON"},
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
					body := map[string]any{}
					if cmd.IsSet("customer_id") {
						body["customer_id"] = cmd.String("customer_id")
					}
					if cmd.IsSet("recurrence_name") {
						body["recurrence_name"] = cmd.String("recurrence_name")
					}
					if cmd.IsSet("start_date") {
						body["start_date"] = cmd.String("start_date")
					}
					if cmd.IsSet("end_date") {
						body["end_date"] = cmd.String("end_date")
					}
					if cmd.IsSet("recurrence_frequency") {
						body["recurrence_frequency"] = cmd.String("recurrence_frequency")
					}
					if cmd.IsSet("repeat_every") {
						body["repeat_every"] = cmd.Int("repeat_every")
					}
					if err := internal.MergeJSON(cmd, body); err != nil {
						return err
					}
					raw, err := c.Request("PUT", c.BooksBase+"/recurringinvoices", &zohttp.RequestOpts{
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
				Usage: "List recurring invoices",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "status", Usage: "Filter"},
					&cli.StringFlag{Name: "customer-id", Usage: "Filter"},
					&cli.StringFlag{Name: "page", Usage: "Filter"},
					&cli.StringFlag{Name: "per-page", Usage: "Filter"},
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
					params := orgParams(orgID)
					if v := cmd.String("status"); v != "" {
						params["status"] = v
					}
					if v := cmd.String("customer-id"); v != "" {
						params["customer_id"] = v
					}
					if v := cmd.String("page"); v != "" {
						params["page"] = v
					}
					if v := cmd.String("per-page"); v != "" {
						params["per_page"] = v
					}
					raw, err := c.Request("GET", c.BooksBase+"/recurringinvoices", &zohttp.RequestOpts{Params: params})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "update",
				Usage:     "Update a recurring invoice",
				ArgsUsage: "<recurringinvoice-id>",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "customer_id", Usage: "Customer ID"},
					&cli.StringFlag{Name: "recurrence_name", Usage: "Recurrence name"},
					&cli.StringFlag{Name: "start_date", Usage: "Start date (YYYY-MM-DD)"},
					&cli.StringFlag{Name: "end_date", Usage: "End date (YYYY-MM-DD)"},
					&cli.StringFlag{Name: "recurrence_frequency", Usage: "Recurrence frequency"},
					&cli.IntFlag{Name: "repeat_every", Usage: "Repeat interval"},
					&cli.StringFlag{Name: "json", Usage: "Additional fields as JSON"},
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
					body := map[string]any{}
					if cmd.IsSet("customer_id") {
						body["customer_id"] = cmd.String("customer_id")
					}
					if cmd.IsSet("recurrence_name") {
						body["recurrence_name"] = cmd.String("recurrence_name")
					}
					if cmd.IsSet("start_date") {
						body["start_date"] = cmd.String("start_date")
					}
					if cmd.IsSet("end_date") {
						body["end_date"] = cmd.String("end_date")
					}
					if cmd.IsSet("recurrence_frequency") {
						body["recurrence_frequency"] = cmd.String("recurrence_frequency")
					}
					if cmd.IsSet("repeat_every") {
						body["repeat_every"] = cmd.Int("repeat_every")
					}
					if err := internal.MergeJSON(cmd, body); err != nil {
						return err
					}
					raw, err := c.Request("PUT", c.BooksBase+"/recurringinvoices/"+cmd.Args().First(), &zohttp.RequestOpts{
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
				ArgsUsage: "<recurringinvoice-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("GET", c.BooksBase+"/recurringinvoices/"+cmd.Args().First(), &zohttp.RequestOpts{Params: orgParams(orgID)})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "delete",
				Usage:     "Delete a recurring invoice",
				ArgsUsage: "<recurringinvoice-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("DELETE", c.BooksBase+"/recurringinvoices/"+cmd.Args().First(), &zohttp.RequestOpts{Params: orgParams(orgID)})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "stop",
				Usage:     "Stop a recurring invoice",
				ArgsUsage: "<recurringinvoice-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("POST", c.BooksBase+"/recurringinvoices/"+cmd.Args().First()+"/status/stop", &zohttp.RequestOpts{Params: orgParams(orgID)})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "resume",
				Usage:     "Resume a recurring invoice",
				ArgsUsage: "<recurringinvoice-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("POST", c.BooksBase+"/recurringinvoices/"+cmd.Args().First()+"/status/resume", &zohttp.RequestOpts{Params: orgParams(orgID)})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "update-template",
				Usage:     "Update template of a recurring invoice",
				ArgsUsage: "<recurringinvoice-id> <template-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("PUT", c.BooksBase+"/recurringinvoices/"+cmd.Args().First()+"/templates/"+cmd.Args().Get(1), &zohttp.RequestOpts{Params: orgParams(orgID)})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "list-history",
				Usage:     "List history of a recurring invoice",
				ArgsUsage: "<recurringinvoice-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("GET", c.BooksBase+"/recurringinvoices/"+cmd.Args().First()+"/comments", &zohttp.RequestOpts{Params: orgParams(orgID)})
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
				Name:  "create",
				Usage: "Create a credit note",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "customer_id", Required: true, Usage: "Customer ID"},
					&cli.StringFlag{Name: "date", Usage: "Credit note date (YYYY-MM-DD)"},
					&cli.StringFlag{Name: "reference_number", Usage: "Reference number"},
					&cli.StringFlag{Name: "json", Usage: "Additional fields as JSON"},
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
					body := map[string]any{}
					body["customer_id"] = cmd.String("customer_id")
					if cmd.IsSet("date") {
						body["date"] = cmd.String("date")
					}
					if cmd.IsSet("reference_number") {
						body["reference_number"] = cmd.String("reference_number")
					}
					if err := internal.MergeJSON(cmd, body); err != nil {
						return err
					}
					raw, err := c.Request("POST", c.BooksBase+"/creditnotes", &zohttp.RequestOpts{
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
				Usage: "Update a credit note by custom field",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "customer_id", Usage: "Customer ID"},
					&cli.StringFlag{Name: "date", Usage: "Credit note date (YYYY-MM-DD)"},
					&cli.StringFlag{Name: "reference_number", Usage: "Reference number"},
					&cli.StringFlag{Name: "json", Usage: "Additional fields as JSON"},
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
					body := map[string]any{}
					if cmd.IsSet("customer_id") {
						body["customer_id"] = cmd.String("customer_id")
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
					raw, err := c.Request("PUT", c.BooksBase+"/creditnotes", &zohttp.RequestOpts{
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
				Usage: "List credit notes",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "status", Usage: "Filter"},
					&cli.StringFlag{Name: "customer-id", Usage: "Filter"},
					&cli.StringFlag{Name: "page", Usage: "Filter"},
					&cli.StringFlag{Name: "per-page", Usage: "Filter"},
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
					params := orgParams(orgID)
					if v := cmd.String("status"); v != "" {
						params["status"] = v
					}
					if v := cmd.String("customer-id"); v != "" {
						params["customer_id"] = v
					}
					if v := cmd.String("page"); v != "" {
						params["page"] = v
					}
					if v := cmd.String("per-page"); v != "" {
						params["per_page"] = v
					}
					raw, err := c.Request("GET", c.BooksBase+"/creditnotes", &zohttp.RequestOpts{Params: params})
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
					&cli.StringFlag{Name: "date", Usage: "Credit note date (YYYY-MM-DD)"},
					&cli.StringFlag{Name: "reference_number", Usage: "Reference number"},
					&cli.StringFlag{Name: "json", Usage: "Additional fields as JSON"},
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
					body := map[string]any{}
					if cmd.IsSet("customer_id") {
						body["customer_id"] = cmd.String("customer_id")
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
					raw, err := c.Request("PUT", c.BooksBase+"/creditnotes/"+cmd.Args().First(), &zohttp.RequestOpts{
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
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("GET", c.BooksBase+"/creditnotes/"+cmd.Args().First(), &zohttp.RequestOpts{Params: orgParams(orgID)})
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
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("DELETE", c.BooksBase+"/creditnotes/"+cmd.Args().First(), &zohttp.RequestOpts{Params: orgParams(orgID)})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "email",
				Usage:     "Email a credit note",
				ArgsUsage: "<creditnote-id>",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "to_mail_ids", Required: true, Usage: "Recipient email addresses (comma-separated)"},
					&cli.StringFlag{Name: "subject", Required: true, Usage: "Email subject"},
					&cli.StringFlag{Name: "body", Required: true, Usage: "Email body"},
					&cli.StringFlag{Name: "cc_mail_ids", Usage: "CC email addresses (comma-separated)"},
					&cli.StringFlag{Name: "json", Usage: "Additional fields as JSON"},
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
					raw, err := c.Request("POST", c.BooksBase+"/creditnotes/"+cmd.Args().First()+"/email", &zohttp.RequestOpts{
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
				Usage:     "Get email content of a credit note",
				ArgsUsage: "<creditnote-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("GET", c.BooksBase+"/creditnotes/"+cmd.Args().First()+"/email", &zohttp.RequestOpts{Params: orgParams(orgID)})
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
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("POST", c.BooksBase+"/creditnotes/"+cmd.Args().First()+"/status/void", &zohttp.RequestOpts{Params: orgParams(orgID)})
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
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("POST", c.BooksBase+"/creditnotes/"+cmd.Args().First()+"/status/draft", &zohttp.RequestOpts{Params: orgParams(orgID)})
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
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("POST", c.BooksBase+"/creditnotes/"+cmd.Args().First()+"/status/open", &zohttp.RequestOpts{Params: orgParams(orgID)})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "submit-for-approval",
				Usage:     "Submit a credit note for approval",
				ArgsUsage: "<creditnote-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("POST", c.BooksBase+"/creditnotes/"+cmd.Args().First()+"/submit", &zohttp.RequestOpts{Params: orgParams(orgID)})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "approve",
				Usage:     "Approve a credit note",
				ArgsUsage: "<creditnote-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("POST", c.BooksBase+"/creditnotes/"+cmd.Args().First()+"/approve", &zohttp.RequestOpts{Params: orgParams(orgID)})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "email-history",
				Usage:     "Get email history of a credit note",
				ArgsUsage: "<creditnote-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("GET", c.BooksBase+"/creditnotes/"+cmd.Args().First()+"/emailhistory", &zohttp.RequestOpts{Params: orgParams(orgID)})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "update-billing-address",
				Usage:     "Update billing address of a credit note",
				ArgsUsage: "<creditnote-id>",
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
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
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
					raw, err := c.Request("PUT", c.BooksBase+"/creditnotes/"+cmd.Args().First()+"/address/billing", &zohttp.RequestOpts{
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
				Usage:     "Update shipping address of a credit note",
				ArgsUsage: "<creditnote-id>",
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
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
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
					raw, err := c.Request("PUT", c.BooksBase+"/creditnotes/"+cmd.Args().First()+"/address/shipping", &zohttp.RequestOpts{
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
				Usage: "List credit note templates",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("GET", c.BooksBase+"/creditnotes"+"/templates", &zohttp.RequestOpts{Params: orgParams(orgID)})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "update-template",
				Usage:     "Update template of a credit note",
				ArgsUsage: "<creditnote-id> <template-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("PUT", c.BooksBase+"/creditnotes/"+cmd.Args().First()+"/templates/"+cmd.Args().Get(1), &zohttp.RequestOpts{Params: orgParams(orgID)})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "credit-to-invoice",
				Usage:     "Apply credit note to invoices",
				ArgsUsage: "<creditnote-id>",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "invoice_id", Required: true, Usage: "Invoice ID to apply credit to"},
					&cli.FloatFlag{Name: "amount_applied", Required: true, Usage: "Amount to apply"},
					&cli.StringFlag{Name: "json", Usage: "Additional fields as JSON"},
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
					body := map[string]any{}
					body["invoice_id"] = cmd.String("invoice_id")
					body["amount_applied"] = cmd.Float("amount_applied")
					if err := internal.MergeJSON(cmd, body); err != nil {
						return err
					}
					raw, err := c.Request("POST", c.BooksBase+"/creditnotes/"+cmd.Args().First()+"/invoices", &zohttp.RequestOpts{
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
				Name:      "list-invoices-credited",
				Usage:     "List invoices credited",
				ArgsUsage: "<creditnote-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("GET", c.BooksBase+"/creditnotes/"+cmd.Args().First()+"/invoices", &zohttp.RequestOpts{Params: orgParams(orgID)})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "delete-invoices-credited",
				Usage:     "Delete invoices credited",
				ArgsUsage: "<creditnote-id> <credited-invoice-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("DELETE", c.BooksBase+"/creditnotes/"+cmd.Args().First()+"/invoices/"+cmd.Args().Get(1), &zohttp.RequestOpts{Params: orgParams(orgID)})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "add-comment",
				Usage:     "Add a comment to a credit note",
				ArgsUsage: "<creditnote-id>",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "description", Required: true, Usage: "Comment text"},
					&cli.StringFlag{Name: "json", Usage: "Additional fields as JSON"},
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
					body := map[string]any{}
					body["description"] = cmd.String("description")
					if err := internal.MergeJSON(cmd, body); err != nil {
						return err
					}
					raw, err := c.Request("POST", c.BooksBase+"/creditnotes/"+cmd.Args().First()+"/comments", &zohttp.RequestOpts{
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
				Usage:     "List comments of a credit note",
				ArgsUsage: "<creditnote-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("GET", c.BooksBase+"/creditnotes/"+cmd.Args().First()+"/comments", &zohttp.RequestOpts{Params: orgParams(orgID)})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "delete-comment",
				Usage:     "Delete a comment on a credit note",
				ArgsUsage: "<creditnote-id> <comment-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("DELETE", c.BooksBase+"/creditnotes/"+cmd.Args().First()+"/comments/"+cmd.Args().Get(1), &zohttp.RequestOpts{Params: orgParams(orgID)})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "refund",
				Usage:     "Refund a credit note",
				ArgsUsage: "<creditnote-id>",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "date", Required: true, Usage: "Refund date (YYYY-MM-DD)"},
					&cli.StringFlag{Name: "refund_mode", Usage: "Refund mode"},
					&cli.FloatFlag{Name: "amount", Required: true, Usage: "Refund amount"},
					&cli.StringFlag{Name: "json", Usage: "Additional fields as JSON"},
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
					body := map[string]any{}
					body["date"] = cmd.String("date")
					if cmd.IsSet("refund_mode") {
						body["refund_mode"] = cmd.String("refund_mode")
					}
					body["amount"] = cmd.Float("amount")
					if err := internal.MergeJSON(cmd, body); err != nil {
						return err
					}
					raw, err := c.Request("POST", c.BooksBase+"/creditnotes/"+cmd.Args().First()+"/refunds", &zohttp.RequestOpts{
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
				Usage:     "List refunds of a credit note",
				ArgsUsage: "<creditnote-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("GET", c.BooksBase+"/creditnotes/"+cmd.Args().First()+"/refunds", &zohttp.RequestOpts{Params: orgParams(orgID)})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "update-refund",
				Usage:     "Update a refund of a credit note",
				ArgsUsage: "<creditnote-id> <refund-id>",
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
					c, err := getClient()
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
					raw, err := c.Request("PUT", c.BooksBase+"/creditnotes/"+cmd.Args().First()+"/refunds/"+cmd.Args().Get(1), &zohttp.RequestOpts{
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
				Usage:     "Get a refund of a credit note",
				ArgsUsage: "<creditnote-id> <refund-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("GET", c.BooksBase+"/creditnotes/"+cmd.Args().First()+"/refunds/"+cmd.Args().Get(1), &zohttp.RequestOpts{Params: orgParams(orgID)})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "delete-refund",
				Usage:     "Delete a refund of a credit note",
				ArgsUsage: "<creditnote-id> <refund-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("DELETE", c.BooksBase+"/creditnotes/"+cmd.Args().First()+"/refunds/"+cmd.Args().Get(1), &zohttp.RequestOpts{Params: orgParams(orgID)})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
		},
	}
}

func customerDebitNotesCmd() *cli.Command {
	return &cli.Command{
		Name:  "customer-debit-notes",
		Usage: "Customer debit note operations",
		Commands: []*cli.Command{
			{
				Name:  "create",
				Usage: "Create a customer debit note",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "customer_id", Required: true, Usage: "Customer ID"},
					&cli.StringFlag{Name: "date", Required: true, Usage: "Debit note date (YYYY-MM-DD)"},
					&cli.StringFlag{Name: "invoice_number", Usage: "Debit note number"},
					&cli.StringFlag{Name: "due_date", Usage: "Due date (YYYY-MM-DD)"},
					&cli.StringFlag{Name: "json", Usage: "Additional fields as JSON"},
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
					body := map[string]any{}
					body["customer_id"] = cmd.String("customer_id")
					body["date"] = cmd.String("date")
					if cmd.IsSet("invoice_number") {
						body["invoice_number"] = cmd.String("invoice_number")
					}
					if cmd.IsSet("due_date") {
						body["due_date"] = cmd.String("due_date")
					}
					if err := internal.MergeJSON(cmd, body); err != nil {
						return err
					}
					raw, err := c.Request("POST", c.BooksBase+"/customerdebitnotes", &zohttp.RequestOpts{
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
				Usage: "List customer debit notes",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "page", Usage: "Filter"},
					&cli.StringFlag{Name: "per-page", Usage: "Filter"},
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
					params := orgParams(orgID)
					if v := cmd.String("page"); v != "" {
						params["page"] = v
					}
					if v := cmd.String("per-page"); v != "" {
						params["per_page"] = v
					}
					raw, err := c.Request("GET", c.BooksBase+"/customerdebitnotes", &zohttp.RequestOpts{Params: params})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "update",
				Usage:     "Update a customer debit note",
				ArgsUsage: "<debitnote-id>",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "customer_id", Usage: "Customer ID"},
					&cli.StringFlag{Name: "date", Usage: "Debit note date (YYYY-MM-DD)"},
					&cli.StringFlag{Name: "invoice_number", Usage: "Debit note number"},
					&cli.StringFlag{Name: "due_date", Usage: "Due date (YYYY-MM-DD)"},
					&cli.StringFlag{Name: "json", Usage: "Additional fields as JSON"},
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
					body := map[string]any{}
					if cmd.IsSet("customer_id") {
						body["customer_id"] = cmd.String("customer_id")
					}
					if cmd.IsSet("date") {
						body["date"] = cmd.String("date")
					}
					if cmd.IsSet("invoice_number") {
						body["invoice_number"] = cmd.String("invoice_number")
					}
					if cmd.IsSet("due_date") {
						body["due_date"] = cmd.String("due_date")
					}
					if err := internal.MergeJSON(cmd, body); err != nil {
						return err
					}
					raw, err := c.Request("PUT", c.BooksBase+"/customerdebitnotes/"+cmd.Args().First(), &zohttp.RequestOpts{
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
				Usage:     "Get a customer debit note",
				ArgsUsage: "<debitnote-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("GET", c.BooksBase+"/customerdebitnotes/"+cmd.Args().First(), &zohttp.RequestOpts{Params: orgParams(orgID)})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "delete",
				Usage:     "Delete a customer debit note",
				ArgsUsage: "<debitnote-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("DELETE", c.BooksBase+"/customerdebitnotes/"+cmd.Args().First(), &zohttp.RequestOpts{Params: orgParams(orgID)})
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
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
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
					raw, err := c.Request("POST", c.BooksBase+"/customerpayments", &zohttp.RequestOpts{
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
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
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
					raw, err := c.Request("PUT", c.BooksBase+"/customerpayments", &zohttp.RequestOpts{
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
					&cli.StringFlag{Name: "page", Usage: "Filter"},
					&cli.StringFlag{Name: "per-page", Usage: "Filter"},
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
					params := orgParams(orgID)
					if v := cmd.String("customer-id"); v != "" {
						params["customer_id"] = v
					}
					if v := cmd.String("page"); v != "" {
						params["page"] = v
					}
					if v := cmd.String("per-page"); v != "" {
						params["per_page"] = v
					}
					raw, err := c.Request("GET", c.BooksBase+"/customerpayments", &zohttp.RequestOpts{Params: params})
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
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
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
					raw, err := c.Request("PUT", c.BooksBase+"/customerpayments/"+cmd.Args().First(), &zohttp.RequestOpts{
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
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("GET", c.BooksBase+"/customerpayments/"+cmd.Args().First(), &zohttp.RequestOpts{Params: orgParams(orgID)})
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
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("DELETE", c.BooksBase+"/customerpayments/"+cmd.Args().First(), &zohttp.RequestOpts{Params: orgParams(orgID)})
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
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
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
					raw, err := c.Request("POST", c.BooksBase+"/customerpayments/"+cmd.Args().First()+"/refunds", &zohttp.RequestOpts{
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
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("GET", c.BooksBase+"/customerpayments/"+cmd.Args().First()+"/refunds", &zohttp.RequestOpts{Params: orgParams(orgID)})
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
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
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
					raw, err := c.Request("PUT", c.BooksBase+"/customerpayments/"+cmd.Args().First()+"/customfields", &zohttp.RequestOpts{
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
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
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
					raw, err := c.Request("PUT", c.BooksBase+"/customerpayments/"+cmd.Args().First()+"/refunds/"+cmd.Args().Get(1), &zohttp.RequestOpts{
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
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("GET", c.BooksBase+"/customerpayments/"+cmd.Args().First()+"/refunds/"+cmd.Args().Get(1), &zohttp.RequestOpts{Params: orgParams(orgID)})
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
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("DELETE", c.BooksBase+"/customerpayments/"+cmd.Args().First()+"/refunds/"+cmd.Args().Get(1), &zohttp.RequestOpts{Params: orgParams(orgID)})
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
					c, err := getClient()
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
					c, err := getClient()
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
					&cli.StringFlag{Name: "page", Usage: "Filter"},
					&cli.StringFlag{Name: "per-page", Usage: "Filter"},
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
					if v := cmd.String("page"); v != "" {
						params["page"] = v
					}
					if v := cmd.String("per-page"); v != "" {
						params["per_page"] = v
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
					c, err := getClient()
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
					c, err := getClient()
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
					c, err := getClient()
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
					c, err := getClient()
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
					c, err := getClient()
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
					c, err := getClient()
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
					c, err := getClient()
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
					c, err := getClient()
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
					c, err := getClient()
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
					c, err := getClient()
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
					c, err := getClient()
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
					c, err := getClient()
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
					c, err := getClient()
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
					c, err := getClient()
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
					&cli.StringFlag{Name: "page", Usage: "Filter"},
					&cli.StringFlag{Name: "per-page", Usage: "Filter"},
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
					params := orgParams(orgID)
					if v := cmd.String("status"); v != "" {
						params["status"] = v
					}
					if v := cmd.String("page"); v != "" {
						params["page"] = v
					}
					if v := cmd.String("per-page"); v != "" {
						params["per_page"] = v
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
					c, err := getClient()
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
					c, err := getClient()
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
					c, err := getClient()
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
					c, err := getClient()
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
					c, err := getClient()
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
					c, err := getClient()
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
					c, err := getClient()
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

func retainerInvoicesCmd() *cli.Command {
	return &cli.Command{
		Name:  "retainer-invoices",
		Usage: "Retainer invoice operations",
		Commands: []*cli.Command{
			{
				Name:  "create",
				Usage: "Create a retainer invoice",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "customer_id", Required: true, Usage: "Customer ID"},
					&cli.StringFlag{Name: "date", Usage: "Retainer invoice date (YYYY-MM-DD)"},
					&cli.StringFlag{Name: "reference_number", Usage: "Reference number"},
					&cli.StringFlag{Name: "json", Usage: "Additional fields as JSON"},
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
					body := map[string]any{}
					body["customer_id"] = cmd.String("customer_id")
					if cmd.IsSet("date") {
						body["date"] = cmd.String("date")
					}
					if cmd.IsSet("reference_number") {
						body["reference_number"] = cmd.String("reference_number")
					}
					if err := internal.MergeJSON(cmd, body); err != nil {
						return err
					}
					raw, err := c.Request("POST", c.BooksBase+"/retainerinvoices", &zohttp.RequestOpts{
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
				Usage: "List retainer invoices",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "status", Usage: "Filter"},
					&cli.StringFlag{Name: "customer-id", Usage: "Filter"},
					&cli.StringFlag{Name: "page", Usage: "Filter"},
					&cli.StringFlag{Name: "per-page", Usage: "Filter"},
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
					params := orgParams(orgID)
					if v := cmd.String("status"); v != "" {
						params["status"] = v
					}
					if v := cmd.String("customer-id"); v != "" {
						params["customer_id"] = v
					}
					if v := cmd.String("page"); v != "" {
						params["page"] = v
					}
					if v := cmd.String("per-page"); v != "" {
						params["per_page"] = v
					}
					raw, err := c.Request("GET", c.BooksBase+"/retainerinvoices", &zohttp.RequestOpts{Params: params})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "update",
				Usage:     "Update a retainer invoice",
				ArgsUsage: "<retainerinvoice-id>",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "customer_id", Usage: "Customer ID"},
					&cli.StringFlag{Name: "date", Usage: "Retainer invoice date (YYYY-MM-DD)"},
					&cli.StringFlag{Name: "reference_number", Usage: "Reference number"},
					&cli.StringFlag{Name: "json", Usage: "Additional fields as JSON"},
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
					body := map[string]any{}
					if cmd.IsSet("customer_id") {
						body["customer_id"] = cmd.String("customer_id")
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
					raw, err := c.Request("PUT", c.BooksBase+"/retainerinvoices/"+cmd.Args().First(), &zohttp.RequestOpts{
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
				ArgsUsage: "<retainerinvoice-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("GET", c.BooksBase+"/retainerinvoices/"+cmd.Args().First(), &zohttp.RequestOpts{Params: orgParams(orgID)})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "delete",
				Usage:     "Delete a retainer invoice",
				ArgsUsage: "<retainerinvoice-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("DELETE", c.BooksBase+"/retainerinvoices/"+cmd.Args().First(), &zohttp.RequestOpts{Params: orgParams(orgID)})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "mark-sent",
				Usage:     "Mark a retainer invoice as sent",
				ArgsUsage: "<retainerinvoice-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("POST", c.BooksBase+"/retainerinvoices/"+cmd.Args().First()+"/status/sent", &zohttp.RequestOpts{Params: orgParams(orgID)})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "void",
				Usage:     "Void a retainer invoice",
				ArgsUsage: "<retainerinvoice-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("POST", c.BooksBase+"/retainerinvoices/"+cmd.Args().First()+"/status/void", &zohttp.RequestOpts{Params: orgParams(orgID)})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "mark-draft",
				Usage:     "Mark a retainer invoice as draft",
				ArgsUsage: "<retainerinvoice-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("POST", c.BooksBase+"/retainerinvoices/"+cmd.Args().First()+"/status/draft", &zohttp.RequestOpts{Params: orgParams(orgID)})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "submit-for-approval",
				Usage:     "Submit a retainer invoice for approval",
				ArgsUsage: "<retainerinvoice-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("POST", c.BooksBase+"/retainerinvoices/"+cmd.Args().First()+"/submit", &zohttp.RequestOpts{Params: orgParams(orgID)})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "approve",
				Usage:     "Approve a retainer invoice",
				ArgsUsage: "<retainerinvoice-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("POST", c.BooksBase+"/retainerinvoices/"+cmd.Args().First()+"/approve", &zohttp.RequestOpts{Params: orgParams(orgID)})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "email",
				Usage:     "Email a retainer invoice",
				ArgsUsage: "<retainerinvoice-id>",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "to_mail_ids", Required: true, Usage: "Recipient email addresses (comma-separated)"},
					&cli.StringFlag{Name: "cc_mail_ids", Usage: "CC email addresses (comma-separated)"},
					&cli.StringFlag{Name: "subject", Required: true, Usage: "Email subject"},
					&cli.StringFlag{Name: "body", Required: true, Usage: "Email body"},
					&cli.StringFlag{Name: "json", Usage: "Additional fields as JSON"},
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
					raw, err := c.Request("POST", c.BooksBase+"/retainerinvoices/"+cmd.Args().First()+"/email", &zohttp.RequestOpts{
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
				Usage:     "Get email content of a retainer invoice",
				ArgsUsage: "<retainerinvoice-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("GET", c.BooksBase+"/retainerinvoices/"+cmd.Args().First()+"/email", &zohttp.RequestOpts{Params: orgParams(orgID)})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "update-billing-address",
				Usage:     "Update billing address of a retainer invoice",
				ArgsUsage: "<retainerinvoice-id>",
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
					c, err := getClient()
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
					raw, err := c.Request("PUT", c.BooksBase+"/retainerinvoices/"+cmd.Args().First()+"/address/billing", &zohttp.RequestOpts{
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
				Usage: "List retainer invoice templates",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("GET", c.BooksBase+"/retainerinvoices"+"/templates", &zohttp.RequestOpts{Params: orgParams(orgID)})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "update-template",
				Usage:     "Update template of a retainer invoice",
				ArgsUsage: "<retainerinvoice-id> <template-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("PUT", c.BooksBase+"/retainerinvoices/"+cmd.Args().First()+"/templates/"+cmd.Args().Get(1), &zohttp.RequestOpts{Params: orgParams(orgID)})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "add-attachment",
				Usage:     "Add attachment to a retainer invoice",
				ArgsUsage: "<retainerinvoice-id>",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "file", Required: true, Usage: "File path to upload"},
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
					data, err := os.ReadFile(cmd.String("file"))
					if err != nil {
						return fmt.Errorf("reading file: %w", err)
					}
					raw, err := c.Request("POST", c.BooksBase+"/retainerinvoices/"+cmd.Args().First()+"/attachment", &zohttp.RequestOpts{
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
				Usage:     "Get attachment of a retainer invoice",
				ArgsUsage: "<retainerinvoice-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("GET", c.BooksBase+"/retainerinvoices/"+cmd.Args().First()+"/attachment", &zohttp.RequestOpts{Params: orgParams(orgID)})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "delete-attachment",
				Usage:     "Delete attachment of a retainer invoice",
				ArgsUsage: "<retainerinvoice-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("DELETE", c.BooksBase+"/retainerinvoices/"+cmd.Args().First()+"/attachment", &zohttp.RequestOpts{Params: orgParams(orgID)})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "add-comment",
				Usage:     "Add a comment to a retainer invoice",
				ArgsUsage: "<retainerinvoice-id>",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "description", Required: true, Usage: "Description"},
					&cli.StringFlag{Name: "json", Usage: "Additional fields as JSON"},
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
					body := map[string]any{}
					body["description"] = cmd.String("description")
					if err := internal.MergeJSON(cmd, body); err != nil {
						return err
					}
					raw, err := c.Request("POST", c.BooksBase+"/retainerinvoices/"+cmd.Args().First()+"/comments", &zohttp.RequestOpts{
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
				Usage:     "List comments of a retainer invoice",
				ArgsUsage: "<retainerinvoice-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("GET", c.BooksBase+"/retainerinvoices/"+cmd.Args().First()+"/comments", &zohttp.RequestOpts{Params: orgParams(orgID)})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "update-comment",
				Usage:     "Update a comment on a retainer invoice",
				ArgsUsage: "<retainerinvoice-id> <comment-id>",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "description", Required: true, Usage: "Comment text"},
					&cli.StringFlag{Name: "json", Usage: "Additional fields as JSON"},
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
					body := map[string]any{}
					body["description"] = cmd.String("description")
					if err := internal.MergeJSON(cmd, body); err != nil {
						return err
					}
					raw, err := c.Request("PUT", c.BooksBase+"/retainerinvoices/"+cmd.Args().First()+"/comments/"+cmd.Args().Get(1), &zohttp.RequestOpts{
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
				Usage:     "Delete a comment on a retainer invoice",
				ArgsUsage: "<retainerinvoice-id> <comment-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("DELETE", c.BooksBase+"/retainerinvoices/"+cmd.Args().First()+"/comments/"+cmd.Args().Get(1), &zohttp.RequestOpts{Params: orgParams(orgID)})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
		},
	}
}

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
					c, err := getClient()
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
					c, err := getClient()
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
					&cli.StringFlag{Name: "page", Usage: "Filter"},
					&cli.StringFlag{Name: "per-page", Usage: "Filter"},
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
					params := orgParams(orgID)
					if v := cmd.String("status"); v != "" {
						params["status"] = v
					}
					if v := cmd.String("vendor-id"); v != "" {
						params["vendor_id"] = v
					}
					if v := cmd.String("page"); v != "" {
						params["page"] = v
					}
					if v := cmd.String("per-page"); v != "" {
						params["per_page"] = v
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
					c, err := getClient()
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
					c, err := getClient()
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
					c, err := getClient()
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
					c, err := getClient()
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
					c, err := getClient()
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
					c, err := getClient()
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
					c, err := getClient()
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
					c, err := getClient()
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
					c, err := getClient()
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
					c, err := getClient()
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
					c, err := getClient()
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
					c, err := getClient()
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
					c, err := getClient()
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
					c, err := getClient()
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
					c, err := getClient()
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
					c, err := getClient()
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
					c, err := getClient()
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
					c, err := getClient()
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
					c, err := getClient()
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
					c, err := getClient()
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
					c, err := getClient()
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
					c, err := getClient()
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
					c, err := getClient()
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
					c, err := getClient()
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
					c, err := getClient()
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
					&cli.StringFlag{Name: "page", Usage: "Filter"},
					&cli.StringFlag{Name: "per-page", Usage: "Filter"},
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
					params := orgParams(orgID)
					if v := cmd.String("status"); v != "" {
						params["status"] = v
					}
					if v := cmd.String("vendor-id"); v != "" {
						params["vendor_id"] = v
					}
					if v := cmd.String("page"); v != "" {
						params["page"] = v
					}
					if v := cmd.String("per-page"); v != "" {
						params["per_page"] = v
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
					c, err := getClient()
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
					c, err := getClient()
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
					c, err := getClient()
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
					c, err := getClient()
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
					c, err := getClient()
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
					c, err := getClient()
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
					c, err := getClient()
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
					c, err := getClient()
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
					c, err := getClient()
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
					c, err := getClient()
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
					c, err := getClient()
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
					c, err := getClient()
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
					c, err := getClient()
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
					c, err := getClient()
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
					c, err := getClient()
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
					c, err := getClient()
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
					c, err := getClient()
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
					c, err := getClient()
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
					c, err := getClient()
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
					c, err := getClient()
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
					&cli.StringFlag{Name: "page", Usage: "Filter"},
					&cli.StringFlag{Name: "per-page", Usage: "Filter"},
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
					params := orgParams(orgID)
					if v := cmd.String("status"); v != "" {
						params["status"] = v
					}
					if v := cmd.String("page"); v != "" {
						params["page"] = v
					}
					if v := cmd.String("per-page"); v != "" {
						params["per_page"] = v
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
					c, err := getClient()
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
					c, err := getClient()
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
					c, err := getClient()
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
					c, err := getClient()
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
					c, err := getClient()
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
					c, err := getClient()
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
					c, err := getClient()
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
					&cli.StringFlag{Name: "page", Usage: "Filter"},
					&cli.StringFlag{Name: "per-page", Usage: "Filter"},
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
					params := orgParams(orgID)
					if v := cmd.String("status"); v != "" {
						params["status"] = v
					}
					if v := cmd.String("vendor-id"); v != "" {
						params["vendor_id"] = v
					}
					if v := cmd.String("page"); v != "" {
						params["page"] = v
					}
					if v := cmd.String("per-page"); v != "" {
						params["per_page"] = v
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
					c, err := getClient()
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
					c, err := getClient()
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
					c, err := getClient()
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
					c, err := getClient()
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
					c, err := getClient()
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
					c, err := getClient()
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
					c, err := getClient()
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
					c, err := getClient()
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
					c, err := getClient()
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
					c, err := getClient()
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
					c, err := getClient()
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
					c, err := getClient()
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
					c, err := getClient()
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
					c, err := getClient()
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
					c, err := getClient()
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
					c, err := getClient()
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
					c, err := getClient()
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
					c, err := getClient()
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
					c, err := getClient()
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
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
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
					raw, err := c.Request("POST", c.BooksBase+"/vendorpayments", &zohttp.RequestOpts{
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
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
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
					raw, err := c.Request("PUT", c.BooksBase+"/vendorpayments", &zohttp.RequestOpts{
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
					&cli.StringFlag{Name: "page", Usage: "Filter"},
					&cli.StringFlag{Name: "per-page", Usage: "Filter"},
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
					params := orgParams(orgID)
					if v := cmd.String("vendor-id"); v != "" {
						params["vendor_id"] = v
					}
					if v := cmd.String("page"); v != "" {
						params["page"] = v
					}
					if v := cmd.String("per-page"); v != "" {
						params["per_page"] = v
					}
					raw, err := c.Request("GET", c.BooksBase+"/vendorpayments", &zohttp.RequestOpts{Params: params})
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
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
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
					raw, err := c.Request("PUT", c.BooksBase+"/vendorpayments/"+cmd.Args().First(), &zohttp.RequestOpts{
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
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("GET", c.BooksBase+"/vendorpayments/"+cmd.Args().First(), &zohttp.RequestOpts{Params: orgParams(orgID)})
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
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("DELETE", c.BooksBase+"/vendorpayments/"+cmd.Args().First(), &zohttp.RequestOpts{Params: orgParams(orgID)})
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
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
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
					raw, err := c.Request("POST", c.BooksBase+"/vendorpayments/"+cmd.Args().First()+"/refunds", &zohttp.RequestOpts{
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
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("GET", c.BooksBase+"/vendorpayments/"+cmd.Args().First()+"/refunds", &zohttp.RequestOpts{Params: orgParams(orgID)})
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
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
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
					raw, err := c.Request("PUT", c.BooksBase+"/vendorpayments/"+cmd.Args().First()+"/refunds/"+cmd.Args().Get(1), &zohttp.RequestOpts{
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
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("GET", c.BooksBase+"/vendorpayments/"+cmd.Args().First()+"/refunds/"+cmd.Args().Get(1), &zohttp.RequestOpts{Params: orgParams(orgID)})
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
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("DELETE", c.BooksBase+"/vendorpayments/"+cmd.Args().First()+"/refunds/"+cmd.Args().Get(1), &zohttp.RequestOpts{Params: orgParams(orgID)})
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
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
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
					raw, err := c.Request("POST", c.BooksBase+"/vendorpayments/"+cmd.Args().First()+"/email", &zohttp.RequestOpts{
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
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("GET", c.BooksBase+"/vendorpayments/"+cmd.Args().First()+"/email", &zohttp.RequestOpts{Params: orgParams(orgID)})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
		},
	}
}

func customModulesCmd() *cli.Command {
	return &cli.Command{
		Name:  "custom-modules",
		Usage: "Custom module operations",
		Commands: []*cli.Command{
			{
				Name:  "create",
				Usage: "Create a custom module record",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "module", Required: true, Usage: "Module name"},
					&cli.StringFlag{Name: "record_name", Required: true, Usage: "Record name"},
					&cli.StringFlag{Name: "json", Usage: "Additional fields as JSON"},
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
					body := map[string]any{}
					body["record_name"] = cmd.String("record_name")
					if err := internal.MergeJSON(cmd, body); err != nil {
						return err
					}
					raw, err := c.Request("POST", c.BooksBase+"/custommodules/"+cmd.String("module"), &zohttp.RequestOpts{
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
				Name:  "bulk-update",
				Usage: "Bulk update custom module records",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "module", Required: true, Usage: "Module name"},
					&cli.StringFlag{Name: "record_name", Usage: "Record name"},
					&cli.StringFlag{Name: "json", Usage: "Additional fields as JSON"},
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
					body := map[string]any{}
					if cmd.IsSet("record_name") {
						body["record_name"] = cmd.String("record_name")
					}
					if err := internal.MergeJSON(cmd, body); err != nil {
						return err
					}
					raw, err := c.Request("PUT", c.BooksBase+"/custommodules/"+cmd.String("module"), &zohttp.RequestOpts{
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
				Usage: "List custom module records",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "module", Required: true, Usage: "Module name"},
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
					params := orgParams(orgID)
					if v := cmd.String("page"); v != "" {
						params["page"] = v
					}
					if v := cmd.String("per-page"); v != "" {
						params["per_page"] = v
					}
					raw, err := c.Request("GET", c.BooksBase+"/custommodules/"+cmd.String("module"), &zohttp.RequestOpts{Params: params})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "update-record",
				Usage:     "Update a custom module record",
				ArgsUsage: "<record-id>",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "module", Required: true, Usage: "Module name"},
					&cli.StringFlag{Name: "record_name", Usage: "Record name"},
					&cli.StringFlag{Name: "json", Usage: "Additional fields as JSON"},
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
					body := map[string]any{}
					if cmd.IsSet("record_name") {
						body["record_name"] = cmd.String("record_name")
					}
					if err := internal.MergeJSON(cmd, body); err != nil {
						return err
					}
					raw, err := c.Request("PUT", c.BooksBase+"/custommodules/"+cmd.String("module")+"/"+cmd.Args().First(), &zohttp.RequestOpts{
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
				Name:      "get-record",
				Usage:     "Get a custom module record",
				ArgsUsage: "<record-id>",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "module", Required: true, Usage: "Module name"},
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
					raw, err := c.Request("GET", c.BooksBase+"/custommodules/"+cmd.String("module")+"/"+cmd.Args().First(), &zohttp.RequestOpts{Params: orgParams(orgID)})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "delete-record",
				Usage:     "Delete a custom module record",
				ArgsUsage: "<record-id>",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "module", Required: true, Usage: "Module name"},
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
					raw, err := c.Request("DELETE", c.BooksBase+"/custommodules/"+cmd.String("module")+"/"+cmd.Args().First(), &zohttp.RequestOpts{Params: orgParams(orgID)})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
		},
	}
}

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
					c, err := getClient()
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
					&cli.StringFlag{Name: "page", Usage: "Filter"},
					&cli.StringFlag{Name: "per-page", Usage: "Filter"},
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
					params := orgParams(orgID)
					if v := cmd.String("page"); v != "" {
						params["page"] = v
					}
					if v := cmd.String("per-page"); v != "" {
						params["per_page"] = v
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
					c, err := getClient()
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
					c, err := getClient()
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
					c, err := getClient()
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
					c, err := getClient()
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
					c, err := getClient()
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
					c, err := getClient()
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
					c, err := getClient()
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
					c, err := getClient()
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
					c, err := getClient()
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
					&cli.StringFlag{Name: "page", Usage: "Filter"},
					&cli.StringFlag{Name: "per-page", Usage: "Filter"},
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
					params := orgParams(orgID)
					if v := cmd.String("account-id"); v != "" {
						params["account_id"] = v
					}
					if v := cmd.String("status"); v != "" {
						params["status"] = v
					}
					if v := cmd.String("page"); v != "" {
						params["page"] = v
					}
					if v := cmd.String("per-page"); v != "" {
						params["per_page"] = v
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
					c, err := getClient()
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
					c, err := getClient()
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
					c, err := getClient()
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
					c, err := getClient()
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
					c, err := getClient()
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
					c, err := getClient()
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
					c, err := getClient()
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
					c, err := getClient()
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
					c, err := getClient()
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
					c, err := getClient()
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
					c, err := getClient()
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
					c, err := getClient()
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
					c, err := getClient()
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
					c, err := getClient()
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
					c, err := getClient()
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
					c, err := getClient()
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
					c, err := getClient()
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
					c, err := getClient()
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
					c, err := getClient()
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
					c, err := getClient()
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
					c, err := getClient()
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
					c, err := getClient()
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

func chartOfAccountsCmd() *cli.Command {
	return &cli.Command{
		Name:  "chart-of-accounts",
		Usage: "Chart of accounts operations",
		Commands: []*cli.Command{
			{
				Name:  "create",
				Usage: "Create a chart of account",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "account_name", Required: true, Usage: "Account name"},
					&cli.StringFlag{Name: "account_type", Required: true, Usage: "Account type"},
					&cli.StringFlag{Name: "description", Usage: "Description"},
					&cli.StringFlag{Name: "parent_account_id", Usage: "Parent account ID"},
					&cli.BoolFlag{Name: "is_sub_account", Usage: "Is sub account"},
					&cli.StringFlag{Name: "json", Usage: "Additional fields as JSON"},
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
					body := map[string]any{}
					body["account_name"] = cmd.String("account_name")
					body["account_type"] = cmd.String("account_type")
					if cmd.IsSet("description") {
						body["description"] = cmd.String("description")
					}
					if cmd.IsSet("parent_account_id") {
						body["parent_account_id"] = cmd.String("parent_account_id")
					}
					if cmd.IsSet("is_sub_account") {
						body["is_sub_account"] = cmd.Bool("is_sub_account")
					}
					if err := internal.MergeJSON(cmd, body); err != nil {
						return err
					}
					raw, err := c.Request("POST", c.BooksBase+"/chartofaccounts", &zohttp.RequestOpts{
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
				Usage: "List chart of accounts",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "sort-column", Usage: "Filter"},
					&cli.StringFlag{Name: "page", Usage: "Filter"},
					&cli.StringFlag{Name: "per-page", Usage: "Filter"},
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
					params := orgParams(orgID)
					if v := cmd.String("sort-column"); v != "" {
						params["sort_column"] = v
					}
					if v := cmd.String("page"); v != "" {
						params["page"] = v
					}
					if v := cmd.String("per-page"); v != "" {
						params["per_page"] = v
					}
					raw, err := c.Request("GET", c.BooksBase+"/chartofaccounts", &zohttp.RequestOpts{Params: params})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "update",
				Usage:     "Update a chart of account",
				ArgsUsage: "<account-id>",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "account_name", Usage: "Account name"},
					&cli.StringFlag{Name: "account_type", Usage: "Account type"},
					&cli.StringFlag{Name: "description", Usage: "Description"},
					&cli.StringFlag{Name: "parent_account_id", Usage: "Parent account ID"},
					&cli.BoolFlag{Name: "is_sub_account", Usage: "Is sub account"},
					&cli.StringFlag{Name: "json", Usage: "Additional fields as JSON"},
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
					body := map[string]any{}
					if cmd.IsSet("account_name") {
						body["account_name"] = cmd.String("account_name")
					}
					if cmd.IsSet("account_type") {
						body["account_type"] = cmd.String("account_type")
					}
					if cmd.IsSet("description") {
						body["description"] = cmd.String("description")
					}
					if cmd.IsSet("parent_account_id") {
						body["parent_account_id"] = cmd.String("parent_account_id")
					}
					if cmd.IsSet("is_sub_account") {
						body["is_sub_account"] = cmd.Bool("is_sub_account")
					}
					if err := internal.MergeJSON(cmd, body); err != nil {
						return err
					}
					raw, err := c.Request("PUT", c.BooksBase+"/chartofaccounts/"+cmd.Args().First(), &zohttp.RequestOpts{
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
				Usage:     "Get a chart of account",
				ArgsUsage: "<account-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("GET", c.BooksBase+"/chartofaccounts/"+cmd.Args().First(), &zohttp.RequestOpts{Params: orgParams(orgID)})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "delete",
				Usage:     "Delete a chart of account",
				ArgsUsage: "<account-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("DELETE", c.BooksBase+"/chartofaccounts/"+cmd.Args().First(), &zohttp.RequestOpts{Params: orgParams(orgID)})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "mark-active",
				Usage:     "Mark a chart of account as active",
				ArgsUsage: "<account-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("POST", c.BooksBase+"/chartofaccounts/"+cmd.Args().First()+"/active", &zohttp.RequestOpts{Params: orgParams(orgID)})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "mark-inactive",
				Usage:     "Mark a chart of account as inactive",
				ArgsUsage: "<account-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("POST", c.BooksBase+"/chartofaccounts/"+cmd.Args().First()+"/inactive", &zohttp.RequestOpts{Params: orgParams(orgID)})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "list-transactions",
				Usage:     "List transactions of a chart of account",
				ArgsUsage: "<account-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("GET", c.BooksBase+"/chartofaccounts/"+cmd.Args().First()+"/transactions", &zohttp.RequestOpts{Params: orgParams(orgID)})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "delete-transaction",
				Usage:     "Delete a transaction",
				ArgsUsage: "<transaction-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("DELETE", c.BooksBase+"/chartofaccounts/transactions/"+cmd.Args().First(), &zohttp.RequestOpts{Params: orgParams(orgID)})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
		},
	}
}

func journalsCmd() *cli.Command {
	return &cli.Command{
		Name:  "journals",
		Usage: "Journal operations",
		Commands: []*cli.Command{
			{
				Name:  "create",
				Usage: "Create a journal",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "reference_number", Required: true, Usage: "Reference number"},
					&cli.StringFlag{Name: "journal_date", Required: true, Usage: "Journal date (YYYY-MM-DD)"},
					&cli.StringFlag{Name: "currency_id", Usage: "Currency ID"},
					&cli.StringFlag{Name: "json", Usage: "Additional fields as JSON"},
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
					body := map[string]any{}
					body["reference_number"] = cmd.String("reference_number")
					body["journal_date"] = cmd.String("journal_date")
					if cmd.IsSet("currency_id") {
						body["currency_id"] = cmd.String("currency_id")
					}
					if err := internal.MergeJSON(cmd, body); err != nil {
						return err
					}
					raw, err := c.Request("POST", c.BooksBase+"/journals", &zohttp.RequestOpts{
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
				Usage: "List journals",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "date-start", Usage: "Filter"},
					&cli.StringFlag{Name: "date-end", Usage: "Filter"},
					&cli.StringFlag{Name: "page", Usage: "Filter"},
					&cli.StringFlag{Name: "per-page", Usage: "Filter"},
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
					params := orgParams(orgID)
					if v := cmd.String("date-start"); v != "" {
						params["date_start"] = v
					}
					if v := cmd.String("date-end"); v != "" {
						params["date_end"] = v
					}
					if v := cmd.String("page"); v != "" {
						params["page"] = v
					}
					if v := cmd.String("per-page"); v != "" {
						params["per_page"] = v
					}
					raw, err := c.Request("GET", c.BooksBase+"/journals", &zohttp.RequestOpts{Params: params})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "update",
				Usage:     "Update a journal",
				ArgsUsage: "<journal-id>",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "reference_number", Usage: "Reference number"},
					&cli.StringFlag{Name: "journal_date", Usage: "Journal date (YYYY-MM-DD)"},
					&cli.StringFlag{Name: "currency_id", Usage: "Currency ID"},
					&cli.StringFlag{Name: "json", Usage: "Additional fields as JSON"},
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
					body := map[string]any{}
					if cmd.IsSet("reference_number") {
						body["reference_number"] = cmd.String("reference_number")
					}
					if cmd.IsSet("journal_date") {
						body["journal_date"] = cmd.String("journal_date")
					}
					if cmd.IsSet("currency_id") {
						body["currency_id"] = cmd.String("currency_id")
					}
					if err := internal.MergeJSON(cmd, body); err != nil {
						return err
					}
					raw, err := c.Request("PUT", c.BooksBase+"/journals/"+cmd.Args().First(), &zohttp.RequestOpts{
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
				Usage:     "Get a journal",
				ArgsUsage: "<journal-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("GET", c.BooksBase+"/journals/"+cmd.Args().First(), &zohttp.RequestOpts{Params: orgParams(orgID)})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "delete",
				Usage:     "Delete a journal",
				ArgsUsage: "<journal-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("DELETE", c.BooksBase+"/journals/"+cmd.Args().First(), &zohttp.RequestOpts{Params: orgParams(orgID)})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "mark-published",
				Usage:     "Mark a journal as published",
				ArgsUsage: "<journal-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("POST", c.BooksBase+"/journals/"+cmd.Args().First()+"/publish", &zohttp.RequestOpts{Params: orgParams(orgID)})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "add-attachment",
				Usage:     "Add attachment to a journal",
				ArgsUsage: "<journal-id>",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "file", Required: true, Usage: "File path to upload"},
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
					data, err := os.ReadFile(cmd.String("file"))
					if err != nil {
						return fmt.Errorf("reading file: %w", err)
					}
					raw, err := c.Request("POST", c.BooksBase+"/journals/"+cmd.Args().First()+"/attachment", &zohttp.RequestOpts{
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
				Name:      "add-comment",
				Usage:     "Add a comment to a journal",
				ArgsUsage: "<journal-id>",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "description", Required: true, Usage: "Description"},
					&cli.StringFlag{Name: "json", Usage: "Additional fields as JSON"},
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
					body := map[string]any{}
					body["description"] = cmd.String("description")
					if err := internal.MergeJSON(cmd, body); err != nil {
						return err
					}
					raw, err := c.Request("POST", c.BooksBase+"/journals/"+cmd.Args().First()+"/comments", &zohttp.RequestOpts{
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
				Usage:     "Delete a comment on a journal",
				ArgsUsage: "<journal-id> <comment-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("DELETE", c.BooksBase+"/journals/"+cmd.Args().First()+"/comments/"+cmd.Args().Get(1), &zohttp.RequestOpts{Params: orgParams(orgID)})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
		},
	}
}

func fixedAssetsCmd() *cli.Command {
	return &cli.Command{
		Name:  "fixed-assets",
		Usage: "Fixed asset operations",
		Commands: []*cli.Command{
			{
				Name:  "create",
				Usage: "Create a fixed asset",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "asset_name", Required: true, Usage: "Asset name"},
					&cli.StringFlag{Name: "asset_type_id", Required: true, Usage: "Asset type ID"},
					&cli.StringFlag{Name: "purchase_date", Usage: "Purchase date (YYYY-MM-DD)"},
					&cli.FloatFlag{Name: "purchase_price", Usage: "Purchase price"},
					&cli.StringFlag{Name: "account_id", Usage: "Account ID"},
					&cli.StringFlag{Name: "depreciation_method", Usage: "Depreciation method"},
					&cli.FloatFlag{Name: "depreciation_rate", Usage: "Depreciation rate"},
					&cli.StringFlag{Name: "depreciation_account_id", Usage: "Depreciation account ID"},
					&cli.StringFlag{Name: "description", Usage: "Description"},
					&cli.StringFlag{Name: "json", Usage: "Additional fields as JSON"},
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
					body := map[string]any{}
					body["asset_name"] = cmd.String("asset_name")
					body["asset_type_id"] = cmd.String("asset_type_id")
					if cmd.IsSet("purchase_date") {
						body["purchase_date"] = cmd.String("purchase_date")
					}
					if cmd.IsSet("purchase_price") {
						body["purchase_price"] = cmd.Float("purchase_price")
					}
					if cmd.IsSet("account_id") {
						body["account_id"] = cmd.String("account_id")
					}
					if cmd.IsSet("depreciation_method") {
						body["depreciation_method"] = cmd.String("depreciation_method")
					}
					if cmd.IsSet("depreciation_rate") {
						body["depreciation_rate"] = cmd.Float("depreciation_rate")
					}
					if cmd.IsSet("depreciation_account_id") {
						body["depreciation_account_id"] = cmd.String("depreciation_account_id")
					}
					if cmd.IsSet("description") {
						body["description"] = cmd.String("description")
					}
					if err := internal.MergeJSON(cmd, body); err != nil {
						return err
					}
					raw, err := c.Request("POST", c.BooksBase+"/fixedassets", &zohttp.RequestOpts{
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
				Usage: "List fixed assets",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "page", Usage: "Filter"},
					&cli.StringFlag{Name: "per-page", Usage: "Filter"},
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
					params := orgParams(orgID)
					if v := cmd.String("page"); v != "" {
						params["page"] = v
					}
					if v := cmd.String("per-page"); v != "" {
						params["per_page"] = v
					}
					raw, err := c.Request("GET", c.BooksBase+"/fixedassets", &zohttp.RequestOpts{Params: params})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "update",
				Usage:     "Update a fixed asset",
				ArgsUsage: "<asset-id>",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "asset_name", Usage: "Asset name"},
					&cli.StringFlag{Name: "asset_type_id", Usage: "Asset type ID"},
					&cli.StringFlag{Name: "purchase_date", Usage: "Purchase date (YYYY-MM-DD)"},
					&cli.FloatFlag{Name: "purchase_price", Usage: "Purchase price"},
					&cli.StringFlag{Name: "account_id", Usage: "Account ID"},
					&cli.StringFlag{Name: "depreciation_method", Usage: "Depreciation method"},
					&cli.FloatFlag{Name: "depreciation_rate", Usage: "Depreciation rate"},
					&cli.StringFlag{Name: "depreciation_account_id", Usage: "Depreciation account ID"},
					&cli.StringFlag{Name: "description", Usage: "Description"},
					&cli.StringFlag{Name: "json", Usage: "Additional fields as JSON"},
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
					body := map[string]any{}
					if cmd.IsSet("asset_name") {
						body["asset_name"] = cmd.String("asset_name")
					}
					if cmd.IsSet("asset_type_id") {
						body["asset_type_id"] = cmd.String("asset_type_id")
					}
					if cmd.IsSet("purchase_date") {
						body["purchase_date"] = cmd.String("purchase_date")
					}
					if cmd.IsSet("purchase_price") {
						body["purchase_price"] = cmd.Float("purchase_price")
					}
					if cmd.IsSet("account_id") {
						body["account_id"] = cmd.String("account_id")
					}
					if cmd.IsSet("depreciation_method") {
						body["depreciation_method"] = cmd.String("depreciation_method")
					}
					if cmd.IsSet("depreciation_rate") {
						body["depreciation_rate"] = cmd.Float("depreciation_rate")
					}
					if cmd.IsSet("depreciation_account_id") {
						body["depreciation_account_id"] = cmd.String("depreciation_account_id")
					}
					if cmd.IsSet("description") {
						body["description"] = cmd.String("description")
					}
					if err := internal.MergeJSON(cmd, body); err != nil {
						return err
					}
					raw, err := c.Request("PUT", c.BooksBase+"/fixedassets/"+cmd.Args().First(), &zohttp.RequestOpts{
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
				Usage:     "Get a fixed asset",
				ArgsUsage: "<asset-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("GET", c.BooksBase+"/fixedassets/"+cmd.Args().First(), &zohttp.RequestOpts{Params: orgParams(orgID)})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "delete",
				Usage:     "Delete a fixed asset",
				ArgsUsage: "<asset-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("DELETE", c.BooksBase+"/fixedassets/"+cmd.Args().First(), &zohttp.RequestOpts{Params: orgParams(orgID)})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "get-history",
				Usage:     "Get history of a fixed asset",
				ArgsUsage: "<asset-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("GET", c.BooksBase+"/fixedassets/"+cmd.Args().First()+"/comments", &zohttp.RequestOpts{Params: orgParams(orgID)})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "get-forecast-depreciation",
				Usage:     "Get forecast depreciation",
				ArgsUsage: "<asset-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("GET", c.BooksBase+"/fixedassets/"+cmd.Args().First()+"/depreciation", &zohttp.RequestOpts{Params: orgParams(orgID)})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "mark-active",
				Usage:     "Mark a fixed asset as active",
				ArgsUsage: "<asset-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("POST", c.BooksBase+"/fixedassets/"+cmd.Args().First()+"/active", &zohttp.RequestOpts{Params: orgParams(orgID)})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "cancel",
				Usage:     "Cancel a fixed asset",
				ArgsUsage: "<asset-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("POST", c.BooksBase+"/fixedassets/"+cmd.Args().First()+"/cancel", &zohttp.RequestOpts{Params: orgParams(orgID)})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "mark-draft",
				Usage:     "Mark a fixed asset as draft",
				ArgsUsage: "<asset-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("POST", c.BooksBase+"/fixedassets/"+cmd.Args().First()+"/draft", &zohttp.RequestOpts{Params: orgParams(orgID)})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "write-off",
				Usage:     "Write off a fixed asset",
				ArgsUsage: "<asset-id>",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "date", Required: true, Usage: "Date (YYYY-MM-DD)"},
					&cli.StringFlag{Name: "account_id", Required: true, Usage: "Account ID"},
					&cli.StringFlag{Name: "description", Usage: "Description"},
					&cli.StringFlag{Name: "notes", Usage: "Notes"},
					&cli.StringFlag{Name: "json", Usage: "Additional fields as JSON"},
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
					body := map[string]any{}
					body["date"] = cmd.String("date")
					body["account_id"] = cmd.String("account_id")
					if cmd.IsSet("description") {
						body["description"] = cmd.String("description")
					}
					if cmd.IsSet("notes") {
						body["notes"] = cmd.String("notes")
					}
					if err := internal.MergeJSON(cmd, body); err != nil {
						return err
					}
					raw, err := c.Request("POST", c.BooksBase+"/fixedassets/"+cmd.Args().First()+"/writeoff", &zohttp.RequestOpts{
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
				Name:      "sell",
				Usage:     "Sell a fixed asset",
				ArgsUsage: "<asset-id>",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "date", Required: true, Usage: "Date (YYYY-MM-DD)"},
					&cli.FloatFlag{Name: "selling_price", Required: true, Usage: "Selling price"},
					&cli.StringFlag{Name: "account_id", Required: true, Usage: "Account ID"},
					&cli.StringFlag{Name: "customer_id", Usage: "Customer ID"},
					&cli.StringFlag{Name: "description", Usage: "Description"},
					&cli.StringFlag{Name: "notes", Usage: "Notes"},
					&cli.StringFlag{Name: "json", Usage: "Additional fields as JSON"},
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
					body := map[string]any{}
					body["date"] = cmd.String("date")
					body["selling_price"] = cmd.Float("selling_price")
					body["account_id"] = cmd.String("account_id")
					if cmd.IsSet("customer_id") {
						body["customer_id"] = cmd.String("customer_id")
					}
					if cmd.IsSet("description") {
						body["description"] = cmd.String("description")
					}
					if cmd.IsSet("notes") {
						body["notes"] = cmd.String("notes")
					}
					if err := internal.MergeJSON(cmd, body); err != nil {
						return err
					}
					raw, err := c.Request("POST", c.BooksBase+"/fixedassets/"+cmd.Args().First()+"/sell", &zohttp.RequestOpts{
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
				Name:      "add-comment",
				Usage:     "Add a comment to a fixed asset",
				ArgsUsage: "<asset-id>",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "description", Required: true, Usage: "Description"},
					&cli.StringFlag{Name: "json", Usage: "Additional fields as JSON"},
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
					body := map[string]any{}
					body["description"] = cmd.String("description")
					if err := internal.MergeJSON(cmd, body); err != nil {
						return err
					}
					raw, err := c.Request("POST", c.BooksBase+"/fixedassets/"+cmd.Args().First()+"/comments", &zohttp.RequestOpts{
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
				Usage:     "Delete a comment on a fixed asset",
				ArgsUsage: "<asset-id> <comment-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("DELETE", c.BooksBase+"/fixedassets/"+cmd.Args().First()+"/comments/"+cmd.Args().Get(1), &zohttp.RequestOpts{Params: orgParams(orgID)})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:  "create-type",
				Usage: "Create a fixed asset type",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "asset_type_name", Required: true, Usage: "Asset type name"},
					&cli.StringFlag{Name: "asset_account_id", Usage: "Asset account ID"},
					&cli.StringFlag{Name: "depreciation_account_id", Usage: "Depreciation account ID"},
					&cli.StringFlag{Name: "depreciation_method", Usage: "Depreciation method"},
					&cli.FloatFlag{Name: "depreciation_rate", Usage: "Depreciation rate"},
					&cli.StringFlag{Name: "description", Usage: "Description"},
					&cli.StringFlag{Name: "json", Usage: "Additional fields as JSON"},
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
					body := map[string]any{}
					body["asset_type_name"] = cmd.String("asset_type_name")
					if cmd.IsSet("asset_account_id") {
						body["asset_account_id"] = cmd.String("asset_account_id")
					}
					if cmd.IsSet("depreciation_account_id") {
						body["depreciation_account_id"] = cmd.String("depreciation_account_id")
					}
					if cmd.IsSet("depreciation_method") {
						body["depreciation_method"] = cmd.String("depreciation_method")
					}
					if cmd.IsSet("depreciation_rate") {
						body["depreciation_rate"] = cmd.Float("depreciation_rate")
					}
					if cmd.IsSet("description") {
						body["description"] = cmd.String("description")
					}
					if err := internal.MergeJSON(cmd, body); err != nil {
						return err
					}
					raw, err := c.Request("POST", c.BooksBase+"/fixedassets"+"/types", &zohttp.RequestOpts{
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
				Name:  "list-types",
				Usage: "List fixed asset types",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("GET", c.BooksBase+"/fixedassets"+"/types", &zohttp.RequestOpts{Params: orgParams(orgID)})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "update-type",
				Usage:     "Update a fixed asset type",
				ArgsUsage: "<type-id>",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "asset_type_name", Usage: "Asset type name"},
					&cli.StringFlag{Name: "asset_account_id", Usage: "Asset account ID"},
					&cli.StringFlag{Name: "depreciation_account_id", Usage: "Depreciation account ID"},
					&cli.StringFlag{Name: "depreciation_method", Usage: "Depreciation method"},
					&cli.FloatFlag{Name: "depreciation_rate", Usage: "Depreciation rate"},
					&cli.StringFlag{Name: "description", Usage: "Description"},
					&cli.StringFlag{Name: "json", Usage: "Additional fields as JSON"},
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
					body := map[string]any{}
					if cmd.IsSet("asset_type_name") {
						body["asset_type_name"] = cmd.String("asset_type_name")
					}
					if cmd.IsSet("asset_account_id") {
						body["asset_account_id"] = cmd.String("asset_account_id")
					}
					if cmd.IsSet("depreciation_account_id") {
						body["depreciation_account_id"] = cmd.String("depreciation_account_id")
					}
					if cmd.IsSet("depreciation_method") {
						body["depreciation_method"] = cmd.String("depreciation_method")
					}
					if cmd.IsSet("depreciation_rate") {
						body["depreciation_rate"] = cmd.Float("depreciation_rate")
					}
					if cmd.IsSet("description") {
						body["description"] = cmd.String("description")
					}
					if err := internal.MergeJSON(cmd, body); err != nil {
						return err
					}
					raw, err := c.Request("PUT", c.BooksBase+"/fixedassets"+"/types/"+cmd.Args().First(), &zohttp.RequestOpts{
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
				Name:      "delete-type",
				Usage:     "Delete a fixed asset type",
				ArgsUsage: "<type-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("DELETE", c.BooksBase+"/fixedassets"+"/types/"+cmd.Args().First(), &zohttp.RequestOpts{Params: orgParams(orgID)})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
		},
	}
}

func baseCurrencyAdjCmd() *cli.Command {
	return &cli.Command{
		Name:  "base-currency-adjustment",
		Usage: "Base currency adjustment operations",
		Commands: []*cli.Command{
			{
				Name:  "create",
				Usage: "Create a base currency adjustment",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "adjustment_date", Required: true, Usage: "Adjustment date (YYYY-MM-DD)"},
					&cli.StringFlag{Name: "currency_id", Required: true, Usage: "Currency ID"},
					&cli.FloatFlag{Name: "exchange_rate", Required: true, Usage: "Exchange rate"},
					&cli.StringFlag{Name: "notes", Usage: "Notes"},
					&cli.StringFlag{Name: "json", Usage: "Additional fields as JSON"},
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
					body := map[string]any{}
					body["adjustment_date"] = cmd.String("adjustment_date")
					body["currency_id"] = cmd.String("currency_id")
					body["exchange_rate"] = cmd.Float("exchange_rate")
					if cmd.IsSet("notes") {
						body["notes"] = cmd.String("notes")
					}
					if err := internal.MergeJSON(cmd, body); err != nil {
						return err
					}
					raw, err := c.Request("POST", c.BooksBase+"/basecurrencyadjustment", &zohttp.RequestOpts{
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
				Usage: "List base currency adjustments",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "page", Usage: "Filter"},
					&cli.StringFlag{Name: "per-page", Usage: "Filter"},
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
					params := orgParams(orgID)
					if v := cmd.String("page"); v != "" {
						params["page"] = v
					}
					if v := cmd.String("per-page"); v != "" {
						params["per_page"] = v
					}
					raw, err := c.Request("GET", c.BooksBase+"/basecurrencyadjustment", &zohttp.RequestOpts{Params: params})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "get",
				Usage:     "Get a base currency adjustment",
				ArgsUsage: "<adjustment-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("GET", c.BooksBase+"/basecurrencyadjustment/"+cmd.Args().First(), &zohttp.RequestOpts{Params: orgParams(orgID)})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "delete",
				Usage:     "Delete a base currency adjustment",
				ArgsUsage: "<adjustment-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("DELETE", c.BooksBase+"/basecurrencyadjustment/"+cmd.Args().First(), &zohttp.RequestOpts{Params: orgParams(orgID)})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:  "list-account-details",
				Usage: "List account details for adjustment",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("GET", c.BooksBase+"/basecurrencyadjustment"+"/accounts", &zohttp.RequestOpts{Params: orgParams(orgID)})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
		},
	}
}

func booksProjectsCmd() *cli.Command {
	return &cli.Command{
		Name:  "projects",
		Usage: "Project operations",
		Commands: []*cli.Command{
			{
				Name:  "create",
				Usage: "Create a project",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "project_name", Required: true, Usage: "Project name"},
					&cli.StringFlag{Name: "customer_id", Required: true, Usage: "Customer ID"},
					&cli.StringFlag{Name: "currency_id", Required: true, Usage: "Currency ID"},
					&cli.StringFlag{Name: "billing_type", Required: true, Usage: "Billing type"},
					&cli.StringFlag{Name: "description", Usage: "Description"},
					&cli.FloatFlag{Name: "rate", Usage: "Rate"},
					&cli.StringFlag{Name: "budget_type", Usage: "Budget type"},
					&cli.FloatFlag{Name: "budget_hours", Usage: "Budget hours"},
					&cli.FloatFlag{Name: "budget_amount", Usage: "Budget amount"},
					&cli.FloatFlag{Name: "cost_budget_amount", Usage: "Cost budget amount"},
					&cli.StringFlag{Name: "user_id", Usage: "User ID"},
					&cli.StringFlag{Name: "json", Usage: "Additional fields as JSON"},
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
					body := map[string]any{}
					body["project_name"] = cmd.String("project_name")
					body["customer_id"] = cmd.String("customer_id")
					body["currency_id"] = cmd.String("currency_id")
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
						body["budget_hours"] = cmd.Float("budget_hours")
					}
					if cmd.IsSet("budget_amount") {
						body["budget_amount"] = cmd.Float("budget_amount")
					}
					if cmd.IsSet("cost_budget_amount") {
						body["cost_budget_amount"] = cmd.Float("cost_budget_amount")
					}
					if cmd.IsSet("user_id") {
						body["user_id"] = cmd.String("user_id")
					}
					if err := internal.MergeJSON(cmd, body); err != nil {
						return err
					}
					raw, err := c.Request("POST", c.BooksBase+"/projects", &zohttp.RequestOpts{
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
				Usage: "Update a project by custom field",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "project_name", Required: true, Usage: "Project name"},
					&cli.StringFlag{Name: "customer_id", Required: true, Usage: "Customer ID"},
					&cli.StringFlag{Name: "currency_id", Required: true, Usage: "Currency ID"},
					&cli.StringFlag{Name: "billing_type", Required: true, Usage: "Billing type"},
					&cli.StringFlag{Name: "description", Usage: "Description"},
					&cli.FloatFlag{Name: "rate", Usage: "Rate"},
					&cli.StringFlag{Name: "budget_type", Usage: "Budget type"},
					&cli.FloatFlag{Name: "budget_hours", Usage: "Budget hours"},
					&cli.FloatFlag{Name: "budget_amount", Usage: "Budget amount"},
					&cli.FloatFlag{Name: "cost_budget_amount", Usage: "Cost budget amount"},
					&cli.StringFlag{Name: "user_id", Usage: "User ID"},
					&cli.StringFlag{Name: "json", Usage: "Additional fields as JSON"},
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
					body := map[string]any{}
					body["project_name"] = cmd.String("project_name")
					body["customer_id"] = cmd.String("customer_id")
					body["currency_id"] = cmd.String("currency_id")
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
						body["budget_hours"] = cmd.Float("budget_hours")
					}
					if cmd.IsSet("budget_amount") {
						body["budget_amount"] = cmd.Float("budget_amount")
					}
					if cmd.IsSet("cost_budget_amount") {
						body["cost_budget_amount"] = cmd.Float("cost_budget_amount")
					}
					if cmd.IsSet("user_id") {
						body["user_id"] = cmd.String("user_id")
					}
					if err := internal.MergeJSON(cmd, body); err != nil {
						return err
					}
					raw, err := c.Request("PUT", c.BooksBase+"/projects", &zohttp.RequestOpts{
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
				Usage: "List projects",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "status", Usage: "Filter"},
					&cli.StringFlag{Name: "customer-id", Usage: "Filter"},
					&cli.StringFlag{Name: "page", Usage: "Filter"},
					&cli.StringFlag{Name: "per-page", Usage: "Filter"},
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
					params := orgParams(orgID)
					if v := cmd.String("status"); v != "" {
						params["status"] = v
					}
					if v := cmd.String("customer-id"); v != "" {
						params["customer_id"] = v
					}
					if v := cmd.String("page"); v != "" {
						params["page"] = v
					}
					if v := cmd.String("per-page"); v != "" {
						params["per_page"] = v
					}
					raw, err := c.Request("GET", c.BooksBase+"/projects", &zohttp.RequestOpts{Params: params})
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
					&cli.StringFlag{Name: "currency_id", Usage: "Currency ID"},
					&cli.StringFlag{Name: "billing_type", Usage: "Billing type"},
					&cli.StringFlag{Name: "description", Usage: "Description"},
					&cli.FloatFlag{Name: "rate", Usage: "Rate"},
					&cli.StringFlag{Name: "budget_type", Usage: "Budget type"},
					&cli.FloatFlag{Name: "budget_hours", Usage: "Budget hours"},
					&cli.FloatFlag{Name: "budget_amount", Usage: "Budget amount"},
					&cli.FloatFlag{Name: "cost_budget_amount", Usage: "Cost budget amount"},
					&cli.StringFlag{Name: "user_id", Usage: "User ID"},
					&cli.StringFlag{Name: "json", Usage: "Additional fields as JSON"},
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
					body := map[string]any{}
					if cmd.IsSet("project_name") {
						body["project_name"] = cmd.String("project_name")
					}
					if cmd.IsSet("customer_id") {
						body["customer_id"] = cmd.String("customer_id")
					}
					if cmd.IsSet("currency_id") {
						body["currency_id"] = cmd.String("currency_id")
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
						body["budget_hours"] = cmd.Float("budget_hours")
					}
					if cmd.IsSet("budget_amount") {
						body["budget_amount"] = cmd.Float("budget_amount")
					}
					if cmd.IsSet("cost_budget_amount") {
						body["cost_budget_amount"] = cmd.Float("cost_budget_amount")
					}
					if cmd.IsSet("user_id") {
						body["user_id"] = cmd.String("user_id")
					}
					if err := internal.MergeJSON(cmd, body); err != nil {
						return err
					}
					raw, err := c.Request("PUT", c.BooksBase+"/projects/"+cmd.Args().First(), &zohttp.RequestOpts{
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
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("GET", c.BooksBase+"/projects/"+cmd.Args().First(), &zohttp.RequestOpts{Params: orgParams(orgID)})
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
					raw, err := c.Request("DELETE", c.BooksBase+"/projects/"+cmd.Args().First(), &zohttp.RequestOpts{Params: orgParams(orgID)})
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
					raw, err := c.Request("POST", c.BooksBase+"/projects/"+cmd.Args().First()+"/active", &zohttp.RequestOpts{Params: orgParams(orgID)})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "inactivate",
				Usage:     "Inactivate a project",
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
					raw, err := c.Request("POST", c.BooksBase+"/projects/"+cmd.Args().First()+"/inactive", &zohttp.RequestOpts{Params: orgParams(orgID)})
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
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "project_name", Required: true, Usage: "Project name"},
					&cli.StringFlag{Name: "start_date", Usage: "Start date (YYYY-MM-DD)"},
					&cli.StringFlag{Name: "end_date", Usage: "End date (YYYY-MM-DD)"},
					&cli.BoolFlag{Name: "clone_tasks", Usage: "Clone tasks"},
					&cli.BoolFlag{Name: "clone_users", Usage: "Clone users"},
					&cli.StringFlag{Name: "json", Usage: "Additional fields as JSON"},
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
					body := map[string]any{}
					body["project_name"] = cmd.String("project_name")
					if cmd.IsSet("start_date") {
						body["start_date"] = cmd.String("start_date")
					}
					if cmd.IsSet("end_date") {
						body["end_date"] = cmd.String("end_date")
					}
					if cmd.IsSet("clone_tasks") {
						body["clone_tasks"] = cmd.Bool("clone_tasks")
					}
					if cmd.IsSet("clone_users") {
						body["clone_users"] = cmd.Bool("clone_users")
					}
					if err := internal.MergeJSON(cmd, body); err != nil {
						return err
					}
					opts := &zohttp.RequestOpts{Params: orgParams(orgID), JSON: body}
					raw, err := c.Request("POST", c.BooksBase+"/projects/"+cmd.Args().First()+"/clone", opts)
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "assign-users",
				Usage:     "Assign users to a project",
				ArgsUsage: "<project-id>",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "user_id", Required: true, Usage: "User ID"},
					&cli.FloatFlag{Name: "rate", Usage: "Rate"},
					&cli.FloatFlag{Name: "budget_hours", Usage: "Budget hours"},
					&cli.FloatFlag{Name: "cost_rate", Usage: "Cost rate"},
					&cli.StringFlag{Name: "user_role", Usage: "User role"},
					&cli.StringFlag{Name: "json", Usage: "Additional fields as JSON"},
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
					body := map[string]any{}
					body["user_id"] = cmd.String("user_id")
					if cmd.IsSet("rate") {
						body["rate"] = cmd.Float("rate")
					}
					if cmd.IsSet("budget_hours") {
						body["budget_hours"] = cmd.Float("budget_hours")
					}
					if cmd.IsSet("cost_rate") {
						body["cost_rate"] = cmd.Float("cost_rate")
					}
					if cmd.IsSet("user_role") {
						body["user_role"] = cmd.String("user_role")
					}
					if err := internal.MergeJSON(cmd, body); err != nil {
						return err
					}
					raw, err := c.Request("POST", c.BooksBase+"/projects/"+cmd.Args().First()+"/users", &zohttp.RequestOpts{
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
				Name:      "list-users",
				Usage:     "List users of a project",
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
					raw, err := c.Request("GET", c.BooksBase+"/projects/"+cmd.Args().First()+"/users", &zohttp.RequestOpts{Params: orgParams(orgID)})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "invite-user",
				Usage:     "Invite a user to a project",
				ArgsUsage: "<project-id>",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "user_id", Usage: "User ID"},
					&cli.StringFlag{Name: "email", Usage: "Email address"},
					&cli.StringFlag{Name: "user_name", Usage: "User name"},
					&cli.StringFlag{Name: "user_role", Usage: "User role"},
					&cli.FloatFlag{Name: "rate", Usage: "Rate"},
					&cli.FloatFlag{Name: "budget_hours", Usage: "Budget hours"},
					&cli.FloatFlag{Name: "cost_rate", Usage: "Cost rate"},
					&cli.StringFlag{Name: "json", Usage: "Additional fields as JSON"},
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
					body := map[string]any{}
					if cmd.IsSet("user_id") {
						body["user_id"] = cmd.String("user_id")
					}
					if cmd.IsSet("email") {
						body["email"] = cmd.String("email")
					}
					if cmd.IsSet("user_name") {
						body["user_name"] = cmd.String("user_name")
					}
					if cmd.IsSet("user_role") {
						body["user_role"] = cmd.String("user_role")
					}
					if cmd.IsSet("rate") {
						body["rate"] = cmd.Float("rate")
					}
					if cmd.IsSet("budget_hours") {
						body["budget_hours"] = cmd.Float("budget_hours")
					}
					if cmd.IsSet("cost_rate") {
						body["cost_rate"] = cmd.Float("cost_rate")
					}
					if err := internal.MergeJSON(cmd, body); err != nil {
						return err
					}
					raw, err := c.Request("POST", c.BooksBase+"/projects/"+cmd.Args().First()+"/users/invite", &zohttp.RequestOpts{
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
				Name:      "update-user",
				Usage:     "Update a user in a project",
				ArgsUsage: "<project-id> <user-id>",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "user_name", Usage: "User name"},
					&cli.StringFlag{Name: "user_role", Usage: "User role"},
					&cli.FloatFlag{Name: "rate", Usage: "Rate"},
					&cli.FloatFlag{Name: "budget_hours", Usage: "Budget hours"},
					&cli.FloatFlag{Name: "cost_rate", Usage: "Cost rate"},
					&cli.StringFlag{Name: "json", Usage: "Additional fields as JSON"},
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
					body := map[string]any{}
					if cmd.IsSet("user_name") {
						body["user_name"] = cmd.String("user_name")
					}
					if cmd.IsSet("user_role") {
						body["user_role"] = cmd.String("user_role")
					}
					if cmd.IsSet("rate") {
						body["rate"] = cmd.Float("rate")
					}
					if cmd.IsSet("budget_hours") {
						body["budget_hours"] = cmd.Float("budget_hours")
					}
					if cmd.IsSet("cost_rate") {
						body["cost_rate"] = cmd.Float("cost_rate")
					}
					if err := internal.MergeJSON(cmd, body); err != nil {
						return err
					}
					raw, err := c.Request("PUT", c.BooksBase+"/projects/"+cmd.Args().First()+"/users/"+cmd.Args().Get(1), &zohttp.RequestOpts{
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
				Name:      "get-user",
				Usage:     "Get a user in a project",
				ArgsUsage: "<project-id> <user-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("GET", c.BooksBase+"/projects/"+cmd.Args().First()+"/users/"+cmd.Args().Get(1), &zohttp.RequestOpts{Params: orgParams(orgID)})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "delete-user",
				Usage:     "Delete a user from a project",
				ArgsUsage: "<project-id> <user-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("DELETE", c.BooksBase+"/projects/"+cmd.Args().First()+"/users/"+cmd.Args().Get(1), &zohttp.RequestOpts{Params: orgParams(orgID)})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "post-comment",
				Usage:     "Post a comment on a project",
				ArgsUsage: "<project-id>",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "description", Required: true, Usage: "Description"},
					&cli.StringFlag{Name: "json", Usage: "Additional fields as JSON"},
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
					body := map[string]any{}
					body["description"] = cmd.String("description")
					if err := internal.MergeJSON(cmd, body); err != nil {
						return err
					}
					raw, err := c.Request("POST", c.BooksBase+"/projects/"+cmd.Args().First()+"/comments", &zohttp.RequestOpts{
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
				Usage:     "List comments of a project",
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
					raw, err := c.Request("GET", c.BooksBase+"/projects/"+cmd.Args().First()+"/comments", &zohttp.RequestOpts{Params: orgParams(orgID)})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "delete-comment",
				Usage:     "Delete a comment on a project",
				ArgsUsage: "<project-id> <comment-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("DELETE", c.BooksBase+"/projects/"+cmd.Args().First()+"/comments/"+cmd.Args().Get(1), &zohttp.RequestOpts{Params: orgParams(orgID)})
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
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("GET", c.BooksBase+"/projects/"+cmd.Args().First()+"/invoices", &zohttp.RequestOpts{Params: orgParams(orgID)})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
		},
	}
}

func tasksCmd() *cli.Command {
	return &cli.Command{
		Name:  "tasks",
		Usage: "Task operations",
		Commands: []*cli.Command{
			{
				Name:  "create",
				Usage: "Create a task",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "project-id", Required: true, Usage: "Project ID"},
					&cli.StringFlag{Name: "task_name", Required: true, Usage: "Task name"},
					&cli.StringFlag{Name: "description", Usage: "Description"},
					&cli.FloatFlag{Name: "rate", Usage: "Rate"},
					&cli.FloatFlag{Name: "budget_hours", Usage: "Budget hours"},
					&cli.StringFlag{Name: "json", Usage: "Additional fields as JSON"},
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
					body := map[string]any{}
					body["task_name"] = cmd.String("task_name")
					if cmd.IsSet("description") {
						body["description"] = cmd.String("description")
					}
					if cmd.IsSet("rate") {
						body["rate"] = cmd.Float("rate")
					}
					if cmd.IsSet("budget_hours") {
						body["budget_hours"] = cmd.Float("budget_hours")
					}
					if err := internal.MergeJSON(cmd, body); err != nil {
						return err
					}
					raw, err := c.Request("POST", c.BooksBase+"/projects/"+cmd.String("project-id")+"/tasks", &zohttp.RequestOpts{
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
				Usage: "List tasks",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "project-id", Required: true, Usage: "Project ID"},
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
					raw, err := c.Request("GET", c.BooksBase+"/projects/"+cmd.String("project-id")+"/tasks", &zohttp.RequestOpts{Params: orgParams(orgID)})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "update",
				Usage:     "Update a task",
				ArgsUsage: "<task-id>",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "project-id", Required: true, Usage: "Project ID"},
					&cli.StringFlag{Name: "task_name", Usage: "Task name"},
					&cli.StringFlag{Name: "description", Usage: "Description"},
					&cli.FloatFlag{Name: "rate", Usage: "Rate"},
					&cli.FloatFlag{Name: "budget_hours", Usage: "Budget hours"},
					&cli.StringFlag{Name: "json", Usage: "Additional fields as JSON"},
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
					body := map[string]any{}
					if cmd.IsSet("task_name") {
						body["task_name"] = cmd.String("task_name")
					}
					if cmd.IsSet("description") {
						body["description"] = cmd.String("description")
					}
					if cmd.IsSet("rate") {
						body["rate"] = cmd.Float("rate")
					}
					if cmd.IsSet("budget_hours") {
						body["budget_hours"] = cmd.Float("budget_hours")
					}
					if err := internal.MergeJSON(cmd, body); err != nil {
						return err
					}
					raw, err := c.Request("PUT", c.BooksBase+"/projects/"+cmd.String("project-id")+"/tasks/"+cmd.Args().First(), &zohttp.RequestOpts{
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
				Usage:     "Get a task",
				ArgsUsage: "<task-id>",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "project-id", Required: true, Usage: "Project ID"},
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
					raw, err := c.Request("GET", c.BooksBase+"/projects/"+cmd.String("project-id")+"/tasks/"+cmd.Args().First(), &zohttp.RequestOpts{Params: orgParams(orgID)})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "delete",
				Usage:     "Delete a task",
				ArgsUsage: "<task-id>",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "project-id", Required: true, Usage: "Project ID"},
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
					raw, err := c.Request("DELETE", c.BooksBase+"/projects/"+cmd.String("project-id")+"/tasks/"+cmd.Args().First(), &zohttp.RequestOpts{Params: orgParams(orgID)})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
		},
	}
}

func timeEntriesCmd() *cli.Command {
	return &cli.Command{
		Name:  "time-entries",
		Usage: "Time entry operations",
		Commands: []*cli.Command{
			{
				Name:  "log",
				Usage: "Log a time entry",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "project-id", Required: true, Usage: "Project ID"},
					&cli.StringFlag{Name: "user_id", Required: true, Usage: "User ID"},
					&cli.StringFlag{Name: "task_id", Usage: "Task ID"},
					&cli.StringFlag{Name: "log_date", Required: true, Usage: "Log date (YYYY-MM-DD)"},
					&cli.StringFlag{Name: "log_time", Usage: "Log time (HH:mm)"},
					&cli.StringFlag{Name: "begin_time", Usage: "Begin time (HH:mm)"},
					&cli.StringFlag{Name: "end_time", Usage: "End time (HH:mm)"},
					&cli.BoolFlag{Name: "is_billable", Usage: "Is billable"},
					&cli.StringFlag{Name: "notes", Usage: "Notes"},
					&cli.BoolFlag{Name: "start_timer", Usage: "Start timer"},
					&cli.FloatFlag{Name: "cost_rate", Usage: "Cost rate"},
					&cli.StringFlag{Name: "json", Usage: "Additional fields as JSON"},
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
					body := map[string]any{}
					body["user_id"] = cmd.String("user_id")
					if cmd.IsSet("task_id") {
						body["task_id"] = cmd.String("task_id")
					}
					body["log_date"] = cmd.String("log_date")
					if cmd.IsSet("log_time") {
						body["log_time"] = cmd.String("log_time")
					}
					if cmd.IsSet("begin_time") {
						body["begin_time"] = cmd.String("begin_time")
					}
					if cmd.IsSet("end_time") {
						body["end_time"] = cmd.String("end_time")
					}
					if cmd.IsSet("is_billable") {
						body["is_billable"] = cmd.Bool("is_billable")
					}
					if cmd.IsSet("notes") {
						body["notes"] = cmd.String("notes")
					}
					if cmd.IsSet("start_timer") {
						body["start_timer"] = cmd.Bool("start_timer")
					}
					if cmd.IsSet("cost_rate") {
						body["cost_rate"] = cmd.Float("cost_rate")
					}
					if err := internal.MergeJSON(cmd, body); err != nil {
						return err
					}
					raw, err := c.Request("POST", c.BooksBase+"/projects/"+cmd.String("project-id")+"/timeentries", &zohttp.RequestOpts{
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
				Usage: "List time entries",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "project-id", Required: true, Usage: "Project ID"},
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
					params := orgParams(orgID)
					if v := cmd.String("page"); v != "" {
						params["page"] = v
					}
					if v := cmd.String("per-page"); v != "" {
						params["per_page"] = v
					}
					raw, err := c.Request("GET", c.BooksBase+"/projects/"+cmd.String("project-id")+"/timeentries", &zohttp.RequestOpts{Params: params})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "update",
				Usage:     "Update a time entry",
				ArgsUsage: "<timeentry-id>",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "project-id", Required: true, Usage: "Project ID"},
					&cli.StringFlag{Name: "user_id", Usage: "User ID"},
					&cli.StringFlag{Name: "task_id", Usage: "Task ID"},
					&cli.StringFlag{Name: "log_date", Usage: "Log date (YYYY-MM-DD)"},
					&cli.StringFlag{Name: "log_time", Usage: "Log time (HH:mm)"},
					&cli.StringFlag{Name: "begin_time", Usage: "Begin time (HH:mm)"},
					&cli.StringFlag{Name: "end_time", Usage: "End time (HH:mm)"},
					&cli.BoolFlag{Name: "is_billable", Usage: "Is billable"},
					&cli.StringFlag{Name: "notes", Usage: "Notes"},
					&cli.BoolFlag{Name: "start_timer", Usage: "Start timer"},
					&cli.FloatFlag{Name: "cost_rate", Usage: "Cost rate"},
					&cli.StringFlag{Name: "json", Usage: "Additional fields as JSON"},
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
					body := map[string]any{}
					if cmd.IsSet("user_id") {
						body["user_id"] = cmd.String("user_id")
					}
					if cmd.IsSet("task_id") {
						body["task_id"] = cmd.String("task_id")
					}
					if cmd.IsSet("log_date") {
						body["log_date"] = cmd.String("log_date")
					}
					if cmd.IsSet("log_time") {
						body["log_time"] = cmd.String("log_time")
					}
					if cmd.IsSet("begin_time") {
						body["begin_time"] = cmd.String("begin_time")
					}
					if cmd.IsSet("end_time") {
						body["end_time"] = cmd.String("end_time")
					}
					if cmd.IsSet("is_billable") {
						body["is_billable"] = cmd.Bool("is_billable")
					}
					if cmd.IsSet("notes") {
						body["notes"] = cmd.String("notes")
					}
					if cmd.IsSet("start_timer") {
						body["start_timer"] = cmd.Bool("start_timer")
					}
					if cmd.IsSet("cost_rate") {
						body["cost_rate"] = cmd.Float("cost_rate")
					}
					if err := internal.MergeJSON(cmd, body); err != nil {
						return err
					}
					raw, err := c.Request("PUT", c.BooksBase+"/projects/"+cmd.String("project-id")+"/timeentries/"+cmd.Args().First(), &zohttp.RequestOpts{
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
				Usage:     "Get a time entry",
				ArgsUsage: "<timeentry-id>",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "project-id", Required: true, Usage: "Project ID"},
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
					raw, err := c.Request("GET", c.BooksBase+"/projects/"+cmd.String("project-id")+"/timeentries/"+cmd.Args().First(), &zohttp.RequestOpts{Params: orgParams(orgID)})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "delete",
				Usage:     "Delete a time entry",
				ArgsUsage: "<timeentry-id>",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "project-id", Required: true, Usage: "Project ID"},
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
					raw, err := c.Request("DELETE", c.BooksBase+"/projects/"+cmd.String("project-id")+"/timeentries/"+cmd.Args().First(), &zohttp.RequestOpts{Params: orgParams(orgID)})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "start-timer",
				Usage:     "Start a timer for a time entry",
				ArgsUsage: "<timeentry-id>",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "project-id", Required: true, Usage: "Project ID"},
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
					raw, err := c.Request("POST", c.BooksBase+"/projects/"+cmd.String("project-id")+"/timeentries/"+cmd.Args().First()+"/timer/start", &zohttp.RequestOpts{Params: orgParams(orgID)})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "stop-timer",
				Usage:     "Stop a timer for a time entry",
				ArgsUsage: "<timeentry-id>",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "project-id", Required: true, Usage: "Project ID"},
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
					raw, err := c.Request("POST", c.BooksBase+"/projects/"+cmd.String("project-id")+"/timeentries/"+cmd.Args().First()+"/timer/stop", &zohttp.RequestOpts{Params: orgParams(orgID)})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:  "get-timer",
				Usage: "Get current running timer",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "project-id", Required: true, Usage: "Project ID"},
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
					raw, err := c.Request("GET", c.BooksBase+"/projects/"+cmd.String("project-id")+"/timeentries/timer", &zohttp.RequestOpts{Params: orgParams(orgID)})
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
				Name:  "create",
				Usage: "Create a user",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "user_name", Required: true, Usage: "User name"},
					&cli.StringFlag{Name: "email", Required: true, Usage: "Email address"},
					&cli.StringFlag{Name: "user_role", Required: true, Usage: "User role"},
					&cli.StringFlag{Name: "status", Usage: "Status"},
					&cli.StringFlag{Name: "json", Usage: "Additional fields as JSON"},
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
					body := map[string]any{}
					body["user_name"] = cmd.String("user_name")
					body["email"] = cmd.String("email")
					body["user_role"] = cmd.String("user_role")
					if cmd.IsSet("status") {
						body["status"] = cmd.String("status")
					}
					if err := internal.MergeJSON(cmd, body); err != nil {
						return err
					}
					raw, err := c.Request("POST", c.BooksBase+"/users", &zohttp.RequestOpts{
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
				Usage: "List users",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "page", Usage: "Filter"},
					&cli.StringFlag{Name: "per-page", Usage: "Filter"},
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
					params := orgParams(orgID)
					if v := cmd.String("page"); v != "" {
						params["page"] = v
					}
					if v := cmd.String("per-page"); v != "" {
						params["per_page"] = v
					}
					raw, err := c.Request("GET", c.BooksBase+"/users", &zohttp.RequestOpts{Params: params})
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
					&cli.StringFlag{Name: "user_name", Usage: "User name"},
					&cli.StringFlag{Name: "email", Usage: "Email address"},
					&cli.StringFlag{Name: "user_role", Usage: "User role"},
					&cli.StringFlag{Name: "status", Usage: "Status"},
					&cli.StringFlag{Name: "json", Usage: "Additional fields as JSON"},
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
					body := map[string]any{}
					if cmd.IsSet("user_name") {
						body["user_name"] = cmd.String("user_name")
					}
					if cmd.IsSet("email") {
						body["email"] = cmd.String("email")
					}
					if cmd.IsSet("user_role") {
						body["user_role"] = cmd.String("user_role")
					}
					if cmd.IsSet("status") {
						body["status"] = cmd.String("status")
					}
					if err := internal.MergeJSON(cmd, body); err != nil {
						return err
					}
					raw, err := c.Request("PUT", c.BooksBase+"/users/"+cmd.Args().First(), &zohttp.RequestOpts{
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
					raw, err := c.Request("GET", c.BooksBase+"/users/"+cmd.Args().First(), &zohttp.RequestOpts{Params: orgParams(orgID)})
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
					raw, err := c.Request("DELETE", c.BooksBase+"/users/"+cmd.Args().First(), &zohttp.RequestOpts{Params: orgParams(orgID)})
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
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("GET", c.BooksBase+"/users"+"/me", &zohttp.RequestOpts{Params: orgParams(orgID)})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "invite",
				Usage:     "Invite a user",
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
					raw, err := c.Request("POST", c.BooksBase+"/users/"+cmd.Args().First()+"/invite", &zohttp.RequestOpts{Params: orgParams(orgID)})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "mark-active",
				Usage:     "Mark a user as active",
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
					raw, err := c.Request("POST", c.BooksBase+"/users/"+cmd.Args().First()+"/active", &zohttp.RequestOpts{Params: orgParams(orgID)})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "mark-inactive",
				Usage:     "Mark a user as inactive",
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
					raw, err := c.Request("POST", c.BooksBase+"/users/"+cmd.Args().First()+"/inactive", &zohttp.RequestOpts{Params: orgParams(orgID)})
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
				Name:  "create",
				Usage: "Create an item",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "name", Required: true, Usage: "Item name"},
					&cli.FloatFlag{Name: "rate", Required: true, Usage: "Item rate"},
					&cli.StringFlag{Name: "description", Usage: "Item description"},
					&cli.StringFlag{Name: "sku", Usage: "Item SKU"},
					&cli.StringFlag{Name: "json", Usage: "Additional fields as JSON"},
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
					body := map[string]any{}
					body["name"] = cmd.String("name")
					body["rate"] = cmd.Float("rate")
					if cmd.IsSet("description") {
						body["description"] = cmd.String("description")
					}
					if cmd.IsSet("sku") {
						body["sku"] = cmd.String("sku")
					}
					if err := internal.MergeJSON(cmd, body); err != nil {
						return err
					}
					raw, err := c.Request("POST", c.BooksBase+"/items", &zohttp.RequestOpts{
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
				Usage: "Update an item by custom field",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "name", Usage: "Item name"},
					&cli.StringFlag{Name: "description", Usage: "Description"},
					&cli.FloatFlag{Name: "rate", Usage: "Rate"},
					&cli.FloatFlag{Name: "purchase_rate", Usage: "Purchase rate"},
					&cli.StringFlag{Name: "sku", Usage: "Item SKU"},
					&cli.StringFlag{Name: "item_type", Usage: "Item type"},
					&cli.StringFlag{Name: "account_id", Usage: "Account ID"},
					&cli.StringFlag{Name: "tax_id", Usage: "Tax ID"},
					&cli.StringFlag{Name: "json", Usage: "Additional fields as JSON"},
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
					body := map[string]any{}
					if cmd.IsSet("name") {
						body["name"] = cmd.String("name")
					}
					if cmd.IsSet("description") {
						body["description"] = cmd.String("description")
					}
					if cmd.IsSet("rate") {
						body["rate"] = cmd.Float("rate")
					}
					if cmd.IsSet("purchase_rate") {
						body["purchase_rate"] = cmd.Float("purchase_rate")
					}
					if cmd.IsSet("sku") {
						body["sku"] = cmd.String("sku")
					}
					if cmd.IsSet("item_type") {
						body["item_type"] = cmd.String("item_type")
					}
					if cmd.IsSet("account_id") {
						body["account_id"] = cmd.String("account_id")
					}
					if cmd.IsSet("tax_id") {
						body["tax_id"] = cmd.String("tax_id")
					}
					if err := internal.MergeJSON(cmd, body); err != nil {
						return err
					}
					raw, err := c.Request("PUT", c.BooksBase+"/items", &zohttp.RequestOpts{
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
				Usage: "List items",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "name", Usage: "Filter"},
					&cli.StringFlag{Name: "status", Usage: "Filter"},
					&cli.StringFlag{Name: "page", Usage: "Filter"},
					&cli.StringFlag{Name: "per-page", Usage: "Filter"},
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
					params := orgParams(orgID)
					if v := cmd.String("name"); v != "" {
						params["name"] = v
					}
					if v := cmd.String("status"); v != "" {
						params["status"] = v
					}
					if v := cmd.String("page"); v != "" {
						params["page"] = v
					}
					if v := cmd.String("per-page"); v != "" {
						params["per_page"] = v
					}
					raw, err := c.Request("GET", c.BooksBase+"/items", &zohttp.RequestOpts{Params: params})
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
					&cli.StringFlag{Name: "sku", Usage: "Item SKU"},
					&cli.StringFlag{Name: "json", Usage: "Additional fields as JSON"},
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
					if cmd.IsSet("sku") {
						body["sku"] = cmd.String("sku")
					}
					if err := internal.MergeJSON(cmd, body); err != nil {
						return err
					}
					raw, err := c.Request("PUT", c.BooksBase+"/items/"+cmd.Args().First(), &zohttp.RequestOpts{
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
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("GET", c.BooksBase+"/items/"+cmd.Args().First(), &zohttp.RequestOpts{Params: orgParams(orgID)})
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
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("DELETE", c.BooksBase+"/items/"+cmd.Args().First(), &zohttp.RequestOpts{Params: orgParams(orgID)})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "update-custom-fields",
				Usage:     "Update custom fields of an item",
				ArgsUsage: "<item-id>",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "customfield_id", Usage: "Custom field ID"},
					&cli.StringFlag{Name: "value", Usage: "Custom field value"},
					&cli.StringFlag{Name: "json", Usage: "Additional fields as JSON"},
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
					raw, err := c.Request("PUT", c.BooksBase+"/items/"+cmd.Args().First()+"/customfields", &zohttp.RequestOpts{
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
				Name:      "mark-active",
				Usage:     "Mark an item as active",
				ArgsUsage: "<item-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("POST", c.BooksBase+"/items/"+cmd.Args().First()+"/active", &zohttp.RequestOpts{Params: orgParams(orgID)})
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
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("POST", c.BooksBase+"/items/"+cmd.Args().First()+"/inactive", &zohttp.RequestOpts{Params: orgParams(orgID)})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
		},
	}
}

func locationsCmd() *cli.Command {
	return &cli.Command{
		Name:  "locations",
		Usage: "Location operations",
		Commands: []*cli.Command{
			{
				Name:  "enable",
				Usage: "Enable locations",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("POST", c.BooksBase+"/locations"+"/enable", &zohttp.RequestOpts{Params: orgParams(orgID)})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:  "create",
				Usage: "Create a location",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "location_name", Required: true, Usage: "Location name"},
					&cli.StringFlag{Name: "location_code", Usage: "Location code"},
					&cli.StringFlag{Name: "address", Usage: "Address"},
					&cli.StringFlag{Name: "city", Usage: "City"},
					&cli.StringFlag{Name: "state", Usage: "State"},
					&cli.StringFlag{Name: "zip", Usage: "Zip/postal code"},
					&cli.StringFlag{Name: "country", Usage: "Country"},
					&cli.StringFlag{Name: "phone", Usage: "Phone number"},
					&cli.StringFlag{Name: "fax", Usage: "Fax number"},
					&cli.BoolFlag{Name: "is_primary", Usage: "Is primary location"},
					&cli.StringFlag{Name: "json", Usage: "Additional fields as JSON"},
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
					body := map[string]any{}
					body["location_name"] = cmd.String("location_name")
					if cmd.IsSet("location_code") {
						body["location_code"] = cmd.String("location_code")
					}
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
					if cmd.IsSet("phone") {
						body["phone"] = cmd.String("phone")
					}
					if cmd.IsSet("fax") {
						body["fax"] = cmd.String("fax")
					}
					if cmd.IsSet("is_primary") {
						body["is_primary"] = cmd.Bool("is_primary")
					}
					if err := internal.MergeJSON(cmd, body); err != nil {
						return err
					}
					raw, err := c.Request("POST", c.BooksBase+"/locations", &zohttp.RequestOpts{
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
				Usage: "List locations",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("GET", c.BooksBase+"/locations", &zohttp.RequestOpts{Params: orgParams(orgID)})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "update",
				Usage:     "Update a location",
				ArgsUsage: "<location-id>",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "location_name", Usage: "Location name"},
					&cli.StringFlag{Name: "location_code", Usage: "Location code"},
					&cli.StringFlag{Name: "address", Usage: "Address"},
					&cli.StringFlag{Name: "city", Usage: "City"},
					&cli.StringFlag{Name: "state", Usage: "State"},
					&cli.StringFlag{Name: "zip", Usage: "Zip/postal code"},
					&cli.StringFlag{Name: "country", Usage: "Country"},
					&cli.StringFlag{Name: "phone", Usage: "Phone number"},
					&cli.StringFlag{Name: "fax", Usage: "Fax number"},
					&cli.BoolFlag{Name: "is_primary", Usage: "Is primary location"},
					&cli.StringFlag{Name: "json", Usage: "Additional fields as JSON"},
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
					body := map[string]any{}
					if cmd.IsSet("location_name") {
						body["location_name"] = cmd.String("location_name")
					}
					if cmd.IsSet("location_code") {
						body["location_code"] = cmd.String("location_code")
					}
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
					if cmd.IsSet("phone") {
						body["phone"] = cmd.String("phone")
					}
					if cmd.IsSet("fax") {
						body["fax"] = cmd.String("fax")
					}
					if cmd.IsSet("is_primary") {
						body["is_primary"] = cmd.Bool("is_primary")
					}
					if err := internal.MergeJSON(cmd, body); err != nil {
						return err
					}
					raw, err := c.Request("PUT", c.BooksBase+"/locations/"+cmd.Args().First(), &zohttp.RequestOpts{
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
				Usage:     "Delete a location",
				ArgsUsage: "<location-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("DELETE", c.BooksBase+"/locations/"+cmd.Args().First(), &zohttp.RequestOpts{Params: orgParams(orgID)})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "mark-active",
				Usage:     "Mark a location as active",
				ArgsUsage: "<location-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("POST", c.BooksBase+"/locations/"+cmd.Args().First()+"/active", &zohttp.RequestOpts{Params: orgParams(orgID)})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "mark-inactive",
				Usage:     "Mark a location as inactive",
				ArgsUsage: "<location-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("POST", c.BooksBase+"/locations/"+cmd.Args().First()+"/inactive", &zohttp.RequestOpts{Params: orgParams(orgID)})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "mark-primary",
				Usage:     "Mark a location as primary",
				ArgsUsage: "<location-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("POST", c.BooksBase+"/locations/"+cmd.Args().First()+"/primary", &zohttp.RequestOpts{Params: orgParams(orgID)})
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
				Name:  "create",
				Usage: "Create a currency",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "currency_code", Required: true, Usage: "Currency code"},
					&cli.StringFlag{Name: "currency_symbol", Usage: "Currency symbol"},
					&cli.StringFlag{Name: "json", Usage: "Additional fields as JSON"},
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
					body := map[string]any{}
					body["currency_code"] = cmd.String("currency_code")
					if cmd.IsSet("currency_symbol") {
						body["currency_symbol"] = cmd.String("currency_symbol")
					}
					if err := internal.MergeJSON(cmd, body); err != nil {
						return err
					}
					raw, err := c.Request("POST", c.BooksBase+"/settings/currencies", &zohttp.RequestOpts{
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
					raw, err := c.Request("GET", c.BooksBase+"/settings/currencies", &zohttp.RequestOpts{Params: orgParams(orgID)})
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
					&cli.StringFlag{Name: "json", Usage: "Additional fields as JSON"},
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
					body := map[string]any{}
					if cmd.IsSet("currency_code") {
						body["currency_code"] = cmd.String("currency_code")
					}
					if cmd.IsSet("currency_symbol") {
						body["currency_symbol"] = cmd.String("currency_symbol")
					}
					if err := internal.MergeJSON(cmd, body); err != nil {
						return err
					}
					raw, err := c.Request("PUT", c.BooksBase+"/settings/currencies/"+cmd.Args().First(), &zohttp.RequestOpts{
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
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("GET", c.BooksBase+"/settings/currencies/"+cmd.Args().First(), &zohttp.RequestOpts{Params: orgParams(orgID)})
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
					raw, err := c.Request("DELETE", c.BooksBase+"/settings/currencies/"+cmd.Args().First(), &zohttp.RequestOpts{Params: orgParams(orgID)})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "create-exchange-rate",
				Usage:     "Create an exchange rate",
				ArgsUsage: "<currency-id>",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "effective_date", Required: true, Usage: "Effective date (YYYY-MM-DD)"},
					&cli.FloatFlag{Name: "rate", Required: true, Usage: "Rate"},
					&cli.StringFlag{Name: "json", Usage: "Additional fields as JSON"},
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
					body := map[string]any{}
					body["effective_date"] = cmd.String("effective_date")
					body["rate"] = cmd.Float("rate")
					if err := internal.MergeJSON(cmd, body); err != nil {
						return err
					}
					raw, err := c.Request("POST", c.BooksBase+"/settings/currencies/"+cmd.Args().First()+"/exchangerates", &zohttp.RequestOpts{
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
				Name:      "list-exchange-rates",
				Usage:     "List exchange rates",
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
					raw, err := c.Request("GET", c.BooksBase+"/settings/currencies/"+cmd.Args().First()+"/exchangerates", &zohttp.RequestOpts{Params: orgParams(orgID)})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "update-exchange-rate",
				Usage:     "Update an exchange rate",
				ArgsUsage: "<currency-id> <exchange-rate-id>",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "effective_date", Usage: "Effective date (YYYY-MM-DD)"},
					&cli.FloatFlag{Name: "rate", Usage: "Rate"},
					&cli.StringFlag{Name: "json", Usage: "Additional fields as JSON"},
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
					body := map[string]any{}
					if cmd.IsSet("effective_date") {
						body["effective_date"] = cmd.String("effective_date")
					}
					if cmd.IsSet("rate") {
						body["rate"] = cmd.Float("rate")
					}
					if err := internal.MergeJSON(cmd, body); err != nil {
						return err
					}
					raw, err := c.Request("PUT", c.BooksBase+"/settings/currencies/"+cmd.Args().First()+"/exchangerates/"+cmd.Args().Get(1), &zohttp.RequestOpts{
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
				Name:      "get-exchange-rate",
				Usage:     "Get an exchange rate",
				ArgsUsage: "<currency-id> <exchange-rate-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("GET", c.BooksBase+"/settings/currencies/"+cmd.Args().First()+"/exchangerates/"+cmd.Args().Get(1), &zohttp.RequestOpts{Params: orgParams(orgID)})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "delete-exchange-rate",
				Usage:     "Delete an exchange rate",
				ArgsUsage: "<currency-id> <exchange-rate-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("DELETE", c.BooksBase+"/settings/currencies/"+cmd.Args().First()+"/exchangerates/"+cmd.Args().Get(1), &zohttp.RequestOpts{Params: orgParams(orgID)})
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
				Name:  "create",
				Usage: "Create a tax",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "tax_name", Required: true, Usage: "Tax name"},
					&cli.FloatFlag{Name: "tax_percentage", Required: true, Usage: "Tax percentage"},
					&cli.StringFlag{Name: "json", Usage: "Additional fields as JSON"},
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
					body := map[string]any{}
					body["tax_name"] = cmd.String("tax_name")
					body["tax_percentage"] = cmd.Float("tax_percentage")
					if err := internal.MergeJSON(cmd, body); err != nil {
						return err
					}
					raw, err := c.Request("POST", c.BooksBase+"/settings/taxes", &zohttp.RequestOpts{
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
					raw, err := c.Request("GET", c.BooksBase+"/settings/taxes", &zohttp.RequestOpts{Params: orgParams(orgID)})
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
					&cli.StringFlag{Name: "json", Usage: "Additional fields as JSON"},
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
					body := map[string]any{}
					if cmd.IsSet("tax_name") {
						body["tax_name"] = cmd.String("tax_name")
					}
					if cmd.IsSet("tax_percentage") {
						body["tax_percentage"] = cmd.Float("tax_percentage")
					}
					if err := internal.MergeJSON(cmd, body); err != nil {
						return err
					}
					raw, err := c.Request("PUT", c.BooksBase+"/settings/taxes/"+cmd.Args().First(), &zohttp.RequestOpts{
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
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("GET", c.BooksBase+"/settings/taxes/"+cmd.Args().First(), &zohttp.RequestOpts{Params: orgParams(orgID)})
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
					raw, err := c.Request("DELETE", c.BooksBase+"/settings/taxes/"+cmd.Args().First(), &zohttp.RequestOpts{Params: orgParams(orgID)})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:  "create-group",
				Usage: "Create a tax group",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "tax_group_name", Required: true, Usage: "Tax group name"},
					&cli.StringFlag{Name: "taxes", Usage: "Tax IDs (comma-separated)"},
					&cli.StringFlag{Name: "json", Usage: "Additional fields as JSON"},
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
					body := map[string]any{}
					body["tax_group_name"] = cmd.String("tax_group_name")
					if cmd.IsSet("taxes") {
						var parsed any
						if err := json.Unmarshal([]byte(cmd.String("taxes")), &parsed); err != nil {
							return err
						}
						body["taxes"] = parsed
					}
					if err := internal.MergeJSON(cmd, body); err != nil {
						return err
					}
					raw, err := c.Request("POST", c.BooksBase+"/settings/taxgroups", &zohttp.RequestOpts{
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
				Name:  "list-groups",
				Usage: "List tax groups",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("GET", c.BooksBase+"/settings/taxgroups", &zohttp.RequestOpts{Params: orgParams(orgID)})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "update-group",
				Usage:     "Update a tax group",
				ArgsUsage: "<group-id>",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "tax_group_name", Usage: "Tax group name"},
					&cli.StringFlag{Name: "taxes", Usage: "Tax IDs (comma-separated)"},
					&cli.StringFlag{Name: "json", Usage: "Additional fields as JSON"},
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
					body := map[string]any{}
					if cmd.IsSet("tax_group_name") {
						body["tax_group_name"] = cmd.String("tax_group_name")
					}
					if cmd.IsSet("taxes") {
						var parsed any
						if err := json.Unmarshal([]byte(cmd.String("taxes")), &parsed); err != nil {
							return err
						}
						body["taxes"] = parsed
					}
					if err := internal.MergeJSON(cmd, body); err != nil {
						return err
					}
					raw, err := c.Request("PUT", c.BooksBase+"/settings/taxgroups/"+cmd.Args().First(), &zohttp.RequestOpts{
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
					raw, err := c.Request("GET", c.BooksBase+"/settings/taxgroups/"+cmd.Args().First(), &zohttp.RequestOpts{Params: orgParams(orgID)})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "delete-group",
				Usage:     "Delete a tax group",
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
					raw, err := c.Request("DELETE", c.BooksBase+"/settings/taxgroups/"+cmd.Args().First(), &zohttp.RequestOpts{Params: orgParams(orgID)})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:  "create-authority",
				Usage: "Create a tax authority",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "tax_authority_name", Required: true, Usage: "Tax authority name"},
					&cli.StringFlag{Name: "country", Usage: "Country"},
					&cli.StringFlag{Name: "state", Usage: "State"},
					&cli.StringFlag{Name: "jurisdiction", Usage: "Jurisdiction"},
					&cli.StringFlag{Name: "json", Usage: "Additional fields as JSON"},
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
					body := map[string]any{}
					body["tax_authority_name"] = cmd.String("tax_authority_name")
					if cmd.IsSet("country") {
						body["country"] = cmd.String("country")
					}
					if cmd.IsSet("state") {
						body["state"] = cmd.String("state")
					}
					if cmd.IsSet("jurisdiction") {
						body["jurisdiction"] = cmd.String("jurisdiction")
					}
					if err := internal.MergeJSON(cmd, body); err != nil {
						return err
					}
					raw, err := c.Request("POST", c.BooksBase+"/settings/taxauthorities", &zohttp.RequestOpts{
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
				Name:  "list-authorities",
				Usage: "List tax authorities",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("GET", c.BooksBase+"/settings/taxauthorities", &zohttp.RequestOpts{Params: orgParams(orgID)})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "update-authority",
				Usage:     "Update a tax authority",
				ArgsUsage: "<authority-id>",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "tax_authority_name", Usage: "Tax authority name"},
					&cli.StringFlag{Name: "country", Usage: "Country"},
					&cli.StringFlag{Name: "state", Usage: "State"},
					&cli.StringFlag{Name: "jurisdiction", Usage: "Jurisdiction"},
					&cli.StringFlag{Name: "json", Usage: "Additional fields as JSON"},
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
					body := map[string]any{}
					if cmd.IsSet("tax_authority_name") {
						body["tax_authority_name"] = cmd.String("tax_authority_name")
					}
					if cmd.IsSet("country") {
						body["country"] = cmd.String("country")
					}
					if cmd.IsSet("state") {
						body["state"] = cmd.String("state")
					}
					if cmd.IsSet("jurisdiction") {
						body["jurisdiction"] = cmd.String("jurisdiction")
					}
					if err := internal.MergeJSON(cmd, body); err != nil {
						return err
					}
					raw, err := c.Request("PUT", c.BooksBase+"/settings/taxauthorities/"+cmd.Args().First(), &zohttp.RequestOpts{
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
				Name:      "get-authority",
				Usage:     "Get a tax authority",
				ArgsUsage: "<authority-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("GET", c.BooksBase+"/settings/taxauthorities/"+cmd.Args().First(), &zohttp.RequestOpts{Params: orgParams(orgID)})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "delete-authority",
				Usage:     "Delete a tax authority",
				ArgsUsage: "<authority-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("DELETE", c.BooksBase+"/settings/taxauthorities/"+cmd.Args().First(), &zohttp.RequestOpts{Params: orgParams(orgID)})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:  "create-exemption",
				Usage: "Create a tax exemption",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "tax_exemption_name", Required: true, Usage: "Tax exemption name"},
					&cli.StringFlag{Name: "tax_exemption_code", Usage: "Tax exemption code"},
					&cli.StringFlag{Name: "description", Usage: "Description"},
					&cli.StringFlag{Name: "json", Usage: "Additional fields as JSON"},
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
					body := map[string]any{}
					body["tax_exemption_name"] = cmd.String("tax_exemption_name")
					if cmd.IsSet("tax_exemption_code") {
						body["tax_exemption_code"] = cmd.String("tax_exemption_code")
					}
					if cmd.IsSet("description") {
						body["description"] = cmd.String("description")
					}
					if err := internal.MergeJSON(cmd, body); err != nil {
						return err
					}
					raw, err := c.Request("POST", c.BooksBase+"/settings/taxexemptions", &zohttp.RequestOpts{
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
				Name:  "list-exemptions",
				Usage: "List tax exemptions",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("GET", c.BooksBase+"/settings/taxexemptions", &zohttp.RequestOpts{Params: orgParams(orgID)})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "update-exemption",
				Usage:     "Update a tax exemption",
				ArgsUsage: "<exemption-id>",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "tax_exemption_name", Usage: "Tax exemption name"},
					&cli.StringFlag{Name: "tax_exemption_code", Usage: "Tax exemption code"},
					&cli.StringFlag{Name: "description", Usage: "Description"},
					&cli.StringFlag{Name: "json", Usage: "Additional fields as JSON"},
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
					body := map[string]any{}
					if cmd.IsSet("tax_exemption_name") {
						body["tax_exemption_name"] = cmd.String("tax_exemption_name")
					}
					if cmd.IsSet("tax_exemption_code") {
						body["tax_exemption_code"] = cmd.String("tax_exemption_code")
					}
					if cmd.IsSet("description") {
						body["description"] = cmd.String("description")
					}
					if err := internal.MergeJSON(cmd, body); err != nil {
						return err
					}
					raw, err := c.Request("PUT", c.BooksBase+"/settings/taxexemptions/"+cmd.Args().First(), &zohttp.RequestOpts{
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
				Name:      "get-exemption",
				Usage:     "Get a tax exemption",
				ArgsUsage: "<exemption-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("GET", c.BooksBase+"/settings/taxexemptions/"+cmd.Args().First(), &zohttp.RequestOpts{Params: orgParams(orgID)})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "delete-exemption",
				Usage:     "Delete a tax exemption",
				ArgsUsage: "<exemption-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("DELETE", c.BooksBase+"/settings/taxexemptions/"+cmd.Args().First(), &zohttp.RequestOpts{Params: orgParams(orgID)})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
		},
	}
}

func openingBalanceCmd() *cli.Command {
	return &cli.Command{
		Name:  "opening-balance",
		Usage: "Opening balance operations",
		Commands: []*cli.Command{
			{
				Name:  "create",
				Usage: "Create opening balance",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "account_id", Required: true, Usage: "Account ID"},
					&cli.FloatFlag{Name: "opening_balance", Required: true, Usage: "Opening balance amount"},
					&cli.StringFlag{Name: "json", Usage: "Additional fields as JSON"},
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
					body := map[string]any{}
					body["account_id"] = cmd.String("account_id")
					body["opening_balance"] = cmd.Float("opening_balance")
					if err := internal.MergeJSON(cmd, body); err != nil {
						return err
					}
					raw, err := c.Request("POST", c.BooksBase+"/settings/openingbalances", &zohttp.RequestOpts{
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
				Name:  "update",
				Usage: "Update opening balance",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "account_id", Usage: "Account ID"},
					&cli.FloatFlag{Name: "opening_balance", Usage: "Opening balance amount"},
					&cli.StringFlag{Name: "json", Usage: "Additional fields as JSON"},
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
					body := map[string]any{}
					if cmd.IsSet("account_id") {
						body["account_id"] = cmd.String("account_id")
					}
					if cmd.IsSet("opening_balance") {
						body["opening_balance"] = cmd.Float("opening_balance")
					}
					if err := internal.MergeJSON(cmd, body); err != nil {
						return err
					}
					raw, err := c.Request("PUT", c.BooksBase+"/settings/openingbalances", &zohttp.RequestOpts{
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
				Name:  "get",
				Usage: "Get opening balance",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("GET", c.BooksBase+"/settings/openingbalances", &zohttp.RequestOpts{Params: orgParams(orgID)})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:  "delete",
				Usage: "Delete opening balance",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("DELETE", c.BooksBase+"/settings/openingbalances", &zohttp.RequestOpts{Params: orgParams(orgID)})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
		},
	}
}

func crmIntegrationCmd() *cli.Command {
	return &cli.Command{
		Name:  "crm-integration",
		Usage: "CRM integration operations",
		Commands: []*cli.Command{
			{
				Name:  "import-contact",
				Usage: "Import contact from CRM",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "contact_ids", Required: true, Usage: "CRM contact IDs (comma-separated)"},
					&cli.StringFlag{Name: "json", Usage: "Additional fields as JSON"},
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
					body := map[string]any{}
					var parsed0 any
					if err := json.Unmarshal([]byte(cmd.String("contact_ids")), &parsed0); err != nil {
						return err
					}
					body["contact_ids"] = parsed0
					if err := internal.MergeJSON(cmd, body); err != nil {
						return err
					}
					raw, err := c.Request("POST", c.BooksBase+"/crm"+"/contacts/importfromcrm", &zohttp.RequestOpts{
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
				Name:  "import-item",
				Usage: "Import item from CRM",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "item_ids", Required: true, Usage: "CRM item IDs (comma-separated)"},
					&cli.StringFlag{Name: "json", Usage: "Additional fields as JSON"},
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
					body := map[string]any{}
					var parsed0 any
					if err := json.Unmarshal([]byte(cmd.String("item_ids")), &parsed0); err != nil {
						return err
					}
					body["item_ids"] = parsed0
					if err := internal.MergeJSON(cmd, body); err != nil {
						return err
					}
					raw, err := c.Request("POST", c.BooksBase+"/crm"+"/items/importfromcrm", &zohttp.RequestOpts{
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

func reportingTagsCmd() *cli.Command {
	return &cli.Command{
		Name:  "reporting-tags",
		Usage: "Reporting tag operations",
		Commands: []*cli.Command{
			{
				Name:  "create",
				Usage: "Create a reporting tag",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "tag_name", Required: true, Usage: "Reporting tag name"},
					&cli.StringFlag{Name: "json", Usage: "Additional fields as JSON"},
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
					body := map[string]any{}
					body["tag_name"] = cmd.String("tag_name")
					if err := internal.MergeJSON(cmd, body); err != nil {
						return err
					}
					raw, err := c.Request("POST", c.BooksBase+"/settings/reportingtags", &zohttp.RequestOpts{
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
					raw, err := c.Request("GET", c.BooksBase+"/settings/reportingtags", &zohttp.RequestOpts{Params: orgParams(orgID)})
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
					&cli.StringFlag{Name: "tag_name", Usage: "Reporting tag name"},
					&cli.StringFlag{Name: "json", Usage: "Additional fields as JSON"},
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
					body := map[string]any{}
					if cmd.IsSet("tag_name") {
						body["tag_name"] = cmd.String("tag_name")
					}
					if err := internal.MergeJSON(cmd, body); err != nil {
						return err
					}
					raw, err := c.Request("PUT", c.BooksBase+"/settings/reportingtags/"+cmd.Args().First(), &zohttp.RequestOpts{
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
					raw, err := c.Request("DELETE", c.BooksBase+"/settings/reportingtags/"+cmd.Args().First(), &zohttp.RequestOpts{Params: orgParams(orgID)})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "mark-default-option",
				Usage:     "Mark default option for a reporting tag",
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
					raw, err := c.Request("PUT", c.BooksBase+"/settings/reportingtags/"+cmd.Args().First()+"/default/"+cmd.Args().Get(1), &zohttp.RequestOpts{Params: orgParams(orgID)})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "update-options",
				Usage:     "Update options of a reporting tag",
				ArgsUsage: "<tag-id>",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "tag_option_name", Required: true, Usage: "Tag option name"},
					&cli.StringFlag{Name: "json", Usage: "Additional fields as JSON"},
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
					body := map[string]any{}
					body["tag_option_name"] = cmd.String("tag_option_name")
					if err := internal.MergeJSON(cmd, body); err != nil {
						return err
					}
					raw, err := c.Request("PUT", c.BooksBase+"/settings/reportingtags/"+cmd.Args().First()+"/options", &zohttp.RequestOpts{
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
				Name:      "update-visibility",
				Usage:     "Update visibility of a reporting tag",
				ArgsUsage: "<tag-id>",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "visibility", Required: true, Usage: "Visibility"},
					&cli.StringFlag{Name: "json", Usage: "Additional fields as JSON"},
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
					body := map[string]any{}
					body["visibility"] = cmd.String("visibility")
					if err := internal.MergeJSON(cmd, body); err != nil {
						return err
					}
					raw, err := c.Request("PUT", c.BooksBase+"/settings/reportingtags/"+cmd.Args().First()+"/visibility", &zohttp.RequestOpts{
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
				Name:      "mark-active",
				Usage:     "Mark a reporting tag as active",
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
					raw, err := c.Request("POST", c.BooksBase+"/settings/reportingtags/"+cmd.Args().First()+"/active", &zohttp.RequestOpts{Params: orgParams(orgID)})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "mark-inactive",
				Usage:     "Mark a reporting tag as inactive",
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
					raw, err := c.Request("POST", c.BooksBase+"/settings/reportingtags/"+cmd.Args().First()+"/inactive", &zohttp.RequestOpts{Params: orgParams(orgID)})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "mark-option-active",
				Usage:     "Mark a reporting tag option as active",
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
					raw, err := c.Request("POST", c.BooksBase+"/settings/reportingtags/"+cmd.Args().First()+"/options/"+cmd.Args().Get(1)+"/active", &zohttp.RequestOpts{Params: orgParams(orgID)})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "mark-option-inactive",
				Usage:     "Mark a reporting tag option as inactive",
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
					raw, err := c.Request("POST", c.BooksBase+"/settings/reportingtags/"+cmd.Args().First()+"/options/"+cmd.Args().Get(1)+"/inactive", &zohttp.RequestOpts{Params: orgParams(orgID)})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "get-options-detail",
				Usage:     "Get options detail of a reporting tag",
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
					raw, err := c.Request("GET", c.BooksBase+"/settings/reportingtags/"+cmd.Args().First()+"/options", &zohttp.RequestOpts{Params: orgParams(orgID)})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:  "get-all-options",
				Usage: "Get all reporting tag options",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("GET", c.BooksBase+"/settings/reportingtags"+"/options", &zohttp.RequestOpts{Params: orgParams(orgID)})
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
					&cli.StringFlag{Name: "tag_order", Required: true, Usage: "Tag order JSON"},
					&cli.StringFlag{Name: "json", Usage: "Additional fields as JSON"},
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
					body := map[string]any{}
					body["tag_order"] = cmd.String("tag_order")
					if err := internal.MergeJSON(cmd, body); err != nil {
						return err
					}
					raw, err := c.Request("PUT", c.BooksBase+"/settings/reportingtags"+"/reorder", &zohttp.RequestOpts{
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
