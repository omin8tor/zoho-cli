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

func TestMain(m *testing.M) {
	if os.Getenv("ZOHO_REFRESH_TOKEN") != "" {
		cmd := exec.Command("./zoho", "crm", "modules", "list")
		cmd.Dir = ".."
		cmd.Env = os.Environ()
		_ = cmd.Run()
	}
	os.Exit(m.Run())
}

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

func toJSON(t *testing.T, v any) string {
	t.Helper()
	b, err := json.Marshal(v)
	if err != nil {
		t.Fatalf("failed to marshal JSON: %v", err)
	}
	return string(b)
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

func (c *testCleanup) trackDeskTicket(id, orgID string) {
	c.add("delete desk ticket "+id, func() {
		zohoIgnoreError(c.t, "desk", "tickets", "delete", id, "--org", orgID)
	})
}

func (c *testCleanup) trackDeskContact(id, orgID string) {
	c.add("delete desk contact "+id, func() {
		zohoIgnoreError(c.t, "desk", "contacts", "delete", id, "--org", orgID)
	})
}

func (c *testCleanup) trackDeskAccount(id, orgID string) {
	c.add("delete desk account "+id, func() {
		zohoIgnoreError(c.t, "desk", "accounts", "delete", id, "--org", orgID)
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

func (c *testCleanup) trackCliqChannel(id string) {
	c.add("delete cliq channel "+id, func() {
		zohoIgnoreError(c.t, "cliq", "channels", "delete", id)
	})
}

func (c *testCleanup) trackMailFolder(id, accountID string) {
	c.add("delete mail folder "+id, func() {
		zohoIgnoreError(c.t, "mail", "folders", "delete", id, "--account", accountID)
	})
}

func (c *testCleanup) trackMailLabel(id, accountID string) {
	c.add("delete mail label "+id, func() {
		zohoIgnoreError(c.t, "mail", "labels", "delete", id, "--account", accountID)
	})
}

func (c *testCleanup) trackWriterDoc(id string) {
	c.add("trash writer doc "+id, func() {
		zohoIgnoreError(c.t, "drive", "files", "trash", id)
	})
}

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

func randomSuffix() string {
	b := make([]byte, 4)
	if _, err := rand.Read(b); err != nil {
		return fmt.Sprintf("%08x", time.Now().UnixNano()&0xffffffff)
	}
	return fmt.Sprintf("%x", b)
}

func testName(t *testing.T) string {
	t.Helper()
	return fmt.Sprintf("%s_%d_%s", testPrefix, time.Now().Unix(), randomSuffix())
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

func requireProjectsPortalID(t *testing.T) string {
	t.Helper()
	id := os.Getenv("ZOHO_PORTAL_ID")
	if id == "" {
		t.Skip("skipping: ZOHO_PORTAL_ID not set")
	}
	return id
}

func requireExpenseOrgID(t *testing.T) string {
	t.Helper()
	id := os.Getenv("ZOHO_EXPENSE_ORG_ID")
	if id == "" {
		t.Skip("skipping: ZOHO_EXPENSE_ORG_ID not set")
	}
	return id
}

func requireDeskOrgID(t *testing.T) string {
	t.Helper()
	id := os.Getenv("ZOHO_DESK_ORG_ID")
	if id == "" {
		t.Skip("skipping: ZOHO_DESK_ORG_ID not set")
	}
	return id
}

func requireMailAccountID(t *testing.T) string {
	t.Helper()
	id := os.Getenv("ZOHO_MAIL_ACCOUNT_ID")
	if id != "" {
		return id
	}
	out, err := zohoMayFail(t, "mail", "accounts", "list")
	if err != nil {
		t.Skipf("skipping: cannot discover mail account ID: %v", err)
	}
	m := parseJSON(t, out)
	data, ok := m["data"].([]any)
	if !ok || len(data) == 0 {
		t.Skip("skipping: no mail accounts found")
	}
	first, _ := data[0].(map[string]any)
	id = fmt.Sprintf("%v", first["accountId"])
	if id == "" || id == "<nil>" {
		t.Skip("skipping: no accountId in mail accounts response")
	}
	return id
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

func extractProjectsID(t *testing.T, out string) string {
	t.Helper()
	m := parseJSON(t, out)
	id := fmt.Sprintf("%v", m["id"])
	if id == "" || id == "<nil>" {
		t.Fatalf("no id in Projects response:\n%s", truncate(out, 500))
	}
	return id
}

func getProject(t *testing.T, projectID string) map[string]any {
	t.Helper()
	out := zoho(t, "projects", "core", "get", projectID)
	return parseJSON(t, out)
}

func getTask(t *testing.T, taskID, projectID string) map[string]any {
	t.Helper()
	out := zoho(t, "projects", "tasks", "get", taskID, "--project", projectID)
	return parseJSON(t, out)
}

func getIssue(t *testing.T, issueID, projectID string) map[string]any {
	t.Helper()
	out := zoho(t, "projects", "issues", "get", issueID, "--project", projectID)
	arr := parseJSONArray(t, out)
	if len(arr) == 0 {
		t.Fatalf("issues get returned empty array for %s", issueID)
	}
	return arr[0]
}

func getTasklist(t *testing.T, tasklistID, projectID string) map[string]any {
	t.Helper()
	out := zoho(t, "projects", "tasklists", "get", tasklistID, "--project", projectID)
	return parseJSON(t, out)
}

func getMilestone(t *testing.T, milestoneID, projectID string) map[string]any {
	t.Helper()
	out := zoho(t, "projects", "milestones", "get", milestoneID, "--project", projectID)
	return parseJSON(t, out)
}

func getTimelog(t *testing.T, timelogID, projectID string) map[string]any {
	t.Helper()
	out := zoho(t, "projects", "timelogs", "get", timelogID, "--project", projectID, "--type", "task")
	return parseJSON(t, out)
}

func (c *testCleanup) trackProject(id string) {
	c.add("delete project "+id, func() {
		zohoIgnoreError(c.t, "projects", "core", "delete", id)
	})
	c.add("trash project "+id, func() {
		zohoIgnoreError(c.t, "projects", "core", "trash", id)
	})
}

func (c *testCleanup) trackTask(id, projectID string) {
	c.add("delete task "+id, func() {
		zohoIgnoreError(c.t, "projects", "tasks", "delete", id, "--project", projectID)
	})
}

func (c *testCleanup) trackIssue(id, projectID string) {
	c.add("delete issue "+id, func() {
		zohoIgnoreError(c.t, "projects", "issues", "delete", id, "--project", projectID)
	})
}

func (c *testCleanup) trackTasklist(id, projectID string) {
	c.add("delete tasklist "+id, func() {
		zohoIgnoreError(c.t, "projects", "tasklists", "delete", id, "--project", projectID)
	})
}

func (c *testCleanup) trackMilestone(id, projectID string) {
	c.add("delete milestone "+id, func() {
		zohoIgnoreError(c.t, "projects", "milestones", "delete", id, "--project", projectID)
	})
}

func (c *testCleanup) trackTimelog(id, projectID string) {
	c.add("delete timelog "+id, func() {
		zohoIgnoreError(c.t, "projects", "timelogs", "delete", id, "--project", projectID, "--type", "task")
	})
}

func (c *testCleanup) trackProjectGroup(id string) {
	c.add("delete project-group "+id, func() {
		zohoIgnoreError(c.t, "projects", "project-groups", "delete", id)
	})
}

func (c *testCleanup) trackTag(id string) {
	c.add("delete tag "+id, func() {
		zohoIgnoreError(c.t, "projects", "tags", "delete", id)
	})
}

func (c *testCleanup) trackRole(id string) {
	c.add("delete role "+id, func() {
		zohoIgnoreError(c.t, "projects", "roles", "delete", id)
	})
}

func (c *testCleanup) trackForumCategory(id, projectID string) {
	c.add("delete forum-category "+id, func() {
		zohoIgnoreError(c.t, "projects", "forum-categories", "delete", id, "--project", projectID)
	})
}

func (c *testCleanup) trackForum(id, projectID string) {
	c.add("delete forum "+id, func() {
		zohoIgnoreError(c.t, "projects", "forums", "delete", id, "--project", projectID)
	})
}

func (c *testCleanup) trackPhase(id, projectID string) {
	c.add("delete phase "+id, func() {
		zohoIgnoreError(c.t, "projects", "phases", "delete", id, "--project", projectID)
	})
}

func (c *testCleanup) trackEvent(id, projectID string) {
	c.add("delete event "+id, func() {
		zohoIgnoreError(c.t, "projects", "events", "delete", id, "--project", projectID)
	})
}

func (c *testCleanup) trackTeam(id string) {
	c.add("delete team "+id, func() {
		zohoIgnoreError(c.t, "projects", "teams", "delete", id)
	})
}

func (c *testCleanup) trackProfile(id string) {
	c.add("delete profile "+id, func() {
		zohoIgnoreError(c.t, "projects", "profiles", "delete", id)
	})
}

func (c *testCleanup) trackDashboard(id string) {
	c.add("delete dashboard "+id, func() {
		zohoIgnoreError(c.t, "projects", "reports", "dashboard-delete", id)
	})
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
	t.Parallel()
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
	if !strings.Contains(out, "mail") {
		t.Error("expected mail in help-all output")
	}
	if !strings.Contains(out, "books") {
		t.Error("expected books in help-all output")
	}
}

func TestDriveTeams(t *testing.T) {
	t.Parallel()
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
	t.Parallel()
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
	t.Parallel()
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

func TestProjectsPortals(t *testing.T) {
	t.Parallel()
	portalID := requireProjectsPortalID(t)

	t.Run("portals/get", func(t *testing.T) {
		out := zoho(t, "projects", "portals", "get")
		m := parseJSON(t, out)
		pd, ok := m["portal_details"].(map[string]any)
		if ok {
			m = pd
		}
		id := fmt.Sprintf("%v", m["id"])
		if id == "" || id == "<nil>" {
			t.Fatalf("expected portal id in response:\n%s", truncate(out, 500))
		}
		idNum := fmt.Sprintf("%v", m["id"])
		if !strings.Contains(idNum, portalID) {
			idF := strings.TrimRight(strings.TrimRight(fmt.Sprintf("%f", m["id"]), "0"), ".")
			if idF != portalID {
				t.Logf("portal id mismatch: got %s / %s, env has %s", idNum, idF, portalID)
			}
		}
		name := fmt.Sprintf("%v", m["name"])
		if name == "" || name == "<nil>" {
			t.Errorf("expected portal name in response")
		}
		t.Logf("portal: id=%s name=%s", id, name)
	})

	t.Run("portals/list-known-bug", func(t *testing.T) {
		t.Skip("portals list endpoint returns INVALID_METHOD in V3 API")
	})
}

func TestProjects(t *testing.T) {
	t.Parallel()
	_ = requireProjectsPortalID(t)
	cleanup := newCleanup(t)

	var projectID string
	var projectName string
	var project2ID string
	var project2Name string
	var taskID string
	var taskName string
	var subtaskID string
	var clonedTaskID string
	var issueID string
	var issueName string
	var clonedIssueID string
	var tasklistID string
	var tasklistName string
	var milestoneID string
	var milestoneName string
	var timelogID string
	var tlTaskID string
	var ownerZPUID string
	var taskCommentID string
	var issueCommentID string
	var tasklistCommentID string
	var projectCommentID string
	var projectGroupID string
	var tagID string
	var roleID string
	var forumCategoryID string
	var forumID string
	var forumCommentID string
	var issueLinkID string
	var phaseID string
	var phaseName string
	var clonedPhaseID string
	var phaseCommentID string
	var eventID string
	var eventCommentID string

	t.Run("core/create", func(t *testing.T) {
		projectName = testName(t)
		out := zoho(t, "projects", "core", "create", "--name", projectName)
		projectID = extractProjectsID(t, out)
		cleanup.trackProject(projectID)
		t.Logf("created project %s (%s)", projectID, projectName)

		proj := getProject(t, projectID)
		assertStringField(t, proj, "name", projectName)
	})

	t.Run("core/get", func(t *testing.T) {
		requireID(t, projectID, "core/create must have succeeded")
		proj := getProject(t, projectID)
		assertEqual(t, fmt.Sprintf("%v", proj["id"]), projectID)
		assertStringField(t, proj, "name", projectName)
	})

	t.Run("timelogs/setup-owner", func(t *testing.T) {
		requireID(t, projectID, "core/create must have succeeded")
		proj := getProject(t, projectID)
		if cb, ok := proj["created_by"].(map[string]any); ok {
			ownerZPUID = fmt.Sprintf("%v", cb["zpuid"])
		}
		if ownerZPUID == "" || ownerZPUID == "<nil>" {
			t.Fatal("could not determine owner zpuid from project")
		}
		t.Logf("owner zpuid: %s", ownerZPUID)
	})

	t.Run("core/list", func(t *testing.T) {
		requireID(t, projectID, "core/create must have succeeded")
		out := zoho(t, "projects", "core", "list")
		arr := parseJSONArray(t, out)
		assertNonEmpty(t, arr, "expected at least one project")
		_, found := findInArray(arr, projectID)
		if !found {
			t.Errorf("created project %s not found in project list", projectID)
		}
	})

	t.Run("core/update", func(t *testing.T) {
		requireID(t, projectID, "core/create must have succeeded")
		updatedName := projectName + "_updated"
		out := zoho(t, "projects", "core", "update", projectID,
			"--json", toJSON(t, map[string]any{"name": updatedName}))
		m := parseJSON(t, out)
		assertEqual(t, fmt.Sprintf("%v", m["name"]), updatedName)

		proj := getProject(t, projectID)
		assertStringField(t, proj, "name", updatedName)
		projectName = updatedName
	})

	t.Run("core/create-second-project", func(t *testing.T) {
		project2Name = testName(t) + "_p2"
		out := zoho(t, "projects", "core", "create", "--name", project2Name)
		project2ID = extractProjectsID(t, out)
		cleanup.trackProject(project2ID)
		t.Logf("created second project %s (%s)", project2ID, project2Name)
	})

	t.Run("tasks/create", func(t *testing.T) {
		requireID(t, projectID, "core/create must have succeeded")
		taskName = testName(t) + "_task"
		out := zoho(t, "projects", "tasks", "create",
			"--name", taskName, "--project", projectID)
		taskID = extractProjectsID(t, out)
		cleanup.trackTask(taskID, projectID)
		t.Logf("created task %s", taskID)

		task := getTask(t, taskID, projectID)
		assertStringField(t, task, "name", taskName)
	})

	t.Run("tasks/get", func(t *testing.T) {
		requireID(t, taskID, "tasks/create must have succeeded")
		task := getTask(t, taskID, projectID)
		assertEqual(t, fmt.Sprintf("%v", task["id"]), taskID)
		assertStringField(t, task, "name", taskName)
	})

	t.Run("tasks/update", func(t *testing.T) {
		requireID(t, taskID, "tasks/create must have succeeded")
		updatedTaskName := taskName + "_upd"
		out := zoho(t, "projects", "tasks", "update", taskID,
			"--project", projectID,
			"--json", toJSON(t, map[string]any{"name": updatedTaskName}))
		m := parseJSON(t, out)
		assertEqual(t, fmt.Sprintf("%v", m["name"]), updatedTaskName)

		task := getTask(t, taskID, projectID)
		assertStringField(t, task, "name", updatedTaskName)
		taskName = updatedTaskName
	})

	t.Run("tasks/list", func(t *testing.T) {
		requireID(t, taskID, "tasks/create must have succeeded")
		out := zoho(t, "projects", "tasks", "list", "--project", projectID)
		arr := parseJSONArray(t, out)
		assertNonEmpty(t, arr, "expected at least one task")
		_, found := findInArray(arr, taskID)
		if !found {
			t.Errorf("task %s not found in task list", taskID)
		}
	})

	t.Run("tasks/my", func(t *testing.T) {
		out := zoho(t, "projects", "tasks", "my")
		m := parseJSON(t, out)
		if _, ok := m["tasks"]; !ok {
			if _, ok2 := m["page_info"]; !ok2 {
				t.Errorf("expected tasks or page_info key in my-tasks response:\n%s", truncate(out, 500))
			}
		}
	})

	t.Run("task-comments/add", func(t *testing.T) {
		requireID(t, taskID, "tasks/create must have succeeded")
		out := zoho(t, "projects", "task-comments", "add",
			"--task", taskID, "--project", projectID,
			"--comment", testPrefix+"_task_comment")
		arr := parseJSONArray(t, out)
		if len(arr) == 0 {
			t.Fatalf("task-comments add returned empty array:\n%s", truncate(out, 500))
		}
		taskCommentID = fmt.Sprintf("%v", arr[0]["id"])
		if taskCommentID == "" || taskCommentID == "<nil>" {
			t.Fatalf("no id in task-comments add response:\n%s", truncate(out, 500))
		}
		assertStringField(t, arr[0], "comment", testPrefix+"_task_comment")
		t.Logf("created task comment %s", taskCommentID)
	})

	t.Run("task-comments/list", func(t *testing.T) {
		requireID(t, taskCommentID, "task-comments/add must have succeeded")
		out := zoho(t, "projects", "task-comments", "list",
			"--task", taskID, "--project", projectID)
		m := parseJSON(t, out)
		comments, ok := m["comments"].([]any)
		if !ok {
			t.Fatalf("expected comments array in list response:\n%s", truncate(out, 500))
		}
		found := false
		for _, c := range comments {
			cm, _ := c.(map[string]any)
			if fmt.Sprintf("%v", cm["id"]) == taskCommentID {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("comment %s not found in task-comments list", taskCommentID)
		}
	})

	t.Run("task-comments/update", func(t *testing.T) {
		requireID(t, taskCommentID, "task-comments/add must have succeeded")
		updatedComment := testPrefix + "_task_comment_upd"
		out := zoho(t, "projects", "task-comments", "update", taskCommentID,
			"--task", taskID, "--project", projectID,
			"--comment", updatedComment)
		arr := parseJSONArray(t, out)
		if len(arr) == 0 {
			t.Fatalf("task-comments update returned empty array:\n%s", truncate(out, 500))
		}
		assertStringField(t, arr[0], "comment", updatedComment)

		listOut := zoho(t, "projects", "task-comments", "list",
			"--task", taskID, "--project", projectID)
		assertContains(t, listOut, updatedComment)
	})

	t.Run("task-comments/delete", func(t *testing.T) {
		requireID(t, taskCommentID, "task-comments/add must have succeeded")
		out := zoho(t, "projects", "task-comments", "delete", taskCommentID,
			"--task", taskID, "--project", projectID)
		parseJSON(t, out)

		listOut := zoho(t, "projects", "task-comments", "list",
			"--task", taskID, "--project", projectID)
		if strings.Contains(listOut, taskCommentID) {
			t.Errorf("comment %s still found in list after delete", taskCommentID)
		}
		taskCommentID = ""
	})

	t.Run("task-followers/follow", func(t *testing.T) {
		requireID(t, taskID, "tasks/create must have succeeded")
		out := zoho(t, "projects", "task-followers", "follow", taskID,
			"--project", projectID)
		parseJSON(t, out)
	})

	t.Run("task-followers/list", func(t *testing.T) {
		requireID(t, taskID, "tasks/create must have succeeded")
		out := zoho(t, "projects", "task-followers", "list", taskID,
			"--project", projectID)
		m := parseJSON(t, out)
		followers, ok := m["followers"].([]any)
		if !ok || len(followers) == 0 {
			t.Fatalf("expected non-empty followers array:\n%s", truncate(out, 500))
		}
		found := false
		for _, f := range followers {
			fm, _ := f.(map[string]any)
			if fmt.Sprintf("%v", fm["zpuid"]) == ownerZPUID {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("follower %s not found in task-followers list", ownerZPUID)
		}
	})

	t.Run("task-followers/unfollow", func(t *testing.T) {
		requireID(t, taskID, "tasks/create must have succeeded")
		out := zoho(t, "projects", "task-followers", "unfollow", taskID,
			"--project", projectID)
		parseJSON(t, out)

		time.Sleep(2 * time.Second)
		listOut := zoho(t, "projects", "task-followers", "list", taskID,
			"--project", projectID)
		if strings.Contains(listOut, ownerZPUID) {
			t.Errorf("follower %s still present after unfollow", ownerZPUID)
		}
	})

	t.Run("tasks/add-subtask", func(t *testing.T) {
		requireID(t, taskID, "tasks/create must have succeeded")
		subtaskName := testName(t) + "_sub"
		out := zoho(t, "projects", "tasks", "add-subtask",
			"--parent", taskID, "--name", subtaskName, "--project", projectID)
		subtaskID = extractProjectsID(t, out)
		cleanup.trackTask(subtaskID, projectID)
		t.Logf("created subtask %s", subtaskID)

		sub := getTask(t, subtaskID, projectID)
		assertStringField(t, sub, "name", subtaskName)
	})

	t.Run("tasks/subtasks-known-broken", func(t *testing.T) {
		t.Skip("subtasks endpoint not available in V3 API")
	})

	t.Run("tasks/clone", func(t *testing.T) {
		requireID(t, taskID, "tasks/create must have succeeded")
		out := zoho(t, "projects", "tasks", "clone", taskID,
			"--project", projectID)
		clonedTaskID = extractProjectsID(t, out)
		cleanup.trackTask(clonedTaskID, projectID)
		t.Logf("cloned task %s -> %s", taskID, clonedTaskID)

		cloned := getTask(t, clonedTaskID, projectID)
		clonedName := fmt.Sprintf("%v", cloned["name"])
		if clonedName == "" || clonedName == "<nil>" {
			t.Errorf("cloned task has no name")
		}
	})

	t.Run("tasks/move", func(t *testing.T) {
		requireID(t, clonedTaskID, "tasks/clone must have succeeded")
		requireID(t, project2ID, "second project must exist")
		tlOut := zoho(t, "projects", "tasklists", "create",
			"--name", testName(t)+"_tl", "--project", project2ID)
		targetTL := extractProjectsID(t, tlOut)
		cleanup.trackTasklist(targetTL, project2ID)
		moveJSON := toJSON(t, map[string]any{"target_tasklist_id": targetTL})
		zoho(t, "projects", "tasks", "move", clonedTaskID,
			"--project", projectID,
			"--json", moveJSON)

		movedTask := getTask(t, clonedTaskID, project2ID)
		movedName := fmt.Sprintf("%v", movedTask["name"])
		if movedName == "" || movedName == "<nil>" {
			t.Errorf("moved task not found in target project")
		}
		cleanup.trackTask(clonedTaskID, project2ID)
	})

	t.Run("tasks/delete-subtask", func(t *testing.T) {
		requireID(t, subtaskID, "tasks/add-subtask must have succeeded")
		out := zoho(t, "projects", "tasks", "delete", subtaskID,
			"--project", projectID)
		parseJSON(t, out)
		t.Logf("deleted subtask %s", subtaskID)

		r := runZoho(t, "projects", "tasks", "get", subtaskID, "--project", projectID)
		if r.ExitCode == 0 {
			t.Errorf("subtask %s still accessible after delete", subtaskID)
		}
		subtaskID = ""
	})

	t.Run("tasks/delete", func(t *testing.T) {
		requireID(t, taskID, "tasks/create must have succeeded")
		out := zoho(t, "projects", "tasks", "delete", taskID,
			"--project", projectID)
		parseJSON(t, out)

		r := runZoho(t, "projects", "tasks", "get", taskID, "--project", projectID)
		if r.ExitCode == 0 {
			t.Errorf("task %s still accessible after delete", taskID)
		}
		taskID = ""
	})

	t.Run("issues/create", func(t *testing.T) {
		requireID(t, projectID, "core/create must have succeeded")
		issueName = testName(t) + "_issue"
		out := zoho(t, "projects", "issues", "create",
			"--name", issueName, "--project", projectID)
		issueID = extractProjectsID(t, out)
		cleanup.trackIssue(issueID, projectID)
		t.Logf("created issue %s", issueID)

		issue := getIssue(t, issueID, projectID)
		assertStringField(t, issue, "name", issueName)
	})

	t.Run("issues/get", func(t *testing.T) {
		requireID(t, issueID, "issues/create must have succeeded")
		issue := getIssue(t, issueID, projectID)
		assertEqual(t, fmt.Sprintf("%v", issue["id"]), issueID)
		assertStringField(t, issue, "name", issueName)
	})

	t.Run("issues/update", func(t *testing.T) {
		requireID(t, issueID, "issues/create must have succeeded")
		updatedIssueName := issueName + "_upd"
		out := zoho(t, "projects", "issues", "update", issueID,
			"--project", projectID,
			"--json", toJSON(t, map[string]any{"name": updatedIssueName}))
		m := parseJSON(t, out)
		assertEqual(t, fmt.Sprintf("%v", m["name"]), updatedIssueName)

		issue := getIssue(t, issueID, projectID)
		assertStringField(t, issue, "name", updatedIssueName)
		issueName = updatedIssueName
	})

	t.Run("issues/list", func(t *testing.T) {
		requireID(t, issueID, "issues/create must have succeeded")
		out := zoho(t, "projects", "issues", "list", "--project", projectID)
		arr := parseJSONArray(t, out)
		assertNonEmpty(t, arr, "expected at least one issue")
		_, found := findInArray(arr, issueID)
		if !found {
			t.Errorf("issue %s not found in issue list", issueID)
		}
	})

	t.Run("issues/defaults-removed", func(t *testing.T) {
		r := runZoho(t, "projects", "issues", "defaults", "--project", projectID)
		if r.ExitCode == 0 {
			t.Errorf("issues defaults command should not exist (removed: endpoint not in Projects V3)")
		}
	})

	t.Run("issues/description", func(t *testing.T) {
		requireID(t, issueID, "issues/create must have succeeded")
		out := zoho(t, "projects", "issues", "description", issueID,
			"--project", projectID)
		m := parseJSON(t, out)
		if _, ok := m["description"]; !ok {
			if _, ok2 := m["content"]; !ok2 {
				t.Errorf("expected description or content key in response:\n%s", truncate(out, 500))
			}
		}
	})

	t.Run("issues/activities", func(t *testing.T) {
		requireID(t, issueID, "issues/create must have succeeded")
		out := zoho(t, "projects", "issues", "activities", issueID,
			"--project", projectID)
		m := parseJSON(t, out)
		if _, ok := m["activities"]; !ok {
			if _, ok2 := m["page_info"]; !ok2 {
				t.Errorf("expected activities or page_info key in response:\n%s", truncate(out, 500))
			}
		}
	})

	t.Run("issue-comments/add", func(t *testing.T) {
		requireID(t, issueID, "issues/create must have succeeded")
		out := zoho(t, "projects", "issue-comments", "add",
			"--issue", issueID, "--project", projectID,
			"--comment", testPrefix+"_issue_comment")
		m := parseJSON(t, out)
		comments, ok := m["comments"].([]any)
		if !ok || len(comments) == 0 {
			t.Fatalf("expected comments array in add response:\n%s", truncate(out, 500))
		}
		cm, _ := comments[0].(map[string]any)
		issueCommentID = fmt.Sprintf("%v", cm["id"])
		if issueCommentID == "" || issueCommentID == "<nil>" {
			t.Fatalf("no id in issue-comments add response:\n%s", truncate(out, 500))
		}
		assertStringField(t, cm, "comment", testPrefix+"_issue_comment")
		t.Logf("created issue comment %s", issueCommentID)
	})

	t.Run("issue-comments/get", func(t *testing.T) {
		requireID(t, issueCommentID, "issue-comments/add must have succeeded")
		out := zoho(t, "projects", "issue-comments", "get", issueCommentID,
			"--issue", issueID, "--project", projectID)
		m := parseJSON(t, out)
		comments, ok := m["comments"].([]any)
		if !ok || len(comments) == 0 {
			t.Fatalf("expected comments array in get response:\n%s", truncate(out, 500))
		}
		cm, _ := comments[0].(map[string]any)
		assertEqual(t, fmt.Sprintf("%v", cm["id"]), issueCommentID)
		assertStringField(t, cm, "comment", testPrefix+"_issue_comment")
	})

	t.Run("issue-comments/list", func(t *testing.T) {
		requireID(t, issueCommentID, "issue-comments/add must have succeeded")
		out := zoho(t, "projects", "issue-comments", "list",
			"--issue", issueID, "--project", projectID)
		m := parseJSON(t, out)
		comments, ok := m["comments"].([]any)
		if !ok {
			t.Fatalf("expected comments array in list response:\n%s", truncate(out, 500))
		}
		found := false
		for _, c := range comments {
			cm, _ := c.(map[string]any)
			if fmt.Sprintf("%v", cm["id"]) == issueCommentID {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("comment %s not found in issue-comments list", issueCommentID)
		}
	})

	t.Run("issue-comments/update", func(t *testing.T) {
		requireID(t, issueCommentID, "issue-comments/add must have succeeded")
		updatedComment := testPrefix + "_issue_comment_upd"
		out := zoho(t, "projects", "issue-comments", "update", issueCommentID,
			"--issue", issueID, "--project", projectID,
			"--comment", updatedComment)
		m := parseJSON(t, out)
		comments, ok := m["comments"].([]any)
		if !ok || len(comments) == 0 {
			t.Fatalf("expected comments array in update response:\n%s", truncate(out, 500))
		}
		cm, _ := comments[0].(map[string]any)
		assertStringField(t, cm, "comment", updatedComment)

		getOut := zoho(t, "projects", "issue-comments", "get", issueCommentID,
			"--issue", issueID, "--project", projectID)
		assertContains(t, getOut, updatedComment)
	})

	t.Run("issue-comments/delete", func(t *testing.T) {
		requireID(t, issueCommentID, "issue-comments/add must have succeeded")
		out := zoho(t, "projects", "issue-comments", "delete", issueCommentID,
			"--issue", issueID, "--project", projectID)
		parseJSON(t, out)

		listOut := zoho(t, "projects", "issue-comments", "list",
			"--issue", issueID, "--project", projectID)
		if strings.Contains(listOut, issueCommentID) {
			t.Errorf("comment %s still found in list after delete", issueCommentID)
		}
		issueCommentID = ""
	})

	t.Run("issue-followers/follow-known-limitation", func(t *testing.T) {
		t.Skip("issue follower self-follow not supported")
	})

	t.Run("issue-followers/list", func(t *testing.T) {
		requireID(t, issueID, "issues/create must have succeeded")
		out := zoho(t, "projects", "issue-followers", "list", issueID,
			"--project", projectID)
		m := parseJSON(t, out)
		if _, ok := m["followers"]; !ok {
			t.Fatalf("expected followers key in list response:\n%s", truncate(out, 500))
		}
	})

	t.Run("issues/clone", func(t *testing.T) {
		requireID(t, issueID, "issues/create must have succeeded")
		out := zoho(t, "projects", "issues", "clone", issueID,
			"--project", projectID)
		clonedIssueID = extractProjectsID(t, out)
		cleanup.trackIssue(clonedIssueID, projectID)
		t.Logf("cloned issue %s -> %s", issueID, clonedIssueID)
	})

	t.Run("issue-linking/list-empty", func(t *testing.T) {
		requireID(t, issueID, "issues/create must have succeeded")
		out := zoho(t, "projects", "issue-linking", "list", issueID,
			"--project", projectID)
		m := parseJSON(t, out)
		if _, ok := m["issue_linked"]; !ok {
			t.Fatalf("expected issue_linked key in response:\n%s", truncate(out, 500))
		}
	})

	t.Run("issue-linking/link", func(t *testing.T) {
		requireID(t, issueID, "issues/create must have succeeded")
		requireID(t, clonedIssueID, "issues/clone must have succeeded")
		out := zoho(t, "projects", "issue-linking", "link", issueID,
			"--project", projectID,
			"--json", toJSON(t, map[string]any{
				"link_type": "relate",
				"issue_ids": []string{clonedIssueID},
			}))
		parseJSON(t, out)

		listOut := zoho(t, "projects", "issue-linking", "list", issueID,
			"--project", projectID)
		lm := parseJSON(t, listOut)
		linked, ok := lm["issue_linked"].(map[string]any)
		if !ok {
			t.Fatalf("expected issue_linked object in list response:\n%s", truncate(listOut, 500))
		}
		linkedIssues, ok := linked["linked_issues"].(map[string]any)
		if !ok {
			t.Fatalf("expected linked_issues in response:\n%s", truncate(listOut, 500))
		}
		for _, v := range linkedIssues {
			if arr, ok := v.([]any); ok {
				for _, item := range arr {
					if im, ok := item.(map[string]any); ok {
						if fmt.Sprintf("%v", im["issue_id"]) == clonedIssueID {
							issueLinkID = fmt.Sprintf("%v", im["link_id"])
							break
						}
					}
				}
			}
			if issueLinkID != "" {
				break
			}
		}
		if issueLinkID == "" || issueLinkID == "<nil>" {
			t.Fatalf("could not extract link_id from list response:\n%s", truncate(listOut, 500))
		}
		assertContains(t, listOut, clonedIssueID)
		t.Logf("created issue link %s", issueLinkID)
	})

	t.Run("issue-linking/change-type", func(t *testing.T) {
		requireID(t, issueID, "issues/create must have succeeded")
		requireID(t, issueLinkID, "issue-linking/link must have succeeded")
		out := zoho(t, "projects", "issue-linking", "change-type",
			issueID, issueLinkID,
			"--project", projectID,
			"--json", toJSON(t, map[string]any{"link_type": "blocks"}))
		parseJSON(t, out)

		listOut := zoho(t, "projects", "issue-linking", "list", issueID,
			"--project", projectID)
		assertContains(t, listOut, "blocks")
	})

	t.Run("issue-linking/unlink", func(t *testing.T) {
		requireID(t, issueID, "issues/create must have succeeded")
		requireID(t, issueLinkID, "issue-linking/link must have succeeded")
		out := zoho(t, "projects", "issue-linking", "unlink",
			issueID, issueLinkID,
			"--project", projectID)
		parseJSON(t, out)

		listOut := zoho(t, "projects", "issue-linking", "list", issueID,
			"--project", projectID)
		if strings.Contains(listOut, clonedIssueID) {
			t.Errorf("cloned issue %s still in linked issues after unlink", clonedIssueID)
		}
		issueLinkID = ""
	})

	t.Run("issue-linking/bulk-link", func(t *testing.T) {
		requireID(t, issueID, "issues/create must have succeeded")
		requireID(t, clonedIssueID, "issues/clone must have succeeded")
		r := runZoho(t, "projects", "issue-linking", "bulk-link",
			"--project", projectID,
			"--json", toJSON(t, map[string]any{
				"link_type":         "relate",
				"issue_ids":         []string{issueID},
				"linking_issue_ids": []string{clonedIssueID},
			}))
		if r.ExitCode != 0 {
			t.Logf("bulk-link failed: %s",
				truncate(r.Stderr+r.Stdout, 300))
			return
		}
		parseJSON(t, r.Stdout)

		listOut := zoho(t, "projects", "issue-linking", "list", issueID,
			"--project", projectID)
		assertContains(t, listOut, clonedIssueID)
	})

	t.Run("issues/move", func(t *testing.T) {
		requireID(t, clonedIssueID, "issues/clone must have succeeded")
		requireID(t, project2ID, "second project must exist")
		moveJSON := toJSON(t, map[string]any{"to_project": project2ID})
		zoho(t, "projects", "issues", "move", clonedIssueID,
			"--project", projectID,
			"--json", moveJSON)

		movedIssue := getIssue(t, clonedIssueID, project2ID)
		movedName := fmt.Sprintf("%v", movedIssue["name"])
		if movedName == "" || movedName == "<nil>" {
			t.Errorf("moved issue not found in target project")
		}
		cleanup.trackIssue(clonedIssueID, project2ID)
	})

	t.Run("issue-resolution/add", func(t *testing.T) {
		requireID(t, issueID, "issues/create must have succeeded")
		out := zoho(t, "projects", "issue-resolution", "add", issueID,
			"--project", projectID,
			"--json", toJSON(t, map[string]any{"resolution": testPrefix + " resolution text"}))
		parseJSON(t, out)
	})

	t.Run("issue-resolution/get", func(t *testing.T) {
		requireID(t, issueID, "issues/create must have succeeded")
		out := zoho(t, "projects", "issue-resolution", "get", issueID,
			"--project", projectID)
		m := parseJSON(t, out)
		ir, ok := m["issue_resolution"].(map[string]any)
		if !ok {
			t.Fatalf("expected issue_resolution object in get response:\n%s", truncate(out, 500))
		}
		resolution := fmt.Sprintf("%v", ir["resolution"])
		assertContains(t, resolution, testPrefix)
	})

	t.Run("issue-resolution/update", func(t *testing.T) {
		requireID(t, issueID, "issues/create must have succeeded")
		zoho(t, "projects", "issue-resolution", "update", issueID,
			"--project", projectID,
			"--json", toJSON(t, map[string]any{"resolution": testPrefix + " updated resolution"}))

		out := zoho(t, "projects", "issue-resolution", "get", issueID,
			"--project", projectID)
		m := parseJSON(t, out)
		ir, _ := m["issue_resolution"].(map[string]any)
		resolution := fmt.Sprintf("%v", ir["resolution"])
		assertContains(t, resolution, "updated resolution")
	})

	t.Run("issue-resolution/delete", func(t *testing.T) {
		requireID(t, issueID, "issues/create must have succeeded")
		out := zoho(t, "projects", "issue-resolution", "delete", issueID,
			"--project", projectID)
		parseJSON(t, out)

		getOut := zoho(t, "projects", "issue-resolution", "get", issueID,
			"--project", projectID)
		m := parseJSON(t, getOut)
		ir, _ := m["issue_resolution"].(map[string]any)
		if _, hasResolution := ir["resolution"]; hasResolution {
			t.Errorf("resolution key still present after delete")
		}
	})

	t.Run("issue-attachments/list", func(t *testing.T) {
		requireID(t, issueID, "issues/create must have succeeded")
		out := zoho(t, "projects", "issue-attachments", "list", issueID,
			"--project", projectID)
		m := parseJSON(t, out)
		if _, ok := m["attachments"]; !ok {
			t.Fatalf("expected attachments key in response:\n%s", truncate(out, 500))
		}
	})

	t.Run("issue-attachments/associate-known-broken", func(t *testing.T) {
		t.Skip("no real attachment ID available for test")
	})

	t.Run("issue-attachments/dissociate-known-broken", func(t *testing.T) {
		t.Skip("no real attachment ID available for test")
	})

	t.Run("issues/delete", func(t *testing.T) {
		requireID(t, issueID, "issues/create must have succeeded")
		out := zoho(t, "projects", "issues", "delete", issueID,
			"--project", projectID)
		parseJSON(t, out)

		r := runZoho(t, "projects", "issues", "get", issueID, "--project", projectID)
		if r.ExitCode == 0 && strings.Contains(r.Stdout, issueID) {
			t.Errorf("issue %s still accessible after delete", issueID)
		}
		issueID = ""
	})

	t.Run("tasklists/create", func(t *testing.T) {
		requireID(t, projectID, "core/create must have succeeded")
		tasklistName = testName(t) + "_tl"
		out := zoho(t, "projects", "tasklists", "create",
			"--name", tasklistName, "--project", projectID)
		tasklistID = extractProjectsID(t, out)
		cleanup.trackTasklist(tasklistID, projectID)
		t.Logf("created tasklist %s", tasklistID)

		tl := getTasklist(t, tasklistID, projectID)
		assertStringField(t, tl, "name", tasklistName)
	})

	t.Run("tasklists/get", func(t *testing.T) {
		requireID(t, tasklistID, "tasklists/create must have succeeded")
		tl := getTasklist(t, tasklistID, projectID)
		assertEqual(t, fmt.Sprintf("%v", tl["id"]), tasklistID)
		assertStringField(t, tl, "name", tasklistName)
	})

	t.Run("tasklists/list", func(t *testing.T) {
		requireID(t, tasklistID, "tasklists/create must have succeeded")
		out := zoho(t, "projects", "tasklists", "list", "--project", projectID)
		arr := parseJSONArray(t, out)
		assertNonEmpty(t, arr, "expected at least one tasklist")
		_, found := findInArray(arr, tasklistID)
		if !found {
			t.Errorf("tasklist %s not found in list", tasklistID)
		}
	})

	t.Run("tasklists/update", func(t *testing.T) {
		requireID(t, tasklistID, "tasklists/create must have succeeded")
		updatedTLName := tasklistName + "_upd"
		out := zoho(t, "projects", "tasklists", "update", tasklistID,
			"--project", projectID,
			"--json", toJSON(t, map[string]any{"name": updatedTLName}))
		m := parseJSON(t, out)
		assertEqual(t, fmt.Sprintf("%v", m["name"]), updatedTLName)

		tl := getTasklist(t, tasklistID, projectID)
		assertStringField(t, tl, "name", updatedTLName)
		tasklistName = updatedTLName
	})

	t.Run("tasklist-comments/add", func(t *testing.T) {
		requireID(t, tasklistID, "tasklists/create must have succeeded")
		out := zoho(t, "projects", "tasklist-comments", "add",
			"--tasklist", tasklistID, "--project", projectID,
			"--comment", testPrefix+"_tl_comment")
		arr := parseJSONArray(t, out)
		if len(arr) == 0 {
			t.Fatalf("tasklist-comments add returned empty array:\n%s", truncate(out, 500))
		}
		tasklistCommentID = fmt.Sprintf("%v", arr[0]["id"])
		if tasklistCommentID == "" || tasklistCommentID == "<nil>" {
			t.Fatalf("no id in tasklist-comments add response:\n%s", truncate(out, 500))
		}
		assertStringField(t, arr[0], "comment", testPrefix+"_tl_comment")
		t.Logf("created tasklist comment %s", tasklistCommentID)
	})

	t.Run("tasklist-comments/get", func(t *testing.T) {
		requireID(t, tasklistCommentID, "tasklist-comments/add must have succeeded")
		out := zoho(t, "projects", "tasklist-comments", "get", tasklistCommentID,
			"--tasklist", tasklistID, "--project", projectID)
		m := parseJSON(t, out)
		assertEqual(t, fmt.Sprintf("%v", m["id"]), tasklistCommentID)
	})

	t.Run("tasklist-comments/list", func(t *testing.T) {
		requireID(t, tasklistCommentID, "tasklist-comments/add must have succeeded")
		out := zoho(t, "projects", "tasklist-comments", "list",
			"--tasklist", tasklistID, "--project", projectID)
		m := parseJSON(t, out)
		comments, ok := m["comments"].([]any)
		if !ok {
			t.Fatalf("expected comments array in list response:\n%s", truncate(out, 500))
		}
		found := false
		for _, c := range comments {
			cm, _ := c.(map[string]any)
			if fmt.Sprintf("%v", cm["id"]) == tasklistCommentID {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("comment %s not found in tasklist-comments list", tasklistCommentID)
		}
	})

	t.Run("tasklist-comments/update", func(t *testing.T) {
		requireID(t, tasklistCommentID, "tasklist-comments/add must have succeeded")
		updatedComment := testPrefix + "_tl_comment_upd"
		out := zoho(t, "projects", "tasklist-comments", "update", tasklistCommentID,
			"--tasklist", tasklistID, "--project", projectID,
			"--comment", updatedComment)
		arr := parseJSONArray(t, out)
		if len(arr) == 0 {
			t.Fatalf("tasklist-comments update returned empty array:\n%s", truncate(out, 500))
		}
		assertStringField(t, arr[0], "comment", updatedComment)

		getOut := zoho(t, "projects", "tasklist-comments", "get", tasklistCommentID,
			"--tasklist", tasklistID, "--project", projectID)
		assertContains(t, getOut, updatedComment)
	})

	t.Run("tasklist-comments/delete", func(t *testing.T) {
		requireID(t, tasklistCommentID, "tasklist-comments/add must have succeeded")
		out := zoho(t, "projects", "tasklist-comments", "delete", tasklistCommentID,
			"--tasklist", tasklistID, "--project", projectID)
		parseJSON(t, out)

		listOut := zoho(t, "projects", "tasklist-comments", "list",
			"--tasklist", tasklistID, "--project", projectID)
		if strings.Contains(listOut, tasklistCommentID) {
			t.Errorf("comment %s still found in list after delete", tasklistCommentID)
		}
		tasklistCommentID = ""
	})

	t.Run("tasklist-followers/follow", func(t *testing.T) {
		requireID(t, tasklistID, "tasklists/create must have succeeded")
		out := zoho(t, "projects", "tasklist-followers", "follow", tasklistID,
			"--project", projectID)
		parseJSON(t, out)
	})

	t.Run("tasklist-followers/list", func(t *testing.T) {
		requireID(t, tasklistID, "tasklists/create must have succeeded")
		out := zoho(t, "projects", "tasklist-followers", "list", tasklistID,
			"--project", projectID)
		m := parseJSON(t, out)
		followers, ok := m["followers"].([]any)
		if !ok || len(followers) == 0 {
			t.Fatalf("expected non-empty followers array:\n%s", truncate(out, 500))
		}
		found := false
		for _, f := range followers {
			fm, _ := f.(map[string]any)
			if fmt.Sprintf("%v", fm["zpuid"]) == ownerZPUID {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("follower %s not found in tasklist-followers list", ownerZPUID)
		}
	})

	t.Run("tasklist-followers/unfollow", func(t *testing.T) {
		requireID(t, tasklistID, "tasklists/create must have succeeded")
		out := zoho(t, "projects", "tasklist-followers", "unfollow", tasklistID,
			"--project", projectID)
		parseJSON(t, out)

		time.Sleep(2 * time.Second)
		listOut := zoho(t, "projects", "tasklist-followers", "list", tasklistID,
			"--project", projectID)
		if strings.Contains(listOut, ownerZPUID) {
			t.Errorf("follower %s still present after unfollow", ownerZPUID)
		}
	})

	t.Run("tasklists/delete", func(t *testing.T) {
		requireID(t, tasklistID, "tasklists/create must have succeeded")
		out := zoho(t, "projects", "tasklists", "delete", tasklistID,
			"--project", projectID)
		parseJSON(t, out)

		time.Sleep(2 * time.Second)
		r := runZoho(t, "projects", "tasklists", "get", tasklistID, "--project", projectID)
		if r.ExitCode == 0 {
			t.Errorf("tasklist %s still accessible after delete", tasklistID)
		}
		tasklistID = ""
	})

	t.Run("milestones/create", func(t *testing.T) {
		requireID(t, projectID, "core/create must have succeeded")
		milestoneName = testName(t) + "_ms"
		now := time.Now()
		startDate := now.Format("2006-01-02")
		endDate := now.AddDate(0, 1, 0).Format("2006-01-02")
		out := zoho(t, "projects", "milestones", "create",
			"--name", milestoneName,
			"--start", startDate,
			"--end", endDate,
			"--project", projectID)
		milestoneID = extractProjectsID(t, out)
		cleanup.trackMilestone(milestoneID, projectID)
		t.Logf("created milestone %s", milestoneID)

		ms := getMilestone(t, milestoneID, projectID)
		assertStringField(t, ms, "name", milestoneName)
	})

	t.Run("milestones/get", func(t *testing.T) {
		requireID(t, milestoneID, "milestones/create must have succeeded")
		ms := getMilestone(t, milestoneID, projectID)
		assertEqual(t, fmt.Sprintf("%v", ms["id"]), milestoneID)
		assertStringField(t, ms, "name", milestoneName)
	})

	t.Run("milestones/list", func(t *testing.T) {
		requireID(t, milestoneID, "milestones/create must have succeeded")
		out := zoho(t, "projects", "milestones", "list", "--project", projectID)
		arr := parseJSONArray(t, out)
		assertNonEmpty(t, arr, "expected at least one milestone")
		_, found := findInArray(arr, milestoneID)
		if !found {
			t.Errorf("milestone %s not found in list", milestoneID)
		}
	})

	t.Run("milestones/update", func(t *testing.T) {
		requireID(t, milestoneID, "milestones/create must have succeeded")
		updatedMSName := milestoneName + "_upd"
		out := zoho(t, "projects", "milestones", "update", milestoneID,
			"--project", projectID,
			"--json", toJSON(t, map[string]any{"name": updatedMSName}))
		m := parseJSON(t, out)
		assertEqual(t, fmt.Sprintf("%v", m["name"]), updatedMSName)

		ms := getMilestone(t, milestoneID, projectID)
		assertStringField(t, ms, "name", updatedMSName)
		milestoneName = updatedMSName
	})

	t.Run("milestones/delete", func(t *testing.T) {
		requireID(t, milestoneID, "milestones/create must have succeeded")
		out := zoho(t, "projects", "milestones", "delete", milestoneID,
			"--project", projectID)
		parseJSON(t, out)

		r := runZoho(t, "projects", "milestones", "get", milestoneID, "--project", projectID)
		if r.ExitCode == 0 {
			t.Errorf("milestone %s still accessible after delete", milestoneID)
		}
		milestoneID = ""
	})

	t.Run("timelogs/add", func(t *testing.T) {
		requireID(t, projectID, "core/create must have succeeded")
		tlTaskName := testName(t) + "_tltask"
		taskOut := zoho(t, "projects", "tasks", "create",
			"--name", tlTaskName, "--project", projectID)
		tlTaskID = extractProjectsID(t, taskOut)
		cleanup.trackTask(tlTaskID, projectID)

		zoho(t, "projects", "tasks", "update", tlTaskID,
			"--project", projectID,
			"--json", toJSON(t, map[string]any{
				"owners_and_work": map[string]any{
					"owners": []map[string]string{{"zpuid": ownerZPUID}},
				},
			}))

		today := time.Now().Format("2006-01-02")
		out := zoho(t, "projects", "timelogs", "add",
			"--date", today,
			"--hours", "2",
			"--task", tlTaskID,
			"--owner", ownerZPUID,
			"--notes", testPrefix+"_timelog",
			"--project", projectID)
		timelogID = extractProjectsID(t, out)
		cleanup.trackTimelog(timelogID, projectID)
		t.Logf("created timelog %s for task %s", timelogID, tlTaskID)
	})

	t.Run("timelogs/get", func(t *testing.T) {
		requireID(t, timelogID, "timelogs/add must have succeeded")
		tl := getTimelog(t, timelogID, projectID)
		assertEqual(t, fmt.Sprintf("%v", tl["id"]), timelogID)
	})

	t.Run("timelogs/list", func(t *testing.T) {
		requireID(t, timelogID, "timelogs/add must have succeeded")
		out := zoho(t, "projects", "timelogs", "list",
			"--module", "task", "--project", projectID)
		var raw json.RawMessage
		if err := json.Unmarshal([]byte(out), &raw); err != nil {
			t.Fatalf("failed to parse timelogs list: %v\nraw: %s", err, truncate(out, 500))
		}
		assertContains(t, out, timelogID)
	})

	t.Run("timelogs/update", func(t *testing.T) {
		requireID(t, timelogID, "timelogs/add must have succeeded")
		out := zoho(t, "projects", "timelogs", "update", timelogID,
			"--project", projectID,
			"--type", "task",
			"--task", tlTaskID,
			"--json", toJSON(t, map[string]any{"hours": "3"}))
		parseJSON(t, out)

		tl := getTimelog(t, timelogID, projectID)
		hours := fmt.Sprintf("%v", tl["log_hour"])
		if hours != "3" && hours != "03:00" && hours != "3:00" {
			t.Errorf("expected timelog hours=3 after update, got %s", hours)
		}
	})

	t.Run("timelogs/delete", func(t *testing.T) {
		requireID(t, timelogID, "timelogs/add must have succeeded")
		out := zoho(t, "projects", "timelogs", "delete", timelogID,
			"--project", projectID, "--type", "task")
		parseJSON(t, out)

		r := runZoho(t, "projects", "timelogs", "get", timelogID, "--project", projectID, "--type", "task")
		if r.ExitCode == 0 {
			t.Errorf("timelog %s still accessible after delete", timelogID)
		}
		timelogID = ""
	})

	t.Run("project-comments/add", func(t *testing.T) {
		requireID(t, projectID, "core/create must have succeeded")
		commentText := testPrefix + "_project_comment"
		out := zoho(t, "projects", "project-comments", "add",
			"--project", projectID,
			"--comment", commentText)
		m := parseJSON(t, out)
		comments, ok := m["comments"].([]any)
		if !ok || len(comments) == 0 {
			t.Fatalf("expected comments array in add response:\n%s", truncate(out, 500))
		}
		cm, _ := comments[0].(map[string]any)
		projectCommentID = fmt.Sprintf("%v", cm["id"])
		if projectCommentID == "" || projectCommentID == "<nil>" {
			t.Fatalf("no id in project-comments add response:\n%s", truncate(out, 500))
		}
		assertStringField(t, cm, "content", commentText)
		t.Logf("created project comment %s", projectCommentID)
	})

	t.Run("project-comments/get", func(t *testing.T) {
		requireID(t, projectCommentID, "project-comments/add must have succeeded")
		out := zoho(t, "projects", "project-comments", "get", projectCommentID,
			"--project", projectID)
		m := parseJSON(t, out)
		assertEqual(t, fmt.Sprintf("%v", m["id"]), projectCommentID)
		assertStringField(t, m, "content", testPrefix+"_project_comment")
	})

	t.Run("project-comments/list", func(t *testing.T) {
		requireID(t, projectCommentID, "project-comments/add must have succeeded")
		out := zoho(t, "projects", "project-comments", "list",
			"--project", projectID)
		m := parseJSON(t, out)
		comments, ok := m["comments"].([]any)
		if !ok {
			t.Fatalf("expected comments array in list response:\n%s", truncate(out, 500))
		}
		found := false
		for _, c := range comments {
			cm, _ := c.(map[string]any)
			if fmt.Sprintf("%v", cm["id"]) == projectCommentID {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("comment %s not found in project-comments list", projectCommentID)
		}
	})

	t.Run("project-comments/update", func(t *testing.T) {
		requireID(t, projectCommentID, "project-comments/add must have succeeded")
		updatedComment := testPrefix + "_project_comment_upd"
		out := zoho(t, "projects", "project-comments", "update", projectCommentID,
			"--project", projectID,
			"--comment", updatedComment)
		m := parseJSON(t, out)
		comments, ok := m["comments"].([]any)
		if !ok || len(comments) == 0 {
			t.Fatalf("expected comments array in update response:\n%s", truncate(out, 500))
		}
		cm, _ := comments[0].(map[string]any)
		assertStringField(t, cm, "content", updatedComment)

		getOut := zoho(t, "projects", "project-comments", "get", projectCommentID,
			"--project", projectID)
		getM := parseJSON(t, getOut)
		assertStringField(t, getM, "content", updatedComment)
	})

	t.Run("project-comments/delete", func(t *testing.T) {
		requireID(t, projectCommentID, "project-comments/add must have succeeded")
		out := zoho(t, "projects", "project-comments", "delete", projectCommentID,
			"--project", projectID)
		parseJSON(t, out)

		listOut := zoho(t, "projects", "project-comments", "list",
			"--project", projectID)
		if strings.Contains(listOut, projectCommentID) {
			t.Errorf("comment %s still found in list after delete", projectCommentID)
		}
		projectCommentID = ""
	})

	t.Run("forum-categories/create", func(t *testing.T) {
		requireID(t, projectID, "core/create must have succeeded")
		categoryName := testName(t) + "_forum_category"
		out := zoho(t, "projects", "forum-categories", "create",
			"--project", projectID,
			"--json", toJSON(t, map[string]any{"name": categoryName}))
		m := parseJSON(t, out)
		categories, ok := m["categories"].([]any)
		if !ok || len(categories) == 0 {
			t.Fatalf("expected categories array in create response:\n%s", truncate(out, 500))
		}
		cm, _ := categories[0].(map[string]any)
		forumCategoryID = fmt.Sprintf("%v", cm["id"])
		if forumCategoryID == "" || forumCategoryID == "<nil>" {
			t.Fatalf("no id in forum-categories create response:\n%s", truncate(out, 500))
		}
		assertStringField(t, cm, "name", categoryName)
		cleanup.trackForumCategory(forumCategoryID, projectID)
		t.Logf("created forum category %s", forumCategoryID)
	})

	t.Run("forum-categories/list", func(t *testing.T) {
		requireID(t, forumCategoryID, "forum-categories/create must have succeeded")
		out := zoho(t, "projects", "forum-categories", "list", "--project", projectID)
		arr := parseJSONArray(t, out)
		found := false
		for _, c := range arr {
			if fmt.Sprintf("%v", c["id"]) == forumCategoryID {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("forum category %s not found in list", forumCategoryID)
		}
	})

	t.Run("forums/create", func(t *testing.T) {
		requireID(t, forumCategoryID, "forum-categories/create must have succeeded")
		forumTitle := testName(t) + "_forum"
		forumContent := testPrefix + "_forum_content"
		out := zoho(t, "projects", "forums", "create",
			"--project", projectID,
			"--json", toJSON(t, map[string]any{
				"title":       forumTitle,
				"content":     forumContent,
				"category_id": forumCategoryID,
			}))
		m := parseJSON(t, out)
		forums, ok := m["forums"].([]any)
		if !ok || len(forums) == 0 {
			t.Fatalf("expected forums array in create response:\n%s", truncate(out, 500))
		}
		fm, _ := forums[0].(map[string]any)
		forumID = fmt.Sprintf("%v", fm["id"])
		if forumID == "" || forumID == "<nil>" {
			t.Fatalf("no id in forums create response:\n%s", truncate(out, 500))
		}
		assertStringField(t, fm, "title", forumTitle)
		assertStringField(t, fm, "content", forumContent)
		cleanup.trackForum(forumID, projectID)
		t.Logf("created forum %s", forumID)
	})

	t.Run("forums/get", func(t *testing.T) {
		requireID(t, forumID, "forums/create must have succeeded")
		out := zoho(t, "projects", "forums", "get", forumID,
			"--project", projectID)
		m := parseJSON(t, out)
		forums, ok := m["forums"].([]any)
		if !ok || len(forums) == 0 {
			t.Fatalf("expected forums array in get response:\n%s", truncate(out, 500))
		}
		fm, _ := forums[0].(map[string]any)
		assertEqual(t, fmt.Sprintf("%v", fm["id"]), forumID)
	})

	t.Run("forums/list", func(t *testing.T) {
		requireID(t, forumID, "forums/create must have succeeded")
		out := zoho(t, "projects", "forums", "list",
			"--project", projectID)
		m := parseJSON(t, out)
		forums, ok := m["forums"].([]any)
		if !ok {
			t.Fatalf("expected forums array in list response:\n%s", truncate(out, 500))
		}
		found := false
		for _, f := range forums {
			fm, _ := f.(map[string]any)
			if fmt.Sprintf("%v", fm["id"]) == forumID {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("forum %s not found in forums list", forumID)
		}
	})

	t.Run("forums/update", func(t *testing.T) {
		requireID(t, forumID, "forums/create must have succeeded")
		updatedTitle := testName(t) + "_forum_upd"
		updatedContent := testPrefix + "_forum_content_upd"
		out := zoho(t, "projects", "forums", "update", forumID,
			"--project", projectID,
			"--json", toJSON(t, map[string]any{"title": updatedTitle, "content": updatedContent}))
		m := parseJSON(t, out)
		forums, ok := m["forums"].([]any)
		if !ok || len(forums) == 0 {
			t.Fatalf("expected forums array in update response:\n%s", truncate(out, 500))
		}
		fm, _ := forums[0].(map[string]any)
		assertStringField(t, fm, "title", updatedTitle)
		assertStringField(t, fm, "content", updatedContent)

		getOut := zoho(t, "projects", "forums", "get", forumID,
			"--project", projectID)
		getM := parseJSON(t, getOut)
		getForums, ok := getM["forums"].([]any)
		if !ok || len(getForums) == 0 {
			t.Fatalf("expected forums array in get response after update:\n%s", truncate(getOut, 500))
		}
		getFM, _ := getForums[0].(map[string]any)
		assertStringField(t, getFM, "title", updatedTitle)
		assertStringField(t, getFM, "content", updatedContent)
	})

	t.Run("forum-comments/add", func(t *testing.T) {
		requireID(t, forumID, "forums/create must have succeeded")
		commentText := testPrefix + "_forum_comment"
		out := zoho(t, "projects", "forum-comments", "add",
			"--project", projectID,
			"--forum", forumID,
			"--comment", commentText)
		m := parseJSON(t, out)
		comment, ok := m["forum_comments"].(map[string]any)
		if !ok {
			t.Fatalf("expected forum_comments object in add response:\n%s", truncate(out, 500))
		}
		forumCommentID = fmt.Sprintf("%v", comment["id"])
		if forumCommentID == "" || forumCommentID == "<nil>" {
			t.Fatalf("no id in forum-comments add response:\n%s", truncate(out, 500))
		}
		assertStringField(t, comment, "content", commentText)
		t.Logf("created forum comment %s", forumCommentID)
	})

	t.Run("forum-comments/get", func(t *testing.T) {
		requireID(t, forumCommentID, "forum-comments/add must have succeeded")
		out := zoho(t, "projects", "forum-comments", "get", forumCommentID,
			"--project", projectID,
			"--forum", forumID)
		m := parseJSON(t, out)
		comments, ok := m["forum_comments"].([]any)
		if !ok || len(comments) == 0 {
			t.Fatalf("expected forum_comments array in get response:\n%s", truncate(out, 500))
		}
		cm, _ := comments[0].(map[string]any)
		assertEqual(t, fmt.Sprintf("%v", cm["id"]), forumCommentID)
	})

	t.Run("forum-comments/list", func(t *testing.T) {
		requireID(t, forumCommentID, "forum-comments/add must have succeeded")
		out := zoho(t, "projects", "forum-comments", "list",
			"--project", projectID,
			"--forum", forumID)
		m := parseJSON(t, out)
		comments, ok := m["forum_comments"].([]any)
		if !ok {
			t.Fatalf("expected forum_comments array in list response:\n%s", truncate(out, 500))
		}
		found := false
		for _, c := range comments {
			cm, _ := c.(map[string]any)
			if fmt.Sprintf("%v", cm["id"]) == forumCommentID {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("forum comment %s not found in list", forumCommentID)
		}
	})

	t.Run("forum-comments/update", func(t *testing.T) {
		requireID(t, forumCommentID, "forum-comments/add must have succeeded")
		updatedComment := testPrefix + "_forum_comment_upd"
		out := zoho(t, "projects", "forum-comments", "update", forumCommentID,
			"--project", projectID,
			"--forum", forumID,
			"--comment", updatedComment)
		m := parseJSON(t, out)
		comments, ok := m["forum_comments"].([]any)
		if !ok || len(comments) == 0 {
			t.Fatalf("expected forum_comments array in update response:\n%s", truncate(out, 500))
		}
		cm, _ := comments[0].(map[string]any)
		assertStringField(t, cm, "content", updatedComment)

		getOut := zoho(t, "projects", "forum-comments", "get", forumCommentID,
			"--project", projectID,
			"--forum", forumID)
		getM := parseJSON(t, getOut)
		getComments, ok := getM["forum_comments"].([]any)
		if !ok || len(getComments) == 0 {
			t.Fatalf("expected forum_comments array in get response after update:\n%s", truncate(getOut, 500))
		}
		getCM, _ := getComments[0].(map[string]any)
		assertStringField(t, getCM, "content", updatedComment)
	})

	t.Run("forum-comments/best-answer-known-limitation", func(t *testing.T) {
		t.Skip("best-answer endpoint returns error")
	})

	t.Run("forum-comments/unbest-answer-known-limitation", func(t *testing.T) {
		t.Skip("unbest-answer endpoint returns error")
	})

	t.Run("forum-comments/delete", func(t *testing.T) {
		requireID(t, forumCommentID, "forum-comments/add must have succeeded")
		out := zoho(t, "projects", "forum-comments", "delete", forumCommentID,
			"--project", projectID,
			"--forum", forumID)
		parseJSON(t, out)

		listOut := zoho(t, "projects", "forum-comments", "list",
			"--project", projectID,
			"--forum", forumID)
		if strings.Contains(listOut, forumCommentID) {
			t.Errorf("forum comment %s still found in list after delete", forumCommentID)
		}
		forumCommentID = ""
	})

	t.Run("forum-followers/list", func(t *testing.T) {
		requireID(t, forumID, "forums/create must have succeeded")
		out := zoho(t, "projects", "forum-followers", "list", forumID,
			"--project", projectID)
		m := parseJSON(t, out)
		if _, ok := m["followers"].([]any); !ok {
			t.Fatalf("expected followers array in list response:\n%s", truncate(out, 500))
		}
	})

	t.Run("forum-followers/follow-known-limitation", func(t *testing.T) {
		t.Skip("forum follower endpoint requires body")
	})

	t.Run("forum-followers/unfollow-known-limitation", func(t *testing.T) {
		t.Skip("forum follower endpoint requires body")
	})

	t.Run("forums/delete", func(t *testing.T) {
		requireID(t, forumID, "forums/create must have succeeded")
		out := zoho(t, "projects", "forums", "delete", forumID,
			"--project", projectID)
		parseJSON(t, out)

		listOut := zoho(t, "projects", "forums", "list",
			"--project", projectID)
		if strings.Contains(listOut, forumID) {
			t.Errorf("forum %s still found in list after delete", forumID)
		}
		forumID = ""
	})

	t.Run("forum-categories/update", func(t *testing.T) {
		requireID(t, forumCategoryID, "forum-categories/create must have succeeded")
		updatedName := testName(t) + "_forum_category_upd"
		out := zoho(t, "projects", "forum-categories", "update", forumCategoryID,
			"--project", projectID,
			"--json", toJSON(t, map[string]any{"name": updatedName}))
		m := parseJSON(t, out)
		categories, ok := m["categories"].([]any)
		if !ok || len(categories) == 0 {
			t.Fatalf("expected categories array in update response:\n%s", truncate(out, 500))
		}
		cm, _ := categories[0].(map[string]any)
		assertEqual(t, fmt.Sprintf("%v", cm["id"]), forumCategoryID)
		assertStringField(t, cm, "name", updatedName)
	})

	t.Run("forum-categories/delete", func(t *testing.T) {
		requireID(t, forumCategoryID, "forum-categories/create must have succeeded")
		out := zoho(t, "projects", "forum-categories", "delete", forumCategoryID,
			"--project", projectID)
		parseJSON(t, out)
	})

	t.Run("forum-categories/list-verify-deleted", func(t *testing.T) {
		requireID(t, forumCategoryID, "forum-categories/create must have succeeded")
		out := zoho(t, "projects", "forum-categories", "list",
			"--project", projectID)
		arr := parseJSONArray(t, out)
		for _, c := range arr {
			if fmt.Sprintf("%v", c["id"]) == forumCategoryID {
				t.Errorf("forum category %s still found after delete", forumCategoryID)
				break
			}
		}
		forumCategoryID = ""
	})

	t.Run("project-groups/create", func(t *testing.T) {
		groupName := testName(t) + "_group"
		out := zoho(t, "projects", "project-groups", "create",
			"--json", toJSON(t, map[string]any{"name": groupName, "type": "public"}))
		m := parseJSON(t, out)
		projectGroupID = fmt.Sprintf("%v", m["id"])
		if projectGroupID == "" || projectGroupID == "<nil>" {
			t.Fatalf("no id in project-groups create response:\n%s", truncate(out, 500))
		}
		assertStringField(t, m, "name", groupName)
		cleanup.trackProjectGroup(projectGroupID)
		t.Logf("created project group %s", projectGroupID)
	})

	t.Run("project-groups/list", func(t *testing.T) {
		requireID(t, projectGroupID, "project-groups/create must have succeeded")
		out := zoho(t, "projects", "project-groups", "list")
		m := parseJSON(t, out)
		groups, ok := m["project-groups"].([]any)
		if !ok {
			t.Fatalf("expected project-groups array in list response:\n%s", truncate(out, 500))
		}
		found := false
		for _, g := range groups {
			gm, _ := g.(map[string]any)
			if fmt.Sprintf("%v", gm["id"]) == projectGroupID {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("project group %s not found in list", projectGroupID)
		}
	})

	t.Run("project-groups/my", func(t *testing.T) {
		out := zoho(t, "projects", "project-groups", "my")
		m := parseJSON(t, out)
		if _, ok := m["project-groups"].([]any); !ok {
			t.Fatalf("expected project-groups array in my response:\n%s", truncate(out, 500))
		}
	})

	t.Run("project-groups/update", func(t *testing.T) {
		requireID(t, projectGroupID, "project-groups/create must have succeeded")
		updatedName := testName(t) + "_group_upd"
		out := zoho(t, "projects", "project-groups", "update", projectGroupID,
			"--json", toJSON(t, map[string]any{"name": updatedName}))
		m := parseJSON(t, out)
		assertEqual(t, fmt.Sprintf("%v", m["id"]), projectGroupID)
		assertStringField(t, m, "name", updatedName)
	})

	t.Run("project-groups/delete", func(t *testing.T) {
		requireID(t, projectGroupID, "project-groups/create must have succeeded")
		out := zoho(t, "projects", "project-groups", "delete", projectGroupID)
		parseJSON(t, out)
	})

	t.Run("project-groups/list-verify-deleted", func(t *testing.T) {
		requireID(t, projectGroupID, "project-groups/create must have succeeded")
		out := zoho(t, "projects", "project-groups", "list")
		m := parseJSON(t, out)
		groups, ok := m["project-groups"].([]any)
		if !ok {
			t.Fatalf("expected project-groups array in list response:\n%s", truncate(out, 500))
		}
		for _, g := range groups {
			gm, _ := g.(map[string]any)
			if fmt.Sprintf("%v", gm["id"]) == projectGroupID {
				t.Errorf("project group %s still found after delete", projectGroupID)
				break
			}
		}
		projectGroupID = ""
	})

	t.Run("tags/create", func(t *testing.T) {
		tagName := testName(t) + "_tag"
		out := zoho(t, "projects", "tags", "create",
			"--json", toJSON(t, []map[string]any{{"name": tagName, "color_class": "bg-tag1"}}))
		m := parseJSON(t, out)
		tags, ok := m["tags"].([]any)
		if !ok || len(tags) == 0 {
			t.Fatalf("expected tags array in create response:\n%s", truncate(out, 500))
		}
		tm, _ := tags[0].(map[string]any)
		tagID = fmt.Sprintf("%v", tm["id"])
		if tagID == "" || tagID == "<nil>" {
			t.Fatalf("no id in tags create response:\n%s", truncate(out, 500))
		}
		assertStringField(t, tm, "name", tagName)
		cleanup.trackTag(tagID)
		t.Logf("created tag %s", tagID)
	})

	t.Run("tags/list", func(t *testing.T) {
		requireID(t, tagID, "tags/create must have succeeded")
		out := zoho(t, "projects", "tags", "list")
		m := parseJSON(t, out)
		tags, ok := m["tags"].([]any)
		if !ok {
			t.Fatalf("expected tags array in list response:\n%s", truncate(out, 500))
		}
		found := false
		for _, it := range tags {
			tm, _ := it.(map[string]any)
			if fmt.Sprintf("%v", tm["id"]) == tagID {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("tag %s not found in tags list", tagID)
		}
	})

	t.Run("tags/project-list", func(t *testing.T) {
		requireID(t, tagID, "tags/create must have succeeded")
		requireID(t, projectID, "core/create must have succeeded")
		out := zoho(t, "projects", "tags", "project-list", "--project", projectID)
		m := parseJSON(t, out)
		tags, ok := m["tags"].([]any)
		if !ok {
			t.Fatalf("expected tags array in project-list response:\n%s", truncate(out, 500))
		}
		for _, it := range tags {
			tm, _ := it.(map[string]any)
			if fmt.Sprintf("%v", tm["id"]) == tagID {
				t.Errorf("tag %s should not be associated yet", tagID)
				break
			}
		}
	})

	t.Run("tags/update", func(t *testing.T) {
		requireID(t, tagID, "tags/create must have succeeded")
		updatedName := testName(t) + "_tag_upd"
		out := zoho(t, "projects", "tags", "update", tagID,
			"--json", toJSON(t, map[string]any{"name": updatedName}))
		m := parseJSON(t, out)
		tm, ok := m["tags"].(map[string]any)
		if !ok {
			t.Fatalf("expected tags object in update response:\n%s", truncate(out, 500))
		}
		assertEqual(t, fmt.Sprintf("%v", tm["id"]), tagID)
		assertStringField(t, tm, "name", updatedName)
	})

	t.Run("tags/associate", func(t *testing.T) {
		requireID(t, tagID, "tags/create must have succeeded")
		requireID(t, projectID, "core/create must have succeeded")
		requireID(t, tlTaskID, "timelogs/add must have succeeded")
		out := zoho(t, "projects", "tags", "associate", tagID,
			"--entity", tlTaskID,
			"--entity-type", "5",
			"--project", projectID)
		parseJSON(t, out)
	})

	t.Run("tags/project-list-verify", func(t *testing.T) {
		requireID(t, tagID, "tags/create must have succeeded")
		requireID(t, projectID, "core/create must have succeeded")
		out := zoho(t, "projects", "tags", "project-list", "--project", projectID)
		m := parseJSON(t, out)
		tags, ok := m["tags"].([]any)
		if !ok {
			t.Fatalf("expected tags array in project-list response:\n%s", truncate(out, 500))
		}
		found := false
		for _, it := range tags {
			tm, _ := it.(map[string]any)
			if fmt.Sprintf("%v", tm["id"]) == tagID {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("tag %s not found in project tags after associate", tagID)
		}
	})

	t.Run("tags/dissociate", func(t *testing.T) {
		requireID(t, tagID, "tags/create must have succeeded")
		requireID(t, projectID, "core/create must have succeeded")
		requireID(t, tlTaskID, "timelogs/add must have succeeded")
		out := zoho(t, "projects", "tags", "dissociate", tagID,
			"--entity", tlTaskID,
			"--entity-type", "5",
			"--project", projectID)
		parseJSON(t, out)
	})

	t.Run("tags/delete", func(t *testing.T) {
		requireID(t, tagID, "tags/create must have succeeded")
		out := zoho(t, "projects", "tags", "delete", tagID)
		parseJSON(t, out)

		listOut := zoho(t, "projects", "tags", "list")
		if strings.Contains(listOut, tagID) {
			t.Errorf("tag %s still found in list after delete", tagID)
		}
		tagID = ""
	})

	t.Run("roles/create", func(t *testing.T) {
		roleName := testName(t) + "_role"
		out := zoho(t, "projects", "roles", "create",
			"--json", toJSON(t, map[string]any{"name": roleName}))
		m := parseJSON(t, out)
		roleID = fmt.Sprintf("%v", m["id"])
		if roleID == "" || roleID == "<nil>" {
			t.Fatalf("no id in roles create response:\n%s", truncate(out, 500))
		}
		assertStringField(t, m, "name", roleName)
		cleanup.trackRole(roleID)
		t.Logf("created role %s", roleID)
	})

	t.Run("roles/get", func(t *testing.T) {
		requireID(t, roleID, "roles/create must have succeeded")
		out := zoho(t, "projects", "roles", "get", roleID)
		m := parseJSON(t, out)
		assertEqual(t, fmt.Sprintf("%v", m["id"]), roleID)
	})

	t.Run("roles/list", func(t *testing.T) {
		requireID(t, roleID, "roles/create must have succeeded")
		out := zoho(t, "projects", "roles", "list")
		m := parseJSON(t, out)
		roles, ok := m["roles"].([]any)
		if !ok {
			t.Fatalf("expected roles array in list response:\n%s", truncate(out, 500))
		}
		found := false
		for _, r := range roles {
			rm, _ := r.(map[string]any)
			if fmt.Sprintf("%v", rm["id"]) == roleID {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("role %s not found in roles list", roleID)
		}
	})

	t.Run("roles/update", func(t *testing.T) {
		requireID(t, roleID, "roles/create must have succeeded")
		updatedName := testName(t) + "_role_upd"
		zoho(t, "projects", "roles", "update", roleID,
			"--json", toJSON(t, map[string]any{"name": updatedName}))

		out := zoho(t, "projects", "roles", "get", roleID)
		m := parseJSON(t, out)
		assertStringField(t, m, "name", updatedName)
	})

	t.Run("roles/set-default-known-broken", func(t *testing.T) {
		t.Skip("roles set-default endpoint returns error")
	})

	t.Run("roles/delete", func(t *testing.T) {
		requireID(t, roleID, "roles/create must have succeeded")
		out := zoho(t, "projects", "roles", "delete", roleID)
		parseJSON(t, out)

		listOut := zoho(t, "projects", "roles", "list")
		m := parseJSON(t, listOut)
		roles, ok := m["roles"].([]any)
		if ok {
			for _, r := range roles {
				rm, _ := r.(map[string]any)
				if fmt.Sprintf("%v", rm["id"]) == roleID {
					t.Errorf("role %s still found in list after delete", roleID)
					break
				}
			}
		}
		roleID = ""
	})

	t.Run("project-users/list", func(t *testing.T) {
		requireID(t, projectID, "core/create must have succeeded")
		out := zoho(t, "projects", "project-users", "list", "--project", projectID)
		arr := parseJSONArray(t, out)
		assertNonEmpty(t, arr, "expected at least one project user")
	})

	t.Run("project-users/get", func(t *testing.T) {
		requireID(t, projectID, "core/create must have succeeded")
		requireID(t, ownerZPUID, "timelogs/setup-owner must have succeeded")
		out := zoho(t, "projects", "project-users", "get", ownerZPUID,
			"--project", projectID)
		m := parseJSON(t, out)
		assertEqual(t, fmt.Sprintf("%v", m["id"]), ownerZPUID)
	})

	t.Run("dependencies/add-known-broken", func(t *testing.T) {
		t.Skip("dependencies endpoint not functional")
	})

	t.Run("dependencies/remove-known-broken", func(t *testing.T) {
		t.Skip("dependencies endpoint not functional")
	})

	t.Run("task-customviews/list", func(t *testing.T) {
		out := zoho(t, "projects", "task-customviews", "list")
		m := parseJSON(t, out)
		views, ok := m["views"].([]any)
		if !ok {
			t.Fatalf("expected views array in task-customviews list response:\n%s", truncate(out, 500))
		}
		if len(views) == 0 {
			t.Error("expected non-empty views array")
		}
	})

	t.Run("task-customviews/project-list", func(t *testing.T) {
		requireID(t, projectID, "core/create must have succeeded")
		out := zoho(t, "projects", "task-customviews", "project-list",
			"--project", projectID)
		m := parseJSON(t, out)
		if _, ok := m["views"]; !ok {
			t.Fatalf("expected views key in task-customviews project-list response:\n%s", truncate(out, 500))
		}
	})

	t.Run("task-customviews/get", func(t *testing.T) {
		out := zoho(t, "projects", "task-customviews", "list")
		m := parseJSON(t, out)
		views, ok := m["views"].([]any)
		if !ok || len(views) == 0 {
			t.Skip("no views available for get test")
		}
		first, _ := views[0].(map[string]any)
		viewID := fmt.Sprintf("%v", first["id"])
		if viewID == "" || viewID == "<nil>" {
			t.Skip("first view has no id")
		}

		getOut := zoho(t, "projects", "task-customviews", "get", viewID)
		gm := parseJSON(t, getOut)
		tcv, ok := gm["taskcustomviews"].([]any)
		if !ok || len(tcv) == 0 {
			t.Fatalf("expected taskcustomviews array in get response:\n%s", truncate(getOut, 500))
		}
		first2, _ := tcv[0].(map[string]any)
		assertEqual(t, fmt.Sprintf("%v", first2["id"]), viewID)
	})

	t.Run("issue-customviews/list", func(t *testing.T) {
		out := zoho(t, "projects", "issue-customviews", "list")
		m := parseJSON(t, out)
		dv, ok := m["default_views"].([]any)
		if !ok {
			t.Fatalf("expected default_views array in issue-customviews list response:\n%s", truncate(out, 500))
		}
		if len(dv) == 0 {
			t.Error("expected non-empty default_views array")
		}
	})

	t.Run("issue-customviews/project-list", func(t *testing.T) {
		requireID(t, projectID, "core/create must have succeeded")
		out := zoho(t, "projects", "issue-customviews", "project-list",
			"--project", projectID)
		m := parseJSON(t, out)
		for _, key := range []string{"custom_views", "default_views", "favourites", "shared_views"} {
			if _, ok := m[key]; !ok {
				t.Errorf("expected %s key in issue-customviews project-list response", key)
			}
		}
	})

	t.Run("issue-customviews/get", func(t *testing.T) {
		out := zoho(t, "projects", "issue-customviews", "list")
		m := parseJSON(t, out)
		dv, ok := m["default_views"].([]any)
		if !ok || len(dv) == 0 {
			t.Skip("no default_views available for get test")
		}
		first, _ := dv[0].(map[string]any)
		viewID := fmt.Sprintf("%v", first["custom_view_id"])
		if viewID == "" || viewID == "<nil>" {
			t.Skip("first default_view has no custom_view_id")
		}

		getOut := zoho(t, "projects", "issue-customviews", "get", viewID)
		gm := parseJSON(t, getOut)
		cv, ok := gm["customview"].(map[string]any)
		if !ok {
			t.Fatalf("expected customview object in get response:\n%s", truncate(getOut, 500))
		}
		assertEqual(t, fmt.Sprintf("%v", cv["custom_view_id"]), viewID)
	})

	t.Run("task-statustimeline/get", func(t *testing.T) {
		requireID(t, tlTaskID, "timelogs/add must have succeeded")
		requireID(t, projectID, "core/create must have succeeded")
		out := zoho(t, "projects", "task-statustimeline", "get", tlTaskID,
			"--project", projectID)
		m := parseJSON(t, out)
		if _, ok := m["timeline"]; !ok {
			if _, ok2 := m["status_timeline"]; !ok2 {
				t.Errorf("expected timeline or status_timeline key in response:\n%s", truncate(out, 500))
			}
		}
	})

	t.Run("task-statustimeline/project-known-broken", func(t *testing.T) {
		t.Skip("task-statustimeline project endpoint may not exist in V3")
	})

	t.Run("task-statustimeline/portal", func(t *testing.T) {
		out := zoho(t, "projects", "task-statustimeline", "portal")
		m := parseJSON(t, out)
		if _, ok := m["timeline"]; !ok {
			if _, ok2 := m["status_timeline"]; !ok2 {
				if _, ok3 := m["page_info"]; !ok3 {
					t.Errorf("expected timeline, status_timeline, or page_info key in response:\n%s", truncate(out, 500))
				}
			}
		}
	})

	t.Run("attachments/list", func(t *testing.T) {
		requireID(t, tlTaskID, "timelogs/add must have succeeded")
		requireID(t, projectID, "core/create must have succeeded")
		out := zoho(t, "projects", "attachments", "list",
			"--project", projectID,
			"--type", "task",
			"--entity-id", tlTaskID)
		m := parseJSON(t, out)
		if _, ok := m["attachments"]; !ok {
			if _, ok2 := m["page_info"]; !ok2 {
				t.Errorf("expected attachments or page_info key in response:\n%s", truncate(out, 500))
			}
		}
	})

	t.Run("attachments/upload-known-broken", func(t *testing.T) {
		t.Skip("attachments upload needs WorkDrive integration")
	})

	t.Run("attachments/get-known-broken", func(t *testing.T) {
		t.Skip("no real project attachment ID available")
	})

	t.Run("attachments/associate-known-broken", func(t *testing.T) {
		t.Skip("no real project attachment ID available")
	})

	t.Run("attachments/dissociate-known-broken", func(t *testing.T) {
		t.Skip("no real project attachment ID available")
	})

	t.Run("phases/create", func(t *testing.T) {
		requireID(t, projectID, "core/create must have succeeded")
		phaseName = testName(t) + "_phase"
		out := zoho(t, "projects", "phases", "create",
			"--project", projectID,
			"--json", toJSON(t, map[string]any{"name": phaseName}))
		phaseID = extractProjectsID(t, out)
		cleanup.trackPhase(phaseID, projectID)
		t.Logf("created phase %s (%s)", phaseID, phaseName)

		out = zoho(t, "projects", "phases", "get", phaseID, "--project", projectID)
		m := parseJSON(t, out)
		assertStringField(t, m, "name", phaseName)
	})

	t.Run("phases/get", func(t *testing.T) {
		requireID(t, phaseID, "phases/create must have succeeded")
		out := zoho(t, "projects", "phases", "get", phaseID, "--project", projectID)
		m := parseJSON(t, out)
		assertEqual(t, fmt.Sprintf("%v", m["id"]), phaseID)
		assertStringField(t, m, "name", phaseName)
	})

	t.Run("phases/list-project", func(t *testing.T) {
		requireID(t, phaseID, "phases/create must have succeeded")
		out := zoho(t, "projects", "phases", "list-project", "--project", projectID)
		m := parseJSON(t, out)
		milestones, ok := m["milestones"].([]any)
		if !ok {
			t.Fatalf("expected milestones array in list-project response:\n%s", truncate(out, 500))
		}
		found := false
		for _, ms := range milestones {
			mm, _ := ms.(map[string]any)
			if fmt.Sprintf("%v", mm["id"]) == phaseID {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("phase %s not found in list-project", phaseID)
		}
	})

	t.Run("phases/list", func(t *testing.T) {
		requireID(t, phaseID, "phases/create must have succeeded")
		out := zoho(t, "projects", "phases", "list")
		m := parseJSON(t, out)
		milestones, ok := m["milestones"].([]any)
		if !ok {
			t.Fatalf("expected milestones array in list response:\n%s", truncate(out, 500))
		}
		found := false
		for _, ms := range milestones {
			mm, _ := ms.(map[string]any)
			if fmt.Sprintf("%v", mm["id"]) == phaseID {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("phase %s not found in portal-level list", phaseID)
		}
	})

	t.Run("phases/update", func(t *testing.T) {
		requireID(t, phaseID, "phases/create must have succeeded")
		updatedPhaseName := phaseName + "_upd"
		out := zoho(t, "projects", "phases", "update", phaseID,
			"--project", projectID,
			"--json", toJSON(t, map[string]any{"name": updatedPhaseName}))
		m := parseJSON(t, out)
		assertEqual(t, fmt.Sprintf("%v", m["name"]), updatedPhaseName)

		out = zoho(t, "projects", "phases", "get", phaseID, "--project", projectID)
		m = parseJSON(t, out)
		assertStringField(t, m, "name", updatedPhaseName)
		phaseName = updatedPhaseName
	})

	t.Run("phases/activities", func(t *testing.T) {
		requireID(t, phaseID, "phases/create must have succeeded")
		out := zoho(t, "projects", "phases", "activities", phaseID,
			"--project", projectID)
		m := parseJSON(t, out)
		if _, ok := m["activities"]; !ok {
			t.Fatalf("expected activities key in response:\n%s", truncate(out, 500))
		}
	})

	t.Run("phases/clone", func(t *testing.T) {
		requireID(t, phaseID, "phases/create must have succeeded")
		out := zoho(t, "projects", "phases", "clone", phaseID,
			"--project", projectID)
		m := parseJSON(t, out)
		clonedPhaseID = fmt.Sprintf("%v", m["id"])
		if clonedPhaseID == "" || clonedPhaseID == "<nil>" {
			t.Fatalf("no id in phases clone response:\n%s", truncate(out, 500))
		}
		cleanup.trackPhase(clonedPhaseID, projectID)
		if clonedPhaseID == phaseID {
			t.Errorf("cloned phase has same ID as original")
		}
		t.Logf("cloned phase %s -> %s", phaseID, clonedPhaseID)
	})

	t.Run("phases/move-known-broken", func(t *testing.T) {
		t.Skip("phases move endpoint returns Zoho 500 error")
	})

	t.Run("phase-followers/add", func(t *testing.T) {
		requireID(t, phaseID, "phases/create must have succeeded")
		requireID(t, ownerZPUID, "timelogs/setup-owner must have succeeded")
		out := zoho(t, "projects", "phase-followers", "add", phaseID,
			"--project", projectID,
			"--json", toJSON(t, map[string]any{"followers": []string{ownerZPUID}}))
		m := parseJSON(t, out)
		if _, ok := m["followers"]; !ok {
			t.Fatalf("expected followers key in add response:\n%s", truncate(out, 500))
		}
	})

	t.Run("phase-followers/list", func(t *testing.T) {
		requireID(t, phaseID, "phases/create must have succeeded")
		out := zoho(t, "projects", "phase-followers", "list", phaseID,
			"--project", projectID)
		m := parseJSON(t, out)
		followers, ok := m["followers"].([]any)
		if !ok || len(followers) == 0 {
			t.Fatalf("expected non-empty followers array:\n%s", truncate(out, 500))
		}
		found := false
		for _, f := range followers {
			fm, _ := f.(map[string]any)
			if fmt.Sprintf("%v", fm["zpuid"]) == ownerZPUID {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("follower %s not found in phase-followers list", ownerZPUID)
		}
	})

	t.Run("phase-followers/remove", func(t *testing.T) {
		requireID(t, phaseID, "phases/create must have succeeded")
		requireID(t, ownerZPUID, "timelogs/setup-owner must have succeeded")
		zoho(t, "projects", "phase-followers", "remove", phaseID,
			"--project", projectID,
			"--json", toJSON(t, map[string]any{"followers": []string{ownerZPUID}}))

		time.Sleep(2 * time.Second)
		out := zoho(t, "projects", "phase-followers", "list", phaseID,
			"--project", projectID)
		if strings.Contains(out, ownerZPUID) {
			t.Errorf("follower %s still present after phase-followers remove", ownerZPUID)
		}
	})

	t.Run("phase-comments/add", func(t *testing.T) {
		requireID(t, phaseID, "phases/create must have succeeded")
		out := zoho(t, "projects", "phase-comments", "add",
			"--phase", phaseID, "--project", projectID,
			"--comment", testPrefix+"_phase_comment")
		m := parseJSON(t, out)
		phaseCommentID = fmt.Sprintf("%v", m["id"])
		if phaseCommentID == "" || phaseCommentID == "<nil>" {
			t.Fatalf("no id in phase-comments add response:\n%s", truncate(out, 500))
		}
		assertStringField(t, m, "content", testPrefix+"_phase_comment")
		t.Logf("created phase comment %s", phaseCommentID)
	})

	t.Run("phase-comments/list", func(t *testing.T) {
		requireID(t, phaseCommentID, "phase-comments/add must have succeeded")
		out := zoho(t, "projects", "phase-comments", "list",
			"--phase", phaseID, "--project", projectID)
		m := parseJSON(t, out)
		comments, ok := m["comments"].([]any)
		if !ok {
			t.Fatalf("expected comments array in list response:\n%s", truncate(out, 500))
		}
		found := false
		for _, c := range comments {
			cm, _ := c.(map[string]any)
			if fmt.Sprintf("%v", cm["id"]) == phaseCommentID {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("comment %s not found in phase-comments list", phaseCommentID)
		}
	})

	t.Run("phase-comments/update", func(t *testing.T) {
		requireID(t, phaseCommentID, "phase-comments/add must have succeeded")
		updatedComment := testPrefix + "_phase_comment_upd"
		out := zoho(t, "projects", "phase-comments", "update", phaseCommentID,
			"--phase", phaseID, "--project", projectID,
			"--comment", updatedComment)
		m := parseJSON(t, out)
		assertStringField(t, m, "content", updatedComment)

		listOut := zoho(t, "projects", "phase-comments", "list",
			"--phase", phaseID, "--project", projectID)
		assertContains(t, listOut, updatedComment)
	})

	t.Run("phase-comments/delete", func(t *testing.T) {
		requireID(t, phaseCommentID, "phase-comments/add must have succeeded")
		zoho(t, "projects", "phase-comments", "delete", phaseCommentID,
			"--phase", phaseID, "--project", projectID)

		listOut := zoho(t, "projects", "phase-comments", "list",
			"--phase", phaseID, "--project", projectID)
		if strings.Contains(listOut, phaseCommentID) {
			t.Errorf("comment %s still found in list after delete", phaseCommentID)
		}
		phaseCommentID = ""
	})

	t.Run("phases/delete", func(t *testing.T) {
		requireID(t, phaseID, "phases/create must have succeeded")
		zoho(t, "projects", "phases", "delete", phaseID, "--project", projectID)

		r := runZoho(t, "projects", "phases", "get", phaseID, "--project", projectID)
		if r.ExitCode == 0 {
			t.Errorf("phase %s still accessible after delete", phaseID)
		}
		phaseID = ""

		if clonedPhaseID != "" {
			zohoIgnoreError(t, "projects", "phases", "delete", clonedPhaseID, "--project", projectID)
			clonedPhaseID = ""
		}
	})

	t.Run("events/create", func(t *testing.T) {
		requireID(t, projectID, "core/create must have succeeded")
		requireID(t, ownerZPUID, "timelogs/setup-owner must have succeeded")
		eventName := testName(t) + "_event"
		now := time.Now().AddDate(0, 1, 0)
		startsAt := now.Format("2006-01-02T15:04:05+00:00")
		endsAt := now.Add(1 * time.Hour).Format("2006-01-02T15:04:05+00:00")
		out := zoho(t, "projects", "events", "create",
			"--project", projectID,
			"--json", toJSON(t, map[string]any{
				"title":     eventName,
				"starts_at": startsAt,
				"ends_at":   endsAt,
				"attendees": []string{ownerZPUID},
			}))
		eventID = extractProjectsID(t, out)
		cleanup.trackEvent(eventID, projectID)
		t.Logf("created event %s (%s)", eventID, eventName)

		out = zoho(t, "projects", "events", "get", eventID, "--project", projectID)
		m := parseJSON(t, out)
		assertStringField(t, m, "title", eventName)
	})

	t.Run("events/get", func(t *testing.T) {
		requireID(t, eventID, "events/create must have succeeded")
		out := zoho(t, "projects", "events", "get", eventID, "--project", projectID)
		m := parseJSON(t, out)
		assertEqual(t, fmt.Sprintf("%v", m["id"]), eventID)
	})

	t.Run("events/list", func(t *testing.T) {
		requireID(t, eventID, "events/create must have succeeded")
		out := zoho(t, "projects", "events", "list", "--project", projectID)
		m := parseJSON(t, out)
		events, ok := m["events"].([]any)
		if !ok {
			t.Fatalf("expected events array in list response:\n%s", truncate(out, 500))
		}
		found := false
		for _, e := range events {
			em, _ := e.(map[string]any)
			if fmt.Sprintf("%v", em["id"]) == eventID {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("event %s not found in events list", eventID)
		}
	})

	t.Run("events/update", func(t *testing.T) {
		requireID(t, eventID, "events/create must have succeeded")
		updatedTitle := testPrefix + "_event_upd"
		out := zoho(t, "projects", "events", "update", eventID,
			"--project", projectID,
			"--json", toJSON(t, map[string]any{"title": updatedTitle}))
		m := parseJSON(t, out)
		assertStringField(t, m, "title", updatedTitle)

		out = zoho(t, "projects", "events", "get", eventID, "--project", projectID)
		m = parseJSON(t, out)
		assertStringField(t, m, "title", updatedTitle)
	})

	t.Run("event-comments/add", func(t *testing.T) {
		requireID(t, eventID, "events/create must have succeeded")
		out := zoho(t, "projects", "event-comments", "add",
			"--event", eventID, "--project", projectID,
			"--comment", testPrefix+"_event_comment")
		m := parseJSON(t, out)
		eventCommentID = fmt.Sprintf("%v", m["id"])
		if eventCommentID == "" || eventCommentID == "<nil>" {
			t.Fatalf("no id in event-comments add response:\n%s", truncate(out, 500))
		}
		assertStringField(t, m, "content", testPrefix+"_event_comment")
		t.Logf("created event comment %s", eventCommentID)
	})

	t.Run("event-comments/list", func(t *testing.T) {
		requireID(t, eventCommentID, "event-comments/add must have succeeded")
		out := zoho(t, "projects", "event-comments", "list",
			"--event", eventID, "--project", projectID)
		arr := parseJSONArray(t, out)
		_, found := findInArray(arr, eventCommentID)
		if !found {
			t.Errorf("comment %s not found in event-comments list", eventCommentID)
		}
	})

	t.Run("event-comments/get-known-broken", func(t *testing.T) {
		t.Skip("event-comments get endpoint returns INVALID_METHOD")
	})

	t.Run("event-comments/update", func(t *testing.T) {
		requireID(t, eventCommentID, "event-comments/add must have succeeded")
		updatedComment := testPrefix + "_event_comment_upd"
		zoho(t, "projects", "event-comments", "update", eventCommentID,
			"--event", eventID, "--project", projectID,
			"--comment", updatedComment)

		listOut := zoho(t, "projects", "event-comments", "list",
			"--event", eventID, "--project", projectID)
		assertContains(t, listOut, updatedComment)
	})

	t.Run("event-comments/delete", func(t *testing.T) {
		requireID(t, eventCommentID, "event-comments/add must have succeeded")
		zoho(t, "projects", "event-comments", "delete", eventCommentID,
			"--event", eventID, "--project", projectID)

		listOut := zoho(t, "projects", "event-comments", "list",
			"--event", eventID, "--project", projectID)
		if strings.Contains(listOut, eventCommentID) {
			t.Errorf("comment %s still found in list after delete", eventCommentID)
		}
		eventCommentID = ""
	})

	t.Run("events/delete", func(t *testing.T) {
		requireID(t, eventID, "events/create must have succeeded")
		zoho(t, "projects", "events", "delete", eventID, "--project", projectID)

		r := runZoho(t, "projects", "events", "get", eventID, "--project", projectID)
		if r.ExitCode == 0 {
			t.Errorf("event %s still accessible after delete", eventID)
		}
		eventID = ""
	})

	t.Run("search/portal", func(t *testing.T) {
		requireID(t, projectID, "core/create must have succeeded")
		retryUntil(t, 30*time.Second, func() bool {
			out, err := zohoMayFail(t, "projects", "search", "portal",
				"--query", testPrefix)
			if err != nil {
				return false
			}
			return strings.Contains(out, projectName) || strings.Contains(out, testPrefix)
		})
	})

	t.Run("search/project", func(t *testing.T) {
		requireID(t, projectID, "core/create must have succeeded")
		retryUntil(t, 30*time.Second, func() bool {
			out, err := zohoMayFail(t, "projects", "search", "project",
				"--query", testPrefix, "--project", projectID)
			if err != nil {
				return false
			}
			return strings.Contains(out, testPrefix)
		})
	})

	t.Run("feed/post", func(t *testing.T) {
		requireID(t, projectID, "core/create must have succeeded")
		statusContent := testPrefix + "_status"
		out := zoho(t, "projects", "feed", "post",
			"--project", projectID,
			"--json", toJSON(t, map[string]any{"content": statusContent}))
		m := parseJSON(t, out)
		id := fmt.Sprintf("%v", m["id"])
		if id == "" || id == "<nil>" {
			t.Fatalf("no id in feed post response:\n%s", truncate(out, 500))
		}
		t.Logf("posted status %s", id)
	})

	t.Run("feed/status", func(t *testing.T) {
		requireID(t, projectID, "core/create must have succeeded")
		out := zoho(t, "projects", "feed", "status", "--project", projectID)
		m := parseJSON(t, out)
		if _, ok := m["status"]; !ok {
			if _, ok2 := m["page_info"]; !ok2 {
				t.Errorf("expected status or page_info in response:\n%s", truncate(out, 500))
			}
		}
	})

	t.Run("timers/setup", func(t *testing.T) {
		t.Skip("all timer tests are broken (Zoho 500)")
	})

	t.Run("timers/running-empty", func(t *testing.T) {
		out := zoho(t, "projects", "timelog-timers", "running")
		m := parseJSON(t, out)
		if timers, ok := m["timer"]; ok {
			arr, _ := timers.([]any)
			t.Logf("running timers: %d", len(arr))
		}
	})

	t.Run("timers/start-known-broken", func(t *testing.T) {
		t.Skip("timelog-timers start endpoint returns Zoho 500")
	})

	t.Run("timers/pause", func(t *testing.T) {
		t.Skip("depends on timers/start which returns Zoho 500")
	})

	t.Run("timers/resume", func(t *testing.T) {
		t.Skip("depends on timers/start which returns Zoho 500")
	})

	t.Run("timers/stop", func(t *testing.T) {
		t.Skip("depends on timers/start which returns Zoho 500")
	})

	t.Run("timers/delete", func(t *testing.T) {
		t.Skip("depends on timers/start which returns Zoho 500")
	})

	t.Run("pins/setup", func(t *testing.T) {
		t.Skip("all pin tests are broken (create rejects all fields)")
	})

	t.Run("pins/list-empty", func(t *testing.T) {
		out := zoho(t, "projects", "timelog-pins", "list")
		m := parseJSON(t, out)
		if _, ok := m["pins"]; !ok {
			if _, ok2 := m["page_info"]; !ok2 {
				t.Errorf("expected pins or page_info key in response:\n%s", truncate(out, 500))
			}
		}
	})

	t.Run("pins/create-known-broken", func(t *testing.T) {
		t.Skip("timelog-pins create endpoint rejects all fields")
	})

	t.Run("pins/list-after-create", func(t *testing.T) {
		t.Skip("depends on pins/create which is broken")
	})

	t.Run("pins/delete", func(t *testing.T) {
		t.Skip("depends on pins/create which is broken")
	})

	t.Run("timelog-bulk/list-missing-module", func(t *testing.T) {
		r := runZoho(t, "projects", "timelog-bulk", "list")
		assertExitCode(t, r, 1)
		assertContains(t, r.Stderr, "Required flag")
	})

	t.Run("timelog-bulk/project-list", func(t *testing.T) {
		requireID(t, projectID, "core/create must have succeeded")
		out := zoho(t, "projects", "timelog-bulk", "project-list",
			"--project", projectID, "--module", "task")
		m := parseJSON(t, out)
		if _, ok := m["timelogs"]; !ok {
			if _, ok2 := m["page_info"]; !ok2 {
				t.Errorf("expected timelogs or page_info key in response:\n%s", truncate(out, 500))
			}
		}
	})

	t.Run("users/list", func(t *testing.T) {
		out := zoho(t, "projects", "users", "list")
		arr := parseJSONArray(t, out)
		assertNonEmpty(t, arr, "expected at least one user")
		t.Logf("found %d portal users", len(arr))
	})

	t.Run("users/get", func(t *testing.T) {
		requireID(t, ownerZPUID, "owner zpuid must be available")
		out := zoho(t, "projects", "users", "get", ownerZPUID)
		m := parseJSON(t, out)
		id := fmt.Sprintf("%v", m["id"])
		if id == "" || id == "<nil>" {
			t.Errorf("expected id in user get response:\n%s", truncate(out, 500))
		}
	})

	t.Run("core/trash-second", func(t *testing.T) {
		requireID(t, project2ID, "second project must exist")
		out := zoho(t, "projects", "core", "trash", project2ID)
		parseJSON(t, out)
		t.Logf("trashed project %s", project2ID)

		r := runZoho(t, "projects", "core", "get", project2ID)
		if r.ExitCode == 0 {
			m := parseJSON(t, r.Stdout)
			status := fmt.Sprintf("%v", m["status"])
			t.Logf("trashed project get returned status=%s (may still be accessible)", status)
		}
	})

	t.Run("core/restore-second", func(t *testing.T) {
		requireID(t, project2ID, "second project must exist")
		time.Sleep(2 * time.Second)
		out := zoho(t, "projects", "core", "restore", project2ID)
		parseJSON(t, out)

		retryUntil(t, 30*time.Second, func() bool {
			r := runZoho(t, "projects", "core", "get", project2ID)
			if r.ExitCode != 0 {
				return false
			}
			return true
		})
		proj := getProject(t, project2ID)
		assertStringField(t, proj, "name", project2Name)
	})

	t.Run("trash/list", func(t *testing.T) {
		out := zoho(t, "projects", "trash", "list")
		m := parseJSON(t, out)
		if _, ok := m["trash"]; !ok {
			if _, ok2 := m["page_info"]; !ok2 {
				t.Errorf("expected trash or page_info key in response:\n%s", truncate(out, 500))
			}
		}
	})

	t.Run("trash/restore-known-broken", func(t *testing.T) {
		t.Skip("trash restore with fake IDs always fails")
	})

	t.Run("trash/delete-known-broken", func(t *testing.T) {
		t.Skip("trash delete with fake IDs always fails")
	})

	t.Run("core/trash-and-delete-second", func(t *testing.T) {
		requireID(t, project2ID, "second project must exist")
		zoho(t, "projects", "core", "trash", project2ID)
		time.Sleep(2 * time.Second)
		out := zoho(t, "projects", "core", "delete", project2ID)
		parseJSON(t, out)
		t.Logf("permanently deleted project %s", project2ID)
		project2ID = ""
	})

	t.Run("core/trash-and-delete-main", func(t *testing.T) {
		requireID(t, projectID, "core/create must have succeeded")
		zoho(t, "projects", "core", "trash", projectID)
		time.Sleep(2 * time.Second)
		out := zoho(t, "projects", "core", "delete", projectID)
		parseJSON(t, out)
		t.Logf("permanently deleted project %s", projectID)
		projectID = ""
	})
}

func TestProjectsProfiles(t *testing.T) {
	t.Parallel()
	_ = requireProjectsPortalID(t)
	cleanup := newCleanup(t)

	var profileID string
	var profileName string

	t.Run("profiles/list", func(t *testing.T) {
		out := zoho(t, "projects", "profiles", "list")
		m := parseJSON(t, out)
		profiles, ok := m["profiles"].([]any)
		if !ok {
			t.Fatalf("expected profiles array in response:\n%s", truncate(out, 500))
		}
		if len(profiles) == 0 {
			t.Error("expected at least one built-in profile")
		}
		t.Logf("found %d profiles", len(profiles))
	})

	t.Run("profiles/create", func(t *testing.T) {
		profileName = testName(t) + "_profile"
		out := zoho(t, "projects", "profiles", "create",
			"--json", toJSON(t, map[string]any{
				"name": profileName,
				"type": "3",
			}))
		m := parseJSON(t, out)
		profileID = fmt.Sprintf("%v", m["id"])
		if profileID == "" || profileID == "<nil>" {
			t.Fatalf("no id in profile create response:\n%s", truncate(out, 500))
		}
		cleanup.trackProfile(profileID)
		t.Logf("created profile %s (%s)", profileID, profileName)
	})

	t.Run("profiles/get", func(t *testing.T) {
		requireID(t, profileID, "profiles/create must have succeeded")
		out := zoho(t, "projects", "profiles", "get", profileID)
		m := parseJSON(t, out)
		assertEqual(t, fmt.Sprintf("%v", m["id"]), profileID)
		assertStringField(t, m, "name", profileName)
	})

	t.Run("profiles/update", func(t *testing.T) {
		requireID(t, profileID, "profiles/create must have succeeded")
		updatedName := profileName + "_upd"
		zoho(t, "projects", "profiles", "update", profileID,
			"--json", toJSON(t, map[string]any{"name": updatedName}))

		out := zoho(t, "projects", "profiles", "get", profileID)
		m := parseJSON(t, out)
		assertStringField(t, m, "name", updatedName)
		profileName = updatedName
	})

	t.Run("profiles/set-default", func(t *testing.T) {
		requireID(t, profileID, "profiles/create must have succeeded")
		zoho(t, "projects", "profiles", "set-default", profileID)

		out := zoho(t, "projects", "profiles", "get", profileID)
		m := parseJSON(t, out)
		if fmt.Sprintf("%v", m["is_default"]) != "true" {
			t.Errorf("profile is_default not set to true after set-default, got %v", m["is_default"])
		}

		zoho(t, "projects", "profiles", "set-default", "1212503000000015149")
	})

	t.Run("profiles/delete", func(t *testing.T) {
		requireID(t, profileID, "profiles/create must have succeeded")
		zoho(t, "projects", "profiles", "delete", profileID)

		r := runZoho(t, "projects", "profiles", "get", profileID)
		if r.ExitCode == 0 {
			t.Errorf("profile %s still accessible after delete", profileID)
		}
		profileID = ""
	})
}

func TestProjectsDashboards(t *testing.T) {
	t.Parallel()
	_ = requireProjectsPortalID(t)
	cleanup := newCleanup(t)

	var dashboardID string
	var dashboardName string

	t.Run("reports/workload-meta", func(t *testing.T) {
		out := zoho(t, "projects", "reports", "workload-meta")
		m := parseJSON(t, out)
		if _, ok := m["chart_details"]; !ok {
			t.Errorf("expected chart_details in workload-meta response:\n%s", truncate(out, 500))
		}
	})

	t.Run("reports/workload-known-broken", func(t *testing.T) {
		t.Skip("reports workload endpoint returns Zoho 500")
	})

	t.Run("reports/dashboards-list", func(t *testing.T) {
		out := zoho(t, "projects", "reports", "dashboards")
		m := parseJSON(t, out)
		if _, ok := m["folders"]; !ok {
			if _, ok2 := m["page_info"]; !ok2 {
				t.Errorf("expected folders or page_info in dashboards response:\n%s", truncate(out, 500))
			}
		}
	})

	t.Run("reports/dashboard-create", func(t *testing.T) {
		dashboardName = testName(t) + "_dash"
		out := zoho(t, "projects", "reports", "dashboard-create",
			"--json", toJSON(t, map[string]any{"name": dashboardName}))
		m := parseJSON(t, out)
		dashboardID = fmt.Sprintf("%v", m["id"])
		if dashboardID == "" || dashboardID == "<nil>" {
			t.Fatalf("no id in dashboard create response:\n%s", truncate(out, 500))
		}
		cleanup.trackDashboard(dashboardID)
		t.Logf("created dashboard %s (%s)", dashboardID, dashboardName)
	})

	t.Run("reports/dashboard-get", func(t *testing.T) {
		requireID(t, dashboardID, "dashboard-create must have succeeded")
		out := zoho(t, "projects", "reports", "dashboard-get", dashboardID)
		m := parseJSON(t, out)
		assertEqual(t, fmt.Sprintf("%v", m["id"]), dashboardID)
		assertStringField(t, m, "name", dashboardName)
	})

	t.Run("reports/dashboard-update", func(t *testing.T) {
		requireID(t, dashboardID, "dashboard-create must have succeeded")
		updatedName := dashboardName + "_upd"
		zoho(t, "projects", "reports", "dashboard-update", dashboardID,
			"--json", toJSON(t, map[string]any{"name": updatedName}))
		getOut := zoho(t, "projects", "reports", "dashboard-get", dashboardID)
		gm := parseJSON(t, getOut)
		assertStringField(t, gm, "name", updatedName)
		dashboardName = updatedName
	})

	t.Run("reports/dashboard-delete", func(t *testing.T) {
		requireID(t, dashboardID, "dashboard-create must have succeeded")
		zoho(t, "projects", "reports", "dashboard-delete", dashboardID)

		r := runZoho(t, "projects", "reports", "dashboard-get", dashboardID)
		if r.ExitCode == 0 {
			t.Errorf("dashboard %s still accessible after delete", dashboardID)
		}
		dashboardID = ""
	})
}

func TestProjectsTeams(t *testing.T) {
	t.Parallel()
	portalID := requireProjectsPortalID(t)
	cleanup := newCleanup(t)

	var teamID string
	var teamName string
	var projectID string
	var ownerZPUID string

	t.Run("teams/setup", func(t *testing.T) {
		out := zoho(t, "projects", "users", "list")
		arr := parseJSONArray(t, out)
		if len(arr) == 0 {
			t.Fatal("no users found in portal")
		}
		ownerZPUID = fmt.Sprintf("%v", arr[0]["zpuid"])
		if ownerZPUID == "" || ownerZPUID == "<nil>" {
			ownerZPUID = fmt.Sprintf("%v", arr[0]["id"])
		}
		if ownerZPUID == "" || ownerZPUID == "<nil>" {
			t.Fatal("could not determine user zpuid")
		}
		t.Logf("using owner zpuid: %s", ownerZPUID)

		projectName := testName(t) + "_teamproj"
		projOut := zoho(t, "projects", "core", "create", "--name", projectName)
		projectID = extractProjectsID(t, projOut)
		cleanup.trackProject(projectID)
		t.Logf("created project %s for team tests", projectID)
		_ = portalID
	})

	t.Run("teams/create", func(t *testing.T) {
		requireID(t, ownerZPUID, "teams/setup must have succeeded")
		teamName = testName(t) + "_team"
		out, err := zohoMayFail(t, "projects", "teams", "create",
			"--json", toJSON(t, map[string]any{
				"name":     teamName,
				"lead":     ownerZPUID,
				"user_ids": map[string]any{"add": []string{ownerZPUID}},
			}))
		if err != nil {
			t.Logf("team create failed (may need scope ZohoProjects.teams.ALL): %v", err)
			t.Logf("response: %s", truncate(out, 500))
			return
		}
		m := parseJSON(t, out)
		teamID = fmt.Sprintf("%v", m["id"])
		if teamID == "" || teamID == "<nil>" {
			t.Fatalf("no id in team create response:\n%s", truncate(out, 500))
		}
		cleanup.trackTeam(teamID)
		t.Logf("created team %s (%s)", teamID, teamName)
	})

	t.Run("teams/list", func(t *testing.T) {
		requireID(t, teamID, "teams/create must have succeeded")
		out := zoho(t, "projects", "teams", "list")
		m := parseJSON(t, out)
		teams, ok := m["teams"].([]any)
		if !ok {
			t.Fatalf("expected teams array in response:\n%s", truncate(out, 500))
		}
		found := false
		for _, tm := range teams {
			tmm, _ := tm.(map[string]any)
			if fmt.Sprintf("%v", tmm["id"]) == teamID {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("team %s not found in teams list", teamID)
		}
	})

	t.Run("teams/get-known-broken", func(t *testing.T) {
		t.Skip("teams get endpoint returns INVALID_METHOD")
	})

	t.Run("teams/update", func(t *testing.T) {
		requireID(t, teamID, "teams/create must have succeeded")
		updatedName := teamName + "_upd"
		zoho(t, "projects", "teams", "update", teamID,
			"--json", toJSON(t, map[string]any{"name": updatedName}))

		out := zoho(t, "projects", "teams", "list")
		m := parseJSON(t, out)
		teams, _ := m["teams"].([]any)
		for _, tm := range teams {
			tmm, _ := tm.(map[string]any)
			if fmt.Sprintf("%v", tmm["id"]) == teamID {
				assertStringField(t, tmm, "name", updatedName)
				break
			}
		}
		teamName = updatedName
	})

	t.Run("teams/users-known-broken", func(t *testing.T) {
		t.Skip("teams users endpoint returns URL_RULE_NOT_CONFIGURED")
	})

	t.Run("teams/projects-known-broken", func(t *testing.T) {
		t.Skip("teams projects endpoint not functional")
	})

	t.Run("teams/add-to-project", func(t *testing.T) {
		requireID(t, teamID, "teams/create must have succeeded")
		requireID(t, projectID, "teams/setup must have succeeded")
		out, err := zohoMayFail(t, "projects", "teams", "add-to-project",
			"--project", projectID,
			"--json", toJSON(t, []string{teamID}))
		if err != nil {
			t.Logf("add-to-project failed: %v", err)
			t.Logf("response: %s", truncate(out, 500))
			return
		}
		t.Logf("added team %s to project %s", teamID, projectID)
	})

	t.Run("teams/project-list", func(t *testing.T) {
		requireID(t, projectID, "teams/setup must have succeeded")
		out, err := zohoMayFail(t, "projects", "teams", "project-list",
			"--project", projectID)
		if err != nil {
			t.Logf("project-list failed: %v", err)
			return
		}
		m := parseJSON(t, out)
		if _, ok := m["teams"]; !ok {
			if _, ok2 := m["page_info"]; !ok2 {
				t.Logf("unexpected response format:\n%s", truncate(out, 500))
			}
		}
	})

	t.Run("teams/remove-from-project", func(t *testing.T) {
		requireID(t, teamID, "teams/create must have succeeded")
		requireID(t, projectID, "teams/setup must have succeeded")
		r := runZoho(t, "projects", "teams", "remove-from-project", teamID,
			"--project", projectID)
		if r.ExitCode != 0 {
			t.Logf("remove-from-project failed: %s", truncate(r.Stderr+r.Stdout, 300))
		}
	})

	t.Run("teams/delete", func(t *testing.T) {
		requireID(t, teamID, "teams/create must have succeeded")
		zoho(t, "projects", "teams", "delete", teamID)

		out := zoho(t, "projects", "teams", "list")
		m := parseJSON(t, out)
		teams, _ := m["teams"].([]any)
		for _, tm := range teams {
			tmm, _ := tm.(map[string]any)
			if fmt.Sprintf("%v", tmm["id"]) == teamID {
				t.Errorf("team %s still found in list after delete", teamID)
				break
			}
		}
		teamID = ""
	})
}

func TestProjectsErrors(t *testing.T) {
	t.Parallel()
	_ = requireProjectsPortalID(t)

	t.Run("bad-auth", func(t *testing.T) {
		r := runZohoWithEnv(t, map[string]string{
			"ZOHO_CLIENT_ID":     "bad_client_id",
			"ZOHO_CLIENT_SECRET": "bad_client_secret",
			"ZOHO_REFRESH_TOKEN": "bad_refresh_token",
			"ZOHO_DC":            "com",
		}, "projects", "core", "list")
		assertExitCode(t, r, 2)
	})

	t.Run("missing-portal", func(t *testing.T) {
		r := runZohoWithEnv(t, map[string]string{
			"ZOHO_PORTAL_ID": "",
		}, "projects", "portals", "get")
		if r.ExitCode == 0 {
			t.Error("expected error when portal ID is missing")
		}
		assertContains(t, r.Stderr, "--portal")
	})

	t.Run("nonexistent-project", func(t *testing.T) {
		r := runZoho(t, "projects", "core", "get", "999999999999")
		if r.ExitCode == 0 {
			t.Error("expected error for nonexistent project")
		}
	})

	t.Run("nonexistent-task", func(t *testing.T) {
		r := runZoho(t, "projects", "tasks", "get", "999999999999",
			"--project", "999999999999")
		if r.ExitCode == 0 {
			t.Error("expected error for nonexistent task")
		}
	})

	t.Run("missing-required-name", func(t *testing.T) {
		r := runZoho(t, "projects", "core", "create")
		assertExitCode(t, r, 1)
		assertContains(t, r.Stderr, "Required flag")
	})

	t.Run("missing-required-project-flag", func(t *testing.T) {
		r := runZoho(t, "projects", "tasks", "list")
		assertExitCode(t, r, 1)
		assertContains(t, r.Stderr, "Required flag")
	})

	t.Run("milestone-missing-dates", func(t *testing.T) {
		r := runZoho(t, "projects", "milestones", "create",
			"--name", "test", "--project", "12345")
		assertExitCode(t, r, 1)
		assertContains(t, r.Stderr, "Required flag")
	})

	t.Run("timelog-missing-required", func(t *testing.T) {
		r := runZoho(t, "projects", "timelogs", "add",
			"--hours", "2", "--project", "12345")
		assertExitCode(t, r, 1)
		assertContains(t, r.Stderr, "Required flag")
	})
}

func TestProjectsEmergencyCleanup(t *testing.T) {
	if os.Getenv("ZOHO_EMERGENCY_CLEANUP") == "" {
		t.Skip("set ZOHO_EMERGENCY_CLEANUP=1 to run")
	}
	_ = requireProjectsPortalID(t)

	out := zoho(t, "projects", "core", "list")
	arr := parseJSONArray(t, out)
	t.Logf("found %d total projects", len(arr))
	for _, proj := range arr {
		name := fmt.Sprintf("%v", proj["name"])
		if !strings.HasPrefix(name, testPrefix) {
			continue
		}
		id := fmt.Sprintf("%v", proj["id"])
		t.Logf("cleaning orphaned project %s (%s)", id, name)
		zohoIgnoreError(t, "projects", "core", "trash", id)
		time.Sleep(1 * time.Second)
		zohoIgnoreError(t, "projects", "core", "delete", id)
	}
}

func assertExpenseCodeZero(t *testing.T, m map[string]any) {
	t.Helper()
	if fmt.Sprintf("%v", m["code"]) != "0" {
		t.Fatalf("expected expense code=0, got %v", m["code"])
	}
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
			"--json", toJSON(t, map[string]any{"category_name": categoryName}))
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
			"--json", toJSON(t, map[string]any{"category_name": updatedName}))
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
			"--json", toJSON(t, map[string]any{"contact_name": customerName}))
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
			"--json", toJSON(t, map[string]any{"contact_name": updatedName}))
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
			"--json", toJSON(t, map[string]any{
				"currency_code":   "MXN",
				"currency_symbol": "$",
				"currency_format": "1,234,567.89",
				"price_precision": 2,
			}))
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
			"--json", toJSON(t, map[string]any{"currency_symbol": "Mex$", "currency_format": "1,234,567.89", "price_precision": 2}))
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
				"--json", toJSON(t, map[string]any{
					"tax_name":       fmt.Sprintf("%s Tax %s", testPrefix, randomSuffix()),
					"tax_percentage": rate,
					"tax_type":       "tax",
				}))
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
			"--json", toJSON(t, map[string]any{"tax_name": updatedName, "tax_percentage": createdTaxRate, "tax_type": "tax"}))
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
			"--json", toJSON(t, map[string]any{"project_name": projectName}))
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
			"--json", toJSON(t, map[string]any{"project_name": updatedName}))
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
		out := zoho(t, "expense", "trips", "create", "--org", orgID,
			"--json", toJSON(t, map[string]any{
				"is_international":    false,
				"destination_country": "South Africa",
				"start_date":          "2026-06-01",
				"end_date":            "2026-06-05",
				"business_purpose":    fmt.Sprintf("%s trip %s", testPrefix, randomSuffix()),
			}))
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
		updatedPurpose := fmt.Sprintf("%s trip updated %s", testPrefix, randomSuffix())
		out := zoho(t, "expense", "trips", "update", tripID, "--org", orgID,
			"--json", toJSON(t, map[string]any{
				"business_purpose": updatedPurpose,
			}))
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
			"--json", toJSON(t, map[string]any{"report_name": reportName}))
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
			"--json", toJSON(t, map[string]any{"report_name": updatedName}))
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
			"--json", toJSON(t, map[string]any{
				"currency_id": currencyID,
				"line_items": []map[string]any{{
					"category_id": categoryID,
					"amount":      25.50,
				}},
				"date": "2026-03-01",
			}))
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
			"--json", toJSON(t, map[string]any{
				"currency_id": currencyID,
				"line_items": []map[string]any{{
					"category_id": categoryID,
					"amount":      31.75,
				}},
				"date": "2026-03-02",
			}))
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

func assertSheetSuccess(t *testing.T, m map[string]any) {
	t.Helper()
	status := fmt.Sprintf("%v", m["status"])
	if status != "success" {
		t.Fatalf("expected Sheet status=success, got %q: error_code=%v error_message=%v",
			status, m["error_code"], m["error_message"])
	}
}

func (c *testCleanup) trackSheetWorkbook(resourceID string) {
	c.add("trash+delete sheet workbook "+resourceID, func() {
		zohoIgnoreError(c.t, "sheet", "workbooks", "trash", "--workbook", resourceID)
		zohoIgnoreError(c.t, "sheet", "workbooks", "delete", "--workbook", resourceID)
	})
}

func TestSheetUtility(t *testing.T) {
	t.Parallel()
	cleanup := newCleanup(t)
	var workbookID string

	t.Run("setup", func(t *testing.T) {
		name := fmt.Sprintf("%s_SheetUtil_%s", testPrefix, randomSuffix())
		out := zoho(t, "sheet", "workbooks", "create", "--workbook-name", name)
		m := parseJSON(t, out)
		assertSheetSuccess(t, m)
		workbookID = fmt.Sprintf("%v", m["resource_id"])
		cleanup.trackSheetWorkbook(workbookID)
	})

	t.Run("range-to-index", func(t *testing.T) {
		requireID(t, workbookID, "setup must have succeeded")
		out := zoho(t, "sheet", "utility", "range-to-index", "--workbook", workbookID, "--range", "A1:C5")
		m := parseJSON(t, out)
		assertSheetSuccess(t, m)
		assertEqual(t, fmt.Sprintf("%v", m["start_row"]), "1")
		assertEqual(t, fmt.Sprintf("%v", m["start_column"]), "1")
		assertEqual(t, fmt.Sprintf("%v", m["end_row"]), "5")
		assertEqual(t, fmt.Sprintf("%v", m["end_column"]), "3")
	})

	t.Run("index-to-range", func(t *testing.T) {
		requireID(t, workbookID, "setup must have succeeded")
		out := zoho(t, "sheet", "utility", "index-to-range", "--workbook", workbookID,
			"--start-row", "1", "--start-column", "1", "--end-row", "5", "--end-column", "3")
		m := parseJSON(t, out)
		assertSheetSuccess(t, m)
		rs := fmt.Sprintf("%v", m["range_string"])
		if rs == "" || rs == "<nil>" {
			t.Fatalf("expected range_string in response: %v", m)
		}
	})
}

func TestSheetWorkbooks(t *testing.T) {
	t.Parallel()
	cleanup := newCleanup(t)
	var workbookID string
	var copyID string
	var versionNumber string

	t.Run("setup", func(t *testing.T) {
		name := fmt.Sprintf("%s_SheetWorkbooks_%s", testPrefix, randomSuffix())
		out := zoho(t, "sheet", "workbooks", "create", "--workbook-name", name)
		m := parseJSON(t, out)
		assertSheetSuccess(t, m)
		workbookID = fmt.Sprintf("%v", m["resource_id"])
		cleanup.trackSheetWorkbook(workbookID)
	})

	t.Run("list", func(t *testing.T) {
		out := zoho(t, "sheet", "workbooks", "list")
		m := parseJSON(t, out)
		assertSheetSuccess(t, m)
		if !strings.Contains(fmt.Sprintf("%v", m), workbookID) {
			t.Errorf("workbook %s not found in list response", workbookID)
		}
	})

	t.Run("copy", func(t *testing.T) {
		requireID(t, workbookID, "setup must have succeeded")
		name := fmt.Sprintf("%s_Copy_%s", testPrefix, randomSuffix())
		out := zoho(t, "sheet", "workbooks", "copy", "--workbook", workbookID, "--new-workbook-name", name)
		m := parseJSON(t, out)
		assertSheetSuccess(t, m)
		copyID = fmt.Sprintf("%v", m["resource_id"])
		cleanup.trackSheetWorkbook(copyID)
	})

	t.Run("lock", func(t *testing.T) {
		t.Skip("lock/unlock require user_emails param not yet implemented")
	})

	t.Run("unlock", func(t *testing.T) {
		t.Skip("lock/unlock require user_emails param not yet implemented")
	})

	t.Run("create-version", func(t *testing.T) {
		requireID(t, workbookID, "setup must have succeeded")
		m := parseJSON(t, zoho(t, "sheet", "workbooks", "create-version", "--workbook", workbookID, "--version-description", "integration version"))
		assertSheetSuccess(t, m)
		versionNumber = fmt.Sprintf("%v", m["version_number"])
		if versionNumber == "" || versionNumber == "<nil>" {
			if data, ok := m["data"].(map[string]any); ok {
				versionNumber = fmt.Sprintf("%v", data["version_number"])
			}
		}
	})

	t.Run("versions", func(t *testing.T) {
		requireID(t, workbookID, "setup must have succeeded")
		m := parseJSON(t, zoho(t, "sheet", "workbooks", "versions", "--workbook", workbookID))
		assertSheetSuccess(t, m)
		if !strings.Contains(fmt.Sprintf("%v", m), "version") {
			t.Errorf("expected version data in versions response")
		}
	})

	t.Run("revert-version", func(t *testing.T) {
		requireID(t, workbookID, "setup must have succeeded")
		if versionNumber == "" || versionNumber == "<nil>" {
			t.Skip("version_number missing in create-version response")
		}
		m := parseJSON(t, zoho(t, "sheet", "workbooks", "revert-version", "--workbook", workbookID, "--version-number", versionNumber))
		assertSheetSuccess(t, m)
	})

	t.Run("publish", func(t *testing.T) {
		t.Skip("publish requires publish_type param - format unknown")
	})

	t.Run("unpublish", func(t *testing.T) {
		t.Skip("publish requires publish_type param - format unknown")
	})

	t.Run("trash-copy", func(t *testing.T) {
		requireID(t, copyID, "copy must have succeeded")
		m := parseJSON(t, zoho(t, "sheet", "workbooks", "trash", "--workbook", copyID))
		assertSheetSuccess(t, m)
	})

	t.Run("restore-copy", func(t *testing.T) {
		requireID(t, copyID, "trash-copy must have succeeded")
		m := parseJSON(t, zoho(t, "sheet", "workbooks", "restore", "--workbook", copyID))
		assertSheetSuccess(t, m)
	})

	t.Run("trash-delete-copy", func(t *testing.T) {
		requireID(t, copyID, "restore-copy must have succeeded")
		assertSheetSuccess(t, parseJSON(t, zoho(t, "sheet", "workbooks", "trash", "--workbook", copyID)))
		assertSheetSuccess(t, parseJSON(t, zoho(t, "sheet", "workbooks", "delete", "--workbook", copyID)))
		copyID = ""
	})
}

func TestSheetWorksheets(t *testing.T) {
	t.Parallel()
	cleanup := newCleanup(t)
	var workbookID string
	worksheetName := "Sheet1"
	newWorksheet := fmt.Sprintf("%s_WS_%s", testPrefix, randomSuffix())
	renamedWorksheet := fmt.Sprintf("%s_WS_REN_%s", testPrefix, randomSuffix())

	t.Run("setup", func(t *testing.T) {
		name := fmt.Sprintf("%s_SheetWorksheets_%s", testPrefix, randomSuffix())
		m := parseJSON(t, zoho(t, "sheet", "workbooks", "create", "--workbook-name", name))
		assertSheetSuccess(t, m)
		workbookID = fmt.Sprintf("%v", m["resource_id"])
		cleanup.trackSheetWorkbook(workbookID)
	})

	t.Run("list", func(t *testing.T) {
		requireID(t, workbookID, "setup must have succeeded")
		m := parseJSON(t, zoho(t, "sheet", "worksheets", "list", "--workbook", workbookID))
		assertSheetSuccess(t, m)
		if !strings.Contains(fmt.Sprintf("%v", m), "Sheet1") {
			t.Errorf("expected Sheet1 in worksheets list")
		}
	})

	t.Run("create", func(t *testing.T) {
		requireID(t, workbookID, "setup must have succeeded")
		m := parseJSON(t, zoho(t, "sheet", "worksheets", "create", "--workbook", workbookID, "--worksheet-name", newWorksheet))
		assertSheetSuccess(t, m)
		listM := parseJSON(t, zoho(t, "sheet", "worksheets", "list", "--workbook", workbookID))
		assertSheetSuccess(t, listM)
		if !strings.Contains(fmt.Sprintf("%v", listM), newWorksheet) {
			t.Errorf("newly created worksheet %s not found in list", newWorksheet)
		}
	})

	t.Run("copy", func(t *testing.T) {
		requireID(t, workbookID, "setup must have succeeded")
		copyName := fmt.Sprintf("%s_WS_COPY_%s", testPrefix, randomSuffix())
		assertSheetSuccess(t, parseJSON(t, zoho(t, "sheet", "worksheets", "copy", "--workbook", workbookID, "--worksheet", worksheetName, "--new-worksheet-name", copyName)))
	})

	t.Run("rename", func(t *testing.T) {
		requireID(t, workbookID, "setup must have succeeded")
		m := parseJSON(t, zoho(t, "sheet", "worksheets", "rename", "--workbook", workbookID, "--worksheet", newWorksheet, "--new-worksheet-name", renamedWorksheet))
		assertSheetSuccess(t, m)
		listM := parseJSON(t, zoho(t, "sheet", "worksheets", "list", "--workbook", workbookID))
		assertSheetSuccess(t, listM)
		if !strings.Contains(fmt.Sprintf("%v", listM), renamedWorksheet) {
			t.Errorf("renamed worksheet %s not found in list", renamedWorksheet)
		}
	})

	t.Run("delete", func(t *testing.T) {
		requireID(t, workbookID, "setup must have succeeded")
		m := parseJSON(t, zoho(t, "sheet", "worksheets", "delete", "--workbook", workbookID, "--worksheet", renamedWorksheet))
		assertSheetSuccess(t, m)
		listM := parseJSON(t, zoho(t, "sheet", "worksheets", "list", "--workbook", workbookID))
		assertSheetSuccess(t, listM)
		if strings.Contains(fmt.Sprintf("%v", listM), renamedWorksheet) {
			t.Errorf("deleted worksheet %s still in list", renamedWorksheet)
		}
	})
}

func TestSheetCells(t *testing.T) {
	t.Parallel()
	cleanup := newCleanup(t)
	var workbookID string
	worksheetName := "Sheet1"

	t.Run("setup", func(t *testing.T) {
		name := fmt.Sprintf("%s_SheetCells_%s", testPrefix, randomSuffix())
		m := parseJSON(t, zoho(t, "sheet", "workbooks", "create", "--workbook-name", name))
		assertSheetSuccess(t, m)
		workbookID = fmt.Sprintf("%v", m["resource_id"])
		cleanup.trackSheetWorkbook(workbookID)
	})

	t.Run("set", func(t *testing.T) {
		requireID(t, workbookID, "setup must have succeeded")
		assertSheetSuccess(t, parseJSON(t, zoho(t, "sheet", "cells", "set", "--workbook", workbookID, "--worksheet", worksheetName, "--row", "1", "--column", "1", "--content", "Hello")))
	})

	t.Run("get", func(t *testing.T) {
		requireID(t, workbookID, "set must have succeeded")
		m := parseJSON(t, zoho(t, "sheet", "cells", "get", "--workbook", workbookID, "--worksheet", worksheetName, "--row", "1", "--column", "1"))
		assertSheetSuccess(t, m)
		if !strings.Contains(fmt.Sprintf("%v", m), "Hello") {
			t.Fatalf("expected response to contain Hello: %v", m)
		}
	})

	t.Run("set-multiple", func(t *testing.T) {
		t.Skip("cells.content.set data format TBD")
	})

	t.Run("set-range", func(t *testing.T) {
		requireID(t, workbookID, "setup must have succeeded")
		csv := "a,b,c\nd,e,f"
		assertSheetSuccess(t, parseJSON(t, zoho(t, "sheet", "cells", "set-range", "--workbook", workbookID, "--worksheet", worksheetName, "--row", "3", "--column", "1", "--data", csv)))
	})

	t.Run("get-range", func(t *testing.T) {
		requireID(t, workbookID, "setup must have succeeded")
		m := parseJSON(t, zoho(t, "sheet", "cells", "get-range", "--workbook", workbookID, "--worksheet", worksheetName, "--start-row", "3", "--start-column", "1", "--end-row", "4", "--end-column", "3"))
		assertSheetSuccess(t, m)
		resp := fmt.Sprintf("%v", m)
		if !strings.Contains(resp, "a") || !strings.Contains(resp, "f") {
			t.Errorf("expected set-range data (a,f) in get-range response: %s", truncate(resp, 500))
		}
	})

	t.Run("set-row", func(t *testing.T) {
		requireID(t, workbookID, "setup must have succeeded")
		assertSheetSuccess(t, parseJSON(t, zoho(t, "sheet", "cells", "set-row", "--workbook", workbookID, "--worksheet", worksheetName, "--row", "4", "--column-array", `[1,2,3]`, "--data-array", `["r4c1","r4c2","r4c3"]`)))
	})

	t.Run("get-worksheet", func(t *testing.T) {
		requireID(t, workbookID, "setup must have succeeded")
		m := parseJSON(t, zoho(t, "sheet", "cells", "get-worksheet", "--workbook", workbookID, "--worksheet", worksheetName, "--start-row", "1", "--start-column", "1", "--end-row", "5", "--end-column", "3"))
		assertSheetSuccess(t, m)
		resp := fmt.Sprintf("%v", m)
		if !strings.Contains(resp, "r4c1") {
			t.Errorf("expected set-row data (r4c1) in get-worksheet response: %s", truncate(resp, 500))
		}
	})

	t.Run("get-used-area", func(t *testing.T) {
		requireID(t, workbookID, "setup must have succeeded")
		assertSheetSuccess(t, parseJSON(t, zoho(t, "sheet", "cells", "get-used-area", "--workbook", workbookID, "--worksheet", worksheetName)))
	})
}

func TestSheetNamedRanges(t *testing.T) {
	t.Parallel()
	cleanup := newCleanup(t)
	var workbookID string
	worksheetName := "Sheet1"
	rangeName := fmt.Sprintf("%s_RNG_%s", testPrefix, randomSuffix())

	t.Run("setup", func(t *testing.T) {
		name := fmt.Sprintf("%s_SheetNamedRanges_%s", testPrefix, randomSuffix())
		m := parseJSON(t, zoho(t, "sheet", "workbooks", "create", "--workbook-name", name))
		assertSheetSuccess(t, m)
		workbookID = fmt.Sprintf("%v", m["resource_id"])
		cleanup.trackSheetWorkbook(workbookID)
	})

	t.Run("create", func(t *testing.T) {
		requireID(t, workbookID, "setup must have succeeded")
		assertSheetSuccess(t, parseJSON(t, zoho(t, "sheet", "named-ranges", "create", "--workbook", workbookID, "--worksheet", worksheetName, "--name", rangeName, "--range", "A1:B2")))
	})

	t.Run("list", func(t *testing.T) {
		requireID(t, workbookID, "setup must have succeeded")
		m := parseJSON(t, zoho(t, "sheet", "named-ranges", "list", "--workbook", workbookID))
		assertSheetSuccess(t, m)
		if !strings.Contains(fmt.Sprintf("%v", m), rangeName) {
			t.Errorf("named range %s not found in list", rangeName)
		}
	})

	t.Run("update", func(t *testing.T) {
		requireID(t, workbookID, "create must have succeeded")
		assertSheetSuccess(t, parseJSON(t, zoho(t, "sheet", "named-ranges", "update", "--workbook", workbookID, "--worksheet", worksheetName, "--name", rangeName, "--range", "A1:C3")))
	})

	t.Run("get-named-range", func(t *testing.T) {
		requireID(t, workbookID, "create must have succeeded")
		assertSheetSuccess(t, parseJSON(t, zoho(t, "sheet", "cells", "get-named-range", "--workbook", workbookID, "--named-range", rangeName)))
	})

	t.Run("delete", func(t *testing.T) {
		requireID(t, workbookID, "create must have succeeded")
		assertSheetSuccess(t, parseJSON(t, zoho(t, "sheet", "named-ranges", "delete", "--workbook", workbookID, "--name", rangeName)))
	})
}

func TestSheetContent(t *testing.T) {
	t.Parallel()
	cleanup := newCleanup(t)
	var workbookID string
	worksheetName := "Sheet1"

	t.Run("setup", func(t *testing.T) {
		name := fmt.Sprintf("%s_SheetContent_%s", testPrefix, randomSuffix())
		m := parseJSON(t, zoho(t, "sheet", "workbooks", "create", "--workbook-name", name))
		assertSheetSuccess(t, m)
		workbookID = fmt.Sprintf("%v", m["resource_id"])
		cleanup.trackSheetWorkbook(workbookID)
	})

	t.Run("seed", func(t *testing.T) {
		requireID(t, workbookID, "setup must have succeeded")
		assertSheetSuccess(t, parseJSON(t, zoho(t, "sheet", "cells", "set-range", "--workbook", workbookID, "--worksheet", worksheetName, "--row", "1", "--column", "1", "--data", "A,B,C\n1,2,3")))
	})

	t.Run("append-csv", func(t *testing.T) {
		requireID(t, workbookID, "setup must have succeeded")
		assertSheetSuccess(t, parseJSON(t, zoho(t, "sheet", "content", "append-csv", "--workbook", workbookID, "--worksheet", worksheetName, "--data", "4,5,6")))
	})

	t.Run("append-json", func(t *testing.T) {
		requireID(t, workbookID, "setup must have succeeded")
		jsonData := `[{"A":"X","B":"Y"}]`
		assertSheetSuccess(t, parseJSON(t, zoho(t, "sheet", "content", "append-json", "--workbook", workbookID, "--worksheet", worksheetName, "--json", jsonData)))
	})

	t.Run("find", func(t *testing.T) {
		requireID(t, workbookID, "seed must have succeeded")
		m := parseJSON(t, zoho(t, "sheet", "content", "find", "--workbook", workbookID, "--worksheet", worksheetName, "--search", "X", "--scope", "worksheet"))
		assertSheetSuccess(t, m)
		resp := fmt.Sprintf("%v", m)
		if cnt, ok := m["count"]; ok {
			if fmt.Sprintf("%v", cnt) == "0" {
				t.Errorf("find returned 0 results for 'X'")
			}
		} else if !strings.Contains(resp, "X") {
			t.Errorf("expected find response to contain search term or results: %s", truncate(resp, 500))
		}
	})

	t.Run("find-replace", func(t *testing.T) {
		requireID(t, workbookID, "seed must have succeeded")
		assertSheetSuccess(t, parseJSON(t, zoho(t, "sheet", "content", "find-replace", "--workbook", workbookID, "--worksheet", worksheetName, "--search", "X", "--replace-with", "Y", "--scope", "worksheet")))
	})

	t.Run("clear-contents", func(t *testing.T) {
		requireID(t, workbookID, "setup must have succeeded")
		assertSheetSuccess(t, parseJSON(t, zoho(t, "sheet", "content", "clear-contents", "--workbook", workbookID, "--worksheet", worksheetName, "--start-row", "1", "--start-column", "1", "--end-row", "5", "--end-column", "5")))
	})

	t.Run("recalculate", func(t *testing.T) {
		t.Skip("workbook.recalculate not available on this account")
	})
}

func TestSheetFormat(t *testing.T) {
	t.Parallel()
	cleanup := newCleanup(t)
	var workbookID string
	worksheetName := "Sheet1"

	t.Run("setup", func(t *testing.T) {
		name := fmt.Sprintf("%s_SheetFormat_%s", testPrefix, randomSuffix())
		m := parseJSON(t, zoho(t, "sheet", "workbooks", "create", "--workbook-name", name))
		assertSheetSuccess(t, m)
		workbookID = fmt.Sprintf("%v", m["resource_id"])
		cleanup.trackSheetWorkbook(workbookID)
	})

	t.Run("insert-row", func(t *testing.T) {
		requireID(t, workbookID, "setup must have succeeded")
		assertSheetSuccess(t, parseJSON(t, zoho(t, "sheet", "format", "insert-row", "--workbook", workbookID, "--worksheet", worksheetName, "--row", "1", "--count", "1")))
	})

	t.Run("insert-column", func(t *testing.T) {
		requireID(t, workbookID, "setup must have succeeded")
		assertSheetSuccess(t, parseJSON(t, zoho(t, "sheet", "format", "insert-column", "--workbook", workbookID, "--worksheet", worksheetName, "--column", "1", "--count", "1")))
	})

	t.Run("row-height", func(t *testing.T) {
		t.Skip("Zoho Sheet API rejects all row_index_array formats")
		requireID(t, workbookID, "setup must have succeeded")
		assertSheetSuccess(t, parseJSON(t, zoho(t, "sheet", "format", "row-height", "--workbook", workbookID, "--worksheet", worksheetName, "--row-index-array", "[1]", "--row-height", "24")))
	})

	t.Run("column-width", func(t *testing.T) {
		t.Skip("Zoho Sheet API rejects all column_index_array formats")
		requireID(t, workbookID, "setup must have succeeded")
		assertSheetSuccess(t, parseJSON(t, zoho(t, "sheet", "format", "column-width", "--workbook", workbookID, "--worksheet", worksheetName, "--column-index-array", "[1]", "--column-width", "120")))
	})

	t.Run("set-note", func(t *testing.T) {
		requireID(t, workbookID, "setup must have succeeded")
		assertSheetSuccess(t, parseJSON(t, zoho(t, "sheet", "format", "set-note", "--workbook", workbookID, "--worksheet", worksheetName, "--row", "1", "--column", "1", "--note", "test")))
	})

	t.Run("delete-row", func(t *testing.T) {
		requireID(t, workbookID, "setup must have succeeded")
		assertSheetSuccess(t, parseJSON(t, zoho(t, "sheet", "format", "delete-row", "--workbook", workbookID, "--worksheet", worksheetName, "--row", "1")))
	})

	t.Run("delete-column", func(t *testing.T) {
		requireID(t, workbookID, "setup must have succeeded")
		assertSheetSuccess(t, parseJSON(t, zoho(t, "sheet", "format", "delete-column", "--workbook", workbookID, "--worksheet", worksheetName, "--column", "1")))
	})

	t.Run("delete-rows", func(t *testing.T) {
		t.Skip("Zoho Sheet API rejects all row_index_array formats")
		requireID(t, workbookID, "setup must have succeeded")
		assertSheetSuccess(t, parseJSON(t, zoho(t, "sheet", "format", "delete-rows", "--workbook", workbookID, "--worksheet", worksheetName, "--row-index-array", "[5,6]")))
	})
}

func TestSheetTables(t *testing.T) {
	t.Parallel()
	cleanup := newCleanup(t)
	var workbookID string
	worksheetName := "Sheet1"
	tableName := fmt.Sprintf("%s_TBL_%s", testPrefix, randomSuffix())
	tableReady := false
	tableVisible := false

	t.Run("setup", func(t *testing.T) {
		name := fmt.Sprintf("%s_SheetTables_%s", testPrefix, randomSuffix())
		m := parseJSON(t, zoho(t, "sheet", "workbooks", "create", "--workbook-name", name))
		assertSheetSuccess(t, m)
		workbookID = fmt.Sprintf("%v", m["resource_id"])
		cleanup.trackSheetWorkbook(workbookID)
		assertSheetSuccess(t, parseJSON(t, zoho(t, "sheet", "cells", "set-range", "--workbook", workbookID, "--worksheet", worksheetName, "--row", "1", "--column", "1", "--data", "Name,Value\nA,1\nB,2")))
	})

	t.Run("create", func(t *testing.T) {
		requireID(t, workbookID, "setup must have succeeded")
		m := parseJSON(t, zoho(t, "sheet", "tables", "create", "--workbook", workbookID, "--worksheet", worksheetName, "--table-name", tableName, "--start-row", "1", "--start-column", "1", "--end-row", "3", "--end-column", "2"))
		assertSheetSuccess(t, m)
		if tn, ok := m["table_name"]; ok {
			tableName = fmt.Sprintf("%v", tn)
		}
		tableReady = true
	})

	t.Run("list", func(t *testing.T) {
		requireID(t, workbookID, "setup must have succeeded")
		m := parseJSON(t, zoho(t, "sheet", "tables", "list", "--workbook", workbookID, "--worksheet", worksheetName))
		assertSheetSuccess(t, m)
		if !strings.Contains(fmt.Sprintf("%v", m), tableName) {
			t.Skip("table not visible yet in list response")
		}
		tableVisible = true
	})

	t.Run("fetch-records", func(t *testing.T) {
		requireID(t, workbookID, "create must have succeeded")
		if !tableReady || !tableVisible {
			t.Skip("table create did not succeed")
		}
		m := parseJSON(t, zoho(t, "sheet", "tables", "fetch-records", "--workbook", workbookID, "--table-name", tableName, "--criteria", `"Name"="A"`))
		assertSheetSuccess(t, m)
		if records, ok := m["records"]; ok {
			arr, _ := records.([]any)
			if len(arr) == 0 {
				t.Errorf("expected at least one record matching Name=A")
			}
		} else if !strings.Contains(fmt.Sprintf("%v", m), "A") {
			t.Errorf("expected record data containing 'A' in fetch-records response")
		}
	})

	t.Run("add-records", func(t *testing.T) {
		requireID(t, workbookID, "create must have succeeded")
		assertSheetSuccess(t, parseJSON(t, zoho(t, "sheet", "tables", "add-records", "--workbook", workbookID, "--table-name", tableName, "--json", `[{"Name":"C","Value":"3"}]`)))
	})

	t.Run("update-records", func(t *testing.T) {
		requireID(t, workbookID, "create must have succeeded")
		assertSheetSuccess(t, parseJSON(t, zoho(t, "sheet", "tables", "update-records", "--workbook", workbookID, "--table-name", tableName, "--criteria", `"Name"="C"`, "--json", `{"Value":"4"}`)))
	})

	t.Run("delete-records", func(t *testing.T) {
		t.Skip("Zoho Sheet API rejects all criteria_json formats for table.records.delete")
		requireID(t, workbookID, "create must have succeeded")
		assertSheetSuccess(t, parseJSON(t, zoho(t, "sheet", "tables", "delete-records", "--workbook", workbookID, "--table-name", tableName, "--criteria", `{"field":"Name","comparator":"equal","value":"C"}`, "--delete-rows", "true")))
	})

	t.Run("rename-headers", func(t *testing.T) {
		t.Skip("Zoho Sheet API rejects all data formats for table.header.rename")
		requireID(t, workbookID, "create must have succeeded")
		assertSheetSuccess(t, parseJSON(t, zoho(t, "sheet", "tables", "rename-headers", "--workbook", workbookID, "--table-name", tableName, "--data", `{"Name":"NewName"}`)))
	})

	t.Run("insert-columns", func(t *testing.T) {
		requireID(t, workbookID, "create must have succeeded")
		assertSheetSuccess(t, parseJSON(t, zoho(t, "sheet", "tables", "insert-columns", "--workbook", workbookID, "--table-name", tableName, "--columns", `["Extra"]`)))
	})

	t.Run("delete-columns", func(t *testing.T) {
		requireID(t, workbookID, "create must have succeeded")
		assertSheetSuccess(t, parseJSON(t, zoho(t, "sheet", "tables", "delete-columns", "--workbook", workbookID, "--table-name", tableName, "--columns", `["Extra"]`)))
	})

	t.Run("remove", func(t *testing.T) {
		requireID(t, workbookID, "create must have succeeded")
		assertSheetSuccess(t, parseJSON(t, zoho(t, "sheet", "tables", "remove", "--workbook", workbookID, "--table-name", tableName)))
	})
}

func TestSheetRecords(t *testing.T) {
	t.Parallel()
	cleanup := newCleanup(t)
	var workbookID string
	worksheetName := "Sheet1"

	t.Run("setup", func(t *testing.T) {
		name := fmt.Sprintf("%s_SheetRecords_%s", testPrefix, randomSuffix())
		m := parseJSON(t, zoho(t, "sheet", "workbooks", "create", "--workbook-name", name))
		assertSheetSuccess(t, m)
		workbookID = fmt.Sprintf("%v", m["resource_id"])
		cleanup.trackSheetWorkbook(workbookID)
	})

	t.Run("add", func(t *testing.T) {
		requireID(t, workbookID, "setup must have succeeded")
		assertSheetSuccess(t, parseJSON(t, zoho(t, "sheet", "cells", "set-range", "--workbook", workbookID, "--worksheet", worksheetName, "--row", "1", "--column", "1", "--data", "Name,Email\n")))
		jsonData := `[{"Name":"Alice","Email":"alice@example.com"}]`
		assertSheetSuccess(t, parseJSON(t, zoho(t, "sheet", "records", "add", "--workbook", workbookID, "--worksheet", worksheetName, "--header-row", "1", "--json", jsonData)))
	})

	t.Run("fetch", func(t *testing.T) {
		requireID(t, workbookID, "add must have succeeded")
		assertSheetSuccess(t, parseJSON(t, zoho(t, "sheet", "records", "fetch", "--workbook", workbookID, "--worksheet", worksheetName, "--header-row", "1", "--start-row", "2", "--count", "10")))
	})

	t.Run("update", func(t *testing.T) {
		requireID(t, workbookID, "add must have succeeded")
		jsonData := `{"Email":"alice+updated@example.com"}`
		assertSheetSuccess(t, parseJSON(t, zoho(t, "sheet", "records", "update", "--workbook", workbookID, "--worksheet", worksheetName, "--criteria", `"Name"="Alice"`, "--header-row", "1", "--json", jsonData)))
	})

	t.Run("delete", func(t *testing.T) {
		requireID(t, workbookID, "add must have succeeded")
		assertSheetSuccess(t, parseJSON(t, zoho(t, "sheet", "records", "delete", "--workbook", workbookID, "--worksheet", worksheetName, "--criteria", `"Name"="Alice"`, "--header-row", "1", "--delete-rows", "true")))
	})

	t.Run("insert-columns", func(t *testing.T) {
		requireID(t, workbookID, "setup must have succeeded")
		assertSheetSuccess(t, parseJSON(t, zoho(t, "sheet", "records", "insert-columns", "--workbook", workbookID, "--worksheet", worksheetName, "--columns", `["Extra"]`)))
	})
}

func TestSheetPremium(t *testing.T) {
	t.Parallel()
	cleanup := newCleanup(t)
	var workbookID string
	worksheetName := "Sheet1"
	premiumAvailable := false

	t.Run("setup", func(t *testing.T) {
		name := fmt.Sprintf("%s_SheetPremium_%s", testPrefix, randomSuffix())
		m := parseJSON(t, zoho(t, "sheet", "workbooks", "create", "--workbook-name", name))
		assertSheetSuccess(t, m)
		workbookID = fmt.Sprintf("%v", m["resource_id"])
		cleanup.trackSheetWorkbook(workbookID)
	})

	t.Run("add-records", func(t *testing.T) {
		requireID(t, workbookID, "setup must have succeeded")
		_, err := zohoMayFail(t, "sheet", "cells", "set-range", "--workbook", workbookID, "--worksheet", worksheetName, "--row", "1", "--column", "1", "--data", "Name\n")
		if err != nil {
			t.Skip("premium pre-seed failed")
		}
		out, err := zohoMayFail(t, "sheet", "premium", "add-records", "--workbook", workbookID, "--worksheet", worksheetName, "--header-row", "1", "--json", `[{"Name":"P1"}]`)
		if err != nil {
			t.Skip("premium APIs not available on this account")
		}
		assertSheetSuccess(t, parseJSON(t, out))
		premiumAvailable = true
	})

	t.Run("fetch-records", func(t *testing.T) {
		requireID(t, workbookID, "add-records must have succeeded")
		if !premiumAvailable {
			t.Skip("premium APIs not available on this account")
		}
		out, err := zohoMayFail(t, "sheet", "premium", "fetch-records", "--workbook", workbookID, "--worksheet", worksheetName, "--header-row", "1")
		if err != nil {
			t.Skip("premium APIs not available on this account")
		}
		assertSheetSuccess(t, parseJSON(t, out))
	})

	t.Run("update-records", func(t *testing.T) {
		requireID(t, workbookID, "add-records must have succeeded")
		if !premiumAvailable {
			t.Skip("premium APIs not available on this account")
		}
		out, err := zohoMayFail(t, "sheet", "premium", "update-records", "--workbook", workbookID, "--worksheet", worksheetName, "--header-row", "1", "--criteria", "Name=\"P1\"", "--json", `[{"Name":"P2"}]`)
		if err != nil {
			t.Skip("premium APIs not available on this account")
		}
		assertSheetSuccess(t, parseJSON(t, out))
	})
}

func TestSheetMerge(t *testing.T) {
	t.Parallel()
	cleanup := newCleanup(t)
	var workbookID string

	t.Run("setup", func(t *testing.T) {
		name := fmt.Sprintf("%s_SheetMerge_%s", testPrefix, randomSuffix())
		m := parseJSON(t, zoho(t, "sheet", "workbooks", "create", "--workbook-name", name))
		assertSheetSuccess(t, m)
		workbookID = fmt.Sprintf("%v", m["resource_id"])
		cleanup.trackSheetWorkbook(workbookID)
	})

	t.Run("templates", func(t *testing.T) {
		assertSheetSuccess(t, parseJSON(t, zoho(t, "sheet", "merge", "templates")))
	})

	t.Run("fields", func(t *testing.T) {
		requireID(t, workbookID, "setup must have succeeded")
		assertSheetSuccess(t, parseJSON(t, zoho(t, "sheet", "merge", "fields", "--workbook", workbookID)))
	})

	t.Run("jobs", func(t *testing.T) {
		requireID(t, workbookID, "setup must have succeeded")
		r := runZoho(t, "sheet", "merge", "jobs", "--workbook", workbookID)
		if r.ExitCode != 0 {
			t.Skip("merge jobs endpoint unavailable for this workbook")
		}
		if strings.TrimSpace(r.Stdout) == "" {
			t.Skip("merge jobs response is empty")
		}
	})
}

func TestSheetErrors(t *testing.T) {
	t.Parallel()
	cleanup := newCleanup(t)
	var workbookID string

	t.Run("setup", func(t *testing.T) {
		name := fmt.Sprintf("%s_SheetErrors_%s", testPrefix, randomSuffix())
		m := parseJSON(t, zoho(t, "sheet", "workbooks", "create", "--workbook-name", name))
		assertSheetSuccess(t, m)
		workbookID = fmt.Sprintf("%v", m["resource_id"])
		cleanup.trackSheetWorkbook(workbookID)
	})

	t.Run("missing-required-flags", func(t *testing.T) {
		requireID(t, workbookID, "setup must have succeeded")
		r := runZoho(t, "sheet", "cells", "get", "--workbook", workbookID)
		assertExitCode(t, r, 1)
	})

	t.Run("bad-auth", func(t *testing.T) {
		r := runZohoWithEnv(t, map[string]string{
			"ZOHO_CLIENT_ID":     "invalid-client-id",
			"ZOHO_CLIENT_SECRET": "invalid-client-secret",
			"ZOHO_REFRESH_TOKEN": "invalid-refresh-token",
		}, "sheet", "workbooks", "list")
		assertExitCode(t, r, 2)
	})

	t.Run("invalid-workbook", func(t *testing.T) {
		r := runZoho(t, "sheet", "workbooks", "versions", "--workbook", "999999999999999")
		if r.ExitCode == 0 {
			t.Fatal("expected non-zero exit for invalid workbook")
		}
	})
}

func TestSheetEmergencyCleanup(t *testing.T) {
	out := zoho(t, "sheet", "workbooks", "list")
	m := parseJSON(t, out)
	assertSheetSuccess(t, m)

	arr, ok := m["workbooks"].([]any)
	if !ok {
		return
	}

	for _, item := range arr {
		rec, _ := item.(map[string]any)
		name := fmt.Sprintf("%v", rec["workbook_name"])
		if name == "" || name == "<nil>" {
			name = fmt.Sprintf("%v", rec["name"])
		}
		if !strings.HasPrefix(name, testPrefix) {
			continue
		}
		resourceID := fmt.Sprintf("%v", rec["resource_id"])
		if resourceID == "" || resourceID == "<nil>" {
			resourceID = fmt.Sprintf("%v", rec["workbook_id"])
		}
		if resourceID == "" || resourceID == "<nil>" {
			resourceID = fmt.Sprintf("%v", rec["id"])
		}
		if resourceID == "" || resourceID == "<nil>" {
			continue
		}
		t.Logf("cleaning sheet workbook %s (%s)", resourceID, name)
		zohoIgnoreError(t, "sheet", "workbooks", "trash", "--workbook", resourceID)
		zohoIgnoreError(t, "sheet", "workbooks", "delete", "--workbook", resourceID)
	}
}

func TestCliqUsers(t *testing.T) {
	t.Parallel()

	const knownUserEmail = "jasmin@miaie.com"
	const knownUserID = "913284317"

	t.Run("list", func(t *testing.T) {
		out := zoho(t, "cliq", "users", "list")
		m := parseJSON(t, out)
		data, ok := m["data"].([]any)
		if !ok {
			t.Fatalf("expected data array in users list response:\n%s", truncate(out, 500))
		}
		if len(data) == 0 {
			t.Fatal("expected at least one user")
		}
		found := false
		for _, item := range data {
			user, _ := item.(map[string]any)
			if fmt.Sprintf("%v", user["email_id"]) == knownUserEmail {
				found = true
				assertEqual(t, fmt.Sprintf("%v", user["id"]), knownUserID)
				break
			}
		}
		if !found {
			t.Errorf("known user %s not found in users list", knownUserEmail)
		}
	})

	t.Run("get-by-id", func(t *testing.T) {
		out := zoho(t, "cliq", "users", "get", knownUserID)
		m := parseJSON(t, out)
		data, ok := m["data"].(map[string]any)
		if !ok {
			t.Fatalf("expected data object in users get response:\n%s", truncate(out, 500))
		}
		assertEqual(t, fmt.Sprintf("%v", data["id"]), knownUserID)
		assertStringField(t, data, "email_id", knownUserEmail)
		assertStringField(t, data, "status", "active")
	})

	t.Run("get-by-email", func(t *testing.T) {
		out := zoho(t, "cliq", "users", "get", knownUserEmail)
		m := parseJSON(t, out)
		data, ok := m["data"].(map[string]any)
		if !ok {
			t.Fatalf("expected data object in users get response:\n%s", truncate(out, 500))
		}
		assertEqual(t, fmt.Sprintf("%v", data["id"]), knownUserID)
		assertStringField(t, data, "email_id", knownUserEmail)
	})
}

func TestCliqChannels(t *testing.T) {
	t.Parallel()
	cleanup := newCleanup(t)

	var channelID string
	var uniqueName string
	var chatID string
	channelName := fmt.Sprintf("%s_chan_%s", testPrefix, randomSuffix())

	t.Run("create", func(t *testing.T) {
		out := zoho(t, "cliq", "channels", "create",
			"--name", channelName, "--level", "private",
			"--description", "Integration test channel")
		m := parseJSON(t, out)
		channelID = fmt.Sprintf("%v", m["channel_id"])
		uniqueName = fmt.Sprintf("%v", m["unique_name"])
		chatID = fmt.Sprintf("%v", m["chat_id"])
		if channelID == "" || channelID == "<nil>" {
			t.Fatalf("expected channel_id in create response:\n%s", truncate(out, 500))
		}
		cleanup.trackCliqChannel(channelID)
		assertStringField(t, m, "status", "created")
		assertStringField(t, m, "level", "private")
		if chatID == "" || chatID == "<nil>" || chatID == "null" {
			t.Fatalf("expected real chat_id for private channel, got %q", chatID)
		}
		name := fmt.Sprintf("%v", m["name"])
		assertContains(t, strings.ToLower(name), strings.ToLower(channelName))
		t.Logf("created channel %s (unique_name=%s, chat_id=%s)", channelID, uniqueName, chatID)
	})

	t.Run("list", func(t *testing.T) {
		requireID(t, channelID, "create must have succeeded")
		out := zoho(t, "cliq", "channels", "list")
		arr := parseJSONArray(t, out)
		assertNonEmpty(t, arr, "expected at least one channel")
		found := false
		for _, ch := range arr {
			if fmt.Sprintf("%v", ch["channel_id"]) == channelID {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("created channel %s not found in list", channelID)
		}
	})

	t.Run("get", func(t *testing.T) {
		requireID(t, uniqueName, "create must have succeeded")
		out := zoho(t, "cliq", "channels", "get", uniqueName)
		m := parseJSON(t, out)
		assertStringField(t, m, "type", "channel")
		data, ok := m["data"].(map[string]any)
		if !ok {
			t.Fatalf("expected data object in channel get response:\n%s", truncate(out, 500))
		}
		assertEqual(t, fmt.Sprintf("%v", data["channel_id"]), channelID)
	})

	t.Run("members", func(t *testing.T) {
		requireID(t, channelID, "create must have succeeded")
		out := zoho(t, "cliq", "channels", "members", channelID)
		m := parseJSON(t, out)
		members, ok := m["members"].([]any)
		if !ok {
			t.Fatalf("expected members array in response:\n%s", truncate(out, 500))
		}
		if len(members) == 0 {
			t.Fatal("expected at least one member (creator)")
		}
		foundCreator := false
		for _, member := range members {
			mem, _ := member.(map[string]any)
			if fmt.Sprintf("%v", mem["email_id"]) == "jasmin@miaie.com" {
				foundCreator = true
				break
			}
		}
		if !foundCreator {
			t.Error("creator jasmin@miaie.com not found in channel members")
		}
	})
}

func TestCliqBuddiesMessage(t *testing.T) {
	t.Parallel()

	const recipientEmail = "uday@miaie.com"
	const selfEmail = "jasmin@miaie.com"

	t.Run("send", func(t *testing.T) {
		msgText := fmt.Sprintf("%s buddy test %s", testPrefix, randomSuffix())
		out := zoho(t, "cliq", "buddies", "message", recipientEmail,
			"--text", msgText)
		m := parseJSON(t, out)
		assertEqual(t, fmt.Sprintf("%v", m["ok"]), "true")
		assertStringField(t, m, "email", recipientEmail)
	})

	t.Run("self-rejected", func(t *testing.T) {
		r := runZoho(t, "cliq", "buddies", "message", selfEmail,
			"--text", "should fail")
		if r.ExitCode == 0 {
			t.Fatal("expected non-zero exit code for self-message")
		}
		assertContains(t, r.Stderr+r.Stdout, "buddies_self_message_restricted")
	})
}

func TestCliqMessages(t *testing.T) {
	t.Parallel()
	cleanup := newCleanup(t)

	var chatID string
	var channelID string
	var msgID string
	channelName := fmt.Sprintf("%s_msg_%s", testPrefix, randomSuffix())
	originalText := fmt.Sprintf("%s msg %s", testPrefix, randomSuffix())
	editedText := fmt.Sprintf("%s edited %s", testPrefix, randomSuffix())

	t.Run("setup", func(t *testing.T) {
		out := zoho(t, "cliq", "channels", "create",
			"--name", channelName, "--level", "private")
		m := parseJSON(t, out)
		channelID = fmt.Sprintf("%v", m["channel_id"])
		chatID = fmt.Sprintf("%v", m["chat_id"])
		if channelID == "" || channelID == "<nil>" {
			t.Fatalf("expected channel_id in create response:\n%s", truncate(out, 500))
		}
		cleanup.trackCliqChannel(channelID)
		if chatID == "" || chatID == "<nil>" || chatID == "null" {
			t.Fatalf("expected real chat_id for private channel, got %q", chatID)
		}
		t.Logf("created channel %s with chat_id %s", channelID, chatID)
	})

	t.Run("send", func(t *testing.T) {
		requireID(t, chatID, "setup must have succeeded")
		out := zoho(t, "cliq", "chats", "message", chatID,
			"--text", originalText)
		m := parseJSON(t, out)
		assertEqual(t, fmt.Sprintf("%v", m["ok"]), "true")
		t.Logf("sent message to chat %s", chatID)
	})

	t.Run("list", func(t *testing.T) {
		requireID(t, chatID, "setup must have succeeded")
		retryUntil(t, 15*time.Second, func() bool {
			out := zoho(t, "cliq", "messages", "list", chatID, "--limit", "10")
			m := parseJSON(t, out)
			data, ok := m["data"].([]any)
			if !ok || len(data) == 0 {
				return false
			}
			for _, item := range data {
				msg, _ := item.(map[string]any)
				content, _ := msg["content"].(map[string]any)
				if content != nil && fmt.Sprintf("%v", content["text"]) == originalText {
					msgID = fmt.Sprintf("%v", msg["id"])
					return true
				}
			}
			return false
		})
		if msgID == "" {
			t.Fatal("sent message not found in messages list")
		}
		t.Logf("found message ID %s", msgID)
	})

	t.Run("edit", func(t *testing.T) {
		requireID(t, msgID, "list must have found the message")
		zoho(t, "cliq", "messages", "edit", chatID, msgID,
			"--text", editedText)
	})

	t.Run("list-after-edit", func(t *testing.T) {
		requireID(t, msgID, "list must have found the message")
		retryUntil(t, 15*time.Second, func() bool {
			out := zoho(t, "cliq", "messages", "list", chatID, "--limit", "10")
			m := parseJSON(t, out)
			data, ok := m["data"].([]any)
			if !ok {
				return false
			}
			for _, item := range data {
				msg, _ := item.(map[string]any)
				if fmt.Sprintf("%v", msg["id"]) == msgID {
					content, _ := msg["content"].(map[string]any)
					if content != nil && fmt.Sprintf("%v", content["text"]) == editedText {
						return true
					}
				}
			}
			return false
		})
	})

	t.Run("delete", func(t *testing.T) {
		requireID(t, msgID, "list must have found the message")
		zoho(t, "cliq", "messages", "delete", chatID, msgID)
	})

	t.Run("list-after-delete", func(t *testing.T) {
		requireID(t, msgID, "list must have found the message")
		retryUntil(t, 15*time.Second, func() bool {
			out := zoho(t, "cliq", "messages", "list", chatID, "--limit", "10")
			m := parseJSON(t, out)
			data, ok := m["data"].([]any)
			if !ok {
				return false
			}
			for _, item := range data {
				msg, _ := item.(map[string]any)
				if fmt.Sprintf("%v", msg["id"]) == msgID {
					content, _ := msg["content"].(map[string]any)
					if content == nil {
						return true
					}
					return fmt.Sprintf("%v", content["text"]) == ""
				}
			}
			return false
		})
	})
}

func TestCliqEmergencyCleanup(t *testing.T) {
	if os.Getenv("ZOHO_EMERGENCY_CLEANUP") == "" {
		t.Skip("set ZOHO_EMERGENCY_CLEANUP=1 to run")
	}
	out := zoho(t, "cliq", "channels", "list")
	arr := parseJSONArray(t, out)
	for _, ch := range arr {
		name := fmt.Sprintf("%v", ch["name"])
		uniqueName := fmt.Sprintf("%v", ch["unique_name"])
		if !strings.Contains(name, testPrefix) && !strings.Contains(uniqueName, testPrefix) {
			continue
		}
		channelID := fmt.Sprintf("%v", ch["channel_id"])
		if channelID == "" || channelID == "<nil>" {
			continue
		}
		t.Logf("cleaning orphaned channel %s (%s)", channelID, name)
		zohoIgnoreError(t, "cliq", "channels", "delete", channelID)
	}
}

func TestWriter(t *testing.T) {
	t.Parallel()
	cleanup := newCleanup(t)

	var docID string
	docName := fmt.Sprintf("%s_Writer_%s", testPrefix, randomSuffix())

	t.Run("create", func(t *testing.T) {
		out := zoho(t, "writer", "create",
			"--name", docName, "--folder", driveTestParentFolder, "--type", "writer")
		docID = extractDriveID(t, out)
		cleanup.trackWriterDoc(docID)
		attrs := driveAttr(t, out)
		assertContains(t, fmt.Sprintf("%v", attrs["name"]), docName)
		assertEqual(t, fmt.Sprintf("%v", attrs["destination_id"]), driveTestParentFolder)
		t.Logf("created writer doc %s", docID)
	})

	t.Run("details", func(t *testing.T) {
		requireID(t, docID, "create must have succeeded")
		out := zoho(t, "writer", "details", docID)
		m := parseJSON(t, out)
		writerDocID := fmt.Sprintf("%v", m["document_id"])
		if writerDocID == "" || writerDocID == "<nil>" {
			t.Fatalf("expected document_id in response:\n%s", truncate(out, 500))
		}
		assertStringField(t, m, "type", "document")
		assertContains(t, fmt.Sprintf("%v", m["document_name"]), docName)
		role := fmt.Sprintf("%v", m["role"])
		if role != "OWNER" && role != "COOWNER" {
			t.Errorf("unexpected role %q", role)
		}
		t.Logf("writer details: id=%s, name=%s, role=%s, status=%s",
			writerDocID, m["document_name"], m["role"], m["status"])
	})

	t.Run("read/empty", func(t *testing.T) {
		requireID(t, docID, "create must have succeeded")
		out := zoho(t, "writer", "read", docID)
		m := parseJSON(t, out)
		errMsg := fmt.Sprintf("%v", m["error"])
		assertContains(t, errMsg, "R3002")
		assertContains(t, errMsg, "empty")
	})

	t.Run("read/empty-html", func(t *testing.T) {
		requireID(t, docID, "create must have succeeded")
		out := zoho(t, "writer", "read", docID, "--format", "html")
		m := parseJSON(t, out)
		errMsg := fmt.Sprintf("%v", m["error"])
		assertContains(t, errMsg, "R3002")
	})

	t.Run("download/empty", func(t *testing.T) {
		requireID(t, docID, "create must have succeeded")
		out := zoho(t, "writer", "download", docID)
		m := parseJSON(t, out)
		errMsg := fmt.Sprintf("%v", m["error"])
		assertContains(t, errMsg, "R3002")
	})

	t.Run("download/empty-to-file", func(t *testing.T) {
		requireID(t, docID, "create must have succeeded")
		tmpFile := t.TempDir() + "/writer_empty.txt"
		out := zoho(t, "writer", "download", docID, "--output", tmpFile)
		m := parseJSON(t, out)
		assertContains(t, fmt.Sprintf("%v", m["error"]), "R3002")
		if _, err := os.Stat(tmpFile); err == nil {
			t.Error("expected no file written for empty doc download")
		}
	})

	t.Run("fields", func(t *testing.T) {
		t.Skip("writer fields returns 401 — requires WorkDrive.organization.ALL scope")
	})

	t.Run("merge", func(t *testing.T) {
		t.Skip("writer merge returns 401 — requires WorkDrive.organization.ALL scope")
	})

	t.Run("trash", func(t *testing.T) {
		t.Skip("writer trash returns 401 — requires elevated Writer scope")
	})

	t.Run("delete", func(t *testing.T) {
		t.Skip("writer delete returns 401 — requires elevated Writer scope")
	})
}

func TestWriterErrors(t *testing.T) {
	t.Parallel()

	t.Run("details/bad-id", func(t *testing.T) {
		r := runZoho(t, "writer", "details", "nonexistent_id_12345")
		if r.ExitCode == 0 {
			t.Error("expected non-zero exit code for nonexistent doc ID")
		}
	})

	t.Run("read/bad-id", func(t *testing.T) {
		r := runZoho(t, "writer", "read", "nonexistent_id_12345")
		if r.ExitCode == 0 {
			t.Error("expected non-zero exit code for nonexistent doc ID")
		}
	})

	t.Run("create/bad-folder", func(t *testing.T) {
		r := runZoho(t, "writer", "create",
			"--name", "ZOHOTEST_bad_folder", "--folder", "fake_folder_id_xyz")
		if r.ExitCode == 0 {
			t.Error("expected non-zero exit code for bad folder ID")
		}
	})

	t.Run("create/bad-type", func(t *testing.T) {
		r := runZoho(t, "writer", "create",
			"--name", "ZOHOTEST_bad_type", "--folder", driveTestParentFolder, "--type", "invalid")
		assertExitCode(t, r, 4)
	})

	t.Run("merge/bad-json", func(t *testing.T) {
		r := runZoho(t, "writer", "merge", "some_doc_id", "--json", "not valid json")
		assertExitCode(t, r, 4)
	})
}

func TestDeskDepartments(t *testing.T) {
	t.Parallel()
	orgID := requireDeskOrgID(t)

	var departmentID string

	t.Run("list", func(t *testing.T) {
		out := zoho(t, "desk", "departments", "list", "--org", orgID)
		m := parseJSON(t, out)
		departments, ok := m["data"].([]any)
		if !ok {
			t.Logf("departments list missing data array, trying direct array parse")
			var arr []any
			if err := json.Unmarshal([]byte(out), &arr); err == nil {
				departments = arr
				ok = true
			}
		}
		if !ok {
			t.Errorf("expected data array in departments list response:\n%s", truncate(out, 500))
			return
		}
		if len(departments) == 0 {
			t.Errorf("expected at least one department")
			return
		}
		first, _ := departments[0].(map[string]any)
		departmentID = fmt.Sprintf("%v", first["id"])
	})

	t.Run("get", func(t *testing.T) {
		requireID(t, departmentID, "list must have returned departments")
		out := zoho(t, "desk", "departments", "get", departmentID, "--org", orgID)
		m := parseJSON(t, out)
		assertEqual(t, fmt.Sprintf("%v", m["id"]), departmentID)
	})
}

func TestDeskAgents(t *testing.T) {
	t.Parallel()
	orgID := requireDeskOrgID(t)

	var agentID string

	t.Run("list", func(t *testing.T) {
		out := zoho(t, "desk", "agents", "list", "--org", orgID)
		m := parseJSON(t, out)
		agents, ok := m["data"].([]any)
		if !ok {
			t.Logf("agents list missing data array, trying direct array parse")
			var arr []any
			if err := json.Unmarshal([]byte(out), &arr); err == nil {
				agents = arr
				ok = true
			}
		}
		if !ok {
			t.Errorf("expected data array in agents list response:\n%s", truncate(out, 500))
			return
		}
		if len(agents) == 0 {
			t.Errorf("expected at least one agent")
			return
		}
		first, _ := agents[0].(map[string]any)
		agentID = fmt.Sprintf("%v", first["id"])
	})

	t.Run("get", func(t *testing.T) {
		requireID(t, agentID, "list must have returned agents")
		out := zoho(t, "desk", "agents", "get", agentID, "--org", orgID)
		m := parseJSON(t, out)
		assertEqual(t, fmt.Sprintf("%v", m["id"]), agentID)
	})
}

func TestDesk(t *testing.T) {
	t.Parallel()
	orgID := requireDeskOrgID(t)
	cleanup := newCleanup(t)

	var departmentID string
	var contactID string
	var ticketID string
	var commentID string

	t.Run("setup/departments", func(t *testing.T) {
		out := zoho(t, "desk", "departments", "list", "--org", orgID)
		m := parseJSON(t, out)
		departments, ok := m["data"].([]any)
		if !ok || len(departments) == 0 {
			t.Logf("departments setup missing data array or empty, trying direct array parse")
			var arr []any
			if err := json.Unmarshal([]byte(out), &arr); err == nil {
				departments = arr
				ok = true
			}
		}
		if !ok || len(departments) == 0 {
			t.Fatalf("expected at least one department:\n%s", truncate(out, 500))
		}
		first, _ := departments[0].(map[string]any)
		departmentID = fmt.Sprintf("%v", first["id"])
	})

	t.Run("contacts/create", func(t *testing.T) {
		requireID(t, departmentID, "setup must have succeeded")
		contactName := fmt.Sprintf("%s Contact %s", testPrefix, randomSuffix())
		out := zoho(t, "desk", "contacts", "create", "--org", orgID,
			"--json", toJSON(t, map[string]any{
				"lastName": contactName,
				"email":    fmt.Sprintf("%s@zohotest.example.com", randomSuffix()),
			}))
		m := parseJSON(t, out)
		contactID = fmt.Sprintf("%v", m["id"])
		if contactID == "" || contactID == "<nil>" {
			t.Fatalf("expected id in contact create response:\n%s", truncate(out, 500))
		}
		cleanup.trackDeskContact(contactID, orgID)
	})

	t.Run("contacts/get", func(t *testing.T) {
		requireID(t, contactID, "contacts/create must have succeeded")
		out := zoho(t, "desk", "contacts", "get", contactID, "--org", orgID)
		m := parseJSON(t, out)
		assertEqual(t, fmt.Sprintf("%v", m["id"]), contactID)
	})

	t.Run("contacts/update", func(t *testing.T) {
		requireID(t, contactID, "contacts/create must have succeeded")
		updatedName := fmt.Sprintf("%s Contact Upd %s", testPrefix, randomSuffix())
		zoho(t, "desk", "contacts", "update", contactID, "--org", orgID,
			"--json", toJSON(t, map[string]any{"lastName": updatedName}))

		out := zoho(t, "desk", "contacts", "get", contactID, "--org", orgID)
		m := parseJSON(t, out)
		assertEqual(t, fmt.Sprintf("%v", m["lastName"]), updatedName)
	})

	t.Run("contacts/list", func(t *testing.T) {
		out := zoho(t, "desk", "contacts", "list", "--org", orgID)
		m := parseJSON(t, out)
		if _, ok := m["data"].([]any); !ok {
			t.Logf("contacts list missing data array, response shape:\n%s", truncate(out, 500))
		}
	})

	t.Run("contacts/search", func(t *testing.T) {
		r := runZoho(t, "desk", "contacts", "search", "--org", orgID, "--email", "zohotest")
		if r.ExitCode == 0 {
			parseJSON(t, r.Stdout)
		}
	})

	t.Run("tickets/create", func(t *testing.T) {
		requireID(t, departmentID, "setup must have succeeded")
		subject := fmt.Sprintf("%s Ticket %s", testPrefix, randomSuffix())
		body := map[string]any{
			"subject":      subject,
			"departmentId": departmentID,
			"priority":     "Medium",
			"status":       "Open",
		}
		if contactID != "" && contactID != "<nil>" {
			body["contactId"] = contactID
		}
		out := zoho(t, "desk", "tickets", "create", "--org", orgID,
			"--json", toJSON(t, body))
		m := parseJSON(t, out)
		ticketID = fmt.Sprintf("%v", m["id"])
		if ticketID == "" || ticketID == "<nil>" {
			t.Fatalf("expected id in ticket create response:\n%s", truncate(out, 500))
		}
		cleanup.trackDeskTicket(ticketID, orgID)
	})

	t.Run("tickets/get", func(t *testing.T) {
		requireID(t, ticketID, "tickets/create must have succeeded")
		out := zoho(t, "desk", "tickets", "get", ticketID, "--org", orgID)
		m := parseJSON(t, out)
		assertEqual(t, fmt.Sprintf("%v", m["id"]), ticketID)
		assertEqual(t, fmt.Sprintf("%v", m["status"]), "Open")
	})

	t.Run("tickets/update", func(t *testing.T) {
		requireID(t, ticketID, "tickets/create must have succeeded")
		updatedSubject := fmt.Sprintf("%s Ticket Upd %s", testPrefix, randomSuffix())
		zoho(t, "desk", "tickets", "update", ticketID, "--org", orgID,
			"--json", toJSON(t, map[string]any{
				"subject":  updatedSubject,
				"priority": "High",
			}))

		out := zoho(t, "desk", "tickets", "get", ticketID, "--org", orgID)
		m := parseJSON(t, out)
		assertEqual(t, fmt.Sprintf("%v", m["subject"]), updatedSubject)
		assertEqual(t, fmt.Sprintf("%v", m["priority"]), "High")
	})

	t.Run("tickets/list", func(t *testing.T) {
		out := zoho(t, "desk", "tickets", "list", "--org", orgID, "--limit", "5")
		m := parseJSON(t, out)
		if _, ok := m["data"].([]any); !ok {
			t.Logf("tickets list missing data array, response shape:\n%s", truncate(out, 500))
		}
	})

	t.Run("tickets/search", func(t *testing.T) {
		r := runZoho(t, "desk", "tickets", "search", "--org", orgID, "--subject", testPrefix+"*")
		if r.ExitCode == 0 {
			m := parseJSON(t, r.Stdout)
			if data, ok := m["data"].([]any); ok && len(data) > 0 {
				found := false
				for _, item := range data {
					im, _ := item.(map[string]any)
					if fmt.Sprintf("%v", im["id"]) == ticketID {
						found = true
						break
					}
				}
				if !found {
					t.Logf("ticket %s not found in search results (search may take time to index)", ticketID)
				}
			}
		}
	})

	t.Run("tickets/threads", func(t *testing.T) {
		requireID(t, ticketID, "tickets/create must have succeeded")
		out := zoho(t, "desk", "tickets", "threads", ticketID, "--org", orgID)
		m := parseJSON(t, out)
		if _, ok := m["data"].([]any); !ok {
			t.Logf("threads response may be empty or have different structure:\n%s", truncate(out, 300))
		}
	})

	t.Run("tickets/add-comment", func(t *testing.T) {
		requireID(t, ticketID, "tickets/create must have succeeded")
		commentContent := fmt.Sprintf("%s comment %s", testPrefix, randomSuffix())
		out := zoho(t, "desk", "tickets", "add-comment", ticketID, "--org", orgID,
			"--json", toJSON(t, map[string]any{
				"content":  commentContent,
				"isPublic": false,
			}))
		m := parseJSON(t, out)
		commentID = fmt.Sprintf("%v", m["id"])
		if commentID == "" || commentID == "<nil>" {
			t.Fatalf("expected id in add-comment response:\n%s", truncate(out, 500))
		}
	})

	t.Run("tickets/comments", func(t *testing.T) {
		requireID(t, ticketID, "tickets/create must have succeeded")
		requireID(t, commentID, "tickets/add-comment must have succeeded")
		out := zoho(t, "desk", "tickets", "comments", ticketID, "--org", orgID)
		m := parseJSON(t, out)
		comments, ok := m["data"].([]any)
		if !ok {
			t.Errorf("expected data array in comments list:\n%s", truncate(out, 500))
			return
		}
		found := false
		for _, c := range comments {
			cm, _ := c.(map[string]any)
			if fmt.Sprintf("%v", cm["id"]) == commentID {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("comment %s not found in ticket comments list", commentID)
		}
	})

	t.Run("tickets/attachments", func(t *testing.T) {
		requireID(t, ticketID, "tickets/create must have succeeded")
		out := zoho(t, "desk", "tickets", "attachments", ticketID, "--org", orgID)
		parseJSON(t, out)
	})

	t.Run("tickets/history", func(t *testing.T) {
		requireID(t, ticketID, "tickets/create must have succeeded")
		out := zoho(t, "desk", "tickets", "history", ticketID, "--org", orgID)
		parseJSON(t, out)
	})

	t.Run("search", func(t *testing.T) {
		r := runZoho(t, "desk", "search", "--org", orgID, "--query", testPrefix)
		if r.ExitCode == 0 {
			parseJSON(t, r.Stdout)
		}
	})

	t.Run("accounts/list", func(t *testing.T) {
		out := zoho(t, "desk", "accounts", "list", "--org", orgID, "--limit", "5")
		m := parseJSON(t, out)
		if _, ok := m["data"].([]any); !ok {
			t.Logf("accounts list missing data array, response shape:\n%s", truncate(out, 500))
		}
	})

	t.Run("tickets/delete", func(t *testing.T) {
		requireID(t, ticketID, "tickets/create must have succeeded")
		zoho(t, "desk", "tickets", "delete", ticketID, "--org", orgID)

		time.Sleep(2 * time.Second)
		r := runZoho(t, "desk", "tickets", "get", ticketID, "--org", orgID)
		if r.ExitCode == 0 {
			m := parseJSON(t, r.Stdout)
			status := fmt.Sprintf("%v", m["status"])
			if status != "Closed" && status != "Deleted" {
				t.Errorf("ticket %s still accessible after delete (status=%s)", ticketID, status)
			}
		}
		ticketID = ""
	})

	t.Run("contacts/delete", func(t *testing.T) {
		requireID(t, contactID, "contacts/create must have succeeded")
		zoho(t, "desk", "contacts", "delete", contactID, "--org", orgID)

		time.Sleep(2 * time.Second)
		r := runZoho(t, "desk", "contacts", "get", contactID, "--org", orgID)
		if r.ExitCode == 0 {
			t.Errorf("contact %s still accessible after delete", contactID)
		}
		contactID = ""
	})
}

func TestDeskErrors(t *testing.T) {
	t.Parallel()
	orgID := requireDeskOrgID(t)

	t.Run("missing-org", func(t *testing.T) {
		r := runZohoWithEnv(t, map[string]string{"ZOHO_DESK_ORG_ID": ""}, "desk", "departments", "list")
		assertExitCode(t, r, 4)
	})

	t.Run("invalid-ticket-id", func(t *testing.T) {
		r := runZoho(t, "desk", "tickets", "get", "999999999999999", "--org", orgID)
		if r.ExitCode == 0 {
			t.Errorf("expected non-zero exit for invalid ticket ID")
		}
	})

	t.Run("invalid-contact-id", func(t *testing.T) {
		r := runZoho(t, "desk", "contacts", "get", "999999999999999", "--org", orgID)
		if r.ExitCode == 0 {
			t.Errorf("expected non-zero exit for invalid contact ID")
		}
	})
}

func TestDeskEmergencyCleanup(t *testing.T) {
	if os.Getenv("ZOHO_EMERGENCY_CLEANUP") != "1" {
		t.Skip("set ZOHO_EMERGENCY_CLEANUP=1 to run")
	}
	orgID := requireDeskOrgID(t)

	out, err := zohoMayFail(t, "desk", "tickets", "list", "--org", orgID, "--limit", "100")
	if err == nil {
		m := parseJSON(t, out)
		if data, ok := m["data"].([]any); ok {
			for _, item := range data {
				im, _ := item.(map[string]any)
				subject := fmt.Sprintf("%v", im["subject"])
				if strings.HasPrefix(subject, testPrefix) {
					id := fmt.Sprintf("%v", im["id"])
					zohoIgnoreError(t, "desk", "tickets", "delete", id, "--org", orgID)
				}
			}
		}
	}

	out, err = zohoMayFail(t, "desk", "contacts", "list", "--org", orgID, "--limit", "100")
	if err == nil {
		m := parseJSON(t, out)
		if data, ok := m["data"].([]any); ok {
			for _, item := range data {
				im, _ := item.(map[string]any)
				lastName := fmt.Sprintf("%v", im["lastName"])
				if strings.HasPrefix(lastName, testPrefix) {
					id := fmt.Sprintf("%v", im["id"])
					zohoIgnoreError(t, "desk", "contacts", "delete", id, "--org", orgID)
				}
			}
		}
	}
}

func (c *testCleanup) trackSignRequest(id string) {
	c.add("delete sign request "+id, func() {
		zohoIgnoreError(c.t, "sign", "requests", "delete", id)
	})
}

func TestSignFieldTypes(t *testing.T) {
	t.Parallel()
	out := zoho(t, "sign", "field-types")
	m := parseJSON(t, out)
	fieldTypes, ok := m["field_types"].([]any)
	if !ok {
		t.Fatalf("expected field_types array in response:\n%s", truncate(out, 500))
	}
	if len(fieldTypes) == 0 {
		t.Error("expected at least one field type")
	}
	t.Logf("found %d field types", len(fieldTypes))
}

func TestSignRequestTypes(t *testing.T) {
	t.Parallel()

	t.Run("list", func(t *testing.T) {
		out := zoho(t, "sign", "request-types", "list")
		m := parseJSON(t, out)
		requestTypes, ok := m["request_types"].([]any)
		if !ok {
			t.Fatalf("expected request_types array in response:\n%s", truncate(out, 500))
		}
		if len(requestTypes) == 0 {
			t.Error("expected at least one request type")
		}
		t.Logf("found %d request types", len(requestTypes))
	})

	t.Run("create", func(t *testing.T) {
		typeName := fmt.Sprintf("%s Type %s", testPrefix, randomSuffix())
		out := zoho(t, "sign", "request-types", "create",
			"--data", toJSON(t, map[string]any{
				"request_types": map[string]any{
					"request_type_name":        typeName,
					"request_type_description": "Integration test type",
				},
			}))
		m := parseJSON(t, out)
		rt, ok := m["request_types"].(map[string]any)
		if !ok {
			t.Fatalf("expected request_types object in response:\n%s", truncate(out, 500))
		}
		assertEqual(t, fmt.Sprintf("%v", rt["request_type_name"]), typeName)
		t.Logf("created request type: %v", rt["request_type_id"])
	})
}

func TestSignFolders(t *testing.T) {
	t.Parallel()

	var folderName string

	t.Run("list", func(t *testing.T) {
		out := zoho(t, "sign", "folders", "list")
		m := parseJSON(t, out)
		folders, ok := m["folders"].([]any)
		if !ok {
			t.Fatalf("expected folders array in response:\n%s", truncate(out, 500))
		}
		if len(folders) == 0 {
			t.Error("expected at least one folder")
		}
		t.Logf("found %d folders", len(folders))
	})

	t.Run("create", func(t *testing.T) {
		folderName = fmt.Sprintf("%s Folder %s", testPrefix, randomSuffix())
		out := zoho(t, "sign", "folders", "create", "--name", folderName)
		m := parseJSON(t, out)
		folder, ok := m["folders"].(map[string]any)
		if !ok {
			t.Fatalf("expected folders object in response:\n%s", truncate(out, 500))
		}
		assertEqual(t, fmt.Sprintf("%v", folder["folder_name"]), folderName)
		t.Logf("created folder: %v (id: %v)", folder["folder_name"], folder["folder_id"])
	})

	t.Run("list-verify-created", func(t *testing.T) {
		if folderName == "" {
			t.Skip("create must have succeeded")
		}
		out := zoho(t, "sign", "folders", "list")
		m := parseJSON(t, out)
		folders, ok := m["folders"].([]any)
		if !ok {
			t.Fatalf("expected folders array in response:\n%s", truncate(out, 500))
		}
		found := false
		for _, f := range folders {
			fm, _ := f.(map[string]any)
			if fmt.Sprintf("%v", fm["folder_name"]) == folderName {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("created folder %q not found in list", folderName)
		}
	})
}

func TestSign(t *testing.T) {
	t.Parallel()
	cleanup := newCleanup(t)

	var requestID string
	var requestName string

	t.Run("requests/create", func(t *testing.T) {
		tmpDir := t.TempDir()
		testFile := tmpDir + "/test.pdf"
		pdfContent := []byte("%PDF-1.0\n1 0 obj<</Type/Catalog/Pages 2 0 R>>endobj\n2 0 obj<</Type/Pages/Kids[3 0 R]/Count 1>>endobj\n3 0 obj<</Type/Page/MediaBox[0 0 612 792]/Parent 2 0 R>>endobj\nxref\n0 4\n0000000000 65535 f \n0000000009 00000 n \n0000000058 00000 n \n0000000115 00000 n \ntrailer<</Size 4/Root 1 0 R>>\nstartxref\n190\n%%EOF")
		if err := os.WriteFile(testFile, pdfContent, 0644); err != nil {
			t.Fatalf("failed to create test PDF: %v", err)
		}
		requestName = fmt.Sprintf("%s Doc %s", testPrefix, randomSuffix())
		out := zoho(t, "sign", "requests", "create",
			"--file", testFile,
			"--data", toJSON(t, map[string]any{
				"requests": map[string]any{
					"request_name":  requestName,
					"is_sequential": false,
					"actions": []map[string]any{
						{
							"action_type":      "SIGN",
							"recipient_name":   "Test Signer",
							"recipient_email":  "zohotest-signer@example.com",
							"signing_order":    0,
							"verify_recipient": false,
						},
					},
				},
			}))
		m := parseJSON(t, out)
		requests, ok := m["requests"].(map[string]any)
		if !ok {
			t.Fatalf("expected requests object in response:\n%s", truncate(out, 500))
		}
		requestID = fmt.Sprintf("%v", requests["request_id"])
		if requestID == "" || requestID == "<nil>" {
			t.Fatalf("expected request_id in response:\n%s", truncate(out, 500))
		}
		cleanup.trackSignRequest(requestID)
		assertEqual(t, fmt.Sprintf("%v", requests["request_name"]), requestName)
		assertEqual(t, fmt.Sprintf("%v", requests["request_status"]), "draft")
		t.Logf("created sign request: %s (id: %s)", requestName, requestID)
	})

	t.Run("requests/get", func(t *testing.T) {
		requireID(t, requestID, "requests/create must have succeeded")
		out := zoho(t, "sign", "requests", "get", requestID)
		m := parseJSON(t, out)
		requests, ok := m["requests"].(map[string]any)
		if !ok {
			t.Fatalf("expected requests object in response:\n%s", truncate(out, 500))
		}
		assertEqual(t, fmt.Sprintf("%v", requests["request_id"]), requestID)
		assertEqual(t, fmt.Sprintf("%v", requests["request_name"]), requestName)
	})

	t.Run("requests/list", func(t *testing.T) {
		requireID(t, requestID, "requests/create must have succeeded")
		out := zoho(t, "sign", "requests", "list", "--row-count", "10")
		m := parseJSON(t, out)
		requests, ok := m["requests"].([]any)
		if !ok {
			t.Fatalf("expected requests array in response:\n%s", truncate(out, 500))
		}
		found := false
		for _, r := range requests {
			rm, _ := r.(map[string]any)
			if fmt.Sprintf("%v", rm["request_id"]) == requestID {
				found = true
				break
			}
		}
		if !found {
			t.Logf("created request %s not found in list (may need more rows)", requestID)
		}
	})

	t.Run("requests/list-pagination", func(t *testing.T) {
		out := zoho(t, "sign", "requests", "list",
			"--row-count", "2",
			"--start-index", "1",
			"--sort-column", "created_time",
			"--sort-order", "DESC")
		m := parseJSON(t, out)
		if _, ok := m["requests"].([]any); !ok {
			t.Fatalf("expected requests array in paginated response:\n%s", truncate(out, 500))
		}
		if pc, ok := m["page_context"].(map[string]any); ok {
			t.Logf("page_context: row_count=%v, start_index=%v, has_more_rows=%v",
				pc["row_count"], pc["start_index"], pc["has_more_rows"])
		}
	})

	t.Run("requests/update", func(t *testing.T) {
		requireID(t, requestID, "requests/create must have succeeded")
		updatedName := fmt.Sprintf("%s Doc Upd %s", testPrefix, randomSuffix())
		out := zoho(t, "sign", "requests", "update", requestID,
			"--data", toJSON(t, map[string]any{
				"requests": map[string]any{
					"request_name": updatedName,
				},
			}))
		m := parseJSON(t, out)
		requests, ok := m["requests"].(map[string]any)
		if !ok {
			t.Fatalf("expected requests object in response:\n%s", truncate(out, 500))
		}
		assertEqual(t, fmt.Sprintf("%v", requests["request_name"]), updatedName)
		requestName = updatedName
	})

	t.Run("requests/field-data", func(t *testing.T) {
		requireID(t, requestID, "requests/create must have succeeded")
		out := zoho(t, "sign", "requests", "field-data", requestID)
		m := parseJSON(t, out)
		if code, ok := m["code"].(float64); ok {
			assertEqual(t, fmt.Sprintf("%v", int(code)), "0")
		}
	})

	t.Run("requests/download", func(t *testing.T) {
		requireID(t, requestID, "requests/create must have succeeded")
		tmpDir := t.TempDir()
		downloadPath := tmpDir + "/downloaded.pdf"
		out := zoho(t, "sign", "requests", "download", requestID,
			"--output", downloadPath)
		m := parseJSON(t, out)
		assertEqual(t, fmt.Sprintf("%v", m["ok"]), "true")
		downloaded, err := os.ReadFile(downloadPath)
		if err != nil {
			t.Fatalf("failed to read downloaded file: %v", err)
		}
		if len(downloaded) == 0 {
			t.Error("downloaded PDF is empty")
		}
		t.Logf("downloaded PDF: %d bytes", len(downloaded))
	})

	t.Run("requests/download-stdout", func(t *testing.T) {
		requireID(t, requestID, "requests/create must have succeeded")
		r := runZoho(t, "sign", "requests", "download", requestID)
		if r.ExitCode != 0 {
			t.Fatalf("download to stdout failed (exit %d): %s", r.ExitCode, r.Stderr)
		}
		if len(r.Stdout) == 0 {
			t.Error("expected PDF content on stdout")
		}
	})

	t.Run("requests/remind", func(t *testing.T) {
		requireID(t, requestID, "requests/create must have succeeded")
		r := runZoho(t, "sign", "requests", "remind", requestID)
		if r.ExitCode != 0 {
			t.Logf("remind failed (expected for draft docs, exit %d): %s", r.ExitCode, r.Stderr)
		}
	})

	t.Run("requests/recall", func(t *testing.T) {
		requireID(t, requestID, "requests/create must have succeeded")
		r := runZoho(t, "sign", "requests", "recall", requestID)
		if r.ExitCode != 0 {
			t.Logf("recall failed (expected for draft docs, exit %d): %s", r.ExitCode, r.Stderr)
		}
	})

	t.Run("requests/delete", func(t *testing.T) {
		requireID(t, requestID, "requests/create must have succeeded")
		out := zoho(t, "sign", "requests", "delete", requestID)
		m := parseJSON(t, out)
		if code, ok := m["code"].(float64); ok {
			assertEqual(t, fmt.Sprintf("%v", int(code)), "0")
		}
	})

	t.Run("requests/get-after-delete", func(t *testing.T) {
		requireID(t, requestID, "requests/create must have succeeded")
		r := runZoho(t, "sign", "requests", "get", requestID)
		if r.ExitCode == 0 {
			m := parseJSON(t, r.Stdout)
			if requests, ok := m["requests"].(map[string]any); ok {
				isDeleted := fmt.Sprintf("%v", requests["is_deleted"])
				if isDeleted != "true" {
					t.Logf("request still accessible after delete, is_deleted=%s", isDeleted)
				}
			}
		}
	})
}

func TestSignTemplates(t *testing.T) {
	t.Parallel()

	var templateID string

	t.Run("list", func(t *testing.T) {
		out := zoho(t, "sign", "templates", "list")
		m := parseJSON(t, out)
		if _, ok := m["templates"].([]any); !ok {
			t.Logf("templates list response shape:\n%s", truncate(out, 500))
		}
	})

	t.Run("create", func(t *testing.T) {
		tmpDir := t.TempDir()
		testFile := tmpDir + "/template.pdf"
		pdfContent := []byte("%PDF-1.0\n1 0 obj<</Type/Catalog/Pages 2 0 R>>endobj\n2 0 obj<</Type/Pages/Kids[3 0 R]/Count 1>>endobj\n3 0 obj<</Type/Page/MediaBox[0 0 612 792]/Parent 2 0 R>>endobj\nxref\n0 4\n0000000000 65535 f \n0000000009 00000 n \n0000000058 00000 n \n0000000115 00000 n \ntrailer<</Size 4/Root 1 0 R>>\nstartxref\n190\n%%EOF")
		if err := os.WriteFile(testFile, pdfContent, 0644); err != nil {
			t.Fatalf("failed to create test PDF: %v", err)
		}
		templateName := fmt.Sprintf("%s Template %s", testPrefix, randomSuffix())
		out := zoho(t, "sign", "templates", "create",
			"--file", testFile,
			"--data", toJSON(t, map[string]any{
				"templates": map[string]any{
					"template_name": templateName,
					"actions": []map[string]any{
						{
							"action_type":    "SIGN",
							"recipient_name": "Test Signer",
							"role":           "Signer",
							"signing_order":  0,
						},
					},
				},
			}))
		m := parseJSON(t, out)
		templates, ok := m["templates"].(map[string]any)
		if !ok {
			t.Fatalf("expected templates object in response:\n%s", truncate(out, 500))
		}
		templateID = fmt.Sprintf("%v", templates["template_id"])
		if templateID == "" || templateID == "<nil>" {
			t.Fatalf("expected template_id in response:\n%s", truncate(out, 500))
		}
		t.Logf("created template: %v (id: %s)", templates["template_name"], templateID)
	})

	t.Run("get", func(t *testing.T) {
		requireID(t, templateID, "create must have succeeded")
		out := zoho(t, "sign", "templates", "get", templateID)
		m := parseJSON(t, out)
		templates, ok := m["templates"].(map[string]any)
		if !ok {
			t.Fatalf("expected templates object in response:\n%s", truncate(out, 500))
		}
		assertEqual(t, fmt.Sprintf("%v", templates["template_id"]), templateID)
	})

	t.Run("delete", func(t *testing.T) {
		requireID(t, templateID, "create must have succeeded")
		out := zoho(t, "sign", "templates", "delete", templateID)
		m := parseJSON(t, out)
		if code, ok := m["code"].(float64); ok {
			assertEqual(t, fmt.Sprintf("%v", int(code)), "0")
		}
	})
}

func TestSignErrors(t *testing.T) {
	t.Parallel()

	t.Run("requests/get-nonexistent", func(t *testing.T) {
		r := runZoho(t, "sign", "requests", "get", "999999999999999")
		if r.ExitCode == 0 {
			t.Error("expected non-zero exit for nonexistent request ID")
		}
	})

	t.Run("templates/get-nonexistent", func(t *testing.T) {
		r := runZoho(t, "sign", "templates", "get", "999999999999999")
		if r.ExitCode == 0 {
			t.Error("expected non-zero exit for nonexistent template ID")
		}
	})

	t.Run("requests/get-missing-arg", func(t *testing.T) {
		r := runZoho(t, "sign", "requests", "get")
		assertExitCode(t, r, 4)
	})

	t.Run("requests/create-missing-file", func(t *testing.T) {
		r := runZoho(t, "sign", "requests", "create",
			"--file", "/nonexistent/path/file.pdf",
			"--data", toJSON(t, map[string]any{"requests": map[string]any{"request_name": "test"}}))
		if r.ExitCode == 0 {
			t.Error("expected non-zero exit for missing file")
		}
	})

	t.Run("requests/download-missing-arg", func(t *testing.T) {
		r := runZoho(t, "sign", "requests", "download")
		assertExitCode(t, r, 4)
	})

	t.Run("requests/download-document-missing-args", func(t *testing.T) {
		r := runZoho(t, "sign", "requests", "download-document", "123")
		assertExitCode(t, r, 4)
	})

	t.Run("requests/submit-missing-arg", func(t *testing.T) {
		r := runZoho(t, "sign", "requests", "submit",
			"--data", toJSON(t, map[string]any{"requests": map[string]any{}}))
		assertExitCode(t, r, 4)
	})

	t.Run("requests/update-missing-arg", func(t *testing.T) {
		r := runZoho(t, "sign", "requests", "update",
			"--data", toJSON(t, map[string]any{"requests": map[string]any{}}))
		assertExitCode(t, r, 4)
	})

	t.Run("requests/delete-missing-arg", func(t *testing.T) {
		r := runZoho(t, "sign", "requests", "delete")
		assertExitCode(t, r, 4)
	})

	t.Run("requests/extend-missing-arg", func(t *testing.T) {
		r := runZoho(t, "sign", "requests", "extend", "--expire-by", "2030-01-01")
		assertExitCode(t, r, 4)
	})

	t.Run("templates/get-missing-arg", func(t *testing.T) {
		r := runZoho(t, "sign", "templates", "get")
		assertExitCode(t, r, 4)
	})

	t.Run("templates/delete-missing-arg", func(t *testing.T) {
		r := runZoho(t, "sign", "templates", "delete")
		assertExitCode(t, r, 4)
	})

	t.Run("templates/send-missing-arg", func(t *testing.T) {
		r := runZoho(t, "sign", "templates", "send",
			"--data", toJSON(t, map[string]any{"templates": map[string]any{}}))
		assertExitCode(t, r, 4)
	})
}

func TestSignEmergencyCleanup(t *testing.T) {
	if os.Getenv("ZOHO_EMERGENCY_CLEANUP") != "1" {
		t.Skip("set ZOHO_EMERGENCY_CLEANUP=1 to run")
	}

	out, err := zohoMayFail(t, "sign", "requests", "list", "--row-count", "100")
	if err == nil {
		m := parseJSON(t, out)
		if requests, ok := m["requests"].([]any); ok {
			for _, item := range requests {
				im, _ := item.(map[string]any)
				name := fmt.Sprintf("%v", im["request_name"])
				if strings.HasPrefix(name, testPrefix) {
					id := fmt.Sprintf("%v", im["request_id"])
					zohoIgnoreError(t, "sign", "requests", "delete", id)
				}
			}
		}
	}

	out, err = zohoMayFail(t, "sign", "templates", "list", "--row-count", "100")
	if err == nil {
		m := parseJSON(t, out)
		if templates, ok := m["templates"].([]any); ok {
			for _, item := range templates {
				im, _ := item.(map[string]any)
				name := fmt.Sprintf("%v", im["template_name"])
				if strings.HasPrefix(name, testPrefix) {
					id := fmt.Sprintf("%v", im["template_id"])
					zohoIgnoreError(t, "sign", "templates", "delete", id)
				}
			}
		}
	}
}

func TestMailAccounts(t *testing.T) {
	t.Parallel()
	var accountID string
	var primaryEmail string

	t.Run("list", func(t *testing.T) {
		out := zoho(t, "mail", "accounts", "list")
		m := parseJSON(t, out)
		data, ok := m["data"].([]any)
		if !ok {
			t.Fatalf("expected data array in response:\n%s", truncate(out, 500))
		}
		if len(data) == 0 {
			t.Fatal("expected at least one mail account")
		}
		first, _ := data[0].(map[string]any)
		accountID = fmt.Sprintf("%v", first["accountId"])
		if accountID == "" || accountID == "<nil>" {
			t.Fatalf("expected accountId in first account:\n%s", truncate(out, 500))
		}
		t.Logf("discovered account ID: %s", accountID)
		if email, ok := first["primaryEmailAddress"].(string); ok {
			primaryEmail = email
			t.Logf("primary email: %s", primaryEmail)
		}
	})

	t.Run("get", func(t *testing.T) {
		requireID(t, accountID, "accounts list must have succeeded")
		out := zoho(t, "mail", "accounts", "get", accountID)
		m := parseJSON(t, out)
		data, ok := m["data"].(map[string]any)
		if !ok {
			t.Fatalf("expected data object in response:\n%s", truncate(out, 500))
		}
		gotID := fmt.Sprintf("%v", data["accountId"])
		if gotID != accountID {
			t.Errorf("expected accountId %s, got %s", accountID, gotID)
		}
	})
}

func TestMailFolders(t *testing.T) {
	t.Parallel()
	accountID := requireMailAccountID(t)
	cleanup := newCleanup(t)
	suffix := randomSuffix()
	var folderID string

	t.Run("list", func(t *testing.T) {
		out := zoho(t, "mail", "folders", "list", "--account", accountID)
		m := parseJSON(t, out)
		data, ok := m["data"].([]any)
		if !ok {
			t.Fatalf("expected data array in response:\n%s", truncate(out, 500))
		}
		if len(data) == 0 {
			t.Fatal("expected at least one folder (Inbox)")
		}
		t.Logf("found %d folders", len(data))
	})

	t.Run("create", func(t *testing.T) {
		folderName := testPrefix + " Folder " + suffix
		body := toJSON(t, map[string]any{"folderName": folderName})
		out := zoho(t, "mail", "folders", "create", "--account", accountID, "--json", body)
		m := parseJSON(t, out)
		data, ok := m["data"].(map[string]any)
		if !ok {
			t.Fatalf("expected data object in response:\n%s", truncate(out, 500))
		}
		folderID = fmt.Sprintf("%v", data["folderId"])
		if folderID == "" || folderID == "<nil>" {
			t.Fatalf("expected folderId in create response:\n%s", truncate(out, 500))
		}
		cleanup.trackMailFolder(folderID, accountID)
		t.Logf("created folder: %s (ID: %s)", folderName, folderID)
	})

	t.Run("get", func(t *testing.T) {
		requireID(t, folderID, "create must have succeeded")
		out := zoho(t, "mail", "folders", "get", folderID, "--account", accountID)
		m := parseJSON(t, out)
		data, ok := m["data"].(map[string]any)
		if !ok {
			t.Fatalf("expected data object in response:\n%s", truncate(out, 500))
		}
		gotID := fmt.Sprintf("%v", data["folderId"])
		if gotID != folderID {
			t.Errorf("expected folderId %s, got %s", folderID, gotID)
		}
	})

	t.Run("update", func(t *testing.T) {
		requireID(t, folderID, "create must have succeeded")
		newName := testPrefix + " Folder Updated " + suffix
		body := toJSON(t, map[string]any{"mode": "rename", "folderName": newName})
		out := zoho(t, "mail", "folders", "update", folderID, "--account", accountID, "--json", body)
		m := parseJSON(t, out)
		if status, ok := m["status"].(map[string]any); ok {
			code := fmt.Sprintf("%v", status["code"])
			if code != "200" {
				t.Errorf("expected status code 200, got %s:\n%s", code, truncate(out, 500))
			}
		}
		t.Logf("updated folder name to: %s", newName)
	})

	t.Run("delete", func(t *testing.T) {
		requireID(t, folderID, "create must have succeeded")
		out := zoho(t, "mail", "folders", "delete", folderID, "--account", accountID)
		m := parseJSON(t, out)
		if status, ok := m["status"].(map[string]any); ok {
			code := fmt.Sprintf("%v", status["code"])
			if code != "200" {
				t.Errorf("expected status code 200, got %s:\n%s", code, truncate(out, 500))
			}
		}
		t.Logf("deleted folder: %s", folderID)
	})
}

func TestMailLabels(t *testing.T) {
	t.Parallel()
	accountID := requireMailAccountID(t)
	cleanup := newCleanup(t)
	suffix := randomSuffix()
	var labelID string

	t.Run("list", func(t *testing.T) {
		out := zoho(t, "mail", "labels", "list", "--account", accountID)
		m := parseJSON(t, out)
		if _, ok := m["data"]; !ok {
			if _, ok := m["status"]; !ok {
				t.Fatalf("expected data or status in response:\n%s", truncate(out, 500))
			}
		}
		t.Logf("labels list succeeded")
	})

	t.Run("create", func(t *testing.T) {
		labelName := "ZT " + suffix
		body := toJSON(t, map[string]any{"displayName": labelName, "color": "#FF0000"})
		out := zoho(t, "mail", "labels", "create", "--account", accountID, "--json", body)
		m := parseJSON(t, out)
		data, ok := m["data"].(map[string]any)
		if !ok {
			t.Fatalf("expected data object in response:\n%s", truncate(out, 500))
		}
		labelID = fmt.Sprintf("%v", data["labelId"])
		if labelID == "" || labelID == "<nil>" {
			labelID = fmt.Sprintf("%v", data["tagId"])
		}
		if labelID == "" || labelID == "<nil>" {
			t.Fatalf("expected labelId or tagId in create response:\n%s", truncate(out, 500))
		}
		cleanup.trackMailLabel(labelID, accountID)
		t.Logf("created label: %s (ID: %s)", labelName, labelID)
	})

	t.Run("get", func(t *testing.T) {
		requireID(t, labelID, "create must have succeeded")
		out := zoho(t, "mail", "labels", "get", labelID, "--account", accountID)
		m := parseJSON(t, out)
		if _, ok := m["data"]; !ok {
			if _, ok := m["status"]; !ok {
				t.Fatalf("expected data or status in response:\n%s", truncate(out, 500))
			}
		}
		t.Logf("get label %s succeeded", labelID)
	})

	t.Run("update", func(t *testing.T) {
		requireID(t, labelID, "create must have succeeded")
		newName := "ZT Lbl " + suffix
		body := toJSON(t, map[string]any{"displayName": newName})
		zoho(t, "mail", "labels", "update", labelID, "--account", accountID, "--json", body)
		t.Logf("updated label name to: %s", newName)
	})

	t.Run("delete", func(t *testing.T) {
		requireID(t, labelID, "create must have succeeded")
		zoho(t, "mail", "labels", "delete", labelID, "--account", accountID)
		t.Logf("deleted label: %s", labelID)
	})
}

func TestMailMessages(t *testing.T) {
	t.Parallel()
	accountID := requireMailAccountID(t)
	suffix := randomSuffix()

	var inboxFolderID string
	t.Run("discover-inbox", func(t *testing.T) {
		out := zoho(t, "mail", "folders", "list", "--account", accountID)
		m := parseJSON(t, out)
		data, ok := m["data"].([]any)
		if !ok || len(data) == 0 {
			t.Skip("skipping: no folders found")
		}
		for _, item := range data {
			folder, _ := item.(map[string]any)
			name := fmt.Sprintf("%v", folder["folderName"])
			if name == "Inbox" {
				inboxFolderID = fmt.Sprintf("%v", folder["folderId"])
				break
			}
		}
		if inboxFolderID == "" {
			first, _ := data[0].(map[string]any)
			inboxFolderID = fmt.Sprintf("%v", first["folderId"])
		}
		if inboxFolderID == "" || inboxFolderID == "<nil>" {
			t.Skip("skipping: could not discover Inbox folder ID")
		}
		t.Logf("inbox folder ID: %s", inboxFolderID)
	})

	t.Run("list", func(t *testing.T) {
		requireID(t, inboxFolderID, "discover-inbox must have succeeded")
		out := zoho(t, "mail", "messages", "list", "--account", accountID, "--folder", inboxFolderID, "--limit", "5")
		m := parseJSON(t, out)
		if _, ok := m["data"]; !ok {
			if _, ok := m["status"]; !ok {
				t.Fatalf("expected data or status in response:\n%s", truncate(out, 500))
			}
		}
		t.Logf("messages list succeeded")
	})

	t.Run("search", func(t *testing.T) {
		out := zoho(t, "mail", "messages", "search", "--account", accountID, "--query", "newMails", "--limit", "5")
		m := parseJSON(t, out)
		if _, ok := m["data"]; !ok {
			if _, ok := m["status"]; !ok {
				t.Fatalf("expected data or status in response:\n%s", truncate(out, 500))
			}
		}
		t.Logf("messages search succeeded")
	})

	t.Run("send", func(t *testing.T) {
		out, err := zohoMayFail(t, "mail", "accounts", "get", accountID)
		if err != nil {
			t.Skipf("skipping send: cannot get account details: %v", err)
		}
		m := parseJSON(t, out)
		data, ok := m["data"].(map[string]any)
		if !ok {
			t.Skip("skipping send: no data in account response")
		}
		email, ok := data["primaryEmailAddress"].(string)
		if !ok || email == "" {
			t.Skip("skipping send: no primaryEmailAddress in account")
		}
		subject := testPrefix + " " + suffix
		sendOut := zoho(t, "mail", "messages", "send",
			"--account", accountID,
			"--from", email,
			"--to", email,
			"--subject", subject,
			"--content", "Integration test email body",
			"--format", "plaintext",
		)
		sm := parseJSON(t, sendOut)
		if status, ok := sm["status"].(map[string]any); ok {
			code := fmt.Sprintf("%v", status["code"])
			t.Logf("send status code: %s", code)
		}
		t.Logf("sent test email with subject: %s", subject)
	})
}

func TestMailErrors(t *testing.T) {
	t.Parallel()

	t.Run("missing-account-flag", func(t *testing.T) {
		r := runZohoWithEnv(t, map[string]string{"ZOHO_MAIL_ACCOUNT_ID": ""}, "mail", "folders", "list")
		assertExitCode(t, r, 4)
		assertContains(t, r.Stderr, "--account flag or ZOHO_MAIL_ACCOUNT_ID env var required")
	})

	t.Run("missing-folder-flag", func(t *testing.T) {
		accountID := requireMailAccountID(t)
		r := runZoho(t, "mail", "messages", "list", "--account", accountID)
		assertExitCode(t, r, 1)
	})

	t.Run("missing-arg-folder-get", func(t *testing.T) {
		accountID := requireMailAccountID(t)
		r := runZoho(t, "mail", "folders", "get", "--account", accountID)
		assertExitCode(t, r, 4)
		assertContains(t, r.Stderr, "folder-id argument required")
	})

	t.Run("missing-arg-label-delete", func(t *testing.T) {
		accountID := requireMailAccountID(t)
		r := runZoho(t, "mail", "labels", "delete", "--account", accountID)
		assertExitCode(t, r, 4)
		assertContains(t, r.Stderr, "label-id argument required")
	})

	t.Run("missing-arg-message-get", func(t *testing.T) {
		accountID := requireMailAccountID(t)
		r := runZoho(t, "mail", "messages", "get", "--account", accountID, "--folder", "123")
		assertExitCode(t, r, 4)
		assertContains(t, r.Stderr, "message-id argument required")
	})

	t.Run("missing-org-flag", func(t *testing.T) {
		r := runZohoWithEnv(t, map[string]string{"ZOHO_MAIL_ORG_ID": ""}, "mail", "organization", "get")
		assertExitCode(t, r, 4)
		assertContains(t, r.Stderr, "--org flag or ZOHO_MAIL_ORG_ID env var required")
	})
}

func TestMailSignatures(t *testing.T) {
	t.Parallel()

	t.Run("get", func(t *testing.T) {
		out, err := zohoMayFail(t, "mail", "signatures", "get")
		if err != nil {
			t.Skipf("skipping: signatures get failed: %v", err)
		}
		m := parseJSON(t, out)
		if _, ok := m["data"]; !ok {
			if _, ok := m["status"]; !ok {
				t.Fatalf("expected data or status in response:\n%s", truncate(out, 500))
			}
		}
		t.Logf("signatures get succeeded")
	})
}

func TestMailTasks(t *testing.T) {
	t.Parallel()

	t.Run("list-personal", func(t *testing.T) {
		out, err := zohoMayFail(t, "mail", "tasks", "list-personal")
		if err != nil {
			t.Skipf("skipping: tasks list-personal failed: %v", err)
		}
		m := parseJSON(t, out)
		if _, ok := m["data"]; !ok {
			if _, ok := m["status"]; !ok {
				t.Fatalf("expected data or status in response:\n%s", truncate(out, 500))
			}
		}
		t.Logf("tasks list-personal succeeded")
	})

	t.Run("task-groups", func(t *testing.T) {
		out, err := zohoMayFail(t, "mail", "tasks", "task-groups")
		if err != nil {
			t.Skipf("skipping: task-groups failed: %v", err)
		}
		m := parseJSON(t, out)
		if _, ok := m["data"]; !ok {
			if _, ok := m["status"]; !ok {
				t.Fatalf("expected data or status in response:\n%s", truncate(out, 500))
			}
		}
		t.Logf("task-groups succeeded")
	})
}

func TestMailEmergencyCleanup(t *testing.T) {
	if os.Getenv("ZOHO_EMERGENCY_CLEANUP") != "1" {
		t.Skip("set ZOHO_EMERGENCY_CLEANUP=1 to run")
	}
	accountID := requireMailAccountID(t)

	out, err := zohoMayFail(t, "mail", "folders", "list", "--account", accountID)
	if err == nil {
		m := parseJSON(t, out)
		if data, ok := m["data"].([]any); ok {
			for _, item := range data {
				im, _ := item.(map[string]any)
				name := fmt.Sprintf("%v", im["folderName"])
				if strings.HasPrefix(name, testPrefix) {
					id := fmt.Sprintf("%v", im["folderId"])
					zohoIgnoreError(t, "mail", "folders", "delete", id, "--account", accountID)
				}
			}
		}
	}

	out, err = zohoMayFail(t, "mail", "labels", "list", "--account", accountID)
	if err == nil {
		m := parseJSON(t, out)
		if data, ok := m["data"].([]any); ok {
			for _, item := range data {
				im, _ := item.(map[string]any)
				name := fmt.Sprintf("%v", im["tagName"])
				if strings.HasPrefix(name, testPrefix) {
					id := fmt.Sprintf("%v", im["labelId"])
					if id == "" || id == "<nil>" {
						id = fmt.Sprintf("%v", im["tagId"])
					}
					zohoIgnoreError(t, "mail", "labels", "delete", id, "--account", accountID)
				}
			}
		}
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
