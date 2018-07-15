package api_test

import (
	"testing"

	"errors"

	"log"

	"github.com/it-chain/it-chain-Engine/blockchain"
	"github.com/it-chain/it-chain-Engine/blockchain/api"
	"github.com/it-chain/it-chain-Engine/blockchain/test/mock"
	"github.com/stretchr/testify/assert"
)

func TestBlockApi_Synchronize(t *testing.T) {

	//given
	tests := map[string]struct {
		input struct {
			peerService struct {
				peer blockchain.Peer
				err  error
			}
			SyncWithPeerErr error
			err             error
		}
		err error
	}{
		"success": {
			input: struct {
				peerService struct {
					peer blockchain.Peer
					err  error
				}
				SyncWithPeerErr error
				err             error
			}{
				peerService: struct {
					peer blockchain.Peer
					err  error
				}{
					peer: blockchain.Peer{
						IpAddress: "junksound",
					},
					err: nil,
				},
				SyncWithPeerErr: nil,
			},
			err: nil,
		},
		"fail 1: without Peer": {
			input: struct {
				peerService struct {
					peer blockchain.Peer
					err  error
				}
				SyncWithPeerErr error
				err             error
			}{
				peerService: struct {
					peer blockchain.Peer
					err  error
				}{
					peer: blockchain.Peer{},
					err:  errors.New("error without peer"),
				},
			},
			err: api.ErrGetRandomPeer,
		},

		"fail 2: without Peer": {
			input: struct {
				peerService struct {
					peer blockchain.Peer
					err  error
				}
				SyncWithPeerErr error
				err             error
			}{
				peerService: struct {
					peer blockchain.Peer
					err  error
				}{
					peer: blockchain.Peer{
						IpAddress: "junksound",
					},
					err: nil,
				},
				SyncWithPeerErr: errors.New("error when SyncWithPeer"),
			},
			err: api.ErrSyncWithPeer,
		},
	}

	syncService := &mock.SyncService{}

	peerService := &mock.PeerService{}

	blockQueryService := &mock.BlockQueryService{}

	repo := &mock.EventRepository{}

	bApi, err := api.NewBlockApi(syncService, peerService, blockQueryService, repo, "junksound")

	if err != nil {
		log.Println(err.Error())
	}

	for testName, test := range tests {
		t.Logf("running test case %s", testName)

		peerService.GetRandomPeerFunc = func() (blockchain.Peer, error) {
			return test.input.peerService.peer, test.input.peerService.err
		}

		syncService.SyncWithPeerFunc = func(peer blockchain.Peer) error {
			return test.input.SyncWithPeerErr
		}

		//when
		err = bApi.Synchronize()

		//then
		assert.Equal(t, test.err, err)

	}
}

func TestBlockApi_StageBlock(t *testing.T) {
	tests := map[string]struct {
		input struct {
			block blockchain.Block
		}
	}{
		"success": {
			input: struct {
				block blockchain.Block
			}{block: &blockchain.DefaultBlock{
				Height: uint64(11),
			}},
		},
	}

	syncService := &mock.SyncService{}

	peerService := &mock.PeerService{}

	blockQueryService := &mock.BlockQueryService{}

	repo := &mock.EventRepository{}
	publisherId := "zf"

	blockApi, _ := api.NewBlockApi(syncService, peerService, blockQueryService, repo, publisherId)

	for testName, test := range tests {
		t.Logf("running test case %s", testName)

		blockApi.StageBlock(test.input.block)
	}
}

func TestBlockApi_CommitBlockFromPoolOrSync(t *testing.T) {
	tests := map[string]struct {
		input struct {
			height  blockchain.BlockHeight
			blockId string
		}
		err error
	}{
		"success": {
			input: struct {
				height  blockchain.BlockHeight
				blockId string
			}{height: blockchain.BlockHeight(12), blockId: "zf"},
			err: nil,
		},
	}
	// When

	syncService := &mock.SyncService{}

	peerService := &mock.PeerService{}

	blockQueryService := &mock.BlockQueryService{}

	blockQueryService.GetStagedBlockByIdFunc = func(blockId string) (blockchain.Block, error) {
		assert.Equal(t, "zf", blockId)

		return &blockchain.DefaultBlock{
			Height: blockchain.BlockHeight(12),
		}, nil
	}
	blockQueryService.GetLastCommitedBlockFunc = func() (blockchain.Block, error) {
		return &blockchain.DefaultBlock{
			Height: blockchain.BlockHeight(12),
		}, nil
	}

	repo := &mock.EventRepository{}
	publisherId := "zf"

	// When
	blockApi, _ := api.NewBlockApi(syncService, peerService, blockQueryService, repo, publisherId)

	for testName, test := range tests {
		t.Logf("running test case %s", testName)

		// When
		err := blockApi.CommitBlockFromPoolOrSync(test.input.blockId)

		// Then
		assert.Equal(t, test.err, err)
	}
}

func TestBlockApi_SyncedCheck(t *testing.T) {
	// TODO:
}

func TestBlockApi_SyncIsProgressing(t *testing.T) {
	// when
	syncService := &mock.SyncService{}

	peerService := &mock.PeerService{}

	blockQueryService := &mock.BlockQueryService{}

	repo := &mock.EventRepository{}

	publisherId := "zf"

	// when
	blockApi, _ := api.NewBlockApi(syncService, peerService, blockQueryService, repo, publisherId)

	// then
	state := blockApi.SyncIsProgressing()
	assert.Equal(t, blockchain.DONE, state)
}
