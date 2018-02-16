package service

import (
	"it-chain/network/comm"
	pb "it-chain/network/protos"
	"it-chain/domain"
)

type PeerServiceImpl struct {
	peerTable *domain.PeerTable
	comm      comm.ConnectionManager
}

func NewPeerServiceImpl(peerTable *domain.PeerTable,comm comm.ConnectionManager) PeerService{

	return &PeerServiceImpl{
		peerTable: peerTable,
		comm: comm,
	}
}

func (ps *PeerServiceImpl) GetPeerTable() domain.PeerTable{
	return *ps.peerTable
}

//peer info 찾기
func (ps *PeerServiceImpl) GetPeerByPeerID(peerID string) (*domain.Peer){

	Peer := ps.peerTable.FindPeerByPeerID(peerID)
	return Peer
}

//peer info
func (ps *PeerServiceImpl) PushPeerTable(peerIDs []string){

}

//주기적으로 handle 함수가 콜 된다.
//주기적으로 peerTable의 peerlist에게 peerTable을 전송한다.
//todo struct to grpc proto의 변환 문제
func (ps *PeerServiceImpl) BroadCastPeerTable(interface{}){
	//logger.Println("pushing peer table")

	Peers, err := ps.peerTable.SelectRandomPeers(0.5)

	if err != nil{
		logger.Println("no peer exist")
		return
	}

	ps.peerTable.IncrementHeartBeat()

	message := &pb.StreamMessage{}
	message.Content = &pb.StreamMessage_PeerTable{
		PeerTable: domain.ToProtoPeerTable(*ps.peerTable),
	}

	if err !=nil{
		logger.Println("fail to serialize message")
	}

	errorCallBack := func(onError error) {
		logger.Println("fail to send message error:", onError.Error())
	}

	for _,Peer := range Peers{
		ps.comm.SendStream(message,errorCallBack, Peer.PeerID)
	}
}

func (ps *PeerServiceImpl) UpdatePeerTable(peerTable domain.PeerTable){

	ps.peerTable.Lock()
	defer ps.peerTable.Unlock()

	for id, Peer := range peerTable.PeerMap{
		peer,ok := ps.peerTable.PeerMap[id]
		if ok{
			peer.Update(Peer)
		}else{
			ps.AddPeer(Peer)
		}
	}

	ps.peerTable.UpdateTimeStamp()
}

func (ps *PeerServiceImpl) AddPeer(Peer *domain.Peer){

	if Peer.PeerID == ""{
		logger.Error("failed to connect with", Peer)
		return
	}

	if Peer.GetEndPoint() == ""{
		logger.Error("failed to connect with", Peer)
		return
	}

	err := ps.comm.CreateStreamClientConn(Peer.PeerID,Peer.GetEndPoint(), nil)

	if err != nil{
		logger.Error("failed to connect with", Peer)
		return
	}

	ps.peerTable.AddPeer(Peer)
}

func (ps *PeerServiceImpl) RequestPeer(host string, port string) *domain.Peer{

	return nil
}

func (ps *PeerServiceImpl) GetLeader() *domain.Peer{
	return ps.peerTable.Leader
}