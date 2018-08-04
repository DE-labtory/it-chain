/*
 * Copyright 2018 It-chain
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
	"errors"
	"fmt"

	"encoding/json"
)

type Stage string

const (
	IDLE_STAGE       Stage = "IdleStage"
	PREPREPARE_STAGE Stage = "PrePrepareStage"
	PREPARE_STAGE    Stage = "PrepareStage"
	COMMIT_STAGE     Stage = "CommitStage"
)

var ErrDecodingEmptyBlock = errors.New("Empty Block decoding failed")
var ErrPrepareMsgNil = errors.New("Prepare msg is nil")
var ErrBlockHashNil = errors.New("Block hash is nil")
var ErrCommitMsgNil = errors.New("Commit msg is nil")
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

type RepresentativeID string

type Representative struct {
	ID RepresentativeID
}

func (r Representative) GetID() string {
	return string(r.ID)
}

func NewRepresentative(ID string) *Representative {
	return &Representative{ID: RepresentativeID(ID)}
}

type PrePrepareMsg struct {
	StateID        StateID
	SenderID       string
	Representative []*Representative
	ProposedBlock  ProposedBlock
}

func NewPrePrepareMsg(s *State) *PrePrepareMsg {
	return &PrePrepareMsg{
		StateID:        s.StateID,
		Representative: s.Representatives,
		ProposedBlock:  s.Block,
	}
}

func (pp PrePrepareMsg) ToByte() ([]byte, error) {
	data, err := json.Marshal(pp)
	if err != nil {
		return nil, err
	}
	return data, nil
}

type PrepareMsg struct {
	StateID   StateID
	SenderID  string
	BlockHash []byte
}

func NewPrepareMsg(s *State) *PrepareMsg {
	return &PrepareMsg{
		StateID:   s.StateID,
		BlockHash: s.Block.Seal,
	}
}

func (p PrepareMsg) ToByte() ([]byte, error) {
	data, err := json.Marshal(p)
	if err != nil {
		return nil, err
	}
	return data, nil
}

type CommitMsg struct {
	StateID  StateID
	SenderID string
}

func NewCommitMsg(s *State) *CommitMsg {
	return &CommitMsg{
		StateID: s.StateID,
	}
}

func (c CommitMsg) ToByte() ([]byte, error) {
	data, err := json.Marshal(c)
	if err != nil {
		return nil, err
	}
	return data, nil
}

type PrepareMsgPool struct {
	messages []PrepareMsg
}

func NewPrepareMsgPool() PrepareMsgPool {
	return PrepareMsgPool{
		messages: make([]PrepareMsg, 0),
	}
}

func (p *PrepareMsgPool) Save(prepareMsg *PrepareMsg) error {
	if prepareMsg == nil {
		return ErrPrepareMsgNil
	}

	senderID := prepareMsg.SenderID
	index := p.findIndexOfPrepareMsg(senderID)

	if index != -1 {
		return errors.New(fmt.Sprintf("Already exist member [%s]", senderID))
	}

	blockHash := prepareMsg.BlockHash

	if blockHash == nil {
		return ErrBlockHashNil
	}

	p.messages = append(p.messages, *prepareMsg)

	return nil
}

func (p *PrepareMsgPool) RemoveAllMsgs() {
	p.messages = make([]PrepareMsg, 0)
}

func (p *PrepareMsgPool) Get() []PrepareMsg {
	return p.messages
}

func (p *PrepareMsgPool) findIndexOfPrepareMsg(senderID string) int {
	for i, msg := range p.messages {
		if msg.SenderID == senderID {
			return i
		}
	}

	return -1
}

type CommitMsgPool struct {
	messages []CommitMsg
}

func NewCommitMsgPool() CommitMsgPool {
	return CommitMsgPool{
		messages: make([]CommitMsg, 0),
	}
}

func (c *CommitMsgPool) Save(commitMsg *CommitMsg) error {
	if commitMsg == nil {
		return ErrCommitMsgNil
	}

	senderID := commitMsg.SenderID
	index := c.findIndexOfCommitMsg(senderID)

	if index != -1 {
		return errors.New(fmt.Sprintf("Already exist member [%s]", senderID))
	}

	c.messages = append(c.messages, *commitMsg)

	return nil
}

func (c *CommitMsgPool) RemoveAllMsgs() {
	c.messages = make([]CommitMsg, 0)
}

func (c *CommitMsgPool) Get() []CommitMsg {
	return c.messages
}

func (c *CommitMsgPool) findIndexOfCommitMsg(senderID string) int {
	for i, msg := range c.messages {
		if msg.SenderID == senderID {
			return i
		}
	}

	return -1
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
	StateID         StateID
	Representatives []*Representative
	Block           ProposedBlock
	CurrentStage    Stage
	PrepareMsgPool  PrepareMsgPool
	CommitMsgPool   CommitMsgPool
}

func (s *State) GetID() string {
	return s.StateID.ID
}

func (s *State) Start() {
	s.CurrentStage = PREPREPARE_STAGE
}

func (s *State) IsPrepareStage() bool {

	if s.CurrentStage == PREPARE_STAGE {
		return true
	}
	return false
}

func (s *State) IsCommitStage() bool {

	if s.CurrentStage == COMMIT_STAGE {
		return true
	}
	return false
}

func (s *State) ToPrepareStage() {
	s.CurrentStage = PREPARE_STAGE
}

func (s *State) ToCommitStage() {
	s.CurrentStage = COMMIT_STAGE
}

func (s *State) ToIdleStage() {
	s.CurrentStage = IDLE_STAGE
}

func (s *State) SavePrepareMsg(prepareMsg *PrepareMsg) error {
	if s.StateID.ID != prepareMsg.StateID.ID {
		return ErrStateIdNotSame
	}

	return s.PrepareMsgPool.Save(prepareMsg)
}

func (s *State) SaveCommitMsg(commitMsg *CommitMsg) error {
	if s.StateID.ID != commitMsg.StateID.ID {
		return ErrStateIdNotSame
	}

	return s.CommitMsgPool.Save(commitMsg)
}
func (s *State) CheckPrepareCondition() bool {
	representativeNum := len(s.Representatives)
	commitMsgNum := len(s.PrepareMsgPool.Get())
	satisfyNum := representativeNum / 3

	if commitMsgNum > (satisfyNum + 1) {
		return true
	}
	return false
}
func (s *State) CheckCommitCondition() bool {
	representativeNum := len(s.Representatives)
	commitMsgNum := len(s.CommitMsgPool.Get())
	satisfyNum := representativeNum / 3

	if commitMsgNum > (satisfyNum + 1) {
		return true
	}
	return false
}
