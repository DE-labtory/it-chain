package comm

import (
	"google.golang.org/grpc"
	"it-chain/network/protos"
	"sync"
	"golang.org/x/net/context"
	"sync/atomic"
)

//직접적으로 grpc를 보내는 역활 수행
//todo client 와 server connection을 합칠 것인지 분리 할 것인지 생각 지금은 client만을 고려한 구조체
type Connection struct {
	conn         *grpc.ClientConn
	client       message.MessageServiceClient
	clientStream message.MessageService_StreamClient
	cancel       context.CancelFunc
	stopFlag     int32
	sync.RWMutex
}

func NewConnection(conn *grpc.ClientConn) (*Connection,error){

	ctx, cf := context.WithCancel(context.Background())
	client := message.NewMessageServiceClient(conn)
	clientStream, err := client.Stream(ctx)

	if err != nil{
		conn.Close()
		return nil, err
	}

	return &Connection{
		clientStream: clientStream,
		cancel: cf,
		client: client,
		conn: conn,
	}, nil
}

func (conn *Connection) toDie() bool {
	return atomic.LoadInt32(&(conn.stopFlag)) == int32(1)
}

func (conn *Connection) SendWithStream(envelop *message.Envelope) error{

	err := conn.clientStream.Send(envelop)
	if err != nil{

		return err
	}

	return nil
}


func (conn *Connection) Close(){

	conn.Lock()

	if conn.conn != nil {
		conn.conn.Close()
	}

	if conn.clientStream != nil{
		conn.clientStream.CloseSend()
	}

	if conn.cancel != nil{
		conn.cancel()
	}

	conn.Unlock()
}