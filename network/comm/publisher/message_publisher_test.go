package publisher

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"time"
	"fmt"
	"github.com/it-chain/it-chain-Engine/network/protos"
	pb "github.com/it-chain/it-chain-Engine/network/protos"
	"github.com/golang/protobuf/proto"
	"github.com/stretchr/testify/mock"
	"github.com/it-chain/it-chain-Engine/auth"
	"github.com/it-chain/it-chain-Engine/network/comm/msg"
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

func (mc MockCrypto) GetKey() (pri, pub auth.Key, err error){
	return nil,nil,nil
}

func (mc MockCrypto) LoadKey() (pri, pub auth.Key, err error){
	return nil,nil,nil
}

func MakeEnvelopeHavingPeerTable() *message.Envelope{

	peerTable := &pb.StreamMessage_PeerTable{}
	peerTable.PeerTable = &pb.PeerTable{}

	message := &message.StreamMessage{}
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
	messageHandler := NewMessagePublisher(crpyto)

	assert.NotNil(t,messageHandler.subscribers)
	assert.NotNil(t,messageHandler.crpyto)
}

func TestMessagePublisher_AddSubscriber(t *testing.T) {
	crpyto := MockCrypto{}
	messageHandler := NewMessagePublisher(crpyto)

	subfunc := func(message msg.OutterMessage){

	}

	err := messageHandler.AddSubscriber("mocksub",subfunc)

	assert.Equal(t,1,len(messageHandler.subscribers))
	assert.NoError(t,err)
}

func TestMessagePublisher_ReceivedMessageHandle(t *testing.T) {

	count := 1
	crpyto := MockCrypto{}
	messageHandler := NewMessagePublisher(crpyto)

	subfunc := func(message msg.OutterMessage){
		count ++
		fmt.Println("published")
	}

	err := messageHandler.AddSubscriber("mocksub",subfunc)

	envelop := MakeEnvelopeHavingPeerTable()

	messageHandler.ReceivedMessageHandle(msg.OutterMessage{Envelope:envelop})

	assert.NoError(t,err)

	time.Sleep(1*time.Second)

	assert.Equal(t,count,2)
}

func TestMessagePublisher_MultipleReceivedMessageHandle(t *testing.T) {

	count := 1
	crpyto := MockCrypto{}
	messageHandler := NewMessagePublisher(crpyto)

	subfunc := func(message msg.OutterMessage){
		count ++
		fmt.Println("published")
	}

	subfunc2 := func(message msg.OutterMessage){
		count ++
		fmt.Println("published")
	}

	err := messageHandler.AddSubscriber("mocksub",subfunc)
	err = messageHandler.AddSubscriber("mocksub2",subfunc2)

	envelop := MakeEnvelopeHavingPeerTable()

	messageHandler.ReceivedMessageHandle(msg.OutterMessage{Envelope:envelop})

	assert.NoError(t,err)

	time.Sleep(1*time.Second)

	assert.Equal(t,count,3)
}
//
func TestMessagePublisher_ReceivedMessageHandleError(t *testing.T) {

	count := 1
	crpyto := MockCrypto{}
	messageHandler := NewMessagePublisher(crpyto)

	subfunc := func(message msg.OutterMessage){
		count ++
		fmt.Println("published")
	}

	err := messageHandler.AddSubscriber("mocksub",subfunc)
	err = messageHandler.AddSubscriber("mocksub",subfunc)

	assert.Error(t,err)
}
