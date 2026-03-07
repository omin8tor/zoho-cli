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


func assertExpenseCodeZero(t *testing.T, m map[string]any) {
	t.Helper()
	if fmt.Sprintf("%v", m["code"]) != "0" {
		t.Fatalf("expected expense code=0, got %v", m["code"])
	}
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


func (c *testCleanup) trackSignRequest(id string) {
	c.add("delete sign request "+id, func() {
		zohoIgnoreError(c.t, "sign", "requests", "delete", id)
	})
}


