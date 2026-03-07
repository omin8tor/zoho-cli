//go:build integration

package internal_test

import (
	"fmt"
	"os"
	"strings"
	"testing"
)

func requireCreatorOwner(t *testing.T) string {
	t.Helper()
	owner := os.Getenv("ZOHO_CREATOR_OWNER")
	if owner == "" {
		t.Skip("skipping: ZOHO_CREATOR_OWNER not set")
	}
	return owner
}

func requireCreatorApp(t *testing.T) string {
	t.Helper()
	app := os.Getenv("ZOHO_CREATOR_APP")
	if app == "" {
		t.Skip("skipping: ZOHO_CREATOR_APP not set")
	}
	return app
}

func requireCreatorReport(t *testing.T) string {
	t.Helper()
	report := os.Getenv("ZOHO_CREATOR_REPORT")
	if report == "" {
		t.Skip("skipping: ZOHO_CREATOR_REPORT not set")
	}
	return report
}

func requireCreatorForm(t *testing.T) string {
	t.Helper()
	form := os.Getenv("ZOHO_CREATOR_FORM")
	if form == "" {
		t.Skip("skipping: ZOHO_CREATOR_FORM not set")
	}
	return form
}

func (c *testCleanup) trackCreatorRecord(id, owner, app, report string) {
	c.add("delete creator record "+id, func() {
		zohoIgnoreError(c.t, "creator", "records", "delete", id,
			"--owner", owner, "--app", app, "--report", report)
	})
}

func TestCreatorApplications(t *testing.T) {
	t.Parallel()
	owner := requireCreatorOwner(t)

	t.Run("list", func(t *testing.T) {
		out := zoho(t, "creator", "applications", "list", "--owner", owner)
		m := parseJSON(t, out)
		if _, ok := m["applications"]; ok {
			arr, ok := m["applications"].([]any)
			if !ok {
				t.Fatalf("expected applications array:\n%s", truncate(out, 500))
			}
			if len(arr) == 0 {
				t.Log("warning: no applications found")
			}
		}
	})
}

func TestCreatorReports(t *testing.T) {
	t.Parallel()
	owner := requireCreatorOwner(t)
	app := requireCreatorApp(t)

	t.Run("list", func(t *testing.T) {
		out := zoho(t, "creator", "reports", "list", "--owner", owner, "--app", app)
		m := parseJSON(t, out)
		if reports, ok := m["reports"].([]any); ok {
			if len(reports) == 0 {
				t.Log("warning: no reports found")
			}
		}
	})
}

func TestCreatorForms(t *testing.T) {
	t.Parallel()
	owner := requireCreatorOwner(t)
	app := requireCreatorApp(t)

	t.Run("list", func(t *testing.T) {
		out := zoho(t, "creator", "forms", "list", "--owner", owner, "--app", app)
		m := parseJSON(t, out)
		if forms, ok := m["forms"].([]any); ok {
			if len(forms) == 0 {
				t.Log("warning: no forms found")
			}
		}
	})
}

func TestCreatorFields(t *testing.T) {
	t.Parallel()
	owner := requireCreatorOwner(t)
	app := requireCreatorApp(t)
	form := requireCreatorForm(t)

	t.Run("list", func(t *testing.T) {
		out := zoho(t, "creator", "fields", "list", "--owner", owner, "--app", app, "--form", form)
		m := parseJSON(t, out)
		if fields, ok := m["fields"].([]any); ok {
			if len(fields) == 0 {
				t.Log("warning: no fields found")
			}
		}
	})
}

func TestCreatorPages(t *testing.T) {
	t.Parallel()
	owner := requireCreatorOwner(t)
	app := requireCreatorApp(t)

	t.Run("list", func(t *testing.T) {
		out := zoho(t, "creator", "pages", "list", "--owner", owner, "--app", app)
		_ = parseJSON(t, out)
	})
}

func TestCreatorSections(t *testing.T) {
	t.Parallel()
	owner := requireCreatorOwner(t)
	app := requireCreatorApp(t)

	t.Run("list", func(t *testing.T) {
		out := zoho(t, "creator", "sections", "list", "--owner", owner, "--app", app)
		_ = parseJSON(t, out)
	})
}

func TestCreatorRecords(t *testing.T) {
	t.Parallel()
	owner := requireCreatorOwner(t)
	app := requireCreatorApp(t)
	report := requireCreatorReport(t)
	form := requireCreatorForm(t)
	cleanup := newCleanup(t)

	var recordID string

	t.Run("list", func(t *testing.T) {
		out := zoho(t, "creator", "records", "list",
			"--owner", owner, "--app", app, "--report", report)
		m := parseJSON(t, out)
		if data, ok := m["data"].([]any); ok {
			if len(data) == 0 {
				t.Log("warning: no records found")
			}
		}
	})

	t.Run("add", func(t *testing.T) {
		name := fmt.Sprintf("%s Rec %s", testPrefix, randomSuffix())
		out := zoho(t, "creator", "records", "add",
			"--owner", owner, "--app", app, "--form", form,
			"--json", toJSON(t, map[string]any{
				"data": map[string]any{
					"Name": map[string]any{
						"first_name": name,
						"last_name":  "Test",
					},
				},
			}))
		m := parseJSON(t, out)
		if data, ok := m["data"].([]any); ok && len(data) > 0 {
			rec := data[0].(map[string]any)
			if id, ok := rec["Added_Id"]; ok {
				recordID = fmt.Sprintf("%v", id)
			} else if id, ok := rec["ID"]; ok {
				recordID = fmt.Sprintf("%v", id)
			}
		} else if id, ok := m["ID"]; ok {
			recordID = fmt.Sprintf("%v", id)
		}
		if recordID != "" && recordID != "<nil>" {
			cleanup.trackCreatorRecord(recordID, owner, app, report)
		}
	})

	t.Run("get", func(t *testing.T) {
		requireID(t, recordID, "add must have succeeded")
		out := zoho(t, "creator", "records", "get", recordID,
			"--owner", owner, "--app", app, "--report", report)
		m := parseJSON(t, out)
		if data, ok := m["data"].(map[string]any); ok {
			id := fmt.Sprintf("%v", data["ID"])
			assertEqual(t, id, recordID)
		}
	})

	t.Run("update", func(t *testing.T) {
		requireID(t, recordID, "add must have succeeded")
		updatedName := fmt.Sprintf("%s RecUpd %s", testPrefix, randomSuffix())
		out := zoho(t, "creator", "records", "update", recordID,
			"--owner", owner, "--app", app, "--report", report,
			"--json", toJSON(t, map[string]any{
				"data": map[string]any{
					"Name": map[string]any{
						"first_name": updatedName,
						"last_name":  "Updated",
					},
				},
			}))
		_ = parseJSON(t, out)
	})

	t.Run("delete", func(t *testing.T) {
		requireID(t, recordID, "add must have succeeded")
		out := zoho(t, "creator", "records", "delete", recordID,
			"--owner", owner, "--app", app, "--report", report)
		_ = parseJSON(t, out)
		recordID = ""
	})
}

func TestCreatorRecordsListPagination(t *testing.T) {
	t.Parallel()
	owner := requireCreatorOwner(t)
	app := requireCreatorApp(t)
	report := requireCreatorReport(t)

	t.Run("with-max-records", func(t *testing.T) {
		out := zoho(t, "creator", "records", "list",
			"--owner", owner, "--app", app, "--report", report,
			"--max-records", "1")
		m := parseJSON(t, out)
		if data, ok := m["data"].([]any); ok {
			if len(data) > 1 {
				t.Errorf("expected at most 1 record, got %d", len(data))
			}
		}
	})
}

func TestCreatorErrors(t *testing.T) {
	t.Parallel()

	t.Run("missing-owner", func(t *testing.T) {
		r := runZohoWithEnv(t, map[string]string{"ZOHO_CREATOR_OWNER": ""}, "creator", "applications", "list")
		if r.ExitCode == 0 {
			t.Error("expected non-zero exit code when --owner missing")
		}
		if !strings.Contains(r.Stderr, "ZOHO_CREATOR_OWNER") {
			t.Errorf("expected error mentioning ZOHO_CREATOR_OWNER, got: %s", r.Stderr)
		}
	})

	t.Run("missing-app", func(t *testing.T) {
		owner := os.Getenv("ZOHO_CREATOR_OWNER")
		if owner == "" {
			t.Skip("ZOHO_CREATOR_OWNER not set")
		}
		r := runZohoWithEnv(t, map[string]string{"ZOHO_CREATOR_APP": ""},
			"creator", "reports", "list", "--owner", owner)
		if r.ExitCode == 0 {
			t.Error("expected non-zero exit code when --app missing")
		}
		if !strings.Contains(r.Stderr, "ZOHO_CREATOR_APP") {
			t.Errorf("expected error mentioning ZOHO_CREATOR_APP, got: %s", r.Stderr)
		}
	})

	t.Run("missing-record-id", func(t *testing.T) {
		owner := os.Getenv("ZOHO_CREATOR_OWNER")
		app := os.Getenv("ZOHO_CREATOR_APP")
		if owner == "" || app == "" {
			t.Skip("ZOHO_CREATOR_OWNER or ZOHO_CREATOR_APP not set")
		}
		r := runZoho(t, "creator", "records", "get", "--owner", owner, "--app", app, "--report", "All_Tasks")
		if r.ExitCode == 0 {
			t.Error("expected non-zero exit code when record-id missing")
		}
	})

	t.Run("missing-json-flag", func(t *testing.T) {
		owner := os.Getenv("ZOHO_CREATOR_OWNER")
		app := os.Getenv("ZOHO_CREATOR_APP")
		if owner == "" || app == "" {
			t.Skip("ZOHO_CREATOR_OWNER or ZOHO_CREATOR_APP not set")
		}
		r := runZoho(t, "creator", "records", "add", "--owner", owner, "--app", app, "--form", "Test_Form")
		if r.ExitCode == 0 {
			t.Error("expected non-zero exit code when --json missing")
		}
	})

	t.Run("invalid-owner", func(t *testing.T) {
		r := runZoho(t, "creator", "applications", "list", "--owner", "nonexistent_owner_99999")
		if r.ExitCode == 0 {
			t.Log("warning: invalid owner did not cause error (API may be lenient)")
		}
	})
}

func TestCreatorEmergencyCleanup(t *testing.T) {
	t.Parallel()
	if os.Getenv("ZOHO_EMERGENCY_CLEANUP") != "1" {
		t.Skip("set ZOHO_EMERGENCY_CLEANUP=1 to run")
	}
	owner := requireCreatorOwner(t)
	app := requireCreatorApp(t)
	report := requireCreatorReport(t)

	out, err := zohoMayFail(t, "creator", "records", "list",
		"--owner", owner, "--app", app, "--report", report)
	if err == nil {
		m := parseJSON(t, out)
		if data, ok := m["data"].([]any); ok {
			for _, item := range data {
				rec, _ := item.(map[string]any)
				name := ""
				if n, ok := rec["Name"].(map[string]any); ok {
					name = fmt.Sprintf("%v", n["first_name"])
				} else if n, ok := rec["Name"].(string); ok {
					name = n
				} else if n, ok := rec["Single_Line"].(string); ok {
					name = n
				}
				if strings.HasPrefix(name, testPrefix) {
					id := fmt.Sprintf("%v", rec["ID"])
					zohoIgnoreError(t, "creator", "records", "delete", id,
						"--owner", owner, "--app", app, "--report", report)
				}
			}
		}
	}
}
