from __future__ import annotations

from dataclasses import dataclass
from typing import Annotated, Any

import cappa

from zoho_cli.http.client import ZohoClient, get_client
from zoho_cli.output import output


@cappa.command(name="list", help="List channels")
@dataclass
class ChannelsList:
    def __call__(self, client: Annotated[ZohoClient, cappa.Dep(get_client)]) -> None:
        url = f"{client.cliq_base}/api/v2/channels"
        data = client.request("GET", url)
        channels = data.get("data", data.get("channels", []))
        if isinstance(channels, list):
            output(channels)
        else:
            output(data)


@cappa.command(name="message", help="Send a message to a channel")
@dataclass
class ChannelsMessage:
    channel: Annotated[str, cappa.Arg(help="Channel name")]
    text: Annotated[str, cappa.Arg(long="--text", help="Message text")]
    bot: Annotated[
        str | None, cappa.Arg(long="--bot", default=None, help="Bot name to send as")
    ] = None

    def __call__(self, client: Annotated[ZohoClient, cappa.Dep(get_client)]) -> None:
        url = f"{client.cliq_base}/api/v2/channelsbyname/{self.channel}/message"
        body: dict[str, Any] = {"text": self.text}
        if self.bot:
            body["bot"] = {"name": self.bot}
        data = client.request("POST", url, json=body)
        output({"ok": True, "channel": self.channel, "response": data})


@cappa.command(name="channels", help="Cliq channel operations")
@dataclass
class Channels:
    subcommand: cappa.Subcommands[ChannelsList | ChannelsMessage]


@cappa.command(name="message", help="Send a direct message")
@dataclass
class ChatsMessage:
    chat_id: Annotated[str, cappa.Arg(help="Chat ID")]
    text: Annotated[str, cappa.Arg(long="--text", help="Message text")]

    def __call__(self, client: Annotated[ZohoClient, cappa.Dep(get_client)]) -> None:
        url = f"{client.cliq_base}/api/v2/chats/{self.chat_id}/message"
        body: dict[str, Any] = {"text": self.text}
        data = client.request("POST", url, json=body)
        output({"ok": True, "chat_id": self.chat_id, "response": data})


@cappa.command(name="chats", help="Cliq chat operations")
@dataclass
class Chats:
    subcommand: cappa.Subcommands[ChatsMessage]


@cappa.command(name="cliq", help="Zoho Cliq operations")
@dataclass
class Cliq:
    subcommand: cappa.Subcommands[Channels | Chats]
