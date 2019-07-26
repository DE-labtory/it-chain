/*
 * Copyright 2018 DE-labtory
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
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewParliament(t *testing.T) {

	new_p := NewParliament()
	new_r := make(map[string]Representative)
	assert.Equal(t, new_p.Leader, Leader{})
	assert.Equal(t, new_p.Representatives, new_r)
}

func TestNewRepresentative(t *testing.T) {

	new_r := NewRepresentative("hero")
	assert.Equal(t, new_r.ID, "hero")

}

func TestRepresentative_GetID(t *testing.T) {

	new_r := NewRepresentative("hero")
	myid := new_r.GetID()
	assert.Equal(t, myid, "hero")

}

func TestLeader_GetID(t *testing.T) {

	// 1.make parliament
	new_p := NewParliament()

	// 2. set leader
	new_p.SetLeader("hero")

	// 3. get leader id
	ourLeader := new_p.GetLeader()
	ourLeader.LeaderId = "hero"
	leaderId := ourLeader.GetID()
	assert.Equal(t, leaderId, "hero")
}
