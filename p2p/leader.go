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
	"errors"
	"fmt"

	"github.com/it-chain/engine/common/event"
	"github.com/it-chain/engine/core/eventstore"
	"github.com/it-chain/midgard"
)

type Leader struct {
	LeaderId LeaderId
}

type LeaderId struct {
	Id string
}

func (lid LeaderId) ToString() string {
	return string(lid.Id)
}

func (l Leader) GetID() string {
	return l.LeaderId.ToString()
}

func (l *Leader) On(e midgard.Event) error {

	switch v := e.(type) {

	case event.LeaderChanged:
		l.LeaderId = LeaderId{v.GetID()}

	default:
		return errors.New(fmt.Sprintf("unhandled event [%s]", v))
	}

	return nil
}

func UpdateLeader(peer Peer) error {

	leader := Leader{
		LeaderId: LeaderId{Id: peer.PeerId.Id},
	}

	if leader.LeaderId.Id == "" {
		return ErrEmptyLeaderId
	}

	events := make([]midgard.Event, 0)

	leaderUpdatedEvent := event.LeaderUpdated{
		EventModel: midgard.EventModel{
			ID:   leader.LeaderId.ToString(),
			Type: "leader.update",
		},
	}

	leader.On(leaderUpdatedEvent)

	events = append(events, leaderUpdatedEvent)

	return eventstore.Save(leaderUpdatedEvent.GetID(), events...)

}
