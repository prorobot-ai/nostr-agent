package core

import "github.com/nbd-wtf/go-nostr"

type ContentStructure struct {
	Content string `json:"content"`
	Kind    string `json:"kind"`
}

type DefaultProvider struct {
	Relay      *nostr.Relay
	ChannelId  string
	PublicKey  string
	PrivateKey string
}

type Message struct {
	Content string `json:"content"`
	Kind    string `json:"kind"`
}

type OutgoingMessage struct {
	ReceiverPubKey string `json:"receiver_pub_key,omitempty"` // For direct messages (DMs)
	ChannelID      string `json:"channel_id,omitempty"`       // For group/channel messages
	Content        string `json:"content"`                    // The message content
	Timestamp      int64  `json:"timestamp"`                  // When the message was created
}

type RelayProvider interface {
	GetRelay() *nostr.Relay
	GetChannelId() string
	GetPrivateKey() string
	GetPublicKey() string
}

// EventType defines a type for all supported event types
type EventType string

// Enum-like constants for Event Types
const (
	DMMessageEvent     EventType = "dm_message"
	DMResponseEvent    EventType = "dm_response"
	GroupMessageEvent  EventType = "group_message"
	GroupResponseEvent EventType = "group_response"
	BotStartedEvent    EventType = "bot_started"
	BotStoppedEvent    EventType = "bot_stopped"
)
