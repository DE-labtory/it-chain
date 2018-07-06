package adapter_test

import (
	"testing"

	"github.com/it-chain/it-chain-Engine/blockchain"
	"github.com/it-chain/it-chain-Engine/blockchain/infra/adapter"

	"errors"

	"github.com/it-chain/midgard"
	"github.com/stretchr/testify/assert"
)

type MockPeerRepository struct {
	AddFunc    func(peer blockchain.Peer) error
	RemoveFunc func(peer blockchain.PeerId) error
}

func (nr MockPeerRepository) Add(peer blockchain.Peer) error {
	return nr.AddFunc(peer)
}
func (nr MockPeerRepository) Remove(peerId blockchain.PeerId) error {
	return nr.RemoveFunc(peerId)
}

type MockEventListenerBlockApi struct {
}

func (api MockEventListenerBlockApi) AddBlockToPool(block blockchain.Block) {
	return
}

func (api MockEventListenerBlockApi) CheckAndSaveBlockFromPool(height blockchain.Block) error {
	return nil
}

func TestRepositoryProjector_HandleNodeCreatedEvent(t *testing.T) {
	tests := map[string]struct {
		input struct {
			ID      string
			peerId  string
			address string
			rpErr   error
		}
		err error
	}{
		"success": {
			input: struct {
				ID      string
				peerId  string
				address string
				rpErr   error
			}{ID: string("zf"), peerId: string("zf2"), address: string("11.22.33.44"), rpErr: nil},
			err: nil,
		},
		"empty eventId test": {
			input: struct {
				ID      string
				peerId  string
				address string
				rpErr   error
			}{ID: string(""), peerId: string("zf2"), address: string("11.22.33.44"), rpErr: nil},
			err: adapter.ErrEmptyEventId,
		},
		"repository error test": {
			input: struct {
				ID      string
				peerId  string
				address string
				rpErr   error
			}{ID: string("zf"), peerId: string("zf2"), address: string("11.22.33.44"), rpErr: errors.New("repository error")},
			err: adapter.ErrNodeApi,
		},
	}

	for testName, test := range tests {
		t.Logf("running test case %s", testName)

		// When
		nr := MockPeerRepository{}
		nr.AddFunc = func(peer blockchain.Peer) error {
			assert.Equal(t, peer.PeerId.Id, string("zf2"))
			assert.Equal(t, peer.IpAddress, string("11.22.33.44"))
			return test.input.rpErr
		}

		repositoryProjector := adapter.RepositoryProjector{PeerRepository: nr}

		event := blockchain.NodeCreatedEvent{
			EventModel: midgard.EventModel{
				ID: test.input.ID,
			},
			Peer: blockchain.Peer{
				PeerId: blockchain.PeerId{
					test.input.peerId,
				},
				IpAddress: test.input.address,
			},
		}

		// When
		err := repositoryProjector.HandleNodeCreatedEvent(event)

		// Then
		assert.Equal(t, err, test.err)
	}
}

func TestEventHandler_HandleNodeDeletedEvent(t *testing.T) {
	tests := map[string]struct {
		input struct {
			ID     string
			peerId string
			rpErr  error
		}
		err error
	}{
		"success": {
			input: struct {
				ID     string
				peerId string
				rpErr  error
			}{ID: string("zf"), peerId: string("zf2"), rpErr: nil},
			err: nil,
		},
		"empty eventId test": {
			input: struct {
				ID     string
				peerId string
				rpErr  error
			}{ID: string(""), peerId: string("zf2"), rpErr: nil},
			err: adapter.ErrEmptyEventId,
		},
		"repository error test": {
			input: struct {
				ID     string
				peerId string
				rpErr  error
			}{ID: string("zf"), peerId: string("zf2"), rpErr: errors.New("repository error")},
			err: adapter.ErrNodeApi,
		},
	}

	for testName, test := range tests {
		t.Logf("running test case %s", testName)

		// When
		nr := MockPeerRepository{}
		nr.RemoveFunc = func(peerId blockchain.PeerId) error {
			assert.Equal(t, peerId.Id, string("zf2"))
			return test.input.rpErr
		}

		repositoryProjector := adapter.RepositoryProjector{PeerRepository: nr}

		event := blockchain.NodeDeletedEvent{
			EventModel: midgard.EventModel{
				ID: test.input.ID,
			},
			Peer: blockchain.Peer{
				PeerId: blockchain.PeerId{
					test.input.peerId,
				},
			},
		}

		// When
		err := repositoryProjector.HandleNodeDeletedEvent(event)

		// Then
		assert.Equal(t, err, test.err)
	}
}

type MockEventRepository struct{}

func (er MockEventRepository) Load(aggregate midgard.Aggregate, aggregateID string) error { return nil }
func (er MockEventRepository) Save(aggregateID string, events ...midgard.Event) error     { return nil }
func (er MockEventRepository) Close()                                                     {}

//todo eventstore를 활용한 testcase재 작성필요
//func TestEventHandler_HandleBlockAddToPoolEvent(t *testing.T) {
//	tests := map[string]struct {
//		input struct {
//			blockchain.BlockAddToPoolEvent
//		}
//		err error
//	}{
//		"success": {
//			input: struct {
//				blockchain.BlockAddToPoolEvent
//			}{BlockAddToPoolEvent: blockchain.BlockAddToPoolEvent{
//				Block: &blockchain.DefaultBlock{
//					Height: uint64(12),
//				},
//			}},
//			err: nil,
//		},
//		"block nil test": {
//			input: struct {
//				blockchain.BlockAddToPoolEvent
//			}{BlockAddToPoolEvent: blockchain.BlockAddToPoolEvent{
//				Block: nil,
//			}},
//			err: adapter.ErrBlockNil,
//		},
//	}
//
//	// When
//	nr := MockPeerRepository{}
//	er := MockEventRepository{}
//	rp := adapter.RepositoryProjector{
//		PeerRepository:  nr,
//		EventRepository: er,
//	}
//
//	blockApi := MockEventListenerBlockApi{}
//
//	eventHandler := adapter.NewEventHandler(rp, blockApi)
//
//	for testName, test := range tests {
//		t.Logf("running test case %s", testName)
//
//		// When
//		err := eventHandler.HandleBlockAddToPoolEvent(test.input.BlockAddToPoolEvent)
//
//		// Then
//		assert.Equal(t, err, test.err)
//	}
//}
