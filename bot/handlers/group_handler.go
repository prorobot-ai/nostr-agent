package handlers

import (
	"agent/bot"
	"agent/core"
	"agent/services/weather"
	"log"
	"strings"
)

// GroupHandler handles group chat commands
type GroupHandler struct {
	ChannelID string
	EventBus  *bot.EventBus
}

func (h *GroupHandler) Subscribe() {
	log.Println("✅ Subscribed")
	h.EventBus.Subscribe(core.GroupMessageEvent, h.respondToMessage)
}

func (h *GroupHandler) respondToMessage(message *core.OutgoingMessage) {
	switch {
	case strings.Contains(message.Content, "!weather"):
		weatherReport := weather.GetReport()

		reply := &core.OutgoingMessage{
			Content:   weatherReport,
			ChannelID: h.ChannelID,
		}
		h.EventBus.Publish(core.GroupResponseEvent, reply)
	}
}
