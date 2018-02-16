package mock

import (
	"net"
	"google.golang.org/grpc/reflection"
	pb "it-chain/network/protos"
	"io"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"log"
)

var counter = 0

type MockServerHandler func(stream pb.StreamService_StreamServer, envelope *pb.Envelope)

type Mockserver struct {
	handler MockServerHandler
}

func (s *Mockserver) Stream(stream pb.StreamService_StreamServer) (error) {

	for {
		envelope,err := stream.Recv()

		if err == io.EOF {
			return nil
		}

		if err != nil {
			return err
		}

		s.handler(stream,envelope)
	}
}

func (s *Mockserver) Ping(ctx context.Context, in *pb.Empty) (*pb.Empty, error) {
	return &pb.Empty{}, nil
}

func ListenMockServer(ipAddress string, handler MockServerHandler) (*grpc.Server,net.Listener){

	lis, err := net.Listen("tcp", ipAddress)

	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	mockServer := &Mockserver{}
	mockServer.handler = handler

	s := grpc.NewServer()
	pb.RegisterStreamServiceServer(s, mockServer)
	// Register reflection service on gRPC server.
	reflection.Register(s)

	go func(){
		if err := s.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
			s.Stop()
			lis.Close()
		}
	}()

	return s,lis
}

