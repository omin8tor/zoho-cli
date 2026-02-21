from __future__ import annotations

from zoho_cli.pagination import _has_next


def test_has_next_true_bool():
    assert _has_next({"has_next_page": True}) is True


def test_has_next_false_bool():
    assert _has_next({"has_next_page": False}) is False


def test_has_next_true_string():
    assert _has_next({"has_next_page": "true"}) is True


def test_has_next_false_string():
    assert _has_next({"has_next_page": "false"}) is False


def test_has_next_none():
    assert _has_next({"has_next_page": None}) is False


def test_has_next_missing():
    assert _has_next({}) is False


def test_has_next_not_dict():
    assert _has_next("invalid") is False
