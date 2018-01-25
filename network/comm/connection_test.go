package comm

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"net"
	"log"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc"
	pb "it-chain/network/protos"
	"golang.org/x/net/context"
	"fmt"
	"io"
	"time"
)

const (
	ipaddress = "127.0.0.1:5555"
)

var counter = 0

type Mockserver struct {}

func (s *Mockserver) Stream(stream pb.MessageService_StreamServer) (error) {

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
	pb.RegisterMessageServiceServer(s, &Mockserver{})
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

	conn,err  := NewConnection(grpc_conn,nil,"1")

	defer conn.Close()

	if err != nil{
		assert.Fail(t,err.Error())
	}

	assert.NotNil(t,conn)
	assert.NotNil(t,conn.conn)
	assert.NotNil(t,conn.cancel)
	assert.NotNil(t,conn.client)
	assert.NotNil(t,conn.clientStream)

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

	conn,err  := NewConnection(grpc_conn, nil, "1")

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
