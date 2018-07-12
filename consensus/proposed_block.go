package consensus

import (
	"encoding/json"
	"errors"
)

var ErrDecodingEmptyBlock = errors.New("Empty Block decoding failed")

type ProposedBlock struct {
	Seal []byte
	body []byte
}

func (block *ProposedBlock) Serialize() ([]byte, error) {
	data, err := json.Marshal(block)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (block *ProposedBlock) Deserialize(serializedBlock []byte) error {
	if len(serializedBlock) == 0 {
		return ErrDecodingEmptyBlock
	}

	err := json.Unmarshal(serializedBlock, block)
	if err != nil {
		return err
	}

	return nil
}
