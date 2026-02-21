from __future__ import annotations

from dataclasses import dataclass

import cappa


@cappa.command(name="create", help="Create a new Writer document")
@dataclass
class WriterCreate:
    def __call__(self) -> None:
        raise NotImplementedError


@cappa.command(name="read", help="Read document content as text")
@dataclass
class WriterRead:
    doc_id: str

    def __call__(self) -> None:
        raise NotImplementedError


@cappa.command(name="download", help="Download a document")
@dataclass
class WriterDownload:
    doc_id: str

    def __call__(self) -> None:
        raise NotImplementedError


@cappa.command(name="writer", help="Zoho Writer operations")
@dataclass
class Writer:
    subcommand: cappa.Subcommands[WriterCreate | WriterRead | WriterDownload]
