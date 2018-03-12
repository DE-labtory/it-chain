package main

import (
	"gopkg.in/jeevatkm/go-model.v1"
	"github.com/it-chain/it-chain-Engine/domain"
	pb "github.com/it-chain/it-chain-Engine/network/protos"
	"fmt"
)

func main(){

	peer1 := &domain.Peer{}
	peer1.PeerID ="1"
	peer1.IpAddress = "127.0.0.1"
	peer1.PubKey = []byte("hello world!")

	peerTable := &domain.PeerTable{}
	peerTable.MyID = "123"
	peerTable.PeerMap = make(map[string]*domain.Peer)
	peerTable.AddPeer(peer1)

	pbPeerTable := &pb.PeerTable{}

	//srchResMap, err := model.Map(peerTable)

	err := model.Copy(pbPeerTable, peerTable)

	if err != nil{
		fmt.Println(err)
	}

	fmt.Println(pbPeerTable)
}