from __future__ import annotations

from dataclasses import dataclass

import cappa


@cappa.command(name="modules", help="CRM module operations")
@dataclass
class Modules:
    subcommand: cappa.Subcommands[ModulesList | ModulesFields]


@cappa.command(name="list", help="List available CRM modules")
@dataclass
class ModulesList:
    def __call__(self) -> None:
        raise NotImplementedError


@cappa.command(name="fields", help="List fields for a CRM module")
@dataclass
class ModulesFields:
    module: str

    def __call__(self) -> None:
        raise NotImplementedError


@cappa.command(name="records", help="CRM record operations")
@dataclass
class Records:
    subcommand: cappa.Subcommands[
        RecordsList | RecordsGet | RecordsCreate | RecordsUpdate | RecordsDelete | RecordsSearch
    ]


@cappa.command(name="list", help="List records in a module")
@dataclass
class RecordsList:
    module: str

    def __call__(self) -> None:
        raise NotImplementedError


@cappa.command(name="get", help="Get a single record")
@dataclass
class RecordsGet:
    module: str
    record_id: str

    def __call__(self) -> None:
        raise NotImplementedError


@cappa.command(name="create", help="Create a record")
@dataclass
class RecordsCreate:
    module: str

    def __call__(self) -> None:
        raise NotImplementedError


@cappa.command(name="update", help="Update a record")
@dataclass
class RecordsUpdate:
    module: str
    record_id: str

    def __call__(self) -> None:
        raise NotImplementedError


@cappa.command(name="delete", help="Delete a record")
@dataclass
class RecordsDelete:
    module: str
    record_id: str

    def __call__(self) -> None:
        raise NotImplementedError


@cappa.command(name="search", help="Search records in a module")
@dataclass
class RecordsSearch:
    module: str

    def __call__(self) -> None:
        raise NotImplementedError


@cappa.command(name="notes", help="CRM record notes")
@dataclass
class Notes:
    subcommand: cappa.Subcommands[NotesList | NotesAdd | NotesUpdate | NotesDelete]


@cappa.command(name="list", help="List notes on a record")
@dataclass
class NotesList:
    module: str
    record_id: str

    def __call__(self) -> None:
        raise NotImplementedError


@cappa.command(name="add", help="Add a note to a record")
@dataclass
class NotesAdd:
    module: str
    record_id: str

    def __call__(self) -> None:
        raise NotImplementedError


@cappa.command(name="update", help="Update a note")
@dataclass
class NotesUpdate:
    note_id: str

    def __call__(self) -> None:
        raise NotImplementedError


@cappa.command(name="delete", help="Delete a note")
@dataclass
class NotesDelete:
    note_id: str

    def __call__(self) -> None:
        raise NotImplementedError


@cappa.command(name="related", help="CRM related records")
@dataclass
class Related:
    subcommand: cappa.Subcommands[RelatedList]


@cappa.command(name="list", help="List related records")
@dataclass
class RelatedList:
    module: str
    record_id: str
    related_list: str

    def __call__(self) -> None:
        raise NotImplementedError


@cappa.command(name="users", help="CRM users")
@dataclass
class Users:
    subcommand: cappa.Subcommands[UsersList]


@cappa.command(name="list", help="List CRM users")
@dataclass
class UsersList:
    def __call__(self) -> None:
        raise NotImplementedError


@cappa.command(name="owner", help="Change record owner")
@dataclass
class Owner:
    subcommand: cappa.Subcommands[OwnerChange]


@cappa.command(name="change", help="Change record owner")
@dataclass
class OwnerChange:
    module: str
    record_id: str

    def __call__(self) -> None:
        raise NotImplementedError


@cappa.command(name="crm", help="Zoho CRM operations")
@dataclass
class Crm:
    subcommand: cappa.Subcommands[Modules | Records | Notes | Related | Users | Owner]
