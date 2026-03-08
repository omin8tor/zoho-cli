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

func (c *testCleanup) trackBiginRecord(module, id string) {
	c.add("delete bigin "+module+" "+id, func() {
		zohoIgnoreError(c.t, "bigin", "records", "delete", module, id)
	})
}

func (c *testCleanup) trackBiginNote(module, recordID, noteID string) {
	c.add("delete bigin note "+noteID, func() {
		zohoIgnoreError(c.t, "bigin", "notes", "delete", module, recordID, noteID)
	})
}

func (c *testCleanup) trackBiginAttachment(module, recordID, attID string) {
	c.add("delete bigin attachment "+attID, func() {
		zohoIgnoreError(c.t, "bigin", "attachments", "delete", module, recordID, attID)
	})
}

func skipIfBiginUnavailable(t *testing.T) {
	t.Helper()
	out, err := zohoMayFail(t, "bigin", "modules", "list")
	if err != nil {
		t.Skipf("skipping: Bigin not available on this account: %v", err)
	}
	if strings.Contains(out, "<html") || strings.Contains(out, "OAUTH_SCOPE_MISMATCH") {
		t.Skipf("skipping: Bigin not available on this account")
	}
}

func biginExtractID(t *testing.T, out string) string {
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
		t.Fatalf("failed to parse bigin response for ID extraction: %v\nraw: %s", err, truncate(out, 500))
	}
	if len(resp.Data) == 0 {
		t.Fatalf("no data in bigin response:\n%s", truncate(out, 500))
	}
	if resp.Data[0].Status != "success" {
		t.Fatalf("bigin operation was not successful:\n%s", truncate(out, 500))
	}
	id := resp.Data[0].Details.ID
	if id == "" {
		t.Fatalf("empty ID in bigin response:\n%s", truncate(out, 500))
	}
	return id
}

func biginAssertStatus(t *testing.T, out string, want string) {
	t.Helper()
	var resp struct {
		Data []struct {
			Status string `json:"status"`
		} `json:"data"`
	}
	if err := json.Unmarshal([]byte(out), &resp); err != nil {
		t.Fatalf("failed to parse bigin response: %v\nraw: %s", err, truncate(out, 500))
	}
	if len(resp.Data) == 0 || resp.Data[0].Status != want {
		t.Errorf("expected bigin status %q in response:\n%s", want, truncate(out, 500))
	}
}

func TestBiginModules(t *testing.T) {
	t.Parallel()
	skipIfBiginUnavailable(t)

	t.Run("list", func(t *testing.T) {
		out := zoho(t, "bigin", "modules", "list")
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
			modMap, _ := mod.(map[string]any)
			if n, ok := modMap["api_name"].(string); ok {
				names[n] = true
			}
		}
		for _, want := range []string{"Contacts", "Pipelines"} {
			if !names[want] {
				t.Errorf("expected module %s in list", want)
			}
		}
	})

	t.Run("fields", func(t *testing.T) {
		out := zoho(t, "bigin", "modules", "fields", "Contacts")
		m := parseJSON(t, out)
		fields, ok := m["fields"].([]any)
		if !ok {
			t.Fatalf("expected fields array in response:\n%s", truncate(out, 500))
		}
		if len(fields) == 0 {
			t.Fatal("expected at least one field for Contacts")
		}
	})

	t.Run("layouts", func(t *testing.T) {
		out := zoho(t, "bigin", "modules", "layouts", "Contacts")
		m := parseJSON(t, out)
		layouts, ok := m["layouts"].([]any)
		if !ok {
			t.Fatalf("expected layouts array in response:\n%s", truncate(out, 500))
		}
		if len(layouts) == 0 {
			t.Fatal("expected at least one layout for Contacts")
		}
	})

	t.Run("related-lists", func(t *testing.T) {
		out := zoho(t, "bigin", "modules", "related-lists", "Contacts")
		m := parseJSON(t, out)
		relLists, ok := m["related_lists"].([]any)
		if !ok {
			t.Fatalf("expected related_lists array in response:\n%s", truncate(out, 500))
		}
		if len(relLists) == 0 {
			t.Fatal("expected at least one related list for Contacts")
		}
	})
}

func TestBiginOrg(t *testing.T) {
	t.Parallel()
	skipIfBiginUnavailable(t)

	t.Run("get", func(t *testing.T) {
		out := zoho(t, "bigin", "org", "get")
		m := parseJSON(t, out)
		orgs, ok := m["org"].([]any)
		if !ok {
			t.Fatalf("expected org array in response:\n%s", truncate(out, 500))
		}
		if len(orgs) == 0 {
			t.Fatal("expected at least one organization")
		}
	})
}

func TestBiginRoles(t *testing.T) {
	t.Parallel()
	skipIfBiginUnavailable(t)

	t.Run("list", func(t *testing.T) {
		out := zoho(t, "bigin", "roles", "list")
		m := parseJSON(t, out)
		roles, ok := m["roles"].([]any)
		if !ok {
			t.Fatalf("expected roles array in response:\n%s", truncate(out, 500))
		}
		if len(roles) == 0 {
			t.Fatal("expected at least one role")
		}
	})
}

func TestBiginProfiles(t *testing.T) {
	t.Parallel()
	skipIfBiginUnavailable(t)

	t.Run("list", func(t *testing.T) {
		out := zoho(t, "bigin", "profiles", "list")
		m := parseJSON(t, out)
		profiles, ok := m["profiles"].([]any)
		if !ok {
			t.Fatalf("expected profiles array in response:\n%s", truncate(out, 500))
		}
		if len(profiles) == 0 {
			t.Fatal("expected at least one profile")
		}
	})
}

func TestBiginUsers(t *testing.T) {
	t.Parallel()
	skipIfBiginUnavailable(t)

	var userID string

	t.Run("list", func(t *testing.T) {
		out := zoho(t, "bigin", "users", "list")
		m := parseJSON(t, out)
		users, ok := m["users"].([]any)
		if !ok {
			t.Fatalf("expected users array in response:\n%s", truncate(out, 500))
		}
		if len(users) == 0 {
			t.Fatal("expected at least one user")
		}
		first := users[0].(map[string]any)
		userID = fmt.Sprintf("%v", first["id"])
	})

	t.Run("get", func(t *testing.T) {
		requireID(t, userID, "list must have succeeded")
		out := zoho(t, "bigin", "users", "get", userID)
		m := parseJSON(t, out)
		users, ok := m["users"].([]any)
		if !ok || len(users) == 0 {
			t.Fatalf("expected users array with user in response:\n%s", truncate(out, 500))
		}
		user := users[0].(map[string]any)
		assertEqual(t, fmt.Sprintf("%v", user["id"]), userID)
	})
}

func TestBiginContactsCRUD(t *testing.T) {
	t.Parallel()
	skipIfBiginUnavailable(t)
	cleanup := newCleanup(t)

	var contactID string
	var contactLastName string

	t.Run("create", func(t *testing.T) {
		contactLastName = fmt.Sprintf("%s BiginCt %s", testPrefix, randomSuffix())
		out := zoho(t, "bigin", "records", "create", "Contacts",
			"--json", toJSON(t, map[string]any{
				"Last_Name": contactLastName,
				"Email":     strings.ToLower(fmt.Sprintf("bigin_%s@test.example.com", randomSuffix())),
			}))
		contactID = biginExtractID(t, out)
		cleanup.trackBiginRecord("Contacts", contactID)
	})

	t.Run("get", func(t *testing.T) {
		requireID(t, contactID, "create must have succeeded")
		out := zoho(t, "bigin", "records", "get", "Contacts", contactID)
		m := parseJSON(t, out)
		data, ok := m["data"].([]any)
		if !ok || len(data) == 0 {
			t.Fatalf("expected data array in response:\n%s", truncate(out, 500))
		}
		rec := data[0].(map[string]any)
		assertEqual(t, fmt.Sprintf("%v", rec["id"]), contactID)
		assertEqual(t, fmt.Sprintf("%v", rec["Last_Name"]), contactLastName)
	})

	t.Run("list", func(t *testing.T) {
		out := zoho(t, "bigin", "records", "list", "Contacts", "--limit", "5")
		arr := parseJSONArray(t, out)
		if len(arr) == 0 {
			t.Fatal("expected at least one contact")
		}
	})

	t.Run("update", func(t *testing.T) {
		requireID(t, contactID, "create must have succeeded")
		updatedName := fmt.Sprintf("%s BiginCtUpd %s", testPrefix, randomSuffix())
		out := zoho(t, "bigin", "records", "update", "Contacts", contactID,
			"--json", toJSON(t, map[string]any{"Last_Name": updatedName}))
		biginAssertStatus(t, out, "success")

		getOut := zoho(t, "bigin", "records", "get", "Contacts", contactID)
		gm := parseJSON(t, getOut)
		gdata := gm["data"].([]any)
		rec := gdata[0].(map[string]any)
		assertEqual(t, fmt.Sprintf("%v", rec["Last_Name"]), updatedName)
	})

	t.Run("search-by-criteria", func(t *testing.T) {
		requireID(t, contactID, "create must have succeeded")
		retryUntil(t, 30*time.Second, func() bool {
			out, err := zohoMayFail(t, "bigin", "search", "Contacts",
				"--criteria", fmt.Sprintf("(Last_Name:starts_with:%s)", testPrefix))
			if err != nil {
				return false
			}
			m := parseJSON(t, out)
			data, ok := m["data"].([]any)
			return ok && len(data) > 0
		})
	})

	t.Run("delete", func(t *testing.T) {
		requireID(t, contactID, "create must have succeeded")
		out := zoho(t, "bigin", "records", "delete", "Contacts", contactID)
		biginAssertStatus(t, out, "success")
		contactID = ""
	})
}

func TestBiginPipelinesCRUD(t *testing.T) {
	t.Parallel()
	skipIfBiginUnavailable(t)
	cleanup := newCleanup(t)

	var pipelineRecordID string
	var pipelineRecordName string

	t.Run("list", func(t *testing.T) {
		out := zoho(t, "bigin", "records", "list", "Pipelines", "--limit", "5")
		arr := parseJSONArray(t, out)
		if len(arr) == 0 {
			t.Log("no existing pipeline records")
		}
	})

	t.Run("create", func(t *testing.T) {
		pipelineRecordName = fmt.Sprintf("%s Deal %s", testPrefix, randomSuffix())
		out := zoho(t, "bigin", "records", "create", "Pipelines",
			"--json", toJSON(t, map[string]any{
				"Pipeline_Name":  pipelineRecordName,
				"Pipeline_Stage": "Qualification",
			}))
		pipelineRecordID = biginExtractID(t, out)
		cleanup.trackBiginRecord("Pipelines", pipelineRecordID)
	})

	t.Run("get", func(t *testing.T) {
		requireID(t, pipelineRecordID, "create must have succeeded")
		out := zoho(t, "bigin", "records", "get", "Pipelines", pipelineRecordID)
		m := parseJSON(t, out)
		data, ok := m["data"].([]any)
		if !ok || len(data) == 0 {
			t.Fatalf("expected data array:\n%s", truncate(out, 500))
		}
		rec := data[0].(map[string]any)
		assertEqual(t, fmt.Sprintf("%v", rec["id"]), pipelineRecordID)
	})

	t.Run("update", func(t *testing.T) {
		requireID(t, pipelineRecordID, "create must have succeeded")
		updatedName := fmt.Sprintf("%s DealUpd %s", testPrefix, randomSuffix())
		out := zoho(t, "bigin", "records", "update", "Pipelines", pipelineRecordID,
			"--json", toJSON(t, map[string]any{"Pipeline_Name": updatedName}))
		biginAssertStatus(t, out, "success")
	})

	t.Run("delete", func(t *testing.T) {
		requireID(t, pipelineRecordID, "create must have succeeded")
		out := zoho(t, "bigin", "records", "delete", "Pipelines", pipelineRecordID)
		biginAssertStatus(t, out, "success")
		pipelineRecordID = ""
	})
}

func TestBiginNotes(t *testing.T) {
	t.Parallel()
	skipIfBiginUnavailable(t)
	cleanup := newCleanup(t)

	var contactID string
	var noteID string

	t.Run("create-contact", func(t *testing.T) {
		name := fmt.Sprintf("%s NotesCt %s", testPrefix, randomSuffix())
		out := zoho(t, "bigin", "records", "create", "Contacts",
			"--json", toJSON(t, map[string]any{"Last_Name": name}))
		contactID = biginExtractID(t, out)
		cleanup.trackBiginRecord("Contacts", contactID)
	})

	t.Run("add", func(t *testing.T) {
		requireID(t, contactID, "create-contact must have succeeded")
		out := zoho(t, "bigin", "notes", "add", "Contacts", contactID,
			"--content", "Bigin integration test note",
			"--title", "Test Note")
		noteID = biginExtractID(t, out)
		cleanup.trackBiginNote("Contacts", contactID, noteID)
	})

	t.Run("list", func(t *testing.T) {
		requireID(t, contactID, "create-contact must have succeeded")
		requireID(t, noteID, "add must have succeeded")
		out := zoho(t, "bigin", "notes", "list", "Contacts", contactID)
		m := parseJSON(t, out)
		data, ok := m["data"].([]any)
		if !ok || len(data) == 0 {
			t.Fatalf("expected at least one note:\n%s", truncate(out, 500))
		}
		found := false
		for _, item := range data {
			note := item.(map[string]any)
			if fmt.Sprintf("%v", note["id"]) == noteID {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("note %s not found in list", noteID)
		}
	})

	t.Run("update", func(t *testing.T) {
		requireID(t, contactID, "create-contact must have succeeded")
		requireID(t, noteID, "add must have succeeded")
		out := zoho(t, "bigin", "notes", "update", "Contacts", contactID, noteID,
			"--content", "Updated bigin note content",
			"--title", "Updated Title")
		biginAssertStatus(t, out, "success")
	})

	t.Run("delete", func(t *testing.T) {
		requireID(t, contactID, "create-contact must have succeeded")
		requireID(t, noteID, "add must have succeeded")
		out := zoho(t, "bigin", "notes", "delete", "Contacts", contactID, noteID)
		biginAssertStatus(t, out, "success")
		noteID = ""
	})
}

func TestBiginAttachments(t *testing.T) {
	t.Parallel()
	skipIfBiginUnavailable(t)
	cleanup := newCleanup(t)

	var contactID string
	var attachmentID string
	var testFileContent []byte

	t.Run("create-contact", func(t *testing.T) {
		name := fmt.Sprintf("%s AttCt %s", testPrefix, randomSuffix())
		out := zoho(t, "bigin", "records", "create", "Contacts",
			"--json", toJSON(t, map[string]any{"Last_Name": name}))
		contactID = biginExtractID(t, out)
		cleanup.trackBiginRecord("Contacts", contactID)
	})

	t.Run("upload", func(t *testing.T) {
		requireID(t, contactID, "create-contact must have succeeded")
		tmpDir := t.TempDir()
		testFile := tmpDir + "/bigin-test.txt"
		testFileContent = []byte("bigin-cli integration test " + time.Now().String())
		if err := os.WriteFile(testFile, testFileContent, 0644); err != nil {
			t.Fatalf("failed to create test file: %v", err)
		}
		out := zoho(t, "bigin", "attachments", "upload", "Contacts", contactID, testFile)
		attachmentID = biginExtractID(t, out)
		cleanup.trackBiginAttachment("Contacts", contactID, attachmentID)
	})

	t.Run("list", func(t *testing.T) {
		requireID(t, contactID, "create-contact must have succeeded")
		requireID(t, attachmentID, "upload must have succeeded")
		out := zoho(t, "bigin", "attachments", "list", "Contacts", contactID)
		m := parseJSON(t, out)
		data, ok := m["data"].([]any)
		if !ok || len(data) == 0 {
			t.Fatalf("expected at least one attachment:\n%s", truncate(out, 500))
		}
		found := false
		for _, item := range data {
			att := item.(map[string]any)
			if fmt.Sprintf("%v", att["id"]) == attachmentID {
				found = true
				assertEqual(t, fmt.Sprintf("%v", att["File_Name"]), "bigin-test.txt")
				break
			}
		}
		if !found {
			t.Errorf("attachment %s not found in list", attachmentID)
		}
	})

	t.Run("download", func(t *testing.T) {
		requireID(t, contactID, "create-contact must have succeeded")
		requireID(t, attachmentID, "upload must have succeeded")
		tmpDir := t.TempDir()
		downloadPath := tmpDir + "/downloaded.txt"
		zoho(t, "bigin", "attachments", "download", "Contacts", contactID, attachmentID,
			"--output", downloadPath)
		downloaded, err := os.ReadFile(downloadPath)
		if err != nil {
			t.Fatalf("failed to read downloaded file: %v", err)
		}
		if len(downloaded) == 0 {
			t.Error("downloaded file is empty")
		}
	})

	t.Run("delete", func(t *testing.T) {
		requireID(t, contactID, "create-contact must have succeeded")
		requireID(t, attachmentID, "upload must have succeeded")
		out := zoho(t, "bigin", "attachments", "delete", "Contacts", contactID, attachmentID)
		biginAssertStatus(t, out, "success")
		attachmentID = ""
	})
}

func TestBiginTags(t *testing.T) {
	t.Parallel()
	skipIfBiginUnavailable(t)
	cleanup := newCleanup(t)

	var contactID string

	t.Run("create-contact", func(t *testing.T) {
		name := fmt.Sprintf("%s TagCt %s", testPrefix, randomSuffix())
		out := zoho(t, "bigin", "records", "create", "Contacts",
			"--json", toJSON(t, map[string]any{"Last_Name": name}))
		contactID = biginExtractID(t, out)
		cleanup.trackBiginRecord("Contacts", contactID)
	})

	t.Run("add", func(t *testing.T) {
		requireID(t, contactID, "create-contact must have succeeded")
		out := zoho(t, "bigin", "tags", "add", "Contacts",
			"--ids", contactID, "--tags", "bigintest-tag-a,bigintest-tag-b")
		biginAssertStatus(t, out, "success")
	})

	t.Run("verify-tags", func(t *testing.T) {
		requireID(t, contactID, "create-contact must have succeeded")
		out := zoho(t, "bigin", "records", "get", "Contacts", contactID,
			"--fields", "id,Tag")
		m := parseJSON(t, out)
		data, ok := m["data"].([]any)
		if !ok || len(data) == 0 {
			t.Fatalf("expected data:\n%s", truncate(out, 500))
		}
		rec := data[0].(map[string]any)
		tags, ok := rec["Tag"].([]any)
		if !ok {
			t.Fatalf("expected Tag array in record:\n%s", truncate(out, 500))
		}
		tagNames := make(map[string]bool)
		for _, tag := range tags {
			tagMap := tag.(map[string]any)
			tagNames[fmt.Sprintf("%v", tagMap["name"])] = true
		}
		if !tagNames["bigintest-tag-a"] {
			t.Error("tag 'bigintest-tag-a' not found after add")
		}
		if !tagNames["bigintest-tag-b"] {
			t.Error("tag 'bigintest-tag-b' not found after add")
		}
	})

	t.Run("remove", func(t *testing.T) {
		requireID(t, contactID, "create-contact must have succeeded")
		out := zoho(t, "bigin", "tags", "remove", "Contacts",
			"--ids", contactID, "--tags", "bigintest-tag-a,bigintest-tag-b")
		biginAssertStatus(t, out, "success")
	})
}

func TestBiginUpsert(t *testing.T) {
	t.Parallel()
	skipIfBiginUnavailable(t)
	cleanup := newCleanup(t)

	var contactID string
	email := fmt.Sprintf("bigin_upsert_%s@test.example.com", randomSuffix())

	t.Run("insert", func(t *testing.T) {
		name := fmt.Sprintf("%s UpsertCt %s", testPrefix, randomSuffix())
		out := zoho(t, "bigin", "records", "upsert", "Contacts",
			"--json", toJSON(t, map[string]any{"Last_Name": name, "Email": email}),
			"--duplicate-check", "Email")
		contactID = biginExtractID(t, out)
		cleanup.trackBiginRecord("Contacts", contactID)
	})

	t.Run("update-via-upsert", func(t *testing.T) {
		requireID(t, contactID, "insert must have succeeded")
		updatedName := fmt.Sprintf("%s UpsertUpd %s", testPrefix, randomSuffix())
		out := zoho(t, "bigin", "records", "upsert", "Contacts",
			"--json", toJSON(t, map[string]any{"Last_Name": updatedName, "Email": email}),
			"--duplicate-check", "Email")
		var resp struct {
			Data []struct {
				Action string `json:"action"`
			} `json:"data"`
		}
		if err := json.Unmarshal([]byte(out), &resp); err == nil && len(resp.Data) > 0 {
			if resp.Data[0].Action != "update" {
				t.Logf("expected upsert action 'update', got %q", resp.Data[0].Action)
			}
		}
	})
}

func TestBiginRelated(t *testing.T) {
	t.Parallel()
	skipIfBiginUnavailable(t)
	cleanup := newCleanup(t)

	var contactID string
	var noteID string

	t.Run("create-contact", func(t *testing.T) {
		name := fmt.Sprintf("%s RelCt %s", testPrefix, randomSuffix())
		out := zoho(t, "bigin", "records", "create", "Contacts",
			"--json", toJSON(t, map[string]any{"Last_Name": name}))
		contactID = biginExtractID(t, out)
		cleanup.trackBiginRecord("Contacts", contactID)
	})

	t.Run("add-note", func(t *testing.T) {
		requireID(t, contactID, "create-contact must have succeeded")
		out := zoho(t, "bigin", "notes", "add", "Contacts", contactID,
			"--content", "Related list test note")
		noteID = biginExtractID(t, out)
		cleanup.trackBiginNote("Contacts", contactID, noteID)
	})

	t.Run("list-related-notes", func(t *testing.T) {
		requireID(t, contactID, "create-contact must have succeeded")
		requireID(t, noteID, "add-note must have succeeded")
		out := zoho(t, "bigin", "related", "list", "Contacts", contactID, "Notes")
		m := parseJSON(t, out)
		data, ok := m["data"].([]any)
		if !ok || len(data) == 0 {
			t.Fatalf("expected related notes:\n%s", truncate(out, 500))
		}
		found := false
		for _, item := range data {
			note := item.(map[string]any)
			if fmt.Sprintf("%v", note["id"]) == noteID {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("note %s not found in related list", noteID)
		}
	})
}

func TestBiginCOQL(t *testing.T) {
	t.Parallel()
	skipIfBiginUnavailable(t)
	cleanup := newCleanup(t)

	var contactID string

	t.Run("create-contact", func(t *testing.T) {
		name := fmt.Sprintf("%s CoqlCt %s", testPrefix, randomSuffix())
		out := zoho(t, "bigin", "records", "create", "Contacts",
			"--json", toJSON(t, map[string]any{"Last_Name": name}))
		contactID = biginExtractID(t, out)
		cleanup.trackBiginRecord("Contacts", contactID)
	})

	t.Run("query", func(t *testing.T) {
		requireID(t, contactID, "create-contact must have succeeded")
		retryUntil(t, 15*time.Second, func() bool {
			query := fmt.Sprintf("select id, Last_Name from Contacts where id = '%s'", contactID)
			out, err := zohoMayFail(t, "bigin", "coql", "--query", query)
			if err != nil {
				return false
			}
			m := parseJSON(t, out)
			data, ok := m["data"].([]any)
			return ok && len(data) > 0
		})
	})

	t.Run("like-query", func(t *testing.T) {
		requireID(t, contactID, "create-contact must have succeeded")
		retryUntil(t, 15*time.Second, func() bool {
			query := fmt.Sprintf("select id, Last_Name from Contacts where Last_Name like '%s%%' limit 5", testPrefix)
			out, err := zohoMayFail(t, "bigin", "coql", "--query", query)
			if err != nil {
				return false
			}
			m := parseJSON(t, out)
			data, ok := m["data"].([]any)
			if !ok || len(data) == 0 {
				return false
			}
			for _, item := range data {
				rec := item.(map[string]any)
				name := fmt.Sprintf("%v", rec["Last_Name"])
				if !strings.HasPrefix(name, testPrefix) {
					t.Errorf("LIKE '%s%%' returned record with Last_Name=%q", testPrefix, name)
				}
			}
			return true
		})
	})
}

func TestBiginRecordCount(t *testing.T) {
	t.Parallel()
	skipIfBiginUnavailable(t)

	t.Run("count-contacts", func(t *testing.T) {
		out, err := zohoMayFail(t, "bigin", "records", "count", "Contacts")
		if err != nil {
			t.Skipf("record count may not be supported: %v", err)
		}
		m := parseJSON(t, out)
		if _, ok := m["count"]; !ok {
			t.Logf("count response: %s", truncate(out, 500))
		}
	})
}

func TestBiginBulkDelete(t *testing.T) {
	t.Parallel()
	skipIfBiginUnavailable(t)
	cleanup := newCleanup(t)

	name1 := fmt.Sprintf("%s Bulk1 %s", testPrefix, randomSuffix())
	out1 := zoho(t, "bigin", "records", "create", "Contacts",
		"--json", toJSON(t, map[string]any{"Last_Name": name1}))
	id1 := biginExtractID(t, out1)
	cleanup.trackBiginRecord("Contacts", id1)

	name2 := fmt.Sprintf("%s Bulk2 %s", testPrefix, randomSuffix())
	out2 := zoho(t, "bigin", "records", "create", "Contacts",
		"--json", toJSON(t, map[string]any{"Last_Name": name2}))
	id2 := biginExtractID(t, out2)
	cleanup.trackBiginRecord("Contacts", id2)

	t.Run("bulk-delete", func(t *testing.T) {
		ids := id1 + "," + id2
		out := zoho(t, "bigin", "records", "bulk-delete", "Contacts", ids)
		var resp struct {
			Data []struct {
				Status string `json:"status"`
			} `json:"data"`
		}
		if err := json.Unmarshal([]byte(out), &resp); err != nil {
			t.Fatalf("failed to parse bulk-delete response: %v", err)
		}
		if len(resp.Data) != 2 {
			t.Errorf("expected 2 results in bulk-delete, got %d", len(resp.Data))
		}
		for i, d := range resp.Data {
			if d.Status != "success" {
				t.Errorf("bulk-delete item %d: expected success, got %q", i, d.Status)
			}
		}
	})
}

func TestBiginErrors(t *testing.T) {
	t.Parallel()

	t.Run("missing-module-arg", func(t *testing.T) {
		r := runZoho(t, "bigin", "records", "list")
		assertExitCode(t, r, 4)
		assertContains(t, r.Stderr, "module")
	})

	t.Run("missing-record-id", func(t *testing.T) {
		r := runZoho(t, "bigin", "records", "get", "Contacts")
		assertExitCode(t, r, 4)
	})

	t.Run("missing-json-flag", func(t *testing.T) {
		r := runZoho(t, "bigin", "records", "create", "Contacts")
		if r.ExitCode == 0 {
			t.Error("expected non-zero exit code when --json missing")
		}
	})

	t.Run("invalid-json", func(t *testing.T) {
		r := runZoho(t, "bigin", "records", "create", "Contacts", "--json", "not json")
		assertExitCode(t, r, 4)
		assertContains(t, r.Stderr, "invalid JSON")
	})

	t.Run("nonexistent-record", func(t *testing.T) {
		r := runZoho(t, "bigin", "records", "get", "Contacts", "999999999999999999")
		if r.ExitCode == 0 {
			t.Log("warning: API may return success for nonexistent records")
		}
	})

	t.Run("missing-note-args", func(t *testing.T) {
		r := runZoho(t, "bigin", "notes", "list", "Contacts")
		assertExitCode(t, r, 4)
	})

	t.Run("missing-attachment-args", func(t *testing.T) {
		r := runZoho(t, "bigin", "attachments", "list", "Contacts")
		assertExitCode(t, r, 4)
	})
}

func TestBiginEmergencyCleanup(t *testing.T) {
	t.Parallel()
	if os.Getenv("ZOHO_EMERGENCY_CLEANUP") != "1" {
		t.Skip("set ZOHO_EMERGENCY_CLEANUP=1 to run")
	}

	out, err := zohoMayFail(t, "bigin", "coql", "--query",
		fmt.Sprintf("select id from Contacts where Last_Name like '%s%%'", testPrefix))
	if err == nil {
		var resp struct {
			Data []struct {
				ID string `json:"id"`
			} `json:"data"`
		}
		if jsonErr := json.Unmarshal([]byte(out), &resp); jsonErr == nil {
			t.Logf("found %d orphaned bigin contacts", len(resp.Data))
			for _, rec := range resp.Data {
				t.Logf("deleting orphaned bigin contact %s", rec.ID)
				zohoIgnoreError(t, "bigin", "records", "delete", "Contacts", rec.ID)
			}
		}
	}

	out, err = zohoMayFail(t, "bigin", "coql", "--query",
		fmt.Sprintf("select id from Pipelines where Pipeline_Name like '%s%%'", testPrefix))
	if err == nil {
		var resp struct {
			Data []struct {
				ID string `json:"id"`
			} `json:"data"`
		}
		if jsonErr := json.Unmarshal([]byte(out), &resp); jsonErr == nil {
			t.Logf("found %d orphaned bigin pipeline records", len(resp.Data))
			for _, rec := range resp.Data {
				t.Logf("deleting orphaned bigin pipeline record %s", rec.ID)
				zohoIgnoreError(t, "bigin", "records", "delete", "Pipelines", rec.ID)
			}
		}
	}
}
