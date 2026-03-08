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
					c, err := zohttp.GetClient()
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
					c, err := zohttp.GetClient()
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

					if cmd.Bool("all") || cmd.IsSet("limit") {

						items, err := pagination.Paginate(pagination.PaginationConfig{

							Client: c,

							URL: c.BooksBase + "/custommodules/" + cmd.String("module"),

							Opts: &zohttp.RequestOpts{Params: params},

							ItemsKey: "module_records",

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
					c, err := zohttp.GetClient()
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
					c, err := zohttp.GetClient()
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
					c, err := zohttp.GetClient()
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
					c, err := zohttp.GetClient()
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
					c, err := zohttp.GetClient()
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
					c, err := zohttp.GetClient()
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
					c, err := zohttp.GetClient()
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
					c, err := zohttp.GetClient()
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
					c, err := zohttp.GetClient()
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
					c, err := zohttp.GetClient()
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
					c, err := zohttp.GetClient()
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
					c, err := zohttp.GetClient()
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
					c, err := zohttp.GetClient()
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
					c, err := zohttp.GetClient()
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
					c, err := zohttp.GetClient()
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
					c, err := zohttp.GetClient()
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
					c, err := zohttp.GetClient()
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
