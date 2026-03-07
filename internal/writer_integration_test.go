//go:build integration

package internal_test

import (
	"fmt"
	"os"
	"testing"
)

func (c *testCleanup) trackWriterDoc(id string) {
	c.add("trash writer doc "+id, func() {
		zohoIgnoreError(c.t, "drive", "files", "trash", id)
	})
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


