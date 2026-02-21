from __future__ import annotations

import cappa

from zoho_cli.auth.commands import Auth, Status
from zoho_cli.main import Zoho


def test_parse_auth_status():
    result = cappa.parse(Zoho, argv=["auth", "status"])
    assert isinstance(result.subcommand, Auth)
    assert isinstance(result.subcommand.subcommand, Status)


def test_parse_auth_login():
    result = cappa.parse(
        Zoho,
        argv=["auth", "login", "--client-id", "cid", "--client-secret", "cs"],
    )
    assert isinstance(result.subcommand, Auth)
    login = result.subcommand.subcommand
    assert login.client_id == "cid"
    assert login.client_secret == "cs"
    assert login.dc == "com"
