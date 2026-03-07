//go:build ignore

package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strings"
)

type fileSpec struct {
	name   string
	ranges [][2]int
}

func main() {
	inputFile := "internal/integration_test.go.bak"
	lines, err := readLines(inputFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error reading %s: %v\n", inputFile, err)
		os.Exit(1)
	}

	specs := []fileSpec{
		{
			name: "internal/helpers_test.go",
			ranges: [][2]int{
				{18, 18},     // const testPrefix
				{20, 28},     // TestMain
				{30, 35},     // type Result
				{36, 117},    // runZoho, runZohoWithEnv, zoho, zohoMayFail, zohoIgnoreError
				{119, 275},   // toJSON, parseJSON*, extractID, requireID, retryUntil, assertEqual, assert*, truncate
				{332, 356},   // cleanupEntry, testCleanup, newCleanup, add
				{564, 575},   // randomSuffix, testName
				{1872, 1891}, // TestHelpAll
				{5241, 5247}, // assertExpenseCodeZero (used only by expense, but avoids redecl)
				{6232, 6247}, // assertSheetSuccess, trackSheetWorkbook (used only by sheet)
				{7798, 7803}, // trackSignRequest (used only by sign)
			},
		},
		{
			name: "internal/crm_integration_test.go",
			ranges: [][2]int{
				{277, 330},  // getRecord, getRecordMayFail, getNotes, getAttachments, findInArray, hasTag
				{358, 375},  // trackLead, trackNote, trackAttachment
				{442, 452},  // trackContact, trackAccount
				{508, 562},  // convertResult, extractConvertIDs
				{909, 1870}, // TestCRMModules..TestCRMEmergencyCleanup
			},
		},
		{
			name: "internal/drive_integration_test.go",
			ranges: [][2]int{
				{577, 677},   // driveTestParentFolder, driveAttr, extractDrive*, getDriveFile*, assertDriveAttr, trackDriveFile, trackDriveFolder, requireDriveTeamID
				{1892, 2309}, // TestDriveTeams..TestDriveEmergencyCleanup
			},
		},
		{
			name: "internal/projects_integration_test.go",
			ranges: [][2]int{
				{679, 686},   // requireProjectsPortalID
				{760, 907},   // extractProjectsID, getProject, getTask, getIssue, getTasklist, getMilestone, getTimelog, trackProject..trackDashboard
				{2310, 5240}, // TestProjectsPortals..TestProjectsEmergencyCleanup
			},
		},
		{
			name: "internal/expense_integration_test.go",
			ranges: [][2]int{
				{376, 422},   // trackExpenseCategory..trackExpenseExpense
				{688, 696},   // requireExpenseOrgID
				{5248, 6231}, // TestExpenseOrganizations..TestExpenseEmergencyCleanup
			},
		},
		{
			name: "internal/sheet_integration_test.go",
			ranges: [][2]int{
				{6248, 6977}, // TestSheet* tests
			},
		},
		{
			name: "internal/cliq_integration_test.go",
			ranges: [][2]int{
				{454, 458},   // trackCliqChannel
				{6978, 7280}, // TestCliq* tests
			},
		},
		{
			name: "internal/writer_integration_test.go",
			ranges: [][2]int{
				{472, 476},   // trackWriterDoc
				{7281, 7406}, // TestWriter* tests
			},
		},
		{
			name: "internal/desk_integration_test.go",
			ranges: [][2]int{
				{424, 440},   // trackDeskTicket, trackDeskContact, trackDeskAccount
				{697, 705},   // requireDeskOrgID
				{7407, 7797}, // TestDesk* tests
			},
		},
		{
			name: "internal/sign_integration_test.go",
			ranges: [][2]int{
				{7804, 8283}, // TestSign* tests
			},
		},
		{
			name: "internal/mail_integration_test.go",
			ranges: [][2]int{
				{460, 470},   // trackMailFolder, trackMailLabel
				{706, 727},   // requireMailAccountID
				{8284, 8689}, // TestMail* tests
			},
		},
		{
			name: "internal/books_integration_test.go",
			ranges: [][2]int{
				{478, 506},   // trackBooksContact..trackBooksExpense
				{729, 758},   // requireBooksOrgID, assertBooksCodeZero
				{8690, 9284}, // TestBooks* tests
			},
		},
	}

	for _, spec := range specs {
		writeModuleFile(spec, lines)
	}

	fmt.Println("\nDone! Verify with: go test -tags integration -count=1 -run TestXXX ./internal/")
}

func writeModuleFile(spec fileSpec, allLines []string) {
	usedImports := map[string]bool{}

	for _, r := range spec.ranges {
		start := r[0] - 1
		end := r[1]
		if end > len(allLines) {
			end = len(allLines)
		}
		for i := start; i < end; i++ {
			line := allLines[i]
			if strings.Contains(line, "json.") {
				usedImports["encoding/json"] = true
			}
			if strings.Contains(line, "fmt.") {
				usedImports["fmt"] = true
			}
			if strings.Contains(line, "os.") {
				usedImports["os"] = true
			}
			if strings.Contains(line, "strings.") {
				usedImports["strings"] = true
			}
			if strings.Contains(line, "time.") {
				usedImports["time"] = true
			}
			if strings.Contains(line, "bytes.") {
				usedImports["bytes"] = true
			}
			if strings.Contains(line, "context.") {
				usedImports["context"] = true
			}
			if strings.Contains(line, "exec.") {
				usedImports["os/exec"] = true
			}
			if strings.Contains(line, "rand.") {
				usedImports["crypto/rand"] = true
			}
			if strings.Contains(line, "filepath.") {
				usedImports["path/filepath"] = true
			}
		}
	}
	usedImports["testing"] = true

	f, err := os.Create(spec.name)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error creating %s: %v\n", spec.name, err)
		os.Exit(1)
	}
	defer f.Close()

	fmt.Fprintln(f, "//go:build integration")
	fmt.Fprintln(f, "")
	fmt.Fprintln(f, "package internal_test")
	fmt.Fprintln(f, "")

	imports := make([]string, 0, len(usedImports))
	for imp := range usedImports {
		imports = append(imports, imp)
	}
	sort.Strings(imports)

	if len(imports) == 1 {
		fmt.Fprintf(f, "import %q\n", imports[0])
	} else {
		fmt.Fprintln(f, "import (")
		for _, imp := range imports {
			fmt.Fprintf(f, "\t%q\n", imp)
		}
		fmt.Fprintln(f, ")")
	}
	fmt.Fprintln(f, "")

	for _, r := range spec.ranges {
		start := r[0] - 1
		end := r[1]
		if end > len(allLines) {
			end = len(allLines)
		}
		for i := start; i < end; i++ {
			fmt.Fprintln(f, allLines[i])
		}
		fmt.Fprintln(f, "")
	}

	fmt.Printf("wrote %s\n", spec.name)
}

func readLines(path string) ([]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	var lines []string
	scanner := bufio.NewScanner(f)
	buf := make([]byte, 0, 1024*1024)
	scanner.Buffer(buf, 1024*1024)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}
