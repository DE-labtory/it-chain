package consensus

import "sync"

type PrepareMsgRepository interface {
	DeleteAllPrepareMsg(consensusId ConsensusId)
	InsertPrepareMsg(prepareMsg PrepareMsg)
	FindPrepareMsgsByConsensusID(consensusId ConsensusId) []PrepareMsg
}

type PrepareMsgRepositoryImpl struct {
	PrepareMsgPool map[ConsensusId][]PrepareMsg
	lock           *sync.RWMutex
}

type CommitMsgRepository interface {
	DeleteAllCommitMsg(consensusId ConsensusId)
	InsertCommitMsg(commitMsg CommitMsg)
	FindCommitMsgsByConsensusID(consensusId ConsensusId) []CommitMsg
}

type CommitMsgRepositoryImpl struct {
	CommitMsgPool map[ConsensusId][]CommitMsg
	lock          *sync.RWMutex
}

func NewPrepareMsgRepository() *PrepareMsgRepositoryImpl {
	return &PrepareMsgRepositoryImpl{
		PrepareMsgPool: make(map[ConsensusId][]PrepareMsg),
		lock:           &sync.RWMutex{},
	}
}

func NewCommitMsgRepository() *CommitMsgRepositoryImpl {
	return &CommitMsgRepositoryImpl{
		CommitMsgPool: make(map[ConsensusId][]CommitMsg),
		lock:          &sync.RWMutex{},
	}
}

func (pr *PrepareMsgRepositoryImpl) DeleteAllPrepareMsg(consensusId ConsensusId) {
	pr.lock.Lock()
	defer pr.lock.Unlock()

	delete(pr.PrepareMsgPool, consensusId)
}

func (cr *CommitMsgRepositoryImpl) DeleteAllCommitMsg(consensusId ConsensusId) {
	cr.lock.Lock()
	defer cr.lock.Unlock()

	delete(cr.CommitMsgPool, consensusId)
}

func (pr *PrepareMsgRepositoryImpl) InsertPrepareMsg(prepareMsg PrepareMsg) {
	pr.lock.Lock()
	defer pr.lock.Unlock()

	if prepareMsg.SenderId == "" {
		return
	}

	prepareMsgPool := pr.PrepareMsgPool[prepareMsg.ConsensusId]

	if prepareMsgPool == nil {
		pr.PrepareMsgPool[prepareMsg.ConsensusId] = make([]PrepareMsg, 0)
	}

	var hasItem = func(prepareMsgPool []PrepareMsg, prepareMsg PrepareMsg) bool {
		for _, msg := range prepareMsgPool {
			if msg.SenderId == prepareMsg.SenderId {
				return true
			}
		}
		return false
	}(prepareMsgPool, prepareMsg)

	if hasItem {
		return
	}

	pr.PrepareMsgPool[prepareMsg.ConsensusId] = append(pr.PrepareMsgPool[prepareMsg.ConsensusId], prepareMsg)
}

func (cr *CommitMsgRepositoryImpl) InsertCommitMsg(commitMsg CommitMsg) {
	cr.lock.Lock()
	defer cr.lock.Unlock()

	if commitMsg.SenderId == "" {
		return
	}

	CommitMsgPool := cr.CommitMsgPool[commitMsg.ConsensusId]

	if CommitMsgPool == nil {
		cr.CommitMsgPool[commitMsg.ConsensusId] = make([]CommitMsg, 0)
	}

	var hasItem = func(CommitMsgPool []CommitMsg, CommitMsg CommitMsg) bool {
		for _, msg := range CommitMsgPool {
			if msg.SenderId == CommitMsg.SenderId {
				return true
			}
		}
		return false
	}(CommitMsgPool, commitMsg)

	if hasItem {
		return
	}

	cr.CommitMsgPool[commitMsg.ConsensusId] = append(cr.CommitMsgPool[commitMsg.ConsensusId], commitMsg)
}

func (pr *PrepareMsgRepositoryImpl) FindPrepareMsgsByConsensusID(consensusId ConsensusId) []PrepareMsg {

	return pr.PrepareMsgPool[consensusId]
}

func (cr *CommitMsgRepositoryImpl) FindCommitMsgsByConsensusID(consensusId ConsensusId) []CommitMsg {

	return cr.CommitMsgPool[consensusId]
}