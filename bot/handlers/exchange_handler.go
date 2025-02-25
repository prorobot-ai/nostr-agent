package handlers

import (
	"agent/bot"
	"agent/core"
	"log"
	"strconv"
	"strings"
	"time"
)

type ExchangeHandler struct {
	ChannelID string
	EventBus  *bot.EventBus
}

func (h *ExchangeHandler) Subscribe(eventBus *bot.EventBus) {
	log.Println("✅ Subscribed")
	h.EventBus = eventBus
	h.EventBus.Subscribe(core.GroupMessageEvent, h.HandleMessage)
}

func (h *ExchangeHandler) HandleMessage(message *core.OutgoingMessage) {
	switch {
	case strings.Contains(message.Content, "@"+message.ReceiverPubKey):
		words := h.splitMessageContent(message.Content)

		number, err := strconv.Atoi(words[1])
		if err != nil {
			return
		}

		number++

		// ⏳ Introduce a delay (e.g., 2 seconds)
		time.Sleep(1 * time.Second)

		reply := &core.OutgoingMessage{
			Content:   "@" + message.ReceiverPubKey + " " + strconv.Itoa(number),
			ChannelID: h.ChannelID,
		}
		h.EventBus.Publish(core.GroupResponseEvent, reply)
	}
}

func (h *ExchangeHandler) splitMessageContent(content string) []string {
	words := strings.Split(content, " ")
	return words
}
