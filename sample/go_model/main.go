package main

import (
	"gopkg.in/jeevatkm/go-model.v1"
	"it-chain/domain"
	pb "it-chain/network/protos"
	"fmt"
)

func main(){

	peer1 := &domain.PeerInfo{}
	peer1.PeerID ="1"
	peer1.IpAddress = "127.0.0.1"
	peer1.PubKey = []byte("hello world!")

	peerTable := &domain.PeerTable{}
	peerTable.OwnerID = "123"
	peerTable.PeerMap = make(map[string]*domain.PeerInfo)
	peerTable.AddPeerInfo(peer1)

	pbPeerTable := &pb.PeerTable{}

	//srchResMap, err := model.Map(peerTable)

	err := model.Copy(pbPeerTable, peerTable)

	if err != nil{
		fmt.Println(err)
	}

	fmt.Println(pbPeerTable)
}