package main

import (
	pb "it-chain/service/webhook/proto"
	"google.golang.org/grpc"
	"log"
	"context"
)

func main() {

	address := "127.0.0.1:50070"
	payloadURL := "http://localhost:9000"

	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect : %v", err)
	}

	defer conn.Close()

	client := pb.NewWebhookClient(conn)

	_, err = client.Register(context.Background(), &pb.WebhookRequest{payloadURL})
	if err != nil {
		log.Fatalf("Failed to register webhook : %v", err)
	}

}

