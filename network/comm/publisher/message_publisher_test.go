package publisher

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"it-chain/service/peer/event"
	"it-chain/network/comm"
	"time"
	"fmt"
	"it-chain/network/protos"
	pb "it-chain/network/protos"
	"github.com/golang/protobuf/proto"
)

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
	messageHandler := NewMessagePublisher(event.MessageTypes)

	assert.NotNil(t,messageHandler.bus)
	assert.Equal(t,len(messageHandler.topicMap),1)
	assert.Equal(t,messageHandler.topicMap[event.UpdatePeerTable],event.UpdatePeerTable)
}

func TestMessagePublisher_AddSubscriber(t *testing.T) {
	messageHandler := NewMessagePublisher(event.MessageTypes)

	subfunc := func(message comm.OutterMessage){

	}

	err := messageHandler.AddSubscriber(event.UpdatePeerTable,subfunc,true)
	assert.NoError(t,err)
}

func TestMessagePublisher_ReceivedMessageHandle(t *testing.T) {

	count := 1
	messageHandler := NewMessagePublisher(event.MessageTypes)

	subfunc := func(message comm.OutterMessage){
		count ++
		fmt.Println("published")
	}

	err := messageHandler.AddSubscriber(event.UpdatePeerTable,subfunc,true)

	envelop := MakeEnvelopeHavingPeerTable()

	messageHandler.ReceivedMessageHandle(comm.OutterMessage{Envelope:envelop})

	assert.NoError(t,err)

	time.Sleep(1*time.Second)

	assert.Equal(t,count,2)
}

func TestMessagePublisher_MultipleReceivedMessageHandle(t *testing.T) {

	count := 1
	messageHandler := NewMessagePublisher(event.MessageTypes)

	subfunc := func(message comm.OutterMessage){
		count ++
		fmt.Println("published")
	}

	err := messageHandler.AddSubscriber(event.UpdatePeerTable,subfunc,true)
	err = messageHandler.AddSubscriber(event.UpdatePeerTable,subfunc,true)

	envelop := MakeEnvelopeHavingPeerTable()

	messageHandler.ReceivedMessageHandle(comm.OutterMessage{Envelope:envelop})

	assert.NoError(t,err)

	time.Sleep(1*time.Second)

	assert.Equal(t,count,3)
}

func TestMessagePublisher_ReceivedMessageHandleError(t *testing.T) {

	count := 1
	messageHandler := NewMessagePublisher(event.MessageTypes)

	subfunc := func(message comm.OutterMessage){
		count ++
		fmt.Println("published")
	}

	err := messageHandler.AddSubscriber("hello",subfunc,true)

	assert.Error(t,err)
}