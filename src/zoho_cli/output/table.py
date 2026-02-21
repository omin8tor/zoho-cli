from __future__ import annotations

import sys
from typing import Any

from rich.console import Console
from rich.table import Table


def print_table(data: Any, title: str | None = None) -> None:
    console = Console(file=sys.stdout)

    if isinstance(data, dict):
        table = Table(title=title, show_header=True)
        table.add_column("Key")
        table.add_column("Value")
        for k, v in data.items():
            table.add_row(str(k), str(v))
        console.print(table)
        return

    if isinstance(data, list) and data and isinstance(data[0], dict):
        table = Table(title=title, show_header=True)
        cols = list(data[0].keys())
        for col in cols:
            table.add_column(col)
        for row in data:
            table.add_row(*(str(row.get(c, "")) for c in cols))
        console.print(table)
        return

    console.print(data)
