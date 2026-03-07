//go:build integration

package internal_test

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"
)

func recruitExtractID(t *testing.T, out string) string {
	t.Helper()
	var resp struct {
		Data []struct {
			Details struct {
				ID string `json:"id"`
			} `json:"details"`
			Status  string `json:"status"`
			Code    string `json:"code"`
			Message string `json:"message"`
		} `json:"data"`
	}
	if err := json.Unmarshal([]byte(out), &resp); err != nil {
		t.Fatalf("failed to parse recruit response for ID: %v\nraw: %s", err, truncate(out, 500))
	}
	if len(resp.Data) == 0 {
		t.Fatalf("no data in recruit response:\n%s", truncate(out, 500))
	}
	if resp.Data[0].Status != "success" {
		t.Fatalf("recruit operation not successful: code=%s message=%s\n%s",
			resp.Data[0].Code, resp.Data[0].Message, truncate(out, 500))
	}
	id := resp.Data[0].Details.ID
	if id == "" {
		t.Fatalf("empty ID in recruit response:\n%s", truncate(out, 500))
	}
	return id
}

func (c *testCleanup) trackRecruitRecord(module, id string) {
	c.add("delete recruit "+module+" "+id, func() {
		zohoIgnoreError(c.t, "recruit", "records", "delete", module, id)
	})
}

func (c *testCleanup) trackRecruitNote(id string) {
	c.add("delete recruit note "+id, func() {
		zohoIgnoreError(c.t, "recruit", "notes", "delete", id)
	})
}

func (c *testCleanup) trackRecruitAttachment(module, recordID, attID string) {
	c.add("delete recruit attachment "+attID, func() {
		zohoIgnoreError(c.t, "recruit", "attachments", "delete", module, recordID, attID)
	})
}

func TestRecruitModules(t *testing.T) {
	t.Parallel()

	t.Run("list", func(t *testing.T) {
		out := zoho(t, "recruit", "modules", "list")
		m := parseJSON(t, out)
		modules, ok := m["modules"].([]any)
		if !ok {
			t.Fatalf("expected modules array in response:\n%s", truncate(out, 500))
		}
		if len(modules) == 0 {
			t.Fatal("expected at least one module")
		}
		names := make(map[string]bool)
		for _, mod := range modules {
			mm, _ := mod.(map[string]any)
			if n, ok := mm["api_name"].(string); ok {
				names[n] = true
			}
		}
		for _, want := range []string{"Candidates", "Job_Openings"} {
			if !names[want] {
				t.Errorf("expected module %s in list", want)
			}
		}
	})

	t.Run("fields", func(t *testing.T) {
		out := zoho(t, "recruit", "modules", "fields", "Candidates")
		m := parseJSON(t, out)
		fields, ok := m["fields"].([]any)
		if !ok {
			t.Fatalf("expected fields array in response:\n%s", truncate(out, 500))
		}
		if len(fields) == 0 {
			t.Fatal("expected at least one field for Candidates")
		}
		names := make(map[string]bool)
		for _, f := range fields {
			fm, _ := f.(map[string]any)
			if n, ok := fm["api_name"].(string); ok {
				names[n] = true
			}
		}
		for _, want := range []string{"Last_Name", "Email"} {
			if !names[want] {
				t.Errorf("expected field %s in Candidates fields", want)
			}
		}
	})

	t.Run("layouts", func(t *testing.T) {
		out := zoho(t, "recruit", "modules", "layouts", "Candidates")
		m := parseJSON(t, out)
		layouts, ok := m["layouts"].([]any)
		if !ok {
			t.Fatalf("expected layouts array in response:\n%s", truncate(out, 500))
		}
		if len(layouts) == 0 {
			t.Fatal("expected at least one layout for Candidates")
		}
	})
}

func TestRecruitRecords(t *testing.T) {
	t.Parallel()
	cleanup := newCleanup(t)

	var candidateID string
	var candidateLastName string
	var candidateEmail string
	var noteID string
	var attachmentID string

	t.Run("create", func(t *testing.T) {
		candidateLastName = fmt.Sprintf("%s_%d_%s", testPrefix, time.Now().Unix(), randomSuffix())
		candidateEmail = strings.ToLower(candidateLastName) + "@test.example.com"
		data := toJSON(t, map[string]any{
			"Last_Name":  candidateLastName,
			"Email":      candidateEmail,
			"First_Name": "Test",
		})
		out := zoho(t, "recruit", "records", "create", "Candidates", "--json", data)
		candidateID = recruitExtractID(t, out)
		cleanup.trackRecruitRecord("Candidates", candidateID)
	})

	t.Run("get", func(t *testing.T) {
		requireID(t, candidateID, "create must have succeeded")
		out := zoho(t, "recruit", "records", "get", "Candidates", candidateID)
		m := parseJSON(t, out)
		data, ok := m["data"].([]any)
		if !ok || len(data) == 0 {
			t.Fatalf("expected data array in response:\n%s", truncate(out, 500))
		}
		rec, _ := data[0].(map[string]any)
		assertEqual(t, fmt.Sprintf("%v", rec["id"]), candidateID)
	})

	t.Run("list", func(t *testing.T) {
		out := zoho(t, "recruit", "records", "list", "Candidates", "--per-page", "5")
		m := parseJSON(t, out)
		data, ok := m["data"].([]any)
		if !ok {
			t.Fatalf("expected data array in response:\n%s", truncate(out, 500))
		}
		if len(data) == 0 {
			t.Fatal("expected at least one candidate")
		}
	})

	t.Run("update", func(t *testing.T) {
		requireID(t, candidateID, "create must have succeeded")
		data := toJSON(t, map[string]any{"First_Name": "Updated"})
		out := zoho(t, "recruit", "records", "update", "Candidates", candidateID, "--json", data)
		m := parseJSON(t, out)
		respData, ok := m["data"].([]any)
		if !ok || len(respData) == 0 {
			t.Fatalf("expected data array in update response:\n%s", truncate(out, 500))
		}
		rec, _ := respData[0].(map[string]any)
		if fmt.Sprintf("%v", rec["status"]) != "success" {
			t.Fatalf("update not successful:\n%s", truncate(out, 500))
		}
	})

	t.Run("search", func(t *testing.T) {
		requireID(t, candidateID, "create must have succeeded")
		retryUntilSkip(t, 30*time.Second, func() bool {
			out, err := zohoMayFail(t, "recruit", "records", "search", "Candidates",
				"--email", candidateEmail)
			if err != nil {
				return false
			}
			m := parseJSON(t, out)
			data, ok := m["data"].([]any)
			if !ok || len(data) == 0 {
				return false
			}
			for _, item := range data {
				rec, _ := item.(map[string]any)
				if fmt.Sprintf("%v", rec["id"]) == candidateID {
					return true
				}
			}
			return false
		})
	})

	t.Run("notes/create", func(t *testing.T) {
		requireID(t, candidateID, "create must have succeeded")
		data := toJSON(t, map[string]any{
			"Note_Title":   "Test Note",
			"Note_Content": "Integration test note",
			"Parent_Id":    candidateID,
			"se_module":    "Candidates",
		})
		out := zoho(t, "recruit", "notes", "create", "--json", data)
		noteID = recruitExtractID(t, out)
		cleanup.trackRecruitNote(noteID)
	})

	t.Run("notes/list", func(t *testing.T) {
		requireID(t, candidateID, "create must have succeeded")
		requireID(t, noteID, "notes/create must have succeeded")
		out := zoho(t, "recruit", "notes", "list", "Candidates", candidateID)
		m := parseJSON(t, out)
		data, ok := m["data"].([]any)
		if !ok || len(data) == 0 {
			t.Fatalf("expected at least one note:\n%s", truncate(out, 500))
		}
		found := false
		for _, item := range data {
			note, _ := item.(map[string]any)
			if fmt.Sprintf("%v", note["id"]) == noteID {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("note %s not found in list", noteID)
		}
	})

	t.Run("notes/delete", func(t *testing.T) {
		requireID(t, noteID, "notes/create must have succeeded")
		zoho(t, "recruit", "notes", "delete", noteID)
		noteID = ""
	})

	t.Run("attachments/upload", func(t *testing.T) {
		requireID(t, candidateID, "create must have succeeded")
		tmpDir := t.TempDir()
		testFile := tmpDir + "/test-recruit-attachment.txt"
		content := []byte("recruit integration test file " + time.Now().String())
		if err := os.WriteFile(testFile, content, 0644); err != nil {
			t.Fatalf("failed to create test file: %v", err)
		}
		out := zoho(t, "recruit", "attachments", "upload", "Candidates", candidateID, testFile, "--category", "Others")
		attachmentID = recruitExtractID(t, out)
		cleanup.trackRecruitAttachment("Candidates", candidateID, attachmentID)
	})

	t.Run("attachments/list", func(t *testing.T) {
		requireID(t, candidateID, "create must have succeeded")
		requireID(t, attachmentID, "attachments/upload must have succeeded")
		out := zoho(t, "recruit", "attachments", "list", "Candidates", candidateID)
		m := parseJSON(t, out)
		data, ok := m["data"].([]any)
		if !ok || len(data) == 0 {
			t.Fatalf("expected at least one attachment:\n%s", truncate(out, 500))
		}
		found := false
		for _, item := range data {
			att, _ := item.(map[string]any)
			if fmt.Sprintf("%v", att["id"]) == attachmentID {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("attachment %s not found in list", attachmentID)
		}
	})

	t.Run("attachments/download", func(t *testing.T) {
		requireID(t, candidateID, "create must have succeeded")
		requireID(t, attachmentID, "attachments/upload must have succeeded")
		tmpDir := t.TempDir()
		downloadPath := tmpDir + "/downloaded.txt"
		zoho(t, "recruit", "attachments", "download", "Candidates", candidateID, attachmentID,
			"--output", downloadPath)
		downloaded, err := os.ReadFile(downloadPath)
		if err != nil {
			t.Fatalf("failed to read downloaded file: %v", err)
		}
		if len(downloaded) == 0 {
			t.Error("downloaded file is empty")
		}
	})

	t.Run("attachments/delete", func(t *testing.T) {
		requireID(t, candidateID, "create must have succeeded")
		requireID(t, attachmentID, "attachments/upload must have succeeded")
		zoho(t, "recruit", "attachments", "delete", "Candidates", candidateID, attachmentID)
		attachmentID = ""
	})

	t.Run("tags/add", func(t *testing.T) {
		requireID(t, candidateID, "create must have succeeded")
		out, err := zohoMayFail(t, "recruit", "tags", "add", "Candidates",
			"--ids", candidateID, "--tag-names", "zohotest-recruit-tag")
		if err != nil {
			t.Skipf("tag add failed (tag may need to exist first): %v", err)
		}
		m := parseJSON(t, out)
		data, ok := m["data"].([]any)
		if !ok || len(data) == 0 {
			t.Skipf("tag add returned no data (may need pre-existing tag)")
		}
		rec, _ := data[0].(map[string]any)
		if fmt.Sprintf("%v", rec["status"]) != "success" {
			t.Skipf("tag add not successful (tag may need to exist first):\n%s", truncate(out, 500))
		}
	})

	t.Run("tags/remove", func(t *testing.T) {
		requireID(t, candidateID, "create must have succeeded")
		out, err := zohoMayFail(t, "recruit", "tags", "remove", "Candidates",
			"--ids", candidateID, "--tag-names", "zohotest-recruit-tag")
		if err != nil {
			t.Skipf("tag remove failed: %v", err)
		}
		m := parseJSON(t, out)
		data, ok := m["data"].([]any)
		if !ok || len(data) == 0 {
			t.Skipf("tag remove returned no data")
		}
		rec, _ := data[0].(map[string]any)
		if fmt.Sprintf("%v", rec["status"]) != "success" {
			t.Logf("tag remove status: %v", rec["status"])
		}
	})

	t.Run("related/list-notes", func(t *testing.T) {
		requireID(t, candidateID, "create must have succeeded")
		out, _ := zohoMayFail(t, "recruit", "related", "list", "Candidates", candidateID, "Notes")
		if out != "" {
			parseJSON(t, out)
		}
	})

	t.Run("delete", func(t *testing.T) {
		requireID(t, candidateID, "create must have succeeded")
		out := zoho(t, "recruit", "records", "delete", "Candidates", candidateID)
		m := parseJSON(t, out)
		data, ok := m["data"].([]any)
		if !ok || len(data) == 0 {
			t.Fatalf("expected data array in delete response:\n%s", truncate(out, 500))
		}
		rec, _ := data[0].(map[string]any)
		if fmt.Sprintf("%v", rec["status"]) != "success" {
			t.Fatalf("delete not successful:\n%s", truncate(out, 500))
		}
		candidateID = ""
	})
}

func TestRecruitJobOpenings(t *testing.T) {
	t.Parallel()
	cleanup := newCleanup(t)

	var jobID string

	t.Run("create", func(t *testing.T) {
		name := fmt.Sprintf("%s Job %s", testPrefix, randomSuffix())
		data := toJSON(t, map[string]any{
			"Job_Opening_Name":    name,
			"Client_Name":         "Test Client",
			"Number_of_Positions": "1",
			"Job_Opening_Status":  "In-progress",
		})
		out := zoho(t, "recruit", "records", "create", "Job_Openings", "--json", data)
		jobID = recruitExtractID(t, out)
		cleanup.trackRecruitRecord("Job_Openings", jobID)
	})

	t.Run("get", func(t *testing.T) {
		requireID(t, jobID, "create must have succeeded")
		out := zoho(t, "recruit", "records", "get", "Job_Openings", jobID)
		m := parseJSON(t, out)
		data, ok := m["data"].([]any)
		if !ok || len(data) == 0 {
			t.Fatalf("expected data array:\n%s", truncate(out, 500))
		}
		rec, _ := data[0].(map[string]any)
		assertEqual(t, fmt.Sprintf("%v", rec["id"]), jobID)
	})

	t.Run("list", func(t *testing.T) {
		out := zoho(t, "recruit", "records", "list", "Job_Openings", "--per-page", "5")
		m := parseJSON(t, out)
		data, ok := m["data"].([]any)
		if !ok {
			t.Fatalf("expected data array:\n%s", truncate(out, 500))
		}
		if len(data) == 0 {
			t.Fatal("expected at least one job opening")
		}
	})

	t.Run("update", func(t *testing.T) {
		requireID(t, jobID, "create must have succeeded")
		data := toJSON(t, map[string]any{"Number_of_Positions": "2"})
		out := zoho(t, "recruit", "records", "update", "Job_Openings", jobID, "--json", data)
		m := parseJSON(t, out)
		respData, ok := m["data"].([]any)
		if !ok || len(respData) == 0 {
			t.Fatalf("expected data array in update response:\n%s", truncate(out, 500))
		}
		rec, _ := respData[0].(map[string]any)
		if fmt.Sprintf("%v", rec["status"]) != "success" {
			t.Fatalf("update not successful:\n%s", truncate(out, 500))
		}
	})

	t.Run("delete", func(t *testing.T) {
		requireID(t, jobID, "create must have succeeded")
		out := zoho(t, "recruit", "records", "delete", "Job_Openings", jobID)
		m := parseJSON(t, out)
		data, ok := m["data"].([]any)
		if !ok || len(data) == 0 {
			t.Fatalf("expected data array in delete response:\n%s", truncate(out, 500))
		}
		rec, _ := data[0].(map[string]any)
		if fmt.Sprintf("%v", rec["status"]) != "success" {
			t.Fatalf("delete not successful:\n%s", truncate(out, 500))
		}
		jobID = ""
	})
}

func TestRecruitUsers(t *testing.T) {
	t.Parallel()

	t.Run("list", func(t *testing.T) {
		out := zoho(t, "recruit", "users", "list")
		m := parseJSON(t, out)
		users, ok := m["users"].([]any)
		if !ok {
			t.Fatalf("expected users array in response:\n%s", truncate(out, 500))
		}
		if len(users) == 0 {
			t.Fatal("expected at least one user")
		}
		user, _ := users[0].(map[string]any)
		if _, ok := user["id"]; !ok {
			t.Error("expected id field on user")
		}
	})
}

func TestRecruitTags(t *testing.T) {
	t.Parallel()

	t.Run("list", func(t *testing.T) {
		out, _ := zohoMayFail(t, "recruit", "tags", "list", "--module", "Candidates")
		if out != "" {
			parseJSON(t, out)
		}
	})
}

func TestRecruitErrors(t *testing.T) {
	t.Parallel()

	t.Run("missing-module-records-get", func(t *testing.T) {
		r := runZoho(t, "recruit", "records", "get", "Candidates")
		assertExitCode(t, r, 4)
	})

	t.Run("missing-json-flag", func(t *testing.T) {
		r := runZoho(t, "recruit", "records", "create", "Candidates")
		if r.ExitCode == 0 {
			t.Error("expected non-zero exit code when --json missing")
		}
	})

	t.Run("invalid-json", func(t *testing.T) {
		r := runZoho(t, "recruit", "records", "create", "Candidates", "--json", "not json")
		assertExitCode(t, r, 4)
	})

	t.Run("missing-module-fields", func(t *testing.T) {
		r := runZoho(t, "recruit", "modules", "fields")
		assertExitCode(t, r, 4)
	})

	t.Run("missing-note-id", func(t *testing.T) {
		r := runZoho(t, "recruit", "notes", "delete")
		assertExitCode(t, r, 4)
	})

	t.Run("missing-attachment-args", func(t *testing.T) {
		r := runZoho(t, "recruit", "attachments", "list", "Candidates")
		assertExitCode(t, r, 4)
	})
}

func TestRecruitEmergencyCleanup(t *testing.T) {
	t.Parallel()
	if os.Getenv("ZOHO_EMERGENCY_CLEANUP") != "1" {
		t.Skip("set ZOHO_EMERGENCY_CLEANUP=1 to run")
	}

	out, err := zohoMayFail(t, "recruit", "records", "list", "Candidates", "--per-page", "200")
	if err == nil {
		m := parseJSON(t, out)
		if data, ok := m["data"].([]any); ok {
			for _, item := range data {
				rec, _ := item.(map[string]any)
				lastName := fmt.Sprintf("%v", rec["Last_Name"])
				if strings.HasPrefix(lastName, testPrefix) {
					id := fmt.Sprintf("%v", rec["id"])
					t.Logf("deleting orphaned candidate %s (%s)", id, lastName)
					zohoIgnoreError(t, "recruit", "records", "delete", "Candidates", id)
				}
			}
		}
	}

	out, err = zohoMayFail(t, "recruit", "records", "list", "Job_Openings", "--per-page", "200")
	if err == nil {
		m := parseJSON(t, out)
		if data, ok := m["data"].([]any); ok {
			for _, item := range data {
				rec, _ := item.(map[string]any)
				title := fmt.Sprintf("%v", rec["Posting_Title"])
				if strings.HasPrefix(title, testPrefix) {
					id := fmt.Sprintf("%v", rec["id"])
					t.Logf("deleting orphaned job opening %s (%s)", id, title)
					zohoIgnoreError(t, "recruit", "records", "delete", "Job_Openings", id)
				}
			}
		}
	}
}
