package sheet

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	internal "github.com/omin8tor/zoho-cli/internal"
	zohttp "github.com/omin8tor/zoho-cli/internal/http"
	"github.com/omin8tor/zoho-cli/internal/output"
	"github.com/omin8tor/zoho-cli/internal/pagination"
	"github.com/urfave/cli/v3"
)

var workbookFlag = &cli.StringFlag{Name: "workbook", Required: true, Usage: "Workbook resource ID"}

var worksheetFlag = &cli.StringFlag{Name: "worksheet", Usage: "Worksheet name"}

func Commands() *cli.Command {
	return &cli.Command{
		Name:  "sheet",
		Usage: "Zoho Sheet operations",
		Commands: []*cli.Command{
			workbooksCmd(),
			worksheetsCmd(),
			tablesCmd(),
			recordsCmd(),
			cellsCmd(),
			contentCmd(),
			formatCmd(),
			namedRangesCmd(),
			mergeCmd(),
			premiumCmd(),
			utilityCmd(),
		},
	}
}

func workbooksCmd() *cli.Command {
	return &cli.Command{
		Name:  "workbooks",
		Usage: "Workbook operations",
		Commands: []*cli.Command{
			{
				Name:  "list",
				Usage: "List all workbooks",
				Flags: []cli.Flag{
					&cli.BoolFlag{Name: "all", Usage: "Fetch all records"},
					&cli.IntFlag{Name: "limit", Usage: "Max total records to fetch"},
					&cli.StringFlag{Name: "sort-option", Usage: "Sort option"},
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					params := map[string]string{"method": "workbook.list"}
					if v := cmd.String("sort-option"); v != "" {
						params["sort_option"] = v
					}
					if cmd.Bool("all") || cmd.IsSet("limit") {
						setPage := func(state *pagination.PageState, p map[string]string) {
							p["start_index"] = fmt.Sprintf("%d", state.Offset)
							p["count"] = "100"
						}
						items, err := pagination.Paginate(ctx, pagination.PaginationConfig{
							Client:   c,
							Method:   "POST",
							URL:      c.SheetBase + "/workbooks",
							Opts:     &zohttp.RequestOpts{Params: params},
							ItemsKey: "workbooks",
							PageSize: 100,
							Limit:    cmd.Int("limit"),
							SetPage:  setPage,
							HasMore:  pagination.HasMoreByCount,
						})
						if err != nil {
							return err
						}
						return output.JSON(items)
					}
					raw, err := c.Request(ctx, "POST", c.SheetBase+"/workbooks", &zohttp.RequestOpts{Params: params})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:  "templates",
				Usage: "List all templates",
				Flags: []cli.Flag{
					&cli.IntFlag{Name: "start-index", Usage: "Start index"},
					&cli.IntFlag{Name: "count", Usage: "Number of templates"},
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					params := map[string]string{"method": "template.list"}
					if v := cmd.Int("start-index"); v > 0 {
						params["start_index"] = fmt.Sprintf("%d", v)
					}
					if v := cmd.Int("count"); v > 0 {
						params["count"] = fmt.Sprintf("%d", v)
					}
					raw, err := c.Request(ctx, "POST", c.SheetBase+"/templates", &zohttp.RequestOpts{Params: params})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:  "versions",
				Usage: "List all versions",
				Flags: []cli.Flag{
					workbookFlag,
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					params := map[string]string{"method": "workbook.version.list"}
					raw, err := c.Request(ctx, "POST", c.SheetBase+"/"+cmd.String("workbook"), &zohttp.RequestOpts{Params: params})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:  "create",
				Usage: "Create workbook",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "workbook-name", Required: true, Usage: "Workbook name"},
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					params := map[string]string{
						"method":        "workbook.create",
						"workbook_name": cmd.String("workbook-name"),
					}
					raw, err := c.Request(ctx, "POST", c.SheetBase+"/create", &zohttp.RequestOpts{Params: params})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:  "create-from-template",
				Usage: "Create workbook from template",
				Flags: []cli.Flag{
					workbookFlag,
					&cli.StringFlag{Name: "workbook-name", Required: true, Usage: "New workbook name"},
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					params := map[string]string{
						"method":        "workbook.createfromtemplate",
						"resource_id":   cmd.String("workbook"),
						"workbook_name": cmd.String("workbook-name"),
					}
					raw, err := c.Request(ctx, "POST", c.SheetBase+"/createfromtemplate", &zohttp.RequestOpts{Params: params})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:  "upload",
				Usage: "Upload workbook",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "file", Required: true, Usage: "Path to file"},
					&cli.StringFlag{Name: "workbook-name", Usage: "Workbook name"},
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					filePath := cmd.String("file")
					data, err := os.ReadFile(filePath)
					if err != nil {
						return fmt.Errorf("failed to read file: %w", err)
					}
					name := filepath.Base(filePath)
					params := map[string]string{"method": "workbook.upload"}
					form := map[string]string{}
					if v := cmd.String("workbook-name"); v != "" {
						form["workbook_name"] = v
					}
					raw, err := c.Request(ctx, "POST", c.SheetBase+"/upload", &zohttp.RequestOpts{
						Params: params,
						Files:  map[string]zohttp.FileUpload{"file": {Filename: name, Data: data}},
						Form:   form,
					})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:  "download",
				Usage: "Download workbook",
				Flags: []cli.Flag{
					workbookFlag,
					&cli.StringFlag{Name: "format", Required: true, Usage: "Download format (xlsx/csv/tsv/ods/pdf/html)"},
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					params := map[string]string{
						"method":          "workbook.download",
						"download_format": cmd.String("format"),
					}
					raw, err := c.Request(ctx, "POST", c.SheetBase+"/download/"+cmd.String("workbook"), &zohttp.RequestOpts{Params: params})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:  "insert-images",
				Usage: "Insert images into workbook",
				Flags: []cli.Flag{
					workbookFlag,
					&cli.StringFlag{Name: "image-json", Required: true, Usage: "Image JSON configuration"},
					&cli.StringFlag{Name: "file", Usage: "Path to image file"},
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					params := map[string]string{
						"method":      "workbook.images.insert",
						"resource_id": cmd.String("workbook"),
						"image_json":  cmd.String("image-json"),
					}
					opts := &zohttp.RequestOpts{Params: params}
					if filePath := cmd.String("file"); filePath != "" {
						data, err := os.ReadFile(filePath)
						if err != nil {
							return fmt.Errorf("failed to read file: %w", err)
						}
						name := filepath.Base(filePath)
						opts.Files = map[string]zohttp.FileUpload{"imagefiles": {Filename: name, Data: data}}
					}
					raw, err := c.Request(ctx, "POST", c.SheetBase+"/insertimages", opts)
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:  "copy",
				Usage: "Copy workbook",
				Flags: []cli.Flag{
					workbookFlag,
					&cli.StringFlag{Name: "new-workbook-name", Required: true, Usage: "New workbook name"},
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					params := map[string]string{
						"method":            "workbook.copy",
						"resource_id":       cmd.String("workbook"),
						"new_workbook_name": cmd.String("new-workbook-name"),
					}
					raw, err := c.Request(ctx, "POST", c.SheetBase+"/copy", &zohttp.RequestOpts{Params: params})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:  "share",
				Usage: "Share workbook",
				Flags: []cli.Flag{
					workbookFlag,
					&cli.StringFlag{Name: "email", Required: true, Usage: "Email address to share with"},
					&cli.StringFlag{Name: "role", Usage: "Role for shared user"},
					&cli.StringFlag{Name: "notify", Usage: "Notify user (true/false)"},
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					params := map[string]string{
						"method":      "workbook.share",
						"resource_id": cmd.String("workbook"),
						"email_id":    cmd.String("email"),
					}
					if v := cmd.String("role"); v != "" {
						params["role"] = v
					}
					if v := cmd.String("notify"); v != "" {
						params["notify"] = v
					}
					raw, err := c.Request(ctx, "POST", c.SheetBase+"/share", &zohttp.RequestOpts{Params: params})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:  "create-version",
				Usage: "Create a version",
				Flags: []cli.Flag{
					workbookFlag,
					&cli.StringFlag{Name: "version-description", Required: true, Usage: "Version description"},
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					params := map[string]string{
						"method":              "workbook.version.create",
						"version_description": cmd.String("version-description"),
					}
					raw, err := c.Request(ctx, "POST", c.SheetBase+"/"+cmd.String("workbook"), &zohttp.RequestOpts{Params: params})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:  "revert-version",
				Usage: "Revert to a version",
				Flags: []cli.Flag{
					workbookFlag,
					&cli.StringFlag{Name: "version-number", Required: true, Usage: "Version number"},
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					params := map[string]string{
						"method":         "workbook.version.revert",
						"version_number": cmd.String("version-number"),
					}
					raw, err := c.Request(ctx, "POST", c.SheetBase+"/"+cmd.String("workbook"), &zohttp.RequestOpts{Params: params})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:  "trash",
				Usage: "Trash workbook",
				Flags: []cli.Flag{
					workbookFlag,
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					params := map[string]string{"method": "workbook.trash", "resource_ids": fmt.Sprintf(`["%s"]`, cmd.String("workbook"))}
					raw, err := c.Request(ctx, "POST", c.SheetBase+"/trash", &zohttp.RequestOpts{Params: params})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:  "restore",
				Usage: "Restore workbook",
				Flags: []cli.Flag{
					workbookFlag,
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					params := map[string]string{"method": "workbook.restore", "resource_ids": fmt.Sprintf(`["%s"]`, cmd.String("workbook"))}
					raw, err := c.Request(ctx, "POST", c.SheetBase+"/restore", &zohttp.RequestOpts{Params: params})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:  "delete",
				Usage: "Delete workbook",
				Flags: []cli.Flag{
					workbookFlag,
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					params := map[string]string{"method": "workbook.delete", "resource_ids": fmt.Sprintf(`["%s"]`, cmd.String("workbook"))}
					raw, err := c.Request(ctx, "POST", c.SheetBase+"/delete", &zohttp.RequestOpts{Params: params})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:  "publish",
				Usage: "Publish workbook",
				Flags: []cli.Flag{
					workbookFlag,
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					params := map[string]string{"method": "workbook.publish", "resource_id": cmd.String("workbook")}
					raw, err := c.Request(ctx, "POST", c.SheetBase+"/publish", &zohttp.RequestOpts{Params: params})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:  "unpublish",
				Usage: "Remove publish from workbook",
				Flags: []cli.Flag{
					workbookFlag,
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					params := map[string]string{"method": "workbook.publish.remove", "resource_id": cmd.String("workbook")}
					raw, err := c.Request(ctx, "POST", c.SheetBase+"/publish", &zohttp.RequestOpts{Params: params})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:  "lock",
				Usage: "Lock workbook",
				Flags: []cli.Flag{
					workbookFlag,
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					params := map[string]string{"method": "lock"}
					raw, err := c.Request(ctx, "POST", c.SheetBase+"/"+cmd.String("workbook"), &zohttp.RequestOpts{Params: params})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:  "unlock",
				Usage: "Unlock workbook",
				Flags: []cli.Flag{
					workbookFlag,
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					params := map[string]string{"method": "unlock"}
					raw, err := c.Request(ctx, "POST", c.SheetBase+"/"+cmd.String("workbook"), &zohttp.RequestOpts{Params: params})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
		},
	}
}

func worksheetsCmd() *cli.Command {
	return &cli.Command{
		Name:  "worksheets",
		Usage: "Worksheet operations",
		Commands: []*cli.Command{
			{
				Name:  "list",
				Usage: "List all worksheets",
				Flags: []cli.Flag{
					workbookFlag,
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					params := map[string]string{"method": "worksheet.list"}
					raw, err := c.Request(ctx, "POST", c.SheetBase+"/"+cmd.String("workbook"), &zohttp.RequestOpts{Params: params})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:  "create",
				Usage: "Create worksheet",
				Flags: []cli.Flag{
					workbookFlag,
					&cli.StringFlag{Name: "worksheet-name", Required: true, Usage: "Worksheet name"},
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					params := map[string]string{
						"method":         "worksheet.insert",
						"worksheet_name": cmd.String("worksheet-name"),
					}
					raw, err := c.Request(ctx, "POST", c.SheetBase+"/"+cmd.String("workbook"), &zohttp.RequestOpts{Params: params})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:  "copy",
				Usage: "Copy worksheet within same workbook",
				Flags: []cli.Flag{
					workbookFlag,
					worksheetFlag,
					&cli.StringFlag{Name: "new-worksheet-name", Required: true, Usage: "New worksheet name"},
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					params := map[string]string{
						"method":             "worksheet.copy",
						"worksheet_name":     cmd.String("worksheet"),
						"new_worksheet_name": cmd.String("new-worksheet-name"),
					}
					raw, err := c.Request(ctx, "POST", c.SheetBase+"/"+cmd.String("workbook"), &zohttp.RequestOpts{Params: params})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:  "copy-to",
				Usage: "Copy worksheet to another workbook",
				Flags: []cli.Flag{
					workbookFlag,
					worksheetFlag,
					&cli.StringFlag{Name: "dest-workbook", Required: true, Usage: "Destination workbook resource ID"},
					&cli.StringFlag{Name: "new-worksheet-name", Usage: "New worksheet name"},
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					params := map[string]string{
						"method":           "worksheet.copy.otherdoc",
						"worksheet_name":   cmd.String("worksheet"),
						"dest_resource_id": cmd.String("dest-workbook"),
					}
					if v := cmd.String("new-worksheet-name"); v != "" {
						params["new_worksheet_name"] = v
					}
					raw, err := c.Request(ctx, "POST", c.SheetBase+"/"+cmd.String("workbook"), &zohttp.RequestOpts{Params: params})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:  "rename",
				Usage: "Rename worksheet",
				Flags: []cli.Flag{
					workbookFlag,
					worksheetFlag,
					&cli.StringFlag{Name: "new-worksheet-name", Required: true, Usage: "New worksheet name"},
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					params := map[string]string{
						"method":   "worksheet.rename",
						"old_name": cmd.String("worksheet"),
						"new_name": cmd.String("new-worksheet-name"),
					}
					raw, err := c.Request(ctx, "POST", c.SheetBase+"/"+cmd.String("workbook"), &zohttp.RequestOpts{Params: params})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:  "delete",
				Usage: "Delete worksheet",
				Flags: []cli.Flag{
					workbookFlag,
					worksheetFlag,
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					params := map[string]string{
						"method":         "worksheet.delete",
						"worksheet_name": cmd.String("worksheet"),
					}
					raw, err := c.Request(ctx, "POST", c.SheetBase+"/"+cmd.String("workbook"), &zohttp.RequestOpts{Params: params})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:  "delete-multiple",
				Usage: "Delete multiple worksheets",
				Flags: []cli.Flag{
					workbookFlag,
					&cli.StringFlag{Name: "worksheet-names", Required: true, Usage: "JSON array of worksheet names"},
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					params := map[string]string{
						"method":          "worksheets.delete",
						"worksheet_names": cmd.String("worksheet-names"),
					}
					raw, err := c.Request(ctx, "POST", c.SheetBase+"/"+cmd.String("workbook"), &zohttp.RequestOpts{Params: params})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
		},
	}
}

func tablesCmd() *cli.Command {
	return &cli.Command{
		Name:  "tables",
		Usage: "Table operations",
		Commands: []*cli.Command{
			{
				Name:  "list",
				Usage: "List all tables",
				Flags: []cli.Flag{
					workbookFlag,
					worksheetFlag,
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					params := map[string]string{"method": "table.list"}
					if v := cmd.String("worksheet"); v != "" {
						params["worksheet_name"] = v
					}
					raw, err := c.Request(ctx, "POST", c.SheetBase+"/"+cmd.String("workbook"), &zohttp.RequestOpts{Params: params})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:  "create",
				Usage: "Create table",
				Flags: []cli.Flag{
					workbookFlag,
					worksheetFlag,
					&cli.StringFlag{Name: "table-name", Required: true, Usage: "Table name"},
					&cli.IntFlag{Name: "start-row", Required: true, Usage: "Start row"},
					&cli.IntFlag{Name: "start-column", Required: true, Usage: "Start column"},
					&cli.IntFlag{Name: "end-row", Required: true, Usage: "End row"},
					&cli.IntFlag{Name: "end-column", Required: true, Usage: "End column"},
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					params := map[string]string{
						"method":       "table.create",
						"table_name":   cmd.String("table-name"),
						"start_row":    fmt.Sprintf("%d", cmd.Int("start-row")),
						"start_column": fmt.Sprintf("%d", cmd.Int("start-column")),
						"end_row":      fmt.Sprintf("%d", cmd.Int("end-row")),
						"end_column":   fmt.Sprintf("%d", cmd.Int("end-column")),
					}
					if v := cmd.String("worksheet"); v != "" {
						params["worksheet_name"] = v
					}
					raw, err := c.Request(ctx, "POST", c.SheetBase+"/"+cmd.String("workbook"), &zohttp.RequestOpts{Params: params})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:  "remove",
				Usage: "Remove table",
				Flags: []cli.Flag{
					workbookFlag,
					&cli.StringFlag{Name: "table-name", Required: true, Usage: "Table name"},
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					params := map[string]string{
						"method":     "table.remove",
						"table_name": cmd.String("table-name"),
					}
					raw, err := c.Request(ctx, "POST", c.SheetBase+"/"+cmd.String("workbook"), &zohttp.RequestOpts{Params: params})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:  "rename-headers",
				Usage: "Rename headers of table",
				Flags: []cli.Flag{
					workbookFlag,
					&cli.StringFlag{Name: "table-name", Required: true, Usage: "Table name"},
					&cli.StringFlag{Name: "data", Required: true, Usage: "Header rename data JSON"},
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					params := map[string]string{
						"method":     "table.header.rename",
						"table_name": cmd.String("table-name"),
						"data":       cmd.String("data"),
					}
					raw, err := c.Request(ctx, "POST", c.SheetBase+"/"+cmd.String("workbook"), &zohttp.RequestOpts{Params: params})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:  "fetch-records",
				Usage: "Fetch records from table",
				Flags: []cli.Flag{
					workbookFlag,
					&cli.StringFlag{Name: "table-name", Required: true, Usage: "Table name"},
					&cli.StringFlag{Name: "criteria", Usage: "Filter criteria"},
					&cli.StringFlag{Name: "criteria-json", Usage: "Filter criteria JSON"},
					&cli.IntFlag{Name: "start-index", Usage: "Start index"},
					&cli.IntFlag{Name: "count", Usage: "Number of records"},
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					params := map[string]string{
						"method":     "table.records.fetch",
						"table_name": cmd.String("table-name"),
					}
					if v := cmd.String("criteria"); v != "" {
						params["criteria"] = v
					}
					if v := cmd.String("criteria-json"); v != "" {
						params["criteria"] = v
					}
					if v := cmd.Int("start-index"); v > 0 {
						params["start_index"] = fmt.Sprintf("%d", v)
					}
					if v := cmd.Int("count"); v > 0 {
						params["count"] = fmt.Sprintf("%d", v)
					}
					raw, err := c.Request(ctx, "POST", c.SheetBase+"/"+cmd.String("workbook"), &zohttp.RequestOpts{Params: params})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:  "add-records",
				Usage: "Add records to table",
				Flags: []cli.Flag{
					workbookFlag,
					&cli.StringFlag{Name: "table-name", Required: true, Usage: "Table name"},
					&cli.StringFlag{Name: "data", Required: true, Usage: "Record data as JSON array"},
					&cli.StringFlag{Name: "json", Usage: "Additional fields as JSON"},
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					params := map[string]string{
						"method":     "table.records.add",
						"table_name": cmd.String("table-name"),
						"json_data":  cmd.String("data"),
					}
					if err := internal.MergeJSONForm(cmd, params); err != nil {
						return err
					}
					raw, err := c.Request(ctx, "POST", c.SheetBase+"/"+cmd.String("workbook"), &zohttp.RequestOpts{Params: params})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:  "update-records",
				Usage: "Update records in table",
				Flags: []cli.Flag{
					workbookFlag,
					&cli.StringFlag{Name: "table-name", Required: true, Usage: "Table name"},
					&cli.StringFlag{Name: "criteria", Required: true, Usage: "Filter criteria"},
					&cli.StringFlag{Name: "data", Required: true, Usage: "Update data as JSON"},
					&cli.StringFlag{Name: "json", Usage: "Additional fields as JSON"},
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					params := map[string]string{
						"method":     "table.records.update",
						"table_name": cmd.String("table-name"),
						"criteria":   cmd.String("criteria"),
						"data":       cmd.String("data"),
					}
					if err := internal.MergeJSONForm(cmd, params); err != nil {
						return err
					}
					raw, err := c.Request(ctx, "POST", c.SheetBase+"/"+cmd.String("workbook"), &zohttp.RequestOpts{Params: params})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:  "delete-records",
				Usage: "Delete records from table",
				Flags: []cli.Flag{
					workbookFlag,
					&cli.StringFlag{Name: "table-name", Required: true, Usage: "Table name"},
					&cli.StringFlag{Name: "criteria", Required: true, Usage: "Filter criteria"},
					&cli.StringFlag{Name: "delete-rows", Usage: "Delete entire rows (true/false)"},
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					params := map[string]string{
						"method":        "table.records.delete",
						"table_name":    cmd.String("table-name"),
						"criteria_json": cmd.String("criteria"),
					}
					if v := cmd.String("delete-rows"); v != "" {
						params["delete_rows"] = v
					}
					raw, err := c.Request(ctx, "POST", c.SheetBase+"/"+cmd.String("workbook"), &zohttp.RequestOpts{Params: params})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:  "insert-columns",
				Usage: "Insert columns to table",
				Flags: []cli.Flag{
					workbookFlag,
					&cli.StringFlag{Name: "table-name", Required: true, Usage: "Table name"},
					&cli.StringFlag{Name: "columns", Required: true, Usage: "Columns JSON"},
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					params := map[string]string{
						"method":       "table.columns.insert",
						"table_name":   cmd.String("table-name"),
						"column_names": cmd.String("columns"),
					}
					raw, err := c.Request(ctx, "POST", c.SheetBase+"/"+cmd.String("workbook"), &zohttp.RequestOpts{Params: params})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:  "delete-columns",
				Usage: "Delete columns from table",
				Flags: []cli.Flag{
					workbookFlag,
					&cli.StringFlag{Name: "table-name", Required: true, Usage: "Table name"},
					&cli.StringFlag{Name: "columns", Required: true, Usage: "Columns JSON"},
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					params := map[string]string{
						"method":       "table.columns.delete",
						"table_name":   cmd.String("table-name"),
						"column_names": cmd.String("columns"),
					}
					raw, err := c.Request(ctx, "POST", c.SheetBase+"/"+cmd.String("workbook"), &zohttp.RequestOpts{Params: params})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
		},
	}
}

func recordsCmd() *cli.Command {
	return &cli.Command{
		Name:  "records",
		Usage: "Worksheet record operations",
		Commands: []*cli.Command{
			{
				Name:  "fetch",
				Usage: "Fetch records from worksheet",
				Flags: []cli.Flag{
					workbookFlag,
					worksheetFlag,
					&cli.StringFlag{Name: "criteria", Usage: "Filter criteria"},
					&cli.IntFlag{Name: "header-row", Usage: "Header row number"},
					&cli.IntFlag{Name: "start-row", Usage: "Start row number"},
					&cli.IntFlag{Name: "count", Usage: "Number of records"},
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					params := map[string]string{"method": "worksheet.records.fetch"}
					if v := cmd.String("worksheet"); v != "" {
						params["worksheet_name"] = v
					}
					if v := cmd.String("criteria"); v != "" {
						params["criteria"] = v
					}
					if v := cmd.Int("header-row"); v > 0 {
						params["header_row"] = fmt.Sprintf("%d", v)
					}
					if v := cmd.Int("start-row"); v > 0 {
						params["start_row"] = fmt.Sprintf("%d", v)
					}
					if v := cmd.Int("count"); v > 0 {
						params["count"] = fmt.Sprintf("%d", v)
					}
					raw, err := c.Request(ctx, "POST", c.SheetBase+"/"+cmd.String("workbook"), &zohttp.RequestOpts{Params: params})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:  "add",
				Usage: "Add records to worksheet",
				Flags: []cli.Flag{
					workbookFlag,
					worksheetFlag,
					&cli.IntFlag{Name: "header-row", Usage: "Header row number"},
					&cli.StringFlag{Name: "data", Required: true, Usage: "Record data as JSON array"},
					&cli.StringFlag{Name: "json", Usage: "Additional fields as JSON"},
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					params := map[string]string{
						"method":    "worksheet.records.add",
						"json_data": cmd.String("data"),
					}
					if v := cmd.String("worksheet"); v != "" {
						params["worksheet_name"] = v
					}
					if v := cmd.Int("header-row"); v > 0 {
						params["header_row"] = fmt.Sprintf("%d", v)
					}
					if err := internal.MergeJSONForm(cmd, params); err != nil {
						return err
					}
					raw, err := c.Request(ctx, "POST", c.SheetBase+"/"+cmd.String("workbook"), &zohttp.RequestOpts{Params: params})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:  "update",
				Usage: "Update records in worksheet",
				Flags: []cli.Flag{
					workbookFlag,
					worksheetFlag,
					&cli.StringFlag{Name: "criteria", Required: true, Usage: "Filter criteria"},
					&cli.IntFlag{Name: "header-row", Usage: "Header row number"},
					&cli.StringFlag{Name: "data", Required: true, Usage: "Update data as JSON"},
					&cli.StringFlag{Name: "json", Usage: "Additional fields as JSON"},
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					params := map[string]string{
						"method":   "worksheet.records.update",
						"criteria": cmd.String("criteria"),
						"data":     cmd.String("data"),
					}
					if v := cmd.String("worksheet"); v != "" {
						params["worksheet_name"] = v
					}
					if v := cmd.Int("header-row"); v > 0 {
						params["header_row"] = fmt.Sprintf("%d", v)
					}
					if err := internal.MergeJSONForm(cmd, params); err != nil {
						return err
					}
					raw, err := c.Request(ctx, "POST", c.SheetBase+"/"+cmd.String("workbook"), &zohttp.RequestOpts{Params: params})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:  "delete",
				Usage: "Delete records from worksheet",
				Flags: []cli.Flag{
					workbookFlag,
					worksheetFlag,
					&cli.StringFlag{Name: "criteria", Required: true, Usage: "Filter criteria"},
					&cli.IntFlag{Name: "header-row", Usage: "Header row number"},
					&cli.StringFlag{Name: "delete-rows", Usage: "Delete entire rows (true/false)"},
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					params := map[string]string{
						"method":   "worksheet.records.delete",
						"criteria": cmd.String("criteria"),
					}
					if v := cmd.String("worksheet"); v != "" {
						params["worksheet_name"] = v
					}
					if v := cmd.Int("header-row"); v > 0 {
						params["header_row"] = fmt.Sprintf("%d", v)
					}
					if v := cmd.String("delete-rows"); v != "" {
						params["delete_rows"] = v
					}
					raw, err := c.Request(ctx, "POST", c.SheetBase+"/"+cmd.String("workbook"), &zohttp.RequestOpts{Params: params})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:  "insert-columns",
				Usage: "Insert columns to worksheet records",
				Flags: []cli.Flag{
					workbookFlag,
					worksheetFlag,
					&cli.StringFlag{Name: "columns", Required: true, Usage: "Columns JSON"},
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					params := map[string]string{
						"method":       "records.columns.insert",
						"column_names": cmd.String("columns"),
					}
					if v := cmd.String("worksheet"); v != "" {
						params["worksheet_name"] = v
					}
					raw, err := c.Request(ctx, "POST", c.SheetBase+"/"+cmd.String("workbook"), &zohttp.RequestOpts{Params: params})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
		},
	}
}

func cellsCmd() *cli.Command {
	return &cli.Command{
		Name:  "cells",
		Usage: "Cell and range content operations",
		Commands: []*cli.Command{
			{
				Name:  "get",
				Usage: "Get content of cell",
				Flags: []cli.Flag{
					workbookFlag,
					worksheetFlag,
					&cli.IntFlag{Name: "row", Required: true, Usage: "Row index"},
					&cli.IntFlag{Name: "column", Required: true, Usage: "Column index"},
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					params := map[string]string{
						"method": "cell.content.get",
						"row":    fmt.Sprintf("%d", cmd.Int("row")),
						"column": fmt.Sprintf("%d", cmd.Int("column")),
					}
					if v := cmd.String("worksheet"); v != "" {
						params["worksheet_name"] = v
					}
					raw, err := c.Request(ctx, "POST", c.SheetBase+"/"+cmd.String("workbook"), &zohttp.RequestOpts{Params: params})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:  "get-range",
				Usage: "Get content of range",
				Flags: []cli.Flag{
					workbookFlag,
					worksheetFlag,
					&cli.IntFlag{Name: "start-row", Required: true, Usage: "Start row index"},
					&cli.IntFlag{Name: "start-column", Required: true, Usage: "Start column index"},
					&cli.IntFlag{Name: "end-row", Required: true, Usage: "End row index"},
					&cli.IntFlag{Name: "end-column", Required: true, Usage: "End column index"},
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					params := map[string]string{
						"method":       "range.content.get",
						"start_row":    fmt.Sprintf("%d", cmd.Int("start-row")),
						"start_column": fmt.Sprintf("%d", cmd.Int("start-column")),
						"end_row":      fmt.Sprintf("%d", cmd.Int("end-row")),
						"end_column":   fmt.Sprintf("%d", cmd.Int("end-column")),
					}
					if v := cmd.String("worksheet"); v != "" {
						params["worksheet_name"] = v
					}
					raw, err := c.Request(ctx, "POST", c.SheetBase+"/"+cmd.String("workbook"), &zohttp.RequestOpts{Params: params})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:  "get-named-range",
				Usage: "Get content of named range",
				Flags: []cli.Flag{
					workbookFlag,
					&cli.StringFlag{Name: "named-range", Required: true, Usage: "Named range"},
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					params := map[string]string{
						"method":        "namedrange.content.get",
						"name_of_range": cmd.String("named-range"),
					}
					raw, err := c.Request(ctx, "POST", c.SheetBase+"/"+cmd.String("workbook"), &zohttp.RequestOpts{Params: params})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:  "get-worksheet",
				Usage: "Get content of worksheet area",
				Flags: []cli.Flag{
					workbookFlag,
					worksheetFlag,
					&cli.IntFlag{Name: "start-row", Usage: "Start row"},
					&cli.IntFlag{Name: "start-column", Usage: "Start column"},
					&cli.IntFlag{Name: "end-row", Usage: "End row"},
					&cli.IntFlag{Name: "end-column", Usage: "End column"},
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					params := map[string]string{"method": "worksheet.content.get"}
					if v := cmd.String("worksheet"); v != "" {
						params["worksheet_name"] = v
					}
					if v := cmd.Int("start-row"); v > 0 {
						params["start_row"] = fmt.Sprintf("%d", v)
					}
					if v := cmd.Int("start-column"); v > 0 {
						params["start_column"] = fmt.Sprintf("%d", v)
					}
					if v := cmd.Int("end-row"); v > 0 {
						params["end_row"] = fmt.Sprintf("%d", v)
					}
					if v := cmd.Int("end-column"); v > 0 {
						params["end_column"] = fmt.Sprintf("%d", v)
					}
					raw, err := c.Request(ctx, "POST", c.SheetBase+"/"+cmd.String("workbook"), &zohttp.RequestOpts{Params: params})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:  "get-used-area",
				Usage: "Get used area of worksheet",
				Flags: []cli.Flag{
					workbookFlag,
					worksheetFlag,
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					params := map[string]string{"method": "worksheet.usedarea"}
					if v := cmd.String("worksheet"); v != "" {
						params["worksheet_name"] = v
					}
					raw, err := c.Request(ctx, "POST", c.SheetBase+"/"+cmd.String("workbook"), &zohttp.RequestOpts{Params: params})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:  "set",
				Usage: "Set content to cell",
				Flags: []cli.Flag{
					workbookFlag,
					worksheetFlag,
					&cli.IntFlag{Name: "row", Required: true, Usage: "Row index"},
					&cli.IntFlag{Name: "column", Required: true, Usage: "Column index"},
					&cli.StringFlag{Name: "content", Required: true, Usage: "Cell content"},
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					params := map[string]string{
						"method":  "cell.content.set",
						"row":     fmt.Sprintf("%d", cmd.Int("row")),
						"column":  fmt.Sprintf("%d", cmd.Int("column")),
						"content": cmd.String("content"),
					}
					if v := cmd.String("worksheet"); v != "" {
						params["worksheet_name"] = v
					}
					raw, err := c.Request(ctx, "POST", c.SheetBase+"/"+cmd.String("workbook"), &zohttp.RequestOpts{Params: params})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:  "set-multiple",
				Usage: "Set content to multiple cells",
				Flags: []cli.Flag{
					workbookFlag,
					worksheetFlag,
					&cli.StringFlag{Name: "data", Required: true, Usage: "Cell data JSON"},
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					params := map[string]string{
						"method": "cells.content.set",
						"data":   cmd.String("data"),
					}
					if v := cmd.String("worksheet"); v != "" {
						params["worksheet_name"] = v
					}
					raw, err := c.Request(ctx, "POST", c.SheetBase+"/"+cmd.String("workbook"), &zohttp.RequestOpts{Params: params})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:  "set-row",
				Usage: "Set content to row",
				Flags: []cli.Flag{
					workbookFlag,
					worksheetFlag,
					&cli.IntFlag{Name: "row", Required: true, Usage: "Row number"},
					&cli.StringFlag{Name: "column-array", Required: true, Usage: "Column array JSON"},
					&cli.StringFlag{Name: "data-array", Required: true, Usage: "Data array JSON"},
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					params := map[string]string{
						"method":       "row.content.set",
						"row":          fmt.Sprintf("%d", cmd.Int("row")),
						"column_array": cmd.String("column-array"),
						"data_array":   cmd.String("data-array"),
					}
					if v := cmd.String("worksheet"); v != "" {
						params["worksheet_name"] = v
					}
					raw, err := c.Request(ctx, "POST", c.SheetBase+"/"+cmd.String("workbook"), &zohttp.RequestOpts{Params: params})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:  "set-range",
				Usage: "Set content to range",
				Flags: []cli.Flag{
					workbookFlag,
					worksheetFlag,
					&cli.IntFlag{Name: "row", Required: true, Usage: "Row index"},
					&cli.IntFlag{Name: "column", Required: true, Usage: "Column index"},
					&cli.StringFlag{Name: "data", Required: true, Usage: "Data JSON 2D array"},
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					params := map[string]string{
						"method": "worksheet.csvdata.set",
						"row":    fmt.Sprintf("%d", cmd.Int("row")),
						"column": fmt.Sprintf("%d", cmd.Int("column")),
						"data":   cmd.String("data"),
					}
					if v := cmd.String("worksheet"); v != "" {
						params["worksheet_name"] = v
					}
					raw, err := c.Request(ctx, "POST", c.SheetBase+"/"+cmd.String("workbook"), &zohttp.RequestOpts{Params: params})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
		},
	}
}

func contentCmd() *cli.Command {
	return &cli.Command{
		Name:  "content",
		Usage: "Content operations",
		Commands: []*cli.Command{
			{
				Name:  "append-csv",
				Usage: "Append rows with CSV data",
				Flags: []cli.Flag{
					workbookFlag,
					worksheetFlag,
					&cli.StringFlag{Name: "data", Required: true, Usage: "CSV data string"},
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					params := map[string]string{
						"method": "worksheet.csvdata.append",
						"data":   cmd.String("data"),
					}
					if v := cmd.String("worksheet"); v != "" {
						params["worksheet_name"] = v
					}
					raw, err := c.Request(ctx, "POST", c.SheetBase+"/"+cmd.String("workbook"), &zohttp.RequestOpts{Params: params})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:  "append-json",
				Usage: "Append rows with JSON data",
				Flags: []cli.Flag{
					workbookFlag,
					worksheetFlag,
					&cli.IntFlag{Name: "header-row", Usage: "Header row number"},
					&cli.StringFlag{Name: "data", Required: true, Usage: "Row data as JSON array"},
					&cli.StringFlag{Name: "json", Usage: "Additional fields as JSON"},
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					params := map[string]string{
						"method":    "worksheet.jsondata.append",
						"json_data": cmd.String("data"),
					}
					if v := cmd.String("worksheet"); v != "" {
						params["worksheet_name"] = v
					}
					if v := cmd.Int("header-row"); v > 0 {
						params["header_row"] = fmt.Sprintf("%d", v)
					}
					if err := internal.MergeJSONForm(cmd, params); err != nil {
						return err
					}
					raw, err := c.Request(ctx, "POST", c.SheetBase+"/"+cmd.String("workbook"), &zohttp.RequestOpts{Params: params})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:  "update-json",
				Usage: "Update rows with JSON data",
				Flags: []cli.Flag{
					workbookFlag,
					worksheetFlag,
					&cli.IntFlag{Name: "header-row", Usage: "Header row number"},
					&cli.StringFlag{Name: "criteria", Required: true, Usage: "Filter criteria"},
					&cli.StringFlag{Name: "data", Required: true, Usage: "Update data as JSON"},
					&cli.StringFlag{Name: "json", Usage: "Additional fields as JSON"},
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					params := map[string]string{
						"method":    "worksheet.jsondata.set",
						"criteria":  cmd.String("criteria"),
						"json_data": cmd.String("data"),
					}
					if v := cmd.String("worksheet"); v != "" {
						params["worksheet_name"] = v
					}
					if v := cmd.Int("header-row"); v > 0 {
						params["header_row"] = fmt.Sprintf("%d", v)
					}
					if err := internal.MergeJSONForm(cmd, params); err != nil {
						return err
					}
					raw, err := c.Request(ctx, "POST", c.SheetBase+"/"+cmd.String("workbook"), &zohttp.RequestOpts{Params: params})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:  "insert-json",
				Usage: "Insert row with JSON data",
				Flags: []cli.Flag{
					workbookFlag,
					worksheetFlag,
					&cli.IntFlag{Name: "header-row", Usage: "Header row number"},
					&cli.IntFlag{Name: "row-index", Required: true, Usage: "Row index to insert at"},
					&cli.StringFlag{Name: "data", Required: true, Usage: "Row data as JSON array"},
					&cli.StringFlag{Name: "json", Usage: "Additional fields as JSON"},
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					params := map[string]string{
						"method":    "worksheet.jsondata.insert",
						"row_index": fmt.Sprintf("%d", cmd.Int("row-index")),
						"json_data": cmd.String("data"),
					}
					if v := cmd.String("worksheet"); v != "" {
						params["worksheet_name"] = v
					}
					if v := cmd.Int("header-row"); v > 0 {
						params["header_row"] = fmt.Sprintf("%d", v)
					}
					if err := internal.MergeJSONForm(cmd, params); err != nil {
						return err
					}
					raw, err := c.Request(ctx, "POST", c.SheetBase+"/"+cmd.String("workbook"), &zohttp.RequestOpts{Params: params})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:  "clear-contents",
				Usage: "Clear contents of range",
				Flags: []cli.Flag{
					workbookFlag,
					worksheetFlag,
					&cli.IntFlag{Name: "start-row", Required: true, Usage: "Start row index"},
					&cli.IntFlag{Name: "start-column", Required: true, Usage: "Start column index"},
					&cli.IntFlag{Name: "end-row", Required: true, Usage: "End row index"},
					&cli.IntFlag{Name: "end-column", Required: true, Usage: "End column index"},
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					params := map[string]string{
						"method":       "range.content.clear",
						"start_row":    fmt.Sprintf("%d", cmd.Int("start-row")),
						"start_column": fmt.Sprintf("%d", cmd.Int("start-column")),
						"end_row":      fmt.Sprintf("%d", cmd.Int("end-row")),
						"end_column":   fmt.Sprintf("%d", cmd.Int("end-column")),
					}
					if v := cmd.String("worksheet"); v != "" {
						params["worksheet_name"] = v
					}
					raw, err := c.Request(ctx, "POST", c.SheetBase+"/"+cmd.String("workbook"), &zohttp.RequestOpts{Params: params})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:  "clear-range",
				Usage: "Clear range",
				Flags: []cli.Flag{
					workbookFlag,
					worksheetFlag,
					&cli.IntFlag{Name: "start-row", Required: true, Usage: "Start row index"},
					&cli.IntFlag{Name: "start-column", Required: true, Usage: "Start column index"},
					&cli.IntFlag{Name: "end-row", Required: true, Usage: "End row index"},
					&cli.IntFlag{Name: "end-column", Required: true, Usage: "End column index"},
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					params := map[string]string{
						"method":       "range.clear",
						"start_row":    fmt.Sprintf("%d", cmd.Int("start-row")),
						"start_column": fmt.Sprintf("%d", cmd.Int("start-column")),
						"end_row":      fmt.Sprintf("%d", cmd.Int("end-row")),
						"end_column":   fmt.Sprintf("%d", cmd.Int("end-column")),
					}
					if v := cmd.String("worksheet"); v != "" {
						params["worksheet_name"] = v
					}
					raw, err := c.Request(ctx, "POST", c.SheetBase+"/"+cmd.String("workbook"), &zohttp.RequestOpts{Params: params})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:  "clear-filters",
				Usage: "Clear filters",
				Flags: []cli.Flag{
					workbookFlag,
					worksheetFlag,
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					params := map[string]string{"method": "worksheet.filter.clear"}
					if v := cmd.String("worksheet"); v != "" {
						params["worksheet_name"] = v
					}
					raw, err := c.Request(ctx, "POST", c.SheetBase+"/"+cmd.String("workbook"), &zohttp.RequestOpts{Params: params})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:  "find",
				Usage: "Find content",
				Flags: []cli.Flag{
					workbookFlag,
					worksheetFlag,
					&cli.StringFlag{Name: "search", Required: true, Usage: "Value to search for"},
					&cli.StringFlag{Name: "scope", Required: true, Usage: "Search scope"},
					&cli.BoolFlag{Name: "is-case-sensitive", Usage: "Case sensitive search"},
					&cli.BoolFlag{Name: "is-exact-match", Usage: "Exact match search"},
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					params := map[string]string{
						"method": "find",
						"search": cmd.String("search"),
						"scope":  cmd.String("scope"),
					}
					if v := cmd.String("worksheet"); v != "" {
						params["worksheet_name"] = v
					}
					if cmd.IsSet("is-case-sensitive") {
						params["is_case_sensitive"] = fmt.Sprintf("%t", cmd.Bool("is-case-sensitive"))
					}
					if cmd.IsSet("is-exact-match") {
						params["is_exact_match"] = fmt.Sprintf("%t", cmd.Bool("is-exact-match"))
					}
					raw, err := c.Request(ctx, "POST", c.SheetBase+"/"+cmd.String("workbook"), &zohttp.RequestOpts{Params: params})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:  "find-replace",
				Usage: "Find and replace content",
				Flags: []cli.Flag{
					workbookFlag,
					worksheetFlag,
					&cli.StringFlag{Name: "search", Required: true, Usage: "Value to search for"},
					&cli.StringFlag{Name: "replace-with", Required: true, Usage: "Replacement value"},
					&cli.StringFlag{Name: "scope", Required: true, Usage: "Search scope"},
					&cli.BoolFlag{Name: "is-case-sensitive", Usage: "Case sensitive search"},
					&cli.BoolFlag{Name: "is-exact-match", Usage: "Exact match search"},
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					params := map[string]string{
						"method":       "replace",
						"search":       cmd.String("search"),
						"replace_with": cmd.String("replace-with"),
						"scope":        cmd.String("scope"),
					}
					if v := cmd.String("worksheet"); v != "" {
						params["worksheet_name"] = v
					}
					if cmd.IsSet("is-case-sensitive") {
						params["is_case_sensitive"] = fmt.Sprintf("%t", cmd.Bool("is-case-sensitive"))
					}
					if cmd.IsSet("is-exact-match") {
						params["is_exact_match"] = fmt.Sprintf("%t", cmd.Bool("is-exact-match"))
					}
					raw, err := c.Request(ctx, "POST", c.SheetBase+"/"+cmd.String("workbook"), &zohttp.RequestOpts{Params: params})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:  "recalculate",
				Usage: "Recalculate workbook",
				Flags: []cli.Flag{
					workbookFlag,
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					params := map[string]string{"method": "workbook.recalculate"}
					raw, err := c.Request(ctx, "POST", c.SheetBase+"/"+cmd.String("workbook"), &zohttp.RequestOpts{Params: params})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
		},
	}
}

func formatCmd() *cli.Command {
	return &cli.Command{
		Name:  "format",
		Usage: "Formatting and structure operations",
		Commands: []*cli.Command{
			{
				Name:  "ranges",
				Usage: "Format ranges",
				Flags: []cli.Flag{
					workbookFlag,
					&cli.StringFlag{Name: "format-json", Required: true, Usage: "Format JSON"},
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					params := map[string]string{
						"method":      "ranges.format.set",
						"format_json": cmd.String("format-json"),
					}
					raw, err := c.Request(ctx, "POST", c.SheetBase+"/"+cmd.String("workbook"), &zohttp.RequestOpts{Params: params})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:  "image-fit",
				Usage: "Image fit options",
				Flags: []cli.Flag{
					workbookFlag,
					&cli.StringFlag{Name: "image-json", Required: true, Usage: "Image JSON"},
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					params := map[string]string{
						"method":     "range.images.fit",
						"image_json": cmd.String("image-json"),
					}
					raw, err := c.Request(ctx, "POST", c.SheetBase+"/"+cmd.String("workbook"), &zohttp.RequestOpts{Params: params})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:  "row-height",
				Usage: "Set row height",
				Flags: []cli.Flag{
					workbookFlag,
					worksheetFlag,
					&cli.StringFlag{Name: "row-index-array", Required: true, Usage: "JSON array of row indices"},
					&cli.IntFlag{Name: "row-height", Required: true, Usage: "Height in pixels"},
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					params := map[string]string{
						"method":          "worksheet.rows.height",
						"row_index_array": cmd.String("row-index-array"),
						"row_height":      fmt.Sprintf("%d", cmd.Int("row-height")),
					}
					if v := cmd.String("worksheet"); v != "" {
						params["worksheet_name"] = v
					}
					raw, err := c.Request(ctx, "POST", c.SheetBase+"/"+cmd.String("workbook"), &zohttp.RequestOpts{Params: params})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:  "column-width",
				Usage: "Set column width",
				Flags: []cli.Flag{
					workbookFlag,
					worksheetFlag,
					&cli.StringFlag{Name: "column-index-array", Required: true, Usage: "JSON array of column indices"},
					&cli.IntFlag{Name: "column-width", Required: true, Usage: "Width in pixels"},
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					params := map[string]string{
						"method":             "worksheet.columns.width",
						"column_index_array": cmd.String("column-index-array"),
						"column_width":       fmt.Sprintf("%d", cmd.Int("column-width")),
					}
					if v := cmd.String("worksheet"); v != "" {
						params["worksheet_name"] = v
					}
					raw, err := c.Request(ctx, "POST", c.SheetBase+"/"+cmd.String("workbook"), &zohttp.RequestOpts{Params: params})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:  "insert-row",
				Usage: "Insert row",
				Flags: []cli.Flag{
					workbookFlag,
					worksheetFlag,
					&cli.IntFlag{Name: "row", Required: true, Usage: "Row number"},
					&cli.IntFlag{Name: "count", Usage: "Number of rows to insert"},
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					params := map[string]string{
						"method": "row.insert",
						"row":    fmt.Sprintf("%d", cmd.Int("row")),
					}
					if v := cmd.String("worksheet"); v != "" {
						params["worksheet_name"] = v
					}
					if v := cmd.Int("count"); v > 0 {
						params["count"] = fmt.Sprintf("%d", v)
					}
					raw, err := c.Request(ctx, "POST", c.SheetBase+"/"+cmd.String("workbook"), &zohttp.RequestOpts{Params: params})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:  "insert-column",
				Usage: "Insert column",
				Flags: []cli.Flag{
					workbookFlag,
					worksheetFlag,
					&cli.IntFlag{Name: "column", Required: true, Usage: "Column number"},
					&cli.IntFlag{Name: "count", Usage: "Number of columns to insert"},
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					params := map[string]string{
						"method": "column.insert",
						"column": fmt.Sprintf("%d", cmd.Int("column")),
					}
					if v := cmd.String("worksheet"); v != "" {
						params["worksheet_name"] = v
					}
					if v := cmd.Int("count"); v > 0 {
						params["count"] = fmt.Sprintf("%d", v)
					}
					raw, err := c.Request(ctx, "POST", c.SheetBase+"/"+cmd.String("workbook"), &zohttp.RequestOpts{Params: params})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:  "delete-row",
				Usage: "Delete row",
				Flags: []cli.Flag{
					workbookFlag,
					worksheetFlag,
					&cli.IntFlag{Name: "row", Required: true, Usage: "Row number"},
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					params := map[string]string{
						"method": "row.delete",
						"row":    fmt.Sprintf("%d", cmd.Int("row")),
					}
					if v := cmd.String("worksheet"); v != "" {
						params["worksheet_name"] = v
					}
					raw, err := c.Request(ctx, "POST", c.SheetBase+"/"+cmd.String("workbook"), &zohttp.RequestOpts{Params: params})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:  "delete-rows",
				Usage: "Delete multiple rows",
				Flags: []cli.Flag{
					workbookFlag,
					worksheetFlag,
					&cli.StringFlag{Name: "row-index-array", Required: true, Usage: "JSON array of row indices"},
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					params := map[string]string{
						"method":          "worksheet.rows.delete",
						"row_index_array": cmd.String("row-index-array"),
					}
					if v := cmd.String("worksheet"); v != "" {
						params["worksheet_name"] = v
					}
					raw, err := c.Request(ctx, "POST", c.SheetBase+"/"+cmd.String("workbook"), &zohttp.RequestOpts{Params: params})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:  "delete-column",
				Usage: "Delete column",
				Flags: []cli.Flag{
					workbookFlag,
					worksheetFlag,
					&cli.IntFlag{Name: "column", Required: true, Usage: "Column number"},
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					params := map[string]string{
						"method": "column.delete",
						"column": fmt.Sprintf("%d", cmd.Int("column")),
					}
					if v := cmd.String("worksheet"); v != "" {
						params["worksheet_name"] = v
					}
					raw, err := c.Request(ctx, "POST", c.SheetBase+"/"+cmd.String("workbook"), &zohttp.RequestOpts{Params: params})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:  "set-note",
				Usage: "Set note to cell",
				Flags: []cli.Flag{
					workbookFlag,
					worksheetFlag,
					&cli.IntFlag{Name: "row", Required: true, Usage: "Row index"},
					&cli.IntFlag{Name: "column", Required: true, Usage: "Column index"},
					&cli.StringFlag{Name: "note", Required: true, Usage: "Note text"},
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					params := map[string]string{
						"method": "cell.note.set",
						"row":    fmt.Sprintf("%d", cmd.Int("row")),
						"column": fmt.Sprintf("%d", cmd.Int("column")),
						"note":   cmd.String("note"),
					}
					if v := cmd.String("worksheet"); v != "" {
						params["worksheet_name"] = v
					}
					raw, err := c.Request(ctx, "POST", c.SheetBase+"/"+cmd.String("workbook"), &zohttp.RequestOpts{Params: params})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
		},
	}
}

func namedRangesCmd() *cli.Command {
	return &cli.Command{
		Name:  "named-ranges",
		Usage: "Named range operations",
		Commands: []*cli.Command{
			{
				Name:  "list",
				Usage: "List all named ranges",
				Flags: []cli.Flag{
					workbookFlag,
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					params := map[string]string{"method": "namedrange.list"}
					raw, err := c.Request(ctx, "POST", c.SheetBase+"/"+cmd.String("workbook"), &zohttp.RequestOpts{Params: params})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:  "create",
				Usage: "Create named range",
				Flags: []cli.Flag{
					workbookFlag,
					worksheetFlag,
					&cli.StringFlag{Name: "name", Required: true, Usage: "Named range name"},
					&cli.StringFlag{Name: "range", Required: true, Usage: "Range reference"},
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					params := map[string]string{
						"method":        "namedrange.create",
						"name_of_range": cmd.String("name"),
						"range":         cmd.String("range"),
					}
					if v := cmd.String("worksheet"); v != "" {
						params["worksheet_name"] = v
					}
					raw, err := c.Request(ctx, "POST", c.SheetBase+"/"+cmd.String("workbook"), &zohttp.RequestOpts{Params: params})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:  "update",
				Usage: "Update named range",
				Flags: []cli.Flag{
					workbookFlag,
					worksheetFlag,
					&cli.StringFlag{Name: "name", Required: true, Usage: "Named range name"},
					&cli.StringFlag{Name: "range", Required: true, Usage: "Range reference"},
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					params := map[string]string{
						"method":        "namedrange.update",
						"name_of_range": cmd.String("name"),
						"range":         cmd.String("range"),
					}
					if v := cmd.String("worksheet"); v != "" {
						params["worksheet_name"] = v
					}
					raw, err := c.Request(ctx, "POST", c.SheetBase+"/"+cmd.String("workbook"), &zohttp.RequestOpts{Params: params})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:  "delete",
				Usage: "Delete named range",
				Flags: []cli.Flag{
					workbookFlag,
					&cli.StringFlag{Name: "name", Required: true, Usage: "Named range name"},
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					params := map[string]string{
						"method":        "namedrange.delete",
						"name_of_range": cmd.String("name"),
					}
					raw, err := c.Request(ctx, "POST", c.SheetBase+"/"+cmd.String("workbook"), &zohttp.RequestOpts{Params: params})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
		},
	}
}

func mergeCmd() *cli.Command {
	return &cli.Command{
		Name:  "merge",
		Usage: "Merge template operations",
		Commands: []*cli.Command{
			{
				Name:  "templates",
				Usage: "Get merge templates",
				Action: func(ctx context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					params := map[string]string{"method": "mergetemplate.list"}
					raw, err := c.Request(ctx, "POST", c.SheetBase+"/mergetemplates", &zohttp.RequestOpts{Params: params})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:  "fields",
				Usage: "Get merge fields",
				Flags: []cli.Flag{
					workbookFlag,
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					params := map[string]string{"method": "workbook.mergefield.list"}
					raw, err := c.Request(ctx, "POST", c.SheetBase+"/"+cmd.String("workbook"), &zohttp.RequestOpts{Params: params})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:  "jobs",
				Usage: "Get merge jobs",
				Flags: []cli.Flag{
					workbookFlag,
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					params := map[string]string{"method": "workbook.mergejob.list"}
					raw, err := c.Request(ctx, "POST", c.SheetBase+"/"+cmd.String("workbook"), &zohttp.RequestOpts{Params: params})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:  "job-detail",
				Usage: "Get merge job details",
				Flags: []cli.Flag{
					workbookFlag,
					&cli.StringFlag{Name: "job-id", Required: true, Usage: "Job ID"},
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					params := map[string]string{
						"method": "workbook.mergejob.details",
						"job_id": cmd.String("job-id"),
					}
					raw, err := c.Request(ctx, "POST", c.SheetBase+"/"+cmd.String("workbook"), &zohttp.RequestOpts{Params: params})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:  "save",
				Usage: "Merge and save",
				Flags: []cli.Flag{
					workbookFlag,
					&cli.StringFlag{Name: "merge-data", Required: true, Usage: "Merge data JSON"},
					&cli.StringFlag{Name: "output-settings", Required: true, Usage: "Output settings JSON"},
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					params := map[string]string{
						"method":          "merge.save",
						"resource_id":     cmd.String("workbook"),
						"merge_data":      cmd.String("merge-data"),
						"output_settings": cmd.String("output-settings"),
					}
					raw, err := c.Request(ctx, "POST", c.SheetBase+"/merge", &zohttp.RequestOpts{Params: params})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:  "email",
				Usage: "Merge and email",
				Flags: []cli.Flag{
					workbookFlag,
					&cli.StringFlag{Name: "merge-data", Required: true, Usage: "Merge data JSON"},
					&cli.StringFlag{Name: "email-settings", Required: true, Usage: "Email settings JSON"},
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					params := map[string]string{
						"method":         "merge.email.attachment",
						"resource_id":    cmd.String("workbook"),
						"merge_data":     cmd.String("merge-data"),
						"email_settings": cmd.String("email-settings"),
					}
					raw, err := c.Request(ctx, "POST", c.SheetBase+"/merge", &zohttp.RequestOpts{Params: params})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
		},
	}
}

func premiumCmd() *cli.Command {
	return &cli.Command{
		Name:  "premium",
		Usage: "Premium API operations",
		Commands: []*cli.Command{
			{
				Name:  "fetch-records",
				Usage: "Fetch records (premium)",
				Flags: []cli.Flag{
					workbookFlag,
					worksheetFlag,
					&cli.StringFlag{Name: "criteria", Usage: "Filter criteria"},
					&cli.IntFlag{Name: "header-row", Usage: "Header row number"},
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					params := map[string]string{"method": "premium.records.fetch", "resource_id": cmd.String("workbook")}
					if v := cmd.String("worksheet"); v != "" {
						params["worksheet_name"] = v
					}
					if v := cmd.String("criteria"); v != "" {
						params["criteria"] = v
					}
					if v := cmd.Int("header-row"); v > 0 {
						params["header_row"] = fmt.Sprintf("%d", v)
					}
					raw, err := c.Request(ctx, "POST", c.SheetBase+"/fetchrecords", &zohttp.RequestOpts{Params: params})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:  "add-records",
				Usage: "Add records (premium)",
				Flags: []cli.Flag{
					workbookFlag,
					worksheetFlag,
					&cli.IntFlag{Name: "header-row", Usage: "Header row number"},
					&cli.StringFlag{Name: "data", Required: true, Usage: "Record data as JSON array"},
					&cli.StringFlag{Name: "json", Usage: "Additional fields as JSON"},
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					params := map[string]string{
						"method":      "premium.records.add",
						"resource_id": cmd.String("workbook"),
						"json_data":   cmd.String("data"),
					}
					if v := cmd.String("worksheet"); v != "" {
						params["worksheet_name"] = v
					}
					if v := cmd.Int("header-row"); v > 0 {
						params["header_row"] = fmt.Sprintf("%d", v)
					}
					if err := internal.MergeJSONForm(cmd, params); err != nil {
						return err
					}
					raw, err := c.Request(ctx, "POST", c.SheetBase+"/addrecords", &zohttp.RequestOpts{Params: params})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:  "update-records",
				Usage: "Update records (premium)",
				Flags: []cli.Flag{
					workbookFlag,
					worksheetFlag,
					&cli.IntFlag{Name: "header-row", Usage: "Header row number"},
					&cli.StringFlag{Name: "criteria", Required: true, Usage: "Filter criteria"},
					&cli.StringFlag{Name: "data", Required: true, Usage: "Update data as JSON"},
					&cli.StringFlag{Name: "json", Usage: "Additional fields as JSON"},
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					params := map[string]string{
						"method":      "premium.records.update",
						"resource_id": cmd.String("workbook"),
						"criteria":    cmd.String("criteria"),
						"json_data":   cmd.String("data"),
					}
					if v := cmd.String("worksheet"); v != "" {
						params["worksheet_name"] = v
					}
					if v := cmd.Int("header-row"); v > 0 {
						params["header_row"] = fmt.Sprintf("%d", v)
					}
					if err := internal.MergeJSONForm(cmd, params); err != nil {
						return err
					}
					raw, err := c.Request(ctx, "POST", c.SheetBase+"/updaterecords", &zohttp.RequestOpts{Params: params})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
		},
	}
}

func utilityCmd() *cli.Command {
	return &cli.Command{
		Name:  "utility",
		Usage: "Utility operations",
		Commands: []*cli.Command{
			{
				Name:  "range-to-index",
				Usage: "Convert range to index",
				Flags: []cli.Flag{
					workbookFlag,
					&cli.StringFlag{Name: "range", Required: true, Usage: "Range reference"},
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					params := map[string]string{
						"method":        "range.index.get",
						"range_address": cmd.String("range"),
					}
					raw, err := c.Request(ctx, "POST", c.SheetBase+"/utils", &zohttp.RequestOpts{Params: params})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:  "index-to-range",
				Usage: "Convert index to range",
				Flags: []cli.Flag{
					workbookFlag,
					&cli.IntFlag{Name: "start-row", Required: true, Usage: "Start row index"},
					&cli.IntFlag{Name: "start-column", Required: true, Usage: "Start column index"},
					&cli.IntFlag{Name: "end-row", Usage: "End row index"},
					&cli.IntFlag{Name: "end-column", Usage: "End column index"},
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					params := map[string]string{
						"method":       "range.address.get",
						"start_row":    fmt.Sprintf("%d", cmd.Int("start-row")),
						"start_column": fmt.Sprintf("%d", cmd.Int("start-column")),
					}
					if v := cmd.Int("end-row"); v > 0 {
						params["end_row"] = fmt.Sprintf("%d", v)
					}
					if v := cmd.Int("end-column"); v > 0 {
						params["end_column"] = fmt.Sprintf("%d", v)
					}
					raw, err := c.Request(ctx, "POST", c.SheetBase+"/utils", &zohttp.RequestOpts{Params: params})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
		},
	}
}
