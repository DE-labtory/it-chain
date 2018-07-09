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
	msgPool := pr.PreparePool[prepareMsg.ConsensusId]

	pr.PreparePool[prepareMsg.ConsensusId] = append(msgPool, prepareMsg)
}

func (pr *PrepareMsgRepositoryImpl) Remove(id ConsensusId) {
	delete(pr.PreparePool, id)
}

func (pr *PrepareMsgRepositoryImpl) FindPrepareMsgsByConsensusID(id ConsensusId) []PrepareMsg {
	return pr.PreparePool[id]
}

type CommitMsgRepository interface {
	Save(commitMsg CommitMsg)
	Remove(id ConsensusId)
	FindCommitMsgsByConsensusID(id ConsensusId) []CommitMsg
}

type CommitMsgRepositoryImpl struct {
	CommitPool map[ConsensusId]*[]CommitMsg
	lock       *sync.RWMutex
}

func NewCommitMsgRepository() CommitMsgRepository {
	return &CommitMsgRepositoryImpl{
		CommitPool: make(map[ConsensusId]*[]CommitMsg, 0),
		lock:       &sync.RWMutex{},
	}
}

func (cr *CommitMsgRepositoryImpl) Save(commitMsg CommitMsg) {
	msgPool := *cr.CommitPool[commitMsg.ConsensusId]

	if msgPool == nil {
		msgPool = make([]CommitMsg, 0)
	}

	msgPool = append(msgPool, commitMsg)
}

func (cr *CommitMsgRepositoryImpl) Remove(id ConsensusId) {
	*cr.CommitPool[id] = nil
}

func (cr *CommitMsgRepositoryImpl) FindCommitMsgsByConsensusID(id ConsensusId) []CommitMsg {
	return *cr.CommitPool[id]
}
