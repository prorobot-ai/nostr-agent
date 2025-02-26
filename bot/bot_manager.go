package bot

import (
	"log"
)

type BotManager struct {
	Bots         []*BaseBot
	currentIndex int
}

func (m *BotManager) AddBot(botInstance *BaseBot) {
	m.Bots = append(m.Bots, botInstance)
}

// 🔄 Get the next bot mention in round-robin fashion
func (m *BotManager) GetNextBotMention() string {
	if len(m.Bots) == 0 {
		return ""
	}
	nextIndex := (m.currentIndex + 1) % len(m.Bots)
	nextBot := m.Bots[nextIndex]
	m.currentIndex = nextIndex
	return nextBot.GetPublicKey()
}

func (m *BotManager) GetPeers(exclude string) []string {
	var peers []string
	for _, b := range m.Bots {
		if b.GetPublicKey() != exclude { // Avoid adding itself
			peers = append(peers, b.GetPublicKey())
		}
	}
	return peers
}

func (m *BotManager) StartAll() {
	for _, bot := range m.Bots {
		go bot.Start()
	}
}

// Dynamically assigns different programs to bots
func (m *BotManager) AssignPrograms() {
	// First, gather all bot public keys
	var allPeers []string
	for _, bot := range m.Bots {
		allPeers = append(allPeers, bot.GetPublicKey())
		bot.ResetPrograms()
	}

	// Now assign programs to bots
	for _, bot := range m.Bots {
		if bot.Name == "Yin" {
			log.Printf("🛠 Assigning ChatterProgram to [%s]", bot.Name)
			bot.AddProgram(&ChatterProgram{
				MaxRunCount:     1,
				CurrentRunCount: 0,
				Leader:          true,
				Peers:           filterPeers(allPeers, bot.GetPublicKey()),
			})
		} else if bot.Name == "Yang" {
			log.Printf("🛠 Assigning ResponderProgram to [%s]", bot.Name)
			bot.AddProgram(&ResponderProgram{
				MaxRunCount:     10,
				CurrentRunCount: 0,
				ResponseDelay:   1,
				Peers:           filterPeers(allPeers, bot.GetPublicKey()),
			})
		}
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
