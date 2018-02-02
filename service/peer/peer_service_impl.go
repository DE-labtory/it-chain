package peer

import (
	"it-chain/network/comm"
	"it-chain/common"
	pb "it-chain/network/protos"
	"github.com/golang/protobuf/proto"
)

var logger = common.GetLogger("peer_service.go")

type PeerServiceImpl struct {
	peerTable *PeerTable
	comm      comm.ConnectionManager
}

func NewPeerServiceImpl(peerTable *PeerTable,comm comm.ConnectionManager) PeerService{

	return &PeerServiceImpl{
		peerTable: peerTable,
		comm: comm,
	}
}

func (ps *PeerServiceImpl) GetPeerTable() PeerTable{
	return *ps.peerTable
}

//peer info 찾기
func (ps *PeerServiceImpl) GetPeerInfoByPeerID(peerID string) (*PeerInfo){

	peerInfo := ps.peerTable.FindPeerByPeerID(peerID)
	return peerInfo
}

//peer info
func (ps *PeerServiceImpl) PushPeerTable(peerIDs []string){

}

//주기적으로 handle 함수가 콜 된다.
//주기적으로 peerTable의 peerlist에게 peerTable을 전송한다.
//todo signing이 들어가야함
//todo struct to grpc proto의 변환 문제
func (ps *PeerServiceImpl) BroadCastPeerTable(interface{}){
	logger.Println("pushing peer table")

	peerInfos, err := ps.peerTable.SelectRandomPeerInfos(0.5)

	if err != nil{
		logger.Println("no peer exist")
		return
	}

	logger.Println("pushing peerTable:",ps.peerTable)

	ps.peerTable.IncrementHeartBeat()
	message := &pb.Message{}

	envelope := pb.Envelope{}
	envelope.Payload, err = proto.Marshal(message)

	if err !=nil{
		logger.Println("fail to serialize message")
	}

	errorCallBack := func(onError error) {
		logger.Println("fail to send message error:", onError.Error())
	}

	for _,peerInfo := range peerInfos{
		ps.comm.SendStream(envelope,errorCallBack, peerInfo.PeerID)
	}
}

func (ps *PeerServiceImpl) UpdatePeerTable(peerTable PeerTable){

	ps.peerTable.Lock()
	defer ps.peerTable.Unlock()

	for id, peerInfo := range peerTable.PeerMap{
		peer,ok := ps.peerTable.PeerMap[id]
		if ok{
			peer.Update(peerInfo)
		}else{
			ps.AddPeerInfo(peerInfo)
		}
	}

	ps.peerTable.UpdateTimeStamp()
}

func (ps *PeerServiceImpl) AddPeerInfo(peerInfo *PeerInfo){

	if peerInfo.PeerID == ""{
		logger.Error("failed to connect with", peerInfo)
		return
	}

	if peerInfo.GetEndPoint() == ""{
		logger.Error("failed to connect with", peerInfo)
		return
	}
	err := ps.comm.CreateStreamConn(peerInfo.PeerID,peerInfo.GetEndPoint(), nil)

	if err != nil{
		logger.Error("failed to connect with", peerInfo)
		return
	}

	ps.peerTable.AddPeerInfo(peerInfo)
}

func (ps *PeerServiceImpl) RequestPeerInfo(host string, port string) *PeerInfo{

	return nil
}
