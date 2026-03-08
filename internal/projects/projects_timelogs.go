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
				Action: func(ctx context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					portal, err := requirePortal(cmd)
					if err != nil {
						return err
					}
					url := base(c, portal, cmd.String("project")) + "/timelogs"
					moduleJSON := fmt.Sprintf(`{"type":"%s"}`, cmd.String("module"))
					raw, err := c.Request(ctx, "GET", url, &zohttp.RequestOpts{
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
				Action: func(ctx context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
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
					raw, err := c.Request(ctx, "POST", url, &zohttp.RequestOpts{
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
				Action: func(ctx context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					portal, err := requirePortal(cmd)
					if err != nil {
						return err
					}
					url := base(c, portal, cmd.String("project")) + "/logs/" + cmd.Args().First()
					raw, err := c.Request(ctx, "GET", url, &zohttp.RequestOpts{Params: map[string]string{"type": cmd.String("type")}})
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
					&cli.FloatFlag{Name: "hours", Usage: "Hours (e.g. 2, 1.5, 0:30)"},
					&cli.StringFlag{Name: "date", Usage: "Date (YYYY-MM-DD)"},
					&cli.StringFlag{Name: "notes", Usage: "Notes for time entry"},
					&cli.StringFlag{Name: "bill-status", Usage: "Billable or Non Billable"},
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
					if cmd.IsSet("hours") {
						body["hours"] = cmd.Float("hours")
					}
					if cmd.IsSet("date") {
						body["date"] = cmd.String("date")
					}
					if cmd.IsSet("notes") {
						body["notes"] = cmd.String("notes")
					}
					if cmd.IsSet("bill-status") {
						body["bill_status"] = cmd.String("bill-status")
					}
					if err := internal.MergeJSON(cmd, body); err != nil {
						return err
					}
					if _, ok := body["module"]; !ok {
						mod := map[string]string{"type": cmd.String("type")}
						if t := cmd.String("task"); t != "" {
							mod["id"] = t
						}
						body["module"] = mod
					}
					url := base(c, portal, cmd.String("project")) + "/logs/" + cmd.Args().First()
					raw, err := c.Request(ctx, "PATCH", url, &zohttp.RequestOpts{JSON: body})
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
				Action: func(ctx context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					portal, err := requirePortal(cmd)
					if err != nil {
						return err
					}
					url := base(c, portal, cmd.String("project")) + "/logs/" + cmd.Args().First()
					raw, err := c.Request(ctx, "DELETE", url, &zohttp.RequestOpts{JSON: map[string]string{"module": cmd.String("type")}})
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
				Flags: []cli.Flag{
					portalFlag,
					&cli.StringFlag{Name: "module", Required: true, Usage: "Module type (task or bug)"},
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
					url := c.ProjectsBase + "/portal/" + portal + "/timelogs"
					raw, err := c.Request(ctx, "GET", url, &zohttp.RequestOpts{
						Params: map[string]string{"module": cmd.String("module")},
					})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:  "project-list",
				Usage: "List timelogs (project-level)",
				Flags: []cli.Flag{
					portalFlag, projectFlag,
					&cli.StringFlag{Name: "module", Required: true, Usage: "Module type (task or bug)"},
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
					url := base(c, portal, cmd.String("project")) + "/timelogs"
					raw, err := c.Request(ctx, "GET", url, &zohttp.RequestOpts{
						Params: map[string]string{"module": cmd.String("module")},
					})
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
					&cli.StringFlag{Name: "log-object", Required: true, Usage: "Timelogs as JSON array"},
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
					form := map[string]string{"log_object": cmd.String("log-object")}
					if err := internal.MergeJSONForm(cmd, form); err != nil {
						return err
					}
					url := c.ProjectsBase + "/portal/" + portal + "/addbulktimelogs"
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
				Name:  "delete",
				Usage: "Bulk delete timelogs",
				Flags: []cli.Flag{
					portalFlag,
					&cli.StringFlag{Name: "ids", Required: true, Usage: "JSON array of {id, module} objects"},
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
					url := c.ProjectsBase + "/portal/" + portal + "/timelogs/bulkdelete"
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

func timelogTimersCmd() *cli.Command {
	return &cli.Command{
		Name:  "timelog-timers",
		Usage: "Timer operations",
		Commands: []*cli.Command{
			{
				Name:  "running",
				Usage: "Get running timers",
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
					url := c.ProjectsBase + "/portal/" + portal + "/timesheet/timers"
					raw, err := c.Request(ctx, "GET", url, nil)
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
					portalFlag,
					&cli.StringFlag{Name: "entity_id", Usage: "Entity ID (task/issue)"},
					&cli.StringFlag{Name: "project_id", Usage: "Project ID"},
					&cli.StringFlag{Name: "module_id", Usage: "Module ID"},
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
					if cmd.IsSet("entity_id") {
						body["entity_id"] = cmd.String("entity_id")
					}
					if cmd.IsSet("project_id") {
						body["project_id"] = cmd.String("project_id")
					}
					if cmd.IsSet("module_id") {
						body["module_id"] = cmd.String("module_id")
					}
					if err := internal.MergeJSON(cmd, body); err != nil {
						return err
					}
					url := c.ProjectsBase + "/portal/" + portal + "/timesheet/timers"
					raw, err := c.Request(ctx, "POST", url, &zohttp.RequestOpts{JSON: body})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "get",
				Usage:     "Get a timer",
				ArgsUsage: "<timer-id>",
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
					url := c.ProjectsBase + "/portal/" + portal + "/timesheet/timers/" + cmd.Args().First()
					raw, err := c.Request(ctx, "GET", url, nil)
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "pause",
				Usage:     "Pause a timer",
				ArgsUsage: "<timer-id>",
				Flags: []cli.Flag{
					portalFlag,
					&cli.StringFlag{Name: "notes", Usage: "Timer notes"},
					&cli.StringFlag{Name: "type", Usage: "Timer entity type"},
					&cli.StringFlag{Name: "log_id", Usage: "Timelog ID"},
					&cli.StringFlag{Name: "entity_id", Usage: "Task or issue entity ID"},
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
					if cmd.IsSet("notes") {
						body["notes"] = cmd.String("notes")
					}
					if cmd.IsSet("type") {
						body["type"] = cmd.String("type")
					}
					if cmd.IsSet("log_id") {
						body["log_id"] = cmd.String("log_id")
					}
					if cmd.IsSet("entity_id") {
						body["entity_id"] = cmd.String("entity_id")
					}
					if err := internal.MergeJSON(cmd, body); err != nil {
						return err
					}
					url := c.ProjectsBase + "/portal/" + portal + "/timesheet/timers/" + cmd.Args().First() + "/pause"
					raw, err := c.Request(ctx, "PATCH", url, &zohttp.RequestOpts{JSON: body})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "resume",
				Usage:     "Resume a timer",
				ArgsUsage: "<timer-id>",
				Flags: []cli.Flag{
					portalFlag,
					&cli.StringFlag{Name: "notes", Usage: "Timer notes"},
					&cli.StringFlag{Name: "type", Usage: "Timer entity type"},
					&cli.StringFlag{Name: "log_id", Usage: "Timelog ID"},
					&cli.StringFlag{Name: "entity_id", Usage: "Task or issue entity ID"},
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
					if cmd.IsSet("notes") {
						body["notes"] = cmd.String("notes")
					}
					if cmd.IsSet("type") {
						body["type"] = cmd.String("type")
					}
					if cmd.IsSet("log_id") {
						body["log_id"] = cmd.String("log_id")
					}
					if cmd.IsSet("entity_id") {
						body["entity_id"] = cmd.String("entity_id")
					}
					if err := internal.MergeJSON(cmd, body); err != nil {
						return err
					}
					url := c.ProjectsBase + "/portal/" + portal + "/timesheet/timers/" + cmd.Args().First() + "/resume"
					raw, err := c.Request(ctx, "PATCH", url, &zohttp.RequestOpts{JSON: body})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "stop",
				Usage:     "Stop a timer",
				ArgsUsage: "<timer-id>",
				Flags: []cli.Flag{
					portalFlag,
					&cli.StringFlag{Name: "item_id", Usage: "Item ID"},
					&cli.StringFlag{Name: "log_name", Usage: "Timelog name"},
					&cli.StringFlag{Name: "date", Usage: "Date"},
					&cli.StringFlag{Name: "project_id", Usage: "Project ID"},
					&cli.StringFlag{Name: "type", Usage: "Timer entity type"},
					&cli.FloatFlag{Name: "hours", Usage: "Total hours"},
					&cli.StringFlag{Name: "start_time", Usage: "Start time"},
					&cli.StringFlag{Name: "end_time", Usage: "End time"},
					&cli.StringFlag{Name: "bill_status", Usage: "Bill status"},
					&cli.StringFlag{Name: "notes", Usage: "Timer notes"},
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
					if cmd.IsSet("item_id") {
						body["item_id"] = cmd.String("item_id")
					}
					if cmd.IsSet("log_name") {
						body["log_name"] = cmd.String("log_name")
					}
					if cmd.IsSet("date") {
						body["date"] = cmd.String("date")
					}
					if cmd.IsSet("project_id") {
						body["project_id"] = cmd.String("project_id")
					}
					if cmd.IsSet("type") {
						body["type"] = cmd.String("type")
					}
					if cmd.IsSet("hours") {
						body["hours"] = cmd.Float("hours")
					}
					if cmd.IsSet("start_time") {
						body["start_time"] = cmd.String("start_time")
					}
					if cmd.IsSet("end_time") {
						body["end_time"] = cmd.String("end_time")
					}
					if cmd.IsSet("bill_status") {
						body["bill_status"] = cmd.String("bill_status")
					}
					if cmd.IsSet("notes") {
						body["notes"] = cmd.String("notes")
					}
					if err := internal.MergeJSON(cmd, body); err != nil {
						return err
					}
					url := c.ProjectsBase + "/portal/" + portal + "/timesheet/timers/" + cmd.Args().First() + "/stop"
					raw, err := c.Request(ctx, "PATCH", url, &zohttp.RequestOpts{JSON: body})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "delete",
				Usage:     "Delete a timer",
				ArgsUsage: "<timer-id>",
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
					url := c.ProjectsBase + "/portal/" + portal + "/timesheet/timers/" + cmd.Args().First()
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

func timelogPinsCmd() *cli.Command {
	return &cli.Command{
		Name:  "timelog-pins",
		Usage: "Timelog pin operations",
		Commands: []*cli.Command{
			{
				Name:  "list",
				Usage: "List timelog pins",
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
					url := c.ProjectsBase + "/portal/" + portal + "/timesheet/pin"
					raw, err := c.Request(ctx, "GET", url, nil)
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:  "create",
				Usage: "Pin a timelog",
				Flags: []cli.Flag{
					portalFlag,
					&cli.StringFlag{Name: "project_id", Usage: "Project ID"},
					&cli.StringFlag{Name: "module", Usage: "Module type"},
					&cli.StringFlag{Name: "zpuid", Usage: "User ZPUID"},
					&cli.IntFlag{Name: "sequence", Usage: "Pin sequence"},
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
					if cmd.IsSet("project_id") {
						body["project_id"] = cmd.String("project_id")
					}
					if cmd.IsSet("module") {
						body["module"] = cmd.String("module")
					}
					if cmd.IsSet("zpuid") {
						body["zpuid"] = cmd.String("zpuid")
					}
					if cmd.IsSet("sequence") {
						body["sequence"] = cmd.Int("sequence")
					}
					if err := internal.MergeJSON(cmd, body); err != nil {
						return err
					}
					url := c.ProjectsBase + "/portal/" + portal + "/timesheet/pin"
					raw, err := c.Request(ctx, "POST", url, &zohttp.RequestOpts{JSON: body})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "update",
				Usage:     "Update a timelog pin",
				ArgsUsage: "<pin-id>",
				Flags: []cli.Flag{
					portalFlag,
					&cli.IntFlag{Name: "sequence", Usage: "Pin sequence"},
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
					if cmd.IsSet("sequence") {
						body["sequence"] = cmd.Int("sequence")
					}
					if err := internal.MergeJSON(cmd, body); err != nil {
						return err
					}
					url := c.ProjectsBase + "/portal/" + portal + "/timesheet/pin/" + cmd.Args().First()
					raw, err := c.Request(ctx, "PATCH", url, &zohttp.RequestOpts{JSON: body})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "delete",
				Usage:     "Unpin a timelog",
				ArgsUsage: "<pin-id>",
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
					url := c.ProjectsBase + "/portal/" + portal + "/timesheet/pin/" + cmd.Args().First()
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
