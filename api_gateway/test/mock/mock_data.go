package mock

import (
	"time"

	"github.com/it-chain/engine/blockchain"
	"github.com/it-chain/engine/common/event"
	"github.com/it-chain/yggdrasill/common"
)

func GetNewBlock(prevSeal []byte, height uint64) *blockchain.DefaultBlock {
	validator := &blockchain.DefaultValidator{}
	testingTime := time.Now()
	blockCreator := []byte("testUser")
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
		Status:    1,
		PeerID:    "peer1",
		TimeStamp: time.Now(),
		Jsonrpc:   "jsonrpc1",
		Method:    "mtd1",
		Function:  "fn1",
		Args:      []string{"arg1", "arg2"},
		Signature: []byte{0x4},
	}

	txList = append(txList, tx)

	return txList
}
