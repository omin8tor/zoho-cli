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

func requireDeskOrgID(t *testing.T) string {
	t.Helper()
	id := os.Getenv("ZOHO_DESK_ORG_ID")
	if id == "" {
		t.Skip("skipping: ZOHO_DESK_ORG_ID not set")
	}
	return id
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
		contactEmail := fmt.Sprintf("%s@zohotest.example.com", randomSuffix())
		out := zoho(t, "desk", "contacts", "create", "--org", orgID,
			"--lastName", contactName,
			"--email", contactEmail)
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
			"--lastName", updatedName)

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
		args := []string{"desk", "tickets", "create", "--org", orgID,
			"--subject", fmt.Sprintf("%v", body["subject"]),
			"--departmentId", fmt.Sprintf("%v", body["departmentId"]),
			"--priority", fmt.Sprintf("%v", body["priority"]),
			"--status", fmt.Sprintf("%v", body["status"]),
		}
		if v, ok := body["contactId"]; ok {
			args = append(args, "--contactId", fmt.Sprintf("%v", v))
		}
		out := zoho(t, args...)
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
			"--subject", updatedSubject,
			"--priority", "High")

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
			"--content", commentContent,
			"--isPublic=false")
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
