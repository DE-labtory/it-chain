package blockchain

import (
	"errors"
	"fmt"

	"github.com/it-chain/it-chain-Engine/core/eventstore"
	"github.com/it-chain/midgard"
)

// BlockSyncState Aggregate ID
var BC_SYNC_STATE_AID = "BC_SYNC_STATE_AID"

type ProgressState bool

const (
	PROGRESSING ProgressState = true
	DONE        ProgressState = false
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
		isProgress: DONE,
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
	eventstore.Save(BC_SYNC_STATE_AID, event)
	bss.On(event)
}

func createSyncStartEvent() *SyncStartEvent {
	return &SyncStartEvent{
		EventModel: midgard.EventModel{
			ID:   BC_SYNC_STATE_AID,
			Type: "sync.started",
		},
	}
}

func createSyncDoneEvent() *SyncDoneEvent {
	return &SyncDoneEvent{
		EventModel: midgard.EventModel{
			ID:   BC_SYNC_STATE_AID,
			Type: "sync.done",
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
