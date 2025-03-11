package main

import (
	"agent/bot"
	"agent/bot/handlers"
	"agent/bot/listeners"
	"agent/bot/publishers"
	"agent/core"
	"flag"
	"log"
)

func main() {
	initializeLogging()

	// Parse command-line flags
	configFile := flag.String("config", "", "Path to YAML configuration file for the bot")
	flag.Parse()

	if *configFile == "" {
		log.Fatal("❌ No configuration file provided. Use '--config=your_bot.yaml'")
	}

	// Load the bot configuration from YAML
	botConfigs, err := core.LoadBotConfigs(*configFile)
	if err != nil {
		log.Fatalf("❌ Could not load bot configuration: %v", err)
	}

	// Initialize the shared BotManager
	manager := bot.NewBotManager()

	// Dynamically start bots based on YAML configuration
	for _, botCfg := range botConfigs.Bots {
		startDynamicBot(botCfg, manager)
	}

	// Start all bots concurrently
	manager.StartAll()

	// Keep the program running
	select {}
}

// 🔄 Set up logging format
func initializeLogging() {
	log.SetPrefix("[agent] ")
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
}

// 🚀 Dynamically initialize and start a bot based on config
func startDynamicBot(config core.BotConfig, manager *bot.BotManager) {
	log.Printf("🤖 Starting bot: %s...", config.Name)

	eventBus := bot.NewEventBus()
	if eventBus == nil {
		log.Fatalf("❌ Failed to initialize EventBus for %s", config.Name)
	}

	listener := initializeListener(config.Listener, config.ChannelID)
	publisher := initializePublisher(config.Publisher, config.ChannelID)

	// Initialize the bot
	bot := bot.NewBaseBot(
		config,
		listener,
		publisher,
		eventBus,
	)

	handler := initializeHandler(
		config.Handler,
		config.ChannelID,
		manager,
		bot,
	)

	manager.AddBot(bot)

	handler.Subscribe(eventBus)

	eventBus.Subscribe(getEventType(config.EventType), func(message *core.BusMessage) {
		if err := publisher.Broadcast(bot, message); err != nil {
			log.Printf("❌ [%s] Failed to broadcast message: %v", config.Name, err)
		}
	})
}

//////////////////////////////////////////////////////////////////////////////////////
// ✅ Dynamic Resolver Functions
//////////////////////////////////////////////////////////////////////////////////////

func initializeListener(listenerType, channelID string) bot.EventListener {
	switch listenerType {
	case "DMListener":
		return &listeners.DMListener{}
	case "GroupListener":
		return &listeners.GroupListener{ChannelID: channelID}
	default:
		log.Fatalf("❌ Unknown listener type: %s", listenerType)
		return nil
	}
}

func initializePublisher(publisherType, channelID string) bot.Publisher {
	switch publisherType {
	case "DMPublisher":
		return &publishers.DMPublisher{}
	case "GroupPublisher":
		return &publishers.GroupPublisher{ChannelID: channelID}
	default:
		log.Fatalf("❌ Unknown publisher type: %s", publisherType)
		return nil
	}
}

func initializeHandler(handlerType, channelID string, manager *bot.BotManager, botInstance *bot.BaseBot) bot.EventHandler {
	switch handlerType {
	case "ExchangeHandler":
		return &handlers.ExchangeHandler{
			ChannelID: channelID,
			Manager:   manager,
			Bot:       botInstance,
		}
	case "SupportHandler":
		return &handlers.SupportHandler{}
	case "GroupHandler":
		return &handlers.GroupHandler{ChannelID: channelID}
	case "WelcomeHandler":
		return &handlers.WelcomeHandler{ChannelID: channelID}
	default:
		log.Fatalf("❌ Unknown handler type: %s", handlerType)
		return nil
	}
}

func getEventType(eventType string) core.EventType {
	switch eventType {
	case "DMResponseEvent":
		return core.DMResponseEvent
	case "GroupResponseEvent":
		return core.GroupResponseEvent
	default:
		log.Fatalf("❌ Unknown event type: %s", eventType)
		return ""
	}
}
