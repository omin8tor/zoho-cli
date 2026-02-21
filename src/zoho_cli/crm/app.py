from __future__ import annotations

import json
from dataclasses import dataclass
from typing import Annotated, Any

import cappa

from zoho_cli.http.client import ZohoClient, get_client
from zoho_cli.output import output
from zoho_cli.pagination import paginate_crm

_NOTE_FIELDS = "id,Note_Title,Note_Content,Created_Time,Modified_Time,Parent_Id"

_RELATED_DEFAULT_FIELDS: dict[str, str] = {
    "Notes": _NOTE_FIELDS,
    "Contacts": "id,Full_Name,Email,Phone,Created_Time",
    "Deals": "id,Deal_Name,Stage,Amount,Closing_Date",
    "Accounts": "id,Account_Name,Phone,Industry,Created_Time",
    "Leads": "id,Full_Name,Email,Company,Lead_Status",
    "Tasks": "id,Subject,Status,Due_Date,Priority",
    "Events": "id,Event_Title,Start_DateTime,End_DateTime",
    "Calls": "id,Subject,Call_Type,Call_Start_Time,Call_Duration",
    "Products": "id,Product_Name,Unit_Price,Product_Code",
    "Campaigns": "id,Campaign_Name,Type,Status,Start_Date",
}
_RELATED_FALLBACK_FIELDS = "id,Created_Time,Modified_Time"

_MODULE_DEFAULT_FIELDS: dict[str, str] = {
    "Leads": "id,Full_Name,Email,Company,Lead_Status,Phone,Created_Time",
    "Contacts": "id,Full_Name,Email,Phone,Account_Name,Created_Time",
    "Deals": "id,Deal_Name,Stage,Amount,Closing_Date,Account_Name,Created_Time",
    "Accounts": "id,Account_Name,Phone,Industry,Website,Created_Time",
    "Tasks": "id,Subject,Status,Due_Date,Priority,Created_Time",
    "Events": "id,Event_Title,Start_DateTime,End_DateTime,Created_Time",
    "Calls": "id,Subject,Call_Type,Call_Start_Time,Call_Duration,Created_Time",
}
_MODULE_FALLBACK_FIELDS = "id,Created_Time,Modified_Time"

_RECORD_STRIP = frozenset(
    {
        "Created_By",
        "Modified_By",
        "$editable",
        "$permissions",
        "$review_process",
        "$sharing_permission",
        "$approval",
        "$orchestration",
        "$in_merge",
        "$approval_state",
        "$pathfinder",
        "$review",
        "$currency_symbol",
    }
)


def _slim_module(m: dict[str, Any]) -> dict[str, Any]:
    return {
        "id": m.get("id"),
        "api_name": m.get("api_name"),
        "singular_label": m.get("singular_label"),
        "plural_label": m.get("plural_label"),
        "creatable": m.get("creatable"),
        "editable": m.get("editable"),
        "deletable": m.get("deletable"),
        "searchable": m.get("global_search_supported"),
    }


def _slim_field(f: dict[str, Any]) -> dict[str, Any]:
    return {
        "id": f.get("id"),
        "api_name": f.get("api_name"),
        "display_label": f.get("display_label"),
        "data_type": f.get("data_type"),
        "max_length": f.get("length"),
        "required": f.get("system_mandatory") or f.get("required"),
        "read_only": f.get("read_only"),
    }


def _slim_record(r: dict[str, Any]) -> dict[str, Any]:
    owner = r.get("Owner")
    result: dict[str, Any] = {
        "id": r.get("id"),
        "Owner": owner.get("name") if isinstance(owner, dict) else owner,
        "Created_Time": r.get("Created_Time"),
        "Modified_Time": r.get("Modified_Time"),
    }
    for k, v in r.items():
        if k not in result and k not in _RECORD_STRIP:
            result[k] = v
    return result


def _slim_note(n: dict[str, Any]) -> dict[str, Any]:
    parent = n.get("Parent_Id")
    return {
        "id": n.get("id"),
        "title": n.get("Note_Title"),
        "content": n.get("Note_Content"),
        "created_time": n.get("Created_Time"),
        "modified_time": n.get("Modified_Time"),
        "parent_module": parent.get("module", {}).get("api_name")
        if isinstance(parent, dict)
        else None,
        "parent_id": parent.get("id") if isinstance(parent, dict) else None,
    }


def _slim_user(u: dict[str, Any]) -> dict[str, Any]:
    return {
        "id": u.get("id"),
        "name": u.get("full_name"),
        "email": u.get("email"),
        "role": u.get("role", {}).get("name") if isinstance(u.get("role"), dict) else None,
        "status": u.get("status"),
    }


def _normalize_action_result(data: dict[str, Any], module: str, action: str) -> dict[str, Any]:
    items = data.get("data", [])
    if isinstance(items, list) and items:
        item = items[0]
        details = item.get("details", {})
        return {
            "ok": item.get("code") == "SUCCESS",
            "module": module,
            "record_id": details.get("id"),
            "action": action,
            "message": item.get("message", ""),
        }
    code = data.get("code")
    if code:
        return {
            "ok": code == "SUCCESS",
            "module": module,
            "record_id": data.get("details", {}).get("id"),
            "action": action,
            "message": data.get("message", ""),
        }
    return {"ok": False, "module": module, "action": action, "message": "Unrecognized response"}


@cappa.command(name="list", help="List available CRM modules")
@dataclass
class ModulesList:
    include_hidden: Annotated[
        bool,
        cappa.Arg(long="--include-hidden", default=False, help="Include hidden/system modules"),
    ] = False

    def __call__(self, client: Annotated[ZohoClient, cappa.Dep(get_client)]) -> None:
        url = f"{client.crm_base}/settings/modules"
        data = client.request("GET", url)
        modules = data.get("modules", [])
        if not self.include_hidden:
            modules = [m for m in modules if m.get("show_as_tab", False)]
        output([_slim_module(m) for m in modules])


@cappa.command(name="fields", help="List fields for a CRM module")
@dataclass
class ModulesFields:
    module: Annotated[str, cappa.Arg(help="Module API name (e.g. Leads, Contacts, Deals)")]

    def __call__(self, client: Annotated[ZohoClient, cappa.Dep(get_client)]) -> None:
        url = f"{client.crm_base}/settings/fields"
        data = client.request("GET", url, params={"module": self.module})
        fields = data.get("fields", [])
        output([_slim_field(f) for f in fields])


@cappa.command(name="modules", help="CRM module operations")
@dataclass
class Modules:
    subcommand: cappa.Subcommands[ModulesList | ModulesFields]


@cappa.command(name="list", help="List records in a module")
@dataclass
class RecordsList:
    module: Annotated[str, cappa.Arg(help="Module API name")]
    fields: Annotated[
        str | None,
        cappa.Arg(long="--fields", default=None, help="Comma-separated field API names"),
    ] = None
    sort_by: Annotated[
        str | None,
        cappa.Arg(long="--sort-by", default=None, help="Field to sort by"),
    ] = None
    sort_order: Annotated[
        str | None,
        cappa.Arg(long="--sort-order", default=None, help="asc or desc"),
    ] = None
    page: Annotated[int, cappa.Arg(long="--page", default=1, help="Page number")] = 1
    per_page: Annotated[
        int, cappa.Arg(long="--per-page", default=200, help="Records per page (max 200)")
    ] = 200
    all_pages: Annotated[
        bool, cappa.Arg(long="--all", default=False, help="Auto-paginate all records")
    ] = False

    def __call__(self, client: Annotated[ZohoClient, cappa.Dep(get_client)]) -> None:
        url = f"{client.crm_base}/{self.module}"
        params: dict[str, str] = {}
        if self.fields:
            params["fields"] = self.fields
        else:
            params["fields"] = _MODULE_DEFAULT_FIELDS.get(self.module, _MODULE_FALLBACK_FIELDS)
        if self.sort_by:
            params["sort_by"] = self.sort_by
        if self.sort_order:
            params["sort_order"] = self.sort_order

        if self.all_pages:
            records = paginate_crm(client, url, params=params)
            output([_slim_record(r) for r in records])
        else:
            params["page"] = str(self.page)
            params["per_page"] = str(min(self.per_page, 200))
            data = client.request("GET", url, params=params)
            records = data.get("data", [])
            output([_slim_record(r) for r in records])


@cappa.command(name="get", help="Get a single record")
@dataclass
class RecordsGet:
    module: Annotated[str, cappa.Arg(help="Module API name")]
    record_id: Annotated[str, cappa.Arg(help="Record ID")]
    fields: Annotated[
        str | None,
        cappa.Arg(long="--fields", default=None, help="Comma-separated field API names"),
    ] = None

    def __call__(self, client: Annotated[ZohoClient, cappa.Dep(get_client)]) -> None:
        url = f"{client.crm_base}/{self.module}/{self.record_id}"
        params: dict[str, str] = {}
        if self.fields:
            params["fields"] = self.fields
        data = client.request("GET", url, params=params or None)
        records = data.get("data", [])
        if records:
            output(_slim_record(records[0]))
        else:
            output({"error": f"Record {self.record_id} not found in {self.module}"})


@cappa.command(name="create", help="Create a record")
@dataclass
class RecordsCreate:
    module: Annotated[str, cappa.Arg(help="Module API name")]
    json_data: Annotated[
        str,
        cappa.Arg(long="--json", help="Record data as JSON string"),
    ]
    trigger: Annotated[
        str | None,
        cappa.Arg(
            long="--trigger",
            default=None,
            help="Comma-separated triggers: approval,workflow,blueprint",
        ),
    ] = None

    def __call__(self, client: Annotated[ZohoClient, cappa.Dep(get_client)]) -> None:
        parsed = json.loads(self.json_data)
        body: dict[str, Any] = {"data": [parsed]}
        if self.trigger:
            body["trigger"] = self.trigger.split(",")
        url = f"{client.crm_base}/{self.module}"
        data = client.request("POST", url, json=body)
        output(_normalize_action_result(data, self.module, "create"))


@cappa.command(name="update", help="Update a record")
@dataclass
class RecordsUpdate:
    module: Annotated[str, cappa.Arg(help="Module API name")]
    record_id: Annotated[str, cappa.Arg(help="Record ID")]
    json_data: Annotated[str, cappa.Arg(long="--json", help="Fields to update as JSON")]
    trigger: Annotated[
        str | None,
        cappa.Arg(long="--trigger", default=None, help="Comma-separated triggers"),
    ] = None

    def __call__(self, client: Annotated[ZohoClient, cappa.Dep(get_client)]) -> None:
        parsed = json.loads(self.json_data)
        body: dict[str, Any] = {"data": [parsed]}
        if self.trigger:
            body["trigger"] = self.trigger.split(",")
        url = f"{client.crm_base}/{self.module}/{self.record_id}"
        data = client.request("PUT", url, json=body)
        output(_normalize_action_result(data, self.module, "update"))


@cappa.command(name="delete", help="Delete a record")
@dataclass
class RecordsDelete:
    module: Annotated[str, cappa.Arg(help="Module API name")]
    record_id: Annotated[str, cappa.Arg(help="Record ID")]

    def __call__(self, client: Annotated[ZohoClient, cappa.Dep(get_client)]) -> None:
        url = f"{client.crm_base}/{self.module}/{self.record_id}"
        data = client.request("DELETE", url)
        output(_normalize_action_result(data, self.module, "delete"))


@cappa.command(name="search", help="Search records in a module")
@dataclass
class RecordsSearch:
    module: Annotated[str, cappa.Arg(help="Module API name")]
    word: Annotated[str | None, cappa.Arg(long="--word", default=None, help="Keyword search")] = (
        None
    )
    email: Annotated[str | None, cappa.Arg(long="--email", default=None, help="Email search")] = (
        None
    )
    phone: Annotated[str | None, cappa.Arg(long="--phone", default=None, help="Phone search")] = (
        None
    )
    criteria: Annotated[
        str | None,
        cappa.Arg(long="--criteria", default=None, help="Criteria e.g. (Stage:equals:Closed Won)"),
    ] = None
    fields: Annotated[
        str | None, cappa.Arg(long="--fields", default=None, help="Comma-separated fields")
    ] = None
    page: Annotated[int, cappa.Arg(long="--page", default=1)] = 1
    per_page: Annotated[int, cappa.Arg(long="--per-page", default=200)] = 200

    def __call__(self, client: Annotated[ZohoClient, cappa.Dep(get_client)]) -> None:
        url = f"{client.crm_base}/{self.module}/search"
        params: dict[str, str] = {
            "page": str(self.page),
            "per_page": str(min(self.per_page, 200)),
        }
        if self.word:
            params["word"] = self.word
        elif self.email:
            params["email"] = self.email
        elif self.phone:
            params["phone"] = self.phone
        elif self.criteria:
            params["criteria"] = self.criteria
        if self.fields:
            params["fields"] = self.fields
        data = client.request("GET", url, params=params)
        records = data.get("data", [])
        output([_slim_record(r) for r in records])


@cappa.command(name="records", help="CRM record operations")
@dataclass
class Records:
    subcommand: cappa.Subcommands[
        RecordsList | RecordsGet | RecordsCreate | RecordsUpdate | RecordsDelete | RecordsSearch
    ]


@cappa.command(name="list", help="List notes on a record")
@dataclass
class NotesList:
    module: Annotated[str, cappa.Arg(help="Module API name")]
    record_id: Annotated[str, cappa.Arg(help="Record ID")]
    page: Annotated[int, cappa.Arg(long="--page", default=1)] = 1
    per_page: Annotated[int, cappa.Arg(long="--per-page", default=200)] = 200

    def __call__(self, client: Annotated[ZohoClient, cappa.Dep(get_client)]) -> None:
        url = f"{client.crm_base}/{self.module}/{self.record_id}/Notes"
        params: dict[str, str] = {
            "page": str(self.page),
            "per_page": str(min(self.per_page, 200)),
            "fields": _NOTE_FIELDS,
        }
        data = client.request("GET", url, params=params)
        notes = data.get("data", [])
        output([_slim_note(n) for n in notes])


@cappa.command(name="add", help="Add a note to a record")
@dataclass
class NotesAdd:
    module: Annotated[str, cappa.Arg(help="Module API name")]
    record_id: Annotated[str, cappa.Arg(help="Record ID")]
    content: Annotated[str, cappa.Arg(long="--content", help="Note content")]
    title: Annotated[str | None, cappa.Arg(long="--title", default=None, help="Note title")] = None

    def __call__(self, client: Annotated[ZohoClient, cappa.Dep(get_client)]) -> None:
        url = f"{client.crm_base}/{self.module}/{self.record_id}/Notes"
        note_data: dict[str, str] = {"Note_Content": self.content}
        if self.title:
            note_data["Note_Title"] = self.title
        data = client.request("POST", url, json={"data": [note_data]})
        items = data.get("data", [])
        if isinstance(items, list) and items:
            item = items[0]
            details = item.get("details", {})
            output(
                {
                    "ok": item.get("code") == "SUCCESS",
                    "note_id": details.get("id"),
                    "message": item.get("message", ""),
                }
            )
        else:
            output({"ok": True, "action": "add_note"})


@cappa.command(name="update", help="Update a note")
@dataclass
class NotesUpdate:
    note_id: Annotated[str, cappa.Arg(help="Note ID")]
    title: Annotated[str | None, cappa.Arg(long="--title", default=None)] = None
    content: Annotated[str | None, cappa.Arg(long="--content", default=None)] = None

    def __call__(self, client: Annotated[ZohoClient, cappa.Dep(get_client)]) -> None:
        url = f"{client.crm_base}/Notes/{self.note_id}"
        note_data: dict[str, str] = {}
        if self.title is not None:
            note_data["Note_Title"] = self.title
        if self.content is not None:
            note_data["Note_Content"] = self.content
        data = client.request("PUT", url, json={"data": [note_data]})
        items = data.get("data", [])
        if isinstance(items, list) and items:
            item = items[0]
            output(
                {
                    "ok": item.get("code") == "SUCCESS",
                    "note_id": self.note_id,
                    "message": item.get("message", ""),
                }
            )
        else:
            output({"ok": True, "note_id": self.note_id})


@cappa.command(name="delete", help="Delete a note")
@dataclass
class NotesDelete:
    note_id: Annotated[str, cappa.Arg(help="Note ID")]

    def __call__(self, client: Annotated[ZohoClient, cappa.Dep(get_client)]) -> None:
        url = f"{client.crm_base}/Notes/{self.note_id}"
        data = client.request("DELETE", url)
        items = data.get("data", [])
        if isinstance(items, list) and items:
            item = items[0]
            output(
                {
                    "ok": item.get("code") == "SUCCESS",
                    "note_id": self.note_id,
                    "message": item.get("message", ""),
                }
            )
        else:
            output({"ok": True, "note_id": self.note_id})


@cappa.command(name="notes", help="CRM record notes")
@dataclass
class Notes:
    subcommand: cappa.Subcommands[NotesList | NotesAdd | NotesUpdate | NotesDelete]


@cappa.command(name="list", help="List related records")
@dataclass
class RelatedList:
    module: Annotated[str, cappa.Arg(help="Parent module API name")]
    record_id: Annotated[str, cappa.Arg(help="Parent record ID")]
    related_list: Annotated[str, cappa.Arg(help="Related list API name")]
    fields: Annotated[
        str | None, cappa.Arg(long="--fields", default=None, help="Comma-separated fields")
    ] = None
    page: Annotated[int, cappa.Arg(long="--page", default=1)] = 1
    per_page: Annotated[int, cappa.Arg(long="--per-page", default=200)] = 200

    def __call__(self, client: Annotated[ZohoClient, cappa.Dep(get_client)]) -> None:
        url = f"{client.crm_base}/{self.module}/{self.record_id}/{self.related_list}"
        params: dict[str, str] = {
            "page": str(self.page),
            "per_page": str(min(self.per_page, 200)),
        }
        if self.fields:
            params["fields"] = self.fields
        else:
            params["fields"] = _RELATED_DEFAULT_FIELDS.get(
                self.related_list, _RELATED_FALLBACK_FIELDS
            )
        data = client.request("GET", url, params=params)
        records = data.get("data", [])
        if self.related_list == "Notes":
            output([_slim_note(r) for r in records])
        else:
            output([_slim_record(r) for r in records])


@cappa.command(name="related", help="CRM related records")
@dataclass
class Related:
    subcommand: cappa.Subcommands[RelatedList]


@cappa.command(name="list", help="List CRM users")
@dataclass
class UsersList:
    page: Annotated[int, cappa.Arg(long="--page", default=1)] = 1
    per_page: Annotated[int, cappa.Arg(long="--per-page", default=200)] = 200

    def __call__(self, client: Annotated[ZohoClient, cappa.Dep(get_client)]) -> None:
        url = f"{client.crm_base}/users"
        params: dict[str, str] = {
            "page": str(self.page),
            "per_page": str(min(self.per_page, 200)),
            "type": "AllUsers",
        }
        data = client.request("GET", url, params=params)
        users = data.get("users", [])
        output([_slim_user(u) for u in users])


@cappa.command(name="users", help="CRM users")
@dataclass
class Users:
    subcommand: cappa.Subcommands[UsersList]


@cappa.command(name="change", help="Change record owner")
@dataclass
class OwnerChange:
    module: Annotated[str, cappa.Arg(help="Module API name")]
    record_id: Annotated[str, cappa.Arg(help="Record ID")]
    owner: Annotated[str, cappa.Arg(long="--owner", help="New owner user ID")]
    notify: Annotated[bool, cappa.Arg(long="--notify", default=True, help="Notify new owner")] = (
        True
    )

    def __call__(self, client: Annotated[ZohoClient, cappa.Dep(get_client)]) -> None:
        url = f"{client.crm_base}/{self.module}/{self.record_id}/actions/change_owner"
        body = {"owner": {"id": self.owner}, "notify": self.notify}
        data = client.request("POST", url, json=body)
        output(
            {
                "ok": data.get("code") == "SUCCESS" if data.get("code") else True,
                "module": self.module,
                "record_id": self.record_id,
                "new_owner_id": self.owner,
                "action": "change_owner",
                "message": data.get("message", ""),
            }
        )


@cappa.command(name="owner", help="Change record owner")
@dataclass
class Owner:
    subcommand: cappa.Subcommands[OwnerChange]


@cappa.command(name="crm", help="Zoho CRM operations")
@dataclass
class Crm:
    subcommand: cappa.Subcommands[Modules | Records | Notes | Related | Users | Owner]
