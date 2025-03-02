package bot

import (
	"agent/bot/programs"
	"log"
)

type BotManager struct {
	Bots     []*BaseBot
	Programs map[*BaseBot][]programs.BotProgram
}

func NewBotManager() *BotManager {
	return &BotManager{
		Bots:     []*BaseBot{},
		Programs: make(map[*BaseBot][]programs.BotProgram),
	}
}

func (m *BotManager) AddBot(bot *BaseBot) {
	m.Bots = append(m.Bots, bot)
}

func (m *BotManager) StartAll() {
	m.AssignPrograms()

	for _, bot := range m.Bots {
		go bot.Start()
	}
}

func (m *BotManager) InitializePrograms(bot *BaseBot) {
	buffer := []programs.BotProgram{}

	bot.ResetPrograms()

	var allPeers []string
	for _, b := range m.Bots {
		allPeers = append(allPeers, b.PublicKey)
	}

	if bot.Config.Name == "Yin" {
		log.Printf("🛠 Assigning ChatterProgram to [%s]", bot.Config.Name)
		buffer = append(buffer, &programs.ChatterProgram{
			MaxRunCount:     1,
			CurrentRunCount: 0,
			Leader:          true,
			Peers:           filterPeers(allPeers, bot.PublicKey),
		})
	} else if bot.Config.Name == "Yang" {
		log.Printf("🛠 Assigning ResponderProgram to [%s]", bot.Config.Name)
		buffer = append(buffer, &programs.ResponderProgram{
			MaxRunCount:   10,
			ResponseDelay: 1,
			Peers:         filterPeers(allPeers, bot.PublicKey),
		})
	} else if bot.Config.Name == "HypeWizard" {
		log.Printf("🛠 Assigning ConductorProgram to [%s]", bot.Config.Name)

		conductor := &programs.ConductorProgram{
			MaxRunCount:   10,
			ResponseDelay: 1,
			Url:           bot.Config.ProgramConfig.Url,
			Peers:         filterPeers(allPeers, bot.PublicKey),
		}

		conductor.InitCrawlerClient(bot.Config.ProgramConfig.Address)
		buffer = append(buffer, conductor)
	}

	bot.AssignPrograms(buffer)
	m.Programs[bot] = buffer
}

func (m *BotManager) AssignPrograms() {
	for _, bot := range m.Bots {
		m.InitializePrograms(bot)
	}
}

// **filterPeers** removes the bot's own public key from the peer list
func filterPeers(peers []string, exclude string) []string {
	var filtered []string
	for _, peer := range peers {
		if peer != exclude {
			filtered = append(filtered, peer)
		}
	}
	return filtered
}
