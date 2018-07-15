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

	repo := &mock.EventRepository{}

	bApi, err := api.NewBlockApi(syncService, peerService, repo, "junksound")

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
