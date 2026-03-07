//go:build integration

package internal_test

import (
	"fmt"
	"os"
	"strings"
	"testing"
)

func (c *testCleanup) trackBooksContact(id, orgID string) {
	c.add("delete books contact "+id, func() {
		zohoIgnoreError(c.t, "books", "contacts", "delete", id, "--org", orgID)
	})
}

func (c *testCleanup) trackBooksItem(id, orgID string) {
	c.add("delete books item "+id, func() {
		zohoIgnoreError(c.t, "books", "items", "delete", id, "--org", orgID)
	})
}

func (c *testCleanup) trackBooksInvoice(id, orgID string) {
	c.add("delete books invoice "+id, func() {
		zohoIgnoreError(c.t, "books", "invoices", "delete", id, "--org", orgID)
	})
}

func (c *testCleanup) trackBooksEstimate(id, orgID string) {
	c.add("delete books estimate "+id, func() {
		zohoIgnoreError(c.t, "books", "estimates", "delete", id, "--org", orgID)
	})
}

func (c *testCleanup) trackBooksExpense(id, orgID string) {
	c.add("delete books expense "+id, func() {
		zohoIgnoreError(c.t, "books", "expenses", "delete", id, "--org", orgID)
	})
}

func requireBooksOrgID(t *testing.T) string {
	t.Helper()
	id := os.Getenv("ZOHO_BOOKS_ORG_ID")
	if id != "" {
		return id
	}
	out, err := zohoMayFail(t, "books", "organizations", "list")
	if err != nil {
		t.Skipf("skipping: cannot discover books org ID: %v", err)
	}
	m := parseJSON(t, out)
	orgs, ok := m["organizations"].([]any)
	if !ok || len(orgs) == 0 {
		t.Skip("skipping: no books organizations found")
	}
	org := orgs[0].(map[string]any)
	orgID := fmt.Sprintf("%v", org["organization_id"])
	if orgID == "" || orgID == "<nil>" {
		t.Skip("skipping: books organization_id is empty")
	}
	return orgID
}

func assertBooksCodeZero(t *testing.T, m map[string]any) {
	t.Helper()
	if code, ok := m["code"].(float64); ok && code != 0 {
		msg := fmt.Sprintf("%v", m["message"])
		t.Fatalf("books API error: code=%.0f message=%s", code, msg)
	}
}

func TestBooksOrganizations(t *testing.T) {
	t.Parallel()
	orgID := requireBooksOrgID(t)

	t.Run("list", func(t *testing.T) {
		out := zoho(t, "books", "organizations", "list")
		m := parseJSON(t, out)
		assertBooksCodeZero(t, m)
		arr, ok := m["organizations"].([]any)
		if !ok {
			t.Fatalf("expected organizations array in response:\n%s", truncate(out, 500))
		}
		if len(arr) == 0 {
			t.Fatal("expected at least one organization")
		}
	})

	t.Run("get", func(t *testing.T) {
		out := zoho(t, "books", "organizations", "get", orgID)
		m := parseJSON(t, out)
		assertBooksCodeZero(t, m)
		org, ok := m["organization"].(map[string]any)
		if !ok {
			t.Fatalf("expected organization object in response:\n%s", truncate(out, 500))
		}
		id := fmt.Sprintf("%v", org["organization_id"])
		if id == "" || id == "<nil>" {
			t.Fatalf("expected organization_id in response:\n%s", truncate(out, 500))
		}
	})

	_ = orgID
}

func TestBooksContacts(t *testing.T) {
	t.Parallel()
	orgID := requireBooksOrgID(t)
	cleanup := newCleanup(t)

	var contactID string

	t.Run("list", func(t *testing.T) {
		out := zoho(t, "books", "contacts", "list", "--org", orgID)
		m := parseJSON(t, out)
		assertBooksCodeZero(t, m)
	})

	t.Run("create", func(t *testing.T) {
		name := fmt.Sprintf("%s Contact %s", testPrefix, randomSuffix())
		out := zoho(t, "books", "contacts", "create", "--org", orgID,
			"--json", toJSON(t, map[string]any{"contact_name": name, "contact_type": "customer"}))
		m := parseJSON(t, out)
		assertBooksCodeZero(t, m)
		contact, ok := m["contact"].(map[string]any)
		if !ok {
			t.Fatalf("expected contact object:\n%s", truncate(out, 500))
		}
		contactID = fmt.Sprintf("%v", contact["contact_id"])
		if contactID == "" || contactID == "<nil>" {
			t.Fatalf("expected contact_id:\n%s", truncate(out, 500))
		}
		cleanup.trackBooksContact(contactID, orgID)
	})

	t.Run("get", func(t *testing.T) {
		requireID(t, contactID, "create must have succeeded")
		out := zoho(t, "books", "contacts", "get", contactID, "--org", orgID)
		m := parseJSON(t, out)
		assertBooksCodeZero(t, m)
		contact := m["contact"].(map[string]any)
		assertEqual(t, fmt.Sprintf("%v", contact["contact_id"]), contactID)
	})

	t.Run("update", func(t *testing.T) {
		requireID(t, contactID, "create must have succeeded")
		updatedName := fmt.Sprintf("%s Contact Upd %s", testPrefix, randomSuffix())
		out := zoho(t, "books", "contacts", "update", contactID, "--org", orgID,
			"--json", toJSON(t, map[string]any{"contact_name": updatedName}))
		m := parseJSON(t, out)
		assertBooksCodeZero(t, m)
	})

	t.Run("mark-inactive", func(t *testing.T) {
		requireID(t, contactID, "create must have succeeded")
		out := zoho(t, "books", "contacts", "mark-inactive", contactID, "--org", orgID)
		m := parseJSON(t, out)
		assertBooksCodeZero(t, m)
	})

	t.Run("mark-active", func(t *testing.T) {
		requireID(t, contactID, "create must have succeeded")
		out := zoho(t, "books", "contacts", "mark-active", contactID, "--org", orgID)
		m := parseJSON(t, out)
		assertBooksCodeZero(t, m)
	})

	t.Run("delete", func(t *testing.T) {
		requireID(t, contactID, "create must have succeeded")
		out := zoho(t, "books", "contacts", "delete", contactID, "--org", orgID)
		m := parseJSON(t, out)
		assertBooksCodeZero(t, m)
	})
}

func TestBooksItems(t *testing.T) {
	t.Parallel()
	orgID := requireBooksOrgID(t)
	cleanup := newCleanup(t)

	var itemID string

	t.Run("list", func(t *testing.T) {
		out := zoho(t, "books", "items", "list", "--org", orgID)
		m := parseJSON(t, out)
		assertBooksCodeZero(t, m)
	})

	t.Run("create", func(t *testing.T) {
		name := fmt.Sprintf("%s Item %s", testPrefix, randomSuffix())
		out := zoho(t, "books", "items", "create", "--org", orgID,
			"--json", toJSON(t, map[string]any{"name": name, "rate": 100}))
		m := parseJSON(t, out)
		assertBooksCodeZero(t, m)
		item, ok := m["item"].(map[string]any)
		if !ok {
			t.Fatalf("expected item object:\n%s", truncate(out, 500))
		}
		itemID = fmt.Sprintf("%v", item["item_id"])
		if itemID == "" || itemID == "<nil>" {
			t.Fatalf("expected item_id:\n%s", truncate(out, 500))
		}
		cleanup.trackBooksItem(itemID, orgID)
	})

	t.Run("get", func(t *testing.T) {
		requireID(t, itemID, "create must have succeeded")
		out := zoho(t, "books", "items", "get", itemID, "--org", orgID)
		m := parseJSON(t, out)
		assertBooksCodeZero(t, m)
		item := m["item"].(map[string]any)
		assertEqual(t, fmt.Sprintf("%v", item["item_id"]), itemID)
	})

	t.Run("update", func(t *testing.T) {
		requireID(t, itemID, "create must have succeeded")
		updatedName := fmt.Sprintf("%s Item Upd %s", testPrefix, randomSuffix())
		out := zoho(t, "books", "items", "update", itemID, "--org", orgID,
			"--json", toJSON(t, map[string]any{"name": updatedName}))
		m := parseJSON(t, out)
		assertBooksCodeZero(t, m)
	})

	t.Run("mark-inactive", func(t *testing.T) {
		requireID(t, itemID, "create must have succeeded")
		out := zoho(t, "books", "items", "mark-inactive", itemID, "--org", orgID)
		m := parseJSON(t, out)
		assertBooksCodeZero(t, m)
	})

	t.Run("mark-active", func(t *testing.T) {
		requireID(t, itemID, "create must have succeeded")
		out := zoho(t, "books", "items", "mark-active", itemID, "--org", orgID)
		m := parseJSON(t, out)
		assertBooksCodeZero(t, m)
	})

	t.Run("delete", func(t *testing.T) {
		requireID(t, itemID, "create must have succeeded")
		out := zoho(t, "books", "items", "delete", itemID, "--org", orgID)
		m := parseJSON(t, out)
		assertBooksCodeZero(t, m)
	})
}

func TestBooksEstimates(t *testing.T) {
	t.Parallel()
	orgID := requireBooksOrgID(t)
	cleanup := newCleanup(t)

	var contactID string
	var estimateID string

	t.Run("create-contact", func(t *testing.T) {
		name := fmt.Sprintf("%s EstCust %s", testPrefix, randomSuffix())
		out := zoho(t, "books", "contacts", "create", "--org", orgID,
			"--json", toJSON(t, map[string]any{"contact_name": name, "contact_type": "customer"}))
		m := parseJSON(t, out)
		assertBooksCodeZero(t, m)
		contact := m["contact"].(map[string]any)
		contactID = fmt.Sprintf("%v", contact["contact_id"])
		cleanup.trackBooksContact(contactID, orgID)
	})

	t.Run("list", func(t *testing.T) {
		out := zoho(t, "books", "estimates", "list", "--org", orgID)
		m := parseJSON(t, out)
		assertBooksCodeZero(t, m)
	})

	t.Run("create", func(t *testing.T) {
		requireID(t, contactID, "create-contact must have succeeded")
		out := zoho(t, "books", "estimates", "create", "--org", orgID,
			"--json", toJSON(t, map[string]any{
				"customer_id": contactID,
				"line_items": []map[string]any{
					{"description": "Test line item", "rate": 100, "quantity": 1},
				},
			}))
		m := parseJSON(t, out)
		assertBooksCodeZero(t, m)
		estimate, ok := m["estimate"].(map[string]any)
		if !ok {
			t.Fatalf("expected estimate object:\n%s", truncate(out, 500))
		}
		estimateID = fmt.Sprintf("%v", estimate["estimate_id"])
		if estimateID == "" || estimateID == "<nil>" {
			t.Fatalf("expected estimate_id:\n%s", truncate(out, 500))
		}
		cleanup.trackBooksEstimate(estimateID, orgID)
	})

	t.Run("get", func(t *testing.T) {
		requireID(t, estimateID, "create must have succeeded")
		out := zoho(t, "books", "estimates", "get", estimateID, "--org", orgID)
		m := parseJSON(t, out)
		assertBooksCodeZero(t, m)
	})

	t.Run("mark-sent", func(t *testing.T) {
		requireID(t, estimateID, "create must have succeeded")
		out := zoho(t, "books", "estimates", "mark-sent", estimateID, "--org", orgID)
		m := parseJSON(t, out)
		assertBooksCodeZero(t, m)
	})

	t.Run("list-comments", func(t *testing.T) {
		requireID(t, estimateID, "create must have succeeded")
		out := zoho(t, "books", "estimates", "list-comments", estimateID, "--org", orgID)
		m := parseJSON(t, out)
		assertBooksCodeZero(t, m)
	})

	t.Run("delete", func(t *testing.T) {
		requireID(t, estimateID, "create must have succeeded")
		out := zoho(t, "books", "estimates", "delete", estimateID, "--org", orgID)
		m := parseJSON(t, out)
		assertBooksCodeZero(t, m)
	})
}

func TestBooksInvoices(t *testing.T) {
	t.Parallel()
	orgID := requireBooksOrgID(t)
	cleanup := newCleanup(t)

	var contactID string
	var invoiceID string

	t.Run("create-contact", func(t *testing.T) {
		name := fmt.Sprintf("%s InvCust %s", testPrefix, randomSuffix())
		out := zoho(t, "books", "contacts", "create", "--org", orgID,
			"--json", toJSON(t, map[string]any{"contact_name": name, "contact_type": "customer"}))
		m := parseJSON(t, out)
		assertBooksCodeZero(t, m)
		contact := m["contact"].(map[string]any)
		contactID = fmt.Sprintf("%v", contact["contact_id"])
		cleanup.trackBooksContact(contactID, orgID)
	})

	t.Run("list", func(t *testing.T) {
		out := zoho(t, "books", "invoices", "list", "--org", orgID)
		m := parseJSON(t, out)
		assertBooksCodeZero(t, m)
	})

	t.Run("create", func(t *testing.T) {
		requireID(t, contactID, "create-contact must have succeeded")
		out := zoho(t, "books", "invoices", "create", "--org", orgID,
			"--json", toJSON(t, map[string]any{
				"customer_id": contactID,
				"line_items": []map[string]any{
					{"description": "Test service", "rate": 200, "quantity": 1},
				},
			}))
		m := parseJSON(t, out)
		assertBooksCodeZero(t, m)
		invoice, ok := m["invoice"].(map[string]any)
		if !ok {
			t.Fatalf("expected invoice object:\n%s", truncate(out, 500))
		}
		invoiceID = fmt.Sprintf("%v", invoice["invoice_id"])
		if invoiceID == "" || invoiceID == "<nil>" {
			t.Fatalf("expected invoice_id:\n%s", truncate(out, 500))
		}
		cleanup.trackBooksInvoice(invoiceID, orgID)
	})

	t.Run("get", func(t *testing.T) {
		requireID(t, invoiceID, "create must have succeeded")
		out := zoho(t, "books", "invoices", "get", invoiceID, "--org", orgID)
		m := parseJSON(t, out)
		assertBooksCodeZero(t, m)
	})

	t.Run("mark-sent", func(t *testing.T) {
		requireID(t, invoiceID, "create must have succeeded")
		out := zoho(t, "books", "invoices", "mark-sent", invoiceID, "--org", orgID)
		m := parseJSON(t, out)
		assertBooksCodeZero(t, m)
	})

	t.Run("list-payments", func(t *testing.T) {
		requireID(t, invoiceID, "create must have succeeded")
		out := zoho(t, "books", "invoices", "list-payments", invoiceID, "--org", orgID)
		m := parseJSON(t, out)
		assertBooksCodeZero(t, m)
	})

	t.Run("list-comments", func(t *testing.T) {
		requireID(t, invoiceID, "create must have succeeded")
		out := zoho(t, "books", "invoices", "list-comments", invoiceID, "--org", orgID)
		m := parseJSON(t, out)
		assertBooksCodeZero(t, m)
	})

	t.Run("delete", func(t *testing.T) {
		requireID(t, invoiceID, "create must have succeeded")
		out := zoho(t, "books", "invoices", "delete", invoiceID, "--org", orgID)
		m := parseJSON(t, out)
		assertBooksCodeZero(t, m)
	})
}

func TestBooksExpenses(t *testing.T) {
	t.Parallel()
	orgID := requireBooksOrgID(t)
	cleanup := newCleanup(t)

	var accountID string
	var expenseID string

	t.Run("discover-account", func(t *testing.T) {
		out := zoho(t, "books", "chart-of-accounts", "list", "--org", orgID)
		m := parseJSON(t, out)
		assertBooksCodeZero(t, m)
		accounts, ok := m["chartofaccounts"].([]any)
		if !ok || len(accounts) == 0 {
			t.Skip("no chart of accounts found")
		}
		for _, a := range accounts {
			acct := a.(map[string]any)
			aType := fmt.Sprintf("%v", acct["account_type"])
			if aType == "expense" || aType == "cost_of_goods_sold" {
				accountID = fmt.Sprintf("%v", acct["account_id"])
				break
			}
		}
		if accountID == "" || accountID == "<nil>" {
			accountID = fmt.Sprintf("%v", accounts[0].(map[string]any)["account_id"])
		}
		if accountID == "" || accountID == "<nil>" {
			t.Skip("no usable account_id found in chart of accounts")
		}
	})

	t.Run("list", func(t *testing.T) {
		out := zoho(t, "books", "expenses", "list", "--org", orgID)
		m := parseJSON(t, out)
		assertBooksCodeZero(t, m)
	})

	t.Run("create", func(t *testing.T) {
		requireID(t, accountID, "discover-account must have succeeded")
		out := zoho(t, "books", "expenses", "create", "--org", orgID,
			"--json", toJSON(t, map[string]any{
				"account_id": accountID,
				"date":       "2026-03-07",
				"amount":     10.00,
			}))
		m := parseJSON(t, out)
		assertBooksCodeZero(t, m)
		expense, ok := m["expense"].(map[string]any)
		if !ok {
			t.Fatalf("expected expense object:\n%s", truncate(out, 500))
		}
		expenseID = fmt.Sprintf("%v", expense["expense_id"])
		if expenseID == "" || expenseID == "<nil>" {
			t.Fatalf("expected expense_id:\n%s", truncate(out, 500))
		}
		cleanup.trackBooksExpense(expenseID, orgID)
	})

	t.Run("get", func(t *testing.T) {
		requireID(t, expenseID, "create must have succeeded")
		out := zoho(t, "books", "expenses", "get", expenseID, "--org", orgID)
		m := parseJSON(t, out)
		assertBooksCodeZero(t, m)
	})

	t.Run("list-history", func(t *testing.T) {
		requireID(t, expenseID, "create must have succeeded")
		out := zoho(t, "books", "expenses", "list-history", expenseID, "--org", orgID)
		m := parseJSON(t, out)
		assertBooksCodeZero(t, m)
	})

	t.Run("delete", func(t *testing.T) {
		requireID(t, expenseID, "create must have succeeded")
		out := zoho(t, "books", "expenses", "delete", expenseID, "--org", orgID)
		m := parseJSON(t, out)
		assertBooksCodeZero(t, m)
	})
}

func TestBooksCurrencies(t *testing.T) {
	t.Parallel()
	orgID := requireBooksOrgID(t)

	t.Run("list", func(t *testing.T) {
		out := zoho(t, "books", "currencies", "list", "--org", orgID)
		m := parseJSON(t, out)
		assertBooksCodeZero(t, m)
		arr, ok := m["currencies"].([]any)
		if !ok {
			t.Fatalf("expected currencies array:\n%s", truncate(out, 500))
		}
		if len(arr) == 0 {
			t.Fatal("expected at least one currency")
		}
	})
}

func TestBooksTaxes(t *testing.T) {
	t.Parallel()
	orgID := requireBooksOrgID(t)

	t.Run("list", func(t *testing.T) {
		out := zoho(t, "books", "taxes", "list", "--org", orgID)
		m := parseJSON(t, out)
		assertBooksCodeZero(t, m)
	})
}

func TestBooksUsers(t *testing.T) {
	t.Parallel()
	orgID := requireBooksOrgID(t)

	t.Run("list", func(t *testing.T) {
		out := zoho(t, "books", "users", "list", "--org", orgID)
		m := parseJSON(t, out)
		assertBooksCodeZero(t, m)
		arr, ok := m["users"].([]any)
		if !ok {
			t.Fatalf("expected users array:\n%s", truncate(out, 500))
		}
		if len(arr) == 0 {
			t.Fatal("expected at least one user")
		}
	})

	t.Run("get-current", func(t *testing.T) {
		out := zoho(t, "books", "users", "get-current", "--org", orgID)
		m := parseJSON(t, out)
		assertBooksCodeZero(t, m)
	})
}

func TestBooksChartOfAccounts(t *testing.T) {
	t.Parallel()
	orgID := requireBooksOrgID(t)

	t.Run("list", func(t *testing.T) {
		out := zoho(t, "books", "chart-of-accounts", "list", "--org", orgID)
		m := parseJSON(t, out)
		assertBooksCodeZero(t, m)
		arr, ok := m["chartofaccounts"].([]any)
		if !ok {
			t.Fatalf("expected chartofaccounts array:\n%s", truncate(out, 500))
		}
		if len(arr) == 0 {
			t.Fatal("expected at least one chart of account")
		}
	})
}

func TestBooksBankAccounts(t *testing.T) {
	t.Parallel()
	orgID := requireBooksOrgID(t)

	t.Run("list", func(t *testing.T) {
		out := zoho(t, "books", "bank-accounts", "list", "--org", orgID)
		m := parseJSON(t, out)
		assertBooksCodeZero(t, m)
	})
}

func TestBooksErrors(t *testing.T) {
	t.Parallel()

	t.Run("missing-org", func(t *testing.T) {
		r := runZohoWithEnv(t, map[string]string{"ZOHO_BOOKS_ORG_ID": ""}, "books", "contacts", "list")
		if r.ExitCode == 0 {
			t.Error("expected non-zero exit code when --org missing")
		}
		if !strings.Contains(r.Stderr, "ZOHO_BOOKS_ORG_ID") {
			t.Errorf("expected error mentioning ZOHO_BOOKS_ORG_ID, got: %s", r.Stderr)
		}
	})

	t.Run("invalid-org-id", func(t *testing.T) {
		r := runZoho(t, "books", "contacts", "list", "--org", "invalid-org-id-12345")
		if r.ExitCode == 0 {
			t.Log("warning: invalid org ID did not cause error (API may be lenient)")
		}
	})

	t.Run("missing-json-flag", func(t *testing.T) {
		orgID := os.Getenv("ZOHO_BOOKS_ORG_ID")
		if orgID == "" {
			t.Skip("ZOHO_BOOKS_ORG_ID not set")
		}
		r := runZoho(t, "books", "contacts", "create", "--org", orgID)
		if r.ExitCode == 0 {
			t.Error("expected non-zero exit code when --json missing")
		}
	})
}

func TestBooksEmergencyCleanup(t *testing.T) {
	t.Parallel()
	if os.Getenv("ZOHO_EMERGENCY_CLEANUP") != "1" {
		t.Skip("set ZOHO_EMERGENCY_CLEANUP=1 to run")
	}
	orgID := requireBooksOrgID(t)

	out, err := zohoMayFail(t, "books", "contacts", "list", "--org", orgID)
	if err == nil {
		m := parseJSON(t, out)
		if contacts, ok := m["contacts"].([]any); ok {
			for _, item := range contacts {
				c, _ := item.(map[string]any)
				name := fmt.Sprintf("%v", c["contact_name"])
				if strings.HasPrefix(name, testPrefix) {
					id := fmt.Sprintf("%v", c["contact_id"])
					zohoIgnoreError(t, "books", "contacts", "delete", id, "--org", orgID)
				}
			}
		}
	}

	out, err = zohoMayFail(t, "books", "items", "list", "--org", orgID)
	if err == nil {
		m := parseJSON(t, out)
		if items, ok := m["items"].([]any); ok {
			for _, item := range items {
				it, _ := item.(map[string]any)
				name := fmt.Sprintf("%v", it["name"])
				if strings.HasPrefix(name, testPrefix) {
					id := fmt.Sprintf("%v", it["item_id"])
					zohoIgnoreError(t, "books", "items", "delete", id, "--org", orgID)
				}
			}
		}
	}

	out, err = zohoMayFail(t, "books", "estimates", "list", "--org", orgID)
	if err == nil {
		m := parseJSON(t, out)
		if estimates, ok := m["estimates"].([]any); ok {
			for _, item := range estimates {
				est, _ := item.(map[string]any)
				ref := fmt.Sprintf("%v", est["reference_number"])
				if strings.HasPrefix(ref, testPrefix) {
					id := fmt.Sprintf("%v", est["estimate_id"])
					zohoIgnoreError(t, "books", "estimates", "delete", id, "--org", orgID)
				}
			}
		}
	}

	out, err = zohoMayFail(t, "books", "invoices", "list", "--org", orgID)
	if err == nil {
		m := parseJSON(t, out)
		if invoices, ok := m["invoices"].([]any); ok {
			for _, item := range invoices {
				inv, _ := item.(map[string]any)
				ref := fmt.Sprintf("%v", inv["reference_number"])
				if strings.HasPrefix(ref, testPrefix) {
					id := fmt.Sprintf("%v", inv["invoice_id"])
					zohoIgnoreError(t, "books", "invoices", "delete", id, "--org", orgID)
				}
			}
		}
	}
}

