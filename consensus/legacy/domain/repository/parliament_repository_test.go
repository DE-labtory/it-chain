package repository

import (
	"testing"

	"fmt"

	"github.com/it-chain/it-chain-Engine/consensus/legacy/domain/model/parliament"
	"github.com/stretchr/testify/assert"
)

func TestParlimentRepository_impl_Get(t *testing.T) {
	pr := NewPaliamentRepository()

	p := pr.Get()

	assert.Nil(t, nil, p.Leader)
	assert.Equal(t, 0, len(p.Members))
}

func TestParlimentRepository_impl_Save(t *testing.T) {

	pr := NewPaliamentRepository()

	p := pr.Get()
	m := &parliament.Member{ID: parliament.PeerID{"mem1"}}
	p.AddMember(m)

	err := pr.Save(p)

	if err != nil {

	}

	p2 := pr.Get()

	fmt.Println(p2.Members[0].ID)
	assert.Equal(t, p.Members, p2.Members)

}
