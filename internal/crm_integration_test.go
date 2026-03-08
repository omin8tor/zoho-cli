//go:build integration

package internal_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"
)

func getRecord(t *testing.T, module, id, fields string) map[string]any {
	t.Helper()
	out := zoho(t, "crm", "records", "get", module, id, "--fields", fields)
	return parseJSON(t, out)
}

func getRecordMayFail(t *testing.T, module, id string) (map[string]any, error) {
	t.Helper()
	r := runZoho(t, "crm", "records", "get", module, id, "--fields", "id")
	if r.ExitCode != 0 {
		return nil, fmt.Errorf("exit %d: %s", r.ExitCode, r.Stderr)
	}
	parsed := parseJSON(t, r.Stdout)
	if _, ok := parsed["id"]; !ok {
		return nil, fmt.Errorf("record not found (no id in response): %s", truncate(r.Stdout, 200))
	}
	return parsed, nil
}

func getNotes(t *testing.T, module, recordID string) []map[string]any {
	t.Helper()
	out := zoho(t, "crm", "notes", "list", module, recordID,
		"--fields", "id,Note_Title,Note_Content")
	return parseJSONArray(t, out)
}

func getAttachments(t *testing.T, module, recordID string) []map[string]any {
	t.Helper()
	out := zoho(t, "crm", "attachments", "list", module, recordID)
	return parseJSONArray(t, out)
}

func findInArray(arr []map[string]any, id string) (map[string]any, bool) {
	for _, item := range arr {
		if fmt.Sprintf("%v", item["id"]) == id {
			return item, true
		}
	}
	return nil, false
}

func hasTag(rec map[string]any, tagName string) bool {
	tags, ok := rec["Tag"].([]any)
	if !ok {
		return false
	}
	for _, tag := range tags {
		tagMap, _ := tag.(map[string]any)
		if fmt.Sprintf("%v", tagMap["name"]) == tagName {
			return true
		}
	}
	return false
}

func (c *testCleanup) trackLead(id string) {
	c.add("delete lead "+id, func() {
		zohoIgnoreError(c.t, "crm", "records", "delete", "Leads", id)
	})
}

func (c *testCleanup) trackNote(id string) {
	c.add("delete note "+id, func() {
		zohoIgnoreError(c.t, "crm", "notes", "delete", id)
	})
}

func (c *testCleanup) trackAttachment(module, recordID, attID string) {
	c.add("delete attachment "+attID, func() {
		zohoIgnoreError(c.t, "crm", "attachments", "delete", module, recordID, attID)
	})
}

func (c *testCleanup) trackContact(id string) {
	c.add("delete contact "+id, func() {
		zohoIgnoreError(c.t, "crm", "records", "delete", "Contacts", id)
	})
}

func (c *testCleanup) trackAccount(id string) {
	c.add("delete account "+id, func() {
		zohoIgnoreError(c.t, "crm", "records", "delete", "Accounts", id)
	})
}

type convertResult struct {
	ContactID string
	AccountID string
	DealID    string
}

func extractConvertIDs(t *testing.T, out string) convertResult {
	t.Helper()
	var resp struct {
		Data []struct {
			Code    string `json:"code"`
			Status  string `json:"status"`
			Message string `json:"message"`
			Details struct {
				Contacts *struct {
					ID   string `json:"id"`
					Name string `json:"name"`
				} `json:"Contacts"`
				Accounts *struct {
					ID   string `json:"id"`
					Name string `json:"name"`
				} `json:"Accounts"`
				Deals *struct {
					ID   string `json:"id"`
					Name string `json:"name"`
				} `json:"Deals"`
			} `json:"details"`
		} `json:"data"`
	}
	if err := json.Unmarshal([]byte(out), &resp); err != nil {
		t.Fatalf("failed to parse convert response: %v\nraw: %s", err, truncate(out, 500))
	}
	if len(resp.Data) == 0 {
		t.Fatalf("no data in convert response:\n%s", truncate(out, 500))
	}
	if resp.Data[0].Status != "success" {
		t.Fatalf("convert not successful: code=%s message=%s\n%s",
			resp.Data[0].Code, resp.Data[0].Message, truncate(out, 500))
	}
	d := resp.Data[0].Details
	if d.Contacts == nil || d.Contacts.ID == "" {
		t.Fatalf("no Contact ID in convert response:\n%s", truncate(out, 500))
	}
	if d.Accounts == nil || d.Accounts.ID == "" {
		t.Fatalf("no Account ID in convert response:\n%s", truncate(out, 500))
	}
	result := convertResult{
		ContactID: d.Contacts.ID,
		AccountID: d.Accounts.ID,
	}
	if d.Deals != nil {
		result.DealID = d.Deals.ID
	}
	return result
}

func TestCRMModules(t *testing.T) {
	t.Parallel()
	t.Run("list", func(t *testing.T) {
		out := zoho(t, "crm", "modules", "list")
		arr := parseJSONArray(t, out)
		assertNonEmpty(t, arr, "expected modules list")
		names := make(map[string]bool)
		for _, m := range arr {
			if n, ok := m["api_name"].(string); ok {
				names[n] = true
			}
		}
		for _, want := range []string{"Leads", "Contacts", "Accounts", "Deals"} {
			if !names[want] {
				t.Errorf("expected module %s in list", want)
			}
		}
	})

	t.Run("list-include-hidden", func(t *testing.T) {
		allOut := zoho(t, "crm", "modules", "list", "--include-hidden")
		allModules := parseJSONArray(t, allOut)
		visibleOut := zoho(t, "crm", "modules", "list")
		visibleModules := parseJSONArray(t, visibleOut)
		if len(allModules) <= len(visibleModules) {
			t.Errorf("--include-hidden should show more modules: all=%d visible=%d",
				len(allModules), len(visibleModules))
		}
	})

	t.Run("fields", func(t *testing.T) {
		out := zoho(t, "crm", "modules", "fields", "Leads")
		arr := parseJSONArray(t, out)
		assertNonEmpty(t, arr, "expected fields for Leads")
		names := make(map[string]bool)
		for _, f := range arr {
			if n, ok := f["api_name"].(string); ok {
				names[n] = true
			}
		}
		for _, want := range []string{"Last_Name", "Company", "Email"} {
			if !names[want] {
				t.Errorf("expected field %s in Leads fields list", want)
			}
		}
	})

	t.Run("related-lists", func(t *testing.T) {
		out := zoho(t, "crm", "modules", "related-lists", "Leads")
		arr := parseJSONArray(t, out)
		assertNonEmpty(t, arr, "expected related lists for Leads")
		names := make(map[string]bool)
		for _, rl := range arr {
			if _, ok := rl["id"]; !ok {
				t.Errorf("related list missing 'id' key: %v", rl)
			}
			if n, ok := rl["api_name"].(string); ok {
				names[n] = true
			} else {
				t.Errorf("related list missing 'api_name' key: %v", rl)
			}
		}
		if !names["Notes"] {
			t.Error("expected 'Notes' in related lists for Leads")
		}
		if !names["Attachments"] {
			t.Error("expected 'Attachments' in related lists for Leads")
		}
	})

	t.Run("layouts", func(t *testing.T) {
		out := zoho(t, "crm", "modules", "layouts", "Leads")
		arr := parseJSONArray(t, out)
		assertNonEmpty(t, arr, "expected layouts for Leads")
		foundStandard := false
		for _, layout := range arr {
			if _, ok := layout["id"]; !ok {
				t.Errorf("layout missing 'id' key")
			}
			if _, ok := layout["name"]; !ok {
				t.Errorf("layout missing 'name' key")
			}
			if fmt.Sprintf("%v", layout["name"]) == "Standard" {
				foundStandard = true
			}
		}
		if !foundStandard {
			t.Error("expected 'Standard' layout for Leads")
		}
	})

	t.Run("custom-views", func(t *testing.T) {
		out := zoho(t, "crm", "modules", "custom-views", "Leads")
		arr := parseJSONArray(t, out)
		assertNonEmpty(t, arr, "expected custom views for Leads")
		names := make(map[string]bool)
		for _, cv := range arr {
			if _, ok := cv["id"]; !ok {
				t.Errorf("custom view missing 'id' key")
			}
			if _, ok := cv["display_value"]; !ok {
				t.Errorf("custom view missing 'display_value' key")
			}
			if n, ok := cv["display_value"].(string); ok {
				names[n] = true
			}
		}
		if !names["All Leads"] {
			t.Error("expected 'All Leads' custom view for Leads")
		}
	})
}

func TestCRM(t *testing.T) {
	t.Parallel()
	cleanup := newCleanup(t)

	var leadID string
	var leadName string
	var leadEmail string
	var leadPhone string
	var upsertLeadID string
	var upsertName string
	var upsertEmail string
	var noteID string
	var attachmentID string
	var testFileContent []byte

	t.Run("users/list", func(t *testing.T) {
		out := zoho(t, "crm", "users", "list")
		arr := parseJSONArray(t, out)
		assertNonEmpty(t, arr, "expected at least one user")
		t.Logf("found user ID: %v", arr[0]["id"])
	})

	t.Run("records/get-nonexistent", func(t *testing.T) {
		_, err := getRecordMayFail(t, "Leads", "999999999999999999")
		if err == nil {
			t.Error("expected error for nonexistent record, but got a valid record")
		}
	})

	t.Run("records/create", func(t *testing.T) {
		leadName = testName(t)
		leadEmail = strings.ToLower(leadName) + "@test.example.com"
		leadPhone = fmt.Sprintf("555%07d", time.Now().UnixNano()%10000000)

		data := toJSON(t, map[string]any{"Last_Name": leadName, "Company": "TestCorp", "Email": leadEmail, "Phone": leadPhone})
		out := zoho(t, "crm", "records", "create", "Leads", "--json", data)
		leadID = extractID(t, out)
		cleanup.trackLead(leadID)

		rec := getRecord(t, "Leads", leadID, "id,Last_Name,Company,Email,Phone")
		assertEqual(t, fmt.Sprintf("%v", rec["id"]), leadID)
		assertStringField(t, rec, "Last_Name", leadName)
		assertStringField(t, rec, "Company", "TestCorp")
		assertStringField(t, rec, "Email", leadEmail)
		assertEqual(t, fmt.Sprintf("%v", rec["Phone"]), leadPhone)

		retryUntil(t, 10*time.Second, func() bool {
			query := fmt.Sprintf("select id, Last_Name from Leads where id = '%s'", leadID)
			coqlOut, coqlErr := zohoMayFail(t, "crm", "coql", "--query", query)
			if coqlErr != nil {
				return false
			}
			coqlParsed := parseJSON(t, coqlOut)
			coqlData, ok := coqlParsed["data"].([]any)
			if !ok || len(coqlData) == 0 {
				return false
			}
			coqlRec, _ := coqlData[0].(map[string]any)
			return fmt.Sprintf("%v", coqlRec["Last_Name"]) == leadName
		})
	})

	t.Run("records/get", func(t *testing.T) {
		requireID(t, leadID, "create must have succeeded")
		rec := getRecord(t, "Leads", leadID, "id,Last_Name,Company,Email")
		assertEqual(t, fmt.Sprintf("%v", rec["id"]), leadID)
		assertEqual(t, rec["Last_Name"], leadName)
		assertEqual(t, rec["Company"], "TestCorp")
		assertEqual(t, rec["Email"], leadEmail)
	})

	t.Run("records/list", func(t *testing.T) {
		out := zoho(t, "crm", "records", "list", "Leads",
			"--fields", "id,Last_Name,Created_Time",
			"--sort-by", "Created_Time",
			"--sort-order", "desc",
			"--limit", "5")
		arr := parseJSONArray(t, out)
		assertNonEmpty(t, arr, "expected at least one lead in list")

		if len(arr) > 5 {
			t.Errorf("--limit 5 but got %d records", len(arr))
		}

		for i := 1; i < len(arr); i++ {
			prev := fmt.Sprintf("%v", arr[i-1]["Created_Time"])
			curr := fmt.Sprintf("%v", arr[i]["Created_Time"])
			if curr > prev {
				t.Errorf("sort order violated: record[%d] Created_Time=%s > record[%d] Created_Time=%s",
					i, curr, i-1, prev)
			}
		}

		for _, rec := range arr {
			for key := range rec {
				switch key {
				case "id", "Last_Name", "Created_Time":
				default:
					t.Errorf("unexpected field %q in response with --fields id,Last_Name,Created_Time", key)
				}
			}
		}
	})

	t.Run("records/list-default-fields", func(t *testing.T) {
		out := zoho(t, "crm", "records", "list", "Leads", "--limit", "1")
		arr := parseJSONArray(t, out)
		assertNonEmpty(t, arr, "expected at least one lead")
		rec := arr[0]
		for _, want := range []string{"id", "Created_Time", "Modified_Time"} {
			if _, ok := rec[want]; !ok {
				t.Errorf("expected default field %q in response without --fields", want)
			}
		}
	})

	t.Run("records/list-all", func(t *testing.T) {
		outAll := zoho(t, "crm", "records", "list", "Leads",
			"--fields", "id", "--all")
		all := parseJSONArray(t, outAll)
		assertNonEmpty(t, all, "expected at least one lead with --all")

		outOne := zoho(t, "crm", "records", "list", "Leads",
			"--fields", "id", "--limit", "1")
		page1 := parseJSONArray(t, outOne)
		if len(page1) != 1 {
			t.Fatalf("expected 1 record with --limit 1, got %d", len(page1))
		}
		if len(all) < len(page1) {
			t.Errorf("--all should return at least as many as --limit 1: all=%d page1=%d",
				len(all), len(page1))
		}
	})

	t.Run("records/update", func(t *testing.T) {
		requireID(t, leadID, "create must have succeeded")

		before := getRecord(t, "Leads", leadID, "id,Last_Name,Company,Email")
		assertEqual(t, before["Company"], "TestCorp")
		assertEqual(t, before["Last_Name"], leadName)

		data := `{"Company":"UpdatedCorp"}`
		out := zoho(t, "crm", "records", "update", "Leads", leadID, "--json", data)
		assertStatus(t, out, "success")

		after := getRecord(t, "Leads", leadID, "id,Last_Name,Company,Email")
		assertStringField(t, after, "Company", "UpdatedCorp")
		assertStringField(t, after, "Last_Name", leadName)
		assertStringField(t, after, "Email", leadEmail)
	})

	t.Run("records/search-by-criteria", func(t *testing.T) {
		requireID(t, leadID, "create must have succeeded")
		retryUntil(t, 30*time.Second, func() bool {
			out, err := zohoMayFail(t, "crm", "records", "search", "Leads",
				"--criteria", fmt.Sprintf("(Last_Name:equals:%s)", leadName),
				"--fields", "id,Last_Name,Company")
			if err != nil {
				return false
			}
			var arr []map[string]any
			json.Unmarshal([]byte(out), &arr)
			if len(arr) == 0 {
				return false
			}
			rec, found := findInArray(arr, leadID)
			if !found {
				return false
			}
			assertEqual(t, rec["Last_Name"], leadName)
			assertEqual(t, rec["Company"], "UpdatedCorp")
			return true
		})
	})

	t.Run("records/search-by-word", func(t *testing.T) {
		requireID(t, leadID, "create must have succeeded")
		retryUntil(t, 30*time.Second, func() bool {
			out, err := zohoMayFail(t, "crm", "records", "search", "Leads",
				"--word", leadName, "--fields", "id,Last_Name")
			if err != nil {
				return false
			}
			var arr []map[string]any
			json.Unmarshal([]byte(out), &arr)
			if len(arr) == 0 {
				return false
			}
			rec, found := findInArray(arr, leadID)
			if !found {
				return false
			}
			return fmt.Sprintf("%v", rec["Last_Name"]) == leadName
		})
	})

	t.Run("records/search-by-email", func(t *testing.T) {
		requireID(t, leadID, "create must have succeeded")
		retryUntil(t, 30*time.Second, func() bool {
			out, err := zohoMayFail(t, "crm", "records", "search", "Leads",
				"--email", leadEmail, "--fields", "id,Email")
			if err != nil {
				return false
			}
			var arr []map[string]any
			json.Unmarshal([]byte(out), &arr)
			if len(arr) == 0 {
				return false
			}
			rec, found := findInArray(arr, leadID)
			if !found {
				return false
			}
			return fmt.Sprintf("%v", rec["Email"]) == leadEmail
		})
	})

	t.Run("records/search-by-phone", func(t *testing.T) {
		requireID(t, leadID, "create must have succeeded")
		retryUntil(t, 30*time.Second, func() bool {
			out, err := zohoMayFail(t, "crm", "records", "search", "Leads",
				"--phone", leadPhone, "--fields", "id,Phone")
			if err != nil {
				return false
			}
			var arr []map[string]any
			json.Unmarshal([]byte(out), &arr)
			if len(arr) == 0 {
				return false
			}
			_, found := findInArray(arr, leadID)
			return found
		})
	})

	t.Run("records/upsert-insert", func(t *testing.T) {
		upsertName = testName(t)
		upsertEmail = strings.ToLower(upsertName) + "@test.example.com"
		data := toJSON(t, map[string]any{"Last_Name": upsertName, "Company": "UpsertCorp", "Email": upsertEmail})
		out := zoho(t, "crm", "records", "upsert", "Leads", "--json", data, "--duplicate-check", "Email")
		upsertLeadID = extractID(t, out)
		cleanup.trackLead(upsertLeadID)
		assertAction(t, out, "insert")

		rec := getRecord(t, "Leads", upsertLeadID, "id,Last_Name,Company,Email")
		assertStringField(t, rec, "Last_Name", upsertName)
		assertStringField(t, rec, "Company", "UpsertCorp")
		assertStringField(t, rec, "Email", upsertEmail)
	})

	t.Run("records/upsert-update", func(t *testing.T) {
		requireID(t, upsertLeadID, "upsert-insert must have succeeded")

		before := getRecord(t, "Leads", upsertLeadID, "id,Last_Name,Company,Email")
		assertEqual(t, before["Company"], "UpsertCorp")

		data := toJSON(t, map[string]any{"Last_Name": "UpdatedViaUpsert", "Company": "UpsertCorpV2", "Email": upsertEmail})
		out := zoho(t, "crm", "records", "upsert", "Leads", "--json", data, "--duplicate-check", "Email")
		assertAction(t, out, "update")

		after := getRecord(t, "Leads", upsertLeadID, "id,Last_Name,Company,Email")
		assertEqual(t, fmt.Sprintf("%v", after["id"]), upsertLeadID)
		assertStringField(t, after, "Last_Name", "UpdatedViaUpsert")
		assertStringField(t, after, "Company", "UpsertCorpV2")
		assertStringField(t, after, "Email", upsertEmail)
	})

	t.Run("notes/add", func(t *testing.T) {
		requireID(t, leadID, "create must have succeeded")
		out := zoho(t, "crm", "notes", "add", "Leads", leadID,
			"--content", "Integration test note content",
			"--title", "Test Note")
		noteID = extractID(t, out)
		cleanup.trackNote(noteID)

		notes := getNotes(t, "Leads", leadID)
		note, found := findInArray(notes, noteID)
		if !found {
			t.Fatalf("note %s not found on Zoho after add", noteID)
		}
		assertStringField(t, note, "Note_Title", "Test Note")
		assertStringField(t, note, "Note_Content", "Integration test note content")
	})

	t.Run("notes/list", func(t *testing.T) {
		requireID(t, leadID, "create must have succeeded")
		requireID(t, noteID, "notes/add must have succeeded")
		notes := getNotes(t, "Leads", leadID)
		assertNonEmpty(t, notes, "expected at least one note")

		note, found := findInArray(notes, noteID)
		if !found {
			t.Fatalf("note %s not found in list", noteID)
		}
		assertStringField(t, note, "Note_Title", "Test Note")
		assertStringField(t, note, "Note_Content", "Integration test note content")
	})

	t.Run("notes/update", func(t *testing.T) {
		requireID(t, noteID, "notes/add must have succeeded")
		out := zoho(t, "crm", "notes", "update", noteID,
			"--content", "Updated note content",
			"--title", "Updated Note Title")
		assertStatus(t, out, "success")

		notes := getNotes(t, "Leads", leadID)
		note, found := findInArray(notes, noteID)
		if !found {
			t.Fatalf("note %s not found on Zoho after update", noteID)
		}
		assertStringField(t, note, "Note_Title", "Updated Note Title")
		assertStringField(t, note, "Note_Content", "Updated note content")
	})

	t.Run("related/list", func(t *testing.T) {
		requireID(t, leadID, "create must have succeeded")
		requireID(t, noteID, "notes/add must have succeeded")
		out := zoho(t, "crm", "related", "list", "Leads", leadID, "Notes",
			"--fields", "id,Note_Title")
		arr := parseJSONArray(t, out)
		assertNonEmpty(t, arr, "expected at least one related note")

		_, found := findInArray(arr, noteID)
		if !found {
			t.Errorf("note %s not found in related list for lead %s", noteID, leadID)
		}
	})

	t.Run("notes/delete", func(t *testing.T) {
		requireID(t, noteID, "notes/add must have succeeded")
		out := zoho(t, "crm", "notes", "delete", noteID)
		assertStatus(t, out, "success")

		notes := getNotes(t, "Leads", leadID)
		_, found := findInArray(notes, noteID)
		if found {
			t.Errorf("note %s still found on Zoho after delete", noteID)
		}
		noteID = ""
	})

	t.Run("attachments/upload", func(t *testing.T) {
		requireID(t, leadID, "create must have succeeded")
		tmpDir := t.TempDir()
		testFile := tmpDir + "/test-attachment.txt"
		testFileContent = []byte("zoho-cli integration test file " + time.Now().String())
		if err := os.WriteFile(testFile, testFileContent, 0644); err != nil {
			t.Fatalf("failed to create test file: %v", err)
		}
		out := zoho(t, "crm", "attachments", "upload", "Leads", leadID, testFile)
		attachmentID = extractID(t, out)
		cleanup.trackAttachment("Leads", leadID, attachmentID)

		atts := getAttachments(t, "Leads", leadID)
		att, found := findInArray(atts, attachmentID)
		if !found {
			t.Fatalf("attachment %s not found on Zoho after upload", attachmentID)
		}
		assertStringField(t, att, "File_Name", "test-attachment.txt")
	})

	t.Run("attachments/list", func(t *testing.T) {
		requireID(t, leadID, "create must have succeeded")
		requireID(t, attachmentID, "attachments/upload must have succeeded")
		atts := getAttachments(t, "Leads", leadID)
		assertNonEmpty(t, atts, "expected at least one attachment")

		att, found := findInArray(atts, attachmentID)
		if !found {
			t.Fatalf("attachment %s not found in list", attachmentID)
		}
		assertStringField(t, att, "File_Name", "test-attachment.txt")
	})

	t.Run("attachments/download", func(t *testing.T) {
		requireID(t, leadID, "create must have succeeded")
		requireID(t, attachmentID, "attachments/upload must have succeeded")
		tmpDir := t.TempDir()
		downloadPath := tmpDir + "/downloaded.txt"
		zoho(t, "crm", "attachments", "download", "Leads", leadID, attachmentID,
			"--output", downloadPath)
		downloaded, err := os.ReadFile(downloadPath)
		if err != nil {
			t.Fatalf("failed to read downloaded file: %v", err)
		}
		if !bytes.Equal(downloaded, testFileContent) {
			t.Errorf("downloaded content doesn't match: got %d bytes, want %d bytes",
				len(downloaded), len(testFileContent))
		}
	})

	t.Run("attachments/download-stdout", func(t *testing.T) {
		requireID(t, leadID, "create must have succeeded")
		requireID(t, attachmentID, "attachments/upload must have succeeded")
		r := runZoho(t, "crm", "attachments", "download", "Leads", leadID, attachmentID)
		if r.ExitCode != 0 {
			t.Fatalf("download to stdout failed (exit %d): %s", r.ExitCode, r.Stderr)
		}
		if !bytes.Equal([]byte(r.Stdout), testFileContent) {
			t.Errorf("stdout content doesn't match: got %d bytes, want %d bytes",
				len(r.Stdout), len(testFileContent))
		}
	})

	t.Run("attachments/special-filename", func(t *testing.T) {
		requireID(t, leadID, "create must have succeeded")
		tmpDir := t.TempDir()
		testFile := tmpDir + "/test file (1).txt"
		content := []byte("special filename test")
		if err := os.WriteFile(testFile, content, 0644); err != nil {
			t.Fatalf("failed to create test file: %v", err)
		}
		out := zoho(t, "crm", "attachments", "upload", "Leads", leadID, testFile)
		attID := extractID(t, out)

		atts := getAttachments(t, "Leads", leadID)
		att, found := findInArray(atts, attID)
		if !found {
			t.Fatalf("attachment %s not found after upload", attID)
		}
		assertStringField(t, att, "File_Name", "test file (1).txt")

		delOut := zoho(t, "crm", "attachments", "delete", "Leads", leadID, attID)
		assertStatus(t, delOut, "success")
	})

	t.Run("attachments/delete", func(t *testing.T) {
		requireID(t, leadID, "create must have succeeded")
		requireID(t, attachmentID, "attachments/upload must have succeeded")
		out := zoho(t, "crm", "attachments", "delete", "Leads", leadID, attachmentID)
		assertStatus(t, out, "success")

		atts := getAttachments(t, "Leads", leadID)
		_, found := findInArray(atts, attachmentID)
		if found {
			t.Errorf("attachment %s still found on Zoho after delete", attachmentID)
		}
		attachmentID = ""
	})

	t.Run("tags/add-multiple", func(t *testing.T) {
		requireID(t, leadID, "create must have succeeded")
		out := zoho(t, "crm", "tags", "add", "Leads",
			"--ids", leadID, "--tags", "zohotest-tag-a,zohotest-tag-b")
		assertStatus(t, out, "success")

		rec := getRecord(t, "Leads", leadID, "id,Tag")
		if !hasTag(rec, "zohotest-tag-a") {
			t.Errorf("tag 'zohotest-tag-a' not found after add; got: %v", rec["Tag"])
		}
		if !hasTag(rec, "zohotest-tag-b") {
			t.Errorf("tag 'zohotest-tag-b' not found after add; got: %v", rec["Tag"])
		}
	})

	t.Run("tags/remove-one-of-two", func(t *testing.T) {
		requireID(t, leadID, "create must have succeeded")
		out := zoho(t, "crm", "tags", "remove", "Leads",
			"--ids", leadID, "--tags", "zohotest-tag-a")
		assertStatus(t, out, "success")

		rec := getRecord(t, "Leads", leadID, "id,Tag")
		if hasTag(rec, "zohotest-tag-a") {
			t.Errorf("tag 'zohotest-tag-a' still present after remove")
		}
		if !hasTag(rec, "zohotest-tag-b") {
			t.Errorf("tag 'zohotest-tag-b' should still be present; got: %v", rec["Tag"])
		}
	})

	t.Run("tags/remove-remaining", func(t *testing.T) {
		requireID(t, leadID, "create must have succeeded")
		out := zoho(t, "crm", "tags", "remove", "Leads",
			"--ids", leadID, "--tags", "zohotest-tag-b")
		assertStatus(t, out, "success")

		rec := getRecord(t, "Leads", leadID, "id,Tag")
		if hasTag(rec, "zohotest-tag-b") {
			t.Errorf("tag 'zohotest-tag-b' still present after remove")
		}
	})

	t.Run("tags/add-multi-records", func(t *testing.T) {
		name1 := testName(t)
		data1 := toJSON(t, map[string]any{"Last_Name": name1, "Company": "TagCorp1"})
		out1 := zoho(t, "crm", "records", "create", "Leads", "--json", data1)
		id1 := extractID(t, out1)
		cleanup.trackLead(id1)

		name2 := testName(t)
		data2 := toJSON(t, map[string]any{"Last_Name": name2, "Company": "TagCorp2"})
		out2 := zoho(t, "crm", "records", "create", "Leads", "--json", data2)
		id2 := extractID(t, out2)
		cleanup.trackLead(id2)

		ids := id1 + "," + id2
		out := zoho(t, "crm", "tags", "add", "Leads",
			"--ids", ids, "--tags", "zohotest-multi-tag")
		assertStatus(t, out, "success")

		rec1 := getRecord(t, "Leads", id1, "id,Tag")
		if !hasTag(rec1, "zohotest-multi-tag") {
			t.Errorf("tag not found on record %s", id1)
		}
		rec2 := getRecord(t, "Leads", id2, "id,Tag")
		if !hasTag(rec2, "zohotest-multi-tag") {
			t.Errorf("tag not found on record %s", id2)
		}

		zoho(t, "crm", "tags", "remove", "Leads",
			"--ids", ids, "--tags", "zohotest-multi-tag")
	})

	t.Run("coql", func(t *testing.T) {
		requireID(t, leadID, "create must have succeeded")
		query := fmt.Sprintf("select id, Last_Name, Company from Leads where id = '%s'", leadID)
		out := zoho(t, "crm", "coql", "--query", query)
		parsed := parseJSON(t, out)
		data, ok := parsed["data"].([]any)
		if !ok || len(data) == 0 {
			t.Fatalf("COQL returned no results:\n%s", truncate(out, 500))
		}
		rec, _ := data[0].(map[string]any)
		assertEqual(t, fmt.Sprintf("%v", rec["id"]), leadID)
		assertEqual(t, rec["Last_Name"], leadName)
		assertEqual(t, rec["Company"], "UpdatedCorp")
	})

	t.Run("coql/order-by", func(t *testing.T) {
		requireID(t, leadID, "create must have succeeded")
		requireID(t, upsertLeadID, "upsert-insert must have succeeded")
		query := fmt.Sprintf(
			"select id, Last_Name from Leads where id in ('%s','%s') order by Last_Name asc",
			leadID, upsertLeadID)
		out := zoho(t, "crm", "coql", "--query", query)
		parsed := parseJSON(t, out)
		data, ok := parsed["data"].([]any)
		if !ok || len(data) < 2 {
			t.Fatalf("COQL ORDER BY returned fewer than 2 results:\n%s", truncate(out, 500))
		}
		first, _ := data[0].(map[string]any)
		second, _ := data[1].(map[string]any)
		name1 := fmt.Sprintf("%v", first["Last_Name"])
		name2 := fmt.Sprintf("%v", second["Last_Name"])
		if name1 > name2 {
			t.Errorf("ORDER BY asc violated: %q > %q", name1, name2)
		}
	})

	t.Run("coql/limit", func(t *testing.T) {
		query := fmt.Sprintf("select id from Leads where Last_Name like '%s%%' limit 2", testPrefix)
		out := zoho(t, "crm", "coql", "--query", query)
		parsed := parseJSON(t, out)
		data, ok := parsed["data"].([]any)
		if !ok {
			t.Fatalf("COQL LIMIT returned no data:\n%s", truncate(out, 500))
		}
		if len(data) > 2 {
			t.Errorf("COQL LIMIT 2 returned %d records, expected at most 2", len(data))
		}
	})

	t.Run("coql/like-operator", func(t *testing.T) {
		requireID(t, leadID, "create must have succeeded")
		query := fmt.Sprintf("select id, Last_Name from Leads where Last_Name like '%s%%'", testPrefix)
		out := zoho(t, "crm", "coql", "--query", query)
		parsed := parseJSON(t, out)
		data, ok := parsed["data"].([]any)
		if !ok || len(data) == 0 {
			t.Fatalf("COQL LIKE returned no results:\n%s", truncate(out, 500))
		}
		for _, item := range data {
			rec, _ := item.(map[string]any)
			name := fmt.Sprintf("%v", rec["Last_Name"])
			if !strings.HasPrefix(name, testPrefix) {
				t.Errorf("LIKE '%s%%' returned record with Last_Name=%q", testPrefix, name)
			}
		}
	})

	t.Run("coql/multi-field-types", func(t *testing.T) {
		requireID(t, leadID, "create must have succeeded")
		query := fmt.Sprintf(
			"select id, Last_Name, Company, Email, Created_Time from Leads where id = '%s'", leadID)
		out := zoho(t, "crm", "coql", "--query", query)
		parsed := parseJSON(t, out)
		data, ok := parsed["data"].([]any)
		if !ok || len(data) == 0 {
			t.Fatalf("COQL multi-field returned no data:\n%s", truncate(out, 500))
		}
		rec, _ := data[0].(map[string]any)
		for _, field := range []string{"id", "Last_Name", "Company", "Email", "Created_Time"} {
			if _, ok := rec[field]; !ok {
				t.Errorf("expected field %q in COQL result", field)
			}
		}
		if _, ok := rec["id"].(string); !ok {
			t.Errorf("id should be string, got %T", rec["id"])
		}
		if _, ok := rec["Created_Time"].(string); !ok {
			t.Errorf("Created_Time should be string, got %T", rec["Created_Time"])
		}
	})

	t.Run("search/global", func(t *testing.T) {
		requireID(t, leadID, "create must have succeeded")
		retryUntil(t, 30*time.Second, func() bool {
			out, err := zohoMayFail(t, "crm", "search-global", leadName)
			if err != nil {
				return false
			}
			var envelope struct {
				Data []map[string]any `json:"data"`
			}
			if jsonErr := json.Unmarshal([]byte(out), &envelope); jsonErr != nil {
				return false
			}
			for _, r := range envelope.Data {
				if fmt.Sprintf("%v", r["id"]) != leadID {
					continue
				}
				if setype, ok := r["$setype"].(string); !ok || setype != "Leads" {
					t.Errorf("search-global $setype: got %q, want %q", setype, "Leads")
				}
				return true
			}
			return false
		})
	})

	t.Run("records/bulk-delete", func(t *testing.T) {
		name1 := testName(t)
		data1 := toJSON(t, map[string]any{"Last_Name": name1, "Company": "BulkCorp1"})
		out1 := zoho(t, "crm", "records", "create", "Leads", "--json", data1)
		id1 := extractID(t, out1)

		name2 := testName(t)
		data2 := toJSON(t, map[string]any{"Last_Name": name2, "Company": "BulkCorp2"})
		out2 := zoho(t, "crm", "records", "create", "Leads", "--json", data2)
		id2 := extractID(t, out2)

		cleanup.trackLead(id1)
		cleanup.trackLead(id2)

		rec1 := getRecord(t, "Leads", id1, "id,Last_Name")
		assertEqual(t, rec1["Last_Name"], name1)
		rec2 := getRecord(t, "Leads", id2, "id,Last_Name")
		assertEqual(t, rec2["Last_Name"], name2)

		ids := id1 + "," + id2
		out := zoho(t, "crm", "records", "bulk-delete", "Leads", ids)
		assertStatus(t, out, "success")
		var bulkResp struct {
			Data []struct {
				Status string `json:"status"`
			} `json:"data"`
		}
		if err := json.Unmarshal([]byte(out), &bulkResp); err != nil {
			t.Fatalf("failed to parse bulk-delete response: %v", err)
		}
		if len(bulkResp.Data) != 2 {
			t.Errorf("expected 2 results in bulk-delete, got %d", len(bulkResp.Data))
		}
		for i, d := range bulkResp.Data {
			if d.Status != "success" {
				t.Errorf("bulk-delete item %d: expected success, got %q", i, d.Status)
			}
		}

		_, err1 := getRecordMayFail(t, "Leads", id1)
		if err1 == nil {
			t.Errorf("lead %s still exists on Zoho after bulk delete", id1)
		}
		_, err2 := getRecordMayFail(t, "Leads", id2)
		if err2 == nil {
			t.Errorf("lead %s still exists on Zoho after bulk delete", id2)
		}
	})

	t.Run("records/delete", func(t *testing.T) {
		requireID(t, leadID, "create must have succeeded")
		out := zoho(t, "crm", "records", "delete", "Leads", leadID)
		assertStatus(t, out, "success")

		_, err := getRecordMayFail(t, "Leads", leadID)
		if err == nil {
			t.Errorf("lead %s still accessible on Zoho after delete", leadID)
		}
		leadID = ""
	})

}

func TestCRMConvert(t *testing.T) {
	t.Parallel()
	cleanup := newCleanup(t)

	leadName := testName(t)
	leadCompany := testPrefix + "ConvertCorp_" + randomSuffix()
	leadEmail := strings.ToLower(leadName) + "@test.example.com"

	data := toJSON(t, map[string]any{"Last_Name": leadName, "Company": leadCompany, "Email": leadEmail})
	createOut := zoho(t, "crm", "records", "create", "Leads", "--json", data)
	leadID := extractID(t, createOut)
	cleanup.trackLead(leadID)

	convertOut := zoho(t, "crm", "convert", leadID)
	ids := extractConvertIDs(t, convertOut)

	cleanup.trackContact(ids.ContactID)
	cleanup.trackAccount(ids.AccountID)

	t.Run("contact-exists", func(t *testing.T) {
		retryUntil(t, 15*time.Second, func() bool {
			rec, err := getRecordMayFail(t, "Contacts", ids.ContactID)
			if err != nil {
				return false
			}
			return fmt.Sprintf("%v", rec["id"]) == ids.ContactID
		})
		rec := getRecord(t, "Contacts", ids.ContactID, "id,Last_Name,Email")
		assertStringField(t, rec, "Last_Name", leadName)
		assertStringField(t, rec, "Email", leadEmail)
	})

	t.Run("account-exists", func(t *testing.T) {
		retryUntil(t, 15*time.Second, func() bool {
			rec, err := getRecordMayFail(t, "Accounts", ids.AccountID)
			if err != nil {
				return false
			}
			return fmt.Sprintf("%v", rec["id"]) == ids.AccountID
		})
		rec := getRecord(t, "Accounts", ids.AccountID, "id,Account_Name")
		assertStringField(t, rec, "Account_Name", leadCompany)
	})

	t.Run("lead-gone", func(t *testing.T) {
		retryUntil(t, 15*time.Second, func() bool {
			_, err := getRecordMayFail(t, "Leads", leadID)
			return err != nil
		})
	})

	t.Run("no-deal-created", func(t *testing.T) {
		if ids.DealID != "" {
			t.Errorf("expected no deal from simple conversion, got deal ID %s", ids.DealID)
		}
	})
}

func TestCRMErrors(t *testing.T) {
	t.Parallel()
	t.Run("bad-auth", func(t *testing.T) {
		r := runZohoWithEnv(t, map[string]string{
			"ZOHO_CLIENT_ID":     "bad_client_id",
			"ZOHO_CLIENT_SECRET": "bad_client_secret",
			"ZOHO_REFRESH_TOKEN": "bad_refresh_token",
			"ZOHO_DC":            "com",
		}, "crm", "records", "list", "Leads", "--fields", "id")
		assertExitCode(t, r, 2)
		if !strings.Contains(r.Stderr, "invalid_client") && !strings.Contains(r.Stderr, "Token refresh") {
			t.Errorf("expected auth error in stderr, got: %s", truncate(r.Stderr, 500))
		}
	})

	t.Run("invalid-module", func(t *testing.T) {
		r := runZoho(t, "crm", "records", "list", "FakeModule", "--fields", "id")
		assertExitCode(t, r, 1)
		assertContains(t, r.Stderr, "INVALID_MODULE")
	})

	t.Run("invalid-json", func(t *testing.T) {
		r := runZoho(t, "crm", "records", "create", "Leads", "--json", "not json")
		assertExitCode(t, r, 1)
		assertContains(t, r.Stderr, "INVALID_DATA")
	})

	t.Run("invalid-coql", func(t *testing.T) {
		r := runZoho(t, "crm", "coql", "--query", "select broken")
		assertExitCode(t, r, 1)
		assertContains(t, r.Stderr, "SYNTAX_ERROR")
	})

	t.Run("missing-required-flag", func(t *testing.T) {
		r := runZoho(t, "crm", "records", "create", "Leads")
		assertExitCode(t, r, 1)
		assertContains(t, r.Stderr, "Required flag")
	})

	t.Run("nonexistent-record", func(t *testing.T) {
		_, err := getRecordMayFail(t, "Leads", "999999999999999999")
		if err == nil {
			t.Error("expected error for nonexistent record")
		}
	})

	t.Run("invalid-coql-no-from", func(t *testing.T) {
		r := runZoho(t, "crm", "coql", "--query", "select id")
		assertExitCode(t, r, 1)
	})

	t.Run("invalid-coql-bad-field", func(t *testing.T) {
		r := runZoho(t, "crm", "coql", "--query",
			"select Nonexistent_Field_XYZ from Leads limit 1")
		assertExitCode(t, r, 1)
	})
}

func TestCRMEmergencyCleanup(t *testing.T) {
	if os.Getenv("ZOHO_EMERGENCY_CLEANUP") == "" {
		t.Skip("set ZOHO_EMERGENCY_CLEANUP=1 to run")
	}
	out := zoho(t, "crm", "coql", "--query",
		"select id from Leads where Last_Name like 'ZOHOTEST%'")
	var resp struct {
		Data []struct {
			ID string `json:"id"`
		} `json:"data"`
	}
	if err := json.Unmarshal([]byte(out), &resp); err != nil {
		t.Fatalf("failed to parse COQL response: %v", err)
	}
	t.Logf("found %d orphaned test leads", len(resp.Data))
	for _, rec := range resp.Data {
		t.Logf("deleting orphaned lead %s", rec.ID)
		zohoIgnoreError(t, "crm", "records", "delete", "Leads", rec.ID)
	}
}
