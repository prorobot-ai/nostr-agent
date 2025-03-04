package programs

import "agent/core"

// **BotProgram** defines a contract for all bot programs
type BotProgram interface {
	Run(bot Bot, message *core.Message) string
	ShouldRun(message *core.Message) bool
	IsActive() bool
}

// **Bot** allows programs to interact with any bot
type Bot interface {
	GetAliases() []string
	GetPublicKey() string
	GetNextReceiver(p *ChatterProgram) string
	Publish(message *core.Message)
}
