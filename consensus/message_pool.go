package consensus

import (
	"errors"
	"fmt"
)

type PrepareMsgPool struct {
	messages []PrepareMsg
}

func (pmp *PrepareMsgPool) Save(prepareMsg *PrepareMsg) error {
	if prepareMsg == nil {
		return errors.New("Prepare msg is nil")
	}

	senderID := prepareMsg.SenderId
	index := pmp.findIndexOfPrepareMsg(senderID)

	if index != -1 {
		return errors.New(fmt.Sprintf("Already exist member [%s]", senderID))
	}

	proposedBlock := prepareMsg.ProposedBlock

	if proposedBlock == nil {
		return errors.New("Proposed block is nil")
	}

	pmp.messages = append(pmp.messages, *prepareMsg)

	return nil
}

func (pmp *PrepareMsgPool) RemoveAllMsgs() {
	pmp.messages = make([]PrepareMsg, 0)
}

func (pmp *PrepareMsgPool) Get() []PrepareMsg {
	return pmp.messages
}

func (pmp *PrepareMsgPool) findIndexOfPrepareMsg(senderID string) int {
	for i, msg := range pmp.messages {
		if msg.SenderId == senderID {
			return i
		}
	}

	return -1
}

type CommitMsgPool struct {
	messages []CommitMsg
}

func (cmp *CommitMsgPool) Save(commitMsg *CommitMsg) error {
	if commitMsg == nil {
		return errors.New("Commit msg is nil")
	}

	senderID := commitMsg.SenderId
	index := cmp.findIndexOfCommitMsg(senderID)

	if index != -1 {
		return errors.New(fmt.Sprintf("Already exist member [%s]", senderID))
	}

	cmp.messages = append(cmp.messages, *commitMsg)

	return nil
}

func (cmp *CommitMsgPool) RemoveAllMsgs() {
	cmp.messages = make([]CommitMsg, 0)
}

func (cmp *CommitMsgPool) Get() []CommitMsg {
	return cmp.messages
}

func (cmp *CommitMsgPool) findIndexOfCommitMsg(senderID string) int {
	for i, msg := range cmp.messages {
		if msg.SenderId == senderID {
			return i
		}
	}

	return -1
}
