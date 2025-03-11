package programs

import (
	"agent/core"
	"log"
	"strconv"
	"time"

	"github.com/nbd-wtf/go-nostr/nip19"
)

// **ResponderProgram** - Handles responding when mentioned
type ResponderProgram struct {
	IsRunning       bool
	CurrentRunCount int

	ProgramConfig core.ProgramConfig

	Peers []string
}

// ✅ **Check if the program is active**
func (p *ResponderProgram) IsActive() bool {
	return p.IsRunning
}

// ✅ **Should this program run?**
func (p *ResponderProgram) ShouldRun(message *core.BusMessage) bool {
	return true
}

// ✅ **Run Responder Logic**
func (p *ResponderProgram) Run(bot Bot, message *core.BusMessage) string {
	log.Printf("🏃 [%s] [ResponderProgram] [%d]", bot.GetPublicKey(), p.CurrentRunCount)

	if p.CurrentRunCount >= p.ProgramConfig.MaxRunCount {
		log.Printf("🛑 [%s] [ResponderProgram] reached max run count. Terminating...", bot.GetPublicKey())
		p.IsRunning = false
		return "🔴"
	}

	if !p.IsRunning {
		p.IsRunning = true
		p.CurrentRunCount = 0
	}

	p.CurrentRunCount++

	text := message.Payload.Text

	mention := core.ExtractMention(text)
	receiver := bot.GetPublicKey()

	encodedPublicKey, err := nip19.EncodePublicKey(receiver)
	if err != nil {
		log.Printf("❌ Error encoding public key: %v", err)
		return "🔴"
	}

	if mention == "" || mention != encodedPublicKey {
		return "🟠 No valid mention"
	}

	words := core.SplitMessageContent(text)
	if len(words) < 2 {
		log.Println("⚠️ Malformed message, missing number.")
		return "🟠"
	}

	number, err := strconv.Atoi(words[1])
	if err != nil {
		log.Println("❌ Could not parse number:", err)
		return "🟠"
	}
	number++

	time.Sleep(time.Duration(p.ProgramConfig.ResponseDelay) * time.Second)

	encodedPublicKey, err = nip19.EncodePublicKey(message.SenderPublicKey)
	if err != nil {
		log.Printf("❌ Error encoding public key: %v", err)
		return "🔴"
	}

	reply := &core.BusMessage{
		ChannelID:         message.ChannelID,
		ReceiverPublicKey: bot.GetPublicKey(),

		Payload: core.ContentStructure{
			Kind: "message",
			Text: core.SerializeContent("@"+encodedPublicKey+" "+strconv.Itoa(number), "message"),
		},
	}

	bot.Publish(reply)
	return "🟢"
}
