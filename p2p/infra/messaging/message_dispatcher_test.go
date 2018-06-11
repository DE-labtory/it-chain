package messaging

import (
	"testing"
	"github.com/magiconair/properties/assert"
	"errors"
	"reflect"
	"github.com/it-chain/it-chain-Engine/p2p"
)


func TestMessageDispatcher_RequestLeaderInfo(t *testing.T) {
	publisher := Publisher(func(exchange string, topic string, data interface{}) error {{
		assert.Equal(t, exchange, "Command")
		assert.Equal(t, topic, "message.deliver")
		assert.Equal(t, reflect.TypeOf(data).String(), "p2p.MessageDeliverCommand")

		return nil
	}})
	messageDispatcher := MessageDispatcher{
		publisher:publisher,
	}

	nodeId := p2p.NodeId{
		"1",
	}

	err := messageDispatcher.RequestLeaderInfo(nodeId)

	assert.Equal(t, err, errors.New("Empty nodeid proposed"))
}


func TestMessageDispatcher_DeliverLeaderInfo(t *testing.T) {
	publisher := Publisher(func(exchange string, topic string, data interface{}) error {{
		assert.Equal(t, exchange, "Command")
		assert.Equal(t, topic, "message.deliver")
		assert.Equal(t, reflect.TypeOf(data).String(), "p2p.MessageDeliverCommand")

		return nil
	}})
	messageDispatcher := MessageDispatcher{
		publisher:publisher,
	}
	leader := p2p.Leader{
		LeaderId: p2p.LeaderId{"1"},
	}
	nodeId := p2p.NodeId{
		"1",
	}
	err := messageDispatcher.DeliverLeaderInfo(nodeId,leader)

	assert.Equal(t, err, errors.New("Empty nodeid proposed"))
}

func TestMessageDispatcher_RequestNodeList(t *testing.T) {
	publisher := Publisher(func(exchange string, topic string, data interface{}) error {{
		assert.Equal(t, exchange, "Command")
		assert.Equal(t, topic, "message.deliver")
		assert.Equal(t, reflect.TypeOf(data).String(), "p2p.MessageDeliverCommand")

		return nil
	}})
	messageDispatcher := MessageDispatcher{
		publisher:publisher,
	}
	nodeId := p2p.NodeId{
		"1",
	}
	err := messageDispatcher.RequestNodeList(nodeId)

	assert.Equal(t, err, errors.New("Empty nodeid proposed"))
}

func TestMessageDispatcher_DeliverNodeList(t *testing.T) {
	publisher := Publisher(func(exchange string, topic string, data interface{}) error {{
		assert.Equal(t, exchange, "Command")
		assert.Equal(t, topic, "message.deliver")
		assert.Equal(t, reflect.TypeOf(data).String(), "p2p.MessageDeliverCommand")

		return nil
	}})
	messageDispatcher := MessageDispatcher{
		publisher:publisher,
	}
	nodeId := p2p.NodeId{
		"1",
	}
	nodeList := []p2p.Node{}
	err := messageDispatcher.DeliverNodeList(nodeId, nodeList)

	assert.Equal(t, err, errors.New("Empty nodeid proposed"))
	assert.Equal(t, err, errors.New("Empty node list proposed"))
}
