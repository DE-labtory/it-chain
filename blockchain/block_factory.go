package blockchain

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"time"

	"github.com/it-chain/it-chain-Engine/core/eventstore"
	"github.com/it-chain/midgard"
	"github.com/rs/xid"
)

var ErrGetConfig = errors.New("error when get Config")
var ErrBuildingSeal = errors.New("error when building Seal")
var ErrBuildingTxSeal = errors.New("error when building TxSeal")

func CreateGenesisBlock(genesisconfFilePath string) (Block, error) {

	GenesisBlock := &DefaultBlock{}
	validator := DefaultValidator{}

	byteValue, err := configFromJson(genesisconfFilePath)

	if err != nil {
		return nil, ErrGetConfig
	}

	json.Unmarshal(byteValue, &GenesisBlock)

	GenesisBlock.SetTimestamp((time.Now()).Round(0))

	Seal, err := validator.BuildSeal(GenesisBlock)

	if err != nil {
		return nil, ErrBuildingSeal
	}

	createEvent := BlockCreatedEvent{
		EventModel: midgard.EventModel{
			ID:   xid.New().String(),
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

	GenesisBlock.On(&createEvent)

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

	Seal, err := validator.BuildSeal(ProposedBlock)

	if err != nil {
		return nil, ErrBuildingSeal
	}

	createEvent := BlockCreatedEvent{
		EventModel: midgard.EventModel{
			ID:   xid.New().String(),
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

	ProposedBlock.On(&createEvent)

	return ProposedBlock, nil
}

func configFromJson(filePath string) ([]uint8, error) {
	jsonFile, err := os.Open(filePath)
	defer jsonFile.Close()
	if err != nil {
		return nil, err
	}

	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		return nil, err
	}
	return byteValue, nil
}
