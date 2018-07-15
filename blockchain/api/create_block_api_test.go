package api_test

import (
	"os"
	"testing"
	"time"

	"github.com/it-chain/it-chain-Engine/api_gateway"
	"github.com/it-chain/it-chain-Engine/blockchain"
	"github.com/it-chain/it-chain-Engine/blockchain/api"
	"github.com/it-chain/it-chain-Engine/blockchain/test/mock"
	"github.com/stretchr/testify/assert"
)

func TestCreateBlockApi_CreateBlock(t *testing.T) {
	dbPath := "./.db"
	blockRepository, err := api_gateway.NewCommitedBlockRepositoryImpl(dbPath)
	assert.NoError(t, err)

	blockService := mock.BlockService{}
	blockService.ExecuteBlockFunc = func(block blockchain.Block) error {
		assert.Equal(t, blockchain.BlockHeight(2), block.GetHeight())
		return nil
	}

	defer func() {
		blockRepository.Close()
		os.RemoveAll(dbPath)
	}()

	// when
	block1 := getNewBlock([]byte("genesis"), 0)
	err = blockRepository.AddBlock(block1)
	assert.NoError(t, err)

	// when
	block2 := getNewBlock(block1.GetSeal(), 1)
	err = blockRepository.AddBlock(block2)
	assert.NoError(t, err)

	// when
	queryApi := api_gateway.BlockQueryApi{}
	queryApi.CommitedBlockRepository = blockRepository

	// when
	blockApi := api.NewCreateBlockApi(queryApi, blockService, "zf")

	txList := blockchain.ConvertTxTypeToTransaction(getTxList(time.Now()))

	// when
	err = blockApi.CreateBlock(txList)

	// then
	assert.NoError(t, err)

}

// TODO: move to mock package
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

// TODO: move to mock package
func newEmptyBlock(prevSeal []byte, height uint64, creator []byte) *blockchain.DefaultBlock {
	block := &blockchain.DefaultBlock{}

	block.SetPrevSeal(prevSeal)
	block.SetHeight(height)
	block.SetCreator(creator)

	return block
}

// TODO: move to mock package
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
					Function: "function04",
					Args:     []string{"arg1", "arg2"},
				},
				ID: "txdata04",
			},
		},
	}
}

// TODO: move to mock package
func getTime() time.Time {
	testingTime, _ := time.Parse("Jan 2, 2006 at 3:04pm (MST)", "Feb 3, 2013 at 7:54pm (UTC)")
	return testingTime
}

// TODO: move to mock package
func convertTxListType(txList []*blockchain.DefaultTransaction) []blockchain.Transaction {
	convTxList := make([]blockchain.Transaction, 0)
	for _, tx := range txList {
		convTxList = append(convTxList, tx)
	}

	return convTxList
}
