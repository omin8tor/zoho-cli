from __future__ import annotations

import json
import os
from datetime import UTC, datetime, timedelta
from unittest.mock import patch

import pytest

from zoho_cli.auth.config import AuthConfig, _from_env, resolve_auth, save_tokens
from zoho_cli.errors import AuthError


def test_auth_config_token_valid_when_fresh():
    config = AuthConfig(
        client_id="cid",
        client_secret="cs",
        refresh_token="rt",
        access_token="at",
        access_token_expires_at=datetime.now(UTC) + timedelta(hours=1),
    )
    assert config.token_valid


def test_auth_config_token_invalid_when_expired():
    config = AuthConfig(
        client_id="cid",
        client_secret="cs",
        refresh_token="rt",
        access_token="at",
        access_token_expires_at=datetime.now(UTC) - timedelta(minutes=1),
    )
    assert not config.token_valid


def test_auth_config_token_invalid_within_buffer():
    config = AuthConfig(
        client_id="cid",
        client_secret="cs",
        refresh_token="rt",
        access_token="at",
        access_token_expires_at=datetime.now(UTC) + timedelta(minutes=3),
    )
    assert not config.token_valid


def test_auth_config_token_invalid_when_missing():
    config = AuthConfig(
        client_id="cid",
        client_secret="cs",
        refresh_token="rt",
    )
    assert not config.token_valid


def test_from_env_returns_config_when_all_set():
    env = {
        "ZOHO_CLIENT_ID": "cid",
        "ZOHO_CLIENT_SECRET": "cs",
        "ZOHO_REFRESH_TOKEN": "rt",
        "ZOHO_DC": "eu",
    }
    with patch.dict(os.environ, env, clear=False):
        config = _from_env()
    assert config is not None
    assert config.client_id == "cid"
    assert config.dc == "eu"
    assert config.source == "env"


def test_from_env_returns_none_when_missing():
    env = {"ZOHO_CLIENT_ID": "cid"}
    with patch.dict(os.environ, env, clear=True):
        config = _from_env()
    assert config is None


def test_resolve_auth_raises_when_no_config():
    with patch.dict(os.environ, {}, clear=True):
        with patch("zoho_cli.auth.config._from_config_file", return_value=None):
            with pytest.raises(AuthError, match="Not authenticated"):
                resolve_auth()


def test_save_tokens_creates_file(tmp_path):
    tokens_file = tmp_path / "tokens.json"
    with (
        patch("zoho_cli.auth.config.CONFIG_DIR", tmp_path),
        patch("zoho_cli.auth.config.TOKENS_FILE", tokens_file),
    ):
        save_tokens(
            refresh_token="rt",
            access_token="at",
            expires_in=3600,
            dc="com",
            accounts_url="https://accounts.zoho.com",
            api_domain="https://www.zohoapis.com",
        )
    assert tokens_file.exists()
    data = json.loads(tokens_file.read_text())
    assert data["refresh_token"] == "rt"
    assert data["access_token"] == "at"
    assert data["dc"] == "com"
