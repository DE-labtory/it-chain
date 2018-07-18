package blockchain

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"time"

	"github.com/it-chain/it-chain-Engine/common"
	"github.com/it-chain/it-chain-Engine/core/eventstore"
	"github.com/it-chain/midgard"
)

func CreateGenesisBlock(genesisconfFilePath string) (Block, error) {

	//declare
	GenesisBlock := &DefaultBlock{}
	validator := DefaultValidator{}
	TimeStamp := (time.Now()).Round(0)

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
	createEvent, err := createBlockCreatedEvent(Seal, GenesisBlock.PrevSeal, GenesisBlock.Height, convertTxType(GenesisBlock.TxList), GenesisBlock.TxSeal, TimeStamp, GenesisBlock.Creator)
	if err != nil {
		return nil, ErrCreatingEvent
	}

	//save
	eventstore.Save(createEvent.GetID(), createEvent)

	//on
	err = GenesisBlock.On(createEvent)
	if err != nil {
		return nil, ErrOnEvent
	}

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

func createBlockCreatedEvent(seal []byte, prevSeal []byte, height uint64, txList []Transaction, txSeal [][]byte, timeStamp time.Time, creator []byte) (*BlockCreatedEvent, error) {
	txListBytes, err := common.Serialize(txList)

	if err != nil {
		return &BlockCreatedEvent{}, err
	}

	return &BlockCreatedEvent{
		EventModel: midgard.EventModel{
			ID:   string(seal),
			Type: "block.created",
		},
		Seal:      seal,
		PrevSeal:  prevSeal,
		Height:    height,
		TxList:    txListBytes,
		TxSeal:    txSeal,
		Timestamp: timeStamp,
		Creator:   creator,
	}, nil
}

func CreateProposedBlock(prevSeal []byte, height uint64, txList []Transaction, Creator []byte) (Block, error) {

	//declare
	ProposedBlock := &DefaultBlock{}
	validator := DefaultValidator{}
	TimeStamp := (time.Now()).Round(0)

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
	createEvent, err := createBlockCreatedEvent(Seal, prevSeal, height, txList, txSeal, TimeStamp, Creator)
	if err != nil {
		return nil, ErrCreatingEvent
	}

	//save
	eventstore.Save(createEvent.GetID(), createEvent)

	//on
	err = ProposedBlock.On(createEvent)
	if err != nil {
		return nil, ErrOnEvent
	}

	return ProposedBlock, nil
}
