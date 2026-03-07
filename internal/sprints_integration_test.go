//go:build integration

package internal_test

import (
	"fmt"
	"os"
	"strings"
	"testing"
)

func requireSprintsTeamID(t *testing.T) string {
	t.Helper()
	id := os.Getenv("ZOHO_SPRINTS_TEAM_ID")
	if id == "" {
		t.Skip("skipping: ZOHO_SPRINTS_TEAM_ID not set")
	}
	return id
}

func (c *testCleanup) trackSprintsProject(teamID, projectID string) {
	c.add("delete sprints project "+projectID, func() {
		zohoIgnoreError(c.t, "sprints", "projects", "delete", projectID, "--team", teamID)
	})
}

func (c *testCleanup) trackSprintsItem(teamID, projectID, sprintID, itemID string) {
	c.add("delete sprints item "+itemID, func() {
		zohoIgnoreError(c.t, "sprints", "items", "delete", projectID, sprintID, itemID, "--team", teamID)
	})
}

func (c *testCleanup) trackSprintsEpic(teamID, projectID, epicID string) {
	c.add("delete sprints epic "+epicID, func() {
		zohoIgnoreError(c.t, "sprints", "epics", "delete", projectID, epicID, "--team", teamID)
	})
}

func TestSprintsTeams(t *testing.T) {
	t.Parallel()
	_ = requireSprintsTeamID(t)

	t.Run("list", func(t *testing.T) {
		out := zoho(t, "sprints", "teams", "list")
		m := parseJSON(t, out)
		status := fmt.Sprintf("%v", m["status"])
		if status != "success" {
			t.Fatalf("expected status=success, got %q:\n%s", status, truncate(out, 500))
		}
	})
}

func TestSprintsProjects(t *testing.T) {
	t.Parallel()
	teamID := requireSprintsTeamID(t)
	cleanup := newCleanup(t)

	var projectID string
	var ownerID, projGroupID string

	t.Run("list", func(t *testing.T) {
		out := zoho(t, "sprints", "projects", "list", "--team", teamID)
		m := parseJSON(t, out)
		if _, ok := m["projectJObj"]; !ok {
			if status := fmt.Sprintf("%v", m["status"]); status != "success" {
				t.Fatalf("unexpected response:\n%s", truncate(out, 500))
			}
		}
		props, _ := m["project_prop"].(map[string]any)
		projObj, _ := m["projectJObj"].(map[string]any)
		if props != nil && projObj != nil {
			ownerIdx := int(props["owner"].(float64))
			groupIdx := int(props["groupId"].(float64))
			for _, v := range projObj {
				arr, _ := v.([]any)
				if arr != nil && len(arr) > ownerIdx && len(arr) > groupIdx {
					ownerID = fmt.Sprintf("%v", arr[ownerIdx])
					projGroupID = fmt.Sprintf("%v", arr[groupIdx])
					break
				}
			}
		}
	})

	t.Run("create", func(t *testing.T) {
		if ownerID == "" || projGroupID == "" {
			t.Skip("could not discover owner/projgroup from existing projects")
		}
		name := fmt.Sprintf("%s Proj %s", testPrefix, randomSuffix())
		out := zoho(t, "sprints", "projects", "create", "--team", teamID,
			"--name", name,
			"--owner", ownerID,
			"--projgroup", projGroupID)
		m := parseJSON(t, out)
		if pid, ok := m["projectId"]; ok {
			projectID = fmt.Sprintf("%v", pid)
		} else if pids, ok := m["projectIds"]; ok {
			arr, _ := pids.([]any)
			if len(arr) > 0 {
				projectID = fmt.Sprintf("%v", arr[0])
			}
		}
		if projectID == "" || projectID == "<nil>" {
			for k, v := range m {
				if strings.Contains(strings.ToLower(k), "project") && strings.Contains(strings.ToLower(k), "id") {
					projectID = fmt.Sprintf("%v", v)
					break
				}
			}
		}
		if projectID != "" && projectID != "<nil>" {
			cleanup.trackSprintsProject(teamID, projectID)
		}
		status := fmt.Sprintf("%v", m["status"])
		if status != "success" {
			t.Fatalf("expected status=success, got %q:\n%s", status, truncate(out, 500))
		}
	})

	t.Run("get", func(t *testing.T) {
		requireID(t, projectID, "create must have succeeded")
		out := zoho(t, "sprints", "projects", "get", projectID, "--team", teamID)
		_ = parseJSON(t, out)
	})

	t.Run("update", func(t *testing.T) {
		requireID(t, projectID, "create must have succeeded")
		updName := fmt.Sprintf("%s ProjUpd %s", testPrefix, randomSuffix())
		out := zoho(t, "sprints", "projects", "update", projectID, "--team", teamID,
			"--name", updName)
		m := parseJSON(t, out)
		status := fmt.Sprintf("%v", m["status"])
		if status != "success" {
			t.Logf("update response: %s", truncate(out, 500))
		}
	})

	t.Run("delete", func(t *testing.T) {
		requireID(t, projectID, "create must have succeeded")
		out, err := zohoMayFail(t, "sprints", "projects", "delete", projectID, "--team", teamID)
		if err != nil {
			t.Logf("delete may have failed (expected for some project types): %v\n%s", err, truncate(out, 500))
		}
	})
}

func TestSprintsProjectMetadata(t *testing.T) {
	t.Parallel()
	teamID := requireSprintsTeamID(t)

	var projectID string

	t.Run("discover-project", func(t *testing.T) {
		out := zoho(t, "sprints", "projects", "list", "--team", teamID)
		m := parseJSON(t, out)
		if projObj, ok := m["projectJObj"].(map[string]any); ok {
			for id := range projObj {
				projectID = id
				break
			}
		}
		if projectID == "" {
			if projIDs, ok := m["projectIds"].([]any); ok && len(projIDs) > 0 {
				projectID = fmt.Sprintf("%v", projIDs[0])
			}
		}
		if projectID == "" || projectID == "<nil>" {
			t.Skip("no projects found to test metadata endpoints")
		}
	})

	t.Run("statuses", func(t *testing.T) {
		requireID(t, projectID, "discover-project must have succeeded")
		out := zoho(t, "sprints", "statuses", "list", projectID, "--team", teamID)
		_ = parseJSON(t, out)
	})

	t.Run("item-types", func(t *testing.T) {
		requireID(t, projectID, "discover-project must have succeeded")
		out := zoho(t, "sprints", "item-types", "list", projectID, "--team", teamID)
		_ = parseJSON(t, out)
	})

	t.Run("priorities", func(t *testing.T) {
		requireID(t, projectID, "discover-project must have succeeded")
		out := zoho(t, "sprints", "priorities", "list", projectID, "--team", teamID)
		_ = parseJSON(t, out)
	})

	t.Run("members", func(t *testing.T) {
		out := zoho(t, "sprints", "members", "list", "--team", teamID)
		_ = parseJSON(t, out)
	})
}

func TestSprintsSprints(t *testing.T) {
	t.Parallel()
	teamID := requireSprintsTeamID(t)

	var projectID string

	t.Run("discover-project", func(t *testing.T) {
		out := zoho(t, "sprints", "projects", "list", "--team", teamID)
		m := parseJSON(t, out)
		if projObj, ok := m["projectJObj"].(map[string]any); ok {
			for id := range projObj {
				projectID = id
				break
			}
		}
		if projectID == "" {
			if projIDs, ok := m["projectIds"].([]any); ok && len(projIDs) > 0 {
				projectID = fmt.Sprintf("%v", projIDs[0])
			}
		}
		if projectID == "" || projectID == "<nil>" {
			t.Skip("no projects found to test sprints")
		}
	})

	t.Run("list", func(t *testing.T) {
		requireID(t, projectID, "discover-project must have succeeded")
		out := zoho(t, "sprints", "sprints", "list", projectID, "--team", teamID)
		_ = parseJSON(t, out)
	})

	t.Run("list-active", func(t *testing.T) {
		requireID(t, projectID, "discover-project must have succeeded")
		out := zoho(t, "sprints", "sprints", "list", projectID, "--team", teamID, "--type", "2")
		_ = parseJSON(t, out)
	})
}

func TestSprintsItems(t *testing.T) {
	t.Parallel()
	teamID := requireSprintsTeamID(t)
	cleanup := newCleanup(t)

	var projectID string
	var backlogID string

	t.Run("discover-project-and-backlog", func(t *testing.T) {
		out := zoho(t, "sprints", "projects", "list", "--team", teamID)
		m := parseJSON(t, out)
		if projObj, ok := m["projectJObj"].(map[string]any); ok {
			for id, v := range projObj {
				projectID = id
				if arr, ok := v.([]any); ok {
					for i, val := range arr {
						s := fmt.Sprintf("%v", val)
						if strings.Contains(s, "backlog") || i == 0 {
							_ = s
						}
					}
				}
				break
			}
		}
		if projectID == "" {
			if projIDs, ok := m["projectIds"].([]any); ok && len(projIDs) > 0 {
				projectID = fmt.Sprintf("%v", projIDs[0])
			}
		}
		if projectID == "" || projectID == "<nil>" {
			t.Skip("no projects found to test items")
		}

		sprintsOut := zoho(t, "sprints", "sprints", "list", projectID, "--team", teamID)
		sm := parseJSON(t, sprintsOut)
		if sprintIDs, ok := sm["sprintIds"].([]any); ok {
			for _, sid := range sprintIDs {
				idStr := fmt.Sprintf("%v", sid)
				if sprintObj, ok := sm["sprintJObj"].(map[string]any); ok {
					if data, ok := sprintObj[idStr]; ok {
						if arr, ok := data.([]any); ok && len(arr) > 0 {
							name := fmt.Sprintf("%v", arr[0])
							if strings.Contains(strings.ToLower(name), "backlog") {
								backlogID = idStr
								break
							}
						}
					}
				}
			}
			if backlogID == "" && len(sprintIDs) > 0 {
				backlogID = fmt.Sprintf("%v", sprintIDs[0])
			}
		}
		if backlogID == "" || backlogID == "<nil>" {
			t.Skip("no sprints/backlog found")
		}
	})

	t.Run("list", func(t *testing.T) {
		requireID(t, projectID, "discover must have succeeded")
		requireID(t, backlogID, "discover must have succeeded")
		out := zoho(t, "sprints", "items", "list", projectID, backlogID, "--team", teamID)
		_ = parseJSON(t, out)
	})

	var itemTypeID, priorityID string

	t.Run("discover-item-metadata", func(t *testing.T) {
		requireID(t, projectID, "discover must have succeeded")

		itOut := zoho(t, "sprints", "item-types", "list", projectID, "--team", teamID)
		itm := parseJSON(t, itOut)
		if ids, ok := itm["projItemTypeIds"].([]any); ok && len(ids) > 0 {
			itemTypeID = fmt.Sprintf("%v", ids[0])
		}
		if itemTypeID == "" || itemTypeID == "<nil>" {
			t.Skip("no item types found")
		}

		prOut := zoho(t, "sprints", "priorities", "list", projectID, "--team", teamID)
		prm := parseJSON(t, prOut)
		if ids, ok := prm["projPriorityIds"].([]any); ok && len(ids) > 0 {
			priorityID = fmt.Sprintf("%v", ids[0])
		}
		if priorityID == "" || priorityID == "<nil>" {
			t.Skip("no priorities found")
		}
	})

	var itemID string

	t.Run("create", func(t *testing.T) {
		requireID(t, projectID, "discover must have succeeded")
		requireID(t, backlogID, "discover must have succeeded")
		requireID(t, itemTypeID, "discover-item-metadata must have succeeded")
		requireID(t, priorityID, "discover-item-metadata must have succeeded")
		name := fmt.Sprintf("%s Item %s", testPrefix, randomSuffix())
		out := zoho(t, "sprints", "items", "create", projectID, backlogID, "--team", teamID,
			"--name", name,
			"--projitemtypeid", itemTypeID,
			"--projpriorityid", priorityID)
		m := parseJSON(t, out)
		if id, ok := m["itemId"]; ok {
			itemID = fmt.Sprintf("%v", id)
		} else if ids, ok := m["itemIds"]; ok {
			if arr, ok := ids.([]any); ok && len(arr) > 0 {
				itemID = fmt.Sprintf("%v", arr[0])
			}
		}
		if itemID != "" && itemID != "<nil>" {
			cleanup.trackSprintsItem(teamID, projectID, backlogID, itemID)
		}
		status := fmt.Sprintf("%v", m["status"])
		if status != "success" {
			t.Logf("create response: %s", truncate(out, 500))
		}
	})

	t.Run("get", func(t *testing.T) {
		requireID(t, itemID, "create must have succeeded")
		out := zoho(t, "sprints", "items", "get", projectID, backlogID, itemID, "--team", teamID)
		_ = parseJSON(t, out)
	})

	t.Run("update", func(t *testing.T) {
		requireID(t, itemID, "create must have succeeded")
		updName := fmt.Sprintf("%s ItemUpd %s", testPrefix, randomSuffix())
		out, err := zohoMayFail(t, "sprints", "items", "update", projectID, backlogID, itemID, "--team", teamID,
			"--name", updName)
		if err != nil {
			t.Logf("update may have failed: %v\n%s", err, truncate(out, 500))
		}
	})

	t.Run("delete", func(t *testing.T) {
		requireID(t, itemID, "create must have succeeded")
		out, err := zohoMayFail(t, "sprints", "items", "delete", projectID, backlogID, itemID, "--team", teamID)
		if err != nil {
			t.Logf("delete may have failed: %v\n%s", err, truncate(out, 500))
		}
	})
}

func TestSprintsEpics(t *testing.T) {
	t.Parallel()
	teamID := requireSprintsTeamID(t)
	cleanup := newCleanup(t)

	var projectID string

	t.Run("discover-project", func(t *testing.T) {
		out := zoho(t, "sprints", "projects", "list", "--team", teamID)
		m := parseJSON(t, out)
		if projObj, ok := m["projectJObj"].(map[string]any); ok {
			for id := range projObj {
				projectID = id
				break
			}
		}
		if projectID == "" {
			if projIDs, ok := m["projectIds"].([]any); ok && len(projIDs) > 0 {
				projectID = fmt.Sprintf("%v", projIDs[0])
			}
		}
		if projectID == "" || projectID == "<nil>" {
			t.Skip("no projects found to test epics")
		}
	})

	t.Run("list", func(t *testing.T) {
		requireID(t, projectID, "discover-project must have succeeded")
		out := zoho(t, "sprints", "epics", "list", projectID, "--team", teamID)
		_ = parseJSON(t, out)
	})

	var epicID string

	t.Run("create", func(t *testing.T) {
		requireID(t, projectID, "discover-project must have succeeded")
		name := fmt.Sprintf("%s Epic %s", testPrefix, randomSuffix())
		out := zoho(t, "sprints", "epics", "create", projectID, "--team", teamID,
			"--name", name)
		m := parseJSON(t, out)
		if id, ok := m["epicId"]; ok {
			epicID = fmt.Sprintf("%v", id)
		} else if ids, ok := m["epicIds"]; ok {
			if arr, ok := ids.([]any); ok && len(arr) > 0 {
				epicID = fmt.Sprintf("%v", arr[0])
			}
		}
		if epicID != "" && epicID != "<nil>" {
			cleanup.trackSprintsEpic(teamID, projectID, epicID)
		}
		status := fmt.Sprintf("%v", m["status"])
		if status != "success" {
			t.Logf("epic create response: %s", truncate(out, 500))
		}
	})

	t.Run("update", func(t *testing.T) {
		requireID(t, epicID, "create must have succeeded")
		updName := fmt.Sprintf("%s EpicUpd %s", testPrefix, randomSuffix())
		out, err := zohoMayFail(t, "sprints", "epics", "update", projectID, epicID, "--team", teamID,
			"--name", updName)
		if err != nil {
			t.Logf("epic update may have failed: %v\n%s", err, truncate(out, 500))
		}
	})

	t.Run("delete", func(t *testing.T) {
		requireID(t, epicID, "create must have succeeded")
		out, err := zohoMayFail(t, "sprints", "epics", "delete", projectID, epicID, "--team", teamID)
		if err != nil {
			t.Logf("epic delete may have failed: %v\n%s", err, truncate(out, 500))
		}
	})
}

func TestSprintsErrors(t *testing.T) {
	t.Parallel()

	t.Run("missing-team", func(t *testing.T) {
		r := runZohoWithEnv(t, map[string]string{"ZOHO_SPRINTS_TEAM_ID": ""}, "sprints", "projects", "list")
		if r.ExitCode == 0 {
			t.Error("expected non-zero exit code when --team missing")
		}
		if !strings.Contains(r.Stderr, "ZOHO_SPRINTS_TEAM_ID") {
			t.Errorf("expected error mentioning ZOHO_SPRINTS_TEAM_ID, got: %s", r.Stderr)
		}
	})

	t.Run("missing-project-id-for-sprints-list", func(t *testing.T) {
		teamID := os.Getenv("ZOHO_SPRINTS_TEAM_ID")
		if teamID == "" {
			t.Skip("ZOHO_SPRINTS_TEAM_ID not set")
		}
		r := runZoho(t, "sprints", "sprints", "list", "--team", teamID)
		assertExitCode(t, r, 4)
	})

	t.Run("missing-args-items-create", func(t *testing.T) {
		teamID := os.Getenv("ZOHO_SPRINTS_TEAM_ID")
		if teamID == "" {
			t.Skip("ZOHO_SPRINTS_TEAM_ID not set")
		}
		r := runZoho(t, "sprints", "items", "create", "--team", teamID, "--name", "test")
		assertExitCode(t, r, 1)
	})

	t.Run("missing-args-items-get", func(t *testing.T) {
		teamID := os.Getenv("ZOHO_SPRINTS_TEAM_ID")
		if teamID == "" {
			t.Skip("ZOHO_SPRINTS_TEAM_ID not set")
		}
		r := runZoho(t, "sprints", "items", "get", "--team", teamID)
		assertExitCode(t, r, 4)
	})
}

func TestSprintsEmergencyCleanup(t *testing.T) {
	t.Parallel()
	if os.Getenv("ZOHO_EMERGENCY_CLEANUP") != "1" {
		t.Skip("set ZOHO_EMERGENCY_CLEANUP=1 to run")
	}
	teamID := requireSprintsTeamID(t)

	out, err := zohoMayFail(t, "sprints", "projects", "list", "--team", teamID)
	if err != nil {
		return
	}
	m := parseJSON(t, out)
	projObj, ok := m["projectJObj"].(map[string]any)
	if !ok {
		return
	}
	for projectID, v := range projObj {
		arr, ok := v.([]any)
		if !ok || len(arr) == 0 {
			continue
		}
		name := fmt.Sprintf("%v", arr[0])
		if strings.HasPrefix(name, testPrefix) {
			zohoIgnoreError(t, "sprints", "projects", "delete", projectID, "--team", teamID)
		}
	}
}
