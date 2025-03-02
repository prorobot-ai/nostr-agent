package grpcclient

import (
	"context"
	"log"
	"time"

	pb "github.com/prorobot-ai/grpc-protos/gen/crawler"
	"google.golang.org/grpc"
)

type CrawlerClient struct {
	conn   *grpc.ClientConn
	client pb.CrawlerServiceClient
}

// ✅ Initialize the gRPC Client
func NewCrawlerClient(serverAddr string) (*CrawlerClient, error) {
	conn, err := grpc.Dial(serverAddr, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		return nil, err
	}

	client := pb.NewCrawlerServiceClient(conn)
	return &CrawlerClient{conn: conn, client: client}, nil
}

// ✅ Send a Crawl Request
func (c *CrawlerClient) StartCrawl(url, jobID string) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	stream, err := c.client.StartCrawl(ctx, &pb.CrawlRequest{
		Url:   url,
		JobId: jobID,
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

// ✅ Fetch Job Status
func (c *CrawlerClient) GetJobStatus(jobID string) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	stream, err := c.client.GetJobStatus(ctx, &pb.JobStatusRequest{JobId: jobID})
	if err != nil {
		log.Fatalf("❌ Failed to get job status: %v", err)
	}

	// Read streaming response
	for {
		resp, err := stream.Recv()
		if err != nil {
			break
		}
		log.Printf("📡 Job [%s] Status: %s", resp.JobId, resp.Status)
	}
}

// ✅ Close the gRPC connection
func (c *CrawlerClient) Close() {
	c.conn.Close()
}
