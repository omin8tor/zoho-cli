package projects

import (
	"context"
	"github.com/omin8tor/zoho-cli/internal"
	zohttp "github.com/omin8tor/zoho-cli/internal/http"
	"github.com/omin8tor/zoho-cli/internal/output"
	"github.com/urfave/cli/v3"
	"io"
	"os"
	"path/filepath"
)

func eventsCmd() *cli.Command {
	return &cli.Command{
		Name:  "events",
		Usage: "Event operations",
		Commands: []*cli.Command{
			{
				Name:  "list",
				Usage: "List events",
				Flags: []cli.Flag{portalFlag, projectFlag},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					portal, err := requirePortal(cmd)
					if err != nil {
						return err
					}
					url := base(c, portal, cmd.String("project")) + "/events"
					raw, err := c.Request(ctx, "GET", url, nil)
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "get",
				Usage:     "Get an event",
				ArgsUsage: "<event-id>",
				Flags:     []cli.Flag{portalFlag, projectFlag},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					portal, err := requirePortal(cmd)
					if err != nil {
						return err
					}
					url := base(c, portal, cmd.String("project")) + "/events/" + cmd.Args().First()
					raw, err := c.Request(ctx, "GET", url, nil)
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:  "create",
				Usage: "Create an event",
				Flags: []cli.Flag{
					portalFlag, projectFlag,
					&cli.StringFlag{Name: "title", Required: true, Usage: "Event title"},
					&cli.StringFlag{Name: "starts-at", Usage: "Start datetime (ISO 8601)"},
					&cli.StringFlag{Name: "ends-at", Usage: "End datetime (ISO 8601)"},
					&cli.StringFlag{Name: "json", Usage: "Additional fields as JSON"},
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					portal, err := requirePortal(cmd)
					if err != nil {
						return err
					}
					body := map[string]any{"title": cmd.String("title")}
					if cmd.IsSet("starts-at") {
						body["starts_at"] = cmd.String("starts-at")
					}
					if cmd.IsSet("ends-at") {
						body["ends_at"] = cmd.String("ends-at")
					}
					if err := internal.MergeJSON(cmd, body); err != nil {
						return err
					}
					url := base(c, portal, cmd.String("project")) + "/events"
					raw, err := c.Request(ctx, "POST", url, &zohttp.RequestOpts{JSON: body})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "update",
				Usage:     "Update an event",
				ArgsUsage: "<event-id>",
				Flags: []cli.Flag{
					portalFlag, projectFlag,
					&cli.StringFlag{Name: "title", Usage: "Event title"},
					&cli.StringFlag{Name: "starts-at", Usage: "Start datetime (ISO 8601)"},
					&cli.StringFlag{Name: "ends-at", Usage: "End datetime (ISO 8601)"},
					&cli.StringFlag{Name: "json", Usage: "Additional fields as JSON"},
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					portal, err := requirePortal(cmd)
					if err != nil {
						return err
					}
					body := map[string]any{}
					if cmd.IsSet("title") {
						body["title"] = cmd.String("title")
					}
					if cmd.IsSet("starts-at") {
						body["starts_at"] = cmd.String("starts-at")
					}
					if cmd.IsSet("ends-at") {
						body["ends_at"] = cmd.String("ends-at")
					}
					if err := internal.MergeJSON(cmd, body); err != nil {
						return err
					}
					url := base(c, portal, cmd.String("project")) + "/events/" + cmd.Args().First()
					raw, err := c.Request(ctx, "PATCH", url, &zohttp.RequestOpts{JSON: body})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "delete",
				Usage:     "Delete an event",
				ArgsUsage: "<event-id>",
				Flags:     []cli.Flag{portalFlag, projectFlag},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					portal, err := requirePortal(cmd)
					if err != nil {
						return err
					}
					url := base(c, portal, cmd.String("project")) + "/events/" + cmd.Args().First()
					raw, err := c.Request(ctx, "DELETE", url, nil)
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
		},
	}
}

func eventCommentsCmd() *cli.Command {
	return &cli.Command{
		Name:  "event-comments",
		Usage: "Event comment operations",
		Commands: []*cli.Command{
			{
				Name:  "list",
				Usage: "List event comments",
				Flags: []cli.Flag{
					portalFlag, projectFlag,
					&cli.StringFlag{Name: "event", Required: true, Usage: "Event ID"},
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					portal, err := requirePortal(cmd)
					if err != nil {
						return err
					}
					url := base(c, portal, cmd.String("project")) + "/events/" + cmd.String("event") + "/comments"
					raw, err := c.Request(ctx, "GET", url, nil)
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "get",
				Usage:     "Get an event comment",
				ArgsUsage: "<comment-id>",
				Flags: []cli.Flag{
					portalFlag, projectFlag,
					&cli.StringFlag{Name: "event", Required: true, Usage: "Event ID"},
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					portal, err := requirePortal(cmd)
					if err != nil {
						return err
					}
					url := base(c, portal, cmd.String("project")) + "/events/" + cmd.String("event") + "/comments/" + cmd.Args().First()
					raw, err := c.Request(ctx, "GET", url, nil)
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:  "add",
				Usage: "Add an event comment",
				Flags: []cli.Flag{
					portalFlag, projectFlag,
					&cli.StringFlag{Name: "event", Required: true, Usage: "Event ID"},
					&cli.StringFlag{Name: "comment", Required: true, Usage: "Comment text"},
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					portal, err := requirePortal(cmd)
					if err != nil {
						return err
					}
					url := base(c, portal, cmd.String("project")) + "/events/" + cmd.String("event") + "/comments"
					raw, err := c.Request(ctx, "POST", url, &zohttp.RequestOpts{JSON: map[string]string{"content": cmd.String("comment")}})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "update",
				Usage:     "Update an event comment",
				ArgsUsage: "<comment-id>",
				Flags: []cli.Flag{
					portalFlag, projectFlag,
					&cli.StringFlag{Name: "event", Required: true, Usage: "Event ID"},
					&cli.StringFlag{Name: "comment", Required: true, Usage: "Updated comment text"},
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					portal, err := requirePortal(cmd)
					if err != nil {
						return err
					}
					url := base(c, portal, cmd.String("project")) + "/events/" + cmd.String("event") + "/comments/" + cmd.Args().First()
					raw, err := c.Request(ctx, "PATCH", url, &zohttp.RequestOpts{JSON: map[string]string{"content": cmd.String("comment")}})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "delete",
				Usage:     "Delete an event comment",
				ArgsUsage: "<comment-id>",
				Flags: []cli.Flag{
					portalFlag, projectFlag,
					&cli.StringFlag{Name: "event", Required: true, Usage: "Event ID"},
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					portal, err := requirePortal(cmd)
					if err != nil {
						return err
					}
					url := base(c, portal, cmd.String("project")) + "/events/" + cmd.String("event") + "/comments/" + cmd.Args().First()
					raw, err := c.Request(ctx, "DELETE", url, nil)
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
		Usage: "Project attachment operations",
		Commands: []*cli.Command{
			{
				Name:  "list",
				Usage: "List project attachments",
				Flags: []cli.Flag{
					portalFlag, projectFlag,
					&cli.StringFlag{Name: "type", Required: true, Usage: "Entity type (task, bug, forum)"},
					&cli.StringFlag{Name: "entity-id", Required: true, Usage: "Entity ID"},
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					portal, err := requirePortal(cmd)
					if err != nil {
						return err
					}
					url := base(c, portal, cmd.String("project")) + "/attachments"
					raw, err := c.Request(ctx, "GET", url, &zohttp.RequestOpts{
						Params: map[string]string{
							"entity_type": cmd.String("type"),
							"entity_id":   cmd.String("entity-id"),
						},
					})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "get",
				Usage:     "Get a project attachment",
				ArgsUsage: "<attachment-id>",
				Flags:     []cli.Flag{portalFlag, projectFlag},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					portal, err := requirePortal(cmd)
					if err != nil {
						return err
					}
					url := base(c, portal, cmd.String("project")) + "/attachments/" + cmd.Args().First()
					raw, err := c.Request(ctx, "GET", url, nil)
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:  "upload",
				Usage: "Upload a file attachment",
				Flags: []cli.Flag{
					portalFlag,
					&cli.StringFlag{Name: "file", Required: true, Usage: "File path to upload"},
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					portal, err := requirePortal(cmd)
					if err != nil {
						return err
					}
					path := cmd.String("file")
					f, err := os.Open(path)
					if err != nil {
						return err
					}
					defer f.Close()
					fileBytes, err := io.ReadAll(f)
					if err != nil {
						return err
					}
					url := c.ProjectsBase + "/portal/" + portal + "/attachments"
					raw, err := c.Request(ctx, "POST", url, &zohttp.RequestOpts{
						Files: map[string]zohttp.FileUpload{
							"file": {Filename: filepath.Base(path), Data: fileBytes},
						},
					})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "associate",
				Usage:     "Associate an attachment to an entity",
				ArgsUsage: "<attachment-id>",
				Flags: []cli.Flag{
					portalFlag, projectFlag,
					&cli.StringFlag{Name: "type", Required: true, Usage: "Entity type (task, issue, etc.)"},
					&cli.StringFlag{Name: "entity-id", Required: true, Usage: "Entity ID"},
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					portal, err := requirePortal(cmd)
					if err != nil {
						return err
					}
					url := base(c, portal, cmd.String("project")) + "/attachments/" + cmd.Args().First()
					raw, err := c.Request(ctx, "POST", url, &zohttp.RequestOpts{
						Form: map[string]string{
							"entity_type": cmd.String("type"),
							"entity_id":   cmd.String("entity-id"),
						},
					})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "dissociate",
				Usage:     "Dissociate an attachment from a project",
				ArgsUsage: "<attachment-id>",
				Flags:     []cli.Flag{portalFlag, projectFlag},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					portal, err := requirePortal(cmd)
					if err != nil {
						return err
					}
					url := base(c, portal, cmd.String("project")) + "/attachments/" + cmd.Args().First()
					raw, err := c.Request(ctx, "DELETE", url, nil)
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
		},
	}
}
