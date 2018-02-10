package comm

import (
	pb "it-chain/network/protos"
	"sync"
	"it-chain/common"
	"it-chain/auth"
	"encoding/json"
)

var commLogger = common.GetLogger("connection_manager_impl.go")

type ConnectionManagerImpl struct {
	connectionMap map[string]Connection
	crpyto        auth.Crypto
	sync.RWMutex
}

func NewConnectionManagerImpl(crpyto auth.Crypto) *ConnectionManagerImpl{
	return &ConnectionManagerImpl{
		connectionMap: make(map[string]Connection),
		crpyto: crpyto,
	}
}

func (comm *ConnectionManagerImpl) CreateStreamConn(connectionID string, ip string, handler ReceiveMessageHandle) error{

	//Peer의 ipAddress로 connection을 연결
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

func (comm *ConnectionManagerImpl) SendStream(data interface{}, errorCallBack OnError, connectionID string){


	payload, err := json.Marshal(data)

	if err != nil{
		return
	}

	_, pub, err  := comm.crpyto.LoadKey()

	if err != nil{
		return
	}

	sig, err  := comm.crpyto.Sign(payload,nil)

	if err != nil{
		return
	}

	envelope := &pb.Envelope{}
	envelope.Payload = payload
	envelope.Pubkey = pub.SKI()
	envelope.Signature = sig
	
	conn, ok := comm.connectionMap[connectionID]

	if ok{
		conn.Send(envelope,errorCallBack)
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

func (comm *ConnectionManagerImpl) Close(connectionID string){

	conn, ok := comm.connectionMap[connectionID]

	if ok{
		commLogger.Println("connection:",connectionID, "is closing")
		conn.Close()
		delete(comm.connectionMap,connectionID)
	}
}

func (comm *ConnectionManagerImpl) Size() int{
	return len(comm.connectionMap)
}