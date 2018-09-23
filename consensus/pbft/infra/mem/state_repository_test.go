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
package mem_test

import (
	"testing"

	"github.com/it-chain/engine/consensus/pbft"
	"github.com/it-chain/engine/consensus/pbft/infra/mem"
	"github.com/stretchr/testify/assert"
)

func TestConsensusRepository_Save(t *testing.T) {

	mock1 := pbft.State{
		StateID: pbft.StateID{"mock1"},
	}
	repo := mem.NewStateRepository()
	err := repo.Save(mock1)
	assert.Equal(t, nil, err)
	mock2 := pbft.State{
		StateID: pbft.StateID{"mock2"},
	}
	err2 := repo.Save(mock2)
	assert.Equal(t, pbft.ErrInvalidSave, err2)

}

func TestConsensusRepository_Load(t *testing.T) {

	repo := mem.NewStateRepository()
	_, err := repo.Load()
	// case1 : Repository has no consensus
	assert.Equal(t, err, pbft.ErrEmptyRepo)

	// case2 : Repository has consensus
	mockConsensus := pbft.State{
		StateID: pbft.StateID{"hihi"},
	}
	repo.Save(mockConsensus)

	_, err2 := repo.Load()
	assert.Nil(t, err2)

}
func TestConsensusRepository_Remove(t *testing.T) {
	repo := mem.NewStateRepository()
	mockConsensus := pbft.State{
		StateID: pbft.StateID{"hihi"},
	}
	repo.Save(mockConsensus)
	repo.Remove()
	_, err := repo.Load()
	assert.Equal(t, pbft.ErrEmptyRepo, err)

}
