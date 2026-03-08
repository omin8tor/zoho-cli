package projects

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/omin8tor/zoho-cli/internal"
	zohttp "github.com/omin8tor/zoho-cli/internal/http"
	"github.com/omin8tor/zoho-cli/internal/output"
	"github.com/urfave/cli/v3"
)

func leavesCmd() *cli.Command {
	return &cli.Command{
		Name:  "leaves",
		Usage: "Leave operations",
		Commands: []*cli.Command{
			{
				Name:  "list",
				Usage: "List leaves",
				Flags: []cli.Flag{portalFlag},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					portal, err := requirePortal(cmd)
					if err != nil {
						return err
					}
					url := c.ProjectsBase + "/portal/" + portal + "/leave"
					raw, err := c.Request(ctx, "GET", url, nil)
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "get",
				Usage:     "Get a leave",
				ArgsUsage: "<leave-id>",
				Flags:     []cli.Flag{portalFlag},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					portal, err := requirePortal(cmd)
					if err != nil {
						return err
					}
					url := c.ProjectsBase + "/portal/" + portal + "/leave/" + cmd.Args().First()
					raw, err := c.Request(ctx, "GET", url, nil)
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:  "create",
				Usage: "Create a leave",
				Flags: []cli.Flag{
					portalFlag,
					&cli.StringFlag{Name: "from-date", Usage: "From date (ISO 8601)"},
					&cli.StringFlag{Name: "to-date", Usage: "To date (ISO 8601)"},
					&cli.StringFlag{Name: "reason", Usage: "Reason for leave"},
					&cli.StringFlag{Name: "leave-type", Usage: "Leave type"},
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
					if cmd.IsSet("reason") {
						body["reason"] = cmd.String("reason")
					}
					if cmd.IsSet("from-date") {
						body["from_date"] = cmd.String("from-date")
					}
					if cmd.IsSet("to-date") {
						body["to_date"] = cmd.String("to-date")
					}
					if cmd.IsSet("leave-type") {
						body["type"] = cmd.String("leave-type")
					}
					if err := internal.MergeJSON(cmd, body); err != nil {
						return err
					}
					url := c.ProjectsBase + "/portal/" + portal + "/leave"
					raw, err := c.Request(ctx, "POST", url, &zohttp.RequestOpts{JSON: body})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "update",
				Usage:     "Update a leave",
				ArgsUsage: "<leave-id>",
				Flags: []cli.Flag{
					portalFlag,
					&cli.StringFlag{Name: "from-date", Usage: "From date (ISO 8601)"},
					&cli.StringFlag{Name: "to-date", Usage: "To date (ISO 8601)"},
					&cli.StringFlag{Name: "reason", Usage: "Reason for leave"},
					&cli.StringFlag{Name: "leave-type", Usage: "Leave type"},
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
					if cmd.IsSet("reason") {
						body["reason"] = cmd.String("reason")
					}
					if cmd.IsSet("from-date") {
						body["from_date"] = cmd.String("from-date")
					}
					if cmd.IsSet("to-date") {
						body["to_date"] = cmd.String("to-date")
					}
					if cmd.IsSet("leave-type") {
						body["type"] = cmd.String("leave-type")
					}
					if err := internal.MergeJSON(cmd, body); err != nil {
						return err
					}
					url := c.ProjectsBase + "/portal/" + portal + "/leave/" + cmd.Args().First()
					raw, err := c.Request(ctx, "PATCH", url, &zohttp.RequestOpts{JSON: body})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "delete",
				Usage:     "Delete a leave",
				ArgsUsage: "<leave-id>",
				Flags:     []cli.Flag{portalFlag},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					portal, err := requirePortal(cmd)
					if err != nil {
						return err
					}
					url := c.ProjectsBase + "/portal/" + portal + "/leave/" + cmd.Args().First()
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

func tagsCmd() *cli.Command {
	return &cli.Command{
		Name:  "tags",
		Usage: "Tag operations",
		Commands: []*cli.Command{
			{
				Name:  "list",
				Usage: "List tags (portal-level)",
				Flags: []cli.Flag{portalFlag},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					portal, err := requirePortal(cmd)
					if err != nil {
						return err
					}
					url := c.ProjectsBase + "/portal/" + portal + "/tags"
					raw, err := c.Request(ctx, "GET", url, nil)
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:  "project-list",
				Usage: "List tags (project-level)",
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
					url := base(c, portal, cmd.String("project")) + "/tags"
					raw, err := c.Request(ctx, "GET", url, nil)
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:  "create",
				Usage: "Create a tag",
				Flags: []cli.Flag{
					portalFlag,
					&cli.StringFlag{Name: "tags", Required: true, Usage: "Tags as JSON array"},
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
					form := map[string]string{"tags": cmd.String("tags")}
					if err := internal.MergeJSONForm(cmd, form); err != nil {
						return err
					}
					url := c.ProjectsBase + "/portal/" + portal + "/tags"
					raw, err := c.Request(ctx, "POST", url, &zohttp.RequestOpts{Form: form})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "update",
				Usage:     "Update a tag",
				ArgsUsage: "<tag-id>",
				Flags: []cli.Flag{
					portalFlag,
					&cli.StringFlag{Name: "name", Usage: "Tag name"},
					&cli.StringFlag{Name: "color-class", Usage: "Color class (e.g. bg-tag1)"},
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
					if cmd.IsSet("name") {
						body["name"] = cmd.String("name")
					}
					if cmd.IsSet("color-class") {
						body["color_class"] = cmd.String("color-class")
					}
					if err := internal.MergeJSON(cmd, body); err != nil {
						return err
					}
					url := c.ProjectsBase + "/portal/" + portal + "/tags/" + cmd.Args().First()
					raw, err := c.Request(ctx, "PATCH", url, &zohttp.RequestOpts{JSON: body})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "delete",
				Usage:     "Delete a tag",
				ArgsUsage: "<tag-id>",
				Flags:     []cli.Flag{portalFlag},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					portal, err := requirePortal(cmd)
					if err != nil {
						return err
					}
					url := c.ProjectsBase + "/portal/" + portal + "/tags/" + cmd.Args().First()
					raw, err := c.Request(ctx, "DELETE", url, nil)
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "associate",
				Usage:     "Associate a tag to entities",
				ArgsUsage: "<tag-id>",
				Flags: []cli.Flag{
					portalFlag, projectFlag,
					&cli.StringFlag{Name: "entity", Required: true, Usage: "Entity ID"},
					&cli.StringFlag{Name: "entity-type", Required: true, Usage: "Entity type (1=Project,2=Milestone,3=Tasklist,5=Task,6=Issue)"},
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
					url := base(c, portal, cmd.String("project")) + "/tags/associate"
					raw, err := c.Request(ctx, "POST", url, &zohttp.RequestOpts{Form: map[string]string{
						"tag_id":     cmd.Args().First(),
						"entity_id":  cmd.String("entity"),
						"entityType": cmd.String("entity-type"),
					}})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "dissociate",
				Usage:     "Dissociate a tag from entities",
				ArgsUsage: "<tag-id>",
				Flags: []cli.Flag{
					portalFlag, projectFlag,
					&cli.StringFlag{Name: "entity", Required: true, Usage: "Entity ID"},
					&cli.StringFlag{Name: "entity-type", Required: true, Usage: "Entity type (1=Project,2=Milestone,3=Tasklist,5=Task,6=Issue)"},
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
					url := base(c, portal, cmd.String("project")) + "/tags/dissociate"
					raw, err := c.Request(ctx, "POST", url, &zohttp.RequestOpts{Form: map[string]string{
						"tag_id":     cmd.Args().First(),
						"entity_id":  cmd.String("entity"),
						"entityType": cmd.String("entity-type"),
					}})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
		},
	}
}

func portalsCmd() *cli.Command {
	return &cli.Command{
		Name:  "portals",
		Usage: "Portal operations",
		Commands: []*cli.Command{
			{
				Name:  "get",
				Usage: "Get a portal",
				Flags: []cli.Flag{portalFlag},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					portal, err := requirePortal(cmd)
					if err != nil {
						return err
					}
					url := c.ProjectsBase + "/portal/" + portal
					raw, err := c.Request(ctx, "GET", url, nil)
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
		},
	}
}

func trashCmd() *cli.Command {
	return &cli.Command{
		Name:  "trash",
		Usage: "Trash operations",
		Commands: []*cli.Command{
			{
				Name:  "list",
				Usage: "List trash items",
				Flags: []cli.Flag{portalFlag},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					portal, err := requirePortal(cmd)
					if err != nil {
						return err
					}
					url := c.ProjectsBase + "/portal/" + portal + "/bin"
					raw, err := c.Request(ctx, "GET", url, nil)
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:  "delete",
				Usage: "Permanently delete trash items",
				Flags: []cli.Flag{
					portalFlag,
					&cli.StringFlag{Name: "ids", Required: true, Usage: "Record IDs as JSON array"},
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
					var body any
					if err := json.Unmarshal([]byte(cmd.String("ids")), &body); err != nil {
						return internal.NewValidationError(fmt.Sprintf("--ids: invalid JSON: %v", err))
					}
					url := c.ProjectsBase + "/portal/" + portal + "/bin"
					raw, err := c.Request(ctx, "DELETE", url, &zohttp.RequestOpts{JSON: body})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:  "restore",
				Usage: "Restore trash items",
				Flags: []cli.Flag{
					portalFlag,
					&cli.StringFlag{Name: "ids", Required: true, Usage: "Record IDs as JSON array"},
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
					var body any
					if err := json.Unmarshal([]byte(cmd.String("ids")), &body); err != nil {
						return internal.NewValidationError(fmt.Sprintf("--ids: invalid JSON: %v", err))
					}
					url := c.ProjectsBase + "/portal/" + portal + "/bin/restore"
					raw, err := c.Request(ctx, "POST", url, &zohttp.RequestOpts{JSON: body})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:  "empty",
				Usage: "Empty all trash",
				Flags: []cli.Flag{portalFlag},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					portal, err := requirePortal(cmd)
					if err != nil {
						return err
					}
					url := c.ProjectsBase + "/portal/" + portal + "/empty-bin"
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

func searchCmd() *cli.Command {
	return &cli.Command{
		Name:  "search",
		Usage: "Search operations",
		Commands: []*cli.Command{
			{
				Name:  "portal",
				Usage: "Search across portal",
				Flags: []cli.Flag{
					portalFlag,
					&cli.StringFlag{Name: "query", Required: true, Usage: "Search query"},
					&cli.StringFlag{Name: "module", Value: "all", Usage: "Module to search"},
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
					url := c.ProjectsBase + "/portal/" + portal + "/search"
					raw, err := c.Request(ctx, "GET", url, &zohttp.RequestOpts{
						Params: map[string]string{
							"search_term": cmd.String("query"),
							"module":      cmd.String("module"),
						},
					})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:  "project",
				Usage: "Search within a project",
				Flags: []cli.Flag{
					portalFlag, projectFlag,
					&cli.StringFlag{Name: "query", Required: true, Usage: "Search query"},
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
					url := base(c, portal, cmd.String("project")) + "/search"
					raw, err := c.Request(ctx, "GET", url, &zohttp.RequestOpts{
						Params: map[string]string{
							"search_term": cmd.String("query"),
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

func feedCmd() *cli.Command {
	return &cli.Command{
		Name:  "feed",
		Usage: "Feed/status operations",
		Commands: []*cli.Command{
			{
				Name:  "status",
				Usage: "Get project status feed",
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
					url := base(c, portal, cmd.String("project")) + "/status"
					raw, err := c.Request(ctx, "GET", url, nil)
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:  "post",
				Usage: "Post a status update",
				Flags: []cli.Flag{
					portalFlag, projectFlag,
					&cli.StringFlag{Name: "content", Required: true, Usage: "Status content"},
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
					body := map[string]any{"content": cmd.String("content")}
					if err := internal.MergeJSON(cmd, body); err != nil {
						return err
					}
					url := base(c, portal, cmd.String("project")) + "/status"
					raw, err := c.Request(ctx, "POST", url, &zohttp.RequestOpts{JSON: body})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
		},
	}
}

func projectCommentsCmd() *cli.Command {
	return &cli.Command{
		Name:  "project-comments",
		Usage: "Project comment operations",
		Commands: []*cli.Command{
			{
				Name:  "list",
				Usage: "List project comments",
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
					url := base(c, portal, cmd.String("project")) + "/comments"
					raw, err := c.Request(ctx, "GET", url, nil)
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "get",
				Usage:     "Get a project comment",
				ArgsUsage: "<comment-id>",
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
					url := base(c, portal, cmd.String("project")) + "/comments/" + cmd.Args().First()
					raw, err := c.Request(ctx, "GET", url, nil)
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:  "add",
				Usage: "Add a project comment",
				Flags: []cli.Flag{
					portalFlag, projectFlag,
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
					url := base(c, portal, cmd.String("project")) + "/comments"
					raw, err := c.Request(ctx, "POST", url, &zohttp.RequestOpts{JSON: map[string]string{"content": cmd.String("comment")}})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "update",
				Usage:     "Update a project comment",
				ArgsUsage: "<comment-id>",
				Flags: []cli.Flag{
					portalFlag, projectFlag,
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
					url := base(c, portal, cmd.String("project")) + "/comments/" + cmd.Args().First()
					raw, err := c.Request(ctx, "PATCH", url, &zohttp.RequestOpts{JSON: map[string]string{"content": cmd.String("comment")}})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "delete",
				Usage:     "Delete a project comment",
				ArgsUsage: "<comment-id>",
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
					url := base(c, portal, cmd.String("project")) + "/comments/" + cmd.Args().First()
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

func projectGroupsCmd() *cli.Command {
	return &cli.Command{
		Name:  "project-groups",
		Usage: "Project group operations",
		Commands: []*cli.Command{
			{
				Name:  "list",
				Usage: "List project groups",
				Flags: []cli.Flag{portalFlag},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					portal, err := requirePortal(cmd)
					if err != nil {
						return err
					}
					url := c.ProjectsBase + "/portal/" + portal + "/project-groups"
					raw, err := c.Request(ctx, "GET", url, nil)
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:  "my",
				Usage: "List my project groups",
				Flags: []cli.Flag{portalFlag},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					portal, err := requirePortal(cmd)
					if err != nil {
						return err
					}
					url := c.ProjectsBase + "/portal/" + portal + "/my-project-groups"
					raw, err := c.Request(ctx, "GET", url, nil)
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:  "create",
				Usage: "Create a project group",
				Flags: []cli.Flag{
					portalFlag,
					&cli.StringFlag{Name: "name", Required: true, Usage: "Group name"},
					&cli.StringFlag{Name: "group-type", Usage: "Group type (public or private)"},
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
					body := map[string]any{"name": cmd.String("name")}
					if cmd.IsSet("group-type") {
						body["type"] = cmd.String("group-type")
					}
					if err := internal.MergeJSON(cmd, body); err != nil {
						return err
					}
					url := c.ProjectsBase + "/portal/" + portal + "/project-groups"
					raw, err := c.Request(ctx, "POST", url, &zohttp.RequestOpts{JSON: body})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "update",
				Usage:     "Update a project group",
				ArgsUsage: "<group-id>",
				Flags: []cli.Flag{
					portalFlag,
					&cli.StringFlag{Name: "name", Usage: "Group name"},
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
					if cmd.IsSet("name") {
						body["name"] = cmd.String("name")
					}
					if err := internal.MergeJSON(cmd, body); err != nil {
						return err
					}
					url := c.ProjectsBase + "/portal/" + portal + "/project-groups/" + cmd.Args().First()
					raw, err := c.Request(ctx, "PATCH", url, &zohttp.RequestOpts{JSON: body})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "delete",
				Usage:     "Delete a project group",
				ArgsUsage: "<group-id>",
				Flags:     []cli.Flag{portalFlag},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					portal, err := requirePortal(cmd)
					if err != nil {
						return err
					}
					url := c.ProjectsBase + "/portal/" + portal + "/project-groups/" + cmd.Args().First()
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

func teamsCmd() *cli.Command {
	return &cli.Command{
		Name:  "teams",
		Usage: "Team operations",
		Commands: []*cli.Command{
			{
				Name:  "list",
				Usage: "List all teams",
				Flags: []cli.Flag{portalFlag},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					portal, err := requirePortal(cmd)
					if err != nil {
						return err
					}
					url := c.ProjectsBase + "/portal/" + portal + "/teams"
					raw, err := c.Request(ctx, "GET", url, nil)
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "get",
				Usage:     "Get a team",
				ArgsUsage: "<team-id>",
				Flags:     []cli.Flag{portalFlag},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					portal, err := requirePortal(cmd)
					if err != nil {
						return err
					}
					url := c.ProjectsBase + "/portal/" + portal + "/teams/" + cmd.Args().First()
					raw, err := c.Request(ctx, "GET", url, nil)
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:  "project-list",
				Usage: "List teams in a project",
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
					url := base(c, portal, cmd.String("project")) + "/teams"
					raw, err := c.Request(ctx, "GET", url, nil)
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "users",
				Usage:     "List team users",
				ArgsUsage: "<team-id>",
				Flags:     []cli.Flag{portalFlag},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					portal, err := requirePortal(cmd)
					if err != nil {
						return err
					}
					url := c.ProjectsBase + "/portal/" + portal + "/teams/" + cmd.Args().First() + "/users"
					raw, err := c.Request(ctx, "GET", url, nil)
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "projects",
				Usage:     "List team projects",
				ArgsUsage: "<team-id>",
				Flags:     []cli.Flag{portalFlag},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					portal, err := requirePortal(cmd)
					if err != nil {
						return err
					}
					url := c.ProjectsBase + "/portal/" + portal + "/teams/" + cmd.Args().First() + "/projects"
					raw, err := c.Request(ctx, "GET", url, nil)
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:  "create",
				Usage: "Create a team",
				Flags: []cli.Flag{
					portalFlag,
					&cli.StringFlag{Name: "name", Required: true, Usage: "Team name"},
					&cli.StringFlag{Name: "lead", Usage: "Team lead ZPUID"},
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
					body := map[string]any{"name": cmd.String("name")}
					if cmd.IsSet("lead") {
						body["lead"] = cmd.String("lead")
					}
					if err := internal.MergeJSON(cmd, body); err != nil {
						return err
					}
					url := c.ProjectsBase + "/portal/" + portal + "/teams"
					raw, err := c.Request(ctx, "POST", url, &zohttp.RequestOpts{JSON: body})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "update",
				Usage:     "Update a team",
				ArgsUsage: "<team-id>",
				Flags: []cli.Flag{
					portalFlag,
					&cli.StringFlag{Name: "name", Usage: "Team name"},
					&cli.StringFlag{Name: "lead", Usage: "Team lead ZPUID"},
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
					if cmd.IsSet("name") {
						body["name"] = cmd.String("name")
					}
					if cmd.IsSet("lead") {
						body["lead"] = cmd.String("lead")
					}
					if err := internal.MergeJSON(cmd, body); err != nil {
						return err
					}
					url := c.ProjectsBase + "/portal/" + portal + "/teams/" + cmd.Args().First()
					raw, err := c.Request(ctx, "PATCH", url, &zohttp.RequestOpts{JSON: body})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "delete",
				Usage:     "Delete a team",
				ArgsUsage: "<team-id>",
				Flags:     []cli.Flag{portalFlag},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					portal, err := requirePortal(cmd)
					if err != nil {
						return err
					}
					url := c.ProjectsBase + "/portal/" + portal + "/teams/" + cmd.Args().First()
					raw, err := c.Request(ctx, "DELETE", url, nil)
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:  "add-to-project",
				Usage: "Add a team to a project",
				Flags: []cli.Flag{
					portalFlag, projectFlag,
					&cli.StringFlag{Name: "team-ids", Required: true, Usage: "Team IDs as JSON array"},
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
					form := map[string]string{"team_ids": cmd.String("team-ids")}
					if err := internal.MergeJSONForm(cmd, form); err != nil {
						return err
					}
					url := base(c, portal, cmd.String("project")) + "/associate-teams"
					raw, err := c.Request(ctx, "POST", url, &zohttp.RequestOpts{
						Form: form,
					})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "remove-from-project",
				Usage:     "Remove a team from a project",
				ArgsUsage: "<team-id>",
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
					url := base(c, portal, cmd.String("project")) + "/dissociate-teams/" + cmd.Args().First()
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

func profilesCmd() *cli.Command {
	return &cli.Command{
		Name:  "profiles",
		Usage: "Profile operations",
		Commands: []*cli.Command{
			{
				Name:  "list",
				Usage: "List profiles",
				Flags: []cli.Flag{portalFlag},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					portal, err := requirePortal(cmd)
					if err != nil {
						return err
					}
					url := c.ProjectsBase + "/portal/" + portal + "/profiles"
					raw, err := c.Request(ctx, "GET", url, nil)
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "get",
				Usage:     "Get a profile",
				ArgsUsage: "<profile-id>",
				Flags:     []cli.Flag{portalFlag},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					portal, err := requirePortal(cmd)
					if err != nil {
						return err
					}
					url := c.ProjectsBase + "/portal/" + portal + "/profiles/" + cmd.Args().First()
					raw, err := c.Request(ctx, "GET", url, nil)
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:  "create",
				Usage: "Create a profile",
				Flags: []cli.Flag{
					portalFlag,
					&cli.StringFlag{Name: "name", Required: true, Usage: "Profile name"},
					&cli.StringFlag{Name: "profile-type", Usage: "Profile type"},
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
					body := map[string]any{"name": cmd.String("name")}
					if cmd.IsSet("profile-type") {
						body["type"] = cmd.String("profile-type")
					}
					if err := internal.MergeJSON(cmd, body); err != nil {
						return err
					}
					url := c.ProjectsBase + "/portal/" + portal + "/profiles"
					raw, err := c.Request(ctx, "POST", url, &zohttp.RequestOpts{JSON: body})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "update",
				Usage:     "Update a profile",
				ArgsUsage: "<profile-id>",
				Flags: []cli.Flag{
					portalFlag,
					&cli.StringFlag{Name: "name", Usage: "Profile name"},
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
					if cmd.IsSet("name") {
						body["name"] = cmd.String("name")
					}
					if err := internal.MergeJSON(cmd, body); err != nil {
						return err
					}
					url := c.ProjectsBase + "/portal/" + portal + "/profiles/" + cmd.Args().First()
					raw, err := c.Request(ctx, "PATCH", url, &zohttp.RequestOpts{JSON: body})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "set-default",
				Usage:     "Set a profile as default",
				ArgsUsage: "<profile-id>",
				Flags:     []cli.Flag{portalFlag},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					portal, err := requirePortal(cmd)
					if err != nil {
						return err
					}
					url := c.ProjectsBase + "/portal/" + portal + "/profiles/" + cmd.Args().First() + "/setprimary"
					raw, err := c.Request(ctx, "POST", url, nil)
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "delete",
				Usage:     "Delete a profile",
				ArgsUsage: "<profile-id>",
				Flags:     []cli.Flag{portalFlag},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					portal, err := requirePortal(cmd)
					if err != nil {
						return err
					}
					url := c.ProjectsBase + "/portal/" + portal + "/profiles/" + cmd.Args().First()
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

func rolesCmd() *cli.Command {
	return &cli.Command{
		Name:  "roles",
		Usage: "Role operations",
		Commands: []*cli.Command{
			{
				Name:  "list",
				Usage: "List roles",
				Flags: []cli.Flag{portalFlag},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					portal, err := requirePortal(cmd)
					if err != nil {
						return err
					}
					url := c.ProjectsBase + "/portal/" + portal + "/roles"
					raw, err := c.Request(ctx, "GET", url, nil)
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "get",
				Usage:     "Get a role",
				ArgsUsage: "<role-id>",
				Flags:     []cli.Flag{portalFlag},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					portal, err := requirePortal(cmd)
					if err != nil {
						return err
					}
					url := c.ProjectsBase + "/portal/" + portal + "/roles/" + cmd.Args().First()
					raw, err := c.Request(ctx, "GET", url, nil)
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:  "create",
				Usage: "Create a role",
				Flags: []cli.Flag{
					portalFlag,
					&cli.StringFlag{Name: "name", Required: true, Usage: "Role name"},
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
					body := map[string]any{"name": cmd.String("name")}
					if err := internal.MergeJSON(cmd, body); err != nil {
						return err
					}
					url := c.ProjectsBase + "/portal/" + portal + "/roles"
					raw, err := c.Request(ctx, "POST", url, &zohttp.RequestOpts{JSON: body})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "update",
				Usage:     "Update a role",
				ArgsUsage: "<role-id>",
				Flags: []cli.Flag{
					portalFlag,
					&cli.StringFlag{Name: "name", Usage: "Role name"},
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
					if cmd.IsSet("name") {
						body["name"] = cmd.String("name")
					}
					if err := internal.MergeJSON(cmd, body); err != nil {
						return err
					}
					url := c.ProjectsBase + "/portal/" + portal + "/roles/" + cmd.Args().First()
					raw, err := c.Request(ctx, "PATCH", url, &zohttp.RequestOpts{JSON: body})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "set-default",
				Usage:     "Set a role as default",
				ArgsUsage: "<role-id>",
				Flags:     []cli.Flag{portalFlag},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					portal, err := requirePortal(cmd)
					if err != nil {
						return err
					}
					url := c.ProjectsBase + "/portal/" + portal + "/roles/" + cmd.Args().First() + "/default"
					raw, err := c.Request(ctx, "POST", url, nil)
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "delete",
				Usage:     "Delete a role",
				ArgsUsage: "<role-id>",
				Flags:     []cli.Flag{portalFlag},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					portal, err := requirePortal(cmd)
					if err != nil {
						return err
					}
					url := c.ProjectsBase + "/portal/" + portal + "/roles/" + cmd.Args().First()
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

func customRecordsCmd() *cli.Command {
	return &cli.Command{
		Name:  "custom-records",
		Usage: "Custom module record operations",
		Commands: []*cli.Command{
			{
				Name:  "list",
				Usage: "List custom records",
				Flags: []cli.Flag{
					portalFlag, projectFlag,
					&cli.StringFlag{Name: "module", Required: true, Usage: "Module name"},
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
					url := base(c, portal, cmd.String("project")) + "/modules/" + cmd.String("module") + "/records"
					raw, err := c.Request(ctx, "GET", url, nil)
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "get",
				Usage:     "Get a custom record",
				ArgsUsage: "<record-id>",
				Flags: []cli.Flag{
					portalFlag, projectFlag,
					&cli.StringFlag{Name: "module", Required: true, Usage: "Module name"},
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
					url := base(c, portal, cmd.String("project")) + "/modules/" + cmd.String("module") + "/records/" + cmd.Args().First()
					raw, err := c.Request(ctx, "GET", url, nil)
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:  "create",
				Usage: "Create a custom record",
				Flags: []cli.Flag{
					portalFlag, projectFlag,
					&cli.StringFlag{Name: "module", Required: true, Usage: "Module name"},
					&cli.StringFlag{Name: "name", Usage: "Record name"},
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
					if cmd.IsSet("name") {
						body["name"] = cmd.String("name")
					}
					if err := internal.MergeJSON(cmd, body); err != nil {
						return err
					}
					url := base(c, portal, cmd.String("project")) + "/modules/" + cmd.String("module") + "/records"
					raw, err := c.Request(ctx, "POST", url, &zohttp.RequestOpts{JSON: body})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "update",
				Usage:     "Update a custom record",
				ArgsUsage: "<record-id>",
				Flags: []cli.Flag{
					portalFlag, projectFlag,
					&cli.StringFlag{Name: "module", Required: true, Usage: "Module name"},
					&cli.StringFlag{Name: "name", Usage: "Record name"},
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
					if cmd.IsSet("name") {
						body["name"] = cmd.String("name")
					}
					if err := internal.MergeJSON(cmd, body); err != nil {
						return err
					}
					url := base(c, portal, cmd.String("project")) + "/modules/" + cmd.String("module") + "/records/" + cmd.Args().First()
					raw, err := c.Request(ctx, "PATCH", url, &zohttp.RequestOpts{JSON: body})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "trash",
				Usage:     "Move a custom record to trash",
				ArgsUsage: "<record-id>",
				Flags: []cli.Flag{
					portalFlag, projectFlag,
					&cli.StringFlag{Name: "module", Required: true, Usage: "Module name"},
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
					url := base(c, portal, cmd.String("project")) + "/modules/" + cmd.String("module") + "/records/" + cmd.Args().First() + "/trash"
					raw, err := c.Request(ctx, "POST", url, nil)
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "restore",
				Usage:     "Restore a custom record from trash",
				ArgsUsage: "<record-id>",
				Flags: []cli.Flag{
					portalFlag, projectFlag,
					&cli.StringFlag{Name: "module", Required: true, Usage: "Module name"},
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
					url := base(c, portal, cmd.String("project")) + "/modules/" + cmd.String("module") + "/records/" + cmd.Args().First() + "/restore"
					raw, err := c.Request(ctx, "POST", url, nil)
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "delete",
				Usage:     "Delete a custom record",
				ArgsUsage: "<record-id>",
				Flags: []cli.Flag{
					portalFlag, projectFlag,
					&cli.StringFlag{Name: "module", Required: true, Usage: "Module name"},
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
					url := base(c, portal, cmd.String("project")) + "/modules/" + cmd.String("module") + "/records/" + cmd.Args().First()
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

func reportsCmd() *cli.Command {
	return &cli.Command{
		Name:  "reports",
		Usage: "Report and dashboard operations",
		Commands: []*cli.Command{
			{
				Name:  "workload-meta",
				Usage: "Get workload report metadata",
				Flags: []cli.Flag{portalFlag},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					portal, err := requirePortal(cmd)
					if err != nil {
						return err
					}
					url := c.ProjectsBase + "/portal/" + portal + "/reports/workload/meta"
					raw, err := c.Request(ctx, "GET", url, nil)
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:  "workload",
				Usage: "Get workload report",
				Flags: []cli.Flag{portalFlag},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					portal, err := requirePortal(cmd)
					if err != nil {
						return err
					}
					url := c.ProjectsBase + "/portal/" + portal + "/reports/workload"
					raw, err := c.Request(ctx, "GET", url, nil)
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:  "dashboards",
				Usage: "List dashboards",
				Flags: []cli.Flag{portalFlag},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					portal, err := requirePortal(cmd)
					if err != nil {
						return err
					}
					url := c.ProjectsBase + "/portal/" + portal + "/dashboards"
					raw, err := c.Request(ctx, "GET", url, nil)
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "dashboard-get",
				Usage:     "Get a dashboard",
				ArgsUsage: "<dashboard-id>",
				Flags:     []cli.Flag{portalFlag},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					portal, err := requirePortal(cmd)
					if err != nil {
						return err
					}
					url := c.ProjectsBase + "/portal/" + portal + "/dashboards/" + cmd.Args().First()
					raw, err := c.Request(ctx, "GET", url, nil)
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:  "dashboard-create",
				Usage: "Create a dashboard",
				Flags: []cli.Flag{
					portalFlag,
					&cli.StringFlag{Name: "name", Required: true, Usage: "Dashboard name"},
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
					body := map[string]any{"name": cmd.String("name")}
					if err := internal.MergeJSON(cmd, body); err != nil {
						return err
					}
					url := c.ProjectsBase + "/portal/" + portal + "/dashboards"
					raw, err := c.Request(ctx, "POST", url, &zohttp.RequestOpts{JSON: body})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "dashboard-update",
				Usage:     "Update a dashboard",
				ArgsUsage: "<dashboard-id>",
				Flags: []cli.Flag{
					portalFlag,
					&cli.StringFlag{Name: "name", Usage: "Dashboard name"},
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
					if cmd.IsSet("name") {
						body["name"] = cmd.String("name")
					}
					if err := internal.MergeJSON(cmd, body); err != nil {
						return err
					}
					url := c.ProjectsBase + "/portal/" + portal + "/dashboards/" + cmd.Args().First()
					raw, err := c.Request(ctx, "PATCH", url, &zohttp.RequestOpts{JSON: body})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "dashboard-delete",
				Usage:     "Delete a dashboard",
				ArgsUsage: "<dashboard-id>",
				Flags:     []cli.Flag{portalFlag},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					portal, err := requirePortal(cmd)
					if err != nil {
						return err
					}
					url := c.ProjectsBase + "/portal/" + portal + "/dashboards/" + cmd.Args().First()
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
