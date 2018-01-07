package blockchainleveldb

import (
	"fmt"
	"bytes"
	"encoding/gob"
	"it-chain/service/blockchain"
)

func Serialize(block *blockchain.Block) ([]byte, error) {
	var b bytes.Buffer

	encoder := gob.NewEncoder(&b)
	err := encoder.Encode(block)

	if err != nil {
		panic(fmt.Sprintf("Error encoding block : %s", err))
	}

	return b.Bytes(), err
}

func Deserialize(serializedBytes []byte) (*blockchain.Block, error) {
	var b bytes.Buffer
	block := &blockchain.Block{}

	b.Write(serializedBytes)
	decoder := gob.NewDecoder(&b)
	err := decoder.Decode(block)

	if err != nil {
		panic(fmt.Sprintf("Error decoding block : %s", err))
	}

	return block, err
}