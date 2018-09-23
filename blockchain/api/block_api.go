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
	"github.com/it-chain/engine/blockchain/infra/mem"
	"github.com/it-chain/engine/common/event"
	"github.com/it-chain/iLogger"
)

type BlockApi struct {
	publisherId     string
	blockRepository blockchain.BlockRepository
	eventService    blockchain.EventService
	BlockPool       *mem.BlockPool
}

func NewBlockApi(publisherId string, blockRepository blockchain.BlockRepository, eventService blockchain.EventService, blockPool *mem.BlockPool) (*BlockApi, error) {
	return &BlockApi{
		publisherId:     publisherId,
		blockRepository: blockRepository,
		eventService:    eventService,
		BlockPool:       blockPool,
	}, nil
}

func (bApi BlockApi) CheckAndSaveBlockFromPool(height blockchain.BlockHeight) error {
	return nil
}

func (bApi BlockApi) SyncIsProgressing() blockchain.ProgressState {
	return blockchain.DONE
}

func (api BlockApi) ConsentBlock(consensusType string, block blockchain.DefaultBlock) error {

	switch consensusType {
	case "solo":
		return api.CommitBlock(block)

	case "pbft":
		event, err := createBlockCreatedEvent(block)
		if err != nil {
			return err
		}
		return api.eventService.Publish("block.consent", event)

	default:
		iLogger.Error(nil, fmt.Sprintf("[blockchain] undefined consensus type: %v", consensusType))
		return ErrUndefinedConsensusType
	}

}

func (bApi BlockApi) CommitGenesisBlock(GenesisConfPath string) error {
	iLogger.Info(nil, "[Blockchain] Committing genesis block")

	// create
	GenesisBlock, err := blockchain.CreateGenesisBlock(GenesisConfPath)

	if err != nil {
		return err
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

	iLogger.Info(nil, fmt.Sprintf("[Blockchain] Genesis block has Committed - seal: [%x], height: [%d]", GenesisBlock.Seal, GenesisBlock.Height))

	return bApi.eventService.Publish("block.committed", commitEvent)
}

/**
set state to 'committed'
publish block committed event
*/
func (bApi BlockApi) CommitBlock(block blockchain.DefaultBlock) error {
	iLogger.Info(nil, "[Blockchain] Committing proposed block")

	// save(commit)
	block.SetState(blockchain.Committed)

	err := bApi.blockRepository.Save(block)

	if err != nil {
		return ErrSaveBlock
	}

	// publish
	commitEvent, err := createBlockCommittedEvent(block)

	if err != nil {
		return ErrCreateEvent
	}

	iLogger.Info(nil, fmt.Sprintf("[Blockchain] Proposed block has Committed - seal: [%x],  height: [%d]", block.Seal, block.Height))

	return bApi.eventService.Publish("block.committed", commitEvent)
}

func (bApi BlockApi) StageBlock(block blockchain.DefaultBlock) error {
	block.SetState(blockchain.Staged)
	bApi.BlockPool.Add(block)

	return nil
}

func (api BlockApi) CreateProposedBlock(txList []*blockchain.DefaultTransaction) (blockchain.DefaultBlock, error) {
	iLogger.Info(nil, "[Blockchain] Create proposed block")
	lastBlock, err := api.blockRepository.FindLast()
	if err != nil {
		return blockchain.DefaultBlock{}, ErrGetLastBlock
	}

	prevSeal := lastBlock.GetSeal()
	height := lastBlock.GetHeight() + 1
	creator := api.publisherId

	block, err := blockchain.CreateProposedBlock(prevSeal, height, txList, creator)
	if err != nil {
		return blockchain.DefaultBlock{}, err
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
		State:     blockchain.Committed,
	}, nil
}

func createBlockCreatedEvent(block blockchain.DefaultBlock) (event.BlockCreated, error) {

	txList := blockchain.ConvBackFromTransactionList(block.TxList)

	return event.BlockCreated{
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
