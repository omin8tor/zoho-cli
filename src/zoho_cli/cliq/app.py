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


@cappa.command(name="get", help="Get channel info")
@dataclass
class ChannelsGet:
    channel: Annotated[str, cappa.Arg(help="Channel unique name or ID")]

    def __call__(self, client: Annotated[ZohoClient, cappa.Dep(get_client)]) -> None:
        url = f"{client.cliq_base}/api/v2/channelsbyname/{self.channel}"
        data = client.request("GET", url)
        output(data)


@cappa.command(name="create", help="Create a channel")
@dataclass
class ChannelsCreate:
    name: Annotated[str, cappa.Arg(long="--name", help="Channel name")]
    description: Annotated[
        str | None, cappa.Arg(long="--description", default=None, help="Channel description")
    ] = None

    def __call__(self, client: Annotated[ZohoClient, cappa.Dep(get_client)]) -> None:
        url = f"{client.cliq_base}/api/v2/channels"
        body: dict[str, Any] = {"name": self.name}
        if self.description:
            body["description"] = self.description
        data = client.request("POST", url, json=body)
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


@cappa.command(name="members", help="List channel members")
@dataclass
class ChannelsMembers:
    channel: Annotated[str, cappa.Arg(help="Channel unique name or ID")]

    def __call__(self, client: Annotated[ZohoClient, cappa.Dep(get_client)]) -> None:
        url = f"{client.cliq_base}/api/v2/channelsbyname/{self.channel}/members"
        data = client.request("GET", url)
        output(data)


@cappa.command(name="channels", help="Cliq channel operations")
@dataclass
class Channels:
    subcommand: cappa.Subcommands[
        ChannelsList | ChannelsGet | ChannelsCreate | ChannelsMessage | ChannelsMembers
    ]


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


@cappa.command(name="message", help="Send a DM by email address")
@dataclass
class BuddiesMessage:
    email: Annotated[str, cappa.Arg(help="Recipient email address")]
    text: Annotated[str, cappa.Arg(long="--text", help="Message text")]

    def __call__(self, client: Annotated[ZohoClient, cappa.Dep(get_client)]) -> None:
        url = f"{client.cliq_base}/api/v2/buddies/{self.email}/message"
        body: dict[str, Any] = {"text": self.text}
        data = client.request("POST", url, json=body)
        output({"ok": True, "email": self.email, "response": data})


@cappa.command(name="buddies", help="Cliq buddy/DM operations")
@dataclass
class Buddies:
    subcommand: cappa.Subcommands[BuddiesMessage]


@cappa.command(name="list", help="List messages in a chat")
@dataclass
class MessagesList:
    chat_id: Annotated[str, cappa.Arg(help="Chat or channel ID")]
    limit: Annotated[int, cappa.Arg(long="--limit", default=50, help="Number of messages")] = 50

    def __call__(self, client: Annotated[ZohoClient, cappa.Dep(get_client)]) -> None:
        url = f"{client.cliq_base}/api/v2/chats/{self.chat_id}/messages"
        data = client.request("GET", url, params={"limit": str(self.limit)})
        output(data)


@cappa.command(name="edit", help="Edit a message")
@dataclass
class MessagesEdit:
    chat_id: Annotated[str, cappa.Arg(help="Chat ID")]
    message_id: Annotated[str, cappa.Arg(help="Message ID")]
    text: Annotated[str, cappa.Arg(long="--text", help="New message text")]

    def __call__(self, client: Annotated[ZohoClient, cappa.Dep(get_client)]) -> None:
        url = f"{client.cliq_base}/api/v2/chats/{self.chat_id}/messages/{self.message_id}"
        data = client.request("PUT", url, json={"text": self.text})
        output(data)


@cappa.command(name="delete", help="Delete a message")
@dataclass
class MessagesDelete:
    chat_id: Annotated[str, cappa.Arg(help="Chat ID")]
    message_id: Annotated[str, cappa.Arg(help="Message ID")]

    def __call__(self, client: Annotated[ZohoClient, cappa.Dep(get_client)]) -> None:
        url = f"{client.cliq_base}/api/v2/chats/{self.chat_id}/messages/{self.message_id}"
        data = client.request("DELETE", url)
        output(data)


@cappa.command(name="messages", help="Cliq message operations")
@dataclass
class Messages:
    subcommand: cappa.Subcommands[MessagesList | MessagesEdit | MessagesDelete]


@cappa.command(name="list", help="List users in the organization")
@dataclass
class CliqUsersList:
    def __call__(self, client: Annotated[ZohoClient, cappa.Dep(get_client)]) -> None:
        url = f"{client.cliq_base}/api/v2/users"
        data = client.request("GET", url)
        output(data)


@cappa.command(name="get", help="Get user details")
@dataclass
class CliqUsersGet:
    user_id: Annotated[str, cappa.Arg(help="User ID or email")]

    def __call__(self, client: Annotated[ZohoClient, cappa.Dep(get_client)]) -> None:
        url = f"{client.cliq_base}/api/v2/users/{self.user_id}"
        data = client.request("GET", url)
        output(data)


@cappa.command(name="users", help="Cliq user operations")
@dataclass
class CliqUsers:
    subcommand: cappa.Subcommands[CliqUsersList | CliqUsersGet]


@cappa.command(name="cliq", help="Zoho Cliq operations")
@dataclass
class Cliq:
    subcommand: cappa.Subcommands[Channels | Chats | Buddies | Messages | CliqUsers]
