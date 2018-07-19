package consensus

import (
	"testing"

	"github.com/it-chain/engine/consensus/test/mock"
	"github.com/it-chain/engine/core/eventstore"
	"github.com/it-chain/midgard"
	"github.com/stretchr/testify/assert"
)

func TestNewParliament(t *testing.T) {
	p := NewParliament()

	assert.Equal(t, "0", p.GetID())
	assert.Equal(t, 0, len(p.Members))
	assert.Nil(t, p.Leader)
}

func TestParliament_IsNeedConsensus(t *testing.T) {
	// case 1
	p := NewParliament()

	flag := p.IsNeedConsensus()

	assert.False(t, flag)

	// case 2
	p.Members = append(p.Members, &Member{})

	flag = p.IsNeedConsensus()

	assert.True(t, flag)
}

func TestParliament_HasLeader(t *testing.T) {
	// case 1
	p := NewParliament()

	flag := p.HasLeader()

	assert.False(t, flag)

	// case2
	p.Leader = &Leader{}

	flag = p.HasLeader()

	assert.True(t, flag)
}

func TestParliament_ChangeLeader(t *testing.T) {
	// given
	p := NewParliament()
	l := &Leader{LeaderId: LeaderId{"leader"}}
	eventRepository := mock.MockEventRepository{}

	// when
	eventRepository.SaveFunc = func(aggregateID string, events ...midgard.Event) error {
		assert.Equal(t, "leader", events[0].(LeaderChangedEvent).LeaderId)
		return nil
	}
	eventstore.InitForMock(eventRepository)

	err := p.ChangeLeader(l)

	// then
	assert.Nil(t, err)
	assert.Equal(t, l.GetID(), p.Leader.GetID())
}

func TestParliament_AddMember(t *testing.T) {
	// given
	p := NewParliament()
	m := &Member{MemberId: MemberId{"member"}}
	eventRepository := mock.MockEventRepository{}

	// when
	eventRepository.SaveFunc = func(aggregateID string, events ...midgard.Event) error {
		assert.Equal(t, "member", events[0].(MemberJoinedEvent).MemberId)
		return nil
	}
	eventstore.InitForMock(eventRepository)

	err := p.AddMember(m)

	// then
	assert.Nil(t, err)
	assert.Equal(t, 1, len(p.Members))
}

func TestParliament_RemoveMember(t *testing.T) {
	// given
	p := NewParliament()
	m := &Member{MemberId: MemberId{"member"}}
	eventRepository := mock.MockEventRepository{}
	eventRepository.SaveFunc = func(aggregateID string, events ...midgard.Event) error {
		return nil
	}
	eventstore.InitForMock(eventRepository)
	p.AddMember(m)

	// case 1
	eventRepository.SaveFunc = func(aggregateID string, events ...midgard.Event) error {
		assert.NotEqual(t, m.MemberId.Id, events[0].(MemberRemovedEvent).MemberId)
		return nil
	}
	eventstore.InitForMock(eventRepository)

	err := p.RemoveMember(MemberId{"nonmember"})

	// then
	assert.Nil(t, err)
	assert.Equal(t, 1, len(p.Members))

	// case2
	eventRepository.SaveFunc = func(aggregateID string, events ...midgard.Event) error {
		assert.Equal(t, m.MemberId.Id, events[0].(MemberRemovedEvent).MemberId)
		return nil
	}
	eventstore.InitForMock(eventRepository)

	err = p.RemoveMember(m.MemberId)

	// then
	assert.Nil(t, err)
	assert.Equal(t, 0, len(p.Members))
}

func TestParliament_ValidateRepresentative(t *testing.T) {
	p := NewParliament()
	eventRepository := mock.MockEventRepository{}
	eventRepository.SaveFunc = func(aggregateID string, events ...midgard.Event) error {
		return nil
	}
	eventstore.InitForMock(eventRepository)

	// case 1
	var representatives1 []*Representative
	for i := 0; i < 3; i++ {
		p.AddMember(&Member{
			MemberId: MemberId{string(i)},
		})

		representatives1 = append(representatives1, &Representative{
			Id: RepresentativeId(string(i)),
		})
	}

	flag := p.ValidateRepresentative(representatives1)

	assert.True(t, flag)

	// case 2
	var representatives2 []*Representative
	for i := 3; i < 6; i++ {
		representatives2 = append(representatives2, &Representative{
			Id: RepresentativeId(string(i)),
		})
	}

	flag = p.ValidateRepresentative(representatives2)

	assert.False(t, flag)
}

func TestParliament_FindByPeerID(t *testing.T) {
	p := NewParliament()
	m := &Member{MemberId: MemberId{"member"}}
	eventRepository := mock.MockEventRepository{}
	eventRepository.SaveFunc = func(aggregateID string, events ...midgard.Event) error {
		return nil
	}
	eventstore.InitForMock(eventRepository)
	p.AddMember(m)

	// case 1
	member := p.FindByPeerID("member")

	assert.Equal(t, "member", member.GetID())

	// case 2
	member = p.FindByPeerID("nonmember")

	assert.Nil(t, member)
}
