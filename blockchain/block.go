package blockchain

import (
	"encoding/json"
	"time"

	"bytes"

	"github.com/it-chain/midgard"
	ygg "github.com/it-chain/yggdrasill/common"

	"errors"
	"fmt"
)

type Block = ygg.Block

type BlockHeight = uint64

type DefaultBlock struct {
	Seal      []byte
	PrevSeal  []byte
	Height    uint64
	TxList    []*DefaultTransaction
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
	convTx, ok := transaction.(*DefaultTransaction)
	if ok {
		if block.TxList == nil {
			block.TxList = make([]*DefaultTransaction, 0)
		}
		block.TxList = append(block.TxList, convTx)

		return nil
	}

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

	default:
		return errors.New(fmt.Sprintf("unhandled event [%s]", v))
	}

	return nil
}

//ToDo: 삭제 제안 - junk_sound
////interface of api gateway query api
//type BlockQueryApi interface {
//	GetBlockByHeight(blockHeight uint64) (Block, error)
//	GetBlockBySeal(seal []byte) (Block, error)
//	GetBlockByTxID(txid string) (Block, error)
//	GetLastBlock() (Block, error)
//}

// ToDo: 미완성 부분에 대한 주석처리
//type Action interface {
//	DoAction(block Block) error
//}
//
//// TODO: Write test case
//func CreateSaveOrSyncAction(checkResult int64) Action {
//	if checkResult > 0 {
//		return NewSyncAction()
//	} else if checkResult == 0 {
//		return NewSaveAction()
//	} else {
//		return NewDefaultAction()
//	}
//}
//
//type SyncAction struct{}
//
//func NewSyncAction() *SyncAction {
//	return &SyncAction{}
//}

//ToDo: 삭제 제안 - junk_sound
//func NewEmptyBlock(prevSeal []byte, height uint64, creator []byte) *DefaultBlock {
//	block := &DefaultBlock{}
//	return block
//}

// ToDo: 미완성 부분에 대한 주석처리
//func (syncAction *SyncAction) DoAction(block Block) error {
//	// TODO: Start synchronize
//	return nil
//}
//
//type SaveAction struct {
//	blockPool BlockPool
//}
//
//func NewSaveAction() *SaveAction {
//	return &SaveAction{}
//}
//
//// TODO: Write test case
//func (saveAction *SaveAction) DoAction(block Block) error {
//	event, err := createBlockCommittedEvent(block)
//	if err != nil {
//		return err
//	}
//	blockId := string(block.GetSeal())
//	eventstore.Save(blockId, event)
//	return nil
//}
//
//func createBlockCommittedEvent(block Block) (BlockCommittedEvent, error) {
//	seal := string(block.GetSeal())
//	return BlockCommittedEvent{
//		EventModel: midgard.EventModel{
//			ID: seal,
//		},
//		Seal: seal,
//	}, nil
//}
//
//type DefaultAction struct{}
//
//func NewDefaultAction() *DefaultAction {
//	return &DefaultAction{}
//}
//
//// TODO: Write test case
//func (defaultAction *DefaultAction) DoAction(block Block) error {
//	log.Printf("got shorter height block [%v]", block.GetHeight())
//	return nil
//}
