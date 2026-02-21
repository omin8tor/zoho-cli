from __future__ import annotations

import json

from zoho_cli.output.json import print_json


def test_print_json_dict(capsys):
    print_json({"key": "value"})
    captured = capsys.readouterr()
    data = json.loads(captured.out)
    assert data == {"key": "value"}


def test_print_json_list(capsys):
    print_json([1, 2, 3])
    captured = capsys.readouterr()
    data = json.loads(captured.out)
    assert data == [1, 2, 3]
