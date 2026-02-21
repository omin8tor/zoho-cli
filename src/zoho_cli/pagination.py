from __future__ import annotations

from typing import TYPE_CHECKING, Any

if TYPE_CHECKING:
    from zoho_cli.http.client import ZohoClient


def _has_next(page_info: Any) -> bool:
    if isinstance(page_info, dict):
        val = page_info.get("has_next_page")
        if val is None:
            return False
        if isinstance(val, bool):
            return val
        return str(val).lower() == "true"
    return False


def paginate_projects(
    client: ZohoClient,
    url: str,
    items_key: str | None,
    *,
    params: dict[str, Any] | None = None,
    max_pages: int = 20,
    per_page: int = 100,
) -> list[dict[str, Any]]:
    all_items: list[dict[str, Any]] = []
    p: dict[str, Any] = dict(params or {})
    page = 1

    for _ in range(max_pages):
        p["page"] = page
        p["per_page"] = per_page
        raw: Any = client.request("GET", url, params={k: str(v) for k, v in p.items()})

        if isinstance(raw, list):
            all_items.extend(raw)
            break

        items = raw.get(items_key, []) if items_key else (raw if isinstance(raw, list) else [])
        if not isinstance(items, list):
            break
        all_items.extend(items)

        page_info = raw.get("page_info", {})
        if isinstance(page_info, list):
            break
        if not _has_next(page_info):
            break
        if not items:
            break
        page += 1

    return all_items


def paginate_workdrive(
    client: ZohoClient,
    url: str,
    *,
    params: dict[str, Any] | None = None,
    max_pages: int = 10,
    per_page: int = 50,
) -> list[dict[str, Any]]:
    all_items: list[dict[str, Any]] = []
    p: dict[str, Any] = dict(params or {})
    p.setdefault("page[limit]", str(per_page))

    for _ in range(max_pages):
        p["page[offset]"] = str(len(all_items))
        data = client.request("GET", url, params={k: str(v) for k, v in p.items()})
        items = data.get("data", [])
        if not isinstance(items, list):
            break
        all_items.extend(items)

        if not data.get("meta", {}).get("has_next", False):
            break
        if not items:
            break

    return all_items


def paginate_crm(
    client: ZohoClient,
    url: str,
    *,
    params: dict[str, Any] | None = None,
    max_pages: int = 20,
    per_page: int = 200,
) -> list[dict[str, Any]]:
    all_items: list[dict[str, Any]] = []
    p: dict[str, Any] = dict(params or {})
    p.setdefault("per_page", str(per_page))
    page = 1

    for _ in range(max_pages):
        p["page"] = str(page)
        data = client.request("GET", url, params={k: str(v) for k, v in p.items()})
        items = data.get("data", [])
        if not isinstance(items, list):
            break
        all_items.extend(items)

        info = data.get("info", {})
        if not info.get("more_records", False):
            break
        page += 1

        if not items:
            break

    return all_items
