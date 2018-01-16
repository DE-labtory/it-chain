package comm

import (
	pb "it-chain/network/protos"
	"sync"
	"it-chain/service/domain"
)
type CommImpl struct{
	connectionMap map[string]*Connection
	sync.RWMutex
}

func NewCommImpl() *CommImpl{
	return &CommImpl{
		connectionMap: make(map[string]*Connection),
	}
}

func (comm *CommImpl)CreateConn(peerInfo domain.PeerInfo) error{

	//peerInfo의 ipAddress로 connection을 연결
	endpoint := peerInfo.GetEndPoint()
	grpcConnection,err := NewConnectionWithAddress(endpoint,false,nil)

	if err != nil{
		return err
	}

	conn,err  := NewConnection(grpcConnection)

	if err != nil{
		return err
	}

	comm.Lock()
	comm.connectionMap[peerInfo.PeerID] = conn
	comm.Unlock()

	return nil
}

func (comm *CommImpl) Send(envelop pb.Envelope, errorCallBack onError, peerInfos ...domain.PeerInfo){

	for _, peerInfo := range peerInfos{
		conn, ok := comm.connectionMap[peerInfo.PeerID]
		if ok{
			conn.SendWithStream(&envelop,errorCallBack)
		}else{
			//todo 처리
		}
	}
}

func (comm *CommImpl) Stop(){

	for id, conn := range comm.connectionMap{
		conn.Close()
		delete(comm.connectionMap,id)
	}
}

func (comm *CommImpl) Close(peerInfo domain.PeerInfo){

	conn, ok := comm.connectionMap[peerInfo.PeerID]

	if ok{
		conn.Close()
		delete(comm.connectionMap,peerInfo.PeerID)
	}
}

func (comm *CommImpl) Size() int{
	return len(comm.connectionMap)
}