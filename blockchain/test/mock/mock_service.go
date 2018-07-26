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

type BlockQueryService struct {
	GetStagedBlockByHeightFunc   func(height blockchain.BlockHeight) (blockchain.DefaultBlock, error)
	GetStagedBlockByIdFunc       func(blockId string) (blockchain.DefaultBlock, error)
	GetLastCommitedBlockFunc     func() (blockchain.DefaultBlock, error)
	GetCommitedBlockByHeightFunc func(height blockchain.BlockHeight) (blockchain.DefaultBlock, error)
}

func (s BlockQueryService) GetStagedBlockByHeight(height blockchain.BlockHeight) (blockchain.DefaultBlock, error) {
	return s.GetStagedBlockByHeightFunc(height)
}
func (s BlockQueryService) GetStagedBlockById(blockId string) (blockchain.DefaultBlock, error) {
	return s.GetStagedBlockByIdFunc(blockId)
}
func (s BlockQueryService) GetLastCommitedBlock() (blockchain.DefaultBlock, error) {
	return s.GetLastCommitedBlockFunc()
}
func (s BlockQueryService) GetCommitedBlockByHeight(height blockchain.BlockHeight) (blockchain.DefaultBlock, error) {
	return s.GetCommitedBlockByHeightFunc(height)
}

type BlockExecuteService struct {
	ExecuteBlockFunc func(block blockchain.Block) error
}

func (s BlockExecuteService) ExecuteBlock(block blockchain.Block) error {
	return s.ExecuteBlockFunc(block)
}
