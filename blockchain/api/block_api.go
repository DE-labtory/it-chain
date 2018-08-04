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

package api

import (
	"errors"

	"github.com/it-chain/engine/blockchain"
)

var ErrGetLastBlock = errors.New("Error in getting last block")
var ErrCreateProposedBlock = errors.New("Error in creating proposed block")
var ErrFailBlockTypeCasting = errors.New("Error failed type casting block")

type BlockApi struct {
	publisherId       string
	blockQueryService blockchain.BlockQueryService
	blockRepository   blockchain.BlockRepository
}

func NewBlockApi(publisherId string, blockQueryService blockchain.BlockQueryService) (BlockApi, error) {
	return BlockApi{
		publisherId:       publisherId,
		blockQueryService: blockQueryService,
	}, nil
}

// TODO: Check 과정에서 임의의 노드에게서 받은 blockchain 정보로 동기화 되었는지 확인한다.
func (bApi BlockApi) SyncedCheck(block blockchain.Block) error {
	return nil
}

// 받은 block을 block pool에 추가한다.
func (bApi BlockApi) AddBlockToPool(block blockchain.Block) error {
	return nil
}

func (bApi BlockApi) CheckAndSaveBlockFromPool(height blockchain.BlockHeight) error {
	return nil
}

func (bApi BlockApi) SyncIsProgressing() blockchain.ProgressState {
	return blockchain.DONE
}

func (bApi BlockApi) CreateProposedBlock(txList []blockchain.Transaction) (blockchain.DefaultBlock, error) {

	lastBlock, err := bApi.blockRepository.FindLast()

	if err != nil {
		return blockchain.DefaultBlock{}, ErrGetLastBlock
	}

	prevSeal := lastBlock.GetSeal()
	height := lastBlock.GetHeight() + 1
	defaultTxList := blockchain.GetBackTxType(txList)
	creator := bApi.publisherId

	block, err := blockchain.CreateProposedBlock(prevSeal, height, defaultTxList, []byte(creator))

	if err != nil {
		return blockchain.DefaultBlock{}, ErrCreateProposedBlock
	}

	defaultBlock, ok := block.(*blockchain.DefaultBlock)
	if !ok {
		return blockchain.DefaultBlock{}, ErrFailBlockTypeCasting
	}

	return *defaultBlock, nil
}

func (bApi BlockApi) CreateGenesisBlock(GenesisConfPath string) error {

	GenesisBlock, err := blockchain.CreateGenesisBlock(GenesisConfPath)

	if err != nil {
		return ErrCreateGenesisBlock
	}

	err = blockchain.CommitBlock(GenesisBlock)

	if err != nil {
		return ErrCommitBlock
	}

	return nil
}
