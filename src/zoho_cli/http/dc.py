from __future__ import annotations

DC_MAP: dict[str, dict[str, str]] = {
    "com": {
        "accounts": "https://accounts.zoho.com",
        "cliq": "https://cliq.zoho.com",
        "api": "https://www.zohoapis.com",
        "crm": "https://zohoapis.com",
        "projects": "https://projectsapi.zoho.com",
        "workdrive": "https://workdrive.zoho.com",
        "writer": "https://www.zohoapis.com/writer",
        "download": "https://download.zoho.com",
    },
    "eu": {
        "accounts": "https://accounts.zoho.eu",
        "cliq": "https://cliq.zoho.eu",
        "api": "https://www.zohoapis.eu",
        "crm": "https://zohoapis.eu",
        "projects": "https://projectsapi.zoho.eu",
        "workdrive": "https://workdrive.zoho.eu",
        "writer": "https://www.zohoapis.eu/writer",
        "download": "https://download.zoho.eu",
    },
    "in": {
        "accounts": "https://accounts.zoho.in",
        "cliq": "https://cliq.zoho.in",
        "api": "https://www.zohoapis.in",
        "crm": "https://zohoapis.in",
        "projects": "https://projectsapi.zoho.in",
        "workdrive": "https://workdrive.zoho.in",
        "writer": "https://www.zohoapis.in/writer",
        "download": "https://download.zoho.in",
    },
    "com.au": {
        "accounts": "https://accounts.zoho.com.au",
        "cliq": "https://cliq.zoho.com.au",
        "api": "https://www.zohoapis.com.au",
        "crm": "https://zohoapis.com.au",
        "projects": "https://projectsapi.zoho.com.au",
        "workdrive": "https://workdrive.zoho.com.au",
        "writer": "https://www.zohoapis.com.au/writer",
        "download": "https://download.zoho.com.au",
    },
    "jp": {
        "accounts": "https://accounts.zoho.jp",
        "cliq": "https://cliq.zoho.jp",
        "api": "https://www.zohoapis.jp",
        "crm": "https://zohoapis.jp",
        "projects": "https://projectsapi.zoho.jp",
        "workdrive": "https://workdrive.zoho.jp",
        "writer": "https://www.zohoapis.jp/writer",
        "download": "https://download.zoho.jp",
    },
    "ca": {
        "accounts": "https://accounts.zohocloud.ca",
        "cliq": "https://cliq.zohocloud.ca",
        "api": "https://www.zohoapis.ca",
        "crm": "https://zohoapis.ca",
        "projects": "https://projectsapi.zohocloud.ca",
        "workdrive": "https://workdrive.zohocloud.ca",
        "writer": "https://www.zohoapis.ca/writer",
        "download": "https://download.zohocloud.ca",
    },
    "sa": {
        "accounts": "https://accounts.zoho.sa",
        "cliq": "https://cliq.zoho.sa",
        "api": "https://www.zohoapis.sa",
        "crm": "https://zohoapis.sa",
        "projects": "https://projectsapi.zoho.sa",
        "workdrive": "https://workdrive.zoho.sa",
        "writer": "https://www.zohoapis.sa/writer",
        "download": "https://download.zoho.sa",
    },
    "uk": {
        "accounts": "https://accounts.zoho.uk",
        "cliq": "https://cliq.zoho.uk",
        "api": "https://www.zohoapis.uk",
        "crm": "https://zohoapis.uk",
        "projects": "https://projectsapi.zoho.uk",
        "workdrive": "https://workdrive.zoho.uk",
        "writer": "https://www.zohoapis.uk/writer",
        "download": "https://download.zoho.uk",
    },
    "com.cn": {
        "accounts": "https://accounts.zoho.com.cn",
        "cliq": "https://cliq.zoho.com.cn",
        "api": "https://www.zohoapis.com.cn",
        "crm": "https://zohoapis.com.cn",
        "projects": "https://projectsapi.zoho.com.cn",
        "workdrive": "https://workdrive.zoho.com.cn",
        "writer": "https://www.zohoapis.com.cn/writer",
        "download": "https://download.zoho.com.cn",
    },
}

_DEFAULT_DC = "com"

VALID_DCS = frozenset(DC_MAP.keys())


def _get(dc: str, service: str) -> str:
    return DC_MAP.get(dc, DC_MAP[_DEFAULT_DC])[service]


def accounts_url(dc: str) -> str:
    return _get(dc, "accounts")


def cliq_url(dc: str) -> str:
    return _get(dc, "cliq")


def crm_url(dc: str) -> str:
    return _get(dc, "crm")


def projects_url(dc: str) -> str:
    return _get(dc, "projects")


def workdrive_url(dc: str) -> str:
    return _get(dc, "workdrive")


def writer_url(dc: str) -> str:
    return _get(dc, "writer")


def download_url(dc: str) -> str:
    return _get(dc, "download")
