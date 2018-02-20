package comm

import (
	"google.golang.org/grpc"
	"it-chain/network/protos"
	"sync"
	"golang.org/x/net/context"
	"sync/atomic"
	"errors"
	"it-chain/common"
	"time"
	"google.golang.org/grpc/credentials"
	pb "it-chain/network/protos"
	"it-chain/network/comm/msg"
)

var logger_comm = common.GetLogger("connection.go")

const defaultTimeout = time.Second * 3

type ReceiveMessageHandle func(message msg.OutterMessage)

//직접적으로 grpc를 보내고 받는 역활 수행
type ConnectionImpl struct {
	conn           *grpc.ClientConn
	client         message.StreamServiceClient
	clientStream   message.StreamService_StreamClient
	serverStream   message.StreamService_StreamServer
	cancel         context.CancelFunc
	stopFlag       int32
	connectionID   string
	handle         ReceiveMessageHandle
	outChannl      chan *msg.InnerMessage
	readChannel    chan *message.Envelope
	stopChannel    chan struct{}
	sync.RWMutex
}

type stream interface{
	Send(*pb.Envelope) error
	Recv() (*pb.Envelope, error)
}

//get time from config
//timeOut := viper.GetInt("grpc.timeout")
func NewConnectionWithAddress(peerAddress string,  tslEnabled bool, creds credentials.TransportCredentials) (*grpc.ClientConn, error){

	var opts []grpc.DialOption

	if tslEnabled {
		opts = append(opts, grpc.WithTransportCredentials(creds))
	} else {
		opts = append(opts, grpc.WithInsecure())
	}

	opts = append(opts, grpc.WithTimeout(defaultTimeout))
	conn, err := grpc.Dial(peerAddress, opts...)
	if err != nil {
		return nil, err
	}

	return conn, err
}

//todo channel의 buffer size를 config에서 읽기
func NewConnection(clientStream pb.StreamService_StreamClient, serverStream pb.StreamService_StreamServer, conn *grpc.ClientConn,
			client pb.StreamServiceClient, handle ReceiveMessageHandle, connectionID string, cf context.CancelFunc) (Connection,error){

	connection := &ConnectionImpl{
		clientStream: clientStream,
		serverStream: serverStream,
		cancel: cf,
		client: client,
		conn: conn,
		outChannl: make(chan *msg.InnerMessage,200),
		readChannel: make(chan *message.Envelope,200),
		stopChannel: make(chan struct{},1),
		handle: handle,
		connectionID: connectionID,
	}

	go connection.listen()

	return connection, nil
}

func (conn *ConnectionImpl) toDie() bool {
	return atomic.LoadInt32(&(conn.stopFlag)) == int32(1)
}

func (conn *ConnectionImpl) Send(envelope *message.Envelope, errCallBack func(error)){

	conn.Lock()
	defer conn.Unlock()

	m := &msg.InnerMessage{
		Envelope: envelope,
		OnErr:    errCallBack,
	}

	conn.outChannl <- m
}

func (conn *ConnectionImpl) Close(){

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

func (conn *ConnectionImpl) listen() error{
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
		case message := <-conn.readChannel:
			if conn.handle != nil{
				conn.handle(msg.OutterMessage{Envelope:message,Conn:conn,ConnectionID:conn.connectionID})
			}
		}
	}

	return nil
}

func (conn *ConnectionImpl) getStream() stream{

	conn.Lock()
	defer conn.Unlock()

	if conn.clientStream != nil {
		return conn.clientStream
	}

	if conn.serverStream != nil{
		return conn.serverStream
	}

	return nil
}

func (conn *ConnectionImpl) ReadStream(errChan chan error){

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

		logger_comm.Println("received:",envelope)

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

func (conn *ConnectionImpl) WriteStream(){

	for !conn.toDie() {
		stream := conn.getStream()
		if stream == nil {
			logger_comm.Error(conn.connectionID, "Stream is nil, aborting!")
			return
		}
		select {
		case m := <-conn.outChannl:
			logger_comm.Println("sending", m.Envelope)
			err := stream.Send(m.Envelope)
			if err != nil {
				go m.OnErr(err)
				return
			}
		case stop := <-conn.stopChannel:
			logger_comm.Debug("Closing writing to stream")
			conn.stopChannel <- stop
			return
		}
	}
}