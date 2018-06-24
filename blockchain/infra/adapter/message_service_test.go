package adapter_test

import (
	"testing"

	"reflect"

	"github.com/it-chain/it-chain-Engine/blockchain/infra/adapter"
	"github.com/it-chain/it-chain-Engine/p2p"
	"github.com/stretchr/testify/assert"
)

func TestMessageService_RequestBlock(t *testing.T) {

	tests := map[string]struct {
		input struct {
			nodeId p2p.NodeId
			height uint64
		}
		err error
	}{
		"success: request block": {
			input: struct {
				nodeId p2p.NodeId
				height uint64
			}{
				nodeId: p2p.NodeId{
					Id: "1",
				},
				height: uint64(0),
			},
			err: nil,
		},
		"fail: empty node id": {
			input: struct {
				nodeId p2p.NodeId
				height uint64
			}{
				nodeId: p2p.NodeId{},
				height: uint64(0),
			},
			err: adapter.ErrEmptyNodeId,
		},
	}

	publish := func(exchange string, topic string, data interface{}) error {
		assert.Equal(t, exchange, "Command")
		assert.Equal(t, topic, "message.deliver")
		assert.Equal(t, reflect.TypeOf(data).String(), "blockchain.MessageDeliverCommand")

		return nil
	}

	messageService := adapter.NewMessageService(publish)

	for testName, test := range tests {
		t.Logf("running test case %s", testName)
		err := messageService.RequestBlock(test.input.nodeId, test.input.height)
		assert.Equal(t, err, test.err)
	}

}
