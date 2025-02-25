package handlers

import (
	"agent/bot"
	"agent/core"
	"encoding/json"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/nbd-wtf/go-nostr/nip19"
)

type ExchangeHandler struct {
	ChannelID     string
	EventBus      *bot.EventBus
	Manager       *bot.BotManager
	Bot           *bot.BaseBot
	Peers         []string
	CurrentIndex  int
	ExchangeCount int
	Chatter       core.ChatterConfig
	IsActive      bool
}

// ✅ Subscribe to events
func (h *ExchangeHandler) Subscribe(eventBus *bot.EventBus) {
	log.Println("✅ Subscribed to Group Messages")
	h.EventBus = eventBus
	h.EventBus.Subscribe(core.GroupMessageEvent, h.HandleMessage)
}

// 🚀 Dynamically load peers from BotManager
func (h *ExchangeHandler) LoadPeers(exclude string) {
	h.Peers = h.Manager.GetPeers(exclude)
	log.Printf("📡 Peers loaded dynamically: %v", h.Peers)
}

// 🔄 Handle incoming messages and manage exchanges
func (h *ExchangeHandler) HandleMessage(message *core.OutgoingMessage) {
	if h.ExchangeCount >= h.Chatter.MaxExchanges {
		log.Println("🔇 Maximum exchanges reached, stopping chatter.")
		h.IsActive = false
		return
	}

	mention := h.extractMention(message.Content)

	switch {
	case strings.Contains(message.Content, "🧮") && h.Chatter.Leader:
		if h.IsActive {
			log.Println("🔄 Chatter already active.")
			return
		}
		log.Println("👑 Leader initiating chatter!")
		h.IsActive = true
		h.ExchangeCount = 0
		h.CurrentIndex = 0

		currentBotPubKey := h.Bot.GetPublicKey()

		// 🔥 Load peers dynamically excluding itself
		h.LoadPeers(currentBotPubKey)

		time.Sleep(time.Duration(h.Chatter.InitialDelay) * time.Second)
		h.startChatter()

	case h.IsActive || h.wasMentioned(mention):
		// Activate the bot if it's mentioned for the first time
		if !h.IsActive {
			log.Printf("🔔 %s is now active!", mention)
			h.IsActive = true
		}

		npub, _ := nip19.EncodePublicKey(message.ReceiverPublicKey)
		if npub != mention {
			return
		}

		words := h.splitMessageContent(message.Content)
		number, err := strconv.Atoi(words[1])
		if err != nil {
			return
		}

		number++

		// ⏳ Add a response delay
		time.Sleep(time.Duration(h.Chatter.ResponseDelay) * time.Second)

		npub, _ = nip19.EncodePublicKey(message.SenderPublicKey)

		reply := &core.OutgoingMessage{
			Content:           h.createMessage("@" + npub + " " + strconv.Itoa(number)),
			ChannelID:         h.ChannelID,
			ReceiverPublicKey: h.Bot.PublicKey,
		}
		h.ExchangeCount++
		h.EventBus.Publish(core.GroupResponseEvent, reply)
	}
}

func (h *ExchangeHandler) createMessage(text string) string {
	message := core.ContentStructure{
		Content: text,
		Kind:    "message",
	}

	jsonData, err := json.Marshal(message)
	if err != nil {
		fmt.Println("Error marshalling JSON:", err)
		return ""
	}

	return string(jsonData)
}

// 🚀 Starts the initial chatter
func (h *ExchangeHandler) startChatter() {
	if len(h.Peers) > 0 {
		receiverPublicKey := h.Manager.GetNextBotMention()
		npub, _ := nip19.EncodePublicKey(receiverPublicKey)

		startMessage := &core.OutgoingMessage{
			Content:           h.createMessage("@" + npub + " 0"),
			ChannelID:         h.ChannelID,
			ReceiverPublicKey: h.Bot.PublicKey,
		}
		h.EventBus.Publish(core.GroupResponseEvent, startMessage)
	}
}

// 🔍 Extract mention dynamically using regex
func (h *ExchangeHandler) extractMention(content string) string {
	re := regexp.MustCompile(`@([a-zA-Z0-9]+)`)
	match := re.FindStringSubmatch(content)
	if len(match) > 1 {
		return match[1]
	}
	return ""
}

func (h *ExchangeHandler) wasMentioned(mention string) bool {
	npub, _ := nip19.EncodePublicKey(h.Bot.PublicKey)
	log.Printf("⚪️ %s", mention)
	return mention == npub
}

// 🛠️ Utility: Split message for processing
func (h *ExchangeHandler) splitMessageContent(content string) []string {
	return strings.Split(content, " ")
}
