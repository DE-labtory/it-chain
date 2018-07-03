package blockchain

import (
	"errors"
	"fmt"
	"github.com/it-chain/it-chain-Engine/core/eventstore"
	"github.com/it-chain/midgard"
)

type BlockPool interface {
	Add(block Block)
	Get(height BlockHeight) Block
	Delete(height Block)
}


// block queued events Aggregate id
// BlockQueuedEvent들을 모아놓은 aggregate id 이다.
// struct는 존재하지 않는다.
var BLOCK_QUEUED_EVENTS_AID = "BLOCK_QUEUED_EVENTS_AID"


var BLOCK_POOL_AID = "BLOCK_POOL_AID"

type BlockPoolModel struct {
	midgard.AggregateModel
	pool map[BlockHeight] Block
}

func NewBlockPool() *BlockPoolModel {
	return &BlockPoolModel{
		AggregateModel: midgard.AggregateModel{
			ID: BLOCK_POOL_AID,
		},
		pool: make(map[BlockHeight] Block),
	}
}

func (p *BlockPoolModel) Add(block Block) {
	height := block.GetHeight()
	p.pool[height] = block

	addEvent := createBlockAddToPoolEvent(block)
	eventstore.Save(BLOCK_POOL_AID, addEvent)

	qEvent := createBlockQueuedEvent(block)
	eventstore.Save(BLOCK_QUEUED_EVENTS_AID, qEvent)
}

func (p *BlockPoolModel) Get(height BlockHeight) Block {
	return p.pool[height]
}

func (p *BlockPoolModel) Delete(block Block) {
	delete(p.pool, block.GetHeight())

	event := createBlockRemoveFromPoolEvent(block)
	eventstore.Save(BLOCK_POOL_AID, event)
}

func createBlockQueuedEvent(block Block) BlockQueuedEvent {
	return BlockQueuedEvent{
		EventModel: midgard.EventModel{
			ID: BLOCK_QUEUED_EVENTS_AID,
		},
		Block: block,
	}
}

func createBlockAddToPoolEvent(block Block) BlockAddToPoolEvent {
	return BlockAddToPoolEvent{
		EventModel: midgard.EventModel{
			ID: BLOCK_POOL_AID,
		},
		Block: block,
	}
}

func createBlockRemoveFromPoolEvent(block Block) BlockRemoveFromPoolEvent {
	return BlockRemoveFromPoolEvent{
		EventModel: midgard.EventModel{
			ID: BLOCK_POOL_AID,
		},
		Block: block,
	}
}

func (p *BlockPoolModel) GetID() string {
	return BLOCK_POOL_AID
}

func (p *BlockPoolModel) On(event midgard.Event) error {
	switch v := event.(type) {

	case *BlockAddToPoolEvent:
		block := event.(BlockAddToPoolEvent).Block
		p.pool[block.GetHeight()] = block

	case *BlockRemoveFromPoolEvent:
		block := event.(BlockRemoveFromPoolEvent).Block
		delete(p.pool, block.GetHeight())

	default:
		return errors.New(fmt.Sprintf("unhandled event [%s]", v))
	}
	return nil
}


// BlockSyncState Aggregate ID
var BC_SYNC_STATE_AID = "BC_SYNC_STATE_AID"

type ProgressState bool

const (
	PROGRESSING ProgressState = true
	DONE ProgressState = false
)

// 현재 블록 동기화가 진행 중인지 정보를 가진다.
type BlockSyncState struct {
	midgard.AggregateModel
	isProgress ProgressState
}

func NewBlockSyncState() *BlockSyncState {
	return &BlockSyncState{
		AggregateModel: midgard.AggregateModel{
			ID: BC_SYNC_STATE_AID,
		},
		isProgress:DONE,
	}
}

func (state *BlockSyncState) GetID() string {
	return BC_SYNC_STATE_AID
}

func (state *BlockSyncState) IsProgressing() ProgressState {
	return state.isProgress
}

func (state *BlockSyncState) On(event midgard.Event) error {
	switch v := event.(type) {

	case *SyncStartEvent:
		state.isProgress = PROGRESSING

	case *SyncDoneEvent:
		state.isProgress = DONE

	default:
		return errors.New(fmt.Sprintf("unhandled event [%s]", v))
	}

	return nil
}


