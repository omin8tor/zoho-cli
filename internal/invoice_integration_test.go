//go:build integration

package internal_test

import (
	"fmt"
	"os"
	"strings"
	"testing"
)

func (c *testCleanup) trackInvoiceContact(id, orgID string) {
	c.add("delete invoice contact "+id, func() {
		zohoIgnoreError(c.t, "invoice", "contacts", "delete", id, "--org", orgID)
	})
}

func (c *testCleanup) trackInvoiceItem(id, orgID string) {
	c.add("delete invoice item "+id, func() {
		zohoIgnoreError(c.t, "invoice", "items", "delete", id, "--org", orgID)
	})
}

func (c *testCleanup) trackInvoiceInvoice(id, orgID string) {
	c.add("delete invoice invoice "+id, func() {
		zohoIgnoreError(c.t, "invoice", "invoices", "delete", id, "--org", orgID)
	})
}

func (c *testCleanup) trackInvoiceEstimate(id, orgID string) {
	c.add("delete invoice estimate "+id, func() {
		zohoIgnoreError(c.t, "invoice", "estimates", "delete", id, "--org", orgID)
	})
}

func requireInvoiceOrgID(t *testing.T) string {
	t.Helper()
	id := os.Getenv("ZOHO_BOOKS_ORG_ID")
	if id != "" {
		return id
	}
	out, err := zohoMayFail(t, "invoice", "organizations", "list")
	if err != nil {
		t.Skipf("skipping: cannot discover invoice org ID: %v", err)
	}
	m := parseJSON(t, out)
	orgs, ok := m["organizations"].([]any)
	if !ok || len(orgs) == 0 {
		t.Skip("skipping: no invoice organizations found")
	}
	org := orgs[0].(map[string]any)
	orgID := fmt.Sprintf("%v", org["organization_id"])
	if orgID == "" || orgID == "<nil>" {
		t.Skip("skipping: invoice organization_id is empty")
	}
	return orgID
}

func assertInvoiceCodeZero(t *testing.T, m map[string]any) {
	t.Helper()
	if code, ok := m["code"].(float64); ok && code != 0 {
		msg := fmt.Sprintf("%v", m["message"])
		t.Fatalf("invoice API error: code=%.0f message=%s", code, msg)
	}
}

func TestInvoiceOrganizations(t *testing.T) {
	t.Parallel()
	orgID := requireInvoiceOrgID(t)

	t.Run("list", func(t *testing.T) {
		out := zoho(t, "invoice", "organizations", "list")
		m := parseJSON(t, out)
		assertInvoiceCodeZero(t, m)
		arr, ok := m["organizations"].([]any)
		if !ok {
			t.Fatalf("expected organizations array in response:\n%s", truncate(out, 500))
		}
		if len(arr) == 0 {
			t.Fatal("expected at least one organization")
		}
	})

	t.Run("get", func(t *testing.T) {
		out := zoho(t, "invoice", "organizations", "get", orgID)
		m := parseJSON(t, out)
		assertInvoiceCodeZero(t, m)
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

func TestInvoiceContacts(t *testing.T) {
	t.Parallel()
	orgID := requireInvoiceOrgID(t)
	cleanup := newCleanup(t)

	var contactID string

	t.Run("list", func(t *testing.T) {
		out := zoho(t, "invoice", "contacts", "list", "--org", orgID)
		m := parseJSON(t, out)
		assertInvoiceCodeZero(t, m)
	})

	t.Run("create", func(t *testing.T) {
		name := fmt.Sprintf("%s Contact %s", testPrefix, randomSuffix())
		out := zoho(t, "invoice", "contacts", "create", "--org", orgID,
			"--contact_name", name, "--contact_type", "customer")
		m := parseJSON(t, out)
		assertInvoiceCodeZero(t, m)
		contact, ok := m["contact"].(map[string]any)
		if !ok {
			t.Fatalf("expected contact object:\n%s", truncate(out, 500))
		}
		contactID = fmt.Sprintf("%v", contact["contact_id"])
		if contactID == "" || contactID == "<nil>" {
			t.Fatalf("expected contact_id:\n%s", truncate(out, 500))
		}
		cleanup.trackInvoiceContact(contactID, orgID)
	})

	t.Run("get", func(t *testing.T) {
		requireID(t, contactID, "create must have succeeded")
		out := zoho(t, "invoice", "contacts", "get", contactID, "--org", orgID)
		m := parseJSON(t, out)
		assertInvoiceCodeZero(t, m)
		contact := m["contact"].(map[string]any)
		assertEqual(t, fmt.Sprintf("%v", contact["contact_id"]), contactID)
	})

	t.Run("update", func(t *testing.T) {
		requireID(t, contactID, "create must have succeeded")
		updatedName := fmt.Sprintf("%s Contact Upd %s", testPrefix, randomSuffix())
		out := zoho(t, "invoice", "contacts", "update", contactID, "--org", orgID,
			"--contact_name", updatedName)
		m := parseJSON(t, out)
		assertInvoiceCodeZero(t, m)
	})

	t.Run("mark-inactive", func(t *testing.T) {
		requireID(t, contactID, "create must have succeeded")
		out := zoho(t, "invoice", "contacts", "mark-inactive", contactID, "--org", orgID)
		m := parseJSON(t, out)
		assertInvoiceCodeZero(t, m)
	})

	t.Run("mark-active", func(t *testing.T) {
		requireID(t, contactID, "create must have succeeded")
		out := zoho(t, "invoice", "contacts", "mark-active", contactID, "--org", orgID)
		m := parseJSON(t, out)
		assertInvoiceCodeZero(t, m)
	})

	t.Run("delete", func(t *testing.T) {
		requireID(t, contactID, "create must have succeeded")
		out := zoho(t, "invoice", "contacts", "delete", contactID, "--org", orgID)
		m := parseJSON(t, out)
		assertInvoiceCodeZero(t, m)
	})
}

func TestInvoiceItems(t *testing.T) {
	t.Parallel()
	orgID := requireInvoiceOrgID(t)
	cleanup := newCleanup(t)

	var itemID string

	t.Run("list", func(t *testing.T) {
		out := zoho(t, "invoice", "items", "list", "--org", orgID)
		m := parseJSON(t, out)
		assertInvoiceCodeZero(t, m)
	})

	t.Run("create", func(t *testing.T) {
		name := fmt.Sprintf("%s Item %s", testPrefix, randomSuffix())
		out := zoho(t, "invoice", "items", "create", "--org", orgID,
			"--name", name, "--rate", "100")
		m := parseJSON(t, out)
		assertInvoiceCodeZero(t, m)
		item, ok := m["item"].(map[string]any)
		if !ok {
			t.Fatalf("expected item object:\n%s", truncate(out, 500))
		}
		itemID = fmt.Sprintf("%v", item["item_id"])
		if itemID == "" || itemID == "<nil>" {
			t.Fatalf("expected item_id:\n%s", truncate(out, 500))
		}
		cleanup.trackInvoiceItem(itemID, orgID)
	})

	t.Run("get", func(t *testing.T) {
		requireID(t, itemID, "create must have succeeded")
		out := zoho(t, "invoice", "items", "get", itemID, "--org", orgID)
		m := parseJSON(t, out)
		assertInvoiceCodeZero(t, m)
		item := m["item"].(map[string]any)
		assertEqual(t, fmt.Sprintf("%v", item["item_id"]), itemID)
	})

	t.Run("update", func(t *testing.T) {
		requireID(t, itemID, "create must have succeeded")
		updatedName := fmt.Sprintf("%s Item Upd %s", testPrefix, randomSuffix())
		out := zoho(t, "invoice", "items", "update", itemID, "--org", orgID,
			"--name", updatedName)
		m := parseJSON(t, out)
		assertInvoiceCodeZero(t, m)
	})

	t.Run("mark-inactive", func(t *testing.T) {
		requireID(t, itemID, "create must have succeeded")
		out := zoho(t, "invoice", "items", "mark-inactive", itemID, "--org", orgID)
		m := parseJSON(t, out)
		assertInvoiceCodeZero(t, m)
	})

	t.Run("mark-active", func(t *testing.T) {
		requireID(t, itemID, "create must have succeeded")
		out := zoho(t, "invoice", "items", "mark-active", itemID, "--org", orgID)
		m := parseJSON(t, out)
		assertInvoiceCodeZero(t, m)
	})

	t.Run("delete", func(t *testing.T) {
		requireID(t, itemID, "create must have succeeded")
		out := zoho(t, "invoice", "items", "delete", itemID, "--org", orgID)
		m := parseJSON(t, out)
		assertInvoiceCodeZero(t, m)
	})
}

func TestInvoiceInvoices(t *testing.T) {
	t.Parallel()
	orgID := requireInvoiceOrgID(t)
	cleanup := newCleanup(t)

	var contactID string
	var itemID string
	var invoiceID string

	t.Run("create-contact", func(t *testing.T) {
		name := fmt.Sprintf("%s InvCust %s", testPrefix, randomSuffix())
		out := zoho(t, "invoice", "contacts", "create", "--org", orgID,
			"--contact_name", name, "--contact_type", "customer")
		m := parseJSON(t, out)
		assertInvoiceCodeZero(t, m)
		contact := m["contact"].(map[string]any)
		contactID = fmt.Sprintf("%v", contact["contact_id"])
		cleanup.trackInvoiceContact(contactID, orgID)
	})

	t.Run("create-item", func(t *testing.T) {
		name := fmt.Sprintf("%s InvItem %s", testPrefix, randomSuffix())
		out := zoho(t, "invoice", "items", "create", "--org", orgID,
			"--name", name, "--rate", "250")
		m := parseJSON(t, out)
		assertInvoiceCodeZero(t, m)
		item := m["item"].(map[string]any)
		itemID = fmt.Sprintf("%v", item["item_id"])
		cleanup.trackInvoiceItem(itemID, orgID)
	})

	t.Run("list", func(t *testing.T) {
		out := zoho(t, "invoice", "invoices", "list", "--org", orgID)
		m := parseJSON(t, out)
		assertInvoiceCodeZero(t, m)
	})

	t.Run("create", func(t *testing.T) {
		requireID(t, contactID, "create-contact must have succeeded")
		requireID(t, itemID, "create-item must have succeeded")
		out := zoho(t, "invoice", "invoices", "create", "--org", orgID,
			"--customer_id", contactID,
			"--json", toJSON(t, map[string]any{
				"line_items": []map[string]any{
					{"item_id": itemID, "rate": 250, "quantity": 1},
				},
			}))
		m := parseJSON(t, out)
		assertInvoiceCodeZero(t, m)
		inv, ok := m["invoice"].(map[string]any)
		if !ok {
			t.Fatalf("expected invoice object:\n%s", truncate(out, 500))
		}
		invoiceID = fmt.Sprintf("%v", inv["invoice_id"])
		if invoiceID == "" || invoiceID == "<nil>" {
			t.Fatalf("expected invoice_id:\n%s", truncate(out, 500))
		}
		cleanup.trackInvoiceInvoice(invoiceID, orgID)
	})

	t.Run("get", func(t *testing.T) {
		requireID(t, invoiceID, "create must have succeeded")
		out := zoho(t, "invoice", "invoices", "get", invoiceID, "--org", orgID)
		m := parseJSON(t, out)
		assertInvoiceCodeZero(t, m)
	})

	t.Run("update", func(t *testing.T) {
		requireID(t, invoiceID, "create must have succeeded")
		out := zoho(t, "invoice", "invoices", "update", invoiceID, "--org", orgID,
			"--json", toJSON(t, map[string]any{
				"line_items": []map[string]any{
					{"item_id": itemID, "rate": 300, "quantity": 2},
				},
			}))
		m := parseJSON(t, out)
		assertInvoiceCodeZero(t, m)
	})

	t.Run("mark-sent", func(t *testing.T) {
		requireID(t, invoiceID, "create must have succeeded")
		out := zoho(t, "invoice", "invoices", "mark-sent", invoiceID, "--org", orgID)
		m := parseJSON(t, out)
		assertInvoiceCodeZero(t, m)
	})

	t.Run("list-payments", func(t *testing.T) {
		requireID(t, invoiceID, "create must have succeeded")
		out := zoho(t, "invoice", "invoices", "list-payments", invoiceID, "--org", orgID)
		m := parseJSON(t, out)
		assertInvoiceCodeZero(t, m)
	})

	t.Run("list-comments", func(t *testing.T) {
		requireID(t, invoiceID, "create must have succeeded")
		out := zoho(t, "invoice", "invoices", "list-comments", invoiceID, "--org", orgID)
		m := parseJSON(t, out)
		assertInvoiceCodeZero(t, m)
	})

	t.Run("delete", func(t *testing.T) {
		requireID(t, invoiceID, "create must have succeeded")
		out := zoho(t, "invoice", "invoices", "delete", invoiceID, "--org", orgID)
		m := parseJSON(t, out)
		assertInvoiceCodeZero(t, m)
	})
}

func TestInvoiceEstimates(t *testing.T) {
	t.Parallel()
	orgID := requireInvoiceOrgID(t)
	cleanup := newCleanup(t)

	var contactID string
	var estimateID string

	t.Run("create-contact", func(t *testing.T) {
		name := fmt.Sprintf("%s EstCust %s", testPrefix, randomSuffix())
		out := zoho(t, "invoice", "contacts", "create", "--org", orgID,
			"--contact_name", name, "--contact_type", "customer")
		m := parseJSON(t, out)
		assertInvoiceCodeZero(t, m)
		contact := m["contact"].(map[string]any)
		contactID = fmt.Sprintf("%v", contact["contact_id"])
		cleanup.trackInvoiceContact(contactID, orgID)
	})

	t.Run("list", func(t *testing.T) {
		out := zoho(t, "invoice", "estimates", "list", "--org", orgID)
		m := parseJSON(t, out)
		assertInvoiceCodeZero(t, m)
	})

	t.Run("create", func(t *testing.T) {
		requireID(t, contactID, "create-contact must have succeeded")
		out := zoho(t, "invoice", "estimates", "create", "--org", orgID,
			"--customer_id", contactID,
			"--json", toJSON(t, map[string]any{
				"line_items": []map[string]any{
					{"description": "Test line item", "rate": 100, "quantity": 1},
				},
			}))
		m := parseJSON(t, out)
		assertInvoiceCodeZero(t, m)
		estimate, ok := m["estimate"].(map[string]any)
		if !ok {
			t.Fatalf("expected estimate object:\n%s", truncate(out, 500))
		}
		estimateID = fmt.Sprintf("%v", estimate["estimate_id"])
		if estimateID == "" || estimateID == "<nil>" {
			t.Fatalf("expected estimate_id:\n%s", truncate(out, 500))
		}
		cleanup.trackInvoiceEstimate(estimateID, orgID)
	})

	t.Run("get", func(t *testing.T) {
		requireID(t, estimateID, "create must have succeeded")
		out := zoho(t, "invoice", "estimates", "get", estimateID, "--org", orgID)
		m := parseJSON(t, out)
		assertInvoiceCodeZero(t, m)
	})

	t.Run("list-comments", func(t *testing.T) {
		requireID(t, estimateID, "create must have succeeded")
		out := zoho(t, "invoice", "estimates", "list-comments", estimateID, "--org", orgID)
		m := parseJSON(t, out)
		assertInvoiceCodeZero(t, m)
	})

	t.Run("delete", func(t *testing.T) {
		requireID(t, estimateID, "create must have succeeded")
		out := zoho(t, "invoice", "estimates", "delete", estimateID, "--org", orgID)
		m := parseJSON(t, out)
		assertInvoiceCodeZero(t, m)
	})
}

func TestInvoiceCurrencies(t *testing.T) {
	t.Parallel()
	orgID := requireInvoiceOrgID(t)

	t.Run("list", func(t *testing.T) {
		out := zoho(t, "invoice", "currencies", "list", "--org", orgID)
		m := parseJSON(t, out)
		assertInvoiceCodeZero(t, m)
		arr, ok := m["currencies"].([]any)
		if !ok {
			t.Fatalf("expected currencies array:\n%s", truncate(out, 500))
		}
		if len(arr) == 0 {
			t.Fatal("expected at least one currency")
		}
	})
}

func TestInvoiceTaxes(t *testing.T) {
	t.Parallel()
	orgID := requireInvoiceOrgID(t)

	t.Run("list", func(t *testing.T) {
		out := zoho(t, "invoice", "taxes", "list", "--org", orgID)
		m := parseJSON(t, out)
		assertInvoiceCodeZero(t, m)
	})
}

func TestInvoiceUsers(t *testing.T) {
	t.Parallel()
	orgID := requireInvoiceOrgID(t)

	t.Run("list", func(t *testing.T) {
		out := zoho(t, "invoice", "users", "list", "--org", orgID)
		m := parseJSON(t, out)
		assertInvoiceCodeZero(t, m)
		arr, ok := m["users"].([]any)
		if !ok {
			t.Fatalf("expected users array:\n%s", truncate(out, 500))
		}
		if len(arr) == 0 {
			t.Fatal("expected at least one user")
		}
	})

	t.Run("get-current", func(t *testing.T) {
		out := zoho(t, "invoice", "users", "get-current", "--org", orgID)
		m := parseJSON(t, out)
		assertInvoiceCodeZero(t, m)
	})
}

func TestInvoiceErrors(t *testing.T) {
	t.Parallel()

	t.Run("missing-org", func(t *testing.T) {
		r := runZohoWithEnv(t, map[string]string{"ZOHO_BOOKS_ORG_ID": ""}, "invoice", "contacts", "list")
		if r.ExitCode == 0 {
			t.Error("expected non-zero exit code when --org missing")
		}
		if !strings.Contains(r.Stderr, "ZOHO_BOOKS_ORG_ID") {
			t.Errorf("expected error mentioning ZOHO_BOOKS_ORG_ID, got: %s", r.Stderr)
		}
	})

	t.Run("invalid-org-id", func(t *testing.T) {
		r := runZoho(t, "invoice", "contacts", "list", "--org", "invalid-org-id-12345")
		if r.ExitCode == 0 {
			t.Log("warning: invalid org ID did not cause error (API may be lenient)")
		}
	})

	t.Run("missing-required-flag", func(t *testing.T) {
		orgID := os.Getenv("ZOHO_BOOKS_ORG_ID")
		if orgID == "" {
			t.Skip("ZOHO_BOOKS_ORG_ID not set")
		}
		r := runZoho(t, "invoice", "contacts", "create", "--org", orgID)
		if r.ExitCode == 0 {
			t.Error("expected non-zero exit code when --contact_name missing")
		}
	})

	t.Run("bad-contact-id", func(t *testing.T) {
		orgID := os.Getenv("ZOHO_BOOKS_ORG_ID")
		if orgID == "" {
			t.Skip("ZOHO_BOOKS_ORG_ID not set")
		}
		r := runZoho(t, "invoice", "contacts", "get", "999999999999999", "--org", orgID)
		if r.ExitCode == 0 {
			t.Log("warning: bad contact ID did not cause error (API may be lenient)")
		}
	})

	t.Run("bad-invoice-id", func(t *testing.T) {
		orgID := os.Getenv("ZOHO_BOOKS_ORG_ID")
		if orgID == "" {
			t.Skip("ZOHO_BOOKS_ORG_ID not set")
		}
		r := runZoho(t, "invoice", "invoices", "get", "999999999999999", "--org", orgID)
		if r.ExitCode == 0 {
			t.Log("warning: bad invoice ID did not cause error (API may be lenient)")
		}
	})
}

func TestInvoiceEmergencyCleanup(t *testing.T) {
	t.Parallel()
	if os.Getenv("ZOHO_EMERGENCY_CLEANUP") != "1" {
		t.Skip("set ZOHO_EMERGENCY_CLEANUP=1 to run")
	}
	orgID := requireInvoiceOrgID(t)

	out, err := zohoMayFail(t, "invoice", "invoices", "list", "--org", orgID)
	if err == nil {
		m := parseJSON(t, out)
		if invoices, ok := m["invoices"].([]any); ok {
			for _, item := range invoices {
				inv, _ := item.(map[string]any)
				ref := fmt.Sprintf("%v", inv["reference_number"])
				if strings.HasPrefix(ref, testPrefix) {
					id := fmt.Sprintf("%v", inv["invoice_id"])
					zohoIgnoreError(t, "invoice", "invoices", "delete", id, "--org", orgID)
				}
			}
		}
	}

	out, err = zohoMayFail(t, "invoice", "estimates", "list", "--org", orgID)
	if err == nil {
		m := parseJSON(t, out)
		if estimates, ok := m["estimates"].([]any); ok {
			for _, item := range estimates {
				est, _ := item.(map[string]any)
				ref := fmt.Sprintf("%v", est["reference_number"])
				if strings.HasPrefix(ref, testPrefix) {
					id := fmt.Sprintf("%v", est["estimate_id"])
					zohoIgnoreError(t, "invoice", "estimates", "delete", id, "--org", orgID)
				}
			}
		}
	}

	out, err = zohoMayFail(t, "invoice", "items", "list", "--org", orgID)
	if err == nil {
		m := parseJSON(t, out)
		if items, ok := m["items"].([]any); ok {
			for _, item := range items {
				it, _ := item.(map[string]any)
				name := fmt.Sprintf("%v", it["name"])
				if strings.HasPrefix(name, testPrefix) {
					id := fmt.Sprintf("%v", it["item_id"])
					zohoIgnoreError(t, "invoice", "items", "delete", id, "--org", orgID)
				}
			}
		}
	}

	out, err = zohoMayFail(t, "invoice", "contacts", "list", "--org", orgID)
	if err == nil {
		m := parseJSON(t, out)
		if contacts, ok := m["contacts"].([]any); ok {
			for _, item := range contacts {
				c, _ := item.(map[string]any)
				name := fmt.Sprintf("%v", c["contact_name"])
				if strings.HasPrefix(name, testPrefix) {
					id := fmt.Sprintf("%v", c["contact_id"])
					zohoIgnoreError(t, "invoice", "contacts", "delete", id, "--org", orgID)
				}
			}
		}
	}
}
