package recruit

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/omin8tor/zoho-cli/internal"
	zohttp "github.com/omin8tor/zoho-cli/internal/http"
	"github.com/omin8tor/zoho-cli/internal/output"
	"github.com/urfave/cli/v3"
)

func Commands() *cli.Command {
	return &cli.Command{
		Name:  "recruit",
		Usage: "Zoho Recruit operations",
		Commands: []*cli.Command{
			modulesCmd(),
			recordsCmd(),
			notesCmd(),
			attachmentsCmd(),
			tagsCmd(),
			associateCmd(),
			relatedCmd(),
			usersCmd(),
		},
	}
}

func modulesCmd() *cli.Command {
	return &cli.Command{
		Name:  "modules",
		Usage: "Recruit module metadata",
		Commands: []*cli.Command{
			{
				Name:  "list",
				Usage: "List available modules",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					raw, err := c.Request("GET", c.RecruitBase+"/settings/modules", nil)
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "fields",
				Usage:     "List fields for a module",
				ArgsUsage: "<module>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					if cmd.Args().Len() < 1 {
						return internal.NewValidationError("module name required")
					}
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					raw, err := c.Request("GET", c.RecruitBase+"/settings/fields", &zohttp.RequestOpts{
						Params: map[string]string{"module": cmd.Args().First()},
					})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "layouts",
				Usage:     "List layouts for a module",
				ArgsUsage: "<module>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					if cmd.Args().Len() < 1 {
						return internal.NewValidationError("module name required")
					}
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					raw, err := c.Request("GET", c.RecruitBase+"/settings/layouts", &zohttp.RequestOpts{
						Params: map[string]string{"module": cmd.Args().First()},
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

func recordsCmd() *cli.Command {
	return &cli.Command{
		Name:  "records",
		Usage: "Recruit record operations",
		Commands: []*cli.Command{
			{
				Name:      "list",
				Usage:     "List records in a module",
				ArgsUsage: "<module>",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "fields", Usage: "Comma-separated field API names"},
					&cli.StringFlag{Name: "sort-by", Usage: "Field to sort by"},
					&cli.StringFlag{Name: "sort-order", Usage: "asc or desc"},
					&cli.IntFlag{Name: "page", Value: 1, Usage: "Page number"},
					&cli.IntFlag{Name: "per-page", Value: 200, Usage: "Records per page"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					if cmd.Args().Len() < 1 {
						return internal.NewValidationError("module name required")
					}
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					params := map[string]string{
						"page":     fmt.Sprintf("%d", cmd.Int("page")),
						"per_page": fmt.Sprintf("%d", cmd.Int("per-page")),
					}
					if f := cmd.String("fields"); f != "" {
						params["fields"] = f
					}
					if s := cmd.String("sort-by"); s != "" {
						params["sort_by"] = s
					}
					if s := cmd.String("sort-order"); s != "" {
						params["sort_order"] = s
					}
					raw, err := c.Request("GET", c.RecruitBase+"/"+cmd.Args().First(), &zohttp.RequestOpts{Params: params})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "get",
				Usage:     "Get a single record",
				ArgsUsage: "<module> <record-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					if cmd.Args().Len() < 2 {
						return internal.NewValidationError("module name and record ID required")
					}
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					raw, err := c.Request("GET", c.RecruitBase+"/"+cmd.Args().Get(0)+"/"+cmd.Args().Get(1), nil)
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "create",
				Usage:     "Create a record",
				ArgsUsage: "<module>",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "json", Required: true, Usage: "Record data as JSON"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					if cmd.Args().Len() < 1 {
						return internal.NewValidationError("module name required")
					}
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					var parsed any
					if err := json.Unmarshal([]byte(cmd.String("json")), &parsed); err != nil {
						return internal.NewValidationError(fmt.Sprintf("invalid JSON: %v", err))
					}
					body := map[string]any{"data": []any{parsed}}
					raw, err := c.Request("POST", c.RecruitBase+"/"+cmd.Args().First(), &zohttp.RequestOpts{JSON: body})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "update",
				Usage:     "Update a record",
				ArgsUsage: "<module> <record-id>",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "json", Required: true, Usage: "Fields to update as JSON"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					if cmd.Args().Len() < 2 {
						return internal.NewValidationError("module name and record ID required")
					}
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					var parsed any
					if err := json.Unmarshal([]byte(cmd.String("json")), &parsed); err != nil {
						return internal.NewValidationError(fmt.Sprintf("invalid JSON: %v", err))
					}
					body := map[string]any{"data": []any{parsed}}
					raw, err := c.Request("PUT", c.RecruitBase+"/"+cmd.Args().Get(0)+"/"+cmd.Args().Get(1), &zohttp.RequestOpts{JSON: body})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "delete",
				Usage:     "Delete a record",
				ArgsUsage: "<module> <record-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					if cmd.Args().Len() < 2 {
						return internal.NewValidationError("module name and record ID required")
					}
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					module, recordID := cmd.Args().Get(0), cmd.Args().Get(1)
					raw, err := c.Request("DELETE", c.RecruitBase+"/"+module, &zohttp.RequestOpts{
						Params: map[string]string{"ids": recordID},
					})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "search",
				Usage:     "Search records in a module",
				ArgsUsage: "<module>",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "criteria", Usage: "Search criteria"},
					&cli.StringFlag{Name: "email", Usage: "Email search"},
					&cli.StringFlag{Name: "phone", Usage: "Phone search"},
					&cli.StringFlag{Name: "word", Usage: "Keyword search"},
					&cli.IntFlag{Name: "page", Value: 1},
					&cli.IntFlag{Name: "per-page", Value: 200},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					if cmd.Args().Len() < 1 {
						return internal.NewValidationError("module name required")
					}
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					params := map[string]string{
						"page":     fmt.Sprintf("%d", cmd.Int("page")),
						"per_page": fmt.Sprintf("%d", cmd.Int("per-page")),
					}
					if v := cmd.String("criteria"); v != "" {
						params["criteria"] = v
					} else if v := cmd.String("email"); v != "" {
						params["email"] = v
					} else if v := cmd.String("phone"); v != "" {
						params["phone"] = v
					} else if v := cmd.String("word"); v != "" {
						params["word"] = v
					}
					raw, err := c.Request("GET", c.RecruitBase+"/"+cmd.Args().First()+"/search", &zohttp.RequestOpts{Params: params})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
		},
	}
}

func notesCmd() *cli.Command {
	return &cli.Command{
		Name:  "notes",
		Usage: "Recruit record notes",
		Commands: []*cli.Command{
			{
				Name:      "list",
				Usage:     "List notes on a record",
				ArgsUsage: "<module> <record-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					if cmd.Args().Len() < 2 {
						return internal.NewValidationError("module name and record ID required")
					}
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					module, recordID := cmd.Args().Get(0), cmd.Args().Get(1)
					raw, err := c.Request("GET", c.RecruitBase+"/"+module+"/"+recordID+"/Notes", nil)
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:  "create",
				Usage: "Create a note",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "json", Required: true, Usage: "Note data as JSON"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					var parsed any
					if err := json.Unmarshal([]byte(cmd.String("json")), &parsed); err != nil {
						return internal.NewValidationError(fmt.Sprintf("invalid JSON: %v", err))
					}
					body := map[string]any{"data": []any{parsed}}
					raw, err := c.Request("POST", c.RecruitBase+"/Notes", &zohttp.RequestOpts{JSON: body})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "delete",
				Usage:     "Delete a note",
				ArgsUsage: "<note-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					if cmd.Args().Len() < 1 {
						return internal.NewValidationError("note ID required")
					}
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					raw, err := c.Request("DELETE", c.RecruitBase+"/Notes/"+cmd.Args().First(), nil)
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
		},
	}
}

func attachmentsCmd() *cli.Command {
	return &cli.Command{
		Name:  "attachments",
		Usage: "Recruit record attachments",
		Commands: []*cli.Command{
			{
				Name:      "list",
				Usage:     "List attachments on a record",
				ArgsUsage: "<module> <record-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					if cmd.Args().Len() < 2 {
						return internal.NewValidationError("module name and record ID required")
					}
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					module, recordID := cmd.Args().Get(0), cmd.Args().Get(1)
					raw, err := c.Request("GET", c.RecruitBase+"/"+module+"/"+recordID+"/Attachments", nil)
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "upload",
				Usage:     "Upload an attachment to a record",
				ArgsUsage: "<module> <record-id> <file-path>",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "category", Usage: "Attachment category (e.g. Resume, Cover Letter, Offer, Others)"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					if cmd.Args().Len() < 3 {
						return internal.NewValidationError("module name, record ID, and file path required")
					}
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					module, recordID, filePath := cmd.Args().Get(0), cmd.Args().Get(1), cmd.Args().Get(2)
					data, err := os.ReadFile(filePath)
					if err != nil {
						return fmt.Errorf("cannot read file: %w", err)
					}
					name := filePath
					if idx := strings.LastIndexAny(filePath, "/\\"); idx >= 0 {
						name = filePath[idx+1:]
					}
					params := map[string]string{}
					if v := cmd.String("category"); v != "" {
						params["attachments_category"] = v
					}
					raw, err := c.Request("POST", c.RecruitBase+"/"+module+"/"+recordID+"/Attachments", &zohttp.RequestOpts{
						Params: params,
						Files:  map[string]zohttp.FileUpload{"file": {Filename: name, Data: data}},
					})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "download",
				Usage:     "Download an attachment",
				ArgsUsage: "<module> <record-id> <attachment-id>",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "output", Usage: "Output file path"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					if cmd.Args().Len() < 3 {
						return internal.NewValidationError("module name, record ID, and attachment ID required")
					}
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					module, recordID, attID := cmd.Args().Get(0), cmd.Args().Get(1), cmd.Args().Get(2)
					url := c.RecruitBase + "/" + module + "/" + recordID + "/Attachments/" + attID
					body, _, _, err := c.RequestRaw("GET", url, nil)
					if err != nil {
						return err
					}
					if out := cmd.String("output"); out != "" {
						if err := os.WriteFile(out, body, 0600); err != nil {
							return err
						}
						return output.JSON(map[string]any{"ok": true, "path": out, "size": len(body)})
					}
					os.Stdout.Write(body)
					return nil
				},
			},
			{
				Name:      "delete",
				Usage:     "Delete an attachment",
				ArgsUsage: "<module> <record-id> <attachment-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					if cmd.Args().Len() < 3 {
						return internal.NewValidationError("module name, record ID, and attachment ID required")
					}
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					module, recordID, attID := cmd.Args().Get(0), cmd.Args().Get(1), cmd.Args().Get(2)
					raw, err := c.Request("DELETE", c.RecruitBase+"/"+module+"/"+recordID+"/Attachments/"+attID, nil)
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
		},
	}
}

func tagsCmd() *cli.Command {
	return &cli.Command{
		Name:  "tags",
		Usage: "Recruit tag operations",
		Commands: []*cli.Command{
			{
				Name:  "list",
				Usage: "List tags for a module",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "module", Required: true, Usage: "Module API name"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					raw, err := c.Request("GET", c.RecruitBase+"/settings/tags", &zohttp.RequestOpts{
						Params: map[string]string{"module": cmd.String("module")},
					})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:  "create",
				Usage: "Create tags for a module",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "module", Required: true, Usage: "Module API name"},
					&cli.StringFlag{Name: "json", Required: true, Usage: "Tags data as JSON"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					var parsed any
					if err := json.Unmarshal([]byte(cmd.String("json")), &parsed); err != nil {
						return internal.NewValidationError(fmt.Sprintf("invalid JSON: %v", err))
					}
					body := map[string]any{"tags": []any{parsed}}
					raw, err := c.Request("POST", c.RecruitBase+"/settings/tags", &zohttp.RequestOpts{
						Params: map[string]string{"module": cmd.String("module")},
						JSON:   body,
					})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "add",
				Usage:     "Add tags to records",
				ArgsUsage: "<module>",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "ids", Required: true, Usage: "Comma-separated record IDs"},
					&cli.StringFlag{Name: "tag-names", Required: true, Usage: "Comma-separated tag names"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					if cmd.Args().Len() < 1 {
						return internal.NewValidationError("module name required")
					}
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					raw, err := c.Request("POST", c.RecruitBase+"/"+cmd.Args().First()+"/actions/add_tags", &zohttp.RequestOpts{
						Params: map[string]string{
							"tag_names": cmd.String("tag-names"),
							"ids":       cmd.String("ids"),
						},
					})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "remove",
				Usage:     "Remove tags from records",
				ArgsUsage: "<module>",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "ids", Required: true, Usage: "Comma-separated record IDs"},
					&cli.StringFlag{Name: "tag-names", Required: true, Usage: "Comma-separated tag names"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					if cmd.Args().Len() < 1 {
						return internal.NewValidationError("module name required")
					}
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					raw, err := c.Request("POST", c.RecruitBase+"/"+cmd.Args().First()+"/actions/remove_tags", &zohttp.RequestOpts{
						Params: map[string]string{
							"tag_names": cmd.String("tag-names"),
							"ids":       cmd.String("ids"),
						},
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

func associateCmd() *cli.Command {
	return &cli.Command{
		Name:  "associate",
		Usage: "Recruit record association",
		Commands: []*cli.Command{
			{
				Name:      "list",
				Usage:     "List associated records",
				ArgsUsage: "<module> <record-id>",
				Flags: []cli.Flag{
					&cli.IntFlag{Name: "page", Value: 1},
					&cli.IntFlag{Name: "per-page", Value: 200},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					if cmd.Args().Len() < 2 {
						return internal.NewValidationError("module name and record ID required")
					}
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					module, recordID := cmd.Args().Get(0), cmd.Args().Get(1)
					params := map[string]string{
						"page":     fmt.Sprintf("%d", cmd.Int("page")),
						"per_page": fmt.Sprintf("%d", cmd.Int("per-page")),
					}
					raw, err := c.Request("GET", c.RecruitBase+"/"+module+"/"+recordID+"/associate", &zohttp.RequestOpts{Params: params})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
		},
	}
}

func relatedCmd() *cli.Command {
	return &cli.Command{
		Name:  "related",
		Usage: "Recruit related records",
		Commands: []*cli.Command{
			{
				Name:      "list",
				Usage:     "List related records",
				ArgsUsage: "<module> <record-id> <related-module>",
				Flags: []cli.Flag{
					&cli.IntFlag{Name: "page", Value: 1},
					&cli.IntFlag{Name: "per-page", Value: 200},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					if cmd.Args().Len() < 3 {
						return internal.NewValidationError("module name, record ID, and related module required")
					}
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					module, recordID, relModule := cmd.Args().Get(0), cmd.Args().Get(1), cmd.Args().Get(2)
					params := map[string]string{
						"page":     fmt.Sprintf("%d", cmd.Int("page")),
						"per_page": fmt.Sprintf("%d", cmd.Int("per-page")),
					}
					raw, err := c.Request("GET", c.RecruitBase+"/"+module+"/"+recordID+"/"+relModule, &zohttp.RequestOpts{Params: params})
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
		Usage: "Recruit users",
		Commands: []*cli.Command{
			{
				Name:  "list",
				Usage: "List Recruit users",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "type", Value: "AllUsers", Usage: "User type filter"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					raw, err := c.Request("GET", c.RecruitBase+"/users", &zohttp.RequestOpts{
						Params: map[string]string{"type": cmd.String("type")},
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
