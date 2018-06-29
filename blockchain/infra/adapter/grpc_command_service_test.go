package adapter_test

import (
	"testing"

	"reflect"

	"time"

	"github.com/it-chain/it-chain-Engine/blockchain/infra/adapter"
	"github.com/it-chain/it-chain-Engine/p2p"
	"github.com/it-chain/yggdrasill/common"
	"github.com/stretchr/testify/assert"
)

type MockBlock struct {
	seal []byte
}

func (MockBlock) SetSeal(seal []byte) {
	panic("implement me")
}
func (MockBlock) SetPrevSeal(prevseal []byte) {
	panic("implement me")
}
func (MockBlock) SetHeight(height uint64) {
	panic("implement me")
}
func (MockBlock) PutTx(tx common.Transaction) error {
	panic("implement me")
}
func (MockBlock) SetTxSeal(txSeal [][]byte) {
	panic("implement me")
}
func (MockBlock) SetCreator(creator []byte) {
	panic("implement me")
}
func (MockBlock) SetTimestamp(currentTime time.Time) {
	panic("implement me")
}
func (m MockBlock) GetSeal() []byte {
	return m.seal
}
func (MockBlock) GetPrevSeal() []byte {
	panic("implement me")
}
func (MockBlock) GetHeight() uint64 {
	panic("implement me")
}
func (MockBlock) GetTxList() []common.Transaction {
	panic("implement me")
}
func (MockBlock) GetTxSeal() [][]byte {
	panic("implement me")
}
func (MockBlock) GetCreator() []byte {
	panic("implement me")
}
func (MockBlock) GetTimestamp() time.Time {
	panic("implement me")
}
func (MockBlock) Serialize() ([]byte, error) {
	panic("implement me")
}
func (MockBlock) Deserialize(serializedBlock []byte) error {
	panic("implement me")
}
func (MockBlock) IsReadyToPublish() bool {
	panic("implement me")
}
func (MockBlock) IsPrev(serializedPrevBlock []byte) bool {
	panic("implement me")
}

func TestGrpcCommandService_RequestBlock(t *testing.T) {

	tests := map[string]struct {
		input struct {
			peerId p2p.PeerId
			height uint64
		}
		err error
	}{
		"success: request block": {
			input: struct {
				peerId p2p.PeerId
				height uint64
			}{
				peerId: p2p.PeerId{
					Id: "1",
				},
				height: uint64(0),
			},
			err: nil,
		},
		"fail: empty node id": {
			input: struct {
				peerId p2p.PeerId
				height uint64
			}{
				peerId: p2p.PeerId{},
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
		err := GrpcCommandService.RequestBlock(test.input.peerId, test.input.height)
		assert.Equal(t, err, test.err)
	}

}

func TestGrpcCommandService_ResponseBlock(t *testing.T) {

	tests := map[string]struct {
		input struct {
			peerId p2p.PeerId
			block  MockBlock
		}
		err error
	}{
		"success: request block": {
			input: struct {
				peerId p2p.PeerId
				block  MockBlock
			}{
				peerId: p2p.PeerId{
					Id: "1",
				},
				block: MockBlock{
					seal: []byte("seal"),
				},
			},
			err: nil,
		},
		"fail: empty node id": {
			input: struct {
				peerId p2p.PeerId
				block  MockBlock
			}{
				peerId: p2p.PeerId{},
				block: MockBlock{
					seal: []byte("seal"),
				},
			},
			err: adapter.ErrEmptyNodeId,
		},
		"fail: empty block seal": {
			input: struct {
				peerId p2p.PeerId
				block  MockBlock
			}{
				peerId: p2p.PeerId{
					"1",
				},
				block: MockBlock{
					seal: nil,
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
		err := GrpcCommandService.ResponseBlock(test.input.peerId, test.input.block)
		assert.Equal(t, err, test.err)
	}

}
