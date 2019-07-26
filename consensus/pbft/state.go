/*
 * Copyright 2018 DE-labtory
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * https://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package pbft

import (
	"encoding/json"
	"errors"

	"github.com/rs/xid"
)

type Stage string

const (
	IDLE_STAGE      Stage = "IdleStage"
	PROPOSE_STAGE   Stage = "ProposeStage"
	PREVOTE_STAGE   Stage = "PrevoteStage"
	PRECOMMIT_STAGE Stage = "PreCommitStage"
)

var ErrDecodingEmptyBlock = errors.New("Empty Block decoding failed")
var ErrPrevoteMsgNil = errors.New("ReceivePrevote msg is nil")
var ErrBlockHashNil = errors.New("Block hash is nil")
var ErrPreCommitMsgNil = errors.New("ReceivePreCommit msg is nil")
var ErrStateIdNotSame = errors.New("State ID is not same")

type ProposedBlock struct {
	Seal []byte
	Body []byte
}

func (block *ProposedBlock) Serialize() ([]byte, error) {
	data, err := json.Marshal(block)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (block *ProposedBlock) Deserialize(serializedBlock []byte) error {
	if len(serializedBlock) == 0 {
		return ErrDecodingEmptyBlock
	}

	err := json.Unmarshal(serializedBlock, block)
	if err != nil {
		return err
	}

	return nil
}

type MemberID string

func (m MemberID) ToString() string {
	return string(m)
}

type ProposeMsg struct {
	StateID        StateID
	SenderID       string
	Representative []Representative
	ProposedBlock  ProposedBlock
}

func NewProposeMsg(s *State, senderID string) *ProposeMsg {
	return &ProposeMsg{
		StateID:        s.StateID,
		SenderID:       senderID,
		Representative: s.Representatives,
		ProposedBlock:  s.Block,
	}
}

func (pp ProposeMsg) ToByte() ([]byte, error) {
	data, err := json.Marshal(pp)
	if err != nil {
		return nil, err
	}
	return data, nil
}

type PrevoteMsg struct {
	MsgID     string
	StateID   StateID
	SenderID  string
	BlockHash []byte
}

func NewPrevoteMsg(s *State, senderID string) *PrevoteMsg {
	return &PrevoteMsg{
		MsgID:     xid.New().String(),
		StateID:   s.StateID,
		SenderID:  senderID,
		BlockHash: s.Block.Seal,
	}
}

func (p PrevoteMsg) ToByte() ([]byte, error) {
	data, err := json.Marshal(p)
	if err != nil {
		return nil, err
	}
	return data, nil
}

type PreCommitMsg struct {
	MsgID    string
	StateID  StateID
	SenderID string
}

func NewPreCommitMsg(s *State, senderID string) *PreCommitMsg {
	return &PreCommitMsg{
		MsgID:    xid.New().String(),
		StateID:  s.StateID,
		SenderID: senderID,
	}
}

func (c PreCommitMsg) ToByte() ([]byte, error) {
	data, err := json.Marshal(c)
	if err != nil {
		return nil, err
	}
	return data, nil
}

type PrevoteMsgPool struct {
	messages map[string]PrevoteMsg
}

func NewPrevoteMsgPool() PrevoteMsgPool {
	return PrevoteMsgPool{
		messages: make(map[string]PrevoteMsg),
	}
}

func (p *PrevoteMsgPool) Save(prevoteMsg *PrevoteMsg) error {
	if prevoteMsg == nil {
		return ErrPrevoteMsgNil
	}

	senderID := prevoteMsg.SenderID
	if p.checkSenderExisting(senderID) {
		return errors.New("Already exist member [" + senderID + "]")
	}

	blockHash := prevoteMsg.BlockHash

	if blockHash == nil {
		return ErrBlockHashNil
	}

	id := prevoteMsg.MsgID
	p.messages[id] = *prevoteMsg
	return nil
}

func (p *PrevoteMsgPool) RemoveAllMsgs() {
	p.messages = make(map[string]PrevoteMsg)
}

func (p *PrevoteMsgPool) Remove(msgID string) {
	delete(p.messages, msgID)
}

func (p *PrevoteMsgPool) FindAll() []PrevoteMsg {
	mList := make([]PrevoteMsg, 0)

	for _, msg := range p.messages {
		mList = append(mList, msg)
	}

	return mList
}

func (p *PrevoteMsgPool) FindById(msgID string) PrevoteMsg {
	return p.messages[msgID]
}

func (p *PrevoteMsgPool) checkSenderExisting(senderID string) bool {
	for _, msg := range p.messages {
		if msg.SenderID == senderID {
			return true
		}
	}

	return false
}

type PreCommitMsgPool struct {
	messages map[string]PreCommitMsg
}

func NewPreCommitMsgPool() PreCommitMsgPool {
	return PreCommitMsgPool{
		messages: make(map[string]PreCommitMsg),
	}
}

func (c *PreCommitMsgPool) Save(precommitMsg *PreCommitMsg) error {
	if precommitMsg == nil {
		return ErrPreCommitMsgNil
	}

	senderID := precommitMsg.SenderID
	if c.checkSenderExisting(senderID) {
		return errors.New("Already exist member [" + senderID + "]")
	}

	id := precommitMsg.MsgID
	c.messages[id] = *precommitMsg
	return nil
}

func (c *PreCommitMsgPool) RemoveAllMsgs() {
	c.messages = make(map[string]PreCommitMsg)
}

func (c *PreCommitMsgPool) Remove(msgID string) {
	delete(c.messages, msgID)
}

func (c *PreCommitMsgPool) FindAll() []PreCommitMsg {
	mList := make([]PreCommitMsg, 0)

	for _, msg := range c.messages {
		mList = append(mList, msg)
	}

	return mList
}

func (c *PreCommitMsgPool) FindById(msgID string) PreCommitMsg {
	return c.messages[msgID]
}

func (c *PreCommitMsgPool) checkSenderExisting(senderID string) bool {
	for _, msg := range c.messages {
		if msg.SenderID == senderID {
			return true
		}
	}

	return false
}

type StateID struct {
	ID string
}

func NewStateID(id string) StateID {
	return StateID{
		ID: id,
	}
}

type State struct {
	StateID          StateID
	Representatives  []Representative
	Block            ProposedBlock
	CurrentStage     Stage
	PrevoteMsgPool   PrevoteMsgPool
	PreCommitMsgPool PreCommitMsgPool
}

func (s *State) GetID() string {
	return s.StateID.ID
}

func (s *State) GetReceipients(publisherID string) []Representative {
	receipients := make([]Representative, 0)

	for _, rep := range s.Representatives {
		if rep.ID != publisherID {
			receipients = append(receipients, rep)
		}
	}
	return receipients

}

func (s *State) Start() {
	s.CurrentStage = PROPOSE_STAGE
}

func (s *State) IsPrevoteStage() bool {

	if s.CurrentStage == PREVOTE_STAGE {
		return true
	}
	return false
}

func (s *State) IsPreCommitStage() bool {

	if s.CurrentStage == PRECOMMIT_STAGE {
		return true
	}
	return false
}

func (s *State) ToPrevoteStage() {
	s.CurrentStage = PREVOTE_STAGE
}

func (s *State) ToPreCommitStage() {
	s.CurrentStage = PRECOMMIT_STAGE
}

func (s *State) ToIdleStage() {
	s.CurrentStage = IDLE_STAGE
}

func (s *State) SavePrevoteMsg(prevoteMsg *PrevoteMsg) error {
	if s.StateID.ID != prevoteMsg.StateID.ID {
		return ErrStateIdNotSame
	}

	return s.PrevoteMsgPool.Save(prevoteMsg)
}

func (s *State) SavePreCommitMsg(precommitMsg *PreCommitMsg) error {
	if s.StateID.ID != precommitMsg.StateID.ID {
		return ErrStateIdNotSame
	}

	return s.PreCommitMsgPool.Save(precommitMsg)
}
func (s *State) CheckPrevoteCondition() bool {
	representativeNum := len(s.Representatives)
	prevoteMsgNum := len(s.PrevoteMsgPool.FindAll())
	satisfyNum := representativeNum / 3

	if prevoteMsgNum >= (satisfyNum + 1) {
		return true
	}
	return false
}
func (s *State) CheckPreCommitCondition() bool {
	representativeNum := len(s.Representatives)
	commitMsgNum := len(s.PreCommitMsgPool.FindAll())
	satisfyNum := representativeNum / 3

	if commitMsgNum >= (satisfyNum + 1) {
		return true
	}
	return false
}

type StateRepository interface {
	Save(state State) error
	Load() (State, error)
	Remove()
}
