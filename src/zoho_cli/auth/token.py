from __future__ import annotations

import httpx

from zoho_cli.auth.config import AuthConfig, save_tokens
from zoho_cli.errors import AuthError


def _normalize_expires_in(value: int | float | str | None) -> int:
    if value is None:
        return 3600
    try:
        expires = int(value)
    except TypeError, ValueError:
        return 3600
    if expires > 86400:
        return max(1, expires // 1000)
    return max(1, expires)


def refresh_access_token(config: AuthConfig) -> str:
    resp = httpx.post(
        f"{config.accounts_url}/oauth/v2/token",
        params={
            "client_id": config.client_id,
            "client_secret": config.client_secret,
            "grant_type": "refresh_token",
            "refresh_token": config.refresh_token,
        },
        timeout=30.0,
    )
    resp.raise_for_status()
    data = resp.json()

    if "error" in data:
        raise AuthError(f"Token refresh failed: {data['error']}")

    access_token: str = data["access_token"]
    expires_in = _normalize_expires_in(data.get("expires_in"))
    api_domain = data.get("api_domain", config.api_domain)

    config.access_token = access_token
    config.api_domain = api_domain

    if config.source == "config":
        save_tokens(
            refresh_token=config.refresh_token,
            access_token=access_token,
            expires_in=expires_in,
            dc=config.dc,
            accounts_url=config.accounts_url,
            api_domain=api_domain,
            scopes=config.scopes,
        )

    return access_token


def ensure_access_token(config: AuthConfig, force_refresh: bool = False) -> str:
    if not force_refresh and config.token_valid:
        return config.access_token
    return refresh_access_token(config)
