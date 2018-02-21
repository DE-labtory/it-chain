package main

import (
	pb "it-chain/service/webhook/proto"
	"google.golang.org/grpc"
	"log"
	"context"
)

func main() {

	address := "127.0.0.1:50070"
	payloadURL := "127.0.0.1:8080"

	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	defer conn.Close()

	client := pb.NewWebhookClient(conn)

	res, err := client.Register(context.Background(), &pb.WebhookRequest{payloadURL})
	if err != nil {
		log.Fatalf("failed to register webhook : %v", err)
	}

	res, err = client.Remove(context.Background(), &pb.WebhookRequest{payloadURL})
	if err != nil {
		log.Fatalf("failed to remove webhook : %v", err)
	}

}

