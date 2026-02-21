from __future__ import annotations

import sys

EXIT_SUCCESS = 0
EXIT_ERROR = 1
EXIT_AUTH = 2
EXIT_NOT_FOUND = 3
EXIT_VALIDATION = 4


class ZohoCliError(Exception):
    exit_code: int = EXIT_ERROR


class AuthError(ZohoCliError):
    exit_code = EXIT_AUTH


class NotFoundError(ZohoCliError):
    exit_code = EXIT_NOT_FOUND


class ValidationError(ZohoCliError):
    exit_code = EXIT_VALIDATION


class ZohoAPIError(ZohoCliError):
    def __init__(self, message: str, status_code: int | None = None):
        super().__init__(message)
        self.status_code = status_code
        if status_code == 401:
            self.exit_code = EXIT_AUTH
        elif status_code == 404:
            self.exit_code = EXIT_NOT_FOUND


def err(message: str) -> None:
    print(message, file=sys.stderr)
