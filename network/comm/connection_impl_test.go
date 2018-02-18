package comm

import (
	"testing"
	"github.com/stretchr/testify/assert"
	pb "it-chain/network/protos"
	"fmt"
	"time"
	"golang.org/x/net/context"
	"it-chain/network/comm/mock"
)

const ipaddress = "127.0.0.1:5555"

func TestNewConnection(t *testing.T) {

	handler := func (streamServer pb.StreamService_StreamServer,envelope *pb.Envelope){

	}

	mockServer := &mock.Mockserver{}
	mockServer.Handler = handler

	server, listner := mock.ListenMockServer(mockServer,ipaddress)

	grpc_conn, err := NewConnectionWithAddress(ipaddress,false,nil)

	if err != nil{
		assert.Fail(t,"fail to create connection")
	}

	ctx, cf := context.WithCancel(context.Background())
	client := pb.NewStreamServiceClient(grpc_conn)
	clientStream, err := client.Stream(ctx)

	//serverStream should be nil
	conn,err := NewConnection(clientStream,nil,
		grpc_conn,client,nil,"1",cf)

	defer conn.Close()

	if err != nil{
		assert.Fail(t,err.Error())
	}

	assert.NotNil(t,conn)

	server.Stop()
	listner.Close()
}

func TestConnection_SendWithStream(t *testing.T) {

	counter := 0
	handler := func (streamServer pb.StreamService_StreamServer,envelope *pb.Envelope){
		counter++
	}

	mockServer := &mock.Mockserver{}
	mockServer.Handler = handler

	server, listner := mock.ListenMockServer(mockServer,ipaddress)

	grpc_conn, err := NewConnectionWithAddress(ipaddress,false,nil)

	if err != nil{
		assert.Fail(t,"fail to create connection")
	}

	ctx, cf := context.WithCancel(context.Background())
	client := pb.NewStreamServiceClient(grpc_conn)
	clientStream, err := client.Stream(ctx)

	//serverStream should be nil
	conn,err := NewConnection(clientStream,nil,
		grpc_conn,client,nil,"1",cf)

	defer conn.Close()

	if err != nil{
		assert.Fail(t,err.Error())
	}

	envelope := &pb.Envelope{Signature:[]byte("123")}


	conn.Send(envelope,nil)
	conn.Send(envelope, nil)
	conn.Send(envelope, nil)

	time.Sleep(3*time.Second)

	server.Stop()
	listner.Close()

	assert.Equal(t,3,counter)
}

func TestNewConnectionWithAddress(t *testing.T) {

	conn,err := NewConnectionWithAddress("127.0.0.1:8080",false,nil)
	defer conn.Close()

	if err != nil{
		assert.Fail(t,"fail to connect")
	}

	assert.NotNil(t,conn)
}

func TestConnectionImpl_ReadStream(t *testing.T) {

	counter := 0
	handler := func (streamServer pb.StreamService_StreamServer,envelope *pb.Envelope){
		counter++
	}

	mockServer := &mock.Mockserver{}
	mockServer.Handler = handler

	server, listner := mock.ListenMockServer(mockServer,ipaddress)

	grpc_conn, err := NewConnectionWithAddress(ipaddress,false,nil)

	if err != nil{
		assert.Fail(t,"fail to create connection")
	}

	ctx, cf := context.WithCancel(context.Background())
	client := pb.NewStreamServiceClient(grpc_conn)
	clientStream, err := client.Stream(ctx)

	var receivedMessageCounter = 0

	var MockMessageHandle = func(message OutterMessage){
		receivedMessageCounter ++
	}

	//serverStream should be nil
	conn,err := NewConnection(clientStream,nil,
		grpc_conn,client,MockMessageHandle,"1",cf)

	defer conn.Close()

	if err != nil{
		assert.Fail(t,err.Error())
	}

	envelope := &pb.Envelope{Signature:[]byte("123")}

	fmt.Println(counter)

	conn.Send(envelope,nil)
	conn.Send(envelope, nil)
	conn.Send(envelope, nil)

	time.Sleep(3*time.Second)

	server.Stop()
	listner.Close()

	assert.Equal(t,3,counter)
}
