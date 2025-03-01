package publishers

import (
	"agent/bot"
	"agent/core"
	"context"
	"log"

	"github.com/nbd-wtf/go-nostr"
	"github.com/nbd-wtf/go-nostr/nip04"
)

// DMPublisher handles sending encrypted direct messages (DM) between bots
type DMPublisher struct{}

// Publish sends an encrypted direct message to the receiver
func (publisher *DMPublisher) Broadcast(b *bot.BaseBot, message *core.OutgoingMessage) error {
	receiverPubKey := message.ReceiverPublicKey

	sk := b.SecretKey
	// Compute the shared secret
	shared, err := nip04.ComputeSharedSecret(receiverPubKey, sk)
	if err != nil {
		log.Printf("❌ Failed to compute shared secret: %v", err)
		return err
	}

	// Encrypt the message
	encryptedMessage, err := nip04.Encrypt(message.Content, shared)
	if err != nil {
		log.Printf("❌ Failed to encrypt message: %v", err)
		return err
	}

	// Create the event
	ev := nostr.Event{
		PubKey:    b.PublicKey,
		CreatedAt: nostr.Now(),
		Kind:      nostr.KindEncryptedDirectMessage,
		Content:   encryptedMessage,
		Tags:      nostr.Tags{{"p", receiverPubKey}},
	}

	// Sign the event
	if err := ev.Sign(sk); err != nil {
		log.Printf("❌ Failed to sign event: %v", err)
		return err
	}

	// Publish the message via the relay
	relay := b.Relay
	if relay == nil {
		log.Printf("❌ Relay connection is not established")
		return err
	}

	if err := relay.Publish(context.Background(), ev); err != nil {
		log.Printf("❌ Failed to publish DM: %v", err)
		return err
	}

	log.Printf("✉️ DM sent to %s: %s", receiverPubKey, message.Content)
	return nil
}
