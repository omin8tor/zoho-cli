package projects

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/omin8tor/zoho-cli/internal/auth"
	zohttp "github.com/omin8tor/zoho-cli/internal/http"
	"github.com/omin8tor/zoho-cli/internal/output"
	"github.com/omin8tor/zoho-cli/internal/pagination"
	"github.com/urfave/cli/v3"
)

func convertDate(s string) string {
	if t, err := time.Parse("2006-01-02", s); err == nil {
		return t.Format("01-02-2006")
	}
	return s
}

func getClient() (*zohttp.Client, error) {
	config, err := auth.ResolveAuth()
	if err != nil {
		return nil, err
	}
	return zohttp.NewClient(config)
}

func base(c *zohttp.Client, portal, project string) string {
	return c.ProjectsBase + "/portal/" + portal + "/projects/" + project
}

var portalFlag = &cli.StringFlag{Name: "portal", Usage: "Portal ID", Sources: cli.EnvVars("ZOHO_PORTAL_ID")}
var projectFlag = &cli.StringFlag{Name: "project", Required: true, Usage: "Project ID"}

func requirePortal(cmd *cli.Command) (string, error) {
	v := cmd.String("portal")
	if v == "" {
		return "", fmt.Errorf("--portal is required (or set ZOHO_PORTAL_ID env var)")
	}
	return v, nil
}

func Commands() *cli.Command {
	return &cli.Command{
		Name:  "projects",
		Usage: "Zoho Projects operations",
		Commands: []*cli.Command{
			projectsCoreCmd(),
			tasksCmd(),
			taskCommentsCmd(),
			taskFollowersCmd(),
			taskCustomViewsCmd(),
			taskStatusTimelineCmd(),
			issuesCmd(),
			issueCommentsCmd(),
			issueFollowersCmd(),
			issueLinkingCmd(),
			issueResolutionCmd(),
			issueAttachmentsCmd(),
			issueCustomViewsCmd(),
			tasklistsCmd(),
			tasklistCommentsCmd(),
			tasklistFollowersCmd(),
			timelogsCmd(),
			timelogBulkCmd(),
			timelogTimersCmd(),
			timelogPinsCmd(),
			usersCmd(),
			projectUsersCmd(),
			milestonesCmd(),
			phasesCmd(),
			phaseFollowersCmd(),
			phaseCommentsCmd(),
			dependenciesCmd(),
			forumsCmd(),
			forumCommentsCmd(),
			forumCategoriesCmd(),
			forumFollowersCmd(),
			eventsCmd(),
			eventCommentsCmd(),
			attachmentsCmd(),
			leavesCmd(),
			tagsCmd(),
			portalsCmd(),
			trashCmd(),
			searchCmd(),
			feedCmd(),
			projectCommentsCmd(),
			projectGroupsCmd(),
			teamsCmd(),
			profilesCmd(),
			rolesCmd(),
			customRecordsCmd(),
			reportsCmd(),
		},
	}
}

func projectsCoreCmd() *cli.Command {
	return &cli.Command{
		Name:  "core",
		Usage: "Project CRUD operations",
		Commands: []*cli.Command{
			{
				Name:  "list",
				Usage: "List all projects",
				Flags: []cli.Flag{portalFlag},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					portal, err := requirePortal(cmd)
					if err != nil {
						return err
					}
					url := c.ProjectsBase + "/portal/" + portal + "/projects"
					items, err := pagination.PaginateProjects(c, url, "", nil, 0)
					if err != nil {
						return err
					}
					return output.JSON(items)
				},
			},
			{
				Name:      "get",
				Usage:     "Get a single project",
				ArgsUsage: "<project-id>",
				Flags:     []cli.Flag{portalFlag},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					portal, err := requirePortal(cmd)
					if err != nil {
						return err
					}
					url := c.ProjectsBase + "/portal/" + portal + "/projects/" + cmd.Args().First()
					raw, err := c.Request("GET", url, nil)
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:  "search",
				Usage: "Search projects",
				Flags: []cli.Flag{
					portalFlag,
					&cli.StringFlag{Name: "query", Required: true, Usage: "Search query"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					portal, err := requirePortal(cmd)
					if err != nil {
						return err
					}
					url := c.ProjectsBase + "/portal/" + portal + "/search"
					raw, err := c.Request("GET", url, &zohttp.RequestOpts{
						Params: map[string]string{
							"search_term": cmd.String("query"),
							"module":      "all",
							"status":      "all",
						},
					})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:  "create",
				Usage: "Create a project",
				Flags: []cli.Flag{
					portalFlag,
					&cli.StringFlag{Name: "name", Required: true, Usage: "Project name"},
					&cli.StringFlag{Name: "json", Usage: "Additional fields as JSON"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					portal, err := requirePortal(cmd)
					if err != nil {
						return err
					}
					body := map[string]any{"name": cmd.String("name")}
					if j := cmd.String("json"); j != "" {
						var extra map[string]any
						json.Unmarshal([]byte(j), &extra)
						for k, v := range extra {
							body[k] = v
						}
					}
					url := c.ProjectsBase + "/portal/" + portal + "/projects"
					raw, err := c.Request("POST", url, &zohttp.RequestOpts{JSON: body})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "update",
				Usage:     "Update a project",
				ArgsUsage: "<project-id>",
				Flags: []cli.Flag{
					portalFlag,
					&cli.StringFlag{Name: "json", Required: true, Usage: "Fields to update as JSON"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					portal, err := requirePortal(cmd)
					if err != nil {
						return err
					}
					var parsed map[string]any
					json.Unmarshal([]byte(cmd.String("json")), &parsed)
					url := c.ProjectsBase + "/portal/" + portal + "/projects/" + cmd.Args().First()
					raw, err := c.Request("PATCH", url, &zohttp.RequestOpts{JSON: parsed})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "delete",
				Usage:     "Delete a project",
				ArgsUsage: "<project-id>",
				Flags:     []cli.Flag{portalFlag},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					portal, err := requirePortal(cmd)
					if err != nil {
						return err
					}
					url := c.ProjectsBase + "/portal/" + portal + "/projects/" + cmd.Args().First()
					raw, err := c.Request("DELETE", url, nil)
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "trash",
				Usage:     "Move a project to trash",
				ArgsUsage: "<project-id>",
				Flags:     []cli.Flag{portalFlag},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					portal, err := requirePortal(cmd)
					if err != nil {
						return err
					}
					url := c.ProjectsBase + "/portal/" + portal + "/projects/" + cmd.Args().First() + "/trash"
					raw, err := c.Request("POST", url, nil)
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "restore",
				Usage:     "Restore a project from trash",
				ArgsUsage: "<project-id>",
				Flags:     []cli.Flag{portalFlag},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					portal, err := requirePortal(cmd)
					if err != nil {
						return err
					}
					url := c.ProjectsBase + "/portal/" + portal + "/projects/" + cmd.Args().First() + "/restore"
					raw, err := c.Request("POST", url, nil)
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
		},
	}
}

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
					&cli.StringFlag{Name: "status", Usage: "Filter: open, closed, in progress"},
					&cli.StringFlag{Name: "priority", Usage: "Filter: none, low, medium, high"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
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
					items, err := pagination.PaginateProjects(c, url, "tasks", params, 0)
					if err != nil {
						return err
					}
					return output.JSON(items)
				},
			},
			{
				Name:  "my",
				Usage: "List my tasks across all projects",
				Flags: []cli.Flag{
					portalFlag,
					&cli.StringFlag{Name: "status", Usage: "Filter: open, closed, in progress"},
					&cli.StringFlag{Name: "priority", Usage: "Filter: none, low, medium, high"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
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
					items, err := pagination.PaginateProjects(c, url, "tasks", params, 0)
					if err != nil {
						return err
					}
					return output.JSON(items)
				},
			},
			{
				Name:      "get",
				Usage:     "Get a single task",
				ArgsUsage: "<task-id>",
				Flags:     []cli.Flag{portalFlag, projectFlag},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					portal, err := requirePortal(cmd)
					if err != nil {
						return err
					}
					url := base(c, portal, cmd.String("project")) + "/tasks/" + cmd.Args().First()
					raw, err := c.Request("GET", url, nil)
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
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					portal, err := requirePortal(cmd)
					if err != nil {
						return err
					}
					body := map[string]any{"name": cmd.String("name")}
					if j := cmd.String("json"); j != "" {
						var extra map[string]any
						json.Unmarshal([]byte(j), &extra)
						for k, v := range extra {
							body[k] = v
						}
					}
					url := base(c, portal, cmd.String("project")) + "/tasks"
					raw, err := c.Request("POST", url, &zohttp.RequestOpts{JSON: body})
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
					&cli.StringFlag{Name: "json", Required: true, Usage: "Fields to update as JSON"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					portal, err := requirePortal(cmd)
					if err != nil {
						return err
					}
					var parsed map[string]any
					json.Unmarshal([]byte(cmd.String("json")), &parsed)
					url := base(c, portal, cmd.String("project")) + "/tasks/" + cmd.Args().First()
					raw, err := c.Request("PATCH", url, &zohttp.RequestOpts{JSON: parsed})
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
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					portal, err := requirePortal(cmd)
					if err != nil {
						return err
					}
					url := base(c, portal, cmd.String("project")) + "/tasks/" + cmd.Args().First()
					raw, err := c.Request("DELETE", url, nil)
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
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					portal, err := requirePortal(cmd)
					if err != nil {
						return err
					}
					url := base(c, portal, cmd.String("project")) + "/tasks/" + cmd.Args().First() + "/subtasks"
					raw, err := c.Request("GET", url, nil)
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
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
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
					if j := cmd.String("json"); j != "" {
						var extra map[string]any
						json.Unmarshal([]byte(j), &extra)
						for k, v := range extra {
							body[k] = v
						}
					}
					url := base(c, portal, cmd.String("project")) + "/tasks"
					raw, err := c.Request("POST", url, &zohttp.RequestOpts{JSON: body})
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
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					portal, err := requirePortal(cmd)
					if err != nil {
						return err
					}
					url := base(c, portal, cmd.String("project")) + "/tasks/" + cmd.Args().First() + "/copy"
					raw, err := c.Request("POST", url, &zohttp.RequestOpts{
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
					&cli.StringFlag{Name: "json", Required: true, Usage: "Move details as JSON"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					portal, err := requirePortal(cmd)
					if err != nil {
						return err
					}
					var body any
					json.Unmarshal([]byte(cmd.String("json")), &body)
					url := base(c, portal, cmd.String("project")) + "/tasks/" + cmd.Args().First() + "/move"
					raw, err := c.Request("POST", url, &zohttp.RequestOpts{JSON: body})
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
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					portal, err := requirePortal(cmd)
					if err != nil {
						return err
					}
					url := base(c, portal, cmd.String("project")) + "/tasks/" + cmd.String("task") + "/comments"
					raw, err := c.Request("GET", url, nil)
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
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					portal, err := requirePortal(cmd)
					if err != nil {
						return err
					}
					url := base(c, portal, cmd.String("project")) + "/tasks/" + cmd.String("task") + "/comments"
					raw, err := c.Request("POST", url, &zohttp.RequestOpts{JSON: map[string]string{"comment": cmd.String("comment")}})
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
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					portal, err := requirePortal(cmd)
					if err != nil {
						return err
					}
					url := base(c, portal, cmd.String("project")) + "/tasks/" + cmd.String("task") + "/comments/" + cmd.Args().First()
					raw, err := c.Request("PATCH", url, &zohttp.RequestOpts{JSON: map[string]string{"comment": cmd.String("comment")}})
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
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					portal, err := requirePortal(cmd)
					if err != nil {
						return err
					}
					url := base(c, portal, cmd.String("project")) + "/tasks/" + cmd.String("task") + "/comments/" + cmd.Args().First()
					raw, err := c.Request("DELETE", url, nil)
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
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					portal, err := requirePortal(cmd)
					if err != nil {
						return err
					}
					url := base(c, portal, cmd.String("project")) + "/tasks/" + cmd.Args().First() + "/followers"
					raw, err := c.Request("GET", url, nil)
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
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					portal, err := requirePortal(cmd)
					if err != nil {
						return err
					}
					url := base(c, portal, cmd.String("project")) + "/tasks/" + cmd.Args().First() + "/followers"
					raw, err := c.Request("POST", url, nil)
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
					&cli.StringFlag{Name: "json", Required: true, Usage: "Follower details as JSON"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					portal, err := requirePortal(cmd)
					if err != nil {
						return err
					}
					var body any
					json.Unmarshal([]byte(cmd.String("json")), &body)
					url := base(c, portal, cmd.String("project")) + "/tasks/" + cmd.Args().First() + "/followers"
					raw, err := c.Request("POST", url, &zohttp.RequestOpts{JSON: body})
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
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					portal, err := requirePortal(cmd)
					if err != nil {
						return err
					}
					url := base(c, portal, cmd.String("project")) + "/tasks/" + cmd.Args().First() + "/followers"
					raw, err := c.Request("DELETE", url, nil)
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
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					portal, err := requirePortal(cmd)
					if err != nil {
						return err
					}
					url := c.ProjectsBase + "/portal/" + portal + "/tasks/customviews"
					raw, err := c.Request("GET", url, nil)
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
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					portal, err := requirePortal(cmd)
					if err != nil {
						return err
					}
					url := base(c, portal, cmd.String("project")) + "/tasks/customviews"
					raw, err := c.Request("GET", url, nil)
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
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					portal, err := requirePortal(cmd)
					if err != nil {
						return err
					}
					url := c.ProjectsBase + "/portal/" + portal + "/tasks/customviews/" + cmd.Args().First()
					raw, err := c.Request("GET", url, nil)
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
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					portal, err := requirePortal(cmd)
					if err != nil {
						return err
					}
					url := base(c, portal, cmd.String("project")) + "/tasks/" + cmd.Args().First() + "/statustimeline"
					raw, err := c.Request("GET", url, nil)
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
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					portal, err := requirePortal(cmd)
					if err != nil {
						return err
					}
					url := base(c, portal, cmd.String("project")) + "/tasks/statustimeline"
					raw, err := c.Request("GET", url, nil)
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
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					portal, err := requirePortal(cmd)
					if err != nil {
						return err
					}
					url := c.ProjectsBase + "/portal/" + portal + "/tasks/statustimeline"
					raw, err := c.Request("GET", url, nil)
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
		},
	}
}

func issuesCmd() *cli.Command {
	return &cli.Command{
		Name:  "issues",
		Usage: "Project issue operations",
		Commands: []*cli.Command{
			{
				Name:  "list",
				Usage: "List issues in a project",
				Flags: []cli.Flag{portalFlag, projectFlag},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					portal, err := requirePortal(cmd)
					if err != nil {
						return err
					}
					url := base(c, portal, cmd.String("project")) + "/issues"
					items, err := pagination.PaginateProjects(c, url, "issues", nil, 0)
					if err != nil {
						return err
					}
					return output.JSON(items)
				},
			},
			{
				Name:      "get",
				Usage:     "Get a single issue",
				ArgsUsage: "<issue-id>",
				Flags:     []cli.Flag{portalFlag, projectFlag},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					portal, err := requirePortal(cmd)
					if err != nil {
						return err
					}
					url := base(c, portal, cmd.String("project")) + "/issues/" + cmd.Args().First()
					raw, err := c.Request("GET", url, nil)
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:  "create",
				Usage: "Create an issue",
				Flags: []cli.Flag{
					portalFlag, projectFlag,
					&cli.StringFlag{Name: "name", Required: true, Usage: "Issue title"},
					&cli.StringFlag{Name: "json", Usage: "Additional fields as JSON"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					portal, err := requirePortal(cmd)
					if err != nil {
						return err
					}
					body := map[string]any{"name": cmd.String("name")}
					if j := cmd.String("json"); j != "" {
						var extra map[string]any
						json.Unmarshal([]byte(j), &extra)
						for k, v := range extra {
							body[k] = v
						}
					}
					url := base(c, portal, cmd.String("project")) + "/issues"
					raw, err := c.Request("POST", url, &zohttp.RequestOpts{JSON: body})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "update",
				Usage:     "Update an issue",
				ArgsUsage: "<issue-id>",
				Flags: []cli.Flag{
					portalFlag, projectFlag,
					&cli.StringFlag{Name: "json", Required: true, Usage: "Fields to update as JSON"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					portal, err := requirePortal(cmd)
					if err != nil {
						return err
					}
					var parsed map[string]any
					json.Unmarshal([]byte(cmd.String("json")), &parsed)
					url := base(c, portal, cmd.String("project")) + "/issues/" + cmd.Args().First()
					raw, err := c.Request("PATCH", url, &zohttp.RequestOpts{JSON: parsed})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "delete",
				Usage:     "Delete an issue",
				ArgsUsage: "<issue-id>",
				Flags:     []cli.Flag{portalFlag, projectFlag},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					portal, err := requirePortal(cmd)
					if err != nil {
						return err
					}
					url := base(c, portal, cmd.String("project")) + "/issues/" + cmd.Args().First()
					raw, err := c.Request("DELETE", url, nil)
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},

			{
				Name:      "description",
				Usage:     "Get issue description",
				ArgsUsage: "<issue-id>",
				Flags:     []cli.Flag{portalFlag, projectFlag},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					portal, err := requirePortal(cmd)
					if err != nil {
						return err
					}
					url := base(c, portal, cmd.String("project")) + "/issues/" + cmd.Args().First() + "/description"
					raw, err := c.Request("GET", url, nil)
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "move",
				Usage:     "Move an issue",
				ArgsUsage: "<issue-id>",
				Flags: []cli.Flag{
					portalFlag, projectFlag,
					&cli.StringFlag{Name: "json", Required: true, Usage: "Move details as JSON"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					portal, err := requirePortal(cmd)
					if err != nil {
						return err
					}
					var body any
					json.Unmarshal([]byte(cmd.String("json")), &body)
					url := base(c, portal, cmd.String("project")) + "/issues/" + cmd.Args().First() + "/move"
					raw, err := c.Request("POST", url, &zohttp.RequestOpts{JSON: body})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "clone",
				Usage:     "Clone an issue",
				ArgsUsage: "<issue-id>",
				Flags:     []cli.Flag{portalFlag, projectFlag},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					portal, err := requirePortal(cmd)
					if err != nil {
						return err
					}
					url := base(c, portal, cmd.String("project")) + "/issues/" + cmd.Args().First() + "/clone"
					raw, err := c.Request("POST", url, nil)
					if err != nil {
						return err
					}
					var envelope map[string]json.RawMessage
					if json.Unmarshal(raw, &envelope) == nil {
						if bugs, ok := envelope["bugs"]; ok {
							var arr []json.RawMessage
							if json.Unmarshal(bugs, &arr) == nil && len(arr) > 0 {
								return output.JSONRaw(arr[0])
							}
						}
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "activities",
				Usage:     "Get issue activities",
				ArgsUsage: "<issue-id>",
				Flags:     []cli.Flag{portalFlag, projectFlag},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					portal, err := requirePortal(cmd)
					if err != nil {
						return err
					}
					url := base(c, portal, cmd.String("project")) + "/issues/" + cmd.Args().First() + "/activities"
					raw, err := c.Request("GET", url, nil)
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
		},
	}
}

func issueCommentsCmd() *cli.Command {
	return &cli.Command{
		Name:  "issue-comments",
		Usage: "Issue comment operations",
		Commands: []*cli.Command{
			{
				Name:  "list",
				Usage: "List issue comments",
				Flags: []cli.Flag{
					portalFlag, projectFlag,
					&cli.StringFlag{Name: "issue", Required: true, Usage: "Issue ID"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					portal, err := requirePortal(cmd)
					if err != nil {
						return err
					}
					url := base(c, portal, cmd.String("project")) + "/issues/" + cmd.String("issue") + "/comments"
					raw, err := c.Request("GET", url, nil)
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "get",
				Usage:     "Get an issue comment",
				ArgsUsage: "<comment-id>",
				Flags: []cli.Flag{
					portalFlag, projectFlag,
					&cli.StringFlag{Name: "issue", Required: true, Usage: "Issue ID"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					portal, err := requirePortal(cmd)
					if err != nil {
						return err
					}
					url := base(c, portal, cmd.String("project")) + "/issues/" + cmd.String("issue") + "/comments/" + cmd.Args().First()
					raw, err := c.Request("GET", url, nil)
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:  "add",
				Usage: "Add an issue comment",
				Flags: []cli.Flag{
					portalFlag, projectFlag,
					&cli.StringFlag{Name: "issue", Required: true, Usage: "Issue ID"},
					&cli.StringFlag{Name: "comment", Required: true, Usage: "Comment text"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					portal, err := requirePortal(cmd)
					if err != nil {
						return err
					}
					url := base(c, portal, cmd.String("project")) + "/issues/" + cmd.String("issue") + "/comments"
					raw, err := c.Request("POST", url, &zohttp.RequestOpts{JSON: map[string]string{"comment": cmd.String("comment")}})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "update",
				Usage:     "Update an issue comment",
				ArgsUsage: "<comment-id>",
				Flags: []cli.Flag{
					portalFlag, projectFlag,
					&cli.StringFlag{Name: "issue", Required: true, Usage: "Issue ID"},
					&cli.StringFlag{Name: "comment", Required: true, Usage: "Updated comment text"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					portal, err := requirePortal(cmd)
					if err != nil {
						return err
					}
					url := base(c, portal, cmd.String("project")) + "/issues/" + cmd.String("issue") + "/comments/" + cmd.Args().First()
					raw, err := c.Request("PATCH", url, &zohttp.RequestOpts{JSON: map[string]string{"comment": cmd.String("comment")}})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "delete",
				Usage:     "Delete an issue comment",
				ArgsUsage: "<comment-id>",
				Flags: []cli.Flag{
					portalFlag, projectFlag,
					&cli.StringFlag{Name: "issue", Required: true, Usage: "Issue ID"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					portal, err := requirePortal(cmd)
					if err != nil {
						return err
					}
					url := base(c, portal, cmd.String("project")) + "/issues/" + cmd.String("issue") + "/comments/" + cmd.Args().First()
					raw, err := c.Request("DELETE", url, nil)
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
		},
	}
}

func issueFollowersCmd() *cli.Command {
	return &cli.Command{
		Name:  "issue-followers",
		Usage: "Issue follower operations",
		Commands: []*cli.Command{
			{
				Name:      "list",
				Usage:     "List issue followers",
				ArgsUsage: "<issue-id>",
				Flags:     []cli.Flag{portalFlag, projectFlag},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					portal, err := requirePortal(cmd)
					if err != nil {
						return err
					}
					url := base(c, portal, cmd.String("project")) + "/issues/" + cmd.Args().First() + "/followers"
					raw, err := c.Request("GET", url, nil)
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "follow",
				Usage:     "Follow an issue",
				ArgsUsage: "<issue-id>",
				Flags:     []cli.Flag{portalFlag, projectFlag},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					portal, err := requirePortal(cmd)
					if err != nil {
						return err
					}
					url := base(c, portal, cmd.String("project")) + "/issues/" + cmd.Args().First() + "/followers"
					raw, err := c.Request("POST", url, nil)
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "remove",
				Usage:     "Remove issue followers",
				ArgsUsage: "<issue-id>",
				Flags: []cli.Flag{
					portalFlag, projectFlag,
					&cli.StringFlag{Name: "json", Required: true, Usage: "Follower IDs as JSON"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					portal, err := requirePortal(cmd)
					if err != nil {
						return err
					}
					var body any
					json.Unmarshal([]byte(cmd.String("json")), &body)
					url := base(c, portal, cmd.String("project")) + "/issues/" + cmd.Args().First() + "/followers"
					raw, err := c.Request("DELETE", url, &zohttp.RequestOpts{JSON: body})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
		},
	}
}

func issueLinkingCmd() *cli.Command {
	return &cli.Command{
		Name:  "issue-linking",
		Usage: "Issue linking operations",
		Commands: []*cli.Command{
			{
				Name:      "list",
				Usage:     "List linked issues",
				ArgsUsage: "<issue-id>",
				Flags:     []cli.Flag{portalFlag, projectFlag},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					portal, err := requirePortal(cmd)
					if err != nil {
						return err
					}
					url := base(c, portal, cmd.String("project")) + "/issues/" + cmd.Args().First() + "/linkedissues"
					raw, err := c.Request("GET", url, nil)
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "link",
				Usage:     "Link an issue",
				ArgsUsage: "<issue-id>",
				Flags: []cli.Flag{
					portalFlag, projectFlag,
					&cli.StringFlag{Name: "json", Required: true, Usage: "Link details as JSON"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					portal, err := requirePortal(cmd)
					if err != nil {
						return err
					}
					var body any
					json.Unmarshal([]byte(cmd.String("json")), &body)
					url := base(c, portal, cmd.String("project")) + "/issues/" + cmd.Args().First() + "/linkedissues"
					raw, err := c.Request("POST", url, &zohttp.RequestOpts{JSON: body})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:  "bulk-link",
				Usage: "Bulk link issues",
				Flags: []cli.Flag{
					portalFlag, projectFlag,
					&cli.StringFlag{Name: "json", Required: true, Usage: "Bulk link details as JSON"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					portal, err := requirePortal(cmd)
					if err != nil {
						return err
					}
					var body any
					json.Unmarshal([]byte(cmd.String("json")), &body)
					url := base(c, portal, cmd.String("project")) + "/issues/bulk-link-bugs"
					raw, err := c.Request("POST", url, &zohttp.RequestOpts{JSON: body})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "change-type",
				Usage:     "Change link type",
				ArgsUsage: "<issue-id> <link-id>",
				Flags: []cli.Flag{
					portalFlag, projectFlag,
					&cli.StringFlag{Name: "json", Required: true, Usage: "Link type as JSON"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					portal, err := requirePortal(cmd)
					if err != nil {
						return err
					}
					var body any
					json.Unmarshal([]byte(cmd.String("json")), &body)
					url := base(c, portal, cmd.String("project")) + "/issues/" + cmd.Args().Get(0) + "/linkedissues/" + cmd.Args().Get(1)
					raw, err := c.Request("PATCH", url, &zohttp.RequestOpts{JSON: body})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "unlink",
				Usage:     "Unlink an issue",
				ArgsUsage: "<issue-id> <link-id>",
				Flags:     []cli.Flag{portalFlag, projectFlag},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					portal, err := requirePortal(cmd)
					if err != nil {
						return err
					}
					url := base(c, portal, cmd.String("project")) + "/issues/" + cmd.Args().Get(0) + "/linkedissues/" + cmd.Args().Get(1)
					raw, err := c.Request("DELETE", url, nil)
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
		},
	}
}

func issueResolutionCmd() *cli.Command {
	return &cli.Command{
		Name:  "issue-resolution",
		Usage: "Issue resolution operations",
		Commands: []*cli.Command{
			{
				Name:      "get",
				Usage:     "Get issue resolution",
				ArgsUsage: "<issue-id>",
				Flags:     []cli.Flag{portalFlag, projectFlag},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					portal, err := requirePortal(cmd)
					if err != nil {
						return err
					}
					url := base(c, portal, cmd.String("project")) + "/issues/" + cmd.Args().First() + "/resolution"
					raw, err := c.Request("GET", url, nil)
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "add",
				Usage:     "Add issue resolution",
				ArgsUsage: "<issue-id>",
				Flags: []cli.Flag{
					portalFlag, projectFlag,
					&cli.StringFlag{Name: "json", Required: true, Usage: "Resolution details as JSON"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					portal, err := requirePortal(cmd)
					if err != nil {
						return err
					}
					var body any
					json.Unmarshal([]byte(cmd.String("json")), &body)
					url := base(c, portal, cmd.String("project")) + "/issues/" + cmd.Args().First() + "/resolution"
					raw, err := c.Request("POST", url, &zohttp.RequestOpts{JSON: body})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "update",
				Usage:     "Update issue resolution",
				ArgsUsage: "<issue-id>",
				Flags: []cli.Flag{
					portalFlag, projectFlag,
					&cli.StringFlag{Name: "json", Required: true, Usage: "Resolution details as JSON"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					portal, err := requirePortal(cmd)
					if err != nil {
						return err
					}
					var body any
					json.Unmarshal([]byte(cmd.String("json")), &body)
					url := base(c, portal, cmd.String("project")) + "/issues/" + cmd.Args().First() + "/resolution"
					raw, err := c.Request("PATCH", url, &zohttp.RequestOpts{JSON: body})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "delete",
				Usage:     "Delete issue resolution",
				ArgsUsage: "<issue-id>",
				Flags:     []cli.Flag{portalFlag, projectFlag},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					portal, err := requirePortal(cmd)
					if err != nil {
						return err
					}
					url := base(c, portal, cmd.String("project")) + "/issues/" + cmd.Args().First() + "/resolution"
					raw, err := c.Request("DELETE", url, nil)
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
		},
	}
}

func issueAttachmentsCmd() *cli.Command {
	return &cli.Command{
		Name:  "issue-attachments",
		Usage: "Issue attachment operations",
		Commands: []*cli.Command{
			{
				Name:      "list",
				Usage:     "List issue attachments",
				ArgsUsage: "<issue-id>",
				Flags:     []cli.Flag{portalFlag, projectFlag},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					portal, err := requirePortal(cmd)
					if err != nil {
						return err
					}
					url := base(c, portal, cmd.String("project")) + "/issues/" + cmd.Args().First() + "/attachments"
					raw, err := c.Request("GET", url, nil)
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "associate",
				Usage:     "Associate attachments to an issue",
				ArgsUsage: "<issue-id>",
				Flags: []cli.Flag{
					portalFlag, projectFlag,
					&cli.StringFlag{Name: "json", Required: true, Usage: "Attachment details as JSON"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					portal, err := requirePortal(cmd)
					if err != nil {
						return err
					}
					var body any
					json.Unmarshal([]byte(cmd.String("json")), &body)
					url := base(c, portal, cmd.String("project")) + "/issues/" + cmd.Args().First() + "/attachments"
					raw, err := c.Request("POST", url, &zohttp.RequestOpts{JSON: body})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "dissociate",
				Usage:     "Dissociate an attachment from an issue",
				ArgsUsage: "<issue-id> <attachment-id>",
				Flags:     []cli.Flag{portalFlag, projectFlag},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					portal, err := requirePortal(cmd)
					if err != nil {
						return err
					}
					url := base(c, portal, cmd.String("project")) + "/issues/" + cmd.Args().Get(0) + "/attachments/" + cmd.Args().Get(1)
					raw, err := c.Request("DELETE", url, nil)
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
		},
	}
}

func issueCustomViewsCmd() *cli.Command {
	return &cli.Command{
		Name:  "issue-customviews",
		Usage: "Issue custom view operations",
		Commands: []*cli.Command{
			{
				Name:  "list",
				Usage: "List issue custom views (portal-level)",
				Flags: []cli.Flag{portalFlag},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					portal, err := requirePortal(cmd)
					if err != nil {
						return err
					}
					url := c.ProjectsBase + "/portal/" + portal + "/issues/customviews"
					raw, err := c.Request("GET", url, nil)
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:  "project-list",
				Usage: "List issue custom views (project-level)",
				Flags: []cli.Flag{portalFlag, projectFlag},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					portal, err := requirePortal(cmd)
					if err != nil {
						return err
					}
					url := base(c, portal, cmd.String("project")) + "/issues/customviews"
					raw, err := c.Request("GET", url, nil)
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "get",
				Usage:     "Get an issue custom view",
				ArgsUsage: "<view-id>",
				Flags:     []cli.Flag{portalFlag},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					portal, err := requirePortal(cmd)
					if err != nil {
						return err
					}
					url := c.ProjectsBase + "/portal/" + portal + "/issues/customviews/" + cmd.Args().First()
					raw, err := c.Request("GET", url, nil)
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
		},
	}
}

func tasklistsCmd() *cli.Command {
	return &cli.Command{
		Name:  "tasklists",
		Usage: "Project tasklist operations",
		Commands: []*cli.Command{
			{
				Name:  "list",
				Usage: "List tasklists",
				Flags: []cli.Flag{portalFlag, projectFlag},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					portal, err := requirePortal(cmd)
					if err != nil {
						return err
					}
					url := base(c, portal, cmd.String("project")) + "/tasklists"
					items, err := pagination.PaginateProjects(c, url, "tasklists", nil, 0)
					if err != nil {
						return err
					}
					return output.JSON(items)
				},
			},
			{
				Name:      "get",
				Usage:     "Get a tasklist",
				ArgsUsage: "<tasklist-id>",
				Flags:     []cli.Flag{portalFlag, projectFlag},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					portal, err := requirePortal(cmd)
					if err != nil {
						return err
					}
					url := base(c, portal, cmd.String("project")) + "/tasklists/" + cmd.Args().First()
					raw, err := c.Request("GET", url, nil)
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:  "create",
				Usage: "Create a tasklist",
				Flags: []cli.Flag{
					portalFlag, projectFlag,
					&cli.StringFlag{Name: "name", Required: true, Usage: "Tasklist name"},
					&cli.StringFlag{Name: "json", Usage: "Additional fields as JSON"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					portal, err := requirePortal(cmd)
					if err != nil {
						return err
					}
					body := map[string]any{"name": cmd.String("name")}
					if j := cmd.String("json"); j != "" {
						var extra map[string]any
						json.Unmarshal([]byte(j), &extra)
						for k, v := range extra {
							body[k] = v
						}
					}
					url := base(c, portal, cmd.String("project")) + "/tasklists"
					raw, err := c.Request("POST", url, &zohttp.RequestOpts{JSON: body})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "update",
				Usage:     "Update a tasklist",
				ArgsUsage: "<tasklist-id>",
				Flags: []cli.Flag{
					portalFlag, projectFlag,
					&cli.StringFlag{Name: "json", Required: true, Usage: "Fields to update as JSON"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					portal, err := requirePortal(cmd)
					if err != nil {
						return err
					}
					var body any
					json.Unmarshal([]byte(cmd.String("json")), &body)
					url := base(c, portal, cmd.String("project")) + "/tasklists/" + cmd.Args().First()
					raw, err := c.Request("PATCH", url, &zohttp.RequestOpts{JSON: body})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "delete",
				Usage:     "Delete a tasklist",
				ArgsUsage: "<tasklist-id>",
				Flags:     []cli.Flag{portalFlag, projectFlag},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					portal, err := requirePortal(cmd)
					if err != nil {
						return err
					}
					url := base(c, portal, cmd.String("project")) + "/tasklists/" + cmd.Args().First()
					raw, err := c.Request("DELETE", url, nil)
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
		},
	}
}

func tasklistCommentsCmd() *cli.Command {
	return &cli.Command{
		Name:  "tasklist-comments",
		Usage: "Tasklist comment operations",
		Commands: []*cli.Command{
			{
				Name:  "list",
				Usage: "List tasklist comments",
				Flags: []cli.Flag{
					portalFlag, projectFlag,
					&cli.StringFlag{Name: "tasklist", Required: true, Usage: "Tasklist ID"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					portal, err := requirePortal(cmd)
					if err != nil {
						return err
					}
					url := base(c, portal, cmd.String("project")) + "/tasklists" + "/" + cmd.String("tasklist") + "/comments"
					raw, err := c.Request("GET", url, nil)
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "get",
				Usage:     "Get a tasklist comment",
				ArgsUsage: "<comment-id>",
				Flags: []cli.Flag{
					portalFlag, projectFlag,
					&cli.StringFlag{Name: "tasklist", Required: true, Usage: "Tasklist ID"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					portal, err := requirePortal(cmd)
					if err != nil {
						return err
					}
					url := base(c, portal, cmd.String("project")) + "/tasklists" + "/" + cmd.String("tasklist") + "/comments/" + cmd.Args().First()
					raw, err := c.Request("GET", url, nil)
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:  "add",
				Usage: "Add a tasklist comment",
				Flags: []cli.Flag{
					portalFlag, projectFlag,
					&cli.StringFlag{Name: "tasklist", Required: true, Usage: "Tasklist ID"},
					&cli.StringFlag{Name: "comment", Required: true, Usage: "Comment text"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					portal, err := requirePortal(cmd)
					if err != nil {
						return err
					}
					url := base(c, portal, cmd.String("project")) + "/tasklists" + "/" + cmd.String("tasklist") + "/comments"
					raw, err := c.Request("POST", url, &zohttp.RequestOpts{JSON: map[string]string{"comment": cmd.String("comment")}})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "update",
				Usage:     "Update a tasklist comment",
				ArgsUsage: "<comment-id>",
				Flags: []cli.Flag{
					portalFlag, projectFlag,
					&cli.StringFlag{Name: "tasklist", Required: true, Usage: "Tasklist ID"},
					&cli.StringFlag{Name: "comment", Required: true, Usage: "Updated comment text"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					portal, err := requirePortal(cmd)
					if err != nil {
						return err
					}
					url := base(c, portal, cmd.String("project")) + "/tasklists" + "/" + cmd.String("tasklist") + "/comments/" + cmd.Args().First()
					raw, err := c.Request("PATCH", url, &zohttp.RequestOpts{JSON: map[string]string{"comment": cmd.String("comment")}})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "delete",
				Usage:     "Delete a tasklist comment",
				ArgsUsage: "<comment-id>",
				Flags: []cli.Flag{
					portalFlag, projectFlag,
					&cli.StringFlag{Name: "tasklist", Required: true, Usage: "Tasklist ID"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					portal, err := requirePortal(cmd)
					if err != nil {
						return err
					}
					url := base(c, portal, cmd.String("project")) + "/tasklists" + "/" + cmd.String("tasklist") + "/comments/" + cmd.Args().First()
					raw, err := c.Request("DELETE", url, nil)
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
		},
	}
}

func tasklistFollowersCmd() *cli.Command {
	return &cli.Command{
		Name:  "tasklist-followers",
		Usage: "Tasklist follower operations",
		Commands: []*cli.Command{
			{
				Name:      "list",
				Usage:     "List tasklist followers",
				ArgsUsage: "<tasklist-id>",
				Flags:     []cli.Flag{portalFlag, projectFlag},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					portal, err := requirePortal(cmd)
					if err != nil {
						return err
					}
					url := base(c, portal, cmd.String("project")) + "/tasklists/" + cmd.Args().First() + "/followers"
					raw, err := c.Request("GET", url, nil)
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "follow",
				Usage:     "Follow a tasklist",
				ArgsUsage: "<tasklist-id>",
				Flags:     []cli.Flag{portalFlag, projectFlag},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					portal, err := requirePortal(cmd)
					if err != nil {
						return err
					}
					url := base(c, portal, cmd.String("project")) + "/tasklists/" + cmd.Args().First() + "/followers"
					raw, err := c.Request("POST", url, nil)
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "unfollow",
				Usage:     "Unfollow a tasklist",
				ArgsUsage: "<tasklist-id>",
				Flags:     []cli.Flag{portalFlag, projectFlag},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					portal, err := requirePortal(cmd)
					if err != nil {
						return err
					}
					url := base(c, portal, cmd.String("project")) + "/tasklists/" + cmd.Args().First() + "/followers"
					raw, err := c.Request("DELETE", url, nil)
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
		},
	}
}

func timelogsCmd() *cli.Command {
	return &cli.Command{
		Name:  "timelogs",
		Usage: "Project timelog operations",
		Commands: []*cli.Command{
			{
				Name:  "list",
				Usage: "List project timelogs",
				Flags: []cli.Flag{
					portalFlag, projectFlag,
					&cli.StringFlag{Name: "module", Value: "general", Usage: "task, issue, or general"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					portal, err := requirePortal(cmd)
					if err != nil {
						return err
					}
					url := base(c, portal, cmd.String("project")) + "/timelogs"
					moduleJSON := fmt.Sprintf(`{"type":"%s"}`, cmd.String("module"))
					raw, err := c.Request("GET", url, &zohttp.RequestOpts{
						Params: map[string]string{
							"module":    moduleJSON,
							"view_type": "projectspan",
						},
					})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:  "add",
				Usage: "Add a timelog",
				Flags: []cli.Flag{
					portalFlag, projectFlag,
					&cli.StringFlag{Name: "date", Required: true, Usage: "Date (YYYY-MM-DD)"},
					&cli.StringFlag{Name: "hours", Required: true, Usage: "Hours (e.g. 2, 1.5, 0:30)"},
					&cli.StringFlag{Name: "task", Usage: "Task ID"},
					&cli.StringFlag{Name: "owner", Usage: "Owner ZPUID"},
					&cli.StringFlag{Name: "bill-status", Value: "Billable", Usage: "Billable or Non Billable"},
					&cli.StringFlag{Name: "notes", Usage: "Notes for time entry"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					portal, err := requirePortal(cmd)
					if err != nil {
						return err
					}
					logItem := map[string]string{
						"project_id":  cmd.String("project"),
						"date":        cmd.String("date"),
						"hours":       cmd.String("hours"),
						"bill_status": cmd.String("bill-status"),
						"log_name":    "Time log",
						"type":        "general",
					}
					if n := cmd.String("notes"); n != "" {
						logItem["notes"] = n
						logItem["log_name"] = n
					}
					if t := cmd.String("task"); t != "" {
						logItem["type"] = "task"
						logItem["item_id"] = t
					}
					if o := cmd.String("owner"); o != "" {
						logItem["owner_zpuid"] = o
					}
					logBytes, err := json.Marshal([]map[string]string{logItem})
					if err != nil {
						return err
					}
					url := c.ProjectsBase + "/portal/" + portal + "/addbulktimelogs"
					raw, err := c.Request("POST", url, &zohttp.RequestOpts{
						Form: map[string]string{
							"log_object": string(logBytes),
						},
					})
					if err != nil {
						return err
					}
					var envelope struct {
						TimeLogs []struct {
							LogDetails []json.RawMessage `json:"log_details"`
						} `json:"time_logs"`
					}
					if json.Unmarshal(raw, &envelope) == nil &&
						len(envelope.TimeLogs) > 0 &&
						len(envelope.TimeLogs[0].LogDetails) > 0 {
						return output.JSONRaw(envelope.TimeLogs[0].LogDetails[0])
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "get",
				Usage:     "Get a timelog",
				ArgsUsage: "<timelog-id>",
				Flags: []cli.Flag{
					portalFlag, projectFlag,
					&cli.StringFlag{Name: "type", Value: "task", Usage: "task, issue, or general"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					portal, err := requirePortal(cmd)
					if err != nil {
						return err
					}
					url := base(c, portal, cmd.String("project")) + "/logs/" + cmd.Args().First()
					raw, err := c.Request("GET", url, &zohttp.RequestOpts{Params: map[string]string{"type": cmd.String("type")}})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "update",
				Usage:     "Update a timelog",
				ArgsUsage: "<timelog-id>",
				Flags: []cli.Flag{
					portalFlag, projectFlag,
					&cli.StringFlag{Name: "type", Value: "task", Usage: "task, issue, or general"},
					&cli.StringFlag{Name: "task", Usage: "Task ID (for module)"},
					&cli.StringFlag{Name: "json", Required: true, Usage: "Fields to update as JSON"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					portal, err := requirePortal(cmd)
					if err != nil {
						return err
					}
					var body map[string]any
					json.Unmarshal([]byte(cmd.String("json")), &body)
					if body == nil {
						body = map[string]any{}
					}
					if _, ok := body["module"]; !ok {
						mod := map[string]string{"type": cmd.String("type")}
						if t := cmd.String("task"); t != "" {
							mod["id"] = t
						}
						body["module"] = mod
					}
					url := base(c, portal, cmd.String("project")) + "/logs/" + cmd.Args().First()
					raw, err := c.Request("PATCH", url, &zohttp.RequestOpts{JSON: body})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "delete",
				Usage:     "Delete a timelog",
				ArgsUsage: "<timelog-id>",
				Flags: []cli.Flag{
					portalFlag, projectFlag,
					&cli.StringFlag{Name: "type", Value: "task", Usage: "task, issue, or general"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					portal, err := requirePortal(cmd)
					if err != nil {
						return err
					}
					url := base(c, portal, cmd.String("project")) + "/logs/" + cmd.Args().First()
					raw, err := c.Request("DELETE", url, &zohttp.RequestOpts{JSON: map[string]string{"module": cmd.String("type")}})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
		},
	}
}

func timelogBulkCmd() *cli.Command {
	return &cli.Command{
		Name:  "timelog-bulk",
		Usage: "Bulk timelog operations",
		Commands: []*cli.Command{
			{
				Name:  "list",
				Usage: "List timelogs (portal-level)",
				Flags: []cli.Flag{portalFlag},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					portal, err := requirePortal(cmd)
					if err != nil {
						return err
					}
					url := c.ProjectsBase + "/portal/" + portal + "/timelogs"
					raw, err := c.Request("GET", url, nil)
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:  "project-list",
				Usage: "List timelogs (project-level)",
				Flags: []cli.Flag{portalFlag, projectFlag},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					portal, err := requirePortal(cmd)
					if err != nil {
						return err
					}
					url := base(c, portal, cmd.String("project")) + "/timelogs"
					raw, err := c.Request("GET", url, nil)
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:  "add",
				Usage: "Bulk add timelogs",
				Flags: []cli.Flag{
					portalFlag,
					&cli.StringFlag{Name: "json", Required: true, Usage: "Timelogs as JSON (log_object array)"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					portal, err := requirePortal(cmd)
					if err != nil {
						return err
					}
					url := c.ProjectsBase + "/portal/" + portal + "/addbulktimelogs"
					raw, err := c.Request("POST", url, &zohttp.RequestOpts{
						Form: map[string]string{
							"log_object": cmd.String("json"),
						},
					})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:  "delete",
				Usage: "Bulk delete timelogs",
				Flags: []cli.Flag{
					portalFlag,
					&cli.StringFlag{Name: "json", Required: true, Usage: "Array of {id, module} objects as JSON"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					portal, err := requirePortal(cmd)
					if err != nil {
						return err
					}
					var body any
					json.Unmarshal([]byte(cmd.String("json")), &body)
					url := c.ProjectsBase + "/portal/" + portal + "/timelogs/bulkdelete"
					raw, err := c.Request("DELETE", url, &zohttp.RequestOpts{JSON: body})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
		},
	}
}

func timelogTimersCmd() *cli.Command {
	return &cli.Command{
		Name:  "timelog-timers",
		Usage: "Timer operations",
		Commands: []*cli.Command{
			{
				Name:  "running",
				Usage: "Get running timers",
				Flags: []cli.Flag{portalFlag},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					portal, err := requirePortal(cmd)
					if err != nil {
						return err
					}
					url := c.ProjectsBase + "/portal/" + portal + "/timers"
					raw, err := c.Request("GET", url, nil)
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:  "start",
				Usage: "Start a timer",
				Flags: []cli.Flag{
					portalFlag, projectFlag,
					&cli.StringFlag{Name: "json", Required: true, Usage: "Timer details as JSON"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					portal, err := requirePortal(cmd)
					if err != nil {
						return err
					}
					var body any
					json.Unmarshal([]byte(cmd.String("json")), &body)
					url := base(c, portal, cmd.String("project")) + "/timers"
					raw, err := c.Request("POST", url, &zohttp.RequestOpts{JSON: body})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "get",
				Usage:     "Get timer for a timelog",
				ArgsUsage: "<timelog-id>",
				Flags:     []cli.Flag{portalFlag, projectFlag},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					portal, err := requirePortal(cmd)
					if err != nil {
						return err
					}
					url := base(c, portal, cmd.String("project")) + "/timelogs/" + cmd.Args().First() + "/timer"
					raw, err := c.Request("GET", url, nil)
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "pause",
				Usage:     "Pause a timer",
				ArgsUsage: "<timelog-id>",
				Flags:     []cli.Flag{portalFlag, projectFlag},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					portal, err := requirePortal(cmd)
					if err != nil {
						return err
					}
					url := base(c, portal, cmd.String("project")) + "/timelogs/" + cmd.Args().First() + "/timer/pause"
					raw, err := c.Request("POST", url, nil)
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "resume",
				Usage:     "Resume a timer",
				ArgsUsage: "<timelog-id>",
				Flags:     []cli.Flag{portalFlag, projectFlag},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					portal, err := requirePortal(cmd)
					if err != nil {
						return err
					}
					url := base(c, portal, cmd.String("project")) + "/timelogs/" + cmd.Args().First() + "/timer/resume"
					raw, err := c.Request("POST", url, nil)
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "stop",
				Usage:     "Stop a timer",
				ArgsUsage: "<timelog-id>",
				Flags:     []cli.Flag{portalFlag, projectFlag},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					portal, err := requirePortal(cmd)
					if err != nil {
						return err
					}
					url := base(c, portal, cmd.String("project")) + "/timelogs/" + cmd.Args().First() + "/timer/stop"
					raw, err := c.Request("POST", url, nil)
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "delete",
				Usage:     "Delete a timer",
				ArgsUsage: "<timelog-id>",
				Flags:     []cli.Flag{portalFlag, projectFlag},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					portal, err := requirePortal(cmd)
					if err != nil {
						return err
					}
					url := base(c, portal, cmd.String("project")) + "/timelogs/" + cmd.Args().First() + "/timer"
					raw, err := c.Request("DELETE", url, nil)
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
		},
	}
}

func timelogPinsCmd() *cli.Command {
	return &cli.Command{
		Name:  "timelog-pins",
		Usage: "Timelog pin operations",
		Commands: []*cli.Command{
			{
				Name:  "list",
				Usage: "List timelog pins (portal-level)",
				Flags: []cli.Flag{portalFlag},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					portal, err := requirePortal(cmd)
					if err != nil {
						return err
					}
					url := c.ProjectsBase + "/portal/" + portal + "/timelogs/pins"
					raw, err := c.Request("GET", url, nil)
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:  "project-list",
				Usage: "List timelog pins (project-level)",
				Flags: []cli.Flag{portalFlag, projectFlag},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					portal, err := requirePortal(cmd)
					if err != nil {
						return err
					}
					url := base(c, portal, cmd.String("project")) + "/timelogs/pins"
					raw, err := c.Request("GET", url, nil)
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "create",
				Usage:     "Pin a timelog",
				ArgsUsage: "<timelog-id>",
				Flags:     []cli.Flag{portalFlag, projectFlag},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					portal, err := requirePortal(cmd)
					if err != nil {
						return err
					}
					url := base(c, portal, cmd.String("project")) + "/timelogs/" + cmd.Args().First() + "/pin"
					raw, err := c.Request("POST", url, nil)
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "update",
				Usage:     "Update a timelog pin",
				ArgsUsage: "<timelog-id>",
				Flags: []cli.Flag{
					portalFlag, projectFlag,
					&cli.StringFlag{Name: "json", Required: true, Usage: "Pin details as JSON"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					portal, err := requirePortal(cmd)
					if err != nil {
						return err
					}
					var body any
					json.Unmarshal([]byte(cmd.String("json")), &body)
					url := base(c, portal, cmd.String("project")) + "/timelogs/" + cmd.Args().First() + "/pin"
					raw, err := c.Request("PATCH", url, &zohttp.RequestOpts{JSON: body})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "delete",
				Usage:     "Unpin a timelog",
				ArgsUsage: "<timelog-id>",
				Flags:     []cli.Flag{portalFlag, projectFlag},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					portal, err := requirePortal(cmd)
					if err != nil {
						return err
					}
					url := base(c, portal, cmd.String("project")) + "/timelogs/" + cmd.Args().First() + "/pin"
					raw, err := c.Request("DELETE", url, nil)
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
		Usage: "Portal user operations",
		Commands: []*cli.Command{
			{
				Name:  "list",
				Usage: "List portal users",
				Flags: []cli.Flag{portalFlag},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					portal, err := requirePortal(cmd)
					if err != nil {
						return err
					}
					url := c.ProjectsBase + "/portal/" + portal + "/users"
					items, err := pagination.PaginateProjects(c, url, "users", nil, 0)
					if err != nil {
						return err
					}
					return output.JSON(items)
				},
			},
			{
				Name:      "get",
				Usage:     "Get a user",
				ArgsUsage: "<user-id>",
				Flags:     []cli.Flag{portalFlag},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					portal, err := requirePortal(cmd)
					if err != nil {
						return err
					}
					url := c.ProjectsBase + "/portal/" + portal + "/users/" + cmd.Args().First()
					raw, err := c.Request("GET", url, nil)
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:  "add",
				Usage: "Add a user to portal",
				Flags: []cli.Flag{
					portalFlag,
					&cli.StringFlag{Name: "json", Required: true, Usage: "User details as JSON"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					portal, err := requirePortal(cmd)
					if err != nil {
						return err
					}
					var body any
					json.Unmarshal([]byte(cmd.String("json")), &body)
					url := c.ProjectsBase + "/portal/" + portal + "/users"
					raw, err := c.Request("POST", url, &zohttp.RequestOpts{JSON: body})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "activate",
				Usage:     "Activate a user",
				ArgsUsage: "<user-id>",
				Flags:     []cli.Flag{portalFlag},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					portal, err := requirePortal(cmd)
					if err != nil {
						return err
					}
					url := c.ProjectsBase + "/portal/" + portal + "/users/" + cmd.Args().First() + "/active"
					raw, err := c.Request("POST", url, nil)
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "deactivate",
				Usage:     "Deactivate a user",
				ArgsUsage: "<user-id>",
				Flags:     []cli.Flag{portalFlag},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					portal, err := requirePortal(cmd)
					if err != nil {
						return err
					}
					url := c.ProjectsBase + "/portal/" + portal + "/users/" + cmd.Args().First() + "/inactive"
					raw, err := c.Request("POST", url, nil)
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "delete",
				Usage:     "Delete a user from portal",
				ArgsUsage: "<user-id>",
				Flags:     []cli.Flag{portalFlag},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					portal, err := requirePortal(cmd)
					if err != nil {
						return err
					}
					url := c.ProjectsBase + "/portal/" + portal + "/users/" + cmd.Args().First()
					raw, err := c.Request("DELETE", url, nil)
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
		},
	}
}

func projectUsersCmd() *cli.Command {
	return &cli.Command{
		Name:  "project-users",
		Usage: "Project-scoped user operations",
		Commands: []*cli.Command{
			{
				Name:  "list",
				Usage: "List project users",
				Flags: []cli.Flag{portalFlag, projectFlag},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					portal, err := requirePortal(cmd)
					if err != nil {
						return err
					}
					url := base(c, portal, cmd.String("project")) + "/users"
					items, err := pagination.PaginateProjects(c, url, "users", nil, 0)
					if err != nil {
						return err
					}
					return output.JSON(items)
				},
			},
			{
				Name:      "get",
				Usage:     "Get a project user",
				ArgsUsage: "<user-id>",
				Flags:     []cli.Flag{portalFlag, projectFlag},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					portal, err := requirePortal(cmd)
					if err != nil {
						return err
					}
					url := base(c, portal, cmd.String("project")) + "/users/" + cmd.Args().First()
					raw, err := c.Request("GET", url, nil)
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:  "add",
				Usage: "Add a user to project",
				Flags: []cli.Flag{
					portalFlag, projectFlag,
					&cli.StringFlag{Name: "json", Required: true, Usage: "User details as JSON"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					portal, err := requirePortal(cmd)
					if err != nil {
						return err
					}
					var body any
					json.Unmarshal([]byte(cmd.String("json")), &body)
					url := base(c, portal, cmd.String("project")) + "/users"
					raw, err := c.Request("POST", url, &zohttp.RequestOpts{JSON: body})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "update",
				Usage:     "Update a project user",
				ArgsUsage: "<user-id>",
				Flags: []cli.Flag{
					portalFlag, projectFlag,
					&cli.StringFlag{Name: "json", Required: true, Usage: "User details as JSON"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					portal, err := requirePortal(cmd)
					if err != nil {
						return err
					}
					var body any
					json.Unmarshal([]byte(cmd.String("json")), &body)
					url := base(c, portal, cmd.String("project")) + "/users/" + cmd.Args().First()
					raw, err := c.Request("PATCH", url, &zohttp.RequestOpts{JSON: body})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "delete",
				Usage:     "Remove a user from project",
				ArgsUsage: "<user-id>",
				Flags:     []cli.Flag{portalFlag, projectFlag},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					portal, err := requirePortal(cmd)
					if err != nil {
						return err
					}
					url := base(c, portal, cmd.String("project")) + "/users/" + cmd.Args().First()
					raw, err := c.Request("DELETE", url, nil)
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
		},
	}
}

func milestonesCmd() *cli.Command {
	return &cli.Command{
		Name:  "milestones",
		Usage: "Project milestone operations",
		Commands: []*cli.Command{
			{
				Name:  "list",
				Usage: "List milestones",
				Flags: []cli.Flag{portalFlag, projectFlag},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					portal, err := requirePortal(cmd)
					if err != nil {
						return err
					}
					url := base(c, portal, cmd.String("project")) + "/milestones"
					items, err := pagination.PaginateProjects(c, url, "milestones", nil, 0)
					if err != nil {
						return err
					}
					return output.JSON(items)
				},
			},
			{
				Name:      "get",
				Usage:     "Get a milestone",
				ArgsUsage: "<milestone-id>",
				Flags:     []cli.Flag{portalFlag, projectFlag},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					portal, err := requirePortal(cmd)
					if err != nil {
						return err
					}
					url := base(c, portal, cmd.String("project")) + "/milestones/" + cmd.Args().First()
					raw, err := c.Request("GET", url, nil)
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
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
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
					if j := cmd.String("json"); j != "" {
						var extra map[string]any
						json.Unmarshal([]byte(j), &extra)
						for k, v := range extra {
							body[k] = v
						}
					}
					url := base(c, portal, cmd.String("project")) + "/milestones"
					raw, err := c.Request("POST", url, &zohttp.RequestOpts{JSON: body})
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
					&cli.StringFlag{Name: "json", Required: true, Usage: "Fields to update as JSON"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					portal, err := requirePortal(cmd)
					if err != nil {
						return err
					}
					var parsed map[string]any
					json.Unmarshal([]byte(cmd.String("json")), &parsed)
					url := base(c, portal, cmd.String("project")) + "/milestones/" + cmd.Args().First()
					raw, err := c.Request("POST", url, &zohttp.RequestOpts{JSON: parsed})
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
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					portal, err := requirePortal(cmd)
					if err != nil {
						return err
					}
					url := base(c, portal, cmd.String("project")) + "/milestones/" + cmd.Args().First()
					raw, err := c.Request("DELETE", url, nil)
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
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					portal, err := requirePortal(cmd)
					if err != nil {
						return err
					}
					url := c.ProjectsBase + "/portal/" + portal + "/phases"
					raw, err := c.Request("GET", url, nil)
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
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					portal, err := requirePortal(cmd)
					if err != nil {
						return err
					}
					url := base(c, portal, cmd.String("project")) + "/phases"
					raw, err := c.Request("GET", url, nil)
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
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					portal, err := requirePortal(cmd)
					if err != nil {
						return err
					}
					url := base(c, portal, cmd.String("project")) + "/phases/" + cmd.Args().First()
					raw, err := c.Request("GET", url, nil)
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
					&cli.StringFlag{Name: "json", Required: true, Usage: "Phase details as JSON"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					portal, err := requirePortal(cmd)
					if err != nil {
						return err
					}
					var body any
					json.Unmarshal([]byte(cmd.String("json")), &body)
					url := base(c, portal, cmd.String("project")) + "/phases"
					raw, err := c.Request("POST", url, &zohttp.RequestOpts{JSON: body})
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
					&cli.StringFlag{Name: "json", Required: true, Usage: "Fields to update as JSON"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					portal, err := requirePortal(cmd)
					if err != nil {
						return err
					}
					var body any
					json.Unmarshal([]byte(cmd.String("json")), &body)
					url := base(c, portal, cmd.String("project")) + "/phases/" + cmd.Args().First()
					raw, err := c.Request("PATCH", url, &zohttp.RequestOpts{JSON: body})
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
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					portal, err := requirePortal(cmd)
					if err != nil {
						return err
					}
					url := base(c, portal, cmd.String("project")) + "/phases/" + cmd.Args().First()
					raw, err := c.Request("DELETE", url, nil)
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
					&cli.StringFlag{Name: "json", Required: true, Usage: "Move details as JSON"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					portal, err := requirePortal(cmd)
					if err != nil {
						return err
					}
					var body any
					json.Unmarshal([]byte(cmd.String("json")), &body)
					url := base(c, portal, cmd.String("project")) + "/phases/" + cmd.Args().First() + "/move"
					raw, err := c.Request("POST", url, &zohttp.RequestOpts{JSON: body})
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
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					portal, err := requirePortal(cmd)
					if err != nil {
						return err
					}
					url := base(c, portal, cmd.String("project")) + "/phases/" + cmd.Args().First() + "/clone"
					raw, err := c.Request("POST", url, nil)
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
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					portal, err := requirePortal(cmd)
					if err != nil {
						return err
					}
					url := base(c, portal, cmd.String("project")) + "/phases/" + cmd.Args().First() + "/activities"
					raw, err := c.Request("GET", url, nil)
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
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					portal, err := requirePortal(cmd)
					if err != nil {
						return err
					}
					url := base(c, portal, cmd.String("project")) + "/phases/" + cmd.Args().First() + "/followers"
					raw, err := c.Request("GET", url, nil)
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
					&cli.StringFlag{Name: "json", Required: true, Usage: "Follower details as JSON"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					portal, err := requirePortal(cmd)
					if err != nil {
						return err
					}
					var body any
					json.Unmarshal([]byte(cmd.String("json")), &body)
					url := base(c, portal, cmd.String("project")) + "/phases/" + cmd.Args().First() + "/followers"
					raw, err := c.Request("POST", url, &zohttp.RequestOpts{JSON: body})
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
					&cli.StringFlag{Name: "json", Required: true, Usage: "Follower IDs as JSON"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					portal, err := requirePortal(cmd)
					if err != nil {
						return err
					}
					var body any
					json.Unmarshal([]byte(cmd.String("json")), &body)
					url := base(c, portal, cmd.String("project")) + "/phases/" + cmd.Args().First() + "/followers"
					raw, err := c.Request("DELETE", url, &zohttp.RequestOpts{JSON: body})
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
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					portal, err := requirePortal(cmd)
					if err != nil {
						return err
					}
					url := base(c, portal, cmd.String("project")) + "/phases/" + cmd.String("phase") + "/comments"
					raw, err := c.Request("GET", url, nil)
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
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					portal, err := requirePortal(cmd)
					if err != nil {
						return err
					}
					url := base(c, portal, cmd.String("project")) + "/phases/" + cmd.String("phase") + "/comments"
					raw, err := c.Request("POST", url, &zohttp.RequestOpts{JSON: map[string]string{"comment": cmd.String("comment")}})
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
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					portal, err := requirePortal(cmd)
					if err != nil {
						return err
					}
					url := base(c, portal, cmd.String("project")) + "/phases/" + cmd.String("phase") + "/comments/" + cmd.Args().First()
					raw, err := c.Request("PATCH", url, &zohttp.RequestOpts{JSON: map[string]string{"comment": cmd.String("comment")}})
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
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					portal, err := requirePortal(cmd)
					if err != nil {
						return err
					}
					url := base(c, portal, cmd.String("project")) + "/phases/" + cmd.String("phase") + "/comments/" + cmd.Args().First()
					raw, err := c.Request("DELETE", url, nil)
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
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
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
					raw, err := c.Request("POST", url, &zohttp.RequestOpts{JSON: body})
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
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					portal, err := requirePortal(cmd)
					if err != nil {
						return err
					}
					url := base(c, portal, cmd.String("project")) + "/tasks/" + cmd.Args().Get(0) + "/dependencies/" + cmd.Args().Get(1)
					raw, err := c.Request("DELETE", url, nil)
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
		},
	}
}

func forumsCmd() *cli.Command {
	return &cli.Command{
		Name:  "forums",
		Usage: "Forum operations",
		Commands: []*cli.Command{
			{
				Name:  "list",
				Usage: "List forums",
				Flags: []cli.Flag{portalFlag, projectFlag},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					portal, err := requirePortal(cmd)
					if err != nil {
						return err
					}
					url := base(c, portal, cmd.String("project")) + "/forums"
					raw, err := c.Request("GET", url, nil)
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "get",
				Usage:     "Get a forum",
				ArgsUsage: "<forum-id>",
				Flags:     []cli.Flag{portalFlag, projectFlag},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					portal, err := requirePortal(cmd)
					if err != nil {
						return err
					}
					url := base(c, portal, cmd.String("project")) + "/forums/" + cmd.Args().First()
					raw, err := c.Request("GET", url, nil)
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:  "create",
				Usage: "Create a forum",
				Flags: []cli.Flag{
					portalFlag, projectFlag,
					&cli.StringFlag{Name: "json", Required: true, Usage: "Forum details as JSON"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					portal, err := requirePortal(cmd)
					if err != nil {
						return err
					}
					var body any
					json.Unmarshal([]byte(cmd.String("json")), &body)
					url := base(c, portal, cmd.String("project")) + "/forums"
					raw, err := c.Request("POST", url, &zohttp.RequestOpts{JSON: body})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "update",
				Usage:     "Update a forum",
				ArgsUsage: "<forum-id>",
				Flags: []cli.Flag{
					portalFlag, projectFlag,
					&cli.StringFlag{Name: "json", Required: true, Usage: "Fields to update as JSON"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					portal, err := requirePortal(cmd)
					if err != nil {
						return err
					}
					var body any
					json.Unmarshal([]byte(cmd.String("json")), &body)
					url := base(c, portal, cmd.String("project")) + "/forums/" + cmd.Args().First()
					raw, err := c.Request("PATCH", url, &zohttp.RequestOpts{JSON: body})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "delete",
				Usage:     "Delete a forum",
				ArgsUsage: "<forum-id>",
				Flags:     []cli.Flag{portalFlag, projectFlag},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					portal, err := requirePortal(cmd)
					if err != nil {
						return err
					}
					url := base(c, portal, cmd.String("project")) + "/forums/" + cmd.Args().First()
					raw, err := c.Request("DELETE", url, nil)
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
		},
	}
}

func forumCommentsCmd() *cli.Command {
	return &cli.Command{
		Name:  "forum-comments",
		Usage: "Forum comment operations",
		Commands: []*cli.Command{
			{
				Name:  "list",
				Usage: "List forum comments",
				Flags: []cli.Flag{
					portalFlag, projectFlag,
					&cli.StringFlag{Name: "forum", Required: true, Usage: "Forum ID"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					portal, err := requirePortal(cmd)
					if err != nil {
						return err
					}
					url := base(c, portal, cmd.String("project")) + "/forums/" + cmd.String("forum") + "/comments"
					raw, err := c.Request("GET", url, nil)
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "get",
				Usage:     "Get a forum comment",
				ArgsUsage: "<comment-id>",
				Flags: []cli.Flag{
					portalFlag, projectFlag,
					&cli.StringFlag{Name: "forum", Required: true, Usage: "Forum ID"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					portal, err := requirePortal(cmd)
					if err != nil {
						return err
					}
					url := base(c, portal, cmd.String("project")) + "/forums/" + cmd.String("forum") + "/comments/" + cmd.Args().First()
					raw, err := c.Request("GET", url, nil)
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:  "add",
				Usage: "Add a forum comment",
				Flags: []cli.Flag{
					portalFlag, projectFlag,
					&cli.StringFlag{Name: "forum", Required: true, Usage: "Forum ID"},
					&cli.StringFlag{Name: "comment", Required: true, Usage: "Comment text"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					portal, err := requirePortal(cmd)
					if err != nil {
						return err
					}
					url := base(c, portal, cmd.String("project")) + "/forums/" + cmd.String("forum") + "/comments"
					raw, err := c.Request("POST", url, &zohttp.RequestOpts{JSON: map[string]string{"comment": cmd.String("comment")}})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "update",
				Usage:     "Update a forum comment",
				ArgsUsage: "<comment-id>",
				Flags: []cli.Flag{
					portalFlag, projectFlag,
					&cli.StringFlag{Name: "forum", Required: true, Usage: "Forum ID"},
					&cli.StringFlag{Name: "comment", Required: true, Usage: "Updated comment text"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					portal, err := requirePortal(cmd)
					if err != nil {
						return err
					}
					url := base(c, portal, cmd.String("project")) + "/forums/" + cmd.String("forum") + "/comments/" + cmd.Args().First()
					raw, err := c.Request("PATCH", url, &zohttp.RequestOpts{JSON: map[string]string{"comment": cmd.String("comment")}})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "delete",
				Usage:     "Delete a forum comment",
				ArgsUsage: "<comment-id>",
				Flags: []cli.Flag{
					portalFlag, projectFlag,
					&cli.StringFlag{Name: "forum", Required: true, Usage: "Forum ID"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					portal, err := requirePortal(cmd)
					if err != nil {
						return err
					}
					url := base(c, portal, cmd.String("project")) + "/forums/" + cmd.String("forum") + "/comments/" + cmd.Args().First()
					raw, err := c.Request("DELETE", url, nil)
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "best-answer",
				Usage:     "Mark as best answer",
				ArgsUsage: "<comment-id>",
				Flags: []cli.Flag{
					portalFlag, projectFlag,
					&cli.StringFlag{Name: "forum", Required: true, Usage: "Forum ID"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					portal, err := requirePortal(cmd)
					if err != nil {
						return err
					}
					url := base(c, portal, cmd.String("project")) + "/forums/" + cmd.String("forum") + "/comments/" + cmd.Args().First() + "/bestanswer"
					raw, err := c.Request("POST", url, nil)
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "unbest-answer",
				Usage:     "Unmark best answer",
				ArgsUsage: "<comment-id>",
				Flags: []cli.Flag{
					portalFlag, projectFlag,
					&cli.StringFlag{Name: "forum", Required: true, Usage: "Forum ID"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					portal, err := requirePortal(cmd)
					if err != nil {
						return err
					}
					url := base(c, portal, cmd.String("project")) + "/forums/" + cmd.String("forum") + "/comments/" + cmd.Args().First() + "/bestanswer"
					raw, err := c.Request("DELETE", url, nil)
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
		},
	}
}

func forumCategoriesCmd() *cli.Command {
	return &cli.Command{
		Name:  "forum-categories",
		Usage: "Forum category operations",
		Commands: []*cli.Command{
			{
				Name:  "list",
				Usage: "List forum categories",
				Flags: []cli.Flag{portalFlag, projectFlag},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					portal, err := requirePortal(cmd)
					if err != nil {
						return err
					}
					url := base(c, portal, cmd.String("project")) + "/forums/categories"
					raw, err := c.Request("GET", url, nil)
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:  "create",
				Usage: "Create a forum category",
				Flags: []cli.Flag{
					portalFlag, projectFlag,
					&cli.StringFlag{Name: "json", Required: true, Usage: "Category details as JSON"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					portal, err := requirePortal(cmd)
					if err != nil {
						return err
					}
					var body any
					json.Unmarshal([]byte(cmd.String("json")), &body)
					url := base(c, portal, cmd.String("project")) + "/forums/categories"
					raw, err := c.Request("POST", url, &zohttp.RequestOpts{JSON: body})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "update",
				Usage:     "Update a forum category",
				ArgsUsage: "<category-id>",
				Flags: []cli.Flag{
					portalFlag, projectFlag,
					&cli.StringFlag{Name: "json", Required: true, Usage: "Category details as JSON"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					portal, err := requirePortal(cmd)
					if err != nil {
						return err
					}
					var body any
					json.Unmarshal([]byte(cmd.String("json")), &body)
					url := base(c, portal, cmd.String("project")) + "/forums/categories/" + cmd.Args().First()
					raw, err := c.Request("PATCH", url, &zohttp.RequestOpts{JSON: body})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "delete",
				Usage:     "Delete a forum category",
				ArgsUsage: "<category-id>",
				Flags:     []cli.Flag{portalFlag, projectFlag},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					portal, err := requirePortal(cmd)
					if err != nil {
						return err
					}
					url := base(c, portal, cmd.String("project")) + "/forums/categories/" + cmd.Args().First()
					raw, err := c.Request("DELETE", url, nil)
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
		},
	}
}

func forumFollowersCmd() *cli.Command {
	return &cli.Command{
		Name:  "forum-followers",
		Usage: "Forum follower operations",
		Commands: []*cli.Command{
			{
				Name:      "list",
				Usage:     "List forum followers",
				ArgsUsage: "<forum-id>",
				Flags:     []cli.Flag{portalFlag, projectFlag},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					portal, err := requirePortal(cmd)
					if err != nil {
						return err
					}
					url := base(c, portal, cmd.String("project")) + "/forums/" + cmd.Args().First() + "/followers"
					raw, err := c.Request("GET", url, nil)
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "follow",
				Usage:     "Follow a forum",
				ArgsUsage: "<forum-id>",
				Flags:     []cli.Flag{portalFlag, projectFlag},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					portal, err := requirePortal(cmd)
					if err != nil {
						return err
					}
					url := base(c, portal, cmd.String("project")) + "/forums/" + cmd.Args().First() + "/followers"
					raw, err := c.Request("POST", url, nil)
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "unfollow",
				Usage:     "Unfollow a forum",
				ArgsUsage: "<forum-id>",
				Flags:     []cli.Flag{portalFlag, projectFlag},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					portal, err := requirePortal(cmd)
					if err != nil {
						return err
					}
					url := base(c, portal, cmd.String("project")) + "/forums/" + cmd.Args().First() + "/followers"
					raw, err := c.Request("DELETE", url, nil)
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
		},
	}
}

func eventsCmd() *cli.Command {
	return &cli.Command{
		Name:  "events",
		Usage: "Event operations",
		Commands: []*cli.Command{
			{
				Name:  "list",
				Usage: "List events",
				Flags: []cli.Flag{portalFlag, projectFlag},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					portal, err := requirePortal(cmd)
					if err != nil {
						return err
					}
					url := base(c, portal, cmd.String("project")) + "/events"
					raw, err := c.Request("GET", url, nil)
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
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					portal, err := requirePortal(cmd)
					if err != nil {
						return err
					}
					url := base(c, portal, cmd.String("project")) + "/events/" + cmd.Args().First()
					raw, err := c.Request("GET", url, nil)
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
					&cli.StringFlag{Name: "json", Required: true, Usage: "Event details as JSON"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					portal, err := requirePortal(cmd)
					if err != nil {
						return err
					}
					var body any
					json.Unmarshal([]byte(cmd.String("json")), &body)
					url := base(c, portal, cmd.String("project")) + "/events"
					raw, err := c.Request("POST", url, &zohttp.RequestOpts{JSON: body})
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
					&cli.StringFlag{Name: "json", Required: true, Usage: "Fields to update as JSON"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					portal, err := requirePortal(cmd)
					if err != nil {
						return err
					}
					var body any
					json.Unmarshal([]byte(cmd.String("json")), &body)
					url := base(c, portal, cmd.String("project")) + "/events/" + cmd.Args().First()
					raw, err := c.Request("PATCH", url, &zohttp.RequestOpts{JSON: body})
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
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					portal, err := requirePortal(cmd)
					if err != nil {
						return err
					}
					url := base(c, portal, cmd.String("project")) + "/events/" + cmd.Args().First()
					raw, err := c.Request("DELETE", url, nil)
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
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					portal, err := requirePortal(cmd)
					if err != nil {
						return err
					}
					url := base(c, portal, cmd.String("project")) + "/events/" + cmd.String("event") + "/comments"
					raw, err := c.Request("GET", url, nil)
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
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					portal, err := requirePortal(cmd)
					if err != nil {
						return err
					}
					url := base(c, portal, cmd.String("project")) + "/events/" + cmd.String("event") + "/comments/" + cmd.Args().First()
					raw, err := c.Request("GET", url, nil)
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
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					portal, err := requirePortal(cmd)
					if err != nil {
						return err
					}
					url := base(c, portal, cmd.String("project")) + "/events/" + cmd.String("event") + "/comments"
					raw, err := c.Request("POST", url, &zohttp.RequestOpts{JSON: map[string]string{"comment": cmd.String("comment")}})
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
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					portal, err := requirePortal(cmd)
					if err != nil {
						return err
					}
					url := base(c, portal, cmd.String("project")) + "/events/" + cmd.String("event") + "/comments/" + cmd.Args().First()
					raw, err := c.Request("PATCH", url, &zohttp.RequestOpts{JSON: map[string]string{"comment": cmd.String("comment")}})
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
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					portal, err := requirePortal(cmd)
					if err != nil {
						return err
					}
					url := base(c, portal, cmd.String("project")) + "/events/" + cmd.String("event") + "/comments/" + cmd.Args().First()
					raw, err := c.Request("DELETE", url, nil)
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
					&cli.StringFlag{Name: "type", Required: true, Usage: "Entity type (project, task, issue, forum, etc.)"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					portal, err := requirePortal(cmd)
					if err != nil {
						return err
					}
					url := base(c, portal, cmd.String("project")) + "/attachments"
					raw, err := c.Request("GET", url, &zohttp.RequestOpts{
						Params: map[string]string{
							"entity_type": cmd.String("type"),
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
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					portal, err := requirePortal(cmd)
					if err != nil {
						return err
					}
					url := base(c, portal, cmd.String("project")) + "/attachments/" + cmd.Args().First()
					raw, err := c.Request("GET", url, nil)
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
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
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
					raw, err := c.Request("POST", url, &zohttp.RequestOpts{
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
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					portal, err := requirePortal(cmd)
					if err != nil {
						return err
					}
					url := base(c, portal, cmd.String("project")) + "/attachments/" + cmd.Args().First()
					raw, err := c.Request("POST", url, &zohttp.RequestOpts{
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
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					portal, err := requirePortal(cmd)
					if err != nil {
						return err
					}
					url := base(c, portal, cmd.String("project")) + "/attachments/" + cmd.Args().First()
					raw, err := c.Request("DELETE", url, nil)
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
		},
	}
}

func leavesCmd() *cli.Command {
	return &cli.Command{
		Name:  "leaves",
		Usage: "Leave operations",
		Commands: []*cli.Command{
			{
				Name:  "list",
				Usage: "List leaves",
				Flags: []cli.Flag{portalFlag},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					portal, err := requirePortal(cmd)
					if err != nil {
						return err
					}
					url := c.ProjectsBase + "/portal/" + portal + "/leaves"
					raw, err := c.Request("GET", url, nil)
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
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					portal, err := requirePortal(cmd)
					if err != nil {
						return err
					}
					url := c.ProjectsBase + "/portal/" + portal + "/leaves/" + cmd.Args().First()
					raw, err := c.Request("GET", url, nil)
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
					&cli.StringFlag{Name: "json", Required: true, Usage: "Leave details as JSON"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					portal, err := requirePortal(cmd)
					if err != nil {
						return err
					}
					var body any
					json.Unmarshal([]byte(cmd.String("json")), &body)
					url := c.ProjectsBase + "/portal/" + portal + "/leaves"
					raw, err := c.Request("POST", url, &zohttp.RequestOpts{JSON: body})
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
					&cli.StringFlag{Name: "json", Required: true, Usage: "Leave details as JSON"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					portal, err := requirePortal(cmd)
					if err != nil {
						return err
					}
					var body any
					json.Unmarshal([]byte(cmd.String("json")), &body)
					url := c.ProjectsBase + "/portal/" + portal + "/leaves/" + cmd.Args().First()
					raw, err := c.Request("PATCH", url, &zohttp.RequestOpts{JSON: body})
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
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					portal, err := requirePortal(cmd)
					if err != nil {
						return err
					}
					url := c.ProjectsBase + "/portal/" + portal + "/leaves/" + cmd.Args().First()
					raw, err := c.Request("DELETE", url, nil)
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
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					portal, err := requirePortal(cmd)
					if err != nil {
						return err
					}
					url := c.ProjectsBase + "/portal/" + portal + "/tags"
					raw, err := c.Request("GET", url, nil)
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
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					portal, err := requirePortal(cmd)
					if err != nil {
						return err
					}
					url := base(c, portal, cmd.String("project")) + "/tags"
					raw, err := c.Request("GET", url, nil)
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
					&cli.StringFlag{Name: "json", Required: true, Usage: "Tag details as JSON"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					portal, err := requirePortal(cmd)
					if err != nil {
						return err
					}
					var body any
					json.Unmarshal([]byte(cmd.String("json")), &body)
					url := c.ProjectsBase + "/portal/" + portal + "/tags"
					raw, err := c.Request("POST", url, &zohttp.RequestOpts{JSON: body})
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
					&cli.StringFlag{Name: "json", Required: true, Usage: "Tag details as JSON"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					portal, err := requirePortal(cmd)
					if err != nil {
						return err
					}
					var body any
					json.Unmarshal([]byte(cmd.String("json")), &body)
					url := c.ProjectsBase + "/portal/" + portal + "/tags/" + cmd.Args().First()
					raw, err := c.Request("PATCH", url, &zohttp.RequestOpts{JSON: body})
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
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					portal, err := requirePortal(cmd)
					if err != nil {
						return err
					}
					url := c.ProjectsBase + "/portal/" + portal + "/tags/" + cmd.Args().First()
					raw, err := c.Request("DELETE", url, nil)
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
					&cli.StringFlag{Name: "json", Required: true, Usage: "Entity IDs as JSON"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					portal, err := requirePortal(cmd)
					if err != nil {
						return err
					}
					var body any
					json.Unmarshal([]byte(cmd.String("json")), &body)
					url := base(c, portal, cmd.String("project")) + "/tags/" + cmd.Args().First() + "/associate"
					raw, err := c.Request("POST", url, &zohttp.RequestOpts{JSON: body})
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
					&cli.StringFlag{Name: "json", Required: true, Usage: "Entity IDs as JSON"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					portal, err := requirePortal(cmd)
					if err != nil {
						return err
					}
					var body any
					json.Unmarshal([]byte(cmd.String("json")), &body)
					url := base(c, portal, cmd.String("project")) + "/tags/" + cmd.Args().First() + "/dissociate"
					raw, err := c.Request("DELETE", url, &zohttp.RequestOpts{JSON: body})
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
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					portal, err := requirePortal(cmd)
					if err != nil {
						return err
					}
					url := c.ProjectsBase + "/portal/" + portal
					raw, err := c.Request("GET", url, nil)
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
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					portal, err := requirePortal(cmd)
					if err != nil {
						return err
					}
					url := c.ProjectsBase + "/portal/" + portal + "/bin"
					raw, err := c.Request("GET", url, nil)
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
					&cli.StringFlag{Name: "json", Required: true, Usage: "Record IDs as JSON"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					portal, err := requirePortal(cmd)
					if err != nil {
						return err
					}
					var body any
					json.Unmarshal([]byte(cmd.String("json")), &body)
					url := c.ProjectsBase + "/portal/" + portal + "/bin"
					raw, err := c.Request("DELETE", url, &zohttp.RequestOpts{JSON: body})
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
					&cli.StringFlag{Name: "json", Required: true, Usage: "Record IDs as JSON"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					portal, err := requirePortal(cmd)
					if err != nil {
						return err
					}
					var body any
					json.Unmarshal([]byte(cmd.String("json")), &body)
					url := c.ProjectsBase + "/portal/" + portal + "/bin/restore"
					raw, err := c.Request("POST", url, &zohttp.RequestOpts{JSON: body})
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
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					portal, err := requirePortal(cmd)
					if err != nil {
						return err
					}
					url := c.ProjectsBase + "/portal/" + portal + "/empty-bin"
					raw, err := c.Request("DELETE", url, nil)
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
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					portal, err := requirePortal(cmd)
					if err != nil {
						return err
					}
					url := c.ProjectsBase + "/portal/" + portal + "/search"
					raw, err := c.Request("GET", url, &zohttp.RequestOpts{
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
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					portal, err := requirePortal(cmd)
					if err != nil {
						return err
					}
					url := base(c, portal, cmd.String("project")) + "/search"
					raw, err := c.Request("GET", url, &zohttp.RequestOpts{
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
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					portal, err := requirePortal(cmd)
					if err != nil {
						return err
					}
					url := base(c, portal, cmd.String("project")) + "/status"
					raw, err := c.Request("GET", url, nil)
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
					&cli.StringFlag{Name: "json", Required: true, Usage: "Status details as JSON"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					portal, err := requirePortal(cmd)
					if err != nil {
						return err
					}
					var body any
					json.Unmarshal([]byte(cmd.String("json")), &body)
					url := base(c, portal, cmd.String("project")) + "/status"
					raw, err := c.Request("POST", url, &zohttp.RequestOpts{JSON: body})
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
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					portal, err := requirePortal(cmd)
					if err != nil {
						return err
					}
					url := base(c, portal, cmd.String("project")) + "/comments"
					raw, err := c.Request("GET", url, nil)
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
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					portal, err := requirePortal(cmd)
					if err != nil {
						return err
					}
					url := base(c, portal, cmd.String("project")) + "/comments/" + cmd.Args().First()
					raw, err := c.Request("GET", url, nil)
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
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					portal, err := requirePortal(cmd)
					if err != nil {
						return err
					}
					url := base(c, portal, cmd.String("project")) + "/comments"
					raw, err := c.Request("POST", url, &zohttp.RequestOpts{JSON: map[string]string{"comment": cmd.String("comment")}})
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
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					portal, err := requirePortal(cmd)
					if err != nil {
						return err
					}
					url := base(c, portal, cmd.String("project")) + "/comments/" + cmd.Args().First()
					raw, err := c.Request("PATCH", url, &zohttp.RequestOpts{JSON: map[string]string{"comment": cmd.String("comment")}})
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
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					portal, err := requirePortal(cmd)
					if err != nil {
						return err
					}
					url := base(c, portal, cmd.String("project")) + "/comments/" + cmd.Args().First()
					raw, err := c.Request("DELETE", url, nil)
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
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					portal, err := requirePortal(cmd)
					if err != nil {
						return err
					}
					url := c.ProjectsBase + "/portal/" + portal + "/projectgroups"
					raw, err := c.Request("GET", url, nil)
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
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					portal, err := requirePortal(cmd)
					if err != nil {
						return err
					}
					url := c.ProjectsBase + "/portal/" + portal + "/projectgroups/mygroups"
					raw, err := c.Request("GET", url, nil)
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
					&cli.StringFlag{Name: "json", Required: true, Usage: "Group details as JSON"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					portal, err := requirePortal(cmd)
					if err != nil {
						return err
					}
					var body any
					json.Unmarshal([]byte(cmd.String("json")), &body)
					url := c.ProjectsBase + "/portal/" + portal + "/projectgroups"
					raw, err := c.Request("POST", url, &zohttp.RequestOpts{JSON: body})
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
					&cli.StringFlag{Name: "json", Required: true, Usage: "Group details as JSON"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					portal, err := requirePortal(cmd)
					if err != nil {
						return err
					}
					var body any
					json.Unmarshal([]byte(cmd.String("json")), &body)
					url := c.ProjectsBase + "/portal/" + portal + "/projectgroups/" + cmd.Args().First()
					raw, err := c.Request("PATCH", url, &zohttp.RequestOpts{JSON: body})
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
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					portal, err := requirePortal(cmd)
					if err != nil {
						return err
					}
					url := c.ProjectsBase + "/portal/" + portal + "/projectgroups/" + cmd.Args().First()
					raw, err := c.Request("DELETE", url, nil)
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
				Name:      "get",
				Usage:     "Get a team",
				ArgsUsage: "<team-id>",
				Flags:     []cli.Flag{portalFlag},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					portal, err := requirePortal(cmd)
					if err != nil {
						return err
					}
					url := c.ProjectsBase + "/portal/" + portal + "/teams/" + cmd.Args().First()
					raw, err := c.Request("GET", url, nil)
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
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					portal, err := requirePortal(cmd)
					if err != nil {
						return err
					}
					url := base(c, portal, cmd.String("project")) + "/teams"
					raw, err := c.Request("GET", url, nil)
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
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					portal, err := requirePortal(cmd)
					if err != nil {
						return err
					}
					url := c.ProjectsBase + "/portal/" + portal + "/teams/" + cmd.Args().First() + "/users"
					raw, err := c.Request("GET", url, nil)
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
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					portal, err := requirePortal(cmd)
					if err != nil {
						return err
					}
					url := c.ProjectsBase + "/portal/" + portal + "/teams/" + cmd.Args().First() + "/projects"
					raw, err := c.Request("GET", url, nil)
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
					&cli.StringFlag{Name: "json", Required: true, Usage: "Team details as JSON"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					portal, err := requirePortal(cmd)
					if err != nil {
						return err
					}
					var body any
					json.Unmarshal([]byte(cmd.String("json")), &body)
					url := c.ProjectsBase + "/portal/" + portal + "/teams"
					raw, err := c.Request("POST", url, &zohttp.RequestOpts{JSON: body})
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
					&cli.StringFlag{Name: "json", Required: true, Usage: "Team details as JSON"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					portal, err := requirePortal(cmd)
					if err != nil {
						return err
					}
					var body any
					json.Unmarshal([]byte(cmd.String("json")), &body)
					url := c.ProjectsBase + "/portal/" + portal + "/teams/" + cmd.Args().First()
					raw, err := c.Request("PATCH", url, &zohttp.RequestOpts{JSON: body})
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
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					portal, err := requirePortal(cmd)
					if err != nil {
						return err
					}
					url := c.ProjectsBase + "/portal/" + portal + "/teams/" + cmd.Args().First()
					raw, err := c.Request("DELETE", url, nil)
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
					&cli.StringFlag{Name: "json", Required: true, Usage: "Team details as JSON"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					portal, err := requirePortal(cmd)
					if err != nil {
						return err
					}
					var body any
					json.Unmarshal([]byte(cmd.String("json")), &body)
					url := base(c, portal, cmd.String("project")) + "/teams"
					raw, err := c.Request("POST", url, &zohttp.RequestOpts{JSON: body})
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
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					portal, err := requirePortal(cmd)
					if err != nil {
						return err
					}
					url := base(c, portal, cmd.String("project")) + "/teams/" + cmd.Args().First()
					raw, err := c.Request("DELETE", url, nil)
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
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					portal, err := requirePortal(cmd)
					if err != nil {
						return err
					}
					url := c.ProjectsBase + "/portal/" + portal + "/profiles"
					raw, err := c.Request("GET", url, nil)
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
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					portal, err := requirePortal(cmd)
					if err != nil {
						return err
					}
					url := c.ProjectsBase + "/portal/" + portal + "/profiles/" + cmd.Args().First()
					raw, err := c.Request("GET", url, nil)
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
					&cli.StringFlag{Name: "json", Required: true, Usage: "Profile details as JSON"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					portal, err := requirePortal(cmd)
					if err != nil {
						return err
					}
					var body any
					json.Unmarshal([]byte(cmd.String("json")), &body)
					url := c.ProjectsBase + "/portal/" + portal + "/profiles"
					raw, err := c.Request("POST", url, &zohttp.RequestOpts{JSON: body})
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
					&cli.StringFlag{Name: "json", Required: true, Usage: "Profile details as JSON"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					portal, err := requirePortal(cmd)
					if err != nil {
						return err
					}
					var body any
					json.Unmarshal([]byte(cmd.String("json")), &body)
					url := c.ProjectsBase + "/portal/" + portal + "/profiles/" + cmd.Args().First()
					raw, err := c.Request("PATCH", url, &zohttp.RequestOpts{JSON: body})
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
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					portal, err := requirePortal(cmd)
					if err != nil {
						return err
					}
					url := c.ProjectsBase + "/portal/" + portal + "/profiles/" + cmd.Args().First() + "/default"
					raw, err := c.Request("POST", url, nil)
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
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					portal, err := requirePortal(cmd)
					if err != nil {
						return err
					}
					url := c.ProjectsBase + "/portal/" + portal + "/profiles/" + cmd.Args().First()
					raw, err := c.Request("DELETE", url, nil)
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
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					portal, err := requirePortal(cmd)
					if err != nil {
						return err
					}
					url := c.ProjectsBase + "/portal/" + portal + "/roles"
					raw, err := c.Request("GET", url, nil)
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
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					portal, err := requirePortal(cmd)
					if err != nil {
						return err
					}
					url := c.ProjectsBase + "/portal/" + portal + "/roles/" + cmd.Args().First()
					raw, err := c.Request("GET", url, nil)
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
					&cli.StringFlag{Name: "json", Required: true, Usage: "Role details as JSON"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					portal, err := requirePortal(cmd)
					if err != nil {
						return err
					}
					var body any
					json.Unmarshal([]byte(cmd.String("json")), &body)
					url := c.ProjectsBase + "/portal/" + portal + "/roles"
					raw, err := c.Request("POST", url, &zohttp.RequestOpts{JSON: body})
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
					&cli.StringFlag{Name: "json", Required: true, Usage: "Role details as JSON"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					portal, err := requirePortal(cmd)
					if err != nil {
						return err
					}
					var body any
					json.Unmarshal([]byte(cmd.String("json")), &body)
					url := c.ProjectsBase + "/portal/" + portal + "/roles/" + cmd.Args().First()
					raw, err := c.Request("PATCH", url, &zohttp.RequestOpts{JSON: body})
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
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					portal, err := requirePortal(cmd)
					if err != nil {
						return err
					}
					url := c.ProjectsBase + "/portal/" + portal + "/roles/" + cmd.Args().First() + "/default"
					raw, err := c.Request("POST", url, nil)
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
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					portal, err := requirePortal(cmd)
					if err != nil {
						return err
					}
					url := c.ProjectsBase + "/portal/" + portal + "/roles/" + cmd.Args().First()
					raw, err := c.Request("DELETE", url, nil)
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
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					portal, err := requirePortal(cmd)
					if err != nil {
						return err
					}
					url := base(c, portal, cmd.String("project")) + "/modules/" + cmd.String("module") + "/records"
					raw, err := c.Request("GET", url, nil)
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
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					portal, err := requirePortal(cmd)
					if err != nil {
						return err
					}
					url := base(c, portal, cmd.String("project")) + "/modules/" + cmd.String("module") + "/records/" + cmd.Args().First()
					raw, err := c.Request("GET", url, nil)
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
					&cli.StringFlag{Name: "json", Required: true, Usage: "Record details as JSON"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					portal, err := requirePortal(cmd)
					if err != nil {
						return err
					}
					var body any
					json.Unmarshal([]byte(cmd.String("json")), &body)
					url := base(c, portal, cmd.String("project")) + "/modules/" + cmd.String("module") + "/records"
					raw, err := c.Request("POST", url, &zohttp.RequestOpts{JSON: body})
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
					&cli.StringFlag{Name: "json", Required: true, Usage: "Record details as JSON"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					portal, err := requirePortal(cmd)
					if err != nil {
						return err
					}
					var body any
					json.Unmarshal([]byte(cmd.String("json")), &body)
					url := base(c, portal, cmd.String("project")) + "/modules/" + cmd.String("module") + "/records/" + cmd.Args().First()
					raw, err := c.Request("PATCH", url, &zohttp.RequestOpts{JSON: body})
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
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					portal, err := requirePortal(cmd)
					if err != nil {
						return err
					}
					url := base(c, portal, cmd.String("project")) + "/modules/" + cmd.String("module") + "/records/" + cmd.Args().First() + "/trash"
					raw, err := c.Request("POST", url, nil)
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
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					portal, err := requirePortal(cmd)
					if err != nil {
						return err
					}
					url := base(c, portal, cmd.String("project")) + "/modules/" + cmd.String("module") + "/records/" + cmd.Args().First() + "/restore"
					raw, err := c.Request("POST", url, nil)
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
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					portal, err := requirePortal(cmd)
					if err != nil {
						return err
					}
					url := base(c, portal, cmd.String("project")) + "/modules/" + cmd.String("module") + "/records/" + cmd.Args().First()
					raw, err := c.Request("DELETE", url, nil)
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
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					portal, err := requirePortal(cmd)
					if err != nil {
						return err
					}
					url := c.ProjectsBase + "/portal/" + portal + "/reports/workload/meta"
					raw, err := c.Request("GET", url, nil)
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
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					portal, err := requirePortal(cmd)
					if err != nil {
						return err
					}
					url := c.ProjectsBase + "/portal/" + portal + "/reports/workload"
					raw, err := c.Request("GET", url, nil)
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:  "allocated",
				Usage: "Get allocated workload report",
				Flags: []cli.Flag{portalFlag},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					portal, err := requirePortal(cmd)
					if err != nil {
						return err
					}
					url := c.ProjectsBase + "/portal/" + portal + "/reports/workload/allocated"
					raw, err := c.Request("GET", url, nil)
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
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					portal, err := requirePortal(cmd)
					if err != nil {
						return err
					}
					url := c.ProjectsBase + "/portal/" + portal + "/dashboards"
					raw, err := c.Request("GET", url, nil)
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
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					portal, err := requirePortal(cmd)
					if err != nil {
						return err
					}
					url := c.ProjectsBase + "/portal/" + portal + "/dashboards/" + cmd.Args().First()
					raw, err := c.Request("GET", url, nil)
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
					&cli.StringFlag{Name: "json", Required: true, Usage: "Dashboard details as JSON"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					portal, err := requirePortal(cmd)
					if err != nil {
						return err
					}
					var body any
					json.Unmarshal([]byte(cmd.String("json")), &body)
					url := c.ProjectsBase + "/portal/" + portal + "/dashboards"
					raw, err := c.Request("POST", url, &zohttp.RequestOpts{JSON: body})
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
					&cli.StringFlag{Name: "json", Required: true, Usage: "Dashboard details as JSON"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					portal, err := requirePortal(cmd)
					if err != nil {
						return err
					}
					var body any
					json.Unmarshal([]byte(cmd.String("json")), &body)
					url := c.ProjectsBase + "/portal/" + portal + "/dashboards/" + cmd.Args().First()
					raw, err := c.Request("PATCH", url, &zohttp.RequestOpts{JSON: body})
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
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					portal, err := requirePortal(cmd)
					if err != nil {
						return err
					}
					url := c.ProjectsBase + "/portal/" + portal + "/dashboards/" + cmd.Args().First()
					raw, err := c.Request("DELETE", url, nil)
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
		},
	}
}
