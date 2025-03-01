package bot

import (
	"agent/core"
	"context"
	"log"
	"sync"

	"github.com/nbd-wtf/go-nostr"
	"github.com/nbd-wtf/go-nostr/nip19"
)

// BaseBot implements the core bot functionalities
type BaseBot struct {
	mu sync.Mutex

	Name             string
	Aliases          []string
	RelayURL         string
	SecretKey        string
	PublicKey        string
	Context          context.Context
	CancelFunc       context.CancelFunc
	Relay            *nostr.Relay
	IsActiveListener bool
	Listener         EventListener
	Publisher        Publisher
	EventBus         *EventBus
	Programs         []BotProgram
}

// NewBaseBot initializes a new instance of BaseBot
func NewBaseBot(relayURL, nsec string, listener EventListener, publisher Publisher, eventBus *EventBus) *BaseBot {
	ctx, cancel := context.WithCancel(context.Background())

	_, sk, _ := nip19.Decode(nsec)
	pk, _ := nostr.GetPublicKey(sk.(string))

	return &BaseBot{
		RelayURL:         relayURL,
		SecretKey:        sk.(string),
		PublicKey:        pk,
		Context:          ctx,
		CancelFunc:       cancel,
		IsActiveListener: false,
		Listener:         listener,
		Publisher:        publisher,
		EventBus:         eventBus,
	}
}

// Starts the bot
func (b *BaseBot) Start() {
	err := b.connectToRelay()
	if err != nil {
		log.Fatalf("❌ Failed to connect: %v", err)
	}
	b.Listener.StartListening(b)
}

func (b *BaseBot) SetName(name string) {
	b.Name = name
}

func (b *BaseBot) SetAliases(aliases []string) {
	b.Aliases = aliases
}

// Connects to the relay
func (b *BaseBot) connectToRelay() error {
	relay, err := nostr.RelayConnect(b.Context, b.RelayURL)
	if err != nil {
		return err
	}

	b.Relay = relay
	// 📝 Check if there are aliases before logging them
	if len(b.Aliases) > 0 {
		log.Printf("✅ [%s] %v Connected 📡 Aliases: %s", b.Name, b.Aliases, b.RelayURL)
	} else {
		log.Printf("✅ [%s] Connected 📡 [%s]", b.Name, b.RelayURL)
	}
	return nil
}

// Stops the bot gracefully
func (b *BaseBot) Stop() {
	b.CancelFunc()
	if b.Relay != nil {
		b.Relay.Close()
	}
	log.Println("🛑 Bot stopped gracefully")
}

// Publishes a message using the Publisher interface
func (b *BaseBot) Publish(message *core.OutgoingMessage) {
	b.Publisher.Broadcast(b, message)
}

// Bot interface implementations
func (b *BaseBot) GetPublicKey() string {
	return b.PublicKey
}

func (b *BaseBot) GetSecretKey() string {
	return b.SecretKey
}

func (b *BaseBot) IsReady() bool {
	return b.IsActiveListener
}

// AssignProgram adds a new program to the bot
func (bot *BaseBot) AddProgram(program BotProgram) {
	bot.mu.Lock()
	defer bot.mu.Unlock()

	bot.Programs = append(bot.Programs, program)
	log.Printf("🚀 [%s] Assigned new program", bot.Name)
}

// ExecutePrograms runs all active programs for a bot
func (bot *BaseBot) ExecutePrograms(message *core.OutgoingMessage) {
	bot.mu.Lock()
	defer bot.mu.Unlock()

	log.Printf("⚙️ [%s] Executing [%d] Programs on message: %s", bot.Name, len(bot.Programs), message.Content)

	var programsToRemove []int

	for i, program := range bot.Programs {
		if program.ShouldRun(message) {
			result := program.Run(bot, message)
			log.Printf("[%s] [%s]", bot.Name, result)

			if !program.IsActive() {
				programsToRemove = append(programsToRemove, i)
			}
		}
	}

	// ✅ Now remove all completed programs
	for i := len(programsToRemove) - 1; i >= 0; i-- {
		index := programsToRemove[i]
		log.Printf("🗑️ [%s] Removing completed program: %T", bot.Name, bot.Programs[index])
		bot.Programs = append(bot.Programs[:index], bot.Programs[index+1:]...)
	}
}

// GetNextReceiver picks the next peer for a given program
func (bot *BaseBot) GetNextReceiver(program *ChatterProgram) string {
	// bot.mu.Lock()
	// defer bot.mu.Unlock()

	if len(program.Peers) == 0 {
		log.Println("⚠️ No peers available.")
		return ""
	}

	index := program.CurrentRunCount % len(program.Peers)
	return program.Peers[index]
}

func (bot *BaseBot) ResetPrograms() {
	bot.Programs = []BotProgram{}
}

// BotProgram defines the interface for all programs a bot can execute
type BotProgram interface {
	Run(b *BaseBot, message *core.OutgoingMessage) string // 🚀 What should the bot do?
	ShouldRun(message *core.OutgoingMessage) bool         // 🔄 When should this program activate?
	IsActive() bool
}
