from __future__ import annotations

import json
from dataclasses import dataclass
from typing import Annotated, Any

import cappa

from zoho_cli.http.client import ZohoClient, get_client
from zoho_cli.output import output
from zoho_cli.pagination import paginate_projects


def _base(client: ZohoClient, portal_id: str, project_id: str) -> str:
    return f"{client.projects_base}/portal/{portal_id}/projects/{project_id}"


@cappa.command(name="list", help="List tasks in a project")
@dataclass
class TasksList:
    project: Annotated[str, cappa.Arg(long="--project", help="Project ID")]
    portal: Annotated[str, cappa.Arg(long="--portal", help="Portal ID")]
    status: Annotated[
        str | None,
        cappa.Arg(long="--status", default=None, help="Filter: open, closed, in progress"),
    ] = None
    priority: Annotated[
        str | None,
        cappa.Arg(long="--priority", default=None, help="Filter: none, low, medium, high"),
    ] = None

    def __call__(self, client: Annotated[ZohoClient, cappa.Dep(get_client)]) -> None:
        url = f"{_base(client, self.portal, self.project)}/tasks"
        params: dict[str, str] = {}
        if self.status:
            params["status"] = self.status
        if self.priority:
            params["priority"] = self.priority
        tasks = paginate_projects(client, url, "tasks", params=params)
        output(tasks)


@cappa.command(name="my", help="List my tasks across all projects")
@dataclass
class TasksMy:
    portal: Annotated[str, cappa.Arg(long="--portal", help="Portal ID")]
    status: Annotated[str | None, cappa.Arg(long="--status", default=None)] = None
    priority: Annotated[str | None, cappa.Arg(long="--priority", default=None)] = None

    def __call__(self, client: Annotated[ZohoClient, cappa.Dep(get_client)]) -> None:
        url = f"{client.projects_base}/portal/{self.portal}/tasks"
        params: dict[str, str] = {}
        if self.status:
            params["status"] = self.status
        if self.priority:
            params["priority"] = self.priority
        tasks = paginate_projects(client, url, "tasks", params=params)
        output(tasks)


@cappa.command(name="get", help="Get a single task")
@dataclass
class TasksGet:
    task_id: Annotated[str, cappa.Arg(help="Task ID")]
    project: Annotated[str, cappa.Arg(long="--project", help="Project ID")]
    portal: Annotated[str, cappa.Arg(long="--portal", help="Portal ID")]

    def __call__(self, client: Annotated[ZohoClient, cappa.Dep(get_client)]) -> None:
        url = f"{_base(client, self.portal, self.project)}/tasks/{self.task_id}"
        data = client.request("GET", url)
        output(data)


@cappa.command(name="create", help="Create a task")
@dataclass
class TasksCreate:
    name: Annotated[str, cappa.Arg(long="--name", help="Task name")]
    project: Annotated[str, cappa.Arg(long="--project", help="Project ID")]
    portal: Annotated[str, cappa.Arg(long="--portal", help="Portal ID")]
    json_data: Annotated[
        str | None, cappa.Arg(long="--json", default=None, help="Additional fields as JSON")
    ] = None

    def __call__(self, client: Annotated[ZohoClient, cappa.Dep(get_client)]) -> None:
        body: dict[str, Any] = {"name": self.name}
        if self.json_data:
            body.update(json.loads(self.json_data))
        url = f"{_base(client, self.portal, self.project)}/tasks"
        data = client.request("POST", url, json=body)
        output(data)


@cappa.command(name="update", help="Update a task")
@dataclass
class TasksUpdate:
    task_id: Annotated[str, cappa.Arg(help="Task ID")]
    project: Annotated[str, cappa.Arg(long="--project", help="Project ID")]
    portal: Annotated[str, cappa.Arg(long="--portal", help="Portal ID")]
    json_data: Annotated[str, cappa.Arg(long="--json", help="Fields to update as JSON")]

    def __call__(self, client: Annotated[ZohoClient, cappa.Dep(get_client)]) -> None:
        parsed = json.loads(self.json_data)
        url = f"{_base(client, self.portal, self.project)}/tasks/{self.task_id}"
        data = client.request("PATCH", url, json=parsed)
        output(data)


@cappa.command(name="delete", help="Delete a task")
@dataclass
class TasksDelete:
    task_id: Annotated[str, cappa.Arg(help="Task ID")]
    project: Annotated[str, cappa.Arg(long="--project", help="Project ID")]
    portal: Annotated[str, cappa.Arg(long="--portal", help="Portal ID")]

    def __call__(self, client: Annotated[ZohoClient, cappa.Dep(get_client)]) -> None:
        url = f"{_base(client, self.portal, self.project)}/tasks/{self.task_id}"
        data = client.request("DELETE", url)
        output(data)


@cappa.command(name="subtasks", help="List subtasks of a task")
@dataclass
class TasksSubtasksList:
    task_id: Annotated[str, cappa.Arg(help="Parent task ID")]
    project: Annotated[str, cappa.Arg(long="--project", help="Project ID")]
    portal: Annotated[str, cappa.Arg(long="--portal", help="Portal ID")]

    def __call__(self, client: Annotated[ZohoClient, cappa.Dep(get_client)]) -> None:
        url = f"{_base(client, self.portal, self.project)}/tasks/{self.task_id}/subtasks"
        data = client.request("GET", url)
        output(data)


@cappa.command(name="add-subtask", help="Create a subtask")
@dataclass
class TasksSubtaskCreate:
    name: Annotated[str, cappa.Arg(long="--name", help="Subtask name")]
    task_id: Annotated[str, cappa.Arg(long="--parent", help="Parent task ID")]
    project: Annotated[str, cappa.Arg(long="--project", help="Project ID")]
    portal: Annotated[str, cappa.Arg(long="--portal", help="Portal ID")]
    json_data: Annotated[
        str | None, cappa.Arg(long="--json", default=None, help="Additional fields as JSON")
    ] = None

    def __call__(self, client: Annotated[ZohoClient, cappa.Dep(get_client)]) -> None:
        body: dict[str, Any] = {"name": self.name}
        if self.json_data:
            body.update(json.loads(self.json_data))
        url = f"{_base(client, self.portal, self.project)}/tasks/{self.task_id}/subtasks"
        data = client.request("POST", url, json=body)
        output(data)


@cappa.command(name="tasks", help="Project task operations")
@dataclass
class Tasks:
    subcommand: cappa.Subcommands[
        TasksList
        | TasksMy
        | TasksGet
        | TasksCreate
        | TasksUpdate
        | TasksDelete
        | TasksSubtasksList
        | TasksSubtaskCreate
    ]


@cappa.command(name="list", help="List issues in a project")
@dataclass
class IssuesList:
    project: Annotated[str, cappa.Arg(long="--project", help="Project ID")]
    portal: Annotated[str, cappa.Arg(long="--portal", help="Portal ID")]

    def __call__(self, client: Annotated[ZohoClient, cappa.Dep(get_client)]) -> None:
        url = f"{_base(client, self.portal, self.project)}/issues"
        issues = paginate_projects(client, url, "issues")
        output(issues)


@cappa.command(name="create", help="Create an issue")
@dataclass
class IssuesCreate:
    name: Annotated[str, cappa.Arg(long="--name", help="Issue title")]
    project: Annotated[str, cappa.Arg(long="--project", help="Project ID")]
    portal: Annotated[str, cappa.Arg(long="--portal", help="Portal ID")]
    json_data: Annotated[
        str | None, cappa.Arg(long="--json", default=None, help="Additional fields as JSON")
    ] = None

    def __call__(self, client: Annotated[ZohoClient, cappa.Dep(get_client)]) -> None:
        body: dict[str, Any] = {"name": self.name}
        if self.json_data:
            body.update(json.loads(self.json_data))
        url = f"{_base(client, self.portal, self.project)}/issues"
        data = client.request("POST", url, json=body)
        output(data)


@cappa.command(name="update", help="Update an issue")
@dataclass
class IssuesUpdate:
    issue_id: Annotated[str, cappa.Arg(help="Issue ID")]
    project: Annotated[str, cappa.Arg(long="--project", help="Project ID")]
    portal: Annotated[str, cappa.Arg(long="--portal", help="Portal ID")]
    json_data: Annotated[str, cappa.Arg(long="--json", help="Fields to update as JSON")]

    def __call__(self, client: Annotated[ZohoClient, cappa.Dep(get_client)]) -> None:
        parsed = json.loads(self.json_data)
        url = f"{_base(client, self.portal, self.project)}/issues/{self.issue_id}"
        data = client.request("PATCH", url, json=parsed)
        output(data)


@cappa.command(name="get", help="Get a single issue")
@dataclass
class IssuesGet:
    issue_id: Annotated[str, cappa.Arg(help="Issue ID")]
    project: Annotated[str, cappa.Arg(long="--project", help="Project ID")]
    portal: Annotated[str, cappa.Arg(long="--portal", help="Portal ID")]

    def __call__(self, client: Annotated[ZohoClient, cappa.Dep(get_client)]) -> None:
        url = f"{_base(client, self.portal, self.project)}/issues/{self.issue_id}"
        data = client.request("GET", url)
        output(data)


@cappa.command(name="delete", help="Delete an issue")
@dataclass
class IssuesDelete:
    issue_id: Annotated[str, cappa.Arg(help="Issue ID")]
    project: Annotated[str, cappa.Arg(long="--project", help="Project ID")]
    portal: Annotated[str, cappa.Arg(long="--portal", help="Portal ID")]

    def __call__(self, client: Annotated[ZohoClient, cappa.Dep(get_client)]) -> None:
        url = f"{_base(client, self.portal, self.project)}/issues/{self.issue_id}"
        data = client.request("DELETE", url)
        output(data)


@cappa.command(name="defaults", help="Get issue default fields (statuses, severities, etc.)")
@dataclass
class IssuesDefaults:
    project: Annotated[str, cappa.Arg(long="--project", help="Project ID")]
    portal: Annotated[str, cappa.Arg(long="--portal", help="Portal ID")]

    def __call__(self, client: Annotated[ZohoClient, cappa.Dep(get_client)]) -> None:
        url = f"{_base(client, self.portal, self.project)}/issues/defaultfields"
        data = client.request("GET", url)
        output(data)


@cappa.command(name="issues", help="Project issue operations")
@dataclass
class Issues:
    subcommand: cappa.Subcommands[
        IssuesList | IssuesCreate | IssuesUpdate | IssuesGet | IssuesDelete | IssuesDefaults
    ]


@cappa.command(name="list", help="List issue comments")
@dataclass
class IssueCommentsList:
    issue_id: Annotated[str, cappa.Arg(long="--issue", help="Issue ID")]
    project: Annotated[str, cappa.Arg(long="--project", help="Project ID")]
    portal: Annotated[str, cappa.Arg(long="--portal", help="Portal ID")]

    def __call__(self, client: Annotated[ZohoClient, cappa.Dep(get_client)]) -> None:
        url = f"{_base(client, self.portal, self.project)}/issues/{self.issue_id}/comments"
        data = client.request("GET", url)
        output(data)


@cappa.command(name="add", help="Add an issue comment")
@dataclass
class IssueCommentsAdd:
    comment: Annotated[str, cappa.Arg(long="--comment", help="Comment text")]
    issue_id: Annotated[str, cappa.Arg(long="--issue", help="Issue ID")]
    project: Annotated[str, cappa.Arg(long="--project", help="Project ID")]
    portal: Annotated[str, cappa.Arg(long="--portal", help="Portal ID")]

    def __call__(self, client: Annotated[ZohoClient, cappa.Dep(get_client)]) -> None:
        url = f"{_base(client, self.portal, self.project)}/issues/{self.issue_id}/comments"
        data = client.request("POST", url, json={"comment": self.comment})
        output(data)


@cappa.command(name="issue-comments", help="Issue comment operations")
@dataclass
class IssueComments:
    subcommand: cappa.Subcommands[IssueCommentsList | IssueCommentsAdd]


@cappa.command(name="list", help="List task comments")
@dataclass
class CommentsList:
    task_id: Annotated[str, cappa.Arg(long="--task", help="Task ID")]
    project: Annotated[str, cappa.Arg(long="--project", help="Project ID")]
    portal: Annotated[str, cappa.Arg(long="--portal", help="Portal ID")]

    def __call__(self, client: Annotated[ZohoClient, cappa.Dep(get_client)]) -> None:
        url = f"{_base(client, self.portal, self.project)}/tasks/{self.task_id}/comments"
        data = client.request("GET", url)
        output(data)


@cappa.command(name="add", help="Add a task comment")
@dataclass
class CommentsAdd:
    comment: Annotated[str, cappa.Arg(long="--comment", help="Comment text")]
    task_id: Annotated[str, cappa.Arg(long="--task", help="Task ID")]
    project: Annotated[str, cappa.Arg(long="--project", help="Project ID")]
    portal: Annotated[str, cappa.Arg(long="--portal", help="Portal ID")]

    def __call__(self, client: Annotated[ZohoClient, cappa.Dep(get_client)]) -> None:
        url = f"{_base(client, self.portal, self.project)}/tasks/{self.task_id}/comments"
        data = client.request("POST", url, json={"comment": self.comment})
        output(data)


@cappa.command(name="update", help="Update a task comment")
@dataclass
class CommentsUpdate:
    comment_id: Annotated[str, cappa.Arg(help="Comment ID")]
    comment: Annotated[str, cappa.Arg(long="--comment", help="Updated comment text")]
    task_id: Annotated[str, cappa.Arg(long="--task", help="Task ID")]
    project: Annotated[str, cappa.Arg(long="--project", help="Project ID")]
    portal: Annotated[str, cappa.Arg(long="--portal", help="Portal ID")]

    def __call__(self, client: Annotated[ZohoClient, cappa.Dep(get_client)]) -> None:
        url = (
            f"{_base(client, self.portal, self.project)}"
            f"/tasks/{self.task_id}/comments/{self.comment_id}"
        )
        data = client.request("PATCH", url, json={"comment": self.comment})
        output(data)


@cappa.command(name="delete", help="Delete a task comment")
@dataclass
class CommentsDelete:
    comment_id: Annotated[str, cappa.Arg(help="Comment ID")]
    task_id: Annotated[str, cappa.Arg(long="--task", help="Task ID")]
    project: Annotated[str, cappa.Arg(long="--project", help="Project ID")]
    portal: Annotated[str, cappa.Arg(long="--portal", help="Portal ID")]

    def __call__(self, client: Annotated[ZohoClient, cappa.Dep(get_client)]) -> None:
        url = (
            f"{_base(client, self.portal, self.project)}"
            f"/tasks/{self.task_id}/comments/{self.comment_id}"
        )
        data = client.request("DELETE", url)
        output(data)


@cappa.command(name="comments", help="Task comment operations")
@dataclass
class Comments:
    subcommand: cappa.Subcommands[CommentsList | CommentsAdd | CommentsUpdate | CommentsDelete]


@cappa.command(name="list", help="List tasklists")
@dataclass
class TasklistsList:
    project: Annotated[str, cappa.Arg(long="--project", help="Project ID")]
    portal: Annotated[str, cappa.Arg(long="--portal", help="Portal ID")]

    def __call__(self, client: Annotated[ZohoClient, cappa.Dep(get_client)]) -> None:
        url = f"{_base(client, self.portal, self.project)}/tasklists"
        tasklists = paginate_projects(client, url, "tasklists")
        output(tasklists)


@cappa.command(name="create", help="Create a tasklist")
@dataclass
class TasklistsCreate:
    name: Annotated[str, cappa.Arg(long="--name", help="Tasklist name")]
    project: Annotated[str, cappa.Arg(long="--project", help="Project ID")]
    portal: Annotated[str, cappa.Arg(long="--portal", help="Portal ID")]
    json_data: Annotated[
        str | None, cappa.Arg(long="--json", default=None, help="Additional fields as JSON")
    ] = None

    def __call__(self, client: Annotated[ZohoClient, cappa.Dep(get_client)]) -> None:
        body: dict[str, Any] = {"name": self.name}
        if self.json_data:
            body.update(json.loads(self.json_data))
        url = f"{_base(client, self.portal, self.project)}/tasklists"
        data = client.request("POST", url, json=body)
        output(data)


@cappa.command(name="update", help="Update a tasklist")
@dataclass
class TasklistsUpdate:
    tasklist_id: Annotated[str, cappa.Arg(help="Tasklist ID")]
    project: Annotated[str, cappa.Arg(long="--project", help="Project ID")]
    portal: Annotated[str, cappa.Arg(long="--portal", help="Portal ID")]
    json_data: Annotated[str, cappa.Arg(long="--json", help="Fields to update as JSON")]

    def __call__(self, client: Annotated[ZohoClient, cappa.Dep(get_client)]) -> None:
        parsed = json.loads(self.json_data)
        url = f"{_base(client, self.portal, self.project)}/tasklists/{self.tasklist_id}"
        data = client.request("PATCH", url, json=parsed)
        output(data)


@cappa.command(name="delete", help="Delete a tasklist")
@dataclass
class TasklistsDelete:
    tasklist_id: Annotated[str, cappa.Arg(help="Tasklist ID")]
    project: Annotated[str, cappa.Arg(long="--project", help="Project ID")]
    portal: Annotated[str, cappa.Arg(long="--portal", help="Portal ID")]

    def __call__(self, client: Annotated[ZohoClient, cappa.Dep(get_client)]) -> None:
        url = f"{_base(client, self.portal, self.project)}/tasklists/{self.tasklist_id}"
        data = client.request("DELETE", url)
        output(data)


@cappa.command(name="tasklists", help="Project tasklist operations")
@dataclass
class Tasklists:
    subcommand: cappa.Subcommands[
        TasklistsList | TasklistsCreate | TasklistsUpdate | TasklistsDelete
    ]


@cappa.command(name="list", help="List project timelogs")
@dataclass
class TimelogsList:
    project: Annotated[str, cappa.Arg(long="--project", help="Project ID")]
    portal: Annotated[str, cappa.Arg(long="--portal", help="Portal ID")]
    module_type: Annotated[
        str, cappa.Arg(long="--module", default="general", help="task, issue, or general")
    ] = "general"

    def __call__(self, client: Annotated[ZohoClient, cappa.Dep(get_client)]) -> None:
        url = f"{_base(client, self.portal, self.project)}/timelogs"
        params: dict[str, str] = {
            "module": json.dumps({"type": self.module_type}),
            "view_type": "projectspan",
        }
        data = client.request("GET", url, params=params)
        output(data)


@cappa.command(name="add", help="Add a timelog")
@dataclass
class TimelogsAdd:
    date: Annotated[str, cappa.Arg(long="--date", help="Date (YYYY-MM-DD)")]
    hours: Annotated[str, cappa.Arg(long="--hours", help="Hours (e.g. 2, 1.5, 0:30)")]
    project: Annotated[str, cappa.Arg(long="--project", help="Project ID")]
    portal: Annotated[str, cappa.Arg(long="--portal", help="Portal ID")]
    task_id: Annotated[str | None, cappa.Arg(long="--task", default=None, help="Task ID")] = None
    bill_status: Annotated[
        str, cappa.Arg(long="--bill-status", default="Billable", help="Billable or Non Billable")
    ] = "Billable"
    notes: Annotated[
        str | None, cappa.Arg(long="--notes", default=None, help="Notes for time entry")
    ] = None

    def __call__(self, client: Annotated[ZohoClient, cappa.Dep(get_client)]) -> None:
        body: dict[str, Any] = {
            "date": self.date,
            "hours": self.hours,
            "bill_status": self.bill_status,
            "log_name": self.notes or "Time log",
        }
        if self.task_id:
            body["module"] = {"type": "task", "id": self.task_id}
        else:
            body["module"] = {"type": "general"}
        if self.notes:
            body["notes"] = self.notes
        url = f"{_base(client, self.portal, self.project)}/log"
        data = client.request("POST", url, json=body)
        output(data)


@cappa.command(name="timelogs", help="Project timelog operations")
@dataclass
class Timelogs:
    subcommand: cappa.Subcommands[TimelogsList | TimelogsAdd]


@cappa.command(name="list", help="List project users")
@dataclass
class ProjectUsersList:
    project: Annotated[str, cappa.Arg(long="--project", help="Project ID")]
    portal: Annotated[str, cappa.Arg(long="--portal", help="Portal ID")]

    def __call__(self, client: Annotated[ZohoClient, cappa.Dep(get_client)]) -> None:
        url = f"{_base(client, self.portal, self.project)}/users"
        users = paginate_projects(client, url, "users")
        output(users)


@cappa.command(name="users", help="Project user operations")
@dataclass
class ProjectUsers:
    subcommand: cappa.Subcommands[ProjectUsersList]


@cappa.command(name="list", help="List all projects")
@dataclass
class ProjectsList:
    portal: Annotated[str, cappa.Arg(long="--portal", help="Portal ID")]

    def __call__(self, client: Annotated[ZohoClient, cappa.Dep(get_client)]) -> None:
        url = f"{client.projects_base}/portal/{self.portal}/projects"
        projects = paginate_projects(client, url, None)
        output(projects)


@cappa.command(name="get", help="Get a single project")
@dataclass
class ProjectsGet:
    project_id: Annotated[str, cappa.Arg(help="Project ID")]
    portal: Annotated[str, cappa.Arg(long="--portal", help="Portal ID")]

    def __call__(self, client: Annotated[ZohoClient, cappa.Dep(get_client)]) -> None:
        url = f"{client.projects_base}/portal/{self.portal}/projects/{self.project_id}"
        data = client.request("GET", url)
        output(data)


@cappa.command(name="search", help="Search projects")
@dataclass
class ProjectsSearch:
    query: Annotated[str, cappa.Arg(long="--query", help="Search query")]
    portal: Annotated[str, cappa.Arg(long="--portal", help="Portal ID")]

    def __call__(self, client: Annotated[ZohoClient, cappa.Dep(get_client)]) -> None:
        url = f"{client.projects_base}/portal/{self.portal}/search"
        params = {"search_term": self.query, "module": "all", "status": "all"}
        data = client.request("GET", url, params=params)
        output(data)


@cappa.command(name="create", help="Create a project")
@dataclass
class ProjectsCreate:
    name: Annotated[str, cappa.Arg(long="--name", help="Project name")]
    portal: Annotated[str, cappa.Arg(long="--portal", help="Portal ID")]
    json_data: Annotated[
        str | None, cappa.Arg(long="--json", default=None, help="Additional fields as JSON")
    ] = None

    def __call__(self, client: Annotated[ZohoClient, cappa.Dep(get_client)]) -> None:
        body: dict[str, Any] = {"name": self.name}
        if self.json_data:
            body.update(json.loads(self.json_data))
        url = f"{client.projects_base}/portal/{self.portal}/projects"
        data = client.request("POST", url, json=body)
        output(data)


@cappa.command(name="update", help="Update a project")
@dataclass
class ProjectsUpdate:
    project_id: Annotated[str, cappa.Arg(help="Project ID")]
    portal: Annotated[str, cappa.Arg(long="--portal", help="Portal ID")]
    json_data: Annotated[str, cappa.Arg(long="--json", help="Fields to update as JSON")]

    def __call__(self, client: Annotated[ZohoClient, cappa.Dep(get_client)]) -> None:
        parsed = json.loads(self.json_data)
        url = f"{client.projects_base}/portal/{self.portal}/projects/{self.project_id}"
        data = client.request("PATCH", url, json=parsed)
        output(data)


@cappa.command(name="list", help="List milestones")
@dataclass
class MilestonesList:
    project: Annotated[str, cappa.Arg(long="--project", help="Project ID")]
    portal: Annotated[str, cappa.Arg(long="--portal", help="Portal ID")]

    def __call__(self, client: Annotated[ZohoClient, cappa.Dep(get_client)]) -> None:
        url = f"{_base(client, self.portal, self.project)}/milestones"
        data = paginate_projects(client, url, "milestones")
        output(data)


@cappa.command(name="get", help="Get a milestone")
@dataclass
class MilestonesGet:
    milestone_id: Annotated[str, cappa.Arg(help="Milestone ID")]
    project: Annotated[str, cappa.Arg(long="--project", help="Project ID")]
    portal: Annotated[str, cappa.Arg(long="--portal", help="Portal ID")]

    def __call__(self, client: Annotated[ZohoClient, cappa.Dep(get_client)]) -> None:
        url = f"{_base(client, self.portal, self.project)}/milestones/{self.milestone_id}"
        data = client.request("GET", url)
        output(data)


@cappa.command(name="create", help="Create a milestone")
@dataclass
class MilestonesCreate:
    name: Annotated[str, cappa.Arg(long="--name", help="Milestone name")]
    start_date: Annotated[str, cappa.Arg(long="--start", help="Start date (YYYY-MM-DD)")]
    end_date: Annotated[str, cappa.Arg(long="--end", help="End date (YYYY-MM-DD)")]
    project: Annotated[str, cappa.Arg(long="--project", help="Project ID")]
    portal: Annotated[str, cappa.Arg(long="--portal", help="Portal ID")]
    json_data: Annotated[
        str | None, cappa.Arg(long="--json", default=None, help="Additional fields as JSON")
    ] = None

    def __call__(self, client: Annotated[ZohoClient, cappa.Dep(get_client)]) -> None:
        body: dict[str, Any] = {
            "name": self.name,
            "start_date": self.start_date,
            "end_date": self.end_date,
        }
        if self.json_data:
            body.update(json.loads(self.json_data))
        url = f"{_base(client, self.portal, self.project)}/milestones"
        data = client.request("POST", url, json=body)
        output(data)


@cappa.command(name="update", help="Update a milestone")
@dataclass
class MilestonesUpdate:
    milestone_id: Annotated[str, cappa.Arg(help="Milestone ID")]
    project: Annotated[str, cappa.Arg(long="--project", help="Project ID")]
    portal: Annotated[str, cappa.Arg(long="--portal", help="Portal ID")]
    json_data: Annotated[str, cappa.Arg(long="--json", help="Fields to update as JSON")]

    def __call__(self, client: Annotated[ZohoClient, cappa.Dep(get_client)]) -> None:
        parsed = json.loads(self.json_data)
        url = f"{_base(client, self.portal, self.project)}/milestones/{self.milestone_id}"
        data = client.request("PATCH", url, json=parsed)
        output(data)


@cappa.command(name="delete", help="Delete a milestone")
@dataclass
class MilestonesDelete:
    milestone_id: Annotated[str, cappa.Arg(help="Milestone ID")]
    project: Annotated[str, cappa.Arg(long="--project", help="Project ID")]
    portal: Annotated[str, cappa.Arg(long="--portal", help="Portal ID")]

    def __call__(self, client: Annotated[ZohoClient, cappa.Dep(get_client)]) -> None:
        url = f"{_base(client, self.portal, self.project)}/milestones/{self.milestone_id}"
        data = client.request("DELETE", url)
        output(data)


@cappa.command(name="milestones", help="Project milestone operations")
@dataclass
class Milestones:
    subcommand: cappa.Subcommands[
        MilestonesList | MilestonesGet | MilestonesCreate | MilestonesUpdate | MilestonesDelete
    ]


@cappa.command(name="add", help="Add a task dependency")
@dataclass
class DependenciesAdd:
    task_id: Annotated[str, cappa.Arg(help="Task ID")]
    dependency_id: Annotated[str, cappa.Arg(long="--depends-on", help="Dependency task ID")]
    project: Annotated[str, cappa.Arg(long="--project", help="Project ID")]
    portal: Annotated[str, cappa.Arg(long="--portal", help="Portal ID")]
    dep_type: Annotated[str, cappa.Arg(long="--type", default="FS", help="FS, SS, FF, or SF")] = (
        "FS"
    )

    def __call__(self, client: Annotated[ZohoClient, cappa.Dep(get_client)]) -> None:
        url = f"{_base(client, self.portal, self.project)}/tasks/{self.task_id}/dependencies"
        body: dict[str, Any] = {"predecessor": {"id": self.dependency_id, "type": self.dep_type}}
        data = client.request("POST", url, json=body)
        output(data)


@cappa.command(name="remove", help="Remove a task dependency")
@dataclass
class DependenciesRemove:
    task_id: Annotated[str, cappa.Arg(help="Task ID")]
    dependency_id: Annotated[str, cappa.Arg(help="Dependency ID to remove")]
    project: Annotated[str, cappa.Arg(long="--project", help="Project ID")]
    portal: Annotated[str, cappa.Arg(long="--portal", help="Portal ID")]

    def __call__(self, client: Annotated[ZohoClient, cappa.Dep(get_client)]) -> None:
        url = (
            f"{_base(client, self.portal, self.project)}"
            f"/tasks/{self.task_id}/dependencies/{self.dependency_id}"
        )
        data = client.request("DELETE", url)
        output(data)


@cappa.command(name="dependencies", help="Task dependency operations")
@dataclass
class Dependencies:
    subcommand: cappa.Subcommands[DependenciesAdd | DependenciesRemove]


@cappa.command(name="projects", help="Zoho Projects operations")
@dataclass
class Projects:
    subcommand: cappa.Subcommands[
        ProjectsList
        | ProjectsGet
        | ProjectsSearch
        | ProjectsCreate
        | ProjectsUpdate
        | Tasks
        | Issues
        | IssueComments
        | Comments
        | Tasklists
        | Timelogs
        | ProjectUsers
        | Milestones
        | Dependencies
    ]
