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


def _slim_project(p: dict[str, Any]) -> dict[str, Any]:
    task_count = p.get("task_count")
    return {
        "id": p.get("id_string", p.get("id")),
        "name": p.get("name"),
        "status": p.get("status"),
        "owner": p.get("owner_name")
        or (p.get("owner", {}).get("name") if isinstance(p.get("owner"), dict) else None),
        "task_count": task_count.get("open") if isinstance(task_count, dict) else None,
        "start_date": p.get("start_date"),
        "end_date": p.get("end_date"),
    }


def _slim_task(task: dict[str, Any]) -> dict[str, Any]:
    owners = task.get("owners_and_work", {}).get("owners", [])
    if not owners:
        details = task.get("details", {})
        owners = details.get("owners", []) if isinstance(details, dict) else []
    owner_names = [o.get("name", "") for o in owners if isinstance(o, dict)]
    status = task.get("status")
    status_name = status.get("name") if isinstance(status, dict) else status
    tasklist = task.get("tasklist")
    tasklist_name = tasklist.get("name") if isinstance(tasklist, dict) else tasklist
    project = task.get("project")
    project_name = project.get("name") if isinstance(project, dict) else None
    project_id = project.get("id") if isinstance(project, dict) else None
    return {
        "id": task.get("id_string", task.get("id")),
        "prefix": task.get("prefix"),
        "name": task.get("name"),
        "project": project_name,
        "project_id": project_id,
        "status": status_name,
        "is_completed": task.get("is_completed"),
        "priority": task.get("priority"),
        "owners": owner_names,
        "tasklist": tasklist_name,
        "start_date": task.get("start_date"),
        "end_date": task.get("end_date"),
        "completion_percentage": task.get("completion_percentage"),
        "logged_hours": task.get("log_hours", {}).get("total_hours"),
    }


def _slim_comment(c: dict[str, Any]) -> dict[str, Any]:
    author = c.get("added_by") or c.get("posted_by")
    return {
        "id": c.get("id"),
        "content": c.get("content"),
        "author": author.get("name") if isinstance(author, dict) else author,
        "created_time": c.get("created_time"),
    }


def _slim_tasklist(tl: dict[str, Any]) -> dict[str, Any]:
    return {
        "id": tl.get("id_string", tl.get("id")),
        "name": tl.get("name"),
        "completed": tl.get("completed"),
        "open_tasks": tl.get("open_task_count"),
        "closed_tasks": tl.get("closed_task_count"),
    }


def _slim_timelog(log: dict[str, Any]) -> dict[str, Any]:
    owner = log.get("owner") or log.get("added_by")
    task = log.get("task")
    return {
        "id": log.get("id"),
        "date": log.get("date") or log.get("log_date"),
        "hours": log.get("hours") or log.get("total_hours"),
        "bill_status": log.get("bill_status"),
        "notes": log.get("notes"),
        "owner": owner.get("name") if isinstance(owner, dict) else owner,
        "task_name": task.get("name") if isinstance(task, dict) else None,
        "task_id": task.get("id") if isinstance(task, dict) else None,
    }


def _slim_issue(issue: dict[str, Any]) -> dict[str, Any]:
    reporter = issue.get("reported_person") or issue.get("reporter")
    assignee = issue.get("assignee")
    module = issue.get("module")
    return {
        "id": issue.get("id_string", issue.get("id")),
        "title": issue.get("title"),
        "status": issue.get("status", {}).get("type")
        if isinstance(issue.get("status"), dict)
        else issue.get("status"),
        "severity": issue.get("severity", {}).get("type")
        if isinstance(issue.get("severity"), dict)
        else issue.get("severity"),
        "reporter": reporter.get("name") if isinstance(reporter, dict) else reporter,
        "assignee": assignee.get("name") if isinstance(assignee, dict) else assignee,
        "module": module.get("name") if isinstance(module, dict) else module,
        "created_time": issue.get("created_time"),
    }


def _slim_user(u: dict[str, Any]) -> dict[str, Any]:
    return {
        "id": u.get("id"),
        "name": u.get("name"),
        "email": u.get("email"),
        "role": u.get("role"),
    }


def _slim_search_result(r: dict[str, Any]) -> dict[str, Any]:
    project = r.get("project")
    return {
        "id": r.get("id"),
        "title": r.get("title") or r.get("name"),
        "type": r.get("entity_type") or r.get("type"),
        "project": project.get("name") if isinstance(project, dict) else None,
        "project_id": project.get("id") if isinstance(project, dict) else None,
    }


def _flatten_timelogs(raw: Any) -> list[dict[str, Any]]:
    if isinstance(raw, list):
        result: list[dict[str, Any]] = []
        for group in raw:
            if isinstance(group, dict):
                details = group.get("log_details", [])
                if isinstance(details, list):
                    for log in details:
                        if isinstance(log, dict):
                            if "date" not in log and "log_date" not in log:
                                log["date"] = group.get("date")
                            result.append(log)
        return result
    return []


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
        tasks = paginate_projects(client, url, "tasks")
        if self.status:
            sl = self.status.lower()
            tasks = [t for t in tasks if str(_slim_task(t).get("status", "")).lower() == sl]
        if self.priority:
            pl = self.priority.lower()
            tasks = [t for t in tasks if str(t.get("priority", "")).lower() == pl]
        output([_slim_task(t) for t in tasks])


@cappa.command(name="my", help="List my tasks across all projects")
@dataclass
class TasksMy:
    portal: Annotated[str, cappa.Arg(long="--portal", help="Portal ID")]
    status: Annotated[str | None, cappa.Arg(long="--status", default=None)] = None
    priority: Annotated[str | None, cappa.Arg(long="--priority", default=None)] = None

    def __call__(self, client: Annotated[ZohoClient, cappa.Dep(get_client)]) -> None:
        url = f"{client.projects_base}/portal/{self.portal}/tasks"
        tasks = paginate_projects(client, url, "tasks")
        if self.status:
            sl = self.status.lower()
            tasks = [t for t in tasks if str(_slim_task(t).get("status", "")).lower() == sl]
        if self.priority:
            pl = self.priority.lower()
            tasks = [t for t in tasks if str(t.get("priority", "")).lower() == pl]
        output([_slim_task(t) for t in tasks])


@cappa.command(name="get", help="Get a single task")
@dataclass
class TasksGet:
    task_id: Annotated[str, cappa.Arg(help="Task ID")]
    project: Annotated[str, cappa.Arg(long="--project", help="Project ID")]
    portal: Annotated[str, cappa.Arg(long="--portal", help="Portal ID")]

    def __call__(self, client: Annotated[ZohoClient, cappa.Dep(get_client)]) -> None:
        url = f"{_base(client, self.portal, self.project)}/tasks/{self.task_id}"
        data = client.request("GET", url)
        output(_slim_task(data))


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
        output({"ok": True, "task": _slim_task(data)})


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
        output({"ok": True, "task": _slim_task(data)})


@cappa.command(name="tasks", help="Project task operations")
@dataclass
class Tasks:
    subcommand: cappa.Subcommands[TasksList | TasksMy | TasksGet | TasksCreate | TasksUpdate]


@cappa.command(name="list", help="List issues in a project")
@dataclass
class IssuesList:
    project: Annotated[str, cappa.Arg(long="--project", help="Project ID")]
    portal: Annotated[str, cappa.Arg(long="--portal", help="Portal ID")]

    def __call__(self, client: Annotated[ZohoClient, cappa.Dep(get_client)]) -> None:
        url = f"{_base(client, self.portal, self.project)}/issues"
        issues = paginate_projects(client, url, "issues")
        output([_slim_issue(i) for i in issues])


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
        output({"ok": True, "issue": _slim_issue(data)})


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
        output({"ok": True, "issue": _slim_issue(data)})


@cappa.command(name="issues", help="Project issue operations")
@dataclass
class Issues:
    subcommand: cappa.Subcommands[IssuesList | IssuesCreate | IssuesUpdate]


@cappa.command(name="list", help="List task comments")
@dataclass
class CommentsList:
    task_id: Annotated[str, cappa.Arg(long="--task", help="Task ID")]
    project: Annotated[str, cappa.Arg(long="--project", help="Project ID")]
    portal: Annotated[str, cappa.Arg(long="--portal", help="Portal ID")]

    def __call__(self, client: Annotated[ZohoClient, cappa.Dep(get_client)]) -> None:
        url = f"{_base(client, self.portal, self.project)}/tasks/{self.task_id}/comments"
        raw: Any = client.request("GET", url)
        comments = raw if isinstance(raw, list) else raw.get("comments", [])
        output([_slim_comment(c) for c in comments])


@cappa.command(name="add", help="Add a task comment")
@dataclass
class CommentsAdd:
    comment: Annotated[str, cappa.Arg(long="--comment", help="Comment text")]
    task_id: Annotated[str, cappa.Arg(long="--task", help="Task ID")]
    project: Annotated[str, cappa.Arg(long="--project", help="Project ID")]
    portal: Annotated[str, cappa.Arg(long="--portal", help="Portal ID")]

    def __call__(self, client: Annotated[ZohoClient, cappa.Dep(get_client)]) -> None:
        url = f"{_base(client, self.portal, self.project)}/tasks/{self.task_id}/comments"
        raw: Any = client.request("POST", url, json={"comment": self.comment})
        if isinstance(raw, list) and raw:
            comment = raw[0]
        elif isinstance(raw, dict):
            comment = raw
        else:
            comment = {}
        output({"ok": True, "comment": _slim_comment(comment)})


@cappa.command(name="comments", help="Task comment operations")
@dataclass
class Comments:
    subcommand: cappa.Subcommands[CommentsList | CommentsAdd]


@cappa.command(name="list", help="List tasklists")
@dataclass
class TasklistsList:
    project: Annotated[str, cappa.Arg(long="--project", help="Project ID")]
    portal: Annotated[str, cappa.Arg(long="--portal", help="Portal ID")]

    def __call__(self, client: Annotated[ZohoClient, cappa.Dep(get_client)]) -> None:
        url = f"{_base(client, self.portal, self.project)}/tasklists"
        tasklists = paginate_projects(client, url, "tasklists")
        output([_slim_tasklist(tl) for tl in tasklists])


@cappa.command(name="tasklists", help="Project tasklist operations")
@dataclass
class Tasklists:
    subcommand: cappa.Subcommands[TasklistsList]


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
        raw: Any = client.request("GET", url, params=params)
        if isinstance(raw, dict):
            raw_logs = raw.get("time_logs", raw.get("timelogs", []))
        else:
            raw_logs = raw if isinstance(raw, list) else []
        logs = _flatten_timelogs(raw_logs)
        if not logs and isinstance(raw_logs, list):
            logs = [tl for tl in raw_logs if isinstance(tl, dict) and "id" in tl]
        output([_slim_timelog(log) for log in logs])


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
        output({"ok": True, "timelog": _slim_timelog(data)})


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
        output([_slim_user(u) for u in users])


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
        output([_slim_project(p) for p in projects])


@cappa.command(name="get", help="Get a single project")
@dataclass
class ProjectsGet:
    project_id: Annotated[str, cappa.Arg(help="Project ID")]
    portal: Annotated[str, cappa.Arg(long="--portal", help="Portal ID")]

    def __call__(self, client: Annotated[ZohoClient, cappa.Dep(get_client)]) -> None:
        url = f"{client.projects_base}/portal/{self.portal}/projects/{self.project_id}"
        data = client.request("GET", url)
        output(_slim_project(data))


@cappa.command(name="search", help="Search projects")
@dataclass
class ProjectsSearch:
    query: Annotated[str, cappa.Arg(long="--query", help="Search query")]
    portal: Annotated[str, cappa.Arg(long="--portal", help="Portal ID")]

    def __call__(self, client: Annotated[ZohoClient, cappa.Dep(get_client)]) -> None:
        url = f"{client.projects_base}/portal/{self.portal}/search"
        params = {"search_term": self.query, "module": "all", "status": "all"}
        data = client.request("GET", url, params=params)
        results = data.get("results", []) if isinstance(data, dict) else []
        output([_slim_search_result(r) for r in results])


@cappa.command(name="projects", help="Zoho Projects operations")
@dataclass
class Projects:
    subcommand: cappa.Subcommands[
        ProjectsList
        | ProjectsGet
        | ProjectsSearch
        | Tasks
        | Issues
        | Comments
        | Tasklists
        | Timelogs
        | ProjectUsers
    ]
