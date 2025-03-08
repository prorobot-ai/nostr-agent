package programs

import (
	"agent/core"
	"log"
	"strings"
	"time"

	"github.com/nbd-wtf/go-nostr/nip19"
)

// **ChatterProgram** - Handles starting and responding
type ChatterProgram struct {
	IsRunning       bool
	CurrentRunCount int

	ProgramConfig core.ProgramConfig

	Peers []string

	Leader bool
}

// ✅ **Check if the program is active**
func (p *ChatterProgram) IsActive() bool {
	return p.IsRunning
}

// ✅ **Determine if this should run**
func (p *ChatterProgram) ShouldRun(message *core.Message) bool {
	text := message.Payload.Text
	return (strings.Contains(text, "🧮") && p.Leader) || p.IsRunning
}

// ✅ **Run Chatter Logic**
func (p *ChatterProgram) Run(bot Bot, message *core.Message) string {
	log.Printf("🏃 [%s] [ChatterProgram] [%d]", bot.GetPublicKey(), p.CurrentRunCount)

	if p.CurrentRunCount >= p.ProgramConfig.MaxRunCount {
		log.Printf("🛑 [%s] [ChatterProgram] reached max run count. Terminating...", bot.GetPublicKey())
		p.IsRunning = false
		return "🔴"
	}

	if !p.IsRunning {
		p.IsRunning = true
		p.CurrentRunCount = 0
	}

	time.Sleep(time.Duration(p.CurrentRunCount) * time.Second)

	p.CurrentRunCount++

	p.startToMention(bot, message)

	return "🟢"
}

// ✅ **Mention the next bot**
func (p *ChatterProgram) startToMention(bot Bot, message *core.Message) {
	time.Sleep(time.Duration(p.CurrentRunCount) * time.Second)

	receiver := bot.GetNextReceiver(p)

	encodedPublicKey, err := nip19.EncodePublicKey(receiver)
	if err != nil {
		log.Printf("❌ Error encoding public key: %v", err)
		return
	}

	reply := &core.Message{
		ChannelID:         message.ChannelID,
		ReceiverPublicKey: bot.GetPublicKey(),

		Payload: core.ContentStructure{
			Kind: "message",
			Text: core.CreateContent("@"+encodedPublicKey+" 0", "message"),
		},
	}

	bot.Publish(reply)
}
