package main

import (
	"context"
	"fmt"
	"os"

	"github.com/omin8tor/zoho-cli/internal"
	"github.com/omin8tor/zoho-cli/internal/auth"
	"github.com/omin8tor/zoho-cli/internal/cliq"
	"github.com/omin8tor/zoho-cli/internal/crm"
	"github.com/omin8tor/zoho-cli/internal/desk"
	"github.com/omin8tor/zoho-cli/internal/drive"
	"github.com/omin8tor/zoho-cli/internal/expense"
	"github.com/omin8tor/zoho-cli/internal/mail"
	"github.com/omin8tor/zoho-cli/internal/output"
	"github.com/omin8tor/zoho-cli/internal/projects"
	"github.com/omin8tor/zoho-cli/internal/sheet"
	"github.com/omin8tor/zoho-cli/internal/sign"
	"github.com/omin8tor/zoho-cli/internal/writer"
	"github.com/urfave/cli/v3"
)

var version = "dev"

func main() {
	app := &cli.Command{
		Name:    "zoho",
		Usage:   "CLI for Zoho REST APIs (CRM, Projects, WorkDrive, Writer, Cliq, Expense, Sheet, Desk, Sign, Mail)",
		Version: version,
		Flags: []cli.Flag{
			&cli.BoolFlag{Name: "help-all", Usage: "Show help for all commands recursively"},
		},
		Commands: []*cli.Command{
			auth.Commands(),
			crm.Commands(),
			projects.Commands(),
			drive.Commands(),
			writer.Commands(),
			cliq.Commands(),
			expense.Commands(),
			mail.Commands(),
			sheet.Commands(),
			desk.Commands(),
			sign.Commands(),
		},
	}

	app.Action = func(_ context.Context, cmd *cli.Command) error {
		if cmd.Bool("help-all") {
			return output.PrintHelpAll(cmd)
		}
		cmd.Root().Run(context.Background(), []string{"zoho", "--help"})
		return nil
	}

	if err := app.Run(context.Background(), os.Args); err != nil {
		if e, ok := err.(*internal.ZohoCliError); ok {
			fmt.Fprintln(os.Stderr, e.Message)
			os.Exit(e.ExitCode)
		}
		if e, ok := err.(*internal.ZohoAPIError); ok {
			fmt.Fprintln(os.Stderr, e.Message)
			os.Exit(e.ExitCode)
		}
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}
