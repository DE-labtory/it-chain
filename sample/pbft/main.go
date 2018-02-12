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
	"os"
	"strings"
	"github.com/rs/xid"
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

	crypto, err := auth.NewCrypto("./sample/pbft/"+node.myInfo.GetEndPoint())
	_, pub ,err := crypto.GenerateKey(&auth.RSAKeyGenOpts{})

	if err !=nil{
		log.Println(err)
	}

	node.myInfo.PubKey = pub.SKI()

	log.Println(node.myInfo.PubKey)

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
	node.peerService = peerService
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
		log.Println("Received Envelop:",envelope)
	}
}

func (s *Node) Ping(ctx context.Context, in *pb.Empty) (*pb.Empty, error) {
	return &pb.Empty{}, nil
}

func (s *Node) GetPeer(context.Context, *pb.Empty) (*pb.Peer, error){

	peer := &pb.Peer{}
	peer.Port = s.myInfo.Port
	peer.PubKey = s.myInfo.PubKey
	peer.IpAddress = s.myInfo.IpAddress
	peer.HeartBeat = int32(s.myInfo.HeartBeat)
	peer.PeerID = s.myInfo.PeerID

	return peer,nil
}

func (s *Node) listen(){

	lis, err := net.Listen("tcp", s.myInfo.GetEndPoint())
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	server := grpc.NewServer()
	pb.RegisterMessageServiceServer(server, s)
	pb.RegisterPeerServiceServer(server,s)
	// Register reflection service on gRPC server.
	reflection.Register(server)

	if err := server.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
		server.Stop()
		lis.Close()
	}
}

func (s *Node) RequestPeer(address string) *pb.Peer{
	log.Println("request peer Information to boot peer:",address)

	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	defer conn.Close()
	c := pb.NewPeerServiceClient(conn)

	peer, err := c.GetPeer(context.Background(), &pb.Empty{})

	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}

	log.Println("recevied peer Info:",peer)

	return peer
}

func main(){

	app := cli.NewApp()

	var myAddress = ""
	var bootAddress = ""

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "IP, ip",
			Usage:       "hostIP:port",
			Destination: &myAddress,
		},
		cli.StringFlag{
			Name:        "BootIP, bi",
			Usage:       "hostIP:port",
			Destination: &bootAddress,
		},
	}

	app.Action = func(c *cli.Context) error {

		address := strings.Split(myAddress,":")

		peer := &domain.Peer{}
		peer.PeerID = xid.New().String()
		peer.IpAddress = address[0]
		peer.Port = address[1]

		log.Println(myAddress)

		node := NewNode(peer)

		if bootAddress != ""{
			log.Println("searching boot peer...")
			p := node.RequestPeer(bootAddress)
			log.Println(p)
			bootPeer := domain.FromProtoPeer(*p)
			log.Println(bootPeer)
			log.Println(node.peerService)
			node.peerService.AddPeer(bootPeer)
		}

		node.listen()

		return nil
	}

	app.Run(os.Args)
}


