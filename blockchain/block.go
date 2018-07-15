package blockchain

import (
	"encoding/json"
	"time"

	"bytes"
	"errors"
	"fmt"

	"github.com/it-chain/midgard"
	ygg "github.com/it-chain/yggdrasill/common"
)

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
		block.TxList = make([]*DefaultTransaction, 0)
	}

	block.TxList = append(block.TxList, transaction.(*DefaultTransaction))

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

type Action interface {
	DoAction(block Block) error
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
	case *BlockStagedEvent:
		block.State = v.State
	case *BlockCommittedEvent:
		block.State = v.State
	default:
		return errors.New(fmt.Sprintf("unhandled event [%s]", v))
	}

	return nil
}

func IsBlockHasAllProperties(block Block) bool {
	return !(block.GetSeal() == nil || block.GetPrevSeal() == nil || block.GetHeight() == 0 ||
		block.GetTxList() == nil || block.GetCreator() == nil || block.GetTimestamp().IsZero())
}
