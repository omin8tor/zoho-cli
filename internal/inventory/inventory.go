package inventory

import (
	"context"
	"encoding/json"
	"fmt"
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
		Name:  "inventory",
		Usage: "Zoho Inventory operations",
		Flags: []cli.Flag{
			&cli.StringFlag{Name: "org", Usage: "Organization ID (or set ZOHO_BOOKS_ORG_ID)"},
		},
		Commands: []*cli.Command{
			itemsCmd(),
			compositeItemsCmd(),
			itemGroupsCmd(),
			contactsCmd(),
			salesOrdersCmd(),
			invoicesCmd(),
			packagesCmd(),
			shipmentOrdersCmd(),
			purchaseOrdersCmd(),
			billsCmd(),
			vendorCreditsCmd(),
			priceListsCmd(),
			warehousesCmd(),
			transferOrdersCmd(),
			adjustmentsCmd(),
			organizationsCmd(),
			currenciesCmd(),
			taxesCmd(),
			usersCmd(),
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
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("GET", c.InventoryBase+"/items", &zohttp.RequestOpts{Params: orgParams(orgID)})
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
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("GET", c.InventoryBase+"/items/"+cmd.Args().First(), &zohttp.RequestOpts{Params: orgParams(orgID)})
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
					if err := json.Unmarshal([]byte(cmd.String("json")), &body); err != nil {
						return internal.NewValidationError(fmt.Sprintf("invalid JSON: %v", err))
					}
					raw, err := c.Request("POST", c.InventoryBase+"/items", &zohttp.RequestOpts{
						Params: orgParams(orgID),
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
				Usage:     "Update an item",
				ArgsUsage: "<item-id>",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "json", Required: true, Usage: "JSON body"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					if cmd.Args().Len() < 1 {
						return internal.NewValidationError("item ID required")
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
					if err := json.Unmarshal([]byte(cmd.String("json")), &body); err != nil {
						return internal.NewValidationError(fmt.Sprintf("invalid JSON: %v", err))
					}
					raw, err := c.Request("PUT", c.InventoryBase+"/items/"+cmd.Args().First(), &zohttp.RequestOpts{
						Params: orgParams(orgID),
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
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("DELETE", c.InventoryBase+"/items/"+cmd.Args().First(), &zohttp.RequestOpts{Params: orgParams(orgID)})
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
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("POST", c.InventoryBase+"/items/"+cmd.Args().First()+"/active", &zohttp.RequestOpts{Params: orgParams(orgID)})
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
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("POST", c.InventoryBase+"/items/"+cmd.Args().First()+"/inactive", &zohttp.RequestOpts{Params: orgParams(orgID)})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
		},
	}
}

func compositeItemsCmd() *cli.Command {
	return &cli.Command{
		Name:  "composite-items",
		Usage: "Composite item operations",
		Commands: []*cli.Command{
			{
				Name:  "list",
				Usage: "List composite items",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("GET", c.InventoryBase+"/compositeitems", &zohttp.RequestOpts{Params: orgParams(orgID)})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "get",
				Usage:     "Get a composite item",
				ArgsUsage: "<composite-item-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					if cmd.Args().Len() < 1 {
						return internal.NewValidationError("composite item ID required")
					}
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("GET", c.InventoryBase+"/compositeitems/"+cmd.Args().First(), &zohttp.RequestOpts{Params: orgParams(orgID)})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:  "create",
				Usage: "Create a composite item",
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
					if err := json.Unmarshal([]byte(cmd.String("json")), &body); err != nil {
						return internal.NewValidationError(fmt.Sprintf("invalid JSON: %v", err))
					}
					raw, err := c.Request("POST", c.InventoryBase+"/compositeitems", &zohttp.RequestOpts{
						Params: orgParams(orgID),
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
				Usage:     "Update a composite item",
				ArgsUsage: "<composite-item-id>",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "json", Required: true, Usage: "JSON body"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					if cmd.Args().Len() < 1 {
						return internal.NewValidationError("composite item ID required")
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
					if err := json.Unmarshal([]byte(cmd.String("json")), &body); err != nil {
						return internal.NewValidationError(fmt.Sprintf("invalid JSON: %v", err))
					}
					raw, err := c.Request("PUT", c.InventoryBase+"/compositeitems/"+cmd.Args().First(), &zohttp.RequestOpts{
						Params: orgParams(orgID),
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
				Usage:     "Delete a composite item",
				ArgsUsage: "<composite-item-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					if cmd.Args().Len() < 1 {
						return internal.NewValidationError("composite item ID required")
					}
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("DELETE", c.InventoryBase+"/compositeitems/"+cmd.Args().First(), &zohttp.RequestOpts{Params: orgParams(orgID)})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
		},
	}
}

func itemGroupsCmd() *cli.Command {
	return &cli.Command{
		Name:  "item-groups",
		Usage: "Item group operations",
		Commands: []*cli.Command{
			{
				Name:  "list",
				Usage: "List item groups",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("GET", c.InventoryBase+"/itemgroups", &zohttp.RequestOpts{Params: orgParams(orgID)})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "get",
				Usage:     "Get an item group",
				ArgsUsage: "<group-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					if cmd.Args().Len() < 1 {
						return internal.NewValidationError("item group ID required")
					}
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("GET", c.InventoryBase+"/itemgroups/"+cmd.Args().First(), &zohttp.RequestOpts{Params: orgParams(orgID)})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:  "create",
				Usage: "Create an item group",
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
					if err := json.Unmarshal([]byte(cmd.String("json")), &body); err != nil {
						return internal.NewValidationError(fmt.Sprintf("invalid JSON: %v", err))
					}
					raw, err := c.Request("POST", c.InventoryBase+"/itemgroups", &zohttp.RequestOpts{
						Params: orgParams(orgID),
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
				Usage:     "Update an item group",
				ArgsUsage: "<group-id>",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "json", Required: true, Usage: "JSON body"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					if cmd.Args().Len() < 1 {
						return internal.NewValidationError("item group ID required")
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
					if err := json.Unmarshal([]byte(cmd.String("json")), &body); err != nil {
						return internal.NewValidationError(fmt.Sprintf("invalid JSON: %v", err))
					}
					raw, err := c.Request("PUT", c.InventoryBase+"/itemgroups/"+cmd.Args().First(), &zohttp.RequestOpts{
						Params: orgParams(orgID),
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
				Usage:     "Delete an item group",
				ArgsUsage: "<group-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					if cmd.Args().Len() < 1 {
						return internal.NewValidationError("item group ID required")
					}
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("DELETE", c.InventoryBase+"/itemgroups/"+cmd.Args().First(), &zohttp.RequestOpts{Params: orgParams(orgID)})
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
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("GET", c.InventoryBase+"/contacts", &zohttp.RequestOpts{Params: orgParams(orgID)})
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
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("GET", c.InventoryBase+"/contacts/"+cmd.Args().First(), &zohttp.RequestOpts{Params: orgParams(orgID)})
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
					if err := json.Unmarshal([]byte(cmd.String("json")), &body); err != nil {
						return internal.NewValidationError(fmt.Sprintf("invalid JSON: %v", err))
					}
					raw, err := c.Request("POST", c.InventoryBase+"/contacts", &zohttp.RequestOpts{
						Params: orgParams(orgID),
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
				Usage:     "Update a contact",
				ArgsUsage: "<contact-id>",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "json", Required: true, Usage: "JSON body"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					if cmd.Args().Len() < 1 {
						return internal.NewValidationError("contact ID required")
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
					if err := json.Unmarshal([]byte(cmd.String("json")), &body); err != nil {
						return internal.NewValidationError(fmt.Sprintf("invalid JSON: %v", err))
					}
					raw, err := c.Request("PUT", c.InventoryBase+"/contacts/"+cmd.Args().First(), &zohttp.RequestOpts{
						Params: orgParams(orgID),
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
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("DELETE", c.InventoryBase+"/contacts/"+cmd.Args().First(), &zohttp.RequestOpts{Params: orgParams(orgID)})
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
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("POST", c.InventoryBase+"/contacts/"+cmd.Args().First()+"/active", &zohttp.RequestOpts{Params: orgParams(orgID)})
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
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("POST", c.InventoryBase+"/contacts/"+cmd.Args().First()+"/inactive", &zohttp.RequestOpts{Params: orgParams(orgID)})
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
				Name:  "list",
				Usage: "List sales orders",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("GET", c.InventoryBase+"/salesorders", &zohttp.RequestOpts{Params: orgParams(orgID)})
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
					if cmd.Args().Len() < 1 {
						return internal.NewValidationError("sales order ID required")
					}
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("GET", c.InventoryBase+"/salesorders/"+cmd.Args().First(), &zohttp.RequestOpts{Params: orgParams(orgID)})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:  "create",
				Usage: "Create a sales order",
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
					if err := json.Unmarshal([]byte(cmd.String("json")), &body); err != nil {
						return internal.NewValidationError(fmt.Sprintf("invalid JSON: %v", err))
					}
					raw, err := c.Request("POST", c.InventoryBase+"/salesorders", &zohttp.RequestOpts{
						Params: orgParams(orgID),
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
				Usage:     "Update a sales order",
				ArgsUsage: "<salesorder-id>",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "json", Required: true, Usage: "JSON body"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					if cmd.Args().Len() < 1 {
						return internal.NewValidationError("sales order ID required")
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
					if err := json.Unmarshal([]byte(cmd.String("json")), &body); err != nil {
						return internal.NewValidationError(fmt.Sprintf("invalid JSON: %v", err))
					}
					raw, err := c.Request("PUT", c.InventoryBase+"/salesorders/"+cmd.Args().First(), &zohttp.RequestOpts{
						Params: orgParams(orgID),
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
				Usage:     "Delete a sales order",
				ArgsUsage: "<salesorder-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					if cmd.Args().Len() < 1 {
						return internal.NewValidationError("sales order ID required")
					}
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("DELETE", c.InventoryBase+"/salesorders/"+cmd.Args().First(), &zohttp.RequestOpts{Params: orgParams(orgID)})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "mark-confirmed",
				Usage:     "Mark a sales order as confirmed",
				ArgsUsage: "<salesorder-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					if cmd.Args().Len() < 1 {
						return internal.NewValidationError("sales order ID required")
					}
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("POST", c.InventoryBase+"/salesorders/"+cmd.Args().First()+"/status/confirmed", &zohttp.RequestOpts{Params: orgParams(orgID)})
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
					if cmd.Args().Len() < 1 {
						return internal.NewValidationError("sales order ID required")
					}
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("POST", c.InventoryBase+"/salesorders/"+cmd.Args().First()+"/status/void", &zohttp.RequestOpts{Params: orgParams(orgID)})
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
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("GET", c.InventoryBase+"/invoices", &zohttp.RequestOpts{Params: orgParams(orgID)})
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
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("GET", c.InventoryBase+"/invoices/"+cmd.Args().First(), &zohttp.RequestOpts{Params: orgParams(orgID)})
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
					if err := json.Unmarshal([]byte(cmd.String("json")), &body); err != nil {
						return internal.NewValidationError(fmt.Sprintf("invalid JSON: %v", err))
					}
					raw, err := c.Request("POST", c.InventoryBase+"/invoices", &zohttp.RequestOpts{
						Params: orgParams(orgID),
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
					if err := json.Unmarshal([]byte(cmd.String("json")), &body); err != nil {
						return internal.NewValidationError(fmt.Sprintf("invalid JSON: %v", err))
					}
					raw, err := c.Request("PUT", c.InventoryBase+"/invoices/"+cmd.Args().First(), &zohttp.RequestOpts{
						Params: orgParams(orgID),
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
					raw, err := c.Request("DELETE", c.InventoryBase+"/invoices/"+cmd.Args().First(), &zohttp.RequestOpts{Params: orgParams(orgID)})
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
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("POST", c.InventoryBase+"/invoices/"+cmd.Args().First()+"/status/sent", &zohttp.RequestOpts{Params: orgParams(orgID)})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "mark-void",
				Usage:     "Mark an invoice as void",
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
					raw, err := c.Request("POST", c.InventoryBase+"/invoices/"+cmd.Args().First()+"/status/void", &zohttp.RequestOpts{Params: orgParams(orgID)})
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
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("POST", c.InventoryBase+"/invoices/"+cmd.Args().First()+"/status/draft", &zohttp.RequestOpts{Params: orgParams(orgID)})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
		},
	}
}

func packagesCmd() *cli.Command {
	return &cli.Command{
		Name:  "packages",
		Usage: "Package operations",
		Commands: []*cli.Command{
			{
				Name:  "list",
				Usage: "List packages",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("GET", c.InventoryBase+"/packages", &zohttp.RequestOpts{Params: orgParams(orgID)})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "get",
				Usage:     "Get a package",
				ArgsUsage: "<package-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					if cmd.Args().Len() < 1 {
						return internal.NewValidationError("package ID required")
					}
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("GET", c.InventoryBase+"/packages/"+cmd.Args().First(), &zohttp.RequestOpts{Params: orgParams(orgID)})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:  "create",
				Usage: "Create a package",
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
					if err := json.Unmarshal([]byte(cmd.String("json")), &body); err != nil {
						return internal.NewValidationError(fmt.Sprintf("invalid JSON: %v", err))
					}
					raw, err := c.Request("POST", c.InventoryBase+"/packages", &zohttp.RequestOpts{
						Params: orgParams(orgID),
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
				Usage:     "Update a package",
				ArgsUsage: "<package-id>",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "json", Required: true, Usage: "JSON body"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					if cmd.Args().Len() < 1 {
						return internal.NewValidationError("package ID required")
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
					if err := json.Unmarshal([]byte(cmd.String("json")), &body); err != nil {
						return internal.NewValidationError(fmt.Sprintf("invalid JSON: %v", err))
					}
					raw, err := c.Request("PUT", c.InventoryBase+"/packages/"+cmd.Args().First(), &zohttp.RequestOpts{
						Params: orgParams(orgID),
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
				Usage:     "Delete a package",
				ArgsUsage: "<package-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					if cmd.Args().Len() < 1 {
						return internal.NewValidationError("package ID required")
					}
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("DELETE", c.InventoryBase+"/packages/"+cmd.Args().First(), &zohttp.RequestOpts{Params: orgParams(orgID)})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
		},
	}
}

func shipmentOrdersCmd() *cli.Command {
	return &cli.Command{
		Name:  "shipment-orders",
		Usage: "Shipment order operations",
		Commands: []*cli.Command{
			{
				Name:  "list",
				Usage: "List shipment orders",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("GET", c.InventoryBase+"/shipmentorders", &zohttp.RequestOpts{Params: orgParams(orgID)})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "get",
				Usage:     "Get a shipment order",
				ArgsUsage: "<shipmentorder-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					if cmd.Args().Len() < 1 {
						return internal.NewValidationError("shipment order ID required")
					}
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("GET", c.InventoryBase+"/shipmentorders/"+cmd.Args().First(), &zohttp.RequestOpts{Params: orgParams(orgID)})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:  "create",
				Usage: "Create a shipment order",
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
					if err := json.Unmarshal([]byte(cmd.String("json")), &body); err != nil {
						return internal.NewValidationError(fmt.Sprintf("invalid JSON: %v", err))
					}
					raw, err := c.Request("POST", c.InventoryBase+"/shipmentorders", &zohttp.RequestOpts{
						Params: orgParams(orgID),
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
				Usage:     "Delete a shipment order",
				ArgsUsage: "<shipmentorder-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					if cmd.Args().Len() < 1 {
						return internal.NewValidationError("shipment order ID required")
					}
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("DELETE", c.InventoryBase+"/shipmentorders/"+cmd.Args().First(), &zohttp.RequestOpts{Params: orgParams(orgID)})
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
				Name:  "list",
				Usage: "List purchase orders",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("GET", c.InventoryBase+"/purchaseorders", &zohttp.RequestOpts{Params: orgParams(orgID)})
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
					if cmd.Args().Len() < 1 {
						return internal.NewValidationError("purchase order ID required")
					}
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("GET", c.InventoryBase+"/purchaseorders/"+cmd.Args().First(), &zohttp.RequestOpts{Params: orgParams(orgID)})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:  "create",
				Usage: "Create a purchase order",
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
					if err := json.Unmarshal([]byte(cmd.String("json")), &body); err != nil {
						return internal.NewValidationError(fmt.Sprintf("invalid JSON: %v", err))
					}
					raw, err := c.Request("POST", c.InventoryBase+"/purchaseorders", &zohttp.RequestOpts{
						Params: orgParams(orgID),
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
				Usage:     "Update a purchase order",
				ArgsUsage: "<purchaseorder-id>",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "json", Required: true, Usage: "JSON body"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					if cmd.Args().Len() < 1 {
						return internal.NewValidationError("purchase order ID required")
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
					if err := json.Unmarshal([]byte(cmd.String("json")), &body); err != nil {
						return internal.NewValidationError(fmt.Sprintf("invalid JSON: %v", err))
					}
					raw, err := c.Request("PUT", c.InventoryBase+"/purchaseorders/"+cmd.Args().First(), &zohttp.RequestOpts{
						Params: orgParams(orgID),
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
				Usage:     "Delete a purchase order",
				ArgsUsage: "<purchaseorder-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					if cmd.Args().Len() < 1 {
						return internal.NewValidationError("purchase order ID required")
					}
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("DELETE", c.InventoryBase+"/purchaseorders/"+cmd.Args().First(), &zohttp.RequestOpts{Params: orgParams(orgID)})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "mark-issued",
				Usage:     "Mark a purchase order as issued",
				ArgsUsage: "<purchaseorder-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					if cmd.Args().Len() < 1 {
						return internal.NewValidationError("purchase order ID required")
					}
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("POST", c.InventoryBase+"/purchaseorders/"+cmd.Args().First()+"/status/issued", &zohttp.RequestOpts{Params: orgParams(orgID)})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "mark-cancelled",
				Usage:     "Mark a purchase order as cancelled",
				ArgsUsage: "<purchaseorder-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					if cmd.Args().Len() < 1 {
						return internal.NewValidationError("purchase order ID required")
					}
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("POST", c.InventoryBase+"/purchaseorders/"+cmd.Args().First()+"/status/cancelled", &zohttp.RequestOpts{Params: orgParams(orgID)})
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
				Name:  "list",
				Usage: "List bills",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("GET", c.InventoryBase+"/bills", &zohttp.RequestOpts{Params: orgParams(orgID)})
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
					if cmd.Args().Len() < 1 {
						return internal.NewValidationError("bill ID required")
					}
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("GET", c.InventoryBase+"/bills/"+cmd.Args().First(), &zohttp.RequestOpts{Params: orgParams(orgID)})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:  "create",
				Usage: "Create a bill",
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
					if err := json.Unmarshal([]byte(cmd.String("json")), &body); err != nil {
						return internal.NewValidationError(fmt.Sprintf("invalid JSON: %v", err))
					}
					raw, err := c.Request("POST", c.InventoryBase+"/bills", &zohttp.RequestOpts{
						Params: orgParams(orgID),
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
				Usage:     "Update a bill",
				ArgsUsage: "<bill-id>",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "json", Required: true, Usage: "JSON body"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					if cmd.Args().Len() < 1 {
						return internal.NewValidationError("bill ID required")
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
					if err := json.Unmarshal([]byte(cmd.String("json")), &body); err != nil {
						return internal.NewValidationError(fmt.Sprintf("invalid JSON: %v", err))
					}
					raw, err := c.Request("PUT", c.InventoryBase+"/bills/"+cmd.Args().First(), &zohttp.RequestOpts{
						Params: orgParams(orgID),
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
				Usage:     "Delete a bill",
				ArgsUsage: "<bill-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					if cmd.Args().Len() < 1 {
						return internal.NewValidationError("bill ID required")
					}
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("DELETE", c.InventoryBase+"/bills/"+cmd.Args().First(), &zohttp.RequestOpts{Params: orgParams(orgID)})
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
					if cmd.Args().Len() < 1 {
						return internal.NewValidationError("bill ID required")
					}
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("POST", c.InventoryBase+"/bills/"+cmd.Args().First()+"/status/open", &zohttp.RequestOpts{Params: orgParams(orgID)})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "mark-void",
				Usage:     "Mark a bill as void",
				ArgsUsage: "<bill-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					if cmd.Args().Len() < 1 {
						return internal.NewValidationError("bill ID required")
					}
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("POST", c.InventoryBase+"/bills/"+cmd.Args().First()+"/status/void", &zohttp.RequestOpts{Params: orgParams(orgID)})
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
				Name:  "list",
				Usage: "List vendor credits",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("GET", c.InventoryBase+"/vendorcredits", &zohttp.RequestOpts{Params: orgParams(orgID)})
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
					if cmd.Args().Len() < 1 {
						return internal.NewValidationError("vendor credit ID required")
					}
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("GET", c.InventoryBase+"/vendorcredits/"+cmd.Args().First(), &zohttp.RequestOpts{Params: orgParams(orgID)})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:  "create",
				Usage: "Create a vendor credit",
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
					if err := json.Unmarshal([]byte(cmd.String("json")), &body); err != nil {
						return internal.NewValidationError(fmt.Sprintf("invalid JSON: %v", err))
					}
					raw, err := c.Request("POST", c.InventoryBase+"/vendorcredits", &zohttp.RequestOpts{
						Params: orgParams(orgID),
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
				Usage:     "Update a vendor credit",
				ArgsUsage: "<vendorcredit-id>",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "json", Required: true, Usage: "JSON body"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					if cmd.Args().Len() < 1 {
						return internal.NewValidationError("vendor credit ID required")
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
					if err := json.Unmarshal([]byte(cmd.String("json")), &body); err != nil {
						return internal.NewValidationError(fmt.Sprintf("invalid JSON: %v", err))
					}
					raw, err := c.Request("PUT", c.InventoryBase+"/vendorcredits/"+cmd.Args().First(), &zohttp.RequestOpts{
						Params: orgParams(orgID),
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
				Usage:     "Delete a vendor credit",
				ArgsUsage: "<vendorcredit-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					if cmd.Args().Len() < 1 {
						return internal.NewValidationError("vendor credit ID required")
					}
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("DELETE", c.InventoryBase+"/vendorcredits/"+cmd.Args().First(), &zohttp.RequestOpts{Params: orgParams(orgID)})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "mark-open",
				Usage:     "Convert a vendor credit to open",
				ArgsUsage: "<vendorcredit-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					if cmd.Args().Len() < 1 {
						return internal.NewValidationError("vendor credit ID required")
					}
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("POST", c.InventoryBase+"/vendorcredits/"+cmd.Args().First()+"/status/open", &zohttp.RequestOpts{Params: orgParams(orgID)})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "mark-void",
				Usage:     "Void a vendor credit",
				ArgsUsage: "<vendorcredit-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					if cmd.Args().Len() < 1 {
						return internal.NewValidationError("vendor credit ID required")
					}
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("POST", c.InventoryBase+"/vendorcredits/"+cmd.Args().First()+"/status/void", &zohttp.RequestOpts{Params: orgParams(orgID)})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
		},
	}
}

func priceListsCmd() *cli.Command {
	return &cli.Command{
		Name:  "price-lists",
		Usage: "Price list operations",
		Commands: []*cli.Command{
			{
				Name:  "list",
				Usage: "List price lists",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("GET", c.InventoryBase+"/pricebooks", &zohttp.RequestOpts{Params: orgParams(orgID)})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "get",
				Usage:     "Get a price list",
				ArgsUsage: "<pricebook-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					if cmd.Args().Len() < 1 {
						return internal.NewValidationError("price list ID required")
					}
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("GET", c.InventoryBase+"/pricebooks/"+cmd.Args().First(), &zohttp.RequestOpts{Params: orgParams(orgID)})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:  "create",
				Usage: "Create a price list",
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
					if err := json.Unmarshal([]byte(cmd.String("json")), &body); err != nil {
						return internal.NewValidationError(fmt.Sprintf("invalid JSON: %v", err))
					}
					raw, err := c.Request("POST", c.InventoryBase+"/pricebooks", &zohttp.RequestOpts{
						Params: orgParams(orgID),
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
				Usage:     "Update a price list",
				ArgsUsage: "<pricebook-id>",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "json", Required: true, Usage: "JSON body"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					if cmd.Args().Len() < 1 {
						return internal.NewValidationError("price list ID required")
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
					if err := json.Unmarshal([]byte(cmd.String("json")), &body); err != nil {
						return internal.NewValidationError(fmt.Sprintf("invalid JSON: %v", err))
					}
					raw, err := c.Request("PUT", c.InventoryBase+"/pricebooks/"+cmd.Args().First(), &zohttp.RequestOpts{
						Params: orgParams(orgID),
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
				Usage:     "Delete a price list",
				ArgsUsage: "<pricebook-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					if cmd.Args().Len() < 1 {
						return internal.NewValidationError("price list ID required")
					}
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("DELETE", c.InventoryBase+"/pricebooks/"+cmd.Args().First(), &zohttp.RequestOpts{Params: orgParams(orgID)})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "mark-active",
				Usage:     "Mark a price list as active",
				ArgsUsage: "<pricebook-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					if cmd.Args().Len() < 1 {
						return internal.NewValidationError("price list ID required")
					}
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("POST", c.InventoryBase+"/pricebooks/"+cmd.Args().First()+"/active", &zohttp.RequestOpts{Params: orgParams(orgID)})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "mark-inactive",
				Usage:     "Mark a price list as inactive",
				ArgsUsage: "<pricebook-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					if cmd.Args().Len() < 1 {
						return internal.NewValidationError("price list ID required")
					}
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("POST", c.InventoryBase+"/pricebooks/"+cmd.Args().First()+"/inactive", &zohttp.RequestOpts{Params: orgParams(orgID)})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
		},
	}
}

func warehousesCmd() *cli.Command {
	return &cli.Command{
		Name:  "warehouses",
		Usage: "Warehouse operations",
		Commands: []*cli.Command{
			{
				Name:  "list",
				Usage: "List warehouses",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("GET", c.InventoryBase+"/settings/warehouses", &zohttp.RequestOpts{Params: orgParams(orgID)})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "get",
				Usage:     "Get a warehouse",
				ArgsUsage: "<warehouse-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					if cmd.Args().Len() < 1 {
						return internal.NewValidationError("warehouse ID required")
					}
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("GET", c.InventoryBase+"/settings/warehouses/"+cmd.Args().First(), &zohttp.RequestOpts{Params: orgParams(orgID)})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:  "create",
				Usage: "Create a warehouse",
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
					if err := json.Unmarshal([]byte(cmd.String("json")), &body); err != nil {
						return internal.NewValidationError(fmt.Sprintf("invalid JSON: %v", err))
					}
					raw, err := c.Request("POST", c.InventoryBase+"/settings/warehouses", &zohttp.RequestOpts{
						Params: orgParams(orgID),
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
				Usage:     "Update a warehouse",
				ArgsUsage: "<warehouse-id>",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "json", Required: true, Usage: "JSON body"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					if cmd.Args().Len() < 1 {
						return internal.NewValidationError("warehouse ID required")
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
					if err := json.Unmarshal([]byte(cmd.String("json")), &body); err != nil {
						return internal.NewValidationError(fmt.Sprintf("invalid JSON: %v", err))
					}
					raw, err := c.Request("PUT", c.InventoryBase+"/settings/warehouses/"+cmd.Args().First(), &zohttp.RequestOpts{
						Params: orgParams(orgID),
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
				Usage:     "Delete a warehouse",
				ArgsUsage: "<warehouse-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					if cmd.Args().Len() < 1 {
						return internal.NewValidationError("warehouse ID required")
					}
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("DELETE", c.InventoryBase+"/settings/warehouses/"+cmd.Args().First(), &zohttp.RequestOpts{Params: orgParams(orgID)})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "mark-active",
				Usage:     "Mark a warehouse as active",
				ArgsUsage: "<warehouse-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					if cmd.Args().Len() < 1 {
						return internal.NewValidationError("warehouse ID required")
					}
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("POST", c.InventoryBase+"/settings/warehouses/"+cmd.Args().First()+"/active", &zohttp.RequestOpts{Params: orgParams(orgID)})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "mark-inactive",
				Usage:     "Mark a warehouse as inactive",
				ArgsUsage: "<warehouse-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					if cmd.Args().Len() < 1 {
						return internal.NewValidationError("warehouse ID required")
					}
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("POST", c.InventoryBase+"/settings/warehouses/"+cmd.Args().First()+"/inactive", &zohttp.RequestOpts{Params: orgParams(orgID)})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
		},
	}
}

func transferOrdersCmd() *cli.Command {
	return &cli.Command{
		Name:  "transfer-orders",
		Usage: "Transfer order operations",
		Commands: []*cli.Command{
			{
				Name:  "list",
				Usage: "List transfer orders",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("GET", c.InventoryBase+"/transferorders", &zohttp.RequestOpts{Params: orgParams(orgID)})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "get",
				Usage:     "Get a transfer order",
				ArgsUsage: "<transferorder-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					if cmd.Args().Len() < 1 {
						return internal.NewValidationError("transfer order ID required")
					}
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("GET", c.InventoryBase+"/transferorders/"+cmd.Args().First(), &zohttp.RequestOpts{Params: orgParams(orgID)})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:  "create",
				Usage: "Create a transfer order",
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
					if err := json.Unmarshal([]byte(cmd.String("json")), &body); err != nil {
						return internal.NewValidationError(fmt.Sprintf("invalid JSON: %v", err))
					}
					raw, err := c.Request("POST", c.InventoryBase+"/transferorders", &zohttp.RequestOpts{
						Params: orgParams(orgID),
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
				Usage:     "Update a transfer order",
				ArgsUsage: "<transferorder-id>",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "json", Required: true, Usage: "JSON body"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					if cmd.Args().Len() < 1 {
						return internal.NewValidationError("transfer order ID required")
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
					if err := json.Unmarshal([]byte(cmd.String("json")), &body); err != nil {
						return internal.NewValidationError(fmt.Sprintf("invalid JSON: %v", err))
					}
					raw, err := c.Request("PUT", c.InventoryBase+"/transferorders/"+cmd.Args().First(), &zohttp.RequestOpts{
						Params: orgParams(orgID),
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
				Usage:     "Delete a transfer order",
				ArgsUsage: "<transferorder-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					if cmd.Args().Len() < 1 {
						return internal.NewValidationError("transfer order ID required")
					}
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("DELETE", c.InventoryBase+"/transferorders/"+cmd.Args().First(), &zohttp.RequestOpts{Params: orgParams(orgID)})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "mark-received",
				Usage:     "Mark a transfer order as received",
				ArgsUsage: "<transferorder-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					if cmd.Args().Len() < 1 {
						return internal.NewValidationError("transfer order ID required")
					}
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("POST", c.InventoryBase+"/transferorders/"+cmd.Args().First()+"/received", &zohttp.RequestOpts{Params: orgParams(orgID)})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
		},
	}
}

func adjustmentsCmd() *cli.Command {
	return &cli.Command{
		Name:  "adjustments",
		Usage: "Inventory adjustment operations",
		Commands: []*cli.Command{
			{
				Name:  "list",
				Usage: "List inventory adjustments",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("GET", c.InventoryBase+"/inventoryadjustments", &zohttp.RequestOpts{Params: orgParams(orgID)})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "get",
				Usage:     "Get an inventory adjustment",
				ArgsUsage: "<adjustment-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					if cmd.Args().Len() < 1 {
						return internal.NewValidationError("adjustment ID required")
					}
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("GET", c.InventoryBase+"/inventoryadjustments/"+cmd.Args().First(), &zohttp.RequestOpts{Params: orgParams(orgID)})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:  "create",
				Usage: "Create an inventory adjustment",
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
					if err := json.Unmarshal([]byte(cmd.String("json")), &body); err != nil {
						return internal.NewValidationError(fmt.Sprintf("invalid JSON: %v", err))
					}
					raw, err := c.Request("POST", c.InventoryBase+"/inventoryadjustments", &zohttp.RequestOpts{
						Params: orgParams(orgID),
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
				Usage:     "Delete an inventory adjustment",
				ArgsUsage: "<adjustment-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					if cmd.Args().Len() < 1 {
						return internal.NewValidationError("adjustment ID required")
					}
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("DELETE", c.InventoryBase+"/inventoryadjustments/"+cmd.Args().First(), &zohttp.RequestOpts{Params: orgParams(orgID)})
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
				Usage: "List organizations",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					raw, err := c.Request("GET", c.InventoryBase+"/organizations", nil)
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
					c, err := getClient()
					if err != nil {
						return err
					}
					raw, err := c.Request("GET", c.InventoryBase+"/organizations/"+cmd.Args().First(), nil)
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
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("GET", c.InventoryBase+"/settings/currencies", &zohttp.RequestOpts{Params: orgParams(orgID)})
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
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("GET", c.InventoryBase+"/settings/currencies/"+cmd.Args().First(), &zohttp.RequestOpts{Params: orgParams(orgID)})
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
					if err := json.Unmarshal([]byte(cmd.String("json")), &body); err != nil {
						return internal.NewValidationError(fmt.Sprintf("invalid JSON: %v", err))
					}
					raw, err := c.Request("POST", c.InventoryBase+"/settings/currencies", &zohttp.RequestOpts{
						Params: orgParams(orgID),
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
					if err := json.Unmarshal([]byte(cmd.String("json")), &body); err != nil {
						return internal.NewValidationError(fmt.Sprintf("invalid JSON: %v", err))
					}
					raw, err := c.Request("PUT", c.InventoryBase+"/settings/currencies/"+cmd.Args().First(), &zohttp.RequestOpts{
						Params: orgParams(orgID),
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
					raw, err := c.Request("DELETE", c.InventoryBase+"/settings/currencies/"+cmd.Args().First(), &zohttp.RequestOpts{Params: orgParams(orgID)})
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
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("GET", c.InventoryBase+"/settings/taxes", &zohttp.RequestOpts{Params: orgParams(orgID)})
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
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("GET", c.InventoryBase+"/settings/taxes/"+cmd.Args().First(), &zohttp.RequestOpts{Params: orgParams(orgID)})
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
					if err := json.Unmarshal([]byte(cmd.String("json")), &body); err != nil {
						return internal.NewValidationError(fmt.Sprintf("invalid JSON: %v", err))
					}
					raw, err := c.Request("POST", c.InventoryBase+"/settings/taxes", &zohttp.RequestOpts{
						Params: orgParams(orgID),
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
					if err := json.Unmarshal([]byte(cmd.String("json")), &body); err != nil {
						return internal.NewValidationError(fmt.Sprintf("invalid JSON: %v", err))
					}
					raw, err := c.Request("PUT", c.InventoryBase+"/settings/taxes/"+cmd.Args().First(), &zohttp.RequestOpts{
						Params: orgParams(orgID),
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
					raw, err := c.Request("DELETE", c.InventoryBase+"/settings/taxes/"+cmd.Args().First(), &zohttp.RequestOpts{Params: orgParams(orgID)})
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
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("GET", c.InventoryBase+"/users", &zohttp.RequestOpts{Params: orgParams(orgID)})
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
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("GET", c.InventoryBase+"/users/"+cmd.Args().First(), &zohttp.RequestOpts{Params: orgParams(orgID)})
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
					raw, err := c.Request("GET", c.InventoryBase+"/users/me", &zohttp.RequestOpts{Params: orgParams(orgID)})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
		},
	}
}
