//go:build integration

package internal_test

import (
	"fmt"
	"os"
	"strings"
	"testing"
)

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
		out := zoho(t, "sign", "requests", "list", "--limit", "10")
		requests := parseJSONArray(t, out)
		found := false
		for _, r := range requests {
			if fmt.Sprintf("%v", r["request_id"]) == requestID {
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
			"--limit", "2",
			"--sort-column", "created_time",
			"--sort-order", "DESC")
		arr := parseJSONArray(t, out)
		if len(arr) > 2 {
			t.Fatalf("expected at most 2 requests with --limit 2, got %d", len(arr))
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

	out, err := zohoMayFail(t, "sign", "requests", "list", "--limit", "100")
	if err == nil {
		arr := parseJSONArray(t, out)
		for _, im := range arr {
			name := fmt.Sprintf("%v", im["request_name"])
			if strings.HasPrefix(name, testPrefix) {
				id := fmt.Sprintf("%v", im["request_id"])
				zohoIgnoreError(t, "sign", "requests", "delete", id)
			}
		}
	}

	out, err = zohoMayFail(t, "sign", "templates", "list", "--limit", "100")
	if err == nil {
		arr := parseJSONArray(t, out)
		for _, im := range arr {
			name := fmt.Sprintf("%v", im["template_name"])
			if strings.HasPrefix(name, testPrefix) {
				id := fmt.Sprintf("%v", im["template_id"])
				zohoIgnoreError(t, "sign", "templates", "delete", id)
			}
		}
	}
}
