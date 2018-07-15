package api_gateway_test

import (
	"testing"

	"os"

	"time"

	"github.com/it-chain/it-chain-Engine/api_gateway"
	"github.com/it-chain/it-chain-Engine/blockchain"
	"github.com/it-chain/yggdrasill/common"
	"github.com/stretchr/testify/assert"
)

func TestBlockPoolRepositoryImpl_AddCreatedBlock(t *testing.T) {
	bpr := api_gateway.NewBlockpoolRepositoryImpl()

	// when
	block1 := &blockchain.DefaultBlock{
		Seal:   []byte{0x1},
		Height: blockchain.BlockHeight(1),
	}
	// when
	bpr.AddCreatedBlock(block1)

	// then
	assert.Equal(t, 1, len(bpr.Blocks))
	assert.Equal(t, []byte{0x1}, bpr.Blocks[0].GetSeal())
}

func TestBlockPoolRepositoryImpl_GetStagedBlockByHeight(t *testing.T) {
	bpr := api_gateway.NewBlockpoolRepositoryImpl()

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
	bpr := api_gateway.NewBlockpoolRepositoryImpl()

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
	testingTime := getTime()
	blockCreator := []byte("testUser")
	txList := getTxList(testingTime)
	block := newEmptyBlock(prevSeal, height, blockCreator)
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

func newEmptyBlock(prevSeal []byte, height uint64, creator []byte) *blockchain.DefaultBlock {
	block := &blockchain.DefaultBlock{}

	block.SetPrevSeal(prevSeal)
	block.SetHeight(height)
	block.SetCreator(creator)

	return block
}

func getTxList(testingTime time.Time) []*blockchain.DefaultTransaction {
	return []*blockchain.DefaultTransaction{
		{
			PeerID:    "p01",
			ID:        "tx01",
			Status:    0,
			Timestamp: testingTime,
			TxData: &blockchain.TxData{
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
			TxData: &blockchain.TxData{
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
			TxData: &blockchain.TxData{
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
			TxData: &blockchain.TxData{
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

func getTime() time.Time {
	testingTime, _ := time.Parse("Jan 2, 2006 at 3:04pm (MST)", "Feb 3, 2013 at 7:54pm (UTC)")
	return testingTime
}

func convertTxListType(txList []*blockchain.DefaultTransaction) []blockchain.Transaction {
	convTxList := make([]common.Transaction, 0)
	for _, tx := range txList {
		convTxList = append(convTxList, tx)
	}

	return convTxList
}
