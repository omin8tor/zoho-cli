from __future__ import annotations

import json
import sys
from typing import Any


def print_json(data: Any, indent: int = 2) -> None:
    json.dump(data, sys.stdout, indent=indent, default=str, ensure_ascii=False)
    sys.stdout.write("\n")
