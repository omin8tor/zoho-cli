package projects

import (
	"context"
	"github.com/omin8tor/zoho-cli/internal"
	zohttp "github.com/omin8tor/zoho-cli/internal/http"
	"github.com/omin8tor/zoho-cli/internal/output"
	"github.com/omin8tor/zoho-cli/internal/pagination"
	"github.com/urfave/cli/v3"
	"time"
)

func convertDate(s string) string {
	if t, err := time.Parse("2006-01-02", s); err == nil {
		return t.Format("01-02-2006")
	}
	return s
}

func base(c *zohttp.Client, portal, project string) string {
	return c.ProjectsBase + "/portal/" + portal + "/projects/" + project
}

var portalFlag = &cli.StringFlag{Name: "portal", Usage: "Portal ID", Sources: cli.EnvVars("ZOHO_PORTAL_ID")}
var projectFlag = &cli.StringFlag{Name: "project", Required: true, Usage: "Project ID"}

func requirePortal(cmd *cli.Command) (string, error) {
	return internal.RequireFlag(cmd, "portal", "ZOHO_PORTAL_ID")
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
					c, err := zohttp.GetClient()
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
					c, err := zohttp.GetClient()
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
					c, err := zohttp.GetClient()
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
					&cli.StringFlag{Name: "name", Usage: "Project name"},
					&cli.StringFlag{Name: "description", Usage: "Project description"},
					&cli.StringFlag{Name: "status", Usage: "Project status"},
					&cli.StringFlag{Name: "json", Usage: "Additional fields as JSON"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
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
					if cmd.IsSet("status") {
						body["status"] = cmd.String("status")
					}
					if err := internal.MergeJSON(cmd, body); err != nil {
						return err
					}
					url := c.ProjectsBase + "/portal/" + portal + "/projects/" + cmd.Args().First()
					raw, err := c.Request("PATCH", url, &zohttp.RequestOpts{JSON: body})
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
					c, err := zohttp.GetClient()
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
					c, err := zohttp.GetClient()
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
					c, err := zohttp.GetClient()
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
