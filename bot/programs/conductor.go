package programs

import (
	"agent/core"
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	pb "github.com/prorobot-ai/grpc-protos/gen/crawler"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/gorilla/websocket"
)

// **ConductorProgram** - Handles responding when mentioned
type ConductorProgram struct {
	IsRunning       bool
	CurrentRunCount int

	ProgramConfig core.ProgramConfig

	Peers []string

	CrawlerClient pb.CrawlerServiceClient
}

// ✅ **Check if the program is active**
func (p *ConductorProgram) IsActive() bool {
	return p.IsRunning
}

// ✅ **Should this program run?**
func (p *ConductorProgram) ShouldRun(message *core.OutgoingMessage) bool {
	return true
}

// ✅ **Run Responder Logic**
func (p *ConductorProgram) Run(bot Bot, message *core.OutgoingMessage) string {
	log.Printf("🏃 [%s] [ConductorProgram] [%d]", bot.GetPublicKey(), p.CurrentRunCount)

	if p.CurrentRunCount >= p.ProgramConfig.MaxRunCount {
		log.Printf("🛑 [%s] [ConductorProgram] reached max run count. Terminating...", bot.GetPublicKey())
		p.IsRunning = false
		return "🔴"
	}

	if !p.IsRunning {
		p.IsRunning = true
		p.CurrentRunCount = 0
	}

	p.CurrentRunCount++

	mention := core.ExtractMention(message.Content)
	aliases := bot.GetAliases()
	set := createSet(aliases)

	if mention == "" || !set[mention] {
		return "🟠 No valid mention"
	}

	words := core.SplitMessageContent(message.Content)
	if len(words) < 2 {
		log.Println("⚠️ Malformed message, missing number.")
		return "🟠"
	}

	time.Sleep(time.Duration(p.ProgramConfig.ResponseDelay) * time.Second)

	// HTTP
	// err := sendJobRequest(p.Url, words[1]) // send the request to jobs service
	// if err != nil {
	// 	log.Printf("❌ Error sending job: %v", err)
	// 	return "🔴"
	// }

	// GRPC
	reply := &core.OutgoingMessage{
		Content:           core.CreateContent("🧙🏻‍♂️ "+words[1]+" ⚡️", "message"),
		ChannelID:         message.ChannelID,
		ReceiverPublicKey: bot.GetPublicKey(),
	}

	remoteJob := &core.RemoteJob{
		ChannelID: message.ChannelID,
		Payload:   words[1],
	}

	p.StartWorkerJob(bot, *remoteJob)

	bot.Publish(reply)

	return "🟢"
}

// ✅ **Initialize gRPC Client in the Program**
func (p *ConductorProgram) InitCrawlerClient(serverAddr string) {
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()), // Use insecure connection (change for TLS)
	}

	conn, err := grpc.NewClient(serverAddr, opts...)
	if err != nil {
		log.Fatalf("❌ Failed to connect to crawler service: %v", err)
	}

	p.CrawlerClient = pb.NewCrawlerServiceClient(conn)
}

// ✅ **Send Crawl Request**
func (p *ConductorProgram) StartWorkerJob(bot Bot, remoteJob core.RemoteJob) {
	if p.CrawlerClient == nil {
		log.Println("❌ Crawler Client is not initialized")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	stream, err := p.CrawlerClient.StartCrawl(ctx, &pb.CrawlRequest{
		Url:   remoteJob.Payload,
		JobId: "job123",
	})
	if err != nil {
		log.Fatalf("❌ Failed to start crawl: %v", err)
	}

	// ✅ Handle response (WebSocket or direct logs)
	handleWorkerResponse(stream, remoteJob, p.ProgramConfig.HubConfig.Socket)
}

// ✅ **Handles gRPC Crawl Response**
func handleWorkerResponse(stream pb.CrawlerService_StartCrawlClient, remoteJob core.RemoteJob, wsURL string) {
	if wsURL == "" {
		// ✅ Log crawl updates if no WebSocket is configured
		for {
			resp, err := stream.Recv()
			if err != nil {
				break
			}
			log.Printf("🔄 Worker Progress: %s", resp.Message)
		}
	} else {
		// ✅ Forward crawl updates to WebSocket
		forwardToWebSocket(stream, wsURL, remoteJob.ChannelID)
	}
}

// ✅ **Send gRPC responses to WebSocket (Short-Lived Session)**
func forwardToWebSocket(stream pb.CrawlerService_StartCrawlClient, wsURL string, channelID string) {
	// 🔹 Establish WebSocket connection (Single Session)
	conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		log.Printf("❌ Failed to connect to WebSocket: %v", err)
		return
	}
	defer func() {
		log.Println("🔴 Closing WebSocket connection")
		conn.Close()
	}()

	log.Println("✅ WebSocket connection established:", wsURL)

	conn.SetPongHandler(func(appData string) error {
		log.Println("✅ Pong received, client is alive")
		return nil
	})

	// 🔹 Read gRPC stream and send each message to WebSocket
	for {
		resp, err := stream.Recv()
		if err != nil {
			log.Printf("❌ gRPC Stream Closed: %v", err)
			break // ✅ No retry needed, just exit
		}

		log.Printf("🔄 Crawl Progress: %s", resp.Message)

		// 🔹 Format WebSocket message
		wsMessage := map[string]string{
			"type":      "worker_update",
			"channelId": channelID,
			"text":      resp.Message,
		}

		jsonMessage, _ := json.Marshal(wsMessage)

		// 🔹 Send message to WebSocket
		err = conn.WriteMessage(websocket.TextMessage, jsonMessage)
		if err != nil {
			log.Printf("❌ Failed to send message to WebSocket: %v", err)
			break // ✅ Exit without retry
		}
	}

	log.Println("🔴 WebSocket session closed gracefully")
}

type JobRequest struct {
	Query string `json:"query"`
}

func sendJobRequest(url string, query string) error {
	job := JobRequest{Query: query}
	body, err := json.Marshal(job)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		log.Printf("Failed to submit job: %s", resp.Status)
		return err
	}

	log.Println("✅ Job submitted successfully!")
	return nil
}

func createSet(arr []string) map[string]bool {
	set := make(map[string]bool)
	for _, v := range arr {
		set[v] = true
	}
	return set
}
