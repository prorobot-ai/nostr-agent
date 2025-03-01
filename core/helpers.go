package core

import (
	"encoding/json"
	"log"
)

// 🛠️ Convert structured content to JSON string
func CreateContent(text string, kind string) string {
	message := ContentStructure{
		Content: text,
		Kind:    kind,
	}

	jsonData, err := json.Marshal(message)
	if err != nil {
		log.Println("❌ Error marshalling JSON:", err)
		return ""
	}

	return string(jsonData)
}

// 🛠️ Convert structured content to JSON string
func CreateMessage(text string) string {
	message := ContentStructure{
		Content: text,
		Kind:    "message",
	}

	jsonData, err := json.Marshal(message)
	if err != nil {
		log.Println("❌ Error marshalling JSON:", err)
		return ""
	}

	return string(jsonData)
}
