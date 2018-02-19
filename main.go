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
)



type Node struct {
	id                   string
	ip                   string
	host                 string
	port                 string
	blockService         service.BlockService
	peerService          service.PeerService
	consensusService     service.ConsensusService
	smartContractService service.SmartContractService
	connectionManager    comm.ConnectionManager
	crypto               auth.Crypto
}

func NewNode(ip string){

	//ip format (xxx.xxx.xxx.xxx:pppp)

	node := &Node{}

	////set baisc Info
	node.ip = ip
	node.host = strings.Split(ip,":")[0]
	node.port = strings.Split(ip,":")[1]

	//// crpyto
	crpyto, err := auth.NewCrypto(viper.GetString("key.defaultPath"),&auth.RSAKeyGenOpts{})

	if err != nil{
		common.Log.Errorln("crypto create error")
		return
	}

	node.crypto = crpyto
	node.id = node.GenerateID()

	///// comm
	//todo need to set stream server
	connectionManager := comm.NewConnectionManagerImpl(crpyto)
	node.connectionManager = connectionManager

	//// peerService
	peer := &domain.Peer{}
	peer.PeerID = node.id
	peer.Port = node.port
	_, pub, _ := node.crypto.GetKey()
	peer.PubKey = pub.SKI()
	peer.IpAddress = node.ip
	peer.HeartBeat = 0
	peer.TimeStamp = time.Now()

	common.Log.Println(node.id)

	peerTable ,err := domain.NewPeerTable(peer)

	if err != nil{
		common.Log.Errorln("peerTable create error")
		return
	}

	peerService := service.NewPeerServiceImpl(peerTable,connectionManager)
	node.peerService = peerService


	/////blockService
	blockService := service.NewLedger(viper.GetString("ledger.defaultPath"))
	node.blockService = blockService

	///smartContractService
	smartContractService := service.NewSmartContractService(viper.GetString("smartcontract.defaultPath"),viper.GetString("smartContract.githubID"))

	/////consensusService
	consensusService := service.NewPBFTConsensusService(node.connectionManager,node.blockService,node.id,smartContractService)
	node.consensusService = consensusService


}

func (n *Node) GenerateID() string{

	pri, _, _ := n.crypto.GetKey()
	h := sha1.New()
	h.Write(pri.SKI())

	return base64.URLEncoding.EncodeToString(h.Sum(nil))
}

func main() {
	NewNode("127.0.0.1:5555")
}