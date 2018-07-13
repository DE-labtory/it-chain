package blockchain

import (
	"encoding/json"
	"time"

	"bytes"

	"errors"

	ygg "github.com/it-chain/yggdrasill/common"
	"github.com/it-chain/midgard"
	"github.com/it-chain/it-chain-Engine/core/eventstore"
	"github.com/it-chain/it-chain-Engine/common"
)

var ErrDecodingEmptyBlock = errors.New("Empty Block decoding failed")
var ErrTransactionType = errors.New("Wrong transaction type")

type Block = ygg.Block

type BlockHeight = uint64

type DefaultBlock struct {
	Seal      []byte
	PrevSeal  []byte
	Height    uint64
	TxList    []Transaction
	TxSeal    [][]byte
	Timestamp time.Time
	Creator   []byte
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

	return ErrTransactionType
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

// This is from #279 @junk-sound
func deserializeTxList(txList []byte) ([]Transaction, error) {
	DefaultTxList := []*DefaultTransaction{}

	err := common.Deserialize(txList, &DefaultTxList)

	if err != nil {
		return nil, err
	}
	TxList := convertTxType(DefaultTxList)

	return TxList, nil
}

func convertTxType(txList []*DefaultTransaction) []Transaction {
	convTxList := make([]Transaction, 0)

	for _, tx := range txList {
		convTxList = append(convTxList, tx)
	}

	return convTxList
}

// interface of api gateway query api
type BlockQueryApi interface {
	GetBlockByHeight(blockHeight uint64) (Block, error)
	GetBlockBySeal(seal []byte) (Block, error)
	GetBlockByTxID(txid string) (Block, error)
	GetLastBlock() (Block, error)
	GetTransactionByTxID(txid string) (Transaction, error)
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
		Seal: seal,
	}, nil
}
