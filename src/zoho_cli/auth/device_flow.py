from __future__ import annotations

import sys
import time

import httpx

from zoho_cli.auth.config import save_client_config, save_tokens
from zoho_cli.errors import AuthError
from zoho_cli.http.dc import accounts_url as dc_accounts_url

DEFAULT_SCOPES = (
    "ZohoCliq.Webhooks.CREATE,ZohoCliq.Channels.ALL,ZohoCliq.Messages.ALL,"
    "ZohoCliq.Chats.ALL,ZohoCliq.Users.ALL,ZohoCliq.Bots.ALL,"
    "WorkDrive.workspace.ALL,WorkDrive.files.ALL,WorkDrive.files.sharing.ALL,"
    "WorkDrive.links.ALL,WorkDrive.team.ALL,WorkDrive.teamfolders.ALL,"
    "ZohoSearch.securesearch.READ,"
    "ZohoWriter.documentEditor.ALL,ZohoPC.files.ALL,"
    "ZohoProjects.portals.ALL,ZohoProjects.projects.ALL,ZohoProjects.tasks.ALL,"
    "ZohoProjects.tasklists.ALL,ZohoProjects.timesheets.ALL,ZohoProjects.bugs.ALL,"
    "ZohoProjects.events.ALL,ZohoProjects.forums.ALL,ZohoProjects.milestones.ALL,"
    "ZohoProjects.documents.ALL,ZohoProjects.users.ALL,"
    "ZohoCRM.modules.ALL,ZohoCRM.settings.ALL,ZohoCRM.users.ALL,ZohoCRM.org.ALL,"
    "ZohoCRM.change_owner.CREATE"
)

_LOCATION_TO_DC = {
    "us": "com",
    "eu": "eu",
    "in": "in",
    "au": "com.au",
    "jp": "jp",
    "ca": "ca",
    "sa": "sa",
    "uk": "uk",
    "cn": "com.cn",
}


def device_flow_login(
    client_id: str,
    client_secret: str,
    dc: str = "com",
    scopes: str | None = None,
) -> None:
    base_url = dc_accounts_url(dc)
    scopes = scopes or DEFAULT_SCOPES
    http = httpx.Client(timeout=30.0)

    resp = http.post(
        f"{base_url}/oauth/v3/device/code",
        params={
            "client_id": client_id,
            "grant_type": "device_request",
            "scope": scopes,
            "access_type": "offline",
        },
    )
    resp.raise_for_status()
    data = resp.json()

    if "error" in data:
        raise AuthError(f"Device code request failed: {data.get('error')}")

    user_code = data["user_code"]
    device_code = data["device_code"]
    verification_url = data.get("verification_url", f"{base_url}/oauth/v3/device")
    interval = max(data.get("interval", 5), 5)
    expires_in = data.get("expires_in", 300)

    verification_complete = data.get("verification_uri_complete", "")

    print(f"\nTo authenticate, visit: {verification_url}", file=sys.stderr)
    print(f"Enter code: {user_code}", file=sys.stderr)
    if verification_complete:
        print(f"\nOr open directly: {verification_complete}", file=sys.stderr)
    print(f"\nWaiting for authorization (expires in {expires_in}s)...", file=sys.stderr)

    poll_url = base_url
    deadline = time.monotonic() + expires_in

    while time.monotonic() < deadline:
        time.sleep(interval)

        poll_resp = http.post(
            f"{poll_url}/oauth/v3/device/token",
            params={
                "client_id": client_id,
                "client_secret": client_secret,
                "grant_type": "device_token",
                "code": device_code,
            },
        )
        poll_data = poll_resp.json()

        if "error" in poll_data:
            error = poll_data["error"]
            if error == "authorization_pending":
                continue
            if error == "slow_down":
                interval = min(interval + 2, 30)
                continue
            if error == "other_dc":
                new_dc = poll_data.get("dc", "")
                if new_dc:
                    dc = _LOCATION_TO_DC.get(new_dc, new_dc)
                    poll_url = dc_accounts_url(dc)
                    print(f"Redirecting to {dc} data center...", file=sys.stderr)
                continue
            if error == "access_denied":
                raise AuthError("Authorization denied by user.")
            if error == "expired":
                raise AuthError("Authorization code expired. Please try again.")
            raise AuthError(f"Device flow error: {error}")

        access_token = poll_data.get("access_token")
        refresh_token = poll_data.get("refresh_token")
        if not access_token or not refresh_token:
            raise AuthError(f"Unexpected response: {poll_data}")

        api_domain = poll_data.get("api_domain", "")
        expires = poll_data.get("expires_in", 3600)

        save_client_config(client_id, client_secret)
        save_tokens(
            refresh_token=refresh_token,
            access_token=access_token,
            expires_in=int(expires),
            dc=dc,
            accounts_url=dc_accounts_url(dc),
            api_domain=api_domain,
            scopes=scopes,
        )
        print("\nAuthenticated successfully!", file=sys.stderr)
        return

    raise AuthError("Authorization timed out. Please try again.")
