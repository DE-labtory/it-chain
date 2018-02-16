package comm

import (
	"sync"
	"it-chain/common"
	"it-chain/auth"
	"github.com/golang/protobuf/proto"
	"golang.org/x/net/context"
	pb "it-chain/network/protos"
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

func (comm *ConnectionManagerImpl) CreateStreamClientConn(connectionID string, ip string, handler ReceiveMessageHandle) error{

	//Peer의 connectionID로 connection을 연결
	_, ok := comm.connectionMap[connectionID]

	if ok{
		return nil
	}

	grpcConnection,err := NewConnectionWithAddress(ip,false,nil)

	if err != nil{
		return err
	}

	ctx, cf := context.WithCancel(context.Background())
	client := pb.NewStreamServiceClient(grpcConnection)
	clientStream, err := client.Stream(ctx)

	//serverStream should be nil
	conn,err := NewConnection(clientStream,nil,
		grpcConnection,client,handler,connectionID,cf)

	if err != nil{
		return err
	}

	comm.Lock()
	comm.connectionMap[connectionID] = conn
	comm.Unlock()

	commLogger.Println("new connection:",connectionID, "are created")

	return nil
}

func (comm *ConnectionManagerImpl) SendStream(message *pb.StreamMessage, errorCallBack OnError, connectionID string){

	//commLogger.Println("Sending data...")

	payload, err := proto.Marshal(message)

	if err != nil{
		commLogger.Println("Marshal error:", err)
		return
	}

	_, pub, err  := comm.crpyto.LoadKey()

	if err != nil{
		commLogger.Println("Load key error:", err)
		return
	}

	sig, err  := comm.crpyto.Sign(payload,auth.DefaultRSAOption)

	if err != nil{
		commLogger.Println("Signing error:", err)
		return
	}

	envelope := &pb.Envelope{}
	envelope.Payload = payload
	envelope.Pubkey = pub.SKI()
	envelope.Signature = sig
	
	conn, ok := comm.connectionMap[connectionID]

	if ok{
		conn.Send(envelope,errorCallBack)
		//commLogger.Println("Sended Envelope:",envelope)
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

func (comm *ConnectionManagerImpl) Stream(stream pb.StreamService_StreamServer) (error) {
	return nil
}