package api_test

import (
	"os"
	"sync"
	"testing"

	"github.com/it-chain/engine/blockchain"
	"github.com/it-chain/engine/blockchain/api"
	"github.com/it-chain/engine/blockchain/infra/repo"
	"github.com/it-chain/engine/blockchain/test/mock"
	"github.com/it-chain/engine/common"
	"github.com/it-chain/engine/common/event"
	"github.com/it-chain/engine/common/rabbitmq/pubsub"
	"github.com/stretchr/testify/assert"
)

func TestSyncApi_Synchronize(t *testing.T) {

	//given
	block1 := mock.GetNewBlock([]byte("block1"), 0)
	block2 := mock.GetNewBlock(block1.Seal, 1)
	block3 := mock.GetNewBlock(block2.Seal, 2)

	tests := map[string]struct {
		blockchain []*blockchain.DefaultBlock
		peer       blockchain.Peer
		dbPath     string
	}{
		"not synced": {
			blockchain: []*blockchain.DefaultBlock{
				block1,
			},
			peer: blockchain.Peer{
				IpAddress: "PeerIP",
				PeerId: blockchain.PeerId{
					"PeerId",
				},
			},
			dbPath: "./.db1",
		},

		"synced": {
			blockchain: []*blockchain.DefaultBlock{
				block1,
				block2,
				block3,
			},
			peer: blockchain.Peer{
				IpAddress: "PeerIP",
				PeerId: blockchain.PeerId{
					"PeerId",
				},
			},
			dbPath: "./.db2",
		},

		"No Peer In Network(it is the first peer)": {
			blockchain: []*blockchain.DefaultBlock{},
			peer:       blockchain.Peer{},
			dbPath:     "./.db3",
		},
	}

	PeerBlockchain := []*blockchain.DefaultBlock{
		block1,
		block2,
		block3,
	}

	//set subscriber
	var wg sync.WaitGroup
	wg.Add(2)

	subscriber := pubsub.NewTopicSubscriber("", "Event")
	defer subscriber.Close()

	handler := &mock.CommitEventHandler{}
	handler.HandleFunc = func(event event.BlockCommitted) {
		assert.Equal(t, blockchain.Committed, event.State)
		wg.Done()
	}

	subscriber.SubscribeTopic("block.*", handler)

	publisherID := "junksound"

	eventService := common.NewEventService("", "Event")

	for testName, test := range tests {
		t.Logf("Running test case %s", testName)

		dbPath := test.dbPath

		br, err := repo.NewBlockRepository(dbPath)
		assert.Equal(t, nil, err)

		defer func() {
			br.Close()
			os.RemoveAll(dbPath)
		}()

		queryService := mock.QueryService{}

		queryService.GetRandomPeerFunc = func() (blockchain.Peer, error) {
			return test.peer, nil
		}

		queryService.GetBlockByHeightFromPeerFunc = func(peer blockchain.Peer, height blockchain.BlockHeight) (blockchain.DefaultBlock, error) {
			return *PeerBlockchain[height], nil
		}

		queryService.GetLastBlockFromPeerFunc = func(peer blockchain.Peer) (blockchain.DefaultBlock, error) {
			return *block3, nil
		}

		for _, block := range test.blockchain {
			br.AddBlock(block)
		}

		sApi, err := api.NewSyncApi(publisherID, br, eventService, queryService)

		assert.NoError(t, err)

		//when

		err = sApi.Synchronize()

		assert.NoError(t, err)

		lastBlock, err := br.FindLast()

		assert.NoError(t, err)

		if testName == "No Peer In Network(it is the first peer)" {

			assert.Equal(t, uint64(0), lastBlock.Height)

			continue
		}

		assert.Equal(t, PeerBlockchain[2].Height, lastBlock.Height)

	}

	wg.Wait()

}
