package mock

import (
	"time"

	"github.com/it-chain/engine/blockchain"
	"github.com/it-chain/engine/common/command"
)

func GetTxResults() []command.TxResult {
	return []command.TxResult{
		{
			TxId: "txid1",
			Data: map[string]string{
				"key1": "value1",
			},
			Success: true,
		},
		{
			TxId: "txid2",
			Data: map[string]string{
				"key2": "value2",
			},
			Success: true,
		},
	}
}

func GetZeroLengthTxResults() []command.TxResult {
	return []command.TxResult{}
}

func GetFailedTxResults() []command.TxResult {
	return []command.TxResult{
		{
			TxId: "txid1",
			Data: map[string]string{
				"key1": "value1",
			},
			Success: true,
		},
		{
			TxId: "txid2",
			Data: map[string]string{
				"key2": "value2",
			},
			Success: false,
		},
	}
}

func GetStagedBlockWithId(blockId string) blockchain.Block {
	testingTime := time.Now()
	return &blockchain.DefaultBlock{
		Seal:      []byte(blockId),
		PrevSeal:  []byte{0x2},
		Height:    blockchain.BlockHeight(1),
		TxList:    getTxList(testingTime),
		TxSeal:    [][]byte{{0x1}},
		Timestamp: testingTime,
		Creator:   []byte("creator01"),
		State:     blockchain.Staged,
	}
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
