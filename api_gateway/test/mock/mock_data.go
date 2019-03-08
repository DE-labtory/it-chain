package mock

import (
	"time"

	"github.com/DE-labtory/engine/blockchain"
	"github.com/DE-labtory/engine/common/event"
	"github.com/DE-labtory/yggdrasill/common"
)

func GetNewBlock(prevSeal []byte, height uint64) *blockchain.DefaultBlock {
	validator := &blockchain.DefaultValidator{}
	testingTime := time.Now()
	blockCreator := "testUser"
	txList := GetTxList(testingTime)
	block := &blockchain.DefaultBlock{}
	block.SetPrevSeal(prevSeal)
	block.SetHeight(height)
	block.SetCreator(blockCreator)
	block.SetTimestamp(testingTime)
	for _, tx := range txList {
		block.PutTx(tx)
	}
	txSeal, _ := validator.BuildTxSeal(ConvertTxListType(txList))
	block.SetTxSeal(txSeal)

	seal, _ := validator.BuildSeal(block.GetTimestamp(), block.GetPrevSeal(), block.GetTxSeal(), block.GetCreator())
	block.SetSeal(seal)

	return block
}

func GetTxList(testingTime time.Time) []*blockchain.DefaultTransaction {
	return []*blockchain.DefaultTransaction{
		{
			ID:        "tx01",
			ICodeID:   "ICode01",
			PeerID:    "p01",
			Timestamp: testingTime,
			Jsonrpc:   "jsonRPC01",
			Function:  "function01",
			Args:      []string{"arg1", "arg2"},
		},
		{

			ID:        "tx02",
			ICodeID:   "ICode02",
			PeerID:    "p02",
			Timestamp: testingTime,
			Jsonrpc:   "jsonRPC02",
			Function:  "function02",
			Args:      []string{"arg1", "arg2"},
		},
		{
			ID:        "tx03",
			ICodeID:   "ICode03",
			PeerID:    "p03",
			Timestamp: testingTime,
			Jsonrpc:   "jsonRPC03",
			Function:  "function03",
			Args:      []string{"arg1", "arg2"},
		},
		{
			ID:        "tx04",
			ICodeID:   "ICode04",
			PeerID:    "p04",
			Timestamp: testingTime,
			Jsonrpc:   "jsonRPC04",
			Function:  "function04",
			Args:      []string{"arg1", "arg2"},
		},
	}
}

func ConvertTxListType(txList []*blockchain.DefaultTransaction) []blockchain.Transaction {
	convTxList := make([]common.Transaction, 0)
	for _, tx := range txList {
		convTxList = append(convTxList, tx)
	}

	return convTxList
}

func GetEventTxList() []event.Tx {
	txList := make([]event.Tx, 0)

	tx := event.Tx{
		ID:        "1",
		ICodeID:   "icode1",
		PeerID:    "peer1",
		TimeStamp: time.Now(),
		Jsonrpc:   "jsonrpc1",
		Function:  "fn1",
		Args:      []string{"arg1", "arg2"},
		Signature: []byte{0x4},
	}

	txList = append(txList, tx)

	return txList
}
