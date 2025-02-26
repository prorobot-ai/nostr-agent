package handlers

import (
	"agent/bot"
	"agent/core"
	"log"
	"strings"

	"github.com/nbd-wtf/go-nostr/nip19"
)

type ExchangeHandler struct {
	ChannelID        string
	EventBus         *bot.EventBus
	Manager          *bot.BotManager
	Bot              *bot.BaseBot
	encodedPublicKey string
}

// ✅ Subscribe to events
func (h *ExchangeHandler) Subscribe(eventBus *bot.EventBus) {
	if eventBus == nil {
		log.Println("❌ EventBus is not initialized!")
		return
	}
	h.EventBus = eventBus
	h.encodedPublicKey, _ = nip19.EncodePublicKey(h.Bot.PublicKey)

	log.Printf("✅ Subscribed 🚌 [%s]", h.Bot.Name)
	h.EventBus.Subscribe(core.GroupMessageEvent, h.HandleMessage)
}

// 🔄 Forward messages to bot for processing
func (h *ExchangeHandler) HandleMessage(message *core.OutgoingMessage) {
	log.Printf("📩 [%s] Handling Message: %s", h.Bot.Name, message.Content) // ✅ Log every message received

	if strings.Contains(message.Content, "🧮") {
		h.Manager.AssignPrograms()
	}

	// 🚫 Don't process own messages
	if message.SenderPublicKey == h.Bot.PublicKey {
		log.Printf("⏩ [%s] Ignoring its own message.", h.Bot.Name)
		return
	}

	h.Bot.ExecutePrograms(message)
}
