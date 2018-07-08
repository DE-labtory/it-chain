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

	//declare
	GenesisBlock := &DefaultBlock{}
	validator := DefaultValidator{}
	TimeStamp := (time.Now()).Round(100 * time.Millisecond)

	//set
	err := setBlockWithConfig(genesisconfFilePath, GenesisBlock)

	if err != nil {
		return nil, ErrSetConfig
	}

	//build
	Seal, err := validator.BuildSeal(TimeStamp, GenesisBlock.PrevSeal, GenesisBlock.TxSeal, GenesisBlock.Creator)

	if err != nil {
		return nil, ErrBuildingSeal
	}

	//create
	createEvent := createBlockCreatedEvent(Seal, GenesisBlock.PrevSeal, GenesisBlock.Height, GenesisBlock.TxList, GenesisBlock.TxSeal, TimeStamp, GenesisBlock.Creator)

	//save
	eventstore.Save(createEvent.GetID(), createEvent)

	//on
	GenesisBlock.On(createEvent)

	return GenesisBlock, nil
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

func CreateProposedBlock(prevSeal []byte, height uint64, txList []Transaction, Creator []byte) (Block, error) {

	//declare
	ProposedBlock := &DefaultBlock{}
	validator := DefaultValidator{}
	TimeStamp := (time.Now()).Round(100 * time.Millisecond)

	//build
	txSeal, err := validator.BuildTxSeal(txList)

	if err != nil {
		return nil, ErrBuildingTxSeal
	}

	Seal, err := validator.BuildSeal(TimeStamp, prevSeal, txSeal, Creator)

	if err != nil {
		return nil, ErrBuildingSeal
	}

	//create
	createEvent := createBlockCreatedEvent(Seal, prevSeal, height, txList, txSeal, TimeStamp, Creator)

	//save
	eventstore.Save(createEvent.GetID(), createEvent)

	//on
	ProposedBlock.On(createEvent)

	return ProposedBlock, nil
}

func createBlockCreatedEvent(seal []byte, prevSeal []byte, height uint64, txList []Transaction, txSeal [][]byte, timeStamp time.Time, creator []byte) *BlockCreatedEvent {
	return &BlockCreatedEvent{
		EventModel: midgard.EventModel{
			ID:   string(seal),
			Type: "block.created",
		},
		Seal:      seal,
		PrevSeal:  prevSeal,
		Height:    height,
		TxList:    txList,
		TxSeal:    txSeal,
		Timestamp: timeStamp,
		Creator:   creator,
	}
}
