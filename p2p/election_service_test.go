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

package p2p_test

import (
	"testing"

	"github.com/it-chain/engine/p2p"
	"github.com/it-chain/engine/p2p/test/mock"
	"github.com/stretchr/testify/assert"
)

func TestNewElectionService(t *testing.T) {
	election := p2p.NewElection("this.is.ip.addres", 30, p2p.TICKING, 123)
	queryService := mock.MockPeerQueryService{}
	client := mock.MockClient{}

	electionService := p2p.NewElectionService(&election, queryService, client)

	assert.Equal(t, 123, electionService.Election.GetVoteCount())
	assert.Equal(t, 30, electionService.Election.GetLeftTime())
	assert.Equal(t, p2p.TICKING, electionService.Election.GetState())
}

func TestGenRandomInRange(t *testing.T) {
	v1 := p2p.GenRandomInRange(0, 10)
	v2 := p2p.GenRandomInRange(0, 10)
	v3 := p2p.GenRandomInRange(0, 10)

	t.Logf("%v", v1)
	t.Logf("%v", v2)
	t.Logf("%v", v3)
}
