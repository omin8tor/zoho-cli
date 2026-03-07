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


