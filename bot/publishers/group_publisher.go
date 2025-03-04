package publishers

import (
	"agent/bot"
	"agent/bot/handlers"
	"agent/core"
	"log"

	"github.com/nbd-wtf/go-nostr"
)

// GroupPublisher handles sending messages to a group/channel
type GroupPublisher struct {
	ChannelID string
	Handler   *handlers.GroupHandler
}

func (publisher *GroupPublisher) Broadcast(b *bot.BaseBot, message *core.Message) error {

	content := message.Payload.Content

	event := nostr.Event{
		PubKey:    b.PublicKey,
		CreatedAt: nostr.Now(),
		Kind:      nostr.KindChannelMessage,
		Content:   content,
		Tags: nostr.Tags{
			{"e", publisher.ChannelID, b.Relay.URL, "root"},
		},
	}

	event.Sign(b.SecretKey)

	if err := b.Relay.Publish(b.Context, event); err != nil {
		log.Printf("❌ Failed to publish group message: %v", err)
	} else {
		groupID := publisher.ChannelID
		groupID = groupID[len(groupID)-3:]

		log.Printf("[%s] 🗣️ [%s] %s", b.Config.Name, groupID, message.Payload)
	}
	return nil
}
