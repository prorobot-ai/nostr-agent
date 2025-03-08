package handlers

import (
	"agent/bot"
	"agent/core"
	"log"
	"strings"
	"time"

	"github.com/nbd-wtf/go-nostr/nip19"
)

type WelcomeHandler struct {
	ChannelID string
	EventBus  *bot.EventBus
}

func (h *WelcomeHandler) Subscribe(eventBus *bot.EventBus) {
	log.Println("✅ Subscribed")
	h.EventBus = eventBus
	h.EventBus.Subscribe(core.DMMessageEvent, h.HandleMessage)
}

func (h *WelcomeHandler) HandleMessage(message *core.Message) {
	text := message.Payload.Text

	switch {
	case strings.Contains(text, "I'm online."):
		npub, _ := nip19.EncodePublicKey(message.ReceiverPublicKey)

		reply := &core.Message{
			ChannelID: h.ChannelID,
			Payload: core.ContentStructure{
				Kind: "message",
				Text: core.CreateContent(npub, "subscriber"),
			},
		}

		time.Sleep(time.Second)
		h.EventBus.Publish(core.GroupResponseEvent, reply)
	}
}
