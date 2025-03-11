package programs

import "agent/core"

// **BotProgram** defines a contract for all bot programs
type BotProgram interface {
	Run(bot Bot, message *core.BusMessage) string
	ShouldRun(message *core.BusMessage) bool
	IsActive() bool
}

// **Bot** allows programs to interact with any bot
type Bot interface {
	GetName() string
	GetAliases() []string
	GetPublicKey() string
	GetNextReceiver(p *ChatterProgram) string
	Publish(message *core.BusMessage)
}
