package inventory

import (
	"context"

	"github.com/omin8tor/zoho-cli/internal"
	zohttp "github.com/omin8tor/zoho-cli/internal/http"
	"github.com/omin8tor/zoho-cli/internal/output"
	"github.com/omin8tor/zoho-cli/internal/pagination"
	"github.com/urfave/cli/v3"
)

func orgParams(orgID string) map[string]string {
	return map[string]string{"organization_id": orgID}
}

func Commands() *cli.Command {
	return &cli.Command{
		Name:  "inventory",
		Usage: "Zoho Inventory operations",
		Flags: []cli.Flag{
			&cli.StringFlag{Name: "org", Sources: cli.EnvVars("ZOHO_BOOKS_ORG_ID"), Usage: "Organization ID (or set ZOHO_BOOKS_ORG_ID)"},
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
				Flags: []cli.Flag{
					&cli.BoolFlag{Name: "all", Usage: "Fetch all records"},
					&cli.IntFlag{Name: "limit", Usage: "Max total records to fetch"},
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := internal.RequireFlag(cmd, "org", "ZOHO_BOOKS_ORG_ID")
					if err != nil {
						return err
					}
					params := orgParams(orgID)
					if cmd.Bool("all") || cmd.IsSet("limit") {
						items, err := pagination.Paginate(ctx, pagination.PaginationConfig{
							Client:   c,
							URL:      c.InventoryBase+"/items",
							Opts:     &zohttp.RequestOpts{Params: params},
							ItemsKey: "items",
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
					raw, err := c.Request(ctx, "GET", c.InventoryBase+"/items", &zohttp.RequestOpts{Params: params})
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
				Action: func(ctx context.Context, cmd *cli.Command) error {
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
					raw, err := c.Request(ctx, "GET", c.InventoryBase+"/items/"+cmd.Args().First(), &zohttp.RequestOpts{Params: orgParams(orgID)})
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
					&cli.FloatFlag{Name: "rate", Usage: "Selling price"},
					&cli.FloatFlag{Name: "purchase_rate", Usage: "Purchase price"},
					&cli.StringFlag{Name: "sku", Usage: "Stock Keeping Unit"},
					&cli.StringFlag{Name: "unit", Usage: "Unit of measurement"},
					&cli.StringFlag{Name: "item_type", Usage: "Item type (inventory, sales, purchases, sales_and_purchases)"},
					&cli.StringFlag{Name: "product_type", Usage: "Product type (goods, service)"},
					&cli.StringFlag{Name: "description", Usage: "Item description"},
					&cli.StringFlag{Name: "json", Usage: "Additional fields as JSON"},
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
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
					if cmd.IsSet("rate") {
						body["rate"] = cmd.Float("rate")
					}
					if cmd.IsSet("purchase_rate") {
						body["purchase_rate"] = cmd.Float("purchase_rate")
					}
					if cmd.IsSet("sku") {
						body["sku"] = cmd.String("sku")
					}
					if cmd.IsSet("unit") {
						body["unit"] = cmd.String("unit")
					}
					if cmd.IsSet("item_type") {
						body["item_type"] = cmd.String("item_type")
					}
					if cmd.IsSet("product_type") {
						body["product_type"] = cmd.String("product_type")
					}
					if cmd.IsSet("description") {
						body["description"] = cmd.String("description")
					}
					if err := internal.MergeJSON(cmd, body); err != nil {
						return err
					}
					raw, err := c.Request(ctx, "POST", c.InventoryBase+"/items", &zohttp.RequestOpts{
						Params: orgParams(orgID),
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
					&cli.StringFlag{Name: "name", Usage: "Item name"},
					&cli.FloatFlag{Name: "rate", Usage: "Selling price"},
					&cli.FloatFlag{Name: "purchase_rate", Usage: "Purchase price"},
					&cli.StringFlag{Name: "sku", Usage: "Stock Keeping Unit"},
					&cli.StringFlag{Name: "unit", Usage: "Unit of measurement"},
					&cli.StringFlag{Name: "description", Usage: "Item description"},
					&cli.StringFlag{Name: "json", Usage: "Additional fields as JSON"},
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
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
					if cmd.IsSet("purchase_rate") {
						body["purchase_rate"] = cmd.Float("purchase_rate")
					}
					if cmd.IsSet("sku") {
						body["sku"] = cmd.String("sku")
					}
					if cmd.IsSet("unit") {
						body["unit"] = cmd.String("unit")
					}
					if cmd.IsSet("description") {
						body["description"] = cmd.String("description")
					}
					if err := internal.MergeJSON(cmd, body); err != nil {
						return err
					}
					raw, err := c.Request(ctx, "PUT", c.InventoryBase+"/items/"+cmd.Args().First(), &zohttp.RequestOpts{
						Params: orgParams(orgID),
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
				Action: func(ctx context.Context, cmd *cli.Command) error {
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
					raw, err := c.Request(ctx, "DELETE", c.InventoryBase+"/items/"+cmd.Args().First(), &zohttp.RequestOpts{Params: orgParams(orgID)})
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
				Action: func(ctx context.Context, cmd *cli.Command) error {
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
					raw, err := c.Request(ctx, "POST", c.InventoryBase+"/items/"+cmd.Args().First()+"/active", &zohttp.RequestOpts{Params: orgParams(orgID)})
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
				Action: func(ctx context.Context, cmd *cli.Command) error {
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
					raw, err := c.Request(ctx, "POST", c.InventoryBase+"/items/"+cmd.Args().First()+"/inactive", &zohttp.RequestOpts{Params: orgParams(orgID)})
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
				Flags: []cli.Flag{
					&cli.BoolFlag{Name: "all", Usage: "Fetch all records"},
					&cli.IntFlag{Name: "limit", Usage: "Max total records to fetch"},
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := internal.RequireFlag(cmd, "org", "ZOHO_BOOKS_ORG_ID")
					if err != nil {
						return err
					}
					params := orgParams(orgID)
					if cmd.Bool("all") || cmd.IsSet("limit") {
						items, err := pagination.Paginate(ctx, pagination.PaginationConfig{
							Client:   c,
							URL:      c.InventoryBase+"/compositeitems",
							Opts:     &zohttp.RequestOpts{Params: params},
							ItemsKey: "composite_items",
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
					raw, err := c.Request(ctx, "GET", c.InventoryBase+"/compositeitems", &zohttp.RequestOpts{Params: params})
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
				Action: func(ctx context.Context, cmd *cli.Command) error {
					if cmd.Args().Len() < 1 {
						return internal.NewValidationError("composite item ID required")
					}
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := internal.RequireFlag(cmd, "org", "ZOHO_BOOKS_ORG_ID")
					if err != nil {
						return err
					}
					raw, err := c.Request(ctx, "GET", c.InventoryBase+"/compositeitems/"+cmd.Args().First(), &zohttp.RequestOpts{Params: orgParams(orgID)})
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
					&cli.StringFlag{Name: "name", Required: true, Usage: "Composite item name"},
					&cli.FloatFlag{Name: "rate", Usage: "Selling price"},
					&cli.StringFlag{Name: "sku", Usage: "Stock Keeping Unit"},
					&cli.StringFlag{Name: "unit", Usage: "Unit of measurement"},
					&cli.StringFlag{Name: "product_type", Usage: "Product type (goods, service)"},
					&cli.StringFlag{Name: "description", Usage: "Item description"},
					&cli.StringFlag{Name: "json", Usage: "Additional fields as JSON"},
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
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
					if cmd.IsSet("rate") {
						body["rate"] = cmd.Float("rate")
					}
					if cmd.IsSet("sku") {
						body["sku"] = cmd.String("sku")
					}
					if cmd.IsSet("unit") {
						body["unit"] = cmd.String("unit")
					}
					if cmd.IsSet("product_type") {
						body["product_type"] = cmd.String("product_type")
					}
					if cmd.IsSet("description") {
						body["description"] = cmd.String("description")
					}
					if err := internal.MergeJSON(cmd, body); err != nil {
						return err
					}
					raw, err := c.Request(ctx, "POST", c.InventoryBase+"/compositeitems", &zohttp.RequestOpts{
						Params: orgParams(orgID),
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
					&cli.StringFlag{Name: "name", Usage: "Composite item name"},
					&cli.FloatFlag{Name: "rate", Usage: "Selling price"},
					&cli.StringFlag{Name: "sku", Usage: "Stock Keeping Unit"},
					&cli.StringFlag{Name: "unit", Usage: "Unit of measurement"},
					&cli.StringFlag{Name: "description", Usage: "Item description"},
					&cli.StringFlag{Name: "json", Usage: "Additional fields as JSON"},
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					if cmd.Args().Len() < 1 {
						return internal.NewValidationError("composite item ID required")
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
					if cmd.IsSet("sku") {
						body["sku"] = cmd.String("sku")
					}
					if cmd.IsSet("unit") {
						body["unit"] = cmd.String("unit")
					}
					if cmd.IsSet("description") {
						body["description"] = cmd.String("description")
					}
					if err := internal.MergeJSON(cmd, body); err != nil {
						return err
					}
					raw, err := c.Request(ctx, "PUT", c.InventoryBase+"/compositeitems/"+cmd.Args().First(), &zohttp.RequestOpts{
						Params: orgParams(orgID),
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
				Action: func(ctx context.Context, cmd *cli.Command) error {
					if cmd.Args().Len() < 1 {
						return internal.NewValidationError("composite item ID required")
					}
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := internal.RequireFlag(cmd, "org", "ZOHO_BOOKS_ORG_ID")
					if err != nil {
						return err
					}
					raw, err := c.Request(ctx, "DELETE", c.InventoryBase+"/compositeitems/"+cmd.Args().First(), &zohttp.RequestOpts{Params: orgParams(orgID)})
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
				Flags: []cli.Flag{
					&cli.BoolFlag{Name: "all", Usage: "Fetch all records"},
					&cli.IntFlag{Name: "limit", Usage: "Max total records to fetch"},
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := internal.RequireFlag(cmd, "org", "ZOHO_BOOKS_ORG_ID")
					if err != nil {
						return err
					}
					params := orgParams(orgID)
					if cmd.Bool("all") || cmd.IsSet("limit") {
						items, err := pagination.Paginate(ctx, pagination.PaginationConfig{
							Client:   c,
							URL:      c.InventoryBase+"/itemgroups",
							Opts:     &zohttp.RequestOpts{Params: params},
							ItemsKey: "itemgroups",
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
					raw, err := c.Request(ctx, "GET", c.InventoryBase+"/itemgroups", &zohttp.RequestOpts{Params: params})
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
				Action: func(ctx context.Context, cmd *cli.Command) error {
					if cmd.Args().Len() < 1 {
						return internal.NewValidationError("item group ID required")
					}
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := internal.RequireFlag(cmd, "org", "ZOHO_BOOKS_ORG_ID")
					if err != nil {
						return err
					}
					raw, err := c.Request(ctx, "GET", c.InventoryBase+"/itemgroups/"+cmd.Args().First(), &zohttp.RequestOpts{Params: orgParams(orgID)})
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
					&cli.StringFlag{Name: "group_name", Required: true, Usage: "Item group name"},
					&cli.StringFlag{Name: "unit", Required: true, Usage: "Unit of measurement"},
					&cli.StringFlag{Name: "description", Usage: "Group description"},
					&cli.StringFlag{Name: "brand", Usage: "Brand name"},
					&cli.StringFlag{Name: "manufacturer", Usage: "Manufacturer name"},
					&cli.StringFlag{Name: "json", Usage: "Additional fields as JSON"},
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := internal.RequireFlag(cmd, "org", "ZOHO_BOOKS_ORG_ID")
					if err != nil {
						return err
					}
					body := map[string]any{}
					body["group_name"] = cmd.String("group_name")
					body["unit"] = cmd.String("unit")
					if cmd.IsSet("description") {
						body["description"] = cmd.String("description")
					}
					if cmd.IsSet("brand") {
						body["brand"] = cmd.String("brand")
					}
					if cmd.IsSet("manufacturer") {
						body["manufacturer"] = cmd.String("manufacturer")
					}
					if err := internal.MergeJSON(cmd, body); err != nil {
						return err
					}
					raw, err := c.Request(ctx, "POST", c.InventoryBase+"/itemgroups", &zohttp.RequestOpts{
						Params: orgParams(orgID),
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
					&cli.StringFlag{Name: "group_name", Usage: "Item group name"},
					&cli.StringFlag{Name: "unit", Usage: "Unit of measurement"},
					&cli.StringFlag{Name: "description", Usage: "Group description"},
					&cli.StringFlag{Name: "brand", Usage: "Brand name"},
					&cli.StringFlag{Name: "manufacturer", Usage: "Manufacturer name"},
					&cli.StringFlag{Name: "json", Usage: "Additional fields as JSON"},
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					if cmd.Args().Len() < 1 {
						return internal.NewValidationError("item group ID required")
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
					if cmd.IsSet("group_name") {
						body["group_name"] = cmd.String("group_name")
					}
					if cmd.IsSet("unit") {
						body["unit"] = cmd.String("unit")
					}
					if cmd.IsSet("description") {
						body["description"] = cmd.String("description")
					}
					if cmd.IsSet("brand") {
						body["brand"] = cmd.String("brand")
					}
					if cmd.IsSet("manufacturer") {
						body["manufacturer"] = cmd.String("manufacturer")
					}
					if err := internal.MergeJSON(cmd, body); err != nil {
						return err
					}
					raw, err := c.Request(ctx, "PUT", c.InventoryBase+"/itemgroups/"+cmd.Args().First(), &zohttp.RequestOpts{
						Params: orgParams(orgID),
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
				Action: func(ctx context.Context, cmd *cli.Command) error {
					if cmd.Args().Len() < 1 {
						return internal.NewValidationError("item group ID required")
					}
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := internal.RequireFlag(cmd, "org", "ZOHO_BOOKS_ORG_ID")
					if err != nil {
						return err
					}
					raw, err := c.Request(ctx, "DELETE", c.InventoryBase+"/itemgroups/"+cmd.Args().First(), &zohttp.RequestOpts{Params: orgParams(orgID)})
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
				Flags: []cli.Flag{
					&cli.BoolFlag{Name: "all", Usage: "Fetch all records"},
					&cli.IntFlag{Name: "limit", Usage: "Max total records to fetch"},
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := internal.RequireFlag(cmd, "org", "ZOHO_BOOKS_ORG_ID")
					if err != nil {
						return err
					}
					params := orgParams(orgID)
					if cmd.Bool("all") || cmd.IsSet("limit") {
						items, err := pagination.Paginate(ctx, pagination.PaginationConfig{
							Client:   c,
							URL:      c.InventoryBase+"/contacts",
							Opts:     &zohttp.RequestOpts{Params: params},
							ItemsKey: "contacts",
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
					raw, err := c.Request(ctx, "GET", c.InventoryBase+"/contacts", &zohttp.RequestOpts{Params: params})
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
					raw, err := c.Request(ctx, "GET", c.InventoryBase+"/contacts/"+cmd.Args().First(), &zohttp.RequestOpts{Params: orgParams(orgID)})
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
					&cli.StringFlag{Name: "contact_name", Required: true, Usage: "Contact name"},
					&cli.StringFlag{Name: "contact_type", Usage: "Contact type (customer, vendor)"},
					&cli.StringFlag{Name: "company_name", Usage: "Company name"},
					&cli.StringFlag{Name: "website", Usage: "Website URL"},
					&cli.StringFlag{Name: "language_code", Usage: "Language code (en, de, es, fr, etc.)"},
					&cli.StringFlag{Name: "notes", Usage: "Notes about the contact"},
					&cli.StringFlag{Name: "json", Usage: "Additional fields as JSON"},
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
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
					if cmd.IsSet("contact_type") {
						body["contact_type"] = cmd.String("contact_type")
					}
					if cmd.IsSet("company_name") {
						body["company_name"] = cmd.String("company_name")
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
					raw, err := c.Request(ctx, "POST", c.InventoryBase+"/contacts", &zohttp.RequestOpts{
						Params: orgParams(orgID),
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
					&cli.StringFlag{Name: "contact_name", Usage: "Contact name"},
					&cli.StringFlag{Name: "contact_type", Usage: "Contact type (customer, vendor)"},
					&cli.StringFlag{Name: "company_name", Usage: "Company name"},
					&cli.StringFlag{Name: "website", Usage: "Website URL"},
					&cli.StringFlag{Name: "notes", Usage: "Notes about the contact"},
					&cli.StringFlag{Name: "json", Usage: "Additional fields as JSON"},
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
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
					if cmd.IsSet("contact_type") {
						body["contact_type"] = cmd.String("contact_type")
					}
					if cmd.IsSet("company_name") {
						body["company_name"] = cmd.String("company_name")
					}
					if cmd.IsSet("website") {
						body["website"] = cmd.String("website")
					}
					if cmd.IsSet("notes") {
						body["notes"] = cmd.String("notes")
					}
					if err := internal.MergeJSON(cmd, body); err != nil {
						return err
					}
					raw, err := c.Request(ctx, "PUT", c.InventoryBase+"/contacts/"+cmd.Args().First(), &zohttp.RequestOpts{
						Params: orgParams(orgID),
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
				Action: func(ctx context.Context, cmd *cli.Command) error {
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
					raw, err := c.Request(ctx, "DELETE", c.InventoryBase+"/contacts/"+cmd.Args().First(), &zohttp.RequestOpts{Params: orgParams(orgID)})
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
					raw, err := c.Request(ctx, "POST", c.InventoryBase+"/contacts/"+cmd.Args().First()+"/active", &zohttp.RequestOpts{Params: orgParams(orgID)})
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
					raw, err := c.Request(ctx, "POST", c.InventoryBase+"/contacts/"+cmd.Args().First()+"/inactive", &zohttp.RequestOpts{Params: orgParams(orgID)})
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
				Flags: []cli.Flag{
					&cli.BoolFlag{Name: "all", Usage: "Fetch all records"},
					&cli.IntFlag{Name: "limit", Usage: "Max total records to fetch"},
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := internal.RequireFlag(cmd, "org", "ZOHO_BOOKS_ORG_ID")
					if err != nil {
						return err
					}
					params := orgParams(orgID)
					if cmd.Bool("all") || cmd.IsSet("limit") {
						items, err := pagination.Paginate(ctx, pagination.PaginationConfig{
							Client:   c,
							URL:      c.InventoryBase+"/salesorders",
							Opts:     &zohttp.RequestOpts{Params: params},
							ItemsKey: "salesorders",
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
					raw, err := c.Request(ctx, "GET", c.InventoryBase+"/salesorders", &zohttp.RequestOpts{Params: params})
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
				Action: func(ctx context.Context, cmd *cli.Command) error {
					if cmd.Args().Len() < 1 {
						return internal.NewValidationError("sales order ID required")
					}
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := internal.RequireFlag(cmd, "org", "ZOHO_BOOKS_ORG_ID")
					if err != nil {
						return err
					}
					raw, err := c.Request(ctx, "GET", c.InventoryBase+"/salesorders/"+cmd.Args().First(), &zohttp.RequestOpts{Params: orgParams(orgID)})
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
					&cli.StringFlag{Name: "customer_id", Required: true, Usage: "Customer ID"},
					&cli.StringFlag{Name: "salesorder_number", Usage: "Sales order number"},
					&cli.StringFlag{Name: "date", Usage: "Sales order date (YYYY-MM-DD)"},
					&cli.StringFlag{Name: "shipment_date", Usage: "Shipment date (YYYY-MM-DD)"},
					&cli.StringFlag{Name: "reference_number", Usage: "Reference number"},
					&cli.StringFlag{Name: "notes", Usage: "Notes"},
					&cli.StringFlag{Name: "terms", Usage: "Terms and conditions"},
					&cli.StringFlag{Name: "json", Usage: "Additional fields as JSON"},
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
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
					if cmd.IsSet("salesorder_number") {
						body["salesorder_number"] = cmd.String("salesorder_number")
					}
					if cmd.IsSet("date") {
						body["date"] = cmd.String("date")
					}
					if cmd.IsSet("shipment_date") {
						body["shipment_date"] = cmd.String("shipment_date")
					}
					if cmd.IsSet("reference_number") {
						body["reference_number"] = cmd.String("reference_number")
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
					raw, err := c.Request(ctx, "POST", c.InventoryBase+"/salesorders", &zohttp.RequestOpts{
						Params: orgParams(orgID),
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
					&cli.StringFlag{Name: "customer_id", Usage: "Customer ID"},
					&cli.StringFlag{Name: "date", Usage: "Sales order date (YYYY-MM-DD)"},
					&cli.StringFlag{Name: "shipment_date", Usage: "Shipment date (YYYY-MM-DD)"},
					&cli.StringFlag{Name: "reference_number", Usage: "Reference number"},
					&cli.StringFlag{Name: "notes", Usage: "Notes"},
					&cli.StringFlag{Name: "terms", Usage: "Terms and conditions"},
					&cli.StringFlag{Name: "json", Usage: "Additional fields as JSON"},
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					if cmd.Args().Len() < 1 {
						return internal.NewValidationError("sales order ID required")
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
					if cmd.IsSet("date") {
						body["date"] = cmd.String("date")
					}
					if cmd.IsSet("shipment_date") {
						body["shipment_date"] = cmd.String("shipment_date")
					}
					if cmd.IsSet("reference_number") {
						body["reference_number"] = cmd.String("reference_number")
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
					raw, err := c.Request(ctx, "PUT", c.InventoryBase+"/salesorders/"+cmd.Args().First(), &zohttp.RequestOpts{
						Params: orgParams(orgID),
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
				Action: func(ctx context.Context, cmd *cli.Command) error {
					if cmd.Args().Len() < 1 {
						return internal.NewValidationError("sales order ID required")
					}
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := internal.RequireFlag(cmd, "org", "ZOHO_BOOKS_ORG_ID")
					if err != nil {
						return err
					}
					raw, err := c.Request(ctx, "DELETE", c.InventoryBase+"/salesorders/"+cmd.Args().First(), &zohttp.RequestOpts{Params: orgParams(orgID)})
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
				Action: func(ctx context.Context, cmd *cli.Command) error {
					if cmd.Args().Len() < 1 {
						return internal.NewValidationError("sales order ID required")
					}
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := internal.RequireFlag(cmd, "org", "ZOHO_BOOKS_ORG_ID")
					if err != nil {
						return err
					}
					raw, err := c.Request(ctx, "POST", c.InventoryBase+"/salesorders/"+cmd.Args().First()+"/status/confirmed", &zohttp.RequestOpts{Params: orgParams(orgID)})
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
				Action: func(ctx context.Context, cmd *cli.Command) error {
					if cmd.Args().Len() < 1 {
						return internal.NewValidationError("sales order ID required")
					}
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := internal.RequireFlag(cmd, "org", "ZOHO_BOOKS_ORG_ID")
					if err != nil {
						return err
					}
					raw, err := c.Request(ctx, "POST", c.InventoryBase+"/salesorders/"+cmd.Args().First()+"/status/void", &zohttp.RequestOpts{Params: orgParams(orgID)})
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
				Flags: []cli.Flag{
					&cli.BoolFlag{Name: "all", Usage: "Fetch all records"},
					&cli.IntFlag{Name: "limit", Usage: "Max total records to fetch"},
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := internal.RequireFlag(cmd, "org", "ZOHO_BOOKS_ORG_ID")
					if err != nil {
						return err
					}
					params := orgParams(orgID)
					if cmd.Bool("all") || cmd.IsSet("limit") {
						items, err := pagination.Paginate(ctx, pagination.PaginationConfig{
							Client:   c,
							URL:      c.InventoryBase+"/invoices",
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
					raw, err := c.Request(ctx, "GET", c.InventoryBase+"/invoices", &zohttp.RequestOpts{Params: params})
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
				Action: func(ctx context.Context, cmd *cli.Command) error {
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
					raw, err := c.Request(ctx, "GET", c.InventoryBase+"/invoices/"+cmd.Args().First(), &zohttp.RequestOpts{Params: orgParams(orgID)})
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
					&cli.StringFlag{Name: "date", Usage: "Invoice date (YYYY-MM-DD)"},
					&cli.StringFlag{Name: "due_date", Usage: "Due date (YYYY-MM-DD)"},
					&cli.StringFlag{Name: "reference_number", Usage: "Reference number"},
					&cli.StringFlag{Name: "notes", Usage: "Notes"},
					&cli.StringFlag{Name: "terms", Usage: "Terms and conditions"},
					&cli.StringFlag{Name: "json", Usage: "Additional fields as JSON"},
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
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
					if cmd.IsSet("date") {
						body["date"] = cmd.String("date")
					}
					if cmd.IsSet("due_date") {
						body["due_date"] = cmd.String("due_date")
					}
					if cmd.IsSet("reference_number") {
						body["reference_number"] = cmd.String("reference_number")
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
					raw, err := c.Request(ctx, "POST", c.InventoryBase+"/invoices", &zohttp.RequestOpts{
						Params: orgParams(orgID),
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
					&cli.StringFlag{Name: "date", Usage: "Invoice date (YYYY-MM-DD)"},
					&cli.StringFlag{Name: "due_date", Usage: "Due date (YYYY-MM-DD)"},
					&cli.StringFlag{Name: "reference_number", Usage: "Reference number"},
					&cli.StringFlag{Name: "notes", Usage: "Notes"},
					&cli.StringFlag{Name: "terms", Usage: "Terms and conditions"},
					&cli.StringFlag{Name: "json", Usage: "Additional fields as JSON"},
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
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
					if cmd.IsSet("date") {
						body["date"] = cmd.String("date")
					}
					if cmd.IsSet("due_date") {
						body["due_date"] = cmd.String("due_date")
					}
					if cmd.IsSet("reference_number") {
						body["reference_number"] = cmd.String("reference_number")
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
					raw, err := c.Request(ctx, "PUT", c.InventoryBase+"/invoices/"+cmd.Args().First(), &zohttp.RequestOpts{
						Params: orgParams(orgID),
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
				Action: func(ctx context.Context, cmd *cli.Command) error {
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
					raw, err := c.Request(ctx, "DELETE", c.InventoryBase+"/invoices/"+cmd.Args().First(), &zohttp.RequestOpts{Params: orgParams(orgID)})
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
				Action: func(ctx context.Context, cmd *cli.Command) error {
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
					raw, err := c.Request(ctx, "POST", c.InventoryBase+"/invoices/"+cmd.Args().First()+"/status/sent", &zohttp.RequestOpts{Params: orgParams(orgID)})
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
				Action: func(ctx context.Context, cmd *cli.Command) error {
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
					raw, err := c.Request(ctx, "POST", c.InventoryBase+"/invoices/"+cmd.Args().First()+"/status/void", &zohttp.RequestOpts{Params: orgParams(orgID)})
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
				Action: func(ctx context.Context, cmd *cli.Command) error {
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
					raw, err := c.Request(ctx, "POST", c.InventoryBase+"/invoices/"+cmd.Args().First()+"/status/draft", &zohttp.RequestOpts{Params: orgParams(orgID)})
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
				Flags: []cli.Flag{
					&cli.BoolFlag{Name: "all", Usage: "Fetch all records"},
					&cli.IntFlag{Name: "limit", Usage: "Max total records to fetch"},
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := internal.RequireFlag(cmd, "org", "ZOHO_BOOKS_ORG_ID")
					if err != nil {
						return err
					}
					params := orgParams(orgID)
					if cmd.Bool("all") || cmd.IsSet("limit") {
						items, err := pagination.Paginate(ctx, pagination.PaginationConfig{
							Client:   c,
							URL:      c.InventoryBase+"/packages",
							Opts:     &zohttp.RequestOpts{Params: params},
							ItemsKey: "packages",
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
					raw, err := c.Request(ctx, "GET", c.InventoryBase+"/packages", &zohttp.RequestOpts{Params: params})
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
				Action: func(ctx context.Context, cmd *cli.Command) error {
					if cmd.Args().Len() < 1 {
						return internal.NewValidationError("package ID required")
					}
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := internal.RequireFlag(cmd, "org", "ZOHO_BOOKS_ORG_ID")
					if err != nil {
						return err
					}
					raw, err := c.Request(ctx, "GET", c.InventoryBase+"/packages/"+cmd.Args().First(), &zohttp.RequestOpts{Params: orgParams(orgID)})
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
					&cli.StringFlag{Name: "package_number", Required: true, Usage: "Package number"},
					&cli.StringFlag{Name: "date", Required: true, Usage: "Package date (YYYY-MM-DD)"},
					&cli.StringFlag{Name: "salesorder_id", Required: true, Usage: "Sales order ID"},
					&cli.StringFlag{Name: "notes", Usage: "Package notes"},
					&cli.StringFlag{Name: "json", Usage: "Additional fields as JSON"},
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := internal.RequireFlag(cmd, "org", "ZOHO_BOOKS_ORG_ID")
					if err != nil {
						return err
					}
					body := map[string]any{}
					body["package_number"] = cmd.String("package_number")
					body["date"] = cmd.String("date")
					body["salesorder_id"] = cmd.String("salesorder_id")
					if cmd.IsSet("notes") {
						body["notes"] = cmd.String("notes")
					}
					if err := internal.MergeJSON(cmd, body); err != nil {
						return err
					}
					raw, err := c.Request(ctx, "POST", c.InventoryBase+"/packages", &zohttp.RequestOpts{
						Params: orgParams(orgID),
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
					&cli.StringFlag{Name: "package_number", Usage: "Package number"},
					&cli.StringFlag{Name: "date", Usage: "Package date (YYYY-MM-DD)"},
					&cli.StringFlag{Name: "notes", Usage: "Package notes"},
					&cli.StringFlag{Name: "json", Usage: "Additional fields as JSON"},
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					if cmd.Args().Len() < 1 {
						return internal.NewValidationError("package ID required")
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
					if cmd.IsSet("package_number") {
						body["package_number"] = cmd.String("package_number")
					}
					if cmd.IsSet("date") {
						body["date"] = cmd.String("date")
					}
					if cmd.IsSet("notes") {
						body["notes"] = cmd.String("notes")
					}
					if err := internal.MergeJSON(cmd, body); err != nil {
						return err
					}
					raw, err := c.Request(ctx, "PUT", c.InventoryBase+"/packages/"+cmd.Args().First(), &zohttp.RequestOpts{
						Params: orgParams(orgID),
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
				Action: func(ctx context.Context, cmd *cli.Command) error {
					if cmd.Args().Len() < 1 {
						return internal.NewValidationError("package ID required")
					}
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := internal.RequireFlag(cmd, "org", "ZOHO_BOOKS_ORG_ID")
					if err != nil {
						return err
					}
					raw, err := c.Request(ctx, "DELETE", c.InventoryBase+"/packages/"+cmd.Args().First(), &zohttp.RequestOpts{Params: orgParams(orgID)})
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
				Flags: []cli.Flag{
					&cli.BoolFlag{Name: "all", Usage: "Fetch all records"},
					&cli.IntFlag{Name: "limit", Usage: "Max total records to fetch"},
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := internal.RequireFlag(cmd, "org", "ZOHO_BOOKS_ORG_ID")
					if err != nil {
						return err
					}
					params := orgParams(orgID)
					if cmd.Bool("all") || cmd.IsSet("limit") {
						items, err := pagination.Paginate(ctx, pagination.PaginationConfig{
							Client:   c,
							URL:      c.InventoryBase+"/shipmentorders",
							Opts:     &zohttp.RequestOpts{Params: params},
							ItemsKey: "shipmentorders",
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
					raw, err := c.Request(ctx, "GET", c.InventoryBase+"/shipmentorders", &zohttp.RequestOpts{Params: params})
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
				Action: func(ctx context.Context, cmd *cli.Command) error {
					if cmd.Args().Len() < 1 {
						return internal.NewValidationError("shipment order ID required")
					}
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := internal.RequireFlag(cmd, "org", "ZOHO_BOOKS_ORG_ID")
					if err != nil {
						return err
					}
					raw, err := c.Request(ctx, "GET", c.InventoryBase+"/shipmentorders/"+cmd.Args().First(), &zohttp.RequestOpts{Params: orgParams(orgID)})
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
					&cli.StringFlag{Name: "shipment_number", Required: true, Usage: "Shipment number"},
					&cli.StringFlag{Name: "date", Required: true, Usage: "Shipment date (YYYY-MM-DD)"},
					&cli.StringFlag{Name: "delivery_method", Required: true, Usage: "Delivery method"},
					&cli.StringFlag{Name: "tracking_number", Required: true, Usage: "Tracking number"},
					&cli.FloatFlag{Name: "shipping_charge", Usage: "Shipping charges"},
					&cli.StringFlag{Name: "notes", Usage: "Notes"},
					&cli.StringFlag{Name: "json", Usage: "Additional fields as JSON"},
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := internal.RequireFlag(cmd, "org", "ZOHO_BOOKS_ORG_ID")
					if err != nil {
						return err
					}
					body := map[string]any{}
					body["shipment_number"] = cmd.String("shipment_number")
					body["date"] = cmd.String("date")
					body["delivery_method"] = cmd.String("delivery_method")
					body["tracking_number"] = cmd.String("tracking_number")
					if cmd.IsSet("shipping_charge") {
						body["shipping_charge"] = cmd.Float("shipping_charge")
					}
					if cmd.IsSet("notes") {
						body["notes"] = cmd.String("notes")
					}
					if err := internal.MergeJSON(cmd, body); err != nil {
						return err
					}
					raw, err := c.Request(ctx, "POST", c.InventoryBase+"/shipmentorders", &zohttp.RequestOpts{
						Params: orgParams(orgID),
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
				Action: func(ctx context.Context, cmd *cli.Command) error {
					if cmd.Args().Len() < 1 {
						return internal.NewValidationError("shipment order ID required")
					}
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := internal.RequireFlag(cmd, "org", "ZOHO_BOOKS_ORG_ID")
					if err != nil {
						return err
					}
					raw, err := c.Request(ctx, "DELETE", c.InventoryBase+"/shipmentorders/"+cmd.Args().First(), &zohttp.RequestOpts{Params: orgParams(orgID)})
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
				Flags: []cli.Flag{
					&cli.BoolFlag{Name: "all", Usage: "Fetch all records"},
					&cli.IntFlag{Name: "limit", Usage: "Max total records to fetch"},
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := internal.RequireFlag(cmd, "org", "ZOHO_BOOKS_ORG_ID")
					if err != nil {
						return err
					}
					params := orgParams(orgID)
					if cmd.Bool("all") || cmd.IsSet("limit") {
						items, err := pagination.Paginate(ctx, pagination.PaginationConfig{
							Client:   c,
							URL:      c.InventoryBase+"/purchaseorders",
							Opts:     &zohttp.RequestOpts{Params: params},
							ItemsKey: "purchaseorders",
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
					raw, err := c.Request(ctx, "GET", c.InventoryBase+"/purchaseorders", &zohttp.RequestOpts{Params: params})
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
				Action: func(ctx context.Context, cmd *cli.Command) error {
					if cmd.Args().Len() < 1 {
						return internal.NewValidationError("purchase order ID required")
					}
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := internal.RequireFlag(cmd, "org", "ZOHO_BOOKS_ORG_ID")
					if err != nil {
						return err
					}
					raw, err := c.Request(ctx, "GET", c.InventoryBase+"/purchaseorders/"+cmd.Args().First(), &zohttp.RequestOpts{Params: orgParams(orgID)})
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
					&cli.StringFlag{Name: "vendor_id", Required: true, Usage: "Vendor ID"},
					&cli.StringFlag{Name: "purchaseorder_number", Usage: "Purchase order number"},
					&cli.StringFlag{Name: "date", Usage: "Purchase order date (YYYY-MM-DD)"},
					&cli.StringFlag{Name: "delivery_date", Usage: "Delivery date (YYYY-MM-DD)"},
					&cli.StringFlag{Name: "reference_number", Usage: "Reference number"},
					&cli.StringFlag{Name: "notes", Usage: "Notes"},
					&cli.StringFlag{Name: "terms", Usage: "Terms and conditions"},
					&cli.StringFlag{Name: "json", Usage: "Additional fields as JSON"},
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := internal.RequireFlag(cmd, "org", "ZOHO_BOOKS_ORG_ID")
					if err != nil {
						return err
					}
					body := map[string]any{}
					body["vendor_id"] = cmd.String("vendor_id")
					if cmd.IsSet("purchaseorder_number") {
						body["purchaseorder_number"] = cmd.String("purchaseorder_number")
					}
					if cmd.IsSet("date") {
						body["date"] = cmd.String("date")
					}
					if cmd.IsSet("delivery_date") {
						body["delivery_date"] = cmd.String("delivery_date")
					}
					if cmd.IsSet("reference_number") {
						body["reference_number"] = cmd.String("reference_number")
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
					raw, err := c.Request(ctx, "POST", c.InventoryBase+"/purchaseorders", &zohttp.RequestOpts{
						Params: orgParams(orgID),
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
					&cli.StringFlag{Name: "vendor_id", Usage: "Vendor ID"},
					&cli.StringFlag{Name: "date", Usage: "Purchase order date (YYYY-MM-DD)"},
					&cli.StringFlag{Name: "delivery_date", Usage: "Delivery date (YYYY-MM-DD)"},
					&cli.StringFlag{Name: "reference_number", Usage: "Reference number"},
					&cli.StringFlag{Name: "notes", Usage: "Notes"},
					&cli.StringFlag{Name: "terms", Usage: "Terms and conditions"},
					&cli.StringFlag{Name: "json", Usage: "Additional fields as JSON"},
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					if cmd.Args().Len() < 1 {
						return internal.NewValidationError("purchase order ID required")
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
					if cmd.IsSet("vendor_id") {
						body["vendor_id"] = cmd.String("vendor_id")
					}
					if cmd.IsSet("date") {
						body["date"] = cmd.String("date")
					}
					if cmd.IsSet("delivery_date") {
						body["delivery_date"] = cmd.String("delivery_date")
					}
					if cmd.IsSet("reference_number") {
						body["reference_number"] = cmd.String("reference_number")
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
					raw, err := c.Request(ctx, "PUT", c.InventoryBase+"/purchaseorders/"+cmd.Args().First(), &zohttp.RequestOpts{
						Params: orgParams(orgID),
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
				Action: func(ctx context.Context, cmd *cli.Command) error {
					if cmd.Args().Len() < 1 {
						return internal.NewValidationError("purchase order ID required")
					}
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := internal.RequireFlag(cmd, "org", "ZOHO_BOOKS_ORG_ID")
					if err != nil {
						return err
					}
					raw, err := c.Request(ctx, "DELETE", c.InventoryBase+"/purchaseorders/"+cmd.Args().First(), &zohttp.RequestOpts{Params: orgParams(orgID)})
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
				Action: func(ctx context.Context, cmd *cli.Command) error {
					if cmd.Args().Len() < 1 {
						return internal.NewValidationError("purchase order ID required")
					}
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := internal.RequireFlag(cmd, "org", "ZOHO_BOOKS_ORG_ID")
					if err != nil {
						return err
					}
					raw, err := c.Request(ctx, "POST", c.InventoryBase+"/purchaseorders/"+cmd.Args().First()+"/status/issued", &zohttp.RequestOpts{Params: orgParams(orgID)})
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
				Action: func(ctx context.Context, cmd *cli.Command) error {
					if cmd.Args().Len() < 1 {
						return internal.NewValidationError("purchase order ID required")
					}
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := internal.RequireFlag(cmd, "org", "ZOHO_BOOKS_ORG_ID")
					if err != nil {
						return err
					}
					raw, err := c.Request(ctx, "POST", c.InventoryBase+"/purchaseorders/"+cmd.Args().First()+"/status/cancelled", &zohttp.RequestOpts{Params: orgParams(orgID)})
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
				Flags: []cli.Flag{
					&cli.BoolFlag{Name: "all", Usage: "Fetch all records"},
					&cli.IntFlag{Name: "limit", Usage: "Max total records to fetch"},
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := internal.RequireFlag(cmd, "org", "ZOHO_BOOKS_ORG_ID")
					if err != nil {
						return err
					}
					params := orgParams(orgID)
					if cmd.Bool("all") || cmd.IsSet("limit") {
						items, err := pagination.Paginate(ctx, pagination.PaginationConfig{
							Client:   c,
							URL:      c.InventoryBase+"/bills",
							Opts:     &zohttp.RequestOpts{Params: params},
							ItemsKey: "bills",
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
					raw, err := c.Request(ctx, "GET", c.InventoryBase+"/bills", &zohttp.RequestOpts{Params: params})
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
				Action: func(ctx context.Context, cmd *cli.Command) error {
					if cmd.Args().Len() < 1 {
						return internal.NewValidationError("bill ID required")
					}
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := internal.RequireFlag(cmd, "org", "ZOHO_BOOKS_ORG_ID")
					if err != nil {
						return err
					}
					raw, err := c.Request(ctx, "GET", c.InventoryBase+"/bills/"+cmd.Args().First(), &zohttp.RequestOpts{Params: orgParams(orgID)})
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
					&cli.StringFlag{Name: "vendor_id", Required: true, Usage: "Vendor ID"},
					&cli.StringFlag{Name: "bill_number", Usage: "Bill number"},
					&cli.StringFlag{Name: "date", Usage: "Bill date (YYYY-MM-DD)"},
					&cli.StringFlag{Name: "due_date", Usage: "Due date (YYYY-MM-DD)"},
					&cli.StringFlag{Name: "reference_number", Usage: "Reference number"},
					&cli.StringFlag{Name: "notes", Usage: "Notes"},
					&cli.StringFlag{Name: "terms", Usage: "Terms and conditions"},
					&cli.StringFlag{Name: "json", Usage: "Additional fields as JSON"},
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := internal.RequireFlag(cmd, "org", "ZOHO_BOOKS_ORG_ID")
					if err != nil {
						return err
					}
					body := map[string]any{}
					body["vendor_id"] = cmd.String("vendor_id")
					if cmd.IsSet("bill_number") {
						body["bill_number"] = cmd.String("bill_number")
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
					if cmd.IsSet("notes") {
						body["notes"] = cmd.String("notes")
					}
					if cmd.IsSet("terms") {
						body["terms"] = cmd.String("terms")
					}
					if err := internal.MergeJSON(cmd, body); err != nil {
						return err
					}
					raw, err := c.Request(ctx, "POST", c.InventoryBase+"/bills", &zohttp.RequestOpts{
						Params: orgParams(orgID),
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
					&cli.StringFlag{Name: "vendor_id", Usage: "Vendor ID"},
					&cli.StringFlag{Name: "bill_number", Usage: "Bill number"},
					&cli.StringFlag{Name: "date", Usage: "Bill date (YYYY-MM-DD)"},
					&cli.StringFlag{Name: "due_date", Usage: "Due date (YYYY-MM-DD)"},
					&cli.StringFlag{Name: "reference_number", Usage: "Reference number"},
					&cli.StringFlag{Name: "notes", Usage: "Notes"},
					&cli.StringFlag{Name: "terms", Usage: "Terms and conditions"},
					&cli.StringFlag{Name: "json", Usage: "Additional fields as JSON"},
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					if cmd.Args().Len() < 1 {
						return internal.NewValidationError("bill ID required")
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
					if cmd.IsSet("vendor_id") {
						body["vendor_id"] = cmd.String("vendor_id")
					}
					if cmd.IsSet("bill_number") {
						body["bill_number"] = cmd.String("bill_number")
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
					if cmd.IsSet("notes") {
						body["notes"] = cmd.String("notes")
					}
					if cmd.IsSet("terms") {
						body["terms"] = cmd.String("terms")
					}
					if err := internal.MergeJSON(cmd, body); err != nil {
						return err
					}
					raw, err := c.Request(ctx, "PUT", c.InventoryBase+"/bills/"+cmd.Args().First(), &zohttp.RequestOpts{
						Params: orgParams(orgID),
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
				Action: func(ctx context.Context, cmd *cli.Command) error {
					if cmd.Args().Len() < 1 {
						return internal.NewValidationError("bill ID required")
					}
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := internal.RequireFlag(cmd, "org", "ZOHO_BOOKS_ORG_ID")
					if err != nil {
						return err
					}
					raw, err := c.Request(ctx, "DELETE", c.InventoryBase+"/bills/"+cmd.Args().First(), &zohttp.RequestOpts{Params: orgParams(orgID)})
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
				Action: func(ctx context.Context, cmd *cli.Command) error {
					if cmd.Args().Len() < 1 {
						return internal.NewValidationError("bill ID required")
					}
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := internal.RequireFlag(cmd, "org", "ZOHO_BOOKS_ORG_ID")
					if err != nil {
						return err
					}
					raw, err := c.Request(ctx, "POST", c.InventoryBase+"/bills/"+cmd.Args().First()+"/status/open", &zohttp.RequestOpts{Params: orgParams(orgID)})
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
				Action: func(ctx context.Context, cmd *cli.Command) error {
					if cmd.Args().Len() < 1 {
						return internal.NewValidationError("bill ID required")
					}
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := internal.RequireFlag(cmd, "org", "ZOHO_BOOKS_ORG_ID")
					if err != nil {
						return err
					}
					raw, err := c.Request(ctx, "POST", c.InventoryBase+"/bills/"+cmd.Args().First()+"/status/void", &zohttp.RequestOpts{Params: orgParams(orgID)})
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
				Flags: []cli.Flag{
					&cli.BoolFlag{Name: "all", Usage: "Fetch all records"},
					&cli.IntFlag{Name: "limit", Usage: "Max total records to fetch"},
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := internal.RequireFlag(cmd, "org", "ZOHO_BOOKS_ORG_ID")
					if err != nil {
						return err
					}
					params := orgParams(orgID)
					if cmd.Bool("all") || cmd.IsSet("limit") {
						items, err := pagination.Paginate(ctx, pagination.PaginationConfig{
							Client:   c,
							URL:      c.InventoryBase+"/vendorcredits",
							Opts:     &zohttp.RequestOpts{Params: params},
							ItemsKey: "vendorcredits",
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
					raw, err := c.Request(ctx, "GET", c.InventoryBase+"/vendorcredits", &zohttp.RequestOpts{Params: params})
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
				Action: func(ctx context.Context, cmd *cli.Command) error {
					if cmd.Args().Len() < 1 {
						return internal.NewValidationError("vendor credit ID required")
					}
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := internal.RequireFlag(cmd, "org", "ZOHO_BOOKS_ORG_ID")
					if err != nil {
						return err
					}
					raw, err := c.Request(ctx, "GET", c.InventoryBase+"/vendorcredits/"+cmd.Args().First(), &zohttp.RequestOpts{Params: orgParams(orgID)})
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
					&cli.StringFlag{Name: "vendor_id", Required: true, Usage: "Vendor ID"},
					&cli.StringFlag{Name: "vendor_credit_number", Usage: "Vendor credit number"},
					&cli.StringFlag{Name: "date", Usage: "Vendor credit date (YYYY-MM-DD)"},
					&cli.StringFlag{Name: "reference_number", Usage: "Reference number"},
					&cli.StringFlag{Name: "notes", Usage: "Notes"},
					&cli.StringFlag{Name: "json", Usage: "Additional fields as JSON"},
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := internal.RequireFlag(cmd, "org", "ZOHO_BOOKS_ORG_ID")
					if err != nil {
						return err
					}
					body := map[string]any{}
					body["vendor_id"] = cmd.String("vendor_id")
					if cmd.IsSet("vendor_credit_number") {
						body["vendor_credit_number"] = cmd.String("vendor_credit_number")
					}
					if cmd.IsSet("date") {
						body["date"] = cmd.String("date")
					}
					if cmd.IsSet("reference_number") {
						body["reference_number"] = cmd.String("reference_number")
					}
					if cmd.IsSet("notes") {
						body["notes"] = cmd.String("notes")
					}
					if err := internal.MergeJSON(cmd, body); err != nil {
						return err
					}
					raw, err := c.Request(ctx, "POST", c.InventoryBase+"/vendorcredits", &zohttp.RequestOpts{
						Params: orgParams(orgID),
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
					&cli.StringFlag{Name: "vendor_id", Usage: "Vendor ID"},
					&cli.StringFlag{Name: "vendor_credit_number", Usage: "Vendor credit number"},
					&cli.StringFlag{Name: "date", Usage: "Vendor credit date (YYYY-MM-DD)"},
					&cli.StringFlag{Name: "reference_number", Usage: "Reference number"},
					&cli.StringFlag{Name: "notes", Usage: "Notes"},
					&cli.StringFlag{Name: "json", Usage: "Additional fields as JSON"},
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					if cmd.Args().Len() < 1 {
						return internal.NewValidationError("vendor credit ID required")
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
					if cmd.IsSet("vendor_id") {
						body["vendor_id"] = cmd.String("vendor_id")
					}
					if cmd.IsSet("vendor_credit_number") {
						body["vendor_credit_number"] = cmd.String("vendor_credit_number")
					}
					if cmd.IsSet("date") {
						body["date"] = cmd.String("date")
					}
					if cmd.IsSet("reference_number") {
						body["reference_number"] = cmd.String("reference_number")
					}
					if cmd.IsSet("notes") {
						body["notes"] = cmd.String("notes")
					}
					if err := internal.MergeJSON(cmd, body); err != nil {
						return err
					}
					raw, err := c.Request(ctx, "PUT", c.InventoryBase+"/vendorcredits/"+cmd.Args().First(), &zohttp.RequestOpts{
						Params: orgParams(orgID),
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
				Action: func(ctx context.Context, cmd *cli.Command) error {
					if cmd.Args().Len() < 1 {
						return internal.NewValidationError("vendor credit ID required")
					}
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := internal.RequireFlag(cmd, "org", "ZOHO_BOOKS_ORG_ID")
					if err != nil {
						return err
					}
					raw, err := c.Request(ctx, "DELETE", c.InventoryBase+"/vendorcredits/"+cmd.Args().First(), &zohttp.RequestOpts{Params: orgParams(orgID)})
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
				Action: func(ctx context.Context, cmd *cli.Command) error {
					if cmd.Args().Len() < 1 {
						return internal.NewValidationError("vendor credit ID required")
					}
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := internal.RequireFlag(cmd, "org", "ZOHO_BOOKS_ORG_ID")
					if err != nil {
						return err
					}
					raw, err := c.Request(ctx, "POST", c.InventoryBase+"/vendorcredits/"+cmd.Args().First()+"/status/open", &zohttp.RequestOpts{Params: orgParams(orgID)})
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
				Action: func(ctx context.Context, cmd *cli.Command) error {
					if cmd.Args().Len() < 1 {
						return internal.NewValidationError("vendor credit ID required")
					}
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := internal.RequireFlag(cmd, "org", "ZOHO_BOOKS_ORG_ID")
					if err != nil {
						return err
					}
					raw, err := c.Request(ctx, "POST", c.InventoryBase+"/vendorcredits/"+cmd.Args().First()+"/status/void", &zohttp.RequestOpts{Params: orgParams(orgID)})
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
				Flags: []cli.Flag{
					&cli.BoolFlag{Name: "all", Usage: "Fetch all records"},
					&cli.IntFlag{Name: "limit", Usage: "Max total records to fetch"},
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := internal.RequireFlag(cmd, "org", "ZOHO_BOOKS_ORG_ID")
					if err != nil {
						return err
					}
					params := orgParams(orgID)
					if cmd.Bool("all") || cmd.IsSet("limit") {
						items, err := pagination.Paginate(ctx, pagination.PaginationConfig{
							Client:   c,
							URL:      c.InventoryBase+"/pricebooks",
							Opts:     &zohttp.RequestOpts{Params: params},
							ItemsKey: "price_lists",
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
					raw, err := c.Request(ctx, "GET", c.InventoryBase+"/pricebooks", &zohttp.RequestOpts{Params: params})
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
				Action: func(ctx context.Context, cmd *cli.Command) error {
					if cmd.Args().Len() < 1 {
						return internal.NewValidationError("price list ID required")
					}
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := internal.RequireFlag(cmd, "org", "ZOHO_BOOKS_ORG_ID")
					if err != nil {
						return err
					}
					raw, err := c.Request(ctx, "GET", c.InventoryBase+"/pricebooks/"+cmd.Args().First(), &zohttp.RequestOpts{Params: orgParams(orgID)})
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
					&cli.StringFlag{Name: "name", Required: true, Usage: "Price list name"},
					&cli.StringFlag{Name: "description", Usage: "Description"},
					&cli.StringFlag{Name: "currency_id", Usage: "Currency ID"},
					&cli.StringFlag{Name: "type", Usage: "Pricebook type"},
					&cli.StringFlag{Name: "json", Usage: "Additional fields as JSON"},
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
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
					if cmd.IsSet("currency_id") {
						body["currency_id"] = cmd.String("currency_id")
					}
					if cmd.IsSet("type") {
						body["type"] = cmd.String("type")
					}
					if err := internal.MergeJSON(cmd, body); err != nil {
						return err
					}
					raw, err := c.Request(ctx, "POST", c.InventoryBase+"/pricebooks", &zohttp.RequestOpts{
						Params: orgParams(orgID),
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
					&cli.StringFlag{Name: "name", Usage: "Price list name"},
					&cli.StringFlag{Name: "description", Usage: "Description"},
					&cli.StringFlag{Name: "currency_id", Usage: "Currency ID"},
					&cli.StringFlag{Name: "json", Usage: "Additional fields as JSON"},
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					if cmd.Args().Len() < 1 {
						return internal.NewValidationError("price list ID required")
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
					if cmd.IsSet("currency_id") {
						body["currency_id"] = cmd.String("currency_id")
					}
					if err := internal.MergeJSON(cmd, body); err != nil {
						return err
					}
					raw, err := c.Request(ctx, "PUT", c.InventoryBase+"/pricebooks/"+cmd.Args().First(), &zohttp.RequestOpts{
						Params: orgParams(orgID),
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
				Action: func(ctx context.Context, cmd *cli.Command) error {
					if cmd.Args().Len() < 1 {
						return internal.NewValidationError("price list ID required")
					}
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := internal.RequireFlag(cmd, "org", "ZOHO_BOOKS_ORG_ID")
					if err != nil {
						return err
					}
					raw, err := c.Request(ctx, "DELETE", c.InventoryBase+"/pricebooks/"+cmd.Args().First(), &zohttp.RequestOpts{Params: orgParams(orgID)})
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
				Action: func(ctx context.Context, cmd *cli.Command) error {
					if cmd.Args().Len() < 1 {
						return internal.NewValidationError("price list ID required")
					}
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := internal.RequireFlag(cmd, "org", "ZOHO_BOOKS_ORG_ID")
					if err != nil {
						return err
					}
					raw, err := c.Request(ctx, "POST", c.InventoryBase+"/pricebooks/"+cmd.Args().First()+"/active", &zohttp.RequestOpts{Params: orgParams(orgID)})
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
				Action: func(ctx context.Context, cmd *cli.Command) error {
					if cmd.Args().Len() < 1 {
						return internal.NewValidationError("price list ID required")
					}
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := internal.RequireFlag(cmd, "org", "ZOHO_BOOKS_ORG_ID")
					if err != nil {
						return err
					}
					raw, err := c.Request(ctx, "POST", c.InventoryBase+"/pricebooks/"+cmd.Args().First()+"/inactive", &zohttp.RequestOpts{Params: orgParams(orgID)})
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
				Flags: []cli.Flag{
					&cli.BoolFlag{Name: "all", Usage: "Fetch all records"},
					&cli.IntFlag{Name: "limit", Usage: "Max total records to fetch"},
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := internal.RequireFlag(cmd, "org", "ZOHO_BOOKS_ORG_ID")
					if err != nil {
						return err
					}
					params := orgParams(orgID)
					if cmd.Bool("all") || cmd.IsSet("limit") {
						items, err := pagination.Paginate(ctx, pagination.PaginationConfig{
							Client:   c,
							URL:      c.InventoryBase+"/settings/warehouses",
							Opts:     &zohttp.RequestOpts{Params: params},
							ItemsKey: "warehouses",
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
					raw, err := c.Request(ctx, "GET", c.InventoryBase+"/settings/warehouses", &zohttp.RequestOpts{Params: params})
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
				Action: func(ctx context.Context, cmd *cli.Command) error {
					if cmd.Args().Len() < 1 {
						return internal.NewValidationError("warehouse ID required")
					}
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := internal.RequireFlag(cmd, "org", "ZOHO_BOOKS_ORG_ID")
					if err != nil {
						return err
					}
					raw, err := c.Request(ctx, "GET", c.InventoryBase+"/settings/warehouses/"+cmd.Args().First(), &zohttp.RequestOpts{Params: orgParams(orgID)})
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
					&cli.StringFlag{Name: "warehouse_name", Required: true, Usage: "Warehouse name"},
					&cli.StringFlag{Name: "address", Usage: "Street address"},
					&cli.StringFlag{Name: "city", Usage: "City"},
					&cli.StringFlag{Name: "state", Usage: "State"},
					&cli.StringFlag{Name: "country", Usage: "Country"},
					&cli.StringFlag{Name: "zip", Usage: "ZIP/postal code"},
					&cli.StringFlag{Name: "json", Usage: "Additional fields as JSON"},
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := internal.RequireFlag(cmd, "org", "ZOHO_BOOKS_ORG_ID")
					if err != nil {
						return err
					}
					body := map[string]any{}
					body["warehouse_name"] = cmd.String("warehouse_name")
					if cmd.IsSet("address") {
						body["address"] = cmd.String("address")
					}
					if cmd.IsSet("city") {
						body["city"] = cmd.String("city")
					}
					if cmd.IsSet("state") {
						body["state"] = cmd.String("state")
					}
					if cmd.IsSet("country") {
						body["country"] = cmd.String("country")
					}
					if cmd.IsSet("zip") {
						body["zip"] = cmd.String("zip")
					}
					if err := internal.MergeJSON(cmd, body); err != nil {
						return err
					}
					raw, err := c.Request(ctx, "POST", c.InventoryBase+"/settings/warehouses", &zohttp.RequestOpts{
						Params: orgParams(orgID),
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
					&cli.StringFlag{Name: "warehouse_name", Usage: "Warehouse name"},
					&cli.StringFlag{Name: "address", Usage: "Street address"},
					&cli.StringFlag{Name: "city", Usage: "City"},
					&cli.StringFlag{Name: "state", Usage: "State"},
					&cli.StringFlag{Name: "country", Usage: "Country"},
					&cli.StringFlag{Name: "zip", Usage: "ZIP/postal code"},
					&cli.StringFlag{Name: "json", Usage: "Additional fields as JSON"},
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					if cmd.Args().Len() < 1 {
						return internal.NewValidationError("warehouse ID required")
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
					if cmd.IsSet("warehouse_name") {
						body["warehouse_name"] = cmd.String("warehouse_name")
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
					if cmd.IsSet("country") {
						body["country"] = cmd.String("country")
					}
					if cmd.IsSet("zip") {
						body["zip"] = cmd.String("zip")
					}
					if err := internal.MergeJSON(cmd, body); err != nil {
						return err
					}
					raw, err := c.Request(ctx, "PUT", c.InventoryBase+"/settings/warehouses/"+cmd.Args().First(), &zohttp.RequestOpts{
						Params: orgParams(orgID),
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
				Action: func(ctx context.Context, cmd *cli.Command) error {
					if cmd.Args().Len() < 1 {
						return internal.NewValidationError("warehouse ID required")
					}
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := internal.RequireFlag(cmd, "org", "ZOHO_BOOKS_ORG_ID")
					if err != nil {
						return err
					}
					raw, err := c.Request(ctx, "DELETE", c.InventoryBase+"/settings/warehouses/"+cmd.Args().First(), &zohttp.RequestOpts{Params: orgParams(orgID)})
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
				Action: func(ctx context.Context, cmd *cli.Command) error {
					if cmd.Args().Len() < 1 {
						return internal.NewValidationError("warehouse ID required")
					}
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := internal.RequireFlag(cmd, "org", "ZOHO_BOOKS_ORG_ID")
					if err != nil {
						return err
					}
					raw, err := c.Request(ctx, "POST", c.InventoryBase+"/settings/warehouses/"+cmd.Args().First()+"/active", &zohttp.RequestOpts{Params: orgParams(orgID)})
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
				Action: func(ctx context.Context, cmd *cli.Command) error {
					if cmd.Args().Len() < 1 {
						return internal.NewValidationError("warehouse ID required")
					}
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := internal.RequireFlag(cmd, "org", "ZOHO_BOOKS_ORG_ID")
					if err != nil {
						return err
					}
					raw, err := c.Request(ctx, "POST", c.InventoryBase+"/settings/warehouses/"+cmd.Args().First()+"/inactive", &zohttp.RequestOpts{Params: orgParams(orgID)})
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
				Flags: []cli.Flag{
					&cli.BoolFlag{Name: "all", Usage: "Fetch all records"},
					&cli.IntFlag{Name: "limit", Usage: "Max total records to fetch"},
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := internal.RequireFlag(cmd, "org", "ZOHO_BOOKS_ORG_ID")
					if err != nil {
						return err
					}
					params := orgParams(orgID)
					if cmd.Bool("all") || cmd.IsSet("limit") {
						items, err := pagination.Paginate(ctx, pagination.PaginationConfig{
							Client:   c,
							URL:      c.InventoryBase+"/transferorders",
							Opts:     &zohttp.RequestOpts{Params: params},
							ItemsKey: "transferorders",
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
					raw, err := c.Request(ctx, "GET", c.InventoryBase+"/transferorders", &zohttp.RequestOpts{Params: params})
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
				Action: func(ctx context.Context, cmd *cli.Command) error {
					if cmd.Args().Len() < 1 {
						return internal.NewValidationError("transfer order ID required")
					}
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := internal.RequireFlag(cmd, "org", "ZOHO_BOOKS_ORG_ID")
					if err != nil {
						return err
					}
					raw, err := c.Request(ctx, "GET", c.InventoryBase+"/transferorders/"+cmd.Args().First(), &zohttp.RequestOpts{Params: orgParams(orgID)})
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
					&cli.StringFlag{Name: "from_warehouse_id", Required: true, Usage: "Source warehouse ID"},
					&cli.StringFlag{Name: "to_warehouse_id", Required: true, Usage: "Destination warehouse ID"},
					&cli.StringFlag{Name: "transfer_order_number", Usage: "Transfer order number"},
					&cli.StringFlag{Name: "date", Usage: "Transfer order date (YYYY-MM-DD)"},
					&cli.StringFlag{Name: "notes", Usage: "Notes"},
					&cli.StringFlag{Name: "json", Usage: "Additional fields as JSON"},
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := internal.RequireFlag(cmd, "org", "ZOHO_BOOKS_ORG_ID")
					if err != nil {
						return err
					}
					body := map[string]any{}
					body["from_warehouse_id"] = cmd.String("from_warehouse_id")
					body["to_warehouse_id"] = cmd.String("to_warehouse_id")
					if cmd.IsSet("transfer_order_number") {
						body["transfer_order_number"] = cmd.String("transfer_order_number")
					}
					if cmd.IsSet("date") {
						body["date"] = cmd.String("date")
					}
					if cmd.IsSet("notes") {
						body["notes"] = cmd.String("notes")
					}
					if err := internal.MergeJSON(cmd, body); err != nil {
						return err
					}
					raw, err := c.Request(ctx, "POST", c.InventoryBase+"/transferorders", &zohttp.RequestOpts{
						Params: orgParams(orgID),
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
					&cli.StringFlag{Name: "from_warehouse_id", Usage: "Source warehouse ID"},
					&cli.StringFlag{Name: "to_warehouse_id", Usage: "Destination warehouse ID"},
					&cli.StringFlag{Name: "date", Usage: "Transfer order date (YYYY-MM-DD)"},
					&cli.StringFlag{Name: "notes", Usage: "Notes"},
					&cli.StringFlag{Name: "json", Usage: "Additional fields as JSON"},
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					if cmd.Args().Len() < 1 {
						return internal.NewValidationError("transfer order ID required")
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
					if cmd.IsSet("from_warehouse_id") {
						body["from_warehouse_id"] = cmd.String("from_warehouse_id")
					}
					if cmd.IsSet("to_warehouse_id") {
						body["to_warehouse_id"] = cmd.String("to_warehouse_id")
					}
					if cmd.IsSet("date") {
						body["date"] = cmd.String("date")
					}
					if cmd.IsSet("notes") {
						body["notes"] = cmd.String("notes")
					}
					if err := internal.MergeJSON(cmd, body); err != nil {
						return err
					}
					raw, err := c.Request(ctx, "PUT", c.InventoryBase+"/transferorders/"+cmd.Args().First(), &zohttp.RequestOpts{
						Params: orgParams(orgID),
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
				Action: func(ctx context.Context, cmd *cli.Command) error {
					if cmd.Args().Len() < 1 {
						return internal.NewValidationError("transfer order ID required")
					}
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := internal.RequireFlag(cmd, "org", "ZOHO_BOOKS_ORG_ID")
					if err != nil {
						return err
					}
					raw, err := c.Request(ctx, "DELETE", c.InventoryBase+"/transferorders/"+cmd.Args().First(), &zohttp.RequestOpts{Params: orgParams(orgID)})
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
				Action: func(ctx context.Context, cmd *cli.Command) error {
					if cmd.Args().Len() < 1 {
						return internal.NewValidationError("transfer order ID required")
					}
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := internal.RequireFlag(cmd, "org", "ZOHO_BOOKS_ORG_ID")
					if err != nil {
						return err
					}
					raw, err := c.Request(ctx, "POST", c.InventoryBase+"/transferorders/"+cmd.Args().First()+"/received", &zohttp.RequestOpts{Params: orgParams(orgID)})
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
				Flags: []cli.Flag{
					&cli.BoolFlag{Name: "all", Usage: "Fetch all records"},
					&cli.IntFlag{Name: "limit", Usage: "Max total records to fetch"},
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := internal.RequireFlag(cmd, "org", "ZOHO_BOOKS_ORG_ID")
					if err != nil {
						return err
					}
					params := orgParams(orgID)
					if cmd.Bool("all") || cmd.IsSet("limit") {
						items, err := pagination.Paginate(ctx, pagination.PaginationConfig{
							Client:   c,
							URL:      c.InventoryBase+"/inventoryadjustments",
							Opts:     &zohttp.RequestOpts{Params: params},
							ItemsKey: "inventory_adjustments",
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
					raw, err := c.Request(ctx, "GET", c.InventoryBase+"/inventoryadjustments", &zohttp.RequestOpts{Params: params})
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
				Action: func(ctx context.Context, cmd *cli.Command) error {
					if cmd.Args().Len() < 1 {
						return internal.NewValidationError("adjustment ID required")
					}
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := internal.RequireFlag(cmd, "org", "ZOHO_BOOKS_ORG_ID")
					if err != nil {
						return err
					}
					raw, err := c.Request(ctx, "GET", c.InventoryBase+"/inventoryadjustments/"+cmd.Args().First(), &zohttp.RequestOpts{Params: orgParams(orgID)})
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
					&cli.StringFlag{Name: "date", Required: true, Usage: "Adjustment date (YYYY-MM-DD)"},
					&cli.StringFlag{Name: "reason", Required: true, Usage: "Reason for adjustment"},
					&cli.StringFlag{Name: "adjustment_type", Usage: "Adjustment type (quantity, value)"},
					&cli.StringFlag{Name: "description", Usage: "Description"},
					&cli.StringFlag{Name: "reference_number", Usage: "Reference number"},
					&cli.StringFlag{Name: "json", Usage: "Additional fields as JSON"},
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := internal.RequireFlag(cmd, "org", "ZOHO_BOOKS_ORG_ID")
					if err != nil {
						return err
					}
					body := map[string]any{}
					body["date"] = cmd.String("date")
					body["reason"] = cmd.String("reason")
					if cmd.IsSet("adjustment_type") {
						body["adjustment_type"] = cmd.String("adjustment_type")
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
					raw, err := c.Request(ctx, "POST", c.InventoryBase+"/inventoryadjustments", &zohttp.RequestOpts{
						Params: orgParams(orgID),
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
				Action: func(ctx context.Context, cmd *cli.Command) error {
					if cmd.Args().Len() < 1 {
						return internal.NewValidationError("adjustment ID required")
					}
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := internal.RequireFlag(cmd, "org", "ZOHO_BOOKS_ORG_ID")
					if err != nil {
						return err
					}
					raw, err := c.Request(ctx, "DELETE", c.InventoryBase+"/inventoryadjustments/"+cmd.Args().First(), &zohttp.RequestOpts{Params: orgParams(orgID)})
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
				Action: func(ctx context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					raw, err := c.Request(ctx, "GET", c.InventoryBase+"/organizations", nil)
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
				Action: func(ctx context.Context, cmd *cli.Command) error {
					if cmd.Args().Len() < 1 {
						return internal.NewValidationError("organization ID required")
					}
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					raw, err := c.Request(ctx, "GET", c.InventoryBase+"/organizations/"+cmd.Args().First(), nil)
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
				Action: func(ctx context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := internal.RequireFlag(cmd, "org", "ZOHO_BOOKS_ORG_ID")
					if err != nil {
						return err
					}
					raw, err := c.Request(ctx, "GET", c.InventoryBase+"/settings/currencies", &zohttp.RequestOpts{Params: orgParams(orgID)})
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
				Action: func(ctx context.Context, cmd *cli.Command) error {
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
					raw, err := c.Request(ctx, "GET", c.InventoryBase+"/settings/currencies/"+cmd.Args().First(), &zohttp.RequestOpts{Params: orgParams(orgID)})
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
					&cli.StringFlag{Name: "currency_symbol", Usage: "Currency symbol"},
					&cli.StringFlag{Name: "currency_format", Usage: "Currency format"},
					&cli.StringFlag{Name: "json", Usage: "Additional fields as JSON"},
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
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
					if cmd.IsSet("currency_symbol") {
						body["currency_symbol"] = cmd.String("currency_symbol")
					}
					if cmd.IsSet("currency_format") {
						body["currency_format"] = cmd.String("currency_format")
					}
					if err := internal.MergeJSON(cmd, body); err != nil {
						return err
					}
					raw, err := c.Request(ctx, "POST", c.InventoryBase+"/settings/currencies", &zohttp.RequestOpts{
						Params: orgParams(orgID),
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
					&cli.StringFlag{Name: "currency_code", Usage: "Currency code"},
					&cli.StringFlag{Name: "currency_symbol", Usage: "Currency symbol"},
					&cli.StringFlag{Name: "currency_format", Usage: "Currency format"},
					&cli.StringFlag{Name: "json", Usage: "Additional fields as JSON"},
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
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
					if err := internal.MergeJSON(cmd, body); err != nil {
						return err
					}
					raw, err := c.Request(ctx, "PUT", c.InventoryBase+"/settings/currencies/"+cmd.Args().First(), &zohttp.RequestOpts{
						Params: orgParams(orgID),
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
				Action: func(ctx context.Context, cmd *cli.Command) error {
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
					raw, err := c.Request(ctx, "DELETE", c.InventoryBase+"/settings/currencies/"+cmd.Args().First(), &zohttp.RequestOpts{Params: orgParams(orgID)})
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
				Flags: []cli.Flag{
					&cli.BoolFlag{Name: "all", Usage: "Fetch all records"},
					&cli.IntFlag{Name: "limit", Usage: "Max total records to fetch"},
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := internal.RequireFlag(cmd, "org", "ZOHO_BOOKS_ORG_ID")
					if err != nil {
						return err
					}
					params := orgParams(orgID)
					if cmd.Bool("all") || cmd.IsSet("limit") {
						items, err := pagination.Paginate(ctx, pagination.PaginationConfig{
							Client:   c,
							URL:      c.InventoryBase+"/settings/taxes",
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
					raw, err := c.Request(ctx, "GET", c.InventoryBase+"/settings/taxes", &zohttp.RequestOpts{Params: params})
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
				Action: func(ctx context.Context, cmd *cli.Command) error {
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
					raw, err := c.Request(ctx, "GET", c.InventoryBase+"/settings/taxes/"+cmd.Args().First(), &zohttp.RequestOpts{Params: orgParams(orgID)})
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
				Action: func(ctx context.Context, cmd *cli.Command) error {
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
					raw, err := c.Request(ctx, "POST", c.InventoryBase+"/settings/taxes", &zohttp.RequestOpts{
						Params: orgParams(orgID),
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
				Action: func(ctx context.Context, cmd *cli.Command) error {
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
					raw, err := c.Request(ctx, "PUT", c.InventoryBase+"/settings/taxes/"+cmd.Args().First(), &zohttp.RequestOpts{
						Params: orgParams(orgID),
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
				Action: func(ctx context.Context, cmd *cli.Command) error {
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
					raw, err := c.Request(ctx, "DELETE", c.InventoryBase+"/settings/taxes/"+cmd.Args().First(), &zohttp.RequestOpts{Params: orgParams(orgID)})
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
				Action: func(ctx context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := internal.RequireFlag(cmd, "org", "ZOHO_BOOKS_ORG_ID")
					if err != nil {
						return err
					}
					raw, err := c.Request(ctx, "GET", c.InventoryBase+"/users", &zohttp.RequestOpts{Params: orgParams(orgID)})
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
				Action: func(ctx context.Context, cmd *cli.Command) error {
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
					raw, err := c.Request(ctx, "GET", c.InventoryBase+"/users/"+cmd.Args().First(), &zohttp.RequestOpts{Params: orgParams(orgID)})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:  "get-current",
				Usage: "Get current user",
				Action: func(ctx context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := internal.RequireFlag(cmd, "org", "ZOHO_BOOKS_ORG_ID")
					if err != nil {
						return err
					}
					raw, err := c.Request(ctx, "GET", c.InventoryBase+"/users/me", &zohttp.RequestOpts{Params: orgParams(orgID)})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
		},
	}
}
