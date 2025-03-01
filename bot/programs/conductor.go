package programs

import (
	"agent/core"
	"log"
	"time"
)

// **ConductorProgram** - Handles responding when mentioned
type ConductorProgram struct {
	MaxRunCount     int
	CurrentRunCount int
	ResponseDelay   int
	IsRunning       bool
	Peers           []string
}

// ✅ **Check if the program is active**
func (p *ConductorProgram) IsActive() bool {
	return p.IsRunning
}

// ✅ **Should this program run?**
func (p *ConductorProgram) ShouldRun(message *core.OutgoingMessage) bool {
	return true
}

// ✅ **Run Responder Logic**
func (p *ConductorProgram) Run(bot Bot, message *core.OutgoingMessage) string {
	log.Printf("🏃 [%s] [ConductorProgram] [%d]", bot.GetPublicKey(), p.CurrentRunCount)

	if p.CurrentRunCount >= p.MaxRunCount {
		log.Printf("🛑 [%s] [ConductorProgram] reached max run count. Terminating...", bot.GetPublicKey())
		p.IsRunning = false
		return "🔴"
	}

	if !p.IsRunning {
		p.IsRunning = true
		p.CurrentRunCount = 0
	}

	p.CurrentRunCount++

	mention := core.ExtractMention(message.Content)
	aliases := bot.GetAliases()
	set := createSet(aliases)

	if mention == "" || !set[mention] {
		return "🟠 No valid mention"
	}

	words := core.SplitMessageContent(message.Content)
	if len(words) < 2 {
		log.Println("⚠️ Malformed message, missing number.")
		return "🟠"
	}

	time.Sleep(time.Duration(p.ResponseDelay) * time.Second)

	reply := &core.OutgoingMessage{
		Content:           core.CreateContent("🧙🏻‍♂️ "+words[1]+" ⚡️", "message"),
		ChannelID:         message.ChannelID,
		ReceiverPublicKey: bot.GetPublicKey(),
	}

	bot.Publish(reply)
	return "🟢"
}
func createSet(arr []string) map[string]bool {
	set := make(map[string]bool)
	for _, v := range arr {
		set[v] = true
	}
	return set
}
