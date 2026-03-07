package mail

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

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

func resolveAccountID(cmd *cli.Command) (string, error) {
	acct := cmd.String("account")
	if acct == "" {
		acct = os.Getenv("ZOHO_MAIL_ACCOUNT_ID")
	}
	if acct == "" {
		return "", internal.NewValidationError("--account flag or ZOHO_MAIL_ACCOUNT_ID env var required")
	}
	return acct, nil
}

func resolveOrgID(cmd *cli.Command) (string, error) {
	org := cmd.String("org")
	if org == "" {
		org = os.Getenv("ZOHO_MAIL_ORG_ID")
	}
	if org == "" {
		return "", internal.NewValidationError("--org flag or ZOHO_MAIL_ORG_ID env var required")
	}
	return org, nil
}

func Commands() *cli.Command {
	return &cli.Command{
		Name:  "mail",
		Usage: "Zoho Mail operations",
		Flags: []cli.Flag{
			&cli.StringFlag{Name: "account", Usage: "Mail Account ID (or set ZOHO_MAIL_ACCOUNT_ID)"},
			&cli.StringFlag{Name: "org", Usage: "Organization ID (or set ZOHO_MAIL_ORG_ID)"},
		},
		Commands: []*cli.Command{
			accountsCmd(),
			foldersCmd(),
			labelsCmd(),
			messagesCmd(),
			threadsCmd(),
			tasksCmd(),
			bookmarksCmd(),
			notesCmd(),
			organizationCmd(),
			domainsCmd(),
			groupsCmd(),
			usersCmd(),
			policyCmd(),
			logsCmd(),
			antispamCmd(),
			signaturesCmd(),
		},
	}
}

func accountsCmd() *cli.Command {
	return &cli.Command{
		Name:  "accounts",
		Usage: "Mail account operations",
		Commands: []*cli.Command{
			{
				Name:  "list",
				Usage: "List all mail accounts",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					raw, err := c.Request("GET", c.MailBase+"/api/accounts", nil)
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "get",
				Usage:     "Get a mail account",
				ArgsUsage: "<account-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					id := cmd.Args().Get(0)
					if id == "" {
						return internal.NewValidationError("account-id argument required")
					}
					c, err := getClient()
					if err != nil {
						return err
					}
					raw, err := c.Request("GET", c.MailBase+"/api/accounts/"+id, nil)
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
		Usage: "Mail folder operations",
		Commands: []*cli.Command{
			{
				Name:  "list",
				Usage: "List all folders",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					accountID, err := resolveAccountID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("GET", c.MailBase+"/api/accounts/"+accountID+"/folders", nil)
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "get",
				Usage:     "Get a folder",
				ArgsUsage: "<folder-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					id := cmd.Args().Get(0)
					if id == "" {
						return internal.NewValidationError("folder-id argument required")
					}
					c, err := getClient()
					if err != nil {
						return err
					}
					accountID, err := resolveAccountID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("GET", c.MailBase+"/api/accounts/"+accountID+"/folders/"+id, nil)
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
					&cli.StringFlag{Name: "json", Required: true, Usage: "JSON body"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					accountID, err := resolveAccountID(cmd)
					if err != nil {
						return err
					}
					var body any
					json.Unmarshal([]byte(cmd.String("json")), &body)
					raw, err := c.Request("POST", c.MailBase+"/api/accounts/"+accountID+"/folders", &zohttp.RequestOpts{JSON: body})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "update",
				Usage:     "Update a folder",
				ArgsUsage: "<folder-id>",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "json", Required: true, Usage: "JSON body"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					id := cmd.Args().Get(0)
					if id == "" {
						return internal.NewValidationError("folder-id argument required")
					}
					c, err := getClient()
					if err != nil {
						return err
					}
					accountID, err := resolveAccountID(cmd)
					if err != nil {
						return err
					}
					var body any
					json.Unmarshal([]byte(cmd.String("json")), &body)
					raw, err := c.Request("PUT", c.MailBase+"/api/accounts/"+accountID+"/folders/"+id, &zohttp.RequestOpts{JSON: body})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "delete",
				Usage:     "Delete a folder",
				ArgsUsage: "<folder-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					id := cmd.Args().Get(0)
					if id == "" {
						return internal.NewValidationError("folder-id argument required")
					}
					c, err := getClient()
					if err != nil {
						return err
					}
					accountID, err := resolveAccountID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("DELETE", c.MailBase+"/api/accounts/"+accountID+"/folders/"+id, nil)
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
		},
	}
}

func labelsCmd() *cli.Command {
	return &cli.Command{
		Name:  "labels",
		Usage: "Mail label operations",
		Commands: []*cli.Command{
			{
				Name:  "list",
				Usage: "List all labels",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					accountID, err := resolveAccountID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("GET", c.MailBase+"/api/accounts/"+accountID+"/labels", nil)
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "get",
				Usage:     "Get a label",
				ArgsUsage: "<label-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					id := cmd.Args().Get(0)
					if id == "" {
						return internal.NewValidationError("label-id argument required")
					}
					c, err := getClient()
					if err != nil {
						return err
					}
					accountID, err := resolveAccountID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("GET", c.MailBase+"/api/accounts/"+accountID+"/labels/"+id, nil)
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:  "create",
				Usage: "Create a label",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "json", Required: true, Usage: "JSON body"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					accountID, err := resolveAccountID(cmd)
					if err != nil {
						return err
					}
					var body any
					json.Unmarshal([]byte(cmd.String("json")), &body)
					raw, err := c.Request("POST", c.MailBase+"/api/accounts/"+accountID+"/labels", &zohttp.RequestOpts{JSON: body})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "update",
				Usage:     "Update a label",
				ArgsUsage: "<label-id>",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "json", Required: true, Usage: "JSON body"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					id := cmd.Args().Get(0)
					if id == "" {
						return internal.NewValidationError("label-id argument required")
					}
					c, err := getClient()
					if err != nil {
						return err
					}
					accountID, err := resolveAccountID(cmd)
					if err != nil {
						return err
					}
					var body any
					json.Unmarshal([]byte(cmd.String("json")), &body)
					raw, err := c.Request("PUT", c.MailBase+"/api/accounts/"+accountID+"/labels/"+id, &zohttp.RequestOpts{JSON: body})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "delete",
				Usage:     "Delete a label",
				ArgsUsage: "<label-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					id := cmd.Args().Get(0)
					if id == "" {
						return internal.NewValidationError("label-id argument required")
					}
					c, err := getClient()
					if err != nil {
						return err
					}
					accountID, err := resolveAccountID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("DELETE", c.MailBase+"/api/accounts/"+accountID+"/labels/"+id, nil)
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
		},
	}
}

func messagesCmd() *cli.Command {
	return &cli.Command{
		Name:  "messages",
		Usage: "Mail message operations",
		Commands: []*cli.Command{
			{
				Name:  "list",
				Usage: "List emails in a folder",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "folder", Required: true, Usage: "Folder ID"},
					&cli.StringFlag{Name: "start", Usage: "Starting index"},
					&cli.StringFlag{Name: "limit", Usage: "Max messages (max 200)"},
					&cli.StringFlag{Name: "status", Usage: "Filter: read, unread, or all"},
					&cli.StringFlag{Name: "flagid", Usage: "Filter by flag ID"},
					&cli.StringFlag{Name: "labelid", Usage: "Filter by label ID"},
					&cli.StringFlag{Name: "thread-id", Usage: "Filter by thread ID"},
					&cli.StringFlag{Name: "sort-by", Usage: "Sort field: date, messageId, size"},
					&cli.StringFlag{Name: "sort-order", Usage: "true (asc) or false (desc)"},
					&cli.StringFlag{Name: "include-to", Usage: "Include to address info (true/false)"},
					&cli.StringFlag{Name: "include-sent", Usage: "Include sent messages (true/false)"},
					&cli.StringFlag{Name: "include-archive", Usage: "Include archived messages (true/false)"},
					&cli.StringFlag{Name: "attached", Usage: "Filter attached mails (true/false)"},
					&cli.StringFlag{Name: "inlined", Usage: "Filter inlined mails (true/false)"},
					&cli.StringFlag{Name: "flagged", Usage: "Filter flagged mails (true/false)"},
					&cli.StringFlag{Name: "responded", Usage: "Filter responded mails (true/false)"},
					&cli.StringFlag{Name: "threaded", Usage: "Filter threaded mails (true/false)"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					accountID, err := resolveAccountID(cmd)
					if err != nil {
						return err
					}
					params := map[string]string{
						"folderId": cmd.String("folder"),
					}
					if v := cmd.String("start"); v != "" {
						params["start"] = v
					}
					if v := cmd.String("limit"); v != "" {
						params["limit"] = v
					}
					if v := cmd.String("status"); v != "" {
						params["status"] = v
					}
					if v := cmd.String("flagid"); v != "" {
						params["flagid"] = v
					}
					if v := cmd.String("labelid"); v != "" {
						params["labelid"] = v
					}
					if v := cmd.String("thread-id"); v != "" {
						params["threadId"] = v
					}
					if v := cmd.String("sort-by"); v != "" {
						params["sortBy"] = v
					}
					if v := cmd.String("sort-order"); v != "" {
						params["sortorder"] = v
					}
					if v := cmd.String("include-to"); v != "" {
						params["includeto"] = v
					}
					if v := cmd.String("include-sent"); v != "" {
						params["includesent"] = v
					}
					if v := cmd.String("include-archive"); v != "" {
						params["includearchive"] = v
					}
					if v := cmd.String("attached"); v != "" {
						params["attachedMails"] = v
					}
					if v := cmd.String("inlined"); v != "" {
						params["inlinedMails"] = v
					}
					if v := cmd.String("flagged"); v != "" {
						params["flaggedMails"] = v
					}
					if v := cmd.String("responded"); v != "" {
						params["respondedMails"] = v
					}
					if v := cmd.String("threaded"); v != "" {
						params["threadedMails"] = v
					}
					raw, err := c.Request("GET", c.MailBase+"/api/accounts/"+accountID+"/messages/view", &zohttp.RequestOpts{Params: params})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:  "search",
				Usage: "Search emails",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "query", Required: true, Usage: "Search query"},
					&cli.StringFlag{Name: "start", Usage: "Starting index"},
					&cli.StringFlag{Name: "limit", Usage: "Max results"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					accountID, err := resolveAccountID(cmd)
					if err != nil {
						return err
					}
					params := map[string]string{
						"searchKey": cmd.String("query"),
					}
					if v := cmd.String("start"); v != "" {
						params["start"] = v
					}
					if v := cmd.String("limit"); v != "" {
						params["limit"] = v
					}
					raw, err := c.Request("GET", c.MailBase+"/api/accounts/"+accountID+"/messages/search", &zohttp.RequestOpts{Params: params})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "get",
				Usage:     "Get email details",
				ArgsUsage: "<message-id>",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "folder", Required: true, Usage: "Folder ID"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					msgID := cmd.Args().Get(0)
					if msgID == "" {
						return internal.NewValidationError("message-id argument required")
					}
					c, err := getClient()
					if err != nil {
						return err
					}
					accountID, err := resolveAccountID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("GET", c.MailBase+"/api/accounts/"+accountID+"/folders/"+cmd.String("folder")+"/messages/"+msgID+"/details", nil)
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "content",
				Usage:     "Get email body content",
				ArgsUsage: "<message-id>",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "folder", Required: true, Usage: "Folder ID"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					msgID := cmd.Args().Get(0)
					if msgID == "" {
						return internal.NewValidationError("message-id argument required")
					}
					c, err := getClient()
					if err != nil {
						return err
					}
					accountID, err := resolveAccountID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("GET", c.MailBase+"/api/accounts/"+accountID+"/folders/"+cmd.String("folder")+"/messages/"+msgID+"/content", nil)
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "headers",
				Usage:     "Get email headers",
				ArgsUsage: "<message-id>",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "folder", Required: true, Usage: "Folder ID"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					msgID := cmd.Args().Get(0)
					if msgID == "" {
						return internal.NewValidationError("message-id argument required")
					}
					c, err := getClient()
					if err != nil {
						return err
					}
					accountID, err := resolveAccountID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("GET", c.MailBase+"/api/accounts/"+accountID+"/folders/"+cmd.String("folder")+"/messages/"+msgID+"/header", nil)
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "original",
				Usage:     "Get original email message",
				ArgsUsage: "<message-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					msgID := cmd.Args().Get(0)
					if msgID == "" {
						return internal.NewValidationError("message-id argument required")
					}
					c, err := getClient()
					if err != nil {
						return err
					}
					accountID, err := resolveAccountID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("GET", c.MailBase+"/api/accounts/"+accountID+"/messages/"+msgID+"/originalmessage", nil)
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "attachment-info",
				Usage:     "Get attachment info for a message",
				ArgsUsage: "<message-id>",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "folder", Required: true, Usage: "Folder ID"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					msgID := cmd.Args().Get(0)
					if msgID == "" {
						return internal.NewValidationError("message-id argument required")
					}
					c, err := getClient()
					if err != nil {
						return err
					}
					accountID, err := resolveAccountID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("GET", c.MailBase+"/api/accounts/"+accountID+"/folders/"+cmd.String("folder")+"/messages/"+msgID+"/attachmentinfo", nil)
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "attachment",
				Usage:     "Download a message attachment",
				ArgsUsage: "<message-id> <attachment-id>",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "folder", Required: true, Usage: "Folder ID"},
					&cli.StringFlag{Name: "output", Usage: "Output file path"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					msgID := cmd.Args().Get(0)
					attachID := cmd.Args().Get(1)
					if msgID == "" || attachID == "" {
						return internal.NewValidationError("message-id and attachment-id arguments required")
					}
					c, err := getClient()
					if err != nil {
						return err
					}
					accountID, err := resolveAccountID(cmd)
					if err != nil {
						return err
					}
					body, _, _, err := c.RequestRaw("GET", c.MailBase+"/api/accounts/"+accountID+"/folders/"+cmd.String("folder")+"/messages/"+msgID+"/attachments/"+attachID, nil)
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
				Name:  "send",
				Usage: "Send an email",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "from", Required: true, Usage: "Sender email address"},
					&cli.StringFlag{Name: "to", Required: true, Usage: "Recipient email address(es)"},
					&cli.StringFlag{Name: "cc", Usage: "CC address(es)"},
					&cli.StringFlag{Name: "bcc", Usage: "BCC address(es)"},
					&cli.StringFlag{Name: "subject", Usage: "Email subject"},
					&cli.StringFlag{Name: "content", Usage: "Email body content"},
					&cli.StringFlag{Name: "format", Usage: "Mail format: html or plaintext"},
					&cli.StringFlag{Name: "ask-receipt", Usage: "Request read receipt: yes or no"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					accountID, err := resolveAccountID(cmd)
					if err != nil {
						return err
					}
					body := map[string]any{
						"fromAddress": cmd.String("from"),
						"toAddress":   cmd.String("to"),
					}
					if v := cmd.String("cc"); v != "" {
						body["ccAddress"] = v
					}
					if v := cmd.String("bcc"); v != "" {
						body["bccAddress"] = v
					}
					if v := cmd.String("subject"); v != "" {
						body["subject"] = v
					}
					if v := cmd.String("content"); v != "" {
						body["content"] = v
					}
					if v := cmd.String("format"); v != "" {
						body["mailFormat"] = v
					}
					if v := cmd.String("ask-receipt"); v != "" {
						body["askReceipt"] = v
					}
					raw, err := c.Request("POST", c.MailBase+"/api/accounts/"+accountID+"/messages", &zohttp.RequestOpts{JSON: body})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "reply",
				Usage:     "Reply to an email",
				ArgsUsage: "<message-id>",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "json", Required: true, Usage: "JSON body"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					msgID := cmd.Args().Get(0)
					if msgID == "" {
						return internal.NewValidationError("message-id argument required")
					}
					c, err := getClient()
					if err != nil {
						return err
					}
					accountID, err := resolveAccountID(cmd)
					if err != nil {
						return err
					}
					var body any
					json.Unmarshal([]byte(cmd.String("json")), &body)
					raw, err := c.Request("POST", c.MailBase+"/api/accounts/"+accountID+"/messages/"+msgID, &zohttp.RequestOpts{JSON: body})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:  "update",
				Usage: "Update message status (read/unread, move, flag, label, archive, spam)",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "json", Required: true, Usage: "JSON body"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					accountID, err := resolveAccountID(cmd)
					if err != nil {
						return err
					}
					var body any
					json.Unmarshal([]byte(cmd.String("json")), &body)
					raw, err := c.Request("PUT", c.MailBase+"/api/accounts/"+accountID+"/updatemessage", &zohttp.RequestOpts{JSON: body})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:  "upload-attachment",
				Usage: "Upload an attachment for composing",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "file", Required: true, Usage: "Path to file"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					accountID, err := resolveAccountID(cmd)
					if err != nil {
						return err
					}
					filePath := cmd.String("file")
					fileData, err := os.ReadFile(filePath)
					if err != nil {
						return fmt.Errorf("failed to read file: %w", err)
					}
					raw, err := c.Request("POST", c.MailBase+"/api/accounts/"+accountID+"/messages/attachments", &zohttp.RequestOpts{
						Files: map[string]zohttp.FileUpload{"file": {Filename: filepath.Base(filePath), Data: fileData}},
					})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "delete",
				Usage:     "Delete an email",
				ArgsUsage: "<message-id>",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "folder", Required: true, Usage: "Folder ID"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					msgID := cmd.Args().Get(0)
					if msgID == "" {
						return internal.NewValidationError("message-id argument required")
					}
					c, err := getClient()
					if err != nil {
						return err
					}
					accountID, err := resolveAccountID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("DELETE", c.MailBase+"/api/accounts/"+accountID+"/folders/"+cmd.String("folder")+"/messages/"+msgID, nil)
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
		},
	}
}

func threadsCmd() *cli.Command {
	return &cli.Command{
		Name:  "threads",
		Usage: "Mail thread operations",
		Commands: []*cli.Command{
			{
				Name:  "update",
				Usage: "Update thread status (read/unread, move, flag, label, spam)",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "json", Required: true, Usage: "JSON body"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					accountID, err := resolveAccountID(cmd)
					if err != nil {
						return err
					}
					var body any
					json.Unmarshal([]byte(cmd.String("json")), &body)
					raw, err := c.Request("PUT", c.MailBase+"/api/accounts/"+accountID+"/updatethread", &zohttp.RequestOpts{JSON: body})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
		},
	}
}

func signaturesCmd() *cli.Command {
	return &cli.Command{
		Name:  "signatures",
		Usage: "Mail signature operations",
		Commands: []*cli.Command{
			{
				Name:  "get",
				Usage: "Get signatures",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					raw, err := c.Request("GET", c.MailBase+"/api/accounts/signature", nil)
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:  "create",
				Usage: "Create a signature",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "json", Required: true, Usage: "JSON body"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					var body any
					json.Unmarshal([]byte(cmd.String("json")), &body)
					raw, err := c.Request("POST", c.MailBase+"/api/accounts/signature", &zohttp.RequestOpts{JSON: body})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:  "update",
				Usage: "Update a signature",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "json", Required: true, Usage: "JSON body"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					var body any
					json.Unmarshal([]byte(cmd.String("json")), &body)
					raw, err := c.Request("PUT", c.MailBase+"/api/accounts/signature", &zohttp.RequestOpts{JSON: body})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:  "delete",
				Usage: "Delete a signature",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					raw, err := c.Request("DELETE", c.MailBase+"/api/accounts/signature", nil)
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
		},
	}
}

func organizationCmd() *cli.Command {
	return &cli.Command{
		Name:  "organization",
		Usage: "Organization admin operations",
		Commands: []*cli.Command{
			{
				Name:  "get",
				Usage: "Get organization details",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("GET", c.MailBase+"/api/organization/"+orgID, nil)
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:  "storage",
				Usage: "Get organization storage info",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("GET", c.MailBase+"/api/organization/"+orgID+"/storage", nil)
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "user-storage",
				Usage:     "Get storage info for a user",
				ArgsUsage: "<zuid>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					zuid := cmd.Args().Get(0)
					if zuid == "" {
						return internal.NewValidationError("zuid argument required")
					}
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("GET", c.MailBase+"/api/organization/"+orgID+"/storage/"+zuid, nil)
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "update-user-storage",
				Usage:     "Update storage for a user",
				ArgsUsage: "<zuid>",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "json", Required: true, Usage: "JSON body"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					zuid := cmd.Args().Get(0)
					if zuid == "" {
						return internal.NewValidationError("zuid argument required")
					}
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
					raw, err := c.Request("PUT", c.MailBase+"/api/organization/"+orgID+"/storage/"+zuid, &zohttp.RequestOpts{JSON: body})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:  "spam-listing",
				Usage: "Get organization spam listing data",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("GET", c.MailBase+"/api/organization/"+orgID+"/antispam/data", nil)
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:  "update-spam-listing",
				Usage: "Update organization spam listing data",
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
					raw, err := c.Request("PUT", c.MailBase+"/api/organization/"+orgID+"/antispam/data", &zohttp.RequestOpts{JSON: body})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:  "delete-spam-listing",
				Usage: "Delete organization spam listing data",
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
					raw, err := c.Request("DELETE", c.MailBase+"/api/organization/"+orgID+"/antispam/data", &zohttp.RequestOpts{JSON: body})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:  "allowed-ips",
				Usage: "Get allowed IPs",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("GET", c.MailBase+"/api/organization/"+orgID+"/allowedIps", nil)
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:  "add-allowed-ips",
				Usage: "Add allowed IPs",
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
					raw, err := c.Request("POST", c.MailBase+"/api/organization/"+orgID+"/allowedIps", &zohttp.RequestOpts{JSON: body})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:  "delete-allowed-ips",
				Usage: "Delete allowed IPs",
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
					raw, err := c.Request("DELETE", c.MailBase+"/api/organization/"+orgID+"/allowedIps", &zohttp.RequestOpts{JSON: body})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
		},
	}
}

func domainsCmd() *cli.Command {
	return &cli.Command{
		Name:  "domains",
		Usage: "Domain admin operations",
		Commands: []*cli.Command{
			{
				Name:  "list",
				Usage: "List domains",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("GET", c.MailBase+"/api/organization/"+orgID+"/domains", nil)
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "get",
				Usage:     "Get a domain",
				ArgsUsage: "<domain-name>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					name := cmd.Args().Get(0)
					if name == "" {
						return internal.NewValidationError("domain-name argument required")
					}
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("GET", c.MailBase+"/api/organization/"+orgID+"/domains/"+name, nil)
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:  "add",
				Usage: "Add a domain",
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
					raw, err := c.Request("POST", c.MailBase+"/api/organization/"+orgID+"/domains", &zohttp.RequestOpts{JSON: body})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "update",
				Usage:     "Update a domain",
				ArgsUsage: "<domain-name>",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "json", Required: true, Usage: "JSON body"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					name := cmd.Args().Get(0)
					if name == "" {
						return internal.NewValidationError("domain-name argument required")
					}
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
					raw, err := c.Request("PUT", c.MailBase+"/api/organization/"+orgID+"/domains/"+name, &zohttp.RequestOpts{JSON: body})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "delete",
				Usage:     "Delete a domain",
				ArgsUsage: "<domain-name>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					name := cmd.Args().Get(0)
					if name == "" {
						return internal.NewValidationError("domain-name argument required")
					}
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("DELETE", c.MailBase+"/api/organization/"+orgID+"/domains/"+name, nil)
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
		},
	}
}

func groupsCmd() *cli.Command {
	return &cli.Command{
		Name:  "groups",
		Usage: "Mail group admin operations",
		Commands: []*cli.Command{
			{
				Name:  "list",
				Usage: "List groups",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("GET", c.MailBase+"/api/organization/"+orgID+"/groups", nil)
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "get",
				Usage:     "Get a group",
				ArgsUsage: "<group-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					id := cmd.Args().Get(0)
					if id == "" {
						return internal.NewValidationError("group-id argument required")
					}
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("GET", c.MailBase+"/api/organization/"+orgID+"/groups/"+id, nil)
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:  "create",
				Usage: "Create a group",
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
					raw, err := c.Request("POST", c.MailBase+"/api/organization/"+orgID+"/groups", &zohttp.RequestOpts{JSON: body})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "update",
				Usage:     "Update a group",
				ArgsUsage: "<group-id>",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "json", Required: true, Usage: "JSON body"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					id := cmd.Args().Get(0)
					if id == "" {
						return internal.NewValidationError("group-id argument required")
					}
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
					raw, err := c.Request("PUT", c.MailBase+"/api/organization/"+orgID+"/groups/"+id, &zohttp.RequestOpts{JSON: body})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "delete",
				Usage:     "Delete a group",
				ArgsUsage: "<group-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					id := cmd.Args().Get(0)
					if id == "" {
						return internal.NewValidationError("group-id argument required")
					}
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("DELETE", c.MailBase+"/api/organization/"+orgID+"/groups/"+id, nil)
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "moderation-emails",
				Usage:     "Get moderation emails for a group",
				ArgsUsage: "<group-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					id := cmd.Args().Get(0)
					if id == "" {
						return internal.NewValidationError("group-id argument required")
					}
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("GET", c.MailBase+"/api/organization/"+orgID+"/groups/"+id+"/messages", nil)
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "moderation-email-content",
				Usage:     "Get content of a moderation email",
				ArgsUsage: "<group-id> <message-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					groupID := cmd.Args().Get(0)
					msgID := cmd.Args().Get(1)
					if groupID == "" || msgID == "" {
						return internal.NewValidationError("group-id and message-id arguments required")
					}
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("GET", c.MailBase+"/api/organization/"+orgID+"/groups/"+groupID+"/messages/"+msgID, nil)
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "moderate",
				Usage:     "Moderate group messages",
				ArgsUsage: "<group-id>",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "json", Required: true, Usage: "JSON body"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					id := cmd.Args().Get(0)
					if id == "" {
						return internal.NewValidationError("group-id argument required")
					}
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
					raw, err := c.Request("PUT", c.MailBase+"/api/organization/"+orgID+"/groups/"+id+"/messages", &zohttp.RequestOpts{JSON: body})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
		},
	}
}

func usersCmd() *cli.Command {
	return &cli.Command{
		Name:  "users",
		Usage: "Organization user admin operations",
		Commands: []*cli.Command{
			{
				Name:  "list",
				Usage: "List organization users",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("GET", c.MailBase+"/api/organization/"+orgID+"/accounts/", nil)
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "get",
				Usage:     "Get a user",
				ArgsUsage: "<zuid>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					zuid := cmd.Args().Get(0)
					if zuid == "" {
						return internal.NewValidationError("zuid argument required")
					}
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("GET", c.MailBase+"/api/organization/"+orgID+"/accounts/"+zuid, nil)
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:  "add",
				Usage: "Add a user",
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
					raw, err := c.Request("POST", c.MailBase+"/api/organization/"+orgID+"/accounts/", &zohttp.RequestOpts{JSON: body})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "update",
				Usage:     "Update a user",
				ArgsUsage: "<zuid>",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "json", Required: true, Usage: "JSON body"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					zuid := cmd.Args().Get(0)
					if zuid == "" {
						return internal.NewValidationError("zuid argument required")
					}
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
					raw, err := c.Request("PUT", c.MailBase+"/api/organization/"+orgID+"/accounts/"+zuid, &zohttp.RequestOpts{JSON: body})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:  "delete",
				Usage: "Delete users",
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
					raw, err := c.Request("DELETE", c.MailBase+"/api/organization/"+orgID+"/accounts", &zohttp.RequestOpts{JSON: body})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
		},
	}
}

func policyCmd() *cli.Command {
	return &cli.Command{
		Name:  "policy",
		Usage: "Mail policy operations",
		Commands: []*cli.Command{
			{
				Name:  "list",
				Usage: "List policies",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("GET", c.MailBase+"/api/organization/"+orgID+"/policy", nil)
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:  "email-restrictions",
				Usage: "Get email restriction policy",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("GET", c.MailBase+"/api/organization/"+orgID+"/policy/mailRestriction", nil)
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:  "account-restrictions",
				Usage: "Get account restriction policy",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("GET", c.MailBase+"/api/organization/"+orgID+"/policy/accountRestriction", nil)
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:  "access-restrictions",
				Usage: "Get access restriction policy",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("GET", c.MailBase+"/api/organization/"+orgID+"/policy/accessRestriction", nil)
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:  "forward-restrictions",
				Usage: "Get mail forward policy",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("GET", c.MailBase+"/api/organization/"+orgID+"/policy/mailForwardPolicy", nil)
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "users",
				Usage:     "Get users assigned to a policy",
				ArgsUsage: "<policy-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					id := cmd.Args().Get(0)
					if id == "" {
						return internal.NewValidationError("policy-id argument required")
					}
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("GET", c.MailBase+"/api/organization/"+orgID+"/policy/"+id+"/getUsers", nil)
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "groups",
				Usage:     "Get groups assigned to a policy",
				ArgsUsage: "<policy-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					id := cmd.Args().Get(0)
					if id == "" {
						return internal.NewValidationError("policy-id argument required")
					}
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("GET", c.MailBase+"/api/organization/"+orgID+"/policy/"+id+"/getGroups", nil)
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:  "create",
				Usage: "Create a policy",
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
					raw, err := c.Request("POST", c.MailBase+"/api/organization/"+orgID+"/policy", &zohttp.RequestOpts{JSON: body})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "update",
				Usage:     "Update a policy",
				ArgsUsage: "<policy-id>",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "json", Required: true, Usage: "JSON body"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					id := cmd.Args().Get(0)
					if id == "" {
						return internal.NewValidationError("policy-id argument required")
					}
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
					raw, err := c.Request("PUT", c.MailBase+"/api/organization/"+orgID+"/policy/"+id, &zohttp.RequestOpts{JSON: body})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
		},
	}
}

func logsCmd() *cli.Command {
	return &cli.Command{
		Name:  "logs",
		Usage: "Mail log operations",
		Commands: []*cli.Command{
			{
				Name:  "login-history",
				Usage: "Get login history",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("GET", c.MailBase+"/api/organization/"+orgID+"/accounts/reports/loginHistory", nil)
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:  "audit",
				Usage: "Get audit logs",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("GET", c.MailBase+"/api/organization/"+orgID+"/activity", nil)
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:  "smtp",
				Usage: "Get SMTP logs",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("GET", c.MailBase+"/api/organization/"+orgID+"/smtplogs", nil)
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
		},
	}
}

func antispamCmd() *cli.Command {
	return &cli.Command{
		Name:  "antispam",
		Usage: "Anti-spam operations",
		Commands: []*cli.Command{
			{
				Name:  "options",
				Usage: "Get anti-spam options",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := c.Request("GET", c.MailBase+"/api/organization/"+orgID+"/antispam/options", nil)
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:  "update",
				Usage: "Update anti-spam options",
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
					raw, err := c.Request("PUT", c.MailBase+"/api/organization/"+orgID+"/antispam/options", &zohttp.RequestOpts{JSON: body})
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
		Usage: "Mail task operations",
		Commands: []*cli.Command{
			{
				Name:  "list-assigned",
				Usage: "List tasks assigned to me",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					raw, err := c.Request("GET", c.MailBase+"/api/tasks", &zohttp.RequestOpts{Params: map[string]string{"assignedTo": "me"}})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:  "list-created",
				Usage: "List tasks created by me",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					raw, err := c.Request("GET", c.MailBase+"/api/tasks", &zohttp.RequestOpts{Params: map[string]string{"createdBy": "me"}})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:  "list-personal",
				Usage: "List personal tasks",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					raw, err := c.Request("GET", c.MailBase+"/api/tasks/me", nil)
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "list-group",
				Usage:     "List group tasks",
				ArgsUsage: "<group-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					id := cmd.Args().Get(0)
					if id == "" {
						return internal.NewValidationError("group-id argument required")
					}
					c, err := getClient()
					if err != nil {
						return err
					}
					raw, err := c.Request("GET", c.MailBase+"/api/tasks/groups/"+id, nil)
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "get-personal",
				Usage:     "Get a personal task",
				ArgsUsage: "<task-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					id := cmd.Args().Get(0)
					if id == "" {
						return internal.NewValidationError("task-id argument required")
					}
					c, err := getClient()
					if err != nil {
						return err
					}
					raw, err := c.Request("GET", c.MailBase+"/api/tasks/me/"+id, nil)
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "get-group",
				Usage:     "Get a group task",
				ArgsUsage: "<group-id> <task-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					groupID := cmd.Args().Get(0)
					taskID := cmd.Args().Get(1)
					if groupID == "" || taskID == "" {
						return internal.NewValidationError("group-id and task-id arguments required")
					}
					c, err := getClient()
					if err != nil {
						return err
					}
					raw, err := c.Request("GET", c.MailBase+"/api/tasks/groups/"+groupID+"/"+taskID, nil)
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:  "create-personal",
				Usage: "Create a personal task",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "json", Required: true, Usage: "JSON body"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					var body any
					json.Unmarshal([]byte(cmd.String("json")), &body)
					raw, err := c.Request("POST", c.MailBase+"/api/tasks/me", &zohttp.RequestOpts{JSON: body})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "create-group",
				Usage:     "Create a group task",
				ArgsUsage: "<group-id>",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "json", Required: true, Usage: "JSON body"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					id := cmd.Args().Get(0)
					if id == "" {
						return internal.NewValidationError("group-id argument required")
					}
					c, err := getClient()
					if err != nil {
						return err
					}
					var body any
					json.Unmarshal([]byte(cmd.String("json")), &body)
					raw, err := c.Request("POST", c.MailBase+"/api/tasks/groups/"+id, &zohttp.RequestOpts{JSON: body})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "update-personal",
				Usage:     "Update a personal task",
				ArgsUsage: "<task-id>",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "json", Required: true, Usage: "JSON body"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					id := cmd.Args().Get(0)
					if id == "" {
						return internal.NewValidationError("task-id argument required")
					}
					c, err := getClient()
					if err != nil {
						return err
					}
					var body any
					json.Unmarshal([]byte(cmd.String("json")), &body)
					raw, err := c.Request("PUT", c.MailBase+"/api/tasks/me/"+id, &zohttp.RequestOpts{JSON: body})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "update-group",
				Usage:     "Update a group task",
				ArgsUsage: "<group-id> <task-id>",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "json", Required: true, Usage: "JSON body"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					groupID := cmd.Args().Get(0)
					taskID := cmd.Args().Get(1)
					if groupID == "" || taskID == "" {
						return internal.NewValidationError("group-id and task-id arguments required")
					}
					c, err := getClient()
					if err != nil {
						return err
					}
					var body any
					json.Unmarshal([]byte(cmd.String("json")), &body)
					raw, err := c.Request("PUT", c.MailBase+"/api/tasks/groups/"+groupID+"/"+taskID, &zohttp.RequestOpts{JSON: body})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "delete-personal",
				Usage:     "Delete a personal task",
				ArgsUsage: "<task-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					id := cmd.Args().Get(0)
					if id == "" {
						return internal.NewValidationError("task-id argument required")
					}
					c, err := getClient()
					if err != nil {
						return err
					}
					raw, err := c.Request("DELETE", c.MailBase+"/api/tasks/me/"+id, nil)
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "delete-group",
				Usage:     "Delete a group task",
				ArgsUsage: "<group-id> <task-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					groupID := cmd.Args().Get(0)
					taskID := cmd.Args().Get(1)
					if groupID == "" || taskID == "" {
						return internal.NewValidationError("group-id and task-id arguments required")
					}
					c, err := getClient()
					if err != nil {
						return err
					}
					raw, err := c.Request("DELETE", c.MailBase+"/api/tasks/groups/"+groupID+"/"+taskID, nil)
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "subtasks-personal",
				Usage:     "List subtasks of a personal task",
				ArgsUsage: "<task-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					id := cmd.Args().Get(0)
					if id == "" {
						return internal.NewValidationError("task-id argument required")
					}
					c, err := getClient()
					if err != nil {
						return err
					}
					raw, err := c.Request("GET", c.MailBase+"/api/tasks/me/"+id+"/subtasks", nil)
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "subtasks-group",
				Usage:     "List subtasks of a group task",
				ArgsUsage: "<group-id> <task-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					groupID := cmd.Args().Get(0)
					taskID := cmd.Args().Get(1)
					if groupID == "" || taskID == "" {
						return internal.NewValidationError("group-id and task-id arguments required")
					}
					c, err := getClient()
					if err != nil {
						return err
					}
					raw, err := c.Request("GET", c.MailBase+"/api/tasks/groups/"+groupID+"/"+taskID+"/subtasks", nil)
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:  "task-groups",
				Usage: "List task groups",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					raw, err := c.Request("GET", c.MailBase+"/api/tasks/groups", nil)
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "group-members",
				Usage:     "List members of a task group",
				ArgsUsage: "<group-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					id := cmd.Args().Get(0)
					if id == "" {
						return internal.NewValidationError("group-id argument required")
					}
					c, err := getClient()
					if err != nil {
						return err
					}
					raw, err := c.Request("GET", c.MailBase+"/api/tasks/groups/"+id+"/members", nil)
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "projects",
				Usage:     "List projects in a task group",
				ArgsUsage: "<group-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					id := cmd.Args().Get(0)
					if id == "" {
						return internal.NewValidationError("group-id argument required")
					}
					c, err := getClient()
					if err != nil {
						return err
					}
					raw, err := c.Request("GET", c.MailBase+"/api/tasks/groups/"+id+"/projects", nil)
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "project-tasks",
				Usage:     "List tasks in a project",
				ArgsUsage: "<group-id> <project-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					groupID := cmd.Args().Get(0)
					projectID := cmd.Args().Get(1)
					if groupID == "" || projectID == "" {
						return internal.NewValidationError("group-id and project-id arguments required")
					}
					c, err := getClient()
					if err != nil {
						return err
					}
					raw, err := c.Request("GET", c.MailBase+"/api/tasks/groups/"+groupID+"/projects/"+projectID, nil)
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "create-project",
				Usage:     "Create a project in a task group",
				ArgsUsage: "<group-id>",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "json", Required: true, Usage: "JSON body"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					id := cmd.Args().Get(0)
					if id == "" {
						return internal.NewValidationError("group-id argument required")
					}
					c, err := getClient()
					if err != nil {
						return err
					}
					var body any
					json.Unmarshal([]byte(cmd.String("json")), &body)
					raw, err := c.Request("POST", c.MailBase+"/api/tasks/groups/"+id+"/projects", &zohttp.RequestOpts{JSON: body})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "update-project",
				Usage:     "Update a project in a task group",
				ArgsUsage: "<group-id> <project-id>",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "json", Required: true, Usage: "JSON body"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					groupID := cmd.Args().Get(0)
					projectID := cmd.Args().Get(1)
					if groupID == "" || projectID == "" {
						return internal.NewValidationError("group-id and project-id arguments required")
					}
					c, err := getClient()
					if err != nil {
						return err
					}
					var body any
					json.Unmarshal([]byte(cmd.String("json")), &body)
					raw, err := c.Request("PUT", c.MailBase+"/api/tasks/groups/"+groupID+"/projects/"+projectID, &zohttp.RequestOpts{JSON: body})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "delete-project",
				Usage:     "Delete a project from a task group",
				ArgsUsage: "<group-id> <project-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					groupID := cmd.Args().Get(0)
					projectID := cmd.Args().Get(1)
					if groupID == "" || projectID == "" {
						return internal.NewValidationError("group-id and project-id arguments required")
					}
					c, err := getClient()
					if err != nil {
						return err
					}
					raw, err := c.Request("DELETE", c.MailBase+"/api/tasks/groups/"+groupID+"/projects/"+projectID, nil)
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
		},
	}
}

func bookmarksCmd() *cli.Command {
	return &cli.Command{
		Name:  "bookmarks",
		Usage: "Mail bookmark operations",
		Commands: []*cli.Command{
			{Name: "list-personal", Usage: "List personal bookmarks", Action: func(_ context.Context, cmd *cli.Command) error { c, err := getClient(); if err != nil { return err }; raw, err := c.Request("GET", c.MailBase+"/api/links/me", nil); if err != nil { return err }; return output.JSONRaw(raw) }},
			{Name: "list-group", Usage: "List group bookmarks", ArgsUsage: "<group-id>", Action: func(_ context.Context, cmd *cli.Command) error { id := cmd.Args().Get(0); if id == "" { return internal.NewValidationError("group-id argument required") }; c, err := getClient(); if err != nil { return err }; raw, err := c.Request("GET", c.MailBase+"/api/links/groups/"+id, nil); if err != nil { return err }; return output.JSONRaw(raw) }},
			{Name: "get-personal", Usage: "Get a personal bookmark", ArgsUsage: "<bookmark-id>", Action: func(_ context.Context, cmd *cli.Command) error { id := cmd.Args().Get(0); if id == "" { return internal.NewValidationError("bookmark-id argument required") }; c, err := getClient(); if err != nil { return err }; raw, err := c.Request("GET", c.MailBase+"/api/links/me/"+id, nil); if err != nil { return err }; return output.JSONRaw(raw) }},
			{Name: "get-group", Usage: "Get a group bookmark", ArgsUsage: "<group-id> <bookmark-id>", Action: func(_ context.Context, cmd *cli.Command) error { g := cmd.Args().Get(0); b := cmd.Args().Get(1); if g == "" || b == "" { return internal.NewValidationError("group-id and bookmark-id arguments required") }; c, err := getClient(); if err != nil { return err }; raw, err := c.Request("GET", c.MailBase+"/api/links/groups/"+g+"/"+b, nil); if err != nil { return err }; return output.JSONRaw(raw) }},
			{Name: "create-personal", Usage: "Create a personal bookmark", Flags: []cli.Flag{&cli.StringFlag{Name: "json", Required: true, Usage: "JSON body"}}, Action: func(_ context.Context, cmd *cli.Command) error { c, err := getClient(); if err != nil { return err }; var body any; json.Unmarshal([]byte(cmd.String("json")), &body); raw, err := c.Request("POST", c.MailBase+"/api/links/me", &zohttp.RequestOpts{JSON: body}); if err != nil { return err }; return output.JSONRaw(raw) }},
			{Name: "create-group", Usage: "Create a group bookmark", ArgsUsage: "<group-id>", Flags: []cli.Flag{&cli.StringFlag{Name: "json", Required: true, Usage: "JSON body"}}, Action: func(_ context.Context, cmd *cli.Command) error { id := cmd.Args().Get(0); if id == "" { return internal.NewValidationError("group-id argument required") }; c, err := getClient(); if err != nil { return err }; var body any; json.Unmarshal([]byte(cmd.String("json")), &body); raw, err := c.Request("POST", c.MailBase+"/api/links/groups/"+id, &zohttp.RequestOpts{JSON: body}); if err != nil { return err }; return output.JSONRaw(raw) }},
			{Name: "update-personal", Usage: "Update a personal bookmark", ArgsUsage: "<bookmark-id>", Flags: []cli.Flag{&cli.StringFlag{Name: "json", Required: true, Usage: "JSON body"}}, Action: func(_ context.Context, cmd *cli.Command) error { id := cmd.Args().Get(0); if id == "" { return internal.NewValidationError("bookmark-id argument required") }; c, err := getClient(); if err != nil { return err }; var body any; json.Unmarshal([]byte(cmd.String("json")), &body); raw, err := c.Request("PUT", c.MailBase+"/api/links/me/"+id, &zohttp.RequestOpts{JSON: body}); if err != nil { return err }; return output.JSONRaw(raw) }},
			{Name: "update-group", Usage: "Update a group bookmark", ArgsUsage: "<group-id> <bookmark-id>", Flags: []cli.Flag{&cli.StringFlag{Name: "json", Required: true, Usage: "JSON body"}}, Action: func(_ context.Context, cmd *cli.Command) error { g := cmd.Args().Get(0); b := cmd.Args().Get(1); if g == "" || b == "" { return internal.NewValidationError("group-id and bookmark-id arguments required") }; c, err := getClient(); if err != nil { return err }; var body any; json.Unmarshal([]byte(cmd.String("json")), &body); raw, err := c.Request("PUT", c.MailBase+"/api/links/groups/"+g+"/"+b, &zohttp.RequestOpts{JSON: body}); if err != nil { return err }; return output.JSONRaw(raw) }},
			{Name: "delete-personal", Usage: "Delete a personal bookmark", ArgsUsage: "<bookmark-id>", Action: func(_ context.Context, cmd *cli.Command) error { id := cmd.Args().Get(0); if id == "" { return internal.NewValidationError("bookmark-id argument required") }; c, err := getClient(); if err != nil { return err }; raw, err := c.Request("DELETE", c.MailBase+"/api/links/me/"+id, nil); if err != nil { return err }; return output.JSONRaw(raw) }},
			{Name: "delete-group", Usage: "Delete a group bookmark", ArgsUsage: "<group-id> <bookmark-id>", Action: func(_ context.Context, cmd *cli.Command) error { g := cmd.Args().Get(0); b := cmd.Args().Get(1); if g == "" || b == "" { return internal.NewValidationError("group-id and bookmark-id arguments required") }; c, err := getClient(); if err != nil { return err }; raw, err := c.Request("DELETE", c.MailBase+"/api/links/groups/"+g+"/"+b, nil); if err != nil { return err }; return output.JSONRaw(raw) }},
			{Name: "favorites", Usage: "List favorite bookmarks", Action: func(_ context.Context, cmd *cli.Command) error { c, err := getClient(); if err != nil { return err }; raw, err := c.Request("GET", c.MailBase+"/api/links/favorites", nil); if err != nil { return err }; return output.JSONRaw(raw) }},
			{Name: "shared", Usage: "List shared bookmarks", Action: func(_ context.Context, cmd *cli.Command) error { c, err := getClient(); if err != nil { return err }; raw, err := c.Request("GET", c.MailBase+"/api/links", nil); if err != nil { return err }; return output.JSONRaw(raw) }},
			{Name: "trash-personal", Usage: "List trashed personal bookmarks", Action: func(_ context.Context, cmd *cli.Command) error { c, err := getClient(); if err != nil { return err }; raw, err := c.Request("GET", c.MailBase+"/api/links/me/trash", nil); if err != nil { return err }; return output.JSONRaw(raw) }},
			{Name: "trash-group", Usage: "List trashed group bookmarks", ArgsUsage: "<group-id>", Action: func(_ context.Context, cmd *cli.Command) error { id := cmd.Args().Get(0); if id == "" { return internal.NewValidationError("group-id argument required") }; c, err := getClient(); if err != nil { return err }; raw, err := c.Request("GET", c.MailBase+"/api/links/groups/"+id+"/trash", nil); if err != nil { return err }; return output.JSONRaw(raw) }},
			{Name: "restore-group", Usage: "Restore a trashed group bookmark", ArgsUsage: "<group-id> <bookmark-id>", Action: func(_ context.Context, cmd *cli.Command) error { g := cmd.Args().Get(0); b := cmd.Args().Get(1); if g == "" || b == "" { return internal.NewValidationError("group-id and bookmark-id arguments required") }; c, err := getClient(); if err != nil { return err }; raw, err := c.Request("PUT", c.MailBase+"/api/links/groups/"+g+"/"+b+"/restore", nil); if err != nil { return err }; return output.JSONRaw(raw) }},
			{Name: "favorite-personal", Usage: "Favorite a personal bookmark", ArgsUsage: "<bookmark-id>", Action: func(_ context.Context, cmd *cli.Command) error { id := cmd.Args().Get(0); if id == "" { return internal.NewValidationError("bookmark-id argument required") }; c, err := getClient(); if err != nil { return err }; raw, err := c.Request("PUT", c.MailBase+"/api/links/me/"+id+"/favorite", nil); if err != nil { return err }; return output.JSONRaw(raw) }},
			{Name: "favorite-group", Usage: "Favorite a group bookmark", ArgsUsage: "<group-id> <bookmark-id>", Action: func(_ context.Context, cmd *cli.Command) error { g := cmd.Args().Get(0); b := cmd.Args().Get(1); if g == "" || b == "" { return internal.NewValidationError("group-id and bookmark-id arguments required") }; c, err := getClient(); if err != nil { return err }; raw, err := c.Request("PUT", c.MailBase+"/api/links/groups/"+g+"/"+b+"/favorite", nil); if err != nil { return err }; return output.JSONRaw(raw) }},
			{Name: "unfavorite-personal", Usage: "Unfavorite a personal bookmark", ArgsUsage: "<bookmark-id>", Action: func(_ context.Context, cmd *cli.Command) error { id := cmd.Args().Get(0); if id == "" { return internal.NewValidationError("bookmark-id argument required") }; c, err := getClient(); if err != nil { return err }; raw, err := c.Request("DELETE", c.MailBase+"/api/links/me/"+id+"/favorite", nil); if err != nil { return err }; return output.JSONRaw(raw) }},
			{Name: "unfavorite-group", Usage: "Unfavorite a group bookmark", ArgsUsage: "<group-id> <bookmark-id>", Action: func(_ context.Context, cmd *cli.Command) error { g := cmd.Args().Get(0); b := cmd.Args().Get(1); if g == "" || b == "" { return internal.NewValidationError("group-id and bookmark-id arguments required") }; c, err := getClient(); if err != nil { return err }; raw, err := c.Request("DELETE", c.MailBase+"/api/links/groups/"+g+"/"+b+"/favorite", nil); if err != nil { return err }; return output.JSONRaw(raw) }},
			{Name: "link-groups", Usage: "List bookmark groups", Action: func(_ context.Context, cmd *cli.Command) error { c, err := getClient(); if err != nil { return err }; raw, err := c.Request("GET", c.MailBase+"/api/links/groups", nil); if err != nil { return err }; return output.JSONRaw(raw) }},
			{Name: "collections-personal", Usage: "List personal bookmark collections", Action: func(_ context.Context, cmd *cli.Command) error { c, err := getClient(); if err != nil { return err }; raw, err := c.Request("GET", c.MailBase+"/api/links/me/collections", nil); if err != nil { return err }; return output.JSONRaw(raw) }},
			{Name: "collections-group", Usage: "List group bookmark collections", ArgsUsage: "<group-id>", Action: func(_ context.Context, cmd *cli.Command) error { id := cmd.Args().Get(0); if id == "" { return internal.NewValidationError("group-id argument required") }; c, err := getClient(); if err != nil { return err }; raw, err := c.Request("GET", c.MailBase+"/api/links/groups/"+id+"/collections", nil); if err != nil { return err }; return output.JSONRaw(raw) }},
			{Name: "collection-bookmarks-personal", Usage: "List bookmarks in a personal collection", ArgsUsage: "<collection-id>", Action: func(_ context.Context, cmd *cli.Command) error { id := cmd.Args().Get(0); if id == "" { return internal.NewValidationError("collection-id argument required") }; c, err := getClient(); if err != nil { return err }; raw, err := c.Request("GET", c.MailBase+"/api/links/me/collections/"+id, nil); if err != nil { return err }; return output.JSONRaw(raw) }},
			{Name: "collection-bookmarks-group", Usage: "List bookmarks in a group collection", ArgsUsage: "<group-id> <collection-id>", Action: func(_ context.Context, cmd *cli.Command) error { g := cmd.Args().Get(0); col := cmd.Args().Get(1); if g == "" || col == "" { return internal.NewValidationError("group-id and collection-id arguments required") }; c, err := getClient(); if err != nil { return err }; raw, err := c.Request("GET", c.MailBase+"/api/links/groups/"+g+"/collections/"+col, nil); if err != nil { return err }; return output.JSONRaw(raw) }},
			{Name: "create-collection-personal", Usage: "Create a personal bookmark collection", Flags: []cli.Flag{&cli.StringFlag{Name: "json", Required: true, Usage: "JSON body"}}, Action: func(_ context.Context, cmd *cli.Command) error { c, err := getClient(); if err != nil { return err }; var body any; json.Unmarshal([]byte(cmd.String("json")), &body); raw, err := c.Request("POST", c.MailBase+"/api/links/me/collections", &zohttp.RequestOpts{JSON: body}); if err != nil { return err }; return output.JSONRaw(raw) }},
			{Name: "create-collection-group", Usage: "Create a group bookmark collection", ArgsUsage: "<group-id>", Flags: []cli.Flag{&cli.StringFlag{Name: "json", Required: true, Usage: "JSON body"}}, Action: func(_ context.Context, cmd *cli.Command) error { id := cmd.Args().Get(0); if id == "" { return internal.NewValidationError("group-id argument required") }; c, err := getClient(); if err != nil { return err }; var body any; json.Unmarshal([]byte(cmd.String("json")), &body); raw, err := c.Request("POST", c.MailBase+"/api/links/groups/"+id+"/collections", &zohttp.RequestOpts{JSON: body}); if err != nil { return err }; return output.JSONRaw(raw) }},
			{Name: "update-collection-personal", Usage: "Update a personal bookmark collection", ArgsUsage: "<collection-id>", Flags: []cli.Flag{&cli.StringFlag{Name: "json", Required: true, Usage: "JSON body"}}, Action: func(_ context.Context, cmd *cli.Command) error { id := cmd.Args().Get(0); if id == "" { return internal.NewValidationError("collection-id argument required") }; c, err := getClient(); if err != nil { return err }; var body any; json.Unmarshal([]byte(cmd.String("json")), &body); raw, err := c.Request("PUT", c.MailBase+"/api/links/me/collections/"+id, &zohttp.RequestOpts{JSON: body}); if err != nil { return err }; return output.JSONRaw(raw) }},
			{Name: "update-collection-group", Usage: "Update a group bookmark collection", ArgsUsage: "<group-id> <collection-id>", Flags: []cli.Flag{&cli.StringFlag{Name: "json", Required: true, Usage: "JSON body"}}, Action: func(_ context.Context, cmd *cli.Command) error { g := cmd.Args().Get(0); col := cmd.Args().Get(1); if g == "" || col == "" { return internal.NewValidationError("group-id and collection-id arguments required") }; c, err := getClient(); if err != nil { return err }; var body any; json.Unmarshal([]byte(cmd.String("json")), &body); raw, err := c.Request("PUT", c.MailBase+"/api/links/groups/"+g+"/collections/"+col, &zohttp.RequestOpts{JSON: body}); if err != nil { return err }; return output.JSONRaw(raw) }},
			{Name: "delete-collection-personal", Usage: "Delete a personal bookmark collection", ArgsUsage: "<collection-id>", Action: func(_ context.Context, cmd *cli.Command) error { id := cmd.Args().Get(0); if id == "" { return internal.NewValidationError("collection-id argument required") }; c, err := getClient(); if err != nil { return err }; raw, err := c.Request("DELETE", c.MailBase+"/api/links/me/collections/"+id, nil); if err != nil { return err }; return output.JSONRaw(raw) }},
			{Name: "delete-collection-group", Usage: "Delete a group bookmark collection", ArgsUsage: "<group-id> <collection-id>", Action: func(_ context.Context, cmd *cli.Command) error { g := cmd.Args().Get(0); col := cmd.Args().Get(1); if g == "" || col == "" { return internal.NewValidationError("group-id and collection-id arguments required") }; c, err := getClient(); if err != nil { return err }; raw, err := c.Request("DELETE", c.MailBase+"/api/links/groups/"+g+"/collections/"+col, nil); if err != nil { return err }; return output.JSONRaw(raw) }},
			{Name: "all-group-collections", Usage: "List all group bookmark collections", Action: func(_ context.Context, cmd *cli.Command) error { c, err := getClient(); if err != nil { return err }; raw, err := c.Request("GET", c.MailBase+"/api/links/groups/collections", nil); if err != nil { return err }; return output.JSONRaw(raw) }},
		},
	}
}

func notesCmd() *cli.Command {
	return &cli.Command{
		Name:  "notes",
		Usage: "Mail notes operations",
		Commands: []*cli.Command{
			{Name: "list-personal", Usage: "List personal notes", Action: func(_ context.Context, cmd *cli.Command) error { c, err := getClient(); if err != nil { return err }; raw, err := c.Request("GET", c.MailBase+"/api/notes/me", nil); if err != nil { return err }; return output.JSONRaw(raw) }},
			{Name: "list-group", Usage: "List group notes", ArgsUsage: "<group-id>", Action: func(_ context.Context, cmd *cli.Command) error { id := cmd.Args().Get(0); if id == "" { return internal.NewValidationError("group-id argument required") }; c, err := getClient(); if err != nil { return err }; raw, err := c.Request("GET", c.MailBase+"/api/notes/groups/"+id, nil); if err != nil { return err }; return output.JSONRaw(raw) }},
			{Name: "get-personal", Usage: "Get a personal note", ArgsUsage: "<note-id>", Action: func(_ context.Context, cmd *cli.Command) error { id := cmd.Args().Get(0); if id == "" { return internal.NewValidationError("note-id argument required") }; c, err := getClient(); if err != nil { return err }; raw, err := c.Request("GET", c.MailBase+"/api/notes/me/"+id, nil); if err != nil { return err }; return output.JSONRaw(raw) }},
			{Name: "get-group", Usage: "Get a group note", ArgsUsage: "<group-id> <note-id>", Action: func(_ context.Context, cmd *cli.Command) error { g := cmd.Args().Get(0); n := cmd.Args().Get(1); if g == "" || n == "" { return internal.NewValidationError("group-id and note-id arguments required") }; c, err := getClient(); if err != nil { return err }; raw, err := c.Request("GET", c.MailBase+"/api/notes/groups/"+g+"/"+n, nil); if err != nil { return err }; return output.JSONRaw(raw) }},
			{Name: "create-personal", Usage: "Create a personal note", Flags: []cli.Flag{&cli.StringFlag{Name: "json", Required: true, Usage: "JSON body"}}, Action: func(_ context.Context, cmd *cli.Command) error { c, err := getClient(); if err != nil { return err }; var body any; json.Unmarshal([]byte(cmd.String("json")), &body); raw, err := c.Request("POST", c.MailBase+"/api/notes/me", &zohttp.RequestOpts{JSON: body}); if err != nil { return err }; return output.JSONRaw(raw) }},
			{Name: "create-group", Usage: "Create a group note", ArgsUsage: "<group-id>", Flags: []cli.Flag{&cli.StringFlag{Name: "json", Required: true, Usage: "JSON body"}}, Action: func(_ context.Context, cmd *cli.Command) error { id := cmd.Args().Get(0); if id == "" { return internal.NewValidationError("group-id argument required") }; c, err := getClient(); if err != nil { return err }; var body any; json.Unmarshal([]byte(cmd.String("json")), &body); raw, err := c.Request("POST", c.MailBase+"/api/notes/groups/"+id, &zohttp.RequestOpts{JSON: body}); if err != nil { return err }; return output.JSONRaw(raw) }},
			{Name: "update-personal", Usage: "Update a personal note", ArgsUsage: "<note-id>", Flags: []cli.Flag{&cli.StringFlag{Name: "json", Required: true, Usage: "JSON body"}}, Action: func(_ context.Context, cmd *cli.Command) error { id := cmd.Args().Get(0); if id == "" { return internal.NewValidationError("note-id argument required") }; c, err := getClient(); if err != nil { return err }; var body any; json.Unmarshal([]byte(cmd.String("json")), &body); raw, err := c.Request("PUT", c.MailBase+"/api/notes/me/"+id, &zohttp.RequestOpts{JSON: body}); if err != nil { return err }; return output.JSONRaw(raw) }},
			{Name: "update-group", Usage: "Update a group note", ArgsUsage: "<group-id> <note-id>", Flags: []cli.Flag{&cli.StringFlag{Name: "json", Required: true, Usage: "JSON body"}}, Action: func(_ context.Context, cmd *cli.Command) error { g := cmd.Args().Get(0); n := cmd.Args().Get(1); if g == "" || n == "" { return internal.NewValidationError("group-id and note-id arguments required") }; c, err := getClient(); if err != nil { return err }; var body any; json.Unmarshal([]byte(cmd.String("json")), &body); raw, err := c.Request("PUT", c.MailBase+"/api/notes/groups/"+g+"/"+n, &zohttp.RequestOpts{JSON: body}); if err != nil { return err }; return output.JSONRaw(raw) }},
			{Name: "delete-personal", Usage: "Delete a personal note", ArgsUsage: "<note-id>", Action: func(_ context.Context, cmd *cli.Command) error { id := cmd.Args().Get(0); if id == "" { return internal.NewValidationError("note-id argument required") }; c, err := getClient(); if err != nil { return err }; raw, err := c.Request("DELETE", c.MailBase+"/api/notes/me/"+id, nil); if err != nil { return err }; return output.JSONRaw(raw) }},
			{Name: "delete-group", Usage: "Delete a group note", ArgsUsage: "<group-id> <note-id>", Action: func(_ context.Context, cmd *cli.Command) error { g := cmd.Args().Get(0); n := cmd.Args().Get(1); if g == "" || n == "" { return internal.NewValidationError("group-id and note-id arguments required") }; c, err := getClient(); if err != nil { return err }; raw, err := c.Request("DELETE", c.MailBase+"/api/notes/groups/"+g+"/"+n, nil); if err != nil { return err }; return output.JSONRaw(raw) }},
			{Name: "favorites", Usage: "List favorite notes", Action: func(_ context.Context, cmd *cli.Command) error { c, err := getClient(); if err != nil { return err }; raw, err := c.Request("GET", c.MailBase+"/api/notes/favorites", nil); if err != nil { return err }; return output.JSONRaw(raw) }},
			{Name: "shared", Usage: "List notes shared to me", Action: func(_ context.Context, cmd *cli.Command) error { c, err := getClient(); if err != nil { return err }; raw, err := c.Request("GET", c.MailBase+"/api/notes/sharedtome", nil); if err != nil { return err }; return output.JSONRaw(raw) }},
			{Name: "note-groups", Usage: "List note groups", Action: func(_ context.Context, cmd *cli.Command) error { c, err := getClient(); if err != nil { return err }; raw, err := c.Request("GET", c.MailBase+"/api/notes/groups", nil); if err != nil { return err }; return output.JSONRaw(raw) }},
			{Name: "books-personal", Usage: "List personal notebooks", Action: func(_ context.Context, cmd *cli.Command) error { c, err := getClient(); if err != nil { return err }; raw, err := c.Request("GET", c.MailBase+"/api/notes/me/books", nil); if err != nil { return err }; return output.JSONRaw(raw) }},
			{Name: "books-group", Usage: "List group notebooks", ArgsUsage: "<group-id>", Action: func(_ context.Context, cmd *cli.Command) error { id := cmd.Args().Get(0); if id == "" { return internal.NewValidationError("group-id argument required") }; c, err := getClient(); if err != nil { return err }; raw, err := c.Request("GET", c.MailBase+"/api/notes/groups/"+id+"/books", nil); if err != nil { return err }; return output.JSONRaw(raw) }},
			{Name: "book-notes-personal", Usage: "List notes in a personal notebook", ArgsUsage: "<book-id>", Action: func(_ context.Context, cmd *cli.Command) error { id := cmd.Args().Get(0); if id == "" { return internal.NewValidationError("book-id argument required") }; c, err := getClient(); if err != nil { return err }; raw, err := c.Request("GET", c.MailBase+"/api/notes/me/books/"+id, nil); if err != nil { return err }; return output.JSONRaw(raw) }},
			{Name: "book-notes-group", Usage: "List notes in a group notebook", ArgsUsage: "<group-id> <book-id>", Action: func(_ context.Context, cmd *cli.Command) error { g := cmd.Args().Get(0); b := cmd.Args().Get(1); if g == "" || b == "" { return internal.NewValidationError("group-id and book-id arguments required") }; c, err := getClient(); if err != nil { return err }; raw, err := c.Request("GET", c.MailBase+"/api/notes/groups/"+g+"/books/"+b, nil); if err != nil { return err }; return output.JSONRaw(raw) }},
			{Name: "create-book-personal", Usage: "Create a personal notebook", Flags: []cli.Flag{&cli.StringFlag{Name: "json", Required: true, Usage: "JSON body"}}, Action: func(_ context.Context, cmd *cli.Command) error { c, err := getClient(); if err != nil { return err }; var body any; json.Unmarshal([]byte(cmd.String("json")), &body); raw, err := c.Request("POST", c.MailBase+"/api/notes/me/books", &zohttp.RequestOpts{JSON: body}); if err != nil { return err }; return output.JSONRaw(raw) }},
			{Name: "create-book-group", Usage: "Create a group notebook", ArgsUsage: "<group-id>", Flags: []cli.Flag{&cli.StringFlag{Name: "json", Required: true, Usage: "JSON body"}}, Action: func(_ context.Context, cmd *cli.Command) error { id := cmd.Args().Get(0); if id == "" { return internal.NewValidationError("group-id argument required") }; c, err := getClient(); if err != nil { return err }; var body any; json.Unmarshal([]byte(cmd.String("json")), &body); raw, err := c.Request("POST", c.MailBase+"/api/notes/groups/"+id+"/books", &zohttp.RequestOpts{JSON: body}); if err != nil { return err }; return output.JSONRaw(raw) }},
			{Name: "update-book-personal", Usage: "Update a personal notebook", ArgsUsage: "<book-id>", Flags: []cli.Flag{&cli.StringFlag{Name: "json", Required: true, Usage: "JSON body"}}, Action: func(_ context.Context, cmd *cli.Command) error { id := cmd.Args().Get(0); if id == "" { return internal.NewValidationError("book-id argument required") }; c, err := getClient(); if err != nil { return err }; var body any; json.Unmarshal([]byte(cmd.String("json")), &body); raw, err := c.Request("PUT", c.MailBase+"/api/notes/me/books/"+id, &zohttp.RequestOpts{JSON: body}); if err != nil { return err }; return output.JSONRaw(raw) }},
			{Name: "update-book-group", Usage: "Update a group notebook", ArgsUsage: "<group-id> <book-id>", Flags: []cli.Flag{&cli.StringFlag{Name: "json", Required: true, Usage: "JSON body"}}, Action: func(_ context.Context, cmd *cli.Command) error { g := cmd.Args().Get(0); b := cmd.Args().Get(1); if g == "" || b == "" { return internal.NewValidationError("group-id and book-id arguments required") }; c, err := getClient(); if err != nil { return err }; var body any; json.Unmarshal([]byte(cmd.String("json")), &body); raw, err := c.Request("PUT", c.MailBase+"/api/notes/groups/"+g+"/books/"+b, &zohttp.RequestOpts{JSON: body}); if err != nil { return err }; return output.JSONRaw(raw) }},
			{Name: "delete-book-personal", Usage: "Delete a personal notebook", ArgsUsage: "<book-id>", Action: func(_ context.Context, cmd *cli.Command) error { id := cmd.Args().Get(0); if id == "" { return internal.NewValidationError("book-id argument required") }; c, err := getClient(); if err != nil { return err }; raw, err := c.Request("DELETE", c.MailBase+"/api/notes/me/books/"+id, nil); if err != nil { return err }; return output.JSONRaw(raw) }},
			{Name: "delete-book-group", Usage: "Delete a group notebook", ArgsUsage: "<group-id> <book-id>", Action: func(_ context.Context, cmd *cli.Command) error { g := cmd.Args().Get(0); b := cmd.Args().Get(1); if g == "" || b == "" { return internal.NewValidationError("group-id and book-id arguments required") }; c, err := getClient(); if err != nil { return err }; raw, err := c.Request("DELETE", c.MailBase+"/api/notes/groups/"+g+"/books/"+b, nil); if err != nil { return err }; return output.JSONRaw(raw) }},
			{Name: "attachments-personal", Usage: "List attachments of a personal note", ArgsUsage: "<note-id>", Action: func(_ context.Context, cmd *cli.Command) error { id := cmd.Args().Get(0); if id == "" { return internal.NewValidationError("note-id argument required") }; c, err := getClient(); if err != nil { return err }; raw, err := c.Request("GET", c.MailBase+"/api/notes/me/"+id+"/attachments", nil); if err != nil { return err }; return output.JSONRaw(raw) }},
			{Name: "attachments-group", Usage: "List attachments of a group note", ArgsUsage: "<group-id> <note-id>", Action: func(_ context.Context, cmd *cli.Command) error { g := cmd.Args().Get(0); n := cmd.Args().Get(1); if g == "" || n == "" { return internal.NewValidationError("group-id and note-id arguments required") }; c, err := getClient(); if err != nil { return err }; raw, err := c.Request("GET", c.MailBase+"/api/notes/groups/"+g+"/"+n+"/attachments", nil); if err != nil { return err }; return output.JSONRaw(raw) }},
			{
				Name: "attachment-personal", Usage: "Download a personal note attachment", ArgsUsage: "<note-id> <attachment-id>",
				Flags: []cli.Flag{&cli.StringFlag{Name: "output", Usage: "Output file path"}},
				Action: func(_ context.Context, cmd *cli.Command) error {
					noteID := cmd.Args().Get(0); attachID := cmd.Args().Get(1)
					if noteID == "" || attachID == "" { return internal.NewValidationError("note-id and attachment-id arguments required") }
					c, err := getClient(); if err != nil { return err }
					body, _, _, err := c.RequestRaw("GET", c.MailBase+"/api/notes/me/"+noteID+"/attachments/"+attachID, nil)
					if err != nil { return err }
					if out := cmd.String("output"); out != "" { if err := os.WriteFile(out, body, 0644); err != nil { return err }; return output.JSON(map[string]any{"ok": true, "path": out, "size": len(body)}) }
					os.Stdout.Write(body); return nil
				},
			},
			{
				Name: "attachment-group", Usage: "Download a group note attachment", ArgsUsage: "<group-id> <note-id> <attachment-id>",
				Flags: []cli.Flag{&cli.StringFlag{Name: "output", Usage: "Output file path"}},
				Action: func(_ context.Context, cmd *cli.Command) error {
					groupID := cmd.Args().Get(0); noteID := cmd.Args().Get(1); attachID := cmd.Args().Get(2)
					if groupID == "" || noteID == "" || attachID == "" { return internal.NewValidationError("group-id, note-id, and attachment-id arguments required") }
					c, err := getClient(); if err != nil { return err }
					body, _, _, err := c.RequestRaw("GET", c.MailBase+"/api/notes/groups/"+groupID+"/"+noteID+"/attachments/"+attachID, nil)
					if err != nil { return err }
					if out := cmd.String("output"); out != "" { if err := os.WriteFile(out, body, 0644); err != nil { return err }; return output.JSON(map[string]any{"ok": true, "path": out, "size": len(body)}) }
					os.Stdout.Write(body); return nil
				},
			},
			{
				Name: "upload-attachment-personal", Usage: "Upload an attachment to a personal note", ArgsUsage: "<note-id>",
				Flags: []cli.Flag{&cli.StringFlag{Name: "file", Required: true, Usage: "Path to file"}},
				Action: func(_ context.Context, cmd *cli.Command) error {
					noteID := cmd.Args().Get(0); if noteID == "" { return internal.NewValidationError("note-id argument required") }
					c, err := getClient(); if err != nil { return err }
					filePath := cmd.String("file"); fileData, err := os.ReadFile(filePath)
					if err != nil { return fmt.Errorf("failed to read file: %w", err) }
					raw, err := c.Request("POST", c.MailBase+"/api/notes/me/"+noteID+"/attachments", &zohttp.RequestOpts{Files: map[string]zohttp.FileUpload{"file": {Filename: filepath.Base(filePath), Data: fileData}}})
					if err != nil { return err }; return output.JSONRaw(raw)
				},
			},
			{
				Name: "upload-attachment-group", Usage: "Upload an attachment to a group note", ArgsUsage: "<group-id> <note-id>",
				Flags: []cli.Flag{&cli.StringFlag{Name: "file", Required: true, Usage: "Path to file"}},
				Action: func(_ context.Context, cmd *cli.Command) error {
					groupID := cmd.Args().Get(0); noteID := cmd.Args().Get(1)
					if groupID == "" || noteID == "" { return internal.NewValidationError("group-id and note-id arguments required") }
					c, err := getClient(); if err != nil { return err }
					filePath := cmd.String("file"); fileData, err := os.ReadFile(filePath)
					if err != nil { return fmt.Errorf("failed to read file: %w", err) }
					raw, err := c.Request("POST", c.MailBase+"/api/notes/groups/"+groupID+"/"+noteID+"/attachments", &zohttp.RequestOpts{Files: map[string]zohttp.FileUpload{"file": {Filename: filepath.Base(filePath), Data: fileData}}})
					if err != nil { return err }; return output.JSONRaw(raw)
				},
			},
			{Name: "delete-attachment-personal", Usage: "Delete a personal note attachment", ArgsUsage: "<note-id> <attachment-id>", Action: func(_ context.Context, cmd *cli.Command) error { n := cmd.Args().Get(0); a := cmd.Args().Get(1); if n == "" || a == "" { return internal.NewValidationError("note-id and attachment-id arguments required") }; c, err := getClient(); if err != nil { return err }; raw, err := c.Request("DELETE", c.MailBase+"/api/notes/me/"+n+"/attachments/"+a, nil); if err != nil { return err }; return output.JSONRaw(raw) }},
			{Name: "delete-attachment-group", Usage: "Delete a group note attachment", ArgsUsage: "<group-id> <note-id> <attachment-id>", Action: func(_ context.Context, cmd *cli.Command) error { g := cmd.Args().Get(0); n := cmd.Args().Get(1); a := cmd.Args().Get(2); if g == "" || n == "" || a == "" { return internal.NewValidationError("group-id, note-id, and attachment-id arguments required") }; c, err := getClient(); if err != nil { return err }; raw, err := c.Request("DELETE", c.MailBase+"/api/notes/groups/"+g+"/"+n+"/attachments/"+a, nil); if err != nil { return err }; return output.JSONRaw(raw) }},
			{Name: "favorite-personal", Usage: "Favorite a personal note", ArgsUsage: "<note-id>", Action: func(_ context.Context, cmd *cli.Command) error { id := cmd.Args().Get(0); if id == "" { return internal.NewValidationError("note-id argument required") }; c, err := getClient(); if err != nil { return err }; raw, err := c.Request("PUT", c.MailBase+"/api/notes/me/"+id+"/favorite", nil); if err != nil { return err }; return output.JSONRaw(raw) }},
			{Name: "favorite-group", Usage: "Favorite a group note", ArgsUsage: "<group-id> <note-id>", Action: func(_ context.Context, cmd *cli.Command) error { g := cmd.Args().Get(0); n := cmd.Args().Get(1); if g == "" || n == "" { return internal.NewValidationError("group-id and note-id arguments required") }; c, err := getClient(); if err != nil { return err }; raw, err := c.Request("PUT", c.MailBase+"/api/notes/groups/"+g+"/"+n+"/favorite", nil); if err != nil { return err }; return output.JSONRaw(raw) }},
			{Name: "unfavorite-personal", Usage: "Unfavorite a personal note", ArgsUsage: "<note-id>", Action: func(_ context.Context, cmd *cli.Command) error { id := cmd.Args().Get(0); if id == "" { return internal.NewValidationError("note-id argument required") }; c, err := getClient(); if err != nil { return err }; raw, err := c.Request("DELETE", c.MailBase+"/api/notes/me/"+id+"/favorite", nil); if err != nil { return err }; return output.JSONRaw(raw) }},
			{Name: "unfavorite-group", Usage: "Unfavorite a group note", ArgsUsage: "<group-id> <note-id>", Action: func(_ context.Context, cmd *cli.Command) error { g := cmd.Args().Get(0); n := cmd.Args().Get(1); if g == "" || n == "" { return internal.NewValidationError("group-id and note-id arguments required") }; c, err := getClient(); if err != nil { return err }; raw, err := c.Request("DELETE", c.MailBase+"/api/notes/groups/"+g+"/"+n+"/favorite", nil); if err != nil { return err }; return output.JSONRaw(raw) }},
		},
	}
}
