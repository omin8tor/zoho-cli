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

func tasksCmd() *cli.Command {
	return &cli.Command{
		Name:  "tasks",
		Usage: "Project task operations",
		Commands: []*cli.Command{
			{
				Name:  "list",
				Usage: "List tasks in a project",
				Flags: []cli.Flag{
					portalFlag, projectFlag,
					allFlag, limitFlag,
					&cli.StringFlag{Name: "status", Usage: "Filter: open, closed, in progress"},
					&cli.StringFlag{Name: "priority", Usage: "Filter: none, low, medium, high"},
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
					url := base(c, portal, cmd.String("project")) + "/tasks"
					params := map[string]string{}
					if s := cmd.String("status"); s != "" {
						params["status"] = s
					}
					if p := cmd.String("priority"); p != "" {
						params["priority"] = p
					}
					return paginateProjectsList(ctx, c, cmd, url, "tasks", params)
				},
			},
			{
				Name:  "my",
				Usage: "List my tasks across all projects",
				Flags: []cli.Flag{
					portalFlag,
					allFlag, limitFlag,
					&cli.StringFlag{Name: "status", Usage: "Filter: open, closed, in progress"},
					&cli.StringFlag{Name: "priority", Usage: "Filter: none, low, medium, high"},
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
					url := c.ProjectsBase + "/portal/" + portal + "/tasks"
					params := map[string]string{}
					if s := cmd.String("status"); s != "" {
						params["status"] = s
					}
					if p := cmd.String("priority"); p != "" {
						params["priority"] = p
					}
					return paginateProjectsList(ctx, c, cmd, url, "tasks", params)
				},
			},
			{
				Name:      "get",
				Usage:     "Get a single task",
				ArgsUsage: "<task-id>",
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
					url := base(c, portal, cmd.String("project")) + "/tasks/" + cmd.Args().First()
					raw, err := c.Request(ctx, "GET", url, nil)
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:  "create",
				Usage: "Create a task",
				Flags: []cli.Flag{
					portalFlag, projectFlag,
					&cli.StringFlag{Name: "name", Required: true, Usage: "Task name"},
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
					url := base(c, portal, cmd.String("project")) + "/tasks"
					raw, err := c.Request(ctx, "POST", url, &zohttp.RequestOpts{JSON: body})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "update",
				Usage:     "Update a task",
				ArgsUsage: "<task-id>",
				Flags: []cli.Flag{
					portalFlag, projectFlag,
					&cli.StringFlag{Name: "name", Usage: "Task name"},
					&cli.StringFlag{Name: "description", Usage: "Task description"},
					&cli.StringFlag{Name: "priority", Usage: "Priority: none, low, medium, high"},
					&cli.StringFlag{Name: "start-date", Usage: "Start date (YYYY-MM-DD)"},
					&cli.StringFlag{Name: "end-date", Usage: "End date (YYYY-MM-DD)"},
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
					if cmd.IsSet("description") {
						body["description"] = cmd.String("description")
					}
					if cmd.IsSet("priority") {
						body["priority"] = cmd.String("priority")
					}
					if cmd.IsSet("start-date") {
						body["start_date"] = cmd.String("start-date")
					}
					if cmd.IsSet("end-date") {
						body["end_date"] = cmd.String("end-date")
					}
					if err := internal.MergeJSON(cmd, body); err != nil {
						return err
					}
					url := base(c, portal, cmd.String("project")) + "/tasks/" + cmd.Args().First()
					raw, err := c.Request(ctx, "PATCH", url, &zohttp.RequestOpts{JSON: body})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "delete",
				Usage:     "Delete a task",
				ArgsUsage: "<task-id>",
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
					url := base(c, portal, cmd.String("project")) + "/tasks/" + cmd.Args().First()
					raw, err := c.Request(ctx, "DELETE", url, nil)
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "subtasks",
				Usage:     "List subtasks of a task",
				ArgsUsage: "<task-id>",
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
					url := base(c, portal, cmd.String("project")) + "/tasks/" + cmd.Args().First() + "/subtasks"
					raw, err := c.Request(ctx, "GET", url, nil)
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:  "add-subtask",
				Usage: "Create a subtask",
				Flags: []cli.Flag{
					portalFlag, projectFlag,
					&cli.StringFlag{Name: "parent", Required: true, Usage: "Parent task ID"},
					&cli.StringFlag{Name: "name", Required: true, Usage: "Subtask name"},
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
						"name": cmd.String("name"),
						"parental_info": map[string]any{
							"parent_task_id": cmd.String("parent"),
						},
					}
					if err := internal.MergeJSON(cmd, body); err != nil {
						return err
					}
					url := base(c, portal, cmd.String("project")) + "/tasks"
					raw, err := c.Request(ctx, "POST", url, &zohttp.RequestOpts{JSON: body})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "clone",
				Usage:     "Clone a task",
				ArgsUsage: "<task-id>",
				Flags: []cli.Flag{
					portalFlag, projectFlag,
					&cli.StringFlag{Name: "instances", Value: "1", Usage: "Number of copies"},
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
					url := base(c, portal, cmd.String("project")) + "/tasks/" + cmd.Args().First() + "/copy"
					raw, err := c.Request(ctx, "POST", url, &zohttp.RequestOpts{
						Form: map[string]string{"no_of_instances": cmd.String("instances")},
					})
					if err != nil {
						return err
					}
					var envelope map[string]json.RawMessage
					if json.Unmarshal(raw, &envelope) == nil {
						if tasks, ok := envelope["tasks"]; ok {
							var arr []json.RawMessage
							if json.Unmarshal(tasks, &arr) == nil && len(arr) > 0 {
								return output.JSONRaw(arr[0])
							}
						}
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "move",
				Usage:     "Move a task",
				ArgsUsage: "<task-id>",
				Flags: []cli.Flag{
					portalFlag, projectFlag,
					&cli.StringFlag{Name: "target_tasklist_id", Usage: "Target tasklist ID"},
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
					if cmd.IsSet("target_tasklist_id") {
						body["target_tasklist_id"] = cmd.String("target_tasklist_id")
					}
					if err := internal.MergeJSON(cmd, body); err != nil {
						return err
					}
					url := base(c, portal, cmd.String("project")) + "/tasks/" + cmd.Args().First() + "/move"
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

func taskCommentsCmd() *cli.Command {
	return &cli.Command{
		Name:  "task-comments",
		Usage: "Task comment operations",
		Commands: []*cli.Command{
			{
				Name:  "list",
				Usage: "List task comments",
				Flags: []cli.Flag{
					portalFlag, projectFlag,
					&cli.StringFlag{Name: "task", Required: true, Usage: "Task ID"},
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
					url := base(c, portal, cmd.String("project")) + "/tasks/" + cmd.String("task") + "/comments"
					raw, err := c.Request(ctx, "GET", url, nil)
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:  "add",
				Usage: "Add a task comment",
				Flags: []cli.Flag{
					portalFlag, projectFlag,
					&cli.StringFlag{Name: "task", Required: true, Usage: "Task ID"},
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
					url := base(c, portal, cmd.String("project")) + "/tasks/" + cmd.String("task") + "/comments"
					raw, err := c.Request(ctx, "POST", url, &zohttp.RequestOpts{JSON: map[string]string{"comment": cmd.String("comment")}})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "update",
				Usage:     "Update a task comment",
				ArgsUsage: "<comment-id>",
				Flags: []cli.Flag{
					portalFlag, projectFlag,
					&cli.StringFlag{Name: "task", Required: true, Usage: "Task ID"},
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
					url := base(c, portal, cmd.String("project")) + "/tasks/" + cmd.String("task") + "/comments/" + cmd.Args().First()
					raw, err := c.Request(ctx, "PATCH", url, &zohttp.RequestOpts{JSON: map[string]string{"comment": cmd.String("comment")}})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "delete",
				Usage:     "Delete a task comment",
				ArgsUsage: "<comment-id>",
				Flags: []cli.Flag{
					portalFlag, projectFlag,
					&cli.StringFlag{Name: "task", Required: true, Usage: "Task ID"},
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
					url := base(c, portal, cmd.String("project")) + "/tasks/" + cmd.String("task") + "/comments/" + cmd.Args().First()
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

func taskFollowersCmd() *cli.Command {
	return &cli.Command{
		Name:  "task-followers",
		Usage: "Task follower operations",
		Commands: []*cli.Command{
			{
				Name:      "list",
				Usage:     "List task followers",
				ArgsUsage: "<task-id>",
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
					url := base(c, portal, cmd.String("project")) + "/tasks/" + cmd.Args().First() + "/followers"
					raw, err := c.Request(ctx, "GET", url, nil)
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "follow",
				Usage:     "Follow a task",
				ArgsUsage: "<task-id>",
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
					url := base(c, portal, cmd.String("project")) + "/tasks/" + cmd.Args().First() + "/follow"
					raw, err := c.Request(ctx, "POST", url, nil)
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "add",
				Usage:     "Add followers to a task",
				ArgsUsage: "<task-id>",
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
					url := base(c, portal, cmd.String("project")) + "/tasks/" + cmd.Args().First() + "/followers"
					raw, err := c.Request(ctx, "POST", url, &zohttp.RequestOpts{JSON: body})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "unfollow",
				Usage:     "Unfollow a task",
				ArgsUsage: "<task-id>",
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
					url := base(c, portal, cmd.String("project")) + "/tasks/" + cmd.Args().First() + "/unfollow"
					raw, err := c.Request(ctx, "POST", url, nil)
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
		},
	}
}

func taskCustomViewsCmd() *cli.Command {
	return &cli.Command{
		Name:  "task-customviews",
		Usage: "Task custom view operations",
		Commands: []*cli.Command{
			{
				Name:  "list",
				Usage: "List task custom views (portal-level)",
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
					url := c.ProjectsBase + "/portal/" + portal + "/tasks/customviews"
					raw, err := c.Request(ctx, "GET", url, nil)
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:  "project-list",
				Usage: "List task custom views (project-level)",
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
					url := base(c, portal, cmd.String("project")) + "/tasks/customviews"
					raw, err := c.Request(ctx, "GET", url, nil)
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "get",
				Usage:     "Get a task custom view",
				ArgsUsage: "<view-id>",
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
					url := c.ProjectsBase + "/portal/" + portal + "/tasks/customviews/" + cmd.Args().First()
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

func taskStatusTimelineCmd() *cli.Command {
	return &cli.Command{
		Name:  "task-statustimeline",
		Usage: "Task status timeline operations",
		Commands: []*cli.Command{
			{
				Name:      "get",
				Usage:     "Get status timeline for a task",
				ArgsUsage: "<task-id>",
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
					url := base(c, portal, cmd.String("project")) + "/tasks/" + cmd.Args().First() + "/status-timeline"
					raw, err := c.Request(ctx, "GET", url, nil)
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:  "project",
				Usage: "Get status timeline for project tasks",
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
					url := base(c, portal, cmd.String("project")) + "/tasks/status-timeline"
					raw, err := c.Request(ctx, "GET", url, nil)
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:  "portal",
				Usage: "Get status timeline for portal tasks",
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
					url := c.ProjectsBase + "/portal/" + portal + "/taskstatushistory"
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
