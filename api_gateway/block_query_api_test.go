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

	"encoding/hex"

	"github.com/it-chain/engine/api_gateway"
	"github.com/it-chain/engine/api_gateway/test/mock"
	"github.com/it-chain/engine/blockchain"
	"github.com/it-chain/engine/common/event"
	"github.com/it-chain/midgard"
	"github.com/stretchr/testify/assert"
)

func TestBlockQueryApi_GetStagedBlockByHeight(t *testing.T) {
	pool := api_gateway.NewBlockPoolRepository()

	block1 := &blockchain.DefaultBlock{
		Seal:   []byte{0x1},
		Height: blockchain.BlockHeight(1),
		State:  blockchain.Created,
	}
	block2 := &blockchain.DefaultBlock{
		Seal:   []byte{0x2},
		Height: blockchain.BlockHeight(2),
		State:  blockchain.Staged,
	}

	pool.Blocks = append(pool.Blocks, block1)
	pool.Blocks = append(pool.Blocks, block2)

	// when
	qa := api_gateway.NewBlockQueryApi(pool, nil)
	// when
	b1, err1 := qa.GetStagedBlockByHeight(1)
	// then
	assert.Equal(t, api_gateway.ErrNoStagedBlock, err1)
	assert.Equal(t, blockchain.DefaultBlock{}, b1)

	// when
	b2, err2 := qa.GetStagedBlockByHeight(2)
	// then
	assert.Equal(t, nil, err2)
	assert.Equal(t, []byte{0x2}, b2.GetSeal())
	assert.Equal(t, blockchain.Staged, b2.State)
	assert.Equal(t, blockchain.BlockHeight(2), b2.Height)
}

func TestBlockQueryApi_GetStagedBlockById(t *testing.T) {
	pool := api_gateway.NewBlockPoolRepository()

	block1 := &blockchain.DefaultBlock{
		Seal:   []byte{0x1},
		Height: blockchain.BlockHeight(1),
		State:  blockchain.Created,
	}
	block2 := &blockchain.DefaultBlock{
		Seal:   []byte{0x2},
		Height: blockchain.BlockHeight(2),
		State:  blockchain.Staged,
	}

	pool.Blocks = append(pool.Blocks, block1)
	pool.Blocks = append(pool.Blocks, block2)

	// when
	qa := api_gateway.NewBlockQueryApi(pool, nil)
	// when
	bId1 := hex.EncodeToString([]byte{0x1})
	b1, err1 := qa.GetStagedBlockById(bId1)
	// then
	assert.Equal(t, api_gateway.ErrNoStagedBlock, err1)
	assert.Equal(t, blockchain.DefaultBlock{}, b1)

	// when
	bId2 := hex.EncodeToString([]byte{0x2})
	b2, err2 := qa.GetStagedBlockById(bId2)
	// then
	assert.Equal(t, nil, err2)
	assert.Equal(t, []byte{0x2}, b2.GetSeal())
	assert.Equal(t, blockchain.Staged, b2.State)
	assert.Equal(t, blockchain.BlockHeight(2), b2.Height)
}

func TestBlockPoolRepositoryImpl_SaveCreatedBlock(t *testing.T) {
	bpr := api_gateway.NewBlockPoolRepository()

	// when
	block1 := &blockchain.DefaultBlock{
		Seal:   []byte{0x1},
		Height: blockchain.BlockHeight(1),
		State:  blockchain.Created,
	}
	// when
	err := bpr.SaveCreatedBlock(*block1)
	// then
	assert.Equal(t, nil, err)
	assert.Equal(t, 1, len(bpr.Blocks))
	assert.Equal(t, []byte{0x1}, bpr.Blocks[0].GetSeal())

	// when
	block2 := &blockchain.DefaultBlock{
		Seal:   []byte{0x2},
		Height: blockchain.BlockHeight(2),
		State:  blockchain.Staged,
	}
	// when
	err2 := bpr.SaveCreatedBlock(*block2)
	// then
	assert.Equal(t, api_gateway.ErrInvalidStateBlock, err2)

	// when
	block3 := &blockchain.DefaultBlock{
		Seal:     []byte{0x1},
		PrevSeal: []byte{0x5},
		Height:   blockchain.BlockHeight(3),
		State:    blockchain.Created,
	}
	// when
	err3 := bpr.SaveCreatedBlock(*block3)
	// then
	assert.NoError(t, err3)
	assert.Equal(t, 1, len(bpr.Blocks))

	// when
	block := bpr.Blocks[0]
	// then
	assert.Equal(t, []byte{0x1}, block.GetSeal())
	assert.Equal(t, []byte{0x5}, block.GetPrevSeal())

}

func TestBlockPoolRepositoryImpl_SaveStagedBlock(t *testing.T) {
	bpr := api_gateway.NewBlockPoolRepository()

	// when
	block1 := &blockchain.DefaultBlock{
		Seal:   []byte{0x1},
		Height: blockchain.BlockHeight(1),
		State:  blockchain.Staged,
	}
	// when
	err := bpr.SaveStagedBlock(*block1)
	// then
	assert.Equal(t, nil, err)
	assert.Equal(t, 1, len(bpr.Blocks))
	assert.Equal(t, []byte{0x1}, bpr.Blocks[0].GetSeal())

	// when
	block2 := &blockchain.DefaultBlock{
		Seal:   []byte{0x2},
		Height: blockchain.BlockHeight(2),
		State:  blockchain.Created,
	}
	// when
	err2 := bpr.SaveStagedBlock(*block2)
	// then
	assert.Equal(t, api_gateway.ErrInvalidStateBlock, err2)

	// when
	block3 := &blockchain.DefaultBlock{
		Seal:     []byte{0x1},
		PrevSeal: []byte{0x5},
		Height:   blockchain.BlockHeight(3),
		State:    blockchain.Staged,
	}
	// when
	err3 := bpr.SaveStagedBlock(*block3)
	// then
	assert.NoError(t, err3)
	assert.Equal(t, 1, len(bpr.Blocks))

	// when
	block := bpr.Blocks[0]
	// then
	assert.Equal(t, []byte{0x1}, block.GetSeal())
	assert.Equal(t, []byte{0x5}, block.GetPrevSeal())

}

func TestBlockPoolRepositoryImpl_FindCreatedBlockById(t *testing.T) {
	bpr := api_gateway.NewBlockPoolRepository()

	// when
	bpr.Blocks = append(bpr.Blocks, &blockchain.DefaultBlock{
		Seal:     []byte{0x1},
		PrevSeal: []byte{0x1},
		Height:   uint64(1),
		State:    blockchain.Created,
	})

	// when
	bpr.Blocks = append(bpr.Blocks, &blockchain.DefaultBlock{
		Seal:     []byte{0x2},
		PrevSeal: []byte{0x2},
		Height:   uint64(2),
		State:    blockchain.Staged,
	})

	// when
	blockID := hex.EncodeToString([]byte{0x1})
	block, err := bpr.FindCreatedBlockById(blockID)

	// then
	assert.Equal(t, err, nil)
	assert.Equal(t, uint64(1), block.GetHeight())
	assert.Equal(t, []byte{0x1}, block.GetSeal())
	assert.Equal(t, []byte{0x1}, block.GetPrevSeal())

	// when
	_, err2 := bpr.FindCreatedBlockById(string([]byte{0x2}))

	// then
	assert.Equal(t, err2, api_gateway.ErrNoCreatedBlock)

	// when
	_, err3 := bpr.FindCreatedBlockById(string([]byte{0x3}))

	// then
	assert.Equal(t, err3, api_gateway.ErrNoCreatedBlock)

}

func TestBlockPoolRepositoryImpl_SaveCreatedBlock_InvalidStateBlock(t *testing.T) {
	bpr := api_gateway.NewBlockPoolRepository()

	// when
	block1 := &blockchain.DefaultBlock{
		Seal:   []byte{0x1},
		Height: blockchain.BlockHeight(1),
		State:  blockchain.Committed,
	}
	// when
	err := bpr.SaveCreatedBlock(*block1)

	// then
	assert.Equal(t, api_gateway.ErrInvalidStateBlock, err)
	assert.Equal(t, 0, len(bpr.Blocks))
}

func TestBlockPoolRepositoryImpl_FindStagedBlockByHeight(t *testing.T) {
	bpr := api_gateway.NewBlockPoolRepository()

	// when
	bpr.Blocks = append(bpr.Blocks, &blockchain.DefaultBlock{
		Seal:     []byte{0x1},
		PrevSeal: []byte{0x1},
		Height:   uint64(1),
		State:    blockchain.Staged,
	})

	// when
	block, err := bpr.FindStagedBlockByHeight(1)

	// then
	assert.Equal(t, err, nil)
	assert.Equal(t, uint64(1), block.GetHeight())
	assert.Equal(t, []byte{0x1}, block.GetSeal())
	assert.Equal(t, []byte{0x1}, block.GetPrevSeal())

	// when
	_, err2 := bpr.FindStagedBlockByHeight(133)

	// then
	assert.Equal(t, err2, api_gateway.ErrNoStagedBlock)
}

func TestBlockPoolRepositoryImpl_FindStagedBlockById(t *testing.T) {
	bpr := api_gateway.NewBlockPoolRepository()

	// when
	bpr.Blocks = append(bpr.Blocks, &blockchain.DefaultBlock{
		Seal:     []byte{0x1},
		PrevSeal: []byte{0x1},
		Height:   uint64(1),
		State:    blockchain.Staged,
	})

	// when
	bpr.Blocks = append(bpr.Blocks, &blockchain.DefaultBlock{
		Seal:     []byte{0x2},
		PrevSeal: []byte{0x2},
		Height:   uint64(2),
		State:    blockchain.Created,
	})

	// when
	blockID := hex.EncodeToString([]byte{0x1})
	block, err := bpr.FindStagedBlockById(blockID)

	// then
	assert.Equal(t, err, nil)
	assert.Equal(t, uint64(1), block.GetHeight())
	assert.Equal(t, []byte{0x1}, block.GetSeal())
	assert.Equal(t, []byte{0x1}, block.GetPrevSeal())

	// when
	blockID2 := hex.EncodeToString([]byte{0x2})
	_, err2 := bpr.FindStagedBlockById(blockID2)
	// then
	assert.Equal(t, err2, api_gateway.ErrNoStagedBlock)

	// when
	blockID3 := hex.EncodeToString([]byte{0x3})
	_, err3 := bpr.FindStagedBlockById(blockID3)

	// then
	assert.Equal(t, err3, api_gateway.ErrNoStagedBlock)
}

func TestBlockPoolRepositoryImpl_FindFirstStagedBlock_basic(t *testing.T) {
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
	block, err := bpr.FindFirstStagedBlock()

	assert.Equal(t, nil, err)
	assert.Equal(t, uint64(2), block.Height)
	assert.Equal(t, []byte{0x2}, block.Seal)
}

func TestBlockPoolRepositoryImpl_FindFirstStagedBlock_NoStagedBlockFound(t *testing.T) {
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
	block, err := bpr.FindFirstStagedBlock()

	assert.Equal(t, api_gateway.ErrNoStagedBlock, err)
	assert.Equal(t, true, block.IsEmpty())
}

func TestBlockPoolRepositoryImpl_FindFirstStagedBlock_lenIsZero(t *testing.T) {
	bpr := api_gateway.NewBlockPoolRepository()

	// when
	block, err := bpr.FindFirstStagedBlock()

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
	bId1 := hex.EncodeToString([]byte{0x1})
	err := bpr.RemoveBlockById(bId1)

	// then
	assert.Equal(t, nil, err)
	assert.Equal(t, 1, len(bpr.Blocks))
	assert.Equal(t, []byte{0x2}, bpr.Blocks[0].GetSeal())

	// when
	bId2 := hex.EncodeToString([]byte{0x3})
	err2 := bpr.RemoveBlockById(bId2)

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
	bId := hex.EncodeToString([]byte{0x3})
	err := bpr.RemoveBlockById(bId)

	// then
	assert.Equal(t, 2, len(bpr.Blocks))
	assert.Equal(t, api_gateway.ErrFailRemoveBlock, err)
}

func TestBlockQueryApi_FindLastCommitedBlock(t *testing.T) {
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
	block1 := mock.GetNewBlock([]byte("genesis"), 0)
	err = cbr.AddBlock(block1)
	// then
	assert.NoError(t, err)

	// when
	block2 := mock.GetNewBlock(block1.GetSeal(), 1)
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

func TestBlockQueryApi_FindCommitedBlockByHeight(t *testing.T) {
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
	block1 := mock.GetNewBlock([]byte("genesis"), 0)
	err = cbr.AddBlock(block1)
	// then
	assert.NoError(t, err)

	// when
	block2 := mock.GetNewBlock(block1.GetSeal(), 1)
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
	block1 := mock.GetNewBlock([]byte("genesis"), 0)
	err = cbr.AddBlock(block1)
	// then
	assert.NoError(t, err)

	// when
	block2 := mock.GetNewBlock(block1.GetSeal(), 1)
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

	// when
	AllBlock, err4 := cbr.FindAllBlock()

	// then
	assert.NoError(t, err4)
	assert.Equal(t, 2, len(AllBlock))

}

func TestBlockEventListener_HandleBlockCreatedEvent(t *testing.T) {
	blockpoolRepo := api_gateway.NewBlockPoolRepository()
	eh := api_gateway.NewBlockEventListener(blockpoolRepo, nil)

	event1 := event.BlockCreated{
		EventModel: midgard.EventModel{
			ID: "block_id1",
		},
		Seal:      []byte{0x1},
		PrevSeal:  []byte{0x2},
		Height:    uint64(3),
		TxList:    mock.GetEventTxList(),
		TxSeal:    [][]byte{{0x1}},
		Timestamp: time.Now(),
		Creator:   []byte{0x4},
		State:     blockchain.Created,
	}

	// when
	err := eh.HandleBlockCreatedEvent(event1)
	// then
	assert.Equal(t, err, nil)
	assert.Equal(t, 1, len(blockpoolRepo.Blocks))

	// when
	block := blockpoolRepo.Blocks[0]
	// then
	assert.Equal(t, []byte{0x1}, block.GetSeal())
	assert.Equal(t, []byte{0x2}, block.GetPrevSeal())
	assert.Equal(t, uint64(3), block.GetHeight())
	assert.Equal(t, blockchain.Created, (*block.(*blockchain.DefaultBlock)).State)
	assert.Equal(t, "1", block.GetTxList()[0].GetID())
	assert.Equal(t, []byte{0x4}, block.GetTxList()[0].GetSignature())
}

func TestBlockEventListener_HandleBlockStagedEvent(t *testing.T) {
	blockpoolRepo := api_gateway.NewBlockPoolRepository()
	eh := api_gateway.NewBlockEventListener(blockpoolRepo, nil)

	// when
	block1 := &blockchain.DefaultBlock{
		Seal:   []byte{0x1},
		Height: blockchain.BlockHeight(1),
		State:  blockchain.Created,
	}
	block2 := &blockchain.DefaultBlock{
		Seal:   []byte{0x2},
		Height: blockchain.BlockHeight(2),
		State:  blockchain.Staged,
	}

	blockpoolRepo.SaveCreatedBlock(*block1)
	blockpoolRepo.SaveCreatedBlock(*block2)

	// when
	block1ID := hex.EncodeToString([]byte{0x1})
	event1 := event.BlockStaged{
		EventModel: midgard.EventModel{
			ID: block1ID,
		},
		State: blockchain.Staged,
	}
	err1 := eh.HandleBlockStagedEvent(event1)
	// then
	assert.NoError(t, err1)

	// when
	defaultBlock, err2 := blockpoolRepo.FindStagedBlockById(block1ID)
	// then
	assert.NoError(t, err2)
	assert.Equal(t, []byte{0x1}, defaultBlock.Seal)
	assert.Equal(t, uint64(1), defaultBlock.Height)

	// when
	block2ID := hex.EncodeToString([]byte{0x2})
	event2 := event.BlockStaged{
		EventModel: midgard.EventModel{
			ID: block2ID,
		},
		State: blockchain.Staged,
	}
	err3 := eh.HandleBlockStagedEvent(event2)
	// then
	assert.Equal(t, api_gateway.ErrNoCreatedBlock, err3)

	// when
	block3ID := hex.EncodeToString([]byte{0x3})
	event3 := event.BlockStaged{
		EventModel: midgard.EventModel{
			ID: block3ID,
		},
		State: blockchain.Staged,
	}
	err4 := eh.HandleBlockStagedEvent(event3)
	// then
	assert.Equal(t, api_gateway.ErrNoCreatedBlock, err4)

}

func TestBlockEventListener_HandleBlockCommitedEvent(t *testing.T) {
	dbPath := "./.db"

	// when
	poolRepo := api_gateway.NewBlockPoolRepository()
	cbr, err := api_gateway.NewCommitedBlockRepositoryImpl(dbPath)
	// then
	assert.Equal(t, nil, err)

	defer func() {
		cbr.Close()
		os.RemoveAll(dbPath)
	}()

	// when
	block1 := mock.GetNewBlock([]byte("genesis"), 0)
	err = cbr.AddBlock(block1)
	// then
	assert.NoError(t, err)

	eh := api_gateway.NewBlockEventListener(poolRepo, cbr)

	// when
	block2 := mock.GetNewBlock(block1.GetSeal(), 1)
	block2.State = blockchain.Staged
	// when
	poolRepo.Blocks = append(poolRepo.Blocks, block2)
	// when
	block2ID := hex.EncodeToString(block2.GetSeal())
	event1 := event.BlockCommitted{
		EventModel: midgard.EventModel{
			ID: block2ID,
		},
		State: blockchain.Committed,
	}
	// when - Handle BlockCommited event
	err1 := eh.HandleBlockCommitedEvent(event1)
	// then
	assert.NoError(t, err1)

	// when - Test whether save target block to yggdrasill
	block3, err2 := cbr.FindBlockByHeight(1)
	// then
	assert.NoError(t, err2)
	assert.Equal(t, block3.Seal, block2.GetSeal())
	assert.Equal(t, blockchain.Committed, block3.State)

	// when - Test whether target block is removed from block pool
	block4, err3 := poolRepo.FindStagedBlockById(block2ID)
	// then
	assert.Equal(t, api_gateway.ErrNoStagedBlock, err3)
	assert.Equal(t, true, block4.IsEmpty())

}
