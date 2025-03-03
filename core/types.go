package core

type ContentStructure struct {
	Content string `json:"content"`
	Kind    string `json:"kind"`
}

type Message struct {
	Content string `json:"content"`
	Kind    string `json:"kind"`
}

type OutgoingMessage struct {
	ReceiverPublicKey string `json:"receiver_pub_key,omitempty"`
	SenderPublicKey   string `json:"sender_pub_key,omitempty"`
	ChannelID         string `json:"channel_id,omitempty"` // For group/channel messages
	Content           string `json:"content"`              // The message content
	Timestamp         int64  `json:"timestamp"`            // When the message was created
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

type RemoteJob struct {
	ChannelID string
	Payload   string
}
