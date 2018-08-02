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

func GetStagedBlockWithId(blockId string) blockchain.DefaultBlock {
	testingTime := time.Now()
	return blockchain.DefaultBlock{
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
