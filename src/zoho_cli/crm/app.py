from __future__ import annotations

import json
import sys
from dataclasses import dataclass
from pathlib import Path
from typing import Annotated, Any

import cappa

from zoho_cli.http.client import ZohoClient, get_client
from zoho_cli.output import output
from zoho_cli.pagination import paginate_crm

_DEFAULT_FIELDS = "id,Created_Time,Modified_Time"


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
        output(modules)


@cappa.command(name="fields", help="List fields for a CRM module")
@dataclass
class ModulesFields:
    module: Annotated[str, cappa.Arg(help="Module API name (e.g. Leads, Contacts, Deals)")]

    def __call__(self, client: Annotated[ZohoClient, cappa.Dep(get_client)]) -> None:
        url = f"{client.crm_base}/settings/fields"
        data = client.request("GET", url, params={"module": self.module})
        output(data.get("fields", []))


@cappa.command(name="related-lists", help="List related lists for a module")
@dataclass
class ModulesRelatedLists:
    module: Annotated[str, cappa.Arg(help="Module API name")]

    def __call__(self, client: Annotated[ZohoClient, cappa.Dep(get_client)]) -> None:
        url = f"{client.crm_base}/settings/related_lists"
        data = client.request("GET", url, params={"module": self.module})
        output(data.get("related_lists", []))


@cappa.command(name="layouts", help="List layouts for a module")
@dataclass
class ModulesLayouts:
    module: Annotated[str, cappa.Arg(help="Module API name")]

    def __call__(self, client: Annotated[ZohoClient, cappa.Dep(get_client)]) -> None:
        url = f"{client.crm_base}/settings/layouts"
        data = client.request("GET", url, params={"module": self.module})
        output(data.get("layouts", []))


@cappa.command(name="custom-views", help="List custom views for a module")
@dataclass
class ModulesCustomViews:
    module: Annotated[str, cappa.Arg(help="Module API name")]

    def __call__(self, client: Annotated[ZohoClient, cappa.Dep(get_client)]) -> None:
        url = f"{client.crm_base}/settings/custom_views"
        data = client.request("GET", url, params={"module": self.module})
        output(data.get("custom_views", []))


@cappa.command(name="modules", help="CRM module operations")
@dataclass
class Modules:
    subcommand: cappa.Subcommands[
        ModulesList | ModulesFields | ModulesRelatedLists | ModulesLayouts | ModulesCustomViews
    ]


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
        params: dict[str, str] = {"fields": self.fields or _DEFAULT_FIELDS}
        if self.sort_by:
            params["sort_by"] = self.sort_by
        if self.sort_order:
            params["sort_order"] = self.sort_order

        if self.all_pages:
            records = paginate_crm(client, url, params=params)
            output(records)
        else:
            params["page"] = str(self.page)
            params["per_page"] = str(min(self.per_page, 200))
            data = client.request("GET", url, params=params)
            output(data.get("data", []))


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
            output(records[0])
        else:
            output(data)


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
        output(data)


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
        output(data)


@cappa.command(name="delete", help="Delete a record")
@dataclass
class RecordsDelete:
    module: Annotated[str, cappa.Arg(help="Module API name")]
    record_id: Annotated[str, cappa.Arg(help="Record ID")]

    def __call__(self, client: Annotated[ZohoClient, cappa.Dep(get_client)]) -> None:
        url = f"{client.crm_base}/{self.module}/{self.record_id}"
        data = client.request("DELETE", url)
        output(data)


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
        output(data.get("data", []))


@cappa.command(name="upsert", help="Upsert a record (insert or update)")
@dataclass
class RecordsUpsert:
    module: Annotated[str, cappa.Arg(help="Module API name")]
    json_data: Annotated[str, cappa.Arg(long="--json", help="Record data as JSON string")]
    duplicate_check_fields: Annotated[
        str | None,
        cappa.Arg(
            long="--duplicate-check", default=None, help="Comma-separated duplicate check fields"
        ),
    ] = None
    trigger: Annotated[
        str | None,
        cappa.Arg(long="--trigger", default=None, help="Comma-separated triggers"),
    ] = None

    def __call__(self, client: Annotated[ZohoClient, cappa.Dep(get_client)]) -> None:
        parsed = json.loads(self.json_data)
        body: dict[str, Any] = {"data": [parsed]}
        if self.duplicate_check_fields:
            body["duplicate_check_fields"] = self.duplicate_check_fields.split(",")
        if self.trigger:
            body["trigger"] = self.trigger.split(",")
        url = f"{client.crm_base}/{self.module}/upsert"
        data = client.request("POST", url, json=body)
        output(data)


@cappa.command(name="bulk-delete", help="Delete multiple records")
@dataclass
class RecordsBulkDelete:
    module: Annotated[str, cappa.Arg(help="Module API name")]
    ids: Annotated[str, cappa.Arg(help="Comma-separated record IDs")]

    def __call__(self, client: Annotated[ZohoClient, cappa.Dep(get_client)]) -> None:
        url = f"{client.crm_base}/{self.module}"
        data = client.request("DELETE", url, params={"ids": self.ids})
        output(data)


@cappa.command(name="records", help="CRM record operations")
@dataclass
class Records:
    subcommand: cappa.Subcommands[
        RecordsList
        | RecordsGet
        | RecordsCreate
        | RecordsUpdate
        | RecordsDelete
        | RecordsSearch
        | RecordsUpsert
        | RecordsBulkDelete
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
        }
        data = client.request("GET", url, params=params)
        output(data.get("data", []))


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
        output(data)


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
        output(data)


@cappa.command(name="delete", help="Delete a note")
@dataclass
class NotesDelete:
    note_id: Annotated[str, cappa.Arg(help="Note ID")]

    def __call__(self, client: Annotated[ZohoClient, cappa.Dep(get_client)]) -> None:
        url = f"{client.crm_base}/Notes/{self.note_id}"
        data = client.request("DELETE", url)
        output(data)


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
        data = client.request("GET", url, params=params)
        output(data.get("data", []))


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
        output(data.get("users", []))


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
        output(data)


@cappa.command(name="owner", help="Change record owner")
@dataclass
class Owner:
    subcommand: cappa.Subcommands[OwnerChange]


@cappa.command(name="coql", help="Run a COQL query")
@dataclass
class Coql:
    query: Annotated[str, cappa.Arg(long="--query", help="COQL query string")]

    def __call__(self, client: Annotated[ZohoClient, cappa.Dep(get_client)]) -> None:
        url = f"{client.crm_base}/coql"
        data = client.request("POST", url, json={"select_query": self.query})
        output(data)


@cappa.command(name="search-global", help="Search across all CRM modules")
@dataclass
class SearchGlobal:
    word: Annotated[str, cappa.Arg(help="Search keyword")]
    page: Annotated[int, cappa.Arg(long="--page", default=1)] = 1
    per_page: Annotated[int, cappa.Arg(long="--per-page", default=10)] = 10

    def __call__(self, client: Annotated[ZohoClient, cappa.Dep(get_client)]) -> None:
        url = f"{client.crm_base}/search"
        params: dict[str, str] = {
            "word": self.word,
            "page": str(self.page),
            "per_page": str(min(self.per_page, 10)),
        }
        data = client.request("GET", url, params=params)
        output(data)


@cappa.command(name="list", help="List attachments on a record")
@dataclass
class AttachmentsList:
    module: Annotated[str, cappa.Arg(help="Module API name")]
    record_id: Annotated[str, cappa.Arg(help="Record ID")]
    page: Annotated[int, cappa.Arg(long="--page", default=1)] = 1
    per_page: Annotated[int, cappa.Arg(long="--per-page", default=200)] = 200

    def __call__(self, client: Annotated[ZohoClient, cappa.Dep(get_client)]) -> None:
        url = f"{client.crm_base}/{self.module}/{self.record_id}/Attachments"
        params: dict[str, str] = {
            "page": str(self.page),
            "per_page": str(min(self.per_page, 200)),
        }
        data = client.request("GET", url, params=params)
        output(data.get("data", []))


@cappa.command(name="upload", help="Upload an attachment to a record")
@dataclass
class AttachmentsUpload:
    module: Annotated[str, cappa.Arg(help="Module API name")]
    record_id: Annotated[str, cappa.Arg(help="Record ID")]
    file_path: Annotated[str, cappa.Arg(help="Local file path")]

    def __call__(self, client: Annotated[ZohoClient, cappa.Dep(get_client)]) -> None:
        path = Path(self.file_path)
        url = f"{client.crm_base}/{self.module}/{self.record_id}/Attachments"
        data = client.request("POST", url, files={"file": (path.name, path.read_bytes())})
        output(data)


@cappa.command(name="download", help="Download an attachment")
@dataclass
class AttachmentsDownload:
    module: Annotated[str, cappa.Arg(help="Module API name")]
    record_id: Annotated[str, cappa.Arg(help="Record ID")]
    attachment_id: Annotated[str, cappa.Arg(help="Attachment ID")]
    output_path: Annotated[
        str | None, cappa.Arg(long="--output", default=None, help="Output file path")
    ] = None

    def __call__(self, client: Annotated[ZohoClient, cappa.Dep(get_client)]) -> None:
        url = f"{client.crm_base}/{self.module}/{self.record_id}/Attachments/{self.attachment_id}"
        resp = client.request_raw("GET", url)
        if self.output_path:
            Path(self.output_path).write_bytes(resp.content)
            output({"ok": True, "path": self.output_path, "size": len(resp.content)})
        else:
            sys.stdout.buffer.write(resp.content)


@cappa.command(name="delete", help="Delete an attachment")
@dataclass
class AttachmentsDelete:
    module: Annotated[str, cappa.Arg(help="Module API name")]
    record_id: Annotated[str, cappa.Arg(help="Record ID")]
    attachment_id: Annotated[str, cappa.Arg(help="Attachment ID")]

    def __call__(self, client: Annotated[ZohoClient, cappa.Dep(get_client)]) -> None:
        url = f"{client.crm_base}/{self.module}/{self.record_id}/Attachments/{self.attachment_id}"
        data = client.request("DELETE", url)
        output(data)


@cappa.command(name="attachments", help="CRM record attachments")
@dataclass
class Attachments:
    subcommand: cappa.Subcommands[
        AttachmentsList | AttachmentsUpload | AttachmentsDownload | AttachmentsDelete
    ]


@cappa.command(name="convert", help="Convert a lead to contact/account/deal")
@dataclass
class ConvertLead:
    record_id: Annotated[str, cappa.Arg(help="Lead record ID")]
    json_data: Annotated[
        str | None,
        cappa.Arg(long="--json", default=None, help="Conversion options as JSON"),
    ] = None

    def __call__(self, client: Annotated[ZohoClient, cappa.Dep(get_client)]) -> None:
        url = f"{client.crm_base}/Leads/{self.record_id}/actions/convert"
        body: dict[str, Any] = {"data": [json.loads(self.json_data) if self.json_data else {}]}
        data = client.request("POST", url, json=body)
        output(data)


@cappa.command(name="add", help="Add tags to records")
@dataclass
class TagsAdd:
    module: Annotated[str, cappa.Arg(help="Module API name")]
    ids: Annotated[str, cappa.Arg(long="--ids", help="Comma-separated record IDs")]
    tags: Annotated[str, cappa.Arg(long="--tags", help="Comma-separated tag names")]

    def __call__(self, client: Annotated[ZohoClient, cappa.Dep(get_client)]) -> None:
        url = f"{client.crm_base}/{self.module}/actions/add_tags"
        params = {"ids": self.ids, "tag_names": self.tags}
        data = client.request("POST", url, params=params)
        output(data)


@cappa.command(name="remove", help="Remove tags from records")
@dataclass
class TagsRemove:
    module: Annotated[str, cappa.Arg(help="Module API name")]
    ids: Annotated[str, cappa.Arg(long="--ids", help="Comma-separated record IDs")]
    tags: Annotated[str, cappa.Arg(long="--tags", help="Comma-separated tag names")]

    def __call__(self, client: Annotated[ZohoClient, cappa.Dep(get_client)]) -> None:
        url = f"{client.crm_base}/{self.module}/actions/remove_tags"
        params = {"ids": self.ids, "tag_names": self.tags}
        data = client.request("POST", url, params=params)
        output(data)


@cappa.command(name="tags", help="CRM tag operations")
@dataclass
class Tags:
    subcommand: cappa.Subcommands[TagsAdd | TagsRemove]


@cappa.command(name="crm", help="Zoho CRM operations")
@dataclass
class Crm:
    subcommand: cappa.Subcommands[
        Modules
        | Records
        | Notes
        | Related
        | Users
        | Owner
        | Coql
        | SearchGlobal
        | Attachments
        | ConvertLead
        | Tags
    ]
