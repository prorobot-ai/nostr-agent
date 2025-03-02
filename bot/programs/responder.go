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
	MaxRunCount     int
	CurrentRunCount int
	ResponseDelay   int
	IsRunning       bool
	Peers           []string
}

// ✅ **Check if the program is active**
func (p *ResponderProgram) IsActive() bool {
	return p.IsRunning
}

// ✅ **Should this program run?**
func (p *ResponderProgram) ShouldRun(message *core.OutgoingMessage) bool {
	return true
}

// ✅ **Run Responder Logic**
func (p *ResponderProgram) Run(bot Bot, message *core.OutgoingMessage) string {
	log.Printf("🏃 [%s] [ResponderProgram] [%d]", bot.GetPublicKey(), p.CurrentRunCount)

	if p.CurrentRunCount >= p.MaxRunCount {
		log.Printf("🛑 [%s] [ResponderProgram] reached max run count. Terminating...", bot.GetPublicKey())
		p.IsRunning = false
		return "🔴"
	}

	if !p.IsRunning {
		p.IsRunning = true
		p.CurrentRunCount = 0
	}

	p.CurrentRunCount++

	mention := core.ExtractMention(message.Content)
	receiver := bot.GetPublicKey()

	encodedPublicKey, err := nip19.EncodePublicKey(receiver)
	if err != nil {
		log.Printf("❌ Error encoding public key: %v", err)
		return "🔴"
	}

	if mention == "" || mention != encodedPublicKey {
		return "🟠 No valid mention"
	}

	words := core.SplitMessageContent(message.Content)
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

	time.Sleep(time.Duration(p.ResponseDelay) * time.Second)

	encodedPublicKey, err = nip19.EncodePublicKey(message.SenderPublicKey)
	if err != nil {
		log.Printf("❌ Error encoding public key: %v", err)
		return "🔴"
	}

	reply := &core.OutgoingMessage{
		Content:           core.CreateContent("@"+encodedPublicKey+" "+strconv.Itoa(number), "message"),
		ChannelID:         message.ChannelID,
		ReceiverPublicKey: bot.GetPublicKey(),
	}

	bot.Publish(reply)
	return "🟢"
}
