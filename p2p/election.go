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

package p2p

import "sync"

type Election struct {
	leftTime  int64 //left time in millisecond
	state     string
	voteCount int
	mux       sync.Mutex
}

func (election *Election) SetLeftTime() int64 {

	election.mux.Lock()
	defer election.mux.Unlock()

	return election.leftTime
}

func (election *Election) ResetLeftTime() {

	election.mux.Lock()
	defer election.mux.Unlock()

	election.leftTime = genRandomInRange(150, 300)
}

//count down left time by tick millisecond  until 0
func (election *Election) CountDownLeftTimeBy(tick int64) {

	election.mux.Lock()
	defer election.mux.Unlock()

	if election.leftTime == 0 {
		return
	}

	election.leftTime = election.leftTime - tick
}

func (election *Election) SetState(state string) {

	election.mux.Lock()
	defer election.mux.Unlock()

	election.state = state
}

func (election *Election) GetState() string {

	election.mux.Lock()
	defer election.mux.Unlock()

	return election.state
}

func (election *Election) GetLeftTime() int64 {

	election.mux.Lock()
	defer election.mux.Unlock()

	return election.leftTime
}

func (election *Election) GetVoteCount() int {

	election.mux.Lock()
	defer election.mux.Unlock()

	return election.voteCount
}

func (election *Election) ResetVoteCount() {

	election.mux.Lock()
	defer election.mux.Unlock()

	election.voteCount = 0
}

func (election *Election) CountUp() {

	election.mux.Lock()
	defer election.mux.Unlock()

	election.voteCount = election.voteCount + 1
}

type ElectionRepository struct {
	election Election
	mux      sync.Mutex
}

func (er ElectionRepository) GetElection() Election {

	er.mux.Lock()
	defer er.mux.Unlock()

	return er.election
}

func (er ElectionRepository) SetElection(election Election) error {

	er.mux.Lock()
	defer er.mux.Unlock()

	er.election = election

	return nil
}
