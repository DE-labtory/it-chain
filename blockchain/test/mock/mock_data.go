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
	"log"
	"math/rand"
	"time"

	"github.com/it-chain/engine/blockchain"
	"github.com/it-chain/engine/common/command"
	"github.com/it-chain/yggdrasill/common"
	"github.com/rs/xid"
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
		TxList:    GetTxList(testingTime),
		TxSeal:    [][]byte{{0x1}},
		Timestamp: testingTime,
		Creator:   "creator01",
		State:     blockchain.Staged,
	}
}

func GetNewBlock(prevSeal []byte, height uint64) *blockchain.DefaultBlock {
	validator := &blockchain.DefaultValidator{}
	testingTime := time.Now().Round(0)
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

	tree, err := validator.BuildTree(ConvertTxListType(txList))
	if err != nil {
		//ToDo: 좀 더 나은처리가 있을 지.
		log.Println(err)
	}

	block.SetTree(tree)

	txSeal, _ := validator.BuildTxSeal(ConvertTxListType(txList))
	block.SetTxSeal(txSeal)

	seal, _ := validator.BuildSeal(block.GetTimestamp(), block.GetPrevSeal(), block.GetTxSealRoot(), block.GetCreator())
	block.SetSeal(seal)

	return block
}

func GetRandomTxList() []*blockchain.DefaultTransaction {
	randSource := rand.NewSource(time.Now().UnixNano())
	randInstance := rand.New(randSource)
	randomTxListLength := randInstance.Intn(100)

	TxList := []*blockchain.DefaultTransaction{}

	for i := 0; i <= randomTxListLength; i++ {
		Tx := &blockchain.DefaultTransaction{
			ID:        xid.New().String(),
			ICodeID:   xid.New().String(),
			PeerID:    xid.New().String(),
			Timestamp: time.Now().Round(0),
			Jsonrpc:   xid.New().String(),
			Function:  xid.New().String(),
			Args:      nil,
			Signature: []byte(xid.New().String()),
		}

		TxList = append(TxList, Tx)
	}

	return TxList
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
			Signature: []byte("signature01"),
		},
		{

			ID:        "tx02",
			ICodeID:   "ICode02",
			PeerID:    "p02",
			Timestamp: testingTime,
			Jsonrpc:   "jsonRPC02",
			Function:  "function02",
			Args:      []string{"arg1", "arg2"},
			Signature: []byte("signature02"),
		},
		{
			ID:        "tx03",
			ICodeID:   "ICode03",
			PeerID:    "p03",
			Timestamp: testingTime,
			Jsonrpc:   "jsonRPC03",
			Function:  "function03",
			Args:      []string{"arg1", "arg2"},
			Signature: []byte("signature03"),
		},
		{
			ID:        "tx04",
			ICodeID:   "ICode04",
			PeerID:    "p04",
			Timestamp: testingTime,
			Jsonrpc:   "jsonRPC04",
			Function:  "function04",
			Args:      []string{"arg1", "arg2"},
			Signature: []byte("signature04"),
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
