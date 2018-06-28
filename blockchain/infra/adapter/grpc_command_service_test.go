package adapter_test

import (
	"testing"

	"reflect"

	"github.com/it-chain/it-chain-Engine/blockchain/infra/adapter"
	"github.com/it-chain/it-chain-Engine/p2p"
	"github.com/it-chain/yggdrasill/impl"
	"github.com/stretchr/testify/assert"
)

type DefaultBlock = impl.DefaultBlock

func TestGrpcCommandService_RequestBlock(t *testing.T) {

	tests := map[string]struct {
		input struct {
			PeerId p2p.PeerId
			height uint64
		}
		err error
	}{
		"success: request block": {
			input: struct {
				PeerId p2p.PeerId
				height uint64
			}{
				PeerId: p2p.PeerId{
					Id: "1",
				},
				height: uint64(0),
			},
			err: nil,
		},
		"fail: empty node id": {
			input: struct {
				PeerId p2p.PeerId
				height uint64
			}{
				PeerId: p2p.PeerId{},
				height: uint64(0),
			},
			err: adapter.ErrEmptyNodeId,
		},
	}

	publish := func(exchange string, topic string, data interface{}) error {
		assert.Equal(t, exchange, "Command")
		assert.Equal(t, topic, "message.deliver")
		assert.Equal(t, reflect.TypeOf(data).String(), "blockchain.GrpcCommand")
		return nil
	}

	GrpcCommandService := adapter.NewGrpcCommandService(publish)

	for testName, test := range tests {
		t.Logf("running test case %s", testName)
		err := GrpcCommandService.RequestBlock(test.input.PeerId, test.input.height)
		assert.Equal(t, err, test.err)
	}

}

func TestGrpcCommandService_ResponseBlock(t *testing.T) {

	tests := map[string]struct {
		input struct {
			PeerId p2p.PeerId
			block  DefaultBlock
		}
		err error
	}{
		"success: request block": {
			input: struct {
				PeerId p2p.PeerId
				block  DefaultBlock
			}{
				PeerId: p2p.PeerId{
					Id: "1",
				},
				block: DefaultBlock{
					Seal: []byte("seal"),
				},
			},
			err: nil,
		},
		"fail: empty node id": {
			input: struct {
				PeerId p2p.PeerId
				block  DefaultBlock
			}{
				PeerId: p2p.PeerId{},
				block: DefaultBlock{
					Seal: []byte("seal"),
				},
			},
			err: adapter.ErrEmptyNodeId,
		},
		"fail: empty block seal": {
			input: struct {
				PeerId p2p.PeerId
				block  DefaultBlock
			}{
				PeerId: p2p.PeerId{
					"1",
				},
				block: DefaultBlock{
					Seal: nil,
				},
			},
			err: adapter.ErrEmptyBlockSeal,
		},
	}

	publish := func(exchange string, topic string, data interface{}) error {
		assert.Equal(t, exchange, "Command")
		assert.Equal(t, topic, "message.deliver")
		assert.Equal(t, reflect.TypeOf(data).String(), "blockchain.GrpcCommand")

		return nil
	}

	GrpcCommandService := adapter.NewGrpcCommandService(publish)

	for testName, test := range tests {
		t.Logf("running test case %s", testName)
		err := GrpcCommandService.ResponseBlock(test.input.PeerId, &test.input.block)
		assert.Equal(t, err, test.err)
	}

}
