package api

import (
	"testing"

	"github.com/it-chain/it-chain-Engine/consensus/legacy/domain/model/parliament"
	"github.com/it-chain/it-chain-Engine/consensus/legacy/domain/repository"
	"github.com/stretchr/testify/assert"
)

func TestParliamentApi_AddMember(t *testing.T) {

	//given
	pr := repository.NewPaliamentRepository()
	pApi := NewParliamentApi(pr)

	m := parliament.Member{ID: parliament.PeerID{ID: "mem1"}}

	//when
	err := pApi.AddMember(m)
	assert.NoError(t, err)

	//then
	p := pr.Get()
	findMember := p.FindByPeerID("mem1")
	assert.Equal(t, m.ID, findMember.ID)
}

func TestParliamentApi_ChangeLeader(t *testing.T) {

	//given
	pr := repository.NewPaliamentRepository()
	pApi := NewParliamentApi(pr)

	//when
	err := pApi.ChangeLeader(parliament.Leader{ID: parliament.PeerID{"leader1"}})

	//then
	assert.NoError(t, err)
	p := pr.Get()
	leader := p.GetLeader()
	assert.NotNil(t, leader)
	assert.Equal(t, parliament.PeerID{"leader1"}, leader.ID)
}

func TestParliamentApi_RemoveMember(t *testing.T) {

	//given
	pr := repository.NewPaliamentRepository()
	pApi := NewParliamentApi(pr)
	m := parliament.Member{ID: parliament.PeerID{ID: "mem1"}}
	err := pApi.AddMember(m)
	assert.NoError(t, err)

	//when
	err = pApi.RemoveMember(m.ID)

	//then
	assert.NoError(t, err)
	p := pr.Get()
	findMember := p.FindByPeerID("mem1")
	assert.Nil(t, findMember)
}
