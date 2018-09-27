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
	"testing"

	"sync"

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

var block1 = mock.GetNewBlock([]byte("block1"), 0)
var block2 = mock.GetNewBlock(block1.Seal, 1)
var block3 = mock.GetNewBlock(block2.Seal, 2)
var block4 = mock.GetNewBlock(block3.Seal, 3)
var block5 = mock.GetNewBlock(block4.Seal, 4)
var block6 = mock.GetNewBlock(block5.Seal, 5)
var block7 = mock.GetNewBlock(block6.Seal, 6)

var peerForSync = blockchain.Peer{
	PeerID:            "PeerID",
	ApiGatewayAddress: "PeerIP",
}

var peerBlockchain = []*blockchain.DefaultBlock{
	block1,
	block2,
	block3,
}

func getQueryService(targetPeer blockchain.Peer) mock.QueryService {
	queryService := mock.QueryService{}

	queryService.GetRandomPeerFunc = func() (blockchain.Peer, error) {
		return targetPeer, nil
	}
	queryService.GetBlockByHeightFromPeerFunc = func(height blockchain.BlockHeight, peer blockchain.Peer) (blockchain.DefaultBlock, error) {
		return *peerBlockchain[height], nil
	}
	queryService.GetLastBlockFromPeerFunc = func(peer blockchain.Peer) (blockchain.DefaultBlock, error) {
		return *block3, nil
	}

	return queryService
}

func TestSyncApi_Synchronize_NotSynced_BlockPool_Has_Shorter_Heights(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(2)

	subscriber := pubsub.NewTopicSubscriber("", "Event")
	defer subscriber.Close()

	commitEventhandler := &mock.CommitEventHandler{}
	commitEventhandler.HandleFunc = func(event event.BlockCommitted) {
		assert.Equal(t, blockchain.Committed, event.State)
		wg.Done()
	}

	subscriber.SubscribeTopic("block.*", commitEventhandler)

	publisherID := "junksound"

	dbPath := "./.db"
	br, err := repo.NewBlockRepository(dbPath)
	assert.Equal(t, nil, err)

	defer func() {
		br.Close()
		os.RemoveAll(dbPath)
	}()

	ssr := mem.NewSyncStateRepository()

	eventService := common.NewEventService("", "Event")

	queryService := getQueryService(peerForSync)

	blockPool := mem.NewBlockPool()

	// when
	br.AddBlock(block1)
	blockPool.Add(*block1)
	blockPool.Add(*block2)

	sApi, err := api.NewSyncApi(publisherID, br, ssr, eventService, queryService, blockPool)
	assert.NoError(t, err)

	//when
	err = sApi.Synchronize()
	assert.NoError(t, err)

	lastBlock, err := br.FindLast()
	assert.NoError(t, err)

	syncState := ssr.Get()
	assert.Equal(t, false, syncState.SyncProgressing)

	assert.Equal(t, uint64(2), lastBlock.Height)
	assert.Equal(t, 0, blockPool.Size())
	wg.Wait()
}

func TestSyncApi_Synchronize_NotSynced_BlockPool_Has_Target_Heights_ThreeBlocks(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(5)

	subscriber := pubsub.NewTopicSubscriber("", "Event")
	defer subscriber.Close()

	commitEventhandler := &mock.CommitEventHandler{}
	commitEventhandler.HandleFunc = func(event event.BlockCommitted) {
		assert.Equal(t, blockchain.Committed, event.State)
		wg.Done()
	}

	subscriber.SubscribeTopic("block.*", commitEventhandler)

	publisherID := "junksound"

	dbPath := "./.db"
	br, err := repo.NewBlockRepository(dbPath)
	assert.Equal(t, nil, err)

	defer func() {
		br.Close()
		os.RemoveAll(dbPath)
	}()

	ssr := mem.NewSyncStateRepository()

	eventService := common.NewEventService("", "Event")

	queryService := getQueryService(peerForSync)

	blockPool := mem.NewBlockPool()

	// when
	br.AddBlock(block1)
	blockPool.Add(*block4)
	blockPool.Add(*block5)
	blockPool.Add(*block6)

	sApi, err := api.NewSyncApi(publisherID, br, ssr, eventService, queryService, blockPool)
	assert.NoError(t, err)

	//when
	err = sApi.Synchronize()
	assert.NoError(t, err)

	lastBlock, err := br.FindLast()
	assert.NoError(t, err)

	syncState := ssr.Get()
	assert.Equal(t, false, syncState.SyncProgressing)

	assert.Equal(t, uint64(5), lastBlock.Height)
	assert.Equal(t, 0, blockPool.Size())
	wg.Wait()
}

func TestSyncApi_Synchronize_NotSynced_BlockPool_Has_Target_Heights_TwoBlocks_OneLeftInPool(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(4)

	subscriber := pubsub.NewTopicSubscriber("", "Event")
	defer subscriber.Close()

	commitEventhandler := &mock.CommitEventHandler{}
	commitEventhandler.HandleFunc = func(event event.BlockCommitted) {
		assert.Equal(t, blockchain.Committed, event.State)
		wg.Done()
	}

	subscriber.SubscribeTopic("block.*", commitEventhandler)

	publisherID := "junksound"

	dbPath := "./.db"
	br, err := repo.NewBlockRepository(dbPath)
	assert.Equal(t, nil, err)

	defer func() {
		br.Close()
		os.RemoveAll(dbPath)
	}()

	ssr := mem.NewSyncStateRepository()

	eventService := common.NewEventService("", "Event")

	queryService := getQueryService(peerForSync)

	blockPool := mem.NewBlockPool()

	// when
	br.AddBlock(block1)
	blockPool.Add(*block4)
	blockPool.Add(*block5)
	blockPool.Add(*block7)

	sApi, err := api.NewSyncApi(publisherID, br, ssr, eventService, queryService, blockPool)
	assert.NoError(t, err)

	//when
	err = sApi.Synchronize()
	assert.NoError(t, err)

	lastBlock, err := br.FindLast()
	assert.NoError(t, err)

	syncState := ssr.Get()
	assert.Equal(t, false, syncState.SyncProgressing)

	assert.Equal(t, uint64(4), lastBlock.Height)
	assert.Equal(t, 1, blockPool.Size())
	wg.Wait()
}

func TestSyncApi_Synchronize_NotSynced_BlockPool_Has_Higher_Heights(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(2)

	subscriber := pubsub.NewTopicSubscriber("", "Event")
	defer subscriber.Close()

	commitEventhandler := &mock.CommitEventHandler{}
	commitEventhandler.HandleFunc = func(event event.BlockCommitted) {
		assert.Equal(t, blockchain.Committed, event.State)
		wg.Done()
	}

	subscriber.SubscribeTopic("block.*", commitEventhandler)

	publisherID := "junksound"

	dbPath := "./.db"
	br, err := repo.NewBlockRepository(dbPath)
	assert.Equal(t, nil, err)

	defer func() {
		br.Close()
		os.RemoveAll(dbPath)
	}()

	ssr := mem.NewSyncStateRepository()

	eventService := common.NewEventService("", "Event")

	queryService := getQueryService(peerForSync)

	blockPool := mem.NewBlockPool()

	// when
	br.AddBlock(block1)
	blockPool.Add(*block6)
	blockPool.Add(*block7)

	sApi, err := api.NewSyncApi(publisherID, br, ssr, eventService, queryService, blockPool)
	assert.NoError(t, err)

	//when
	err = sApi.Synchronize()
	assert.NoError(t, err)

	lastBlock, err := br.FindLast()
	assert.NoError(t, err)

	syncState := ssr.Get()
	assert.Equal(t, false, syncState.SyncProgressing)

	assert.Equal(t, uint64(2), lastBlock.Height)
	assert.Equal(t, 2, blockPool.Size())
	wg.Wait()
}

func TestSyncApi_Synchronize_NotSynced_BlockPool_Has_Shorter_Target_Higher_Heights(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(3)

	subscriber := pubsub.NewTopicSubscriber("", "Event")
	defer subscriber.Close()

	commitEventhandler := &mock.CommitEventHandler{}
	commitEventhandler.HandleFunc = func(event event.BlockCommitted) {
		assert.Equal(t, blockchain.Committed, event.State)
		wg.Done()
	}

	subscriber.SubscribeTopic("block.*", commitEventhandler)

	publisherID := "junksound"

	dbPath := "./.db"
	br, err := repo.NewBlockRepository(dbPath)
	assert.Equal(t, nil, err)

	defer func() {
		br.Close()
		os.RemoveAll(dbPath)
	}()

	ssr := mem.NewSyncStateRepository()

	eventService := common.NewEventService("", "Event")

	queryService := getQueryService(peerForSync)

	blockPool := mem.NewBlockPool()

	// when
	br.AddBlock(block1)
	blockPool.Add(*block1)
	blockPool.Add(*block2)
	blockPool.Add(*block4)
	blockPool.Add(*block6)
	blockPool.Add(*block7)

	sApi, err := api.NewSyncApi(publisherID, br, ssr, eventService, queryService, blockPool)
	assert.NoError(t, err)

	//when
	err = sApi.Synchronize()
	assert.NoError(t, err)

	lastBlock, err := br.FindLast()
	assert.NoError(t, err)

	syncState := ssr.Get()
	assert.Equal(t, false, syncState.SyncProgressing)

	assert.Equal(t, uint64(3), lastBlock.Height)
	assert.Equal(t, 2, blockPool.Size())
	wg.Wait()
}

func TestSyncApi_Synchronize_Synced_BlockPool_Has_NoTarget_Heights(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(0)

	subscriber := pubsub.NewTopicSubscriber("", "Event")
	defer subscriber.Close()

	commitEventhandler := &mock.CommitEventHandler{}
	commitEventhandler.HandleFunc = func(event event.BlockCommitted) {
		assert.Equal(t, blockchain.Committed, event.State)
		wg.Done()
	}

	subscriber.SubscribeTopic("block.*", commitEventhandler)

	publisherID := "junksound"

	dbPath := "./.db"
	br, err := repo.NewBlockRepository(dbPath)
	assert.Equal(t, nil, err)

	defer func() {
		br.Close()
		os.RemoveAll(dbPath)
	}()

	ssr := mem.NewSyncStateRepository()

	eventService := common.NewEventService("", "Event")

	queryService := getQueryService(peerForSync)

	blockPool := mem.NewBlockPool()

	// when
	br.AddBlock(block1)
	br.AddBlock(block2)
	br.AddBlock(block3)
	blockPool.Add(*block1)
	blockPool.Add(*block2)

	sApi, err := api.NewSyncApi(publisherID, br, ssr, eventService, queryService, blockPool)
	assert.NoError(t, err)

	//when
	err = sApi.Synchronize()
	assert.NoError(t, err)

	lastBlock, err := br.FindLast()
	assert.NoError(t, err)

	syncState := ssr.Get()
	assert.Equal(t, false, syncState.SyncProgressing)

	assert.Equal(t, uint64(2), lastBlock.Height)
	assert.Equal(t, 2, blockPool.Size())
	wg.Wait()
}

func TestSyncApi_Synchronize_Synced_BlockPool_Has_Target_Heights(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(0)

	subscriber := pubsub.NewTopicSubscriber("", "Event")
	defer subscriber.Close()

	commitEventhandler := &mock.CommitEventHandler{}
	commitEventhandler.HandleFunc = func(event event.BlockCommitted) {
		assert.Equal(t, blockchain.Committed, event.State)
		wg.Done()
	}

	subscriber.SubscribeTopic("block.*", commitEventhandler)

	publisherID := "junksound"

	dbPath := "./.db"
	br, err := repo.NewBlockRepository(dbPath)
	assert.Equal(t, nil, err)

	defer func() {
		br.Close()
		os.RemoveAll(dbPath)
	}()

	ssr := mem.NewSyncStateRepository()

	eventService := common.NewEventService("", "Event")

	queryService := getQueryService(peerForSync)

	blockPool := mem.NewBlockPool()

	// when
	br.AddBlock(block1)
	br.AddBlock(block2)
	br.AddBlock(block3)
	blockPool.Add(*block4)
	blockPool.Add(*block5)

	sApi, err := api.NewSyncApi(publisherID, br, ssr, eventService, queryService, blockPool)
	assert.NoError(t, err)

	//when
	err = sApi.Synchronize()
	assert.NoError(t, err)

	lastBlock, err := br.FindLast()
	assert.NoError(t, err)

	syncState := ssr.Get()
	assert.Equal(t, false, syncState.SyncProgressing)

	assert.Equal(t, uint64(2), lastBlock.Height)
	assert.Equal(t, 2, blockPool.Size())
	wg.Wait()
}

func TestSyncApi_Synchronize_When_NoPeer_In_Network(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(0)

	subscriber := pubsub.NewTopicSubscriber("", "Event")
	defer subscriber.Close()

	commitEventhandler := &mock.CommitEventHandler{}
	commitEventhandler.HandleFunc = func(event event.BlockCommitted) {
		assert.Equal(t, blockchain.Committed, event.State)
		wg.Done()
	}

	subscriber.SubscribeTopic("block.*", commitEventhandler)

	publisherID := "junksound"

	dbPath := "./.db"
	br, err := repo.NewBlockRepository(dbPath)
	assert.Equal(t, nil, err)

	defer func() {
		br.Close()
		os.RemoveAll(dbPath)
	}()

	ssr := mem.NewSyncStateRepository()

	eventService := common.NewEventService("", "Event")

	queryService := getQueryService(blockchain.Peer{})

	blockPool := mem.NewBlockPool()

	// when
	br.AddBlock(block1)
	br.AddBlock(block2)
	br.AddBlock(block3)
	blockPool.Add(*block4)
	blockPool.Add(*block5)

	sApi, err := api.NewSyncApi(publisherID, br, ssr, eventService, queryService, blockPool)
	assert.NoError(t, err)

	//when
	err = sApi.Synchronize()
	assert.NoError(t, err)

	lastBlock, err := br.FindLast()
	assert.NoError(t, err)

	syncState := ssr.Get()
	assert.Equal(t, false, syncState.SyncProgressing)

	assert.Equal(t, uint64(2), lastBlock.Height)
	assert.Equal(t, 2, blockPool.Size())
	wg.Wait()
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

	ssr := mem.NewSyncStateRepository()

	syncApi, err := api.NewSyncApi(publisherId, blockRepository, ssr, eventService, queryService, blockPool)
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

	ssr := mem.NewSyncStateRepository()

	syncApi, err := api.NewSyncApi(publisherId, blockRepository, ssr, eventService, queryService, blockPool)
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
