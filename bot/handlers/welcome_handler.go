package handlers

import (
	"agent/bot"
	"agent/core"
	"encoding/json"
	"fmt"
	"log"
	"strings"

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

func (h *WelcomeHandler) HandleMessage(message *core.OutgoingMessage) {
	switch {
	case strings.Contains(message.Content, "I'm online."):
		npub, _ := nip19.EncodePublicKey(message.ReceiverPublicKey)
		reply := &core.OutgoingMessage{
			Content:   h.createMessage(npub),
			ChannelID: h.ChannelID,
		}
		h.EventBus.Publish(core.GroupResponseEvent, reply)
	}
}

func (h *WelcomeHandler) createMessage(content string) string {
	message := core.ContentStructure{
		Content: content,
		Kind:    "subscriber",
	}

	jsonData, err := json.Marshal(message)
	if err != nil {
		fmt.Println("Error marshalling JSON:", err)
		return ""
	}

	return string(jsonData)
}
