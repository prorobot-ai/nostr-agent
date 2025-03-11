package handlers

import (
	"agent/bot"
	"agent/core"
	"agent/services/weather"
	"log"
	"strings"
	"time"
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

func (h *GroupHandler) HandleMessage(message *core.BusMessage) {
	text := message.Payload.Text

	switch {
	case strings.Contains(text, "!weather"):
		weatherReport := weather.GetReport()

		reply := &core.BusMessage{
			ChannelID: h.ChannelID,
			Payload: core.ContentStructure{
				Kind: "message",
				Text: core.SerializeContent(weatherReport, "message"),
			},
		}

		time.Sleep(time.Second)
		h.EventBus.Publish(core.GroupResponseEvent, reply)
	}
}
