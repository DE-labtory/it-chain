package message

import (

	"log"
	"github.com/golang/protobuf/proto"
)

func (envelope *Envelope) GetMessage() (*StreamMessage, error){

	m := &StreamMessage{}

	err := proto.Unmarshal(envelope.Payload,m)

	log.Print(m)

	if err != nil{
		return nil, err
	}

	return m, nil
}

//
//import (
//	"it-chain/service/peer"
//	"fmt"
//	"it-chain/service/peer/event"
//)
//
//func PeerTableToTable(peerTable peer.PeerTable) *Message_PeerTable{
//
//	peerTable_ := &Message_PeerTable{}
//	peerTable_.PeerTable = &PeerTable{}
//	peerTable_.PeerTable.MyID = peerTable.MyID
//	pbPeers := make([]*Peer,0)
//
//	for _, Peer := range peerTable.PeerMap {
//		pbPeer := &Peer{}
//		pbPeer.IpAddress = Peer.IpAddress
//		pbPeer.PeerID = Peer.PeerID
//		pbPeer.HeartBeat = int32(Peer.HeartBeat)
//		pbPeer.Port = Peer.Port
//		pbPeer.PubKey = Peer.PubKey
//		pbPeers = append(pbPeers, pbPeer)
//	}
//
//	peerTable_.PeerTable.Peers = pbPeers
//
//	return peerTable_
//}
//
//func (m *Message) FromPeerTable(peerTable peer.PeerTable){
//	m.GetPeerTable().MyID = peerTable.MyID
//
//}
//
//func (m *Message) SetPeers(Peers []peer.Peer, MyID string){
//
//	fmt.Println(m.GetPeerTable())
//
//	m.GetPeerTable().MyID = MyID
//
//	pbPeers := make([]*Peer,0)
//
//	for _, Peer := range Peers {
//		pbPeer := &Peer{}
//		pbPeer.PeerID = Peer.PeerID
//		pbPeer.HeartBeat = int32(Peer.HeartBeat)
//		pbPeer.Port = Peer.Port
//		pbPeer.PubKey = Peer.PubKey
//		pbPeers = append(pbPeers, pbPeer)
//	}
//
//	m.GetPeerTable().Peers = pbPeers
//
//}
//
//func (m *Message) GetPeerTable_() *peer.PeerTable{
//
//	peerTable := &peer.PeerTable{}
//	peerTable.MyID = m.GetPeerTable().MyID
//	peerTable.PeerMap = make(map[string]*peer.Peer)
//
//	for _, Peer := range m.GetPeerTable().Peers {
//		Peer_ := &peer.Peer{}
//		Peer_.IpAddress = Peer.IpAddress
//		Peer_.PeerID  = Peer.PeerID
//		Peer_.PubKey = Peer.PubKey
//		Peer_.HeartBeat = int(Peer.HeartBeat)
//		Peer_.Port = Peer.Port
//
//		peerTable.PeerMap[Peer.PeerID] = Peer_
//	}
//
//	return peerTable
//}
//
//func (m *Message) IsPeerTableUpdateMessage() bool{
//	return m.GetPeerTable() != nil
//}
//
//
//func (message *Message) GetMessageType() string {
//
//	if message.GetPeerTable() != nil{
//		return network.UpdatePeerTable
//	}
//
//	return "no"
//}