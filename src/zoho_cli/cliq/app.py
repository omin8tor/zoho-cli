from __future__ import annotations

from dataclasses import dataclass

import cappa


@cappa.command(name="list", help="List channels")
@dataclass
class ChannelsList:
    def __call__(self) -> None:
        raise NotImplementedError


@cappa.command(name="message", help="Send a message to a channel")
@dataclass
class ChannelsMessage:
    channel: str

    def __call__(self) -> None:
        raise NotImplementedError


@cappa.command(name="channels", help="Cliq channel operations")
@dataclass
class Channels:
    subcommand: cappa.Subcommands[ChannelsList | ChannelsMessage]


@cappa.command(name="message", help="Send a direct message")
@dataclass
class ChatsMessage:
    chat_id: str

    def __call__(self) -> None:
        raise NotImplementedError


@cappa.command(name="chats", help="Cliq chat operations")
@dataclass
class Chats:
    subcommand: cappa.Subcommands[ChatsMessage]


@cappa.command(name="cliq", help="Zoho Cliq operations")
@dataclass
class Cliq:
    subcommand: cappa.Subcommands[Channels | Chats]
