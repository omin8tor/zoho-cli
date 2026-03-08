package creator

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/omin8tor/zoho-cli/internal"
	zohttp "github.com/omin8tor/zoho-cli/internal/http"
	"github.com/omin8tor/zoho-cli/internal/output"
	"github.com/urfave/cli/v3"
)

func resolveOwner(cmd *cli.Command) (string, error) {
	return internal.RequireFlag(cmd, "owner", "ZOHO_CREATOR_OWNER")
}

func resolveApp(cmd *cli.Command) (string, error) {
	return internal.RequireFlag(cmd, "app", "ZOHO_CREATOR_APP")
}

func dataBase(base, owner, app string) string {
	return fmt.Sprintf("%s/data/%s/%s", base, owner, app)
}

func metaBase(base, owner, app string) string {
	return fmt.Sprintf("%s/meta/%s/%s", base, owner, app)
}

func bulkBase(base, owner, app string) string {
	return fmt.Sprintf("%s/bulk/%s/%s", base, owner, app)
}

func Commands() *cli.Command {
	return &cli.Command{
		Name:  "creator",
		Usage: "Zoho Creator operations",
		Flags: []cli.Flag{
			&cli.StringFlag{Name: "owner", Usage: "Account owner name (or set ZOHO_CREATOR_OWNER)", Sources: cli.EnvVars("ZOHO_CREATOR_OWNER")},
			&cli.StringFlag{Name: "app", Usage: "Application link name (or set ZOHO_CREATOR_APP)", Sources: cli.EnvVars("ZOHO_CREATOR_APP")},
		},
		Commands: []*cli.Command{
			applicationsCmd(),
			recordsCmd(),
			reportsCmd(),
			formsCmd(),
			fieldsCmd(),
			pagesCmd(),
			sectionsCmd(),
			bulkReadCmd(),
			bulkWriteCmd(),
		},
	}
}

func applicationsCmd() *cli.Command {
	return &cli.Command{
		Name:  "applications",
		Usage: "Application operations",
		Commands: []*cli.Command{
			{
				Name:  "list",
				Usage: "List applications",
				Action: func(ctx context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					raw, err := c.Request(ctx, "GET", c.CreatorBase+"/meta/applications", nil)
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
		Usage: "Record operations",
		Commands: []*cli.Command{
			{
				Name:  "list",
				Usage: "List records from a report",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "report", Required: true, Usage: "Report link name"},
					&cli.StringFlag{Name: "max-records", Usage: "Max records per request (up to 1000)"},
					&cli.StringFlag{Name: "field-config", Usage: "Field config: quick_view, detail_view, custom, all"},
					&cli.StringFlag{Name: "fields", Usage: "Comma-separated field names (when field_config=custom)"},
					&cli.StringFlag{Name: "cursor", Usage: "Record cursor for pagination"},
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					owner, err := resolveOwner(cmd)
					if err != nil {
						return err
					}
					app, err := resolveApp(cmd)
					if err != nil {
						return err
					}
					params := map[string]string{}
					if v := cmd.String("max-records"); v != "" {
						params["max_records"] = v
					}
					if v := cmd.String("field-config"); v != "" {
						params["field_config"] = v
					}
					if v := cmd.String("fields"); v != "" {
						params["fields"] = v
					}
					headers := map[string]string{}
					if v := cmd.String("cursor"); v != "" {
						headers["record_cursor"] = v
					}
					raw, err := c.Request(ctx, "GET", dataBase(c.CreatorBase, owner, app)+"/report/"+cmd.String("report"), &zohttp.RequestOpts{
						Params:  params,
						Headers: headers,
					})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "get",
				Usage:     "Get a record by ID",
				ArgsUsage: "<record-id>",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "report", Required: true, Usage: "Report link name"},
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					id := cmd.Args().Get(0)
					if id == "" {
						return internal.NewValidationError("record-id argument required")
					}
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					owner, err := resolveOwner(cmd)
					if err != nil {
						return err
					}
					app, err := resolveApp(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request(ctx, "GET", dataBase(c.CreatorBase, owner, app)+"/report/"+cmd.String("report")+"/"+id, nil)
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:  "add",
				Usage: "Add records to a form",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "form", Required: true, Usage: "Form link name"},
					&cli.StringFlag{Name: "json", Required: true, Usage: "JSON body"},
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					owner, err := resolveOwner(cmd)
					if err != nil {
						return err
					}
					app, err := resolveApp(cmd)
					if err != nil {
						return err
					}
					var body any
					if err := json.Unmarshal([]byte(cmd.String("json")), &body); err != nil {
						return internal.NewValidationError(fmt.Sprintf("invalid JSON: %v", err))
					}
					raw, err := c.Request(ctx, "POST", dataBase(c.CreatorBase, owner, app)+"/form/"+cmd.String("form"), &zohttp.RequestOpts{JSON: body})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "update",
				Usage:     "Update a record by ID",
				ArgsUsage: "<record-id>",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "report", Required: true, Usage: "Report link name"},
					&cli.StringFlag{Name: "json", Required: true, Usage: "JSON body"},
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					id := cmd.Args().Get(0)
					if id == "" {
						return internal.NewValidationError("record-id argument required")
					}
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					owner, err := resolveOwner(cmd)
					if err != nil {
						return err
					}
					app, err := resolveApp(cmd)
					if err != nil {
						return err
					}
					var body any
					if err := json.Unmarshal([]byte(cmd.String("json")), &body); err != nil {
						return internal.NewValidationError(fmt.Sprintf("invalid JSON: %v", err))
					}
					raw, err := c.Request(ctx, "PATCH", dataBase(c.CreatorBase, owner, app)+"/report/"+cmd.String("report")+"/"+id, &zohttp.RequestOpts{JSON: body})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "delete",
				Usage:     "Delete a record by ID",
				ArgsUsage: "<record-id>",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "report", Required: true, Usage: "Report link name"},
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					id := cmd.Args().Get(0)
					if id == "" {
						return internal.NewValidationError("record-id argument required")
					}
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					owner, err := resolveOwner(cmd)
					if err != nil {
						return err
					}
					app, err := resolveApp(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request(ctx, "DELETE", dataBase(c.CreatorBase, owner, app)+"/report/"+cmd.String("report")+"/"+id, nil)
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
		},
	}
}

func reportsCmd() *cli.Command {
	return &cli.Command{
		Name:  "reports",
		Usage: "Report operations",
		Commands: []*cli.Command{
			{
				Name:  "list",
				Usage: "List reports in an application",
				Action: func(ctx context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					owner, err := resolveOwner(cmd)
					if err != nil {
						return err
					}
					app, err := resolveApp(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request(ctx, "GET", metaBase(c.CreatorBase, owner, app)+"/reports", nil)
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
		},
	}
}

func formsCmd() *cli.Command {
	return &cli.Command{
		Name:  "forms",
		Usage: "Form operations",
		Commands: []*cli.Command{
			{
				Name:  "list",
				Usage: "List forms in an application",
				Action: func(ctx context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					owner, err := resolveOwner(cmd)
					if err != nil {
						return err
					}
					app, err := resolveApp(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request(ctx, "GET", metaBase(c.CreatorBase, owner, app)+"/forms", nil)
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
		},
	}
}

func fieldsCmd() *cli.Command {
	return &cli.Command{
		Name:  "fields",
		Usage: "Field operations",
		Commands: []*cli.Command{
			{
				Name:  "list",
				Usage: "List fields of a form",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "form", Required: true, Usage: "Form link name"},
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					owner, err := resolveOwner(cmd)
					if err != nil {
						return err
					}
					app, err := resolveApp(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request(ctx, "GET", metaBase(c.CreatorBase, owner, app)+"/form/"+cmd.String("form")+"/fields", nil)
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
		},
	}
}

func pagesCmd() *cli.Command {
	return &cli.Command{
		Name:  "pages",
		Usage: "Page operations",
		Commands: []*cli.Command{
			{
				Name:  "list",
				Usage: "List pages in an application",
				Action: func(ctx context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					owner, err := resolveOwner(cmd)
					if err != nil {
						return err
					}
					app, err := resolveApp(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request(ctx, "GET", metaBase(c.CreatorBase, owner, app)+"/pages", nil)
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
		},
	}
}

func sectionsCmd() *cli.Command {
	return &cli.Command{
		Name:  "sections",
		Usage: "Section operations",
		Commands: []*cli.Command{
			{
				Name:  "list",
				Usage: "List sections in an application",
				Action: func(ctx context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					owner, err := resolveOwner(cmd)
					if err != nil {
						return err
					}
					app, err := resolveApp(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request(ctx, "GET", metaBase(c.CreatorBase, owner, app)+"/sections", nil)
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
		},
	}
}

func bulkReadCmd() *cli.Command {
	return &cli.Command{
		Name:  "bulk-read",
		Usage: "Bulk read operations",
		Commands: []*cli.Command{
			{
				Name:  "create",
				Usage: "Create a bulk read job",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "report", Required: true, Usage: "Report link name"},
					&cli.StringFlag{Name: "json", Required: true, Usage: "JSON body with query fields/criteria/max_records"},
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					owner, err := resolveOwner(cmd)
					if err != nil {
						return err
					}
					app, err := resolveApp(cmd)
					if err != nil {
						return err
					}
					var body any
					if err := json.Unmarshal([]byte(cmd.String("json")), &body); err != nil {
						return internal.NewValidationError(fmt.Sprintf("invalid JSON: %v", err))
					}
					raw, err := c.Request(ctx, "POST", bulkBase(c.CreatorBase, owner, app)+"/report/"+cmd.String("report")+"/read", &zohttp.RequestOpts{JSON: body})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "status",
				Usage:     "Get bulk read job status",
				ArgsUsage: "<job-id>",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "report", Required: true, Usage: "Report link name"},
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					id := cmd.Args().Get(0)
					if id == "" {
						return internal.NewValidationError("job-id argument required")
					}
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					owner, err := resolveOwner(cmd)
					if err != nil {
						return err
					}
					app, err := resolveApp(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request(ctx, "GET", bulkBase(c.CreatorBase, owner, app)+"/report/"+cmd.String("report")+"/read/"+id, nil)
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
		},
	}
}

func bulkWriteCmd() *cli.Command {
	return &cli.Command{
		Name:  "bulk-write",
		Usage: "Bulk write operations",
		Commands: []*cli.Command{
			{
				Name:  "create",
				Usage: "Create a bulk write job",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "report", Required: true, Usage: "Report link name"},
					&cli.StringFlag{Name: "json", Required: true, Usage: "JSON body with query fields/criteria"},
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					owner, err := resolveOwner(cmd)
					if err != nil {
						return err
					}
					app, err := resolveApp(cmd)
					if err != nil {
						return err
					}
					var body any
					if err := json.Unmarshal([]byte(cmd.String("json")), &body); err != nil {
						return internal.NewValidationError(fmt.Sprintf("invalid JSON: %v", err))
					}
					raw, err := c.Request(ctx, "POST", bulkBase(c.CreatorBase, owner, app)+"/report/"+cmd.String("report")+"/write", &zohttp.RequestOpts{JSON: body})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "status",
				Usage:     "Get bulk write job status",
				ArgsUsage: "<job-id>",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "report", Required: true, Usage: "Report link name"},
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					id := cmd.Args().Get(0)
					if id == "" {
						return internal.NewValidationError("job-id argument required")
					}
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					owner, err := resolveOwner(cmd)
					if err != nil {
						return err
					}
					app, err := resolveApp(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request(ctx, "GET", bulkBase(c.CreatorBase, owner, app)+"/report/"+cmd.String("report")+"/write/"+id, nil)
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
		},
	}
}
