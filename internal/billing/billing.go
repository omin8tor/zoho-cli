package billing

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/omin8tor/zoho-cli/internal"
	zohttp "github.com/omin8tor/zoho-cli/internal/http"
	"github.com/omin8tor/zoho-cli/internal/output"
	"github.com/omin8tor/zoho-cli/internal/pagination"
	"github.com/urfave/cli/v3"
)

func Commands() *cli.Command {
	return &cli.Command{
		Name:  "billing",
		Usage: "Zoho Billing operations",
		Flags: []cli.Flag{
			&cli.StringFlag{Name: "org", Sources: cli.EnvVars("ZOHO_BOOKS_ORG_ID"), Usage: "Organization ID (or set ZOHO_BOOKS_ORG_ID)"},
		},
		Commands: []*cli.Command{
			productsCmd(),
			plansCmd(),
			addonsCmd(),
			couponsCmd(),
			customersCmd(),
			subscriptionsCmd(),
			invoicesCmd(),
			paymentsCmd(),
			creditNotesCmd(),
			hostedPagesCmd(),
			eventsCmd(),
			organizationsCmd(),
			currenciesCmd(),
			taxesCmd(),
			usersCmd(),
		},
	}
}

func productsCmd() *cli.Command {
	return &cli.Command{
		Name:  "products",
		Usage: "Product operations",
		Commands: []*cli.Command{
			{
				Name:  "list",
				Usage: "List all products",
				Flags: []cli.Flag{
					&cli.BoolFlag{Name: "all", Usage: "Fetch all records"},
					&cli.IntFlag{Name: "limit", Usage: "Max total records to fetch"},
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
					params := map[string]string{"organization_id": orgID}
					if cmd.Bool("all") || cmd.IsSet("limit") {
						items, err := pagination.Paginate(pagination.PaginationConfig{
							Client:   c,
							URL:      c.BillingBase+"/products",
							Opts:     &zohttp.RequestOpts{Params: params},
							ItemsKey: "products",
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
					raw, err := c.Request("GET", c.BillingBase+"/products", &zohttp.RequestOpts{Params: params})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "get",
				Usage:     "Retrieve a product",
				ArgsUsage: "<product-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					if cmd.Args().Len() < 1 {
						return internal.NewValidationError("product ID required")
					}
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := internal.RequireFlag(cmd, "org", "ZOHO_BOOKS_ORG_ID")
					if err != nil {
						return err
					}
					raw, err := c.Request("GET", c.BillingBase+"/products/"+cmd.Args().First(), &zohttp.RequestOpts{Params: map[string]string{"organization_id": orgID}})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:  "create",
				Usage: "Create a product",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "name", Required: true, Usage: "Product name"},
					&cli.StringFlag{Name: "description", Usage: "Product description"},
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
					if cmd.IsSet("description") {
						body["description"] = cmd.String("description")
					}
					if err := internal.MergeJSON(cmd, body); err != nil {
						return err
					}
					raw, err := c.Request("POST", c.BillingBase+"/products", &zohttp.RequestOpts{
						Params: map[string]string{"organization_id": orgID},
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
				Usage:     "Update a product",
				ArgsUsage: "<product-id>",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "name", Usage: "Product name"},
					&cli.StringFlag{Name: "description", Usage: "Product description"},
					&cli.StringFlag{Name: "json", Usage: "Additional fields as JSON"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					if cmd.Args().Len() < 1 {
						return internal.NewValidationError("product ID required")
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
					if cmd.IsSet("description") {
						body["description"] = cmd.String("description")
					}
					if err := internal.MergeJSON(cmd, body); err != nil {
						return err
					}
					raw, err := c.Request("PUT", c.BillingBase+"/products/"+cmd.Args().First(), &zohttp.RequestOpts{
						Params: map[string]string{"organization_id": orgID},
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
				Usage:     "Delete a product",
				ArgsUsage: "<product-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					if cmd.Args().Len() < 1 {
						return internal.NewValidationError("product ID required")
					}
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := internal.RequireFlag(cmd, "org", "ZOHO_BOOKS_ORG_ID")
					if err != nil {
						return err
					}
					raw, err := c.Request("DELETE", c.BillingBase+"/products/"+cmd.Args().First(), &zohttp.RequestOpts{Params: map[string]string{"organization_id": orgID}})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "mark-active",
				Usage:     "Mark a product as active",
				ArgsUsage: "<product-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					if cmd.Args().Len() < 1 {
						return internal.NewValidationError("product ID required")
					}
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := internal.RequireFlag(cmd, "org", "ZOHO_BOOKS_ORG_ID")
					if err != nil {
						return err
					}
					raw, err := c.Request("POST", c.BillingBase+"/products/"+cmd.Args().First()+"/markasactive", &zohttp.RequestOpts{Params: map[string]string{"organization_id": orgID}})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "mark-inactive",
				Usage:     "Mark a product as inactive",
				ArgsUsage: "<product-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					if cmd.Args().Len() < 1 {
						return internal.NewValidationError("product ID required")
					}
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := internal.RequireFlag(cmd, "org", "ZOHO_BOOKS_ORG_ID")
					if err != nil {
						return err
					}
					raw, err := c.Request("POST", c.BillingBase+"/products/"+cmd.Args().First()+"/markasinactive", &zohttp.RequestOpts{Params: map[string]string{"organization_id": orgID}})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
		},
	}
}

func plansCmd() *cli.Command {
	return &cli.Command{
		Name:  "plans",
		Usage: "Plan operations",
		Commands: []*cli.Command{
			{
				Name:  "list",
				Usage: "List all plans",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "product-id", Usage: "Filter by product"},
					&cli.BoolFlag{Name: "all", Usage: "Fetch all records"},
					&cli.IntFlag{Name: "limit", Usage: "Max total records to fetch"},
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
					params := map[string]string{"organization_id": orgID}
					if v := cmd.String("product-id"); v != "" {
						params["product_id"] = v
					}
					if cmd.Bool("all") || cmd.IsSet("limit") {
						items, err := pagination.Paginate(pagination.PaginationConfig{
							Client:   c,
							URL:      c.BillingBase+"/plans",
							Opts:     &zohttp.RequestOpts{Params: params},
							ItemsKey: "plans",
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
					raw, err := c.Request("GET", c.BillingBase+"/plans", &zohttp.RequestOpts{Params: params})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "get",
				Usage:     "Retrieve a plan",
				ArgsUsage: "<plan-code>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					if cmd.Args().Len() < 1 {
						return internal.NewValidationError("plan code required")
					}
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := internal.RequireFlag(cmd, "org", "ZOHO_BOOKS_ORG_ID")
					if err != nil {
						return err
					}
					raw, err := c.Request("GET", c.BillingBase+"/plans/"+cmd.Args().First(), &zohttp.RequestOpts{Params: map[string]string{"organization_id": orgID}})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:  "create",
				Usage: "Create a plan",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "plan_code", Required: true, Usage: "Plan code"},
					&cli.StringFlag{Name: "name", Required: true, Usage: "Plan name"},
					&cli.StringFlag{Name: "product_id", Required: true, Usage: "Product ID"},
					&cli.FloatFlag{Name: "recurring_price", Required: true, Usage: "Recurring price"},
					&cli.IntFlag{Name: "interval", Required: true, Usage: "Billing interval"},
					&cli.StringFlag{Name: "interval_unit", Required: true, Usage: "Billing interval unit"},
					&cli.StringFlag{Name: "description", Usage: "Plan description"},
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
					body["plan_code"] = cmd.String("plan_code")
					body["name"] = cmd.String("name")
					body["product_id"] = cmd.String("product_id")
					body["recurring_price"] = cmd.Float("recurring_price")
					body["interval"] = cmd.Int("interval")
					body["interval_unit"] = cmd.String("interval_unit")
					if cmd.IsSet("description") {
						body["description"] = cmd.String("description")
					}
					if err := internal.MergeJSON(cmd, body); err != nil {
						return err
					}
					raw, err := c.Request("POST", c.BillingBase+"/plans", &zohttp.RequestOpts{
						Params: map[string]string{"organization_id": orgID},
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
				Usage:     "Update a plan",
				ArgsUsage: "<plan-code>",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "name", Usage: "Plan name"},
					&cli.StringFlag{Name: "product_id", Usage: "Product ID"},
					&cli.FloatFlag{Name: "recurring_price", Usage: "Recurring price"},
					&cli.IntFlag{Name: "interval", Usage: "Billing interval"},
					&cli.StringFlag{Name: "interval_unit", Usage: "Billing interval unit"},
					&cli.StringFlag{Name: "description", Usage: "Plan description"},
					&cli.StringFlag{Name: "json", Usage: "Additional fields as JSON"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					if cmd.Args().Len() < 1 {
						return internal.NewValidationError("plan code required")
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
					if cmd.IsSet("product_id") {
						body["product_id"] = cmd.String("product_id")
					}
					if cmd.IsSet("recurring_price") {
						body["recurring_price"] = cmd.Float("recurring_price")
					}
					if cmd.IsSet("interval") {
						body["interval"] = cmd.Int("interval")
					}
					if cmd.IsSet("interval_unit") {
						body["interval_unit"] = cmd.String("interval_unit")
					}
					if cmd.IsSet("description") {
						body["description"] = cmd.String("description")
					}
					if err := internal.MergeJSON(cmd, body); err != nil {
						return err
					}
					raw, err := c.Request("PUT", c.BillingBase+"/plans/"+cmd.Args().First(), &zohttp.RequestOpts{
						Params: map[string]string{"organization_id": orgID},
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
				Usage:     "Delete a plan",
				ArgsUsage: "<plan-code>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					if cmd.Args().Len() < 1 {
						return internal.NewValidationError("plan code required")
					}
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := internal.RequireFlag(cmd, "org", "ZOHO_BOOKS_ORG_ID")
					if err != nil {
						return err
					}
					raw, err := c.Request("DELETE", c.BillingBase+"/plans/"+cmd.Args().First(), &zohttp.RequestOpts{Params: map[string]string{"organization_id": orgID}})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "mark-active",
				Usage:     "Mark a plan as active",
				ArgsUsage: "<plan-code>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					if cmd.Args().Len() < 1 {
						return internal.NewValidationError("plan code required")
					}
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := internal.RequireFlag(cmd, "org", "ZOHO_BOOKS_ORG_ID")
					if err != nil {
						return err
					}
					raw, err := c.Request("POST", c.BillingBase+"/plans/"+cmd.Args().First()+"/markasactive", &zohttp.RequestOpts{Params: map[string]string{"organization_id": orgID}})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "mark-inactive",
				Usage:     "Mark a plan as inactive",
				ArgsUsage: "<plan-code>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					if cmd.Args().Len() < 1 {
						return internal.NewValidationError("plan code required")
					}
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := internal.RequireFlag(cmd, "org", "ZOHO_BOOKS_ORG_ID")
					if err != nil {
						return err
					}
					raw, err := c.Request("POST", c.BillingBase+"/plans/"+cmd.Args().First()+"/markasinactive", &zohttp.RequestOpts{Params: map[string]string{"organization_id": orgID}})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
		},
	}
}

func addonsCmd() *cli.Command {
	return &cli.Command{
		Name:  "addons",
		Usage: "Addon operations",
		Commands: []*cli.Command{
			{
				Name:  "list",
				Usage: "List all addons",
				Flags: []cli.Flag{
					&cli.BoolFlag{Name: "all", Usage: "Fetch all records"},
					&cli.IntFlag{Name: "limit", Usage: "Max total records to fetch"},
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
					params := map[string]string{"organization_id": orgID}
					if cmd.Bool("all") || cmd.IsSet("limit") {
						items, err := pagination.Paginate(pagination.PaginationConfig{
							Client:   c,
							URL:      c.BillingBase+"/addons",
							Opts:     &zohttp.RequestOpts{Params: params},
							ItemsKey: "addons",
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
					raw, err := c.Request("GET", c.BillingBase+"/addons", &zohttp.RequestOpts{Params: params})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "get",
				Usage:     "Retrieve an addon",
				ArgsUsage: "<addon-code>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					if cmd.Args().Len() < 1 {
						return internal.NewValidationError("addon code required")
					}
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := internal.RequireFlag(cmd, "org", "ZOHO_BOOKS_ORG_ID")
					if err != nil {
						return err
					}
					raw, err := c.Request("GET", c.BillingBase+"/addons/"+cmd.Args().First(), &zohttp.RequestOpts{Params: map[string]string{"organization_id": orgID}})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:  "create",
				Usage: "Create an addon",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "addon_code", Required: true, Usage: "Addon code"},
					&cli.StringFlag{Name: "name", Required: true, Usage: "Addon name"},
					&cli.StringFlag{Name: "type", Required: true, Usage: "Addon type"},
					&cli.FloatFlag{Name: "price", Required: true, Usage: "Addon price"},
					&cli.StringFlag{Name: "pricing_scheme", Usage: "Pricing scheme"},
					&cli.StringFlag{Name: "interval_unit", Usage: "Interval unit"},
					&cli.StringFlag{Name: "unit_name", Usage: "Unit name"},
					&cli.StringFlag{Name: "product_id", Usage: "Product ID"},
					&cli.StringFlag{Name: "description", Usage: "Addon description"},
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
					body["addon_code"] = cmd.String("addon_code")
					body["name"] = cmd.String("name")
					body["type"] = cmd.String("type")
					body["price"] = cmd.Float("price")
					if cmd.IsSet("pricing_scheme") {
						body["pricing_scheme"] = cmd.String("pricing_scheme")
					}
					if cmd.IsSet("interval_unit") {
						body["interval_unit"] = cmd.String("interval_unit")
					}
					if cmd.IsSet("unit_name") {
						body["unit_name"] = cmd.String("unit_name")
					}
					if cmd.IsSet("product_id") {
						body["product_id"] = cmd.String("product_id")
					}
					if cmd.IsSet("description") {
						body["description"] = cmd.String("description")
					}
					if err := internal.MergeJSON(cmd, body); err != nil {
						return err
					}
					raw, err := c.Request("POST", c.BillingBase+"/addons", &zohttp.RequestOpts{
						Params: map[string]string{"organization_id": orgID},
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
				Usage:     "Update an addon",
				ArgsUsage: "<addon-code>",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "name", Usage: "Addon name"},
					&cli.StringFlag{Name: "type", Usage: "Addon type"},
					&cli.FloatFlag{Name: "price", Usage: "Addon price"},
					&cli.StringFlag{Name: "pricing_scheme", Usage: "Pricing scheme"},
					&cli.StringFlag{Name: "interval_unit", Usage: "Interval unit"},
					&cli.StringFlag{Name: "unit_name", Usage: "Unit name"},
					&cli.StringFlag{Name: "product_id", Usage: "Product ID"},
					&cli.StringFlag{Name: "description", Usage: "Addon description"},
					&cli.StringFlag{Name: "json", Usage: "Additional fields as JSON"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					if cmd.Args().Len() < 1 {
						return internal.NewValidationError("addon code required")
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
					if cmd.IsSet("type") {
						body["type"] = cmd.String("type")
					}
					if cmd.IsSet("price") {
						body["price"] = cmd.Float("price")
					}
					if cmd.IsSet("pricing_scheme") {
						body["pricing_scheme"] = cmd.String("pricing_scheme")
					}
					if cmd.IsSet("interval_unit") {
						body["interval_unit"] = cmd.String("interval_unit")
					}
					if cmd.IsSet("unit_name") {
						body["unit_name"] = cmd.String("unit_name")
					}
					if cmd.IsSet("product_id") {
						body["product_id"] = cmd.String("product_id")
					}
					if cmd.IsSet("description") {
						body["description"] = cmd.String("description")
					}
					if err := internal.MergeJSON(cmd, body); err != nil {
						return err
					}
					raw, err := c.Request("PUT", c.BillingBase+"/addons/"+cmd.Args().First(), &zohttp.RequestOpts{
						Params: map[string]string{"organization_id": orgID},
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
				Usage:     "Delete an addon",
				ArgsUsage: "<addon-code>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					if cmd.Args().Len() < 1 {
						return internal.NewValidationError("addon code required")
					}
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := internal.RequireFlag(cmd, "org", "ZOHO_BOOKS_ORG_ID")
					if err != nil {
						return err
					}
					raw, err := c.Request("DELETE", c.BillingBase+"/addons/"+cmd.Args().First(), &zohttp.RequestOpts{Params: map[string]string{"organization_id": orgID}})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "mark-active",
				Usage:     "Mark an addon as active",
				ArgsUsage: "<addon-code>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					if cmd.Args().Len() < 1 {
						return internal.NewValidationError("addon code required")
					}
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := internal.RequireFlag(cmd, "org", "ZOHO_BOOKS_ORG_ID")
					if err != nil {
						return err
					}
					raw, err := c.Request("POST", c.BillingBase+"/addons/"+cmd.Args().First()+"/markasactive", &zohttp.RequestOpts{Params: map[string]string{"organization_id": orgID}})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "mark-inactive",
				Usage:     "Mark an addon as inactive",
				ArgsUsage: "<addon-code>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					if cmd.Args().Len() < 1 {
						return internal.NewValidationError("addon code required")
					}
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := internal.RequireFlag(cmd, "org", "ZOHO_BOOKS_ORG_ID")
					if err != nil {
						return err
					}
					raw, err := c.Request("POST", c.BillingBase+"/addons/"+cmd.Args().First()+"/markasinactive", &zohttp.RequestOpts{Params: map[string]string{"organization_id": orgID}})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
		},
	}
}

func couponsCmd() *cli.Command {
	return &cli.Command{
		Name:  "coupons",
		Usage: "Coupon operations",
		Commands: []*cli.Command{
			{
				Name:  "list",
				Usage: "List all coupons",
				Flags: []cli.Flag{
					&cli.BoolFlag{Name: "all", Usage: "Fetch all records"},
					&cli.IntFlag{Name: "limit", Usage: "Max total records to fetch"},
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
					params := map[string]string{"organization_id": orgID}
					if cmd.Bool("all") || cmd.IsSet("limit") {
						items, err := pagination.Paginate(pagination.PaginationConfig{
							Client:   c,
							URL:      c.BillingBase+"/coupons",
							Opts:     &zohttp.RequestOpts{Params: params},
							ItemsKey: "coupons",
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
					raw, err := c.Request("GET", c.BillingBase+"/coupons", &zohttp.RequestOpts{Params: params})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "get",
				Usage:     "Retrieve a coupon",
				ArgsUsage: "<coupon-code>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					if cmd.Args().Len() < 1 {
						return internal.NewValidationError("coupon code required")
					}
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := internal.RequireFlag(cmd, "org", "ZOHO_BOOKS_ORG_ID")
					if err != nil {
						return err
					}
					raw, err := c.Request("GET", c.BillingBase+"/coupons/"+cmd.Args().First(), &zohttp.RequestOpts{Params: map[string]string{"organization_id": orgID}})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:  "create",
				Usage: "Create a coupon",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "coupon_code", Required: true, Usage: "Coupon code"},
					&cli.StringFlag{Name: "name", Required: true, Usage: "Coupon name"},
					&cli.StringFlag{Name: "discount_type", Usage: "Discount type"},
					&cli.FloatFlag{Name: "discount", Usage: "Discount value"},
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
					body["coupon_code"] = cmd.String("coupon_code")
					body["name"] = cmd.String("name")
					if cmd.IsSet("discount_type") {
						body["discount_type"] = cmd.String("discount_type")
					}
					if cmd.IsSet("discount") {
						body["discount"] = cmd.Float("discount")
					}
					if err := internal.MergeJSON(cmd, body); err != nil {
						return err
					}
					raw, err := c.Request("POST", c.BillingBase+"/coupons", &zohttp.RequestOpts{
						Params: map[string]string{"organization_id": orgID},
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
				Usage:     "Update a coupon",
				ArgsUsage: "<coupon-code>",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "name", Usage: "Coupon name"},
					&cli.StringFlag{Name: "discount_type", Usage: "Discount type"},
					&cli.FloatFlag{Name: "discount", Usage: "Discount value"},
					&cli.StringFlag{Name: "json", Usage: "Additional fields as JSON"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					if cmd.Args().Len() < 1 {
						return internal.NewValidationError("coupon code required")
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
					if cmd.IsSet("discount_type") {
						body["discount_type"] = cmd.String("discount_type")
					}
					if cmd.IsSet("discount") {
						body["discount"] = cmd.Float("discount")
					}
					if err := internal.MergeJSON(cmd, body); err != nil {
						return err
					}
					raw, err := c.Request("PUT", c.BillingBase+"/coupons/"+cmd.Args().First(), &zohttp.RequestOpts{
						Params: map[string]string{"organization_id": orgID},
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
				Usage:     "Delete a coupon",
				ArgsUsage: "<coupon-code>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					if cmd.Args().Len() < 1 {
						return internal.NewValidationError("coupon code required")
					}
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := internal.RequireFlag(cmd, "org", "ZOHO_BOOKS_ORG_ID")
					if err != nil {
						return err
					}
					raw, err := c.Request("DELETE", c.BillingBase+"/coupons/"+cmd.Args().First(), &zohttp.RequestOpts{Params: map[string]string{"organization_id": orgID}})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "mark-active",
				Usage:     "Mark a coupon as active",
				ArgsUsage: "<coupon-code>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					if cmd.Args().Len() < 1 {
						return internal.NewValidationError("coupon code required")
					}
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := internal.RequireFlag(cmd, "org", "ZOHO_BOOKS_ORG_ID")
					if err != nil {
						return err
					}
					raw, err := c.Request("POST", c.BillingBase+"/coupons/"+cmd.Args().First()+"/markasactive", &zohttp.RequestOpts{Params: map[string]string{"organization_id": orgID}})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "mark-inactive",
				Usage:     "Mark a coupon as inactive",
				ArgsUsage: "<coupon-code>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					if cmd.Args().Len() < 1 {
						return internal.NewValidationError("coupon code required")
					}
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := internal.RequireFlag(cmd, "org", "ZOHO_BOOKS_ORG_ID")
					if err != nil {
						return err
					}
					raw, err := c.Request("POST", c.BillingBase+"/coupons/"+cmd.Args().First()+"/markasinactive", &zohttp.RequestOpts{Params: map[string]string{"organization_id": orgID}})
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
		Usage: "Customer operations",
		Commands: []*cli.Command{
			{
				Name:  "list",
				Usage: "List all customers",
				Flags: []cli.Flag{
					&cli.BoolFlag{Name: "all", Usage: "Fetch all records"},
					&cli.IntFlag{Name: "limit", Usage: "Max total records to fetch"},
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
					params := map[string]string{"organization_id": orgID}
					if cmd.Bool("all") || cmd.IsSet("limit") {
						items, err := pagination.Paginate(pagination.PaginationConfig{
							Client:   c,
							URL:      c.BillingBase+"/customers",
							Opts:     &zohttp.RequestOpts{Params: params},
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
					raw, err := c.Request("GET", c.BillingBase+"/customers", &zohttp.RequestOpts{Params: params})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "get",
				Usage:     "Retrieve a customer",
				ArgsUsage: "<customer-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					if cmd.Args().Len() < 1 {
						return internal.NewValidationError("customer ID required")
					}
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := internal.RequireFlag(cmd, "org", "ZOHO_BOOKS_ORG_ID")
					if err != nil {
						return err
					}
					raw, err := c.Request("GET", c.BillingBase+"/customers/"+cmd.Args().First(), &zohttp.RequestOpts{Params: map[string]string{"organization_id": orgID}})
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
					&cli.StringFlag{Name: "display_name", Required: true, Usage: "Display name"},
					&cli.StringFlag{Name: "email", Usage: "Customer email"},
					&cli.StringFlag{Name: "company_name", Usage: "Company name"},
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
					body["display_name"] = cmd.String("display_name")
					if cmd.IsSet("email") {
						body["email"] = cmd.String("email")
					}
					if cmd.IsSet("company_name") {
						body["company_name"] = cmd.String("company_name")
					}
					if err := internal.MergeJSON(cmd, body); err != nil {
						return err
					}
					raw, err := c.Request("POST", c.BillingBase+"/customers", &zohttp.RequestOpts{
						Params: map[string]string{"organization_id": orgID},
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
				Usage:     "Update a customer",
				ArgsUsage: "<customer-id>",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "display_name", Usage: "Display name"},
					&cli.StringFlag{Name: "email", Usage: "Customer email"},
					&cli.StringFlag{Name: "company_name", Usage: "Company name"},
					&cli.StringFlag{Name: "json", Usage: "Additional fields as JSON"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					if cmd.Args().Len() < 1 {
						return internal.NewValidationError("customer ID required")
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
					if cmd.IsSet("display_name") {
						body["display_name"] = cmd.String("display_name")
					}
					if cmd.IsSet("email") {
						body["email"] = cmd.String("email")
					}
					if cmd.IsSet("company_name") {
						body["company_name"] = cmd.String("company_name")
					}
					if err := internal.MergeJSON(cmd, body); err != nil {
						return err
					}
					raw, err := c.Request("PUT", c.BillingBase+"/customers/"+cmd.Args().First(), &zohttp.RequestOpts{
						Params: map[string]string{"organization_id": orgID},
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
				Usage:     "Delete a customer",
				ArgsUsage: "<customer-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					if cmd.Args().Len() < 1 {
						return internal.NewValidationError("customer ID required")
					}
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := internal.RequireFlag(cmd, "org", "ZOHO_BOOKS_ORG_ID")
					if err != nil {
						return err
					}
					raw, err := c.Request("DELETE", c.BillingBase+"/customers/"+cmd.Args().First(), &zohttp.RequestOpts{Params: map[string]string{"organization_id": orgID}})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "mark-active",
				Usage:     "Mark a customer as active",
				ArgsUsage: "<customer-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					if cmd.Args().Len() < 1 {
						return internal.NewValidationError("customer ID required")
					}
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := internal.RequireFlag(cmd, "org", "ZOHO_BOOKS_ORG_ID")
					if err != nil {
						return err
					}
					raw, err := c.Request("POST", c.BillingBase+"/customers/"+cmd.Args().First()+"/markasactive", &zohttp.RequestOpts{Params: map[string]string{"organization_id": orgID}})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "mark-inactive",
				Usage:     "Mark a customer as inactive",
				ArgsUsage: "<customer-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					if cmd.Args().Len() < 1 {
						return internal.NewValidationError("customer ID required")
					}
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := internal.RequireFlag(cmd, "org", "ZOHO_BOOKS_ORG_ID")
					if err != nil {
						return err
					}
					raw, err := c.Request("POST", c.BillingBase+"/customers/"+cmd.Args().First()+"/markasinactive", &zohttp.RequestOpts{Params: map[string]string{"organization_id": orgID}})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
		},
	}
}

func subscriptionsCmd() *cli.Command {
	return &cli.Command{
		Name:  "subscriptions",
		Usage: "Subscription operations",
		Commands: []*cli.Command{
			{
				Name:  "list",
				Usage: "List all subscriptions",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "customer-id", Usage: "Filter by customer"},
					&cli.BoolFlag{Name: "all", Usage: "Fetch all records"},
					&cli.IntFlag{Name: "limit", Usage: "Max total records to fetch"},
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
					params := map[string]string{"organization_id": orgID}
					if v := cmd.String("customer-id"); v != "" {
						params["customer_id"] = v
					}
					if cmd.Bool("all") || cmd.IsSet("limit") {
						items, err := pagination.Paginate(pagination.PaginationConfig{
							Client:   c,
							URL:      c.BillingBase+"/subscriptions",
							Opts:     &zohttp.RequestOpts{Params: params},
							ItemsKey: "subscriptions",
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
					raw, err := c.Request("GET", c.BillingBase+"/subscriptions", &zohttp.RequestOpts{Params: params})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "get",
				Usage:     "Retrieve a subscription",
				ArgsUsage: "<subscription-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					if cmd.Args().Len() < 1 {
						return internal.NewValidationError("subscription ID required")
					}
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := internal.RequireFlag(cmd, "org", "ZOHO_BOOKS_ORG_ID")
					if err != nil {
						return err
					}
					raw, err := c.Request("GET", c.BillingBase+"/subscriptions/"+cmd.Args().First(), &zohttp.RequestOpts{Params: map[string]string{"organization_id": orgID}})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:  "create",
				Usage: "Create a subscription",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "customer_id", Required: true, Usage: "Customer ID"},
					&cli.StringFlag{Name: "plan_code", Required: true, Usage: "Plan code"},
					&cli.IntFlag{Name: "plan_quantity", Usage: "Plan quantity"},
					&cli.StringFlag{Name: "coupon_code", Usage: "Coupon code"},
					&cli.StringFlag{Name: "starts_at", Usage: "Subscription start date"},
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
					if cmd.IsSet("coupon_code") {
						body["coupon_code"] = cmd.String("coupon_code")
					}
					if cmd.IsSet("starts_at") {
						body["starts_at"] = cmd.String("starts_at")
					}
					plan := map[string]any{}
					plan["plan_code"] = cmd.String("plan_code")
					if cmd.IsSet("plan_quantity") {
						plan["plan_quantity"] = cmd.Int("plan_quantity")
					}
					body["plan"] = plan
					if err := internal.MergeJSON(cmd, body); err != nil {
						return err
					}
					raw, err := c.Request("POST", c.BillingBase+"/subscriptions", &zohttp.RequestOpts{
						Params: map[string]string{"organization_id": orgID},
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
				Usage:     "Update a subscription",
				ArgsUsage: "<subscription-id>",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "plan_code", Usage: "Plan code"},
					&cli.IntFlag{Name: "plan_quantity", Usage: "Plan quantity"},
					&cli.StringFlag{Name: "coupon_code", Usage: "Coupon code"},
					&cli.StringFlag{Name: "starts_at", Usage: "Subscription start date"},
					&cli.StringFlag{Name: "json", Usage: "Additional fields as JSON"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					if cmd.Args().Len() < 1 {
						return internal.NewValidationError("subscription ID required")
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
					if cmd.IsSet("coupon_code") {
						body["coupon_code"] = cmd.String("coupon_code")
					}
					if cmd.IsSet("starts_at") {
						body["starts_at"] = cmd.String("starts_at")
					}
					plan := map[string]any{}
					if cmd.IsSet("plan_code") {
						plan["plan_code"] = cmd.String("plan_code")
					}
					if cmd.IsSet("plan_quantity") {
						plan["plan_quantity"] = cmd.Int("plan_quantity")
					}
					if len(plan) > 0 {
						body["plan"] = plan
					}
					if err := internal.MergeJSON(cmd, body); err != nil {
						return err
					}
					raw, err := c.Request("PUT", c.BillingBase+"/subscriptions/"+cmd.Args().First(), &zohttp.RequestOpts{
						Params: map[string]string{"organization_id": orgID},
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
				Usage:     "Delete a subscription",
				ArgsUsage: "<subscription-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					if cmd.Args().Len() < 1 {
						return internal.NewValidationError("subscription ID required")
					}
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := internal.RequireFlag(cmd, "org", "ZOHO_BOOKS_ORG_ID")
					if err != nil {
						return err
					}
					raw, err := c.Request("DELETE", c.BillingBase+"/subscriptions/"+cmd.Args().First(), &zohttp.RequestOpts{Params: map[string]string{"organization_id": orgID}})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "cancel",
				Usage:     "Cancel a subscription",
				ArgsUsage: "<subscription-id>",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "json", Usage: "JSON body with cancel options"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					if cmd.Args().Len() < 1 {
						return internal.NewValidationError("subscription ID required")
					}
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := internal.RequireFlag(cmd, "org", "ZOHO_BOOKS_ORG_ID")
					if err != nil {
						return err
					}
					opts := &zohttp.RequestOpts{Params: map[string]string{"organization_id": orgID}}
					if j := cmd.String("json"); j != "" {
						var body any
						if err := json.Unmarshal([]byte(j), &body); err != nil {
							return internal.NewValidationError(fmt.Sprintf("--json: invalid JSON: %v", err))
						}
						opts.JSON = body
					}
					raw, err := c.Request("POST", c.BillingBase+"/subscriptions/"+cmd.Args().First()+"/cancel", opts)
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "reactivate",
				Usage:     "Reactivate a subscription",
				ArgsUsage: "<subscription-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					if cmd.Args().Len() < 1 {
						return internal.NewValidationError("subscription ID required")
					}
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := internal.RequireFlag(cmd, "org", "ZOHO_BOOKS_ORG_ID")
					if err != nil {
						return err
					}
					raw, err := c.Request("POST", c.BillingBase+"/subscriptions/"+cmd.Args().First()+"/reactivate", &zohttp.RequestOpts{Params: map[string]string{"organization_id": orgID}})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "scheduled-changes",
				Usage:     "View scheduled changes for a subscription",
				ArgsUsage: "<subscription-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					if cmd.Args().Len() < 1 {
						return internal.NewValidationError("subscription ID required")
					}
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := internal.RequireFlag(cmd, "org", "ZOHO_BOOKS_ORG_ID")
					if err != nil {
						return err
					}
					raw, err := c.Request("GET", c.BillingBase+"/subscriptions/"+cmd.Args().First()+"/scheduledchanges", &zohttp.RequestOpts{Params: map[string]string{"organization_id": orgID}})
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
				Usage: "List all invoices",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "subscription-id", Usage: "Filter by subscription"},
					&cli.StringFlag{Name: "customer-id", Usage: "Filter by customer"},
					&cli.BoolFlag{Name: "all", Usage: "Fetch all records"},
					&cli.IntFlag{Name: "limit", Usage: "Max total records to fetch"},
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
					params := map[string]string{"organization_id": orgID}
					if v := cmd.String("subscription-id"); v != "" {
						params["subscription_id"] = v
					}
					if v := cmd.String("customer-id"); v != "" {
						params["customer_id"] = v
					}
					if cmd.Bool("all") || cmd.IsSet("limit") {
						items, err := pagination.Paginate(pagination.PaginationConfig{
							Client:   c,
							URL:      c.BillingBase+"/invoices",
							Opts:     &zohttp.RequestOpts{Params: params},
							ItemsKey: "invoices",
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
					raw, err := c.Request("GET", c.BillingBase+"/invoices", &zohttp.RequestOpts{Params: params})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "get",
				Usage:     "Retrieve an invoice",
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
					raw, err := c.Request("GET", c.BillingBase+"/invoices/"+cmd.Args().First(), &zohttp.RequestOpts{Params: map[string]string{"organization_id": orgID}})
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
					&cli.StringFlag{Name: "date", Usage: "Invoice date"},
					&cli.StringFlag{Name: "due_date", Usage: "Due date"},
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
					if err := internal.MergeJSON(cmd, body); err != nil {
						return err
					}
					raw, err := c.Request("POST", c.BillingBase+"/invoices", &zohttp.RequestOpts{
						Params: map[string]string{"organization_id": orgID},
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
				Usage:     "Update an invoice",
				ArgsUsage: "<invoice-id>",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "customer_id", Usage: "Customer ID"},
					&cli.StringFlag{Name: "invoice_number", Usage: "Invoice number"},
					&cli.StringFlag{Name: "reference_number", Usage: "Reference number"},
					&cli.StringFlag{Name: "date", Usage: "Invoice date"},
					&cli.StringFlag{Name: "due_date", Usage: "Due date"},
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
					if err := internal.MergeJSON(cmd, body); err != nil {
						return err
					}
					raw, err := c.Request("PUT", c.BillingBase+"/invoices/"+cmd.Args().First(), &zohttp.RequestOpts{
						Params: map[string]string{"organization_id": orgID},
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
					raw, err := c.Request("DELETE", c.BillingBase+"/invoices/"+cmd.Args().First(), &zohttp.RequestOpts{Params: map[string]string{"organization_id": orgID}})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "convert-to-open",
				Usage:     "Convert an invoice to open",
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
					raw, err := c.Request("POST", c.BillingBase+"/invoices/"+cmd.Args().First()+"/converttoopen", &zohttp.RequestOpts{Params: map[string]string{"organization_id": orgID}})
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
					raw, err := c.Request("POST", c.BillingBase+"/invoices/"+cmd.Args().First()+"/void", &zohttp.RequestOpts{Params: map[string]string{"organization_id": orgID}})
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
					&cli.StringFlag{Name: "to_mail_ids", Usage: "Comma-separated recipient emails"},
					&cli.StringFlag{Name: "subject", Usage: "Email subject"},
					&cli.StringFlag{Name: "body", Usage: "Email body"},
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
					if cmd.IsSet("to_mail_ids") {
						body["to_mail_ids"] = cmd.String("to_mail_ids")
					}
					if cmd.IsSet("subject") {
						body["subject"] = cmd.String("subject")
					}
					if cmd.IsSet("body") {
						body["body"] = cmd.String("body")
					}
					if err := internal.MergeJSON(cmd, body); err != nil {
						return err
					}
					raw, err := c.Request("POST", c.BillingBase+"/invoices/"+cmd.Args().First()+"/email", &zohttp.RequestOpts{
						Params: map[string]string{"organization_id": orgID},
						JSON:   body,
					})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "collect-charge",
				Usage:     "Collect charge for an invoice",
				ArgsUsage: "<invoice-id>",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "json", Usage: "JSON body"},
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
					opts := &zohttp.RequestOpts{Params: map[string]string{"organization_id": orgID}}
					if j := cmd.String("json"); j != "" {
						var body any
						if err := json.Unmarshal([]byte(j), &body); err != nil {
							return internal.NewValidationError(fmt.Sprintf("--json: invalid JSON: %v", err))
						}
						opts.JSON = body
					}
					raw, err := c.Request("POST", c.BillingBase+"/invoices/"+cmd.Args().First()+"/collect", opts)
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
					raw, err := c.Request("POST", c.BillingBase+"/invoices/"+cmd.Args().First()+"/writeoff", &zohttp.RequestOpts{Params: map[string]string{"organization_id": orgID}})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "cancel-write-off",
				Usage:     "Cancel write off of an invoice",
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
					raw, err := c.Request("POST", c.BillingBase+"/invoices/"+cmd.Args().First()+"/cancelwriteoff", &zohttp.RequestOpts{Params: map[string]string{"organization_id": orgID}})
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
					&cli.StringFlag{Name: "apply_creditnotes", Required: true, Usage: "Credit notes JSON array [{creditnote_id, amount_applied}]"},
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
					var applyCreditnotes any
					if err := json.Unmarshal([]byte(cmd.String("apply_creditnotes")), &applyCreditnotes); err != nil {
						return internal.NewValidationError(fmt.Sprintf("--apply_creditnotes: invalid JSON: %v", err))
					}
					body["apply_creditnotes"] = applyCreditnotes
					if err := internal.MergeJSON(cmd, body); err != nil {
						return err
					}
					raw, err := c.Request("POST", c.BillingBase+"/invoices/"+cmd.Args().First()+"/credits", &zohttp.RequestOpts{
						Params: map[string]string{"organization_id": orgID},
						JSON:   body,
					})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "add-items",
				Usage:     "Add items to a pending invoice",
				ArgsUsage: "<invoice-id>",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "invoice_items", Required: true, Usage: "Invoice items JSON array [{code, name, price, quantity}]"},
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
					var invoiceItems any
					if err := json.Unmarshal([]byte(cmd.String("invoice_items")), &invoiceItems); err != nil {
						return internal.NewValidationError(fmt.Sprintf("--invoice_items: invalid JSON: %v", err))
					}
					body["invoice_items"] = invoiceItems
					if err := internal.MergeJSON(cmd, body); err != nil {
						return err
					}
					raw, err := c.Request("POST", c.BillingBase+"/invoices/"+cmd.Args().First()+"/lineitems", &zohttp.RequestOpts{
						Params: map[string]string{"organization_id": orgID},
						JSON:   body,
					})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "delete-item",
				Usage:     "Delete an item from a pending invoice",
				ArgsUsage: "<invoice-id> <line-item-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					if cmd.Args().Len() < 2 {
						return internal.NewValidationError("invoice ID and line item ID required")
					}
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := internal.RequireFlag(cmd, "org", "ZOHO_BOOKS_ORG_ID")
					if err != nil {
						return err
					}
					raw, err := c.Request("DELETE", c.BillingBase+"/invoices/"+cmd.Args().First()+"/lineitems/"+cmd.Args().Get(1), &zohttp.RequestOpts{Params: map[string]string{"organization_id": orgID}})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
		},
	}
}

func paymentsCmd() *cli.Command {
	return &cli.Command{
		Name:  "payments",
		Usage: "Payment operations",
		Commands: []*cli.Command{
			{
				Name:  "list",
				Usage: "List all payments",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "customer-id", Usage: "Filter by customer"},
					&cli.BoolFlag{Name: "all", Usage: "Fetch all records"},
					&cli.IntFlag{Name: "limit", Usage: "Max total records to fetch"},
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
					params := map[string]string{"organization_id": orgID}
					if v := cmd.String("customer-id"); v != "" {
						params["customer_id"] = v
					}
					if cmd.Bool("all") || cmd.IsSet("limit") {
						items, err := pagination.Paginate(pagination.PaginationConfig{
							Client:   c,
							URL:      c.BillingBase+"/payments",
							Opts:     &zohttp.RequestOpts{Params: params},
							ItemsKey: "payments",
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
					raw, err := c.Request("GET", c.BillingBase+"/payments", &zohttp.RequestOpts{Params: params})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "get",
				Usage:     "Retrieve a payment",
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
					raw, err := c.Request("GET", c.BillingBase+"/payments/"+cmd.Args().First(), &zohttp.RequestOpts{Params: map[string]string{"organization_id": orgID}})
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
				Usage: "List all credit notes",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "customer-id", Usage: "Filter by customer"},
					&cli.BoolFlag{Name: "all", Usage: "Fetch all records"},
					&cli.IntFlag{Name: "limit", Usage: "Max total records to fetch"},
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
					params := map[string]string{"organization_id": orgID}
					if v := cmd.String("customer-id"); v != "" {
						params["customer_id"] = v
					}
					if cmd.Bool("all") || cmd.IsSet("limit") {
						items, err := pagination.Paginate(pagination.PaginationConfig{
							Client:   c,
							URL:      c.BillingBase+"/creditnotes",
							Opts:     &zohttp.RequestOpts{Params: params},
							ItemsKey: "creditnotes",
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
					raw, err := c.Request("GET", c.BillingBase+"/creditnotes", &zohttp.RequestOpts{Params: params})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "get",
				Usage:     "Retrieve a credit note",
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
					raw, err := c.Request("GET", c.BillingBase+"/creditnotes/"+cmd.Args().First(), &zohttp.RequestOpts{Params: map[string]string{"organization_id": orgID}})
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
					&cli.StringFlag{Name: "reference_number", Usage: "Reference number"},
					&cli.StringFlag{Name: "date", Usage: "Credit note date"},
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
					if err := internal.MergeJSON(cmd, body); err != nil {
						return err
					}
					raw, err := c.Request("POST", c.BillingBase+"/creditnotes", &zohttp.RequestOpts{
						Params: map[string]string{"organization_id": orgID},
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
					raw, err := c.Request("DELETE", c.BillingBase+"/creditnotes/"+cmd.Args().First(), &zohttp.RequestOpts{Params: map[string]string{"organization_id": orgID}})
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
					&cli.StringFlag{Name: "to_mail_ids", Usage: "Comma-separated recipient emails"},
					&cli.StringFlag{Name: "subject", Usage: "Email subject"},
					&cli.StringFlag{Name: "body", Usage: "Email body"},
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
					if cmd.IsSet("to_mail_ids") {
						body["to_mail_ids"] = cmd.String("to_mail_ids")
					}
					if cmd.IsSet("subject") {
						body["subject"] = cmd.String("subject")
					}
					if cmd.IsSet("body") {
						body["body"] = cmd.String("body")
					}
					if err := internal.MergeJSON(cmd, body); err != nil {
						return err
					}
					raw, err := c.Request("POST", c.BillingBase+"/creditnotes/"+cmd.Args().First()+"/email", &zohttp.RequestOpts{
						Params: map[string]string{"organization_id": orgID},
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
					raw, err := c.Request("POST", c.BillingBase+"/creditnotes/"+cmd.Args().First()+"/void", &zohttp.RequestOpts{Params: map[string]string{"organization_id": orgID}})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "open-voided",
				Usage:     "Open a voided credit note",
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
					raw, err := c.Request("POST", c.BillingBase+"/creditnotes/"+cmd.Args().First()+"/open", &zohttp.RequestOpts{Params: map[string]string{"organization_id": orgID}})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "apply-credits",
				Usage:     "Apply credits to multiple invoices",
				ArgsUsage: "<creditnote-id>",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "invoices", Required: true, Usage: "Invoices JSON array [{invoice_id, amount_applied}]"},
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
					var invoices any
					if err := json.Unmarshal([]byte(cmd.String("invoices")), &invoices); err != nil {
						return internal.NewValidationError(fmt.Sprintf("--invoices: invalid JSON: %v", err))
					}
					body["invoices"] = invoices
					if err := internal.MergeJSON(cmd, body); err != nil {
						return err
					}
					raw, err := c.Request("POST", c.BillingBase+"/creditnotes/"+cmd.Args().First()+"/invoices", &zohttp.RequestOpts{
						Params: map[string]string{"organization_id": orgID},
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

func hostedPagesCmd() *cli.Command {
	return &cli.Command{
		Name:  "hosted-pages",
		Usage: "Hosted page operations",
		Commands: []*cli.Command{
			{
				Name:  "list",
				Usage: "List all hosted pages",
				Flags: []cli.Flag{
					&cli.BoolFlag{Name: "all", Usage: "Fetch all records"},
					&cli.IntFlag{Name: "limit", Usage: "Max total records to fetch"},
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
					params := map[string]string{"organization_id": orgID}
					if cmd.Bool("all") || cmd.IsSet("limit") {
						items, err := pagination.Paginate(pagination.PaginationConfig{
							Client:   c,
							URL:      c.BillingBase+"/hostedpages",
							Opts:     &zohttp.RequestOpts{Params: params},
							ItemsKey: "hostedpages",
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
					raw, err := c.Request("GET", c.BillingBase+"/hostedpages", &zohttp.RequestOpts{Params: params})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "get",
				Usage:     "Retrieve a hosted page",
				ArgsUsage: "<hostedpage-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					if cmd.Args().Len() < 1 {
						return internal.NewValidationError("hosted page ID required")
					}
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := internal.RequireFlag(cmd, "org", "ZOHO_BOOKS_ORG_ID")
					if err != nil {
						return err
					}
					raw, err := c.Request("GET", c.BillingBase+"/hostedpages/"+cmd.Args().First(), &zohttp.RequestOpts{Params: map[string]string{"organization_id": orgID}})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:  "create-subscription",
				Usage: "Create a subscription via hosted page",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "plan_code", Required: true, Usage: "Plan code"},
					&cli.StringFlag{Name: "customer_id", Usage: "Customer ID"},
					&cli.StringFlag{Name: "redirect_url", Usage: "Redirect URL"},
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
					if cmd.IsSet("customer_id") {
						body["customer_id"] = cmd.String("customer_id")
					}
					if cmd.IsSet("redirect_url") {
						body["redirect_url"] = cmd.String("redirect_url")
					}
					plan := map[string]any{}
					plan["plan_code"] = cmd.String("plan_code")
					body["plan"] = plan
					if err := internal.MergeJSON(cmd, body); err != nil {
						return err
					}
					raw, err := c.Request("POST", c.BillingBase+"/hostedpages/newsubscription", &zohttp.RequestOpts{
						Params: map[string]string{"organization_id": orgID},
						JSON:   body,
					})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:  "update-subscription",
				Usage: "Update a subscription via hosted page",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "subscription_id", Usage: "Subscription ID"},
					&cli.StringFlag{Name: "plan_code", Usage: "Plan code"},
					&cli.StringFlag{Name: "redirect_url", Usage: "Redirect URL"},
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
					if cmd.IsSet("subscription_id") {
						body["subscription_id"] = cmd.String("subscription_id")
					}
					if cmd.IsSet("redirect_url") {
						body["redirect_url"] = cmd.String("redirect_url")
					}
					plan := map[string]any{}
					if cmd.IsSet("plan_code") {
						plan["plan_code"] = cmd.String("plan_code")
					}
					if len(plan) > 0 {
						body["plan"] = plan
					}
					if err := internal.MergeJSON(cmd, body); err != nil {
						return err
					}
					raw, err := c.Request("POST", c.BillingBase+"/hostedpages/updatesubscription", &zohttp.RequestOpts{
						Params: map[string]string{"organization_id": orgID},
						JSON:   body,
					})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:  "update-card",
				Usage: "Update card via hosted page",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "subscription_id", Required: true, Usage: "Subscription ID"},
					&cli.StringFlag{Name: "redirect_url", Usage: "Redirect URL"},
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
					body["subscription_id"] = cmd.String("subscription_id")
					if cmd.IsSet("redirect_url") {
						body["redirect_url"] = cmd.String("redirect_url")
					}
					if err := internal.MergeJSON(cmd, body); err != nil {
						return err
					}
					raw, err := c.Request("POST", c.BillingBase+"/hostedpages/updatecard", &zohttp.RequestOpts{
						Params: map[string]string{"organization_id": orgID},
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

func eventsCmd() *cli.Command {
	return &cli.Command{
		Name:  "events",
		Usage: "Event operations",
		Commands: []*cli.Command{
			{
				Name:  "list",
				Usage: "List all events",
				Flags: []cli.Flag{
					&cli.BoolFlag{Name: "all", Usage: "Fetch all records"},
					&cli.IntFlag{Name: "limit", Usage: "Max total records to fetch"},
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
					params := map[string]string{"organization_id": orgID}
					if cmd.Bool("all") || cmd.IsSet("limit") {
						items, err := pagination.Paginate(pagination.PaginationConfig{
							Client:   c,
							URL:      c.BillingBase+"/events",
							Opts:     &zohttp.RequestOpts{Params: params},
							ItemsKey: "events",
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
					raw, err := c.Request("GET", c.BillingBase+"/events", &zohttp.RequestOpts{Params: params})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "get",
				Usage:     "Retrieve an event",
				ArgsUsage: "<event-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					if cmd.Args().Len() < 1 {
						return internal.NewValidationError("event ID required")
					}
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := internal.RequireFlag(cmd, "org", "ZOHO_BOOKS_ORG_ID")
					if err != nil {
						return err
					}
					raw, err := c.Request("GET", c.BillingBase+"/events/"+cmd.Args().First(), &zohttp.RequestOpts{Params: map[string]string{"organization_id": orgID}})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
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
				Usage: "List all organizations",
				Flags: []cli.Flag{
					&cli.BoolFlag{Name: "all", Usage: "Fetch all records"},
					&cli.IntFlag{Name: "limit", Usage: "Max total records to fetch"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					params := map[string]string{}
					if cmd.Bool("all") || cmd.IsSet("limit") {
						items, err := pagination.Paginate(pagination.PaginationConfig{
							Client:   c,
							URL:      c.BillingBase+"/organizations",
							Opts:     &zohttp.RequestOpts{Params: params},
							ItemsKey: "organizations",
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
					raw, err := c.Request("GET", c.BillingBase+"/organizations", &zohttp.RequestOpts{Params: params})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "get",
				Usage:     "Retrieve an organization",
				ArgsUsage: "<org-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					if cmd.Args().Len() < 1 {
						return internal.NewValidationError("organization ID required")
					}
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					raw, err := c.Request("GET", c.BillingBase+"/organizations/"+cmd.Args().First(), nil)
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
				Usage: "List all currencies",
				Flags: []cli.Flag{
					&cli.BoolFlag{Name: "all", Usage: "Fetch all records"},
					&cli.IntFlag{Name: "limit", Usage: "Max total records to fetch"},
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
					params := map[string]string{"organization_id": orgID}
					if cmd.Bool("all") || cmd.IsSet("limit") {
						items, err := pagination.Paginate(pagination.PaginationConfig{
							Client:   c,
							URL:      c.BillingBase+"/settings/currencies",
							Opts:     &zohttp.RequestOpts{Params: params},
							ItemsKey: "currencies",
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
					raw, err := c.Request("GET", c.BillingBase+"/settings/currencies", &zohttp.RequestOpts{Params: params})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "get",
				Usage:     "Retrieve a currency",
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
					raw, err := c.Request("GET", c.BillingBase+"/currencies/"+cmd.Args().First(), &zohttp.RequestOpts{Params: map[string]string{"organization_id": orgID}})
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
					&cli.StringFlag{Name: "currency_name", Required: true, Usage: "Currency name"},
					&cli.StringFlag{Name: "currency_symbol", Usage: "Currency symbol"},
					&cli.IntFlag{Name: "price_precision", Usage: "Price precision"},
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
					body["currency_name"] = cmd.String("currency_name")
					if cmd.IsSet("currency_symbol") {
						body["currency_symbol"] = cmd.String("currency_symbol")
					}
					if cmd.IsSet("price_precision") {
						body["price_precision"] = cmd.Int("price_precision")
					}
					if err := internal.MergeJSON(cmd, body); err != nil {
						return err
					}
					raw, err := c.Request("POST", c.BillingBase+"/settings/currencies", &zohttp.RequestOpts{
						Params: map[string]string{"organization_id": orgID},
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
				Usage:     "Update a currency",
				ArgsUsage: "<currency-id>",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "currency_name", Usage: "Currency name"},
					&cli.StringFlag{Name: "currency_symbol", Usage: "Currency symbol"},
					&cli.IntFlag{Name: "price_precision", Usage: "Price precision"},
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
					if cmd.IsSet("currency_name") {
						body["currency_name"] = cmd.String("currency_name")
					}
					if cmd.IsSet("currency_symbol") {
						body["currency_symbol"] = cmd.String("currency_symbol")
					}
					if cmd.IsSet("price_precision") {
						body["price_precision"] = cmd.Int("price_precision")
					}
					if err := internal.MergeJSON(cmd, body); err != nil {
						return err
					}
					raw, err := c.Request("PUT", c.BillingBase+"/currencies/"+cmd.Args().First(), &zohttp.RequestOpts{
						Params: map[string]string{"organization_id": orgID},
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
					raw, err := c.Request("DELETE", c.BillingBase+"/currencies/"+cmd.Args().First(), &zohttp.RequestOpts{Params: map[string]string{"organization_id": orgID}})
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
				Usage: "List all taxes",
				Flags: []cli.Flag{
					&cli.BoolFlag{Name: "all", Usage: "Fetch all records"},
					&cli.IntFlag{Name: "limit", Usage: "Max total records to fetch"},
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
					params := map[string]string{"organization_id": orgID}
					if cmd.Bool("all") || cmd.IsSet("limit") {
						items, err := pagination.Paginate(pagination.PaginationConfig{
							Client:   c,
							URL:      c.BillingBase+"/settings/taxes",
							Opts:     &zohttp.RequestOpts{Params: params},
							ItemsKey: "taxes",
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
					raw, err := c.Request("GET", c.BillingBase+"/settings/taxes", &zohttp.RequestOpts{Params: params})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "get",
				Usage:     "Retrieve a tax",
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
					raw, err := c.Request("GET", c.BillingBase+"/taxes/"+cmd.Args().First(), &zohttp.RequestOpts{Params: map[string]string{"organization_id": orgID}})
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
					orgID, err := internal.RequireFlag(cmd, "org", "ZOHO_BOOKS_ORG_ID")
					if err != nil {
						return err
					}
					body := map[string]any{}
					body["tax_name"] = cmd.String("tax_name")
					body["tax_percentage"] = cmd.Float("tax_percentage")
					if cmd.IsSet("tax_type") {
						body["tax_type"] = cmd.String("tax_type")
					}
					if err := internal.MergeJSON(cmd, body); err != nil {
						return err
					}
					raw, err := c.Request("POST", c.BillingBase+"/settings/taxes", &zohttp.RequestOpts{
						Params: map[string]string{"organization_id": orgID},
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
				Usage:     "Update a tax",
				ArgsUsage: "<tax-id>",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "tax_name", Usage: "Tax name"},
					&cli.FloatFlag{Name: "tax_percentage", Usage: "Tax percentage"},
					&cli.StringFlag{Name: "tax_type", Usage: "Tax type"},
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
					raw, err := c.Request("PUT", c.BillingBase+"/taxes/"+cmd.Args().First(), &zohttp.RequestOpts{
						Params: map[string]string{"organization_id": orgID},
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
					raw, err := c.Request("DELETE", c.BillingBase+"/taxes/"+cmd.Args().First(), &zohttp.RequestOpts{Params: map[string]string{"organization_id": orgID}})
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
				Usage: "List all users",
				Flags: []cli.Flag{
					&cli.BoolFlag{Name: "all", Usage: "Fetch all records"},
					&cli.IntFlag{Name: "limit", Usage: "Max total records to fetch"},
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
					params := map[string]string{"organization_id": orgID}
					if cmd.Bool("all") || cmd.IsSet("limit") {
						items, err := pagination.Paginate(pagination.PaginationConfig{
							Client:   c,
							URL:      c.BillingBase+"/users",
							Opts:     &zohttp.RequestOpts{Params: params},
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
					raw, err := c.Request("GET", c.BillingBase+"/users", &zohttp.RequestOpts{Params: params})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "get",
				Usage:     "Retrieve a user",
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
					raw, err := c.Request("GET", c.BillingBase+"/users/"+cmd.Args().First(), &zohttp.RequestOpts{Params: map[string]string{"organization_id": orgID}})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
		},
	}
}
