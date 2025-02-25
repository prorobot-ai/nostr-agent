package handlers

import (
	"agent/bot"
	"agent/core"
	"fmt"
	"log"
	"strings"
)

type SupportHandler struct {
	EventBus *bot.EventBus
}

func (h *SupportHandler) Subscribe(eventBus *bot.EventBus) {
	log.Println("✅ Subscribed")
	h.EventBus = eventBus
	h.EventBus.Subscribe(core.DMMessageEvent, h.HandleMessage)
}

func (h *SupportHandler) HandleMessage(message *core.OutgoingMessage) {
	switch {
	case strings.Contains(message.Content, "!ping"):
		reply := &core.OutgoingMessage{
			ReceiverPubKey: message.ReceiverPubKey,
			Content:        "🏓 Pong! I'm alive.",
		}
		h.EventBus.Publish(core.DMResponseEvent, reply)

	case strings.Contains(message.Content, "I'm online."):
		reply := &core.OutgoingMessage{
			ReceiverPubKey: message.ReceiverPubKey,
			Content:        "👋 Welcome to Dispatch! Let us know if you need any assistance.",
		}
		h.EventBus.Publish(core.DMResponseEvent, reply)

	case strings.Contains(message.Content, "Hi, I would like to report "):
		reply := &core.OutgoingMessage{
			ReceiverPubKey: message.ReceiverPubKey,
			Content: fmt.Sprintf(
				"Could you elaborate on the problem you're encountering with %s? Additional details would greatly assist in resolving your issue. In the meanwhile, feel free to mute the user if that's necessary.",
				h.ExtractUsername(message.Content),
			),
		}
		h.EventBus.Publish(core.DMResponseEvent, reply)
	}
}

func (h *SupportHandler) ExtractUsername(input string) string {
	input = strings.TrimSuffix(input, ".")
	length := len(input)
	if length > 10 {
		return input[length-10:]
	}
	return input
}
