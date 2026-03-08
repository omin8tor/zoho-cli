package main

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/omin8tor/zoho-cli/internal"
	"github.com/omin8tor/zoho-cli/internal/auth"
	"github.com/omin8tor/zoho-cli/internal/bigin"
	"github.com/omin8tor/zoho-cli/internal/billing"
	"github.com/omin8tor/zoho-cli/internal/books"
	"github.com/omin8tor/zoho-cli/internal/cliq"
	"github.com/omin8tor/zoho-cli/internal/creator"
	"github.com/omin8tor/zoho-cli/internal/crm"
	"github.com/omin8tor/zoho-cli/internal/desk"
	"github.com/omin8tor/zoho-cli/internal/drive"
	"github.com/omin8tor/zoho-cli/internal/expense"
	"github.com/omin8tor/zoho-cli/internal/inventory"
	"github.com/omin8tor/zoho-cli/internal/invoice"
	"github.com/omin8tor/zoho-cli/internal/mail"
	"github.com/omin8tor/zoho-cli/internal/output"
	"github.com/omin8tor/zoho-cli/internal/people"
	"github.com/omin8tor/zoho-cli/internal/projects"
	"github.com/omin8tor/zoho-cli/internal/recruit"
	"github.com/omin8tor/zoho-cli/internal/sheet"
	"github.com/omin8tor/zoho-cli/internal/sign"
	"github.com/omin8tor/zoho-cli/internal/sprints"
	"github.com/omin8tor/zoho-cli/internal/writer"
	"github.com/urfave/cli/v3"
)

var version = "dev"

func main() {
	app := &cli.Command{
		Name:    "zoho",
		Usage:   "CLI for Zoho REST APIs",
		Version: version,
		Flags: []cli.Flag{
			&cli.BoolFlag{Name: "help-all", Usage: "Show help for all commands recursively"},
		},
		Commands: []*cli.Command{
			auth.Commands(),
			bigin.Commands(),
			billing.Commands(),
			books.Commands(),
			cliq.Commands(),
			creator.Commands(),
			crm.Commands(),
			desk.Commands(),
			drive.Commands(),
			expense.Commands(),
			inventory.Commands(),
			invoice.Commands(),
			mail.Commands(),
			people.Commands(),
			projects.Commands(),
			recruit.Commands(),
			sheet.Commands(),
			sign.Commands(),
			sprints.Commands(),
			writer.Commands(),
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
		var apiErr *internal.ZohoAPIError
		var cliErr *internal.ZohoCliError
		if errors.As(err, &apiErr) {
			fmt.Fprintln(os.Stderr, apiErr.Message)
			os.Exit(apiErr.ExitCode)
		}
		if errors.As(err, &cliErr) {
			fmt.Fprintln(os.Stderr, cliErr.Message)
			os.Exit(cliErr.ExitCode)
		}
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}
