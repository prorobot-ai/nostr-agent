package plugins

import (
	"agent/bot"
	"fmt"

	"github.com/nbd-wtf/go-nostr"
)

type LoggingPlugin struct{}

// OnGlobalEvent logs every event received by the bot
func (l *LoggingPlugin) OnEvent(bot bot.Bot, event *nostr.Event) {
	fmt.Printf("📜 [Log] Event received from %s with content: %s\n", event.PubKey, event.Content)
}
