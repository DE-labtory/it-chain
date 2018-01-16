package main

import (
	"io"
	"fmt"
	pb "it-chain/network/protos"
	"golang.org/x/net/context"
	"it-chain/service/peer"
	"time"
	"net"
	"log"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc"
	"it-chain/service/batch"
	"it-chain/network/comm"
	"it-chain/service/domain"
	"github.com/urfave/cli"
	"os"
	"it-chain/common"
)

var logger = common.GetLogger("peer.go")

type MessageServer struct {

}

//message를 받으면 state를 업데이트 한다.
func (s *MessageServer) Stream(stream pb.MessageService_StreamServer) (error) {

	for {
		message,err := stream.Recv()

		if err == io.EOF {
			return nil
		}

		if err != nil {
			return err
		}
		fmt.Printf("Received: %d\n", message.String())
	}
}

func (s *MessageServer) Ping(ctx context.Context, in *pb.Empty) (*pb.Empty, error) {
	return &pb.Empty{}, nil
}

func SetPeer(ipAddress string, peerID string){

	peer1 := &domain.PeerInfo{
		Port: "5555",
		PeerID: "peer1",
		IpAddress: "127.0.0.1",
		HeartBeat: 0,
		TimeStamp: time.Now(),
	}

	peerTable,err := domain.NewPeerTable(peer1)

	if err != nil{

	}

	comm := comm.NewCommImpl()

	peerService := peer.NewPeerServiceImpl(peerTable,comm)

	eventBatcher := batch.NewGRPCMessageBatcher(time.Second*5,peerService,false)

	eventBatcher.Add("push peerTable")

	lis, err := net.Listen("tcp", peer1.GetEndPoint())
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterMessageServiceServer(s, &MessageServer{})
	// Register reflection service on gRPC server.
	reflection.Register(s)

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
		s.Stop()
		lis.Close()
	}

}

func main(){

	var ipAddress = ""
	var peerID = ""

	app := cli.NewApp()
	app.Name = "PEER"
	app.Usage = "fight the loneliness!"

	app.Flags = []cli.Flag {
		cli.StringFlag{
			Name:        "ip, i",
			Usage:       "IP Address of Boot Peer",
			Destination: &ipAddress,
		},
		cli.StringFlag{
			Name:        "peerid, p",
			Usage:       "ID of Peer",
			Destination: &peerID,
		},
	}

	app.Action = func(c *cli.Context) error {



		if ipAddress == "" || peerID == ""{
			logger.Println("initiating boot peer")
		}else{
			logger.Println("initiating peer with ",peerID)
		}

		SetPeer(ipAddress,peerID)

		return nil
	}

	app.Run(os.Args)
}
