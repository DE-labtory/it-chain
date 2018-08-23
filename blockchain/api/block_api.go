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
	"fmt"

	"github.com/it-chain/engine/blockchain"
	"github.com/it-chain/engine/common/event"
	"github.com/it-chain/engine/common/logger"
)

type BlockApi struct {
	publisherId      string
	blockRepository  blockchain.BlockRepository
	eventService     blockchain.EventService
	consensusService blockchain.ConsensusService
}

func NewBlockApi(publisherId string, blockRepository blockchain.BlockRepository, eventService blockchain.EventService, consensusService blockchain.ConsensusService) (BlockApi, error) {
	return BlockApi{
		publisherId:      publisherId,
		blockRepository:  blockRepository,
		eventService:     eventService,
		consensusService: consensusService,
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

func (bApi BlockApi) CommitGenesisBlock(GenesisConfPath string) error {
	logger.Info(nil, "[Blockchain] Committing genesis block")

	// create
	GenesisBlock, err := blockchain.CreateGenesisBlock(GenesisConfPath)

	if err != nil {
		return ErrCreateGenesisBlock
	}

	// save(commit)
	GenesisBlock.SetState(blockchain.Committed)

	err = bApi.blockRepository.Save(GenesisBlock)

	if err != nil {
		return ErrSaveBlock
	}

	// publish
	commitEvent, err := createBlockCommittedEvent(GenesisBlock)

	if err != nil {
		return ErrCreateEvent
	}

	logger.Info(nil, fmt.Sprintf("[Blockchain] Genesis block has Committed - seal: [%x], height: [%d]", GenesisBlock.Seal, GenesisBlock.Height))

	return bApi.eventService.Publish("block.committed", commitEvent)
}

func (api BlockApi) CreateProposedBlock(txList []*blockchain.DefaultTransaction) error {
	logger.Info(nil, "[Blockchain] Creating proposed block")

	proposedBlock, err := api.createBlock(txList)
	if err != nil {
		return err
	}

	return api.consensusService.ConsensusBlock(proposedBlock)
}

func (bApi BlockApi) CommitProposedBlock(txList []*blockchain.DefaultTransaction) error {
	logger.Info(nil, "[Blockchain] Committing proposed block")

	// create
	ProposedBlock, err := bApi.createBlock(txList)
	if err != nil {
		return err
	}

	// save(commit)
	ProposedBlock.SetState(blockchain.Committed)

	err = bApi.blockRepository.Save(ProposedBlock)

	if err != nil {
		return ErrSaveBlock
	}

	// publish
	commitEvent, err := createBlockCommittedEvent(ProposedBlock)

	if err != nil {
		return ErrCreateEvent
	}

	logger.Info(nil, fmt.Sprintf("[Blockchain] Proposed block has Committed - seal: [%x],  height: [%d]", ProposedBlock.Seal, ProposedBlock.Height))

	return bApi.eventService.Publish("block.committed", commitEvent)
}

func (api BlockApi) createBlock(txList []*blockchain.DefaultTransaction) (blockchain.DefaultBlock, error) {
	lastBlock, err := api.blockRepository.FindLast()
	if err != nil {
		return blockchain.DefaultBlock{}, ErrGetLastBlock
	}

	prevSeal := lastBlock.GetSeal()
	height := lastBlock.GetHeight() + 1
	creator := api.publisherId

	block, err := blockchain.CreateProposedBlock(prevSeal, height, txList, []byte(creator))
	if err != nil {
		return blockchain.DefaultBlock{}, ErrCreateProposedBlock
	}

	return block, nil
}

func createBlockCommittedEvent(block blockchain.DefaultBlock) (event.BlockCommitted, error) {

	txList := blockchain.ConvBackFromTransactionList(block.TxList)

	return event.BlockCommitted{
		Seal:      block.GetSeal(),
		PrevSeal:  block.GetPrevSeal(),
		Height:    block.GetHeight(),
		TxList:    txList,
		TxSeal:    block.GetTxSeal(),
		Timestamp: block.GetTimestamp(),
		Creator:   block.GetCreator(),
		State:     block.GetState(),
	}, nil
}
