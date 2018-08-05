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

package blockchain

import (
	"encoding/hex"

	"github.com/it-chain/engine/common/event"
	"github.com/it-chain/engine/core/eventstore"
	"github.com/it-chain/midgard"
)

type BlockQueryService interface {
	BlockQueryInnerService
}

type EventService interface {
	Publish(topic string, event interface{}) error
}

type BlockQueryInnerService interface {
	GetStagedBlockByHeight(height BlockHeight) (DefaultBlock, error)
	GetStagedBlockById(blockId string) (DefaultBlock, error)
	GetLastCommitedBlock() (DefaultBlock, error)
	GetCommitedBlockByHeight(height BlockHeight) (DefaultBlock, error)
}

func CommitBlock(block Block) error {

	committedEvent, err := createBlockCommittedEvent(block)

	if err != nil {
		return err
	}

	blockId := string(block.GetSeal())

	return eventstore.Save(blockId, committedEvent)
}

func createBlockCommittedEvent(block Block) (*event.BlockCommitted, error) {

	AggregateID := hex.EncodeToString(block.GetSeal())

	return &event.BlockCommitted{
		EventModel: midgard.EventModel{
			ID:   AggregateID,
			Type: "block.committed",
		},
		State: Committed,
	}, nil
}


func CreateBlockCommittedEvent(block DefaultBlock) (event.BlockCommitted, error) {

	txList := ConvBackFromTransactionList(block.TxList);

	return event.BlockCommitted{
		Seal: block.GetSeal(),
		PrevSeal: block.GetPrevSeal(),
		Height: block.GetHeight(),
		TxList: txList,
		TxSeal: block.GetTxSeal(),
		Timestamp: block.GetTimestamp(),
		Creator: block.GetCreator(),
		State: block.State,
	}, nil
}

