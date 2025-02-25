package publishers

import (
	"agent/bot"
	"agent/bot/handlers"
	"agent/core"
	"log"

	"github.com/nbd-wtf/go-nostr"
	"github.com/nbd-wtf/go-nostr/nip19"
)

// GroupPublisher handles sending messages to a group/channel
type GroupPublisher struct {
	ChannelID string
	Handler   *handlers.GroupHandler
}

func (publisher *GroupPublisher) Broadcast(b *bot.BaseBot, message *core.OutgoingMessage) error {
	event := nostr.Event{
		PubKey:    b.GetPublicKey(),
		CreatedAt: nostr.Now(),
		Kind:      nostr.KindChannelMessage,
		Content:   message.Content,
		Tags: nostr.Tags{
			{"e", publisher.ChannelID, b.Relay.URL, "root"},
		},
	}

	event.Sign(b.SecretKey)

	if err := b.Relay.Publish(b.Context, event); err != nil {
		log.Printf("❌ Failed to publish group message: %v", err)
	} else {
		npub, _ := nip19.EncodePublicKey(b.PublicKey)
		ID := npub[len(npub)-4:]
		log.Printf("📢 Group message sent to channel %s %s", publisher.ChannelID, ID)
	}

	log.Printf("✉️ Message sent to %s: %s", publisher.ChannelID, message.Content)
	return nil
}
