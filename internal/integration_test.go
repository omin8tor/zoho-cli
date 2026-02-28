//go:build integration

package internal_test

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"testing"
	"time"
)

const testPrefix = "ZOHOTEST"

type Result struct {
	Stdout   string
	Stderr   string
	ExitCode int
}

func runZoho(t *testing.T, args ...string) Result {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, "./zoho", args...)
	cmd.Dir = ".."
	cmd.Env = os.Environ()
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	exitCode := 0
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			exitCode = exitErr.ExitCode()
		} else {
			exitCode = -1
		}
	}
	return Result{
		Stdout:   stdout.String(),
		Stderr:   stderr.String(),
		ExitCode: exitCode,
	}
}

func zoho(t *testing.T, args ...string) string {
	t.Helper()
	r := runZoho(t, args...)
	if r.ExitCode != 0 {
		t.Fatalf("zoho %s failed (exit %d):\nstdout: %s\nstderr: %s",
			strings.Join(args, " "), r.ExitCode, r.Stdout, r.Stderr)
	}
	return r.Stdout
}

func zohoMayFail(t *testing.T, args ...string) (string, error) {
	t.Helper()
	r := runZoho(t, args...)
	if r.ExitCode != 0 {
		return r.Stdout, fmt.Errorf("exit %d: %s", r.ExitCode, r.Stderr)
	}
	return r.Stdout, nil
}

func zohoIgnoreError(t *testing.T, args ...string) string {
	t.Helper()
	r := runZoho(t, args...)
	if r.ExitCode != 0 {
		t.Logf("zoho %s failed (ignored, exit %d): %s",
			strings.Join(args, " "), r.ExitCode, r.Stderr)
	}
	return r.Stdout
}

func parseJSON(t *testing.T, s string) map[string]any {
	t.Helper()
	var m map[string]any
	if err := json.Unmarshal([]byte(s), &m); err != nil {
		t.Fatalf("failed to parse JSON object: %v\nraw: %s", err, truncate(s, 500))
	}
	return m
}

func parseJSONArray(t *testing.T, s string) []map[string]any {
	t.Helper()
	var arr []map[string]any
	if err := json.Unmarshal([]byte(s), &arr); err != nil {
		t.Fatalf("failed to parse JSON array: %v\nraw: %s", err, truncate(s, 500))
	}
	return arr
}

func extractID(t *testing.T, out string) string {
	t.Helper()
	var resp struct {
		Data []struct {
			Details struct {
				ID string `json:"id"`
			} `json:"details"`
			Status string `json:"status"`
		} `json:"data"`
	}
	if err := json.Unmarshal([]byte(out), &resp); err != nil {
		t.Fatalf("failed to parse response for ID extraction: %v\nraw: %s", err, truncate(out, 500))
	}
	if len(resp.Data) == 0 {
		t.Fatalf("no data in response:\n%s", truncate(out, 500))
	}
	if resp.Data[0].Status != "success" {
		t.Fatalf("operation was not successful:\n%s", truncate(out, 500))
	}
	id := resp.Data[0].Details.ID
	if id == "" {
		t.Fatalf("empty ID in response:\n%s", truncate(out, 500))
	}
	return id
}

func requireID(t *testing.T, id string, msg string) {
	t.Helper()
	if id == "" {
		t.Skipf("skipping: %s (no ID available)", msg)
	}
}

func retryUntil(t *testing.T, timeout time.Duration, fn func() bool) {
	t.Helper()
	deadline := time.Now().Add(timeout)
	interval := 2 * time.Second
	for time.Now().Before(deadline) {
		if fn() {
			return
		}
		t.Logf("retrying in %v...", interval)
		time.Sleep(interval)
	}
	t.Error("timed out waiting for condition")
}

func assertEqual(t *testing.T, got any, want any) {
	t.Helper()
	gotStr := fmt.Sprintf("%v", got)
	wantStr := fmt.Sprintf("%v", want)
	if gotStr != wantStr {
		t.Errorf("got %q, want %q", gotStr, wantStr)
	}
}

func assertContains(t *testing.T, s, substr string) {
	t.Helper()
	if !strings.Contains(s, substr) {
		t.Errorf("expected output to contain %q, got:\n%s", substr, truncate(s, 500))
	}
}

func assertNonEmpty(t *testing.T, arr []map[string]any, msg string) {
	t.Helper()
	if len(arr) == 0 {
		t.Error(msg)
	}
}

func assertStatus(t *testing.T, out string, want string) {
	t.Helper()
	var resp struct {
		Data []struct {
			Status string `json:"status"`
		} `json:"data"`
	}
	json.Unmarshal([]byte(out), &resp)
	if len(resp.Data) == 0 || resp.Data[0].Status != want {
		t.Errorf("expected status %q in response:\n%s", want, truncate(out, 500))
	}
}

func assertAction(t *testing.T, out string, want string) {
	t.Helper()
	var resp struct {
		Data []struct {
			Action string `json:"action"`
		} `json:"data"`
	}
	json.Unmarshal([]byte(out), &resp)
	if len(resp.Data) == 0 || resp.Data[0].Action != want {
		t.Errorf("expected action %q in response:\n%s", want, truncate(out, 500))
	}
}

func truncate(s string, n int) string {
	if len(s) > n {
		return s[:n] + "... (truncated)"
	}
	return s
}

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

type cleanupEntry struct {
	label string
	fn    func()
}

type testCleanup struct {
	t       *testing.T
	entries []cleanupEntry
}

func newCleanup(t *testing.T) *testCleanup {
	c := &testCleanup{t: t}
	t.Cleanup(func() {
		for i := len(c.entries) - 1; i >= 0; i-- {
			entry := c.entries[i]
			c.t.Logf("cleanup: %s", entry.label)
			entry.fn()
		}
	})
	return c
}

func (c *testCleanup) add(label string, fn func()) {
	c.entries = append(c.entries, cleanupEntry{label: label, fn: fn})
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

func randomSuffix(n int) string {
	const letters = "abcdefghijklmnopqrstuvwxyz0123456789"
	b := make([]byte, n)
	rand.Read(b)
	for i := range b {
		b[i] = letters[b[i]%byte(len(letters))]
	}
	return string(b)
}

func testName(t *testing.T) string {
	t.Helper()
	return fmt.Sprintf("%s_%d_%s", testPrefix, time.Now().Unix(), randomSuffix(6))
}

func TestCRMModules(t *testing.T) {
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
	})

	t.Run("related-lists", func(t *testing.T) {
		out := zoho(t, "crm", "modules", "related-lists", "Leads")
		arr := parseJSONArray(t, out)
		assertNonEmpty(t, arr, "expected related lists for Leads")
	})

	t.Run("layouts", func(t *testing.T) {
		out := zoho(t, "crm", "modules", "layouts", "Leads")
		arr := parseJSONArray(t, out)
		assertNonEmpty(t, arr, "expected layouts for Leads")
	})

	t.Run("custom-views", func(t *testing.T) {
		out := zoho(t, "crm", "modules", "custom-views", "Leads")
		arr := parseJSONArray(t, out)
		assertNonEmpty(t, arr, "expected custom views for Leads")
	})
}

func TestCRM(t *testing.T) {
	cleanup := newCleanup(t)

	var leadID string
	var leadName string
	var leadEmail string
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

		data := fmt.Sprintf(`{"Last_Name":"%s","Company":"TestCorp","Email":"%s"}`,
			leadName, leadEmail)
		out := zoho(t, "crm", "records", "create", "Leads", "--json", data)
		leadID = extractID(t, out)
		cleanup.trackLead(leadID)

		rec := getRecord(t, "Leads", leadID, "id,Last_Name,Company,Email")
		assertEqual(t, fmt.Sprintf("%v", rec["id"]), leadID)
		assertEqual(t, rec["Last_Name"], leadName)
		assertEqual(t, rec["Company"], "TestCorp")
		assertEqual(t, rec["Email"], leadEmail)
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
			"--per-page", "5")
		arr := parseJSONArray(t, out)
		assertNonEmpty(t, arr, "expected at least one lead in list")

		if len(arr) > 5 {
			t.Errorf("--per-page 5 but got %d records", len(arr))
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

	t.Run("records/update", func(t *testing.T) {
		requireID(t, leadID, "create must have succeeded")

		before := getRecord(t, "Leads", leadID, "id,Last_Name,Company,Email")
		assertEqual(t, before["Company"], "TestCorp")
		assertEqual(t, before["Last_Name"], leadName)

		data := `{"Company":"UpdatedCorp"}`
		out := zoho(t, "crm", "records", "update", "Leads", leadID, "--json", data)
		assertStatus(t, out, "success")

		after := getRecord(t, "Leads", leadID, "id,Last_Name,Company,Email")
		assertEqual(t, after["Company"], "UpdatedCorp")
		assertEqual(t, after["Last_Name"], leadName)
		assertEqual(t, after["Email"], leadEmail)
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

	t.Run("records/upsert-insert", func(t *testing.T) {
		upsertName = testName(t)
		upsertEmail = strings.ToLower(upsertName) + "@test.example.com"
		data := fmt.Sprintf(`{"Last_Name":"%s","Company":"UpsertCorp","Email":"%s"}`,
			upsertName, upsertEmail)
		out := zoho(t, "crm", "records", "upsert", "Leads", "--json", data, "--duplicate-check", "Email")
		upsertLeadID = extractID(t, out)
		cleanup.trackLead(upsertLeadID)
		assertAction(t, out, "insert")

		rec := getRecord(t, "Leads", upsertLeadID, "id,Last_Name,Company,Email")
		assertEqual(t, rec["Last_Name"], upsertName)
		assertEqual(t, rec["Company"], "UpsertCorp")
		assertEqual(t, rec["Email"], upsertEmail)
	})

	t.Run("records/upsert-update", func(t *testing.T) {
		requireID(t, upsertLeadID, "upsert-insert must have succeeded")

		before := getRecord(t, "Leads", upsertLeadID, "id,Last_Name,Company,Email")
		assertEqual(t, before["Company"], "UpsertCorp")

		data := fmt.Sprintf(`{"Last_Name":"UpdatedViaUpsert","Company":"UpsertCorpV2","Email":"%s"}`,
			upsertEmail)
		out := zoho(t, "crm", "records", "upsert", "Leads", "--json", data, "--duplicate-check", "Email")
		assertAction(t, out, "update")

		after := getRecord(t, "Leads", upsertLeadID, "id,Last_Name,Company,Email")
		assertEqual(t, fmt.Sprintf("%v", after["id"]), upsertLeadID)
		assertEqual(t, after["Last_Name"], "UpdatedViaUpsert")
		assertEqual(t, after["Company"], "UpsertCorpV2")
		assertEqual(t, after["Email"], upsertEmail)
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
		assertEqual(t, note["Note_Title"], "Test Note")
		assertEqual(t, note["Note_Content"], "Integration test note content")
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
		assertEqual(t, note["Note_Title"], "Test Note")
		assertEqual(t, note["Note_Content"], "Integration test note content")
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
		assertEqual(t, note["Note_Title"], "Updated Note Title")
		assertEqual(t, note["Note_Content"], "Updated note content")
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
		assertEqual(t, att["File_Name"], "test-attachment.txt")
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
		assertEqual(t, att["File_Name"], "test-attachment.txt")
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

	t.Run("tags/add", func(t *testing.T) {
		requireID(t, leadID, "create must have succeeded")
		out := zoho(t, "crm", "tags", "add", "Leads", "--ids", leadID, "--tags", "zohotest-tag")
		assertStatus(t, out, "success")

		rec := getRecord(t, "Leads", leadID, "id,Tag")
		if !hasTag(rec, "zohotest-tag") {
			t.Errorf("tag 'zohotest-tag' not found on record %s after add; got: %v", leadID, rec["Tag"])
		}
	})

	t.Run("tags/remove", func(t *testing.T) {
		requireID(t, leadID, "create must have succeeded")
		out := zoho(t, "crm", "tags", "remove", "Leads", "--ids", leadID, "--tags", "zohotest-tag")
		assertStatus(t, out, "success")

		rec := getRecord(t, "Leads", leadID, "id,Tag")
		if hasTag(rec, "zohotest-tag") {
			t.Errorf("tag 'zohotest-tag' still on record %s after remove", leadID)
		}
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

	t.Run("search-global", func(t *testing.T) {
		requireID(t, leadID, "create must have succeeded")
		retryUntil(t, 30*time.Second, func() bool {
			out, err := zohoMayFail(t, "crm", "search-global", leadName)
			if err != nil {
				return false
			}
			return strings.Contains(out, leadName) && strings.Contains(out, leadID)
		})
	})

	t.Run("records/bulk-delete", func(t *testing.T) {
		name1 := testName(t)
		data1 := fmt.Sprintf(`{"Last_Name":"%s","Company":"BulkCorp1"}`, name1)
		out1 := zoho(t, "crm", "records", "create", "Leads", "--json", data1)
		id1 := extractID(t, out1)

		name2 := testName(t)
		data2 := fmt.Sprintf(`{"Last_Name":"%s","Company":"BulkCorp2"}`, name2)
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

func TestHelpAll(t *testing.T) {
	out := zoho(t, "--help-all")
	if !strings.Contains(out, "crm") {
		t.Error("expected crm in help-all output")
	}
	if !strings.Contains(out, "projects") {
		t.Error("expected projects in help-all output")
	}
	if !strings.Contains(out, "drive") {
		t.Error("expected drive in help-all output")
	}
}
