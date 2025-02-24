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

func (h *WelcomeHandler) Subscribe() {
	log.Println("✅ Subscribed")
	h.EventBus.Subscribe(core.DMMessageEvent, h.welcomeHandler)
}

func (h *WelcomeHandler) welcomeHandler(message *core.OutgoingMessage) {
	switch {
	case strings.Contains(message.Content, "I’m online."):
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
