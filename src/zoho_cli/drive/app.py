from __future__ import annotations

import sys
from dataclasses import dataclass
from pathlib import Path
from typing import Annotated, Any

import cappa

from zoho_cli.errors import ZohoAPIError
from zoho_cli.http.client import ZohoClient, get_client
from zoho_cli.output import output
from zoho_cli.pagination import paginate_workdrive

_JSONAPI_CT = "application/vnd.api+json"

_SERVICE_TYPE_MAP = {
    "zohowriter": "zw",
    "zohosheet": "zohosheet",
    "zohoshow": "zohoshow",
}

_ROLE_IDS = {
    "viewer": 7,
    "commenter": 6,
    "editor": 5,
    "organizer": 4,
}


def _resource_id(item: dict[str, Any]) -> str | None:
    attrs = item.get("attributes", {})
    return attrs.get("resource_id") or item.get("id") or None


def _slim_file(item: dict[str, Any]) -> dict[str, Any]:
    attrs = item.get("attributes", {})
    storage = attrs.get("storage_info")
    return {
        "id": _resource_id(item),
        "name": attrs.get("name"),
        "type": attrs.get("type"),
        "extension": attrs.get("extension"),
        "size": storage.get("size") if isinstance(storage, dict) else None,
        "modified_time": attrs.get("modified_time"),
        "parent_id": attrs.get("parent_id"),
        "permalink": attrs.get("permalink"),
    }


def _slim_folder(item: dict[str, Any]) -> dict[str, Any]:
    attrs = item.get("attributes", {})
    return {
        "id": _resource_id(item),
        "name": attrs.get("name"),
        "type": "folder",
    }


def _slim_permission(item: dict[str, Any]) -> dict[str, Any]:
    attrs = item.get("attributes", {})
    return {
        "permission_id": item.get("id"),
        "email": attrs.get("email_id"),
        "display_name": attrs.get("display_name"),
        "role": attrs.get("role_name"),
        "role_id": attrs.get("role_id"),
        "shared_by": attrs.get("shared_by"),
        "shared_time": attrs.get("shared_time"),
        "expiration_date": attrs.get("expiration_date"),
    }


def _slim_link(item: dict[str, Any]) -> dict[str, Any]:
    attrs = item.get("attributes", {})
    return {
        "link_id": item.get("id"),
        "link": attrs.get("link"),
        "link_name": attrs.get("link_name"),
        "role_id": attrs.get("role_id"),
        "allow_download": attrs.get("allow_download"),
        "expiration_date": attrs.get("expiration_date"),
    }


def _jsonapi_body(attrs: dict[str, Any], *, resource_id: str | None = None) -> dict[str, Any]:
    item: dict[str, Any] = {"type": "files", "attributes": attrs}
    if resource_id:
        item["id"] = resource_id
    return {"data": item}


def _perm_body(attrs: dict[str, Any]) -> dict[str, Any]:
    return {"data": {"type": "permissions", "attributes": attrs}}


def _link_body(attrs: dict[str, Any]) -> dict[str, Any]:
    return {"data": {"type": "links", "attributes": attrs}}


@cappa.command(name="list", help="List folder contents")
@dataclass
class FilesList:
    folder: Annotated[str, cappa.Arg(long="--folder", help="Folder ID")]
    file_type: Annotated[
        str | None, cappa.Arg(long="--type", default=None, help="Filter: file, folder, image, etc.")
    ] = None

    def __call__(self, client: Annotated[ZohoClient, cappa.Dep(get_client)]) -> None:
        url = f"{client.workdrive_base}/files/{self.folder}/files"
        params: dict[str, str] = {}
        if self.file_type:
            params["filter[type]"] = self.file_type
        items = paginate_workdrive(client, url, params=params)
        output([_slim_file(item) for item in items])


@cappa.command(name="get", help="Get file info")
@dataclass
class FilesGet:
    file_id: Annotated[str, cappa.Arg(help="File or folder ID")]

    def __call__(self, client: Annotated[ZohoClient, cappa.Dep(get_client)]) -> None:
        url = f"{client.workdrive_base}/files/{self.file_id}"
        data = client.request("GET", url)
        item = data.get("data", {})
        attrs = item.get("attributes", {})
        storage = attrs.get("storage_info")

        bc_url = f"{client.workdrive_base}/files/{self.file_id}/breadcrumbs"
        path = None
        try:
            bc_data = client.request("GET", bc_url)
            bc_items = bc_data.get("data", [])
            if bc_items:
                parent_ids = bc_items[0].get("attributes", {}).get("parent_ids", [])
                path_parts = [p.get("name", "") for p in parent_ids]
                path = "/" + "/".join(path_parts) if path_parts else None
        except ZohoAPIError:
            pass

        info: dict[str, Any] = {
            "id": _resource_id(item),
            "name": attrs.get("name"),
            "type": attrs.get("type"),
            "extension": attrs.get("extension"),
            "size": storage.get("size") if isinstance(storage, dict) else None,
            "created_time": attrs.get("created_time"),
            "modified_time": attrs.get("modified_time"),
            "parent_id": attrs.get("parent_id"),
            "permalink": attrs.get("permalink"),
            "path": path,
        }

        file_type = attrs.get("type", "")
        if file_type in ("document", "spreadsheet", "presentation", "zdoc"):
            try:
                doc_url = f"{client.writer_base}/documents/{self.file_id}"
                doc_data = client.request("GET", doc_url)
                info["version"] = doc_data.get("version")
                info["open_url"] = doc_data.get("open_url")
                info["preview_url"] = doc_data.get("preview_url")
                info["download_url"] = doc_data.get("download_url")
            except ZohoAPIError:
                pass

        output(info)


@cappa.command(name="search", help="Search files")
@dataclass
class FilesSearch:
    query: Annotated[str, cappa.Arg(long="--query", help="Search keyword")]
    team: Annotated[str, cappa.Arg(long="--team", help="Team ID")]
    mode: Annotated[str, cappa.Arg(long="--mode", default="all", help="all, name, or content")] = (
        "all"
    )
    file_type: Annotated[
        str | None, cappa.Arg(long="--type", default=None, help="Filter by type")
    ] = None

    def __call__(self, client: Annotated[ZohoClient, cappa.Dep(get_client)]) -> None:
        url = f"{client.workdrive_base}/teams/{self.team}/records"
        params: dict[str, str] = {f"search[{self.mode}]": self.query}
        if self.file_type:
            params["filter[type]"] = self.file_type
        items = paginate_workdrive(client, url, params=params)
        output([_slim_file(item) for item in items])


@cappa.command(name="rename", help="Rename a file")
@dataclass
class FilesRename:
    file_id: Annotated[str, cappa.Arg(help="File ID")]
    name: Annotated[str, cappa.Arg(long="--name", help="New name")]

    def __call__(self, client: Annotated[ZohoClient, cappa.Dep(get_client)]) -> None:
        body = _jsonapi_body({"name": self.name})
        data = client.request(
            "PATCH",
            f"{client.workdrive_base}/files/{self.file_id}",
            json=body,
            headers={"Content-Type": _JSONAPI_CT},
        )
        item = data.get("data", {})
        attrs = item.get("attributes", {})
        output({"ok": True, "id": _resource_id(item), "name": attrs.get("name")})


@cappa.command(name="copy", help="Copy a file")
@dataclass
class FilesCopy:
    file_id: Annotated[str, cappa.Arg(help="Source file ID")]
    to: Annotated[str, cappa.Arg(long="--to", help="Destination folder ID")]

    def __call__(self, client: Annotated[ZohoClient, cappa.Dep(get_client)]) -> None:
        body = _jsonapi_body({"resource_id": self.file_id})
        data = client.request(
            "POST",
            f"{client.workdrive_base}/files/{self.to}/copy",
            json=body,
            headers={"Content-Type": _JSONAPI_CT},
        )
        item = data.get("data", {})
        attrs = item.get("attributes", {})
        output({"ok": True, "id": _resource_id(item), "name": attrs.get("name")})


@cappa.command(name="trash", help="Move a file to trash")
@dataclass
class FilesTrash:
    file_id: Annotated[str, cappa.Arg(help="File ID")]

    def __call__(self, client: Annotated[ZohoClient, cappa.Dep(get_client)]) -> None:
        body = _jsonapi_body({"status": 51})
        client.request(
            "PATCH",
            f"{client.workdrive_base}/files/{self.file_id}",
            json=body,
            headers={"Content-Type": _JSONAPI_CT},
        )
        output({"ok": True, "action": "trash", "id": self.file_id})


@cappa.command(name="delete", help="Permanently delete a file")
@dataclass
class FilesDelete:
    file_id: Annotated[str, cappa.Arg(help="File ID")]

    def __call__(self, client: Annotated[ZohoClient, cappa.Dep(get_client)]) -> None:
        body = _jsonapi_body({"status": 61})
        client.request(
            "PATCH",
            f"{client.workdrive_base}/files/{self.file_id}",
            json=body,
            headers={"Content-Type": _JSONAPI_CT},
        )
        output({"ok": True, "action": "delete", "id": self.file_id})


@cappa.command(name="files", help="WorkDrive file operations")
@dataclass
class Files:
    subcommand: cappa.Subcommands[
        FilesList | FilesGet | FilesSearch | FilesRename | FilesCopy | FilesTrash | FilesDelete
    ]


@cappa.command(name="list", help="List team folders")
@dataclass
class FoldersList:
    team: Annotated[str, cappa.Arg(long="--team", help="Team ID")]

    def __call__(self, client: Annotated[ZohoClient, cappa.Dep(get_client)]) -> None:
        url = f"{client.workdrive_base}/teams/{self.team}/teamfolders"
        folders = paginate_workdrive(client, url)
        output([_slim_folder(f) for f in folders])


@cappa.command(name="create", help="Create a folder or document")
@dataclass
class FoldersCreate:
    name: Annotated[str, cappa.Arg(long="--name", help="Name")]
    parent: Annotated[str, cappa.Arg(long="--parent", help="Parent folder ID")]
    file_type: Annotated[
        str,
        cappa.Arg(long="--type", default="folder", help="folder, zohowriter, zohosheet, zohoshow"),
    ] = "folder"

    def __call__(self, client: Annotated[ZohoClient, cappa.Dep(get_client)]) -> None:
        url = f"{client.workdrive_base}/files"
        if self.file_type == "folder":
            body = _jsonapi_body({"name": self.name, "parent_id": self.parent})
        else:
            service_type = _SERVICE_TYPE_MAP.get(self.file_type, self.file_type)
            body = _jsonapi_body(
                {
                    "name": self.name,
                    "parent_id": self.parent,
                    "service_type": service_type,
                }
            )
        data = client.request("POST", url, json=body, headers={"Content-Type": _JSONAPI_CT})
        item = data.get("data", {})
        attrs = item.get("attributes", {})
        output(
            {
                "id": _resource_id(item),
                "name": attrs.get("name"),
                "type": attrs.get("type"),
                "permalink": attrs.get("permalink"),
            }
        )


@cappa.command(name="breadcrumb", help="Show folder path")
@dataclass
class FoldersBreadcrumb:
    folder_id: Annotated[str, cappa.Arg(help="Folder ID")]

    def __call__(self, client: Annotated[ZohoClient, cappa.Dep(get_client)]) -> None:
        url = f"{client.workdrive_base}/files/{self.folder_id}/breadcrumbs"
        data = client.request("GET", url)
        bc_items = data.get("data", [])
        if bc_items:
            parent_ids = bc_items[0].get("attributes", {}).get("parent_ids", [])
            path_parts = [p.get("name", "") for p in parent_ids]
            output({"path": "/" + "/".join(path_parts), "parts": path_parts})
        else:
            output({"path": "/", "parts": []})


@cappa.command(name="folders", help="WorkDrive folder operations")
@dataclass
class Folders:
    subcommand: cappa.Subcommands[FoldersList | FoldersCreate | FoldersBreadcrumb]


@cappa.command(name="download", help="Download a file")
@dataclass
class Download:
    file_id: Annotated[str, cappa.Arg(help="File ID")]
    output_path: Annotated[
        str | None, cappa.Arg(long="--output", default=None, help="Output file path")
    ] = None
    format: Annotated[
        str, cappa.Arg(long="--format", default="native", help="native, txt, html, pdf, docx")
    ] = "native"

    def __call__(self, client: Annotated[ZohoClient, cappa.Dep(get_client)]) -> None:
        if self.format != "native":
            url = f"{client.writer_base}/download/{self.file_id}"
            resp = client.request_raw("GET", url, params={"format": self.format})
        else:
            url = f"{client.download_base}/v1/workdrive/download/{self.file_id}"
            resp = client.request_raw("GET", url)

        if self.output_path:
            Path(self.output_path).write_bytes(resp.content)
            output({"ok": True, "path": self.output_path, "size": len(resp.content)})
        else:
            sys.stdout.buffer.write(resp.content)


@cappa.command(name="upload", help="Upload a file")
@dataclass
class Upload:
    file_path: Annotated[str, cappa.Arg(help="Local file path to upload")]
    folder: Annotated[str, cappa.Arg(long="--folder", help="Destination folder ID")]
    override: Annotated[
        bool, cappa.Arg(long="--override", default=False, help="Override existing file")
    ] = False

    def __call__(self, client: Annotated[ZohoClient, cappa.Dep(get_client)]) -> None:
        path = Path(self.file_path)
        file_bytes = path.read_bytes()
        url = f"{client.workdrive_base}/upload"
        files = {"content": (path.name, file_bytes)}
        form_data: dict[str, str] = {"parent_id": self.folder, "filename": path.name}
        if self.override:
            form_data["override-name-exist"] = "true"
        data = client.request("POST", url, data=form_data, files=files)
        items = data.get("data", [])
        if isinstance(items, list) and items:
            item = items[0]
            attrs = item.get("attributes", {})
            output(
                {
                    "id": attrs.get("resource_id") or item.get("id"),
                    "name": attrs.get("FileName"),
                    "permalink": attrs.get("Permalink"),
                    "parent_id": attrs.get("parent_id"),
                }
            )
        else:
            output({"ok": True, "raw": data})


@cappa.command(name="permissions", help="List file permissions")
@dataclass
class SharePermissions:
    file_id: Annotated[str, cappa.Arg(help="File ID")]

    def __call__(self, client: Annotated[ZohoClient, cappa.Dep(get_client)]) -> None:
        url = f"{client.workdrive_base}/files/{self.file_id}/permissions"
        data = client.request("GET", url)
        items = data.get("data", [])
        output([_slim_permission(item) for item in items])


@cappa.command(name="add", help="Share a file")
@dataclass
class ShareAdd:
    file_id: Annotated[str, cappa.Arg(help="File ID")]
    email: Annotated[str, cappa.Arg(long="--email", help="Email to share with")]
    role: Annotated[
        str, cappa.Arg(long="--role", default="viewer", help="viewer, commenter, editor, organizer")
    ] = "viewer"

    def __call__(self, client: Annotated[ZohoClient, cappa.Dep(get_client)]) -> None:
        attrs: dict[str, Any] = {
            "resource_id": self.file_id,
            "shared_type": "personal",
            "email_id": self.email,
            "role_id": _ROLE_IDS.get(self.role, 7),
            "send_notification_mail": True,
        }
        url = f"{client.workdrive_base}/permissions"
        data = client.request(
            "POST", url, json=_perm_body(attrs), headers={"Content-Type": _JSONAPI_CT}
        )
        item = data.get("data", {})
        output(
            {"ok": True, "permission_id": item.get("id"), "email": self.email, "role": self.role}
        )


@cappa.command(name="revoke", help="Revoke file access")
@dataclass
class ShareRevoke:
    permission_id: Annotated[str, cappa.Arg(help="Permission ID")]

    def __call__(self, client: Annotated[ZohoClient, cappa.Dep(get_client)]) -> None:
        url = f"{client.workdrive_base}/permissions/{self.permission_id}"
        client.request("DELETE", url)
        output({"ok": True, "revoked": self.permission_id})


@cappa.command(name="link", help="Create/get share link")
@dataclass
class ShareLink:
    file_id: Annotated[str, cappa.Arg(help="File ID")]
    role: Annotated[
        str, cappa.Arg(long="--role", default="viewer", help="viewer, commenter, editor")
    ] = "viewer"
    allow_download: Annotated[bool, cappa.Arg(long="--allow-download", default=True)] = True
    link_name: Annotated[
        str | None, cappa.Arg(long="--name", default=None, help="Link display name")
    ] = None
    expiration: Annotated[
        str | None, cappa.Arg(long="--expiration", default=None, help="Expiry date YYYY-MM-DD")
    ] = None
    password: Annotated[
        str | None, cappa.Arg(long="--password", default=None, help="Link password")
    ] = None

    def __call__(self, client: Annotated[ZohoClient, cappa.Dep(get_client)]) -> None:
        attrs: dict[str, Any] = {
            "resource_id": self.file_id,
            "role_id": _ROLE_IDS.get(self.role, 7),
            "allow_download": self.allow_download,
            "request_user_data": False,
        }
        if self.link_name:
            attrs["link_name"] = self.link_name
        if self.expiration:
            attrs["expiration_date"] = self.expiration
        if self.password:
            attrs["password_text"] = self.password
        url = f"{client.workdrive_base}/links"
        data = client.request(
            "POST", url, json=_link_body(attrs), headers={"Content-Type": _JSONAPI_CT}
        )
        item = data.get("data", {})
        attrs_resp = item.get("attributes", {})
        output(
            {
                "ok": True,
                "link_id": item.get("id"),
                "link": attrs_resp.get("link"),
                "role": self.role,
            }
        )


@cappa.command(name="share", help="File sharing operations")
@dataclass
class Share:
    subcommand: cappa.Subcommands[SharePermissions | ShareAdd | ShareRevoke | ShareLink]


@cappa.command(name="drive", help="Zoho WorkDrive operations")
@dataclass
class Drive:
    subcommand: cappa.Subcommands[Files | Folders | Download | Upload | Share]
