package consensus

import (
	"testing"

	"github.com/it-chain/engine/common/event"
	"github.com/it-chain/engine/consensus/test/mock"
	"github.com/it-chain/engine/core/eventstore"
	"github.com/it-chain/midgard"
	"github.com/stretchr/testify/assert"
)

func TestCreateConsensus(t *testing.T) {
	// given
	p := make([]MemberId, 0)
	l := MemberId("leader")
	m := MemberId("member")
	b := ProposedBlock{
		Seal: make([]byte, 0),
		Body: make([]byte, 0),
	}

	eventRepository := mock.MockEventRepository{}
	eventRepository.SaveFunc = func(aggregateID string, events ...midgard.Event) error {
		return nil
	}
	eventstore.InitForMock(eventRepository)

	p = append(p, l)
	p = append(p, m)

	// when
	eventRepository.SaveFunc = func(aggregateID string, events ...midgard.Event) error {
		assert.Equal(t, len(p), len(events[0].(*event.ConsensusCreated).Representatives))
		assert.NotNil(t, events[0].(*event.ConsensusCreated).Seal)
		return nil
	}
	eventstore.InitForMock(eventRepository)

	c, err := CreateConsensus(p, b)

	// then
	assert.Nil(t, err)
	assert.Equal(t, 2, len(c.Representatives))
	assert.Equal(t, b.Body, c.Block.Body)
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
			Body: make([]byte, 0),
		},
	}

	eventRepository := mock.MockEventRepository{}

	// when
	eventRepository.SaveFunc = func(aggregateID string, events ...midgard.Event) error {
		assert.Equal(t, len(r), len(events[0].(*event.ConsensusCreated).Representatives))
		assert.Equal(t, "consensusID", events[0].(*event.ConsensusCreated).GetID())
		return nil
	}
	eventstore.InitForMock(eventRepository)

	c, err := ConstructConsensus(msg)

	// then
	assert.Nil(t, err)
	assert.Equal(t, "consensusID", c.ConsensusID.Id)
	assert.Equal(t, IDLE_STATE, c.CurrentState)
	assert.Equal(t, 2, len(c.Representatives))
}
