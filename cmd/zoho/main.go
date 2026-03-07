package main

import (
	"context"
	"fmt"
	"os"

	"github.com/omin8tor/zoho-cli/internal"
	"github.com/omin8tor/zoho-cli/internal/analytics"
	"github.com/omin8tor/zoho-cli/internal/assist"
	"github.com/omin8tor/zoho-cli/internal/auth"
	"github.com/omin8tor/zoho-cli/internal/backstage"
	"github.com/omin8tor/zoho-cli/internal/bigin"
	"github.com/omin8tor/zoho-cli/internal/billing"
	"github.com/omin8tor/zoho-cli/internal/bookings"
	"github.com/omin8tor/zoho-cli/internal/books"
	"github.com/omin8tor/zoho-cli/internal/campaigns"
	"github.com/omin8tor/zoho-cli/internal/cliq"
	"github.com/omin8tor/zoho-cli/internal/creator"
	"github.com/omin8tor/zoho-cli/internal/crm"
	"github.com/omin8tor/zoho-cli/internal/desk"
	"github.com/omin8tor/zoho-cli/internal/drive"
	"github.com/omin8tor/zoho-cli/internal/expense"
	"github.com/omin8tor/zoho-cli/internal/inventory"
	"github.com/omin8tor/zoho-cli/internal/invoice"
	"github.com/omin8tor/zoho-cli/internal/learn"
	"github.com/omin8tor/zoho-cli/internal/mail"
	"github.com/omin8tor/zoho-cli/internal/marketingauto"
	"github.com/omin8tor/zoho-cli/internal/meeting"
	"github.com/omin8tor/zoho-cli/internal/output"
	"github.com/omin8tor/zoho-cli/internal/pagesense"
	"github.com/omin8tor/zoho-cli/internal/people"
	"github.com/omin8tor/zoho-cli/internal/projects"
	"github.com/omin8tor/zoho-cli/internal/recruit"
	"github.com/omin8tor/zoho-cli/internal/salesiq"
	"github.com/omin8tor/zoho-cli/internal/sheet"
	"github.com/omin8tor/zoho-cli/internal/showtime"
	"github.com/omin8tor/zoho-cli/internal/sign"
	"github.com/omin8tor/zoho-cli/internal/sprints"
	"github.com/omin8tor/zoho-cli/internal/vault"
	"github.com/omin8tor/zoho-cli/internal/voice"
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
			analytics.Commands(),
			assist.Commands(),
			auth.Commands(),
			backstage.Commands(),
			bigin.Commands(),
			billing.Commands(),
			bookings.Commands(),
			books.Commands(),
			campaigns.Commands(),
			cliq.Commands(),
			creator.Commands(),
			crm.Commands(),
			desk.Commands(),
			drive.Commands(),
			expense.Commands(),
			inventory.Commands(),
			invoice.Commands(),
			learn.Commands(),
			mail.Commands(),
			marketingauto.Commands(),
			meeting.Commands(),
			pagesense.Commands(),
			people.Commands(),
			projects.Commands(),
			recruit.Commands(),
			salesiq.Commands(),
			sheet.Commands(),
			showtime.Commands(),
			sign.Commands(),
			sprints.Commands(),
			vault.Commands(),
			voice.Commands(),
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
