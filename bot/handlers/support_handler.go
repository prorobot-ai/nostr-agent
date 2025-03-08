package handlers

import (
	"agent/bot"
	"agent/core"
	"fmt"
	"log"
	"strings"
	"time"
)

type SupportHandler struct {
	EventBus *bot.EventBus
}

func (h *SupportHandler) Subscribe(eventBus *bot.EventBus) {
	log.Println("✅ Subscribed")
	h.EventBus = eventBus
	h.EventBus.Subscribe(core.DMMessageEvent, h.HandleMessage)
}

func (h *SupportHandler) HandleMessage(message *core.Message) {
	text := message.Payload.Text

	switch {
	case strings.Contains(text, "!ping"):
		reply := &core.Message{
			ReceiverPublicKey: message.ReceiverPublicKey,
			Payload: core.ContentStructure{
				Kind: "message",
				Text: "🏓 Pong! I'm alive.",
			},
		}
		time.Sleep(time.Second)
		h.EventBus.Publish(core.DMResponseEvent, reply)

	case strings.Contains(text, "I'm online."):
		reply := &core.Message{
			ReceiverPublicKey: message.ReceiverPublicKey,
			Payload: core.ContentStructure{
				Kind: "message",
				Text: "👋 Welcome to Dispatch! Let us know if you need any assistance.",
			},
		}
		time.Sleep(time.Second)
		h.EventBus.Publish(core.DMResponseEvent, reply)

	case strings.Contains(text, "Hi, I would like to report "):
		reply := &core.Message{
			ReceiverPublicKey: message.ReceiverPublicKey,
			Payload: core.ContentStructure{
				Kind: "message",
				Text: fmt.Sprintf(
					"Could you elaborate on the problem you're encountering with %s? Additional details would greatly assist in resolving your issue. In the meanwhile, feel free to mute the user if that's necessary.",
					h.ExtractUsername(text),
				),
			},
		}
		time.Sleep(time.Second)
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
