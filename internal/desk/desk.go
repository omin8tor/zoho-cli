package desk

import (
	"context"

	"github.com/omin8tor/zoho-cli/internal"
	zohttp "github.com/omin8tor/zoho-cli/internal/http"
	"github.com/omin8tor/zoho-cli/internal/output"
	"github.com/omin8tor/zoho-cli/internal/pagination"
	"github.com/urfave/cli/v3"
)

func orgHeaders(orgID string) map[string]string {
	return map[string]string{"orgId": orgID}
}

func Commands() *cli.Command {
	return &cli.Command{
		Name:  "desk",
		Usage: "Zoho Desk operations",
		Flags: []cli.Flag{
			&cli.StringFlag{Name: "org", Sources: cli.EnvVars("ZOHO_DESK_ORG_ID"), Usage: "Organization ID (or set ZOHO_DESK_ORG_ID)"},
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
					&cli.BoolFlag{Name: "all", Usage: "Fetch all records"},
					&cli.IntFlag{Name: "limit", Usage: "Max total records to fetch"},
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := internal.RequireFlag(cmd, "org", "ZOHO_DESK_ORG_ID")
					if err != nil {
						return err
					}
					if cmd.Bool("all") || cmd.IsSet("limit") {
						items, err := pagination.Paginate(ctx, pagination.PaginationConfig{
							Client:   c,
							URL:      c.DeskBase + "/departments",
							Opts:     &zohttp.RequestOpts{Headers: orgHeaders(orgID)},
							ItemsKey: "data",
							PageSize: 100,
							Limit:    cmd.Int("limit"),
							SetPage:  pagination.FromLimit(100),
							HasMore:  pagination.HasMoreByCount,
						})
						if err != nil {
							return err
						}
						return output.JSON(items)
					}
					raw, err := c.Request(ctx, "GET", c.DeskBase+"/departments", &zohttp.RequestOpts{
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
				Action: func(ctx context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := internal.RequireFlag(cmd, "org", "ZOHO_DESK_ORG_ID")
					if err != nil {
						return err
					}
					raw, err := c.Request(ctx, "GET", c.DeskBase+"/departments/"+cmd.Args().First(), &zohttp.RequestOpts{
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
					&cli.BoolFlag{Name: "all", Usage: "Fetch all records"},
					&cli.IntFlag{Name: "limit", Usage: "Max total records to fetch"},
					&cli.StringFlag{Name: "department-id", Usage: "Department ID"},
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := internal.RequireFlag(cmd, "org", "ZOHO_DESK_ORG_ID")
					if err != nil {
						return err
					}
					params := map[string]string{}
					if v := cmd.String("department-id"); v != "" {
						params["departmentId"] = v
					}
					if cmd.Bool("all") || cmd.IsSet("limit") {
						items, err := pagination.Paginate(ctx, pagination.PaginationConfig{
							Client:   c,
							URL:      c.DeskBase + "/agents",
							Opts:     &zohttp.RequestOpts{Params: params, Headers: orgHeaders(orgID)},
							ItemsKey: "data",
							PageSize: 100,
							Limit:    cmd.Int("limit"),
							SetPage:  pagination.FromLimit(100),
							HasMore:  pagination.HasMoreByCount,
						})
						if err != nil {
							return err
						}
						return output.JSON(items)
					}
					raw, err := c.Request(ctx, "GET", c.DeskBase+"/agents", &zohttp.RequestOpts{
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
				Action: func(ctx context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := internal.RequireFlag(cmd, "org", "ZOHO_DESK_ORG_ID")
					if err != nil {
						return err
					}
					raw, err := c.Request(ctx, "GET", c.DeskBase+"/agents/"+cmd.Args().First(), &zohttp.RequestOpts{
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
					&cli.BoolFlag{Name: "all", Usage: "Fetch all records"},
					&cli.IntFlag{Name: "limit", Usage: "Max total records to fetch"},
					&cli.StringFlag{Name: "department-ids", Usage: "Comma-separated department IDs"},
					&cli.StringFlag{Name: "status", Usage: "Filter by status"},
					&cli.StringFlag{Name: "priority", Usage: "Filter by priority"},
					&cli.StringFlag{Name: "assignee", Usage: "Assignee ID"},
					&cli.StringFlag{Name: "sort-by", Usage: "Sort field"},
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := internal.RequireFlag(cmd, "org", "ZOHO_DESK_ORG_ID")
					if err != nil {
						return err
					}
					params := map[string]string{}
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
					if cmd.Bool("all") || cmd.IsSet("limit") {
						items, err := pagination.Paginate(ctx, pagination.PaginationConfig{
							Client:   c,
							URL:      c.DeskBase + "/tickets",
							Opts:     &zohttp.RequestOpts{Params: params, Headers: orgHeaders(orgID)},
							ItemsKey: "data",
							PageSize: 100,
							Limit:    cmd.Int("limit"),
							SetPage:  pagination.FromLimit(100),
							HasMore:  pagination.HasMoreByCount,
						})
						if err != nil {
							return err
						}
						return output.JSON(items)
					}
					raw, err := c.Request(ctx, "GET", c.DeskBase+"/tickets", &zohttp.RequestOpts{
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
				Action: func(ctx context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := internal.RequireFlag(cmd, "org", "ZOHO_DESK_ORG_ID")
					if err != nil {
						return err
					}
					raw, err := c.Request(ctx, "GET", c.DeskBase+"/tickets/"+cmd.Args().First(), &zohttp.RequestOpts{
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
					&cli.StringFlag{Name: "subject", Required: true, Usage: "Ticket subject"},
					&cli.StringFlag{Name: "departmentId", Required: true, Usage: "Department ID"},
					&cli.StringFlag{Name: "contactId", Usage: "Contact ID"},
					&cli.StringFlag{Name: "email", Usage: "Contact email"},
					&cli.StringFlag{Name: "phone", Usage: "Contact phone"},
					&cli.StringFlag{Name: "description", Usage: "Ticket description"},
					&cli.StringFlag{Name: "priority", Usage: "Ticket priority"},
					&cli.StringFlag{Name: "status", Usage: "Ticket status"},
					&cli.StringFlag{Name: "channel", Usage: "Ticket channel"},
					&cli.StringFlag{Name: "category", Usage: "Ticket category"},
					&cli.StringFlag{Name: "assigneeId", Usage: "Assignee ID"},
					&cli.StringFlag{Name: "dueDate", Usage: "Due date (ISO 8601)"},
					&cli.StringFlag{Name: "json", Usage: "Additional fields as JSON"},
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := internal.RequireFlag(cmd, "org", "ZOHO_DESK_ORG_ID")
					if err != nil {
						return err
					}
					body := map[string]any{}
					body["subject"] = cmd.String("subject")
					body["departmentId"] = cmd.String("departmentId")
					if cmd.IsSet("contactId") {
						body["contactId"] = cmd.String("contactId")
					}
					if cmd.IsSet("email") {
						body["email"] = cmd.String("email")
					}
					if cmd.IsSet("phone") {
						body["phone"] = cmd.String("phone")
					}
					if cmd.IsSet("description") {
						body["description"] = cmd.String("description")
					}
					if cmd.IsSet("priority") {
						body["priority"] = cmd.String("priority")
					}
					if cmd.IsSet("status") {
						body["status"] = cmd.String("status")
					}
					if cmd.IsSet("channel") {
						body["channel"] = cmd.String("channel")
					}
					if cmd.IsSet("category") {
						body["category"] = cmd.String("category")
					}
					if cmd.IsSet("assigneeId") {
						body["assigneeId"] = cmd.String("assigneeId")
					}
					if cmd.IsSet("dueDate") {
						body["dueDate"] = cmd.String("dueDate")
					}
					if err := internal.MergeJSON(cmd, body); err != nil {
						return err
					}
					raw, err := c.Request(ctx, "POST", c.DeskBase+"/tickets", &zohttp.RequestOpts{
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
					&cli.StringFlag{Name: "subject", Usage: "Ticket subject"},
					&cli.StringFlag{Name: "departmentId", Usage: "Department ID"},
					&cli.StringFlag{Name: "contactId", Usage: "Contact ID"},
					&cli.StringFlag{Name: "email", Usage: "Contact email"},
					&cli.StringFlag{Name: "phone", Usage: "Contact phone"},
					&cli.StringFlag{Name: "description", Usage: "Ticket description"},
					&cli.StringFlag{Name: "priority", Usage: "Ticket priority"},
					&cli.StringFlag{Name: "status", Usage: "Ticket status"},
					&cli.StringFlag{Name: "channel", Usage: "Ticket channel"},
					&cli.StringFlag{Name: "category", Usage: "Ticket category"},
					&cli.StringFlag{Name: "assigneeId", Usage: "Assignee ID"},
					&cli.StringFlag{Name: "dueDate", Usage: "Due date (ISO 8601)"},
					&cli.StringFlag{Name: "resolution", Usage: "Ticket resolution"},
					&cli.StringFlag{Name: "json", Usage: "Additional fields as JSON"},
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := internal.RequireFlag(cmd, "org", "ZOHO_DESK_ORG_ID")
					if err != nil {
						return err
					}
					body := map[string]any{}
					if cmd.IsSet("subject") {
						body["subject"] = cmd.String("subject")
					}
					if cmd.IsSet("departmentId") {
						body["departmentId"] = cmd.String("departmentId")
					}
					if cmd.IsSet("contactId") {
						body["contactId"] = cmd.String("contactId")
					}
					if cmd.IsSet("email") {
						body["email"] = cmd.String("email")
					}
					if cmd.IsSet("phone") {
						body["phone"] = cmd.String("phone")
					}
					if cmd.IsSet("description") {
						body["description"] = cmd.String("description")
					}
					if cmd.IsSet("priority") {
						body["priority"] = cmd.String("priority")
					}
					if cmd.IsSet("status") {
						body["status"] = cmd.String("status")
					}
					if cmd.IsSet("channel") {
						body["channel"] = cmd.String("channel")
					}
					if cmd.IsSet("category") {
						body["category"] = cmd.String("category")
					}
					if cmd.IsSet("assigneeId") {
						body["assigneeId"] = cmd.String("assigneeId")
					}
					if cmd.IsSet("dueDate") {
						body["dueDate"] = cmd.String("dueDate")
					}
					if cmd.IsSet("resolution") {
						body["resolution"] = cmd.String("resolution")
					}
					if err := internal.MergeJSON(cmd, body); err != nil {
						return err
					}
					raw, err := c.Request(ctx, "PATCH", c.DeskBase+"/tickets/"+cmd.Args().First(), &zohttp.RequestOpts{
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
				Action: func(ctx context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := internal.RequireFlag(cmd, "org", "ZOHO_DESK_ORG_ID")
					if err != nil {
						return err
					}
					raw, err := c.Request(ctx, "DELETE", c.DeskBase+"/tickets/"+cmd.Args().First(), &zohttp.RequestOpts{
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
				Action: func(ctx context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := internal.RequireFlag(cmd, "org", "ZOHO_DESK_ORG_ID")
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
					raw, err := c.Request(ctx, "GET", c.DeskBase+"/tickets/search", &zohttp.RequestOpts{
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
				Action: func(ctx context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := internal.RequireFlag(cmd, "org", "ZOHO_DESK_ORG_ID")
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
					raw, err := c.Request(ctx, "GET", c.DeskBase+"/tickets/"+cmd.Args().First()+"/threads", &zohttp.RequestOpts{
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
					&cli.StringFlag{Name: "content", Required: true, Usage: "Reply content"},
					&cli.StringFlag{Name: "contentType", Usage: "Reply content type"},
					&cli.StringFlag{Name: "channel", Usage: "Reply channel"},
					&cli.StringFlag{Name: "to", Usage: "Reply recipient"},
					&cli.StringFlag{Name: "fromEmailAddress", Usage: "Sender email"},
					&cli.BoolFlag{Name: "isForward", Usage: "Send as forward"},
					&cli.StringFlag{Name: "json", Usage: "Additional fields as JSON"},
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := internal.RequireFlag(cmd, "org", "ZOHO_DESK_ORG_ID")
					if err != nil {
						return err
					}
					body := map[string]any{}
					body["content"] = cmd.String("content")
					if cmd.IsSet("contentType") {
						body["contentType"] = cmd.String("contentType")
					}
					if cmd.IsSet("channel") {
						body["channel"] = cmd.String("channel")
					}
					if cmd.IsSet("to") {
						body["to"] = cmd.String("to")
					}
					if cmd.IsSet("fromEmailAddress") {
						body["fromEmailAddress"] = cmd.String("fromEmailAddress")
					}
					if cmd.IsSet("isForward") {
						body["isForward"] = cmd.Bool("isForward")
					}
					if err := internal.MergeJSON(cmd, body); err != nil {
						return err
					}
					raw, err := c.Request(ctx, "POST", c.DeskBase+"/tickets/"+cmd.Args().First()+"/sendReply", &zohttp.RequestOpts{
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
				Action: func(ctx context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := internal.RequireFlag(cmd, "org", "ZOHO_DESK_ORG_ID")
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
					raw, err := c.Request(ctx, "GET", c.DeskBase+"/tickets/"+cmd.Args().First()+"/comments", &zohttp.RequestOpts{
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
					&cli.StringFlag{Name: "content", Required: true, Usage: "Comment content"},
					&cli.BoolFlag{Name: "isPublic", Usage: "Whether the comment is public"},
					&cli.StringFlag{Name: "json", Usage: "Additional fields as JSON"},
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := internal.RequireFlag(cmd, "org", "ZOHO_DESK_ORG_ID")
					if err != nil {
						return err
					}
					body := map[string]any{}
					body["content"] = cmd.String("content")
					if cmd.IsSet("isPublic") {
						body["isPublic"] = cmd.Bool("isPublic")
					}
					if err := internal.MergeJSON(cmd, body); err != nil {
						return err
					}
					raw, err := c.Request(ctx, "POST", c.DeskBase+"/tickets/"+cmd.Args().First()+"/comments", &zohttp.RequestOpts{
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
				Action: func(ctx context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := internal.RequireFlag(cmd, "org", "ZOHO_DESK_ORG_ID")
					if err != nil {
						return err
					}
					raw, err := c.Request(ctx, "GET", c.DeskBase+"/tickets/"+cmd.Args().First()+"/attachments", &zohttp.RequestOpts{
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
				Action: func(ctx context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := internal.RequireFlag(cmd, "org", "ZOHO_DESK_ORG_ID")
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
					raw, err := c.Request(ctx, "GET", c.DeskBase+"/tickets/"+cmd.Args().First()+"/history", &zohttp.RequestOpts{
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
					&cli.BoolFlag{Name: "all", Usage: "Fetch all records"},
					&cli.IntFlag{Name: "limit", Usage: "Max total records to fetch"},
					&cli.StringFlag{Name: "sort-by", Usage: "Sort field"},
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := internal.RequireFlag(cmd, "org", "ZOHO_DESK_ORG_ID")
					if err != nil {
						return err
					}
					params := map[string]string{}
					if v := cmd.String("sort-by"); v != "" {
						params["sortBy"] = v
					}
					if cmd.Bool("all") || cmd.IsSet("limit") {
						items, err := pagination.Paginate(ctx, pagination.PaginationConfig{
							Client:   c,
							URL:      c.DeskBase + "/contacts",
							Opts:     &zohttp.RequestOpts{Params: params, Headers: orgHeaders(orgID)},
							ItemsKey: "data",
							PageSize: 100,
							Limit:    cmd.Int("limit"),
							SetPage:  pagination.FromLimit(100),
							HasMore:  pagination.HasMoreByCount,
						})
						if err != nil {
							return err
						}
						return output.JSON(items)
					}
					raw, err := c.Request(ctx, "GET", c.DeskBase+"/contacts", &zohttp.RequestOpts{
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
				Action: func(ctx context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := internal.RequireFlag(cmd, "org", "ZOHO_DESK_ORG_ID")
					if err != nil {
						return err
					}
					raw, err := c.Request(ctx, "GET", c.DeskBase+"/contacts/"+cmd.Args().First(), &zohttp.RequestOpts{
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
					&cli.StringFlag{Name: "lastName", Required: true, Usage: "Last name"},
					&cli.StringFlag{Name: "email", Required: true, Usage: "Email address"},
					&cli.StringFlag{Name: "firstName", Usage: "First name"},
					&cli.StringFlag{Name: "phone", Usage: "Phone number"},
					&cli.StringFlag{Name: "mobile", Usage: "Mobile number"},
					&cli.StringFlag{Name: "accountId", Usage: "Account ID"},
					&cli.StringFlag{Name: "title", Usage: "Job title"},
					&cli.StringFlag{Name: "city", Usage: "City"},
					&cli.StringFlag{Name: "state", Usage: "State"},
					&cli.StringFlag{Name: "country", Usage: "Country"},
					&cli.StringFlag{Name: "zip", Usage: "ZIP code"},
					&cli.StringFlag{Name: "street", Usage: "Street"},
					&cli.StringFlag{Name: "json", Usage: "Additional fields as JSON"},
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := internal.RequireFlag(cmd, "org", "ZOHO_DESK_ORG_ID")
					if err != nil {
						return err
					}
					body := map[string]any{}
					body["lastName"] = cmd.String("lastName")
					body["email"] = cmd.String("email")
					if cmd.IsSet("firstName") {
						body["firstName"] = cmd.String("firstName")
					}
					if cmd.IsSet("phone") {
						body["phone"] = cmd.String("phone")
					}
					if cmd.IsSet("mobile") {
						body["mobile"] = cmd.String("mobile")
					}
					if cmd.IsSet("accountId") {
						body["accountId"] = cmd.String("accountId")
					}
					if cmd.IsSet("title") {
						body["title"] = cmd.String("title")
					}
					if cmd.IsSet("city") {
						body["city"] = cmd.String("city")
					}
					if cmd.IsSet("state") {
						body["state"] = cmd.String("state")
					}
					if cmd.IsSet("country") {
						body["country"] = cmd.String("country")
					}
					if cmd.IsSet("zip") {
						body["zip"] = cmd.String("zip")
					}
					if cmd.IsSet("street") {
						body["street"] = cmd.String("street")
					}
					if err := internal.MergeJSON(cmd, body); err != nil {
						return err
					}
					raw, err := c.Request(ctx, "POST", c.DeskBase+"/contacts", &zohttp.RequestOpts{
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
					&cli.StringFlag{Name: "lastName", Usage: "Last name"},
					&cli.StringFlag{Name: "email", Usage: "Email address"},
					&cli.StringFlag{Name: "firstName", Usage: "First name"},
					&cli.StringFlag{Name: "phone", Usage: "Phone number"},
					&cli.StringFlag{Name: "mobile", Usage: "Mobile number"},
					&cli.StringFlag{Name: "accountId", Usage: "Account ID"},
					&cli.StringFlag{Name: "title", Usage: "Job title"},
					&cli.StringFlag{Name: "city", Usage: "City"},
					&cli.StringFlag{Name: "state", Usage: "State"},
					&cli.StringFlag{Name: "country", Usage: "Country"},
					&cli.StringFlag{Name: "zip", Usage: "ZIP code"},
					&cli.StringFlag{Name: "street", Usage: "Street"},
					&cli.StringFlag{Name: "json", Usage: "Additional fields as JSON"},
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := internal.RequireFlag(cmd, "org", "ZOHO_DESK_ORG_ID")
					if err != nil {
						return err
					}
					body := map[string]any{}
					if cmd.IsSet("lastName") {
						body["lastName"] = cmd.String("lastName")
					}
					if cmd.IsSet("email") {
						body["email"] = cmd.String("email")
					}
					if cmd.IsSet("firstName") {
						body["firstName"] = cmd.String("firstName")
					}
					if cmd.IsSet("phone") {
						body["phone"] = cmd.String("phone")
					}
					if cmd.IsSet("mobile") {
						body["mobile"] = cmd.String("mobile")
					}
					if cmd.IsSet("accountId") {
						body["accountId"] = cmd.String("accountId")
					}
					if cmd.IsSet("title") {
						body["title"] = cmd.String("title")
					}
					if cmd.IsSet("city") {
						body["city"] = cmd.String("city")
					}
					if cmd.IsSet("state") {
						body["state"] = cmd.String("state")
					}
					if cmd.IsSet("country") {
						body["country"] = cmd.String("country")
					}
					if cmd.IsSet("zip") {
						body["zip"] = cmd.String("zip")
					}
					if cmd.IsSet("street") {
						body["street"] = cmd.String("street")
					}
					if err := internal.MergeJSON(cmd, body); err != nil {
						return err
					}
					raw, err := c.Request(ctx, "PATCH", c.DeskBase+"/contacts/"+cmd.Args().First(), &zohttp.RequestOpts{
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
				Action: func(ctx context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := internal.RequireFlag(cmd, "org", "ZOHO_DESK_ORG_ID")
					if err != nil {
						return err
					}
					raw, err := c.Request(ctx, "DELETE", c.DeskBase+"/contacts/"+cmd.Args().First(), &zohttp.RequestOpts{
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
				Action: func(ctx context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := internal.RequireFlag(cmd, "org", "ZOHO_DESK_ORG_ID")
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
					raw, err := c.Request(ctx, "GET", c.DeskBase+"/contacts/search", &zohttp.RequestOpts{
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
					&cli.BoolFlag{Name: "all", Usage: "Fetch all records"},
					&cli.IntFlag{Name: "limit", Usage: "Max total records to fetch"},
					&cli.StringFlag{Name: "sort-by", Usage: "Sort field"},
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := internal.RequireFlag(cmd, "org", "ZOHO_DESK_ORG_ID")
					if err != nil {
						return err
					}
					params := map[string]string{}
					if v := cmd.String("sort-by"); v != "" {
						params["sortBy"] = v
					}
					if cmd.Bool("all") || cmd.IsSet("limit") {
						items, err := pagination.Paginate(ctx, pagination.PaginationConfig{
							Client:   c,
							URL:      c.DeskBase + "/accounts",
							Opts:     &zohttp.RequestOpts{Params: params, Headers: orgHeaders(orgID)},
							ItemsKey: "data",
							PageSize: 100,
							Limit:    cmd.Int("limit"),
							SetPage:  pagination.FromLimit(100),
							HasMore:  pagination.HasMoreByCount,
						})
						if err != nil {
							return err
						}
						return output.JSON(items)
					}
					raw, err := c.Request(ctx, "GET", c.DeskBase+"/accounts", &zohttp.RequestOpts{
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
				Action: func(ctx context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := internal.RequireFlag(cmd, "org", "ZOHO_DESK_ORG_ID")
					if err != nil {
						return err
					}
					raw, err := c.Request(ctx, "GET", c.DeskBase+"/accounts/"+cmd.Args().First(), &zohttp.RequestOpts{
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
					&cli.StringFlag{Name: "accountName", Required: true, Usage: "Account name"},
					&cli.StringFlag{Name: "email", Usage: "Account email"},
					&cli.StringFlag{Name: "phone", Usage: "Account phone"},
					&cli.StringFlag{Name: "website", Usage: "Account website"},
					&cli.StringFlag{Name: "industry", Usage: "Industry"},
					&cli.StringFlag{Name: "description", Usage: "Account description"},
					&cli.StringFlag{Name: "json", Usage: "Additional fields as JSON"},
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := internal.RequireFlag(cmd, "org", "ZOHO_DESK_ORG_ID")
					if err != nil {
						return err
					}
					body := map[string]any{}
					body["accountName"] = cmd.String("accountName")
					if cmd.IsSet("email") {
						body["email"] = cmd.String("email")
					}
					if cmd.IsSet("phone") {
						body["phone"] = cmd.String("phone")
					}
					if cmd.IsSet("website") {
						body["website"] = cmd.String("website")
					}
					if cmd.IsSet("industry") {
						body["industry"] = cmd.String("industry")
					}
					if cmd.IsSet("description") {
						body["description"] = cmd.String("description")
					}
					if err := internal.MergeJSON(cmd, body); err != nil {
						return err
					}
					raw, err := c.Request(ctx, "POST", c.DeskBase+"/accounts", &zohttp.RequestOpts{
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
					&cli.StringFlag{Name: "accountName", Usage: "Account name"},
					&cli.StringFlag{Name: "email", Usage: "Account email"},
					&cli.StringFlag{Name: "phone", Usage: "Account phone"},
					&cli.StringFlag{Name: "website", Usage: "Account website"},
					&cli.StringFlag{Name: "industry", Usage: "Industry"},
					&cli.StringFlag{Name: "description", Usage: "Account description"},
					&cli.StringFlag{Name: "json", Usage: "Additional fields as JSON"},
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := internal.RequireFlag(cmd, "org", "ZOHO_DESK_ORG_ID")
					if err != nil {
						return err
					}
					body := map[string]any{}
					if cmd.IsSet("accountName") {
						body["accountName"] = cmd.String("accountName")
					}
					if cmd.IsSet("email") {
						body["email"] = cmd.String("email")
					}
					if cmd.IsSet("phone") {
						body["phone"] = cmd.String("phone")
					}
					if cmd.IsSet("website") {
						body["website"] = cmd.String("website")
					}
					if cmd.IsSet("industry") {
						body["industry"] = cmd.String("industry")
					}
					if cmd.IsSet("description") {
						body["description"] = cmd.String("description")
					}
					if err := internal.MergeJSON(cmd, body); err != nil {
						return err
					}
					raw, err := c.Request(ctx, "PATCH", c.DeskBase+"/accounts/"+cmd.Args().First(), &zohttp.RequestOpts{
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
				Action: func(ctx context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := internal.RequireFlag(cmd, "org", "ZOHO_DESK_ORG_ID")
					if err != nil {
						return err
					}
					raw, err := c.Request(ctx, "DELETE", c.DeskBase+"/accounts/"+cmd.Args().First(), &zohttp.RequestOpts{
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
		Action: func(ctx context.Context, cmd *cli.Command) error {
			c, err := zohttp.GetClient()
			if err != nil {
				return err
			}
			orgID, err := internal.RequireFlag(cmd, "org", "ZOHO_DESK_ORG_ID")
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
			raw, err := c.Request(ctx, "GET", c.DeskBase+"/search", &zohttp.RequestOpts{
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
