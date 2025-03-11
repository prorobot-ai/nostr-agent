package bot

import (
	"agent/core"

	"github.com/nbd-wtf/go-nostr"
)

// EventListener defines behavior for listening to events
type EventListener interface {
	StartListening(bot *BaseBot)                   // Starts listening for events
	ProcessEvent(bot *BaseBot, event *nostr.Event) // Processes a received event
	HandleConnectionLoss(bot *BaseBot)             // Handles relay disconnections
}

type EventHandler interface {
	Subscribe(eventBus *EventBus)           // Subscribes to specific events
	HandleMessage(message *core.BusMessage) // Processes incoming messages
}

// Publisher defines how messages should be published
type Publisher interface {
	Broadcast(bot *BaseBot, message *core.BusMessage) error // Publishes a message
}
