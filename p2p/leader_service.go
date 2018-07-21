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

import (
	"github.com/it-chain/engine/core/eventstore"
	"github.com/it-chain/midgard"
)

type ILeaderService interface {
	Set(leader Leader) error
}

type LeaderService struct{}

func NewLeaderService() *LeaderService {

	return &LeaderService{}
}

func (ls *LeaderService) Set(leader Leader) error {

	event := LeaderUpdatedEvent{
		EventModel: midgard.EventModel{
			ID: leader.LeaderId.Id,
		},
	}

	return eventstore.Save(leader.LeaderId.Id, event)
}
