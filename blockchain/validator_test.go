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

package blockchain_test

import (
	"testing"
	"time"

	"github.com/it-chain/engine/blockchain"
	"github.com/it-chain/engine/blockchain/test/mock"
	"github.com/stretchr/testify/assert"
)

func TestDefaultValidator_BuildAndValidateSeal(t *testing.T) {

	//given
	validator := blockchain.DefaultValidator{}
	TimeStamp := time.Now().Round(0)
	PrevSeal := []byte("PrevSeal")
	TxList := mock.GetRandomTxList()
	Creator := "Creator"
	Tree, err := validator.BuildTree(convertTxType(TxList))

	assert.NoError(t, err)

	//when

	Seal, err := validator.BuildSeal(TimeStamp, PrevSeal, Tree.GetTxSealRoot(), Creator)

	//then
	assert.NoError(t, err)

	//given
	block := blockchain.DefaultBlock{
		Seal:      Seal,
		Timestamp: TimeStamp,
		Tree:      blockchain.ConvTreeType(Tree),
		PrevSeal:  PrevSeal,
		Creator:   Creator,
	}

	//when
	result, err := validator.ValidateSeal(Seal, &block)

	//then
	assert.Equal(t, true, result)
	assert.NoError(t, err)

}

func TestDefaultValidator_BuildTreeAndValidate(t *testing.T) {
	validator := blockchain.DefaultValidator{}

	TxList := mock.GetRandomTxList()

	tree, err := validator.BuildTree(convertTxType(TxList))
	assert.NoError(t, err)

	isValidated, err := validator.ValidateTree(tree)
	assert.NoError(t, err)
	assert.Equal(t, true, isValidated)
}

func TestDefaultValidator_ValidateTree_GenesisBlock(t *testing.T) {
	validator := blockchain.DefaultValidator{}

	genesisTxList := []blockchain.Transaction{}

	tree, err := validator.BuildTree(genesisTxList)
	assert.NoError(t, err)

	isValidated, err := validator.ValidateTree(tree)
	assert.NoError(t, err)
	assert.Equal(t, true, isValidated)
}
