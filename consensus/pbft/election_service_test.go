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

package pbft_test

import (
	"testing"

	"github.com/it-chain/engine/consensus/pbft"
	"github.com/magiconair/properties/assert"
)

func TestNewElectionService(t *testing.T) {
	e := pbft.NewElectionService("1", 30, pbft.CANDIDATE, 0)
	e.SetState(pbft.CANDIDATE)
	assert.Equal(t, e.GetState(), pbft.CANDIDATE)
	assert.Equal(t, e.GetLeftTime(), 30)
}

func TestElectionService_CountDownLeftTimeBy(t *testing.T) {
	e := SetElectionService()

	e.CountDownLeftTimeBy(1)

	assert.Equal(t, e.GetLeftTime(), 29)
}

func TestElectionService_CountUpVoteCount(t *testing.T) {

	e := SetElectionService()

	e.CountUpVoteCount()

	assert.Equal(t, e.GetVoteCount(), 1)
}

func TestElectionService_GetCandidate(t *testing.T) {
	e := SetElectionService()
	e.SetCandidate(&pbft.Representative{ID: "1"})

	assert.Equal(t, e.GetCandidate().ID, "1")
}

func TestElectionService_GetLeftTime(t *testing.T) {
	e := SetElectionService()
	assert.Equal(t, e.GetLeftTime(), 30)
}

func TestElectionService_GetState(t *testing.T) {
	e := SetElectionService()
	assert.Equal(t, e.GetState(), pbft.CANDIDATE)
}

func TestElectionService_GetVoteCount(t *testing.T) {
	e := SetElectionService()
	assert.Equal(t, e.GetVoteCount(), 0)
}

//todo
func TestElectionService_InitLeftTime(t *testing.T) {

}

//todo
func TestElectionService_ResetLeftTime(t *testing.T) {

}

func TestElectionService_ResetVoteCount(t *testing.T) {
	e := SetElectionService()

	e.CountUpVoteCount()
	e.ResetVoteCount()
	assert.Equal(t, e.GetVoteCount(), 0)
}

func TestElectionService_SetCandidate(t *testing.T) {
	e := SetElectionService()

	r := &pbft.Representative{
		ID: "1",
	}
	e.SetCandidate(r)

	assert.Equal(t, e.GetCandidate(), r)
}

func TestElectionService_SetLeftTime(t *testing.T) {
	e := SetElectionService()
	e.SetLeftTime(30)
	assert.Equal(t, e.GetLeftTime(), 30)
}

func TestElectionService_SetState(t *testing.T) {
	e := SetElectionService()
	e.SetState(pbft.TICKING)

	assert.Equal(t, e.GetState(), pbft.TICKING)

}

func TestElectionService_SetVoteCount(t *testing.T) {
	e := SetElectionService()
	e.SetVoteCount(100)

	assert.Equal(t, e.GetVoteCount(), 100)
}

func SetElectionService() *pbft.ElectionService {
	return pbft.NewElectionService("1", 30, pbft.CANDIDATE, 0)
}
