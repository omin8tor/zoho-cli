from __future__ import annotations

from dataclasses import dataclass

import cappa


@cappa.command(name="tasks", help="Project task operations")
@dataclass
class Tasks:
    subcommand: cappa.Subcommands[TasksList | TasksMy | TasksGet | TasksCreate | TasksUpdate]


@cappa.command(name="list", help="List tasks in a project")
@dataclass
class TasksList:
    def __call__(self) -> None:
        raise NotImplementedError


@cappa.command(name="my", help="List my tasks across all projects")
@dataclass
class TasksMy:
    def __call__(self) -> None:
        raise NotImplementedError


@cappa.command(name="get", help="Get a single task")
@dataclass
class TasksGet:
    def __call__(self) -> None:
        raise NotImplementedError


@cappa.command(name="create", help="Create a task")
@dataclass
class TasksCreate:
    def __call__(self) -> None:
        raise NotImplementedError


@cappa.command(name="update", help="Update a task")
@dataclass
class TasksUpdate:
    def __call__(self) -> None:
        raise NotImplementedError


@cappa.command(name="issues", help="Project issue operations")
@dataclass
class Issues:
    subcommand: cappa.Subcommands[IssuesList | IssuesCreate | IssuesUpdate]


@cappa.command(name="list", help="List issues in a project")
@dataclass
class IssuesList:
    def __call__(self) -> None:
        raise NotImplementedError


@cappa.command(name="create", help="Create an issue")
@dataclass
class IssuesCreate:
    def __call__(self) -> None:
        raise NotImplementedError


@cappa.command(name="update", help="Update an issue")
@dataclass
class IssuesUpdate:
    def __call__(self) -> None:
        raise NotImplementedError


@cappa.command(name="comments", help="Task comment operations")
@dataclass
class Comments:
    subcommand: cappa.Subcommands[CommentsList | CommentsAdd]


@cappa.command(name="list", help="List task comments")
@dataclass
class CommentsList:
    def __call__(self) -> None:
        raise NotImplementedError


@cappa.command(name="add", help="Add a task comment")
@dataclass
class CommentsAdd:
    def __call__(self) -> None:
        raise NotImplementedError


@cappa.command(name="tasklists", help="Project tasklist operations")
@dataclass
class Tasklists:
    subcommand: cappa.Subcommands[TasklistsList]


@cappa.command(name="list", help="List tasklists")
@dataclass
class TasklistsList:
    def __call__(self) -> None:
        raise NotImplementedError


@cappa.command(name="timelogs", help="Project timelog operations")
@dataclass
class Timelogs:
    subcommand: cappa.Subcommands[TimelogsList | TimelogsAdd]


@cappa.command(name="list", help="List project timelogs")
@dataclass
class TimelogsList:
    def __call__(self) -> None:
        raise NotImplementedError


@cappa.command(name="add", help="Add a timelog")
@dataclass
class TimelogsAdd:
    def __call__(self) -> None:
        raise NotImplementedError


@cappa.command(name="users", help="Project user operations")
@dataclass
class ProjectUsers:
    subcommand: cappa.Subcommands[ProjectUsersList]


@cappa.command(name="list", help="List project users")
@dataclass
class ProjectUsersList:
    def __call__(self) -> None:
        raise NotImplementedError


@cappa.command(name="list", help="List all projects")
@dataclass
class ProjectsList:
    def __call__(self) -> None:
        raise NotImplementedError


@cappa.command(name="get", help="Get a single project")
@dataclass
class ProjectsGet:
    project_id: str

    def __call__(self) -> None:
        raise NotImplementedError


@cappa.command(name="search", help="Search projects")
@dataclass
class ProjectsSearch:
    def __call__(self) -> None:
        raise NotImplementedError


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
