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
	"github.com/it-chain/engine/core/eventstore"
	"github.com/it-chain/midgard"
)

type BlockQueryService interface {
	BlockQueryInnerService
}

type BlockQueryInnerService interface {
	GetStagedBlockByHeight(height BlockHeight)(Block, error)
	GetStagedBlockById(blockId string) (Block, error)
	GetLastCommitedBlock() (Block, error)
	GetCommitedBlockByHeight(height BlockHeight) (Block, error)
}

func CommitBlock(block Block) error {

	event, err := createBlockCommittedEvent(block)

	if err != nil {
		return err
	}

	blockId := string(block.GetSeal())
	eventstore.Save(blockId, event)

	return nil
}

func createBlockCommittedEvent(block Block) (*BlockCommittedEvent, error) {

	aggregateId := string(block.GetSeal())

	return &BlockCommittedEvent{
		EventModel: midgard.EventModel{
			ID: aggregateId,
		},
		State: Committed,
	}, nil
}
