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

func (h *GroupHandler) HandleMessage(message *core.Message) {
	content := message.Payload.Content

	switch {
	case strings.Contains(content, "!weather"):
		weatherReport := weather.GetReport()

		reply := &core.Message{
			ChannelID: h.ChannelID,
			Payload: core.ContentStructure{
				Kind:    "message",
				Content: core.CreateContent(weatherReport, "message"),
			},
		}
		h.EventBus.Publish(core.GroupResponseEvent, reply)
	}
}
