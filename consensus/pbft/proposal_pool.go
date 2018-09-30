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

package pbft

import (
	"errors"
	"time"

	"github.com/it-chain/engine/common/command"
)

var ErrBlockSealNil = errors.New("Start consensus command seal is nil")

type Proposal struct {
	Seal      []byte
	PrevSeal  []byte
	Height    uint64
	TxList    []command.Tx
	TxSeal    [][]byte
	Timestamp time.Time
	Creator   string
	State     string
}

type ProposalPool struct {
	proposals []Proposal
}

func NewProposalPool() ProposalPool {
	return ProposalPool{
		proposals: make([]Proposal, 0),
	}
}

func (p *ProposalPool) Save(startCommand command.StartConsensus) error {
	if startCommand.Seal == nil {
		return ErrBlockSealNil
	}

	proposal := Proposal{
		Seal:      startCommand.Seal,
		PrevSeal:  nil,
		Height:    0,
		TxList:    startCommand.TxList,
		TxSeal:    nil,
		Timestamp: time.Time{},
		Creator:   "",
		State:     "",
	}

	p.proposals = append(p.proposals, proposal)

	return nil
}

func (p *ProposalPool) Pop() Proposal {
	if len(p.proposals) == 0 {
		return Proposal{}
	}

	if len(p.proposals) == 1 {
		proposal := p.proposals[0]
		p.proposals = make([]Proposal, 0)
		return proposal
	}

	proposal := p.proposals[0]
	p.proposals = p.proposals[1:]
	return proposal
}

func (p *ProposalPool) RemoveAllMsg() {
	p.proposals = make([]Proposal, 0)
}
