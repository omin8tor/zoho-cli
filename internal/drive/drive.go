package drive

import (
	"context"
	"fmt"
	"os"

	zohttp "github.com/omin8tor/zoho-cli/internal/http"
	"github.com/omin8tor/zoho-cli/internal/output"
	"github.com/omin8tor/zoho-cli/internal/pagination"
	"github.com/urfave/cli/v3"
)

const jsonapiCT = "application/vnd.api+json"

var roleIDs = map[string]int{
	"viewer":    7,
	"commenter": 6,
	"editor":    5,
	"organizer": 4,
}

var serviceTypeMap = map[string]string{
	"zohowriter": "zw",
	"zohosheet":  "zohosheet",
	"zohoshow":   "zohoshow",
	"writer":     "zw",
	"sheet":      "zohosheet",
	"show":       "zohoshow",
}

func requireTeam(cmd *cli.Command) (string, error) {
	v := cmd.String("team")
	if v == "" {
		return "", fmt.Errorf("--team is required (or set ZOHO_TEAM_ID env var)")
	}
	return v, nil
}

func jsonapiBody(attrs map[string]any) map[string]any {
	return map[string]any{
		"data": map[string]any{
			"type":       "files",
			"attributes": attrs,
		},
	}
}

func jsonapiHeaders() map[string]string {
	return map[string]string{"Content-Type": jsonapiCT}
}

func Commands() *cli.Command {
	return &cli.Command{
		Name:  "drive",
		Usage: "Zoho WorkDrive operations",
		Commands: []*cli.Command{
			filesCmd(),
			foldersCmd(),
			downloadCmd(),
			uploadCmd(),
			shareCmd(),
			teamsCmd(),
		},
	}
}

func filesCmd() *cli.Command {
	return &cli.Command{
		Name:  "files",
		Usage: "WorkDrive file operations",
		Commands: []*cli.Command{
			{
				Name:  "list",
				Usage: "List folder contents",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "folder", Required: true, Usage: "Folder ID"},
					&cli.StringFlag{Name: "type", Usage: "Filter: file, folder, image, etc."},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					url := c.WorkDriveBase + "/files/" + cmd.String("folder") + "/files"
					params := map[string]string{}
					if t := cmd.String("type"); t != "" {
						params["filter[type]"] = t
					}
					items, err := pagination.PaginateWorkDrive(c, url, params, 0)
					if err != nil {
						return err
					}
					return output.JSON(items)
				},
			},
			{
				Name:      "get",
				Usage:     "Get file info",
				ArgsUsage: "<file-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					raw, err := c.Request("GET", c.WorkDriveBase+"/files/"+cmd.Args().First(), nil)
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:  "search",
				Usage: "Search files",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "query", Required: true, Usage: "Search keyword"},
					&cli.StringFlag{Name: "team", Usage: "Team ID", Sources: cli.EnvVars("ZOHO_TEAM_ID")},
					&cli.StringFlag{Name: "mode", Value: "all", Usage: "all, name, or content"},
					&cli.StringFlag{Name: "type", Usage: "Filter by type"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					team, err := requireTeam(cmd)
					if err != nil {
						return err
					}
					url := c.WorkDriveBase + "/teams/" + team + "/records"
					mode := cmd.String("mode")
					params := map[string]string{"search[" + mode + "]": cmd.String("query")}
					if t := cmd.String("type"); t != "" {
						params["filter[type]"] = t
					}
					items, err := pagination.PaginateWorkDrive(c, url, params, 0)
					if err != nil {
						return err
					}
					return output.JSON(items)
				},
			},
			{
				Name:      "rename",
				Usage:     "Rename a file",
				ArgsUsage: "<file-id>",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "name", Required: true, Usage: "New name"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					body := jsonapiBody(map[string]any{"name": cmd.String("name")})
					raw, err := c.Request("PATCH", c.WorkDriveBase+"/files/"+cmd.Args().First(), &zohttp.RequestOpts{
						JSON:    body,
						Headers: jsonapiHeaders(),
					})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "copy",
				Usage:     "Copy a file",
				ArgsUsage: "<file-id>",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "to", Required: true, Usage: "Destination folder ID"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					body := jsonapiBody(map[string]any{"resource_id": cmd.Args().First()})
					raw, err := c.Request("POST", c.WorkDriveBase+"/files/"+cmd.String("to")+"/copy", &zohttp.RequestOpts{
						JSON:    body,
						Headers: jsonapiHeaders(),
					})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "move",
				Usage:     "Move a file to a different folder",
				ArgsUsage: "<file-id>",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "to", Required: true, Usage: "Destination folder ID"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					body := jsonapiBody(map[string]any{"parent_id": cmd.String("to")})
					raw, err := c.Request("PATCH", c.WorkDriveBase+"/files/"+cmd.Args().First(), &zohttp.RequestOpts{
						JSON:    body,
						Headers: jsonapiHeaders(),
					})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "trash",
				Usage:     "Move a file to trash",
				ArgsUsage: "<file-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					body := jsonapiBody(map[string]any{"status": 51})
					raw, err := c.Request("PATCH", c.WorkDriveBase+"/files/"+cmd.Args().First(), &zohttp.RequestOpts{
						JSON:    body,
						Headers: jsonapiHeaders(),
					})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "delete",
				Usage:     "Permanently delete a file",
				ArgsUsage: "<file-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					body := jsonapiBody(map[string]any{"status": 61})
					raw, err := c.Request("PATCH", c.WorkDriveBase+"/files/"+cmd.Args().First(), &zohttp.RequestOpts{
						JSON:    body,
						Headers: jsonapiHeaders(),
					})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "restore",
				Usage:     "Restore a file from trash",
				ArgsUsage: "<file-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					body := jsonapiBody(map[string]any{"status": 1})
					raw, err := c.Request("PATCH", c.WorkDriveBase+"/files/"+cmd.Args().First(), &zohttp.RequestOpts{
						JSON:    body,
						Headers: jsonapiHeaders(),
					})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:  "trash-list",
				Usage: "List trashed files in a team folder",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "team-folder", Required: true, Usage: "Team folder ID"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					url := c.WorkDriveBase + "/teamfolders/" + cmd.String("team-folder") + "/trashedfiles"
					items, err := pagination.PaginateWorkDrive(c, url, nil, 0)
					if err != nil {
						return err
					}
					return output.JSON(items)
				},
			},
			{
				Name:      "versions",
				Usage:     "List file versions",
				ArgsUsage: "<file-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					raw, err := c.Request("GET", c.WorkDriveBase+"/files/"+cmd.Args().First()+"/versions", nil)
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
		Usage: "WorkDrive folder operations",
		Commands: []*cli.Command{
			{
				Name:  "list",
				Usage: "List team folders",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "team", Usage: "Team ID", Sources: cli.EnvVars("ZOHO_TEAM_ID")},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					team, err := requireTeam(cmd)
					if err != nil {
						return err
					}
					url := c.WorkDriveBase + "/teams/" + team + "/teamfolders"
					items, err := pagination.PaginateWorkDrive(c, url, nil, 0)
					if err != nil {
						return err
					}
					return output.JSON(items)
				},
			},
			{
				Name:  "create",
				Usage: "Create a folder or document",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "name", Required: true, Usage: "Name"},
					&cli.StringFlag{Name: "parent", Required: true, Usage: "Parent folder ID"},
					&cli.StringFlag{Name: "type", Value: "folder", Usage: "folder, zohowriter, zohosheet, zohoshow"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					attrs := map[string]any{
						"name":      cmd.String("name"),
						"parent_id": cmd.String("parent"),
					}
					ft := cmd.String("type")
					if ft != "folder" {
						st := ft
						if mapped, ok := serviceTypeMap[ft]; ok {
							st = mapped
						}
						attrs["service_type"] = st
					}
					body := jsonapiBody(attrs)
					raw, err := c.Request("POST", c.WorkDriveBase+"/files", &zohttp.RequestOpts{
						JSON:    body,
						Headers: jsonapiHeaders(),
					})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "breadcrumb",
				Usage:     "Show folder path",
				ArgsUsage: "<folder-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					raw, err := c.Request("GET", c.WorkDriveBase+"/files/"+cmd.Args().First()+"/breadcrumbs", nil)
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
		},
	}
}

func downloadCmd() *cli.Command {
	return &cli.Command{
		Name:      "download",
		Usage:     "Download a file",
		ArgsUsage: "<file-id>",
		Flags: []cli.Flag{
			&cli.StringFlag{Name: "output", Usage: "Output file path"},
			&cli.StringFlag{Name: "format", Value: "native", Usage: "native, txt, html, pdf, docx"},
		},
		Action: func(_ context.Context, cmd *cli.Command) error {
			c, err := zohttp.GetClient()
			if err != nil {
				return err
			}
			fileID := cmd.Args().First()
			var url string
			var params map[string]string
			if f := cmd.String("format"); f != "native" {
				url = c.WriterBase + "/download/" + fileID
				params = map[string]string{"format": f}
			} else {
				url = c.WorkDriveBase + "/download/" + fileID
			}
			body, _, _, err := c.RequestRaw("GET", url, params)
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
	}
}

func uploadCmd() *cli.Command {
	return &cli.Command{
		Name:      "upload",
		Usage:     "Upload a file",
		ArgsUsage: "<file-path>",
		Flags: []cli.Flag{
			&cli.StringFlag{Name: "folder", Required: true, Usage: "Destination folder ID"},
			&cli.BoolFlag{Name: "override", Usage: "Override existing file"},
		},
		Action: func(_ context.Context, cmd *cli.Command) error {
			c, err := zohttp.GetClient()
			if err != nil {
				return err
			}
			filePath := cmd.Args().First()
			data, err := os.ReadFile(filePath)
			if err != nil {
				return err
			}
			name := filePath
			for i := len(filePath) - 1; i >= 0; i-- {
				if filePath[i] == '/' || filePath[i] == '\\' {
					name = filePath[i+1:]
					break
				}
			}
			form := map[string]string{
				"parent_id": cmd.String("folder"),
				"filename":  name,
			}
			if cmd.Bool("override") {
				form["override-name-exist"] = "true"
			}
			raw, err := c.Request("POST", c.WorkDriveBase+"/upload", &zohttp.RequestOpts{
				Files: map[string]zohttp.FileUpload{"content": {Filename: name, Data: data}},
				Form:  form,
			})
			if err != nil {
				return err
			}
			return output.JSONRaw(raw)
		},
	}
}

func shareCmd() *cli.Command {
	return &cli.Command{
		Name:  "share",
		Usage: "File sharing operations",
		Commands: []*cli.Command{
			{
				Name:      "permissions",
				Usage:     "List file permissions",
				ArgsUsage: "<file-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					raw, err := c.Request("GET", c.WorkDriveBase+"/files/"+cmd.Args().First()+"/permissions", nil)
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "add",
				Usage:     "Share a file",
				ArgsUsage: "<file-id>",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "email", Required: true, Usage: "Email to share with"},
					&cli.StringFlag{Name: "role", Value: "viewer", Usage: "viewer, commenter, editor, organizer"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					roleID := 7
					if r, ok := roleIDs[cmd.String("role")]; ok {
						roleID = r
					}
					body := map[string]any{
						"data": map[string]any{
							"type": "permissions",
							"attributes": map[string]any{
								"resource_id":            cmd.Args().First(),
								"shared_type":            "personal",
								"email_id":               cmd.String("email"),
								"role_id":                roleID,
								"send_notification_mail": true,
							},
						},
					}
					raw, err := c.Request("POST", c.WorkDriveBase+"/permissions", &zohttp.RequestOpts{
						JSON:    body,
						Headers: jsonapiHeaders(),
					})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "revoke",
				Usage:     "Revoke file access",
				ArgsUsage: "<permission-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					_, err = c.Request("DELETE", c.WorkDriveBase+"/permissions/"+cmd.Args().First(), nil)
					if err != nil {
						return err
					}
					return output.JSON(map[string]any{"ok": true, "revoked": cmd.Args().First()})
				},
			},
			{
				Name:      "links",
				Usage:     "List share links for a file",
				ArgsUsage: "<file-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					raw, err := c.Request("GET", c.WorkDriveBase+"/files/"+cmd.Args().First()+"/links", nil)
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "link",
				Usage:     "Create/get share link",
				ArgsUsage: "<file-id>",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "role", Value: "viewer", Usage: "viewer, commenter, editor"},
					&cli.BoolFlag{Name: "allow-download", Value: true},
					&cli.StringFlag{Name: "name", Usage: "Link display name"},
					&cli.StringFlag{Name: "expiration", Usage: "Expiry date YYYY-MM-DD"},
					&cli.StringFlag{Name: "password", Usage: "Link password"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					roleID := 7
					if r, ok := roleIDs[cmd.String("role")]; ok {
						roleID = r
					}
					attrs := map[string]any{
						"resource_id":       cmd.Args().First(),
						"role_id":           roleID,
						"allow_download":    cmd.Bool("allow-download"),
						"request_user_data": false,
					}
					if n := cmd.String("name"); n != "" {
						attrs["link_name"] = n
					}
					if e := cmd.String("expiration"); e != "" {
						attrs["expiration_date"] = e
					}
					if p := cmd.String("password"); p != "" {
						attrs["password_text"] = p
					}
					body := map[string]any{
						"data": map[string]any{
							"type":       "links",
							"attributes": attrs,
						},
					}
					raw, err := c.Request("POST", c.WorkDriveBase+"/links", &zohttp.RequestOpts{
						JSON:    body,
						Headers: jsonapiHeaders(),
					})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "unlink",
				Usage:     "Delete a share link",
				ArgsUsage: "<link-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					_, err = c.Request("DELETE", c.WorkDriveBase+"/links/"+cmd.Args().First(), nil)
					if err != nil {
						return err
					}
					return output.JSON(map[string]any{"ok": true, "deleted": cmd.Args().First()})
				},
			},
		},
	}
}

func teamsCmd() *cli.Command {
	return &cli.Command{
		Name:  "teams",
		Usage: "WorkDrive team operations",
		Commands: []*cli.Command{
			{
				Name:  "me",
				Usage: "Get current user info",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					raw, err := c.Request("GET", c.WorkDriveBase+"/users/me", nil)
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "members",
				Usage:     "List team members",
				ArgsUsage: "<team-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					raw, err := c.Request("GET", c.WorkDriveBase+"/teams/"+cmd.Args().First()+"/users", nil)
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
		},
	}
}
