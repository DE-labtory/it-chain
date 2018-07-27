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

package blockchain

import (
	"io/ioutil"
	"os"
	"time"

	"encoding/json"

	"encoding/hex"

	"github.com/it-chain/engine/common/event"
	"github.com/it-chain/engine/core/eventstore"
	"github.com/it-chain/midgard"
)

func CreateGenesisBlock(genesisconfFilePath string) (Block, error) {

	//declare
	GenesisBlock := &DefaultBlock{}
	validator := DefaultValidator{}

	//set
	err := setBlockWithConfig(genesisconfFilePath, GenesisBlock)

	if err != nil {
		return nil, ErrSetConfig
	}

	//build
	Seal, err := validator.BuildSeal(GenesisBlock.Timestamp, GenesisBlock.PrevSeal, GenesisBlock.TxSeal, GenesisBlock.Creator)

	if err != nil {
		return nil, ErrBuildingSeal
	}

	//create
	createEvent, err := createBlockCreatedEvent(Seal, GenesisBlock.PrevSeal, GenesisBlock.Height, GenesisBlock.TxList, GenesisBlock.TxSeal, GenesisBlock.Timestamp, GenesisBlock.Creator)
	if err != nil {
		return nil, ErrCreatingEvent
	}

	//on
	err = GenesisBlock.On(createEvent)
	if err != nil {
		return nil, ErrOnEvent
	}

	//save
	err = eventstore.Save(createEvent.GetID(), createEvent)

	if err != nil {
		return nil, ErrSavingEvent
	}

	return GenesisBlock, nil
}

func setBlockWithConfig(filePath string, block Block) error {

	// load
	jsonFile, err := os.Open(filePath)
	defer jsonFile.Close()

	if err != nil {
		return err
	}

	byteValue, err := ioutil.ReadAll(jsonFile)

	if err != nil {
		return err
	}

	GenesisConfig := &GenesisConfig{}

	err = json.Unmarshal(byteValue, GenesisConfig)
	if err != nil {
		return err
	}

	// set
	const longForm = "Jan 1, 2006 at 0:00am (MST)"

	timeStamp, err := time.Parse(longForm, GenesisConfig.TimeStamp)

	if err != nil {
		return err
	}

	block.SetPrevSeal(make([]byte, 0))
	block.SetHeight(uint64(GenesisConfig.Height))
	block.SetTxSeal(make([][]byte, 0))
	block.SetTimestamp(timeStamp)
	block.SetCreator([]byte(GenesisConfig.Creator))

	return nil
}

type GenesisConfig struct {
	Organization string
	NedworkId    string
	Height       int
	TimeStamp    string
	Creator      string
}

func createBlockCreatedEvent(seal []byte, prevSeal []byte, height uint64, txList []*DefaultTransaction, txSeal [][]byte, timeStamp time.Time, creator []byte) (*event.BlockCreated, error) {

	AggregateID := hex.EncodeToString(seal)

	return &event.BlockCreated{
		EventModel: midgard.EventModel{
			ID:   AggregateID,
			Type: "block.created",
		},
		Seal:      seal,
		PrevSeal:  prevSeal,
		Height:    height,
		TxList:    ConvBackFromTransactionList(txList),
		TxSeal:    txSeal,
		Timestamp: timeStamp,
		Creator:   creator,
		State:     Created,
	}, nil
}

func CreateProposedBlock(prevSeal []byte, height uint64, txList []*DefaultTransaction, Creator []byte) (Block, error) {

	//declare
	ProposedBlock := &DefaultBlock{}
	validator := DefaultValidator{}
	TimeStamp := time.Now().Round(0)

	//build
	txSeal, err := validator.BuildTxSeal(ConvertTxType(txList))

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

	//on
	err = ProposedBlock.On(createEvent)
	if err != nil {
		return nil, ErrOnEvent
	}

	//save
	eventstore.Save(createEvent.GetID(), createEvent)

	return ProposedBlock, nil
}

func CreateRetrievedBlock(retrievedBlock Block) (Block, error) {

	//declare
	RetrievedBlock := &DefaultBlock{}

	Seal := retrievedBlock.GetSeal()
	PrevSeal := retrievedBlock.GetPrevSeal()
	Height := retrievedBlock.GetHeight()
	TxList := retrievedBlock.GetTxList()
	TxSeal := retrievedBlock.GetTxSeal()
	TimeStamp := retrievedBlock.GetTimestamp()
	Creator := retrievedBlock.GetCreator()

	//create
	createEvent, err := createBlockCreatedEvent(Seal, PrevSeal, Height, GetBackTxType(TxList), TxSeal, TimeStamp, Creator)
	if err != nil {
		return nil, ErrCreatingEvent
	}

	//on
	err = RetrievedBlock.On(createEvent)
	if err != nil {
		return nil, ErrOnEvent
	}

	//save
	eventstore.Save(createEvent.GetID(), createEvent)

	return RetrievedBlock, nil
}
