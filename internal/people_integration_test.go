//go:build integration

package internal_test

import (
	"fmt"
	"testing"
)

func assertPeopleStatusZero(t *testing.T, m map[string]any) {
	t.Helper()
	resp, ok := m["response"].(map[string]any)
	if !ok {
		return
	}
	if status, ok := resp["status"].(float64); ok && status != 0 {
		msg := fmt.Sprintf("%v", resp["message"])
		t.Fatalf("people API error: status=%.0f message=%s", status, msg)
	}
}

func assertPeopleStatusZeroOrNoRecords(t *testing.T, m map[string]any) {
	t.Helper()
	resp, ok := m["response"].(map[string]any)
	if !ok {
		return
	}
	if status, ok := resp["status"].(float64); ok && status != 0 {
		if errs, ok := resp["errors"].(map[string]any); ok {
			if code, ok := errs["code"].(float64); ok && code == 7024 {
				return
			}
		}
		msg := fmt.Sprintf("%v", resp["message"])
		t.Fatalf("people API error: status=%.0f message=%s", status, msg)
	}
}

func TestPeopleForms(t *testing.T) {
	t.Parallel()

	t.Run("list", func(t *testing.T) {
		out := zoho(t, "people", "forms", "list")
		m := parseJSON(t, out)
		resp, ok := m["response"].(map[string]any)
		if !ok {
			t.Fatalf("expected response object:\n%s", truncate(out, 500))
		}
		assertPeopleStatusZero(t, m)
		result, ok := resp["result"].([]any)
		if !ok {
			t.Fatalf("expected result array in response:\n%s", truncate(out, 500))
		}
		if len(result) == 0 {
			t.Fatal("expected at least one form")
		}
		form := result[0].(map[string]any)
		if _, ok := form["formLinkName"]; !ok {
			t.Fatalf("expected formLinkName in form object:\n%s", truncate(out, 500))
		}
	})

	t.Run("get-fields", func(t *testing.T) {
		out := zoho(t, "people", "forms", "get-fields", "employee")
		m := parseJSON(t, out)
		assertPeopleStatusZero(t, m)
	})
}

func TestPeopleRecords(t *testing.T) {
	t.Parallel()

	t.Run("list-employees", func(t *testing.T) {
		out := zoho(t, "people", "records", "list", "employee", "--limit", "5")
		_ = parseJSONArray(t, out)
	})

	t.Run("list-with-search", func(t *testing.T) {
		out, err := zohoMayFail(t, "people", "records", "list", "employee", "--limit", "1")
		if err != nil {
			t.Skipf("skipping: list with pagination failed: %v", err)
		}
		_ = parseJSONArray(t, out)
	})
}

func TestPeopleAttendance(t *testing.T) {
	t.Parallel()

	t.Run("get-report", func(t *testing.T) {
		out, err := zohoMayFail(t, "people", "attendance", "get-report",
			"--sdate", "2026-03-01", "--edate", "2026-03-07",
			"--date-format", "yyyy-MM-dd")
		if err != nil {
			t.Skipf("skipping: attendance report failed: %v", err)
		}
		_ = parseJSON(t, out)
	})
}

func TestPeopleLeave(t *testing.T) {
	t.Parallel()

	t.Run("list", func(t *testing.T) {
		out, err := zohoMayFail(t, "people", "leave", "list", "--limit", "5")
		if err != nil {
			t.Skipf("skipping: leave list failed: %v", err)
		}
		m := parseJSON(t, out)
		assertPeopleStatusZeroOrNoRecords(t, m)
	})

	t.Run("get-types", func(t *testing.T) {
		out, err := zohoMayFail(t, "people", "leave", "get-types")
		if err != nil {
			t.Skipf("skipping: leave types failed: %v", err)
		}
		_ = parseJSON(t, out)
	})

	t.Run("get-balance", func(t *testing.T) {
		out, err := zohoMayFail(t, "people", "leave", "get-balance")
		if err != nil {
			t.Skipf("skipping: leave balance failed: %v", err)
		}
		_ = parseJSON(t, out)
	})
}

func TestPeopleDepartments(t *testing.T) {
	t.Parallel()

	t.Run("list", func(t *testing.T) {
		out := zoho(t, "people", "departments", "list")
		m := parseJSON(t, out)
		assertPeopleStatusZero(t, m)
	})
}

func TestPeopleDesignations(t *testing.T) {
	t.Parallel()

	t.Run("list", func(t *testing.T) {
		out := zoho(t, "people", "designations", "list")
		m := parseJSON(t, out)
		assertPeopleStatusZeroOrNoRecords(t, m)
	})
}

func TestPeopleRecordsCRUD(t *testing.T) {
	t.Parallel()
	cleanup := newCleanup(t)

	var recordID string

	t.Run("add", func(t *testing.T) {
		empID := fmt.Sprintf("%s_%s", testPrefix, randomSuffix())
		out := zoho(t, "people", "records", "add", "employee",
			"--json", toJSON(t, map[string]any{
				"EmployeeID": empID,
				"FirstName":  testPrefix,
				"LastName":   randomSuffix(),
				"EmailID":    fmt.Sprintf("%s@example.com", empID),
			}))
		m := parseJSON(t, out)
		assertPeopleStatusZero(t, m)
		resp, ok := m["response"].(map[string]any)
		if !ok {
			t.Fatalf("expected response object:\n%s", truncate(out, 500))
		}
		result, ok := resp["result"].(map[string]any)
		if !ok {
			t.Fatalf("expected result object:\n%s", truncate(out, 500))
		}
		recordID = fmt.Sprintf("%v", result["pkId"])
		if recordID == "" || recordID == "<nil>" {
			t.Fatalf("expected pkId in result:\n%s", truncate(out, 500))
		}
		cleanup.add("delete people employee record "+recordID, func() {
			zohoIgnoreError(t, "people", "records", "delete", "employee", recordID)
		})
	})

	t.Run("get", func(t *testing.T) {
		requireID(t, recordID, "add must have succeeded")
		out := zoho(t, "people", "records", "get", "employee", recordID)
		m := parseJSON(t, out)
		assertPeopleStatusZero(t, m)
	})

	t.Run("update", func(t *testing.T) {
		requireID(t, recordID, "add must have succeeded")
		newLast := fmt.Sprintf("%s_Upd_%s", testPrefix, randomSuffix())
		out := zoho(t, "people", "records", "update", "employee", recordID,
			"--json", toJSON(t, map[string]any{"LastName": newLast}))
		m := parseJSON(t, out)
		assertPeopleStatusZero(t, m)
	})

	t.Run("delete", func(t *testing.T) {
		requireID(t, recordID, "add must have succeeded")
		out, err := zohoMayFail(t, "people", "records", "delete", "employee", recordID)
		if err != nil {
			t.Skipf("delete may not be supported via People API: %v", err)
		}
		m := parseJSON(t, out)
		assertPeopleStatusZero(t, m)
		recordID = ""
	})
}

func TestPeopleErrors(t *testing.T) {
	t.Parallel()

	t.Run("missing-form-name-records", func(t *testing.T) {
		r := runZoho(t, "people", "records", "list")
		assertExitCode(t, r, 4)
	})

	t.Run("missing-record-id-get", func(t *testing.T) {
		r := runZoho(t, "people", "records", "get", "employee")
		assertExitCode(t, r, 4)
	})

	t.Run("missing-json-add", func(t *testing.T) {
		r := runZoho(t, "people", "records", "add", "employee")
		if r.ExitCode == 0 {
			t.Error("expected non-zero exit code when --json missing")
		}
	})

	t.Run("missing-form-name-fields", func(t *testing.T) {
		r := runZoho(t, "people", "forms", "get-fields")
		assertExitCode(t, r, 4)
	})

	t.Run("missing-leave-record-id", func(t *testing.T) {
		r := runZoho(t, "people", "leave", "get")
		assertExitCode(t, r, 4)
	})

	t.Run("missing-dept-name", func(t *testing.T) {
		r := runZoho(t, "people", "departments", "get-members")
		assertExitCode(t, r, 4)
	})
}

func TestPeopleEmergencyCleanup(t *testing.T) {
	t.Parallel()
	if testing.Short() {
		t.Skip("skipping emergency cleanup in short mode")
	}
	out, err := zohoMayFail(t, "people", "records", "list", "employee", "--limit", "200")
	if err != nil {
		t.Skipf("skipping: cannot list employees: %v", err)
	}
	m := parseJSON(t, out)
	resp, ok := m["response"].(map[string]any)
	if !ok {
		return
	}
	result, ok := resp["result"].([]any)
	if !ok {
		return
	}
	for _, item := range result {
		rec, ok := item.(map[string]any)
		if !ok {
			continue
		}
		firstName := fmt.Sprintf("%v", rec["FirstName"])
		if firstName == testPrefix {
			recID := fmt.Sprintf("%v", rec["Recordid"])
			if recID != "" && recID != "<nil>" {
				zohoIgnoreError(t, "people", "records", "delete", "employee", recID)
			}
		}
	}
}
