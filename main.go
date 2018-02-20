package main

import (
	"it-chain/service"
	"it-chain/network/comm"
	"it-chain/auth"
	"github.com/spf13/viper"
	"it-chain/common"
	"strings"
	"it-chain/domain"
	"time"
	"encoding/base64"
	"crypto/sha1"
	"sync"
	"fmt"
	"net"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc"
	pb "it-chain/network/protos"
	"golang.org/x/net/context"
)

type Node struct {
	identity             *domain.Peer
	blockService         service.BlockService
	peerService          service.PeerService
	consensusService     service.ConsensusService
	smartContractService service.SmartContractService
	connectionManager    comm.ConnectionManager
	crypto               auth.Crypto
}

//todo msg publisher를 통한 receive Message분배 로직 추가
func NewNode(ip string) *Node{

	//ip format (xxx.xxx.xxx.xxx:pppp)
	node := &Node{}

	////set baisc Info
	node.identity = &domain.Peer{}
	node.identity.IpAddress = strings.Split(ip,":")[0]
	node.identity.Port = strings.Split(ip,":")[1]

	//// crpyto
	crpyto, err := auth.NewCrypto(viper.GetString("key.defaultPath"),&auth.RSAKeyGenOpts{})

	if err != nil{
		common.Log.Errorln("crypto create error")
		return nil
	}

	node.crypto = crpyto
	node.identity.PeerID = node.GenerateID()

	///// comm
	connectionManager := comm.NewConnectionManagerImpl(crpyto)
	node.connectionManager = connectionManager

	//// peerService
	_, pub, _ := node.crypto.GetKey()
	node.identity.PubKey = pub.SKI()
	node.identity.HeartBeat = 0
	node.identity.TimeStamp = time.Now()

	peerTable ,err := domain.NewPeerTable(node.identity)

	if err != nil{
		common.Log.Errorln("peerTable create error")
		return nil
	}

	peerService := service.NewPeerServiceImpl(peerTable,connectionManager)
	node.peerService = peerService

	///// blockService
	blockService := service.NewLedger(viper.GetString("ledger.defaultPath"))
	node.blockService = blockService

	///// smartContractService
	smartContractService := service.NewSmartContractService(viper.GetString("smartcontract.defaultPath"),viper.GetString("smartContract.githubID"))

	///// consensusService
	consensusService := service.NewPBFTConsensusService(node.connectionManager,node.blockService,node.identity,smartContractService)
	node.consensusService = consensusService

	return node
}

func (n *Node) GenerateID() string{

	pri, _, _ := n.crypto.GetKey()
	h := sha1.New()
	h.Write(pri.SKI())

	return base64.URLEncoding.EncodeToString(h.Sum(nil))
}

func (n *Node) SearchBootNode(bootip string){

	common.Log.Println("Searching peer", bootip)
	peer, err := n.peerService.RequestPeer(bootip)

	if err!=nil{
		return
	}

	n.peerService.AddPeer(peer)
}

func (n *Node) GetPeer(context.Context, *pb.Empty) (*pb.Peer, error){

	pp := domain.ToProtoPeer(*n.identity)

	return pp,nil
}

func (n* Node) Run() {
	common.Log.Println("Run it-chain")

	lis, err := net.Listen("tcp", n.identity.GetEndPoint())

	if err != nil {
		common.Log.Fatalf("failed to listen: %v", err)
	}

	server := grpc.NewServer()
	pb.RegisterStreamServiceServer(server,n.connectionManager)
	pb.RegisterPeerServiceServer(server,n)

	reflection.Register(server)

	common.Log.Println("It-chain peer is on:",n.identity.Port)

	if err := server.Serve(lis); err != nil {
		common.Log.Fatalf("failed to serve: %v", err)
		server.Stop()
		lis.Close()
	}
}

func PrintLogo(){
	fmt.Println(`
	___  _________               ________  ___  ___  ________  ___  ________
	|\  \|\___   ___\            |\   ____\|\  \|\  \|\   __  \|\  \|\   ___  \
	\ \  \|___ \  \_|____________\ \  \___|\ \  \\\  \ \  \|\  \ \  \ \  \\ \  \
	 \ \  \   \ \  \|\____________\ \  \    \ \   __  \ \   __  \ \  \ \  \\ \  \
	  \ \  \   \ \  \|____________|\ \  \____\ \  \ \  \ \  \ \  \ \  \ \  \\ \  \
           \ \__\   \ \__\              \ \_______\ \__\ \__\ \__\ \__\ \__\ \__\\ \__\
	    \|__|    \|__|               \|_______|\|__|\|__|\|__|\|__|\|__|\|__| \|__|
	`)
}

func main() {

	PrintLogo()

	ip := viper.GetString("Node.ip")
	bootIp := viper.GetString("bootNode.ip")

	n := NewNode(ip)
	n.SearchBootNode(bootIp)

	wg := sync.WaitGroup{}
	wg.Add(1)

	go n.Run()

	wg.Wait()
}