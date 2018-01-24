package event

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"it-chain/service/peer/event"
)

func TestNewMessageHandler(t *testing.T) {
	messageHandler := NewMessageHandler(event.MessageTypes)

	assert.NotNil(t,messageHandler.bus)
	assert.Equal(t,len(messageHandler.topicMap),1)
	assert.Equal(t,messageHandler.topicMap[event.UpdatePeerTable],event.UpdatePeerTable)
}