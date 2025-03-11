package bot

import (
	"agent/bot/programs"
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

	Config   core.BotConfig
	Programs []programs.BotProgram

	Context    context.Context
	CancelFunc context.CancelFunc

	RelayURL         string
	SecretKey        string
	PublicKey        string
	IsActiveListener bool

	Relay *nostr.Relay

	Listener  EventListener
	Publisher Publisher
	EventBus  *EventBus
}

// NewBaseBot initializes a new instance of BaseBot
func NewBaseBot(config core.BotConfig, listener EventListener, publisher Publisher, eventBus *EventBus) *BaseBot {
	ctx, cancel := context.WithCancel(context.Background())

	_, sk, _ := nip19.Decode(config.Nsec)
	pk, _ := nostr.GetPublicKey(sk.(string))

	return &BaseBot{
		Config:           config,
		RelayURL:         config.RelayURL,
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

// Connects to the relay
func (b *BaseBot) connectToRelay() error {
	relay, err := nostr.RelayConnect(b.Context, b.RelayURL)
	if err != nil {
		return err
	}

	b.Relay = relay
	// 📝 Check if there are aliases before logging them
	if len(b.Config.Aliases) > 0 {
		log.Printf("📡 [%s] %v connected to [%s] ✅ ", b.Config.Name, b.Config.Aliases, b.RelayURL)
	} else {
		log.Printf("📡 [%s] connected to [%s] ✅ ", b.Config.Name, b.RelayURL)
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

// ============================================================
// 🎛️ Bot interface implementations
// ============================================================
func (b *BaseBot) GetName() string {
	return b.Config.Name
}

func (b *BaseBot) GetAliases() []string {
	return b.Config.Aliases
}

func (b *BaseBot) GetPublicKey() string {
	return b.PublicKey
}

// GetNextReceiver picks the next peer for a given program
func (bot *BaseBot) GetNextReceiver(program *programs.ChatterProgram) string {
	// bot.mu.Lock()
	// defer bot.mu.Unlock()

	if len(program.Peers) == 0 {
		log.Println("⚠️ No peers available.")
		return ""
	}

	index := program.CurrentRunCount % len(program.Peers)
	return program.Peers[index]
}

// Publishes a message using the Publisher interface
func (b *BaseBot) Publish(message *core.BusMessage) {
	b.Publisher.Broadcast(b, message)
}

func (b *BaseBot) IsReady() bool {
	return b.IsActiveListener
}

func (bot *BaseBot) AssignPrograms(p []programs.BotProgram) {
	bot.mu.Lock()
	defer bot.mu.Unlock()

	// ✅ Expand slice `p` into individual elements
	bot.Programs = append(bot.Programs, p...)

	log.Printf("🧮 [%s] received [%d] programs ✅", bot.Config.Name, len(p))
}

func (bot *BaseBot) RemoveProgram(p programs.BotProgram) {
	bot.mu.Lock()
	defer bot.mu.Unlock()

	for i, program := range bot.Programs {
		if program == p {
			bot.Programs = append(bot.Programs[:i], bot.Programs[i+1:]...)
			log.Printf("🗑️ [%s] Removed completed program: %T", bot.Config.Name, p)
			break
		}
	}
}

func (bot *BaseBot) ResetPrograms() {
	bot.Programs = []programs.BotProgram{}
}

// ExecutePrograms runs all active programs for a bot
func (bot *BaseBot) ExecutePrograms(message *core.BusMessage) {
	bot.mu.Lock()
	defer bot.mu.Unlock()

	log.Printf("⚙️ [%s] Executing [%d] Programs on message: %s", bot.Config.Name, len(bot.Programs), message.Payload)

	var programsToRemove []int

	for i, program := range bot.Programs {
		if program.ShouldRun(message) {
			result := program.Run(bot, message)
			log.Printf("[%s] [%s]", bot.Config.Name, result)

			if !program.IsActive() {
				programsToRemove = append(programsToRemove, i)
			}
		}
	}

	// ✅ Now remove all completed programs
	for i := len(programsToRemove) - 1; i >= 0; i-- {
		index := programsToRemove[i]
		log.Printf("🗑️ [%s] Removing completed program: %T", bot.Config.Name, bot.Programs[index])
		bot.Programs = append(bot.Programs[:index], bot.Programs[index+1:]...)
	}
}
