from __future__ import annotations

import sys
from dataclasses import dataclass
from pathlib import Path
from typing import Annotated, Any

import cappa

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
        output(items)


@cappa.command(name="get", help="Get file info")
@dataclass
class FilesGet:
    file_id: Annotated[str, cappa.Arg(help="File or folder ID")]

    def __call__(self, client: Annotated[ZohoClient, cappa.Dep(get_client)]) -> None:
        url = f"{client.workdrive_base}/files/{self.file_id}"
        data = client.request("GET", url)
        output(data)


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
        output(items)


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
        output(data)


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
        output(data)


@cappa.command(name="trash", help="Move a file to trash")
@dataclass
class FilesTrash:
    file_id: Annotated[str, cappa.Arg(help="File ID")]

    def __call__(self, client: Annotated[ZohoClient, cappa.Dep(get_client)]) -> None:
        body = _jsonapi_body({"status": 51})
        data = client.request(
            "PATCH",
            f"{client.workdrive_base}/files/{self.file_id}",
            json=body,
            headers={"Content-Type": _JSONAPI_CT},
        )
        output(data)


@cappa.command(name="delete", help="Permanently delete a file")
@dataclass
class FilesDelete:
    file_id: Annotated[str, cappa.Arg(help="File ID")]

    def __call__(self, client: Annotated[ZohoClient, cappa.Dep(get_client)]) -> None:
        body = _jsonapi_body({"status": 61})
        data = client.request(
            "PATCH",
            f"{client.workdrive_base}/files/{self.file_id}",
            json=body,
            headers={"Content-Type": _JSONAPI_CT},
        )
        output(data)


@cappa.command(name="move", help="Move a file to a different folder")
@dataclass
class FilesMove:
    file_id: Annotated[str, cappa.Arg(help="File ID")]
    to: Annotated[str, cappa.Arg(long="--to", help="Destination folder ID")]

    def __call__(self, client: Annotated[ZohoClient, cappa.Dep(get_client)]) -> None:
        body = _jsonapi_body({"parent_id": self.to})
        data = client.request(
            "PATCH",
            f"{client.workdrive_base}/files/{self.file_id}",
            json=body,
            headers={"Content-Type": _JSONAPI_CT},
        )
        output(data)


@cappa.command(name="restore", help="Restore a file from trash")
@dataclass
class FilesRestore:
    file_id: Annotated[str, cappa.Arg(help="File ID")]

    def __call__(self, client: Annotated[ZohoClient, cappa.Dep(get_client)]) -> None:
        body = _jsonapi_body({"status": 1})
        data = client.request(
            "PATCH",
            f"{client.workdrive_base}/files/{self.file_id}",
            json=body,
            headers={"Content-Type": _JSONAPI_CT},
        )
        output(data)


@cappa.command(name="trash-list", help="List trashed files in a folder")
@dataclass
class FilesTrashList:
    folder: Annotated[str, cappa.Arg(long="--folder", help="Folder ID")]

    def __call__(self, client: Annotated[ZohoClient, cappa.Dep(get_client)]) -> None:
        url = f"{client.workdrive_base}/files/{self.folder}/files"
        items = paginate_workdrive(client, url, params={"filter[status]": "51"})
        output(items)


@cappa.command(name="versions", help="List file versions")
@dataclass
class FilesVersions:
    file_id: Annotated[str, cappa.Arg(help="File ID")]

    def __call__(self, client: Annotated[ZohoClient, cappa.Dep(get_client)]) -> None:
        url = f"{client.workdrive_base}/files/{self.file_id}/versions"
        data = client.request("GET", url)
        output(data)


@cappa.command(name="files", help="WorkDrive file operations")
@dataclass
class Files:
    subcommand: cappa.Subcommands[
        FilesList
        | FilesGet
        | FilesSearch
        | FilesRename
        | FilesCopy
        | FilesMove
        | FilesTrash
        | FilesDelete
        | FilesRestore
        | FilesTrashList
        | FilesVersions
    ]


@cappa.command(name="list", help="List team folders")
@dataclass
class FoldersList:
    team: Annotated[str, cappa.Arg(long="--team", help="Team ID")]

    def __call__(self, client: Annotated[ZohoClient, cappa.Dep(get_client)]) -> None:
        url = f"{client.workdrive_base}/teams/{self.team}/teamfolders"
        folders = paginate_workdrive(client, url)
        output(folders)


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
        output(data)


@cappa.command(name="breadcrumb", help="Show folder path")
@dataclass
class FoldersBreadcrumb:
    folder_id: Annotated[str, cappa.Arg(help="Folder ID")]

    def __call__(self, client: Annotated[ZohoClient, cappa.Dep(get_client)]) -> None:
        url = f"{client.workdrive_base}/files/{self.folder_id}/breadcrumbs"
        data = client.request("GET", url)
        output(data)


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
            url = f"{client.workdrive_base}/download/{self.file_id}"
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
        output(data)


@cappa.command(name="upload-url", help="Upload a file from a URL")
@dataclass
class UploadUrl:
    url_to_fetch: Annotated[str, cappa.Arg(help="URL of file to upload")]
    folder: Annotated[str, cappa.Arg(long="--folder", help="Destination folder ID")]
    name: Annotated[str | None, cappa.Arg(long="--name", default=None, help="File name")] = None

    def __call__(self, client: Annotated[ZohoClient, cappa.Dep(get_client)]) -> None:
        url = f"{client.workdrive_base}/files/{self.folder}/remotefile"
        body: dict[str, Any] = {"data": {"attributes": {"url": self.url_to_fetch}}}
        if self.name:
            body["data"]["attributes"]["name"] = self.name
        data = client.request("POST", url, json=body, headers={"Content-Type": _JSONAPI_CT})
        output(data)


@cappa.command(name="permissions", help="List file permissions")
@dataclass
class SharePermissions:
    file_id: Annotated[str, cappa.Arg(help="File ID")]

    def __call__(self, client: Annotated[ZohoClient, cappa.Dep(get_client)]) -> None:
        url = f"{client.workdrive_base}/files/{self.file_id}/permissions"
        data = client.request("GET", url)
        output(data)


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
        output(data)


@cappa.command(name="revoke", help="Revoke file access")
@dataclass
class ShareRevoke:
    permission_id: Annotated[str, cappa.Arg(help="Permission ID")]

    def __call__(self, client: Annotated[ZohoClient, cappa.Dep(get_client)]) -> None:
        url = f"{client.workdrive_base}/permissions/{self.permission_id}"
        client.request("DELETE", url)
        output({"ok": True, "revoked": self.permission_id})


@cappa.command(name="links", help="List share links for a file")
@dataclass
class ShareLinksList:
    file_id: Annotated[str, cappa.Arg(help="File ID")]

    def __call__(self, client: Annotated[ZohoClient, cappa.Dep(get_client)]) -> None:
        url = f"{client.workdrive_base}/files/{self.file_id}/links"
        data = client.request("GET", url)
        output(data)


@cappa.command(name="unlink", help="Delete a share link")
@dataclass
class ShareUnlink:
    link_id: Annotated[str, cappa.Arg(help="Link ID")]

    def __call__(self, client: Annotated[ZohoClient, cappa.Dep(get_client)]) -> None:
        url = f"{client.workdrive_base}/links/{self.link_id}"
        client.request("DELETE", url)
        output({"ok": True, "deleted": self.link_id})


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
        output(data)


@cappa.command(name="share", help="File sharing operations")
@dataclass
class Share:
    subcommand: cappa.Subcommands[
        SharePermissions | ShareAdd | ShareRevoke | ShareLink | ShareLinksList | ShareUnlink
    ]


@cappa.command(name="me", help="Get current user info")
@dataclass
class TeamsMe:
    def __call__(self, client: Annotated[ZohoClient, cappa.Dep(get_client)]) -> None:
        url = f"{client.workdrive_base}/users/me"
        data = client.request("GET", url)
        output(data)


@cappa.command(name="list", help="List teams")
@dataclass
class TeamsList:
    def __call__(self, client: Annotated[ZohoClient, cappa.Dep(get_client)]) -> None:
        url = f"{client.workdrive_base}/users/me/teams"
        data = client.request("GET", url)
        output(data)


@cappa.command(name="members", help="List team members")
@dataclass
class TeamsMembers:
    team: Annotated[str, cappa.Arg(help="Team ID")]

    def __call__(self, client: Annotated[ZohoClient, cappa.Dep(get_client)]) -> None:
        url = f"{client.workdrive_base}/teams/{self.team}/members"
        data = client.request("GET", url)
        output(data)


@cappa.command(name="teams", help="WorkDrive team operations")
@dataclass
class Teams:
    subcommand: cappa.Subcommands[TeamsMe | TeamsList | TeamsMembers]


@cappa.command(name="drive", help="Zoho WorkDrive operations")
@dataclass
class Drive:
    subcommand: cappa.Subcommands[Files | Folders | Download | Upload | UploadUrl | Share | Teams]
