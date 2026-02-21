from __future__ import annotations

import sys
from dataclasses import dataclass
from pathlib import Path
from typing import Annotated, Any

import cappa

from zoho_cli.http.client import ZohoClient, get_client
from zoho_cli.output import output

_JSONAPI_CT = "application/vnd.api+json"

_SERVICE_TYPE_MAP = {
    "writer": "zw",
    "sheet": "zohosheet",
    "show": "zohoshow",
}


@cappa.command(name="create", help="Create a new Writer document")
@dataclass
class WriterCreate:
    name: Annotated[str, cappa.Arg(long="--name", help="Document name")]
    parent: Annotated[str, cappa.Arg(long="--folder", help="Parent folder ID in WorkDrive")]
    doc_type: Annotated[
        str, cappa.Arg(long="--type", default="writer", help="writer, sheet, or show")
    ] = "writer"

    def __call__(self, client: Annotated[ZohoClient, cappa.Dep(get_client)]) -> None:
        url = f"{client.workdrive_base}/files"
        service_type = _SERVICE_TYPE_MAP.get(self.doc_type, "zw")
        body: dict[str, Any] = {
            "data": {
                "type": "files",
                "attributes": {
                    "name": self.name,
                    "parent_id": self.parent,
                    "service_type": service_type,
                },
            }
        }
        data = client.request("POST", url, json=body, headers={"Content-Type": _JSONAPI_CT})
        item = data.get("data", {})
        attrs = item.get("attributes", {})
        resource_id = attrs.get("resource_id") or item.get("id")
        output(
            {
                "id": resource_id,
                "name": attrs.get("name"),
                "type": attrs.get("type"),
                "permalink": attrs.get("permalink"),
            }
        )


@cappa.command(name="read", help="Read document content as text")
@dataclass
class WriterRead:
    doc_id: Annotated[str, cappa.Arg(help="Document ID")]
    format: Annotated[str, cappa.Arg(long="--format", default="txt", help="txt or html")] = "txt"

    def __call__(self, client: Annotated[ZohoClient, cappa.Dep(get_client)]) -> None:
        url = f"{client.writer_base}/download/{self.doc_id}"
        resp = client.request_raw("GET", url, params={"format": self.format})
        content = resp.text
        if not content:
            output({"error": "Document is empty or could not be read"})
            return
        output({"document_id": self.doc_id, "format": self.format, "content": content})


@cappa.command(name="download", help="Download a document")
@dataclass
class WriterDownload:
    doc_id: Annotated[str, cappa.Arg(help="Document ID")]
    format: Annotated[
        str, cappa.Arg(long="--format", default="txt", help="txt, html, pdf, docx, odt, rtf, epub")
    ] = "txt"
    output_path: Annotated[
        str | None, cappa.Arg(long="--output", default=None, help="Output file path")
    ] = None

    def __call__(self, client: Annotated[ZohoClient, cappa.Dep(get_client)]) -> None:
        url = f"{client.writer_base}/download/{self.doc_id}"
        resp = client.request_raw("GET", url, params={"format": self.format})
        if self.output_path:
            Path(self.output_path).write_bytes(resp.content)
            output({"ok": True, "path": self.output_path, "size": len(resp.content)})
        else:
            sys.stdout.buffer.write(resp.content)


@cappa.command(name="writer", help="Zoho Writer operations")
@dataclass
class Writer:
    subcommand: cappa.Subcommands[WriterCreate | WriterRead | WriterDownload]
