package listeners

import (
	"agent/bot"
	"agent/core"
	"encoding/json"
	"log"

	"github.com/nbd-wtf/go-nostr"
	"github.com/nbd-wtf/go-nostr/nip04"
	"github.com/nbd-wtf/go-nostr/nip19"
)

// DMListener handles direct message events
type DMListener struct{}

// StartListening starts listening for direct messages
func (listener *DMListener) StartListening(b *bot.BaseBot) {
	relay := b.Relay
	filters := listener.Filters(b)

	sub, err := relay.Subscribe(b.Context, filters)
	if err != nil {
		log.Printf("❌ Subscription failed: %v", err)
		return
	}
	defer sub.Unsub()

	var storedEvents []*nostr.Event
	processingStoredEvents := false

	for {
		select {
		case event, ok := <-sub.Events:
			if !ok {
				log.Println("🚫 Subscription closed, reconnecting...")
				b.Relay.Close()
				// return fmt.Errorf("subscription closed")
			}

			if !processingStoredEvents {
				storedEvents = append(storedEvents, event)
			} else if b.IsActiveListener {
				listener.ProcessEvent(b, event)
			}

		case <-sub.EndOfStoredEvents:
			if !processingStoredEvents {
				log.Println("📥 Processing pending events...")
				for i := len(storedEvents) - 1; i >= 0; i-- {
					// bot.handleEvent(storedEvents[i])
				}
				storedEvents = nil
				processingStoredEvents = true
				b.IsActiveListener = true
				log.Printf("👂 [%s] listening ", b.Config.Name)
			}
		case <-relay.Context().Done():
			listener.HandleConnectionLoss(b)
			return
		}
	}
}

// ProcessEvent handles incoming direct message events
func (listener *DMListener) ProcessEvent(b *bot.BaseBot, event *nostr.Event) {
	// 🔑 Decrypt the incoming message
	shared, _ := nip04.ComputeSharedSecret(event.PubKey, b.SecretKey)
	npub, _ := nip19.EncodePublicKey(event.PubKey)

	plaintext, err := nip04.Decrypt(event.Content, shared)
	if err != nil {
		log.Printf("❌ Decryption failed: %v", err)
		return
	}

	log.Printf("🔓 Decrypted message: %s", plaintext)

	var message core.ContentStructure
	if err := json.Unmarshal([]byte(plaintext), &message); err != nil {
		log.Printf("❌ Failed to unmarshal message: %v", err)
		return
	}

	log.Printf("💬 [DM from %s]: %s", npub, message.Text)

	// 📩 Pass the event to EventBus
	b.EventBus.Publish(core.DMMessageEvent, &core.Message{
		ReceiverPublicKey: event.PubKey,
		Payload:           message,
	})
}

func (listener *DMListener) Filters(b *bot.BaseBot) []nostr.Filter {
	tags := map[string][]string{"p": {b.PublicKey}}

	return []nostr.Filter{
		{
			Kinds: []int{nostr.KindEncryptedDirectMessage},
			Tags:  tags,
			Limit: 50,
		},
	}
}

// HandleConnectionLoss handles relay disconnections
func (listener *DMListener) HandleConnectionLoss(bot *bot.BaseBot) {
	log.Println("🔄 Reconnecting DM Listener...")
	bot.Start() // Attempt to restart the bot
}
