package crm

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/omin8tor/zoho-cli/internal"
	zohttp "github.com/omin8tor/zoho-cli/internal/http"
	"github.com/omin8tor/zoho-cli/internal/output"
	"github.com/omin8tor/zoho-cli/internal/pagination"
	"github.com/urfave/cli/v3"
)

func Commands() *cli.Command {
	return &cli.Command{
		Name:  "crm",
		Usage: "Zoho CRM operations",
		Commands: []*cli.Command{
			modulesCmd(),
			recordsCmd(),
			notesCmd(),
			relatedCmd(),
			usersCmd(),
			ownerCmd(),
			coqlCmd(),
			searchGlobalCmd(),
			attachmentsCmd(),
			convertCmd(),
			tagsCmd(),
		},
	}
}

const defaultFields = "id,Created_Time,Modified_Time"

func modulesCmd() *cli.Command {
	return &cli.Command{
		Name:  "modules",
		Usage: "CRM module operations",
		Commands: []*cli.Command{
			{
				Name:      "list",
				Usage:     "List available CRM modules",
				ArgsUsage: "",
				Flags: []cli.Flag{
					&cli.BoolFlag{Name: "include-hidden", Usage: "Include hidden/system modules"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					raw, err := c.Request("GET", c.CRMBase+"/settings/modules", nil)
					if err != nil {
						return err
					}
					var envelope struct {
						Modules []json.RawMessage `json:"modules"`
					}
					json.Unmarshal(raw, &envelope)
					if cmd.Bool("include-hidden") {
						return output.JSON(envelope.Modules)
					}
					var filtered []json.RawMessage
					for _, m := range envelope.Modules {
						var mod map[string]any
						json.Unmarshal(m, &mod)
						if show, _ := mod["show_as_tab"].(bool); show {
							filtered = append(filtered, m)
						}
					}
					return output.JSON(filtered)
				},
			},
			{
				Name:      "fields",
				Usage:     "List fields for a CRM module",
				ArgsUsage: "<module>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					raw, err := c.Request("GET", c.CRMBase+"/settings/fields", &zohttp.RequestOpts{
						Params: map[string]string{"module": cmd.Args().First()},
					})
					if err != nil {
						return err
					}
					var envelope struct {
						Fields []json.RawMessage `json:"fields"`
					}
					json.Unmarshal(raw, &envelope)
					return output.JSON(envelope.Fields)
				},
			},
			{
				Name:      "related-lists",
				Usage:     "List related lists for a module",
				ArgsUsage: "<module>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					raw, err := c.Request("GET", c.CRMBase+"/settings/related_lists", &zohttp.RequestOpts{
						Params: map[string]string{"module": cmd.Args().First()},
					})
					if err != nil {
						return err
					}
					var envelope struct {
						RelatedLists []json.RawMessage `json:"related_lists"`
					}
					json.Unmarshal(raw, &envelope)
					return output.JSON(envelope.RelatedLists)
				},
			},
			{
				Name:      "layouts",
				Usage:     "List layouts for a module",
				ArgsUsage: "<module>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					raw, err := c.Request("GET", c.CRMBase+"/settings/layouts", &zohttp.RequestOpts{
						Params: map[string]string{"module": cmd.Args().First()},
					})
					if err != nil {
						return err
					}
					var envelope struct {
						Layouts []json.RawMessage `json:"layouts"`
					}
					json.Unmarshal(raw, &envelope)
					return output.JSON(envelope.Layouts)
				},
			},
			{
				Name:      "custom-views",
				Usage:     "List custom views for a module",
				ArgsUsage: "<module>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					raw, err := c.Request("GET", c.CRMBase+"/settings/custom_views", &zohttp.RequestOpts{
						Params: map[string]string{"module": cmd.Args().First()},
					})
					if err != nil {
						return err
					}
					var envelope struct {
						CustomViews []json.RawMessage `json:"custom_views"`
					}
					json.Unmarshal(raw, &envelope)
					return output.JSON(envelope.CustomViews)
				},
			},
		},
	}
}

func recordsCmd() *cli.Command {
	return &cli.Command{
		Name:  "records",
		Usage: "CRM record operations",
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
					&cli.BoolFlag{Name: "all", Usage: "Auto-paginate all records"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					module := cmd.Args().First()
					fields := cmd.String("fields")
					if fields == "" {
						fields = defaultFields
					}
					params := map[string]string{"fields": fields}
					if s := cmd.String("sort-by"); s != "" {
						params["sort_by"] = s
					}
					if s := cmd.String("sort-order"); s != "" {
						params["sort_order"] = s
					}
					if cmd.Bool("all") {
						records, err := pagination.PaginateCRM(c, c.CRMBase+"/"+module, params, 0)
						if err != nil {
							return err
						}
						return output.JSON(records)
					}
					params["page"] = fmt.Sprintf("%d", cmd.Int("page"))
					pp := cmd.Int("per-page")
					if pp > 200 {
						pp = 200
					}
					params["per_page"] = fmt.Sprintf("%d", pp)
					raw, err := c.Request("GET", c.CRMBase+"/"+module, &zohttp.RequestOpts{Params: params})
					if err != nil {
						return err
					}
					var envelope struct {
						Data []json.RawMessage `json:"data"`
					}
					json.Unmarshal(raw, &envelope)
					return output.JSON(envelope.Data)
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
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					module, recordID := cmd.Args().Get(0), cmd.Args().Get(1)
					params := map[string]string{}
					if f := cmd.String("fields"); f != "" {
						params["fields"] = f
					}
					var opts *zohttp.RequestOpts
					if len(params) > 0 {
						opts = &zohttp.RequestOpts{Params: params}
					}
					raw, err := c.Request("GET", c.CRMBase+"/"+module+"/"+recordID, opts)
					if err != nil {
						return err
					}
					var envelope struct {
						Data []json.RawMessage `json:"data"`
					}
					json.Unmarshal(raw, &envelope)
					if len(envelope.Data) > 0 {
						return output.JSONRaw(envelope.Data[0])
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
					&cli.StringFlag{Name: "trigger", Usage: "Comma-separated triggers: approval,workflow,blueprint"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					var parsed map[string]any
					if err := json.Unmarshal([]byte(cmd.String("json")), &parsed); err != nil {
						return internal.NewValidationError(fmt.Sprintf("--json: invalid JSON: %v", err))
					}
					body := map[string]any{"data": []any{parsed}}
					if t := cmd.String("trigger"); t != "" {
						body["trigger"] = splitComma(t)
					}
					raw, err := c.Request("POST", c.CRMBase+"/"+cmd.Args().First(), &zohttp.RequestOpts{JSON: body})
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
					&cli.StringFlag{Name: "trigger", Usage: "Comma-separated triggers"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					module, recordID := cmd.Args().Get(0), cmd.Args().Get(1)
					var parsed map[string]any
					if err := json.Unmarshal([]byte(cmd.String("json")), &parsed); err != nil {
						return internal.NewValidationError(fmt.Sprintf("--json: invalid JSON: %v", err))
					}
					body := map[string]any{"data": []any{parsed}}
					if t := cmd.String("trigger"); t != "" {
						body["trigger"] = splitComma(t)
					}
					raw, err := c.Request("PUT", c.CRMBase+"/"+module+"/"+recordID, &zohttp.RequestOpts{JSON: body})
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
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					module, recordID := cmd.Args().Get(0), cmd.Args().Get(1)
					raw, err := c.Request("DELETE", c.CRMBase+"/"+module+"/"+recordID, nil)
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
					&cli.StringFlag{Name: "word", Usage: "Keyword search"},
					&cli.StringFlag{Name: "email", Usage: "Email search"},
					&cli.StringFlag{Name: "phone", Usage: "Phone search"},
					&cli.StringFlag{Name: "criteria", Usage: "Criteria e.g. (Stage:equals:Closed Won)"},
					&cli.StringFlag{Name: "fields", Usage: "Comma-separated fields"},
					&cli.IntFlag{Name: "page", Value: 1},
					&cli.IntFlag{Name: "per-page", Value: 200},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
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
					if f := cmd.String("fields"); f != "" {
						params["fields"] = f
					}
					raw, err := c.Request("GET", c.CRMBase+"/"+cmd.Args().First()+"/search", &zohttp.RequestOpts{Params: params})
					if err != nil {
						return err
					}
					var envelope struct {
						Data []json.RawMessage `json:"data"`
					}
					json.Unmarshal(raw, &envelope)
					return output.JSON(envelope.Data)
				},
			},
			{
				Name:      "upsert",
				Usage:     "Upsert a record (insert or update)",
				ArgsUsage: "<module>",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "json", Required: true, Usage: "Record data as JSON string"},
					&cli.StringFlag{Name: "duplicate-check", Usage: "Comma-separated duplicate check fields"},
					&cli.StringFlag{Name: "trigger", Usage: "Comma-separated triggers"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					var parsed map[string]any
					if err := json.Unmarshal([]byte(cmd.String("json")), &parsed); err != nil {
						return internal.NewValidationError(fmt.Sprintf("--json: invalid JSON: %v", err))
					}
					body := map[string]any{"data": []any{parsed}}
					if d := cmd.String("duplicate-check"); d != "" {
						body["duplicate_check_fields"] = splitComma(d)
					}
					if t := cmd.String("trigger"); t != "" {
						body["trigger"] = splitComma(t)
					}
					raw, err := c.Request("POST", c.CRMBase+"/"+cmd.Args().First()+"/upsert", &zohttp.RequestOpts{JSON: body})
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
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					module, ids := cmd.Args().Get(0), cmd.Args().Get(1)
					raw, err := c.Request("DELETE", c.CRMBase+"/"+module, &zohttp.RequestOpts{
						Params: map[string]string{"ids": ids},
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

func notesCmd() *cli.Command {
	return &cli.Command{
		Name:  "notes",
		Usage: "CRM record notes",
		Commands: []*cli.Command{
			{
				Name:      "list",
				Usage:     "List notes on a record",
				ArgsUsage: "<module> <record-id>",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "fields", Usage: "Comma-separated field API names"},
					&cli.IntFlag{Name: "page", Value: 1},
					&cli.IntFlag{Name: "per-page", Value: 200},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					module, recordID := cmd.Args().Get(0), cmd.Args().Get(1)
					fields := cmd.String("fields")
					if fields == "" {
						fields = "id,Note_Title,Note_Content,Created_Time,Modified_Time"
					}
					params := map[string]string{
						"fields":   fields,
						"page":     fmt.Sprintf("%d", cmd.Int("page")),
						"per_page": fmt.Sprintf("%d", min(cmd.Int("per-page"), 200)),
					}
					raw, err := c.Request("GET", c.CRMBase+"/"+module+"/"+recordID+"/Notes", &zohttp.RequestOpts{Params: params})
					if err != nil {
						return err
					}
					var envelope struct {
						Data []json.RawMessage `json:"data"`
					}
					json.Unmarshal(raw, &envelope)
					return output.JSON(envelope.Data)
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
					raw, err := c.Request("POST", c.CRMBase+"/"+module+"/"+recordID+"/Notes", &zohttp.RequestOpts{JSON: body})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "update",
				Usage:     "Update a note",
				ArgsUsage: "<note-id>",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "title", Usage: "Note title"},
					&cli.StringFlag{Name: "content", Usage: "Note content"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					note := map[string]string{}
					if t := cmd.String("title"); t != "" {
						note["Note_Title"] = t
					}
					if ct := cmd.String("content"); ct != "" {
						note["Note_Content"] = ct
					}
					body := map[string]any{"data": []any{note}}
					raw, err := c.Request("PUT", c.CRMBase+"/Notes/"+cmd.Args().First(), &zohttp.RequestOpts{JSON: body})
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
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					raw, err := c.Request("DELETE", c.CRMBase+"/Notes/"+cmd.Args().First(), nil)
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
		Usage: "CRM related records",
		Commands: []*cli.Command{
			{
				Name:      "list",
				Usage:     "List related records",
				ArgsUsage: "<module> <record-id> <related-list>",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "fields", Usage: "Comma-separated fields"},
					&cli.IntFlag{Name: "page", Value: 1},
					&cli.IntFlag{Name: "per-page", Value: 200},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					module, recordID, relList := cmd.Args().Get(0), cmd.Args().Get(1), cmd.Args().Get(2)
					fields := cmd.String("fields")
					if fields == "" {
						fields = defaultFields
					}
					params := map[string]string{
						"fields":   fields,
						"page":     fmt.Sprintf("%d", cmd.Int("page")),
						"per_page": fmt.Sprintf("%d", min(cmd.Int("per-page"), 200)),
					}
					raw, err := c.Request("GET", c.CRMBase+"/"+module+"/"+recordID+"/"+relList, &zohttp.RequestOpts{Params: params})
					if err != nil {
						return err
					}
					var envelope struct {
						Data []json.RawMessage `json:"data"`
					}
					json.Unmarshal(raw, &envelope)
					return output.JSON(envelope.Data)
				},
			},
		},
	}
}

func usersCmd() *cli.Command {
	return &cli.Command{
		Name:  "users",
		Usage: "CRM users",
		Commands: []*cli.Command{
			{
				Name:  "list",
				Usage: "List CRM users",
				Flags: []cli.Flag{
					&cli.IntFlag{Name: "page", Value: 1},
					&cli.IntFlag{Name: "per-page", Value: 200},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					params := map[string]string{
						"type":     "AllUsers",
						"page":     fmt.Sprintf("%d", cmd.Int("page")),
						"per_page": fmt.Sprintf("%d", min(cmd.Int("per-page"), 200)),
					}
					raw, err := c.Request("GET", c.CRMBase+"/users", &zohttp.RequestOpts{Params: params})
					if err != nil {
						return err
					}
					var envelope struct {
						Users []json.RawMessage `json:"users"`
					}
					json.Unmarshal(raw, &envelope)
					return output.JSON(envelope.Users)
				},
			},
		},
	}
}

func ownerCmd() *cli.Command {
	return &cli.Command{
		Name:  "owner",
		Usage: "Change record owner",
		Commands: []*cli.Command{
			{
				Name:      "change",
				Usage:     "Change record owner",
				ArgsUsage: "<module> <record-id>",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "owner", Required: true, Usage: "New owner user ID"},
					&cli.BoolFlag{Name: "notify", Value: true, Usage: "Notify new owner"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					module, recordID := cmd.Args().Get(0), cmd.Args().Get(1)
					body := map[string]any{
						"owner":  map[string]string{"id": cmd.String("owner")},
						"notify": cmd.Bool("notify"),
					}
					raw, err := c.Request("POST", c.CRMBase+"/"+module+"/"+recordID+"/actions/change_owner", &zohttp.RequestOpts{JSON: body})
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
			raw, err := c.Request("POST", c.CRMBase+"/coql", &zohttp.RequestOpts{JSON: body})
			if err != nil {
				return err
			}
			return output.JSONRaw(raw)
		},
	}
}

func searchGlobalCmd() *cli.Command {
	return &cli.Command{
		Name:      "search-global",
		Usage:     "Search across all CRM modules",
		ArgsUsage: "<word>",
		Flags: []cli.Flag{
			&cli.IntFlag{Name: "page", Value: 1},
			&cli.IntFlag{Name: "per-page", Value: 10},
		},
		Action: func(_ context.Context, cmd *cli.Command) error {
			c, err := zohttp.GetClient()
			if err != nil {
				return err
			}
			params := map[string]string{
				"searchword": cmd.Args().First(),
				"page":       fmt.Sprintf("%d", cmd.Int("page")),
				"per_page":   fmt.Sprintf("%d", min(cmd.Int("per-page"), 10)),
			}
			raw, err := c.Request("GET", c.CRMBase+"/search", &zohttp.RequestOpts{Params: params})
			if err != nil {
				return err
			}
			return output.JSONRaw(raw)
		},
	}
}

func attachmentsCmd() *cli.Command {
	return &cli.Command{
		Name:  "attachments",
		Usage: "CRM record attachments",
		Commands: []*cli.Command{
			{
				Name:      "list",
				Usage:     "List attachments on a record",
				ArgsUsage: "<module> <record-id>",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "fields", Usage: "Comma-separated field API names"},
					&cli.IntFlag{Name: "page", Value: 1},
					&cli.IntFlag{Name: "per-page", Value: 200},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					module, recordID := cmd.Args().Get(0), cmd.Args().Get(1)
					fields := cmd.String("fields")
					if fields == "" {
						fields = "id,File_Name,Size,Created_Time"
					}
					params := map[string]string{
						"fields":   fields,
						"page":     fmt.Sprintf("%d", cmd.Int("page")),
						"per_page": fmt.Sprintf("%d", min(cmd.Int("per-page"), 200)),
					}
					raw, err := c.Request("GET", c.CRMBase+"/"+module+"/"+recordID+"/Attachments", &zohttp.RequestOpts{Params: params})
					if err != nil {
						return err
					}
					var envelope struct {
						Data []json.RawMessage `json:"data"`
					}
					json.Unmarshal(raw, &envelope)
					return output.JSON(envelope.Data)
				},
			},
			{
				Name:      "upload",
				Usage:     "Upload an attachment to a record",
				ArgsUsage: "<module> <record-id> <file-path>",
				Action: func(_ context.Context, cmd *cli.Command) error {
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
					raw, err := c.Request("POST", c.CRMBase+"/"+module+"/"+recordID+"/Attachments", &zohttp.RequestOpts{
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
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					module, recordID, attID := cmd.Args().Get(0), cmd.Args().Get(1), cmd.Args().Get(2)
					url := c.CRMBase + "/" + module + "/" + recordID + "/Attachments/" + attID
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
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					module, recordID, attID := cmd.Args().Get(0), cmd.Args().Get(1), cmd.Args().Get(2)
					raw, err := c.Request("DELETE", c.CRMBase+"/"+module+"/"+recordID+"/Attachments/"+attID, nil)
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
		},
	}
}

func convertCmd() *cli.Command {
	return &cli.Command{
		Name:      "convert",
		Usage:     "Convert a lead to contact/account/deal",
		ArgsUsage: "<record-id>",
		Flags: []cli.Flag{
			&cli.StringFlag{Name: "json", Usage: "Conversion options as JSON"},
		},
		Action: func(_ context.Context, cmd *cli.Command) error {
			c, err := zohttp.GetClient()
			if err != nil {
				return err
			}
			var opts map[string]any
			if j := cmd.String("json"); j != "" {
				if err := json.Unmarshal([]byte(j), &opts); err != nil {
					return internal.NewValidationError(fmt.Sprintf("--json: invalid JSON: %v", err))
				}
			} else {
				opts = map[string]any{}
			}
			body := map[string]any{"data": []any{opts}}
			raw, err := c.Request("POST", c.CRMBase+"/Leads/"+cmd.Args().First()+"/actions/convert", &zohttp.RequestOpts{JSON: body})
			if err != nil {
				return err
			}
			return output.JSONRaw(raw)
		},
	}
}

func tagsCmd() *cli.Command {
	return &cli.Command{
		Name:  "tags",
		Usage: "CRM tag operations",
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
					raw, err := c.Request("POST", c.CRMBase+"/"+cmd.Args().First()+"/actions/add_tags", &zohttp.RequestOpts{JSON: body})
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
					raw, err := c.Request("POST", c.CRMBase+"/"+cmd.Args().First()+"/actions/remove_tags", &zohttp.RequestOpts{JSON: body})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
		},
	}
}

func splitComma(s string) []string {
	var result []string
	start := 0
	for i := range len(s) {
		if s[i] == ',' {
			t := trimSpace(s[start:i])
			if t != "" {
				result = append(result, t)
			}
			start = i + 1
		}
	}
	if t := trimSpace(s[start:]); t != "" {
		result = append(result, t)
	}
	return result
}

func trimSpace(s string) string {
	start, end := 0, len(s)
	for start < end && (s[start] == ' ' || s[start] == '\t') {
		start++
	}
	for end > start && (s[end-1] == ' ' || s[end-1] == '\t') {
		end--
	}
	return s[start:end]
}
