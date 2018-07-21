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

package consensus

import (
	"errors"
	"fmt"

	"github.com/it-chain/engine/common/event"
	"github.com/it-chain/engine/core/eventstore"
	"github.com/it-chain/midgard"
)

type LeaderId struct {
	Id string
}

func (lid LeaderId) ToString() string {
	return string(lid.Id)
}

type Leader struct {
	LeaderId LeaderId
}

func (l *Leader) StringLeaderId() string {
	return l.LeaderId.ToString()
}

func (l Leader) GetID() string {
	return l.StringLeaderId()
}

type MemberId struct {
	Id string
}

func (mid MemberId) ToString() string {
	return string(mid.Id)
}

type Member struct {
	MemberId MemberId
}

func (m *Member) StringMemberId() string {
	return m.MemberId.ToString()
}

func (m Member) GetID() string {
	return m.StringMemberId()
}

type ParliamentId string

func (pId ParliamentId) ToString() string {
	return string(pId)
}

type Parliament struct {
	ParliamentId ParliamentId
	Leader       *Leader
	Members      []*Member
}

func NewParliament() Parliament {
	return Parliament{
		ParliamentId: ParliamentId("0"),
		Members:      make([]*Member, 0),
		Leader:       nil,
	}
}

func (p *Parliament) GetID() string {
	return p.ParliamentId.ToString()
}

func (p *Parliament) IsNeedConsensus() bool {
	numOfMember := 0

	if p.HasLeader() {
		numOfMember = numOfMember + 1
	}

	numOfMember = numOfMember + len(p.Members)

	if numOfMember >= 1 {
		return true
	}

	return false
}

func (p *Parliament) HasLeader() bool {
	if p.Leader == nil {
		return false
	}

	return true
}

func (p *Parliament) ChangeLeader(leader *Leader) error {
	if leader == nil {
		return errors.New("Leader is nil")
	}

	leaderChangedEvent := event.LeaderChanged{
		EventModel: midgard.EventModel{
			ID: p.GetID(),
		},
		LeaderId: leader.GetID(),
	}

	if err := p.On(&leaderChangedEvent); err != nil {
		return err
	}

	if err := eventstore.Save(p.GetID(), leaderChangedEvent); err != nil {
		return err
	}

	return nil
}

func (p *Parliament) AddMember(member *Member) error {
	if member == nil {
		return errors.New("Member is nil")
	}

	if member.GetID() == "" {
		return errors.New(fmt.Sprintf("Need Valid PeerID [%s]", member.GetID()))
	}

	index := p.findIndexOfMember(member.GetID())

	if index != -1 {
		return errors.New(fmt.Sprintf("Already exist member [%s]", member.GetID()))
	}

	memberJoinedEvent := event.MemberJoined{
		EventModel: midgard.EventModel{
			ID: p.GetID(),
		},
		MemberId: member.GetID(),
	}

	if err := p.On(&memberJoinedEvent); err != nil {
		return err
	}

	if err := eventstore.Save(p.GetID(), memberJoinedEvent); err != nil {
		return err
	}

	return nil
}

func (p *Parliament) RemoveMember(memberID MemberId) error {
	index := p.findIndexOfMember(memberID.ToString())

	if index == -1 {
		return nil
	}

	memberRemovedEvent := event.MemberRemoved{
		EventModel: midgard.EventModel{
			ID: p.GetID(),
		},
		MemberId: memberID.ToString(),
	}

	if err := p.On(&memberRemovedEvent); err != nil {
		return err
	}

	if err := eventstore.Save(p.GetID(), memberRemovedEvent); err != nil {
		return err
	}

	return nil
}

func (p *Parliament) ValidateRepresentative(representatives []*Representative) bool {
	for _, representatives := range representatives {
		index := p.findIndexOfMember(representatives.GetID())

		if index == -1 {
			return false
		}
	}

	return true
}

func (p *Parliament) findIndexOfMember(memberID string) int {
	for i, member := range p.Members {
		if member.MemberId.Id == memberID {
			return i
		}
	}

	return -1
}

func (p *Parliament) FindByPeerID(memberID string) *Member {
	index := p.findIndexOfMember(memberID)

	if index == -1 {
		return nil
	}

	return p.Members[index]
}

func (p *Parliament) On(parliamentEvent midgard.Event) error {
	switch v := parliamentEvent.(type) {

	case *event.LeaderChanged:
		p.Leader = &Leader{
			LeaderId: LeaderId{v.LeaderId},
		}

	case *event.MemberJoined:
		p.Members = append(p.Members, &Member{
			MemberId: MemberId{v.MemberId},
		})

	case *event.MemberRemoved:
		index := p.findIndexOfMember(v.MemberId)

		if index != -1 {
			p.Members = append(p.Members[:index], p.Members[index+1:]...)
		}

	default:
		return errors.New(fmt.Sprintf("unhandled event [%s]", v))
	}

	return nil
}
