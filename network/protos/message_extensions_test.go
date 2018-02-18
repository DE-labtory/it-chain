package message

import (
	"testing"
	"github.com/golang/protobuf/proto"
	"github.com/magiconair/properties/assert"
)

func TestEnvelope_GetMessage(t *testing.T) {

	sm := &StreamMessage{}
	sm.Content = &StreamMessage_Peer{
		Peer: &Peer{PeerID:"1"},
	}

	payload, _ := proto.Marshal(sm)

	m := &StreamMessage{}
	proto.Unmarshal(payload,m)


	envelope := &Envelope{}
	envelope.Payload = payload

	message , _ := envelope.GetMessage()

	assert.Equal(t, message.GetPeer().PeerID,"1")
}
