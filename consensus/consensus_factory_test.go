package consensus

import (
	"testing"
	"github.com/it-chain/engine/consensus/test/mock"
	"github.com/it-chain/midgard"
	"github.com/it-chain/engine/core/eventstore"
	"github.com/stretchr/testify/assert"
)

func TestCreateConsensus(t *testing.T) {
	// given
	p := NewParliament()
	l := &Leader{LeaderId: LeaderId{"leader"}}
	m := &Member{MemberId: MemberId{"member"},}
	b := ProposedBlock{
		Seal: make([]byte, 0),
		body: make([]byte, 0),
	}

	eventRepository := mock.MockEventRepository{}
	eventRepository.SaveFunc = func(aggregateID string, events ...midgard.Event) error {
		return nil
	}
	eventstore.InitForMock(eventRepository)

	p.ChangeLeader(l)
	p.AddMember(m)

	// when
	eventRepository.SaveFunc = func(aggregateID string, events ...midgard.Event) error {
		assert.Equal(t, 1+len(p.Members), len(events[0].(ConsensusCreatedEvent).Consensus.Representatives))
		assert.NotNil(t, events[0].(ConsensusCreatedEvent).Consensus.Block.Seal)
		return nil
	}
	eventstore.InitForMock(eventRepository)

	c, err := CreateConsensus(p, b)

	// then
	assert.Nil(t, err)
	assert.Equal(t, 1+len(p.Members), len(c.Representatives))
	assert.Equal(t, b.body, c.Block.body)
}

func TestConstructConsensus(t *testing.T) {
	// given
	l := NewRepresentative("leader")
	m := NewRepresentative("member")

	r := make([]*Representative, 0)
	r = append(r, l, m)

	msg := PrePrepareMsg{
		ConsensusId:    NewConsensusId("consensusID"),
		SenderId:       "me",
		Representative: r,
		ProposedBlock: ProposedBlock{
			Seal: make([]byte, 0),
			body: make([]byte, 0),
		},
	}

	eventRepository := mock.MockEventRepository{}

	// when
	eventRepository.SaveFunc = func(aggregateID string, events ...midgard.Event) error {
		assert.Equal(t, len(r), len(events[0].(ConsensusCreatedEvent).Consensus.Representatives))
		assert.Equal(t, "consensusID", events[0].(ConsensusCreatedEvent).GetID())
		return nil
	}
	eventstore.InitForMock(eventRepository)

	c, err := ConstructConsensus(msg)

	// then
	assert.Nil(t, err)
	assert.Equal(t, "consensusID", c.ConsensusID.Id)
	assert.Equal(t, PREPREPARE_STATE, c.CurrentState)
	assert.Equal(t, 2, len(c.Representatives))
}
