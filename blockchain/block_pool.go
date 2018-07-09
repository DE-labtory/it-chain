package blockchain

import (
	"errors"
	"fmt"
	"github.com/it-chain/it-chain-Engine/core/eventstore"
	"github.com/it-chain/midgard"
	"github.com/it-chain/it-chain-Engine/common"
)

type BlockPool interface {
	Add(block Block) error
	Get(height BlockHeight) Block
	Delete(height Block)
}


var BLOCK_POOL_AID = "BLOCK_POOL_AID"

type BlockPoolModel struct {
	midgard.AggregateModel
	Pool map[BlockHeight] Block
}

func NewBlockPool() *BlockPoolModel {
	return &BlockPoolModel{
		AggregateModel: midgard.AggregateModel{
			ID: BLOCK_POOL_AID,
		},
		Pool: make(map[BlockHeight] Block),
	}
}

func (p *BlockPoolModel) Add(block Block) error {
	event, err := createBlockAddToPoolEvent(block)
	if err != nil {
		return err
	}

	eventstore.Save(BLOCK_POOL_AID, event)

	p.On(&event)

	return nil
}

func (p *BlockPoolModel) Get(height BlockHeight) Block {
	return p.Pool[height]
}

func (p *BlockPoolModel) Delete(block Block) {
	event := createBlockRemoveFromPoolEvent(block)
	eventstore.Save(BLOCK_POOL_AID, event)

	p.On(&event)
}

func createBlockAddToPoolEvent(block Block) (BlockAddToPoolEvent, error) {
	txListBytes, err := common.Serialize(block.GetTxList())
	if err != nil {
		return BlockAddToPoolEvent{}, ErrTxListMarshal
	}

	return BlockAddToPoolEvent{
		EventModel: midgard.EventModel{
			ID: BLOCK_POOL_AID,
		},
		Seal: block.GetSeal(),
		PrevSeal: block.GetPrevSeal(),
		Height: block.GetHeight(),
		TxList: txListBytes,
		TxSeal: block.GetTxSeal(),
		Timestamp: block.GetTimestamp(),
		Creator: block.GetCreator(),
	}, nil
}

var ErrTxListMarshal = errors.New("tx list marshal failed")

func createBlockRemoveFromPoolEvent(block Block) BlockRemoveFromPoolEvent {
	return BlockRemoveFromPoolEvent{
		EventModel: midgard.EventModel{
			ID: BLOCK_POOL_AID,
		},
		Height: block.GetHeight(),
	}
}

func (p *BlockPoolModel) GetID() string {
	return BLOCK_POOL_AID
}

func (p *BlockPoolModel) On(event midgard.Event) error {
	switch v := event.(type) {

	case *BlockAddToPoolEvent:
		block, err := createBlockFromAddToPoolEvent(v)
		if err != nil {
			return err
		}
		(p.Pool)[v.Height] = block

	case *BlockRemoveFromPoolEvent:
		delete(p.Pool, v.Height)

	default:
		return errors.New(fmt.Sprintf("unhandled event [%s]", v))
	}
	return nil
}

func createBlockFromAddToPoolEvent(event *BlockAddToPoolEvent) (Block, error) {
	var txList []Transaction
	err := common.Deserialize(event.TxList, &txList)
	if err != nil {
		return &DefaultBlock{}, ErrTxListUnmarshal
	}
	return &DefaultBlock{
		Seal: event.Seal,
		PrevSeal: event.PrevSeal,
		Height: event.Height,
		TxList: txList,
		TxSeal: event.TxSeal,
		Timestamp: event.Timestamp,
		Creator: event.Creator,
	}, nil
}

var ErrTxListUnmarshal = errors.New("tx list unmarshal failed")


// BlockSyncState Aggregate ID
var BC_SYNC_STATE_AID = "BC_SYNC_STATE_AID"

type ProgressState bool

const (
	PROGRESSING ProgressState = true
	DONE ProgressState = false
)


type SyncState interface {
	SetProgress(state ProgressState)
}

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

func (bss *BlockSyncState) GetID() string {
	return BC_SYNC_STATE_AID
}

func (bss *BlockSyncState) SetProgress(state ProgressState) {
	var event midgard.Event
	if state == PROGRESSING {
		event = createSyncStartEvent()
	} else { // state == DONE
		event = createSyncDoneEvent()
	}
	bss.isProgress = state
	eventstore.Save(BC_SYNC_STATE_AID, event)
	bss.On(event)
}

func createSyncStartEvent() SyncStartEvent {
	return SyncStartEvent{
		EventModel: midgard.EventModel{
			ID: BC_SYNC_STATE_AID,
		},
	}
}

func createSyncDoneEvent() SyncDoneEvent {
	return SyncDoneEvent{
		EventModel: midgard.EventModel{
			ID: BC_SYNC_STATE_AID,
		},
	}
}

func (bss *BlockSyncState) IsProgressing() ProgressState {
	return bss.isProgress
}

func (bss *BlockSyncState) On(event midgard.Event) error {
	switch v := event.(type) {

	case *SyncStartEvent:
		bss.isProgress = PROGRESSING

	case *SyncDoneEvent:
		bss.isProgress = DONE

	default:
		return errors.New(fmt.Sprintf("unhandled event [%s]", v))
	}

	return nil
}


