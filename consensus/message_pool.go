package consensus

import (
	"sync"
)

type MsgPoolId string

func (mpId MsgPoolId) ToString() string {
	return string(mpId)
}

type MsgPool struct {
	MsgPoolId      MsgPoolId
	PrepareMsgPool map[ConsensusId][]PrepareMsg
	CommitMsgPool  map[ConsensusId][]CommitMsg
	lock           sync.RWMutex
}

func (mp *MsgPool) GetID() string {
	return mp.MsgPoolId.ToString()
}

