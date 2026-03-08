package people

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/omin8tor/zoho-cli/internal"
	zohttp "github.com/omin8tor/zoho-cli/internal/http"
	"github.com/omin8tor/zoho-cli/internal/output"
	"github.com/omin8tor/zoho-cli/internal/pagination"
	"github.com/urfave/cli/v3"
)

func Commands() *cli.Command {
	return &cli.Command{
		Name:  "people",
		Usage: "Zoho People operations",
		Commands: []*cli.Command{
			formsCmd(),
			recordsCmd(),
			attendanceCmd(),
			leaveCmd(),
			departmentsCmd(),
			designationsCmd(),
		},
	}
}

func formsCmd() *cli.Command {
	return &cli.Command{
		Name:  "forms",
		Usage: "Form operations",
		Commands: []*cli.Command{
			{
				Name:  "list",
				Usage: "List all forms",
				Action: func(ctx context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					raw, err := c.Request(ctx, "GET", c.PeopleBase+"/forms", nil)
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "get-fields",
				Usage:     "Get fields for a form",
				ArgsUsage: "<form-link-name>",
				Action: func(ctx context.Context, cmd *cli.Command) error {
					if cmd.Args().Len() < 1 {
						return internal.NewValidationError("form link name required")
					}
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					raw, err := c.Request(ctx, "GET", c.PeopleBase+"/forms/"+cmd.Args().First()+"/components", nil)
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
		},
	}
}

func recordsCmd() *cli.Command {
	return &cli.Command{
		Name:  "records",
		Usage: "Form record operations",
		Commands: []*cli.Command{
			{
				Name:      "list",
				Usage:     "List records from a form",
				ArgsUsage: "<form-link-name>",
				Flags: []cli.Flag{
					&cli.BoolFlag{Name: "all", Usage: "Fetch all records"},
					&cli.IntFlag{Name: "limit", Usage: "Max total records to fetch"},
					&cli.StringFlag{Name: "search-column", Usage: "Search column (EMPLOYEEID or EMPLOYEEMAILALIAS)"},
					&cli.StringFlag{Name: "search-value", Usage: "Search value"},
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					if cmd.Args().Len() < 1 {
						return internal.NewValidationError("form link name required")
					}
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					params := map[string]string{}
					if v := cmd.String("search-column"); v != "" {
						params["searchColumn"] = v
					}
					if v := cmd.String("search-value"); v != "" {
						params["searchValue"] = v
					}
					if cmd.Bool("all") || cmd.IsSet("limit") {
						items, err := pagination.Paginate(ctx, pagination.PaginationConfig{
							Client:   c,
							URL:      c.PeopleBase + "/forms/" + cmd.Args().First() + "/getRecords",
							Opts:     &zohttp.RequestOpts{Params: params},
							ItemsKey: "response.result",
							PageSize: 200,
							Limit:    cmd.Int("limit"),
							SetPage:  pagination.SIndexLimit(200),
							HasMore:  pagination.HasMoreByCount,
						})
						if err != nil {
							return err
						}
						return output.JSON(items)
					}
					raw, err := c.Request(ctx, "GET", c.PeopleBase+"/forms/"+cmd.Args().First()+"/getRecords", &zohttp.RequestOpts{
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
				Usage:     "Get a single record by ID",
				ArgsUsage: "<form-link-name> <record-id>",
				Action: func(ctx context.Context, cmd *cli.Command) error {
					if cmd.Args().Len() < 2 {
						return internal.NewValidationError("form link name and record ID required")
					}
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					raw, err := c.Request(ctx, "GET", c.PeopleBase+"/forms/"+cmd.Args().First()+"/getDataByID", &zohttp.RequestOpts{
						Params: map[string]string{"recordId": cmd.Args().Get(1)},
					})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "add",
				Usage:     "Add a record to a form",
				ArgsUsage: "<form-link-name>",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "json", Required: true, Usage: "JSON input data"},
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					if cmd.Args().Len() < 1 {
						return internal.NewValidationError("form link name required")
					}
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					var body any
					if err := json.Unmarshal([]byte(cmd.String("json")), &body); err != nil {
						return internal.NewValidationError(fmt.Sprintf("invalid JSON: %v", err))
					}
					inputData, err := json.Marshal(body)
					if err != nil {
						return err
					}
					raw, err := c.Request(ctx, "POST", c.PeopleBase+"/forms/json/"+cmd.Args().First()+"/insertRecord", &zohttp.RequestOpts{
						Form: map[string]string{"inputData": string(inputData)},
					})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "update",
				Usage:     "Update a record in a form",
				ArgsUsage: "<form-link-name> <record-id>",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "json", Required: true, Usage: "JSON input data"},
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					if cmd.Args().Len() < 2 {
						return internal.NewValidationError("form link name and record ID required")
					}
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					var body any
					if err := json.Unmarshal([]byte(cmd.String("json")), &body); err != nil {
						return internal.NewValidationError(fmt.Sprintf("invalid JSON: %v", err))
					}
					inputData, err := json.Marshal(body)
					if err != nil {
						return err
					}
					raw, err := c.Request(ctx, "POST", c.PeopleBase+"/forms/json/"+cmd.Args().First()+"/updateRecord", &zohttp.RequestOpts{
						Form: map[string]string{
							"inputData": string(inputData),
							"recordId":  cmd.Args().Get(1),
						},
					})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "delete",
				Usage:     "Delete a record from a form",
				ArgsUsage: "<form-link-name> <record-id>",
				Action: func(ctx context.Context, cmd *cli.Command) error {
					if cmd.Args().Len() < 2 {
						return internal.NewValidationError("form link name and record ID required")
					}
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					raw, err := c.Request(ctx, "POST", c.PeopleBase+"/forms/"+cmd.Args().First()+"/deleteRecord", &zohttp.RequestOpts{
						Form: map[string]string{"recordId": cmd.Args().Get(1)},
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

func attendanceCmd() *cli.Command {
	return &cli.Command{
		Name:  "attendance",
		Usage: "Attendance operations",
		Commands: []*cli.Command{
			{
				Name:  "checkin",
				Usage: "Record attendance check-in/check-out",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "emp-id", Usage: "Employee ID"},
					&cli.StringFlag{Name: "email-id", Usage: "Employee email ID"},
					&cli.StringFlag{Name: "check-in", Usage: "Check-in time (dd/MM/yyyy HH:mm:ss)"},
					&cli.StringFlag{Name: "check-out", Usage: "Check-out time (dd/MM/yyyy HH:mm:ss)"},
					&cli.StringFlag{Name: "date-format", Usage: "Date format"},
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					params := map[string]string{}
					if v := cmd.String("emp-id"); v != "" {
						params["empId"] = v
					}
					if v := cmd.String("email-id"); v != "" {
						params["emailId"] = v
					}
					if v := cmd.String("check-in"); v != "" {
						params["checkIn"] = v
					}
					if v := cmd.String("check-out"); v != "" {
						params["checkOut"] = v
					}
					if v := cmd.String("date-format"); v != "" {
						params["dateFormat"] = v
					}
					raw, err := c.Request(ctx, "POST", c.PeopleBase+"/attendance", &zohttp.RequestOpts{
						Params: params,
					})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:  "get-report",
				Usage: "Get attendance report for an employee",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "sdate", Required: true, Usage: "Start date"},
					&cli.StringFlag{Name: "edate", Required: true, Usage: "End date"},
					&cli.StringFlag{Name: "emp-id", Usage: "Employee ID"},
					&cli.StringFlag{Name: "email-id", Usage: "Employee email ID"},
					&cli.StringFlag{Name: "date-format", Usage: "Date format (e.g. yyyy-MM-dd)"},
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					params := map[string]string{
						"sdate": cmd.String("sdate"),
						"edate": cmd.String("edate"),
					}
					if v := cmd.String("emp-id"); v != "" {
						params["empId"] = v
					}
					if v := cmd.String("email-id"); v != "" {
						params["emailId"] = v
					}
					if v := cmd.String("date-format"); v != "" {
						params["dateFormat"] = v
					}
					raw, err := c.Request(ctx, "GET", c.PeopleBase+"/attendance/getUserReport", &zohttp.RequestOpts{
						Params: params,
					})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:  "bulk-import",
				Usage: "Bulk import attendance records",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "json", Required: true, Usage: "JSON array of attendance records"},
					&cli.StringFlag{Name: "date-format", Usage: "Date format"},
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					params := map[string]string{}
					if v := cmd.String("date-format"); v != "" {
						params["dateFormat"] = v
					}
					raw, err := c.Request(ctx, "POST", c.PeopleBase+"/attendance/bulkImport", &zohttp.RequestOpts{
						Params: params,
						Form:   map[string]string{"data": cmd.String("json")},
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

func leaveCmd() *cli.Command {
	return &cli.Command{
		Name:  "leave",
		Usage: "Leave operations",
		Commands: []*cli.Command{
			{
				Name:  "list",
				Usage: "List leave records",
				Flags: []cli.Flag{
					&cli.BoolFlag{Name: "all", Usage: "Fetch all records"},
					&cli.IntFlag{Name: "limit", Usage: "Max total records to fetch"},
					&cli.StringFlag{Name: "search-column", Usage: "Search column"},
					&cli.StringFlag{Name: "search-value", Usage: "Search value"},
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					params := map[string]string{}
					if v := cmd.String("search-column"); v != "" {
						params["searchColumn"] = v
					}
					if v := cmd.String("search-value"); v != "" {
						params["searchValue"] = v
					}
					if cmd.Bool("all") || cmd.IsSet("limit") {
						items, err := pagination.Paginate(ctx, pagination.PaginationConfig{
							Client:   c,
							URL:      c.PeopleBase + "/forms/leave/getRecords",
							Opts:     &zohttp.RequestOpts{Params: params},
							ItemsKey: "response.result",
							PageSize: 200,
							Limit:    cmd.Int("limit"),
							SetPage:  pagination.SIndexLimit(200),
							HasMore:  pagination.HasMoreByCount,
						})
						if err != nil {
							return err
						}
						return output.JSON(items)
					}
					raw, err := c.Request(ctx, "GET", c.PeopleBase+"/forms/leave/getRecords", &zohttp.RequestOpts{
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
				Usage:     "Get a leave record by ID",
				ArgsUsage: "<record-id>",
				Action: func(ctx context.Context, cmd *cli.Command) error {
					if cmd.Args().Len() < 1 {
						return internal.NewValidationError("leave record ID required")
					}
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					raw, err := c.Request(ctx, "GET", c.PeopleBase+"/forms/leave/getDataByID", &zohttp.RequestOpts{
						Params: map[string]string{"recordId": cmd.Args().First()},
					})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:  "add",
				Usage: "Add a leave record",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "json", Required: true, Usage: "JSON input data"},
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					var body any
					if err := json.Unmarshal([]byte(cmd.String("json")), &body); err != nil {
						return internal.NewValidationError(fmt.Sprintf("invalid JSON: %v", err))
					}
					inputData, err := json.Marshal(body)
					if err != nil {
						return err
					}
					raw, err := c.Request(ctx, "POST", c.PeopleBase+"/forms/json/leave/insertRecord", &zohttp.RequestOpts{
						Form: map[string]string{"inputData": string(inputData)},
					})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:  "get-types",
				Usage: "Get leave types",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "user-id", Usage: "User record ID"},
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					params := map[string]string{}
					if v := cmd.String("user-id"); v != "" {
						params["userId"] = v
					}
					raw, err := c.Request(ctx, "GET", c.PeopleBase+"/leave/getLeaveTypeDetails", &zohttp.RequestOpts{
						Params: params,
					})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:  "get-balance",
				Usage: "Get leave balance for an employee",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "user-id", Usage: "User record ID"},
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					params := map[string]string{}
					if v := cmd.String("user-id"); v != "" {
						params["userId"] = v
					}
					raw, err := c.Request(ctx, "GET", c.PeopleBase+"/leave/getLeaveTypeDetails", &zohttp.RequestOpts{
						Params: params,
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

func departmentsCmd() *cli.Command {
	return &cli.Command{
		Name:  "departments",
		Usage: "Department operations",
		Commands: []*cli.Command{
			{
				Name:  "list",
				Usage: "List all departments",
				Action: func(ctx context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					raw, err := c.Request(ctx, "GET", c.PeopleBase+"/forms/department/getRecords", nil)
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "get-members",
				Usage:     "Get department members",
				ArgsUsage: "<department>",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "sindex", Usage: "Starting index"},
					&cli.StringFlag{Name: "limit", Usage: "Max records"},
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					if cmd.Args().Len() < 1 {
						return internal.NewValidationError("department name required")
					}
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					params := map[string]string{
						"department": cmd.Args().First(),
					}
					if v := cmd.String("sindex"); v != "" {
						params["sIndex"] = v
					}
					if v := cmd.String("limit"); v != "" {
						params["limit"] = v
					}
					raw, err := c.Request(ctx, "GET", c.PeopleBase+"/forms/department/getRecords", &zohttp.RequestOpts{
						Params: params,
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

func designationsCmd() *cli.Command {
	return &cli.Command{
		Name:  "designations",
		Usage: "Designation operations",
		Commands: []*cli.Command{
			{
				Name:  "list",
				Usage: "List all designations",
				Action: func(ctx context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					raw, err := c.Request(ctx, "GET", c.PeopleBase+"/forms/designation/getRecords", nil)
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
		},
	}
}
