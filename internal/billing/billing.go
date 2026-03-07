package billing

import (
	"context"
	"encoding/json"
	"os"

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

func Commands() *cli.Command {
	return &cli.Command{
		Name:  "billing",
		Usage: "Zoho Billing operations",
		Flags: []cli.Flag{
			&cli.StringFlag{Name: "org", Usage: "Organization ID (or set ZOHO_BOOKS_ORG_ID)"},
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
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("GET", c.BillingBase+"/products", &zohttp.RequestOpts{Params: orgParams(orgID)})
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
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("GET", c.BillingBase+"/products/"+cmd.Args().First(), &zohttp.RequestOpts{Params: orgParams(orgID)})
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
					&cli.StringFlag{Name: "json", Required: true, Usage: "JSON body"},
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
					var body any
					json.Unmarshal([]byte(cmd.String("json")), &body)
					raw, err := c.Request("POST", c.BillingBase+"/products", &zohttp.RequestOpts{
						Params: orgParams(orgID),
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
					&cli.StringFlag{Name: "json", Required: true, Usage: "JSON body"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					if cmd.Args().Len() < 1 {
						return internal.NewValidationError("product ID required")
					}
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					var body any
					json.Unmarshal([]byte(cmd.String("json")), &body)
					raw, err := c.Request("PUT", c.BillingBase+"/products/"+cmd.Args().First(), &zohttp.RequestOpts{
						Params: orgParams(orgID),
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
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("DELETE", c.BillingBase+"/products/"+cmd.Args().First(), &zohttp.RequestOpts{Params: orgParams(orgID)})
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
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("POST", c.BillingBase+"/products/"+cmd.Args().First()+"/markasactive", &zohttp.RequestOpts{Params: orgParams(orgID)})
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
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("POST", c.BillingBase+"/products/"+cmd.Args().First()+"/markasinactive", &zohttp.RequestOpts{Params: orgParams(orgID)})
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
					if v := cmd.String("product-id"); v != "" {
						params["product_id"] = v
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
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("GET", c.BillingBase+"/plans/"+cmd.Args().First(), &zohttp.RequestOpts{Params: orgParams(orgID)})
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
					&cli.StringFlag{Name: "json", Required: true, Usage: "JSON body"},
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
					var body any
					json.Unmarshal([]byte(cmd.String("json")), &body)
					raw, err := c.Request("POST", c.BillingBase+"/plans", &zohttp.RequestOpts{
						Params: orgParams(orgID),
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
					&cli.StringFlag{Name: "json", Required: true, Usage: "JSON body"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					if cmd.Args().Len() < 1 {
						return internal.NewValidationError("plan code required")
					}
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					var body any
					json.Unmarshal([]byte(cmd.String("json")), &body)
					raw, err := c.Request("PUT", c.BillingBase+"/plans/"+cmd.Args().First(), &zohttp.RequestOpts{
						Params: orgParams(orgID),
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
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("DELETE", c.BillingBase+"/plans/"+cmd.Args().First(), &zohttp.RequestOpts{Params: orgParams(orgID)})
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
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("POST", c.BillingBase+"/plans/"+cmd.Args().First()+"/markasactive", &zohttp.RequestOpts{Params: orgParams(orgID)})
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
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("POST", c.BillingBase+"/plans/"+cmd.Args().First()+"/markasinactive", &zohttp.RequestOpts{Params: orgParams(orgID)})
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
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("GET", c.BillingBase+"/addons", &zohttp.RequestOpts{Params: orgParams(orgID)})
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
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("GET", c.BillingBase+"/addons/"+cmd.Args().First(), &zohttp.RequestOpts{Params: orgParams(orgID)})
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
					&cli.StringFlag{Name: "json", Required: true, Usage: "JSON body"},
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
					var body any
					json.Unmarshal([]byte(cmd.String("json")), &body)
					raw, err := c.Request("POST", c.BillingBase+"/addons", &zohttp.RequestOpts{
						Params: orgParams(orgID),
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
					&cli.StringFlag{Name: "json", Required: true, Usage: "JSON body"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					if cmd.Args().Len() < 1 {
						return internal.NewValidationError("addon code required")
					}
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					var body any
					json.Unmarshal([]byte(cmd.String("json")), &body)
					raw, err := c.Request("PUT", c.BillingBase+"/addons/"+cmd.Args().First(), &zohttp.RequestOpts{
						Params: orgParams(orgID),
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
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("DELETE", c.BillingBase+"/addons/"+cmd.Args().First(), &zohttp.RequestOpts{Params: orgParams(orgID)})
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
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("POST", c.BillingBase+"/addons/"+cmd.Args().First()+"/markasactive", &zohttp.RequestOpts{Params: orgParams(orgID)})
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
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("POST", c.BillingBase+"/addons/"+cmd.Args().First()+"/markasinactive", &zohttp.RequestOpts{Params: orgParams(orgID)})
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
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("GET", c.BillingBase+"/coupons", &zohttp.RequestOpts{Params: orgParams(orgID)})
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
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("GET", c.BillingBase+"/coupons/"+cmd.Args().First(), &zohttp.RequestOpts{Params: orgParams(orgID)})
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
					&cli.StringFlag{Name: "json", Required: true, Usage: "JSON body"},
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
					var body any
					json.Unmarshal([]byte(cmd.String("json")), &body)
					raw, err := c.Request("POST", c.BillingBase+"/coupons", &zohttp.RequestOpts{
						Params: orgParams(orgID),
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
					&cli.StringFlag{Name: "json", Required: true, Usage: "JSON body"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					if cmd.Args().Len() < 1 {
						return internal.NewValidationError("coupon code required")
					}
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					var body any
					json.Unmarshal([]byte(cmd.String("json")), &body)
					raw, err := c.Request("PUT", c.BillingBase+"/coupons/"+cmd.Args().First(), &zohttp.RequestOpts{
						Params: orgParams(orgID),
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
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("DELETE", c.BillingBase+"/coupons/"+cmd.Args().First(), &zohttp.RequestOpts{Params: orgParams(orgID)})
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
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("POST", c.BillingBase+"/coupons/"+cmd.Args().First()+"/markasactive", &zohttp.RequestOpts{Params: orgParams(orgID)})
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
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("POST", c.BillingBase+"/coupons/"+cmd.Args().First()+"/markasinactive", &zohttp.RequestOpts{Params: orgParams(orgID)})
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
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("GET", c.BillingBase+"/customers", &zohttp.RequestOpts{Params: orgParams(orgID)})
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
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("GET", c.BillingBase+"/customers/"+cmd.Args().First(), &zohttp.RequestOpts{Params: orgParams(orgID)})
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
					&cli.StringFlag{Name: "json", Required: true, Usage: "JSON body"},
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
					var body any
					json.Unmarshal([]byte(cmd.String("json")), &body)
					raw, err := c.Request("POST", c.BillingBase+"/customers", &zohttp.RequestOpts{
						Params: orgParams(orgID),
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
					&cli.StringFlag{Name: "json", Required: true, Usage: "JSON body"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					if cmd.Args().Len() < 1 {
						return internal.NewValidationError("customer ID required")
					}
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					var body any
					json.Unmarshal([]byte(cmd.String("json")), &body)
					raw, err := c.Request("PUT", c.BillingBase+"/customers/"+cmd.Args().First(), &zohttp.RequestOpts{
						Params: orgParams(orgID),
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
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("DELETE", c.BillingBase+"/customers/"+cmd.Args().First(), &zohttp.RequestOpts{Params: orgParams(orgID)})
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
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("POST", c.BillingBase+"/customers/"+cmd.Args().First()+"/markasactive", &zohttp.RequestOpts{Params: orgParams(orgID)})
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
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("POST", c.BillingBase+"/customers/"+cmd.Args().First()+"/markasinactive", &zohttp.RequestOpts{Params: orgParams(orgID)})
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
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("GET", c.BillingBase+"/subscriptions/"+cmd.Args().First(), &zohttp.RequestOpts{Params: orgParams(orgID)})
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
					&cli.StringFlag{Name: "json", Required: true, Usage: "JSON body"},
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
					var body any
					json.Unmarshal([]byte(cmd.String("json")), &body)
					raw, err := c.Request("POST", c.BillingBase+"/subscriptions", &zohttp.RequestOpts{
						Params: orgParams(orgID),
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
					&cli.StringFlag{Name: "json", Required: true, Usage: "JSON body"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					if cmd.Args().Len() < 1 {
						return internal.NewValidationError("subscription ID required")
					}
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					var body any
					json.Unmarshal([]byte(cmd.String("json")), &body)
					raw, err := c.Request("PUT", c.BillingBase+"/subscriptions/"+cmd.Args().First(), &zohttp.RequestOpts{
						Params: orgParams(orgID),
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
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("DELETE", c.BillingBase+"/subscriptions/"+cmd.Args().First(), &zohttp.RequestOpts{Params: orgParams(orgID)})
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
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					opts := &zohttp.RequestOpts{Params: orgParams(orgID)}
					if j := cmd.String("json"); j != "" {
						var body any
						json.Unmarshal([]byte(j), &body)
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
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("POST", c.BillingBase+"/subscriptions/"+cmd.Args().First()+"/reactivate", &zohttp.RequestOpts{Params: orgParams(orgID)})
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
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("GET", c.BillingBase+"/subscriptions/"+cmd.Args().First()+"/scheduledchanges", &zohttp.RequestOpts{Params: orgParams(orgID)})
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
					if v := cmd.String("subscription-id"); v != "" {
						params["subscription_id"] = v
					}
					if v := cmd.String("customer-id"); v != "" {
						params["customer_id"] = v
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
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("GET", c.BillingBase+"/invoices/"+cmd.Args().First(), &zohttp.RequestOpts{Params: orgParams(orgID)})
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
					&cli.StringFlag{Name: "json", Required: true, Usage: "JSON body"},
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
					var body any
					json.Unmarshal([]byte(cmd.String("json")), &body)
					raw, err := c.Request("POST", c.BillingBase+"/invoices", &zohttp.RequestOpts{
						Params: orgParams(orgID),
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
					&cli.StringFlag{Name: "json", Required: true, Usage: "JSON body"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					if cmd.Args().Len() < 1 {
						return internal.NewValidationError("invoice ID required")
					}
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					var body any
					json.Unmarshal([]byte(cmd.String("json")), &body)
					raw, err := c.Request("PUT", c.BillingBase+"/invoices/"+cmd.Args().First(), &zohttp.RequestOpts{
						Params: orgParams(orgID),
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
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("DELETE", c.BillingBase+"/invoices/"+cmd.Args().First(), &zohttp.RequestOpts{Params: orgParams(orgID)})
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
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("POST", c.BillingBase+"/invoices/"+cmd.Args().First()+"/converttoopen", &zohttp.RequestOpts{Params: orgParams(orgID)})
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
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("POST", c.BillingBase+"/invoices/"+cmd.Args().First()+"/void", &zohttp.RequestOpts{Params: orgParams(orgID)})
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
					&cli.StringFlag{Name: "json", Required: true, Usage: "JSON body"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					if cmd.Args().Len() < 1 {
						return internal.NewValidationError("invoice ID required")
					}
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					var body any
					json.Unmarshal([]byte(cmd.String("json")), &body)
					raw, err := c.Request("POST", c.BillingBase+"/invoices/"+cmd.Args().First()+"/email", &zohttp.RequestOpts{
						Params: orgParams(orgID),
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
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					opts := &zohttp.RequestOpts{Params: orgParams(orgID)}
					if j := cmd.String("json"); j != "" {
						var body any
						json.Unmarshal([]byte(j), &body)
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
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("POST", c.BillingBase+"/invoices/"+cmd.Args().First()+"/writeoff", &zohttp.RequestOpts{Params: orgParams(orgID)})
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
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("POST", c.BillingBase+"/invoices/"+cmd.Args().First()+"/cancelwriteoff", &zohttp.RequestOpts{Params: orgParams(orgID)})
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
					&cli.StringFlag{Name: "json", Required: true, Usage: "JSON body"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					if cmd.Args().Len() < 1 {
						return internal.NewValidationError("invoice ID required")
					}
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					var body any
					json.Unmarshal([]byte(cmd.String("json")), &body)
					raw, err := c.Request("POST", c.BillingBase+"/invoices/"+cmd.Args().First()+"/credits", &zohttp.RequestOpts{
						Params: orgParams(orgID),
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
					&cli.StringFlag{Name: "json", Required: true, Usage: "JSON body"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					if cmd.Args().Len() < 1 {
						return internal.NewValidationError("invoice ID required")
					}
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					var body any
					json.Unmarshal([]byte(cmd.String("json")), &body)
					raw, err := c.Request("POST", c.BillingBase+"/invoices/"+cmd.Args().First()+"/lineitems", &zohttp.RequestOpts{
						Params: orgParams(orgID),
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
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("DELETE", c.BillingBase+"/invoices/"+cmd.Args().First()+"/lineitems/"+cmd.Args().Get(1), &zohttp.RequestOpts{Params: orgParams(orgID)})
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
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("GET", c.BillingBase+"/payments/"+cmd.Args().First(), &zohttp.RequestOpts{Params: orgParams(orgID)})
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
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("GET", c.BillingBase+"/creditnotes/"+cmd.Args().First(), &zohttp.RequestOpts{Params: orgParams(orgID)})
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
					&cli.StringFlag{Name: "json", Required: true, Usage: "JSON body"},
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
					var body any
					json.Unmarshal([]byte(cmd.String("json")), &body)
					raw, err := c.Request("POST", c.BillingBase+"/creditnotes", &zohttp.RequestOpts{
						Params: orgParams(orgID),
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
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("DELETE", c.BillingBase+"/creditnotes/"+cmd.Args().First(), &zohttp.RequestOpts{Params: orgParams(orgID)})
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
					&cli.StringFlag{Name: "json", Required: true, Usage: "JSON body"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					if cmd.Args().Len() < 1 {
						return internal.NewValidationError("credit note ID required")
					}
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					var body any
					json.Unmarshal([]byte(cmd.String("json")), &body)
					raw, err := c.Request("POST", c.BillingBase+"/creditnotes/"+cmd.Args().First()+"/email", &zohttp.RequestOpts{
						Params: orgParams(orgID),
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
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("POST", c.BillingBase+"/creditnotes/"+cmd.Args().First()+"/void", &zohttp.RequestOpts{Params: orgParams(orgID)})
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
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("POST", c.BillingBase+"/creditnotes/"+cmd.Args().First()+"/open", &zohttp.RequestOpts{Params: orgParams(orgID)})
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
					&cli.StringFlag{Name: "json", Required: true, Usage: "JSON body"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					if cmd.Args().Len() < 1 {
						return internal.NewValidationError("credit note ID required")
					}
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					var body any
					json.Unmarshal([]byte(cmd.String("json")), &body)
					raw, err := c.Request("POST", c.BillingBase+"/creditnotes/"+cmd.Args().First()+"/invoices", &zohttp.RequestOpts{
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

func hostedPagesCmd() *cli.Command {
	return &cli.Command{
		Name:  "hosted-pages",
		Usage: "Hosted page operations",
		Commands: []*cli.Command{
			{
				Name:  "list",
				Usage: "List all hosted pages",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("GET", c.BillingBase+"/hostedpages", &zohttp.RequestOpts{Params: orgParams(orgID)})
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
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("GET", c.BillingBase+"/hostedpages/"+cmd.Args().First(), &zohttp.RequestOpts{Params: orgParams(orgID)})
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
					&cli.StringFlag{Name: "json", Required: true, Usage: "JSON body"},
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
					var body any
					json.Unmarshal([]byte(cmd.String("json")), &body)
					raw, err := c.Request("POST", c.BillingBase+"/hostedpages/newsubscription", &zohttp.RequestOpts{
						Params: orgParams(orgID),
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
					&cli.StringFlag{Name: "json", Required: true, Usage: "JSON body"},
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
					var body any
					json.Unmarshal([]byte(cmd.String("json")), &body)
					raw, err := c.Request("POST", c.BillingBase+"/hostedpages/updatesubscription", &zohttp.RequestOpts{
						Params: orgParams(orgID),
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
					&cli.StringFlag{Name: "json", Required: true, Usage: "JSON body"},
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
					var body any
					json.Unmarshal([]byte(cmd.String("json")), &body)
					raw, err := c.Request("POST", c.BillingBase+"/hostedpages/updatecard", &zohttp.RequestOpts{
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

func eventsCmd() *cli.Command {
	return &cli.Command{
		Name:  "events",
		Usage: "Event operations",
		Commands: []*cli.Command{
			{
				Name:  "list",
				Usage: "List all events",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("GET", c.BillingBase+"/events", &zohttp.RequestOpts{Params: orgParams(orgID)})
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
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("GET", c.BillingBase+"/events/"+cmd.Args().First(), &zohttp.RequestOpts{Params: orgParams(orgID)})
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
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					raw, err := c.Request("GET", c.BillingBase+"/organizations", nil)
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
					c, err := getClient()
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
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("GET", c.BillingBase+"/currencies", &zohttp.RequestOpts{Params: orgParams(orgID)})
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
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("GET", c.BillingBase+"/currencies/"+cmd.Args().First(), &zohttp.RequestOpts{Params: orgParams(orgID)})
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
					&cli.StringFlag{Name: "json", Required: true, Usage: "JSON body"},
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
					var body any
					json.Unmarshal([]byte(cmd.String("json")), &body)
					raw, err := c.Request("POST", c.BillingBase+"/currencies", &zohttp.RequestOpts{
						Params: orgParams(orgID),
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
					&cli.StringFlag{Name: "json", Required: true, Usage: "JSON body"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					if cmd.Args().Len() < 1 {
						return internal.NewValidationError("currency ID required")
					}
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					var body any
					json.Unmarshal([]byte(cmd.String("json")), &body)
					raw, err := c.Request("PUT", c.BillingBase+"/currencies/"+cmd.Args().First(), &zohttp.RequestOpts{
						Params: orgParams(orgID),
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
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("DELETE", c.BillingBase+"/currencies/"+cmd.Args().First(), &zohttp.RequestOpts{Params: orgParams(orgID)})
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
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("GET", c.BillingBase+"/taxes", &zohttp.RequestOpts{Params: orgParams(orgID)})
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
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("GET", c.BillingBase+"/taxes/"+cmd.Args().First(), &zohttp.RequestOpts{Params: orgParams(orgID)})
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
					&cli.StringFlag{Name: "json", Required: true, Usage: "JSON body"},
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
					var body any
					json.Unmarshal([]byte(cmd.String("json")), &body)
					raw, err := c.Request("POST", c.BillingBase+"/taxes", &zohttp.RequestOpts{
						Params: orgParams(orgID),
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
					&cli.StringFlag{Name: "json", Required: true, Usage: "JSON body"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					if cmd.Args().Len() < 1 {
						return internal.NewValidationError("tax ID required")
					}
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					var body any
					json.Unmarshal([]byte(cmd.String("json")), &body)
					raw, err := c.Request("PUT", c.BillingBase+"/taxes/"+cmd.Args().First(), &zohttp.RequestOpts{
						Params: orgParams(orgID),
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
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("DELETE", c.BillingBase+"/taxes/"+cmd.Args().First(), &zohttp.RequestOpts{Params: orgParams(orgID)})
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
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("GET", c.BillingBase+"/users", &zohttp.RequestOpts{Params: orgParams(orgID)})
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
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("GET", c.BillingBase+"/users/"+cmd.Args().First(), &zohttp.RequestOpts{Params: orgParams(orgID)})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
		},
	}
}
