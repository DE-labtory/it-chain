package main

import (
	pb "github.com/it-chain/it-chain-Engine/service/webhook/proto"
	"google.golang.org/grpc"
	"log"
	"context"
)

func main() {

	address := "127.0.0.1:44444"
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

