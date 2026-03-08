package projects

import (
	"context"
	"github.com/omin8tor/zoho-cli/internal"
	zohttp "github.com/omin8tor/zoho-cli/internal/http"
	"github.com/omin8tor/zoho-cli/internal/output"
	"github.com/urfave/cli/v3"
)

func usersCmd() *cli.Command {
	return &cli.Command{
		Name:  "users",
		Usage: "Portal user operations",
		Commands: []*cli.Command{
			{
				Name:  "list",
				Usage: "List portal users",
				Flags: []cli.Flag{portalFlag, allFlag, limitFlag},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					portal, err := requirePortal(cmd)
					if err != nil {
						return err
					}
					url := c.ProjectsBase + "/portal/" + portal + "/users"
					return paginateProjectsList(c, cmd, url, "users", nil)
				},
			},
			{
				Name:      "get",
				Usage:     "Get a user",
				ArgsUsage: "<user-id>",
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
					&cli.StringFlag{Name: "userdetails", Usage: "User details as JSON string"},
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
					form := map[string]string{}
					if cmd.IsSet("userdetails") {
						form["userdetails"] = cmd.String("userdetails")
					}
					if err := internal.MergeJSONForm(cmd, form); err != nil {
						return err
					}
					url := c.ProjectsBase + "/portal/" + portal + "/users"
					raw, err := c.Request("POST", url, &zohttp.RequestOpts{
						Form: form,
					})
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
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					portal, err := requirePortal(cmd)
					if err != nil {
						return err
					}
					url := c.ProjectsBase + "/portal/" + portal + "/users/" + cmd.Args().First() + "/activate"
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
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					portal, err := requirePortal(cmd)
					if err != nil {
						return err
					}
					url := c.ProjectsBase + "/portal/" + portal + "/users/" + cmd.Args().First() + "/deactivate"
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
					c, err := zohttp.GetClient()
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
				Flags: []cli.Flag{portalFlag, projectFlag, allFlag, limitFlag},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					portal, err := requirePortal(cmd)
					if err != nil {
						return err
					}
					url := base(c, portal, cmd.String("project")) + "/users"
					return paginateProjectsList(c, cmd, url, "users", nil)
				},
			},
			{
				Name:      "get",
				Usage:     "Get a project user",
				ArgsUsage: "<user-id>",
				Flags:     []cli.Flag{portalFlag, projectFlag},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
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
					&cli.StringFlag{Name: "zpuid", Usage: "User ZPUID"},
					&cli.StringFlag{Name: "role", Usage: "Role in project"},
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
					if cmd.IsSet("zpuid") {
						body["zpuid"] = cmd.String("zpuid")
					}
					if cmd.IsSet("role") {
						body["role"] = cmd.String("role")
					}
					if err := internal.MergeJSON(cmd, body); err != nil {
						return err
					}
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
					&cli.StringFlag{Name: "role", Usage: "Role in project"},
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
					if cmd.IsSet("role") {
						body["role"] = cmd.String("role")
					}
					if err := internal.MergeJSON(cmd, body); err != nil {
						return err
					}
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
					c, err := zohttp.GetClient()
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
