package programs

import (
	"agent/core"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"time"
)

// **CallbackProgram** - Handles responding when mentioned
type CallbackProgram struct {
	IsRunning       bool
	CurrentRunCount int

	ProgramConfig core.ProgramConfig

	Peers []string
}

// ✅ **Check if the program is active**
func (p *CallbackProgram) IsActive() bool {
	return p.IsRunning
}

// ✅ **Should this program run?**
func (p *CallbackProgram) ShouldRun(message *core.BusMessage) bool {
	return true
}

// ✅ **Run Callback Logic**
func (p *CallbackProgram) Run(bot Bot, message *core.BusMessage) string {
	log.Printf("🏃 [%s] [CallbackProgram] [%d]", bot.GetPublicKey(), p.CurrentRunCount)

	if p.CurrentRunCount >= p.ProgramConfig.MaxRunCount {
		log.Printf("🛑 [%s] [CallbackProgram] reached max run count. Terminating...", bot.GetPublicKey())
		p.IsRunning = false
		return "🔴"
	}

	if !p.IsRunning {
		p.IsRunning = true
		p.CurrentRunCount = 0
	}

	p.CurrentRunCount++

	text := message.Payload.Text
	kind := message.Payload.Kind

	pattern := p.ProgramConfig.Pattern

	log.Printf("✔️ [%s] [%s] [%s]", kind, text, pattern)

	re, err := regexp.Compile(pattern)
	if err != nil {
		log.Printf("Error compiling regex: %v", err)
		return "❌" // Indicate an error
	}

	time.Sleep(time.Duration(p.ProgramConfig.ResponseDelay) * time.Second)

	signal := "🟠"
	if re.MatchString(text) {
		log.Println("text matched")

		log.Println(text)

		data := PostData{
			Message: message.Payload.Text,
		}

		postData(p.ProgramConfig.CallbackUrl, data)
		signal = "🟢"
	}

	if re.MatchString(kind) {
		log.Println("kind matched")

		data := PostData{
			Message: message.Payload.Text,
		}

		postData(p.ProgramConfig.CallbackUrl, data)
		signal = "🟢"
	}

	return signal
}

type PostData struct {
	Message string `json:"message"`
}

func postData(url string, data PostData) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		fmt.Println("Error marshalling JSON:", err)
		return
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error making request:", err)
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return
	}

	fmt.Println("Response Status:", resp.Status)
	fmt.Println("Response Body:", string(body))
}
