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

import "github.com/it-chain/engine/blockchain"

type BlockQueryApi struct {
	GetLastBlockFunc         func() (blockchain.Block, error)
	GetBlockByHeightFunc     func(blockHeight uint64) (blockchain.Block, error)
	GetBlockBySealFunc       func(seal []byte) (blockchain.Block, error)
	GetBlockByTxIDFunc       func(txid string) (blockchain.Block, error)
	GetTransactionByTxIDFunc func(txid string) (blockchain.Transaction, error)
}

func (br BlockQueryApi) GetLastBlock() (blockchain.Block, error) {
	return br.GetLastBlockFunc()
}
func (br BlockQueryApi) GetBlockByHeight(blockHeight uint64) (blockchain.Block, error) {
	return br.GetBlockByHeightFunc(blockHeight)
}
func (br BlockQueryApi) GetBlockBySeal(seal []byte) (blockchain.Block, error) {
	return br.GetBlockBySealFunc(seal)
}
func (br BlockQueryApi) GetBlockByTxID(txid string) (blockchain.Block, error) {
	return br.GetBlockByTxIDFunc(txid)
}
func (br BlockQueryApi) GetTransactionByTxID(txid string) (blockchain.Transaction, error) {
	return br.GetTransactionByTxIDFunc(txid)
}
