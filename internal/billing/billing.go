package billing

import (
	"github.com/omin8tor/zoho-cli/internal/auth"
	zohttp "github.com/omin8tor/zoho-cli/internal/http"
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
		Name:  "billing",
		Usage: "Zoho Billing operations",
		Commands: []*cli.Command{},
	}
}
