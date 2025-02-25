package handlers

import (
	"agent/bot"
	"agent/core"
	"agent/services/weather"
	"log"
	"strings"
)

type GroupHandler struct {
	ChannelID string
	EventBus  *bot.EventBus
}

func (h *GroupHandler) Subscribe(eventBus *bot.EventBus) {
	log.Println("✅ Subscribed")
	h.EventBus = eventBus
	h.EventBus.Subscribe(core.GroupMessageEvent, h.HandleMessage)
}

func (h *GroupHandler) HandleMessage(message *core.OutgoingMessage) {
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
