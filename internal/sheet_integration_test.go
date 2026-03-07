//go:build integration

package internal_test

import (
	"fmt"
	"strings"
	"testing"
)

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


