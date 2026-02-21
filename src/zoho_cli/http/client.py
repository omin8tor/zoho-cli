from __future__ import annotations

from typing import Any

import httpx

from zoho_cli.auth.config import AuthConfig, resolve_auth
from zoho_cli.auth.token import ensure_access_token
from zoho_cli.errors import ZohoAPIError
from zoho_cli.http.dc import (
    cliq_url,
    crm_url,
    download_url,
    projects_url,
    workdrive_url,
    writer_url,
)


class ZohoClient:
    def __init__(self, config: AuthConfig) -> None:
        self.config = config
        self._http = httpx.Client(timeout=60.0)
        self._access_token = ensure_access_token(config)

        dc = config.dc
        self.cliq_base = cliq_url(dc)
        self.crm_base = f"{crm_url(dc)}/crm/v8"
        self.projects_base = f"{projects_url(dc)}/api/v3"
        self.workdrive_base = f"{workdrive_url(dc)}/api/v1"
        self.writer_base = f"{writer_url(dc)}/api/v1"
        self.download_base = download_url(dc)

    def close(self) -> None:
        self._http.close()

    def _headers(self) -> dict[str, str]:
        return {"Authorization": f"Zoho-oauthtoken {self._access_token}"}

    def request(
        self,
        method: str,
        url: str,
        *,
        params: dict[str, str] | None = None,
        json: dict[str, Any] | None = None,
        data: dict[str, Any] | None = None,
        headers: dict[str, str] | None = None,
        files: dict[str, Any] | None = None,
    ) -> dict[str, Any]:
        merged = self._headers()
        if headers:
            merged.update(headers)

        kwargs: dict[str, Any] = {"headers": merged, "params": params, "data": data}
        if files:
            kwargs["files"] = files
        else:
            kwargs["json"] = json

        resp = self._http.request(method, url, **kwargs)

        if resp.status_code == 401:
            body = resp.text
            if "scope_invalid" in body or "scope_mismatch" in body:
                raise ZohoAPIError(
                    f"OAuth scope insufficient — re-authorize with correct scopes: {body}",
                    status_code=401,
                )
            self._access_token = ensure_access_token(self.config, force_refresh=True)
            kwargs["headers"] = self._headers()
            if headers:
                kwargs["headers"].update(headers)
            resp = self._http.request(method, url, **kwargs)
            if resp.status_code == 401:
                raise ZohoAPIError(
                    f"Access token expired or invalid after refresh: {resp.text}",
                    status_code=401,
                )

        if resp.status_code >= 400:
            raise ZohoAPIError(
                f"Zoho API error {resp.status_code}: {resp.text}",
                status_code=resp.status_code,
            )

        if not resp.content or resp.status_code == 204:
            return {}

        return resp.json()  # type: ignore[no-any-return]

    def request_raw(
        self,
        method: str,
        url: str,
        *,
        params: dict[str, str] | None = None,
    ) -> httpx.Response:
        resp = self._http.request(method, url, headers=self._headers(), params=params)

        if resp.status_code == 401:
            body = resp.text
            if "scope_invalid" in body or "scope_mismatch" in body:
                raise ZohoAPIError(
                    f"OAuth scope insufficient — re-authorize with correct scopes: {body}",
                    status_code=401,
                )
            self._access_token = ensure_access_token(self.config, force_refresh=True)
            resp = self._http.request(method, url, headers=self._headers(), params=params)
            if resp.status_code == 401:
                raise ZohoAPIError(
                    f"Access token expired or invalid after refresh: {resp.text}",
                    status_code=401,
                )

        if resp.status_code >= 400:
            raise ZohoAPIError(
                f"Zoho API error {resp.status_code}: {resp.text}",
                status_code=resp.status_code,
            )

        return resp


def get_client() -> ZohoClient:
    config = resolve_auth()
    return ZohoClient(config)
