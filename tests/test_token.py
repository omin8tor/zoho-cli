from __future__ import annotations

from zoho_cli.auth.token import _normalize_expires_in


def test_normalize_none():
    assert _normalize_expires_in(None) == 3600


def test_normalize_normal_seconds():
    assert _normalize_expires_in(3600) == 3600


def test_normalize_milliseconds():
    assert _normalize_expires_in(3600000) == 3600


def test_normalize_string():
    assert _normalize_expires_in("3600") == 3600


def test_normalize_invalid():
    assert _normalize_expires_in("invalid") == 3600


def test_normalize_zero():
    assert _normalize_expires_in(0) == 1
