package messaging

import (
	"testing"
	"github.com/magiconair/properties/assert"
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

	assert.Equal(t, err, ErrEmptyNodeId)
}


func TestMessageDispatcher_DeliverLeaderInfo(t *testing.T) {

	tests := map[string]struct{
		input struct{
			nodeId p2p.NodeId
			leader p2p.Leader
		}
		err error
	}{
		"empty node id test":{
			input:struct{
				nodeId p2p.NodeId
				leader p2p.Leader
			}{
				nodeId: p2p.NodeId{
					Id:"1",
				},
				leader:p2p.Leader{
					LeaderId:p2p.LeaderId{
						Id:"1",
					},
				},
			},
			err: ErrEmptyNodeId,
		},
		"empty leader id test":{
			input: struct {
				nodeId p2p.NodeId
				leader p2p.Leader
			}{
				nodeId: p2p.NodeId{
					Id:"1",
				},
				leader: p2p.Leader{
					LeaderId:p2p.LeaderId{
						Id:"1",
					},
				},
			},
			err:ErrEmptyLeaderId,
		},
	}
	publisher := Publisher(func(exchange string, topic string, data interface{}) error {{
		assert.Equal(t, exchange, "Command")
		assert.Equal(t, topic, "message.deliver")
		assert.Equal(t, reflect.TypeOf(data).String(), "p2p.MessageDeliverCommand")

		return nil
	}})
	messageDispatcher := MessageDispatcher{
		publisher:publisher,
	}

	for testName, test := range tests{
		t.Logf("running test case %s", testName)
		err := messageDispatcher.DeliverLeaderInfo(test.input.nodeId,test.input.leader)
		assert.Equal(t, err, test.err)
	}
}

func TestMessageDispatcher_RequestNodeList(t *testing.T) {
	tests := map[string]struct{
		input p2p.NodeId
		err error
	}{
		"empty node id test":{
			input: p2p.NodeId{
				Id:"1",
			},
			err:ErrEmptyNodeId,
		},
	}
	publisher := Publisher(func(exchange string, topic string, data interface{}) error {{
		assert.Equal(t, exchange, "Command")
		assert.Equal(t, topic, "message.deliver")
		assert.Equal(t, reflect.TypeOf(data).String(), "p2p.MessageDeliverCommand")

		return nil
	}})
	messageDispatcher := MessageDispatcher{
		publisher:publisher,
	}

	for testName, test := range tests{
		t.Logf("running test case %s", testName)
		err := messageDispatcher.RequestNodeList(test.input)
		assert.Equal(t, err, test.err)
	}
}

func TestMessageDispatcher_DeliverNodeList(t *testing.T) {
	tests := map[string]struct{
		input struct{
			nodeId p2p.NodeId
			nodeList []p2p.Node
		}
		err error
	}{
		"empty node list test":{
			input: struct {
				nodeId p2p.NodeId
				nodeList []p2p.Node
			}{
				nodeId: p2p.NodeId{
					Id:"1",
				},
				nodeList:[]p2p.Node{},
			},
			err:ErrEmptyNodeList,
		},
		"empty node id test":{
			input: struct {
				nodeId   p2p.NodeId
				nodeList []p2p.Node
			}{
				nodeId: p2p.NodeId{
					Id:"1",
				},
				nodeList:[]p2p.Node{},
			},
			err:ErrEmptyNodeId,
		},
	}
	publisher := Publisher(func(exchange string, topic string, data interface{}) error {{
		assert.Equal(t, exchange, "Command")
		assert.Equal(t, topic, "message.deliver")
		assert.Equal(t, reflect.TypeOf(data).String(), "p2p.MessageDeliverCommand")

		return nil
	}})
	messageDispatcher := MessageDispatcher{
		publisher:publisher,
	}
	for testName, test := range tests{
		t.Logf("running test case %s", testName)
		err := messageDispatcher.DeliverNodeList(test.input.nodeId, test.input.nodeList)
		assert.Equal(t, err, test.err)
	}

}
