package desk

import (
	"context"
	"encoding/json"
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

func resolveOrgID(cmd *cli.Command) (string, error) {
	org := cmd.String("org")
	if org == "" {
		org = os.Getenv("ZOHO_DESK_ORG_ID")
	}
	if org == "" {
		return "", internal.NewValidationError("--org flag or ZOHO_DESK_ORG_ID env var required")
	}
	return org, nil
}

func orgHeaders(orgID string) map[string]string {
	return map[string]string{"orgId": orgID}
}

func Commands() *cli.Command {
	return &cli.Command{
		Name:  "desk",
		Usage: "Zoho Desk operations",
		Flags: []cli.Flag{
			&cli.StringFlag{Name: "org", Usage: "Organization ID (or set ZOHO_DESK_ORG_ID)"},
		},
		Commands: []*cli.Command{
			ticketsCmd(),
			contactsCmd(),
			accountsCmd(),
			agentsCmd(),
			departmentsCmd(),
			searchCmd(),
		},
	}
}

func departmentsCmd() *cli.Command {
	return &cli.Command{
		Name:  "departments",
		Usage: "Department operations",
		Commands: []*cli.Command{
			{
				Name:  "list",
				Usage: "List departments",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "from", Usage: "Starting index"},
					&cli.StringFlag{Name: "limit", Usage: "Max records"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					params := map[string]string{}
					if v := cmd.String("from"); v != "" {
						params["from"] = v
					}
					if v := cmd.String("limit"); v != "" {
						params["limit"] = v
					}
					raw, err := c.Request("GET", c.DeskBase+"/departments", &zohttp.RequestOpts{
						Params:  params,
						Headers: orgHeaders(orgID),
					})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "get",
				Usage:     "Get a department",
				ArgsUsage: "<department-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("GET", c.DeskBase+"/departments/"+cmd.Args().First(), &zohttp.RequestOpts{
						Headers: orgHeaders(orgID),
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

func agentsCmd() *cli.Command {
	return &cli.Command{
		Name:  "agents",
		Usage: "Agent operations",
		Commands: []*cli.Command{
			{
				Name:  "list",
				Usage: "List agents",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "from", Usage: "Starting index"},
					&cli.StringFlag{Name: "limit", Usage: "Max records"},
					&cli.StringFlag{Name: "department-id", Usage: "Department ID"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					params := map[string]string{}
					if v := cmd.String("from"); v != "" {
						params["from"] = v
					}
					if v := cmd.String("limit"); v != "" {
						params["limit"] = v
					}
					if v := cmd.String("department-id"); v != "" {
						params["departmentId"] = v
					}
					raw, err := c.Request("GET", c.DeskBase+"/agents", &zohttp.RequestOpts{
						Params:  params,
						Headers: orgHeaders(orgID),
					})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "get",
				Usage:     "Get an agent",
				ArgsUsage: "<agent-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("GET", c.DeskBase+"/agents/"+cmd.Args().First(), &zohttp.RequestOpts{
						Headers: orgHeaders(orgID),
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

func ticketsCmd() *cli.Command {
	return &cli.Command{
		Name:  "tickets",
		Usage: "Ticket operations",
		Commands: []*cli.Command{
			{
				Name:  "list",
				Usage: "List tickets",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "from", Usage: "Starting index"},
					&cli.StringFlag{Name: "limit", Usage: "Max records"},
					&cli.StringFlag{Name: "department-ids", Usage: "Comma-separated department IDs"},
					&cli.StringFlag{Name: "status", Usage: "Filter by status"},
					&cli.StringFlag{Name: "priority", Usage: "Filter by priority"},
					&cli.StringFlag{Name: "assignee", Usage: "Assignee ID"},
					&cli.StringFlag{Name: "sort-by", Usage: "Sort field"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					params := map[string]string{}
					if v := cmd.String("from"); v != "" {
						params["from"] = v
					}
					if v := cmd.String("limit"); v != "" {
						params["limit"] = v
					}
					if v := cmd.String("department-ids"); v != "" {
						params["departmentIds"] = v
					}
					if v := cmd.String("status"); v != "" {
						params["status"] = v
					}
					if v := cmd.String("priority"); v != "" {
						params["priority"] = v
					}
					if v := cmd.String("assignee"); v != "" {
						params["assignee"] = v
					}
					if v := cmd.String("sort-by"); v != "" {
						params["sortBy"] = v
					}
					raw, err := c.Request("GET", c.DeskBase+"/tickets", &zohttp.RequestOpts{
						Params:  params,
						Headers: orgHeaders(orgID),
					})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "get",
				Usage:     "Get a ticket",
				ArgsUsage: "<ticket-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("GET", c.DeskBase+"/tickets/"+cmd.Args().First(), &zohttp.RequestOpts{
						Headers: orgHeaders(orgID),
					})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:  "create",
				Usage: "Create a ticket",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "json", Required: true, Usage: "JSON body"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					var body any
					json.Unmarshal([]byte(cmd.String("json")), &body)
					raw, err := c.Request("POST", c.DeskBase+"/tickets", &zohttp.RequestOpts{
						JSON:    body,
						Headers: orgHeaders(orgID),
					})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "update",
				Usage:     "Update a ticket",
				ArgsUsage: "<ticket-id>",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "json", Required: true, Usage: "JSON body"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					var body any
					json.Unmarshal([]byte(cmd.String("json")), &body)
					raw, err := c.Request("PATCH", c.DeskBase+"/tickets/"+cmd.Args().First(), &zohttp.RequestOpts{
						JSON:    body,
						Headers: orgHeaders(orgID),
					})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "delete",
				Usage:     "Delete a ticket",
				ArgsUsage: "<ticket-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("DELETE", c.DeskBase+"/tickets/"+cmd.Args().First(), &zohttp.RequestOpts{
						Headers: orgHeaders(orgID),
					})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:  "search",
				Usage: "Search tickets",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "subject", Usage: "Subject text"},
					&cli.StringFlag{Name: "status", Usage: "Ticket status"},
					&cli.StringFlag{Name: "priority", Usage: "Ticket priority"},
					&cli.StringFlag{Name: "created-time", Usage: "Created time filter"},
					&cli.StringFlag{Name: "from", Usage: "Starting index"},
					&cli.StringFlag{Name: "limit", Usage: "Max records"},
					&cli.StringFlag{Name: "sort-by", Usage: "Sort field"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					params := map[string]string{}
					if v := cmd.String("subject"); v != "" {
						params["subject"] = v
					}
					if v := cmd.String("status"); v != "" {
						params["status"] = v
					}
					if v := cmd.String("priority"); v != "" {
						params["priority"] = v
					}
					if v := cmd.String("created-time"); v != "" {
						params["createdTime"] = v
					}
					if v := cmd.String("from"); v != "" {
						params["from"] = v
					}
					if v := cmd.String("limit"); v != "" {
						params["limit"] = v
					}
					if v := cmd.String("sort-by"); v != "" {
						params["sortBy"] = v
					}
					raw, err := c.Request("GET", c.DeskBase+"/tickets/search", &zohttp.RequestOpts{
						Params:  params,
						Headers: orgHeaders(orgID),
					})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "threads",
				Usage:     "List ticket threads",
				ArgsUsage: "<ticket-id>",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "from", Usage: "Starting index"},
					&cli.StringFlag{Name: "limit", Usage: "Max records"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					params := map[string]string{}
					if v := cmd.String("from"); v != "" {
						params["from"] = v
					}
					if v := cmd.String("limit"); v != "" {
						params["limit"] = v
					}
					raw, err := c.Request("GET", c.DeskBase+"/tickets/"+cmd.Args().First()+"/threads", &zohttp.RequestOpts{
						Params:  params,
						Headers: orgHeaders(orgID),
					})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "reply",
				Usage:     "Send a ticket reply",
				ArgsUsage: "<ticket-id>",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "json", Required: true, Usage: "JSON body"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					var body any
					json.Unmarshal([]byte(cmd.String("json")), &body)
					raw, err := c.Request("POST", c.DeskBase+"/tickets/"+cmd.Args().First()+"/sendReply", &zohttp.RequestOpts{
						JSON:    body,
						Headers: orgHeaders(orgID),
					})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "comments",
				Usage:     "List ticket comments",
				ArgsUsage: "<ticket-id>",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "from", Usage: "Starting index"},
					&cli.StringFlag{Name: "limit", Usage: "Max records"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					params := map[string]string{}
					if v := cmd.String("from"); v != "" {
						params["from"] = v
					}
					if v := cmd.String("limit"); v != "" {
						params["limit"] = v
					}
					raw, err := c.Request("GET", c.DeskBase+"/tickets/"+cmd.Args().First()+"/comments", &zohttp.RequestOpts{
						Params:  params,
						Headers: orgHeaders(orgID),
					})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "add-comment",
				Usage:     "Add a ticket comment",
				ArgsUsage: "<ticket-id>",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "json", Required: true, Usage: "JSON body"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					var body any
					json.Unmarshal([]byte(cmd.String("json")), &body)
					raw, err := c.Request("POST", c.DeskBase+"/tickets/"+cmd.Args().First()+"/comments", &zohttp.RequestOpts{
						JSON:    body,
						Headers: orgHeaders(orgID),
					})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "attachments",
				Usage:     "List ticket attachments",
				ArgsUsage: "<ticket-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("GET", c.DeskBase+"/tickets/"+cmd.Args().First()+"/attachments", &zohttp.RequestOpts{
						Headers: orgHeaders(orgID),
					})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "history",
				Usage:     "List ticket history",
				ArgsUsage: "<ticket-id>",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "from", Usage: "Starting index"},
					&cli.StringFlag{Name: "limit", Usage: "Max records"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					params := map[string]string{}
					if v := cmd.String("from"); v != "" {
						params["from"] = v
					}
					if v := cmd.String("limit"); v != "" {
						params["limit"] = v
					}
					raw, err := c.Request("GET", c.DeskBase+"/tickets/"+cmd.Args().First()+"/history", &zohttp.RequestOpts{
						Params:  params,
						Headers: orgHeaders(orgID),
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

func contactsCmd() *cli.Command {
	return &cli.Command{
		Name:  "contacts",
		Usage: "Contact operations",
		Commands: []*cli.Command{
			{
				Name:  "list",
				Usage: "List contacts",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "from", Usage: "Starting index"},
					&cli.StringFlag{Name: "limit", Usage: "Max records"},
					&cli.StringFlag{Name: "sort-by", Usage: "Sort field"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					params := map[string]string{}
					if v := cmd.String("from"); v != "" {
						params["from"] = v
					}
					if v := cmd.String("limit"); v != "" {
						params["limit"] = v
					}
					if v := cmd.String("sort-by"); v != "" {
						params["sortBy"] = v
					}
					raw, err := c.Request("GET", c.DeskBase+"/contacts", &zohttp.RequestOpts{
						Params:  params,
						Headers: orgHeaders(orgID),
					})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "get",
				Usage:     "Get a contact",
				ArgsUsage: "<contact-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("GET", c.DeskBase+"/contacts/"+cmd.Args().First(), &zohttp.RequestOpts{
						Headers: orgHeaders(orgID),
					})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:  "create",
				Usage: "Create a contact",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "json", Required: true, Usage: "JSON body"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					var body any
					json.Unmarshal([]byte(cmd.String("json")), &body)
					raw, err := c.Request("POST", c.DeskBase+"/contacts", &zohttp.RequestOpts{
						JSON:    body,
						Headers: orgHeaders(orgID),
					})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "update",
				Usage:     "Update a contact",
				ArgsUsage: "<contact-id>",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "json", Required: true, Usage: "JSON body"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					var body any
					json.Unmarshal([]byte(cmd.String("json")), &body)
					raw, err := c.Request("PATCH", c.DeskBase+"/contacts/"+cmd.Args().First(), &zohttp.RequestOpts{
						JSON:    body,
						Headers: orgHeaders(orgID),
					})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "delete",
				Usage:     "Delete a contact",
				ArgsUsage: "<contact-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("DELETE", c.DeskBase+"/contacts/"+cmd.Args().First(), &zohttp.RequestOpts{
						Headers: orgHeaders(orgID),
					})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:  "search",
				Usage: "Search contacts",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "email", Usage: "Email address"},
					&cli.StringFlag{Name: "phone", Usage: "Phone number"},
					&cli.StringFlag{Name: "from", Usage: "Starting index"},
					&cli.StringFlag{Name: "limit", Usage: "Max records"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					params := map[string]string{}
					if v := cmd.String("email"); v != "" {
						params["email"] = v
					}
					if v := cmd.String("phone"); v != "" {
						params["phone"] = v
					}
					if v := cmd.String("from"); v != "" {
						params["from"] = v
					}
					if v := cmd.String("limit"); v != "" {
						params["limit"] = v
					}
					raw, err := c.Request("GET", c.DeskBase+"/contacts/search", &zohttp.RequestOpts{
						Params:  params,
						Headers: orgHeaders(orgID),
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

func accountsCmd() *cli.Command {
	return &cli.Command{
		Name:  "accounts",
		Usage: "Account operations",
		Commands: []*cli.Command{
			{
				Name:  "list",
				Usage: "List accounts",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "from", Usage: "Starting index"},
					&cli.StringFlag{Name: "limit", Usage: "Max records"},
					&cli.StringFlag{Name: "sort-by", Usage: "Sort field"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					params := map[string]string{}
					if v := cmd.String("from"); v != "" {
						params["from"] = v
					}
					if v := cmd.String("limit"); v != "" {
						params["limit"] = v
					}
					if v := cmd.String("sort-by"); v != "" {
						params["sortBy"] = v
					}
					raw, err := c.Request("GET", c.DeskBase+"/accounts", &zohttp.RequestOpts{
						Params:  params,
						Headers: orgHeaders(orgID),
					})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "get",
				Usage:     "Get an account",
				ArgsUsage: "<account-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("GET", c.DeskBase+"/accounts/"+cmd.Args().First(), &zohttp.RequestOpts{
						Headers: orgHeaders(orgID),
					})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:  "create",
				Usage: "Create an account",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "json", Required: true, Usage: "JSON body"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					var body any
					json.Unmarshal([]byte(cmd.String("json")), &body)
					raw, err := c.Request("POST", c.DeskBase+"/accounts", &zohttp.RequestOpts{
						JSON:    body,
						Headers: orgHeaders(orgID),
					})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "update",
				Usage:     "Update an account",
				ArgsUsage: "<account-id>",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "json", Required: true, Usage: "JSON body"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					var body any
					json.Unmarshal([]byte(cmd.String("json")), &body)
					raw, err := c.Request("PATCH", c.DeskBase+"/accounts/"+cmd.Args().First(), &zohttp.RequestOpts{
						JSON:    body,
						Headers: orgHeaders(orgID),
					})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "delete",
				Usage:     "Delete an account",
				ArgsUsage: "<account-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("DELETE", c.DeskBase+"/accounts/"+cmd.Args().First(), &zohttp.RequestOpts{
						Headers: orgHeaders(orgID),
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

func searchCmd() *cli.Command {
	return &cli.Command{
		Name:  "search",
		Usage: "Global desk search",
		Flags: []cli.Flag{
			&cli.StringFlag{Name: "query", Required: true, Usage: "Search query"},
			&cli.StringFlag{Name: "module", Value: "tickets,contacts,accounts", Usage: "Modules to search"},
			&cli.StringFlag{Name: "from", Usage: "Starting index"},
			&cli.StringFlag{Name: "limit", Usage: "Max records"},
		},
		Action: func(_ context.Context, cmd *cli.Command) error {
			c, err := getClient()
			if err != nil {
				return err
			}
			orgID, err := resolveOrgID(cmd)
			if err != nil {
				return err
			}
			params := map[string]string{
				"searchStr": cmd.String("query"),
				"module":    cmd.String("module"),
			}
			if v := cmd.String("from"); v != "" {
				params["from"] = v
			}
			if v := cmd.String("limit"); v != "" {
				params["limit"] = v
			}
			raw, err := c.Request("GET", c.DeskBase+"/search", &zohttp.RequestOpts{
				Params:  params,
				Headers: orgHeaders(orgID),
			})
			if err != nil {
				return err
			}
			return output.JSONRaw(raw)
		},
	}
}
