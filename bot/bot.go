package bot

import (
	"agent/core"
	"context"
	"log"

	"github.com/nbd-wtf/go-nostr"
	"github.com/nbd-wtf/go-nostr/nip19"
)

// BaseBot implements the core bot functionalities
type BaseBot struct {
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
}

// NewBaseBot initializes a new instance of BaseBot
func NewBaseBot(relayURL, nsec string, listener EventListener, publisher Publisher, eventBus *EventBus) *BaseBot {
	ctx, cancel := context.WithCancel(context.Background())

	_, sk, _ := nip19.Decode(nsec)
	pk, _ := nostr.GetPublicKey(sk.(string))
	npub, _ := nip19.EncodePublicKey(pk)

	return &BaseBot{
		RelayURL:         relayURL,
		SecretKey:        nsec,
		PublicKey:        npub,
		Context:          ctx,
		CancelFunc:       cancel,
		IsActiveListener: false,
		Listener:         listener,
		Publisher:        publisher,
		EventBus:         eventBus,
	}
}

// Starts the bot
func (bot *BaseBot) Start() {
	err := bot.connectToRelay()
	if err != nil {
		log.Fatalf("❌ Failed to connect: %v", err)
	}
	bot.Listener.StartListening(bot)
}

// Connects to the relay
func (bot *BaseBot) connectToRelay() error {
	relay, err := nostr.RelayConnect(bot.Context, bot.RelayURL)
	if err != nil {
		return err
	}

	bot.Relay = relay
	log.Println("✅ Connected to relay")
	return nil
}

// Stops the bot gracefully
func (bot *BaseBot) Stop() {
	bot.CancelFunc()
	if bot.Relay != nil {
		bot.Relay.Close()
	}
	log.Println("🛑 Bot stopped gracefully")
}

// Publishes a message using the Publisher interface
func (bot *BaseBot) Publish(message *core.OutgoingMessage) {
	bot.Publisher.Broadcast(bot, message)
}

// Bot interface implementations
func (bot *BaseBot) GetPublicKey() string {
	return bot.PublicKey
}

func (bot *BaseBot) GetSecretKey() string {
	return bot.SecretKey
}

func (bot *BaseBot) GetRelayURL() string {
	return bot.RelayURL
}

func (bot *BaseBot) IsReady() bool {
	return bot.IsActiveListener
}
