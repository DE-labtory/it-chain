package comm

import (
	"google.golang.org/grpc"
	"it-chain/network/protos"
	"sync"
	"golang.org/x/net/context"
	"sync/atomic"
	"errors"
	"it-chain/common"
)

var logger_comm = common.GetLogger("connection.go")

type ReceiveMessageHandle func(message outterMessage)

//직접적으로 grpc를 보내는 역활 수행
//todo client 와 server connection을 합칠 것인지 분리 할 것인지 생각 지금은 client만을 고려한 구조체
type Connection struct {
	conn           *grpc.ClientConn
	client         message.MessageServiceClient
	clientStream   message.MessageService_StreamClient
	cancel         context.CancelFunc
	stopFlag       int32
	connectionID   string
	handle        ReceiveMessageHandle
	outChannl      chan *innerMessage
	readChannel    chan *message.Envelope
	stopChannel    chan struct{}
	sync.RWMutex
}

//todo channel의 buffer size를 config에서 읽기
func NewConnection(conn *grpc.ClientConn, handle ReceiveMessageHandle,connectionID string) (*Connection,error){

	ctx, cf := context.WithCancel(context.Background())
	client := message.NewMessageServiceClient(conn)
	clientStream, err := client.Stream(ctx)

	if err != nil{
		conn.Close()
		return nil, err
	}

	connection := &Connection{
		clientStream: clientStream,
		cancel: cf,
		client: client,
		conn: conn,
		outChannl: make(chan *innerMessage,200),
		readChannel: make(chan *message.Envelope,200),
		stopChannel: make(chan struct{},1),
		handle: handle,
		connectionID: connectionID,
	}

	go connection.listen()

	return connection, nil
}

func (conn *Connection) toDie() bool {
	return atomic.LoadInt32(&(conn.stopFlag)) == int32(1)
}

func (conn *Connection) Send(envelope *message.Envelope, errCallBack func(error)){

	conn.Lock()
	defer conn.Unlock()

	m := &innerMessage{
		envelope: envelope,
		onErr:    errCallBack,
	}

	conn.outChannl <- m
}

func (conn *Connection) Close(){

	if conn.toDie() {
		return
	}

	amIFirst := atomic.CompareAndSwapInt32(&conn.stopFlag, int32(0), int32(1))

	if !amIFirst {
		return
	}

	conn.stopChannel <- struct{}{}

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

func (conn *Connection) listen() error{
	errChan := make(chan error, 1)

	go conn.ReadStream(errChan)
	go conn.WriteStream()

	for !conn.toDie(){
		select{
			case stop := <-conn.stopChannel:
				conn.stopChannel <- stop
				return nil
			case err := <-errChan:
				return err
			case msg := <-conn.readChannel:
				conn.handle(outterMessage{msg,conn.connectionID})
		}
	}

	return nil
}

func (conn *Connection) getStream() message.MessageService_StreamClient{

	conn.Lock()
	defer conn.Unlock()

	if conn.clientStream != nil {
		return conn.clientStream
	}

	return nil
}

func (conn *Connection) ReadStream(errChan chan error){

	defer func() {
		recover()
	}()

	for !conn.toDie() {
		stream := conn.getStream()

		if stream == nil {
			logger_comm.Error(conn.connectionID, "Stream is nil, aborting!")
			errChan <- errors.New("Stream is nil")
			return
		}

		envelope, err := stream.Recv()

		if conn.toDie() {
			logger_comm.Debug(conn.connectionID, "canceling read because closing")
			return
		}

		if err != nil {
			errChan <- err
			logger_comm.Errorln(conn.connectionID, "Got error, aborting:", err)
			return
		}

		conn.readChannel <- envelope
	}
}

func (conn *Connection) WriteStream(){

	for !conn.toDie() {
		stream := conn.getStream()
		if stream == nil {
			logger_comm.Error(conn.connectionID, "Stream is nil, aborting!")
			return
		}
		select {
			case m := <-conn.outChannl:
				err := stream.Send(m.envelope)
				if err != nil {
					go m.onErr(err)
					return
				}
			case stop := <-conn.stopChannel:
				logger_comm.Debug("Closing writing to stream")
				conn.stopChannel <- stop
				return
		}
	}

}

type innerMessage struct{
	envelope *message.Envelope
	onErr    func(error)
}

type outterMessage struct{
	envelope *message.Envelope
	connectionID string
}