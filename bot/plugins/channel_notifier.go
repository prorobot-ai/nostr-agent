package plugins

import (
	"agent/bot"
	"agent/core"
	"fmt"
)

type ChannelNotifierPlugin struct {
	RelayURL  string
	ChannelID string
}

// OnTrigger announces when specific message content is detected
func (c *ChannelNotifierPlugin) OnTrigger(bot bot.Bot, message core.Message, senderPubKey string) {
	if message.Content == "I'm online." {
		fmt.Printf("🔔 Channel %s: User %s has come online!\n", c.ChannelID, senderPubKey)
		core.Announce(
			c.RelayURL,
			c.ChannelID,
			senderPubKey,
			bot.GetSecretKey(),
			bot.GetPublicKey(),
		)
	}
}
