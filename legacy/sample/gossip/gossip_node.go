package main

import (
	"log"
	"net"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	pb "github.com/it-chain/it-chain-Engine/legacy/sample/gossip/grpc"
)

const (
	port = ":50052"
)

type server struct{}


//gossip을 받음
func (s *server) PushGossip(ctx context.Context, in *pb.GossipTable) (*pb.Empty, error) {
	return &pb.Empty{},nil
}

func (s *server) Ping(ctx context.Context, in *pb.Empty) (*pb.Empty, error) {
	return &pb.Empty{},nil
}


// return my peer tables
func (s *server) GetGossipTable(ctx context.Context, in *pb.Empty) (*pb.GossipTable, error) {
	return &pb.GossipTable{
		MyID: "1",
	}, nil
}

func (s *server) PullGossip(ctx context.Context, in *pb.Empty) (*pb.GossipTable, error) {
	return &pb.GossipTable{}, nil
}

func main() {

	lis, err := net.Listen("tcp", port)

	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterGossipServer(s, &server{})
	// Register reflection service on gRPC server.
	reflection.Register(s)

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
