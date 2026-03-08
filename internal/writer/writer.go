package writer

import (
	"context"
	"encoding/json"
	"os"
	"strings"

	"github.com/omin8tor/zoho-cli/internal"
	zohttp "github.com/omin8tor/zoho-cli/internal/http"
	"github.com/omin8tor/zoho-cli/internal/output"
	"github.com/urfave/cli/v3"
)

const jsonapiCT = "application/vnd.api+json"

var serviceTypeMap = map[string]string{
	"writer": "zw",
	"sheet":  "zohosheet",
	"show":   "zohoshow",
}

func Commands() *cli.Command {
	return &cli.Command{
		Name:  "writer",
		Usage: "Zoho Writer operations",
		Commands: []*cli.Command{
			{
				Name:  "create",
				Usage: "Create a new Writer document",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "name", Required: true, Usage: "Document name"},
					&cli.StringFlag{Name: "folder", Required: true, Usage: "Parent folder ID in WorkDrive"},
					&cli.StringFlag{Name: "type", Value: "writer", Usage: "writer, sheet, or show"},
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					mapped, ok := serviceTypeMap[cmd.String("type")]
					if !ok {
						return internal.NewValidationError("invalid --type; expected writer, sheet, or show")
					}
					body := map[string]any{
						"data": map[string]any{
							"type": "files",
							"attributes": map[string]any{
								"name":         cmd.String("name"),
								"parent_id":    cmd.String("folder"),
								"service_type": mapped,
							},
						},
					}
					raw, err := c.Request(ctx, "POST", c.WorkDriveBase+"/files", &zohttp.RequestOpts{
						JSON:    body,
						Headers: map[string]string{"Content-Type": jsonapiCT},
					})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "details",
				Usage:     "Get document metadata",
				ArgsUsage: "<doc-id>",
				Action: func(ctx context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					raw, err := c.Request(ctx, "GET", c.WriterBase+"/documents/"+cmd.Args().First(), nil)
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "fields",
				Usage:     "List merge fields in a document",
				ArgsUsage: "<doc-id>",
				Action: func(ctx context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					raw, err := c.Request(ctx, "GET", c.WriterBase+"/documents/"+cmd.Args().First()+"/fields", nil)
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "merge",
				Usage:     "Merge data into a document template",
				ArgsUsage: "<doc-id>",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "json", Required: true, Usage: "Merge data as JSON"},
					&cli.StringFlag{Name: "format", Value: "pdf", Usage: "pdf, docx, or inline"},
					&cli.StringFlag{Name: "output", Usage: "Output file path"},
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					docID := cmd.Args().First()
					var mergeData any
					if err := json.Unmarshal([]byte(cmd.String("json")), &mergeData); err != nil {
						return internal.NewValidationError("invalid --json: must be valid JSON")
					}

					if cmd.String("format") == "inline" {
						body := map[string]any{
							"merge_data":    mergeData,
							"output_format": "inline",
						}
						raw, err := c.Request(ctx, "POST", c.WriterBase+"/documents/"+docID+"/merge", &zohttp.RequestOpts{JSON: body})
						if err != nil {
							return err
						}
						return output.JSONRaw(raw)
					}

					params := map[string]string{
						"output_format": cmd.String("format"),
					}
					mergeJSON, _ := json.Marshal(mergeData)
					params["merge_data"] = string(mergeJSON)
					body, _, _, err := c.RequestRaw(ctx, "POST", c.WriterBase+"/documents/"+docID+"/merge", params)
					if err != nil {
						return err
					}
					if out := cmd.String("output"); out != "" {
						if err := os.WriteFile(out, body, 0600); err != nil {
							return err
						}
						return output.JSON(map[string]any{"ok": true, "path": out, "size": len(body)})
					}
					os.Stdout.Write(body)
					return nil
				},
			},
			{
				Name:      "trash",
				Usage:     "Move a document to trash",
				ArgsUsage: "<doc-id>",
				Action: func(ctx context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					raw, err := c.Request(ctx, "DELETE", c.WriterBase+"/documents/"+cmd.Args().First()+"/trash", nil)
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "delete",
				Usage:     "Permanently delete a trashed document",
				ArgsUsage: "<doc-id>",
				Action: func(ctx context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					raw, err := c.Request(ctx, "DELETE", c.WriterBase+"/documents/"+cmd.Args().First()+"/delete", nil)
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "read",
				Usage:     "Read document content as text",
				ArgsUsage: "<doc-id>",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "format", Value: "txt", Usage: "txt or html"},
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					body, _, _, err := c.RequestRaw(ctx, "GET", c.WriterBase+"/download/"+cmd.Args().First(), map[string]string{"format": cmd.String("format")})
					if err != nil {
						if strings.Contains(err.Error(), "R3002") {
							return output.JSON(map[string]string{"error": "Document is empty — Zoho cannot export empty documents (R3002)"})
						}
						return err
					}
					if len(body) == 0 {
						return output.JSON(map[string]string{"error": "Document is empty or could not be read"})
					}
					return output.JSON(map[string]any{
						"document_id": cmd.Args().First(),
						"format":      cmd.String("format"),
						"content":     string(body),
					})
				},
			},
			{
				Name:      "download",
				Usage:     "Download a document",
				ArgsUsage: "<doc-id>",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "format", Value: "txt", Usage: "txt, html, pdf, docx, odt, rtf, epub"},
					&cli.StringFlag{Name: "output", Usage: "Output file path"},
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					body, _, _, err := c.RequestRaw(ctx, "GET", c.WriterBase+"/download/"+cmd.Args().First(), map[string]string{"format": cmd.String("format")})
					if err != nil {
						if strings.Contains(err.Error(), "R3002") {
							return output.JSON(map[string]string{"error": "Document is empty — Zoho cannot export empty documents (R3002)"})
						}
						return err
					}
					if out := cmd.String("output"); out != "" {
						if err := os.WriteFile(out, body, 0600); err != nil {
							return err
						}
						return output.JSON(map[string]any{"ok": true, "path": out, "size": len(body)})
					}
					os.Stdout.Write(body)
					return nil
				},
			},
		},
	}
}
