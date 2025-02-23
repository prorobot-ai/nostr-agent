package handlers

import (
	"agent/bot"
	"agent/core"
	"fmt"
	"strings"
)

type SupportHandler struct {
	Plugins []bot.HandlerPlugin // Handler-specific plugins
}

func (h *SupportHandler) HandleMessage(bot bot.Bot, message core.Message, senderPubKey string) {
	switch {
	case message.Content == "🙂":
		bot.PublishEncrypted(senderPubKey, "🙃")

	case strings.Contains(message.Content, "report"):
		reply := fmt.Sprintf("Could you elaborate on the problem you're encountering?")
		bot.PublishEncrypted(senderPubKey, reply)

	case message.Content == "I'm online.":
		welcomeMessage := "👋 Welcome to Dispatch! Let us know if you need any assistance."
		bot.PublishEncrypted(senderPubKey, welcomeMessage)

		// 🔥 Trigger plugins directly for notifications
		h.TriggerPlugins(bot, message, senderPubKey)
	}
}

// TriggerPlugins allows handler plugins to react to specific events
func (h *SupportHandler) TriggerPlugins(bot bot.Bot, message core.Message, senderPubKey string) {
	for _, plugin := range h.Plugins {
		plugin.OnTrigger(bot, message, senderPubKey)
	}
}
