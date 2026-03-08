package auth

import (
	"context"
	"fmt"
	"os"

	"github.com/omin8tor/zoho-cli/internal"
	"github.com/omin8tor/zoho-cli/internal/output"
	"github.com/urfave/cli/v3"
)

func Commands() *cli.Command {
	return &cli.Command{
		Name:  "auth",
		Usage: "Authentication management",
		Commands: []*cli.Command{
			{
				Name:  "login",
				Usage: "Authenticate via device flow OAuth",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "client-id", Required: true, Usage: "Zoho OAuth client ID"},
					&cli.StringFlag{Name: "client-secret", Required: true, Usage: "Zoho OAuth client secret"},
					&cli.StringFlag{Name: "dc", Value: "com", Usage: "Data center (com, eu, in, com.au, jp, ca, sa, uk, com.cn)"},
					&cli.StringFlag{Name: "scopes", Usage: "Comma-separated OAuth scopes (defaults to all)"},
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					return DeviceFlowLogin(
						cmd.String("client-id"),
						cmd.String("client-secret"),
						cmd.String("dc"),
						cmd.String("scopes"),
					)
				},
			},
			{
				Name:  "self-client",
				Usage: "Authenticate via self-client code exchange",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "code", Required: true, Usage: "Self-client authorization code from Zoho API Console"},
					&cli.StringFlag{Name: "client-id", Required: true, Usage: "Zoho OAuth client ID"},
					&cli.StringFlag{Name: "client-secret", Required: true, Usage: "Zoho OAuth client secret"},
					&cli.StringFlag{Name: "dc", Value: "com", Usage: "Data center"},
					&cli.StringFlag{Name: "server", Usage: "Accounts server URL override"},
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					return SelfClientExchange(
						cmd.String("client-id"),
						cmd.String("client-secret"),
						cmd.String("code"),
						cmd.String("dc"),
						cmd.String("server"),
					)
				},
			},
			{
				Name:  "status",
				Usage: "Show current authentication status",
				Action: func(ctx context.Context, cmd *cli.Command) error {
					config, err := ResolveAuth()
					if err != nil {
						return internal.NewAuthError("Not authenticated. Run `zoho auth login` to authenticate.")
					}
					return output.JSON(map[string]any{
						"authenticated": true,
						"source":        config.Source,
						"dc":            config.DC,
						"accounts_url":  config.AccountsURL,
						"token_valid":   config.TokenValid(),
					})
				},
			},
			{
				Name:  "refresh",
				Usage: "Force refresh the access token",
				Action: func(ctx context.Context, cmd *cli.Command) error {
					config, err := ResolveAuth()
					if err != nil {
						return err
					}
					token, err := RefreshAccessToken(config)
					if err != nil {
						return err
					}
					fmt.Fprintln(os.Stderr, "Access token refreshed successfully.")
					truncated := token
					if len(truncated) > 20 {
						truncated = truncated[:20] + "..."
					}
					return output.JSON(map[string]string{"access_token": truncated, "status": "refreshed"})
				},
			},
			{
				Name:  "logout",
				Usage: "Clear stored authentication tokens",
				Action: func(ctx context.Context, cmd *cli.Command) error {
					removed := false
					if _, err := os.Stat(TokensFile()); err == nil {
						os.Remove(TokensFile())
						removed = true
					}
					if _, err := os.Stat(ConfigFile()); err == nil {
						os.Remove(ConfigFile())
						removed = true
					}
					if removed {
						fmt.Fprintln(os.Stderr, "Logged out. Stored credentials removed.")
					} else {
						fmt.Fprintln(os.Stderr, "No stored credentials found.")
					}
					return nil
				},
			},
		},
	}
}
