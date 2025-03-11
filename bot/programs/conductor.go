package programs

import (
	"agent/core"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	pb "github.com/prorobot-ai/grpc-protos/gen/crawler"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// ConductorProgram handles responses when mentioned
type ConductorProgram struct {
	IsRunning       bool
	CurrentRunCount int
	ProgramConfig   core.ProgramConfig
	Peers           []string
	CrawlerClient   pb.CrawlerServiceClient
}

// ✅ **Check if the program is active**
func (p *ConductorProgram) IsActive() bool {
	return p.IsRunning
}

// ✅ **Should this program run?**
func (p *ConductorProgram) ShouldRun(message *core.BusMessage) bool {
	return true
}

// ✅ **Run Responder Logic**
func (p *ConductorProgram) Run(bot Bot, message *core.BusMessage) string {
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

	text := message.Payload.Text
	mention := core.ExtractMention(text)
	aliases := bot.GetAliases()
	set := createSet(aliases)

	if mention == "" || !set[mention] {
		return "🟠 No valid mention"
	}

	words := core.SplitMessageContent(text)
	if len(words) < 2 {
		log.Println("⚠️ Malformed message, missing number.")
		return "🟠"
	}

	time.Sleep(time.Duration(p.ProgramConfig.ResponseDelay) * time.Second)

	remoteJob := &core.RemoteJob{
		ChannelID: message.ChannelID,
		SessionID: message.Payload.Metadata,
		Payload:   words[1],
	}

	p.StartWorkerJob(bot, *remoteJob)

	return "🟢"
}

// ✅ **Initialize gRPC Client in the Program**
func (p *ConductorProgram) InitCrawlerClient(serverAddr string) {
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
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

	// ✅ Initialize Notifier
	var notifier core.Notifier
	if p.ProgramConfig.HubConfig.Socket != "" {
		wsNotifier, err := core.NewWebSocketNotifier(p.ProgramConfig.HubConfig.Socket)
		if err != nil {
			log.Println("❌ WebSocket unavailable, falling back to Logger")
			notifier = &core.LoggerNotifier{}
		} else {
			notifier = wsNotifier
		}
	} else {
		notifier = &core.LoggerNotifier{}
	}

	// ✅ Start Crawl Job
	stream, err := p.CrawlerClient.StartCrawl(ctx, &pb.CrawlRequest{
		Url:   remoteJob.Payload,
		JobId: remoteJob.SessionID,
	})
	if err != nil {
		notifier.SendMessage(core.SocketRequest{
			Type:      "error",
			ChannelID: remoteJob.ChannelID,
			Metadata:  remoteJob.SessionID,
			Text:      "Failed to start crawl: " + err.Error(),
			CreatedAt: time.Now().Unix(),
		})
		return
	}

	// ✅ Handle Response
	p.handleWorkerResponse(bot, stream, remoteJob, notifier)
}

// ✅ **Handles gRPC Crawl Response via Notifier**
func (p *ConductorProgram) handleWorkerResponse(bot Bot, stream pb.CrawlerService_StartCrawlClient, remoteJob core.RemoteJob, notifier core.Notifier) {
	var jobID string
	for {
		resp, err := stream.Recv()
		if err != nil {
			// ✅ Check if stream closed unexpectedly
			if err == io.EOF {
				log.Println("✅ gRPC Stream reached EOF gracefully")
			} else {
				log.Printf("❌ gRPC Stream Closed Unexpectedly: %v", err)
			}
			break
		}

		log.Printf("🔄 Worker Job [%s] Progress: %s", resp.JobId, resp.Message)

		jobID = resp.JobId

		notifier.SendMessage(core.SocketRequest{
			Type:      "worker_update",
			ChannelID: remoteJob.ChannelID,
			Metadata:  remoteJob.SessionID,
			Text:      resp.Message,
			CreatedAt: time.Now().Unix(),
		})
	}

	// ✅ Ensure WebSocket remains open until messages are fully processed
	log.Println("✅ Sending worker_done now...")
	notifier.SendMessage(core.SocketRequest{
		Type:      "agent_update",
		ChannelID: remoteJob.ChannelID,
		Metadata:  remoteJob.SessionID,
		Text:      fmt.Sprintf("[%s] exiting program.", bot.GetName()),
		CreatedAt: time.Now().Unix(),
	})

	// ✅ Small delay to ensure WebSocket sends this message
	time.Sleep(500 * time.Millisecond)

	// ✅ Ensure gRPC stream fully closes before sending `clear_status`
	log.Println("✅ Waiting for gRPC shutdown to complete...")
	time.Sleep(500 * time.Millisecond) // Allow last message to flush

	log.Println("✅ Sending agent_done now...")
	notifier.SendMessage(core.SocketRequest{
		Type:      "agent_done",
		ChannelID: remoteJob.ChannelID,
		Metadata:  remoteJob.SessionID,
		Text:      "",
		CreatedAt: time.Now().Unix(),
	})

	// ✅ Final delay before closing WebSocket
	time.Sleep(500 * time.Millisecond)

	url := fmt.Sprintf("%s/%s", p.ProgramConfig.CallbackUrl, jobID)
	message := fmt.Sprintf("🧙🏻‍♂️⚡️ Finished. See report @ %s.", url)

	reply := &core.BusMessage{
		ChannelID:         remoteJob.ChannelID,
		ReceiverPublicKey: bot.GetPublicKey(),
		Payload: core.ContentStructure{
			Kind:     "message",
			Metadata: remoteJob.SessionID,
			Text:     core.SerializeContent(message, "message"),
		},
	}

	bot.Publish(reply)

	// ✅ Now we can safely close the WebSocket
	notifier.Close()
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
