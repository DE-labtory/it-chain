package publisher

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"it-chain/network/comm"
	"time"
	"fmt"
	"it-chain/network/protos"
	pb "it-chain/network/protos"
	"github.com/golang/protobuf/proto"
	"it-chain/network"
	"github.com/stretchr/testify/mock"
	"it-chain/auth"
)

type MockCrypto struct{
	mock.Mock
}

func (mc MockCrypto) Sign(digest []byte, opts auth.SignerOpts) (signature []byte, err error){

	return []byte("asd"),nil
}

func (mc MockCrypto) Verify(key auth.Key, signature, digest []byte, opts auth.SignerOpts) (valid bool, err error){
	return true,nil
}

func (mc MockCrypto) GenerateKey(opts auth.KeyGenOpts) (pri, pub auth.Key, err error){
	return nil,nil,nil
}

func (mc MockCrypto) LoadKey() (pri, pub auth.Key, err error){
	return nil,nil,nil
}

func MakeEnvelopeHavingPeerTable() *message.Envelope{

	peerTable := &pb.Message_PeerTable{}
	peerTable.PeerTable = &pb.PeerTable{}

	message := &message.Message{}
	message.Content = peerTable

	envelope := &pb.Envelope{}
	payload, err := proto.Marshal(message)

	if err !=nil{

	}

	envelope.Payload = payload

	return envelope
}

func TestNewMessageHandler(t *testing.T) {
	crpyto := MockCrypto{}
	messageHandler := NewMessagePublisher(network.MessageTypes,crpyto)

	assert.NotNil(t,messageHandler.bus)
	assert.Equal(t,len(messageHandler.topicMap),1)
	assert.Equal(t,messageHandler.topicMap[network.UpdatePeerTable],network.UpdatePeerTable)
}

func TestMessagePublisher_AddSubscriber(t *testing.T) {
	crpyto := MockCrypto{}
	messageHandler := NewMessagePublisher(network.MessageTypes,crpyto)

	subfunc := func(message comm.OutterMessage){

	}

	err := messageHandler.AddSubscriber(network.UpdatePeerTable,subfunc,true)
	assert.NoError(t,err)
}

func TestMessagePublisher_ReceivedMessageHandle(t *testing.T) {

	count := 1
	crpyto := MockCrypto{}
	messageHandler := NewMessagePublisher(network.MessageTypes,crpyto)

	subfunc := func(message comm.OutterMessage){
		count ++
		fmt.Println("published")
	}

	err := messageHandler.AddSubscriber(network.UpdatePeerTable,subfunc,true)

	envelop := MakeEnvelopeHavingPeerTable()

	messageHandler.ReceivedMessageHandle(comm.OutterMessage{Envelope:envelop})

	assert.NoError(t,err)

	time.Sleep(1*time.Second)

	assert.Equal(t,count,2)
}

func TestMessagePublisher_MultipleReceivedMessageHandle(t *testing.T) {

	count := 1
	crpyto := MockCrypto{}
	messageHandler := NewMessagePublisher(network.MessageTypes,crpyto)

	subfunc := func(message comm.OutterMessage){
		count ++
		fmt.Println("published")
	}

	err := messageHandler.AddSubscriber(network.UpdatePeerTable,subfunc,true)
	err = messageHandler.AddSubscriber(network.UpdatePeerTable,subfunc,true)

	envelop := MakeEnvelopeHavingPeerTable()

	messageHandler.ReceivedMessageHandle(comm.OutterMessage{Envelope:envelop})

	assert.NoError(t,err)

	time.Sleep(1*time.Second)

	assert.Equal(t,count,3)
}

func TestMessagePublisher_ReceivedMessageHandleError(t *testing.T) {

	count := 1
	crpyto := MockCrypto{}
	messageHandler := NewMessagePublisher(network.MessageTypes,crpyto)

	subfunc := func(message comm.OutterMessage){
		count ++
		fmt.Println("published")
	}

	err := messageHandler.AddSubscriber("hello",subfunc,true)

	assert.Error(t,err)
}