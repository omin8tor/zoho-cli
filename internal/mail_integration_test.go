//go:build integration

package internal_test

import (
	"fmt"
	"os"
	"strings"
	"testing"
)

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


