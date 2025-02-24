package main

import (
	"agent/bot"
	"agent/bot/handlers"
	"agent/bot/listeners"
	"agent/bot/publishers"
	"agent/core"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/docopt/docopt-go"
	"github.com/joho/godotenv"
)

const USAGE = `agent

Usage:
	agent support_bot
	agent weather_bot
	agent welcome_bot

Specify <content> as '-' to make the publish or message command read it
from stdin.
`

func main() {
	initializeLogging()
	loadEnvVariables()

	opts, err := docopt.ParseArgs(USAGE, flag.Args(), "")
	if err != nil {
		log.Fatalf("❌ Failed to parse CLI arguments: %v", err)
	}

	relayURL, nsec, channelID := getEnvVariables()

	// Command Execution
	switch {
	case opts["support_bot"].(bool):
		startSupportBot(relayURL, nsec, channelID)
	case opts["weather_bot"].(bool):
		startWeatherBot(relayURL, nsec, channelID)
	case opts["welcome_bot"].(bool):
		startWelcomeBot(relayURL, nsec, channelID)
	default:
		fmt.Println("❗ Invalid command. Use '--help' for usage instructions.")
	}
}

// ✅ Utility Functions

func initializeLogging() {
	log.SetPrefix("[agent] ")
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
}

func loadEnvVariables() {
	if err := godotenv.Load(); err != nil {
		log.Println("⚠️ Could not load .env file. Using system environment variables...")
	}
}

func getEnvVariables() (string, string, string) {
	relayURL := os.Getenv("DISPATCH_RELAY_URL")
	nsec := os.Getenv("DISPATCH_NSEC")
	channelID := os.Getenv("DISPATCH_CHANNEL_ID")

	if nsec == "" || channelID == "" || relayURL == "" {
		log.Fatal("❌ Missing required environment variables: DISPATCH_RELAY_URL, DISPATCH_NSEC, DISPATCH_CHANNEL_ID")
	}

	return relayURL, nsec, channelID
}

// ✅ Command Execution Function: Starts a basic DM bot
func startSupportBot(relayURL, nsec, channelID string) {
	log.Println("🤖 Starting Direct Message Bot...")

	// 🔄 Initialize EventBus for internal communication
	eventBus := bot.NewEventBus()
	if eventBus == nil {
		log.Fatal("❌ Failed to initialize EventBus")
	} else {
		log.Println("✅ EventBus initialized successfully")
	}

	// 📥 Set up the support handler and subscribe to events
	supportHandler := &handlers.SupportHandler{
		EventBus: eventBus,
	}
	supportHandler.Subscribe()

	// 🎧 Initialize the direct message listener
	listener := &listeners.DMListener{}

	// 📤 Initialize the DM publisher for sending outgoing messages
	publisher := &publishers.DMPublisher{}

	// 🤖 Create the bot instance
	supportBot := bot.NewBaseBot(
		relayURL,
		nsec,
		listener,  // Listens for incoming events
		publisher, // Publishes outgoing messages
		eventBus,  // EventBus for internal communication
	)

	// 🔗 Subscribe to DM responses and broadcast them using the publisher
	eventBus.Subscribe(core.DMResponseEvent, func(message *core.OutgoingMessage) {
		if err := publisher.Broadcast(supportBot, message); err != nil {
			log.Printf("❌ Failed to broadcast message: %v", err)
		}
	})

	// 🚦 Initialize the BotManager for managing concurrent bots
	manager := bot.BotManager{}
	manager.AddBot(supportBot)

	// 🚀 Start all bots concurrently
	manager.StartAll()

	// 🔒 Keep the main thread running to prevent exit
	select {}
}

// ✅ Command Execution Function: Starts a basic Group bot
func startWeatherBot(relayURL, nsec, channelID string) {
	log.Println("🤖 Starting Group Bot...")

	// 🔄 Initialize EventBus for internal communication
	eventBus := bot.NewEventBus()
	if eventBus == nil {
		log.Fatal("❌ Failed to initialize EventBus")
	} else {
		log.Println("✅ EventBus initialized successfully")
	}

	// 📥 Set up the support handler and subscribe to events
	groupHandler := &handlers.GroupHandler{
		ChannelID: channelID,
		EventBus:  eventBus,
	}
	groupHandler.Subscribe()

	// 🎧 Initialize the group listener
	listener := &listeners.GroupListener{
		ChannelID: channelID,
	}

	// 📤 Initialize the Group publisher for sending outgoing messages
	publisher := &publishers.GroupPublisher{
		ChannelID: channelID,
	}

	// 🤖 Create the bot instance
	groupBot := bot.NewBaseBot(
		relayURL,
		nsec,
		listener,  // Listens for incoming events
		publisher, // Publishes outgoing messages
		eventBus,  // EventBus for internal communication
	)

	// 🔗 Subscribe to Group responses and broadcast them using the publisher
	eventBus.Subscribe(core.GroupResponseEvent, func(message *core.OutgoingMessage) {
		if err := publisher.Broadcast(groupBot, message); err != nil {
			log.Printf("❌ Failed to broadcast message: %v", err)
		}
	})

	// 🚦 Initialize the BotManager for managing concurrent bots
	manager := bot.BotManager{}
	manager.AddBot(groupBot)

	// 🚀 Start all bots concurrently
	manager.StartAll()

	// 🔒 Keep the main thread running to prevent exit
	select {}
}

// ✅ Command Execution Function: Starts a basic Group bot
func startWelcomeBot(relayURL, nsec, channelID string) {
	log.Println("🤖 Starting Group Bot...")

	eventBus := bot.NewEventBus()

	welcomeHandler := &handlers.WelcomeHandler{
		ChannelID: channelID,
		EventBus:  eventBus,
	}
	welcomeHandler.Subscribe()

	listener := &listeners.DMListener{}

	publisher := &publishers.GroupPublisher{
		ChannelID: channelID,
	}

	groupBot := bot.NewBaseBot(
		relayURL,
		nsec,
		listener,  // Listens for incoming events
		publisher, // Publishes outgoing messages
		eventBus,  // EventBus for internal communication
	)

	eventBus.Subscribe(core.GroupResponseEvent, func(message *core.OutgoingMessage) {
		if err := publisher.Broadcast(groupBot, message); err != nil {
			log.Printf("❌ Failed to broadcast message: %v", err)
		}
	})

	// 🚦 Initialize the BotManager for managing concurrent bots
	manager := bot.BotManager{}
	manager.AddBot(groupBot)

	// 🚀 Start all bots concurrently
	manager.StartAll()

	// 🔒 Keep the main thread running to prevent exit
	select {}
}
