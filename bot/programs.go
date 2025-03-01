package bot

import (
	"agent/core"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/nbd-wtf/go-nostr/nip19"
)

// 🔥 Chatter Program (Handles starting and responding)
type ChatterProgram struct {
	MaxRunCount     int
	Peers           []string
	IsRunning       bool
	Leader          bool
	CurrentRunCount int
}

func (p *ChatterProgram) IsActive() bool {
	return p.IsRunning
}

// ShouldRun determines if this bot should handle the message
func (p *ChatterProgram) ShouldRun(message *core.OutgoingMessage) bool {
	if strings.Contains(message.Content, "🧮") && p.Leader {
		return true
	}
	return p.IsRunning // Only active bots continue running
}

// Run executes the chatter program
func (p *ChatterProgram) Run(bot *BaseBot, message *core.OutgoingMessage) string {
	log.Printf("🏃 [%s] [ChatterProgram] [%d]", bot.Name, p.CurrentRunCount)

	if p.CurrentRunCount >= p.MaxRunCount {
		log.Printf("🛑 [%s] [%T] reached max run count. Terminating...", bot.Name, p.Run)
		// bot.RemoveProgram(p)
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

// Responds when mentioned
func (p *ChatterProgram) startToMention(bot *BaseBot, message *core.OutgoingMessage) {
	// ⏳ Add a response delay
	time.Sleep(time.Duration(p.CurrentRunCount) * time.Second)

	// 🔄 Get next bot in sequence
	receiver := bot.GetNextReceiver(p)

	// ✅ Encode public key
	encodedPublicKey, err := nip19.EncodePublicKey(receiver)
	if err != nil {
		log.Printf("❌ Error encoding public key: %v", err)
		return
	}

	reply := &core.OutgoingMessage{
		Content:           core.CreateMessage("@" + encodedPublicKey + " 0"),
		ChannelID:         message.ChannelID,
		ReceiverPublicKey: bot.PublicKey,
	}

	// 📨 Publish message
	bot.EventBus.Publish(core.GroupResponseEvent, reply)
}

//

// **ResponderProgram** allows a bot to respond to mentions.
type ResponderProgram struct {
	MaxRunCount     int
	CurrentRunCount int
	ResponseDelay   int
	IsRunning       bool
	Peers           []string
}

func (p *ResponderProgram) IsActive() bool {
	return p.IsRunning
}

// **ShouldRun** checks if the message should trigger a response.
func (p *ResponderProgram) ShouldRun(message *core.OutgoingMessage) bool {
	return true
}

// **Run** handles a response when mentioned.
func (p *ResponderProgram) Run(bot *BaseBot, message *core.OutgoingMessage) string {
	log.Printf("🏃 [%s] [ResponderProgram] [%d]", bot.Name, p.CurrentRunCount)

	if p.CurrentRunCount >= p.MaxRunCount {
		log.Printf("🛑 [%s] [%T] reached max run count. Terminating...", bot.Name, p.Run)
		// bot.RemoveProgram(p)
		p.IsRunning = false
		return "🔴"
	}

	if !p.IsRunning {
		p.IsRunning = true
		p.CurrentRunCount = 0
	}

	time.Sleep(time.Duration(p.ResponseDelay) * time.Second)

	p.CurrentRunCount++

	// ✅ 1. Extract and verify the mention
	mention := extractMention(message.Content)
	receiver := bot.GetPublicKey() // Consider converting to npub format if necessary

	encodedPublicKey, err := nip19.EncodePublicKey(receiver)
	if err != nil {
		log.Printf("❌ Error encoding public key: %v", err)
		return "🔴"
	}

	if mention == "" || mention != encodedPublicKey {
		return "🟠 No valid mention"
	}

	// ✅ 2. Parse number safely
	words := splitMessageContent(message.Content)
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

	// ⏳ 3. Introduce response delay
	time.Sleep(time.Duration(p.ResponseDelay) * time.Second)

	encodedPublicKey, err = nip19.EncodePublicKey(message.SenderPublicKey)
	if err != nil {
		log.Printf("❌ Error encoding public key: %v", err)
		return "🔴"
	}

	// 🎯 4. Construct the response
	reply := &core.OutgoingMessage{
		Content:           core.CreateMessage("@" + encodedPublicKey + " " + strconv.Itoa(number)),
		ChannelID:         message.ChannelID,
		ReceiverPublicKey: bot.GetPublicKey(),
	}

	// 🚀 5. Publish the response
	bot.EventBus.Publish(core.GroupResponseEvent, reply)
	return "🟢"
}

// 🛠️ Split message into words
func splitMessageContent(content string) []string {
	return strings.Split(content, " ")
}

// Extracts mentions
func extractMention(content string) string {
	words := strings.Split(content, " ")
	for _, word := range words {
		if strings.HasPrefix(word, "@") {
			return word[1:]
		}
	}
	return ""
}
