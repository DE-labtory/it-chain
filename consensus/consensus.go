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

package consensus

import (
	"errors"
	"fmt"

	"encoding/json"

	"github.com/it-chain/engine/common/event"
	"github.com/it-chain/engine/core/eventstore"
	"github.com/it-chain/midgard"
)

type State string

const (
	IDLE_STATE       State = "IdleState"
	PREPREPARE_STATE State = "PrePrepareState"
	PREPARE_STATE    State = "PrepareState"
	COMMIT_STATE     State = "CommitState"
)

var ErrDecodingEmptyBlock = errors.New("Empty Block decoding failed")

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

type MemberId string

func (m MemberId) ToString() string {
	return string(m)
}

type RepresentativeId string

type Representative struct {
	Id RepresentativeId
}

func (r Representative) GetID() string {
	return string(r.Id)
}

func NewRepresentative(Id string) *Representative {
	return &Representative{Id: RepresentativeId(Id)}
}

type PrePrepareMsg struct {
	ConsensusId    ConsensusId
	SenderId       string
	Representative []*Representative
	ProposedBlock  ProposedBlock
}

func NewPrePrepareMsg(c *Consensus) *PrePrepareMsg {
	return &PrePrepareMsg{
		ConsensusId:    c.ConsensusID,
		Representative: c.Representatives,
		ProposedBlock:  c.Block,
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
	ConsensusId ConsensusId
	SenderId    string
	BlockHash   []byte
}

func NewPrepareMsg(c *Consensus) *PrepareMsg {
	return &PrepareMsg{
		ConsensusId: c.ConsensusID,
		BlockHash:   c.Block.Seal,
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
	ConsensusId ConsensusId
	SenderId    string
}

func NewCommitMsg(c *Consensus) *CommitMsg {
	return &CommitMsg{
		ConsensusId: c.ConsensusID,
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

func (pmp *PrepareMsgPool) Save(prepareMsg *PrepareMsg) error {
	if prepareMsg == nil {
		return errors.New("Prepare msg is nil")
	}

	senderID := prepareMsg.SenderId
	index := pmp.findIndexOfPrepareMsg(senderID)

	if index != -1 {
		return errors.New(fmt.Sprintf("Already exist member [%s]", senderID))
	}

	blockHash := prepareMsg.BlockHash

	if blockHash == nil {
		return errors.New("Block hash is nil")
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

func NewCommitMsgPool() CommitMsgPool {
	return CommitMsgPool{
		messages: make([]CommitMsg, 0),
	}
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

type ConsensusId struct {
	Id string
}

func NewConsensusId(id string) ConsensusId {
	return ConsensusId{
		Id: id,
	}
}

type Consensus struct {
	ConsensusID     ConsensusId
	Representatives []*Representative
	Block           ProposedBlock
	CurrentState    State
	PrepareMsgPool  PrepareMsgPool
	CommitMsgPool   CommitMsgPool
}

func (c *Consensus) GetID() string {
	return c.ConsensusID.Id
}

func (c *Consensus) Start() {
	c.CurrentState = PREPREPARE_STATE
}

func (c *Consensus) IsPrepareState() bool {

	if c.CurrentState == PREPARE_STATE {
		return true
	}
	return false
}

func (c *Consensus) IsCommitState() bool {

	if c.CurrentState == COMMIT_STATE {
		return true
	}
	return false
}

func (c *Consensus) ToPrepareState() {
	c.CurrentState = PREPARE_STATE
}

func (c *Consensus) ToCommitState() {
	c.CurrentState = COMMIT_STATE
}

func (c *Consensus) ToIdleState() {
	c.CurrentState = IDLE_STATE
}

func (c *Consensus) SavePrepareMsg(prepareMsg *PrepareMsg) error {
	if c.ConsensusID.Id != prepareMsg.ConsensusId.Id {
		return errors.New("Consensus ID is not same")
	}

	prepareMsgAddedEvent := event.PrepareMsgAdded{
		EventModel: midgard.EventModel{
			ID: c.ConsensusID.Id,
		},
		SenderId:  prepareMsg.SenderId,
		BlockHash: prepareMsg.BlockHash,
	}

	if err := OnAndSave(c, &prepareMsgAddedEvent); err != nil {
		return err
	}

	return nil
}

func (c *Consensus) SaveCommitMsg(commitMsg *CommitMsg) error {
	if c.ConsensusID.Id != commitMsg.ConsensusId.Id {
		return errors.New("Consensus ID is not same")
	}

	commitMsgAddedEvent := event.CommitMsgAdded{
		EventModel: midgard.EventModel{
			ID: c.ConsensusID.Id,
		},
		SenderId: commitMsg.SenderId,
	}

	if err := OnAndSave(c, &commitMsgAddedEvent); err != nil {
		return err
	}

	return nil
}

func (c *Consensus) On(consensusEvent midgard.Event) error {
	switch v := consensusEvent.(type) {

	case *event.PrepareMsgAdded:
		c.PrepareMsgPool.Save(&PrepareMsg{
			ConsensusId: ConsensusId{v.GetID()},
			SenderId:    v.SenderId,
			BlockHash:   v.BlockHash,
		})

	case *event.CommitMsgAdded:
		c.CommitMsgPool.Save(&CommitMsg{
			ConsensusId: ConsensusId{v.GetID()},
			SenderId:    v.SenderId,
		})

	case *event.ConsensusCreated:
		c.ConsensusID = ConsensusId{v.ConsensusId}
		for _, rStr := range v.Representatives {
			c.Representatives = append(c.Representatives, &Representative{
				Id: RepresentativeId(*rStr),
			})
		}
		c.Block.Body = v.Body
		c.Block.Seal = v.Seal
		c.CurrentState = State(v.CurrentState)
		c.PrepareMsgPool = NewPrepareMsgPool()
		c.CommitMsgPool = NewCommitMsgPool()

	case *event.ConsensusPrePrepared:
		c.Start()

	case *event.ConsensusPrepared:
		c.ToPrepareState()

	case *event.ConsensusCommitted:
		c.ToCommitState()

	case *event.ConsensusFinished:
		c.ToIdleState()

	default:
		return errors.New(fmt.Sprintf("unhandled event [%s]", v))
	}

	return nil
}

func OnAndSave(aggregate midgard.Aggregate, event midgard.Event) error {
	if err := aggregate.On(event); err != nil {
		return err
	}

	if err := eventstore.Save(event.GetID(), event); err != nil {
		return err
	}

	return nil
}
