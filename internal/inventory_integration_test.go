//go:build integration

package internal_test

import (
	"fmt"
	"os"
	"strings"
	"testing"
)

func (c *testCleanup) trackInventoryItem(id, orgID string) {
	c.add("delete inventory item "+id, func() {
		zohoIgnoreError(c.t, "inventory", "items", "delete", id, "--org", orgID)
	})
}

func (c *testCleanup) trackInventoryContact(id, orgID string) {
	c.add("delete inventory contact "+id, func() {
		zohoIgnoreError(c.t, "inventory", "contacts", "delete", id, "--org", orgID)
	})
}

func requireInventoryOrgID(t *testing.T) string {
	t.Helper()
	id := os.Getenv("ZOHO_BOOKS_ORG_ID")
	if id != "" {
		return id
	}
	out, err := zohoMayFail(t, "inventory", "organizations", "list")
	if err != nil {
		t.Skipf("skipping: cannot discover inventory org ID: %v", err)
	}
	m := parseJSON(t, out)
	orgs, ok := m["organizations"].([]any)
	if !ok || len(orgs) == 0 {
		t.Skip("skipping: no inventory organizations found")
	}
	org := orgs[0].(map[string]any)
	orgID := fmt.Sprintf("%v", org["organization_id"])
	if orgID == "" || orgID == "<nil>" {
		t.Skip("skipping: inventory organization_id is empty")
	}
	return orgID
}

func assertInventoryCodeZero(t *testing.T, m map[string]any) {
	t.Helper()
	if code, ok := m["code"].(float64); ok && code != 0 {
		msg := fmt.Sprintf("%v", m["message"])
		t.Fatalf("inventory API error: code=%.0f message=%s", code, msg)
	}
}

func TestInventoryOrganizations(t *testing.T) {
	t.Parallel()
	orgID := requireInventoryOrgID(t)

	t.Run("list", func(t *testing.T) {
		out := zoho(t, "inventory", "organizations", "list")
		m := parseJSON(t, out)
		assertInventoryCodeZero(t, m)
		arr, ok := m["organizations"].([]any)
		if !ok {
			t.Fatalf("expected organizations array in response:\n%s", truncate(out, 500))
		}
		if len(arr) == 0 {
			t.Fatal("expected at least one organization")
		}
	})

	t.Run("get", func(t *testing.T) {
		out := zoho(t, "inventory", "organizations", "get", orgID)
		m := parseJSON(t, out)
		assertInventoryCodeZero(t, m)
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

func TestInventoryItems(t *testing.T) {
	t.Parallel()
	orgID := requireInventoryOrgID(t)
	cleanup := newCleanup(t)

	var itemID string

	t.Run("list", func(t *testing.T) {
		out := zoho(t, "inventory", "items", "list", "--org", orgID)
		m := parseJSON(t, out)
		assertInventoryCodeZero(t, m)
	})

	t.Run("create", func(t *testing.T) {
		name := fmt.Sprintf("%s Item %s", testPrefix, randomSuffix())
		out := zoho(t, "inventory", "items", "create", "--org", orgID,
			"--json", toJSON(t, map[string]any{"name": name, "rate": 100}))
		m := parseJSON(t, out)
		assertInventoryCodeZero(t, m)
		item, ok := m["item"].(map[string]any)
		if !ok {
			t.Fatalf("expected item object:\n%s", truncate(out, 500))
		}
		itemID = fmt.Sprintf("%v", item["item_id"])
		if itemID == "" || itemID == "<nil>" {
			t.Fatalf("expected item_id:\n%s", truncate(out, 500))
		}
		cleanup.trackInventoryItem(itemID, orgID)
	})

	t.Run("get", func(t *testing.T) {
		requireID(t, itemID, "create must have succeeded")
		out := zoho(t, "inventory", "items", "get", itemID, "--org", orgID)
		m := parseJSON(t, out)
		assertInventoryCodeZero(t, m)
		item := m["item"].(map[string]any)
		assertEqual(t, fmt.Sprintf("%v", item["item_id"]), itemID)
	})

	t.Run("update", func(t *testing.T) {
		requireID(t, itemID, "create must have succeeded")
		updatedName := fmt.Sprintf("%s Item Upd %s", testPrefix, randomSuffix())
		out := zoho(t, "inventory", "items", "update", itemID, "--org", orgID,
			"--json", toJSON(t, map[string]any{"name": updatedName}))
		m := parseJSON(t, out)
		assertInventoryCodeZero(t, m)
	})

	t.Run("mark-inactive", func(t *testing.T) {
		requireID(t, itemID, "create must have succeeded")
		out := zoho(t, "inventory", "items", "mark-inactive", itemID, "--org", orgID)
		m := parseJSON(t, out)
		assertInventoryCodeZero(t, m)
	})

	t.Run("mark-active", func(t *testing.T) {
		requireID(t, itemID, "create must have succeeded")
		out := zoho(t, "inventory", "items", "mark-active", itemID, "--org", orgID)
		m := parseJSON(t, out)
		assertInventoryCodeZero(t, m)
	})

	t.Run("delete", func(t *testing.T) {
		requireID(t, itemID, "create must have succeeded")
		out := zoho(t, "inventory", "items", "delete", itemID, "--org", orgID)
		m := parseJSON(t, out)
		assertInventoryCodeZero(t, m)
	})
}

func TestInventoryContacts(t *testing.T) {
	t.Parallel()
	orgID := requireInventoryOrgID(t)
	cleanup := newCleanup(t)

	var contactID string

	t.Run("list", func(t *testing.T) {
		out := zoho(t, "inventory", "contacts", "list", "--org", orgID)
		m := parseJSON(t, out)
		assertInventoryCodeZero(t, m)
	})

	t.Run("create", func(t *testing.T) {
		name := fmt.Sprintf("%s Contact %s", testPrefix, randomSuffix())
		out := zoho(t, "inventory", "contacts", "create", "--org", orgID,
			"--json", toJSON(t, map[string]any{"contact_name": name, "contact_type": "customer"}))
		m := parseJSON(t, out)
		assertInventoryCodeZero(t, m)
		contact, ok := m["contact"].(map[string]any)
		if !ok {
			t.Fatalf("expected contact object:\n%s", truncate(out, 500))
		}
		contactID = fmt.Sprintf("%v", contact["contact_id"])
		if contactID == "" || contactID == "<nil>" {
			t.Fatalf("expected contact_id:\n%s", truncate(out, 500))
		}
		cleanup.trackInventoryContact(contactID, orgID)
	})

	t.Run("get", func(t *testing.T) {
		requireID(t, contactID, "create must have succeeded")
		out := zoho(t, "inventory", "contacts", "get", contactID, "--org", orgID)
		m := parseJSON(t, out)
		assertInventoryCodeZero(t, m)
		contact := m["contact"].(map[string]any)
		assertEqual(t, fmt.Sprintf("%v", contact["contact_id"]), contactID)
	})

	t.Run("update", func(t *testing.T) {
		requireID(t, contactID, "create must have succeeded")
		updatedName := fmt.Sprintf("%s Contact Upd %s", testPrefix, randomSuffix())
		out := zoho(t, "inventory", "contacts", "update", contactID, "--org", orgID,
			"--json", toJSON(t, map[string]any{"contact_name": updatedName}))
		m := parseJSON(t, out)
		assertInventoryCodeZero(t, m)
	})

	t.Run("mark-inactive", func(t *testing.T) {
		requireID(t, contactID, "create must have succeeded")
		out := zoho(t, "inventory", "contacts", "mark-inactive", contactID, "--org", orgID)
		m := parseJSON(t, out)
		assertInventoryCodeZero(t, m)
	})

	t.Run("mark-active", func(t *testing.T) {
		requireID(t, contactID, "create must have succeeded")
		out := zoho(t, "inventory", "contacts", "mark-active", contactID, "--org", orgID)
		m := parseJSON(t, out)
		assertInventoryCodeZero(t, m)
	})

	t.Run("delete", func(t *testing.T) {
		requireID(t, contactID, "create must have succeeded")
		out := zoho(t, "inventory", "contacts", "delete", contactID, "--org", orgID)
		m := parseJSON(t, out)
		assertInventoryCodeZero(t, m)
	})
}

func TestInventoryCurrencies(t *testing.T) {
	t.Parallel()
	orgID := requireInventoryOrgID(t)

	t.Run("list", func(t *testing.T) {
		out := zoho(t, "inventory", "currencies", "list", "--org", orgID)
		m := parseJSON(t, out)
		assertInventoryCodeZero(t, m)
		arr, ok := m["currencies"].([]any)
		if !ok {
			t.Fatalf("expected currencies array:\n%s", truncate(out, 500))
		}
		if len(arr) == 0 {
			t.Fatal("expected at least one currency")
		}
	})
}

func TestInventoryTaxes(t *testing.T) {
	t.Parallel()
	orgID := requireInventoryOrgID(t)

	t.Run("list", func(t *testing.T) {
		out := zoho(t, "inventory", "taxes", "list", "--org", orgID)
		m := parseJSON(t, out)
		assertInventoryCodeZero(t, m)
	})
}

func TestInventoryUsers(t *testing.T) {
	t.Parallel()
	orgID := requireInventoryOrgID(t)

	t.Run("list", func(t *testing.T) {
		out := zoho(t, "inventory", "users", "list", "--org", orgID)
		m := parseJSON(t, out)
		assertInventoryCodeZero(t, m)
		arr, ok := m["users"].([]any)
		if !ok {
			t.Fatalf("expected users array:\n%s", truncate(out, 500))
		}
		if len(arr) == 0 {
			t.Fatal("expected at least one user")
		}
	})

	t.Run("get-current", func(t *testing.T) {
		out := zoho(t, "inventory", "users", "get-current", "--org", orgID)
		m := parseJSON(t, out)
		assertInventoryCodeZero(t, m)
	})
}

func TestInventoryErrors(t *testing.T) {
	t.Parallel()

	t.Run("missing-org", func(t *testing.T) {
		r := runZohoWithEnv(t, map[string]string{"ZOHO_BOOKS_ORG_ID": ""}, "inventory", "contacts", "list")
		if r.ExitCode == 0 {
			t.Error("expected non-zero exit code when --org missing")
		}
		if !strings.Contains(r.Stderr, "ZOHO_BOOKS_ORG_ID") {
			t.Errorf("expected error mentioning ZOHO_BOOKS_ORG_ID, got: %s", r.Stderr)
		}
	})

	t.Run("invalid-org-id", func(t *testing.T) {
		r := runZoho(t, "inventory", "contacts", "list", "--org", "invalid-org-id-12345")
		if r.ExitCode == 0 {
			t.Log("warning: invalid org ID did not cause error (API may be lenient)")
		}
	})

	t.Run("missing-json-flag", func(t *testing.T) {
		orgID := os.Getenv("ZOHO_BOOKS_ORG_ID")
		if orgID == "" {
			t.Skip("ZOHO_BOOKS_ORG_ID not set")
		}
		r := runZoho(t, "inventory", "items", "create", "--org", orgID)
		if r.ExitCode == 0 {
			t.Error("expected non-zero exit code when --json missing")
		}
	})

	t.Run("bad-item-id", func(t *testing.T) {
		orgID := os.Getenv("ZOHO_BOOKS_ORG_ID")
		if orgID == "" {
			t.Skip("ZOHO_BOOKS_ORG_ID not set")
		}
		r := runZoho(t, "inventory", "items", "get", "999999999999999", "--org", orgID)
		if r.ExitCode == 0 {
			t.Log("warning: bad item ID did not cause error")
		}
	})

	t.Run("bad-contact-id", func(t *testing.T) {
		orgID := os.Getenv("ZOHO_BOOKS_ORG_ID")
		if orgID == "" {
			t.Skip("ZOHO_BOOKS_ORG_ID not set")
		}
		r := runZoho(t, "inventory", "contacts", "get", "999999999999999", "--org", orgID)
		if r.ExitCode == 0 {
			t.Log("warning: bad contact ID did not cause error")
		}
	})

	t.Run("missing-item-id", func(t *testing.T) {
		orgID := os.Getenv("ZOHO_BOOKS_ORG_ID")
		if orgID == "" {
			t.Skip("ZOHO_BOOKS_ORG_ID not set")
		}
		r := runZoho(t, "inventory", "items", "get", "--org", orgID)
		assertExitCode(t, r, 4)
	})

	t.Run("missing-contact-id", func(t *testing.T) {
		orgID := os.Getenv("ZOHO_BOOKS_ORG_ID")
		if orgID == "" {
			t.Skip("ZOHO_BOOKS_ORG_ID not set")
		}
		r := runZoho(t, "inventory", "contacts", "get", "--org", orgID)
		assertExitCode(t, r, 4)
	})
}

func TestInventoryEmergencyCleanup(t *testing.T) {
	t.Parallel()
	if os.Getenv("ZOHO_EMERGENCY_CLEANUP") != "1" {
		t.Skip("set ZOHO_EMERGENCY_CLEANUP=1 to run")
	}
	orgID := requireInventoryOrgID(t)

	out, err := zohoMayFail(t, "inventory", "items", "list", "--org", orgID)
	if err == nil {
		m := parseJSON(t, out)
		if items, ok := m["items"].([]any); ok {
			for _, item := range items {
				it, _ := item.(map[string]any)
				name := fmt.Sprintf("%v", it["name"])
				if strings.HasPrefix(name, testPrefix) {
					id := fmt.Sprintf("%v", it["item_id"])
					zohoIgnoreError(t, "inventory", "items", "delete", id, "--org", orgID)
				}
			}
		}
	}

	out, err = zohoMayFail(t, "inventory", "contacts", "list", "--org", orgID)
	if err == nil {
		m := parseJSON(t, out)
		if contacts, ok := m["contacts"].([]any); ok {
			for _, item := range contacts {
				c, _ := item.(map[string]any)
				name := fmt.Sprintf("%v", c["contact_name"])
				if strings.HasPrefix(name, testPrefix) {
					id := fmt.Sprintf("%v", c["contact_id"])
					zohoIgnoreError(t, "inventory", "contacts", "delete", id, "--org", orgID)
				}
			}
		}
	}
}
