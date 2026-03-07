package sprints

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/omin8tor/zoho-cli/internal"
	"github.com/omin8tor/zoho-cli/internal/auth"
	zohttp "github.com/omin8tor/zoho-cli/internal/http"
	"github.com/omin8tor/zoho-cli/internal/output"
	"github.com/urfave/cli/v3"
)

func getClient() (*zohttp.Client, error) {
	config, err := auth.ResolveAuth()
	if err != nil {
		return nil, err
	}
	return zohttp.NewClient(config)
}

func jsonToForm(body map[string]any) map[string]string {
	form := map[string]string{}
	for k, v := range body {
		switch val := v.(type) {
		case string:
			form[k] = val
		default:
			b, err := json.Marshal(val)
			if err != nil {
				form[k] = fmt.Sprintf("%v", val)
			} else {
				form[k] = string(b)
			}
		}
	}
	return form
}

func resolveTeamID(cmd *cli.Command) (string, error) {
	team := cmd.String("team")
	if team == "" {
		team = os.Getenv("ZOHO_SPRINTS_TEAM_ID")
	}
	if team == "" {
		return "", internal.NewValidationError("--team flag or ZOHO_SPRINTS_TEAM_ID env var required")
	}
	return team, nil
}

func teamBase(c *zohttp.Client, teamID string) string {
	return c.SprintsBase + "/team/" + teamID
}

func Commands() *cli.Command {
	return &cli.Command{
		Name:  "sprints",
		Usage: "Zoho Sprints operations",
		Flags: []cli.Flag{
			&cli.StringFlag{Name: "team", Usage: "Team ID (or set ZOHO_SPRINTS_TEAM_ID)"},
		},
		Commands: []*cli.Command{
			teamsCmd(),
			projectsCmd(),
			sprintsCmd(),
			itemsCmd(),
			epicsCmd(),
			statusesCmd(),
			itemTypesCmd(),
			prioritiesCmd(),
			membersCmd(),
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
				Usage: "List teams",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					raw, err := c.Request("GET", c.SprintsBase+"/teams/", nil)
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
		},
	}
}

func projectsCmd() *cli.Command {
	return &cli.Command{
		Name:  "projects",
		Usage: "Project operations",
		Commands: []*cli.Command{
			{
				Name:  "list",
				Usage: "List projects",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "index", Usage: "Start index", Value: "1"},
					&cli.StringFlag{Name: "range", Usage: "Number of records", Value: "100"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					teamID, err := resolveTeamID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("GET", teamBase(c, teamID)+"/projects/", &zohttp.RequestOpts{
						Params: map[string]string{
							"action": "data",
							"index":  cmd.String("index"),
							"range":  cmd.String("range"),
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
				Usage:     "Get a project",
				ArgsUsage: "<project-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					if cmd.Args().Len() < 1 {
						return internal.NewValidationError("project ID required")
					}
					c, err := getClient()
					if err != nil {
						return err
					}
					teamID, err := resolveTeamID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("GET", teamBase(c, teamID)+"/projects/"+cmd.Args().First()+"/", nil)
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
					&cli.StringFlag{Name: "json", Required: true, Usage: "JSON body (name, description, template, etc.)"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					teamID, err := resolveTeamID(cmd)
					if err != nil {
						return err
					}
					var body map[string]any
					if err := json.Unmarshal([]byte(cmd.String("json")), &body); err != nil {
						return internal.NewValidationError(fmt.Sprintf("invalid JSON: %v", err))
					}
					form := jsonToForm(body)
					raw, err := c.Request("POST", teamBase(c, teamID)+"/projects/", &zohttp.RequestOpts{
						Form: form,
					})
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
					&cli.StringFlag{Name: "json", Required: true, Usage: "JSON body"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					if cmd.Args().Len() < 1 {
						return internal.NewValidationError("project ID required")
					}
					c, err := getClient()
					if err != nil {
						return err
					}
					teamID, err := resolveTeamID(cmd)
					if err != nil {
						return err
					}
					var body map[string]any
					if err := json.Unmarshal([]byte(cmd.String("json")), &body); err != nil {
						return internal.NewValidationError(fmt.Sprintf("invalid JSON: %v", err))
					}
					form := jsonToForm(body)
					raw, err := c.Request("POST", teamBase(c, teamID)+"/projects/"+cmd.Args().First()+"/", &zohttp.RequestOpts{
						Form: form,
					})
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
				Action: func(_ context.Context, cmd *cli.Command) error {
					if cmd.Args().Len() < 1 {
						return internal.NewValidationError("project ID required")
					}
					c, err := getClient()
					if err != nil {
						return err
					}
					teamID, err := resolveTeamID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("DELETE", teamBase(c, teamID)+"/projects/"+cmd.Args().First()+"/", nil)
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
		},
	}
}

func sprintsCmd() *cli.Command {
	return &cli.Command{
		Name:  "sprints",
		Usage: "Sprint operations",
		Commands: []*cli.Command{
			{
				Name:      "list",
				Usage:     "List sprints in a project",
				ArgsUsage: "<project-id>",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "type", Usage: "Sprint type filter: 1=upcoming, 2=active, 3=completed, 4=canceled"},
					&cli.StringFlag{Name: "index", Usage: "Start index", Value: "1"},
					&cli.StringFlag{Name: "range", Usage: "Number of records", Value: "100"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					if cmd.Args().Len() < 1 {
						return internal.NewValidationError("project ID required")
					}
					c, err := getClient()
					if err != nil {
						return err
					}
					teamID, err := resolveTeamID(cmd)
					if err != nil {
						return err
					}
					params := map[string]string{
						"action": "data",
						"index":  cmd.String("index"),
						"range":  cmd.String("range"),
					}
					if v := cmd.String("type"); v != "" {
						params["type"] = "[" + v + "]"
					}
					raw, err := c.Request("GET", teamBase(c, teamID)+"/projects/"+cmd.Args().First()+"/sprints/", &zohttp.RequestOpts{
						Params: params,
					})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "get",
				Usage:     "Get a sprint",
				ArgsUsage: "<project-id> <sprint-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					if cmd.Args().Len() < 2 {
						return internal.NewValidationError("project ID and sprint ID required")
					}
					c, err := getClient()
					if err != nil {
						return err
					}
					teamID, err := resolveTeamID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("GET", teamBase(c, teamID)+"/projects/"+cmd.Args().Get(0)+"/sprints/"+cmd.Args().Get(1)+"/", nil)
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "create",
				Usage:     "Create a sprint",
				ArgsUsage: "<project-id>",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "json", Required: true, Usage: "JSON body (name, startdate, enddate, etc.)"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					if cmd.Args().Len() < 1 {
						return internal.NewValidationError("project ID required")
					}
					c, err := getClient()
					if err != nil {
						return err
					}
					teamID, err := resolveTeamID(cmd)
					if err != nil {
						return err
					}
					var body map[string]any
					if err := json.Unmarshal([]byte(cmd.String("json")), &body); err != nil {
						return internal.NewValidationError(fmt.Sprintf("invalid JSON: %v", err))
					}
					form := jsonToForm(body)
					raw, err := c.Request("POST", teamBase(c, teamID)+"/projects/"+cmd.Args().First()+"/sprints/", &zohttp.RequestOpts{
						Form: form,
					})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "update",
				Usage:     "Update a sprint",
				ArgsUsage: "<project-id> <sprint-id>",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "json", Required: true, Usage: "JSON body"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					if cmd.Args().Len() < 2 {
						return internal.NewValidationError("project ID and sprint ID required")
					}
					c, err := getClient()
					if err != nil {
						return err
					}
					teamID, err := resolveTeamID(cmd)
					if err != nil {
						return err
					}
					var body map[string]any
					if err := json.Unmarshal([]byte(cmd.String("json")), &body); err != nil {
						return internal.NewValidationError(fmt.Sprintf("invalid JSON: %v", err))
					}
					form := jsonToForm(body)
					raw, err := c.Request("POST", teamBase(c, teamID)+"/projects/"+cmd.Args().Get(0)+"/sprints/"+cmd.Args().Get(1)+"/", &zohttp.RequestOpts{
						Form: form,
					})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "delete",
				Usage:     "Delete a sprint",
				ArgsUsage: "<project-id> <sprint-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					if cmd.Args().Len() < 2 {
						return internal.NewValidationError("project ID and sprint ID required")
					}
					c, err := getClient()
					if err != nil {
						return err
					}
					teamID, err := resolveTeamID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("DELETE", teamBase(c, teamID)+"/projects/"+cmd.Args().Get(0)+"/sprints/"+cmd.Args().Get(1)+"/", nil)
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
		},
	}
}

func itemsCmd() *cli.Command {
	return &cli.Command{
		Name:  "items",
		Usage: "Work item operations",
		Commands: []*cli.Command{
			{
				Name:      "list",
				Usage:     "List items in a sprint or backlog",
				ArgsUsage: "<project-id> <sprint-id|backlog-id>",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "index", Usage: "Start index", Value: "1"},
					&cli.StringFlag{Name: "range", Usage: "Number of records", Value: "100"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					if cmd.Args().Len() < 2 {
						return internal.NewValidationError("project ID and sprint/backlog ID required")
					}
					c, err := getClient()
					if err != nil {
						return err
					}
					teamID, err := resolveTeamID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("GET", teamBase(c, teamID)+"/projects/"+cmd.Args().Get(0)+"/sprints/"+cmd.Args().Get(1)+"/item/", &zohttp.RequestOpts{
						Params: map[string]string{
							"action": "data",
							"index":  cmd.String("index"),
							"range":  cmd.String("range"),
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
				Usage:     "Get a specific item",
				ArgsUsage: "<project-id> <sprint-id> <item-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					if cmd.Args().Len() < 3 {
						return internal.NewValidationError("project ID, sprint ID, and item ID required")
					}
					c, err := getClient()
					if err != nil {
						return err
					}
					teamID, err := resolveTeamID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("GET", teamBase(c, teamID)+"/projects/"+cmd.Args().Get(0)+"/sprints/"+cmd.Args().Get(1)+"/item/"+cmd.Args().Get(2)+"/", nil)
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "create",
				Usage:     "Create an item in a sprint or backlog",
				ArgsUsage: "<project-id> <sprint-id|backlog-id>",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "json", Required: true, Usage: "JSON body (name, projitemtypeid, etc.)"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					if cmd.Args().Len() < 2 {
						return internal.NewValidationError("project ID and sprint/backlog ID required")
					}
					c, err := getClient()
					if err != nil {
						return err
					}
					teamID, err := resolveTeamID(cmd)
					if err != nil {
						return err
					}
					var body map[string]any
					if err := json.Unmarshal([]byte(cmd.String("json")), &body); err != nil {
						return internal.NewValidationError(fmt.Sprintf("invalid JSON: %v", err))
					}
					form := jsonToForm(body)
					raw, err := c.Request("POST", teamBase(c, teamID)+"/projects/"+cmd.Args().Get(0)+"/sprints/"+cmd.Args().Get(1)+"/item/", &zohttp.RequestOpts{
						Form: form,
					})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "update",
				Usage:     "Update an item",
				ArgsUsage: "<project-id> <sprint-id> <item-id>",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "json", Required: true, Usage: "JSON body"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					if cmd.Args().Len() < 3 {
						return internal.NewValidationError("project ID, sprint ID, and item ID required")
					}
					c, err := getClient()
					if err != nil {
						return err
					}
					teamID, err := resolveTeamID(cmd)
					if err != nil {
						return err
					}
					var body map[string]any
					if err := json.Unmarshal([]byte(cmd.String("json")), &body); err != nil {
						return internal.NewValidationError(fmt.Sprintf("invalid JSON: %v", err))
					}
					form := jsonToForm(body)
					raw, err := c.Request("POST", teamBase(c, teamID)+"/projects/"+cmd.Args().Get(0)+"/sprints/"+cmd.Args().Get(1)+"/item/"+cmd.Args().Get(2)+"/", &zohttp.RequestOpts{
						Form: form,
					})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "delete",
				Usage:     "Delete an item",
				ArgsUsage: "<project-id> <sprint-id> <item-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					if cmd.Args().Len() < 3 {
						return internal.NewValidationError("project ID, sprint ID, and item ID required")
					}
					c, err := getClient()
					if err != nil {
						return err
					}
					teamID, err := resolveTeamID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("DELETE", teamBase(c, teamID)+"/projects/"+cmd.Args().Get(0)+"/sprints/"+cmd.Args().Get(1)+"/item/"+cmd.Args().Get(2)+"/", nil)
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
		},
	}
}

func epicsCmd() *cli.Command {
	return &cli.Command{
		Name:  "epics",
		Usage: "Epic operations",
		Commands: []*cli.Command{
			{
				Name:      "list",
				Usage:     "List epics in a project",
				ArgsUsage: "<project-id>",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "index", Usage: "Start index", Value: "1"},
					&cli.StringFlag{Name: "range", Usage: "Number of records", Value: "100"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					if cmd.Args().Len() < 1 {
						return internal.NewValidationError("project ID required")
					}
					c, err := getClient()
					if err != nil {
						return err
					}
					teamID, err := resolveTeamID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("GET", teamBase(c, teamID)+"/projects/"+cmd.Args().First()+"/epic/", &zohttp.RequestOpts{
						Params: map[string]string{
							"action": "data",
							"index":  cmd.String("index"),
							"range":  cmd.String("range"),
						},
					})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "create",
				Usage:     "Create an epic",
				ArgsUsage: "<project-id>",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "json", Required: true, Usage: "JSON body (name, description, etc.)"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					if cmd.Args().Len() < 1 {
						return internal.NewValidationError("project ID required")
					}
					c, err := getClient()
					if err != nil {
						return err
					}
					teamID, err := resolveTeamID(cmd)
					if err != nil {
						return err
					}
					var body map[string]any
					if err := json.Unmarshal([]byte(cmd.String("json")), &body); err != nil {
						return internal.NewValidationError(fmt.Sprintf("invalid JSON: %v", err))
					}
					form := jsonToForm(body)
					raw, err := c.Request("POST", teamBase(c, teamID)+"/projects/"+cmd.Args().First()+"/epic/", &zohttp.RequestOpts{
						Form: form,
					})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "update",
				Usage:     "Update an epic",
				ArgsUsage: "<project-id> <epic-id>",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "json", Required: true, Usage: "JSON body"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					if cmd.Args().Len() < 2 {
						return internal.NewValidationError("project ID and epic ID required")
					}
					c, err := getClient()
					if err != nil {
						return err
					}
					teamID, err := resolveTeamID(cmd)
					if err != nil {
						return err
					}
					var body map[string]any
					if err := json.Unmarshal([]byte(cmd.String("json")), &body); err != nil {
						return internal.NewValidationError(fmt.Sprintf("invalid JSON: %v", err))
					}
					form := jsonToForm(body)
					raw, err := c.Request("POST", teamBase(c, teamID)+"/projects/"+cmd.Args().Get(0)+"/epic/"+cmd.Args().Get(1)+"/", &zohttp.RequestOpts{
						Form: form,
					})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "delete",
				Usage:     "Delete an epic",
				ArgsUsage: "<project-id> <epic-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					if cmd.Args().Len() < 2 {
						return internal.NewValidationError("project ID and epic ID required")
					}
					c, err := getClient()
					if err != nil {
						return err
					}
					teamID, err := resolveTeamID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("DELETE", teamBase(c, teamID)+"/projects/"+cmd.Args().Get(0)+"/epic/"+cmd.Args().Get(1)+"/", nil)
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
		},
	}
}

func statusesCmd() *cli.Command {
	return &cli.Command{
		Name:  "statuses",
		Usage: "Item status operations",
		Commands: []*cli.Command{
			{
				Name:      "list",
				Usage:     "List item statuses in a project",
				ArgsUsage: "<project-id>",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "index", Usage: "Start index", Value: "1"},
					&cli.StringFlag{Name: "range", Usage: "Number of records", Value: "50"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					if cmd.Args().Len() < 1 {
						return internal.NewValidationError("project ID required")
					}
					c, err := getClient()
					if err != nil {
						return err
					}
					teamID, err := resolveTeamID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("GET", teamBase(c, teamID)+"/projects/"+cmd.Args().First()+"/itemstatus/", &zohttp.RequestOpts{
						Params: map[string]string{
							"action": "data",
							"index":  cmd.String("index"),
							"range":  cmd.String("range"),
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

func itemTypesCmd() *cli.Command {
	return &cli.Command{
		Name:  "item-types",
		Usage: "Item type operations",
		Commands: []*cli.Command{
			{
				Name:      "list",
				Usage:     "List item types in a project",
				ArgsUsage: "<project-id>",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "index", Usage: "Start index", Value: "1"},
					&cli.StringFlag{Name: "range", Usage: "Number of records", Value: "50"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					if cmd.Args().Len() < 1 {
						return internal.NewValidationError("project ID required")
					}
					c, err := getClient()
					if err != nil {
						return err
					}
					teamID, err := resolveTeamID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("GET", teamBase(c, teamID)+"/projects/"+cmd.Args().First()+"/itemtype/", &zohttp.RequestOpts{
						Params: map[string]string{
							"action": "data",
							"index":  cmd.String("index"),
							"range":  cmd.String("range"),
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

func prioritiesCmd() *cli.Command {
	return &cli.Command{
		Name:  "priorities",
		Usage: "Priority operations",
		Commands: []*cli.Command{
			{
				Name:      "list",
				Usage:     "List priorities in a project",
				ArgsUsage: "<project-id>",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "index", Usage: "Start index", Value: "1"},
					&cli.StringFlag{Name: "range", Usage: "Number of records", Value: "50"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					if cmd.Args().Len() < 1 {
						return internal.NewValidationError("project ID required")
					}
					c, err := getClient()
					if err != nil {
						return err
					}
					teamID, err := resolveTeamID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("GET", teamBase(c, teamID)+"/projects/"+cmd.Args().First()+"/priority/", &zohttp.RequestOpts{
						Params: map[string]string{
							"action": "data",
							"index":  cmd.String("index"),
							"range":  cmd.String("range"),
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

func membersCmd() *cli.Command {
	return &cli.Command{
		Name:  "members",
		Usage: "Project member operations",
		Commands: []*cli.Command{
			{
				Name:      "list",
				Usage:     "List members of a project",
				ArgsUsage: "<project-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					if cmd.Args().Len() < 1 {
						return internal.NewValidationError("project ID required")
					}
					c, err := getClient()
					if err != nil {
						return err
					}
					teamID, err := resolveTeamID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("GET", teamBase(c, teamID)+"/projects/"+cmd.Args().First()+"/members/", nil)
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
		},
	}
}
