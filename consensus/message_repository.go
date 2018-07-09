package consensus

import (
	"sync"
)

type MsgPoolId string

func (mpId MsgPoolId) ToString() string {
	return string(mpId)
}

type PrepareMsgRepository interface {
	Save(prepareMsg PrepareMsg)
	Remove(id ConsensusId)
	FindPrepareMsgsByConsensusID(id ConsensusId) []PrepareMsg
}

type PrepareMsgRepositoryImpl struct {
	PreparePool map[ConsensusId][]PrepareMsg
	lock        *sync.RWMutex
}

func NewPrepareMsgRepository() PrepareMsgRepository {
	return &PrepareMsgRepositoryImpl{
		PreparePool: make(map[ConsensusId][]PrepareMsg, 0),
		lock:        &sync.RWMutex{},
	}
}

func (pr *PrepareMsgRepositoryImpl) Save(prepareMsg PrepareMsg) {
	pr.lock.Lock()
	defer pr.lock.Unlock()

	msgPool := pr.PreparePool[prepareMsg.ConsensusId]

	pr.PreparePool[prepareMsg.ConsensusId] = append(msgPool, prepareMsg)
}

func (pr *PrepareMsgRepositoryImpl) Remove(id ConsensusId) {
	pr.lock.Lock()
	defer pr.lock.Unlock()

	delete(pr.PreparePool, id)
}

func (pr *PrepareMsgRepositoryImpl) FindPrepareMsgsByConsensusID(id ConsensusId) []PrepareMsg {
	pr.lock.Lock()
	defer pr.lock.Unlock()

	return pr.PreparePool[id]
}

type CommitMsgRepository interface {
	Save(commitMsg CommitMsg)
	Remove(id ConsensusId)
	FindCommitMsgsByConsensusID(id ConsensusId) []CommitMsg
}

type CommitMsgRepositoryImpl struct {
	CommitPool map[ConsensusId][]CommitMsg
	lock       *sync.RWMutex
}

func NewCommitMsgRepository() CommitMsgRepository {
	return &CommitMsgRepositoryImpl{
		CommitPool: make(map[ConsensusId][]CommitMsg, 0),
		lock:       &sync.RWMutex{},
	}
}

func (cr *CommitMsgRepositoryImpl) Save(commitMsg CommitMsg) {
	cr.lock.Lock()
	defer cr.lock.Unlock()

	msgPool := cr.CommitPool[commitMsg.ConsensusId]

	cr.CommitPool[commitMsg.ConsensusId] = append(msgPool, commitMsg)
}

func (cr *CommitMsgRepositoryImpl) Remove(id ConsensusId) {
	cr.lock.Lock()
	defer cr.lock.Unlock()

	delete(cr.CommitPool, id)
}

func (cr *CommitMsgRepositoryImpl) FindCommitMsgsByConsensusID(id ConsensusId) []CommitMsg {
	cr.lock.Lock()
	defer cr.lock.Unlock()

	return cr.CommitPool[id]
}
