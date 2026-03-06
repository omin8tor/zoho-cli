package sign

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"

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

func Commands() *cli.Command {
	return &cli.Command{
		Name:  "sign",
		Usage: "Zoho Sign operations",
		Commands: []*cli.Command{
			requestsCmd(),
			templatesCmd(),
			foldersCmd(),
			fieldTypesCmd(),
			requestTypesCmd(),
		},
	}
}

func requestsCmd() *cli.Command {
	return &cli.Command{
		Name:  "requests",
		Usage: "Document request operations",
		Commands: []*cli.Command{
			{
				Name:  "list",
				Usage: "List sign requests",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "start-index", Usage: "Starting index for pagination"},
					&cli.StringFlag{Name: "row-count", Usage: "Number of records per page"},
					&cli.StringFlag{Name: "sort-column", Usage: "Sort field (e.g. created_time)"},
					&cli.StringFlag{Name: "sort-order", Usage: "ASC or DESC"},
					&cli.StringFlag{Name: "data", Usage: "Full JSON data param (overrides other pagination flags)"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					params := map[string]string{}
					if v := cmd.String("data"); v != "" {
						params["data"] = v
					} else {
						pc := map[string]any{}
						if v := cmd.String("start-index"); v != "" {
							if n, err := strconv.Atoi(v); err == nil {
								pc["start_index"] = n
							}
						}
						if v := cmd.String("row-count"); v != "" {
							if n, err := strconv.Atoi(v); err == nil {
								pc["row_count"] = n
							}
						}
						if v := cmd.String("sort-column"); v != "" {
							pc["sort_column"] = v
						}
						if v := cmd.String("sort-order"); v != "" {
							pc["sort_order"] = v
						}
						if len(pc) > 0 {
							j, _ := json.Marshal(map[string]any{"page_context": pc})
							params["data"] = string(j)
						}
					}
					raw, err := c.Request("GET", c.SignBase+"/requests", &zohttp.RequestOpts{
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
				Usage:     "Get a sign request",
				ArgsUsage: "<request-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					id := cmd.Args().Get(0)
					if id == "" {
						return internal.NewValidationError("request-id argument required")
					}
					c, err := getClient()
					if err != nil {
						return err
					}
					raw, err := c.Request("GET", c.SignBase+"/requests/"+id, nil)
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:  "create",
				Usage: "Create a sign request",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "file", Required: true, Usage: "Path to document file"},
					&cli.StringFlag{Name: "data", Required: true, Usage: "JSON data for the request"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					filePath := cmd.String("file")
					fileData, err := os.ReadFile(filePath)
					if err != nil {
						return fmt.Errorf("failed to read file: %w", err)
					}
					raw, err := c.Request("POST", c.SignBase+"/requests", &zohttp.RequestOpts{
						Files: map[string]zohttp.FileUpload{"file": {Filename: filepath.Base(filePath), Data: fileData}},
						Form:  map[string]string{"data": cmd.String("data")},
					})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "update",
				Usage:     "Update a sign request",
				ArgsUsage: "<request-id>",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "data", Required: true, Usage: "JSON data for the update"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					id := cmd.Args().Get(0)
					if id == "" {
						return internal.NewValidationError("request-id argument required")
					}
					c, err := getClient()
					if err != nil {
						return err
					}
					raw, err := c.Request("PUT", c.SignBase+"/requests/"+id, &zohttp.RequestOpts{
						Form: map[string]string{"data": cmd.String("data")},
					})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "submit",
				Usage:     "Submit a sign request for signature",
				ArgsUsage: "<request-id>",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "data", Required: true, Usage: "JSON data with fields and actions"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					id := cmd.Args().Get(0)
					if id == "" {
						return internal.NewValidationError("request-id argument required")
					}
					c, err := getClient()
					if err != nil {
						return err
					}
					raw, err := c.Request("POST", c.SignBase+"/requests/"+id+"/submit", &zohttp.RequestOpts{
						Form: map[string]string{"data": cmd.String("data")},
					})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "delete",
				Usage:     "Delete a sign request",
				ArgsUsage: "<request-id>",
				Flags: []cli.Flag{
					&cli.BoolFlag{Name: "recall-inprogress", Usage: "Recall if document is in progress"},
					&cli.StringFlag{Name: "reason", Usage: "Reason for recalling"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					id := cmd.Args().Get(0)
					if id == "" {
						return internal.NewValidationError("request-id argument required")
					}
					c, err := getClient()
					if err != nil {
						return err
					}
					form := map[string]string{}
					if cmd.Bool("recall-inprogress") {
						form["recall_inprogress"] = "true"
					}
					if v := cmd.String("reason"); v != "" {
						form["reason"] = v
					}
					var opts *zohttp.RequestOpts
					if len(form) > 0 {
						opts = &zohttp.RequestOpts{Form: form}
					}
					raw, err := c.Request("PUT", c.SignBase+"/requests/"+id+"/delete", opts)
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "recall",
				Usage:     "Recall a sign request",
				ArgsUsage: "<request-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					id := cmd.Args().Get(0)
					if id == "" {
						return internal.NewValidationError("request-id argument required")
					}
					c, err := getClient()
					if err != nil {
						return err
					}
					raw, err := c.Request("POST", c.SignBase+"/requests/"+id+"/recall", nil)
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "remind",
				Usage:     "Send reminder for a sign request",
				ArgsUsage: "<request-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					id := cmd.Args().Get(0)
					if id == "" {
						return internal.NewValidationError("request-id argument required")
					}
					c, err := getClient()
					if err != nil {
						return err
					}
					raw, err := c.Request("POST", c.SignBase+"/requests/"+id+"/remind", nil)
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "extend",
				Usage:     "Extend expiration of a sign request",
				ArgsUsage: "<request-id>",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "expire-by", Required: true, Usage: "New expiration date"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					id := cmd.Args().Get(0)
					if id == "" {
						return internal.NewValidationError("request-id argument required")
					}
					c, err := getClient()
					if err != nil {
						return err
					}
					raw, err := c.Request("PUT", c.SignBase+"/requests/"+id+"/extend", &zohttp.RequestOpts{
						Form: map[string]string{"expire_by": cmd.String("expire-by")},
					})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "correct",
				Usage:     "Mark a sign request for correction",
				ArgsUsage: "<request-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					id := cmd.Args().Get(0)
					if id == "" {
						return internal.NewValidationError("request-id argument required")
					}
					c, err := getClient()
					if err != nil {
						return err
					}
					raw, err := c.Request("POST", c.SignBase+"/requests/"+id+"/markforcorrection", nil)
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "download",
				Usage:     "Download sign request PDF",
				ArgsUsage: "<request-id>",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "output", Usage: "Output file path"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					id := cmd.Args().Get(0)
					if id == "" {
						return internal.NewValidationError("request-id argument required")
					}
					c, err := getClient()
					if err != nil {
						return err
					}
					body, _, _, err := c.RequestRaw("GET", c.SignBase+"/requests/"+id+"/pdf", nil)
					if err != nil {
						return err
					}
					if out := cmd.String("output"); out != "" {
						if err := os.WriteFile(out, body, 0644); err != nil {
							return err
						}
						return output.JSON(map[string]any{"ok": true, "path": out, "size": len(body)})
					}
					os.Stdout.Write(body)
					return nil
				},
			},
			{
				Name:      "download-document",
				Usage:     "Download a particular document PDF from a sign request",
				ArgsUsage: "<request-id> <document-id>",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "output", Usage: "Output file path"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					id := cmd.Args().Get(0)
					docID := cmd.Args().Get(1)
					if id == "" || docID == "" {
						return internal.NewValidationError("request-id and document-id arguments required")
					}
					c, err := getClient()
					if err != nil {
						return err
					}
					body, _, _, err := c.RequestRaw("GET", c.SignBase+"/requests/"+id+"/documents/"+docID+"/pdf", nil)
					if err != nil {
						return err
					}
					if out := cmd.String("output"); out != "" {
						if err := os.WriteFile(out, body, 0644); err != nil {
							return err
						}
						return output.JSON(map[string]any{"ok": true, "path": out, "size": len(body)})
					}
					os.Stdout.Write(body)
					return nil
				},
			},
			{
				Name:      "download-certificate",
				Usage:     "Download completion certificate for a sign request",
				ArgsUsage: "<request-id>",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "output", Usage: "Output file path"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					id := cmd.Args().Get(0)
					if id == "" {
						return internal.NewValidationError("request-id argument required")
					}
					c, err := getClient()
					if err != nil {
						return err
					}
					body, _, _, err := c.RequestRaw("GET", c.SignBase+"/requests/"+id+"/completioncertificate", nil)
					if err != nil {
						return err
					}
					if out := cmd.String("output"); out != "" {
						if err := os.WriteFile(out, body, 0644); err != nil {
							return err
						}
						return output.JSON(map[string]any{"ok": true, "path": out, "size": len(body)})
					}
					os.Stdout.Write(body)
					return nil
				},
			},
			{
				Name:      "field-data",
				Usage:     "Get form field data for a sign request",
				ArgsUsage: "<request-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					id := cmd.Args().Get(0)
					if id == "" {
						return internal.NewValidationError("request-id argument required")
					}
					c, err := getClient()
					if err != nil {
						return err
					}
					raw, err := c.Request("GET", c.SignBase+"/requests/"+id+"/fielddata", nil)
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
		},
	}
}

func templatesCmd() *cli.Command {
	return &cli.Command{
		Name:  "templates",
		Usage: "Template operations",
		Commands: []*cli.Command{
			{
				Name:  "list",
				Usage: "List templates",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "start-index", Usage: "Starting index for pagination"},
					&cli.StringFlag{Name: "row-count", Usage: "Number of records per page"},
					&cli.StringFlag{Name: "sort-column", Usage: "Sort field"},
					&cli.StringFlag{Name: "sort-order", Usage: "ASC or DESC"},
					&cli.StringFlag{Name: "data", Usage: "Full JSON data param (overrides other pagination flags)"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					params := map[string]string{}
					if v := cmd.String("data"); v != "" {
						params["data"] = v
					} else {
						pc := map[string]any{}
						if v := cmd.String("start-index"); v != "" {
							if n, err := strconv.Atoi(v); err == nil {
								pc["start_index"] = n
							}
						}
						if v := cmd.String("row-count"); v != "" {
							if n, err := strconv.Atoi(v); err == nil {
								pc["row_count"] = n
							}
						}
						if v := cmd.String("sort-column"); v != "" {
							pc["sort_column"] = v
						}
						if v := cmd.String("sort-order"); v != "" {
							pc["sort_order"] = v
						}
						if len(pc) > 0 {
							j, _ := json.Marshal(map[string]any{"page_context": pc})
							params["data"] = string(j)
						}
					}
					raw, err := c.Request("GET", c.SignBase+"/templates", &zohttp.RequestOpts{
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
				Usage:     "Get a template",
				ArgsUsage: "<template-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					id := cmd.Args().Get(0)
					if id == "" {
						return internal.NewValidationError("template-id argument required")
					}
					c, err := getClient()
					if err != nil {
						return err
					}
					raw, err := c.Request("GET", c.SignBase+"/templates/"+id, nil)
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:  "create",
				Usage: "Create a template",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "file", Required: true, Usage: "Path to document file"},
					&cli.StringFlag{Name: "data", Required: true, Usage: "JSON data for the template"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					filePath := cmd.String("file")
					fileData, err := os.ReadFile(filePath)
					if err != nil {
						return fmt.Errorf("failed to read file: %w", err)
					}
					raw, err := c.Request("POST", c.SignBase+"/templates", &zohttp.RequestOpts{
						Files: map[string]zohttp.FileUpload{"file": {Filename: filepath.Base(filePath), Data: fileData}},
						Form:  map[string]string{"data": cmd.String("data")},
					})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "send",
				Usage:     "Send document for signature using template",
				ArgsUsage: "<template-id>",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "data", Required: true, Usage: "JSON data with recipients and fields"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					id := cmd.Args().Get(0)
					if id == "" {
						return internal.NewValidationError("template-id argument required")
					}
					c, err := getClient()
					if err != nil {
						return err
					}
					raw, err := c.Request("POST", c.SignBase+"/templates/"+id+"/createdocument", &zohttp.RequestOpts{
						Form: map[string]string{"data": cmd.String("data")},
					})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "delete",
				Usage:     "Delete a template",
				ArgsUsage: "<template-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					id := cmd.Args().Get(0)
					if id == "" {
						return internal.NewValidationError("template-id argument required")
					}
					c, err := getClient()
					if err != nil {
						return err
					}
					raw, err := c.Request("PUT", c.SignBase+"/templates/"+id+"/delete", nil)
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
		},
	}
}

func foldersCmd() *cli.Command {
	return &cli.Command{
		Name:  "folders",
		Usage: "Folder operations",
		Commands: []*cli.Command{
			{
				Name:  "list",
				Usage: "List folders",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					raw, err := c.Request("GET", c.SignBase+"/folders", nil)
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:  "create",
				Usage: "Create a folder",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "name", Required: true, Usage: "Folder name"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					raw, err := c.Request("POST", c.SignBase+"/folders", &zohttp.RequestOpts{
						Form: map[string]string{"folder_name": cmd.String("name")},
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

func fieldTypesCmd() *cli.Command {
	return &cli.Command{
		Name:  "field-types",
		Usage: "List available field types",
		Action: func(_ context.Context, cmd *cli.Command) error {
			c, err := getClient()
			if err != nil {
				return err
			}
			raw, err := c.Request("GET", c.SignBase+"/fieldtypes", nil)
			if err != nil {
				return err
			}
			return output.JSONRaw(raw)
		},
	}
}

func requestTypesCmd() *cli.Command {
	return &cli.Command{
		Name:  "request-types",
		Usage: "Document type operations",
		Commands: []*cli.Command{
			{
				Name:  "list",
				Usage: "List request types",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					raw, err := c.Request("GET", c.SignBase+"/requesttypes", nil)
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:  "create",
				Usage: "Create a request type",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "data", Required: true, Usage: "JSON data for the request type"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					raw, err := c.Request("POST", c.SignBase+"/requesttypes", &zohttp.RequestOpts{
						Form: map[string]string{"data": cmd.String("data")},
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
