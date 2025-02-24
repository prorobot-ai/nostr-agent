package bot

import (
	"agent/core"

	"github.com/nbd-wtf/go-nostr"
)

// Bot defines the essential interface for all bots
type Bot interface {
	Start()                                // Starts the bot
	Stop()                                 // Gracefully stops the bot
	Publish(message *core.OutgoingMessage) // Publishes a message to the relay
	IsReady() bool                         // Checks if the bot is ready to process events
	GetPublicKey() string                  // Returns the bot's public key
	GetSecretKey() string                  // Returns the bot's secret key
	GetRelayURL() string                   // Returns the relay URL
}

// EventListener defines behavior for listening to events
type EventListener interface {
	StartListening(bot *BaseBot)              // Starts listening for events
	ProcessEvent(bot Bot, event *nostr.Event) // Processes a received event
	HandleConnectionLoss(bot Bot)             // Handles relay disconnections
}

// Publisher defines how messages should be published
type Publisher interface {
	Broadcast(bot *BaseBot, message *core.OutgoingMessage) error // Publishes a message
}
