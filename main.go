package main

import (
	"agent/bot"
	"agent/bot/handlers"
	"agent/bot/plugins"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/docopt/docopt-go"
	"github.com/joho/godotenv"
)

const USAGE = `agent

Usage:
  agent basic_bot

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
	case opts["basic_bot"].(bool):
		startBot(relayURL, nsec, channelID)
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
	relayURL := os.Getenv("BOT_1_RELAY_URL")
	nsec := os.Getenv("BOT_1_NSEC")
	channelID := os.Getenv("BOT_1_CHANNEL_ID")

	if nsec == "" || channelID == "" || relayURL == "" {
		log.Fatal("❌ Missing required environment variables: BOT_1_RELAY_URL, BOT_1_NSEC, BOT_1_CHANNEL_ID")
	}

	return relayURL, nsec, channelID
}

// ✅ Command Execution Functions
func startBot(relayURL string, nsec string, channelID string) {
	log.Println("🤖 Starting Direct Message Bot...")

	loggingPlugin := &plugins.LoggingPlugin{}

	// Create a channel notifier plugin
	channelNotifier := &plugins.ChannelNotifierPlugin{
		RelayURL:  relayURL,
		ChannelID: channelID,
	}

	// Create a support handler
	supportHandler := &handlers.SupportHandler{
		Plugins: []bot.HandlerPlugin{channelNotifier},
	}

	// Initialize the support bot with the notifier plugin
	supportBot := bot.NewBasicBot(
		relayURL,
		nsec,
		supportHandler,
		[]bot.BotPlugin{loggingPlugin},
	)

	// Initialize the BotManager to handle concurrent bots if needed
	manager := bot.BotManager{}
	manager.AddBot(supportBot)

	// Start all bots
	manager.StartAll()

	// Block main thread to keep the bots running
	select {}
}
