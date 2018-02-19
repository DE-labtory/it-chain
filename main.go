package main

import (
	"it-chain/service"
	"it-chain/smartcontract"
	"it-chain/network/comm"
	"it-chain/auth"
	"hash/fnv"
	"github.com/spf13/viper"
)

type Node struct{
	blockService service.BlockService
	peerService service.PeerService
	consensusService service.ConsensusService
	smartContractService smartcontract.SmartContractService
	comm comm.ConnectionManager
	crypto auth.Crypto
}

func NewNode(){

	node := &Node{}

	//// crpyto
	crpyto, err := auth.NewCrypto(viper.GetString("key.default_path"),&auth.RSAKeyGenOpts{})

	if err != nil{

	}

	node.crypto = crpyto
}

func (n *Node) IDGenerator() uint32{
	pri, _, _ := n.crypto.GetKey()
	h := fnv.New32a()
	h.Write([]byte(pri.SKI()))

	return h.Sum32()
}

func main() {
	NewNode()
}