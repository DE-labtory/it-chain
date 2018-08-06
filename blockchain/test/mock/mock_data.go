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

package mock

import (
	"time"

	"github.com/it-chain/engine/blockchain"
	"github.com/it-chain/engine/common/command"
	"github.com/it-chain/yggdrasill/common"
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

func GetNewBlock(prevSeal []byte, height uint64) *blockchain.DefaultBlock {
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
	txSeal, _ := validator.BuildTxSeal(ConvertTxListType(txList))
	block.SetTxSeal(txSeal)

	seal, _ := validator.BuildSeal(block.GetTimestamp(), block.GetPrevSeal(), block.GetTxSeal(), block.GetCreator())
	block.SetSeal(seal)

	return block
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

func ConvertTxListType(txList []*blockchain.DefaultTransaction) []blockchain.Transaction {
	convTxList := make([]common.Transaction, 0)
	for _, tx := range txList {
		convTxList = append(convTxList, tx)
	}

	return convTxList
}
