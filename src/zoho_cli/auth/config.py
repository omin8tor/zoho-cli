from __future__ import annotations

import contextlib
import json
import os
import tomllib
from dataclasses import dataclass
from datetime import UTC, datetime, timedelta
from pathlib import Path

CONFIG_DIR = Path(os.environ.get("ZOHO_CLI_CONFIG_DIR", "~/.config/zoho-cli")).expanduser()
CONFIG_FILE = CONFIG_DIR / "config.toml"
TOKENS_FILE = CONFIG_DIR / "tokens.json"


@dataclass
class AuthConfig:
    client_id: str
    client_secret: str
    refresh_token: str
    dc: str = "com"
    accounts_url: str = ""
    api_domain: str = ""
    access_token: str = ""
    access_token_expires_at: datetime | None = None
    scopes: str = ""
    source: str = "unknown"

    def __post_init__(self) -> None:
        if not self.accounts_url:
            from zoho_cli.http.dc import accounts_url

            self.accounts_url = accounts_url(self.dc)

    @property
    def token_valid(self) -> bool:
        if not self.access_token or not self.access_token_expires_at:
            return False
        return datetime.now(UTC) + timedelta(minutes=5) < self.access_token_expires_at


@dataclass
class TokenData:
    refresh_token: str = ""
    access_token: str = ""
    access_token_expires_at: str = ""
    dc: str = "com"
    accounts_url: str = ""
    api_domain: str = ""
    scopes: str = ""


def _from_env() -> AuthConfig | None:
    client_id = os.environ.get("ZOHO_CLIENT_ID", "")
    client_secret = os.environ.get("ZOHO_CLIENT_SECRET", "")
    refresh_token = os.environ.get("ZOHO_REFRESH_TOKEN", "")

    if not all([client_id, client_secret, refresh_token]):
        return None

    dc = os.environ.get("ZOHO_DC", "com")
    accounts_url_env = os.environ.get("ZOHO_ACCOUNTS_URL", "")

    return AuthConfig(
        client_id=client_id,
        client_secret=client_secret,
        refresh_token=refresh_token,
        dc=dc,
        accounts_url=accounts_url_env,
        source="env",
    )


def _from_config_file() -> AuthConfig | None:
    if not CONFIG_FILE.exists():
        return None

    with open(CONFIG_FILE, "rb") as f:
        config = tomllib.load(f)

    auth = config.get("auth", {})
    client_id = auth.get("client_id", "")
    client_secret = auth.get("client_secret", "")

    if not client_id or not client_secret:
        return None

    tokens = _load_tokens()
    if not tokens or not tokens.refresh_token:
        return None

    expires_at = None
    if tokens.access_token_expires_at:
        with contextlib.suppress(ValueError):
            expires_at = datetime.fromisoformat(tokens.access_token_expires_at)

    return AuthConfig(
        client_id=client_id,
        client_secret=client_secret,
        refresh_token=tokens.refresh_token,
        dc=tokens.dc or "com",
        accounts_url=tokens.accounts_url,
        api_domain=tokens.api_domain,
        access_token=tokens.access_token,
        access_token_expires_at=expires_at,
        scopes=tokens.scopes,
        source="config",
    )


def _load_tokens() -> TokenData | None:
    if not TOKENS_FILE.exists():
        return None
    try:
        data = json.loads(TOKENS_FILE.read_text())
        return TokenData(**{k: v for k, v in data.items() if k in TokenData.__dataclass_fields__})
    except json.JSONDecodeError, TypeError:
        return None


def save_tokens(
    refresh_token: str,
    access_token: str,
    expires_in: int,
    dc: str,
    accounts_url: str,
    api_domain: str,
    scopes: str = "",
) -> None:
    CONFIG_DIR.mkdir(parents=True, exist_ok=True)
    expires_at = datetime.now(UTC) + timedelta(seconds=expires_in)
    data = {
        "refresh_token": refresh_token,
        "access_token": access_token,
        "access_token_expires_at": expires_at.isoformat(),
        "dc": dc,
        "accounts_url": accounts_url,
        "api_domain": api_domain,
        "scopes": scopes,
    }
    TOKENS_FILE.write_text(json.dumps(data, indent=2) + "\n")
    TOKENS_FILE.chmod(0o600)


def save_client_config(client_id: str, client_secret: str) -> None:
    CONFIG_DIR.mkdir(parents=True, exist_ok=True)
    import tomli_w

    config: dict[str, object] = {}
    if CONFIG_FILE.exists():
        with open(CONFIG_FILE, "rb") as f:
            config = tomllib.load(f)
    config["auth"] = {"client_id": client_id, "client_secret": client_secret}
    CONFIG_FILE.write_bytes(tomli_w.dumps(config).encode())
    CONFIG_FILE.chmod(0o600)


def resolve_auth() -> AuthConfig:
    from zoho_cli.errors import AuthError

    config = _from_env()
    if config:
        return config

    config = _from_config_file()
    if config:
        return config

    raise AuthError("Not authenticated. Run `zoho auth login` to authenticate.")
