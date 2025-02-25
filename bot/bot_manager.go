package bot

type BotManager struct {
	Bots         []Bot
	currentIndex int
}

func (m *BotManager) AddBot(botInstance Bot) {
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
