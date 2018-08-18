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
package mem

import (
	"sync"

	"errors"

	"github.com/it-chain/engine/consensus/pbft"
)

var ErrConsensusAlreadyExist = errors.New("State Already Exist")
var ErrLoadConsensus = errors.New("There is no state for loading")

type StateRepository struct {
	state *pbft.State
	sync.RWMutex
}

func NewStateRepository() StateRepository {
	return StateRepository{
		state:   nil,
		RWMutex: sync.RWMutex{},
	}
}
func (repo *StateRepository) Save(state pbft.State) error {

	repo.Lock()
	defer repo.Unlock()

	if repo.state != nil {
		return ErrConsensusAlreadyExist
	}
	repo.state = &state

	return nil
}
func (repo *StateRepository) Load() (*pbft.State, error) {

	if repo.state == nil {
		return nil, ErrLoadConsensus
	}

	return repo.state, nil
}

func (repo *StateRepository) Remove() {

	if repo.state != nil {
		repo.state = nil
	}
}
