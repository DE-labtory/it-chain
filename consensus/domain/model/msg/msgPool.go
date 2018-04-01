package msg

import (
	"sync"

	cs "github.com/it-chain/it-chain-Engine/consensus/domain/model/consensus"
)

type MsgPool struct {
	PrepareMsgPool map[cs.ConsensusID][]PrepareMsg
	CommitMsgPool  map[cs.ConsensusID][]CommitMsg
	lock           sync.RWMutex
}

func (mp MsgPool) InsertPrepareMsg(prepareMsg PrepareMsg) {

	mp.lock.Lock()
	defer mp.lock.Unlock()

	prepareMsgPool := mp.PrepareMsgPool[prepareMsg.ConsensusID]

	if prepareMsgPool == nil {
		mp.PrepareMsgPool[prepareMsg.ConsensusID] = make([]PrepareMsg, 0)
	}

	mp.PrepareMsgPool[prepareMsg.ConsensusID] = append(mp.PrepareMsgPool[prepareMsg.ConsensusID], prepareMsg)
}

func (mp MsgPool) FindPrepareMsgsByConsensusID(id cs.ConsensusID) []PrepareMsg {

	return mp.PrepareMsgPool[id]
}
