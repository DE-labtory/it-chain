//Msgpool manage voting messages from other peer
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

func (mp MsgPool) DeletePrepareMsg(id cs.ConsensusID) {

	mp.lock.Lock()
	defer mp.lock.Unlock()

	delete(mp.PrepareMsgPool, id)
}

func (mp MsgPool) DeleteCommitMsg(id cs.ConsensusID) {

	mp.lock.Lock()
	defer mp.lock.Unlock()

	delete(mp.CommitMsgPool, id)
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

func (mp MsgPool) InsertCommitMsg(commitMsg CommitMsg) {

	mp.lock.Lock()
	defer mp.lock.Unlock()

	commitMsgPool := mp.CommitMsgPool[commitMsg.ConsensusID]

	if commitMsgPool == nil {
		mp.CommitMsgPool[commitMsg.ConsensusID] = make([]CommitMsg, 0)
	}

	mp.CommitMsgPool[commitMsg.ConsensusID] = append(mp.CommitMsgPool[commitMsg.ConsensusID], commitMsg)
}

func (mp MsgPool) FindPrepareMsgsByConsensusID(id cs.ConsensusID) []PrepareMsg {

	return mp.PrepareMsgPool[id]
}

func (mp MsgPool) FindCommitMsgsByConsensusID(id cs.ConsensusID) []CommitMsg {

	return mp.CommitMsgPool[id]
}
