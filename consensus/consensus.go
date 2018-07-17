package consensus

import (
	"errors"
	"fmt"

	"encoding/json"

	"github.com/it-chain/it-chain-Engine/core/eventstore"
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
	body []byte
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

	prepareMsgAddedEvent := PrepareMsgAddedEvent{
		EventModel: midgard.EventModel{
			ID: c.ConsensusID.Id,
		},
		PrepareMsg: struct {
			ConsensusId ConsensusId
			SenderId    string
			BlockHash   []byte
		}{ConsensusId: prepareMsg.ConsensusId, SenderId: prepareMsg.SenderId, BlockHash: prepareMsg.BlockHash},
	}

	if err := c.On(&prepareMsgAddedEvent); err != nil {
		return err
	}

	if err := eventstore.Save(prepareMsgAddedEvent.ID, prepareMsgAddedEvent); err != nil {
		return err
	}

	return nil
}

func (c *Consensus) SaveCommitMsg(commitMsg *CommitMsg) error {
	if c.ConsensusID.Id != commitMsg.ConsensusId.Id {
		return errors.New("Consensus ID is not same")
	}

	commitMsgAddedEvent := CommitMsgAddedEvent{
		EventModel: midgard.EventModel{
			ID: c.ConsensusID.Id,
		},
		CommitMsg: struct {
			ConsensusId ConsensusId
			SenderId    string
		}{ConsensusId: commitMsg.ConsensusId, SenderId: commitMsg.SenderId},
	}

	if err := c.On(&commitMsgAddedEvent); err != nil {
		return err
	}

	if err := eventstore.Save(c.GetID(), commitMsgAddedEvent); err != nil {
		return err
	}

	return nil
}

func (c *Consensus) On(event midgard.Event) error {
	switch v := event.(type) {

	case *PrepareMsgAddedEvent:
		c.PrepareMsgPool.Save(&PrepareMsg{
			ConsensusId: v.PrepareMsg.ConsensusId,
			SenderId:    v.PrepareMsg.SenderId,
			BlockHash:   v.PrepareMsg.BlockHash,
		})

	case *CommitMsgAddedEvent:
		c.CommitMsgPool.Save(&CommitMsg{
			ConsensusId: v.CommitMsg.ConsensusId,
			SenderId:    v.CommitMsg.SenderId,
		})

	case *ConsensusCreatedEvent:
		c.ConsensusID = v.Consensus.ConsensusID
		c.Representatives = v.Consensus.Representatives
		c.Block = v.Consensus.Block
		c.Start()
		c.PrepareMsgPool = v.Consensus.PrepareMsgPool
		c.CommitMsgPool = v.Consensus.CommitMsgPool

	case *ConsensusPrePreparedEvent:
		c.Start()

	case *ConsensusPreparedEvent:
		c.ToPrepareState()

	case *ConsensusCommittedEvent:
		c.ToCommitState()

	case *ConsensusFinishedEvent:
		c.ToIdleState()

	default:
		return errors.New(fmt.Sprintf("unhandled event [%s]", v))
	}

	return nil
}
