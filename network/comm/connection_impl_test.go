package comm

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"net"
	"log"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc"
	pb "it-chain/network/protos"
	"fmt"
	"time"
	"io"
	"golang.org/x/net/context"
)

const (
	ipaddress = "127.0.0.1:5555"
)

var counter = 0

type Mockserver struct {}

func (s *Mockserver) Stream(stream pb.StreamService_StreamServer) (error) {

	for {
		message,err := stream.Recv()

		if err == io.EOF {
			return nil
		}

		if err != nil {
			return err
		}
		counter += 1
		fmt.Printf("Received: %d\n", message.String())

		err = stream.Send(message)

		if err !=nil{
			fmt.Printf("Send Error: %d\n", message.String())
		}
	}
}

func (s *Mockserver) Ping(ctx context.Context, in *pb.Empty) (*pb.Empty, error) {
	return &pb.Empty{}, nil
}

func ListenMockServer() (*grpc.Server,net.Listener){

	lis, err := net.Listen("tcp", ipaddress)

	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterStreamServiceServer(s, &Mockserver{})
	// Register reflection service on gRPC server.
	reflection.Register(s)

	go func(){
		if err := s.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
			s.Stop()
			lis.Close()
		}
	}()

	return s,lis
}


func TestNewConnection(t *testing.T) {

	server, listner := ListenMockServer()

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

	counter = 0
	server, listner := ListenMockServer()

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

	fmt.Println(counter)

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

	server, listner := ListenMockServer()

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

	assert.Equal(t,3,receivedMessageCounter)
}
