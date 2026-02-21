from __future__ import annotations

import json
import sys
from dataclasses import dataclass
from pathlib import Path
from typing import Annotated, Any

import cappa

from zoho_cli.errors import ZohoAPIError
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
        output(data)


@cappa.command(name="details", help="Get document metadata")
@dataclass
class WriterDetails:
    doc_id: Annotated[str, cappa.Arg(help="Document ID")]

    def __call__(self, client: Annotated[ZohoClient, cappa.Dep(get_client)]) -> None:
        url = f"{client.writer_base}/documents/{self.doc_id}"
        data = client.request("GET", url)
        output(data)


@cappa.command(name="fields", help="List merge fields in a document")
@dataclass
class WriterFields:
    doc_id: Annotated[str, cappa.Arg(help="Document ID")]

    def __call__(self, client: Annotated[ZohoClient, cappa.Dep(get_client)]) -> None:
        url = f"{client.writer_base}/documents/{self.doc_id}/fields"
        data = client.request("GET", url)
        output(data)


@cappa.command(name="merge", help="Merge data into a document template")
@dataclass
class WriterMerge:
    doc_id: Annotated[str, cappa.Arg(help="Document/template ID")]
    json_data: Annotated[str, cappa.Arg(long="--json", help="Merge data as JSON")]
    output_format: Annotated[
        str, cappa.Arg(long="--format", default="pdf", help="pdf, docx, or inline")
    ] = "pdf"
    output_path: Annotated[
        str | None, cappa.Arg(long="--output", default=None, help="Output file path")
    ] = None

    def __call__(self, client: Annotated[ZohoClient, cappa.Dep(get_client)]) -> None:
        url = f"{client.writer_base}/documents/{self.doc_id}/merge"
        merge_data = json.loads(self.json_data)
        body: dict[str, Any] = {
            "merge_data": merge_data,
            "output_format": self.output_format,
        }
        if self.output_format == "inline":
            data = client.request("POST", url, json=body)
            output(data)
        else:
            resp = client.request_raw(
                "POST",
                url,
                params={
                    "output_format": self.output_format,
                    "merge_data": json.dumps(merge_data),
                },
            )
            if self.output_path:
                Path(self.output_path).write_bytes(resp.content)
                output({"ok": True, "path": self.output_path, "size": len(resp.content)})
            else:
                sys.stdout.buffer.write(resp.content)


@cappa.command(name="delete", help="Delete a document")
@dataclass
class WriterDelete:
    doc_id: Annotated[str, cappa.Arg(help="Document ID")]

    def __call__(self, client: Annotated[ZohoClient, cappa.Dep(get_client)]) -> None:
        url = f"{client.writer_base}/documents/{self.doc_id}"
        data = client.request("DELETE", url)
        output(data)


@cappa.command(name="read", help="Read document content as text")
@dataclass
class WriterRead:
    doc_id: Annotated[str, cappa.Arg(help="Document ID")]
    format: Annotated[str, cappa.Arg(long="--format", default="txt", help="txt or html")] = "txt"

    def __call__(self, client: Annotated[ZohoClient, cappa.Dep(get_client)]) -> None:
        url = f"{client.writer_base}/download/{self.doc_id}"
        try:
            resp = client.request_raw("GET", url, params={"format": self.format})
        except ZohoAPIError as e:
            if "R3002" in str(e):
                output({"error": "Document is empty — Zoho cannot export empty documents (R3002)"})
                return
            raise
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
        try:
            resp = client.request_raw("GET", url, params={"format": self.format})
        except ZohoAPIError as e:
            if "R3002" in str(e):
                output({"error": "Document is empty — Zoho cannot export empty documents (R3002)"})
                return
            raise
        if self.output_path:
            Path(self.output_path).write_bytes(resp.content)
            output({"ok": True, "path": self.output_path, "size": len(resp.content)})
        else:
            sys.stdout.buffer.write(resp.content)


@cappa.command(name="writer", help="Zoho Writer operations")
@dataclass
class Writer:
    subcommand: cappa.Subcommands[
        WriterCreate
        | WriterDetails
        | WriterFields
        | WriterMerge
        | WriterDelete
        | WriterRead
        | WriterDownload
    ]
