package books

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/omin8tor/zoho-cli/internal"
	zohttp "github.com/omin8tor/zoho-cli/internal/http"
	"github.com/omin8tor/zoho-cli/internal/output"
	"github.com/omin8tor/zoho-cli/internal/pagination"
	"github.com/urfave/cli/v3"
)

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
					raw, err := c.Request(ctx, "POST", c.BooksBase+"/chartofaccounts", &zohttp.RequestOpts{
						Params: orgParams(orgID),
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
					if v := cmd.String("sort-column"); v != "" {
						params["sort_column"] = v
					}

					if cmd.Bool("all") || cmd.IsSet("limit") {

						items, err := pagination.Paginate(ctx, pagination.PaginationConfig{

							Client: c,

							URL: c.BooksBase + "/chartofaccounts",

							Opts: &zohttp.RequestOpts{Params: params},

							ItemsKey: "chartofaccounts",

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
					raw, err := c.Request(ctx, "GET", c.BooksBase+"/chartofaccounts", &zohttp.RequestOpts{Params: params})
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
					raw, err := c.Request(ctx, "PUT", c.BooksBase+"/chartofaccounts/"+cmd.Args().First(), &zohttp.RequestOpts{
						Params: orgParams(orgID),
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
				Action: func(ctx context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request(ctx, "GET", c.BooksBase+"/chartofaccounts/"+cmd.Args().First(), &zohttp.RequestOpts{Params: orgParams(orgID)})
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
				Action: func(ctx context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request(ctx, "DELETE", c.BooksBase+"/chartofaccounts/"+cmd.Args().First(), &zohttp.RequestOpts{Params: orgParams(orgID)})
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
				Action: func(ctx context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request(ctx, "POST", c.BooksBase+"/chartofaccounts/"+cmd.Args().First()+"/active", &zohttp.RequestOpts{Params: orgParams(orgID)})
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
				Action: func(ctx context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request(ctx, "POST", c.BooksBase+"/chartofaccounts/"+cmd.Args().First()+"/inactive", &zohttp.RequestOpts{Params: orgParams(orgID)})
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
				Action: func(ctx context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request(ctx, "GET", c.BooksBase+"/chartofaccounts/"+cmd.Args().First()+"/transactions", &zohttp.RequestOpts{Params: orgParams(orgID)})
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
				Action: func(ctx context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request(ctx, "DELETE", c.BooksBase+"/chartofaccounts/transactions/"+cmd.Args().First(), &zohttp.RequestOpts{Params: orgParams(orgID)})
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
					body["reference_number"] = cmd.String("reference_number")
					body["journal_date"] = cmd.String("journal_date")
					if cmd.IsSet("currency_id") {
						body["currency_id"] = cmd.String("currency_id")
					}
					if err := internal.MergeJSON(cmd, body); err != nil {
						return err
					}
					raw, err := c.Request(ctx, "POST", c.BooksBase+"/journals", &zohttp.RequestOpts{
						Params: orgParams(orgID),
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
					if v := cmd.String("date-start"); v != "" {
						params["date_start"] = v
					}
					if v := cmd.String("date-end"); v != "" {
						params["date_end"] = v
					}

					if cmd.Bool("all") || cmd.IsSet("limit") {

						items, err := pagination.Paginate(ctx, pagination.PaginationConfig{

							Client: c,

							URL: c.BooksBase + "/journals",

							Opts: &zohttp.RequestOpts{Params: params},

							ItemsKey: "journals",

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
					raw, err := c.Request(ctx, "GET", c.BooksBase+"/journals", &zohttp.RequestOpts{Params: params})
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
					raw, err := c.Request(ctx, "PUT", c.BooksBase+"/journals/"+cmd.Args().First(), &zohttp.RequestOpts{
						Params: orgParams(orgID),
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
				Action: func(ctx context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request(ctx, "GET", c.BooksBase+"/journals/"+cmd.Args().First(), &zohttp.RequestOpts{Params: orgParams(orgID)})
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
				Action: func(ctx context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request(ctx, "DELETE", c.BooksBase+"/journals/"+cmd.Args().First(), &zohttp.RequestOpts{Params: orgParams(orgID)})
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
				Action: func(ctx context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request(ctx, "POST", c.BooksBase+"/journals/"+cmd.Args().First()+"/publish", &zohttp.RequestOpts{Params: orgParams(orgID)})
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
				Action: func(ctx context.Context, cmd *cli.Command) error {
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
					raw, err := c.Request(ctx, "POST", c.BooksBase+"/journals/"+cmd.Args().First()+"/attachment", &zohttp.RequestOpts{
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
					body["description"] = cmd.String("description")
					if err := internal.MergeJSON(cmd, body); err != nil {
						return err
					}
					raw, err := c.Request(ctx, "POST", c.BooksBase+"/journals/"+cmd.Args().First()+"/comments", &zohttp.RequestOpts{
						Params: orgParams(orgID),
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
				Action: func(ctx context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request(ctx, "DELETE", c.BooksBase+"/journals/"+cmd.Args().First()+"/comments/"+cmd.Args().Get(1), &zohttp.RequestOpts{Params: orgParams(orgID)})
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
					raw, err := c.Request(ctx, "POST", c.BooksBase+"/fixedassets", &zohttp.RequestOpts{
						Params: orgParams(orgID),
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

					if cmd.Bool("all") || cmd.IsSet("limit") {

						items, err := pagination.Paginate(ctx, pagination.PaginationConfig{

							Client: c,

							URL: c.BooksBase + "/fixedassets",

							Opts: &zohttp.RequestOpts{Params: params},

							ItemsKey: "fixedassets",

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
					raw, err := c.Request(ctx, "GET", c.BooksBase+"/fixedassets", &zohttp.RequestOpts{Params: params})
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
					raw, err := c.Request(ctx, "PUT", c.BooksBase+"/fixedassets/"+cmd.Args().First(), &zohttp.RequestOpts{
						Params: orgParams(orgID),
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
				Action: func(ctx context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request(ctx, "GET", c.BooksBase+"/fixedassets/"+cmd.Args().First(), &zohttp.RequestOpts{Params: orgParams(orgID)})
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
				Action: func(ctx context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request(ctx, "DELETE", c.BooksBase+"/fixedassets/"+cmd.Args().First(), &zohttp.RequestOpts{Params: orgParams(orgID)})
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
				Action: func(ctx context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request(ctx, "GET", c.BooksBase+"/fixedassets/"+cmd.Args().First()+"/comments", &zohttp.RequestOpts{Params: orgParams(orgID)})
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
				Action: func(ctx context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request(ctx, "GET", c.BooksBase+"/fixedassets/"+cmd.Args().First()+"/depreciation", &zohttp.RequestOpts{Params: orgParams(orgID)})
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
				Action: func(ctx context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request(ctx, "POST", c.BooksBase+"/fixedassets/"+cmd.Args().First()+"/active", &zohttp.RequestOpts{Params: orgParams(orgID)})
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
				Action: func(ctx context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request(ctx, "POST", c.BooksBase+"/fixedassets/"+cmd.Args().First()+"/cancel", &zohttp.RequestOpts{Params: orgParams(orgID)})
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
				Action: func(ctx context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request(ctx, "POST", c.BooksBase+"/fixedassets/"+cmd.Args().First()+"/draft", &zohttp.RequestOpts{Params: orgParams(orgID)})
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
					raw, err := c.Request(ctx, "POST", c.BooksBase+"/fixedassets/"+cmd.Args().First()+"/writeoff", &zohttp.RequestOpts{
						Params: orgParams(orgID),
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
					raw, err := c.Request(ctx, "POST", c.BooksBase+"/fixedassets/"+cmd.Args().First()+"/sell", &zohttp.RequestOpts{
						Params: orgParams(orgID),
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
					body["description"] = cmd.String("description")
					if err := internal.MergeJSON(cmd, body); err != nil {
						return err
					}
					raw, err := c.Request(ctx, "POST", c.BooksBase+"/fixedassets/"+cmd.Args().First()+"/comments", &zohttp.RequestOpts{
						Params: orgParams(orgID),
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
				Action: func(ctx context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request(ctx, "DELETE", c.BooksBase+"/fixedassets/"+cmd.Args().First()+"/comments/"+cmd.Args().Get(1), &zohttp.RequestOpts{Params: orgParams(orgID)})
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
					raw, err := c.Request(ctx, "POST", c.BooksBase+"/fixedassets"+"/types", &zohttp.RequestOpts{
						Params: orgParams(orgID),
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
				Action: func(ctx context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request(ctx, "GET", c.BooksBase+"/fixedassets"+"/types", &zohttp.RequestOpts{Params: orgParams(orgID)})
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
					raw, err := c.Request(ctx, "PUT", c.BooksBase+"/fixedassets"+"/types/"+cmd.Args().First(), &zohttp.RequestOpts{
						Params: orgParams(orgID),
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
				Action: func(ctx context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request(ctx, "DELETE", c.BooksBase+"/fixedassets"+"/types/"+cmd.Args().First(), &zohttp.RequestOpts{Params: orgParams(orgID)})
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
					body["adjustment_date"] = cmd.String("adjustment_date")
					body["currency_id"] = cmd.String("currency_id")
					body["exchange_rate"] = cmd.Float("exchange_rate")
					if cmd.IsSet("notes") {
						body["notes"] = cmd.String("notes")
					}
					if err := internal.MergeJSON(cmd, body); err != nil {
						return err
					}
					raw, err := c.Request(ctx, "POST", c.BooksBase+"/basecurrencyadjustment", &zohttp.RequestOpts{
						Params: orgParams(orgID),
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

					if cmd.Bool("all") || cmd.IsSet("limit") {

						items, err := pagination.Paginate(ctx, pagination.PaginationConfig{

							Client: c,

							URL: c.BooksBase + "/basecurrencyadjustment",

							Opts: &zohttp.RequestOpts{Params: params},

							ItemsKey: "basecurrencyadjustments",

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
					raw, err := c.Request(ctx, "GET", c.BooksBase+"/basecurrencyadjustment", &zohttp.RequestOpts{Params: params})
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
				Action: func(ctx context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request(ctx, "GET", c.BooksBase+"/basecurrencyadjustment/"+cmd.Args().First(), &zohttp.RequestOpts{Params: orgParams(orgID)})
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
				Action: func(ctx context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request(ctx, "DELETE", c.BooksBase+"/basecurrencyadjustment/"+cmd.Args().First(), &zohttp.RequestOpts{Params: orgParams(orgID)})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:  "list-account-details",
				Usage: "List account details for adjustment",
				Action: func(ctx context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request(ctx, "GET", c.BooksBase+"/basecurrencyadjustment"+"/accounts", &zohttp.RequestOpts{Params: orgParams(orgID)})
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
					body["account_id"] = cmd.String("account_id")
					body["opening_balance"] = cmd.Float("opening_balance")
					if err := internal.MergeJSON(cmd, body); err != nil {
						return err
					}
					raw, err := c.Request(ctx, "POST", c.BooksBase+"/settings/openingbalances", &zohttp.RequestOpts{
						Params: orgParams(orgID),
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
					if cmd.IsSet("account_id") {
						body["account_id"] = cmd.String("account_id")
					}
					if cmd.IsSet("opening_balance") {
						body["opening_balance"] = cmd.Float("opening_balance")
					}
					if err := internal.MergeJSON(cmd, body); err != nil {
						return err
					}
					raw, err := c.Request(ctx, "PUT", c.BooksBase+"/settings/openingbalances", &zohttp.RequestOpts{
						Params: orgParams(orgID),
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
				Action: func(ctx context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request(ctx, "GET", c.BooksBase+"/settings/openingbalances", &zohttp.RequestOpts{Params: orgParams(orgID)})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:  "delete",
				Usage: "Delete opening balance",
				Action: func(ctx context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request(ctx, "DELETE", c.BooksBase+"/settings/openingbalances", &zohttp.RequestOpts{Params: orgParams(orgID)})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
		},
	}
}
