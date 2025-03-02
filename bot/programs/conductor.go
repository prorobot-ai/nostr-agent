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
)

// **ConductorProgram** - Handles responding when mentioned
type ConductorProgram struct {
	IsRunning       bool
	CurrentRunCount int

	MaxRunCount   int
	ResponseDelay int
	Url           string
	Address       string

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

	if p.CurrentRunCount >= p.MaxRunCount {
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

	time.Sleep(time.Duration(p.ResponseDelay) * time.Second)

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

	p.StartCrawlJob(bot, reply)

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
func (p *ConductorProgram) StartCrawlJob(bot Bot, message *core.OutgoingMessage) {
	if p.CrawlerClient == nil {
		log.Println("❌ Crawler Client is not initialized")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	stream, err := p.CrawlerClient.StartCrawl(ctx, &pb.CrawlRequest{
		Url:   "https://example.com",
		JobId: "job123",
	})
	if err != nil {
		log.Fatalf("❌ Failed to start crawl: %v", err)
	}

	// Read streaming response
	for {
		resp, err := stream.Recv()
		if err != nil {
			break
		}
		log.Printf("🔄 Crawl Progress: %s", resp.Message)
	}
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
