package blockchain

import (
	"encoding/json"
	"time"

	"bytes"

	"errors"

	"fmt"

	"github.com/it-chain/it-chain-Engine/core/eventstore"
	"github.com/it-chain/midgard"
	ygg "github.com/it-chain/yggdrasill/common"
)

var ErrDecodingEmptyBlock = errors.New("Empty Block decoding failed")

type Block = ygg.Block

type BlockHeight = uint64

type BlockState = string

const (
	Created   BlockState = "Created"
	Staged    BlockState = "Staged"
	Committed BlockState = "Committed"
)

type DefaultBlock struct {
	Seal      []byte
	PrevSeal  []byte
	Height    uint64
	TxList    []Transaction
	TxSeal    [][]byte
	Timestamp time.Time
	Creator   []byte
	State     BlockState
}

func (block *DefaultBlock) On(event midgard.Event) error {

	switch v := event.(type) {

	case *BlockCreatedEvent:
		TxList, err := deserializeTxList(v.TxList)

		if err != nil {
			return ErrDeserializingTxList
		}

		block.Seal = v.Seal
		block.PrevSeal = v.PrevSeal
		block.Height = v.Height
		block.TxList = TxList
		block.TxSeal = v.TxSeal
		block.Timestamp = v.Timestamp
		block.Creator = v.Creator
		block.State = v.State

	default:
		return errors.New(fmt.Sprintf("unhandled event [%s]", v))
	}

	return nil
}

// TODO: Write test case
func (block *DefaultBlock) SetSeal(seal []byte) {
	block.Seal = seal
}

// TODO: Write test case
func (block *DefaultBlock) SetPrevSeal(prevSeal []byte) {
	block.PrevSeal = prevSeal
}

// TODO: Write test case
func (block *DefaultBlock) SetHeight(height uint64) {
	block.Height = height
}

// TODO: Write test case
func (block *DefaultBlock) PutTx(transaction Transaction) error {
	if block.TxList == nil {
		block.TxList = make([]Transaction, 0)
	}

	block.TxList = append(block.TxList, transaction)

	return nil
}

// TODO: Write test case
func (block *DefaultBlock) SetTxSeal(txSeal [][]byte) {
	block.TxSeal = txSeal
}

// TODO: Write test case
func (block *DefaultBlock) SetCreator(creator []byte) {
	block.Creator = creator
}

// TODO: Write test case
func (block *DefaultBlock) SetTimestamp(currentTime time.Time) {
	block.Timestamp = currentTime
}

// TODO: Write test case
func (block *DefaultBlock) GetSeal() []byte {
	return block.Seal
}

// TODO: Write test case
func (block *DefaultBlock) GetPrevSeal() []byte {
	return block.PrevSeal
}

// TODO: Write test case
func (block *DefaultBlock) GetHeight() uint64 {
	return block.Height
}

// TODO: Write test case
func (block *DefaultBlock) GetTxList() []Transaction {
	txList := make([]Transaction, 0)
	for _, tx := range block.TxList {
		txList = append(txList, tx)
	}
	return txList
}

// TODO: Write test case
func (block *DefaultBlock) GetTxSeal() [][]byte {
	return block.TxSeal
}

// TODO: Write test case
func (block *DefaultBlock) GetCreator() []byte {
	return block.Creator
}

// TODO: Write test case
func (block *DefaultBlock) GetTimestamp() time.Time {
	return block.Timestamp
}

// TODO: Write test case
func (block *DefaultBlock) Serialize() ([]byte, error) {
	data, err := json.Marshal(block)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// TODO: Write test case
func (block *DefaultBlock) Deserialize(serializedBlock []byte) error {
	if len(serializedBlock) == 0 {
		return ErrDecodingEmptyBlock
	}

	err := json.Unmarshal(serializedBlock, block)
	if err != nil {
		return err
	}

	return nil
}

// TODO: Write test case
func (block *DefaultBlock) IsReadyToPublish() bool {
	return block.Seal != nil
}

// TODO: Write test case
func (block *DefaultBlock) IsPrev(serializedPrevBlock []byte) bool {
	prevBlock := &DefaultBlock{}
	prevBlock.Deserialize(serializedPrevBlock)

	return bytes.Compare(prevBlock.GetSeal(), block.GetPrevSeal()) == 0
}

func CommitBlock(block Block) error {
	event, err := createBlockCommittedEvent(block)
	if err != nil {
		return err
	}
	blockId := string(block.GetSeal())
	eventstore.Save(blockId, event)
	return nil
}

func createBlockCommittedEvent(block Block) (BlockCommittedEvent, error) {
	seal := string(block.GetSeal())
	return BlockCommittedEvent{
		EventModel: midgard.EventModel{
			ID: seal,
		},
		State: Committed,
	}, nil
}

func StageBlock(block Block) error {
	event, err := createBlockStagedEvent(block)
	if err != nil {
		return err
	}

	blockId := string(block.GetSeal())

	err = eventstore.Save(blockId, event)
	if err != nil {
		return err
	}

	return nil
}

func createBlockStagedEvent(block Block) (BlockStagedEvent, error) {
	return BlockStagedEvent{
		EventModel: midgard.EventModel{
			ID: string(block.GetSeal()),
		},
		State: Staged,
	}, nil
}
