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
	manager := bot.BotManager{}

	// Dynamically start bots based on YAML configuration
	for _, botCfg := range botConfigs.Bots {
		startDynamicBot(botCfg, &manager)
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
func startDynamicBot(cfg core.BotConfig, manager *bot.BotManager) {
	log.Printf("🤖 Starting bot: %s...", cfg.Name)

	eventBus := bot.NewEventBus()
	if eventBus == nil {
		log.Fatalf("❌ Failed to initialize EventBus for %s", cfg.Name)
	}

	listener := resolveListener(cfg.Listener, cfg.ChannelID)
	publisher := resolvePublisher(cfg.Publisher, cfg.ChannelID)

	// Initialize the bot
	botInstance := bot.NewBaseBot(
		cfg.RelayURL,
		cfg.Nsec,
		listener,
		publisher,
		eventBus,
	)

	botInstance.SetName(cfg.Name)
	botInstance.SetAliases(cfg.Aliases)

	handler := resolveHandler(
		cfg.Handler,
		cfg.ChannelID,
		manager,
		botInstance,
	)

	manager.AddBot(botInstance)

	// Subscribe to relevant events
	handler.Subscribe(eventBus)

	eventBus.Subscribe(resolveEventType(cfg.EventType), func(message *core.OutgoingMessage) {
		if err := publisher.Broadcast(botInstance, message); err != nil {
			log.Printf("❌ [%s] Failed to broadcast message: %v", cfg.Name, err)
		}
	})
}

//////////////////////////////////////////////////////////////////////////////////////
// ✅ Dynamic Resolver Functions
//////////////////////////////////////////////////////////////////////////////////////

func resolveListener(listenerType, channelID string) bot.EventListener {
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

func resolvePublisher(publisherType, channelID string) bot.Publisher {
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

func resolveHandler(handlerType, channelID string, manager *bot.BotManager, botInstance *bot.BaseBot) bot.EventHandler {
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

func resolveEventType(eventType string) core.EventType {
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
