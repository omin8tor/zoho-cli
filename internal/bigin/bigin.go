package bigin

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
		Name:  "bigin",
		Usage: "Zoho Bigin operations",
		Commands: []*cli.Command{
			modulesCmd(),
			recordsCmd(),
			notesCmd(),
			attachmentsCmd(),
			tagsCmd(),
			usersCmd(),
			orgCmd(),
			rolesCmd(),
			profilesCmd(),
			relatedCmd(),
			coqlCmd(),
			searchCmd(),
		},
	}
}

func modulesCmd() *cli.Command {
	return &cli.Command{
		Name:  "modules",
		Usage: "Bigin module operations",
		Commands: []*cli.Command{
			{
				Name:  "list",
				Usage: "List available Bigin modules",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					raw, err := c.Request("GET", c.BiginBase+"/settings/modules", nil)
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "get",
				Usage:     "Get a specific module",
				ArgsUsage: "<module>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					if cmd.Args().Len() < 1 {
						return internal.NewValidationError("module API name required")
					}
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					raw, err := c.Request("GET", c.BiginBase+"/settings/modules/"+cmd.Args().First(), nil)
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "fields",
				Usage:     "List fields for a Bigin module",
				ArgsUsage: "<module>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					if cmd.Args().Len() < 1 {
						return internal.NewValidationError("module API name required")
					}
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					raw, err := c.Request("GET", c.BiginBase+"/settings/fields", &zohttp.RequestOpts{
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
						return internal.NewValidationError("module API name required")
					}
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					raw, err := c.Request("GET", c.BiginBase+"/settings/layouts", &zohttp.RequestOpts{
						Params: map[string]string{"module": cmd.Args().First()},
					})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "related-lists",
				Usage:     "List related lists for a module",
				ArgsUsage: "<module>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					if cmd.Args().Len() < 1 {
						return internal.NewValidationError("module API name required")
					}
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					raw, err := c.Request("GET", c.BiginBase+"/settings/related_lists", &zohttp.RequestOpts{
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
		Usage: "Bigin record operations",
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
					&cli.IntFlag{Name: "per-page", Value: 200, Usage: "Records per page (max 200)"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					if cmd.Args().Len() < 1 {
						return internal.NewValidationError("module API name required")
					}
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					params := map[string]string{
						"page":     fmt.Sprintf("%d", cmd.Int("page")),
						"per_page": fmt.Sprintf("%d", min(cmd.Int("per-page"), 200)),
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
					raw, err := c.Request("GET", c.BiginBase+"/"+cmd.Args().First(), &zohttp.RequestOpts{Params: params})
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
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "fields", Usage: "Comma-separated field API names"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					if cmd.Args().Len() < 2 {
						return internal.NewValidationError("module and record ID required")
					}
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					module, recordID := cmd.Args().Get(0), cmd.Args().Get(1)
					var opts *zohttp.RequestOpts
					if f := cmd.String("fields"); f != "" {
						opts = &zohttp.RequestOpts{Params: map[string]string{"fields": f}}
					}
					raw, err := c.Request("GET", c.BiginBase+"/"+module+"/"+recordID, opts)
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
					&cli.StringFlag{Name: "json", Required: true, Usage: "Record data as JSON string"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					if cmd.Args().Len() < 1 {
						return internal.NewValidationError("module API name required")
					}
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					var parsed map[string]any
					if err := json.Unmarshal([]byte(cmd.String("json")), &parsed); err != nil {
						return internal.NewValidationError(fmt.Sprintf("invalid JSON: %v", err))
					}
					body := map[string]any{"data": []any{parsed}}
					raw, err := c.Request("POST", c.BiginBase+"/"+cmd.Args().First(), &zohttp.RequestOpts{JSON: body})
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
						return internal.NewValidationError("module and record ID required")
					}
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					module, recordID := cmd.Args().Get(0), cmd.Args().Get(1)
					var parsed map[string]any
					if err := json.Unmarshal([]byte(cmd.String("json")), &parsed); err != nil {
						return internal.NewValidationError(fmt.Sprintf("invalid JSON: %v", err))
					}
					body := map[string]any{"data": []any{parsed}}
					raw, err := c.Request("PUT", c.BiginBase+"/"+module+"/"+recordID, &zohttp.RequestOpts{JSON: body})
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
						return internal.NewValidationError("module and record ID required")
					}
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					module, recordID := cmd.Args().Get(0), cmd.Args().Get(1)
					raw, err := c.Request("DELETE", c.BiginBase+"/"+module+"/"+recordID, nil)
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "upsert",
				Usage:     "Upsert a record (insert or update)",
				ArgsUsage: "<module>",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "json", Required: true, Usage: "Record data as JSON string"},
					&cli.StringFlag{Name: "duplicate-check", Usage: "Comma-separated duplicate check fields"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					if cmd.Args().Len() < 1 {
						return internal.NewValidationError("module API name required")
					}
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					var parsed map[string]any
					if err := json.Unmarshal([]byte(cmd.String("json")), &parsed); err != nil {
						return internal.NewValidationError(fmt.Sprintf("invalid JSON: %v", err))
					}
					body := map[string]any{"data": []any{parsed}}
					if d := cmd.String("duplicate-check"); d != "" {
						body["duplicate_check_fields"] = splitComma(d)
					}
					raw, err := c.Request("POST", c.BiginBase+"/"+cmd.Args().First()+"/upsert", &zohttp.RequestOpts{JSON: body})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "bulk-delete",
				Usage:     "Delete multiple records",
				ArgsUsage: "<module> <ids>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					if cmd.Args().Len() < 2 {
						return internal.NewValidationError("module and comma-separated IDs required")
					}
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					module, ids := cmd.Args().Get(0), cmd.Args().Get(1)
					raw, err := c.Request("DELETE", c.BiginBase+"/"+module, &zohttp.RequestOpts{
						Params: map[string]string{"ids": ids},
					})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "deleted",
				Usage:     "List deleted records",
				ArgsUsage: "<module>",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "type", Value: "all", Usage: "Type: all, recycle, permanent"},
					&cli.IntFlag{Name: "page", Value: 1},
					&cli.IntFlag{Name: "per-page", Value: 200},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					if cmd.Args().Len() < 1 {
						return internal.NewValidationError("module API name required")
					}
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					params := map[string]string{
						"type":     cmd.String("type"),
						"page":     fmt.Sprintf("%d", cmd.Int("page")),
						"per_page": fmt.Sprintf("%d", min(cmd.Int("per-page"), 200)),
					}
					raw, err := c.Request("GET", c.BiginBase+"/"+cmd.Args().First()+"/deleted", &zohttp.RequestOpts{Params: params})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "count",
				Usage:     "Get record count for a module",
				ArgsUsage: "<module>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					if cmd.Args().Len() < 1 {
						return internal.NewValidationError("module API name required")
					}
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					raw, err := c.Request("GET", c.BiginBase+"/"+cmd.Args().First()+"/actions/count", nil)
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
		Usage: "Bigin record notes",
		Commands: []*cli.Command{
			{
				Name:      "list",
				Usage:     "List notes on a record",
				ArgsUsage: "<module> <record-id>",
				Flags: []cli.Flag{
					&cli.IntFlag{Name: "page", Value: 1},
					&cli.IntFlag{Name: "per-page", Value: 200},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					if cmd.Args().Len() < 2 {
						return internal.NewValidationError("module and record ID required")
					}
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					module, recordID := cmd.Args().Get(0), cmd.Args().Get(1)
					params := map[string]string{
						"page":     fmt.Sprintf("%d", cmd.Int("page")),
						"per_page": fmt.Sprintf("%d", min(cmd.Int("per-page"), 200)),
					}
					raw, err := c.Request("GET", c.BiginBase+"/"+module+"/"+recordID+"/Notes", &zohttp.RequestOpts{Params: params})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "add",
				Usage:     "Add a note to a record",
				ArgsUsage: "<module> <record-id>",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "content", Required: true, Usage: "Note content"},
					&cli.StringFlag{Name: "title", Usage: "Note title"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					if cmd.Args().Len() < 2 {
						return internal.NewValidationError("module and record ID required")
					}
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					module, recordID := cmd.Args().Get(0), cmd.Args().Get(1)
					note := map[string]string{"Note_Content": cmd.String("content")}
					if t := cmd.String("title"); t != "" {
						note["Note_Title"] = t
					}
					body := map[string]any{"data": []any{note}}
					raw, err := c.Request("POST", c.BiginBase+"/"+module+"/"+recordID+"/Notes", &zohttp.RequestOpts{JSON: body})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "update",
				Usage:     "Update a note",
				ArgsUsage: "<module> <record-id> <note-id>",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "title", Usage: "Note title"},
					&cli.StringFlag{Name: "content", Usage: "Note content"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					if cmd.Args().Len() < 3 {
						return internal.NewValidationError("module, record ID, and note ID required")
					}
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					module, recordID, noteID := cmd.Args().Get(0), cmd.Args().Get(1), cmd.Args().Get(2)
					note := map[string]string{}
					if t := cmd.String("title"); t != "" {
						note["Note_Title"] = t
					}
					if ct := cmd.String("content"); ct != "" {
						note["Note_Content"] = ct
					}
					body := map[string]any{"data": []any{note}}
					raw, err := c.Request("PUT", c.BiginBase+"/"+module+"/"+recordID+"/Notes/"+noteID, &zohttp.RequestOpts{JSON: body})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "delete",
				Usage:     "Delete a note",
				ArgsUsage: "<module> <record-id> <note-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					if cmd.Args().Len() < 3 {
						return internal.NewValidationError("module, record ID, and note ID required")
					}
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					module, recordID, noteID := cmd.Args().Get(0), cmd.Args().Get(1), cmd.Args().Get(2)
					raw, err := c.Request("DELETE", c.BiginBase+"/"+module+"/"+recordID+"/Notes/"+noteID, nil)
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
		Usage: "Bigin record attachments",
		Commands: []*cli.Command{
			{
				Name:      "list",
				Usage:     "List attachments on a record",
				ArgsUsage: "<module> <record-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					if cmd.Args().Len() < 2 {
						return internal.NewValidationError("module and record ID required")
					}
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					module, recordID := cmd.Args().Get(0), cmd.Args().Get(1)
					raw, err := c.Request("GET", c.BiginBase+"/"+module+"/"+recordID+"/Attachments", nil)
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
				Action: func(_ context.Context, cmd *cli.Command) error {
					if cmd.Args().Len() < 3 {
						return internal.NewValidationError("module, record ID, and file path required")
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
					for i := len(filePath) - 1; i >= 0; i-- {
						if filePath[i] == '/' || filePath[i] == '\\' {
							name = filePath[i+1:]
							break
						}
					}
					raw, err := c.Request("POST", c.BiginBase+"/"+module+"/"+recordID+"/Attachments", &zohttp.RequestOpts{
						Files: map[string]zohttp.FileUpload{"file": {Filename: name, Data: data}},
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
						return internal.NewValidationError("module, record ID, and attachment ID required")
					}
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					module, recordID, attID := cmd.Args().Get(0), cmd.Args().Get(1), cmd.Args().Get(2)
					url := c.BiginBase + "/" + module + "/" + recordID + "/Attachments/" + attID
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
						return internal.NewValidationError("module, record ID, and attachment ID required")
					}
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					module, recordID, attID := cmd.Args().Get(0), cmd.Args().Get(1), cmd.Args().Get(2)
					raw, err := c.Request("DELETE", c.BiginBase+"/"+module+"/"+recordID+"/Attachments/"+attID, nil)
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
		Usage: "Bigin tag operations",
		Commands: []*cli.Command{
			{
				Name:      "add",
				Usage:     "Add tags to records",
				ArgsUsage: "<module>",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "ids", Required: true, Usage: "Comma-separated record IDs"},
					&cli.StringFlag{Name: "tags", Required: true, Usage: "Comma-separated tag names"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					if cmd.Args().Len() < 1 {
						return internal.NewValidationError("module API name required")
					}
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					tags := splitComma(cmd.String("tags"))
					tagObjs := make([]map[string]string, len(tags))
					for i, t := range tags {
						tagObjs[i] = map[string]string{"name": t}
					}
					body := map[string]any{
						"tags": tagObjs,
						"ids":  splitComma(cmd.String("ids")),
					}
					raw, err := c.Request("POST", c.BiginBase+"/"+cmd.Args().First()+"/actions/add_tags", &zohttp.RequestOpts{JSON: body})
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
					&cli.StringFlag{Name: "tags", Required: true, Usage: "Comma-separated tag names"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					if cmd.Args().Len() < 1 {
						return internal.NewValidationError("module API name required")
					}
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					tags := splitComma(cmd.String("tags"))
					tagObjs := make([]map[string]string, len(tags))
					for i, t := range tags {
						tagObjs[i] = map[string]string{"name": t}
					}
					body := map[string]any{
						"tags": tagObjs,
						"ids":  splitComma(cmd.String("ids")),
					}
					raw, err := c.Request("POST", c.BiginBase+"/"+cmd.Args().First()+"/actions/remove_tags", &zohttp.RequestOpts{JSON: body})
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
		Usage: "Bigin users",
		Commands: []*cli.Command{
			{
				Name:  "list",
				Usage: "List Bigin users",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "type", Value: "AllUsers", Usage: "User type: AllUsers, ActiveUsers, DeactiveUsers, AdminUsers"},
					&cli.IntFlag{Name: "page", Value: 1},
					&cli.IntFlag{Name: "per-page", Value: 200},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					params := map[string]string{
						"type":     cmd.String("type"),
						"page":     fmt.Sprintf("%d", cmd.Int("page")),
						"per_page": fmt.Sprintf("%d", min(cmd.Int("per-page"), 200)),
					}
					raw, err := c.Request("GET", c.BiginBase+"/users", &zohttp.RequestOpts{Params: params})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "get",
				Usage:     "Get a specific user",
				ArgsUsage: "<user-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					if cmd.Args().Len() < 1 {
						return internal.NewValidationError("user ID required")
					}
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					raw, err := c.Request("GET", c.BiginBase+"/users/"+cmd.Args().First(), nil)
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
		},
	}
}

func orgCmd() *cli.Command {
	return &cli.Command{
		Name:  "org",
		Usage: "Bigin organization details",
		Commands: []*cli.Command{
			{
				Name:  "get",
				Usage: "Get organization details",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					raw, err := c.Request("GET", c.BiginBase+"/org", nil)
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
		},
	}
}

func rolesCmd() *cli.Command {
	return &cli.Command{
		Name:  "roles",
		Usage: "Bigin roles",
		Commands: []*cli.Command{
			{
				Name:  "list",
				Usage: "List roles",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					raw, err := c.Request("GET", c.BiginBase+"/settings/roles", nil)
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "get",
				Usage:     "Get a specific role",
				ArgsUsage: "<role-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					if cmd.Args().Len() < 1 {
						return internal.NewValidationError("role ID required")
					}
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					raw, err := c.Request("GET", c.BiginBase+"/settings/roles/"+cmd.Args().First(), nil)
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
		},
	}
}

func profilesCmd() *cli.Command {
	return &cli.Command{
		Name:  "profiles",
		Usage: "Bigin profiles",
		Commands: []*cli.Command{
			{
				Name:  "list",
				Usage: "List profiles",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					raw, err := c.Request("GET", c.BiginBase+"/settings/profiles", nil)
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "get",
				Usage:     "Get a specific profile",
				ArgsUsage: "<profile-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					if cmd.Args().Len() < 1 {
						return internal.NewValidationError("profile ID required")
					}
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					raw, err := c.Request("GET", c.BiginBase+"/settings/profiles/"+cmd.Args().First(), nil)
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
		Usage: "Bigin related records",
		Commands: []*cli.Command{
			{
				Name:      "list",
				Usage:     "List related records",
				ArgsUsage: "<module> <record-id> <related-list>",
				Flags: []cli.Flag{
					&cli.IntFlag{Name: "page", Value: 1},
					&cli.IntFlag{Name: "per-page", Value: 200},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					if cmd.Args().Len() < 3 {
						return internal.NewValidationError("module, record ID, and related list name required")
					}
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					module, recordID, relList := cmd.Args().Get(0), cmd.Args().Get(1), cmd.Args().Get(2)
					params := map[string]string{
						"page":     fmt.Sprintf("%d", cmd.Int("page")),
						"per_page": fmt.Sprintf("%d", min(cmd.Int("per-page"), 200)),
					}
					raw, err := c.Request("GET", c.BiginBase+"/"+module+"/"+recordID+"/"+relList, &zohttp.RequestOpts{Params: params})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
		},
	}
}

func coqlCmd() *cli.Command {
	return &cli.Command{
		Name:  "coql",
		Usage: "Run a COQL query",
		Flags: []cli.Flag{
			&cli.StringFlag{Name: "query", Required: true, Usage: "COQL query string"},
		},
		Action: func(_ context.Context, cmd *cli.Command) error {
			c, err := zohttp.GetClient()
			if err != nil {
				return err
			}
			body := map[string]string{"select_query": cmd.String("query")}
			raw, err := c.Request("POST", c.BiginBase+"/coql", &zohttp.RequestOpts{JSON: body})
			if err != nil {
				return err
			}
			return output.JSONRaw(raw)
		},
	}
}

func searchCmd() *cli.Command {
	return &cli.Command{
		Name:      "search",
		Usage:     "Search records in a module",
		ArgsUsage: "<module>",
		Flags: []cli.Flag{
			&cli.StringFlag{Name: "word", Usage: "Keyword search"},
			&cli.StringFlag{Name: "email", Usage: "Email search"},
			&cli.StringFlag{Name: "phone", Usage: "Phone search"},
			&cli.StringFlag{Name: "criteria", Usage: "Criteria e.g. (Last_Name:equals:Smith)"},
			&cli.IntFlag{Name: "page", Value: 1},
			&cli.IntFlag{Name: "per-page", Value: 200},
		},
		Action: func(_ context.Context, cmd *cli.Command) error {
			if cmd.Args().Len() < 1 {
				return internal.NewValidationError("module API name required")
			}
			c, err := zohttp.GetClient()
			if err != nil {
				return err
			}
			params := map[string]string{
				"page":     fmt.Sprintf("%d", cmd.Int("page")),
				"per_page": fmt.Sprintf("%d", min(cmd.Int("per-page"), 200)),
			}
			if w := cmd.String("word"); w != "" {
				params["word"] = w
			} else if e := cmd.String("email"); e != "" {
				params["email"] = e
			} else if p := cmd.String("phone"); p != "" {
				params["phone"] = p
			} else if cr := cmd.String("criteria"); cr != "" {
				params["criteria"] = cr
			}
			raw, err := c.Request("GET", c.BiginBase+"/"+cmd.Args().First()+"/search", &zohttp.RequestOpts{Params: params})
			if err != nil {
				return err
			}
			return output.JSONRaw(raw)
		},
	}
}

func splitComma(s string) []string {
	var result []string
	for _, part := range strings.Split(s, ",") {
		t := strings.TrimSpace(part)
		if t != "" {
			result = append(result, t)
		}
	}
	return result
}
