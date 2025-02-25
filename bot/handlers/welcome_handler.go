package handlers

import (
	"agent/bot"
	"agent/core"
	"encoding/json"
	"fmt"
	"log"
	"strings"
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

func (h *WelcomeHandler) HandleMessage(message *core.OutgoingMessage) {
	switch {
	case strings.Contains(message.Content, "I'm online."):
		reply := &core.OutgoingMessage{
			Content:   h.createMessage(message.ReceiverPubKey),
			ChannelID: h.ChannelID,
		}
		h.EventBus.Publish(core.GroupResponseEvent, reply)
	}
}

func (h *WelcomeHandler) createMessage(npub string) string {
	message := core.ContentStructure{
		Content: npub,
		Kind:    "subscriber",
	}

	jsonData, err := json.Marshal(message)
	if err != nil {
		fmt.Println("Error marshalling JSON:", err)
		return ""
	}

	return string(jsonData)
}
