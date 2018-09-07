package blockchain_test

import (
	"testing"

	"time"

	"github.com/it-chain/engine/blockchain"
	"github.com/stretchr/testify/assert"
)

func TestDefaultValidator_BuildAndValidateSeal(t *testing.T) {

	//given
	validator := blockchain.DefaultValidator{}
	TimeStamp := time.Now().Round(0)
	PrevSeal := []byte("PrevSeal")
	TxSeal := make([][]byte, 0)
	Creator := "Creator"

	//when

	Seal, err := validator.BuildSeal(TimeStamp, PrevSeal, TxSeal, Creator)

	//then
	assert.NoError(t, err)

	block := blockchain.DefaultBlock{
		Seal:      Seal,
		Timestamp: TimeStamp,
		PrevSeal:  PrevSeal,
		TxSeal:    TxSeal,
		Creator:   Creator,
	}

	//when
	result, err := validator.ValidateSeal(Seal, &block)

	//then
	assert.Equal(t, true, result)

}

func TestDefaultValidator_BuildAndValidateTxSeal(t *testing.T) {
	//given
	validator := blockchain.DefaultValidator{}
	TxList := []*blockchain.DefaultTransaction{
		{
			ID:        "tx01",
			ICodeID:   "Icode01",
			PeerID:    "Peer01",
			Timestamp: time.Now().Round(0),
			Jsonrpc:   "jsonrpc01",
			Function:  "function01",
			Args:      nil,
			Signature: []byte("Signature"),
		},
		{
			ID:        "tx02",
			ICodeID:   "Icode02",
			PeerID:    "Peer02",
			Timestamp: time.Now().Round(0),
			Jsonrpc:   "jsonrpc02",
			Function:  "function02",
			Args:      nil,
			Signature: []byte("Signature"),
		},
	}

	ConvertedTxList := blockchain.ConvertTxType(TxList)

	//when
	txSeal, err := validator.BuildTxSeal(ConvertedTxList)

	//then
	assert.NoError(t, err)

	//when
	result1, err := validator.ValidateTxSeal(txSeal, ConvertedTxList)

	//then
	assert.NoError(t, err)
	assert.Equal(t, true, result1)

	//when
	result2, err := validator.ValidateTransaction(txSeal, TxList[0])

	//then
	assert.NoError(t, err)
	assert.Equal(t, true, result2)

}
