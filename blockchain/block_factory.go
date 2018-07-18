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
	"encoding/json"
	"io/ioutil"
	"os"
	"time"

	"github.com/it-chain/engine/common"
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
	createEvent, err := createBlockCreatedEvent(Seal, GenesisBlock.PrevSeal, GenesisBlock.Height, convertTxType(GenesisBlock.TxList), GenesisBlock.TxSeal, GenesisBlock.Timestamp, GenesisBlock.Creator)
	if err != nil {
		return nil, ErrCreatingEvent
	}

	//on
	err = GenesisBlock.On(createEvent)
	if err != nil {
		return nil, ErrOnEvent
	}

	//save
	eventstore.Save(createEvent.GetID(), createEvent)

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
	createEvent, err := createBlockCreatedEvent(Seal, PrevSeal, Height, TxList, TxSeal, TimeStamp, Creator)
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
