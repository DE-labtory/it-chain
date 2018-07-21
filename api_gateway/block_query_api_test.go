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

package api_gateway_test

import (
	"os"
	"testing"
	"time"

	"github.com/it-chain/engine/api_gateway"
	"github.com/it-chain/engine/blockchain"
	"github.com/it-chain/yggdrasill/common"
	"github.com/stretchr/testify/assert"
)

func TestBlockPoolRepositoryImpl_AddCreatedBlock(t *testing.T) {
	bpr := api_gateway.NewBlockPoolRepository()

	// when
	block1 := &blockchain.DefaultBlock{
		Seal:   []byte{0x1},
		Height: blockchain.BlockHeight(1),
		State:  blockchain.Created,
	}
	// when
	err := bpr.AddCreatedBlock(*block1)

	// then
	assert.Equal(t, nil, err)
	assert.Equal(t, 1, len(bpr.Blocks))
	assert.Equal(t, []byte{0x1}, bpr.Blocks[0].GetSeal())
}

func TestBlockPoolRepositoryImpl_AddCreatedBlock_InvalidStateBlock(t *testing.T) {
	bpr := api_gateway.NewBlockPoolRepository()

	// when
	block1 := &blockchain.DefaultBlock{
		Seal:   []byte{0x1},
		Height: blockchain.BlockHeight(1),
		State:  blockchain.Committed,
	}
	// when
	err := bpr.AddCreatedBlock(*block1)

	// then
	assert.Equal(t, api_gateway.ErrInvalidStateBlock, err)
	assert.Equal(t, 0, len(bpr.Blocks))
}

func TestBlockPoolRepositoryImpl_GetStagedBlockByHeight(t *testing.T) {
	bpr := api_gateway.NewBlockPoolRepository()

	// when
	bpr.Blocks = append(bpr.Blocks, &blockchain.DefaultBlock{
		Seal:     []byte{0x1},
		PrevSeal: []byte{0x1},
		Height:   uint64(1),
	})

	// when
	block, err := bpr.GetStagedBlockByHeight(1)

	// then
	assert.Equal(t, err, nil)
	assert.Equal(t, uint64(1), block.GetHeight())
	assert.Equal(t, []byte{0x1}, block.GetSeal())
	assert.Equal(t, []byte{0x1}, block.GetPrevSeal())

	// when
	_, err2 := bpr.GetStagedBlockByHeight(133)

	// then
	assert.Equal(t, err2, api_gateway.ErrNoStagedBlock)
}

func TestBlockPoolRepositoryImpl_GetStagedBlockById(t *testing.T) {
	bpr := api_gateway.NewBlockPoolRepository()

	// when
	bpr.Blocks = append(bpr.Blocks, &blockchain.DefaultBlock{
		Seal:     []byte{0x1},
		PrevSeal: []byte{0x1},
		Height:   uint64(1),
	})

	// when
	block, err := bpr.GetStagedBlockById(string([]byte{0x1}))

	// then
	assert.Equal(t, err, nil)
	assert.Equal(t, uint64(1), block.GetHeight())
	assert.Equal(t, []byte{0x1}, block.GetSeal())
	assert.Equal(t, []byte{0x1}, block.GetPrevSeal())

	// when
	_, err2 := bpr.GetStagedBlockById(string([]byte{0x2}))

	// then
	assert.Equal(t, err2, api_gateway.ErrNoStagedBlock)
}

func TestBlockPoolRepositoryImpl_GetFirstStagedBlock_basic(t *testing.T) {
	bpr := api_gateway.NewBlockPoolRepository()

	// when
	bpr.Blocks = append(bpr.Blocks, &blockchain.DefaultBlock{
		Seal:     []byte{0x1},
		PrevSeal: []byte{0x1},
		Height:   uint64(1),
		State:    blockchain.Created,
	})

	bpr.Blocks = append(bpr.Blocks, &blockchain.DefaultBlock{
		Seal:     []byte{0x2},
		PrevSeal: []byte{0x2},
		Height:   uint64(2),
		State:    blockchain.Staged,
	})

	// when
	block, err := bpr.GetFirstStagedBlock()

	assert.Equal(t, nil, err)
	assert.Equal(t, uint64(2), block.Height)
	assert.Equal(t, []byte{0x2}, block.Seal)
}

func TestBlockPoolRepositoryImpl_GetFirstStagedBlock_NoStagedBlockFound(t *testing.T) {
	bpr := api_gateway.NewBlockPoolRepository()

	// when
	bpr.Blocks = append(bpr.Blocks, &blockchain.DefaultBlock{
		Seal:     []byte{0x1},
		PrevSeal: []byte{0x1},
		Height:   uint64(1),
		State:    blockchain.Created,
	})

	bpr.Blocks = append(bpr.Blocks, &blockchain.DefaultBlock{
		Seal:     []byte{0x2},
		PrevSeal: []byte{0x2},
		Height:   uint64(2),
		State:    blockchain.Created,
	})

	// when
	block, err := bpr.GetFirstStagedBlock()

	assert.Equal(t, api_gateway.ErrNoStagedBlock, err)
	assert.Equal(t, true, block.IsEmpty())
}

func TestBlockPoolRepositoryImpl_GetFirstStagedBlock_lenIsZero(t *testing.T) {
	bpr := api_gateway.NewBlockPoolRepository()

	// when
	block, err := bpr.GetFirstStagedBlock()

	assert.Equal(t, api_gateway.ErrNoStagedBlock, err)
	assert.Equal(t, uint64(0), block.Height)
	assert.Equal(t, []byte(nil), block.Seal)
}

func TestBlockPoolRepositoryImpl_RemoveById(t *testing.T) {
	bpr := api_gateway.NewBlockPoolRepository()

	// when
	bpr.Blocks = append(bpr.Blocks, &blockchain.DefaultBlock{
		Seal:     []byte{0x1},
		PrevSeal: []byte{0x1},
		Height:   uint64(1),
		State:    blockchain.Created,
	})

	bpr.Blocks = append(bpr.Blocks, &blockchain.DefaultBlock{
		Seal:     []byte{0x2},
		PrevSeal: []byte{0x2},
		Height:   uint64(2),
		State:    blockchain.Created,
	})

	// when
	err := bpr.RemoveById(string(0x1))

	// then
	assert.Equal(t, nil, err)
	assert.Equal(t, 1, len(bpr.Blocks))
	assert.Equal(t, []byte{0x2}, bpr.Blocks[0].GetSeal())

	// when
	err2 := bpr.RemoveById(string(0x3))

	// then
	assert.Equal(t, 1, len(bpr.Blocks))
	assert.Equal(t, api_gateway.ErrFailRemoveBlock, err2)
}

func TestBlockPoolRepositoryImpl_RemoveById_FailRemoving(t *testing.T) {
	bpr := api_gateway.NewBlockPoolRepository()

	// when
	bpr.Blocks = append(bpr.Blocks, &blockchain.DefaultBlock{
		Seal:     []byte{0x1},
		PrevSeal: []byte{0x1},
		Height:   uint64(1),
		State:    blockchain.Created,
	})

	bpr.Blocks = append(bpr.Blocks, &blockchain.DefaultBlock{
		Seal:     []byte{0x2},
		PrevSeal: []byte{0x2},
		Height:   uint64(2),
		State:    blockchain.Created,
	})

	// when
	err := bpr.RemoveById(string(0x3))

	// then
	assert.Equal(t, 2, len(bpr.Blocks))
	assert.Equal(t, api_gateway.ErrFailRemoveBlock, err)
}

func TestBlockQueryApi_GetLastCommitedBlock(t *testing.T) {
	dbPath := "./.db"

	// when
	cbr, err := api_gateway.NewCommitedBlockRepositoryImpl(dbPath)
	// then
	assert.Equal(t, nil, err)

	defer func() {
		cbr.Close()
		os.RemoveAll(dbPath)
	}()

	// when
	block1 := getNewBlock([]byte("genesis"), 0)
	err = cbr.AddBlock(block1)
	// then
	assert.NoError(t, err)

	// when
	block2 := getNewBlock(block1.GetSeal(), 1)
	err = cbr.AddBlock(block2)
	// then
	assert.NoError(t, err)

	blockQueryApi := api_gateway.NewBlockQueryApi(nil, cbr)

	// when
	block3, err := blockQueryApi.GetLastCommitedBlock()
	// then
	assert.NoError(t, err)
	assert.Equal(t, block2.GetSeal(), block3.GetSeal())
	assert.Equal(t, block2.GetHeight(), block3.GetHeight())
	assert.Equal(t, block2.GetPrevSeal(), block3.GetPrevSeal())
}

func TestBlockQueryApi_GetCommitedBlockByHeight(t *testing.T) {
	dbPath := "./.db"

	// when
	cbr, err := api_gateway.NewCommitedBlockRepositoryImpl(dbPath)
	// then
	assert.Equal(t, nil, err)

	defer func() {
		cbr.Close()
		os.RemoveAll(dbPath)
	}()

	// when
	block1 := getNewBlock([]byte("genesis"), 0)
	err = cbr.AddBlock(block1)
	// then
	assert.NoError(t, err)

	// when
	block2 := getNewBlock(block1.GetSeal(), 1)
	err = cbr.AddBlock(block2)
	// then
	assert.NoError(t, err)

	blockQueryApi := api_gateway.NewBlockQueryApi(nil, cbr)

	// when
	block3, err := blockQueryApi.GetCommitedBlockByHeight(blockchain.BlockHeight(1))
	// then
	assert.NoError(t, err)
	assert.Equal(t, block2.GetSeal(), block3.GetSeal())
	assert.Equal(t, block2.GetHeight(), block3.GetHeight())
	assert.Equal(t, block2.GetPrevSeal(), block3.GetPrevSeal())
}

func TestCommitedBlockRepositoryImpl(t *testing.T) {
	dbPath := "./.db"

	// when
	cbr, err := api_gateway.NewCommitedBlockRepositoryImpl(dbPath)

	// then
	assert.Equal(t, nil, err)

	defer func() {
		cbr.Close()
		os.RemoveAll(dbPath)
	}()

	// when
	block1 := getNewBlock([]byte("genesis"), 0)
	err = cbr.AddBlock(block1)
	// then
	assert.NoError(t, err)

	// when
	block2 := getNewBlock(block1.GetSeal(), 1)
	err = cbr.AddBlock(block2)
	// then
	assert.NoError(t, err)

	// when
	block3, err2 := cbr.GetLastBlock()
	// then
	assert.NoError(t, err2)
	assert.Equal(t, blockchain.BlockHeight(1), block3.GetHeight())

	// when
	block4, err3 := cbr.GetBlockByHeight(blockchain.BlockHeight(1))
	assert.NoError(t, err3)
	assert.Equal(t, 4, len(block4.GetTxList()))

}

func getNewBlock(prevSeal []byte, height uint64) *blockchain.DefaultBlock {
	validator := &blockchain.DefaultValidator{}
	testingTime := time.Now()
	blockCreator := []byte("testUser")
	txList := getTxList(testingTime)
	block := &blockchain.DefaultBlock{}
	block.SetPrevSeal(prevSeal)
	block.SetHeight(height)
	block.SetCreator(blockCreator)
	block.SetTimestamp(testingTime)
	for _, tx := range txList {
		block.PutTx(tx)
	}
	txSeal, _ := validator.BuildTxSeal(convertTxListType(txList))
	block.SetTxSeal(txSeal)

	seal, _ := validator.BuildSeal(block.GetTimestamp(), block.GetPrevSeal(), block.GetTxSeal(), block.GetCreator())
	block.SetSeal(seal)

	return block
}

func getTxList(testingTime time.Time) []*blockchain.DefaultTransaction {
	return []*blockchain.DefaultTransaction{
		{
			PeerID:    "p01",
			ID:        "tx01",
			Status:    0,
			Timestamp: testingTime,
			TxData: blockchain.TxData{
				Jsonrpc: "jsonRPC01",
				Method:  "invoke",
				Params: blockchain.Params{
					Type:     0,
					Function: "function01",
					Args:     []string{"arg1", "arg2"},
				},
				ID: "txdata01",
			},
		},
		{
			PeerID:    "p02",
			ID:        "tx02",
			Status:    0,
			Timestamp: testingTime,
			TxData: blockchain.TxData{
				Jsonrpc: "jsonRPC02",
				Method:  "invoke",
				Params: blockchain.Params{
					Type:     0,
					Function: "function02",
					Args:     []string{"arg1", "arg2"},
				},
				ID: "txdata02",
			},
		},
		{
			PeerID:    "p03",
			ID:        "tx03",
			Status:    0,
			Timestamp: testingTime,
			TxData: blockchain.TxData{
				Jsonrpc: "jsonRPC03",
				Method:  "invoke",
				Params: blockchain.Params{
					Type:     0,
					Function: "function03",
					Args:     []string{"arg1", "arg2"},
				},
				ID: "txdata03",
			},
		},
		{
			PeerID:    "p04",
			ID:        "tx04",
			Status:    0,
			Timestamp: testingTime,
			TxData: blockchain.TxData{
				Jsonrpc: "jsonRPC04",
				Method:  "invoke",
				Params: blockchain.Params{
					Type:     0,
					Function: "function04",
					Args:     []string{"arg1", "arg2"},
				},
				ID: "txdata04",
			},
		},
	}
}

func convertTxListType(txList []*blockchain.DefaultTransaction) []blockchain.Transaction {
	convTxList := make([]common.Transaction, 0)
	for _, tx := range txList {
		convTxList = append(convTxList, tx)
	}

	return convTxList
}
