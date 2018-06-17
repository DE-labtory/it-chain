package adapter

import (
	"reflect"
	"testing"

	"github.com/it-chain/it-chain-Engine/p2p"
	"github.com/it-chain/it-chain-Engine/p2p/infra/messaging"
	"github.com/magiconair/properties/assert"
)

func TestMessageDispatcher_RequestLeaderInfo(t *testing.T) {

	tests := map[string]struct {
		input struct {
			nodeId p2p.NodeId
		}
		err error
	}{
		"success": {
			input: struct {
				nodeId p2p.NodeId
			}{
				nodeId: p2p.NodeId{
					Id: "1",
				},
			},
			err: nil,
		},
		"empty node id test": {
			input: struct {
				nodeId p2p.NodeId
			}{
				nodeId: p2p.NodeId{},
			},
			err: messaging.ErrEmptyNodeId,
		},
	}

	publisher := messaging.Publisher(func(exchange string, topic string, data interface{}) error {
		assert.Equal(t, exchange, "Command")
		assert.Equal(t, topic, "message.deliver")
		assert.Equal(t, reflect.TypeOf(data).String(), "p2p.MessageDeliverCommand")

		return nil
	})

	messageDispatcher := messaging.NewMessageDispatcher(publisher)

	for testName, test := range tests {
		t.Logf("running test case %s", testName)
		err := messageDispatcher.RequestLeaderInfo(test.input.nodeId)
		assert.Equal(t, err, test.err)
	}

}

func TestMessageDispatcher_DeliverLeaderInfo(t *testing.T) {

	tests := map[string]struct {
		input struct {
			nodeId p2p.NodeId
			leader p2p.Leader
		}
		err error
	}{
		"empty node id test": {
			input: struct {
				nodeId p2p.NodeId
				leader p2p.Leader
			}{
				nodeId: p2p.NodeId{},
				leader: p2p.Leader{
					LeaderId: p2p.LeaderId{
						Id: "1",
					},
				},
			},
			err: messaging.ErrEmptyNodeId,
		},
		"empty leader id test": {
			input: struct {
				nodeId p2p.NodeId
				leader p2p.Leader
			}{
				nodeId: p2p.NodeId{
					Id: "1",
				},
				leader: p2p.Leader{
					LeaderId: p2p.LeaderId{},
				},
			},
			err: messaging.ErrEmptyLeaderId,
		},
		"success": {
			input: struct {
				nodeId p2p.NodeId
				leader p2p.Leader
			}{
				nodeId: p2p.NodeId{
					Id: "1",
				},
				leader: p2p.Leader{
					LeaderId: p2p.LeaderId{
						Id: "1",
					},
				},
			},
			err: nil,
		},
	}
	publisher := messaging.Publisher(func(exchange string, topic string, data interface{}) error {
		{
			assert.Equal(t, exchange, "Command")
			assert.Equal(t, topic, "message.deliver")
			assert.Equal(t, reflect.TypeOf(data).String(), "p2p.MessageDeliverCommand")

			return nil
		}
	})

	messageDispatcher := messaging.NewMessageDispatcher(publisher)

	for testName, test := range tests {
		t.Logf("running test case %s", testName)
		err := messageDispatcher.DeliverLeaderInfo(test.input.nodeId, test.input.leader)
		assert.Equal(t, err, test.err)
	}
}

func TestMessageDispatcher_RequestNodeList(t *testing.T) {

	tests := map[string]struct {
		input p2p.NodeId
		err   error
	}{
		"empty node id test": {
			input: p2p.NodeId{
				Id: "1",
			},
			err: nil,
		},
		"success": {
			input: p2p.NodeId{},
			err:   messaging.ErrEmptyNodeId,
		},
	}
	publisher := messaging.Publisher(func(exchange string, topic string, data interface{}) error {
		{
			assert.Equal(t, exchange, "Command")
			assert.Equal(t, topic, "message.deliver")
			assert.Equal(t, reflect.TypeOf(data).String(), "p2p.MessageDeliverCommand")

			return nil
		}
	})

	messageDispatcher := messaging.NewMessageDispatcher(publisher)

	for testName, test := range tests {
		t.Logf("running test case %s", testName)
		err := messageDispatcher.RequestNodeList(test.input)
		assert.Equal(t, err, test.err)
	}
}

func TestMessageDispatcher_DeliverNodeList(t *testing.T) {
	tests := map[string]struct {
		input struct {
			nodeId   p2p.NodeId
			nodeList []p2p.Node
		}
		err error
	}{
		"empty node list test": {
			input: struct {
				nodeId   p2p.NodeId
				nodeList []p2p.Node
			}{
				nodeId: p2p.NodeId{
					Id: "1",
				},
				nodeList: []p2p.Node{},
			},
			err: messaging.ErrEmptyNodeList,
		},
		"empty node id test": {
			input: struct {
				nodeId   p2p.NodeId
				nodeList []p2p.Node
			}{
				nodeId:   p2p.NodeId{},
				nodeList: []p2p.Node{},
			},
			err: messaging.ErrEmptyNodeId,
		},
		"success": {
			input: struct {
				nodeId   p2p.NodeId
				nodeList []p2p.Node
			}{
				nodeId: p2p.NodeId{
					Id: "1",
				},
				nodeList: []p2p.Node{
					p2p.Node{
						NodeId: p2p.NodeId{
							Id: "1",
						},
						IpAddress: "123",
					},
				},
			},
			err: nil,
		},
	}
	publisher := messaging.Publisher(func(exchange string, topic string, data interface{}) error {
		{
			assert.Equal(t, exchange, "Command")
			assert.Equal(t, topic, "message.deliver")
			assert.Equal(t, reflect.TypeOf(data).String(), "p2p.MessageDeliverCommand")

			return nil
		}
	})

	messageDispatcher := messaging.NewMessageDispatcher(publisher)

	for testName, test := range tests {
		t.Logf("running test case %s", testName)
		err := messageDispatcher.DeliverNodeList(test.input.nodeId, test.input.nodeList)
		assert.Equal(t, err, test.err)
	}

}
