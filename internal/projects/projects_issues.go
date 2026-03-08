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

func issuesCmd() *cli.Command {
	return &cli.Command{
		Name:  "issues",
		Usage: "Project issue operations",
		Commands: []*cli.Command{
			{
				Name:  "list",
				Usage: "List issues in a project",
				Flags: []cli.Flag{portalFlag, projectFlag, allFlag, limitFlag},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					portal, err := requirePortal(cmd)
					if err != nil {
						return err
					}
					url := base(c, portal, cmd.String("project")) + "/issues"
					return paginateProjectsList(c, cmd, url, "issues", nil)
				},
			},
			{
				Name:      "get",
				Usage:     "Get a single issue",
				ArgsUsage: "<issue-id>",
				Flags:     []cli.Flag{portalFlag, projectFlag},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					portal, err := requirePortal(cmd)
					if err != nil {
						return err
					}
					url := base(c, portal, cmd.String("project")) + "/issues/" + cmd.Args().First()
					raw, err := c.Request("GET", url, nil)
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:  "create",
				Usage: "Create an issue",
				Flags: []cli.Flag{
					portalFlag, projectFlag,
					&cli.StringFlag{Name: "name", Required: true, Usage: "Issue title"},
					&cli.StringFlag{Name: "json", Usage: "Additional fields as JSON"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					portal, err := requirePortal(cmd)
					if err != nil {
						return err
					}
					body := map[string]any{"name": cmd.String("name")}
					if err := internal.MergeJSON(cmd, body); err != nil {
						return err
					}
					url := base(c, portal, cmd.String("project")) + "/issues"
					raw, err := c.Request("POST", url, &zohttp.RequestOpts{JSON: body})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "update",
				Usage:     "Update an issue",
				ArgsUsage: "<issue-id>",
				Flags: []cli.Flag{
					portalFlag, projectFlag,
					&cli.StringFlag{Name: "name", Usage: "Issue title"},
					&cli.StringFlag{Name: "description", Usage: "Issue description"},
					&cli.StringFlag{Name: "priority", Usage: "Priority"},
					&cli.StringFlag{Name: "json", Usage: "Additional fields as JSON"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					portal, err := requirePortal(cmd)
					if err != nil {
						return err
					}
					body := map[string]any{}
					if cmd.IsSet("name") {
						body["name"] = cmd.String("name")
					}
					if cmd.IsSet("description") {
						body["description"] = cmd.String("description")
					}
					if cmd.IsSet("priority") {
						body["priority"] = cmd.String("priority")
					}
					if err := internal.MergeJSON(cmd, body); err != nil {
						return err
					}
					url := base(c, portal, cmd.String("project")) + "/issues/" + cmd.Args().First()
					raw, err := c.Request("PATCH", url, &zohttp.RequestOpts{JSON: body})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "delete",
				Usage:     "Delete an issue",
				ArgsUsage: "<issue-id>",
				Flags:     []cli.Flag{portalFlag, projectFlag},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					portal, err := requirePortal(cmd)
					if err != nil {
						return err
					}
					url := base(c, portal, cmd.String("project")) + "/issues/" + cmd.Args().First()
					raw, err := c.Request("DELETE", url, nil)
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},

			{
				Name:      "description",
				Usage:     "Get issue description",
				ArgsUsage: "<issue-id>",
				Flags:     []cli.Flag{portalFlag, projectFlag},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					portal, err := requirePortal(cmd)
					if err != nil {
						return err
					}
					url := base(c, portal, cmd.String("project")) + "/issues/" + cmd.Args().First() + "/description"
					raw, err := c.Request("GET", url, nil)
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "move",
				Usage:     "Move an issue",
				ArgsUsage: "<issue-id>",
				Flags: []cli.Flag{
					portalFlag, projectFlag,
					&cli.StringFlag{Name: "to_project", Usage: "Target project ID"},
					&cli.StringFlag{Name: "json", Usage: "Additional fields as JSON"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					portal, err := requirePortal(cmd)
					if err != nil {
						return err
					}
					body := map[string]any{}
					if cmd.IsSet("to_project") {
						body["to_project"] = cmd.String("to_project")
					}
					if err := internal.MergeJSON(cmd, body); err != nil {
						return err
					}
					url := base(c, portal, cmd.String("project")) + "/issues/" + cmd.Args().First() + "/move"
					raw, err := c.Request("POST", url, &zohttp.RequestOpts{JSON: body})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "clone",
				Usage:     "Clone an issue",
				ArgsUsage: "<issue-id>",
				Flags:     []cli.Flag{portalFlag, projectFlag},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					portal, err := requirePortal(cmd)
					if err != nil {
						return err
					}
					url := base(c, portal, cmd.String("project")) + "/issues/" + cmd.Args().First() + "/clone"
					raw, err := c.Request("POST", url, nil)
					if err != nil {
						return err
					}
					var envelope map[string]json.RawMessage
					if json.Unmarshal(raw, &envelope) == nil {
						if bugs, ok := envelope["bugs"]; ok {
							var arr []json.RawMessage
							if json.Unmarshal(bugs, &arr) == nil && len(arr) > 0 {
								return output.JSONRaw(arr[0])
							}
						}
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "activities",
				Usage:     "Get issue activities",
				ArgsUsage: "<issue-id>",
				Flags:     []cli.Flag{portalFlag, projectFlag},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					portal, err := requirePortal(cmd)
					if err != nil {
						return err
					}
					url := base(c, portal, cmd.String("project")) + "/issues/" + cmd.Args().First() + "/activities"
					raw, err := c.Request("GET", url, nil)
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
		},
	}
}

func issueCommentsCmd() *cli.Command {
	return &cli.Command{
		Name:  "issue-comments",
		Usage: "Issue comment operations",
		Commands: []*cli.Command{
			{
				Name:  "list",
				Usage: "List issue comments",
				Flags: []cli.Flag{
					portalFlag, projectFlag,
					&cli.StringFlag{Name: "issue", Required: true, Usage: "Issue ID"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					portal, err := requirePortal(cmd)
					if err != nil {
						return err
					}
					url := base(c, portal, cmd.String("project")) + "/issues/" + cmd.String("issue") + "/comments"
					raw, err := c.Request("GET", url, nil)
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "get",
				Usage:     "Get an issue comment",
				ArgsUsage: "<comment-id>",
				Flags: []cli.Flag{
					portalFlag, projectFlag,
					&cli.StringFlag{Name: "issue", Required: true, Usage: "Issue ID"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					portal, err := requirePortal(cmd)
					if err != nil {
						return err
					}
					url := base(c, portal, cmd.String("project")) + "/issues/" + cmd.String("issue") + "/comments/" + cmd.Args().First()
					raw, err := c.Request("GET", url, nil)
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:  "add",
				Usage: "Add an issue comment",
				Flags: []cli.Flag{
					portalFlag, projectFlag,
					&cli.StringFlag{Name: "issue", Required: true, Usage: "Issue ID"},
					&cli.StringFlag{Name: "comment", Required: true, Usage: "Comment text"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					portal, err := requirePortal(cmd)
					if err != nil {
						return err
					}
					url := base(c, portal, cmd.String("project")) + "/issues/" + cmd.String("issue") + "/comments"
					raw, err := c.Request("POST", url, &zohttp.RequestOpts{JSON: map[string]string{"comment": cmd.String("comment")}})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "update",
				Usage:     "Update an issue comment",
				ArgsUsage: "<comment-id>",
				Flags: []cli.Flag{
					portalFlag, projectFlag,
					&cli.StringFlag{Name: "issue", Required: true, Usage: "Issue ID"},
					&cli.StringFlag{Name: "comment", Required: true, Usage: "Updated comment text"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					portal, err := requirePortal(cmd)
					if err != nil {
						return err
					}
					url := base(c, portal, cmd.String("project")) + "/issues/" + cmd.String("issue") + "/comments/" + cmd.Args().First()
					raw, err := c.Request("PATCH", url, &zohttp.RequestOpts{JSON: map[string]string{"comment": cmd.String("comment")}})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "delete",
				Usage:     "Delete an issue comment",
				ArgsUsage: "<comment-id>",
				Flags: []cli.Flag{
					portalFlag, projectFlag,
					&cli.StringFlag{Name: "issue", Required: true, Usage: "Issue ID"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					portal, err := requirePortal(cmd)
					if err != nil {
						return err
					}
					url := base(c, portal, cmd.String("project")) + "/issues/" + cmd.String("issue") + "/comments/" + cmd.Args().First()
					raw, err := c.Request("DELETE", url, nil)
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
		},
	}
}

func issueFollowersCmd() *cli.Command {
	return &cli.Command{
		Name:  "issue-followers",
		Usage: "Issue follower operations",
		Commands: []*cli.Command{
			{
				Name:      "list",
				Usage:     "List issue followers",
				ArgsUsage: "<issue-id>",
				Flags:     []cli.Flag{portalFlag, projectFlag},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					portal, err := requirePortal(cmd)
					if err != nil {
						return err
					}
					url := base(c, portal, cmd.String("project")) + "/issues/" + cmd.Args().First() + "/followers"
					raw, err := c.Request("GET", url, &zohttp.RequestOpts{Params: map[string]string{"page": "1", "per_page": "200"}})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "follow",
				Usage:     "Follow an issue",
				ArgsUsage: "<issue-id>",
				Flags:     []cli.Flag{portalFlag, projectFlag},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					portal, err := requirePortal(cmd)
					if err != nil {
						return err
					}
					url := base(c, portal, cmd.String("project")) + "/issues/" + cmd.Args().First() + "/followers"
					raw, err := c.Request("POST", url, nil)
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "remove",
				Usage:     "Remove issue followers",
				ArgsUsage: "<issue-id>",
				Flags: []cli.Flag{
					portalFlag, projectFlag,
					&cli.StringFlag{Name: "followers", Usage: "Comma-separated follower ZPUIDs"},
					&cli.StringFlag{Name: "json", Usage: "Additional fields as JSON"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					portal, err := requirePortal(cmd)
					if err != nil {
						return err
					}
					body := map[string]any{}
					if cmd.IsSet("followers") {
						var arr any
						if err := json.Unmarshal([]byte(cmd.String("followers")), &arr); err != nil {
							return internal.NewValidationError(fmt.Sprintf("--followers: invalid JSON: %v", err))
						}
						body["followers"] = arr
					}
					if err := internal.MergeJSON(cmd, body); err != nil {
						return err
					}
					url := base(c, portal, cmd.String("project")) + "/issues/" + cmd.Args().First() + "/followers"
					raw, err := c.Request("DELETE", url, &zohttp.RequestOpts{JSON: body})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
		},
	}
}

func issueLinkingCmd() *cli.Command {
	return &cli.Command{
		Name:  "issue-linking",
		Usage: "Issue linking operations",
		Commands: []*cli.Command{
			{
				Name:      "list",
				Usage:     "List linked issues",
				ArgsUsage: "<issue-id>",
				Flags:     []cli.Flag{portalFlag, projectFlag},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					portal, err := requirePortal(cmd)
					if err != nil {
						return err
					}
					url := base(c, portal, cmd.String("project")) + "/issues/" + cmd.Args().First() + "/linkedissues"
					raw, err := c.Request("GET", url, nil)
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "link",
				Usage:     "Link an issue",
				ArgsUsage: "<issue-id>",
				Flags: []cli.Flag{
					portalFlag, projectFlag,
					&cli.StringFlag{Name: "link_type", Usage: "Link type (relate, blocks, is_blocked_by, duplicate)"},
					&cli.StringFlag{Name: "issue_ids", Usage: "Issue IDs as JSON array"},
					&cli.StringFlag{Name: "json", Usage: "Additional fields as JSON"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					portal, err := requirePortal(cmd)
					if err != nil {
						return err
					}
					body := map[string]any{}
					if cmd.IsSet("link_type") {
						body["link_type"] = cmd.String("link_type")
					}
					if cmd.IsSet("issue_ids") {
						var v any
						if err := json.Unmarshal([]byte(cmd.String("issue_ids")), &v); err != nil {
							return internal.NewValidationError(fmt.Sprintf("--issue_ids: invalid JSON: %v", err))
						}
						body["issue_ids"] = v
					}
					if err := internal.MergeJSON(cmd, body); err != nil {
						return err
					}
					url := base(c, portal, cmd.String("project")) + "/issues/" + cmd.Args().First() + "/link"
					raw, err := c.Request("POST", url, &zohttp.RequestOpts{JSON: body})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:  "bulk-link",
				Usage: "Bulk link issues",
				Flags: []cli.Flag{
					portalFlag, projectFlag,
					&cli.StringFlag{Name: "link_type", Usage: "Link type (relate, blocks, is_blocked_by, duplicate)"},
					&cli.StringFlag{Name: "issue_ids", Usage: "Source issue IDs as JSON array"},
					&cli.StringFlag{Name: "linking_issue_ids", Usage: "Target issue IDs as JSON array"},
					&cli.StringFlag{Name: "json", Usage: "Additional fields as JSON"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					portal, err := requirePortal(cmd)
					if err != nil {
						return err
					}
					body := map[string]any{}
					if cmd.IsSet("link_type") {
						body["link_type"] = cmd.String("link_type")
					}
					if cmd.IsSet("issue_ids") {
						var v any
						if err := json.Unmarshal([]byte(cmd.String("issue_ids")), &v); err != nil {
							return internal.NewValidationError(fmt.Sprintf("--issue_ids: invalid JSON: %v", err))
						}
						body["issue_ids"] = v
					}
					if cmd.IsSet("linking_issue_ids") {
						var v any
						if err := json.Unmarshal([]byte(cmd.String("linking_issue_ids")), &v); err != nil {
							return internal.NewValidationError(fmt.Sprintf("--linking_issue_ids: invalid JSON: %v", err))
						}
						body["linking_issue_ids"] = v
					}
					if err := internal.MergeJSON(cmd, body); err != nil {
						return err
					}
					url := base(c, portal, cmd.String("project")) + "/issues/bulk-link-bugs"
					raw, err := c.Request("POST", url, &zohttp.RequestOpts{JSON: body})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "change-type",
				Usage:     "Change link type",
				ArgsUsage: "<issue-id> <link-id>",
				Flags: []cli.Flag{
					portalFlag, projectFlag,
					&cli.StringFlag{Name: "link_type", Usage: "New link type (relate, blocks, is_blocked_by, duplicate)"},
					&cli.StringFlag{Name: "json", Usage: "Additional fields as JSON"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					portal, err := requirePortal(cmd)
					if err != nil {
						return err
					}
					body := map[string]any{}
					if cmd.IsSet("link_type") {
						body["link_type"] = cmd.String("link_type")
					}
					if err := internal.MergeJSON(cmd, body); err != nil {
						return err
					}
					url := base(c, portal, cmd.String("project")) + "/issues/" + cmd.Args().Get(0) + "/link/" + cmd.Args().Get(1)
					raw, err := c.Request("POST", url, &zohttp.RequestOpts{JSON: body})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "unlink",
				Usage:     "Unlink an issue",
				ArgsUsage: "<issue-id> <link-id>",
				Flags:     []cli.Flag{portalFlag, projectFlag},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					portal, err := requirePortal(cmd)
					if err != nil {
						return err
					}
					url := base(c, portal, cmd.String("project")) + "/issues/" + cmd.Args().Get(0) + "/link/" + cmd.Args().Get(1)
					raw, err := c.Request("DELETE", url, nil)
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
		},
	}
}

func issueResolutionCmd() *cli.Command {
	return &cli.Command{
		Name:  "issue-resolution",
		Usage: "Issue resolution operations",
		Commands: []*cli.Command{
			{
				Name:      "get",
				Usage:     "Get issue resolution",
				ArgsUsage: "<issue-id>",
				Flags:     []cli.Flag{portalFlag, projectFlag},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					portal, err := requirePortal(cmd)
					if err != nil {
						return err
					}
					url := base(c, portal, cmd.String("project")) + "/issues/" + cmd.Args().First() + "/resolution"
					raw, err := c.Request("GET", url, nil)
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "add",
				Usage:     "Add issue resolution",
				ArgsUsage: "<issue-id>",
				Flags: []cli.Flag{
					portalFlag, projectFlag,
					&cli.StringFlag{Name: "resolution", Usage: "Resolution text"},
					&cli.StringFlag{Name: "json", Usage: "Additional fields as JSON"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					portal, err := requirePortal(cmd)
					if err != nil {
						return err
					}
					body := map[string]any{}
					if cmd.IsSet("resolution") {
						body["resolution"] = cmd.String("resolution")
					}
					if err := internal.MergeJSON(cmd, body); err != nil {
						return err
					}
					url := base(c, portal, cmd.String("project")) + "/issues/" + cmd.Args().First() + "/resolution"
					raw, err := c.Request("POST", url, &zohttp.RequestOpts{JSON: body})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "update",
				Usage:     "Update issue resolution",
				ArgsUsage: "<issue-id>",
				Flags: []cli.Flag{
					portalFlag, projectFlag,
					&cli.StringFlag{Name: "resolution", Usage: "Resolution text"},
					&cli.StringFlag{Name: "json", Usage: "Additional fields as JSON"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					portal, err := requirePortal(cmd)
					if err != nil {
						return err
					}
					body := map[string]any{}
					if cmd.IsSet("resolution") {
						body["resolution"] = cmd.String("resolution")
					}
					if err := internal.MergeJSON(cmd, body); err != nil {
						return err
					}
					url := base(c, portal, cmd.String("project")) + "/issues/" + cmd.Args().First() + "/resolution"
					raw, err := c.Request("PUT", url, &zohttp.RequestOpts{JSON: body})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "delete",
				Usage:     "Delete issue resolution",
				ArgsUsage: "<issue-id>",
				Flags:     []cli.Flag{portalFlag, projectFlag},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					portal, err := requirePortal(cmd)
					if err != nil {
						return err
					}
					url := base(c, portal, cmd.String("project")) + "/issues/" + cmd.Args().First() + "/resolution"
					raw, err := c.Request("DELETE", url, nil)
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
		},
	}
}

func issueAttachmentsCmd() *cli.Command {
	return &cli.Command{
		Name:  "issue-attachments",
		Usage: "Issue attachment operations",
		Commands: []*cli.Command{
			{
				Name:      "list",
				Usage:     "List issue attachments",
				ArgsUsage: "<issue-id>",
				Flags:     []cli.Flag{portalFlag, projectFlag},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					portal, err := requirePortal(cmd)
					if err != nil {
						return err
					}
					url := base(c, portal, cmd.String("project")) + "/issues/" + cmd.Args().First() + "/attachments"
					raw, err := c.Request("GET", url, nil)
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "associate",
				Usage:     "Associate attachments to an issue",
				ArgsUsage: "<issue-id>",
				Flags: []cli.Flag{
					portalFlag, projectFlag,
					&cli.StringFlag{Name: "attachment-ids", Usage: "Attachment IDs as JSON array"},
					&cli.StringFlag{Name: "json", Usage: "Additional fields as JSON"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					portal, err := requirePortal(cmd)
					if err != nil {
						return err
					}
					form := map[string]string{}
					if cmd.IsSet("attachment-ids") {
						form["attachment_ids"] = cmd.String("attachment-ids")
					}
					if err := internal.MergeJSONForm(cmd, form); err != nil {
						return err
					}
					url := base(c, portal, cmd.String("project")) + "/issues/" + cmd.Args().First() + "/attachments"
					raw, err := c.Request("POST", url, &zohttp.RequestOpts{
						Form: form,
					})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "dissociate",
				Usage:     "Dissociate an attachment from an issue",
				ArgsUsage: "<issue-id> <attachment-id>",
				Flags:     []cli.Flag{portalFlag, projectFlag},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					portal, err := requirePortal(cmd)
					if err != nil {
						return err
					}
					url := base(c, portal, cmd.String("project")) + "/issues/" + cmd.Args().Get(0) + "/attachments/" + cmd.Args().Get(1)
					raw, err := c.Request("DELETE", url, nil)
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
		},
	}
}

func issueCustomViewsCmd() *cli.Command {
	return &cli.Command{
		Name:  "issue-customviews",
		Usage: "Issue custom view operations",
		Commands: []*cli.Command{
			{
				Name:  "list",
				Usage: "List issue custom views (portal-level)",
				Flags: []cli.Flag{portalFlag},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					portal, err := requirePortal(cmd)
					if err != nil {
						return err
					}
					url := c.ProjectsBase + "/portal/" + portal + "/issues/customviews"
					raw, err := c.Request("GET", url, nil)
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:  "project-list",
				Usage: "List issue custom views (project-level)",
				Flags: []cli.Flag{portalFlag, projectFlag},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					portal, err := requirePortal(cmd)
					if err != nil {
						return err
					}
					url := base(c, portal, cmd.String("project")) + "/issues/customviews"
					raw, err := c.Request("GET", url, nil)
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "get",
				Usage:     "Get an issue custom view",
				ArgsUsage: "<view-id>",
				Flags:     []cli.Flag{portalFlag},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					portal, err := requirePortal(cmd)
					if err != nil {
						return err
					}
					url := c.ProjectsBase + "/portal/" + portal + "/issues/customviews/" + cmd.Args().First()
					raw, err := c.Request("GET", url, nil)
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
		},
	}
}
