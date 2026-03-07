//go:build integration

package internal_test

import (
	"fmt"
	"os"
	"strings"
	"testing"
	"time"
)

func (c *testCleanup) trackCliqChannel(id string) {
	c.add("delete cliq channel "+id, func() {
		zohoIgnoreError(c.t, "cliq", "channels", "delete", id)
	})
}

func TestCliqUsers(t *testing.T) {
	t.Parallel()

	const knownUserEmail = "jasmin@miaie.com"
	const knownUserID = "913284317"

	t.Run("list", func(t *testing.T) {
		out := zoho(t, "cliq", "users", "list")
		m := parseJSON(t, out)
		data, ok := m["data"].([]any)
		if !ok {
			t.Fatalf("expected data array in users list response:\n%s", truncate(out, 500))
		}
		if len(data) == 0 {
			t.Fatal("expected at least one user")
		}
		found := false
		for _, item := range data {
			user, _ := item.(map[string]any)
			if fmt.Sprintf("%v", user["email_id"]) == knownUserEmail {
				found = true
				assertEqual(t, fmt.Sprintf("%v", user["id"]), knownUserID)
				break
			}
		}
		if !found {
			t.Errorf("known user %s not found in users list", knownUserEmail)
		}
	})

	t.Run("get-by-id", func(t *testing.T) {
		out := zoho(t, "cliq", "users", "get", knownUserID)
		m := parseJSON(t, out)
		data, ok := m["data"].(map[string]any)
		if !ok {
			t.Fatalf("expected data object in users get response:\n%s", truncate(out, 500))
		}
		assertEqual(t, fmt.Sprintf("%v", data["id"]), knownUserID)
		assertStringField(t, data, "email_id", knownUserEmail)
		assertStringField(t, data, "status", "active")
	})

	t.Run("get-by-email", func(t *testing.T) {
		out := zoho(t, "cliq", "users", "get", knownUserEmail)
		m := parseJSON(t, out)
		data, ok := m["data"].(map[string]any)
		if !ok {
			t.Fatalf("expected data object in users get response:\n%s", truncate(out, 500))
		}
		assertEqual(t, fmt.Sprintf("%v", data["id"]), knownUserID)
		assertStringField(t, data, "email_id", knownUserEmail)
	})
}

func TestCliqChannels(t *testing.T) {
	t.Parallel()
	cleanup := newCleanup(t)

	var channelID string
	var uniqueName string
	var chatID string
	channelName := fmt.Sprintf("%s_chan_%s", testPrefix, randomSuffix())

	t.Run("create", func(t *testing.T) {
		out := zoho(t, "cliq", "channels", "create",
			"--name", channelName, "--level", "private",
			"--description", "Integration test channel")
		m := parseJSON(t, out)
		channelID = fmt.Sprintf("%v", m["channel_id"])
		uniqueName = fmt.Sprintf("%v", m["unique_name"])
		chatID = fmt.Sprintf("%v", m["chat_id"])
		if channelID == "" || channelID == "<nil>" {
			t.Fatalf("expected channel_id in create response:\n%s", truncate(out, 500))
		}
		cleanup.trackCliqChannel(channelID)
		assertStringField(t, m, "status", "created")
		assertStringField(t, m, "level", "private")
		if chatID == "" || chatID == "<nil>" || chatID == "null" {
			t.Fatalf("expected real chat_id for private channel, got %q", chatID)
		}
		name := fmt.Sprintf("%v", m["name"])
		assertContains(t, strings.ToLower(name), strings.ToLower(channelName))
		t.Logf("created channel %s (unique_name=%s, chat_id=%s)", channelID, uniqueName, chatID)
	})

	t.Run("list", func(t *testing.T) {
		requireID(t, channelID, "create must have succeeded")
		out := zoho(t, "cliq", "channels", "list")
		arr := parseJSONArray(t, out)
		assertNonEmpty(t, arr, "expected at least one channel")
		found := false
		for _, ch := range arr {
			if fmt.Sprintf("%v", ch["channel_id"]) == channelID {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("created channel %s not found in list", channelID)
		}
	})

	t.Run("get", func(t *testing.T) {
		requireID(t, uniqueName, "create must have succeeded")
		out := zoho(t, "cliq", "channels", "get", uniqueName)
		m := parseJSON(t, out)
		assertStringField(t, m, "type", "channel")
		data, ok := m["data"].(map[string]any)
		if !ok {
			t.Fatalf("expected data object in channel get response:\n%s", truncate(out, 500))
		}
		assertEqual(t, fmt.Sprintf("%v", data["channel_id"]), channelID)
	})

	t.Run("members", func(t *testing.T) {
		requireID(t, channelID, "create must have succeeded")
		out := zoho(t, "cliq", "channels", "members", channelID)
		m := parseJSON(t, out)
		members, ok := m["members"].([]any)
		if !ok {
			t.Fatalf("expected members array in response:\n%s", truncate(out, 500))
		}
		if len(members) == 0 {
			t.Fatal("expected at least one member (creator)")
		}
		foundCreator := false
		for _, member := range members {
			mem, _ := member.(map[string]any)
			if fmt.Sprintf("%v", mem["email_id"]) == "jasmin@miaie.com" {
				foundCreator = true
				break
			}
		}
		if !foundCreator {
			t.Error("creator jasmin@miaie.com not found in channel members")
		}
	})
}

func TestCliqBuddiesMessage(t *testing.T) {
	t.Parallel()

	const recipientEmail = "uday@miaie.com"
	const selfEmail = "jasmin@miaie.com"

	t.Run("send", func(t *testing.T) {
		msgText := fmt.Sprintf("%s buddy test %s", testPrefix, randomSuffix())
		out := zoho(t, "cliq", "buddies", "message", recipientEmail,
			"--text", msgText)
		m := parseJSON(t, out)
		assertEqual(t, fmt.Sprintf("%v", m["ok"]), "true")
		assertStringField(t, m, "email", recipientEmail)
	})

	t.Run("self-rejected", func(t *testing.T) {
		r := runZoho(t, "cliq", "buddies", "message", selfEmail,
			"--text", "should fail")
		if r.ExitCode == 0 {
			t.Fatal("expected non-zero exit code for self-message")
		}
		assertContains(t, r.Stderr+r.Stdout, "buddies_self_message_restricted")
	})
}

func TestCliqMessages(t *testing.T) {
	t.Parallel()
	cleanup := newCleanup(t)

	var chatID string
	var channelID string
	var msgID string
	channelName := fmt.Sprintf("%s_msg_%s", testPrefix, randomSuffix())
	originalText := fmt.Sprintf("%s msg %s", testPrefix, randomSuffix())
	editedText := fmt.Sprintf("%s edited %s", testPrefix, randomSuffix())

	t.Run("setup", func(t *testing.T) {
		out := zoho(t, "cliq", "channels", "create",
			"--name", channelName, "--level", "private")
		m := parseJSON(t, out)
		channelID = fmt.Sprintf("%v", m["channel_id"])
		chatID = fmt.Sprintf("%v", m["chat_id"])
		if channelID == "" || channelID == "<nil>" {
			t.Fatalf("expected channel_id in create response:\n%s", truncate(out, 500))
		}
		cleanup.trackCliqChannel(channelID)
		if chatID == "" || chatID == "<nil>" || chatID == "null" {
			t.Fatalf("expected real chat_id for private channel, got %q", chatID)
		}
		t.Logf("created channel %s with chat_id %s", channelID, chatID)
	})

	t.Run("send", func(t *testing.T) {
		requireID(t, chatID, "setup must have succeeded")
		out := zoho(t, "cliq", "chats", "message", chatID,
			"--text", originalText)
		m := parseJSON(t, out)
		assertEqual(t, fmt.Sprintf("%v", m["ok"]), "true")
		t.Logf("sent message to chat %s", chatID)
	})

	t.Run("list", func(t *testing.T) {
		requireID(t, chatID, "setup must have succeeded")
		retryUntil(t, 15*time.Second, func() bool {
			out := zoho(t, "cliq", "messages", "list", chatID, "--limit", "10")
			m := parseJSON(t, out)
			data, ok := m["data"].([]any)
			if !ok || len(data) == 0 {
				return false
			}
			for _, item := range data {
				msg, _ := item.(map[string]any)
				content, _ := msg["content"].(map[string]any)
				if content != nil && fmt.Sprintf("%v", content["text"]) == originalText {
					msgID = fmt.Sprintf("%v", msg["id"])
					return true
				}
			}
			return false
		})
		if msgID == "" {
			t.Fatal("sent message not found in messages list")
		}
		t.Logf("found message ID %s", msgID)
	})

	t.Run("edit", func(t *testing.T) {
		requireID(t, msgID, "list must have found the message")
		zoho(t, "cliq", "messages", "edit", chatID, msgID,
			"--text", editedText)
	})

	t.Run("list-after-edit", func(t *testing.T) {
		requireID(t, msgID, "list must have found the message")
		retryUntil(t, 15*time.Second, func() bool {
			out := zoho(t, "cliq", "messages", "list", chatID, "--limit", "10")
			m := parseJSON(t, out)
			data, ok := m["data"].([]any)
			if !ok {
				return false
			}
			for _, item := range data {
				msg, _ := item.(map[string]any)
				if fmt.Sprintf("%v", msg["id"]) == msgID {
					content, _ := msg["content"].(map[string]any)
					if content != nil && fmt.Sprintf("%v", content["text"]) == editedText {
						return true
					}
				}
			}
			return false
		})
	})

	t.Run("delete", func(t *testing.T) {
		requireID(t, msgID, "list must have found the message")
		zoho(t, "cliq", "messages", "delete", chatID, msgID)
	})

	t.Run("list-after-delete", func(t *testing.T) {
		requireID(t, msgID, "list must have found the message")
		retryUntil(t, 15*time.Second, func() bool {
			out := zoho(t, "cliq", "messages", "list", chatID, "--limit", "10")
			m := parseJSON(t, out)
			data, ok := m["data"].([]any)
			if !ok {
				return false
			}
			for _, item := range data {
				msg, _ := item.(map[string]any)
				if fmt.Sprintf("%v", msg["id"]) == msgID {
					content, _ := msg["content"].(map[string]any)
					if content == nil {
						return true
					}
					return fmt.Sprintf("%v", content["text"]) == ""
				}
			}
			return false
		})
	})
}

func TestCliqEmergencyCleanup(t *testing.T) {
	if os.Getenv("ZOHO_EMERGENCY_CLEANUP") == "" {
		t.Skip("set ZOHO_EMERGENCY_CLEANUP=1 to run")
	}
	out := zoho(t, "cliq", "channels", "list")
	arr := parseJSONArray(t, out)
	for _, ch := range arr {
		name := fmt.Sprintf("%v", ch["name"])
		uniqueName := fmt.Sprintf("%v", ch["unique_name"])
		if !strings.Contains(name, testPrefix) && !strings.Contains(uniqueName, testPrefix) {
			continue
		}
		channelID := fmt.Sprintf("%v", ch["channel_id"])
		if channelID == "" || channelID == "<nil>" {
			continue
		}
		t.Logf("cleaning orphaned channel %s (%s)", channelID, name)
		zohoIgnoreError(t, "cliq", "channels", "delete", channelID)
	}
}


