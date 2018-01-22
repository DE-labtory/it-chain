package comm

import (
	pb "it-chain/network/protos"
	"sync"
	"it-chain/service/domain"
	"it-chain/common"
)

var commLogger = common.GetLogger("comm_impl.go")

type ConnectionManagerImpl struct{
	connectionMap map[string]*Connection
	sync.RWMutex
}

func NewConnectionManagerImpl() *ConnectionManagerImpl{
	return &ConnectionManagerImpl{
		connectionMap: make(map[string]*Connection),
	}
}

func (comm *ConnectionManagerImpl) CreateStreamConn(connectionID string, ip string, handler handler) error{

	//peerInfo의 ipAddress로 connection을 연결
	_, ok := comm.connectionMap[connectionID]

	if ok{
		return nil
	}

	grpcConnection,err := NewConnectionWithAddress(ip,false,nil)

	if err != nil{
		return err
	}

	conn,err := NewConnection(grpcConnection,handler,connectionID)

	if err != nil{
		return err
	}

	comm.Lock()
	comm.connectionMap[connectionID] = conn
	comm.Unlock()

	commLogger.Println("new connection:",connectionID, "are created")

	return nil
}

func (comm *ConnectionManagerImpl) SendStream(envelope pb.Envelope, errorCallBack onError, connectionID string){

	conn, ok := comm.connectionMap[connectionID]

	if ok{
		conn.Send(&envelope,errorCallBack)
		//todo 어떤 error일 경우에 conn을 close 할지 정해야함
	}else{
		//todo 처리
	}
}

func (comm *ConnectionManagerImpl) Stop(){
	commLogger.Println("all connections are closing")
	for id, conn := range comm.connectionMap{
		conn.Close()
		delete(comm.connectionMap,id)
	}
}

func (comm *ConnectionManagerImpl) Close(peerInfo domain.PeerInfo){

	conn, ok := comm.connectionMap[peerInfo.PeerID]

	if ok{
		commLogger.Println("connection:",peerInfo.PeerID, "is closing")
		conn.Close()
		delete(comm.connectionMap,peerInfo.PeerID)
	}
}

func (comm *ConnectionManagerImpl) Size() int{
	return len(comm.connectionMap)
}