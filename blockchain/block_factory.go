package blockchain

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"time"

	"github.com/it-chain/it-chain-Engine/core/eventstore"
	"github.com/it-chain/midgard"
)

var ErrSetConfig = errors.New("error when set Config")
var ErrBuildingSeal = errors.New("error when building Seal")
var ErrBuildingTxSeal = errors.New("error when building TxSeal")

func CreateGenesisBlock(genesisconfFilePath string) (Block, error) {

	GenesisBlock := &DefaultBlock{}
	validator := DefaultValidator{}

	err := setBlockWithConfig(genesisconfFilePath, GenesisBlock)

	if err != nil {
		return nil, ErrSetConfig
	}

	GenesisBlock.SetTimestamp((time.Now()).Round(0))

	Seal, err := validator.BuildSeal((time.Now()).Round(0), GenesisBlock.GetPrevSeal(), GenesisBlock.TxSeal, GenesisBlock.Creator)

	if err != nil {
		return nil, ErrBuildingSeal
	}

	createEvent := &BlockCreatedEvent{
		EventModel: midgard.EventModel{
			ID:   string(Seal),
			Type: "block.created",
		},

		Seal:      Seal,
		PrevSeal:  GenesisBlock.PrevSeal,
		Height:    GenesisBlock.Height,
		TxList:    GenesisBlock.TxList,
		TxSeal:    GenesisBlock.TxSeal,
		Timestamp: GenesisBlock.Timestamp,
		Creator:   GenesisBlock.Creator,
	}

	eventstore.Save(createEvent.GetID(), createEvent)

	GenesisBlock.On(createEvent)

	return GenesisBlock, nil
}

func CreateProposedBlock(prevSeal []byte, height uint64, txList []Transaction, Creator []byte) (Block, error) {

	ProposedBlock := &DefaultBlock{}
	validator := DefaultValidator{}

	ProposedBlock.SetPrevSeal(prevSeal)

	ProposedBlock.SetCreator(Creator)

	for _, tx := range txList {
		ProposedBlock.PutTx(tx)
	}

	txSeal, err := validator.BuildTxSeal(ProposedBlock.GetTxList())

	if err != nil {
		return nil, ErrBuildingTxSeal
	}

	ProposedBlock.SetTxSeal(txSeal)

	ProposedBlock.SetTimestamp((time.Now()).Round(0))

	Seal, err := validator.BuildSeal((time.Now()).Round(0), prevSeal, txSeal, Creator)

	if err != nil {
		return nil, ErrBuildingSeal
	}

	createEvent := &BlockCreatedEvent{
		EventModel: midgard.EventModel{
			ID:   string(Seal),
			Type: "block.created",
		},
		Seal:      Seal,
		PrevSeal:  ProposedBlock.PrevSeal,
		Height:    height,
		TxList:    ProposedBlock.TxList,
		TxSeal:    ProposedBlock.TxSeal,
		Timestamp: ProposedBlock.Timestamp,
		Creator:   ProposedBlock.Creator,
	}

	eventstore.Save(createEvent.GetID(), createEvent)

	ProposedBlock.On(createEvent)

	return ProposedBlock, nil
}

func setBlockWithConfig(filePath string, block Block) error {
	jsonFile, err := os.Open(filePath)
	defer jsonFile.Close()
	if err != nil {
		return err
	}

	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		return err
	}

	err = json.Unmarshal(byteValue, block)
	if err != nil {
		return err
	}

	return nil
}

func createBlockCreatedEvent(seal []byte, prevSeal []byte, height uint64, txList []Transaction, txSeal [][]byte, creator []byte) BlockCreatedEvent {
	return BlockCreatedEvent{
		EventModel: midgard.EventModel{
			ID:   string(seal),
			Type: "block.created",
		},
		Seal:      seal,
		PrevSeal:  prevSeal,
		Height:    height,
		TxList:    txList,
		TxSeal:    txSeal,
		Timestamp: (time.Now()).Round(0),
		Creator:   creator,
	}
}
