package core

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

// Notifier interface allows multiple notification mechanisms (WebSocket, Logs, etc.)
type Notifier interface {
	SendMessage(message SocketRequest)
	Close()
}

// SocketRequest represents a standard WebSocket message format
type SocketRequest struct {
	Type      string `json:"type"`
	ChannelID string `json:"channel_id"`
	Metadata  string `json:"metadata"`
	Text      string `json:"text"`
	CreatedAt int64  `json:"created_at"`
}

// WebSocketNotifier implements Notifier for WebSockets
type WebSocketNotifier struct {
	conn *websocket.Conn
}

// NewWebSocketNotifier initializes a WebSocket notifier
func NewWebSocketNotifier(wsURL string) (*WebSocketNotifier, error) {
	timestamp := time.Now().Unix()

	session := fmt.Sprintf("%d", timestamp)

	conn, _, err := websocket.DefaultDialer.Dial(wsURL+"?session="+session, nil)
	if err != nil {
		return nil, err
	}

	log.Println("✅ WebSocket connection established:", wsURL)
	return &WebSocketNotifier{conn: conn}, nil
}

// SendMessage sends a message over WebSocket
func (w *WebSocketNotifier) SendMessage(message SocketRequest) {
	jsonMessage, _ := json.Marshal(message)

	err := w.conn.WriteMessage(websocket.TextMessage, jsonMessage)
	if err != nil {
		log.Printf("❌ Failed to send WebSocket message: %v", err)
	}
}

// Close closes the WebSocket connection
func (w *WebSocketNotifier) Close() {
	if w.conn != nil {
		log.Println("🔴 Closing WebSocket connection")
		w.conn.Close()
	}
}

// LoggerNotifier is a fallback notifier that logs messages instead of WebSocket
type LoggerNotifier struct{}

// SendMessage logs the message when WebSocket isn't available
func (l *LoggerNotifier) SendMessage(message SocketRequest) {
	log.Printf("📢 LOG NOTIFICATION [%s]: %s", message.Type, message.Text)
}

// Close does nothing for LoggerNotifier
func (l *LoggerNotifier) Close() {}
