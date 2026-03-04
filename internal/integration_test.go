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
	return runZohoWithEnv(t, nil, args...)
}

func runZohoWithEnv(t *testing.T, env map[string]string, args ...string) Result {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, "./zoho", args...)
	cmd.Dir = ".."
	baseEnv := os.Environ()
	if len(env) > 0 {
		overridden := make(map[string]bool)
		for k := range env {
			overridden[k] = true
		}
		merged := make([]string, 0, len(baseEnv)+len(env))
		for _, e := range baseEnv {
			key := e
			if idx := strings.IndexByte(e, '='); idx >= 0 {
				key = e[:idx]
			}
			if overridden[key] {
				continue
			}
			merged = append(merged, e)
		}
		for k, v := range env {
			merged = append(merged, k+"="+v)
		}
		cmd.Env = merged
	} else {
		cmd.Env = baseEnv
	}
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

func assertStringField(t *testing.T, rec map[string]any, key, want string) {
	t.Helper()
	got, ok := rec[key]
	if !ok {
		t.Errorf("field %q: missing from record", key)
		return
	}
	gotStr, ok := got.(string)
	if !ok {
		t.Errorf("field %q: expected string, got %T (%v)", key, got, got)
		return
	}
	if gotStr != want {
		t.Errorf("field %q: got %q, want %q", key, gotStr, want)
	}
}

func assertExitCode(t *testing.T, r Result, want int) {
	t.Helper()
	if r.ExitCode != want {
		t.Errorf("expected exit code %d, got %d\nstderr: %s", want, r.ExitCode, truncate(r.Stderr, 500))
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
	if err := json.Unmarshal([]byte(out), &resp); err != nil {
		t.Fatalf("failed to parse response: %v\nraw: %s", err, truncate(out, 500))
	}
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
	if err := json.Unmarshal([]byte(out), &resp); err != nil {
		t.Fatalf("failed to parse response: %v\nraw: %s", err, truncate(out, 500))
	}
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

const driveTestParentFolder = "0any60e555791dc9f472fb1eadfe33100f228"

func driveAttr(t *testing.T, out string) map[string]any {
	t.Helper()
	m := parseJSON(t, out)
	data, ok := m["data"].(map[string]any)
	if !ok {
		t.Fatalf("expected data object in JSON:API response:\n%s", truncate(out, 500))
	}
	attrs, ok := data["attributes"].(map[string]any)
	if !ok {
		t.Fatalf("expected data.attributes in JSON:API response:\n%s", truncate(out, 500))
	}
	return attrs
}

func extractDriveID(t *testing.T, out string) string {
	t.Helper()
	m := parseJSON(t, out)
	data, ok := m["data"].(map[string]any)
	if !ok {
		t.Fatalf("expected data object in JSON:API response:\n%s", truncate(out, 500))
	}
	id, ok := data["id"].(string)
	if !ok || id == "" {
		t.Fatalf("expected non-empty data.id in JSON:API response:\n%s", truncate(out, 500))
	}
	return id
}

func extractDriveUploadID(t *testing.T, out string) string {
	t.Helper()
	m := parseJSON(t, out)
	data, ok := m["data"].([]any)
	if !ok || len(data) == 0 {
		t.Fatalf("expected data array in upload response:\n%s", truncate(out, 500))
	}
	item, ok := data[0].(map[string]any)
	if !ok {
		t.Fatalf("expected object in data[0]:\n%s", truncate(out, 500))
	}
	attrs, ok := item["attributes"].(map[string]any)
	if !ok {
		t.Fatalf("expected attributes in data[0]:\n%s", truncate(out, 500))
	}
	id, ok := attrs["resource_id"].(string)
	if !ok || id == "" {
		t.Fatalf("expected non-empty resource_id in upload response:\n%s", truncate(out, 500))
	}
	return id
}

func getDriveFile(t *testing.T, fileID string) map[string]any {
	t.Helper()
	out := zoho(t, "drive", "files", "get", fileID)
	return parseJSON(t, out)
}

func getDriveFileAttr(t *testing.T, fileID string) map[string]any {
	t.Helper()
	m := getDriveFile(t, fileID)
	data, ok := m["data"].(map[string]any)
	if !ok {
		t.Fatalf("expected data object in file get response")
	}
	attrs, ok := data["attributes"].(map[string]any)
	if !ok {
		t.Fatalf("expected data.attributes in file get response")
	}
	return attrs
}

func assertDriveAttr(t *testing.T, out string, key string, want string) {
	t.Helper()
	attrs := driveAttr(t, out)
	got := fmt.Sprintf("%v", attrs[key])
	if got != want {
		t.Errorf("drive attr %q: got %q, want %q", key, got, want)
	}
}

func (c *testCleanup) trackDriveFile(id string) {
	c.add("trash drive file "+id, func() {
		zohoIgnoreError(c.t, "drive", "files", "trash", id)
	})
}

func (c *testCleanup) trackDriveFolder(id string) {
	c.add("trash drive folder "+id, func() {
		zohoIgnoreError(c.t, "drive", "files", "trash", id)
	})
}

func requireDriveTeamID(t *testing.T) string {
	t.Helper()
	id := os.Getenv("ZOHO_TEAM_ID")
	if id == "" {
		t.Skip("skipping: ZOHO_TEAM_ID not set")
	}
	return id
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

		data := fmt.Sprintf(`{"Last_Name":"%s","Company":"TestCorp","Email":"%s","Phone":"%s"}`,
			leadName, leadEmail, leadPhone)
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

	t.Run("records/list-default-fields", func(t *testing.T) {
		out := zoho(t, "crm", "records", "list", "Leads", "--per-page", "1")
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
			"--fields", "id", "--per-page", "1")
		page1 := parseJSONArray(t, outOne)
		if len(page1) != 1 {
			t.Fatalf("expected 1 record with --per-page 1, got %d", len(page1))
		}
		if len(all) <= len(page1) {
			t.Errorf("--all should return more than --per-page 1: all=%d page1=%d",
				len(all), len(page1))
		}
	})

	t.Run("records/list-page", func(t *testing.T) {
		out1 := zoho(t, "crm", "records", "list", "Leads",
			"--fields", "id", "--page", "1", "--per-page", "1")
		page1 := parseJSONArray(t, out1)
		if len(page1) != 1 {
			t.Fatalf("expected 1 record on page 1, got %d", len(page1))
		}
		id1 := fmt.Sprintf("%v", page1[0]["id"])

		out2 := zoho(t, "crm", "records", "list", "Leads",
			"--fields", "id", "--page", "2", "--per-page", "1")
		page2 := parseJSONArray(t, out2)
		if len(page2) == 1 {
			id2 := fmt.Sprintf("%v", page2[0]["id"])
			if id2 == id1 {
				t.Errorf("page 1 and page 2 returned same record %s", id1)
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
		data := fmt.Sprintf(`{"Last_Name":"%s","Company":"UpsertCorp","Email":"%s"}`,
			upsertName, upsertEmail)
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

		data := fmt.Sprintf(`{"Last_Name":"UpdatedViaUpsert","Company":"UpsertCorpV2","Email":"%s"}`,
			upsertEmail)
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
		data1 := fmt.Sprintf(`{"Last_Name":"%s","Company":"TagCorp1"}`, name1)
		out1 := zoho(t, "crm", "records", "create", "Leads", "--json", data1)
		id1 := extractID(t, out1)
		cleanup.trackLead(id1)

		name2 := testName(t)
		data2 := fmt.Sprintf(`{"Last_Name":"%s","Company":"TagCorp2"}`, name2)
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

	t.Run("search-global", func(t *testing.T) {
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
	cleanup := newCleanup(t)

	leadName := testName(t)
	leadCompany := testPrefix + "ConvertCorp_" + randomSuffix(6)
	leadEmail := strings.ToLower(leadName) + "@test.example.com"

	data := fmt.Sprintf(`{"Last_Name":"%s","Company":"%s","Email":"%s"}`,
		leadName, leadCompany, leadEmail)
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
	t.Run("bad-auth", func(t *testing.T) {
		r := runZohoWithEnv(t, map[string]string{
			"ZOHO_CLIENT_ID":     "bad_client_id",
			"ZOHO_CLIENT_SECRET": "bad_client_secret",
			"ZOHO_REFRESH_TOKEN": "bad_refresh_token",
			"ZOHO_DC":            "com",
		}, "crm", "records", "list", "Leads", "--fields", "id", "--per-page", "1")
		assertExitCode(t, r, 2)
		if !strings.Contains(r.Stderr, "invalid_client") && !strings.Contains(r.Stderr, "Token refresh") {
			t.Errorf("expected auth error in stderr, got: %s", truncate(r.Stderr, 500))
		}
	})

	t.Run("invalid-module", func(t *testing.T) {
		r := runZoho(t, "crm", "records", "list", "FakeModule", "--fields", "id", "--per-page", "1")
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

func TestDriveTeams(t *testing.T) {
	teamID := requireDriveTeamID(t)
	var myUserID string

	t.Run("teams/me", func(t *testing.T) {
		out := zoho(t, "drive", "teams", "me")
		m := parseJSON(t, out)
		data, ok := m["data"].(map[string]any)
		if !ok {
			t.Fatalf("expected data object in response:\n%s", truncate(out, 500))
		}
		id, ok := data["id"].(string)
		if !ok || id == "" {
			t.Fatalf("expected non-empty data.id:\n%s", truncate(out, 500))
		}
		myUserID = id
		attrs, ok := data["attributes"].(map[string]any)
		if !ok {
			t.Fatalf("expected data.attributes:\n%s", truncate(out, 500))
		}
		if _, ok := attrs["email_id"].(string); !ok {
			t.Errorf("expected string email_id in attributes:\n%s", truncate(out, 500))
		}
		t.Logf("current user ID: %s", myUserID)
	})

	t.Run("teams/members", func(t *testing.T) {
		requireID(t, myUserID, "teams/me must have succeeded")
		out := zoho(t, "drive", "teams", "members", teamID)
		m := parseJSON(t, out)
		data, ok := m["data"].([]any)
		if !ok || len(data) == 0 {
			t.Fatalf("expected non-empty data array:\n%s", truncate(out, 500))
		}
		foundSelf := false
		for _, item := range data {
			member, ok := item.(map[string]any)
			if !ok {
				continue
			}
			if _, ok := member["id"].(string); !ok {
				t.Errorf("expected string id on member")
			}
			attrs, _ := member["attributes"].(map[string]any)
			if attrs != nil && fmt.Sprintf("%v", attrs["zuid"]) == myUserID {
				foundSelf = true
			}
		}
		if !foundSelf {
			t.Errorf("current user (zuid=%s) not found in team members list", myUserID)
		}
		t.Logf("found %d team members", len(data))
	})

	t.Run("folders/list", func(t *testing.T) {
		out := zoho(t, "drive", "folders", "list")
		arr := parseJSONArray(t, out)
		assertNonEmpty(t, arr, "expected at least one team folder")
		foundGeneral := false
		for _, item := range arr {
			attrs, ok := item["attributes"].(map[string]any)
			if !ok {
				continue
			}
			if fmt.Sprintf("%v", attrs["name"]) == "General" {
				foundGeneral = true
			}
		}
		if !foundGeneral {
			t.Error("expected 'General' team folder in list")
		}
	})
}

func TestDrive(t *testing.T) {
	teamID := requireDriveTeamID(t)
	cleanup := newCleanup(t)

	parentFolderID := os.Getenv("ZOHO_DRIVE_PARENT_FOLDER")
	if parentFolderID == "" {
		parentFolderID = driveTestParentFolder
	}

	var folderID string
	var folderName string
	var subfolderID string
	var fileID string
	var copyID string
	var testFileContent []byte

	t.Run("folders/create", func(t *testing.T) {
		folderName = testName(t)
		out := zoho(t, "drive", "folders", "create",
			"--name", folderName, "--parent", parentFolderID)
		folderID = extractDriveID(t, out)
		cleanup.trackDriveFolder(folderID)

		attrs := getDriveFileAttr(t, folderID)
		name := fmt.Sprintf("%v", attrs["name"])
		if !strings.HasPrefix(name, folderName) {
			t.Errorf("folder name: got %q, want prefix %q", name, folderName)
		}
		assertEqual(t, fmt.Sprintf("%v", attrs["parent_id"]), parentFolderID)
		t.Logf("created folder %s (%s)", folderID, name)
	})

	t.Run("folders/create-subfolder", func(t *testing.T) {
		requireID(t, folderID, "folders/create must have succeeded")
		subName := testName(t) + "_sub"
		out := zoho(t, "drive", "folders", "create",
			"--name", subName, "--parent", folderID)
		subfolderID = extractDriveID(t, out)
		cleanup.trackDriveFolder(subfolderID)

		attrs := getDriveFileAttr(t, subfolderID)
		assertEqual(t, fmt.Sprintf("%v", attrs["parent_id"]), folderID)
		t.Logf("created subfolder %s", subfolderID)
	})

	t.Run("folders/breadcrumb", func(t *testing.T) {
		requireID(t, subfolderID, "folders/create-subfolder must have succeeded")
		out := zoho(t, "drive", "folders", "breadcrumb", subfolderID)
		m := parseJSON(t, out)
		data, ok := m["data"].([]any)
		if !ok || len(data) == 0 {
			t.Fatalf("expected non-empty breadcrumb data:\n%s", truncate(out, 500))
		}
		item, ok := data[0].(map[string]any)
		if !ok {
			t.Fatalf("expected object in breadcrumb data[0]")
		}
		attrs, ok := item["attributes"].(map[string]any)
		if !ok {
			t.Fatalf("expected attributes in breadcrumb data[0]")
		}
		parentIDs, ok := attrs["parent_ids"].([]any)
		if !ok || len(parentIDs) == 0 {
			t.Fatalf("expected non-empty parent_ids in breadcrumb:\n%s", truncate(out, 500))
		}
		foundFolder := false
		for _, p := range parentIDs {
			entry, ok := p.(map[string]any)
			if !ok {
				continue
			}
			if fmt.Sprintf("%v", entry["resource_id"]) == folderID {
				foundFolder = true
			}
		}
		if !foundFolder {
			t.Errorf("parent folder %s not found in breadcrumb path", folderID)
		}
	})

	t.Run("upload", func(t *testing.T) {
		requireID(t, folderID, "folders/create must have succeeded")
		tmpDir := t.TempDir()
		testFile := tmpDir + "/" + testPrefix + "_drive.txt"
		testFileContent = []byte("ZOHOTEST drive integration " + time.Now().String())
		if err := os.WriteFile(testFile, testFileContent, 0644); err != nil {
			t.Fatalf("failed to create test file: %v", err)
		}
		out := zoho(t, "drive", "upload", testFile, "--folder", folderID)
		fileID = extractDriveUploadID(t, out)
		cleanup.trackDriveFile(fileID)

		attrs := getDriveFileAttr(t, fileID)
		assertEqual(t, fmt.Sprintf("%v", attrs["parent_id"]), folderID)
		t.Logf("uploaded file %s", fileID)
	})

	t.Run("files/list", func(t *testing.T) {
		requireID(t, folderID, "folders/create must have succeeded")
		requireID(t, fileID, "upload must have succeeded")
		out := zoho(t, "drive", "files", "list", "--folder", folderID)
		arr := parseJSONArray(t, out)
		assertNonEmpty(t, arr, "expected at least one file in folder")
		_, found := findInArray(arr, fileID)
		if !found {
			t.Errorf("uploaded file %s not found in folder listing", fileID)
		}
	})

	t.Run("files/get", func(t *testing.T) {
		requireID(t, fileID, "upload must have succeeded")
		m := getDriveFile(t, fileID)
		data := m["data"].(map[string]any)
		assertEqual(t, data["id"], fileID)
		attrs := data["attributes"].(map[string]any)
		assertEqual(t, fmt.Sprintf("%v", attrs["parent_id"]), folderID)
		assertEqual(t, fmt.Sprintf("%v", attrs["status"]), "1")
	})

	t.Run("files/search", func(t *testing.T) {
		requireID(t, fileID, "upload must have succeeded")
		retryUntil(t, 60*time.Second, func() bool {
			out, err := zohoMayFail(t, "drive", "files", "search",
				"--query", testPrefix, "--team", teamID, "--mode", "name")
			if err != nil {
				return false
			}
			var arr []map[string]any
			if jsonErr := json.Unmarshal([]byte(out), &arr); jsonErr != nil || arr == nil {
				return false
			}
			_, found := findInArray(arr, fileID)
			return found
		})
	})

	t.Run("files/rename", func(t *testing.T) {
		requireID(t, fileID, "upload must have succeeded")
		newName := testName(t) + "_renamed.txt"
		out := zoho(t, "drive", "files", "rename", fileID, "--name", newName)
		assertDriveAttr(t, out, "name", newName)

		attrs := getDriveFileAttr(t, fileID)
		assertEqual(t, fmt.Sprintf("%v", attrs["name"]), newName)
	})

	t.Run("files/versions", func(t *testing.T) {
		requireID(t, fileID, "upload must have succeeded")
		out := zoho(t, "drive", "files", "versions", fileID)
		m := parseJSON(t, out)
		data, ok := m["data"].([]any)
		if !ok || len(data) == 0 {
			t.Fatalf("expected non-empty versions data:\n%s", truncate(out, 500))
		}
		first, ok := data[0].(map[string]any)
		if !ok {
			t.Fatalf("expected object in versions data[0]")
		}
		attrs, ok := first["attributes"].(map[string]any)
		if !ok {
			t.Fatalf("expected attributes in version data[0]")
		}
		vn := fmt.Sprintf("%v", attrs["version_number"])
		if vn == "" || vn == "<nil>" {
			t.Errorf("expected version_number in version attributes")
		}
		t.Logf("file has %d version(s), latest version_number=%s", len(data), vn)
	})

	t.Run("files/copy", func(t *testing.T) {
		requireID(t, fileID, "upload must have succeeded")
		requireID(t, subfolderID, "folders/create-subfolder must have succeeded")
		out := zoho(t, "drive", "files", "copy", fileID, "--to", subfolderID)
		copyID = extractDriveID(t, out)
		cleanup.trackDriveFile(copyID)

		attrs := getDriveFileAttr(t, copyID)
		assertEqual(t, fmt.Sprintf("%v", attrs["parent_id"]), subfolderID)
		t.Logf("copied to %s", copyID)
	})

	t.Run("files/move", func(t *testing.T) {
		requireID(t, copyID, "files/copy must have succeeded")
		requireID(t, folderID, "folders/create must have succeeded")
		zoho(t, "drive", "files", "move", copyID, "--to", folderID)

		attrs := getDriveFileAttr(t, copyID)
		assertEqual(t, fmt.Sprintf("%v", attrs["parent_id"]), folderID)

		subOut := zoho(t, "drive", "files", "list", "--folder", subfolderID)
		subArr := parseJSONArray(t, subOut)
		_, found := findInArray(subArr, copyID)
		if found {
			t.Errorf("moved file %s still found in source subfolder", copyID)
		}
	})

	t.Run("download/to-file", func(t *testing.T) {
		requireID(t, fileID, "upload must have succeeded")
		tmpDir := t.TempDir()
		downloadPath := tmpDir + "/downloaded.txt"
		out := zoho(t, "drive", "download", fileID, "--output", downloadPath)
		m := parseJSON(t, out)
		if fmt.Sprintf("%v", m["ok"]) != "true" {
			t.Errorf("expected ok=true in download response:\n%s", truncate(out, 500))
		}
		downloaded, err := os.ReadFile(downloadPath)
		if err != nil {
			t.Fatalf("failed to read downloaded file: %v", err)
		}
		if !bytes.Equal(downloaded, testFileContent) {
			t.Errorf("downloaded content mismatch: got %d bytes, want %d bytes",
				len(downloaded), len(testFileContent))
		}
	})

	t.Run("download/to-stdout", func(t *testing.T) {
		requireID(t, fileID, "upload must have succeeded")
		r := runZoho(t, "drive", "download", fileID)
		if r.ExitCode != 0 {
			t.Fatalf("download to stdout failed (exit %d): %s", r.ExitCode, r.Stderr)
		}
		if !bytes.Equal([]byte(r.Stdout), testFileContent) {
			t.Errorf("stdout content mismatch: got %d bytes, want %d bytes",
				len(r.Stdout), len(testFileContent))
		}
	})

	t.Run("files/trash", func(t *testing.T) {
		requireID(t, copyID, "files/copy must have succeeded")
		out := zoho(t, "drive", "files", "trash", copyID)
		assertDriveAttr(t, out, "status", "51")
	})

	t.Run("files/trash-list", func(t *testing.T) {
		requireID(t, copyID, "files/trash must have succeeded")
		retryUntil(t, 30*time.Second, func() bool {
			out, err := zohoMayFail(t, "drive", "files", "trash-list",
				"--team-folder", parentFolderID)
			if err != nil {
				return false
			}
			var arr []map[string]any
			if jsonErr := json.Unmarshal([]byte(out), &arr); jsonErr != nil {
				return false
			}
			_, found := findInArray(arr, copyID)
			return found
		})
	})

	t.Run("files/restore", func(t *testing.T) {
		requireID(t, copyID, "files/trash must have succeeded")
		out := zoho(t, "drive", "files", "restore", copyID)
		assertDriveAttr(t, out, "status", "1")

		attrs := getDriveFileAttr(t, copyID)
		assertEqual(t, fmt.Sprintf("%v", attrs["status"]), "1")

		retryUntil(t, 30*time.Second, func() bool {
			trashOut, err := zohoMayFail(t, "drive", "files", "trash-list",
				"--team-folder", parentFolderID)
			if err != nil {
				return false
			}
			var arr []map[string]any
			if jsonErr := json.Unmarshal([]byte(trashOut), &arr); jsonErr != nil {
				return false
			}
			_, found := findInArray(arr, copyID)
			return !found
		})
	})

	t.Run("files/trash-cleanup", func(t *testing.T) {
		requireID(t, copyID, "files/restore must have succeeded")
		out := zoho(t, "drive", "files", "trash", copyID)
		assertDriveAttr(t, out, "status", "51")
	})
}

func TestDriveErrors(t *testing.T) {
	t.Run("share-link-skip", func(t *testing.T) {
		t.Skip("share link returns 500 on this account (plan limitation) — zc-8mek")
	})

	t.Run("upload-url-removed", func(t *testing.T) {
		r := runZoho(t, "drive", "upload-url", "https://example.com", "--folder", "fake")
		if r.ExitCode == 0 {
			t.Error("upload-url command should not exist (removed — endpoint never existed)")
		}
	})

	t.Run("bad-file-id", func(t *testing.T) {
		r := runZoho(t, "drive", "files", "get", "nonexistent_id_12345")
		if r.ExitCode == 0 {
			t.Error("expected non-zero exit code for nonexistent file ID")
		}
	})

	t.Run("bad-auth", func(t *testing.T) {
		r := runZohoWithEnv(t, map[string]string{
			"ZOHO_CLIENT_ID":     "bad_client_id",
			"ZOHO_CLIENT_SECRET": "bad_client_secret",
			"ZOHO_REFRESH_TOKEN": "bad_refresh_token",
			"ZOHO_DC":            "com",
		}, "drive", "teams", "me")
		assertExitCode(t, r, 2)
	})
}

func TestDriveEmergencyCleanup(t *testing.T) {
	if os.Getenv("ZOHO_EMERGENCY_CLEANUP") == "" {
		t.Skip("set ZOHO_EMERGENCY_CLEANUP=1 to run")
	}
	teamID := requireDriveTeamID(t)
	out, err := zohoMayFail(t, "drive", "files", "search",
		"--query", testPrefix, "--team", teamID)
	if err != nil {
		t.Fatalf("search for orphaned test files failed: %v", err)
	}
	var arr []map[string]any
	if jsonErr := json.Unmarshal([]byte(out), &arr); jsonErr != nil {
		t.Fatalf("failed to parse search results: %v", jsonErr)
	}
	t.Logf("found %d orphaned test files", len(arr))
	for _, item := range arr {
		id := fmt.Sprintf("%v", item["id"])
		attrs, _ := item["attributes"].(map[string]any)
		name := ""
		if attrs != nil {
			name = fmt.Sprintf("%v", attrs["name"])
		}
		if !strings.HasPrefix(name, testPrefix) {
			continue
		}
		t.Logf("trashing orphaned file %s (%s)", id, name)
		zohoIgnoreError(t, "drive", "files", "trash", id)
	}
}
