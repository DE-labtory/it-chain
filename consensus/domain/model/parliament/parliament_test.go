package parliament

import (
	"testing"

	"github.com/it-chain/it-chain-Engine/consensus/domain/model/consensus"
	"github.com/stretchr/testify/assert"
)

func TestParliament_NewParliament(t *testing.T) {
	p := NewParliament()

	assert.Nil(t, nil, p.Leader)
	assert.Equal(t, 0, len(p.Members))
}

func TestParliament_HasLeader(t *testing.T) {
	//given
	p := NewParliament()

	//when
	flag := p.HasLeader()

	//then
	assert.False(t, flag)

	//given
	p.Leader = &Leader{}

	//when
	flag2 := p.HasLeader()

	//then
	assert.True(t, flag2)
}

func TestParliament_AddMember(t *testing.T) {

	//given
	p := NewParliament()
	m := &Member{ID: PeerID{"member1"}}

	//when
	err := p.AddMember(m)

	//then
	assert.Nil(t, err)

	//when
	err = p.AddMember(&Member{})

	//then
	assert.Error(t, err)
}

func TestParliament_IsNeedConsensus(t *testing.T) {

	//멤버가 없을 경우
	//given
	p := NewParliament()

	//when
	flag := p.IsNeedConsensus()

	//then
	assert.False(t, flag)

	//멤버가 존재할경우
	//given
	err := p.AddMember(&Member{ID: PeerID{"mem1"}})
	assert.NoError(t, err)

	//when
	flag = p.IsNeedConsensus()

	//then
	assert.True(t, flag)
}

func TestParliament_RemoveMember(t *testing.T) {

	//given
	p := NewParliament()
	m := &Member{ID: PeerID{"member1"}}
	err := p.AddMember(m)
	assert.Nil(t, err)

	//when
	p.RemoveMember(PeerID{"member1"})

	//then
	assert.Equal(t, 0, len(p.Members))

	//when
	p.RemoveMember(PeerID{"member2"})

	//then
	assert.Equal(t, 0, len(p.Members))
}

func TestParliament_SetLeader(t *testing.T) {
	//given
	p := NewParliament()
	assert.Nil(t, p.Leader)
	l := &Leader{ID: PeerID{"member1"}}

	//when
	p.SetLeader(l)

	//then
	assert.NotNil(t, p.Leader)
	assert.Equal(t, PeerID{"member1"}, p.Leader.ID)
}

func TestParliament_ValidateRepresentative(t *testing.T) {

	//given
	p := NewParliament()
	m := &Member{ID: PeerID{"member1"}}
	m2 := &Member{ID: PeerID{"member2"}}
	err := p.AddMember(m)
	assert.Nil(t, err)
	err = p.AddMember(m2)
	assert.Nil(t, err)

	r1 := &consensus.Representative{Id: consensus.RepresentativeId(m.ID.ID)}
	r2 := &consensus.Representative{Id: consensus.RepresentativeId(m.ID.ID)}
	r3 := &consensus.Representative{Id: consensus.RepresentativeId("invalidMember")}

	representatives := []*consensus.Representative{r1, r2}

	//when
	flag := p.ValidateRepresentative(representatives)

	//then
	assert.True(t, flag)

	//given
	representatives = []*consensus.Representative{r1, r2, r3}

	//when
	flag = p.ValidateRepresentative(representatives)

	//then
	assert.False(t, flag)
}
