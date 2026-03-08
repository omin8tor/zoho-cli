package books

import (
	"context"

	"github.com/omin8tor/zoho-cli/internal"
	zohttp "github.com/omin8tor/zoho-cli/internal/http"
	"github.com/omin8tor/zoho-cli/internal/output"
	"github.com/urfave/cli/v3"
)

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
					c, err := zohttp.GetClient()
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
					c, err := zohttp.GetClient()
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
					c, err := zohttp.GetClient()
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
