package cliq

import (
	"context"
	"encoding/json"
	"fmt"

	zohttp "github.com/omin8tor/zoho-cli/internal/http"
	"github.com/omin8tor/zoho-cli/internal/output"
	"github.com/urfave/cli/v3"
)

func Commands() *cli.Command {
	return &cli.Command{
		Name:  "cliq",
		Usage: "Zoho Cliq operations",
		Commands: []*cli.Command{
			channelsCmd(),
			chatsCmd(),
			buddiesCmd(),
			messagesCmd(),
			usersCmd(),
		},
	}
}

func channelsCmd() *cli.Command {
	return &cli.Command{
		Name:  "channels",
		Usage: "Cliq channel operations",
		Commands: []*cli.Command{
			{
				Name:  "list",
				Usage: "List channels",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					raw, err := c.Request("GET", c.CliqBase+"/api/v2/channels", nil)
					if err != nil {
						return err
					}
					var data map[string]json.RawMessage
					if err := json.Unmarshal(raw, &data); err == nil {
						if d, ok := data["data"]; ok {
							return output.JSONRaw(d)
						}
						if ch, ok := data["channels"]; ok {
							return output.JSONRaw(ch)
						}
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "get",
				Usage:     "Get channel info",
				ArgsUsage: "<channel-name>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					raw, err := c.Request("GET", c.CliqBase+"/api/v2/channelsbyname/"+cmd.Args().First(), nil)
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:  "create",
				Usage: "Create a channel",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "name", Required: true, Usage: "Channel name"},
					&cli.StringFlag{Name: "description", Usage: "Channel description"},
					&cli.StringFlag{Name: "level", Usage: "Channel level (organization, team, private, external)"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					body := map[string]any{"name": cmd.String("name")}
					if d := cmd.String("description"); d != "" {
						body["description"] = d
					}
					if l := cmd.String("level"); l != "" {
						body["level"] = l
					}
					raw, err := c.Request("POST", c.CliqBase+"/api/v2/channels", &zohttp.RequestOpts{JSON: body})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "message",
				Usage:     "Send a message to a channel",
				ArgsUsage: "<channel-name>",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "text", Required: true, Usage: "Message text"},
					&cli.StringFlag{Name: "bot", Usage: "Bot name to send as"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					body := map[string]any{"text": cmd.String("text")}
					if b := cmd.String("bot"); b != "" {
						body["bot"] = map[string]string{"name": b}
					}
					raw, err := c.Request("POST", c.CliqBase+"/api/v2/channelsbyname/"+cmd.Args().First()+"/message", &zohttp.RequestOpts{JSON: body})
					if err != nil {
						return err
					}
					return output.JSON(map[string]any{"ok": true, "channel": cmd.Args().First(), "response": json.RawMessage(raw)})
				},
			},
			{
				Name:      "members",
				Usage:     "List channel members",
				ArgsUsage: "<channel-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					raw, err := c.Request("GET", c.CliqBase+"/api/v2/channels/"+cmd.Args().First()+"/members", nil)
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "delete",
				Usage:     "Delete a channel",
				ArgsUsage: "<channel-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					raw, err := c.Request("DELETE", c.CliqBase+"/api/v2/channels/"+cmd.Args().First(), nil)
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
		},
	}
}

func chatsCmd() *cli.Command {
	return &cli.Command{
		Name:  "chats",
		Usage: "Cliq chat operations",
		Commands: []*cli.Command{
			{
				Name:      "message",
				Usage:     "Send a direct message",
				ArgsUsage: "<chat-id>",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "text", Required: true, Usage: "Message text"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					raw, err := c.Request("POST", c.CliqBase+"/api/v2/chats/"+cmd.Args().First()+"/message", &zohttp.RequestOpts{
						JSON: map[string]string{"text": cmd.String("text")},
					})
					if err != nil {
						return err
					}
					return output.JSON(map[string]any{"ok": true, "chat_id": cmd.Args().First(), "response": json.RawMessage(raw)})
				},
			},
		},
	}
}

func buddiesCmd() *cli.Command {
	return &cli.Command{
		Name:  "buddies",
		Usage: "Cliq buddy/DM operations",
		Commands: []*cli.Command{
			{
				Name:      "message",
				Usage:     "Send a DM by email address",
				ArgsUsage: "<email>",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "text", Required: true, Usage: "Message text"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					raw, err := c.Request("POST", c.CliqBase+"/api/v2/buddies/"+cmd.Args().First()+"/message", &zohttp.RequestOpts{
						JSON: map[string]string{"text": cmd.String("text")},
					})
					if err != nil {
						return err
					}
					return output.JSON(map[string]any{"ok": true, "email": cmd.Args().First(), "response": json.RawMessage(raw)})
				},
			},
		},
	}
}

func messagesCmd() *cli.Command {
	return &cli.Command{
		Name:  "messages",
		Usage: "Cliq message operations",
		Commands: []*cli.Command{
			{
				Name:      "list",
				Usage:     "List messages in a chat",
				ArgsUsage: "<chat-id>",
				Flags: []cli.Flag{
					&cli.IntFlag{Name: "limit", Value: 50, Usage: "Number of messages"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					raw, err := c.Request("GET", c.CliqBase+"/api/v2/chats/"+cmd.Args().First()+"/messages", &zohttp.RequestOpts{
						Params: map[string]string{"limit": fmt.Sprintf("%d", cmd.Int("limit"))},
					})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "edit",
				Usage:     "Edit a message",
				ArgsUsage: "<chat-id> <message-id>",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "text", Required: true, Usage: "New message text"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					chatID, msgID := cmd.Args().Get(0), cmd.Args().Get(1)
					raw, err := c.Request("PUT", c.CliqBase+"/api/v2/chats/"+chatID+"/messages/"+msgID, &zohttp.RequestOpts{
						JSON: map[string]string{"text": cmd.String("text")},
					})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "delete",
				Usage:     "Delete a message",
				ArgsUsage: "<chat-id> <message-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					chatID, msgID := cmd.Args().Get(0), cmd.Args().Get(1)
					raw, err := c.Request("DELETE", c.CliqBase+"/api/v2/chats/"+chatID+"/messages/"+msgID, nil)
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
		Usage: "Cliq user operations",
		Commands: []*cli.Command{
			{
				Name:  "list",
				Usage: "List users in the organization",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					raw, err := c.Request("GET", c.CliqBase+"/api/v2/users", nil)
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "get",
				Usage:     "Get user details",
				ArgsUsage: "<user-id-or-email>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := zohttp.GetClient()
					if err != nil {
						return err
					}
					raw, err := c.Request("GET", c.CliqBase+"/api/v2/users/"+cmd.Args().First(), nil)
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
		},
	}
}
