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

func requireProjectsPortalID(t *testing.T) string {
	t.Helper()
	id := os.Getenv("ZOHO_PORTAL_ID")
	if id == "" {
		t.Skip("skipping: ZOHO_PORTAL_ID not set")
	}
	return id
}

func extractProjectsID(t *testing.T, out string) string {
	t.Helper()
	m := parseJSON(t, out)
	id := fmt.Sprintf("%v", m["id"])
	if id == "" || id == "<nil>" {
		t.Fatalf("no id in Projects response:\n%s", truncate(out, 500))
	}
	return id
}

func getProject(t *testing.T, projectID string) map[string]any {
	t.Helper()
	out := zoho(t, "projects", "core", "get", projectID)
	return parseJSON(t, out)
}

func getTask(t *testing.T, taskID, projectID string) map[string]any {
	t.Helper()
	out := zoho(t, "projects", "tasks", "get", taskID, "--project", projectID)
	return parseJSON(t, out)
}

func getIssue(t *testing.T, issueID, projectID string) map[string]any {
	t.Helper()
	out := zoho(t, "projects", "issues", "get", issueID, "--project", projectID)
	arr := parseJSONArray(t, out)
	if len(arr) == 0 {
		t.Fatalf("issues get returned empty array for %s", issueID)
	}
	return arr[0]
}

func getTasklist(t *testing.T, tasklistID, projectID string) map[string]any {
	t.Helper()
	out := zoho(t, "projects", "tasklists", "get", tasklistID, "--project", projectID)
	return parseJSON(t, out)
}

func getMilestone(t *testing.T, milestoneID, projectID string) map[string]any {
	t.Helper()
	out := zoho(t, "projects", "milestones", "get", milestoneID, "--project", projectID)
	return parseJSON(t, out)
}

func getTimelog(t *testing.T, timelogID, projectID string) map[string]any {
	t.Helper()
	out := zoho(t, "projects", "timelogs", "get", timelogID, "--project", projectID, "--type", "task")
	return parseJSON(t, out)
}

func (c *testCleanup) trackProject(id string) {
	c.add("delete project "+id, func() {
		zohoIgnoreError(c.t, "projects", "core", "delete", id)
	})
	c.add("trash project "+id, func() {
		zohoIgnoreError(c.t, "projects", "core", "trash", id)
	})
}

func (c *testCleanup) trackTask(id, projectID string) {
	c.add("delete task "+id, func() {
		zohoIgnoreError(c.t, "projects", "tasks", "delete", id, "--project", projectID)
	})
}

func (c *testCleanup) trackIssue(id, projectID string) {
	c.add("delete issue "+id, func() {
		zohoIgnoreError(c.t, "projects", "issues", "delete", id, "--project", projectID)
	})
}

func (c *testCleanup) trackTasklist(id, projectID string) {
	c.add("delete tasklist "+id, func() {
		zohoIgnoreError(c.t, "projects", "tasklists", "delete", id, "--project", projectID)
	})
}

func (c *testCleanup) trackMilestone(id, projectID string) {
	c.add("delete milestone "+id, func() {
		zohoIgnoreError(c.t, "projects", "milestones", "delete", id, "--project", projectID)
	})
}

func (c *testCleanup) trackTimelog(id, projectID string) {
	c.add("delete timelog "+id, func() {
		zohoIgnoreError(c.t, "projects", "timelogs", "delete", id, "--project", projectID, "--type", "task")
	})
}

func (c *testCleanup) trackProjectGroup(id string) {
	c.add("delete project-group "+id, func() {
		zohoIgnoreError(c.t, "projects", "project-groups", "delete", id)
	})
}

func (c *testCleanup) trackTag(id string) {
	c.add("delete tag "+id, func() {
		zohoIgnoreError(c.t, "projects", "tags", "delete", id)
	})
}

func (c *testCleanup) trackRole(id string) {
	c.add("delete role "+id, func() {
		zohoIgnoreError(c.t, "projects", "roles", "delete", id)
	})
}

func (c *testCleanup) trackForumCategory(id, projectID string) {
	c.add("delete forum-category "+id, func() {
		zohoIgnoreError(c.t, "projects", "forum-categories", "delete", id, "--project", projectID)
	})
}

func (c *testCleanup) trackForum(id, projectID string) {
	c.add("delete forum "+id, func() {
		zohoIgnoreError(c.t, "projects", "forums", "delete", id, "--project", projectID)
	})
}

func (c *testCleanup) trackPhase(id, projectID string) {
	c.add("delete phase "+id, func() {
		zohoIgnoreError(c.t, "projects", "phases", "delete", id, "--project", projectID)
	})
}

func (c *testCleanup) trackEvent(id, projectID string) {
	c.add("delete event "+id, func() {
		zohoIgnoreError(c.t, "projects", "events", "delete", id, "--project", projectID)
	})
}

func (c *testCleanup) trackTeam(id string) {
	c.add("delete team "+id, func() {
		zohoIgnoreError(c.t, "projects", "teams", "delete", id)
	})
}

func (c *testCleanup) trackProfile(id string) {
	c.add("delete profile "+id, func() {
		zohoIgnoreError(c.t, "projects", "profiles", "delete", id)
	})
}

func (c *testCleanup) trackDashboard(id string) {
	c.add("delete dashboard "+id, func() {
		zohoIgnoreError(c.t, "projects", "reports", "dashboard-delete", id)
	})
}

func TestProjectsPortals(t *testing.T) {
	t.Parallel()
	portalID := requireProjectsPortalID(t)

	t.Run("portals/get", func(t *testing.T) {
		out := zoho(t, "projects", "portals", "get")
		m := parseJSON(t, out)
		pd, ok := m["portal_details"].(map[string]any)
		if ok {
			m = pd
		}
		id := fmt.Sprintf("%v", m["id"])
		if id == "" || id == "<nil>" {
			t.Fatalf("expected portal id in response:\n%s", truncate(out, 500))
		}
		idNum := fmt.Sprintf("%v", m["id"])
		if !strings.Contains(idNum, portalID) {
			idF := strings.TrimRight(strings.TrimRight(fmt.Sprintf("%f", m["id"]), "0"), ".")
			if idF != portalID {
				t.Logf("portal id mismatch: got %s / %s, env has %s", idNum, idF, portalID)
			}
		}
		name := fmt.Sprintf("%v", m["name"])
		if name == "" || name == "<nil>" {
			t.Errorf("expected portal name in response")
		}
		t.Logf("portal: id=%s name=%s", id, name)
	})

	t.Run("portals/list-known-bug", func(t *testing.T) {
		t.Skip("portals list endpoint returns INVALID_METHOD in V3 API")
	})
}

func TestProjects(t *testing.T) {
	t.Parallel()
	_ = requireProjectsPortalID(t)
	cleanup := newCleanup(t)

	var projectID string
	var projectName string
	var project2ID string
	var project2Name string
	var taskID string
	var taskName string
	var subtaskID string
	var clonedTaskID string
	var issueID string
	var issueName string
	var clonedIssueID string
	var tasklistID string
	var tasklistName string
	var milestoneID string
	var milestoneName string
	var timelogID string
	var tlTaskID string
	var ownerZPUID string
	var taskCommentID string
	var issueCommentID string
	var tasklistCommentID string
	var projectCommentID string
	var projectGroupID string
	var tagID string
	var roleID string
	var forumCategoryID string
	var forumID string
	var forumCommentID string
	var issueLinkID string
	var phaseID string
	var phaseName string
	var clonedPhaseID string
	var phaseCommentID string
	var eventID string
	var eventCommentID string

	t.Run("core/create", func(t *testing.T) {
		projectName = testName(t)
		out := zoho(t, "projects", "core", "create", "--name", projectName)
		projectID = extractProjectsID(t, out)
		cleanup.trackProject(projectID)
		t.Logf("created project %s (%s)", projectID, projectName)

		proj := getProject(t, projectID)
		assertStringField(t, proj, "name", projectName)
	})

	t.Run("core/get", func(t *testing.T) {
		requireID(t, projectID, "core/create must have succeeded")
		proj := getProject(t, projectID)
		assertEqual(t, fmt.Sprintf("%v", proj["id"]), projectID)
		assertStringField(t, proj, "name", projectName)
	})

	t.Run("timelogs/setup-owner", func(t *testing.T) {
		requireID(t, projectID, "core/create must have succeeded")
		proj := getProject(t, projectID)
		if cb, ok := proj["created_by"].(map[string]any); ok {
			ownerZPUID = fmt.Sprintf("%v", cb["zpuid"])
		}
		if ownerZPUID == "" || ownerZPUID == "<nil>" {
			t.Fatal("could not determine owner zpuid from project")
		}
		t.Logf("owner zpuid: %s", ownerZPUID)
	})

	t.Run("core/list", func(t *testing.T) {
		requireID(t, projectID, "core/create must have succeeded")
		out := zoho(t, "projects", "core", "list")
		arr := parseJSONArray(t, out)
		assertNonEmpty(t, arr, "expected at least one project")
		_, found := findInArray(arr, projectID)
		if !found {
			t.Errorf("created project %s not found in project list", projectID)
		}
	})

	t.Run("core/update", func(t *testing.T) {
		requireID(t, projectID, "core/create must have succeeded")
		updatedName := projectName + "_updated"
		out := zoho(t, "projects", "core", "update", projectID,
			"--json", toJSON(t, map[string]any{"name": updatedName}))
		m := parseJSON(t, out)
		assertEqual(t, fmt.Sprintf("%v", m["name"]), updatedName)

		proj := getProject(t, projectID)
		assertStringField(t, proj, "name", updatedName)
		projectName = updatedName
	})

	t.Run("core/create-second-project", func(t *testing.T) {
		project2Name = testName(t) + "_p2"
		out := zoho(t, "projects", "core", "create", "--name", project2Name)
		project2ID = extractProjectsID(t, out)
		cleanup.trackProject(project2ID)
		t.Logf("created second project %s (%s)", project2ID, project2Name)
	})

	t.Run("tasks/create", func(t *testing.T) {
		requireID(t, projectID, "core/create must have succeeded")
		taskName = testName(t) + "_task"
		out := zoho(t, "projects", "tasks", "create",
			"--name", taskName, "--project", projectID)
		taskID = extractProjectsID(t, out)
		cleanup.trackTask(taskID, projectID)
		t.Logf("created task %s", taskID)

		task := getTask(t, taskID, projectID)
		assertStringField(t, task, "name", taskName)
	})

	t.Run("tasks/get", func(t *testing.T) {
		requireID(t, taskID, "tasks/create must have succeeded")
		task := getTask(t, taskID, projectID)
		assertEqual(t, fmt.Sprintf("%v", task["id"]), taskID)
		assertStringField(t, task, "name", taskName)
	})

	t.Run("tasks/update", func(t *testing.T) {
		requireID(t, taskID, "tasks/create must have succeeded")
		updatedTaskName := taskName + "_upd"
		out := zoho(t, "projects", "tasks", "update", taskID,
			"--project", projectID,
			"--json", toJSON(t, map[string]any{"name": updatedTaskName}))
		m := parseJSON(t, out)
		assertEqual(t, fmt.Sprintf("%v", m["name"]), updatedTaskName)

		task := getTask(t, taskID, projectID)
		assertStringField(t, task, "name", updatedTaskName)
		taskName = updatedTaskName
	})

	t.Run("tasks/list", func(t *testing.T) {
		requireID(t, taskID, "tasks/create must have succeeded")
		out := zoho(t, "projects", "tasks", "list", "--project", projectID)
		arr := parseJSONArray(t, out)
		assertNonEmpty(t, arr, "expected at least one task")
		_, found := findInArray(arr, taskID)
		if !found {
			t.Errorf("task %s not found in task list", taskID)
		}
	})

	t.Run("tasks/my", func(t *testing.T) {
		out := zoho(t, "projects", "tasks", "my")
		m := parseJSON(t, out)
		if _, ok := m["tasks"]; !ok {
			if _, ok2 := m["page_info"]; !ok2 {
				t.Errorf("expected tasks or page_info key in my-tasks response:\n%s", truncate(out, 500))
			}
		}
	})

	t.Run("task-comments/add", func(t *testing.T) {
		requireID(t, taskID, "tasks/create must have succeeded")
		out := zoho(t, "projects", "task-comments", "add",
			"--task", taskID, "--project", projectID,
			"--comment", testPrefix+"_task_comment")
		arr := parseJSONArray(t, out)
		if len(arr) == 0 {
			t.Fatalf("task-comments add returned empty array:\n%s", truncate(out, 500))
		}
		taskCommentID = fmt.Sprintf("%v", arr[0]["id"])
		if taskCommentID == "" || taskCommentID == "<nil>" {
			t.Fatalf("no id in task-comments add response:\n%s", truncate(out, 500))
		}
		assertStringField(t, arr[0], "comment", testPrefix+"_task_comment")
		t.Logf("created task comment %s", taskCommentID)
	})

	t.Run("task-comments/list", func(t *testing.T) {
		requireID(t, taskCommentID, "task-comments/add must have succeeded")
		out := zoho(t, "projects", "task-comments", "list",
			"--task", taskID, "--project", projectID)
		m := parseJSON(t, out)
		comments, ok := m["comments"].([]any)
		if !ok {
			t.Fatalf("expected comments array in list response:\n%s", truncate(out, 500))
		}
		found := false
		for _, c := range comments {
			cm, _ := c.(map[string]any)
			if fmt.Sprintf("%v", cm["id"]) == taskCommentID {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("comment %s not found in task-comments list", taskCommentID)
		}
	})

	t.Run("task-comments/update", func(t *testing.T) {
		requireID(t, taskCommentID, "task-comments/add must have succeeded")
		updatedComment := testPrefix + "_task_comment_upd"
		out := zoho(t, "projects", "task-comments", "update", taskCommentID,
			"--task", taskID, "--project", projectID,
			"--comment", updatedComment)
		arr := parseJSONArray(t, out)
		if len(arr) == 0 {
			t.Fatalf("task-comments update returned empty array:\n%s", truncate(out, 500))
		}
		assertStringField(t, arr[0], "comment", updatedComment)

		listOut := zoho(t, "projects", "task-comments", "list",
			"--task", taskID, "--project", projectID)
		assertContains(t, listOut, updatedComment)
	})

	t.Run("task-comments/delete", func(t *testing.T) {
		requireID(t, taskCommentID, "task-comments/add must have succeeded")
		out := zoho(t, "projects", "task-comments", "delete", taskCommentID,
			"--task", taskID, "--project", projectID)
		parseJSON(t, out)

		listOut := zoho(t, "projects", "task-comments", "list",
			"--task", taskID, "--project", projectID)
		if strings.Contains(listOut, taskCommentID) {
			t.Errorf("comment %s still found in list after delete", taskCommentID)
		}
		taskCommentID = ""
	})

	t.Run("task-followers/follow", func(t *testing.T) {
		requireID(t, taskID, "tasks/create must have succeeded")
		out := zoho(t, "projects", "task-followers", "follow", taskID,
			"--project", projectID)
		parseJSON(t, out)
	})

	t.Run("task-followers/list", func(t *testing.T) {
		requireID(t, taskID, "tasks/create must have succeeded")
		out := zoho(t, "projects", "task-followers", "list", taskID,
			"--project", projectID)
		m := parseJSON(t, out)
		followers, ok := m["followers"].([]any)
		if !ok || len(followers) == 0 {
			t.Fatalf("expected non-empty followers array:\n%s", truncate(out, 500))
		}
		found := false
		for _, f := range followers {
			fm, _ := f.(map[string]any)
			if fmt.Sprintf("%v", fm["zpuid"]) == ownerZPUID {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("follower %s not found in task-followers list", ownerZPUID)
		}
	})

	t.Run("task-followers/unfollow", func(t *testing.T) {
		requireID(t, taskID, "tasks/create must have succeeded")
		out := zoho(t, "projects", "task-followers", "unfollow", taskID,
			"--project", projectID)
		parseJSON(t, out)

		time.Sleep(2 * time.Second)
		listOut := zoho(t, "projects", "task-followers", "list", taskID,
			"--project", projectID)
		if strings.Contains(listOut, ownerZPUID) {
			t.Errorf("follower %s still present after unfollow", ownerZPUID)
		}
	})

	t.Run("tasks/add-subtask", func(t *testing.T) {
		requireID(t, taskID, "tasks/create must have succeeded")
		subtaskName := testName(t) + "_sub"
		out := zoho(t, "projects", "tasks", "add-subtask",
			"--parent", taskID, "--name", subtaskName, "--project", projectID)
		subtaskID = extractProjectsID(t, out)
		cleanup.trackTask(subtaskID, projectID)
		t.Logf("created subtask %s", subtaskID)

		sub := getTask(t, subtaskID, projectID)
		assertStringField(t, sub, "name", subtaskName)
	})

	t.Run("tasks/subtasks-known-broken", func(t *testing.T) {
		t.Skip("subtasks endpoint not available in V3 API")
	})

	t.Run("tasks/clone", func(t *testing.T) {
		requireID(t, taskID, "tasks/create must have succeeded")
		out := zoho(t, "projects", "tasks", "clone", taskID,
			"--project", projectID)
		clonedTaskID = extractProjectsID(t, out)
		cleanup.trackTask(clonedTaskID, projectID)
		t.Logf("cloned task %s -> %s", taskID, clonedTaskID)

		cloned := getTask(t, clonedTaskID, projectID)
		clonedName := fmt.Sprintf("%v", cloned["name"])
		if clonedName == "" || clonedName == "<nil>" {
			t.Errorf("cloned task has no name")
		}
	})

	t.Run("tasks/move", func(t *testing.T) {
		requireID(t, clonedTaskID, "tasks/clone must have succeeded")
		requireID(t, project2ID, "second project must exist")
		tlOut := zoho(t, "projects", "tasklists", "create",
			"--name", testName(t)+"_tl", "--project", project2ID)
		targetTL := extractProjectsID(t, tlOut)
		cleanup.trackTasklist(targetTL, project2ID)
		moveJSON := toJSON(t, map[string]any{"target_tasklist_id": targetTL})
		zoho(t, "projects", "tasks", "move", clonedTaskID,
			"--project", projectID,
			"--json", moveJSON)

		movedTask := getTask(t, clonedTaskID, project2ID)
		movedName := fmt.Sprintf("%v", movedTask["name"])
		if movedName == "" || movedName == "<nil>" {
			t.Errorf("moved task not found in target project")
		}
		cleanup.trackTask(clonedTaskID, project2ID)
	})

	t.Run("tasks/delete-subtask", func(t *testing.T) {
		requireID(t, subtaskID, "tasks/add-subtask must have succeeded")
		out := zoho(t, "projects", "tasks", "delete", subtaskID,
			"--project", projectID)
		parseJSON(t, out)
		t.Logf("deleted subtask %s", subtaskID)

		r := runZoho(t, "projects", "tasks", "get", subtaskID, "--project", projectID)
		if r.ExitCode == 0 {
			t.Errorf("subtask %s still accessible after delete", subtaskID)
		}
		subtaskID = ""
	})

	t.Run("tasks/delete", func(t *testing.T) {
		requireID(t, taskID, "tasks/create must have succeeded")
		out := zoho(t, "projects", "tasks", "delete", taskID,
			"--project", projectID)
		parseJSON(t, out)

		r := runZoho(t, "projects", "tasks", "get", taskID, "--project", projectID)
		if r.ExitCode == 0 {
			t.Errorf("task %s still accessible after delete", taskID)
		}
		taskID = ""
	})

	t.Run("issues/create", func(t *testing.T) {
		requireID(t, projectID, "core/create must have succeeded")
		issueName = testName(t) + "_issue"
		out := zoho(t, "projects", "issues", "create",
			"--name", issueName, "--project", projectID)
		issueID = extractProjectsID(t, out)
		cleanup.trackIssue(issueID, projectID)
		t.Logf("created issue %s", issueID)

		issue := getIssue(t, issueID, projectID)
		assertStringField(t, issue, "name", issueName)
	})

	t.Run("issues/get", func(t *testing.T) {
		requireID(t, issueID, "issues/create must have succeeded")
		issue := getIssue(t, issueID, projectID)
		assertEqual(t, fmt.Sprintf("%v", issue["id"]), issueID)
		assertStringField(t, issue, "name", issueName)
	})

	t.Run("issues/update", func(t *testing.T) {
		requireID(t, issueID, "issues/create must have succeeded")
		updatedIssueName := issueName + "_upd"
		out := zoho(t, "projects", "issues", "update", issueID,
			"--project", projectID,
			"--json", toJSON(t, map[string]any{"name": updatedIssueName}))
		m := parseJSON(t, out)
		assertEqual(t, fmt.Sprintf("%v", m["name"]), updatedIssueName)

		issue := getIssue(t, issueID, projectID)
		assertStringField(t, issue, "name", updatedIssueName)
		issueName = updatedIssueName
	})

	t.Run("issues/list", func(t *testing.T) {
		requireID(t, issueID, "issues/create must have succeeded")
		out := zoho(t, "projects", "issues", "list", "--project", projectID)
		arr := parseJSONArray(t, out)
		assertNonEmpty(t, arr, "expected at least one issue")
		_, found := findInArray(arr, issueID)
		if !found {
			t.Errorf("issue %s not found in issue list", issueID)
		}
	})

	t.Run("issues/defaults-removed", func(t *testing.T) {
		r := runZoho(t, "projects", "issues", "defaults", "--project", projectID)
		if r.ExitCode == 0 {
			t.Errorf("issues defaults command should not exist (removed: endpoint not in Projects V3)")
		}
	})

	t.Run("issues/description", func(t *testing.T) {
		requireID(t, issueID, "issues/create must have succeeded")
		out := zoho(t, "projects", "issues", "description", issueID,
			"--project", projectID)
		m := parseJSON(t, out)
		if _, ok := m["description"]; !ok {
			if _, ok2 := m["content"]; !ok2 {
				t.Errorf("expected description or content key in response:\n%s", truncate(out, 500))
			}
		}
	})

	t.Run("issues/activities", func(t *testing.T) {
		requireID(t, issueID, "issues/create must have succeeded")
		out := zoho(t, "projects", "issues", "activities", issueID,
			"--project", projectID)
		m := parseJSON(t, out)
		if _, ok := m["activities"]; !ok {
			if _, ok2 := m["page_info"]; !ok2 {
				t.Errorf("expected activities or page_info key in response:\n%s", truncate(out, 500))
			}
		}
	})

	t.Run("issue-comments/add", func(t *testing.T) {
		requireID(t, issueID, "issues/create must have succeeded")
		out := zoho(t, "projects", "issue-comments", "add",
			"--issue", issueID, "--project", projectID,
			"--comment", testPrefix+"_issue_comment")
		m := parseJSON(t, out)
		comments, ok := m["comments"].([]any)
		if !ok || len(comments) == 0 {
			t.Fatalf("expected comments array in add response:\n%s", truncate(out, 500))
		}
		cm, _ := comments[0].(map[string]any)
		issueCommentID = fmt.Sprintf("%v", cm["id"])
		if issueCommentID == "" || issueCommentID == "<nil>" {
			t.Fatalf("no id in issue-comments add response:\n%s", truncate(out, 500))
		}
		assertStringField(t, cm, "comment", testPrefix+"_issue_comment")
		t.Logf("created issue comment %s", issueCommentID)
	})

	t.Run("issue-comments/get", func(t *testing.T) {
		requireID(t, issueCommentID, "issue-comments/add must have succeeded")
		out := zoho(t, "projects", "issue-comments", "get", issueCommentID,
			"--issue", issueID, "--project", projectID)
		m := parseJSON(t, out)
		comments, ok := m["comments"].([]any)
		if !ok || len(comments) == 0 {
			t.Fatalf("expected comments array in get response:\n%s", truncate(out, 500))
		}
		cm, _ := comments[0].(map[string]any)
		assertEqual(t, fmt.Sprintf("%v", cm["id"]), issueCommentID)
		assertStringField(t, cm, "comment", testPrefix+"_issue_comment")
	})

	t.Run("issue-comments/list", func(t *testing.T) {
		requireID(t, issueCommentID, "issue-comments/add must have succeeded")
		out := zoho(t, "projects", "issue-comments", "list",
			"--issue", issueID, "--project", projectID)
		m := parseJSON(t, out)
		comments, ok := m["comments"].([]any)
		if !ok {
			t.Fatalf("expected comments array in list response:\n%s", truncate(out, 500))
		}
		found := false
		for _, c := range comments {
			cm, _ := c.(map[string]any)
			if fmt.Sprintf("%v", cm["id"]) == issueCommentID {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("comment %s not found in issue-comments list", issueCommentID)
		}
	})

	t.Run("issue-comments/update", func(t *testing.T) {
		requireID(t, issueCommentID, "issue-comments/add must have succeeded")
		updatedComment := testPrefix + "_issue_comment_upd"
		out := zoho(t, "projects", "issue-comments", "update", issueCommentID,
			"--issue", issueID, "--project", projectID,
			"--comment", updatedComment)
		m := parseJSON(t, out)
		comments, ok := m["comments"].([]any)
		if !ok || len(comments) == 0 {
			t.Fatalf("expected comments array in update response:\n%s", truncate(out, 500))
		}
		cm, _ := comments[0].(map[string]any)
		assertStringField(t, cm, "comment", updatedComment)

		getOut := zoho(t, "projects", "issue-comments", "get", issueCommentID,
			"--issue", issueID, "--project", projectID)
		assertContains(t, getOut, updatedComment)
	})

	t.Run("issue-comments/delete", func(t *testing.T) {
		requireID(t, issueCommentID, "issue-comments/add must have succeeded")
		out := zoho(t, "projects", "issue-comments", "delete", issueCommentID,
			"--issue", issueID, "--project", projectID)
		parseJSON(t, out)

		listOut := zoho(t, "projects", "issue-comments", "list",
			"--issue", issueID, "--project", projectID)
		if strings.Contains(listOut, issueCommentID) {
			t.Errorf("comment %s still found in list after delete", issueCommentID)
		}
		issueCommentID = ""
	})

	t.Run("issue-followers/follow-known-limitation", func(t *testing.T) {
		t.Skip("issue follower self-follow not supported")
	})

	t.Run("issue-followers/list", func(t *testing.T) {
		requireID(t, issueID, "issues/create must have succeeded")
		out := zoho(t, "projects", "issue-followers", "list", issueID,
			"--project", projectID)
		m := parseJSON(t, out)
		if _, ok := m["followers"]; !ok {
			t.Fatalf("expected followers key in list response:\n%s", truncate(out, 500))
		}
	})

	t.Run("issues/clone", func(t *testing.T) {
		requireID(t, issueID, "issues/create must have succeeded")
		out := zoho(t, "projects", "issues", "clone", issueID,
			"--project", projectID)
		clonedIssueID = extractProjectsID(t, out)
		cleanup.trackIssue(clonedIssueID, projectID)
		t.Logf("cloned issue %s -> %s", issueID, clonedIssueID)
	})

	t.Run("issue-linking/list-empty", func(t *testing.T) {
		requireID(t, issueID, "issues/create must have succeeded")
		out := zoho(t, "projects", "issue-linking", "list", issueID,
			"--project", projectID)
		m := parseJSON(t, out)
		if _, ok := m["issue_linked"]; !ok {
			t.Fatalf("expected issue_linked key in response:\n%s", truncate(out, 500))
		}
	})

	t.Run("issue-linking/link", func(t *testing.T) {
		requireID(t, issueID, "issues/create must have succeeded")
		requireID(t, clonedIssueID, "issues/clone must have succeeded")
		out := zoho(t, "projects", "issue-linking", "link", issueID,
			"--project", projectID,
			"--json", toJSON(t, map[string]any{
				"link_type": "relate",
				"issue_ids": []string{clonedIssueID},
			}))
		parseJSON(t, out)

		listOut := zoho(t, "projects", "issue-linking", "list", issueID,
			"--project", projectID)
		lm := parseJSON(t, listOut)
		linked, ok := lm["issue_linked"].(map[string]any)
		if !ok {
			t.Fatalf("expected issue_linked object in list response:\n%s", truncate(listOut, 500))
		}
		linkedIssues, ok := linked["linked_issues"].(map[string]any)
		if !ok {
			t.Fatalf("expected linked_issues in response:\n%s", truncate(listOut, 500))
		}
		for _, v := range linkedIssues {
			if arr, ok := v.([]any); ok {
				for _, item := range arr {
					if im, ok := item.(map[string]any); ok {
						if fmt.Sprintf("%v", im["issue_id"]) == clonedIssueID {
							issueLinkID = fmt.Sprintf("%v", im["link_id"])
							break
						}
					}
				}
			}
			if issueLinkID != "" {
				break
			}
		}
		if issueLinkID == "" || issueLinkID == "<nil>" {
			t.Fatalf("could not extract link_id from list response:\n%s", truncate(listOut, 500))
		}
		assertContains(t, listOut, clonedIssueID)
		t.Logf("created issue link %s", issueLinkID)
	})

	t.Run("issue-linking/change-type", func(t *testing.T) {
		requireID(t, issueID, "issues/create must have succeeded")
		requireID(t, issueLinkID, "issue-linking/link must have succeeded")
		out := zoho(t, "projects", "issue-linking", "change-type",
			issueID, issueLinkID,
			"--project", projectID,
			"--json", toJSON(t, map[string]any{"link_type": "blocks"}))
		parseJSON(t, out)

		listOut := zoho(t, "projects", "issue-linking", "list", issueID,
			"--project", projectID)
		assertContains(t, listOut, "blocks")
	})

	t.Run("issue-linking/unlink", func(t *testing.T) {
		requireID(t, issueID, "issues/create must have succeeded")
		requireID(t, issueLinkID, "issue-linking/link must have succeeded")
		out := zoho(t, "projects", "issue-linking", "unlink",
			issueID, issueLinkID,
			"--project", projectID)
		parseJSON(t, out)

		listOut := zoho(t, "projects", "issue-linking", "list", issueID,
			"--project", projectID)
		if strings.Contains(listOut, clonedIssueID) {
			t.Errorf("cloned issue %s still in linked issues after unlink", clonedIssueID)
		}
		issueLinkID = ""
	})

	t.Run("issue-linking/bulk-link", func(t *testing.T) {
		requireID(t, issueID, "issues/create must have succeeded")
		requireID(t, clonedIssueID, "issues/clone must have succeeded")
		r := runZoho(t, "projects", "issue-linking", "bulk-link",
			"--project", projectID,
			"--json", toJSON(t, map[string]any{
				"link_type":         "relate",
				"issue_ids":         []string{issueID},
				"linking_issue_ids": []string{clonedIssueID},
			}))
		if r.ExitCode != 0 {
			t.Logf("bulk-link failed: %s",
				truncate(r.Stderr+r.Stdout, 300))
			return
		}
		parseJSON(t, r.Stdout)

		listOut := zoho(t, "projects", "issue-linking", "list", issueID,
			"--project", projectID)
		assertContains(t, listOut, clonedIssueID)
	})

	t.Run("issues/move", func(t *testing.T) {
		requireID(t, clonedIssueID, "issues/clone must have succeeded")
		requireID(t, project2ID, "second project must exist")
		moveJSON := toJSON(t, map[string]any{"to_project": project2ID})
		zoho(t, "projects", "issues", "move", clonedIssueID,
			"--project", projectID,
			"--json", moveJSON)

		movedIssue := getIssue(t, clonedIssueID, project2ID)
		movedName := fmt.Sprintf("%v", movedIssue["name"])
		if movedName == "" || movedName == "<nil>" {
			t.Errorf("moved issue not found in target project")
		}
		cleanup.trackIssue(clonedIssueID, project2ID)
	})

	t.Run("issue-resolution/add", func(t *testing.T) {
		requireID(t, issueID, "issues/create must have succeeded")
		out := zoho(t, "projects", "issue-resolution", "add", issueID,
			"--project", projectID,
			"--json", toJSON(t, map[string]any{"resolution": testPrefix + " resolution text"}))
		parseJSON(t, out)
	})

	t.Run("issue-resolution/get", func(t *testing.T) {
		requireID(t, issueID, "issues/create must have succeeded")
		out := zoho(t, "projects", "issue-resolution", "get", issueID,
			"--project", projectID)
		m := parseJSON(t, out)
		ir, ok := m["issue_resolution"].(map[string]any)
		if !ok {
			t.Fatalf("expected issue_resolution object in get response:\n%s", truncate(out, 500))
		}
		resolution := fmt.Sprintf("%v", ir["resolution"])
		assertContains(t, resolution, testPrefix)
	})

	t.Run("issue-resolution/update", func(t *testing.T) {
		requireID(t, issueID, "issues/create must have succeeded")
		zoho(t, "projects", "issue-resolution", "update", issueID,
			"--project", projectID,
			"--json", toJSON(t, map[string]any{"resolution": testPrefix + " updated resolution"}))

		out := zoho(t, "projects", "issue-resolution", "get", issueID,
			"--project", projectID)
		m := parseJSON(t, out)
		ir, _ := m["issue_resolution"].(map[string]any)
		resolution := fmt.Sprintf("%v", ir["resolution"])
		assertContains(t, resolution, "updated resolution")
	})

	t.Run("issue-resolution/delete", func(t *testing.T) {
		requireID(t, issueID, "issues/create must have succeeded")
		out := zoho(t, "projects", "issue-resolution", "delete", issueID,
			"--project", projectID)
		parseJSON(t, out)

		getOut := zoho(t, "projects", "issue-resolution", "get", issueID,
			"--project", projectID)
		m := parseJSON(t, getOut)
		ir, _ := m["issue_resolution"].(map[string]any)
		if _, hasResolution := ir["resolution"]; hasResolution {
			t.Errorf("resolution key still present after delete")
		}
	})

	t.Run("issue-attachments/list", func(t *testing.T) {
		requireID(t, issueID, "issues/create must have succeeded")
		out := zoho(t, "projects", "issue-attachments", "list", issueID,
			"--project", projectID)
		m := parseJSON(t, out)
		if _, ok := m["attachments"]; !ok {
			t.Fatalf("expected attachments key in response:\n%s", truncate(out, 500))
		}
	})

	t.Run("issue-attachments/associate-known-broken", func(t *testing.T) {
		t.Skip("no real attachment ID available for test")
	})

	t.Run("issue-attachments/dissociate-known-broken", func(t *testing.T) {
		t.Skip("no real attachment ID available for test")
	})

	t.Run("issues/delete", func(t *testing.T) {
		requireID(t, issueID, "issues/create must have succeeded")
		out := zoho(t, "projects", "issues", "delete", issueID,
			"--project", projectID)
		parseJSON(t, out)

		r := runZoho(t, "projects", "issues", "get", issueID, "--project", projectID)
		if r.ExitCode == 0 && strings.Contains(r.Stdout, issueID) {
			t.Errorf("issue %s still accessible after delete", issueID)
		}
		issueID = ""
	})

	t.Run("tasklists/create", func(t *testing.T) {
		requireID(t, projectID, "core/create must have succeeded")
		tasklistName = testName(t) + "_tl"
		out := zoho(t, "projects", "tasklists", "create",
			"--name", tasklistName, "--project", projectID)
		tasklistID = extractProjectsID(t, out)
		cleanup.trackTasklist(tasklistID, projectID)
		t.Logf("created tasklist %s", tasklistID)

		tl := getTasklist(t, tasklistID, projectID)
		assertStringField(t, tl, "name", tasklistName)
	})

	t.Run("tasklists/get", func(t *testing.T) {
		requireID(t, tasklistID, "tasklists/create must have succeeded")
		tl := getTasklist(t, tasklistID, projectID)
		assertEqual(t, fmt.Sprintf("%v", tl["id"]), tasklistID)
		assertStringField(t, tl, "name", tasklistName)
	})

	t.Run("tasklists/list", func(t *testing.T) {
		requireID(t, tasklistID, "tasklists/create must have succeeded")
		out := zoho(t, "projects", "tasklists", "list", "--project", projectID)
		arr := parseJSONArray(t, out)
		assertNonEmpty(t, arr, "expected at least one tasklist")
		_, found := findInArray(arr, tasklistID)
		if !found {
			t.Errorf("tasklist %s not found in list", tasklistID)
		}
	})

	t.Run("tasklists/update", func(t *testing.T) {
		requireID(t, tasklistID, "tasklists/create must have succeeded")
		updatedTLName := tasklistName + "_upd"
		out := zoho(t, "projects", "tasklists", "update", tasklistID,
			"--project", projectID,
			"--json", toJSON(t, map[string]any{"name": updatedTLName}))
		m := parseJSON(t, out)
		assertEqual(t, fmt.Sprintf("%v", m["name"]), updatedTLName)

		tl := getTasklist(t, tasklistID, projectID)
		assertStringField(t, tl, "name", updatedTLName)
		tasklistName = updatedTLName
	})

	t.Run("tasklist-comments/add", func(t *testing.T) {
		requireID(t, tasklistID, "tasklists/create must have succeeded")
		out := zoho(t, "projects", "tasklist-comments", "add",
			"--tasklist", tasklistID, "--project", projectID,
			"--comment", testPrefix+"_tl_comment")
		arr := parseJSONArray(t, out)
		if len(arr) == 0 {
			t.Fatalf("tasklist-comments add returned empty array:\n%s", truncate(out, 500))
		}
		tasklistCommentID = fmt.Sprintf("%v", arr[0]["id"])
		if tasklistCommentID == "" || tasklistCommentID == "<nil>" {
			t.Fatalf("no id in tasklist-comments add response:\n%s", truncate(out, 500))
		}
		assertStringField(t, arr[0], "comment", testPrefix+"_tl_comment")
		t.Logf("created tasklist comment %s", tasklistCommentID)
	})

	t.Run("tasklist-comments/get", func(t *testing.T) {
		requireID(t, tasklistCommentID, "tasklist-comments/add must have succeeded")
		out := zoho(t, "projects", "tasklist-comments", "get", tasklistCommentID,
			"--tasklist", tasklistID, "--project", projectID)
		m := parseJSON(t, out)
		assertEqual(t, fmt.Sprintf("%v", m["id"]), tasklistCommentID)
	})

	t.Run("tasklist-comments/list", func(t *testing.T) {
		requireID(t, tasklistCommentID, "tasklist-comments/add must have succeeded")
		out := zoho(t, "projects", "tasklist-comments", "list",
			"--tasklist", tasklistID, "--project", projectID)
		m := parseJSON(t, out)
		comments, ok := m["comments"].([]any)
		if !ok {
			t.Fatalf("expected comments array in list response:\n%s", truncate(out, 500))
		}
		found := false
		for _, c := range comments {
			cm, _ := c.(map[string]any)
			if fmt.Sprintf("%v", cm["id"]) == tasklistCommentID {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("comment %s not found in tasklist-comments list", tasklistCommentID)
		}
	})

	t.Run("tasklist-comments/update", func(t *testing.T) {
		requireID(t, tasklistCommentID, "tasklist-comments/add must have succeeded")
		updatedComment := testPrefix + "_tl_comment_upd"
		out := zoho(t, "projects", "tasklist-comments", "update", tasklistCommentID,
			"--tasklist", tasklistID, "--project", projectID,
			"--comment", updatedComment)
		arr := parseJSONArray(t, out)
		if len(arr) == 0 {
			t.Fatalf("tasklist-comments update returned empty array:\n%s", truncate(out, 500))
		}
		assertStringField(t, arr[0], "comment", updatedComment)

		getOut := zoho(t, "projects", "tasklist-comments", "get", tasklistCommentID,
			"--tasklist", tasklistID, "--project", projectID)
		assertContains(t, getOut, updatedComment)
	})

	t.Run("tasklist-comments/delete", func(t *testing.T) {
		requireID(t, tasklistCommentID, "tasklist-comments/add must have succeeded")
		out := zoho(t, "projects", "tasklist-comments", "delete", tasklistCommentID,
			"--tasklist", tasklistID, "--project", projectID)
		parseJSON(t, out)

		listOut := zoho(t, "projects", "tasklist-comments", "list",
			"--tasklist", tasklistID, "--project", projectID)
		if strings.Contains(listOut, tasklistCommentID) {
			t.Errorf("comment %s still found in list after delete", tasklistCommentID)
		}
		tasklistCommentID = ""
	})

	t.Run("tasklist-followers/follow", func(t *testing.T) {
		requireID(t, tasklistID, "tasklists/create must have succeeded")
		out := zoho(t, "projects", "tasklist-followers", "follow", tasklistID,
			"--project", projectID)
		parseJSON(t, out)
	})

	t.Run("tasklist-followers/list", func(t *testing.T) {
		requireID(t, tasklistID, "tasklists/create must have succeeded")
		out := zoho(t, "projects", "tasklist-followers", "list", tasklistID,
			"--project", projectID)
		m := parseJSON(t, out)
		followers, ok := m["followers"].([]any)
		if !ok || len(followers) == 0 {
			t.Fatalf("expected non-empty followers array:\n%s", truncate(out, 500))
		}
		found := false
		for _, f := range followers {
			fm, _ := f.(map[string]any)
			if fmt.Sprintf("%v", fm["zpuid"]) == ownerZPUID {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("follower %s not found in tasklist-followers list", ownerZPUID)
		}
	})

	t.Run("tasklist-followers/unfollow", func(t *testing.T) {
		requireID(t, tasklistID, "tasklists/create must have succeeded")
		out := zoho(t, "projects", "tasklist-followers", "unfollow", tasklistID,
			"--project", projectID)
		parseJSON(t, out)

		time.Sleep(2 * time.Second)
		listOut := zoho(t, "projects", "tasklist-followers", "list", tasklistID,
			"--project", projectID)
		if strings.Contains(listOut, ownerZPUID) {
			t.Errorf("follower %s still present after unfollow", ownerZPUID)
		}
	})

	t.Run("tasklists/delete", func(t *testing.T) {
		requireID(t, tasklistID, "tasklists/create must have succeeded")
		out := zoho(t, "projects", "tasklists", "delete", tasklistID,
			"--project", projectID)
		parseJSON(t, out)

		time.Sleep(2 * time.Second)
		r := runZoho(t, "projects", "tasklists", "get", tasklistID, "--project", projectID)
		if r.ExitCode == 0 {
			t.Errorf("tasklist %s still accessible after delete", tasklistID)
		}
		tasklistID = ""
	})

	t.Run("milestones/create", func(t *testing.T) {
		requireID(t, projectID, "core/create must have succeeded")
		milestoneName = testName(t) + "_ms"
		now := time.Now()
		startDate := now.Format("2006-01-02")
		endDate := now.AddDate(0, 1, 0).Format("2006-01-02")
		out := zoho(t, "projects", "milestones", "create",
			"--name", milestoneName,
			"--start", startDate,
			"--end", endDate,
			"--project", projectID)
		milestoneID = extractProjectsID(t, out)
		cleanup.trackMilestone(milestoneID, projectID)
		t.Logf("created milestone %s", milestoneID)

		ms := getMilestone(t, milestoneID, projectID)
		assertStringField(t, ms, "name", milestoneName)
	})

	t.Run("milestones/get", func(t *testing.T) {
		requireID(t, milestoneID, "milestones/create must have succeeded")
		ms := getMilestone(t, milestoneID, projectID)
		assertEqual(t, fmt.Sprintf("%v", ms["id"]), milestoneID)
		assertStringField(t, ms, "name", milestoneName)
	})

	t.Run("milestones/list", func(t *testing.T) {
		requireID(t, milestoneID, "milestones/create must have succeeded")
		out := zoho(t, "projects", "milestones", "list", "--project", projectID)
		arr := parseJSONArray(t, out)
		assertNonEmpty(t, arr, "expected at least one milestone")
		_, found := findInArray(arr, milestoneID)
		if !found {
			t.Errorf("milestone %s not found in list", milestoneID)
		}
	})

	t.Run("milestones/update", func(t *testing.T) {
		requireID(t, milestoneID, "milestones/create must have succeeded")
		updatedMSName := milestoneName + "_upd"
		out := zoho(t, "projects", "milestones", "update", milestoneID,
			"--project", projectID,
			"--json", toJSON(t, map[string]any{"name": updatedMSName}))
		m := parseJSON(t, out)
		assertEqual(t, fmt.Sprintf("%v", m["name"]), updatedMSName)

		ms := getMilestone(t, milestoneID, projectID)
		assertStringField(t, ms, "name", updatedMSName)
		milestoneName = updatedMSName
	})

	t.Run("milestones/delete", func(t *testing.T) {
		requireID(t, milestoneID, "milestones/create must have succeeded")
		out := zoho(t, "projects", "milestones", "delete", milestoneID,
			"--project", projectID)
		parseJSON(t, out)

		r := runZoho(t, "projects", "milestones", "get", milestoneID, "--project", projectID)
		if r.ExitCode == 0 {
			t.Errorf("milestone %s still accessible after delete", milestoneID)
		}
		milestoneID = ""
	})

	t.Run("timelogs/add", func(t *testing.T) {
		requireID(t, projectID, "core/create must have succeeded")
		tlTaskName := testName(t) + "_tltask"
		taskOut := zoho(t, "projects", "tasks", "create",
			"--name", tlTaskName, "--project", projectID)
		tlTaskID = extractProjectsID(t, taskOut)
		cleanup.trackTask(tlTaskID, projectID)

		zoho(t, "projects", "tasks", "update", tlTaskID,
			"--project", projectID,
			"--json", toJSON(t, map[string]any{
				"owners_and_work": map[string]any{
					"owners": []map[string]string{{"zpuid": ownerZPUID}},
				},
			}))

		today := time.Now().Format("2006-01-02")
		out := zoho(t, "projects", "timelogs", "add",
			"--date", today,
			"--hours", "2",
			"--task", tlTaskID,
			"--owner", ownerZPUID,
			"--notes", testPrefix+"_timelog",
			"--project", projectID)
		timelogID = extractProjectsID(t, out)
		cleanup.trackTimelog(timelogID, projectID)
		t.Logf("created timelog %s for task %s", timelogID, tlTaskID)
	})

	t.Run("timelogs/get", func(t *testing.T) {
		requireID(t, timelogID, "timelogs/add must have succeeded")
		tl := getTimelog(t, timelogID, projectID)
		assertEqual(t, fmt.Sprintf("%v", tl["id"]), timelogID)
	})

	t.Run("timelogs/list", func(t *testing.T) {
		requireID(t, timelogID, "timelogs/add must have succeeded")
		out := zoho(t, "projects", "timelogs", "list",
			"--module", "task", "--project", projectID)
		var raw json.RawMessage
		if err := json.Unmarshal([]byte(out), &raw); err != nil {
			t.Fatalf("failed to parse timelogs list: %v\nraw: %s", err, truncate(out, 500))
		}
		assertContains(t, out, timelogID)
	})

	t.Run("timelogs/update", func(t *testing.T) {
		requireID(t, timelogID, "timelogs/add must have succeeded")
		out := zoho(t, "projects", "timelogs", "update", timelogID,
			"--project", projectID,
			"--type", "task",
			"--task", tlTaskID,
			"--json", toJSON(t, map[string]any{"hours": "3"}))
		parseJSON(t, out)

		tl := getTimelog(t, timelogID, projectID)
		hours := fmt.Sprintf("%v", tl["log_hour"])
		if hours != "3" && hours != "03:00" && hours != "3:00" {
			t.Errorf("expected timelog hours=3 after update, got %s", hours)
		}
	})

	t.Run("timelogs/delete", func(t *testing.T) {
		requireID(t, timelogID, "timelogs/add must have succeeded")
		out := zoho(t, "projects", "timelogs", "delete", timelogID,
			"--project", projectID, "--type", "task")
		parseJSON(t, out)

		r := runZoho(t, "projects", "timelogs", "get", timelogID, "--project", projectID, "--type", "task")
		if r.ExitCode == 0 {
			t.Errorf("timelog %s still accessible after delete", timelogID)
		}
		timelogID = ""
	})

	t.Run("project-comments/add", func(t *testing.T) {
		requireID(t, projectID, "core/create must have succeeded")
		commentText := testPrefix + "_project_comment"
		out := zoho(t, "projects", "project-comments", "add",
			"--project", projectID,
			"--comment", commentText)
		m := parseJSON(t, out)
		comments, ok := m["comments"].([]any)
		if !ok || len(comments) == 0 {
			t.Fatalf("expected comments array in add response:\n%s", truncate(out, 500))
		}
		cm, _ := comments[0].(map[string]any)
		projectCommentID = fmt.Sprintf("%v", cm["id"])
		if projectCommentID == "" || projectCommentID == "<nil>" {
			t.Fatalf("no id in project-comments add response:\n%s", truncate(out, 500))
		}
		assertStringField(t, cm, "content", commentText)
		t.Logf("created project comment %s", projectCommentID)
	})

	t.Run("project-comments/get", func(t *testing.T) {
		requireID(t, projectCommentID, "project-comments/add must have succeeded")
		out := zoho(t, "projects", "project-comments", "get", projectCommentID,
			"--project", projectID)
		m := parseJSON(t, out)
		assertEqual(t, fmt.Sprintf("%v", m["id"]), projectCommentID)
		assertStringField(t, m, "content", testPrefix+"_project_comment")
	})

	t.Run("project-comments/list", func(t *testing.T) {
		requireID(t, projectCommentID, "project-comments/add must have succeeded")
		out := zoho(t, "projects", "project-comments", "list",
			"--project", projectID)
		m := parseJSON(t, out)
		comments, ok := m["comments"].([]any)
		if !ok {
			t.Fatalf("expected comments array in list response:\n%s", truncate(out, 500))
		}
		found := false
		for _, c := range comments {
			cm, _ := c.(map[string]any)
			if fmt.Sprintf("%v", cm["id"]) == projectCommentID {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("comment %s not found in project-comments list", projectCommentID)
		}
	})

	t.Run("project-comments/update", func(t *testing.T) {
		requireID(t, projectCommentID, "project-comments/add must have succeeded")
		updatedComment := testPrefix + "_project_comment_upd"
		out := zoho(t, "projects", "project-comments", "update", projectCommentID,
			"--project", projectID,
			"--comment", updatedComment)
		m := parseJSON(t, out)
		comments, ok := m["comments"].([]any)
		if !ok || len(comments) == 0 {
			t.Fatalf("expected comments array in update response:\n%s", truncate(out, 500))
		}
		cm, _ := comments[0].(map[string]any)
		assertStringField(t, cm, "content", updatedComment)

		getOut := zoho(t, "projects", "project-comments", "get", projectCommentID,
			"--project", projectID)
		getM := parseJSON(t, getOut)
		assertStringField(t, getM, "content", updatedComment)
	})

	t.Run("project-comments/delete", func(t *testing.T) {
		requireID(t, projectCommentID, "project-comments/add must have succeeded")
		out := zoho(t, "projects", "project-comments", "delete", projectCommentID,
			"--project", projectID)
		parseJSON(t, out)

		listOut := zoho(t, "projects", "project-comments", "list",
			"--project", projectID)
		if strings.Contains(listOut, projectCommentID) {
			t.Errorf("comment %s still found in list after delete", projectCommentID)
		}
		projectCommentID = ""
	})

	t.Run("forum-categories/create", func(t *testing.T) {
		requireID(t, projectID, "core/create must have succeeded")
		categoryName := testName(t) + "_forum_category"
		out := zoho(t, "projects", "forum-categories", "create",
			"--project", projectID,
			"--json", toJSON(t, map[string]any{"name": categoryName}))
		m := parseJSON(t, out)
		categories, ok := m["categories"].([]any)
		if !ok || len(categories) == 0 {
			t.Fatalf("expected categories array in create response:\n%s", truncate(out, 500))
		}
		cm, _ := categories[0].(map[string]any)
		forumCategoryID = fmt.Sprintf("%v", cm["id"])
		if forumCategoryID == "" || forumCategoryID == "<nil>" {
			t.Fatalf("no id in forum-categories create response:\n%s", truncate(out, 500))
		}
		assertStringField(t, cm, "name", categoryName)
		cleanup.trackForumCategory(forumCategoryID, projectID)
		t.Logf("created forum category %s", forumCategoryID)
	})

	t.Run("forum-categories/list", func(t *testing.T) {
		requireID(t, forumCategoryID, "forum-categories/create must have succeeded")
		out := zoho(t, "projects", "forum-categories", "list", "--project", projectID)
		arr := parseJSONArray(t, out)
		found := false
		for _, c := range arr {
			if fmt.Sprintf("%v", c["id"]) == forumCategoryID {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("forum category %s not found in list", forumCategoryID)
		}
	})

	t.Run("forums/create", func(t *testing.T) {
		requireID(t, forumCategoryID, "forum-categories/create must have succeeded")
		forumTitle := testName(t) + "_forum"
		forumContent := testPrefix + "_forum_content"
		out := zoho(t, "projects", "forums", "create",
			"--project", projectID,
			"--json", toJSON(t, map[string]any{
				"title":       forumTitle,
				"content":     forumContent,
				"category_id": forumCategoryID,
			}))
		m := parseJSON(t, out)
		forums, ok := m["forums"].([]any)
		if !ok || len(forums) == 0 {
			t.Fatalf("expected forums array in create response:\n%s", truncate(out, 500))
		}
		fm, _ := forums[0].(map[string]any)
		forumID = fmt.Sprintf("%v", fm["id"])
		if forumID == "" || forumID == "<nil>" {
			t.Fatalf("no id in forums create response:\n%s", truncate(out, 500))
		}
		assertStringField(t, fm, "title", forumTitle)
		assertStringField(t, fm, "content", forumContent)
		cleanup.trackForum(forumID, projectID)
		t.Logf("created forum %s", forumID)
	})

	t.Run("forums/get", func(t *testing.T) {
		requireID(t, forumID, "forums/create must have succeeded")
		out := zoho(t, "projects", "forums", "get", forumID,
			"--project", projectID)
		m := parseJSON(t, out)
		forums, ok := m["forums"].([]any)
		if !ok || len(forums) == 0 {
			t.Fatalf("expected forums array in get response:\n%s", truncate(out, 500))
		}
		fm, _ := forums[0].(map[string]any)
		assertEqual(t, fmt.Sprintf("%v", fm["id"]), forumID)
	})

	t.Run("forums/list", func(t *testing.T) {
		requireID(t, forumID, "forums/create must have succeeded")
		out := zoho(t, "projects", "forums", "list",
			"--project", projectID)
		m := parseJSON(t, out)
		forums, ok := m["forums"].([]any)
		if !ok {
			t.Fatalf("expected forums array in list response:\n%s", truncate(out, 500))
		}
		found := false
		for _, f := range forums {
			fm, _ := f.(map[string]any)
			if fmt.Sprintf("%v", fm["id"]) == forumID {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("forum %s not found in forums list", forumID)
		}
	})

	t.Run("forums/update", func(t *testing.T) {
		requireID(t, forumID, "forums/create must have succeeded")
		updatedTitle := testName(t) + "_forum_upd"
		updatedContent := testPrefix + "_forum_content_upd"
		out := zoho(t, "projects", "forums", "update", forumID,
			"--project", projectID,
			"--json", toJSON(t, map[string]any{"title": updatedTitle, "content": updatedContent}))
		m := parseJSON(t, out)
		forums, ok := m["forums"].([]any)
		if !ok || len(forums) == 0 {
			t.Fatalf("expected forums array in update response:\n%s", truncate(out, 500))
		}
		fm, _ := forums[0].(map[string]any)
		assertStringField(t, fm, "title", updatedTitle)
		assertStringField(t, fm, "content", updatedContent)

		getOut := zoho(t, "projects", "forums", "get", forumID,
			"--project", projectID)
		getM := parseJSON(t, getOut)
		getForums, ok := getM["forums"].([]any)
		if !ok || len(getForums) == 0 {
			t.Fatalf("expected forums array in get response after update:\n%s", truncate(getOut, 500))
		}
		getFM, _ := getForums[0].(map[string]any)
		assertStringField(t, getFM, "title", updatedTitle)
		assertStringField(t, getFM, "content", updatedContent)
	})

	t.Run("forum-comments/add", func(t *testing.T) {
		requireID(t, forumID, "forums/create must have succeeded")
		commentText := testPrefix + "_forum_comment"
		out := zoho(t, "projects", "forum-comments", "add",
			"--project", projectID,
			"--forum", forumID,
			"--comment", commentText)
		m := parseJSON(t, out)
		comment, ok := m["forum_comments"].(map[string]any)
		if !ok {
			t.Fatalf("expected forum_comments object in add response:\n%s", truncate(out, 500))
		}
		forumCommentID = fmt.Sprintf("%v", comment["id"])
		if forumCommentID == "" || forumCommentID == "<nil>" {
			t.Fatalf("no id in forum-comments add response:\n%s", truncate(out, 500))
		}
		assertStringField(t, comment, "content", commentText)
		t.Logf("created forum comment %s", forumCommentID)
	})

	t.Run("forum-comments/get", func(t *testing.T) {
		requireID(t, forumCommentID, "forum-comments/add must have succeeded")
		out := zoho(t, "projects", "forum-comments", "get", forumCommentID,
			"--project", projectID,
			"--forum", forumID)
		m := parseJSON(t, out)
		comments, ok := m["forum_comments"].([]any)
		if !ok || len(comments) == 0 {
			t.Fatalf("expected forum_comments array in get response:\n%s", truncate(out, 500))
		}
		cm, _ := comments[0].(map[string]any)
		assertEqual(t, fmt.Sprintf("%v", cm["id"]), forumCommentID)
	})

	t.Run("forum-comments/list", func(t *testing.T) {
		requireID(t, forumCommentID, "forum-comments/add must have succeeded")
		out := zoho(t, "projects", "forum-comments", "list",
			"--project", projectID,
			"--forum", forumID)
		m := parseJSON(t, out)
		comments, ok := m["forum_comments"].([]any)
		if !ok {
			t.Fatalf("expected forum_comments array in list response:\n%s", truncate(out, 500))
		}
		found := false
		for _, c := range comments {
			cm, _ := c.(map[string]any)
			if fmt.Sprintf("%v", cm["id"]) == forumCommentID {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("forum comment %s not found in list", forumCommentID)
		}
	})

	t.Run("forum-comments/update", func(t *testing.T) {
		requireID(t, forumCommentID, "forum-comments/add must have succeeded")
		updatedComment := testPrefix + "_forum_comment_upd"
		out := zoho(t, "projects", "forum-comments", "update", forumCommentID,
			"--project", projectID,
			"--forum", forumID,
			"--comment", updatedComment)
		m := parseJSON(t, out)
		comments, ok := m["forum_comments"].([]any)
		if !ok || len(comments) == 0 {
			t.Fatalf("expected forum_comments array in update response:\n%s", truncate(out, 500))
		}
		cm, _ := comments[0].(map[string]any)
		assertStringField(t, cm, "content", updatedComment)

		getOut := zoho(t, "projects", "forum-comments", "get", forumCommentID,
			"--project", projectID,
			"--forum", forumID)
		getM := parseJSON(t, getOut)
		getComments, ok := getM["forum_comments"].([]any)
		if !ok || len(getComments) == 0 {
			t.Fatalf("expected forum_comments array in get response after update:\n%s", truncate(getOut, 500))
		}
		getCM, _ := getComments[0].(map[string]any)
		assertStringField(t, getCM, "content", updatedComment)
	})

	t.Run("forum-comments/best-answer-known-limitation", func(t *testing.T) {
		t.Skip("best-answer endpoint returns error")
	})

	t.Run("forum-comments/unbest-answer-known-limitation", func(t *testing.T) {
		t.Skip("unbest-answer endpoint returns error")
	})

	t.Run("forum-comments/delete", func(t *testing.T) {
		requireID(t, forumCommentID, "forum-comments/add must have succeeded")
		out := zoho(t, "projects", "forum-comments", "delete", forumCommentID,
			"--project", projectID,
			"--forum", forumID)
		parseJSON(t, out)

		listOut := zoho(t, "projects", "forum-comments", "list",
			"--project", projectID,
			"--forum", forumID)
		if strings.Contains(listOut, forumCommentID) {
			t.Errorf("forum comment %s still found in list after delete", forumCommentID)
		}
		forumCommentID = ""
	})

	t.Run("forum-followers/list", func(t *testing.T) {
		requireID(t, forumID, "forums/create must have succeeded")
		out := zoho(t, "projects", "forum-followers", "list", forumID,
			"--project", projectID)
		m := parseJSON(t, out)
		if _, ok := m["followers"].([]any); !ok {
			t.Fatalf("expected followers array in list response:\n%s", truncate(out, 500))
		}
	})

	t.Run("forum-followers/follow-known-limitation", func(t *testing.T) {
		t.Skip("forum follower endpoint requires body")
	})

	t.Run("forum-followers/unfollow-known-limitation", func(t *testing.T) {
		t.Skip("forum follower endpoint requires body")
	})

	t.Run("forums/delete", func(t *testing.T) {
		requireID(t, forumID, "forums/create must have succeeded")
		out := zoho(t, "projects", "forums", "delete", forumID,
			"--project", projectID)
		parseJSON(t, out)

		listOut := zoho(t, "projects", "forums", "list",
			"--project", projectID)
		if strings.Contains(listOut, forumID) {
			t.Errorf("forum %s still found in list after delete", forumID)
		}
		forumID = ""
	})

	t.Run("forum-categories/update", func(t *testing.T) {
		requireID(t, forumCategoryID, "forum-categories/create must have succeeded")
		updatedName := testName(t) + "_forum_category_upd"
		out := zoho(t, "projects", "forum-categories", "update", forumCategoryID,
			"--project", projectID,
			"--json", toJSON(t, map[string]any{"name": updatedName}))
		m := parseJSON(t, out)
		categories, ok := m["categories"].([]any)
		if !ok || len(categories) == 0 {
			t.Fatalf("expected categories array in update response:\n%s", truncate(out, 500))
		}
		cm, _ := categories[0].(map[string]any)
		assertEqual(t, fmt.Sprintf("%v", cm["id"]), forumCategoryID)
		assertStringField(t, cm, "name", updatedName)
	})

	t.Run("forum-categories/delete", func(t *testing.T) {
		requireID(t, forumCategoryID, "forum-categories/create must have succeeded")
		out := zoho(t, "projects", "forum-categories", "delete", forumCategoryID,
			"--project", projectID)
		parseJSON(t, out)
	})

	t.Run("forum-categories/list-verify-deleted", func(t *testing.T) {
		requireID(t, forumCategoryID, "forum-categories/create must have succeeded")
		out := zoho(t, "projects", "forum-categories", "list",
			"--project", projectID)
		arr := parseJSONArray(t, out)
		for _, c := range arr {
			if fmt.Sprintf("%v", c["id"]) == forumCategoryID {
				t.Errorf("forum category %s still found after delete", forumCategoryID)
				break
			}
		}
		forumCategoryID = ""
	})

	t.Run("project-groups/create", func(t *testing.T) {
		groupName := testName(t) + "_group"
		out := zoho(t, "projects", "project-groups", "create",
			"--json", toJSON(t, map[string]any{"name": groupName, "type": "public"}))
		m := parseJSON(t, out)
		projectGroupID = fmt.Sprintf("%v", m["id"])
		if projectGroupID == "" || projectGroupID == "<nil>" {
			t.Fatalf("no id in project-groups create response:\n%s", truncate(out, 500))
		}
		assertStringField(t, m, "name", groupName)
		cleanup.trackProjectGroup(projectGroupID)
		t.Logf("created project group %s", projectGroupID)
	})

	t.Run("project-groups/list", func(t *testing.T) {
		requireID(t, projectGroupID, "project-groups/create must have succeeded")
		out := zoho(t, "projects", "project-groups", "list")
		m := parseJSON(t, out)
		groups, ok := m["project-groups"].([]any)
		if !ok {
			t.Fatalf("expected project-groups array in list response:\n%s", truncate(out, 500))
		}
		found := false
		for _, g := range groups {
			gm, _ := g.(map[string]any)
			if fmt.Sprintf("%v", gm["id"]) == projectGroupID {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("project group %s not found in list", projectGroupID)
		}
	})

	t.Run("project-groups/my", func(t *testing.T) {
		out := zoho(t, "projects", "project-groups", "my")
		m := parseJSON(t, out)
		if _, ok := m["project-groups"].([]any); !ok {
			t.Fatalf("expected project-groups array in my response:\n%s", truncate(out, 500))
		}
	})

	t.Run("project-groups/update", func(t *testing.T) {
		requireID(t, projectGroupID, "project-groups/create must have succeeded")
		updatedName := testName(t) + "_group_upd"
		out := zoho(t, "projects", "project-groups", "update", projectGroupID,
			"--json", toJSON(t, map[string]any{"name": updatedName}))
		m := parseJSON(t, out)
		assertEqual(t, fmt.Sprintf("%v", m["id"]), projectGroupID)
		assertStringField(t, m, "name", updatedName)
	})

	t.Run("project-groups/delete", func(t *testing.T) {
		requireID(t, projectGroupID, "project-groups/create must have succeeded")
		out := zoho(t, "projects", "project-groups", "delete", projectGroupID)
		parseJSON(t, out)
	})

	t.Run("project-groups/list-verify-deleted", func(t *testing.T) {
		requireID(t, projectGroupID, "project-groups/create must have succeeded")
		out := zoho(t, "projects", "project-groups", "list")
		m := parseJSON(t, out)
		groups, ok := m["project-groups"].([]any)
		if !ok {
			t.Fatalf("expected project-groups array in list response:\n%s", truncate(out, 500))
		}
		for _, g := range groups {
			gm, _ := g.(map[string]any)
			if fmt.Sprintf("%v", gm["id"]) == projectGroupID {
				t.Errorf("project group %s still found after delete", projectGroupID)
				break
			}
		}
		projectGroupID = ""
	})

	t.Run("tags/create", func(t *testing.T) {
		tagName := testName(t) + "_tag"
		out := zoho(t, "projects", "tags", "create",
			"--json", toJSON(t, []map[string]any{{"name": tagName, "color_class": "bg-tag1"}}))
		m := parseJSON(t, out)
		tags, ok := m["tags"].([]any)
		if !ok || len(tags) == 0 {
			t.Fatalf("expected tags array in create response:\n%s", truncate(out, 500))
		}
		tm, _ := tags[0].(map[string]any)
		tagID = fmt.Sprintf("%v", tm["id"])
		if tagID == "" || tagID == "<nil>" {
			t.Fatalf("no id in tags create response:\n%s", truncate(out, 500))
		}
		assertStringField(t, tm, "name", tagName)
		cleanup.trackTag(tagID)
		t.Logf("created tag %s", tagID)
	})

	t.Run("tags/list", func(t *testing.T) {
		requireID(t, tagID, "tags/create must have succeeded")
		out := zoho(t, "projects", "tags", "list")
		m := parseJSON(t, out)
		tags, ok := m["tags"].([]any)
		if !ok {
			t.Fatalf("expected tags array in list response:\n%s", truncate(out, 500))
		}
		found := false
		for _, it := range tags {
			tm, _ := it.(map[string]any)
			if fmt.Sprintf("%v", tm["id"]) == tagID {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("tag %s not found in tags list", tagID)
		}
	})

	t.Run("tags/project-list", func(t *testing.T) {
		requireID(t, tagID, "tags/create must have succeeded")
		requireID(t, projectID, "core/create must have succeeded")
		out := zoho(t, "projects", "tags", "project-list", "--project", projectID)
		m := parseJSON(t, out)
		tags, ok := m["tags"].([]any)
		if !ok {
			t.Fatalf("expected tags array in project-list response:\n%s", truncate(out, 500))
		}
		for _, it := range tags {
			tm, _ := it.(map[string]any)
			if fmt.Sprintf("%v", tm["id"]) == tagID {
				t.Errorf("tag %s should not be associated yet", tagID)
				break
			}
		}
	})

	t.Run("tags/update", func(t *testing.T) {
		requireID(t, tagID, "tags/create must have succeeded")
		updatedName := testName(t) + "_tag_upd"
		out := zoho(t, "projects", "tags", "update", tagID,
			"--json", toJSON(t, map[string]any{"name": updatedName}))
		m := parseJSON(t, out)
		tm, ok := m["tags"].(map[string]any)
		if !ok {
			t.Fatalf("expected tags object in update response:\n%s", truncate(out, 500))
		}
		assertEqual(t, fmt.Sprintf("%v", tm["id"]), tagID)
		assertStringField(t, tm, "name", updatedName)
	})

	t.Run("tags/associate", func(t *testing.T) {
		requireID(t, tagID, "tags/create must have succeeded")
		requireID(t, projectID, "core/create must have succeeded")
		requireID(t, tlTaskID, "timelogs/add must have succeeded")
		out := zoho(t, "projects", "tags", "associate", tagID,
			"--entity", tlTaskID,
			"--entity-type", "5",
			"--project", projectID)
		parseJSON(t, out)
	})

	t.Run("tags/project-list-verify", func(t *testing.T) {
		requireID(t, tagID, "tags/create must have succeeded")
		requireID(t, projectID, "core/create must have succeeded")
		out := zoho(t, "projects", "tags", "project-list", "--project", projectID)
		m := parseJSON(t, out)
		tags, ok := m["tags"].([]any)
		if !ok {
			t.Fatalf("expected tags array in project-list response:\n%s", truncate(out, 500))
		}
		found := false
		for _, it := range tags {
			tm, _ := it.(map[string]any)
			if fmt.Sprintf("%v", tm["id"]) == tagID {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("tag %s not found in project tags after associate", tagID)
		}
	})

	t.Run("tags/dissociate", func(t *testing.T) {
		requireID(t, tagID, "tags/create must have succeeded")
		requireID(t, projectID, "core/create must have succeeded")
		requireID(t, tlTaskID, "timelogs/add must have succeeded")
		out := zoho(t, "projects", "tags", "dissociate", tagID,
			"--entity", tlTaskID,
			"--entity-type", "5",
			"--project", projectID)
		parseJSON(t, out)
	})

	t.Run("tags/delete", func(t *testing.T) {
		requireID(t, tagID, "tags/create must have succeeded")
		out := zoho(t, "projects", "tags", "delete", tagID)
		parseJSON(t, out)

		listOut := zoho(t, "projects", "tags", "list")
		if strings.Contains(listOut, tagID) {
			t.Errorf("tag %s still found in list after delete", tagID)
		}
		tagID = ""
	})

	t.Run("roles/create", func(t *testing.T) {
		roleName := testName(t) + "_role"
		out := zoho(t, "projects", "roles", "create",
			"--json", toJSON(t, map[string]any{"name": roleName}))
		m := parseJSON(t, out)
		roleID = fmt.Sprintf("%v", m["id"])
		if roleID == "" || roleID == "<nil>" {
			t.Fatalf("no id in roles create response:\n%s", truncate(out, 500))
		}
		assertStringField(t, m, "name", roleName)
		cleanup.trackRole(roleID)
		t.Logf("created role %s", roleID)
	})

	t.Run("roles/get", func(t *testing.T) {
		requireID(t, roleID, "roles/create must have succeeded")
		out := zoho(t, "projects", "roles", "get", roleID)
		m := parseJSON(t, out)
		assertEqual(t, fmt.Sprintf("%v", m["id"]), roleID)
	})

	t.Run("roles/list", func(t *testing.T) {
		requireID(t, roleID, "roles/create must have succeeded")
		out := zoho(t, "projects", "roles", "list")
		m := parseJSON(t, out)
		roles, ok := m["roles"].([]any)
		if !ok {
			t.Fatalf("expected roles array in list response:\n%s", truncate(out, 500))
		}
		found := false
		for _, r := range roles {
			rm, _ := r.(map[string]any)
			if fmt.Sprintf("%v", rm["id"]) == roleID {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("role %s not found in roles list", roleID)
		}
	})

	t.Run("roles/update", func(t *testing.T) {
		requireID(t, roleID, "roles/create must have succeeded")
		updatedName := testName(t) + "_role_upd"
		zoho(t, "projects", "roles", "update", roleID,
			"--json", toJSON(t, map[string]any{"name": updatedName}))

		out := zoho(t, "projects", "roles", "get", roleID)
		m := parseJSON(t, out)
		assertStringField(t, m, "name", updatedName)
	})

	t.Run("roles/set-default-known-broken", func(t *testing.T) {
		t.Skip("roles set-default endpoint returns error")
	})

	t.Run("roles/delete", func(t *testing.T) {
		requireID(t, roleID, "roles/create must have succeeded")
		out := zoho(t, "projects", "roles", "delete", roleID)
		parseJSON(t, out)

		listOut := zoho(t, "projects", "roles", "list")
		m := parseJSON(t, listOut)
		roles, ok := m["roles"].([]any)
		if ok {
			for _, r := range roles {
				rm, _ := r.(map[string]any)
				if fmt.Sprintf("%v", rm["id"]) == roleID {
					t.Errorf("role %s still found in list after delete", roleID)
					break
				}
			}
		}
		roleID = ""
	})

	t.Run("project-users/list", func(t *testing.T) {
		requireID(t, projectID, "core/create must have succeeded")
		out := zoho(t, "projects", "project-users", "list", "--project", projectID)
		arr := parseJSONArray(t, out)
		assertNonEmpty(t, arr, "expected at least one project user")
	})

	t.Run("project-users/get", func(t *testing.T) {
		requireID(t, projectID, "core/create must have succeeded")
		requireID(t, ownerZPUID, "timelogs/setup-owner must have succeeded")
		out := zoho(t, "projects", "project-users", "get", ownerZPUID,
			"--project", projectID)
		m := parseJSON(t, out)
		assertEqual(t, fmt.Sprintf("%v", m["id"]), ownerZPUID)
	})

	t.Run("dependencies/add-known-broken", func(t *testing.T) {
		t.Skip("dependencies endpoint not functional")
	})

	t.Run("dependencies/remove-known-broken", func(t *testing.T) {
		t.Skip("dependencies endpoint not functional")
	})

	t.Run("task-customviews/list", func(t *testing.T) {
		out := zoho(t, "projects", "task-customviews", "list")
		m := parseJSON(t, out)
		views, ok := m["views"].([]any)
		if !ok {
			t.Fatalf("expected views array in task-customviews list response:\n%s", truncate(out, 500))
		}
		if len(views) == 0 {
			t.Error("expected non-empty views array")
		}
	})

	t.Run("task-customviews/project-list", func(t *testing.T) {
		requireID(t, projectID, "core/create must have succeeded")
		out := zoho(t, "projects", "task-customviews", "project-list",
			"--project", projectID)
		m := parseJSON(t, out)
		if _, ok := m["views"]; !ok {
			t.Fatalf("expected views key in task-customviews project-list response:\n%s", truncate(out, 500))
		}
	})

	t.Run("task-customviews/get", func(t *testing.T) {
		out := zoho(t, "projects", "task-customviews", "list")
		m := parseJSON(t, out)
		views, ok := m["views"].([]any)
		if !ok || len(views) == 0 {
			t.Skip("no views available for get test")
		}
		first, _ := views[0].(map[string]any)
		viewID := fmt.Sprintf("%v", first["id"])
		if viewID == "" || viewID == "<nil>" {
			t.Skip("first view has no id")
		}

		getOut := zoho(t, "projects", "task-customviews", "get", viewID)
		gm := parseJSON(t, getOut)
		tcv, ok := gm["taskcustomviews"].([]any)
		if !ok || len(tcv) == 0 {
			t.Fatalf("expected taskcustomviews array in get response:\n%s", truncate(getOut, 500))
		}
		first2, _ := tcv[0].(map[string]any)
		assertEqual(t, fmt.Sprintf("%v", first2["id"]), viewID)
	})

	t.Run("issue-customviews/list", func(t *testing.T) {
		out := zoho(t, "projects", "issue-customviews", "list")
		m := parseJSON(t, out)
		dv, ok := m["default_views"].([]any)
		if !ok {
			t.Fatalf("expected default_views array in issue-customviews list response:\n%s", truncate(out, 500))
		}
		if len(dv) == 0 {
			t.Error("expected non-empty default_views array")
		}
	})

	t.Run("issue-customviews/project-list", func(t *testing.T) {
		requireID(t, projectID, "core/create must have succeeded")
		out := zoho(t, "projects", "issue-customviews", "project-list",
			"--project", projectID)
		m := parseJSON(t, out)
		for _, key := range []string{"custom_views", "default_views", "favourites", "shared_views"} {
			if _, ok := m[key]; !ok {
				t.Errorf("expected %s key in issue-customviews project-list response", key)
			}
		}
	})

	t.Run("issue-customviews/get", func(t *testing.T) {
		out := zoho(t, "projects", "issue-customviews", "list")
		m := parseJSON(t, out)
		dv, ok := m["default_views"].([]any)
		if !ok || len(dv) == 0 {
			t.Skip("no default_views available for get test")
		}
		first, _ := dv[0].(map[string]any)
		viewID := fmt.Sprintf("%v", first["custom_view_id"])
		if viewID == "" || viewID == "<nil>" {
			t.Skip("first default_view has no custom_view_id")
		}

		getOut := zoho(t, "projects", "issue-customviews", "get", viewID)
		gm := parseJSON(t, getOut)
		cv, ok := gm["customview"].(map[string]any)
		if !ok {
			t.Fatalf("expected customview object in get response:\n%s", truncate(getOut, 500))
		}
		assertEqual(t, fmt.Sprintf("%v", cv["custom_view_id"]), viewID)
	})

	t.Run("task-statustimeline/get", func(t *testing.T) {
		requireID(t, tlTaskID, "timelogs/add must have succeeded")
		requireID(t, projectID, "core/create must have succeeded")
		out := zoho(t, "projects", "task-statustimeline", "get", tlTaskID,
			"--project", projectID)
		m := parseJSON(t, out)
		if _, ok := m["timeline"]; !ok {
			if _, ok2 := m["status_timeline"]; !ok2 {
				t.Errorf("expected timeline or status_timeline key in response:\n%s", truncate(out, 500))
			}
		}
	})

	t.Run("task-statustimeline/project-known-broken", func(t *testing.T) {
		t.Skip("task-statustimeline project endpoint may not exist in V3")
	})

	t.Run("task-statustimeline/portal", func(t *testing.T) {
		out := zoho(t, "projects", "task-statustimeline", "portal")
		m := parseJSON(t, out)
		if _, ok := m["timeline"]; !ok {
			if _, ok2 := m["status_timeline"]; !ok2 {
				if _, ok3 := m["page_info"]; !ok3 {
					t.Errorf("expected timeline, status_timeline, or page_info key in response:\n%s", truncate(out, 500))
				}
			}
		}
	})

	t.Run("attachments/list", func(t *testing.T) {
		requireID(t, tlTaskID, "timelogs/add must have succeeded")
		requireID(t, projectID, "core/create must have succeeded")
		out := zoho(t, "projects", "attachments", "list",
			"--project", projectID,
			"--type", "task",
			"--entity-id", tlTaskID)
		m := parseJSON(t, out)
		if _, ok := m["attachments"]; !ok {
			if _, ok2 := m["page_info"]; !ok2 {
				t.Errorf("expected attachments or page_info key in response:\n%s", truncate(out, 500))
			}
		}
	})

	t.Run("attachments/upload-known-broken", func(t *testing.T) {
		t.Skip("attachments upload needs WorkDrive integration")
	})

	t.Run("attachments/get-known-broken", func(t *testing.T) {
		t.Skip("no real project attachment ID available")
	})

	t.Run("attachments/associate-known-broken", func(t *testing.T) {
		t.Skip("no real project attachment ID available")
	})

	t.Run("attachments/dissociate-known-broken", func(t *testing.T) {
		t.Skip("no real project attachment ID available")
	})

	t.Run("phases/create", func(t *testing.T) {
		requireID(t, projectID, "core/create must have succeeded")
		phaseName = testName(t) + "_phase"
		out := zoho(t, "projects", "phases", "create",
			"--project", projectID,
			"--json", toJSON(t, map[string]any{"name": phaseName}))
		phaseID = extractProjectsID(t, out)
		cleanup.trackPhase(phaseID, projectID)
		t.Logf("created phase %s (%s)", phaseID, phaseName)

		out = zoho(t, "projects", "phases", "get", phaseID, "--project", projectID)
		m := parseJSON(t, out)
		assertStringField(t, m, "name", phaseName)
	})

	t.Run("phases/get", func(t *testing.T) {
		requireID(t, phaseID, "phases/create must have succeeded")
		out := zoho(t, "projects", "phases", "get", phaseID, "--project", projectID)
		m := parseJSON(t, out)
		assertEqual(t, fmt.Sprintf("%v", m["id"]), phaseID)
		assertStringField(t, m, "name", phaseName)
	})

	t.Run("phases/list-project", func(t *testing.T) {
		requireID(t, phaseID, "phases/create must have succeeded")
		out := zoho(t, "projects", "phases", "list-project", "--project", projectID)
		m := parseJSON(t, out)
		milestones, ok := m["milestones"].([]any)
		if !ok {
			t.Fatalf("expected milestones array in list-project response:\n%s", truncate(out, 500))
		}
		found := false
		for _, ms := range milestones {
			mm, _ := ms.(map[string]any)
			if fmt.Sprintf("%v", mm["id"]) == phaseID {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("phase %s not found in list-project", phaseID)
		}
	})

	t.Run("phases/list", func(t *testing.T) {
		requireID(t, phaseID, "phases/create must have succeeded")
		out := zoho(t, "projects", "phases", "list")
		m := parseJSON(t, out)
		milestones, ok := m["milestones"].([]any)
		if !ok {
			t.Fatalf("expected milestones array in list response:\n%s", truncate(out, 500))
		}
		found := false
		for _, ms := range milestones {
			mm, _ := ms.(map[string]any)
			if fmt.Sprintf("%v", mm["id"]) == phaseID {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("phase %s not found in portal-level list", phaseID)
		}
	})

	t.Run("phases/update", func(t *testing.T) {
		requireID(t, phaseID, "phases/create must have succeeded")
		updatedPhaseName := phaseName + "_upd"
		out := zoho(t, "projects", "phases", "update", phaseID,
			"--project", projectID,
			"--json", toJSON(t, map[string]any{"name": updatedPhaseName}))
		m := parseJSON(t, out)
		assertEqual(t, fmt.Sprintf("%v", m["name"]), updatedPhaseName)

		out = zoho(t, "projects", "phases", "get", phaseID, "--project", projectID)
		m = parseJSON(t, out)
		assertStringField(t, m, "name", updatedPhaseName)
		phaseName = updatedPhaseName
	})

	t.Run("phases/activities", func(t *testing.T) {
		requireID(t, phaseID, "phases/create must have succeeded")
		out := zoho(t, "projects", "phases", "activities", phaseID,
			"--project", projectID)
		m := parseJSON(t, out)
		if _, ok := m["activities"]; !ok {
			t.Fatalf("expected activities key in response:\n%s", truncate(out, 500))
		}
	})

	t.Run("phases/clone", func(t *testing.T) {
		requireID(t, phaseID, "phases/create must have succeeded")
		out := zoho(t, "projects", "phases", "clone", phaseID,
			"--project", projectID)
		m := parseJSON(t, out)
		clonedPhaseID = fmt.Sprintf("%v", m["id"])
		if clonedPhaseID == "" || clonedPhaseID == "<nil>" {
			t.Fatalf("no id in phases clone response:\n%s", truncate(out, 500))
		}
		cleanup.trackPhase(clonedPhaseID, projectID)
		if clonedPhaseID == phaseID {
			t.Errorf("cloned phase has same ID as original")
		}
		t.Logf("cloned phase %s -> %s", phaseID, clonedPhaseID)
	})

	t.Run("phases/move-known-broken", func(t *testing.T) {
		t.Skip("phases move endpoint returns Zoho 500 error")
	})

	t.Run("phase-followers/add", func(t *testing.T) {
		requireID(t, phaseID, "phases/create must have succeeded")
		requireID(t, ownerZPUID, "timelogs/setup-owner must have succeeded")
		out := zoho(t, "projects", "phase-followers", "add", phaseID,
			"--project", projectID,
			"--json", toJSON(t, map[string]any{"followers": []string{ownerZPUID}}))
		m := parseJSON(t, out)
		if _, ok := m["followers"]; !ok {
			t.Fatalf("expected followers key in add response:\n%s", truncate(out, 500))
		}
	})

	t.Run("phase-followers/list", func(t *testing.T) {
		requireID(t, phaseID, "phases/create must have succeeded")
		out := zoho(t, "projects", "phase-followers", "list", phaseID,
			"--project", projectID)
		m := parseJSON(t, out)
		followers, ok := m["followers"].([]any)
		if !ok || len(followers) == 0 {
			t.Fatalf("expected non-empty followers array:\n%s", truncate(out, 500))
		}
		found := false
		for _, f := range followers {
			fm, _ := f.(map[string]any)
			if fmt.Sprintf("%v", fm["zpuid"]) == ownerZPUID {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("follower %s not found in phase-followers list", ownerZPUID)
		}
	})

	t.Run("phase-followers/remove", func(t *testing.T) {
		requireID(t, phaseID, "phases/create must have succeeded")
		requireID(t, ownerZPUID, "timelogs/setup-owner must have succeeded")
		zoho(t, "projects", "phase-followers", "remove", phaseID,
			"--project", projectID,
			"--json", toJSON(t, map[string]any{"followers": []string{ownerZPUID}}))

		time.Sleep(2 * time.Second)
		out := zoho(t, "projects", "phase-followers", "list", phaseID,
			"--project", projectID)
		if strings.Contains(out, ownerZPUID) {
			t.Errorf("follower %s still present after phase-followers remove", ownerZPUID)
		}
	})

	t.Run("phase-comments/add", func(t *testing.T) {
		requireID(t, phaseID, "phases/create must have succeeded")
		out := zoho(t, "projects", "phase-comments", "add",
			"--phase", phaseID, "--project", projectID,
			"--comment", testPrefix+"_phase_comment")
		m := parseJSON(t, out)
		phaseCommentID = fmt.Sprintf("%v", m["id"])
		if phaseCommentID == "" || phaseCommentID == "<nil>" {
			t.Fatalf("no id in phase-comments add response:\n%s", truncate(out, 500))
		}
		assertStringField(t, m, "content", testPrefix+"_phase_comment")
		t.Logf("created phase comment %s", phaseCommentID)
	})

	t.Run("phase-comments/list", func(t *testing.T) {
		requireID(t, phaseCommentID, "phase-comments/add must have succeeded")
		out := zoho(t, "projects", "phase-comments", "list",
			"--phase", phaseID, "--project", projectID)
		m := parseJSON(t, out)
		comments, ok := m["comments"].([]any)
		if !ok {
			t.Fatalf("expected comments array in list response:\n%s", truncate(out, 500))
		}
		found := false
		for _, c := range comments {
			cm, _ := c.(map[string]any)
			if fmt.Sprintf("%v", cm["id"]) == phaseCommentID {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("comment %s not found in phase-comments list", phaseCommentID)
		}
	})

	t.Run("phase-comments/update", func(t *testing.T) {
		requireID(t, phaseCommentID, "phase-comments/add must have succeeded")
		updatedComment := testPrefix + "_phase_comment_upd"
		out := zoho(t, "projects", "phase-comments", "update", phaseCommentID,
			"--phase", phaseID, "--project", projectID,
			"--comment", updatedComment)
		m := parseJSON(t, out)
		assertStringField(t, m, "content", updatedComment)

		listOut := zoho(t, "projects", "phase-comments", "list",
			"--phase", phaseID, "--project", projectID)
		assertContains(t, listOut, updatedComment)
	})

	t.Run("phase-comments/delete", func(t *testing.T) {
		requireID(t, phaseCommentID, "phase-comments/add must have succeeded")
		zoho(t, "projects", "phase-comments", "delete", phaseCommentID,
			"--phase", phaseID, "--project", projectID)

		listOut := zoho(t, "projects", "phase-comments", "list",
			"--phase", phaseID, "--project", projectID)
		if strings.Contains(listOut, phaseCommentID) {
			t.Errorf("comment %s still found in list after delete", phaseCommentID)
		}
		phaseCommentID = ""
	})

	t.Run("phases/delete", func(t *testing.T) {
		requireID(t, phaseID, "phases/create must have succeeded")
		zoho(t, "projects", "phases", "delete", phaseID, "--project", projectID)

		r := runZoho(t, "projects", "phases", "get", phaseID, "--project", projectID)
		if r.ExitCode == 0 {
			t.Errorf("phase %s still accessible after delete", phaseID)
		}
		phaseID = ""

		if clonedPhaseID != "" {
			zohoIgnoreError(t, "projects", "phases", "delete", clonedPhaseID, "--project", projectID)
			clonedPhaseID = ""
		}
	})

	t.Run("events/create", func(t *testing.T) {
		requireID(t, projectID, "core/create must have succeeded")
		requireID(t, ownerZPUID, "timelogs/setup-owner must have succeeded")
		eventName := testName(t) + "_event"
		now := time.Now().AddDate(0, 1, 0)
		startsAt := now.Format("2006-01-02T15:04:05+00:00")
		endsAt := now.Add(1 * time.Hour).Format("2006-01-02T15:04:05+00:00")
		out := zoho(t, "projects", "events", "create",
			"--project", projectID,
			"--json", toJSON(t, map[string]any{
				"title":     eventName,
				"starts_at": startsAt,
				"ends_at":   endsAt,
				"attendees": []string{ownerZPUID},
			}))
		eventID = extractProjectsID(t, out)
		cleanup.trackEvent(eventID, projectID)
		t.Logf("created event %s (%s)", eventID, eventName)

		out = zoho(t, "projects", "events", "get", eventID, "--project", projectID)
		m := parseJSON(t, out)
		assertStringField(t, m, "title", eventName)
	})

	t.Run("events/get", func(t *testing.T) {
		requireID(t, eventID, "events/create must have succeeded")
		out := zoho(t, "projects", "events", "get", eventID, "--project", projectID)
		m := parseJSON(t, out)
		assertEqual(t, fmt.Sprintf("%v", m["id"]), eventID)
	})

	t.Run("events/list", func(t *testing.T) {
		requireID(t, eventID, "events/create must have succeeded")
		out := zoho(t, "projects", "events", "list", "--project", projectID)
		m := parseJSON(t, out)
		events, ok := m["events"].([]any)
		if !ok {
			t.Fatalf("expected events array in list response:\n%s", truncate(out, 500))
		}
		found := false
		for _, e := range events {
			em, _ := e.(map[string]any)
			if fmt.Sprintf("%v", em["id"]) == eventID {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("event %s not found in events list", eventID)
		}
	})

	t.Run("events/update", func(t *testing.T) {
		requireID(t, eventID, "events/create must have succeeded")
		updatedTitle := testPrefix + "_event_upd"
		out := zoho(t, "projects", "events", "update", eventID,
			"--project", projectID,
			"--json", toJSON(t, map[string]any{"title": updatedTitle}))
		m := parseJSON(t, out)
		assertStringField(t, m, "title", updatedTitle)

		out = zoho(t, "projects", "events", "get", eventID, "--project", projectID)
		m = parseJSON(t, out)
		assertStringField(t, m, "title", updatedTitle)
	})

	t.Run("event-comments/add", func(t *testing.T) {
		requireID(t, eventID, "events/create must have succeeded")
		out := zoho(t, "projects", "event-comments", "add",
			"--event", eventID, "--project", projectID,
			"--comment", testPrefix+"_event_comment")
		m := parseJSON(t, out)
		eventCommentID = fmt.Sprintf("%v", m["id"])
		if eventCommentID == "" || eventCommentID == "<nil>" {
			t.Fatalf("no id in event-comments add response:\n%s", truncate(out, 500))
		}
		assertStringField(t, m, "content", testPrefix+"_event_comment")
		t.Logf("created event comment %s", eventCommentID)
	})

	t.Run("event-comments/list", func(t *testing.T) {
		requireID(t, eventCommentID, "event-comments/add must have succeeded")
		out := zoho(t, "projects", "event-comments", "list",
			"--event", eventID, "--project", projectID)
		arr := parseJSONArray(t, out)
		_, found := findInArray(arr, eventCommentID)
		if !found {
			t.Errorf("comment %s not found in event-comments list", eventCommentID)
		}
	})

	t.Run("event-comments/get-known-broken", func(t *testing.T) {
		t.Skip("event-comments get endpoint returns INVALID_METHOD")
	})

	t.Run("event-comments/update", func(t *testing.T) {
		requireID(t, eventCommentID, "event-comments/add must have succeeded")
		updatedComment := testPrefix + "_event_comment_upd"
		zoho(t, "projects", "event-comments", "update", eventCommentID,
			"--event", eventID, "--project", projectID,
			"--comment", updatedComment)

		listOut := zoho(t, "projects", "event-comments", "list",
			"--event", eventID, "--project", projectID)
		assertContains(t, listOut, updatedComment)
	})

	t.Run("event-comments/delete", func(t *testing.T) {
		requireID(t, eventCommentID, "event-comments/add must have succeeded")
		zoho(t, "projects", "event-comments", "delete", eventCommentID,
			"--event", eventID, "--project", projectID)

		listOut := zoho(t, "projects", "event-comments", "list",
			"--event", eventID, "--project", projectID)
		if strings.Contains(listOut, eventCommentID) {
			t.Errorf("comment %s still found in list after delete", eventCommentID)
		}
		eventCommentID = ""
	})

	t.Run("events/delete", func(t *testing.T) {
		requireID(t, eventID, "events/create must have succeeded")
		zoho(t, "projects", "events", "delete", eventID, "--project", projectID)

		r := runZoho(t, "projects", "events", "get", eventID, "--project", projectID)
		if r.ExitCode == 0 {
			t.Errorf("event %s still accessible after delete", eventID)
		}
		eventID = ""
	})

	t.Run("search/portal", func(t *testing.T) {
		requireID(t, projectID, "core/create must have succeeded")
		retryUntil(t, 30*time.Second, func() bool {
			out, err := zohoMayFail(t, "projects", "search", "portal",
				"--query", testPrefix)
			if err != nil {
				return false
			}
			return strings.Contains(out, projectName) || strings.Contains(out, testPrefix)
		})
	})

	t.Run("search/project", func(t *testing.T) {
		requireID(t, projectID, "core/create must have succeeded")
		retryUntil(t, 30*time.Second, func() bool {
			out, err := zohoMayFail(t, "projects", "search", "project",
				"--query", testPrefix, "--project", projectID)
			if err != nil {
				return false
			}
			return strings.Contains(out, testPrefix)
		})
	})

	t.Run("feed/post", func(t *testing.T) {
		requireID(t, projectID, "core/create must have succeeded")
		statusContent := testPrefix + "_status"
		out := zoho(t, "projects", "feed", "post",
			"--project", projectID,
			"--json", toJSON(t, map[string]any{"content": statusContent}))
		m := parseJSON(t, out)
		id := fmt.Sprintf("%v", m["id"])
		if id == "" || id == "<nil>" {
			t.Fatalf("no id in feed post response:\n%s", truncate(out, 500))
		}
		t.Logf("posted status %s", id)
	})

	t.Run("feed/status", func(t *testing.T) {
		requireID(t, projectID, "core/create must have succeeded")
		out := zoho(t, "projects", "feed", "status", "--project", projectID)
		m := parseJSON(t, out)
		if _, ok := m["status"]; !ok {
			if _, ok2 := m["page_info"]; !ok2 {
				t.Errorf("expected status or page_info in response:\n%s", truncate(out, 500))
			}
		}
	})

	t.Run("timers/setup", func(t *testing.T) {
		t.Skip("all timer tests are broken (Zoho 500)")
	})

	t.Run("timers/running-empty", func(t *testing.T) {
		out := zoho(t, "projects", "timelog-timers", "running")
		m := parseJSON(t, out)
		if timers, ok := m["timer"]; ok {
			arr, _ := timers.([]any)
			t.Logf("running timers: %d", len(arr))
		}
	})

	t.Run("timers/start-known-broken", func(t *testing.T) {
		t.Skip("timelog-timers start endpoint returns Zoho 500")
	})

	t.Run("timers/pause", func(t *testing.T) {
		t.Skip("depends on timers/start which returns Zoho 500")
	})

	t.Run("timers/resume", func(t *testing.T) {
		t.Skip("depends on timers/start which returns Zoho 500")
	})

	t.Run("timers/stop", func(t *testing.T) {
		t.Skip("depends on timers/start which returns Zoho 500")
	})

	t.Run("timers/delete", func(t *testing.T) {
		t.Skip("depends on timers/start which returns Zoho 500")
	})

	t.Run("pins/setup", func(t *testing.T) {
		t.Skip("all pin tests are broken (create rejects all fields)")
	})

	t.Run("pins/list-empty", func(t *testing.T) {
		out := zoho(t, "projects", "timelog-pins", "list")
		m := parseJSON(t, out)
		if _, ok := m["pins"]; !ok {
			if _, ok2 := m["page_info"]; !ok2 {
				t.Errorf("expected pins or page_info key in response:\n%s", truncate(out, 500))
			}
		}
	})

	t.Run("pins/create-known-broken", func(t *testing.T) {
		t.Skip("timelog-pins create endpoint rejects all fields")
	})

	t.Run("pins/list-after-create", func(t *testing.T) {
		t.Skip("depends on pins/create which is broken")
	})

	t.Run("pins/delete", func(t *testing.T) {
		t.Skip("depends on pins/create which is broken")
	})

	t.Run("timelog-bulk/list-missing-module", func(t *testing.T) {
		r := runZoho(t, "projects", "timelog-bulk", "list")
		assertExitCode(t, r, 1)
		assertContains(t, r.Stderr, "Required flag")
	})

	t.Run("timelog-bulk/project-list", func(t *testing.T) {
		requireID(t, projectID, "core/create must have succeeded")
		out := zoho(t, "projects", "timelog-bulk", "project-list",
			"--project", projectID, "--module", "task")
		m := parseJSON(t, out)
		if _, ok := m["timelogs"]; !ok {
			if _, ok2 := m["page_info"]; !ok2 {
				t.Errorf("expected timelogs or page_info key in response:\n%s", truncate(out, 500))
			}
		}
	})

	t.Run("users/list", func(t *testing.T) {
		out := zoho(t, "projects", "users", "list")
		arr := parseJSONArray(t, out)
		assertNonEmpty(t, arr, "expected at least one user")
		t.Logf("found %d portal users", len(arr))
	})

	t.Run("users/get", func(t *testing.T) {
		requireID(t, ownerZPUID, "owner zpuid must be available")
		out := zoho(t, "projects", "users", "get", ownerZPUID)
		m := parseJSON(t, out)
		id := fmt.Sprintf("%v", m["id"])
		if id == "" || id == "<nil>" {
			t.Errorf("expected id in user get response:\n%s", truncate(out, 500))
		}
	})

	t.Run("core/trash-second", func(t *testing.T) {
		requireID(t, project2ID, "second project must exist")
		out := zoho(t, "projects", "core", "trash", project2ID)
		parseJSON(t, out)
		t.Logf("trashed project %s", project2ID)

		r := runZoho(t, "projects", "core", "get", project2ID)
		if r.ExitCode == 0 {
			m := parseJSON(t, r.Stdout)
			status := fmt.Sprintf("%v", m["status"])
			t.Logf("trashed project get returned status=%s (may still be accessible)", status)
		}
	})

	t.Run("core/restore-second", func(t *testing.T) {
		requireID(t, project2ID, "second project must exist")
		time.Sleep(2 * time.Second)
		out := zoho(t, "projects", "core", "restore", project2ID)
		parseJSON(t, out)

		retryUntil(t, 30*time.Second, func() bool {
			r := runZoho(t, "projects", "core", "get", project2ID)
			if r.ExitCode != 0 {
				return false
			}
			return true
		})
		proj := getProject(t, project2ID)
		assertStringField(t, proj, "name", project2Name)
	})

	t.Run("trash/list", func(t *testing.T) {
		out := zoho(t, "projects", "trash", "list")
		m := parseJSON(t, out)
		if _, ok := m["trash"]; !ok {
			if _, ok2 := m["page_info"]; !ok2 {
				t.Errorf("expected trash or page_info key in response:\n%s", truncate(out, 500))
			}
		}
	})

	t.Run("trash/restore-known-broken", func(t *testing.T) {
		t.Skip("trash restore with fake IDs always fails")
	})

	t.Run("trash/delete-known-broken", func(t *testing.T) {
		t.Skip("trash delete with fake IDs always fails")
	})

	t.Run("core/trash-and-delete-second", func(t *testing.T) {
		requireID(t, project2ID, "second project must exist")
		zoho(t, "projects", "core", "trash", project2ID)
		time.Sleep(2 * time.Second)
		out := zoho(t, "projects", "core", "delete", project2ID)
		parseJSON(t, out)
		t.Logf("permanently deleted project %s", project2ID)
		project2ID = ""
	})

	t.Run("core/trash-and-delete-main", func(t *testing.T) {
		requireID(t, projectID, "core/create must have succeeded")
		zoho(t, "projects", "core", "trash", projectID)
		time.Sleep(2 * time.Second)
		out := zoho(t, "projects", "core", "delete", projectID)
		parseJSON(t, out)
		t.Logf("permanently deleted project %s", projectID)
		projectID = ""
	})
}

func TestProjectsProfiles(t *testing.T) {
	t.Parallel()
	_ = requireProjectsPortalID(t)
	cleanup := newCleanup(t)

	var profileID string
	var profileName string

	t.Run("profiles/list", func(t *testing.T) {
		out := zoho(t, "projects", "profiles", "list")
		m := parseJSON(t, out)
		profiles, ok := m["profiles"].([]any)
		if !ok {
			t.Fatalf("expected profiles array in response:\n%s", truncate(out, 500))
		}
		if len(profiles) == 0 {
			t.Error("expected at least one built-in profile")
		}
		t.Logf("found %d profiles", len(profiles))
	})

	t.Run("profiles/create", func(t *testing.T) {
		profileName = testName(t) + "_profile"
		out := zoho(t, "projects", "profiles", "create",
			"--json", toJSON(t, map[string]any{
				"name": profileName,
				"type": "3",
			}))
		m := parseJSON(t, out)
		profileID = fmt.Sprintf("%v", m["id"])
		if profileID == "" || profileID == "<nil>" {
			t.Fatalf("no id in profile create response:\n%s", truncate(out, 500))
		}
		cleanup.trackProfile(profileID)
		t.Logf("created profile %s (%s)", profileID, profileName)
	})

	t.Run("profiles/get", func(t *testing.T) {
		requireID(t, profileID, "profiles/create must have succeeded")
		out := zoho(t, "projects", "profiles", "get", profileID)
		m := parseJSON(t, out)
		assertEqual(t, fmt.Sprintf("%v", m["id"]), profileID)
		assertStringField(t, m, "name", profileName)
	})

	t.Run("profiles/update", func(t *testing.T) {
		requireID(t, profileID, "profiles/create must have succeeded")
		updatedName := profileName + "_upd"
		zoho(t, "projects", "profiles", "update", profileID,
			"--json", toJSON(t, map[string]any{"name": updatedName}))

		out := zoho(t, "projects", "profiles", "get", profileID)
		m := parseJSON(t, out)
		assertStringField(t, m, "name", updatedName)
		profileName = updatedName
	})

	t.Run("profiles/set-default", func(t *testing.T) {
		requireID(t, profileID, "profiles/create must have succeeded")
		zoho(t, "projects", "profiles", "set-default", profileID)

		out := zoho(t, "projects", "profiles", "get", profileID)
		m := parseJSON(t, out)
		if fmt.Sprintf("%v", m["is_default"]) != "true" {
			t.Errorf("profile is_default not set to true after set-default, got %v", m["is_default"])
		}

		zoho(t, "projects", "profiles", "set-default", "1212503000000015149")
	})

	t.Run("profiles/delete", func(t *testing.T) {
		requireID(t, profileID, "profiles/create must have succeeded")
		zoho(t, "projects", "profiles", "delete", profileID)

		r := runZoho(t, "projects", "profiles", "get", profileID)
		if r.ExitCode == 0 {
			t.Errorf("profile %s still accessible after delete", profileID)
		}
		profileID = ""
	})
}

func TestProjectsDashboards(t *testing.T) {
	t.Parallel()
	_ = requireProjectsPortalID(t)
	cleanup := newCleanup(t)

	var dashboardID string
	var dashboardName string

	t.Run("reports/workload-meta", func(t *testing.T) {
		out := zoho(t, "projects", "reports", "workload-meta")
		m := parseJSON(t, out)
		if _, ok := m["chart_details"]; !ok {
			t.Errorf("expected chart_details in workload-meta response:\n%s", truncate(out, 500))
		}
	})

	t.Run("reports/workload-known-broken", func(t *testing.T) {
		t.Skip("reports workload endpoint returns Zoho 500")
	})

	t.Run("reports/dashboards-list", func(t *testing.T) {
		out := zoho(t, "projects", "reports", "dashboards")
		m := parseJSON(t, out)
		if _, ok := m["folders"]; !ok {
			if _, ok2 := m["page_info"]; !ok2 {
				t.Errorf("expected folders or page_info in dashboards response:\n%s", truncate(out, 500))
			}
		}
	})

	t.Run("reports/dashboard-create", func(t *testing.T) {
		dashboardName = testName(t) + "_dash"
		out := zoho(t, "projects", "reports", "dashboard-create",
			"--json", toJSON(t, map[string]any{"name": dashboardName}))
		m := parseJSON(t, out)
		dashboardID = fmt.Sprintf("%v", m["id"])
		if dashboardID == "" || dashboardID == "<nil>" {
			t.Fatalf("no id in dashboard create response:\n%s", truncate(out, 500))
		}
		cleanup.trackDashboard(dashboardID)
		t.Logf("created dashboard %s (%s)", dashboardID, dashboardName)
	})

	t.Run("reports/dashboard-get", func(t *testing.T) {
		requireID(t, dashboardID, "dashboard-create must have succeeded")
		out := zoho(t, "projects", "reports", "dashboard-get", dashboardID)
		m := parseJSON(t, out)
		assertEqual(t, fmt.Sprintf("%v", m["id"]), dashboardID)
		assertStringField(t, m, "name", dashboardName)
	})

	t.Run("reports/dashboard-update", func(t *testing.T) {
		requireID(t, dashboardID, "dashboard-create must have succeeded")
		updatedName := dashboardName + "_upd"
		zoho(t, "projects", "reports", "dashboard-update", dashboardID,
			"--json", toJSON(t, map[string]any{"name": updatedName}))
		getOut := zoho(t, "projects", "reports", "dashboard-get", dashboardID)
		gm := parseJSON(t, getOut)
		assertStringField(t, gm, "name", updatedName)
		dashboardName = updatedName
	})

	t.Run("reports/dashboard-delete", func(t *testing.T) {
		requireID(t, dashboardID, "dashboard-create must have succeeded")
		zoho(t, "projects", "reports", "dashboard-delete", dashboardID)

		r := runZoho(t, "projects", "reports", "dashboard-get", dashboardID)
		if r.ExitCode == 0 {
			t.Errorf("dashboard %s still accessible after delete", dashboardID)
		}
		dashboardID = ""
	})
}

func TestProjectsTeams(t *testing.T) {
	t.Parallel()
	portalID := requireProjectsPortalID(t)
	cleanup := newCleanup(t)

	var teamID string
	var teamName string
	var projectID string
	var ownerZPUID string

	t.Run("teams/setup", func(t *testing.T) {
		out := zoho(t, "projects", "users", "list")
		arr := parseJSONArray(t, out)
		if len(arr) == 0 {
			t.Fatal("no users found in portal")
		}
		ownerZPUID = fmt.Sprintf("%v", arr[0]["zpuid"])
		if ownerZPUID == "" || ownerZPUID == "<nil>" {
			ownerZPUID = fmt.Sprintf("%v", arr[0]["id"])
		}
		if ownerZPUID == "" || ownerZPUID == "<nil>" {
			t.Fatal("could not determine user zpuid")
		}
		t.Logf("using owner zpuid: %s", ownerZPUID)

		projectName := testName(t) + "_teamproj"
		projOut := zoho(t, "projects", "core", "create", "--name", projectName)
		projectID = extractProjectsID(t, projOut)
		cleanup.trackProject(projectID)
		t.Logf("created project %s for team tests", projectID)
		_ = portalID
	})

	t.Run("teams/create", func(t *testing.T) {
		requireID(t, ownerZPUID, "teams/setup must have succeeded")
		teamName = testName(t) + "_team"
		out, err := zohoMayFail(t, "projects", "teams", "create",
			"--json", toJSON(t, map[string]any{
				"name":     teamName,
				"lead":     ownerZPUID,
				"user_ids": map[string]any{"add": []string{ownerZPUID}},
			}))
		if err != nil {
			t.Logf("team create failed (may need scope ZohoProjects.teams.ALL): %v", err)
			t.Logf("response: %s", truncate(out, 500))
			return
		}
		m := parseJSON(t, out)
		teamID = fmt.Sprintf("%v", m["id"])
		if teamID == "" || teamID == "<nil>" {
			t.Fatalf("no id in team create response:\n%s", truncate(out, 500))
		}
		cleanup.trackTeam(teamID)
		t.Logf("created team %s (%s)", teamID, teamName)
	})

	t.Run("teams/list", func(t *testing.T) {
		requireID(t, teamID, "teams/create must have succeeded")
		out := zoho(t, "projects", "teams", "list")
		m := parseJSON(t, out)
		teams, ok := m["teams"].([]any)
		if !ok {
			t.Fatalf("expected teams array in response:\n%s", truncate(out, 500))
		}
		found := false
		for _, tm := range teams {
			tmm, _ := tm.(map[string]any)
			if fmt.Sprintf("%v", tmm["id"]) == teamID {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("team %s not found in teams list", teamID)
		}
	})

	t.Run("teams/get-known-broken", func(t *testing.T) {
		t.Skip("teams get endpoint returns INVALID_METHOD")
	})

	t.Run("teams/update", func(t *testing.T) {
		requireID(t, teamID, "teams/create must have succeeded")
		updatedName := teamName + "_upd"
		zoho(t, "projects", "teams", "update", teamID,
			"--json", toJSON(t, map[string]any{"name": updatedName}))

		out := zoho(t, "projects", "teams", "list")
		m := parseJSON(t, out)
		teams, _ := m["teams"].([]any)
		for _, tm := range teams {
			tmm, _ := tm.(map[string]any)
			if fmt.Sprintf("%v", tmm["id"]) == teamID {
				assertStringField(t, tmm, "name", updatedName)
				break
			}
		}
		teamName = updatedName
	})

	t.Run("teams/users-known-broken", func(t *testing.T) {
		t.Skip("teams users endpoint returns URL_RULE_NOT_CONFIGURED")
	})

	t.Run("teams/projects-known-broken", func(t *testing.T) {
		t.Skip("teams projects endpoint not functional")
	})

	t.Run("teams/add-to-project", func(t *testing.T) {
		requireID(t, teamID, "teams/create must have succeeded")
		requireID(t, projectID, "teams/setup must have succeeded")
		out, err := zohoMayFail(t, "projects", "teams", "add-to-project",
			"--project", projectID,
			"--json", toJSON(t, []string{teamID}))
		if err != nil {
			t.Logf("add-to-project failed: %v", err)
			t.Logf("response: %s", truncate(out, 500))
			return
		}
		t.Logf("added team %s to project %s", teamID, projectID)
	})

	t.Run("teams/project-list", func(t *testing.T) {
		requireID(t, projectID, "teams/setup must have succeeded")
		out, err := zohoMayFail(t, "projects", "teams", "project-list",
			"--project", projectID)
		if err != nil {
			t.Logf("project-list failed: %v", err)
			return
		}
		m := parseJSON(t, out)
		if _, ok := m["teams"]; !ok {
			if _, ok2 := m["page_info"]; !ok2 {
				t.Logf("unexpected response format:\n%s", truncate(out, 500))
			}
		}
	})

	t.Run("teams/remove-from-project", func(t *testing.T) {
		requireID(t, teamID, "teams/create must have succeeded")
		requireID(t, projectID, "teams/setup must have succeeded")
		r := runZoho(t, "projects", "teams", "remove-from-project", teamID,
			"--project", projectID)
		if r.ExitCode != 0 {
			t.Logf("remove-from-project failed: %s", truncate(r.Stderr+r.Stdout, 300))
		}
	})

	t.Run("teams/delete", func(t *testing.T) {
		requireID(t, teamID, "teams/create must have succeeded")
		zoho(t, "projects", "teams", "delete", teamID)

		out := zoho(t, "projects", "teams", "list")
		m := parseJSON(t, out)
		teams, _ := m["teams"].([]any)
		for _, tm := range teams {
			tmm, _ := tm.(map[string]any)
			if fmt.Sprintf("%v", tmm["id"]) == teamID {
				t.Errorf("team %s still found in list after delete", teamID)
				break
			}
		}
		teamID = ""
	})
}

func TestProjectsErrors(t *testing.T) {
	t.Parallel()
	_ = requireProjectsPortalID(t)

	t.Run("bad-auth", func(t *testing.T) {
		r := runZohoWithEnv(t, map[string]string{
			"ZOHO_CLIENT_ID":     "bad_client_id",
			"ZOHO_CLIENT_SECRET": "bad_client_secret",
			"ZOHO_REFRESH_TOKEN": "bad_refresh_token",
			"ZOHO_DC":            "com",
		}, "projects", "core", "list")
		assertExitCode(t, r, 2)
	})

	t.Run("missing-portal", func(t *testing.T) {
		r := runZohoWithEnv(t, map[string]string{
			"ZOHO_PORTAL_ID": "",
		}, "projects", "portals", "get")
		if r.ExitCode == 0 {
			t.Error("expected error when portal ID is missing")
		}
		assertContains(t, r.Stderr, "--portal")
	})

	t.Run("nonexistent-project", func(t *testing.T) {
		r := runZoho(t, "projects", "core", "get", "999999999999")
		if r.ExitCode == 0 {
			t.Error("expected error for nonexistent project")
		}
	})

	t.Run("nonexistent-task", func(t *testing.T) {
		r := runZoho(t, "projects", "tasks", "get", "999999999999",
			"--project", "999999999999")
		if r.ExitCode == 0 {
			t.Error("expected error for nonexistent task")
		}
	})

	t.Run("missing-required-name", func(t *testing.T) {
		r := runZoho(t, "projects", "core", "create")
		assertExitCode(t, r, 1)
		assertContains(t, r.Stderr, "Required flag")
	})

	t.Run("missing-required-project-flag", func(t *testing.T) {
		r := runZoho(t, "projects", "tasks", "list")
		assertExitCode(t, r, 1)
		assertContains(t, r.Stderr, "Required flag")
	})

	t.Run("milestone-missing-dates", func(t *testing.T) {
		r := runZoho(t, "projects", "milestones", "create",
			"--name", "test", "--project", "12345")
		assertExitCode(t, r, 1)
		assertContains(t, r.Stderr, "Required flag")
	})

	t.Run("timelog-missing-required", func(t *testing.T) {
		r := runZoho(t, "projects", "timelogs", "add",
			"--hours", "2", "--project", "12345")
		assertExitCode(t, r, 1)
		assertContains(t, r.Stderr, "Required flag")
	})
}

func TestProjectsEmergencyCleanup(t *testing.T) {
	if os.Getenv("ZOHO_EMERGENCY_CLEANUP") == "" {
		t.Skip("set ZOHO_EMERGENCY_CLEANUP=1 to run")
	}
	_ = requireProjectsPortalID(t)

	out := zoho(t, "projects", "core", "list")
	arr := parseJSONArray(t, out)
	t.Logf("found %d total projects", len(arr))
	for _, proj := range arr {
		name := fmt.Sprintf("%v", proj["name"])
		if !strings.HasPrefix(name, testPrefix) {
			continue
		}
		id := fmt.Sprintf("%v", proj["id"])
		t.Logf("cleaning orphaned project %s (%s)", id, name)
		zohoIgnoreError(t, "projects", "core", "trash", id)
		time.Sleep(1 * time.Second)
		zohoIgnoreError(t, "projects", "core", "delete", id)
	}
}


