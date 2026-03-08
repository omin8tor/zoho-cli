package books

import (
	"context"

	"github.com/omin8tor/zoho-cli/internal"
	zohttp "github.com/omin8tor/zoho-cli/internal/http"
	"github.com/omin8tor/zoho-cli/internal/output"
	"github.com/omin8tor/zoho-cli/internal/pagination"
	"github.com/urfave/cli/v3"
)

func booksProjectsCmd() *cli.Command {
	return &cli.Command{
		Name:  "projects",
		Usage: "Project operations",
		Commands: []*cli.Command{
			{
				Name:  "create",
				Usage: "Create a project",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "project_name", Required: true, Usage: "Project name"},
					&cli.StringFlag{Name: "customer_id", Required: true, Usage: "Customer ID"},
					&cli.StringFlag{Name: "currency_id", Required: true, Usage: "Currency ID"},
					&cli.StringFlag{Name: "billing_type", Required: true, Usage: "Billing type"},
					&cli.StringFlag{Name: "description", Usage: "Description"},
					&cli.FloatFlag{Name: "rate", Usage: "Rate"},
					&cli.StringFlag{Name: "budget_type", Usage: "Budget type"},
					&cli.FloatFlag{Name: "budget_hours", Usage: "Budget hours"},
					&cli.FloatFlag{Name: "budget_amount", Usage: "Budget amount"},
					&cli.FloatFlag{Name: "cost_budget_amount", Usage: "Cost budget amount"},
					&cli.StringFlag{Name: "user_id", Usage: "User ID"},
					&cli.StringFlag{Name: "json", Usage: "Additional fields as JSON"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					body := map[string]any{}
					body["project_name"] = cmd.String("project_name")
					body["customer_id"] = cmd.String("customer_id")
					body["currency_id"] = cmd.String("currency_id")
					body["billing_type"] = cmd.String("billing_type")
					if cmd.IsSet("description") {
						body["description"] = cmd.String("description")
					}
					if cmd.IsSet("rate") {
						body["rate"] = cmd.Float("rate")
					}
					if cmd.IsSet("budget_type") {
						body["budget_type"] = cmd.String("budget_type")
					}
					if cmd.IsSet("budget_hours") {
						body["budget_hours"] = cmd.Float("budget_hours")
					}
					if cmd.IsSet("budget_amount") {
						body["budget_amount"] = cmd.Float("budget_amount")
					}
					if cmd.IsSet("cost_budget_amount") {
						body["cost_budget_amount"] = cmd.Float("cost_budget_amount")
					}
					if cmd.IsSet("user_id") {
						body["user_id"] = cmd.String("user_id")
					}
					if err := internal.MergeJSON(cmd, body); err != nil {
						return err
					}
					raw, err := c.Request("POST", c.BooksBase+"/projects", &zohttp.RequestOpts{
						Params: orgParams(orgID),
						JSON:   body,
					})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:  "update-by-custom-field",
				Usage: "Update a project by custom field",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "project_name", Required: true, Usage: "Project name"},
					&cli.StringFlag{Name: "customer_id", Required: true, Usage: "Customer ID"},
					&cli.StringFlag{Name: "currency_id", Required: true, Usage: "Currency ID"},
					&cli.StringFlag{Name: "billing_type", Required: true, Usage: "Billing type"},
					&cli.StringFlag{Name: "description", Usage: "Description"},
					&cli.FloatFlag{Name: "rate", Usage: "Rate"},
					&cli.StringFlag{Name: "budget_type", Usage: "Budget type"},
					&cli.FloatFlag{Name: "budget_hours", Usage: "Budget hours"},
					&cli.FloatFlag{Name: "budget_amount", Usage: "Budget amount"},
					&cli.FloatFlag{Name: "cost_budget_amount", Usage: "Cost budget amount"},
					&cli.StringFlag{Name: "user_id", Usage: "User ID"},
					&cli.StringFlag{Name: "json", Usage: "Additional fields as JSON"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					body := map[string]any{}
					body["project_name"] = cmd.String("project_name")
					body["customer_id"] = cmd.String("customer_id")
					body["currency_id"] = cmd.String("currency_id")
					body["billing_type"] = cmd.String("billing_type")
					if cmd.IsSet("description") {
						body["description"] = cmd.String("description")
					}
					if cmd.IsSet("rate") {
						body["rate"] = cmd.Float("rate")
					}
					if cmd.IsSet("budget_type") {
						body["budget_type"] = cmd.String("budget_type")
					}
					if cmd.IsSet("budget_hours") {
						body["budget_hours"] = cmd.Float("budget_hours")
					}
					if cmd.IsSet("budget_amount") {
						body["budget_amount"] = cmd.Float("budget_amount")
					}
					if cmd.IsSet("cost_budget_amount") {
						body["cost_budget_amount"] = cmd.Float("cost_budget_amount")
					}
					if cmd.IsSet("user_id") {
						body["user_id"] = cmd.String("user_id")
					}
					if err := internal.MergeJSON(cmd, body); err != nil {
						return err
					}
					raw, err := c.Request("PUT", c.BooksBase+"/projects", &zohttp.RequestOpts{
						Params: orgParams(orgID),
						JSON:   body,
					})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:  "list",
				Usage: "List projects",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "status", Usage: "Filter"},
					&cli.StringFlag{Name: "customer-id", Usage: "Filter"},
					&cli.BoolFlag{Name: "all", Usage: "Fetch all records"},
					&cli.IntFlag{Name: "limit", Usage: "Max total records to fetch"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					params := orgParams(orgID)
					if v := cmd.String("status"); v != "" {
						params["status"] = v
					}
					if v := cmd.String("customer-id"); v != "" {
						params["customer_id"] = v
					}

					if cmd.Bool("all") || cmd.IsSet("limit") {

						items, err := pagination.Paginate(pagination.PaginationConfig{

							Client: c,

							URL: c.BooksBase + "/projects",

							Opts: &zohttp.RequestOpts{Params: params},

							ItemsKey: "projects",

							PageSize: 200,

							Limit: int(cmd.Int("limit")),

							SetPage: pagination.PagePerPage(200),

							HasMore: pagination.HasMoreBooks,
						})

						if err != nil {

							return err

						}

						return output.JSON(items)

					}
					raw, err := c.Request("GET", c.BooksBase+"/projects", &zohttp.RequestOpts{Params: params})
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
					&cli.StringFlag{Name: "project_name", Usage: "Project name"},
					&cli.StringFlag{Name: "customer_id", Usage: "Customer ID"},
					&cli.StringFlag{Name: "currency_id", Usage: "Currency ID"},
					&cli.StringFlag{Name: "billing_type", Usage: "Billing type"},
					&cli.StringFlag{Name: "description", Usage: "Description"},
					&cli.FloatFlag{Name: "rate", Usage: "Rate"},
					&cli.StringFlag{Name: "budget_type", Usage: "Budget type"},
					&cli.FloatFlag{Name: "budget_hours", Usage: "Budget hours"},
					&cli.FloatFlag{Name: "budget_amount", Usage: "Budget amount"},
					&cli.FloatFlag{Name: "cost_budget_amount", Usage: "Cost budget amount"},
					&cli.StringFlag{Name: "user_id", Usage: "User ID"},
					&cli.StringFlag{Name: "json", Usage: "Additional fields as JSON"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					body := map[string]any{}
					if cmd.IsSet("project_name") {
						body["project_name"] = cmd.String("project_name")
					}
					if cmd.IsSet("customer_id") {
						body["customer_id"] = cmd.String("customer_id")
					}
					if cmd.IsSet("currency_id") {
						body["currency_id"] = cmd.String("currency_id")
					}
					if cmd.IsSet("billing_type") {
						body["billing_type"] = cmd.String("billing_type")
					}
					if cmd.IsSet("description") {
						body["description"] = cmd.String("description")
					}
					if cmd.IsSet("rate") {
						body["rate"] = cmd.Float("rate")
					}
					if cmd.IsSet("budget_type") {
						body["budget_type"] = cmd.String("budget_type")
					}
					if cmd.IsSet("budget_hours") {
						body["budget_hours"] = cmd.Float("budget_hours")
					}
					if cmd.IsSet("budget_amount") {
						body["budget_amount"] = cmd.Float("budget_amount")
					}
					if cmd.IsSet("cost_budget_amount") {
						body["cost_budget_amount"] = cmd.Float("cost_budget_amount")
					}
					if cmd.IsSet("user_id") {
						body["user_id"] = cmd.String("user_id")
					}
					if err := internal.MergeJSON(cmd, body); err != nil {
						return err
					}
					raw, err := c.Request("PUT", c.BooksBase+"/projects/"+cmd.Args().First(), &zohttp.RequestOpts{
						Params: orgParams(orgID),
						JSON:   body,
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
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("GET", c.BooksBase+"/projects/"+cmd.Args().First(), &zohttp.RequestOpts{Params: orgParams(orgID)})
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
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("DELETE", c.BooksBase+"/projects/"+cmd.Args().First(), &zohttp.RequestOpts{Params: orgParams(orgID)})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "activate",
				Usage:     "Activate a project",
				ArgsUsage: "<project-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("POST", c.BooksBase+"/projects/"+cmd.Args().First()+"/active", &zohttp.RequestOpts{Params: orgParams(orgID)})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "inactivate",
				Usage:     "Inactivate a project",
				ArgsUsage: "<project-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("POST", c.BooksBase+"/projects/"+cmd.Args().First()+"/inactive", &zohttp.RequestOpts{Params: orgParams(orgID)})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "clone",
				Usage:     "Clone a project",
				ArgsUsage: "<project-id>",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "project_name", Required: true, Usage: "Project name"},
					&cli.StringFlag{Name: "start_date", Usage: "Start date (YYYY-MM-DD)"},
					&cli.StringFlag{Name: "end_date", Usage: "End date (YYYY-MM-DD)"},
					&cli.BoolFlag{Name: "clone_tasks", Usage: "Clone tasks"},
					&cli.BoolFlag{Name: "clone_users", Usage: "Clone users"},
					&cli.StringFlag{Name: "json", Usage: "Additional fields as JSON"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					body := map[string]any{}
					body["project_name"] = cmd.String("project_name")
					if cmd.IsSet("start_date") {
						body["start_date"] = cmd.String("start_date")
					}
					if cmd.IsSet("end_date") {
						body["end_date"] = cmd.String("end_date")
					}
					if cmd.IsSet("clone_tasks") {
						body["clone_tasks"] = cmd.Bool("clone_tasks")
					}
					if cmd.IsSet("clone_users") {
						body["clone_users"] = cmd.Bool("clone_users")
					}
					if err := internal.MergeJSON(cmd, body); err != nil {
						return err
					}
					opts := &zohttp.RequestOpts{Params: orgParams(orgID), JSON: body}
					raw, err := c.Request("POST", c.BooksBase+"/projects/"+cmd.Args().First()+"/clone", opts)
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "assign-users",
				Usage:     "Assign users to a project",
				ArgsUsage: "<project-id>",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "user_id", Required: true, Usage: "User ID"},
					&cli.FloatFlag{Name: "rate", Usage: "Rate"},
					&cli.FloatFlag{Name: "budget_hours", Usage: "Budget hours"},
					&cli.FloatFlag{Name: "cost_rate", Usage: "Cost rate"},
					&cli.StringFlag{Name: "user_role", Usage: "User role"},
					&cli.StringFlag{Name: "json", Usage: "Additional fields as JSON"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					body := map[string]any{}
					body["user_id"] = cmd.String("user_id")
					if cmd.IsSet("rate") {
						body["rate"] = cmd.Float("rate")
					}
					if cmd.IsSet("budget_hours") {
						body["budget_hours"] = cmd.Float("budget_hours")
					}
					if cmd.IsSet("cost_rate") {
						body["cost_rate"] = cmd.Float("cost_rate")
					}
					if cmd.IsSet("user_role") {
						body["user_role"] = cmd.String("user_role")
					}
					if err := internal.MergeJSON(cmd, body); err != nil {
						return err
					}
					raw, err := c.Request("POST", c.BooksBase+"/projects/"+cmd.Args().First()+"/users", &zohttp.RequestOpts{
						Params: orgParams(orgID),
						JSON:   body,
					})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "list-users",
				Usage:     "List users of a project",
				ArgsUsage: "<project-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("GET", c.BooksBase+"/projects/"+cmd.Args().First()+"/users", &zohttp.RequestOpts{Params: orgParams(orgID)})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "invite-user",
				Usage:     "Invite a user to a project",
				ArgsUsage: "<project-id>",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "user_id", Usage: "User ID"},
					&cli.StringFlag{Name: "email", Usage: "Email address"},
					&cli.StringFlag{Name: "user_name", Usage: "User name"},
					&cli.StringFlag{Name: "user_role", Usage: "User role"},
					&cli.FloatFlag{Name: "rate", Usage: "Rate"},
					&cli.FloatFlag{Name: "budget_hours", Usage: "Budget hours"},
					&cli.FloatFlag{Name: "cost_rate", Usage: "Cost rate"},
					&cli.StringFlag{Name: "json", Usage: "Additional fields as JSON"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					body := map[string]any{}
					if cmd.IsSet("user_id") {
						body["user_id"] = cmd.String("user_id")
					}
					if cmd.IsSet("email") {
						body["email"] = cmd.String("email")
					}
					if cmd.IsSet("user_name") {
						body["user_name"] = cmd.String("user_name")
					}
					if cmd.IsSet("user_role") {
						body["user_role"] = cmd.String("user_role")
					}
					if cmd.IsSet("rate") {
						body["rate"] = cmd.Float("rate")
					}
					if cmd.IsSet("budget_hours") {
						body["budget_hours"] = cmd.Float("budget_hours")
					}
					if cmd.IsSet("cost_rate") {
						body["cost_rate"] = cmd.Float("cost_rate")
					}
					if err := internal.MergeJSON(cmd, body); err != nil {
						return err
					}
					raw, err := c.Request("POST", c.BooksBase+"/projects/"+cmd.Args().First()+"/users/invite", &zohttp.RequestOpts{
						Params: orgParams(orgID),
						JSON:   body,
					})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "update-user",
				Usage:     "Update a user in a project",
				ArgsUsage: "<project-id> <user-id>",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "user_name", Usage: "User name"},
					&cli.StringFlag{Name: "user_role", Usage: "User role"},
					&cli.FloatFlag{Name: "rate", Usage: "Rate"},
					&cli.FloatFlag{Name: "budget_hours", Usage: "Budget hours"},
					&cli.FloatFlag{Name: "cost_rate", Usage: "Cost rate"},
					&cli.StringFlag{Name: "json", Usage: "Additional fields as JSON"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					body := map[string]any{}
					if cmd.IsSet("user_name") {
						body["user_name"] = cmd.String("user_name")
					}
					if cmd.IsSet("user_role") {
						body["user_role"] = cmd.String("user_role")
					}
					if cmd.IsSet("rate") {
						body["rate"] = cmd.Float("rate")
					}
					if cmd.IsSet("budget_hours") {
						body["budget_hours"] = cmd.Float("budget_hours")
					}
					if cmd.IsSet("cost_rate") {
						body["cost_rate"] = cmd.Float("cost_rate")
					}
					if err := internal.MergeJSON(cmd, body); err != nil {
						return err
					}
					raw, err := c.Request("PUT", c.BooksBase+"/projects/"+cmd.Args().First()+"/users/"+cmd.Args().Get(1), &zohttp.RequestOpts{
						Params: orgParams(orgID),
						JSON:   body,
					})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "get-user",
				Usage:     "Get a user in a project",
				ArgsUsage: "<project-id> <user-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("GET", c.BooksBase+"/projects/"+cmd.Args().First()+"/users/"+cmd.Args().Get(1), &zohttp.RequestOpts{Params: orgParams(orgID)})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "delete-user",
				Usage:     "Delete a user from a project",
				ArgsUsage: "<project-id> <user-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("DELETE", c.BooksBase+"/projects/"+cmd.Args().First()+"/users/"+cmd.Args().Get(1), &zohttp.RequestOpts{Params: orgParams(orgID)})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "post-comment",
				Usage:     "Post a comment on a project",
				ArgsUsage: "<project-id>",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "description", Required: true, Usage: "Description"},
					&cli.StringFlag{Name: "json", Usage: "Additional fields as JSON"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					body := map[string]any{}
					body["description"] = cmd.String("description")
					if err := internal.MergeJSON(cmd, body); err != nil {
						return err
					}
					raw, err := c.Request("POST", c.BooksBase+"/projects/"+cmd.Args().First()+"/comments", &zohttp.RequestOpts{
						Params: orgParams(orgID),
						JSON:   body,
					})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "list-comments",
				Usage:     "List comments of a project",
				ArgsUsage: "<project-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("GET", c.BooksBase+"/projects/"+cmd.Args().First()+"/comments", &zohttp.RequestOpts{Params: orgParams(orgID)})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "delete-comment",
				Usage:     "Delete a comment on a project",
				ArgsUsage: "<project-id> <comment-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("DELETE", c.BooksBase+"/projects/"+cmd.Args().First()+"/comments/"+cmd.Args().Get(1), &zohttp.RequestOpts{Params: orgParams(orgID)})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "list-invoices",
				Usage:     "List invoices of a project",
				ArgsUsage: "<project-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("GET", c.BooksBase+"/projects/"+cmd.Args().First()+"/invoices", &zohttp.RequestOpts{Params: orgParams(orgID)})
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
		Usage: "Task operations",
		Commands: []*cli.Command{
			{
				Name:  "create",
				Usage: "Create a task",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "project-id", Required: true, Usage: "Project ID"},
					&cli.StringFlag{Name: "task_name", Required: true, Usage: "Task name"},
					&cli.StringFlag{Name: "description", Usage: "Description"},
					&cli.FloatFlag{Name: "rate", Usage: "Rate"},
					&cli.FloatFlag{Name: "budget_hours", Usage: "Budget hours"},
					&cli.StringFlag{Name: "json", Usage: "Additional fields as JSON"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					body := map[string]any{}
					body["task_name"] = cmd.String("task_name")
					if cmd.IsSet("description") {
						body["description"] = cmd.String("description")
					}
					if cmd.IsSet("rate") {
						body["rate"] = cmd.Float("rate")
					}
					if cmd.IsSet("budget_hours") {
						body["budget_hours"] = cmd.Float("budget_hours")
					}
					if err := internal.MergeJSON(cmd, body); err != nil {
						return err
					}
					raw, err := c.Request("POST", c.BooksBase+"/projects/"+cmd.String("project-id")+"/tasks", &zohttp.RequestOpts{
						Params: orgParams(orgID),
						JSON:   body,
					})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:  "list",
				Usage: "List tasks",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "project-id", Required: true, Usage: "Project ID"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("GET", c.BooksBase+"/projects/"+cmd.String("project-id")+"/tasks", &zohttp.RequestOpts{Params: orgParams(orgID)})
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
					&cli.StringFlag{Name: "project-id", Required: true, Usage: "Project ID"},
					&cli.StringFlag{Name: "task_name", Usage: "Task name"},
					&cli.StringFlag{Name: "description", Usage: "Description"},
					&cli.FloatFlag{Name: "rate", Usage: "Rate"},
					&cli.FloatFlag{Name: "budget_hours", Usage: "Budget hours"},
					&cli.StringFlag{Name: "json", Usage: "Additional fields as JSON"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					body := map[string]any{}
					if cmd.IsSet("task_name") {
						body["task_name"] = cmd.String("task_name")
					}
					if cmd.IsSet("description") {
						body["description"] = cmd.String("description")
					}
					if cmd.IsSet("rate") {
						body["rate"] = cmd.Float("rate")
					}
					if cmd.IsSet("budget_hours") {
						body["budget_hours"] = cmd.Float("budget_hours")
					}
					if err := internal.MergeJSON(cmd, body); err != nil {
						return err
					}
					raw, err := c.Request("PUT", c.BooksBase+"/projects/"+cmd.String("project-id")+"/tasks/"+cmd.Args().First(), &zohttp.RequestOpts{
						Params: orgParams(orgID),
						JSON:   body,
					})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "get",
				Usage:     "Get a task",
				ArgsUsage: "<task-id>",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "project-id", Required: true, Usage: "Project ID"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("GET", c.BooksBase+"/projects/"+cmd.String("project-id")+"/tasks/"+cmd.Args().First(), &zohttp.RequestOpts{Params: orgParams(orgID)})
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
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "project-id", Required: true, Usage: "Project ID"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("DELETE", c.BooksBase+"/projects/"+cmd.String("project-id")+"/tasks/"+cmd.Args().First(), &zohttp.RequestOpts{Params: orgParams(orgID)})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
		},
	}
}

func timeEntriesCmd() *cli.Command {
	return &cli.Command{
		Name:  "time-entries",
		Usage: "Time entry operations",
		Commands: []*cli.Command{
			{
				Name:  "log",
				Usage: "Log a time entry",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "project-id", Required: true, Usage: "Project ID"},
					&cli.StringFlag{Name: "user_id", Required: true, Usage: "User ID"},
					&cli.StringFlag{Name: "task_id", Usage: "Task ID"},
					&cli.StringFlag{Name: "log_date", Required: true, Usage: "Log date (YYYY-MM-DD)"},
					&cli.StringFlag{Name: "log_time", Usage: "Log time (HH:mm)"},
					&cli.StringFlag{Name: "begin_time", Usage: "Begin time (HH:mm)"},
					&cli.StringFlag{Name: "end_time", Usage: "End time (HH:mm)"},
					&cli.BoolFlag{Name: "is_billable", Usage: "Is billable"},
					&cli.StringFlag{Name: "notes", Usage: "Notes"},
					&cli.BoolFlag{Name: "start_timer", Usage: "Start timer"},
					&cli.FloatFlag{Name: "cost_rate", Usage: "Cost rate"},
					&cli.StringFlag{Name: "json", Usage: "Additional fields as JSON"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					body := map[string]any{}
					body["user_id"] = cmd.String("user_id")
					if cmd.IsSet("task_id") {
						body["task_id"] = cmd.String("task_id")
					}
					body["log_date"] = cmd.String("log_date")
					if cmd.IsSet("log_time") {
						body["log_time"] = cmd.String("log_time")
					}
					if cmd.IsSet("begin_time") {
						body["begin_time"] = cmd.String("begin_time")
					}
					if cmd.IsSet("end_time") {
						body["end_time"] = cmd.String("end_time")
					}
					if cmd.IsSet("is_billable") {
						body["is_billable"] = cmd.Bool("is_billable")
					}
					if cmd.IsSet("notes") {
						body["notes"] = cmd.String("notes")
					}
					if cmd.IsSet("start_timer") {
						body["start_timer"] = cmd.Bool("start_timer")
					}
					if cmd.IsSet("cost_rate") {
						body["cost_rate"] = cmd.Float("cost_rate")
					}
					if err := internal.MergeJSON(cmd, body); err != nil {
						return err
					}
					raw, err := c.Request("POST", c.BooksBase+"/projects/"+cmd.String("project-id")+"/timeentries", &zohttp.RequestOpts{
						Params: orgParams(orgID),
						JSON:   body,
					})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:  "list",
				Usage: "List time entries",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "project-id", Required: true, Usage: "Project ID"},
					&cli.BoolFlag{Name: "all", Usage: "Fetch all records"},
					&cli.IntFlag{Name: "limit", Usage: "Max total records to fetch"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					params := orgParams(orgID)

					if cmd.Bool("all") || cmd.IsSet("limit") {

						items, err := pagination.Paginate(pagination.PaginationConfig{

							Client: c,

							URL: c.BooksBase + "/projects/" + cmd.String("project-id") + "/timeentries",

							Opts: &zohttp.RequestOpts{Params: params},

							ItemsKey: "time_entries",

							PageSize: 200,

							Limit: int(cmd.Int("limit")),

							SetPage: pagination.PagePerPage(200),

							HasMore: pagination.HasMoreBooks,
						})

						if err != nil {

							return err

						}

						return output.JSON(items)

					}
					raw, err := c.Request("GET", c.BooksBase+"/projects/"+cmd.String("project-id")+"/timeentries", &zohttp.RequestOpts{Params: params})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "update",
				Usage:     "Update a time entry",
				ArgsUsage: "<timeentry-id>",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "project-id", Required: true, Usage: "Project ID"},
					&cli.StringFlag{Name: "user_id", Usage: "User ID"},
					&cli.StringFlag{Name: "task_id", Usage: "Task ID"},
					&cli.StringFlag{Name: "log_date", Usage: "Log date (YYYY-MM-DD)"},
					&cli.StringFlag{Name: "log_time", Usage: "Log time (HH:mm)"},
					&cli.StringFlag{Name: "begin_time", Usage: "Begin time (HH:mm)"},
					&cli.StringFlag{Name: "end_time", Usage: "End time (HH:mm)"},
					&cli.BoolFlag{Name: "is_billable", Usage: "Is billable"},
					&cli.StringFlag{Name: "notes", Usage: "Notes"},
					&cli.BoolFlag{Name: "start_timer", Usage: "Start timer"},
					&cli.FloatFlag{Name: "cost_rate", Usage: "Cost rate"},
					&cli.StringFlag{Name: "json", Usage: "Additional fields as JSON"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					body := map[string]any{}
					if cmd.IsSet("user_id") {
						body["user_id"] = cmd.String("user_id")
					}
					if cmd.IsSet("task_id") {
						body["task_id"] = cmd.String("task_id")
					}
					if cmd.IsSet("log_date") {
						body["log_date"] = cmd.String("log_date")
					}
					if cmd.IsSet("log_time") {
						body["log_time"] = cmd.String("log_time")
					}
					if cmd.IsSet("begin_time") {
						body["begin_time"] = cmd.String("begin_time")
					}
					if cmd.IsSet("end_time") {
						body["end_time"] = cmd.String("end_time")
					}
					if cmd.IsSet("is_billable") {
						body["is_billable"] = cmd.Bool("is_billable")
					}
					if cmd.IsSet("notes") {
						body["notes"] = cmd.String("notes")
					}
					if cmd.IsSet("start_timer") {
						body["start_timer"] = cmd.Bool("start_timer")
					}
					if cmd.IsSet("cost_rate") {
						body["cost_rate"] = cmd.Float("cost_rate")
					}
					if err := internal.MergeJSON(cmd, body); err != nil {
						return err
					}
					raw, err := c.Request("PUT", c.BooksBase+"/projects/"+cmd.String("project-id")+"/timeentries/"+cmd.Args().First(), &zohttp.RequestOpts{
						Params: orgParams(orgID),
						JSON:   body,
					})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "get",
				Usage:     "Get a time entry",
				ArgsUsage: "<timeentry-id>",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "project-id", Required: true, Usage: "Project ID"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("GET", c.BooksBase+"/projects/"+cmd.String("project-id")+"/timeentries/"+cmd.Args().First(), &zohttp.RequestOpts{Params: orgParams(orgID)})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "delete",
				Usage:     "Delete a time entry",
				ArgsUsage: "<timeentry-id>",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "project-id", Required: true, Usage: "Project ID"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("DELETE", c.BooksBase+"/projects/"+cmd.String("project-id")+"/timeentries/"+cmd.Args().First(), &zohttp.RequestOpts{Params: orgParams(orgID)})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "start-timer",
				Usage:     "Start a timer for a time entry",
				ArgsUsage: "<timeentry-id>",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "project-id", Required: true, Usage: "Project ID"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("POST", c.BooksBase+"/projects/"+cmd.String("project-id")+"/timeentries/"+cmd.Args().First()+"/timer/start", &zohttp.RequestOpts{Params: orgParams(orgID)})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "stop-timer",
				Usage:     "Stop a timer for a time entry",
				ArgsUsage: "<timeentry-id>",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "project-id", Required: true, Usage: "Project ID"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("POST", c.BooksBase+"/projects/"+cmd.String("project-id")+"/timeentries/"+cmd.Args().First()+"/timer/stop", &zohttp.RequestOpts{Params: orgParams(orgID)})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:  "get-timer",
				Usage: "Get current running timer",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "project-id", Required: true, Usage: "Project ID"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("GET", c.BooksBase+"/projects/"+cmd.String("project-id")+"/timeentries/timer", &zohttp.RequestOpts{Params: orgParams(orgID)})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
		},
	}
}
