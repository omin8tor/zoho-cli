from __future__ import annotations

from dataclasses import dataclass

import cappa


@cappa.command(name="list", help="List folder contents")
@dataclass
class FilesList:
    def __call__(self) -> None:
        raise NotImplementedError


@cappa.command(name="get", help="Get file info")
@dataclass
class FilesGet:
    file_id: str

    def __call__(self) -> None:
        raise NotImplementedError


@cappa.command(name="search", help="Search files")
@dataclass
class FilesSearch:
    def __call__(self) -> None:
        raise NotImplementedError


@cappa.command(name="rename", help="Rename a file")
@dataclass
class FilesRename:
    file_id: str

    def __call__(self) -> None:
        raise NotImplementedError


@cappa.command(name="copy", help="Copy a file")
@dataclass
class FilesCopy:
    file_id: str

    def __call__(self) -> None:
        raise NotImplementedError


@cappa.command(name="trash", help="Move a file to trash")
@dataclass
class FilesTrash:
    file_id: str

    def __call__(self) -> None:
        raise NotImplementedError


@cappa.command(name="delete", help="Permanently delete a file")
@dataclass
class FilesDelete:
    file_id: str

    def __call__(self) -> None:
        raise NotImplementedError


@cappa.command(name="files", help="WorkDrive file operations")
@dataclass
class Files:
    subcommand: cappa.Subcommands[
        FilesList | FilesGet | FilesSearch | FilesRename | FilesCopy | FilesTrash | FilesDelete
    ]


@cappa.command(name="list", help="List team folders")
@dataclass
class FoldersList:
    def __call__(self) -> None:
        raise NotImplementedError


@cappa.command(name="create", help="Create a folder")
@dataclass
class FoldersCreate:
    def __call__(self) -> None:
        raise NotImplementedError


@cappa.command(name="breadcrumb", help="Show folder path")
@dataclass
class FoldersBreadcrumb:
    folder_id: str

    def __call__(self) -> None:
        raise NotImplementedError


@cappa.command(name="folders", help="WorkDrive folder operations")
@dataclass
class Folders:
    subcommand: cappa.Subcommands[FoldersList | FoldersCreate | FoldersBreadcrumb]


@cappa.command(name="download", help="Download a file")
@dataclass
class Download:
    file_id: str

    def __call__(self) -> None:
        raise NotImplementedError


@cappa.command(name="upload", help="Upload a file")
@dataclass
class Upload:
    def __call__(self) -> None:
        raise NotImplementedError


@cappa.command(name="permissions", help="List file permissions")
@dataclass
class SharePermissions:
    file_id: str

    def __call__(self) -> None:
        raise NotImplementedError


@cappa.command(name="add", help="Share a file")
@dataclass
class ShareAdd:
    file_id: str

    def __call__(self) -> None:
        raise NotImplementedError


@cappa.command(name="revoke", help="Revoke file access")
@dataclass
class ShareRevoke:
    file_id: str

    def __call__(self) -> None:
        raise NotImplementedError


@cappa.command(name="link", help="Create/get share link")
@dataclass
class ShareLink:
    file_id: str

    def __call__(self) -> None:
        raise NotImplementedError


@cappa.command(name="share", help="File sharing operations")
@dataclass
class Share:
    subcommand: cappa.Subcommands[SharePermissions | ShareAdd | ShareRevoke | ShareLink]


@cappa.command(name="drive", help="Zoho WorkDrive operations")
@dataclass
class Drive:
    subcommand: cappa.Subcommands[Files | Folders | Download | Upload | Share]
