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

package api

import (
	"github.com/it-chain/engine/consensus"
	"github.com/it-chain/midgard"
)

type ParliamentApi struct {
	eventRepository *midgard.Repository
}

func NewParliamentApi(eventRepository *midgard.Repository) ParliamentApi {
	return ParliamentApi{
		eventRepository: eventRepository,
	}
}

// todo : Implement & Event Sourcing 첨가

func (p ParliamentApi) ChangeLeader(leader consensus.Leader) error {
	return nil
}

func (p ParliamentApi) AddMember(member consensus.Member) error {
	return nil
}

func (p ParliamentApi) RemoveMember(memberId consensus.MemberId) error {
	return nil
}
