/*
 * Copyright 2018 It-chain
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * https://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package api_test

import (
	"os"
	"sync"
	"testing"

	"github.com/it-chain/engine/blockchain"
	"github.com/it-chain/engine/blockchain/api"
	"github.com/it-chain/engine/blockchain/infra/mem"
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
				ApiGatewayAddress: "PeerIP",
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
				ApiGatewayAddress: "PeerIP",
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

		queryService.GetBlockByHeightFromPeerFunc = func(height blockchain.BlockHeight, peer blockchain.Peer) (blockchain.DefaultBlock, error) {
			return *PeerBlockchain[height], nil
		}

		queryService.GetLastBlockFromPeerFunc = func(peer blockchain.Peer) (blockchain.DefaultBlock, error) {
			return *block3, nil
		}

		for _, block := range test.blockchain {
			br.AddBlock(block)
		}

		blockPool := mem.NewBlockPool()

		sApi, err := api.NewSyncApi(publisherID, br, eventService, queryService, blockPool)

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

func TestSyncApi_CommitStagedBlocks_Add_Blocks_To_Repository(t *testing.T) {
	//given
	block1 := mock.GetNewBlock([]byte("genesis"), 0)
	block2 := mock.GetNewBlock(block1.Seal, 1)
	block3 := mock.GetNewBlock(block2.Seal, 2)
	block4 := mock.GetNewBlock(block3.Seal, 3)
	block5 := mock.GetNewBlock(block4.Seal, 4)

	//given
	publisherId := "zf"

	dbPath := "./.test"
	blockRepository, err := repo.NewBlockRepository(dbPath)
	assert.NoError(t, err)

	defer func() {
		blockRepository.Close()
		os.RemoveAll(dbPath)
	}()

	blockRepository.Save(*block1)
	blockRepository.Save(*block2)
	blockRepository.Save(*block3)

	eventService := common.NewEventService("", "Event")
	queryService := mock.QueryService{}

	//given
	blockPool := mem.NewBlockPool()
	blockPool.Add(*block4)
	blockPool.Add(*block5)

	syncApi, err := api.NewSyncApi(publisherId, blockRepository, eventService, queryService, blockPool)
	assert.NoError(t, err)

	// when
	syncApi.CommitStagedBlocks()

	// then
	block, err := blockRepository.FindLast()
	assert.NoError(t, err)
	assert.Equal(t, block.GetHeight(), (*block5).GetHeight())
}

func TestSyncApi_CommitStagedBlocks_Drop_Blocks_From_BlockPool(t *testing.T) {
	//given
	block1 := mock.GetNewBlock([]byte("genesis"), 0)
	block2 := mock.GetNewBlock(block1.Seal, 1)
	block3 := mock.GetNewBlock(block2.Seal, 2)

	//given
	publisherId := "zf"

	dbPath := "./.test"
	blockRepository, err := repo.NewBlockRepository(dbPath)
	assert.NoError(t, err)

	defer func() {
		blockRepository.Close()
		os.RemoveAll(dbPath)
	}()

	blockRepository.Save(*block1)
	blockRepository.Save(*block2)
	blockRepository.Save(*block3)

	eventService := common.NewEventService("", "Event")
	queryService := mock.QueryService{}

	//given
	blockPool := mem.NewBlockPool()
	blockPool.Add(*block2)
	blockPool.Add(*block3)

	syncApi, err := api.NewSyncApi(publisherId, blockRepository, eventService, queryService, blockPool)
	assert.NoError(t, err)

	// when
	syncApi.CommitStagedBlocks()

	// then
	block, err := blockRepository.FindLast()
	assert.NoError(t, err)
	assert.Equal(t, block.GetHeight(), (*block3).GetHeight())

	// then
	assert.Equal(t, blockPool.Size(), 0)
}

func TestSyncApi_CommitStagedBlocks_When_BlockPool_Has_Higher_Height_Block(t *testing.T) {
	//given
	block1 := mock.GetNewBlock([]byte("genesis"), 0)
	block2 := mock.GetNewBlock(block1.Seal, 1)
	block3 := mock.GetNewBlock(block2.Seal, 2)
	block4 := mock.GetNewBlock(block3.Seal, 3)
	block5 := mock.GetNewBlock(block4.Seal, 4)
	block6 := mock.GetNewBlock(block5.Seal, 5)
	block7 := mock.GetNewBlock(block6.Seal, 6)

	//given
	publisherId := "zf"

	dbPath := "./.test"
	blockRepository, err := repo.NewBlockRepository(dbPath)
	assert.NoError(t, err)

	defer func() {
		blockRepository.Close()
		os.RemoveAll(dbPath)
	}()

	blockRepository.Save(*block1)
	blockRepository.Save(*block2)
	blockRepository.Save(*block3)

	eventService := common.NewEventService("", "Event")
	queryService := mock.QueryService{}

	//given
	blockPool := mem.NewBlockPool()
	blockPool.Add(*block7)

	syncApi, err := api.NewSyncApi(publisherId, blockRepository, eventService, queryService, blockPool)
	assert.NoError(t, err)

	// when
	syncApi.CommitStagedBlocks()

	// then
	block, err := blockRepository.FindLast()
	assert.NoError(t, err)
	assert.Equal(t, block.GetHeight(), (*block3).GetHeight())

	// then
	assert.Equal(t, blockPool.GetSortedKeys(), []uint64{uint64(6)})
	assert.Equal(t, blockPool.Size(), 1)
}
