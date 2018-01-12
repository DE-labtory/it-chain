package main

import (
"log"
"golang.org/x/net/context"
"google.golang.org/grpc"
pb "it-chain/sample/gossip/grpc"
)

const (
	address     = "localhost:50052"
	defaultName = "world"
)

func main() {
	// Set up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	defer conn.Close()
	c := pb.NewGossipClient(conn)


	r, err := c.GetGossipTable(context.Background(), &pb.Empty{})

	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}

	log.Printf("Greeting: %s", r.String())
}
