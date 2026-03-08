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

func milestonesCmd() *cli.Command {
	return &cli.Command{
		Name:  "milestones",
		Usage: "Project milestone operations",
		Commands: []*cli.Command{
			{
				Name:  "list",
				Usage: "List milestones",
				Flags: []cli.Flag{portalFlag, projectFlag, allFlag, limitFlag},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					portal, err := requirePortal(cmd)
					if err != nil {
						return err
					}
					url := base(c, portal, cmd.String("project")) + "/milestones"
					return paginateProjectsList(ctx, c, cmd, url, "milestones", nil)
				},
			},
			{
				Name:      "get",
				Usage:     "Get a milestone",
				ArgsUsage: "<milestone-id>",
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
					url := base(c, portal, cmd.String("project")) + "/milestones/" + cmd.Args().First()
					raw, err := c.Request(ctx, "GET", url, nil)
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:  "create",
				Usage: "Create a milestone",
				Flags: []cli.Flag{
					portalFlag, projectFlag,
					&cli.StringFlag{Name: "name", Required: true, Usage: "Milestone name"},
					&cli.StringFlag{Name: "start", Required: true, Usage: "Start date (YYYY-MM-DD)"},
					&cli.StringFlag{Name: "end", Required: true, Usage: "End date (YYYY-MM-DD)"},
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
					body := map[string]any{
						"name":       cmd.String("name"),
						"start_date": convertDate(cmd.String("start")),
						"end_date":   convertDate(cmd.String("end")),
					}
					if err := internal.MergeJSON(cmd, body); err != nil {
						return err
					}
					url := base(c, portal, cmd.String("project")) + "/milestones"
					raw, err := c.Request(ctx, "POST", url, &zohttp.RequestOpts{JSON: body})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "update",
				Usage:     "Update a milestone",
				ArgsUsage: "<milestone-id>",
				Flags: []cli.Flag{
					portalFlag, projectFlag,
					&cli.StringFlag{Name: "name", Usage: "Milestone name"},
					&cli.StringFlag{Name: "start", Usage: "Start date (YYYY-MM-DD)"},
					&cli.StringFlag{Name: "end", Usage: "End date (YYYY-MM-DD)"},
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
					if cmd.IsSet("start") {
						body["start_date"] = convertDate(cmd.String("start"))
					}
					if cmd.IsSet("end") {
						body["end_date"] = convertDate(cmd.String("end"))
					}
					if err := internal.MergeJSON(cmd, body); err != nil {
						return err
					}
					url := base(c, portal, cmd.String("project")) + "/milestones/" + cmd.Args().First()
					raw, err := c.Request(ctx, "POST", url, &zohttp.RequestOpts{JSON: body})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "delete",
				Usage:     "Delete a milestone",
				ArgsUsage: "<milestone-id>",
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
					url := base(c, portal, cmd.String("project")) + "/milestones/" + cmd.Args().First()
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

func phasesCmd() *cli.Command {
	return &cli.Command{
		Name:  "phases",
		Usage: "Phase operations",
		Commands: []*cli.Command{
			{
				Name:  "list",
				Usage: "List phases (portal-level)",
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
					url := c.ProjectsBase + "/portal/" + portal + "/phases"
					raw, err := c.Request(ctx, "GET", url, nil)
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:  "list-project",
				Usage: "List phases (project-level)",
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
					url := base(c, portal, cmd.String("project")) + "/phases"
					raw, err := c.Request(ctx, "GET", url, nil)
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "get",
				Usage:     "Get a phase",
				ArgsUsage: "<phase-id>",
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
					url := base(c, portal, cmd.String("project")) + "/phases/" + cmd.Args().First()
					raw, err := c.Request(ctx, "GET", url, nil)
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:  "create",
				Usage: "Create a phase",
				Flags: []cli.Flag{
					portalFlag, projectFlag,
					&cli.StringFlag{Name: "name", Required: true, Usage: "Phase name"},
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
					url := base(c, portal, cmd.String("project")) + "/phases"
					raw, err := c.Request(ctx, "POST", url, &zohttp.RequestOpts{JSON: body})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "update",
				Usage:     "Update a phase",
				ArgsUsage: "<phase-id>",
				Flags: []cli.Flag{
					portalFlag, projectFlag,
					&cli.StringFlag{Name: "name", Usage: "Phase name"},
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
					url := base(c, portal, cmd.String("project")) + "/phases/" + cmd.Args().First()
					raw, err := c.Request(ctx, "POST", url, &zohttp.RequestOpts{JSON: body})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "delete",
				Usage:     "Delete a phase",
				ArgsUsage: "<phase-id>",
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
					url := base(c, portal, cmd.String("project")) + "/phases/" + cmd.Args().First()
					raw, err := c.Request(ctx, "DELETE", url, nil)
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "move",
				Usage:     "Move a phase",
				ArgsUsage: "<phase-id>",
				Flags: []cli.Flag{
					portalFlag, projectFlag,
					&cli.StringFlag{Name: "to_project", Usage: "Target project ID"},
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
					if cmd.IsSet("to_project") {
						body["to_project"] = cmd.String("to_project")
					}
					if err := internal.MergeJSON(cmd, body); err != nil {
						return err
					}
					url := base(c, portal, cmd.String("project")) + "/phases/" + cmd.Args().First() + "/move"
					raw, err := c.Request(ctx, "POST", url, &zohttp.RequestOpts{JSON: body})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "clone",
				Usage:     "Clone a phase",
				ArgsUsage: "<phase-id>",
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
					url := base(c, portal, cmd.String("project")) + "/phases/" + cmd.Args().First() + "/clone"
					raw, err := c.Request(ctx, "POST", url, nil)
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "activities",
				Usage:     "Get phase activities",
				ArgsUsage: "<phase-id>",
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
					url := base(c, portal, cmd.String("project")) + "/phases/" + cmd.Args().First() + "/activities"
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

func phaseFollowersCmd() *cli.Command {
	return &cli.Command{
		Name:  "phase-followers",
		Usage: "Phase follower operations",
		Commands: []*cli.Command{
			{
				Name:      "list",
				Usage:     "List phase followers",
				ArgsUsage: "<phase-id>",
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
					url := base(c, portal, cmd.String("project")) + "/phases/" + cmd.Args().First() + "/followers"
					raw, err := c.Request(ctx, "GET", url, &zohttp.RequestOpts{Params: map[string]string{"page": "1", "per_page": "200"}})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "add",
				Usage:     "Add phase followers",
				ArgsUsage: "<phase-id>",
				Flags: []cli.Flag{
					portalFlag, projectFlag,
					&cli.StringFlag{Name: "followers", Usage: "Comma-separated follower ZPUIDs"},
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
					if cmd.IsSet("followers") {
						var arr any
						if err := json.Unmarshal([]byte(cmd.String("followers")), &arr); err != nil {
							return internal.NewValidationError(fmt.Sprintf("--followers: invalid JSON: %v", err))
						}
						body["followers"] = arr
					}
					if err := internal.MergeJSON(cmd, body); err != nil {
						return err
					}
					url := base(c, portal, cmd.String("project")) + "/phases/" + cmd.Args().First() + "/follow"
					raw, err := c.Request(ctx, "POST", url, &zohttp.RequestOpts{JSON: body})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "remove",
				Usage:     "Remove phase followers",
				ArgsUsage: "<phase-id>",
				Flags: []cli.Flag{
					portalFlag, projectFlag,
					&cli.StringFlag{Name: "followers", Usage: "Comma-separated follower ZPUIDs"},
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
					if cmd.IsSet("followers") {
						var arr any
						if err := json.Unmarshal([]byte(cmd.String("followers")), &arr); err != nil {
							return internal.NewValidationError(fmt.Sprintf("--followers: invalid JSON: %v", err))
						}
						body["followers"] = arr
					}
					if err := internal.MergeJSON(cmd, body); err != nil {
						return err
					}
					url := base(c, portal, cmd.String("project")) + "/phases/" + cmd.Args().First() + "/unfollow"
					raw, err := c.Request(ctx, "DELETE", url, &zohttp.RequestOpts{JSON: body})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
		},
	}
}

func phaseCommentsCmd() *cli.Command {
	return &cli.Command{
		Name:  "phase-comments",
		Usage: "Phase comment operations",
		Commands: []*cli.Command{
			{
				Name:  "list",
				Usage: "List phase comments",
				Flags: []cli.Flag{
					portalFlag, projectFlag,
					&cli.StringFlag{Name: "phase", Required: true, Usage: "Phase ID"},
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
					url := base(c, portal, cmd.String("project")) + "/phases/" + cmd.String("phase") + "/comments"
					raw, err := c.Request(ctx, "GET", url, nil)
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:  "add",
				Usage: "Add a phase comment",
				Flags: []cli.Flag{
					portalFlag, projectFlag,
					&cli.StringFlag{Name: "phase", Required: true, Usage: "Phase ID"},
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
					url := base(c, portal, cmd.String("project")) + "/phases/" + cmd.String("phase") + "/comments"
					raw, err := c.Request(ctx, "POST", url, &zohttp.RequestOpts{JSON: map[string]string{"content": cmd.String("comment")}})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "update",
				Usage:     "Update a phase comment",
				ArgsUsage: "<comment-id>",
				Flags: []cli.Flag{
					portalFlag, projectFlag,
					&cli.StringFlag{Name: "phase", Required: true, Usage: "Phase ID"},
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
					url := base(c, portal, cmd.String("project")) + "/phases/" + cmd.String("phase") + "/comments/" + cmd.Args().First()
					raw, err := c.Request(ctx, "PATCH", url, &zohttp.RequestOpts{JSON: map[string]string{"content": cmd.String("comment")}})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "delete",
				Usage:     "Delete a phase comment",
				ArgsUsage: "<comment-id>",
				Flags: []cli.Flag{
					portalFlag, projectFlag,
					&cli.StringFlag{Name: "phase", Required: true, Usage: "Phase ID"},
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
					url := base(c, portal, cmd.String("project")) + "/phases/" + cmd.String("phase") + "/comments/" + cmd.Args().First()
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

func dependenciesCmd() *cli.Command {
	return &cli.Command{
		Name:  "dependencies",
		Usage: "Task dependency operations",
		Commands: []*cli.Command{
			{
				Name:      "add",
				Usage:     "Add a task dependency",
				ArgsUsage: "<task-id>",
				Flags: []cli.Flag{
					portalFlag, projectFlag,
					&cli.StringFlag{Name: "depends-on", Required: true, Usage: "Dependency task ID"},
					&cli.StringFlag{Name: "type", Value: "FS", Usage: "FS, SS, FF, or SF"},
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
					url := base(c, portal, cmd.String("project")) + "/tasks/" + cmd.Args().First() + "/dependencies"
					body := map[string]any{
						"predecessor": map[string]string{
							"id":   cmd.String("depends-on"),
							"type": cmd.String("type"),
						},
					}
					raw, err := c.Request(ctx, "POST", url, &zohttp.RequestOpts{JSON: body})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "remove",
				Usage:     "Remove a task dependency",
				ArgsUsage: "<task-id> <dependency-id>",
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
					url := base(c, portal, cmd.String("project")) + "/tasks/" + cmd.Args().Get(0) + "/dependencies/" + cmd.Args().Get(1)
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
