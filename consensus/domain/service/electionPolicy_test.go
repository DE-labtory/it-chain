package service

import (
	"testing"
	"github.com/it-chain/it-chain-Engine/consensus/domain/model/parliament"
	"github.com/stretchr/testify/assert"
)

func TestElect(t *testing.T) {

	p := parliament.Parliament{}

	p.Leader = &parliament.Leader{parliament.PeerID{"leader1"}}
	p.Members = make([]*parliament.Member,0)
	p.Members = append(p.Members, &parliament.Member{parliament.PeerID{"member1"}})


	representatives,err := Elect(p)

	if err != nil{
		assert.Fail(t,err.Error())
	}

	assert.Equal(t,len(representatives),2)

	p.Members = append(p.Members, &parliament.Member{parliament.PeerID{"member2"}})
	representatives,err = Elect(p)

	if err != nil{
		assert.Fail(t,err.Error())
	}

	assert.Equal(t,len(representatives),3)
}