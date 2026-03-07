//go:build integration

package internal_test

import (
	"fmt"
	"os"
	"strings"
	"testing"
)

func requireBillingOrgID(t *testing.T) string {
	t.Helper()
	id := os.Getenv("ZOHO_BOOKS_ORG_ID")
	if id != "" {
		return id
	}
	out, err := zohoMayFail(t, "billing", "organizations", "list")
	if err != nil {
		t.Skipf("skipping: cannot discover billing org ID: %v", err)
	}
	m := parseJSON(t, out)
	orgs, ok := m["organizations"].([]any)
	if !ok || len(orgs) == 0 {
		t.Skip("skipping: no billing organizations found")
	}
	org := orgs[0].(map[string]any)
	orgID := fmt.Sprintf("%v", org["organization_id"])
	if orgID == "" || orgID == "<nil>" {
		t.Skip("skipping: billing organization_id is empty")
	}
	return orgID
}

func assertBillingCodeZero(t *testing.T, m map[string]any) {
	t.Helper()
	if code, ok := m["code"].(float64); ok && code != 0 {
		msg := fmt.Sprintf("%v", m["message"])
		t.Fatalf("billing API error: code=%.0f message=%s", code, msg)
	}
}

func (c *testCleanup) trackBillingProduct(id, orgID string) {
	c.add("delete billing product "+id, func() {
		zohoIgnoreError(c.t, "billing", "products", "delete", id, "--org", orgID)
	})
}

func (c *testCleanup) trackBillingPlan(planCode, orgID string) {
	c.add("delete billing plan "+planCode, func() {
		zohoIgnoreError(c.t, "billing", "plans", "delete", planCode, "--org", orgID)
	})
}

func (c *testCleanup) trackBillingAddon(addonCode, orgID string) {
	c.add("delete billing addon "+addonCode, func() {
		zohoIgnoreError(c.t, "billing", "addons", "delete", addonCode, "--org", orgID)
	})
}

func (c *testCleanup) trackBillingCustomer(id, orgID string) {
	c.add("delete billing customer "+id, func() {
		zohoIgnoreError(c.t, "billing", "customers", "delete", id, "--org", orgID)
	})
}

func TestBillingOrganizations(t *testing.T) {
	t.Parallel()
	orgID := requireBillingOrgID(t)

	t.Run("list", func(t *testing.T) {
		out := zoho(t, "billing", "organizations", "list")
		m := parseJSON(t, out)
		assertBillingCodeZero(t, m)
		arr, ok := m["organizations"].([]any)
		if !ok {
			t.Fatalf("expected organizations array in response:\n%s", truncate(out, 500))
		}
		if len(arr) == 0 {
			t.Fatal("expected at least one organization")
		}
	})

	t.Run("get", func(t *testing.T) {
		out := zoho(t, "billing", "organizations", "get", orgID)
		m := parseJSON(t, out)
		assertBillingCodeZero(t, m)
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

func TestBillingProducts(t *testing.T) {
	t.Parallel()
	orgID := requireBillingOrgID(t)
	cleanup := newCleanup(t)

	var productID string

	t.Run("list", func(t *testing.T) {
		out := zoho(t, "billing", "products", "list", "--org", orgID)
		m := parseJSON(t, out)
		assertBillingCodeZero(t, m)
	})

	t.Run("create", func(t *testing.T) {
		name := fmt.Sprintf("%s Product %s", testPrefix, randomSuffix())
		out := zoho(t, "billing", "products", "create", "--org", orgID,
			"--json", toJSON(t, map[string]any{"name": name}))
		m := parseJSON(t, out)
		assertBillingCodeZero(t, m)
		product, ok := m["product"].(map[string]any)
		if !ok {
			t.Fatalf("expected product object:\n%s", truncate(out, 500))
		}
		productID = fmt.Sprintf("%v", product["product_id"])
		if productID == "" || productID == "<nil>" {
			t.Fatalf("expected product_id:\n%s", truncate(out, 500))
		}
		cleanup.trackBillingProduct(productID, orgID)
	})

	t.Run("get", func(t *testing.T) {
		requireID(t, productID, "create must have succeeded")
		out := zoho(t, "billing", "products", "get", productID, "--org", orgID)
		m := parseJSON(t, out)
		assertBillingCodeZero(t, m)
		product := m["product"].(map[string]any)
		assertEqual(t, fmt.Sprintf("%v", product["product_id"]), productID)
	})

	t.Run("update", func(t *testing.T) {
		requireID(t, productID, "create must have succeeded")
		updatedName := fmt.Sprintf("%s Product Upd %s", testPrefix, randomSuffix())
		out := zoho(t, "billing", "products", "update", productID, "--org", orgID,
			"--json", toJSON(t, map[string]any{"name": updatedName}))
		m := parseJSON(t, out)
		assertBillingCodeZero(t, m)
	})

	t.Run("mark-inactive", func(t *testing.T) {
		requireID(t, productID, "create must have succeeded")
		out := zoho(t, "billing", "products", "mark-inactive", productID, "--org", orgID)
		m := parseJSON(t, out)
		assertBillingCodeZero(t, m)
	})

	t.Run("mark-active", func(t *testing.T) {
		requireID(t, productID, "create must have succeeded")
		out := zoho(t, "billing", "products", "mark-active", productID, "--org", orgID)
		m := parseJSON(t, out)
		assertBillingCodeZero(t, m)
	})

	t.Run("delete", func(t *testing.T) {
		requireID(t, productID, "create must have succeeded")
		out := zoho(t, "billing", "products", "delete", productID, "--org", orgID)
		m := parseJSON(t, out)
		assertBillingCodeZero(t, m)
	})
}

func TestBillingPlans(t *testing.T) {
	t.Parallel()
	orgID := requireBillingOrgID(t)
	cleanup := newCleanup(t)

	var productID string
	var planCode string

	t.Run("create-product", func(t *testing.T) {
		name := fmt.Sprintf("%s PlanProd %s", testPrefix, randomSuffix())
		out := zoho(t, "billing", "products", "create", "--org", orgID,
			"--json", toJSON(t, map[string]any{"name": name}))
		m := parseJSON(t, out)
		assertBillingCodeZero(t, m)
		product := m["product"].(map[string]any)
		productID = fmt.Sprintf("%v", product["product_id"])
		cleanup.trackBillingProduct(productID, orgID)
	})

	t.Run("list", func(t *testing.T) {
		out := zoho(t, "billing", "plans", "list", "--org", orgID)
		m := parseJSON(t, out)
		assertBillingCodeZero(t, m)
	})

	t.Run("create", func(t *testing.T) {
		requireID(t, productID, "create-product must have succeeded")
		planCode = fmt.Sprintf("zohotest-plan-%s", randomSuffix())
		name := fmt.Sprintf("%s Plan %s", testPrefix, randomSuffix())
		out := zoho(t, "billing", "plans", "create", "--org", orgID,
			"--json", toJSON(t, map[string]any{
				"plan_code":       planCode,
				"name":            name,
				"product_id":      productID,
				"recurring_price": 10.00,
				"interval":        1,
				"interval_unit":   "months",
			}))
		m := parseJSON(t, out)
		assertBillingCodeZero(t, m)
		plan, ok := m["plan"].(map[string]any)
		if !ok {
			t.Fatalf("expected plan object:\n%s", truncate(out, 500))
		}
		planCode = fmt.Sprintf("%v", plan["plan_code"])
		if planCode == "" || planCode == "<nil>" {
			t.Fatalf("expected plan_code:\n%s", truncate(out, 500))
		}
		cleanup.trackBillingPlan(planCode, orgID)
	})

	t.Run("get", func(t *testing.T) {
		requireID(t, planCode, "create must have succeeded")
		out := zoho(t, "billing", "plans", "get", planCode, "--org", orgID)
		m := parseJSON(t, out)
		assertBillingCodeZero(t, m)
		plan := m["plan"].(map[string]any)
		assertEqual(t, fmt.Sprintf("%v", plan["plan_code"]), planCode)
	})

	t.Run("update", func(t *testing.T) {
		requireID(t, planCode, "create must have succeeded")
		updatedName := fmt.Sprintf("%s Plan Upd %s", testPrefix, randomSuffix())
		out := zoho(t, "billing", "plans", "update", planCode, "--org", orgID,
			"--json", toJSON(t, map[string]any{"name": updatedName}))
		m := parseJSON(t, out)
		assertBillingCodeZero(t, m)
	})

	t.Run("mark-inactive", func(t *testing.T) {
		requireID(t, planCode, "create must have succeeded")
		out := zoho(t, "billing", "plans", "mark-inactive", planCode, "--org", orgID)
		m := parseJSON(t, out)
		assertBillingCodeZero(t, m)
	})

	t.Run("mark-active", func(t *testing.T) {
		requireID(t, planCode, "create must have succeeded")
		out := zoho(t, "billing", "plans", "mark-active", planCode, "--org", orgID)
		m := parseJSON(t, out)
		assertBillingCodeZero(t, m)
	})

	t.Run("delete", func(t *testing.T) {
		requireID(t, planCode, "create must have succeeded")
		out := zoho(t, "billing", "plans", "delete", planCode, "--org", orgID)
		m := parseJSON(t, out)
		assertBillingCodeZero(t, m)
	})
}

func TestBillingAddons(t *testing.T) {
	t.Parallel()
	orgID := requireBillingOrgID(t)
	cleanup := newCleanup(t)

	var productID string
	var addonCode string

	t.Run("create-product", func(t *testing.T) {
		name := fmt.Sprintf("%s AddonProd %s", testPrefix, randomSuffix())
		out := zoho(t, "billing", "products", "create", "--org", orgID,
			"--json", toJSON(t, map[string]any{"name": name}))
		m := parseJSON(t, out)
		assertBillingCodeZero(t, m)
		product := m["product"].(map[string]any)
		productID = fmt.Sprintf("%v", product["product_id"])
		cleanup.trackBillingProduct(productID, orgID)
	})

	t.Run("list", func(t *testing.T) {
		out := zoho(t, "billing", "addons", "list", "--org", orgID)
		m := parseJSON(t, out)
		assertBillingCodeZero(t, m)
	})

	t.Run("create", func(t *testing.T) {
		requireID(t, productID, "create-product must have succeeded")
		addonCode = fmt.Sprintf("zohotest-addon-%s", randomSuffix())
		name := fmt.Sprintf("%s Addon %s", testPrefix, randomSuffix())
		out := zoho(t, "billing", "addons", "create", "--org", orgID,
			"--json", toJSON(t, map[string]any{
				"addon_code":              addonCode,
				"name":                    name,
				"type":                    "recurring",
				"interval_unit":           "months",
				"pricing_scheme":          "unit",
				"unit_name":               "unit",
				"price_brackets":          []map[string]any{{"start_quantity": 1, "end_quantity": 100, "price": 5}},
				"applicable_to_all_plans": true,
				"product_id":              productID,
			}))
		m := parseJSON(t, out)
		assertBillingCodeZero(t, m)
		addon, ok := m["addon"].(map[string]any)
		if !ok {
			t.Fatalf("expected addon object:\n%s", truncate(out, 500))
		}
		addonCode = fmt.Sprintf("%v", addon["addon_code"])
		if addonCode == "" || addonCode == "<nil>" {
			t.Fatalf("expected addon_code:\n%s", truncate(out, 500))
		}
		cleanup.trackBillingAddon(addonCode, orgID)
	})

	t.Run("get", func(t *testing.T) {
		requireID(t, addonCode, "create must have succeeded")
		out := zoho(t, "billing", "addons", "get", addonCode, "--org", orgID)
		m := parseJSON(t, out)
		assertBillingCodeZero(t, m)
		addon := m["addon"].(map[string]any)
		assertEqual(t, fmt.Sprintf("%v", addon["addon_code"]), addonCode)
	})

	t.Run("update", func(t *testing.T) {
		requireID(t, addonCode, "create must have succeeded")
		updatedName := fmt.Sprintf("%s Addon Upd %s", testPrefix, randomSuffix())
		out := zoho(t, "billing", "addons", "update", addonCode, "--org", orgID,
			"--json", toJSON(t, map[string]any{"name": updatedName}))
		m := parseJSON(t, out)
		assertBillingCodeZero(t, m)
	})

	t.Run("mark-inactive", func(t *testing.T) {
		requireID(t, addonCode, "create must have succeeded")
		out := zoho(t, "billing", "addons", "mark-inactive", addonCode, "--org", orgID)
		m := parseJSON(t, out)
		assertBillingCodeZero(t, m)
	})

	t.Run("mark-active", func(t *testing.T) {
		requireID(t, addonCode, "create must have succeeded")
		out := zoho(t, "billing", "addons", "mark-active", addonCode, "--org", orgID)
		m := parseJSON(t, out)
		assertBillingCodeZero(t, m)
	})

	t.Run("delete", func(t *testing.T) {
		requireID(t, addonCode, "create must have succeeded")
		out := zoho(t, "billing", "addons", "delete", addonCode, "--org", orgID)
		m := parseJSON(t, out)
		assertBillingCodeZero(t, m)
	})
}

func TestBillingCustomers(t *testing.T) {
	t.Parallel()
	orgID := requireBillingOrgID(t)
	cleanup := newCleanup(t)

	var customerID string

	t.Run("list", func(t *testing.T) {
		out := zoho(t, "billing", "customers", "list", "--org", orgID)
		m := parseJSON(t, out)
		assertBillingCodeZero(t, m)
	})

	t.Run("create", func(t *testing.T) {
		suffix := randomSuffix()
		name := fmt.Sprintf("%s Customer %s", testPrefix, suffix)
		out := zoho(t, "billing", "customers", "create", "--org", orgID,
			"--json", toJSON(t, map[string]any{
				"display_name": name,
				"email":        fmt.Sprintf("zohotest_%s@example.com", suffix),
			}))
		m := parseJSON(t, out)
		assertBillingCodeZero(t, m)
		customer, ok := m["customer"].(map[string]any)
		if !ok {
			t.Fatalf("expected customer object:\n%s", truncate(out, 500))
		}
		customerID = fmt.Sprintf("%v", customer["customer_id"])
		if customerID == "" || customerID == "<nil>" {
			t.Fatalf("expected customer_id:\n%s", truncate(out, 500))
		}
		cleanup.trackBillingCustomer(customerID, orgID)
	})

	t.Run("get", func(t *testing.T) {
		requireID(t, customerID, "create must have succeeded")
		out := zoho(t, "billing", "customers", "get", customerID, "--org", orgID)
		m := parseJSON(t, out)
		assertBillingCodeZero(t, m)
		customer := m["customer"].(map[string]any)
		assertEqual(t, fmt.Sprintf("%v", customer["customer_id"]), customerID)
	})

	t.Run("update", func(t *testing.T) {
		requireID(t, customerID, "create must have succeeded")
		updatedName := fmt.Sprintf("%s Customer Upd %s", testPrefix, randomSuffix())
		out := zoho(t, "billing", "customers", "update", customerID, "--org", orgID,
			"--json", toJSON(t, map[string]any{"display_name": updatedName}))
		m := parseJSON(t, out)
		assertBillingCodeZero(t, m)
	})

	t.Run("mark-inactive", func(t *testing.T) {
		requireID(t, customerID, "create must have succeeded")
		out := zoho(t, "billing", "customers", "mark-inactive", customerID, "--org", orgID)
		m := parseJSON(t, out)
		assertBillingCodeZero(t, m)
	})

	t.Run("mark-active", func(t *testing.T) {
		requireID(t, customerID, "create must have succeeded")
		out := zoho(t, "billing", "customers", "mark-active", customerID, "--org", orgID)
		m := parseJSON(t, out)
		assertBillingCodeZero(t, m)
	})

	t.Run("delete", func(t *testing.T) {
		requireID(t, customerID, "create must have succeeded")
		out := zoho(t, "billing", "customers", "delete", customerID, "--org", orgID)
		m := parseJSON(t, out)
		assertBillingCodeZero(t, m)
	})
}

func TestBillingCurrencies(t *testing.T) {
	t.Parallel()
	orgID := requireBillingOrgID(t)

	t.Run("list", func(t *testing.T) {
		out := zoho(t, "billing", "currencies", "list", "--org", orgID)
		m := parseJSON(t, out)
		assertBillingCodeZero(t, m)
		arr, ok := m["currencies"].([]any)
		if !ok {
			t.Fatalf("expected currencies array:\n%s", truncate(out, 500))
		}
		if len(arr) == 0 {
			t.Fatal("expected at least one currency")
		}
	})
}

func TestBillingTaxes(t *testing.T) {
	t.Parallel()
	orgID := requireBillingOrgID(t)

	t.Run("list", func(t *testing.T) {
		out := zoho(t, "billing", "taxes", "list", "--org", orgID)
		m := parseJSON(t, out)
		assertBillingCodeZero(t, m)
	})
}

func TestBillingErrors(t *testing.T) {
	t.Parallel()

	t.Run("missing-org", func(t *testing.T) {
		r := runZohoWithEnv(t, map[string]string{"ZOHO_BOOKS_ORG_ID": ""}, "billing", "customers", "list")
		if r.ExitCode == 0 {
			t.Error("expected non-zero exit code when --org missing")
		}
		if !strings.Contains(r.Stderr, "ZOHO_BOOKS_ORG_ID") {
			t.Errorf("expected error mentioning ZOHO_BOOKS_ORG_ID, got: %s", r.Stderr)
		}
	})

	t.Run("invalid-org-id", func(t *testing.T) {
		r := runZoho(t, "billing", "customers", "list", "--org", "invalid-org-id-12345")
		if r.ExitCode == 0 {
			t.Log("warning: invalid org ID did not cause error (API may be lenient)")
		}
	})

	t.Run("missing-json-flag", func(t *testing.T) {
		orgID := os.Getenv("ZOHO_BOOKS_ORG_ID")
		if orgID == "" {
			t.Skip("ZOHO_BOOKS_ORG_ID not set")
		}
		r := runZoho(t, "billing", "customers", "create", "--org", orgID)
		if r.ExitCode == 0 {
			t.Error("expected non-zero exit code when --json missing")
		}
	})

	t.Run("bad-customer-id", func(t *testing.T) {
		orgID := os.Getenv("ZOHO_BOOKS_ORG_ID")
		if orgID == "" {
			t.Skip("ZOHO_BOOKS_ORG_ID not set")
		}
		r := runZoho(t, "billing", "customers", "get", "000000000000000", "--org", orgID)
		if r.ExitCode == 0 {
			t.Log("warning: bad customer ID did not cause error")
		}
	})

	t.Run("missing-plan-code", func(t *testing.T) {
		r := runZoho(t, "billing", "plans", "get")
		assertExitCode(t, r, 4)
	})

	t.Run("missing-addon-code", func(t *testing.T) {
		r := runZoho(t, "billing", "addons", "get")
		assertExitCode(t, r, 4)
	})

	t.Run("missing-customer-id", func(t *testing.T) {
		r := runZoho(t, "billing", "customers", "get")
		assertExitCode(t, r, 4)
	})

	t.Run("missing-product-id", func(t *testing.T) {
		r := runZoho(t, "billing", "products", "get")
		assertExitCode(t, r, 4)
	})
}

func TestBillingEmergencyCleanup(t *testing.T) {
	t.Parallel()
	if os.Getenv("ZOHO_EMERGENCY_CLEANUP") != "1" {
		t.Skip("set ZOHO_EMERGENCY_CLEANUP=1 to run")
	}
	orgID := requireBillingOrgID(t)

	out, err := zohoMayFail(t, "billing", "customers", "list", "--org", orgID)
	if err == nil {
		m := parseJSON(t, out)
		if customers, ok := m["customers"].([]any); ok {
			for _, item := range customers {
				c, _ := item.(map[string]any)
				name := fmt.Sprintf("%v", c["display_name"])
				if strings.HasPrefix(name, testPrefix) {
					id := fmt.Sprintf("%v", c["customer_id"])
					zohoIgnoreError(t, "billing", "customers", "delete", id, "--org", orgID)
				}
			}
		}
	}

	out, err = zohoMayFail(t, "billing", "addons", "list", "--org", orgID)
	if err == nil {
		m := parseJSON(t, out)
		if addons, ok := m["addons"].([]any); ok {
			for _, item := range addons {
				a, _ := item.(map[string]any)
				name := fmt.Sprintf("%v", a["name"])
				if strings.HasPrefix(name, testPrefix) {
					code := fmt.Sprintf("%v", a["addon_code"])
					zohoIgnoreError(t, "billing", "addons", "delete", code, "--org", orgID)
				}
			}
		}
	}

	out, err = zohoMayFail(t, "billing", "plans", "list", "--org", orgID)
	if err == nil {
		m := parseJSON(t, out)
		if plans, ok := m["plans"].([]any); ok {
			for _, item := range plans {
				p, _ := item.(map[string]any)
				name := fmt.Sprintf("%v", p["name"])
				if strings.HasPrefix(name, testPrefix) {
					code := fmt.Sprintf("%v", p["plan_code"])
					zohoIgnoreError(t, "billing", "plans", "delete", code, "--org", orgID)
				}
			}
		}
	}

	out, err = zohoMayFail(t, "billing", "products", "list", "--org", orgID)
	if err == nil {
		m := parseJSON(t, out)
		if products, ok := m["products"].([]any); ok {
			for _, item := range products {
				p, _ := item.(map[string]any)
				name := fmt.Sprintf("%v", p["name"])
				if strings.HasPrefix(name, testPrefix) {
					id := fmt.Sprintf("%v", p["product_id"])
					zohoIgnoreError(t, "billing", "products", "delete", id, "--org", orgID)
				}
			}
		}
	}
}
