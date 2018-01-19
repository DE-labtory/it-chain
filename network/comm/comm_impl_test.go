package comm

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"time"
	"net"
	"log"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc"
	pb "it-chain/network/protos"
	"it-chain/service/domain"
)

func MockCreateNewPeerInfo(peerID string) *domain.PeerInfo{

	return  &domain.PeerInfo{
		PeerID: peerID,
		Port: "5555",
		IpAddress: "127.0.0.1",
		HeartBeat: 1,
		TimeStamp: time.Now(),
	}
}

func ListenMockServerWithPeer(peer domain.PeerInfo) (*grpc.Server,net.Listener){
	lis, err := net.Listen("tcp", peer.GetEndPoint())
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

func TestCommImpl_CreateConn(t *testing.T) {

	peer1 := MockCreateNewPeerInfo("test1")
	server1, listner1 := ListenMockServerWithPeer(*peer1)

	peer2 := MockCreateNewPeerInfo("test2")
	peer2.Port = "6666"
	server2, listner2 := ListenMockServerWithPeer(*peer2)

	defer func(){
		server1.Stop()
		listner1.Close()
		server2.Stop()
		listner2.Close()
	}()

	comm := NewCommImpl()
	comm.CreateConn(*peer1)
	comm.CreateConn(*peer2)

	assert.NotNil(t,comm)
	assert.Equal(t,2,comm.Size())
}

func TestCommImpl_Send(t *testing.T) {
	counter = 0

	peer1 := MockCreateNewPeerInfo("test1")
	server1, listner1 := ListenMockServerWithPeer(*peer1)

	peer2 := MockCreateNewPeerInfo("test2")
	peer2.Port = "6666"
	server2, listner2 := ListenMockServerWithPeer(*peer2)

	defer func(){
		server1.Stop()
		listner1.Close()
		server2.Stop()
		listner2.Close()
	}()

	comm := NewCommImpl()
	comm.CreateConn(*peer1)
	comm.CreateConn(*peer2)

	envelope := &pb.Envelope{Signature:[]byte("123")}

	comm.Send(*envelope,nil, *peer1, *peer2)
	comm.Send(*envelope,nil, *peer2)

	time.Sleep(3*time.Second)

	assert.Equal(t,3,counter)
}

func TestCommImpl_Stop(t *testing.T) {
	counter = 0
	peer1 := MockCreateNewPeerInfo("test1")
	server1, listner1 := ListenMockServerWithPeer(*peer1)

	peer2 := MockCreateNewPeerInfo("test2")
	peer2.Port = "6666"
	server2, listner2 := ListenMockServerWithPeer(*peer2)

	defer func(){
		server1.Stop()
		listner1.Close()
		server2.Stop()
		listner2.Close()
	}()

	comm := NewCommImpl()
	comm.CreateConn(*peer1)
	comm.CreateConn(*peer2)

	comm.Stop()

	assert.Equal(t,0,comm.Size())
}

func TestCommImpl_Close(t *testing.T) {

	peer1 := MockCreateNewPeerInfo("test1")
	server1, listner1 := ListenMockServerWithPeer(*peer1)

	peer2 := MockCreateNewPeerInfo("test2")
	peer2.Port = "6666"
	server2, listner2 := ListenMockServerWithPeer(*peer2)

	defer func(){
		server1.Stop()
		listner1.Close()
		server2.Stop()
		listner2.Close()
	}()

	comm := NewCommImpl()
	comm.CreateConn(*peer1)
	comm.CreateConn(*peer2)

	comm.Close(*peer1)

	assert.Equal(t,1,comm.Size())
}
