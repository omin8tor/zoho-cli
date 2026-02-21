from __future__ import annotations

import sys

import httpx

from zoho_cli.auth.config import save_client_config, save_tokens
from zoho_cli.auth.token import _normalize_expires_in
from zoho_cli.errors import AuthError
from zoho_cli.http.dc import accounts_url as dc_accounts_url


def self_client_exchange(
    client_id: str,
    client_secret: str,
    code: str,
    dc: str = "com",
    accounts_server: str | None = None,
) -> None:
    base = accounts_server or dc_accounts_url(dc)

    resp = httpx.post(
        f"{base}/oauth/v2/token",
        params={
            "client_id": client_id,
            "client_secret": client_secret,
            "grant_type": "authorization_code",
            "code": code,
        },
        timeout=30.0,
    )
    resp.raise_for_status()
    data = resp.json()

    if "error" in data:
        raise AuthError(f"Token exchange failed: {data['error']}")

    access_token = data.get("access_token", "")
    refresh_token = data.get("refresh_token", "")
    if not access_token:
        raise AuthError(f"No access token in response: {data}")
    if not refresh_token:
        raise AuthError(
            "No refresh token returned. Make sure you generated a Self Client code "
            "with access_type=offline."
        )

    api_domain = data.get("api_domain", "")
    expires_in = _normalize_expires_in(data.get("expires_in"))

    save_client_config(client_id, client_secret)
    save_tokens(
        refresh_token=refresh_token,
        access_token=access_token,
        expires_in=expires_in,
        dc=dc,
        accounts_url=base,
        api_domain=api_domain,
    )
    print("Authenticated successfully!", file=sys.stderr)
