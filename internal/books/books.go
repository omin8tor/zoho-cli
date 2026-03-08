package books

import (
	"github.com/omin8tor/zoho-cli/internal"
	"github.com/urfave/cli/v3"
)

func resolveOrgID(cmd *cli.Command) (string, error) {
	return internal.RequireFlag(cmd, "org", "ZOHO_BOOKS_ORG_ID")
}

func orgParams(orgID string) map[string]string {
	return map[string]string{"organization_id": orgID}
}

func mergeParams(orgID string, extra map[string]string) map[string]string {
	params := orgParams(orgID)
	for k, v := range extra {
		params[k] = v
	}
	return params
}

func Commands() *cli.Command {
	return &cli.Command{
		Name:  "books",
		Usage: "Zoho Books operations",
		Flags: []cli.Flag{
			&cli.StringFlag{Name: "org", Usage: "Organization ID (or set ZOHO_BOOKS_ORG_ID)", Sources: cli.EnvVars("ZOHO_BOOKS_ORG_ID")},
		},
		Commands: []*cli.Command{
			organizationsCmd(),
			contactsCmd(),
			contactPersonsCmd(),
			estimatesCmd(),
			salesOrdersCmd(),
			salesReceiptsCmd(),
			invoicesCmd(),
			recurringInvoicesCmd(),
			creditNotesCmd(),
			customerDebitNotesCmd(),
			customerPaymentsCmd(),
			expensesCmd(),
			recurringExpensesCmd(),
			retainerInvoicesCmd(),
			purchaseOrdersCmd(),
			billsCmd(),
			recurringBillsCmd(),
			vendorCreditsCmd(),
			vendorPaymentsCmd(),
			customModulesCmd(),
			bankAccountsCmd(),
			bankTransactionsCmd(),
			bankRulesCmd(),
			chartOfAccountsCmd(),
			journalsCmd(),
			fixedAssetsCmd(),
			baseCurrencyAdjCmd(),
			booksProjectsCmd(),
			tasksCmd(),
			timeEntriesCmd(),
			usersCmd(),
			itemsCmd(),
			locationsCmd(),
			currenciesCmd(),
			taxesCmd(),
			openingBalanceCmd(),
			crmIntegrationCmd(),
			reportingTagsCmd(),
		},
	}
}
