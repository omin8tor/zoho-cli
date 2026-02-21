from __future__ import annotations

from dataclasses import dataclass
from typing import Annotated

import cappa

from zoho_cli.auth.config import (
    CONFIG_FILE,
    TOKENS_FILE,
    resolve_auth,
)
from zoho_cli.auth.device_flow import device_flow_login
from zoho_cli.auth.self_client import self_client_exchange
from zoho_cli.auth.token import refresh_access_token
from zoho_cli.errors import AuthError, err
from zoho_cli.output import output


@cappa.command(name="login", help="Authenticate via device flow OAuth")
@dataclass
class Login:
    client_id: Annotated[
        str,
        cappa.Arg(long="--client-id", help="Zoho OAuth client ID"),
    ]
    client_secret: Annotated[
        str,
        cappa.Arg(long="--client-secret", help="Zoho OAuth client secret"),
    ]
    dc: Annotated[
        str,
        cappa.Arg(
            long="--dc",
            default="com",
            help="Data center (com, eu, in, com.au, jp, ca, sa, uk, com.cn)",
        ),
    ] = "com"
    scopes: Annotated[
        str | None,
        cappa.Arg(
            long="--scopes", default=None, help="Comma-separated OAuth scopes (defaults to all)"
        ),
    ] = None

    def __call__(self) -> None:
        device_flow_login(
            client_id=self.client_id,
            client_secret=self.client_secret,
            dc=self.dc,
            scopes=self.scopes,
        )


@cappa.command(name="self-client", help="Authenticate via self-client code exchange")
@dataclass
class SelfClient:
    code: Annotated[
        str,
        cappa.Arg(long="--code", help="Self-client authorization code from Zoho API Console"),
    ]
    client_id: Annotated[
        str,
        cappa.Arg(long="--client-id", help="Zoho OAuth client ID"),
    ]
    client_secret: Annotated[
        str,
        cappa.Arg(long="--client-secret", help="Zoho OAuth client secret"),
    ]
    dc: Annotated[
        str,
        cappa.Arg(long="--dc", default="com", help="Data center"),
    ] = "com"
    server: Annotated[
        str | None,
        cappa.Arg(long="--server", default=None, help="Accounts server URL override"),
    ] = None

    def __call__(self) -> None:
        self_client_exchange(
            client_id=self.client_id,
            client_secret=self.client_secret,
            code=self.code,
            dc=self.dc,
            accounts_server=self.server,
        )


@cappa.command(name="status", help="Show current authentication status")
@dataclass
class Status:
    def __call__(self) -> None:
        try:
            config = resolve_auth()
        except AuthError:
            err("Not authenticated. Run `zoho auth login` to authenticate.")
            raise SystemExit(2) from None

        info = {
            "authenticated": True,
            "source": config.source,
            "dc": config.dc,
            "accounts_url": config.accounts_url,
            "token_valid": config.token_valid,
        }
        output(info)


@cappa.command(name="refresh", help="Force refresh the access token")
@dataclass
class Refresh:
    def __call__(self) -> None:
        config = resolve_auth()
        token = refresh_access_token(config)
        err("Access token refreshed successfully.")
        output({"access_token": token[:20] + "...", "status": "refreshed"})


@cappa.command(name="logout", help="Clear stored authentication tokens")
@dataclass
class Logout:
    def __call__(self) -> None:
        removed = False
        if TOKENS_FILE.exists():
            TOKENS_FILE.unlink()
            removed = True
        if CONFIG_FILE.exists():
            CONFIG_FILE.unlink()
            removed = True
        if removed:
            err("Logged out. Stored credentials removed.")
        else:
            err("No stored credentials found.")


@cappa.command(name="auth", help="Authentication management")
@dataclass
class Auth:
    subcommand: cappa.Subcommands[Login | SelfClient | Status | Refresh | Logout]
