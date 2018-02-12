package main

import (
	"it-chain/service"
	"it-chain/domain"
	"it-chain/network/comm"
	pb "it-chain/network/protos"
	"io"
	"github.com/golang/protobuf/proto"
	"golang.org/x/net/context"
	"it-chain/auth"
	"net"
	"google.golang.org/grpc/reflection"
	"log"
	"google.golang.org/grpc"
	"it-chain/network/comm/publisher"
	"time"
	"github.com/urfave/cli"
)

var View = &domain.View{
	ID:"1",
	LeaderID: "peer1",
	PeerID: []string{"peer1","peer2","peer3","peer4"},
}

//pbft testing code
type Node struct {
	myInfo            *domain.Peer
	ip                string
	port              string
	consensusService  service.ConsensusService
	view              *domain.View
	connectionManager comm.ConnectionManager
	blockService      service.BlockService
	peerService       service.PeerService
	messagePublisher  publisher.MessagePublisher
}

func NewNode(peerInfo *domain.Peer) *Node{

	node := &Node{}
	node.myInfo = peerInfo

	crypto, err := auth.NewCrypto("./"+node.myInfo.GetEndPoint())

	if err != nil{
		panic("fail to create keys")
	}

	connectionManager := comm.NewConnectionManagerImpl(crypto)


	//consensusService
	consensusService := service.NewPBFTConsensusService(View,connectionManager,nil)

	//peerService
	peerTable,err := domain.NewPeerTable(node.myInfo)

	if err != nil{
		panic("error set peertable")
	}

	peerService := service.NewPeerServiceImpl(peerTable,connectionManager)

	eventBatcher := service.NewBatchService(time.Second*5,peerService.BroadCastPeerTable,false)
	eventBatcher.Add("push peerTable")

	node.consensusService = consensusService
	node.connectionManager = connectionManager
	node.view = View

	return node
}

func (s *Node) Stream(stream pb.MessageService_StreamServer) (error) {

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
	}
}

func (s *Node) Ping(ctx context.Context, in *pb.Empty) (*pb.Empty, error) {
	return &pb.Empty{}, nil
}

func (s *Node) listen(){

	lis, err := net.Listen("tcp", s.myInfo.GetEndPoint())
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	server := grpc.NewServer()
	pb.RegisterMessageServiceServer(server, s)
	// Register reflection service on gRPC server.
	reflection.Register(server)

	if err := server.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
		server.Stop()
		lis.Close()
	}
}

func main(){


	app := cli.NewApp()

	var myAddress = ""
	var bootAddress = ""

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "My Address",
			Usage:       "hostIP:port",
			Destination: &myAddress,
		},
		cli.StringFlag{
			Name:        "id, I",
			Usage:       "hostIP:port",
			Destination: &bootAddress,
		},
	}

	peer := &domain.Peer{}
	peer.PeerID = "peer1"
	peer.IpAddress = "localhost"
	peer.Port = "4444"

	node := NewNode(peer)
	node.listen()
}


