package listeners

import (
	"agent/bot"
	"agent/core"
	"encoding/json"
	"log"

	"github.com/nbd-wtf/go-nostr"
	"github.com/nbd-wtf/go-nostr/nip19"
)

// GroupListener handles group channel events
type GroupListener struct {
	ChannelID string
}

// StartListening subscribes to group channel events
func (listener *GroupListener) StartListening(b *bot.BaseBot) {
	relay := b.Relay
	filters := listener.Filters(b)

	sub, err := relay.Subscribe(b.Context, filters)
	if err != nil {
		log.Printf("❌ Group subscription failed: %v", err)
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
				log.Println("🚀 Entered active listening mode")
			}
		case <-relay.Context().Done():
			listener.HandleConnectionLoss(b)
			return
		}
	}
}

// ProcessEvent handles group channel messages
func (listener *GroupListener) ProcessEvent(b bot.Bot, event *nostr.Event) {
	var message core.Message
	if err := json.Unmarshal([]byte(event.Content), &message); err != nil {
		log.Printf("Failed to unmarshal message: %v", err)
		return
	}

	// Send event to EventBus for inter-bot communication
	if baseBot, ok := b.(*bot.BaseBot); ok && baseBot.EventBus != nil {
		baseBot.EventBus.Publish(core.GroupMessageEvent, &core.OutgoingMessage{
			ChannelID:         listener.ChannelID,
			ReceiverPublicKey: b.GetPublicKey(),
			SenderPublicKey:   event.PubKey,
			Content:           message.Content,
		})
	}

	npub, _ := nip19.EncodePublicKey(b.GetPublicKey())
	ID := npub[len(npub)-4:]
	log.Printf("🎧 [Group] [%s] [%s]: %s ", ID, listener.ChannelID, message.Content)
}

func (listener *GroupListener) Filters(b bot.Bot) []nostr.Filter {
	return []nostr.Filter{
		{
			Kinds: []int{nostr.KindChannelMessage},
			Tags:  map[string][]string{"e": {listener.ChannelID}},
			Limit: 100,
		},
	}
}

// HandleConnectionLoss reconnects the bot
func (listener *GroupListener) HandleConnectionLoss(bot bot.Bot) {
	log.Println("🔄 Reconnecting Group Listener...")
	bot.Start()
}
