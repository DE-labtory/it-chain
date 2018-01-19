package main

import (
	"io"
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
	"github.com/golang/protobuf/proto"
)

var logger = common.GetLogger("peer.go")

type MessageServer struct {
	peerService peer.PeerService
}

//message를 받으면 state를 업데이트 한다.
func (s *MessageServer) Stream(stream pb.MessageService_StreamServer) (error) {

	for {
		envelope,err := stream.Recv()

		if err == io.EOF {
			return nil
		}

		if err != nil {
			return err
		}

		message := &pb.Message{}

		err = proto.Unmarshal(envelope.Payload,message)

		if err != nil{
			logger.Error("fail to deserialize message")
		}

		logger.Println("recevied peerTable",*message.GetPeerTable_())

		s.peerService.UpdatePeerTable(*message.GetPeerTable_())

		if err != nil{
			logger.Println("fail to deserialize message")
		}
	}
}

func (s *MessageServer) Ping(ctx context.Context, in *pb.Empty) (*pb.Empty, error) {
	return &pb.Empty{}, nil
}

func SetPeer(ipAddress string, peerID string,port string, bootport string,myID string){

	peer1 := &domain.PeerInfo{
		Port: port,
		PeerID: myID,
		IpAddress: "127.0.0.1",
		HeartBeat: 0,
		TimeStamp: time.Now(),
	}

	peerTable,err := domain.NewPeerTable(peer1)

	if err != nil{

	}

	logger.Println(ipAddress)
	logger.Println(peerID)
	logger.Println(bootport)



	comm := comm.NewCommImpl()

	peerService := peer.NewPeerServiceImpl(peerTable,comm)

	if ipAddress != "" && peerID != "" && bootport != ""{
		peer2 := &domain.PeerInfo{
			Port: bootport,
			PeerID: peerID,
			IpAddress: ipAddress,
			HeartBeat: 0,
			TimeStamp: time.Now(),
		}
		logger.Println("boot peer added")
		peerService.AddPeerInfo(peer2)
		logger.Println(peerService.GetPeerTable())
	}

	eventBatcher := batch.NewGRPCMessageBatcher(time.Second*5,peerService,false)

	eventBatcher.Add("push peerTable")

	lis, err := net.Listen("tcp", peer1.GetEndPoint())
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterMessageServiceServer(s, &MessageServer{peerService})
	// Register reflection service on gRPC server.
	reflection.Register(s)

	logger.Println("peer is on",peer1.GetEndPoint())

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
		s.Stop()
		lis.Close()
	}

}


func main(){
	const defaultPost = "5555"
	const defaultID = "peer0"
	var ipAddress = ""
	var peerID = ""
	var port = ""
	var bootport = ""
	var myID = ""

	app := cli.NewApp()
	app.Name = "PEER"
	app.Usage = "check peer's status!"

	app.Flags = []cli.Flag {
		cli.StringFlag{
			Name:        "port, P",
			Usage:       "my port number",
			Destination: &port,
		},
		cli.StringFlag{
			Name:        "id, I",
			Usage:       "my port number",
			Destination: &myID,
		},
		cli.StringFlag{
			Name:        "bootip, bip",
			Usage:       "IP Address of Boot Peer",
			Destination: &ipAddress,
		},
		cli.StringFlag{
			Name:        "bootid, bid",
			Usage:       "ID of Boot Peer",
			Destination: &peerID,
		},cli.StringFlag{
			Name:        "bootport, bport",
			Usage:       "Port of Boot Peer",
			Destination: &bootport,
		},
	}

	app.Action = func(c *cli.Context) error {
		if port == ""{
			port = defaultPost
		}

		if myID == ""{
			myID = defaultID
		}
		if ipAddress == "" || peerID == ""{
			logger.Println("initiating boot peer")
		}else{
			logger.Println("initiating peer with ",peerID)
		}

		SetPeer(ipAddress,peerID,port,bootport,myID)

		return nil
	}

	app.Run(os.Args)
}
