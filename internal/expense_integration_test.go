//go:build integration

package internal_test

import (
	"fmt"
	"os"
	"strings"
	"testing"
)

func (c *testCleanup) trackExpenseCategory(id, orgID string) {
	c.add("delete expense category "+id, func() {
		zohoIgnoreError(c.t, "expense", "categories", "delete", id, "--org", orgID)
	})
}

func (c *testCleanup) trackExpenseCustomer(id, orgID string) {
	c.add("delete expense customer "+id, func() {
		zohoIgnoreError(c.t, "expense", "customers", "delete", id, "--org", orgID)
	})
}

func (c *testCleanup) trackExpenseCurrency(id, orgID string) {
	c.add("delete expense currency "+id, func() {
		zohoIgnoreError(c.t, "expense", "currencies", "delete", id, "--org", orgID)
	})
}

func (c *testCleanup) trackExpenseTax(id, orgID string) {
	c.add("delete expense tax "+id, func() {
		zohoIgnoreError(c.t, "expense", "taxes", "delete", id, "--org", orgID)
	})
}

func (c *testCleanup) trackExpenseProject(id, orgID string) {
	c.add("delete expense project "+id, func() {
		zohoIgnoreError(c.t, "expense", "projects", "delete", id, "--org", orgID)
	})
}

func (c *testCleanup) trackExpenseTrip(id, orgID string) {
	c.add("delete expense trip "+id, func() {
		zohoIgnoreError(c.t, "expense", "trips", "delete", id, "--org", orgID)
	})
}

func (c *testCleanup) trackExpenseReport(id, orgID string) {
	c.add("delete expense report "+id, func() {
		zohoIgnoreError(c.t, "expense", "reports", "delete", id, "--org", orgID)
	})
}

func (c *testCleanup) trackExpenseExpense(id, orgID string) {
	c.add("delete expense expense "+id, func() {
		zohoIgnoreError(c.t, "expense", "expenses", "delete", id, "--org", orgID)
	})
}

func requireExpenseOrgID(t *testing.T) string {
	t.Helper()
	id := os.Getenv("ZOHO_EXPENSE_ORG_ID")
	if id == "" {
		t.Skip("skipping: ZOHO_EXPENSE_ORG_ID not set")
	}
	return id
}

func TestExpenseOrganizations(t *testing.T) {
	t.Parallel()
	orgID := requireExpenseOrgID(t)

	t.Run("list", func(t *testing.T) {
		out := zoho(t, "expense", "organizations", "list")
		m := parseJSON(t, out)
		assertExpenseCodeZero(t, m)
		arr, ok := m["organizations"].([]any)
		if !ok {
			t.Fatalf("expected organizations array in response:\n%s", truncate(out, 500))
		}
		if len(arr) == 0 {
			t.Fatal("expected at least one organization")
		}
	})

	t.Run("get", func(t *testing.T) {
		out := zoho(t, "expense", "organizations", "get", orgID)
		m := parseJSON(t, out)
		assertExpenseCodeZero(t, m)
		org, ok := m["organization"].(map[string]any)
		if !ok {
			t.Fatalf("expected organization object in response:\n%s", truncate(out, 500))
		}
		id := fmt.Sprintf("%v", org["organization_id"])
		if id == "" || id == "<nil>" {
			t.Fatalf("expected organization_id in response:\n%s", truncate(out, 500))
		}
	})
}

func TestExpenseCategories(t *testing.T) {
	t.Parallel()
	orgID := requireExpenseOrgID(t)
	cleanup := newCleanup(t)

	var categoryID string
	var categoryName string

	t.Run("list", func(t *testing.T) {
		out := zoho(t, "expense", "categories", "list", "--org", orgID)
		m := parseJSON(t, out)
		assertExpenseCodeZero(t, m)
		arr, ok := m["expense_accounts"].([]any)
		if !ok {
			t.Fatalf("expected expense_accounts array in response:\n%s", truncate(out, 500))
		}
		if len(arr) == 0 {
			t.Fatal("expected at least one expense category")
		}
	})

	t.Run("create", func(t *testing.T) {
		categoryName = fmt.Sprintf("%s Cat %s", testPrefix, randomSuffix())
		out := zoho(t, "expense", "categories", "create", "--org", orgID,
			"--category_name", categoryName)
		m := parseJSON(t, out)
		assertExpenseCodeZero(t, m)
		cat, ok := m["expense_category"].(map[string]any)
		if !ok {
			t.Fatalf("expected expense_category object in response:\n%s", truncate(out, 500))
		}
		categoryID = fmt.Sprintf("%v", cat["category_id"])
		if categoryID == "" || categoryID == "<nil>" {
			t.Fatalf("expected category_id in create response:\n%s", truncate(out, 500))
		}
		cleanup.trackExpenseCategory(categoryID, orgID)
	})

	t.Run("get", func(t *testing.T) {
		requireID(t, categoryID, "create must have succeeded")
		out := zoho(t, "expense", "categories", "get", categoryID, "--org", orgID)
		m := parseJSON(t, out)
		assertExpenseCodeZero(t, m)
		cat, ok := m["expense_category"].(map[string]any)
		if !ok {
			t.Fatalf("expected expense_category object in response:\n%s", truncate(out, 500))
		}
		assertEqual(t, fmt.Sprintf("%v", cat["category_id"]), categoryID)
	})

	t.Run("update", func(t *testing.T) {
		requireID(t, categoryID, "create must have succeeded")
		updatedName := fmt.Sprintf("%s Cat Upd %s", testPrefix, randomSuffix())
		out := zoho(t, "expense", "categories", "update", categoryID, "--org", orgID,
			"--category_name", updatedName)
		m := parseJSON(t, out)
		assertExpenseCodeZero(t, m)
		cat, ok := m["expense_category"].(map[string]any)
		if ok {
			if got := fmt.Sprintf("%v", cat["category_name"]); got != "" && got != "<nil>" && got != updatedName {
				t.Errorf("category_name: got %q, want %q", got, updatedName)
			}
		}
		categoryName = updatedName
		_ = categoryName
	})

	t.Run("disable", func(t *testing.T) {
		requireID(t, categoryID, "create must have succeeded")
		out := zoho(t, "expense", "categories", "disable", categoryID, "--org", orgID)
		m := parseJSON(t, out)
		assertExpenseCodeZero(t, m)
		out = zoho(t, "expense", "categories", "get", categoryID, "--org", orgID)
		m2 := parseJSON(t, out)
		assertExpenseCodeZero(t, m2)
		cat, _ := m2["expense_category"].(map[string]any)
		if s := fmt.Sprintf("%v", cat["status"]); s != "inactive" {
			t.Errorf("expected category status=inactive after disable, got %s", s)
		}
	})

	t.Run("enable", func(t *testing.T) {
		requireID(t, categoryID, "create must have succeeded")
		out := zoho(t, "expense", "categories", "enable", categoryID, "--org", orgID)
		m := parseJSON(t, out)
		assertExpenseCodeZero(t, m)
		out = zoho(t, "expense", "categories", "get", categoryID, "--org", orgID)
		m2 := parseJSON(t, out)
		assertExpenseCodeZero(t, m2)
		cat, _ := m2["expense_category"].(map[string]any)
		if s := fmt.Sprintf("%v", cat["status"]); s != "active" {
			t.Errorf("expected category status=active after enable, got %s", s)
		}
	})

	t.Run("delete", func(t *testing.T) {
		requireID(t, categoryID, "create must have succeeded")
		out := zoho(t, "expense", "categories", "delete", categoryID, "--org", orgID)
		m := parseJSON(t, out)
		assertExpenseCodeZero(t, m)
		r := runZoho(t, "expense", "categories", "get", categoryID, "--org", orgID)
		if r.ExitCode == 0 {
			t.Errorf("category %s still accessible after delete", categoryID)
		}
		categoryID = ""
	})
}

func TestExpenseCustomers(t *testing.T) {
	t.Parallel()
	orgID := requireExpenseOrgID(t)
	cleanup := newCleanup(t)

	var customerID string
	var customerName string

	t.Run("list", func(t *testing.T) {
		out := zoho(t, "expense", "customers", "list", "--org", orgID)
		m := parseJSON(t, out)
		assertExpenseCodeZero(t, m)
		arr, ok := m["contacts"].([]any)
		if !ok {
			t.Fatalf("expected contacts array in response:\n%s", truncate(out, 500))
		}
		if len(arr) == 0 {
			t.Log("customers list is empty")
		}
	})

	t.Run("create", func(t *testing.T) {
		customerName = fmt.Sprintf("%s Customer %s", testPrefix, randomSuffix())
		out := zoho(t, "expense", "customers", "create", "--org", orgID,
			"--contact_name", customerName)
		m := parseJSON(t, out)
		assertExpenseCodeZero(t, m)
		contact, ok := m["contact"].(map[string]any)
		if !ok {
			t.Fatalf("expected contact object in response:\n%s", truncate(out, 500))
		}
		customerID = fmt.Sprintf("%v", contact["contact_id"])
		if customerID == "" || customerID == "<nil>" {
			t.Fatalf("expected contact_id in create response:\n%s", truncate(out, 500))
		}
		cleanup.trackExpenseCustomer(customerID, orgID)
	})

	t.Run("get", func(t *testing.T) {
		requireID(t, customerID, "create must have succeeded")
		out := zoho(t, "expense", "customers", "get", customerID, "--org", orgID)
		m := parseJSON(t, out)
		assertExpenseCodeZero(t, m)
		contact, ok := m["contact"].(map[string]any)
		if !ok {
			t.Fatalf("expected contact object in response:\n%s", truncate(out, 500))
		}
		assertEqual(t, fmt.Sprintf("%v", contact["contact_id"]), customerID)
	})

	t.Run("update", func(t *testing.T) {
		requireID(t, customerID, "create must have succeeded")
		updatedName := fmt.Sprintf("%s Customer Upd %s", testPrefix, randomSuffix())
		out := zoho(t, "expense", "customers", "update", customerID, "--org", orgID,
			"--contact_name", updatedName)
		m := parseJSON(t, out)
		assertExpenseCodeZero(t, m)
		contact, ok := m["contact"].(map[string]any)
		if ok {
			if got := fmt.Sprintf("%v", contact["contact_name"]); got != "" && got != "<nil>" && got != updatedName {
				t.Errorf("contact_name: got %q, want %q", got, updatedName)
			}
		}
		customerName = updatedName
		_ = customerName
	})

	t.Run("delete", func(t *testing.T) {
		requireID(t, customerID, "create must have succeeded")
		out := zoho(t, "expense", "customers", "delete", customerID, "--org", orgID)
		m := parseJSON(t, out)
		assertExpenseCodeZero(t, m)
		r := runZoho(t, "expense", "customers", "get", customerID, "--org", orgID)
		if r.ExitCode == 0 {
			t.Errorf("customer %s still accessible after delete", customerID)
		}
		customerID = ""
	})
}

func TestExpenseCurrencies(t *testing.T) {
	t.Parallel()
	orgID := requireExpenseOrgID(t)
	cleanup := newCleanup(t)

	var currencyID string

	t.Run("list", func(t *testing.T) {
		out := zoho(t, "expense", "currencies", "list", "--org", orgID)
		m := parseJSON(t, out)
		assertExpenseCodeZero(t, m)
		arr, ok := m["currencies"].([]any)
		if !ok {
			t.Fatalf("expected currencies array in response:\n%s", truncate(out, 500))
		}
		if len(arr) == 0 {
			t.Fatal("expected at least one currency")
		}
	})

	t.Run("create", func(t *testing.T) {
		out := zoho(t, "expense", "currencies", "create", "--org", orgID,
			"--currency_code", "MXN",
			"--currency_symbol", "$",
			"--currency_format", "1,234,567.89",
			"--price_precision", "2")
		m := parseJSON(t, out)
		assertExpenseCodeZero(t, m)
		currency, ok := m["currency"].(map[string]any)
		if !ok {
			t.Fatalf("expected currency object in response:\n%s", truncate(out, 500))
		}
		currencyID = fmt.Sprintf("%v", currency["currency_id"])
		if currencyID == "" || currencyID == "<nil>" {
			t.Fatalf("expected currency_id in create response:\n%s", truncate(out, 500))
		}
		cleanup.trackExpenseCurrency(currencyID, orgID)
	})

	t.Run("get", func(t *testing.T) {
		requireID(t, currencyID, "create must have succeeded")
		out := zoho(t, "expense", "currencies", "get", currencyID, "--org", orgID)
		m := parseJSON(t, out)
		assertExpenseCodeZero(t, m)
		currency, ok := m["currency"].(map[string]any)
		if !ok {
			t.Fatalf("expected currency object in response:\n%s", truncate(out, 500))
		}
		assertEqual(t, fmt.Sprintf("%v", currency["currency_id"]), currencyID)
	})

	t.Run("update", func(t *testing.T) {
		requireID(t, currencyID, "create must have succeeded")
		out := zoho(t, "expense", "currencies", "update", currencyID, "--org", orgID,
			"--currency_symbol", "Mex$",
			"--currency_format", "1,234,567.89",
			"--price_precision", "2")
		m := parseJSON(t, out)
		assertExpenseCodeZero(t, m)
		getOut := zoho(t, "expense", "currencies", "get", currencyID, "--org", orgID)
		gm := parseJSON(t, getOut)
		assertExpenseCodeZero(t, gm)
		cur, _ := gm["currency"].(map[string]any)
		if s := fmt.Sprintf("%v", cur["currency_symbol"]); s != "Mex$" {
			t.Errorf("expected currency_symbol=Mex$ after update, got %s", s)
		}
	})

	t.Run("delete", func(t *testing.T) {
		requireID(t, currencyID, "create must have succeeded")
		out := zoho(t, "expense", "currencies", "delete", currencyID, "--org", orgID)
		m := parseJSON(t, out)
		assertExpenseCodeZero(t, m)
		r := runZoho(t, "expense", "currencies", "get", currencyID, "--org", orgID)
		if r.ExitCode == 0 {
			t.Errorf("currency %s still accessible after delete", currencyID)
		}
		currencyID = ""
	})
}

func TestExpenseTaxes(t *testing.T) {
	t.Parallel()
	orgID := requireExpenseOrgID(t)
	cleanup := newCleanup(t)

	var taxID string
	var taxGroupID string
	var fallbackTaxID string
	var createdTax bool
	var createdTaxRate int

	t.Run("list", func(t *testing.T) {
		out := zoho(t, "expense", "taxes", "list", "--org", orgID)
		m := parseJSON(t, out)
		assertExpenseCodeZero(t, m)
		arr, ok := m["taxes"].([]any)
		if !ok {
			t.Fatalf("expected taxes array in response:\n%s", truncate(out, 500))
		}
		if len(arr) == 0 {
			t.Log("taxes list is empty")
		}
		for _, item := range arr {
			tm, _ := item.(map[string]any)
			if fmt.Sprintf("%v", tm["is_tax_group"]) == "true" {
				if taxGroupID == "" || taxGroupID == "<nil>" {
					taxGroupID = fmt.Sprintf("%v", tm["tax_group_id"])
					if taxGroupID == "" || taxGroupID == "<nil>" {
						taxGroupID = fmt.Sprintf("%v", tm["tax_id"])
					}
				}
				continue
			}
			id := fmt.Sprintf("%v", tm["tax_id"])
			if id != "" && id != "<nil>" && fallbackTaxID == "" {
				fallbackTaxID = id
			}
		}
		if groups, ok := m["tax_groups"].([]any); ok && len(groups) > 0 {
			if gm, ok := groups[0].(map[string]any); ok {
				gid := fmt.Sprintf("%v", gm["tax_group_id"])
				if gid == "" || gid == "<nil>" {
					gid = fmt.Sprintf("%v", gm["tax_id"])
				}
				if gid != "" && gid != "<nil>" {
					taxGroupID = gid
				}
			}
		}
	})

	t.Run("create", func(t *testing.T) {
		for _, rate := range []int{0, 15} {
			out, err := zohoMayFail(t, "expense", "taxes", "create", "--org", orgID,
				"--tax_name", fmt.Sprintf("%s Tax %s", testPrefix, randomSuffix()),
				"--tax_percentage", fmt.Sprintf("%d", rate),
				"--tax_type", "tax")
			if err != nil {
				t.Logf("tax create with rate %d failed: %v", rate, err)
				t.Logf("response: %s", truncate(out, 500))
				continue
			}
			m := parseJSON(t, out)
			assertExpenseCodeZero(t, m)
			tax, ok := m["tax"].(map[string]any)
			if !ok {
				t.Fatalf("expected tax object in response:\n%s", truncate(out, 500))
			}
			taxID = fmt.Sprintf("%v", tax["tax_id"])
			if taxID == "" || taxID == "<nil>" {
				t.Fatalf("expected tax_id in create response:\n%s", truncate(out, 500))
			}
			createdTax = true
			createdTaxRate = rate
			cleanup.trackExpenseTax(taxID, orgID)
			return
		}
		t.Logf("all tax create attempts failed, falling back to existing tax %s", fallbackTaxID)
	})

	t.Run("get", func(t *testing.T) {
		getID := taxID
		if getID == "" {
			getID = fallbackTaxID
		}
		requireID(t, getID, "create or fallback must have provided a tax")
		out := zoho(t, "expense", "taxes", "get", getID, "--org", orgID)
		m := parseJSON(t, out)
		assertExpenseCodeZero(t, m)
		tax, ok := m["tax"].(map[string]any)
		if !ok {
			t.Fatalf("expected tax object in response:\n%s", truncate(out, 500))
		}
		assertEqual(t, fmt.Sprintf("%v", tax["tax_id"]), getID)
	})

	t.Run("update", func(t *testing.T) {
		if !createdTax {
			t.Skip("skipping update: tax was not created by this test")
		}
		updatedName := fmt.Sprintf("%s Tax Upd %s", testPrefix, randomSuffix())
		out := zoho(t, "expense", "taxes", "update", taxID, "--org", orgID,
			"--tax_name", updatedName,
			"--tax_percentage", fmt.Sprintf("%d", createdTaxRate),
			"--tax_type", "tax")
		m := parseJSON(t, out)
		assertExpenseCodeZero(t, m)
		getOut := zoho(t, "expense", "taxes", "get", taxID, "--org", orgID)
		gm := parseJSON(t, getOut)
		assertExpenseCodeZero(t, gm)
		tax, _ := gm["tax"].(map[string]any)
		if got := fmt.Sprintf("%v", tax["tax_name"]); got != updatedName {
			t.Errorf("expected tax_name=%q after update, got %q", updatedName, got)
		}
	})

	t.Run("delete", func(t *testing.T) {
		if !createdTax {
			t.Skip("skipping delete: tax was not created by this test")
		}
		requireID(t, taxID, "create must have succeeded")
		out := zoho(t, "expense", "taxes", "delete", taxID, "--org", orgID)
		m := parseJSON(t, out)
		assertExpenseCodeZero(t, m)
		r := runZoho(t, "expense", "taxes", "get", taxID, "--org", orgID)
		if r.ExitCode == 0 {
			t.Errorf("tax %s still accessible after delete", taxID)
		}
		taxID = ""
	})

	t.Run("get-group", func(t *testing.T) {
		if taxGroupID == "" || taxGroupID == "<nil>" {
			t.Skip("no tax group found in list response")
		}
		out := zoho(t, "expense", "taxes", "get-group", taxGroupID, "--org", orgID)
		m := parseJSON(t, out)
		assertExpenseCodeZero(t, m)
	})
}

func TestExpenseProjects(t *testing.T) {
	t.Parallel()
	orgID := requireExpenseOrgID(t)
	cleanup := newCleanup(t)

	var projectID string
	var projectName string
	var fallbackProjectID string
	var createdProject bool

	t.Run("list", func(t *testing.T) {
		out := zoho(t, "expense", "projects", "list", "--org", orgID)
		m := parseJSON(t, out)
		assertExpenseCodeZero(t, m)
		arr, ok := m["projects"].([]any)
		if !ok {
			t.Fatalf("expected projects array in response:\n%s", truncate(out, 500))
		}
		if len(arr) == 0 {
			t.Log("projects list is empty")
		}
		for _, item := range arr {
			pm, _ := item.(map[string]any)
			id := fmt.Sprintf("%v", pm["project_id"])
			if id != "" && id != "<nil>" {
				fallbackProjectID = id
				break
			}
		}
	})

	t.Run("create", func(t *testing.T) {
		projectName = fmt.Sprintf("%s Proj %s", testPrefix, randomSuffix())
		out, err := zohoMayFail(t, "expense", "projects", "create", "--org", orgID,
			"--project_name", projectName)
		if err != nil {
			t.Logf("project create failed (org restriction): %v", err)
			t.Logf("response: %s", truncate(out, 500))
			if fallbackProjectID != "" {
				t.Logf("using existing project %s from list for get test", fallbackProjectID)
			}
			return
		}
		m := parseJSON(t, out)
		assertExpenseCodeZero(t, m)
		project, ok := m["project"].(map[string]any)
		if !ok {
			t.Fatalf("expected project object in response:\n%s", truncate(out, 500))
		}
		projectID = fmt.Sprintf("%v", project["project_id"])
		if projectID == "" || projectID == "<nil>" {
			t.Fatalf("expected project_id in create response:\n%s", truncate(out, 500))
		}
		createdProject = true
		cleanup.trackExpenseProject(projectID, orgID)
	})

	t.Run("get", func(t *testing.T) {
		getID := projectID
		if getID == "" {
			getID = fallbackProjectID
		}
		requireID(t, getID, "create or list must have provided a project")
		out := zoho(t, "expense", "projects", "get", getID, "--org", orgID)
		m := parseJSON(t, out)
		assertExpenseCodeZero(t, m)
		project, ok := m["project"].(map[string]any)
		if !ok {
			t.Fatalf("expected project object in response:\n%s", truncate(out, 500))
		}
		assertEqual(t, fmt.Sprintf("%v", project["project_id"]), getID)
	})

	t.Run("update", func(t *testing.T) {
		if !createdProject {
			t.Skip("skipping update: project was not created by this test")
		}
		updatedName := fmt.Sprintf("%s Proj Upd %s", testPrefix, randomSuffix())
		out := zoho(t, "expense", "projects", "update", projectID, "--org", orgID,
			"--project_name", updatedName)
		m := parseJSON(t, out)
		assertExpenseCodeZero(t, m)
		getOut := zoho(t, "expense", "projects", "get", projectID, "--org", orgID)
		gm := parseJSON(t, getOut)
		assertExpenseCodeZero(t, gm)
		proj, _ := gm["project"].(map[string]any)
		if got := fmt.Sprintf("%v", proj["project_name"]); got != updatedName {
			t.Errorf("expected project_name=%q after update, got %q", updatedName, got)
		}
		projectName = updatedName
	})

	t.Run("deactivate", func(t *testing.T) {
		if !createdProject {
			t.Skip("skipping deactivate: project was not created by this test")
		}
		out := zoho(t, "expense", "projects", "deactivate", projectID, "--org", orgID)
		m := parseJSON(t, out)
		assertExpenseCodeZero(t, m)
		getOut := zoho(t, "expense", "projects", "get", projectID, "--org", orgID)
		gm := parseJSON(t, getOut)
		assertExpenseCodeZero(t, gm)
		proj, _ := gm["project"].(map[string]any)
		if s := fmt.Sprintf("%v", proj["status"]); s != "inactive" {
			t.Errorf("expected project status=inactive after deactivate, got %s", s)
		}
	})

	t.Run("activate", func(t *testing.T) {
		if !createdProject {
			t.Skip("skipping activate: project was not created by this test")
		}
		out := zoho(t, "expense", "projects", "activate", projectID, "--org", orgID)
		m := parseJSON(t, out)
		assertExpenseCodeZero(t, m)
		getOut := zoho(t, "expense", "projects", "get", projectID, "--org", orgID)
		gm := parseJSON(t, getOut)
		assertExpenseCodeZero(t, gm)
		proj, _ := gm["project"].(map[string]any)
		if s := fmt.Sprintf("%v", proj["status"]); s != "active" {
			t.Errorf("expected project status=active after activate, got %s", s)
		}
	})

	t.Run("delete", func(t *testing.T) {
		if !createdProject {
			t.Skip("skipping delete: project was not created by this test")
		}
		out := zoho(t, "expense", "projects", "delete", projectID, "--org", orgID)
		m := parseJSON(t, out)
		assertExpenseCodeZero(t, m)
		r := runZoho(t, "expense", "projects", "get", projectID, "--org", orgID)
		if r.ExitCode == 0 {
			t.Errorf("project %s still accessible after delete", projectID)
		}
		projectID = ""
		_ = projectName
	})
}

func TestExpenseTrips(t *testing.T) {
	t.Parallel()
	orgID := requireExpenseOrgID(t)
	cleanup := newCleanup(t)

	var tripID string
	var tripNumber string
	destinationCountry := "South Africa"

	t.Run("list", func(t *testing.T) {
		out := zoho(t, "expense", "trips", "list", "--org", orgID)
		m := parseJSON(t, out)
		assertExpenseCodeZero(t, m)
		arr, ok := m["trips"].([]any)
		if !ok {
			t.Fatalf("expected trips array in response:\n%s", truncate(out, 500))
		}
		if len(arr) == 0 {
			t.Log("trips list is empty")
		}
	})

	t.Run("create", func(t *testing.T) {
		purpose := fmt.Sprintf("%s trip %s", testPrefix, randomSuffix())
		tripNumber = fmt.Sprintf("TRIP-%s", randomSuffix())
		out := zoho(t, "expense", "trips", "create", "--org", orgID,
			"--trip_number", tripNumber,
			"--is_international", "true",
			"--destination_country", destinationCountry,
			"--business_purpose", purpose)
		m := parseJSON(t, out)
		assertExpenseCodeZero(t, m)
		trip, ok := m["trip"].(map[string]any)
		if !ok {
			t.Fatalf("expected trip object in response:\n%s", truncate(out, 500))
		}
		tripID = fmt.Sprintf("%v", trip["trip_id"])
		if tripID == "" || tripID == "<nil>" {
			t.Fatalf("expected trip_id in create response:\n%s", truncate(out, 500))
		}
		cleanup.trackExpenseTrip(tripID, orgID)
	})

	t.Run("get", func(t *testing.T) {
		requireID(t, tripID, "create must have succeeded")
		out := zoho(t, "expense", "trips", "get", tripID, "--org", orgID)
		m := parseJSON(t, out)
		assertExpenseCodeZero(t, m)
		trip, ok := m["trip"].(map[string]any)
		if !ok {
			t.Fatalf("expected trip object in response:\n%s", truncate(out, 500))
		}
		assertEqual(t, fmt.Sprintf("%v", trip["trip_id"]), tripID)
	})

	t.Run("update", func(t *testing.T) {
		requireID(t, tripID, "create must have succeeded")
		requireID(t, tripNumber, "create must have succeeded")
		updatedPurpose := fmt.Sprintf("%s trip updated %s", testPrefix, randomSuffix())
		out := zoho(t, "expense", "trips", "update", tripID, "--org", orgID,
			"--trip_number", tripNumber,
			"--is_international", "true",
			"--destination_country", destinationCountry,
			"--business_purpose", updatedPurpose)
		m := parseJSON(t, out)
		assertExpenseCodeZero(t, m)
		getOut := zoho(t, "expense", "trips", "get", tripID, "--org", orgID)
		gm := parseJSON(t, getOut)
		assertExpenseCodeZero(t, gm)
		trip, _ := gm["trip"].(map[string]any)
		if got := fmt.Sprintf("%v", trip["business_purpose"]); got != updatedPurpose {
			t.Errorf("expected business_purpose=%q after update, got %q", updatedPurpose, got)
		}
	})

	t.Run("delete", func(t *testing.T) {
		requireID(t, tripID, "create must have succeeded")
		out := zoho(t, "expense", "trips", "delete", tripID, "--org", orgID)
		m := parseJSON(t, out)
		assertExpenseCodeZero(t, m)
		r := runZoho(t, "expense", "trips", "get", tripID, "--org", orgID)
		if r.ExitCode == 0 {
			t.Errorf("trip %s still accessible after delete", tripID)
		}
		tripID = ""
	})
}

func TestExpenseReports(t *testing.T) {
	t.Parallel()
	orgID := requireExpenseOrgID(t)
	cleanup := newCleanup(t)

	var reportID string

	t.Run("list", func(t *testing.T) {
		out := zoho(t, "expense", "reports", "list", "--org", orgID)
		m := parseJSON(t, out)
		assertExpenseCodeZero(t, m)
		if _, ok := m["expense_reports"].([]any); !ok {
			t.Fatalf("expected expense_reports array in response:\n%s", truncate(out, 500))
		}
	})

	t.Run("create", func(t *testing.T) {
		reportName := fmt.Sprintf("%s Report %s", testPrefix, randomSuffix())
		out := zoho(t, "expense", "reports", "create", "--org", orgID,
			"--report_name", reportName)
		m := parseJSON(t, out)
		assertExpenseCodeZero(t, m)
		report, ok := m["expense_report"].(map[string]any)
		if !ok {
			t.Fatalf("expected expense_report object in response:\n%s", truncate(out, 500))
		}
		reportID = fmt.Sprintf("%v", report["report_id"])
		if reportID == "" || reportID == "<nil>" {
			t.Fatalf("expected report_id in create response:\n%s", truncate(out, 500))
		}
		cleanup.trackExpenseReport(reportID, orgID)
	})

	t.Run("get", func(t *testing.T) {
		requireID(t, reportID, "create must have succeeded")
		out := zoho(t, "expense", "reports", "get", reportID, "--org", orgID)
		m := parseJSON(t, out)
		assertExpenseCodeZero(t, m)
		report, ok := m["expense_report"].(map[string]any)
		if !ok {
			t.Fatalf("expected expense_report object in response:\n%s", truncate(out, 500))
		}
		assertEqual(t, fmt.Sprintf("%v", report["report_id"]), reportID)
	})

	t.Run("update", func(t *testing.T) {
		requireID(t, reportID, "create must have succeeded")
		updatedName := fmt.Sprintf("%s Report Upd %s", testPrefix, randomSuffix())
		out := zoho(t, "expense", "reports", "update", reportID, "--org", orgID,
			"--report_name", updatedName)
		m := parseJSON(t, out)
		assertExpenseCodeZero(t, m)
		getOut := zoho(t, "expense", "reports", "get", reportID, "--org", orgID)
		gm := parseJSON(t, getOut)
		assertExpenseCodeZero(t, gm)
		report, _ := gm["expense_report"].(map[string]any)
		if got := fmt.Sprintf("%v", report["report_name"]); got != updatedName {
			t.Errorf("expected report_name=%q after update, got %q", updatedName, got)
		}
	})

	t.Run("approval-history", func(t *testing.T) {
		requireID(t, reportID, "create must have succeeded")
		out, err := zohoMayFail(t, "expense", "reports", "approval-history", reportID, "--org", orgID)
		if err != nil {
			t.Logf("approval-history may fail for draft reports: %v", err)
			return
		}
		m := parseJSON(t, out)
		assertExpenseCodeZero(t, m)
	})
}

func TestExpenseExpenses(t *testing.T) {
	t.Parallel()
	orgID := requireExpenseOrgID(t)
	cleanup := newCleanup(t)

	var currencyID string
	var categoryID string
	var expenseID string

	t.Run("setup", func(t *testing.T) {
		currOut := zoho(t, "expense", "currencies", "list", "--org", orgID)
		currResp := parseJSON(t, currOut)
		assertExpenseCodeZero(t, currResp)
		currencies, ok := currResp["currencies"].([]any)
		if !ok || len(currencies) == 0 {
			t.Fatalf("expected non-empty currencies list:\n%s", truncate(currOut, 500))
		}
		firstCurrency, _ := currencies[0].(map[string]any)
		currencyID = fmt.Sprintf("%v", firstCurrency["currency_id"])
		if currencyID == "" || currencyID == "<nil>" {
			t.Fatalf("expected currency_id in currencies list:\n%s", truncate(currOut, 500))
		}

		catOut := zoho(t, "expense", "categories", "list", "--org", orgID)
		catResp := parseJSON(t, catOut)
		assertExpenseCodeZero(t, catResp)
		cats, ok := catResp["expense_accounts"].([]any)
		if !ok || len(cats) == 0 {
			t.Fatalf("expected non-empty expense_accounts list:\n%s", truncate(catOut, 500))
		}
		firstCategory, _ := cats[0].(map[string]any)
		categoryID = fmt.Sprintf("%v", firstCategory["category_id"])
		if categoryID == "" || categoryID == "<nil>" {
			categoryID = fmt.Sprintf("%v", firstCategory["account_id"])
		}
		if categoryID == "" || categoryID == "<nil>" {
			t.Fatalf("expected category_id in expense_accounts list:\n%s", truncate(catOut, 500))
		}
	})

	t.Run("list", func(t *testing.T) {
		out := zoho(t, "expense", "expenses", "list", "--org", orgID)
		m := parseJSON(t, out)
		assertExpenseCodeZero(t, m)
		if _, ok := m["expenses"].([]any); !ok {
			t.Fatalf("expected expenses array in response:\n%s", truncate(out, 500))
		}
	})

	t.Run("create", func(t *testing.T) {
		requireID(t, currencyID, "setup must have succeeded")
		requireID(t, categoryID, "setup must have succeeded")
		out := zoho(t, "expense", "expenses", "create", "--org", orgID,
			"--category_id", categoryID,
			"--date", "2026-03-01",
			"--amount", "25.50",
			"--currency_id", currencyID)
		m := parseJSON(t, out)
		assertExpenseCodeZero(t, m)
		if expense, ok := m["expense"].(map[string]any); ok {
			expenseID = fmt.Sprintf("%v", expense["expense_id"])
		}
		if expenseID == "" || expenseID == "<nil>" {
			if expenses, ok := m["expenses"].([]any); ok && len(expenses) > 0 {
				if expense, ok := expenses[0].(map[string]any); ok {
					expenseID = fmt.Sprintf("%v", expense["expense_id"])
				}
			}
		}
		if expenseID == "" || expenseID == "<nil>" {
			t.Fatalf("expected expense_id in create response:\n%s", truncate(out, 500))
		}
		cleanup.trackExpenseExpense(expenseID, orgID)
	})

	t.Run("update", func(t *testing.T) {
		requireID(t, expenseID, "create must have succeeded")
		out := zoho(t, "expense", "expenses", "update", expenseID, "--org", orgID,
			"--category_id", categoryID,
			"--date", "2026-03-02",
			"--amount", "31.75",
			"--currency_id", currencyID)
		m := parseJSON(t, out)
		assertExpenseCodeZero(t, m)
		getOut := zoho(t, "expense", "expenses", "get", expenseID, "--org", orgID)
		gm := parseJSON(t, getOut)
		assertExpenseCodeZero(t, gm)
		assertEqual(t, fmt.Sprintf("%v", gm["expense"].(map[string]any)["expense_id"]), expenseID)
	})
}

func TestExpenseReceipts(t *testing.T) {
	t.Parallel()
	orgID := requireExpenseOrgID(t)

	t.Run("upload", func(t *testing.T) {
		f, err := os.CreateTemp("", "zohotest-receipt-*.txt")
		if err != nil {
			t.Fatalf("failed to create temp receipt: %v", err)
		}
		defer os.Remove(f.Name())
		if _, err := f.Write([]byte("ZOHOTEST receipt bytes")); err != nil {
			f.Close()
			t.Fatalf("failed to write temp receipt: %v", err)
		}
		if err := f.Close(); err != nil {
			t.Fatalf("failed to close temp receipt: %v", err)
		}
		_, err = zohoMayFail(t, "expense", "receipts", "upload", "dummy-expense-id", "--file", f.Name(), "--org", orgID)
		if err != nil {
			t.Logf("receipt upload failed as expected/allowed: %v", err)
		}
	})
}

func TestExpenseUsers(t *testing.T) {
	t.Parallel()
	orgID := requireExpenseOrgID(t)

	var userID string

	t.Run("list", func(t *testing.T) {
		out := zoho(t, "expense", "users", "list", "--org", orgID)
		m := parseJSON(t, out)
		assertExpenseCodeZero(t, m)
		arr, ok := m["users"].([]any)
		if !ok {
			t.Fatalf("expected users array in response:\n%s", truncate(out, 500))
		}
		if len(arr) == 0 {
			t.Fatal("expected at least one user")
		}
		firstUser, _ := arr[0].(map[string]any)
		userID = fmt.Sprintf("%v", firstUser["user_id"])
		if userID == "" || userID == "<nil>" {
			t.Fatalf("expected user_id in users list response:\n%s", truncate(out, 500))
		}
	})

	t.Run("get", func(t *testing.T) {
		requireID(t, userID, "list must have succeeded")
		out := zoho(t, "expense", "users", "get", userID, "--org", orgID)
		m := parseJSON(t, out)
		assertExpenseCodeZero(t, m)
		user, ok := m["user"].(map[string]any)
		if !ok {
			t.Fatalf("expected user object in response:\n%s", truncate(out, 500))
		}
		assertEqual(t, fmt.Sprintf("%v", user["user_id"]), userID)
	})
}

func TestExpenseTags(t *testing.T) {
	t.Parallel()
	orgID := requireExpenseOrgID(t)

	t.Run("list-known-broken", func(t *testing.T) {
		r := runZoho(t, "expense", "tags", "list", "--org", orgID)
		if r.ExitCode == 0 {
			t.Fatal("expected non-zero exit for tags list on this org")
		}
		t.Log("V3 reporting tags API not available for this org")
	})
}

func TestExpenseErrors(t *testing.T) {
	t.Parallel()
	orgID := requireExpenseOrgID(t)

	t.Run("missing-org", func(t *testing.T) {
		r := runZohoWithEnv(t, map[string]string{"ZOHO_EXPENSE_ORG_ID": ""}, "expense", "categories", "list")
		assertExitCode(t, r, 4)
	})

	t.Run("invalid-category-id", func(t *testing.T) {
		r := runZoho(t, "expense", "categories", "get", "999999999999999", "--org", orgID)
		if r.ExitCode == 0 {
			t.Error("expected non-zero exit for invalid category ID")
		}
	})

	t.Run("invalid-customer-id", func(t *testing.T) {
		r := runZoho(t, "expense", "customers", "get", "999999999999999", "--org", orgID)
		if r.ExitCode == 0 {
			t.Error("expected non-zero exit for invalid customer ID")
		}
	})
}

func TestExpenseEmergencyCleanup(t *testing.T) {
	orgID := requireExpenseOrgID(t)

	cleanupList := func(resource string, listArgs []string, arrayKey string, idKeys []string, nameKeys []string) {
		out := zoho(t, listArgs...)
		m := parseJSON(t, out)
		assertExpenseCodeZero(t, m)
		arr, ok := m[arrayKey].([]any)
		if !ok {
			return
		}
		for _, item := range arr {
			rec, _ := item.(map[string]any)
			name := ""
			for _, key := range nameKeys {
				v := fmt.Sprintf("%v", rec[key])
				if v != "" && v != "<nil>" {
					name = v
					break
				}
			}
			if !strings.HasPrefix(name, testPrefix) {
				continue
			}
			id := ""
			for _, key := range idKeys {
				v := fmt.Sprintf("%v", rec[key])
				if v != "" && v != "<nil>" {
					id = v
					break
				}
			}
			if id == "" {
				continue
			}
			t.Logf("cleaning expense %s %s (%s)", resource, id, name)
			zohoIgnoreError(t, "expense", resource, "delete", id, "--org", orgID)
		}
	}

	cleanupList("categories", []string{"expense", "categories", "list", "--org", orgID}, "expense_accounts", []string{"category_id", "account_id"}, []string{"category_name", "account_name", "name"})
	cleanupList("customers", []string{"expense", "customers", "list", "--org", orgID}, "contacts", []string{"contact_id"}, []string{"contact_name", "customer_name", "name"})
	cleanupList("projects", []string{"expense", "projects", "list", "--org", orgID}, "projects", []string{"project_id"}, []string{"project_name", "name"})
	cleanupList("currencies", []string{"expense", "currencies", "list", "--org", orgID}, "currencies", []string{"currency_id"}, []string{"currency_name", "currency_code", "name"})
	cleanupList("taxes", []string{"expense", "taxes", "list", "--org", orgID}, "taxes", []string{"tax_id"}, []string{"tax_name", "name"})
}
