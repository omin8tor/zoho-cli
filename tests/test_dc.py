from __future__ import annotations

from zoho_cli.http.dc import (
    DC_MAP,
    VALID_DCS,
    accounts_url,
    cliq_url,
    crm_url,
    download_url,
    projects_url,
    workdrive_url,
    writer_url,
)


def test_all_dcs_present():
    expected = {"com", "eu", "in", "com.au", "jp", "ca", "sa", "uk", "com.cn"}
    assert set(DC_MAP.keys()) == expected


def test_valid_dcs_frozenset():
    assert VALID_DCS == frozenset(DC_MAP.keys())


def test_com_urls():
    assert accounts_url("com") == "https://accounts.zoho.com"
    assert crm_url("com") == "https://zohoapis.com"
    assert projects_url("com") == "https://projectsapi.zoho.com"
    assert workdrive_url("com") == "https://workdrive.zoho.com"
    assert writer_url("com") == "https://www.zohoapis.com/writer"
    assert download_url("com") == "https://download.zoho.com"
    assert cliq_url("com") == "https://cliq.zoho.com"


def test_ca_irregularities():
    assert accounts_url("ca") == "https://accounts.zohocloud.ca"
    assert projects_url("ca") == "https://projectsapi.zohocloud.ca"
    assert workdrive_url("ca") == "https://workdrive.zohocloud.ca"
    assert crm_url("ca") == "https://zohoapis.ca"
    assert writer_url("ca") == "https://www.zohoapis.ca/writer"
    assert download_url("ca") == "https://download.zohocloud.ca"


def test_unknown_dc_falls_back_to_com():
    assert accounts_url("unknown") == "https://accounts.zoho.com"


def test_each_dc_has_all_services():
    services = {"accounts", "cliq", "api", "crm", "projects", "workdrive", "writer", "download"}
    for dc, urls in DC_MAP.items():
        assert set(urls.keys()) == services, f"DC {dc} missing services"
