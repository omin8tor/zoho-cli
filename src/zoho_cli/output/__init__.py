from __future__ import annotations

from typing import Any

from zoho_cli.output.json import print_json
from zoho_cli.output.table import print_table


def output(data: Any, fmt: str = "json", title: str | None = None) -> None:
    if fmt == "table":
        print_table(data, title=title)
    else:
        print_json(data)
