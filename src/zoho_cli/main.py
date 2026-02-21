from __future__ import annotations

from dataclasses import dataclass

import cappa

from zoho_cli.auth.commands import Auth
from zoho_cli.cliq.app import Cliq
from zoho_cli.crm.app import Crm
from zoho_cli.drive.app import Drive
from zoho_cli.errors import ZohoCliError, err
from zoho_cli.projects.app import Projects
from zoho_cli.writer.app import Writer


@cappa.command(name="zoho", help="CLI for Zoho REST APIs (CRM, Projects, WorkDrive, Writer, Cliq)")
@dataclass
class Zoho:
    subcommand: cappa.Subcommands[Auth | Crm | Projects | Drive | Writer | Cliq]


def main() -> None:
    try:
        cappa.invoke(Zoho)
    except ZohoCliError as e:
        err(str(e))
        raise SystemExit(e.exit_code) from e
    except KeyboardInterrupt:
        err("\nInterrupted.")
        raise SystemExit(130) from None
