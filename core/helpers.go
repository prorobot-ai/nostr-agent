package core

import (
	"encoding/json"
	"log"
	"strings"
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

// 🛠️ Split message into words
func SplitMessageContent(content string) []string {
	return strings.Split(content, " ")
}

// Extracts mentions
func ExtractMention(content string) string {
	words := strings.Split(content, " ")
	for _, word := range words {
		if strings.HasPrefix(word, "@") {
			return word[1:]
		}
	}
	return ""
}
