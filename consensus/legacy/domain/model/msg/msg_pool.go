//Msgpool manage voting messages from other peer
package msg

import (
	"fmt"
	"sync"

	cs "github.com/it-chain/it-chain-Engine/consensus/legacy/domain/model/consensus"
	"errors"
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

//ConsensusID로 PrepareMsg리스트를 찾는다.
//SenderID로 중복 투표를 확인한다.
func (mp MsgPool) InsertPrepareMsg(prepareMsg PrepareMsg) {

	mp.lock.Lock()
	defer mp.lock.Unlock()

	if prepareMsg.SenderID == "" {
		return
	}

	prepareMsgPool := mp.PrepareMsgPool[prepareMsg.ConsensusID]

	if prepareMsgPool == nil {
		mp.PrepareMsgPool[prepareMsg.ConsensusID] = make([]PrepareMsg, 0)
	}

	var hasItem = func(prepareMsgPool []PrepareMsg, prepareMsg PrepareMsg) bool {
		for _, msg := range prepareMsgPool {
			if msg.SenderID == prepareMsg.SenderID {
				return true
			}
		}
		return false
	}(prepareMsgPool, prepareMsg)

	if hasItem {
		return
	}

	mp.PrepareMsgPool[prepareMsg.ConsensusID] = append(mp.PrepareMsgPool[prepareMsg.ConsensusID], prepareMsg)
}

func (mp MsgPool) InsertCommitMsg(commitMsg CommitMsg) error {

	mp.lock.Lock()
	defer mp.lock.Unlock()

	if commitMsg.SenderID == "" {
		return errors.New(fmt.Sprint("Need SenderID [%s]", commitMsg.SenderID))
	}

	commitMsgPool := mp.CommitMsgPool[commitMsg.ConsensusID]

	if commitMsgPool == nil {
		mp.CommitMsgPool[commitMsg.ConsensusID] = make([]CommitMsg, 0)
	}

	var hasItem = func(prepareMsgPool []CommitMsg, commitMsg CommitMsg) bool {
		for _, msg := range prepareMsgPool {
			if msg.SenderID == commitMsg.SenderID {
				return true
			}
		}
		return false
	}(commitMsgPool, commitMsg)

	if hasItem {
		return errors.New(fmt.Sprint("Already has SenderID [%s]", commitMsg.SenderID))
	}

	mp.CommitMsgPool[commitMsg.ConsensusID] = append(mp.CommitMsgPool[commitMsg.ConsensusID], commitMsg)

	return nil
}

func (mp MsgPool) FindPrepareMsgsByConsensusID(id cs.ConsensusID) []PrepareMsg {

	return mp.PrepareMsgPool[id]
}

func (mp MsgPool) FindCommitMsgsByConsensusID(id cs.ConsensusID) []CommitMsg {

	return mp.CommitMsgPool[id]
}
