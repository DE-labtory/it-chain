package comm

import (
	"testing"
	"github.com/stretchr/testify/assert"
	pb "it-chain/network/protos"
	"it-chain/auth"
	"it-chain/network/comm/mock"
	"os"
	"time"
	"log"
	"github.com/golang/protobuf/proto"
	"it-chain/network/comm/msg"
)

//todo connection manager_impl test 모두 수정
func TestCommImpl_CreateStreamClientConn(t *testing.T) {

	counter := 0
	handler := func (streamServer pb.StreamService_StreamServer,envelope *pb.Envelope){
		counter++
	}

	mockServer := &mock.Mockserver{}
	mockServer.Handler = handler

	server1, listner1 := mock.ListenMockServer(mockServer,"127.0.0.1:5555")
	server2, listner2 := mock.ListenMockServer(mockServer,"127.0.0.1:6666")

	cryp, err := auth.NewCrypto("./KeyRepository", &auth.RSAKeyGenOpts{})
	defer os.RemoveAll("./KeyRepository")
	assert.NoError(t, err)

	comm := NewConnectionManagerImpl(cryp)
	comm.CreateStreamClientConn("1","127.0.0.1:5555")
	comm.CreateStreamClientConn("2","127.0.0.1:6666")

	defer func(){
		server1.Stop()
		listner1.Close()
		server2.Stop()
		listner2.Close()
	}()

	assert.NotNil(t,comm)
	assert.Equal(t,2,comm.Size())
}

//
func TestCommImpl_Send(t *testing.T) {

	counter := 0
	handler := func (streamServer pb.StreamService_StreamServer,envelope *pb.Envelope){
		counter++
	}

	mockServer := &mock.Mockserver{}
	mockServer.Handler = handler

	server1, listner1 := mock.ListenMockServer(mockServer,"127.0.0.1:5555")
	server2, listner2 := mock.ListenMockServer(mockServer,"127.0.0.1:6666")

	cryp, err := auth.NewCrypto("./KeyRepository", &auth.RSAKeyGenOpts{})
	defer os.RemoveAll("./KeyRepository")
	assert.NoError(t, err)

	comm := NewConnectionManagerImpl(cryp)
	comm.CreateStreamClientConn("1","127.0.0.1:5555")
	comm.CreateStreamClientConn("2","127.0.0.1:6666")

	message := &pb.StreamMessage{}
	message.Content = &pb.StreamMessage_ConnectionEstablish{
		ConnectionEstablish: &pb.ConnectionEstablish{},
	}

	comm.SendStream(message,nil, "2")
	comm.SendStream(message, nil, "2")

	defer func(){
		server1.Stop()
		listner1.Close()
		server2.Stop()
		listner2.Close()
		comm.Stop()
	}()

	time.Sleep(3*time.Second)

	assert.Equal(t,2,counter)
}

func TestCommImpl_Stop(t *testing.T) {

	counter := 0
	handler := func (streamServer pb.StreamService_StreamServer,envelope *pb.Envelope){
		counter++
	}

	mockServer := &mock.Mockserver{}
	mockServer.Handler = handler

	server1, listner1 := mock.ListenMockServer(mockServer,"127.0.0.1:5555")
	server2, listner2 := mock.ListenMockServer(mockServer,"127.0.0.1:6666")

	cryp, err := auth.NewCrypto("./KeyRepository", &auth.RSAKeyGenOpts{})
	defer os.RemoveAll("./KeyRepository")
	assert.NoError(t, err)

	comm := NewConnectionManagerImpl(cryp)
	comm.CreateStreamClientConn("1","127.0.0.1:5555")
	comm.CreateStreamClientConn("2","127.0.0.1:6666")

	defer func(){
		server1.Stop()
		listner1.Close()
		server2.Stop()
		listner2.Close()
	}()

	comm.Stop()

	assert.Equal(t,0,comm.Size())
}
//

func TestCommImpl_Close(t *testing.T) {

	counter := 0
	handler := func (streamServer pb.StreamService_StreamServer,envelope *pb.Envelope){
		counter++
	}

	mockServer := &mock.Mockserver{}
	mockServer.Handler = handler

	server1, listner1 := mock.ListenMockServer(mockServer,"127.0.0.1:5555")
	server2, listner2 := mock.ListenMockServer(mockServer,"127.0.0.1:6666")

	cryp, err := auth.NewCrypto("./KeyRepository", &auth.RSAKeyGenOpts{})
	defer os.RemoveAll("./KeyRepository")
	assert.NoError(t, err)

	comm := NewConnectionManagerImpl(cryp)
	comm.CreateStreamClientConn("1","127.0.0.1:5555")
	comm.CreateStreamClientConn("2","127.0.0.1:6666")

	defer func(){
		server1.Stop()
		listner1.Close()
		server2.Stop()
		listner2.Close()
	}()

	comm.Close("1")

	assert.Equal(t,1,comm.Size())

	comm.Stop()
}

func TestConnectionManagerImpl_Stream(t *testing.T) {

	cryp, err := auth.NewCrypto("./KeyRepository", &auth.RSAKeyGenOpts{})
	defer os.RemoveAll("./KeyRepository")
	assert.NoError(t, err)

	comm1 := NewConnectionManagerImpl(cryp)
	comm := NewConnectionManagerImpl(cryp)

	var onConnectionHandler = func(conn Connection, peer pb.Peer){
		log.Print("Successfully create connection")
		assert.Equal(t,"1",peer.PeerID)
		assert.Equal(t,comm1.Size(),1)
		log.Print("End")
	}

	comm.SetOnConnectHandler(onConnectionHandler)
	comm1.SetOnConnectHandler(onConnectionHandler)

	server1, listner1 := mock.ListenMockServer(comm,"127.0.0.1:5555")
	server2, listner2 := mock.ListenMockServer(comm1,"127.0.0.1:6666")


	defer func(){
		server1.Stop()
		listner1.Close()
		server2.Stop()
		listner2.Close()
	}()

	var receiveHandler = func(message msg.OutterMessage){

		log.Println("receivedHandler got message")

		sm := &pb.StreamMessage{}
		sm.Content = &pb.StreamMessage_Peer{
			Peer: &pb.Peer{PeerID:"1"},
		}

		payload, err := proto.Marshal(sm)

		if err != nil{
			log.Fatalln("error")
		}

		envelope := &pb.Envelope{}
		envelope.Payload = payload

		var errorCallback = func (err error){
			log.Println(err.Error())
		}

		message.Respond(envelope,errorCallback)
		log.Println("respond message")
	}

	comm.Subscribe("mockReceiver",receiveHandler)

	comm.CreateStreamClientConn("1","127.0.0.1:6666")

	time.Sleep(3*time.Second)
}