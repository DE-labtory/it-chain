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
	"sync"

	"github.com/it-chain/engine/common/logger"
)

const (
	CANDIDATE ElectionState = "CANDIDATE"
	TICKING   ElectionState = "TICKING"
)

type ElectionState string

type ElectionService struct {
	NodeId    string
	candidate *Representative // candidate peer to be leader later
	leftTime  int             //left time in millisecond
	state     ElectionState
	voteCount int
	mux       sync.Mutex
}

func NewElectionService(id string, leftTime int, state ElectionState, voteCount int) *ElectionService {

	return &ElectionService{
		NodeId: id,
		candidate: &Representative{
			ID: "",
		},
		leftTime:  leftTime,
		state:     state,
		voteCount: voteCount,
		mux:       sync.Mutex{},
	}
}

func (e *ElectionService) SetLeftTime(time int) error {

	e.mux.Lock()
	defer e.mux.Unlock()

	e.leftTime = time
	return nil
}

func (e *ElectionService) SetVoteCount(count int) error {
	e.mux.Lock()
	defer e.mux.Unlock()

	e.voteCount = count
	return nil
}

func (e *ElectionService) InitLeftTime() {

	e.mux.Lock()
	defer e.mux.Unlock()

	e.leftTime = GenRandomInRange(150, 300)
}

func (e *ElectionService) ResetLeftTime() {

	e.mux.Lock()
	defer e.mux.Unlock()

	e.leftTime = GenRandomInRange(290, 300)
}

//count down left time by tick millisecond  until 0
func (e *ElectionService) CountDownLeftTimeBy(tick int) {

	if e.leftTime == 0 {
		return
	}

	e.leftTime = e.leftTime - tick
}

func (e *ElectionService) SetState(state ElectionState) {

	e.mux.Lock()
	defer e.mux.Unlock()

	logger.Infof(nil, "[consensus] set state to: %s", state)

	e.state = state
}

func (e *ElectionService) GetState() ElectionState {

	e.mux.Lock()
	defer e.mux.Unlock()

	return e.state
}

func (e *ElectionService) GetLeftTime() int {

	return e.leftTime
}

func (e *ElectionService) GetVoteCount() int {

	e.mux.Lock()
	defer e.mux.Unlock()

	return e.voteCount
}

func (e *ElectionService) ResetVoteCount() {

	e.mux.Lock()
	defer e.mux.Unlock()

	e.voteCount = 0
}

func (e *ElectionService) CountUpVoteCount() {

	e.mux.Lock()
	defer e.mux.Unlock()

	e.voteCount = e.voteCount + 1
}

func (e *ElectionService) SetCandidate(representative *Representative) {
	e.candidate = representative
}

func (e *ElectionService) GetCandidate() *Representative {
	return e.candidate
}
