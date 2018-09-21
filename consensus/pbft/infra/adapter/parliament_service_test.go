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

package adapter_test

import (
	"testing"

	"github.com/it-chain/engine/api_gateway"
	"github.com/it-chain/engine/consensus/pbft"
	"github.com/it-chain/engine/consensus/pbft/infra/adapter"
	"github.com/it-chain/engine/p2p"
	"github.com/it-chain/engine/p2p/infra/mem"
	"github.com/stretchr/testify/assert"
)

func TestParliamentService_RequestLeader(t *testing.T) {
	// given (case 1 : no leader)
	peerRepository := mem.NewPeerReopository()

	peerRepository.Save(p2p.Peer{
		IpAddress: "1.1.1.1",
		PeerId:    p2p.PeerId{"p1"},
	})
	parliament := pbft.NewParliament()

	ps := adapter.NewParliamentService(parliament, api_gateway.NewPeerQueryApi(peerRepository))

	// when
	l, _ := ps.RequestLeader()

	// then
	assert.Equal(t, "", l.ToString())

	// given (case 2 : good case)
	peerRepository.SetLeader(p2p.Leader{
		LeaderId: p2p.LeaderId{Id: "leader"},
	})

	// when
	l, err := ps.RequestLeader()

	// then
	assert.Equal(t, "leader", l.ToString())
	assert.Nil(t, err)
}

func TestParliamentService_RequestPeerList(t *testing.T) {
	// given
	peerRepository := mem.NewPeerReopository()

	p1 := p2p.Peer{
		IpAddress: "1.1.1.1",
		PeerId:    p2p.PeerId{"p1"},
	}

	p2 := p2p.Peer{
		IpAddress: "2.2.2.2",
		PeerId:    p2p.PeerId{"p2"},
	}

	peerRepository.Save(p1)
	peerRepository.Save(p2)

	parliament := pbft.NewParliament()

	ps := adapter.NewParliamentService(parliament, api_gateway.NewPeerQueryApi(peerRepository))

	// when
	peerList, err := ps.RequestPeerList()

	// then
	assert.Equal(t, 2, len(peerList))
	assert.Nil(t, err)
}

func TestParliamentService_IsNeedConsensus(t *testing.T) {
	// given (case 1 : no member)
	peerRepository := mem.NewPeerReopository()
	parliament := pbft.NewParliament()
	ps := adapter.NewParliamentService(parliament, api_gateway.NewPeerQueryApi(peerRepository))

	// when
	flag := ps.IsNeedConsensus()

	// then
	assert.Equal(t, false, flag)

	// given (case 2 : less than 4 members)
	p1 := p2p.Peer{
		IpAddress: "1.1.1.1",
		PeerId:    p2p.PeerId{"p1"},
	}

	p2 := p2p.Peer{
		IpAddress: "2.2.2.2",
		PeerId:    p2p.PeerId{"p2"},
	}

	p3 := p2p.Peer{
		IpAddress: "3.3.3.3",
		PeerId:    p2p.PeerId{"p3"},
	}

	peerRepository.Save(p1)
	peerRepository.Save(p2)
	peerRepository.Save(p3)

	// when
	flag = ps.IsNeedConsensus()

	// then
	assert.Equal(t, false, flag)

	// given (case 3 : equal or moro than 4 members)
	p4 := p2p.Peer{
		IpAddress: "4.4.4.4",
		PeerId:    p2p.PeerId{"p4"},
	}

	peerRepository.Save(p4)

	// when
	flag = ps.IsNeedConsensus()

	// then
	assert.Equal(t, true, flag)
}

func TestParliamentService_Build(t *testing.T) {
	p := SetParliamentService()

	p.Build()

	assert.Equal(t, p.GetRepresentativeById("1").GetID(), "1")
	assert.Equal(t, p.GetRepresentativeById("2").GetID(), "2")
	assert.Equal(t, p.GetRepresentativeById("3").GetID(), "3")
}

func TestParliamentService_FindRepresentativeByIpAddress(t *testing.T) {
	p := SetParliamentService()

	rep := p.FindRepresentativeByIpAddress("1")

	assert.Equal(t, rep.IpAddress, "1")
}

func TestParliamentService_GetLeader(t *testing.T) {
	p := SetParliamentService()

	p.SetLeader(&pbft.Representative{
		IpAddress: "1",
		ID:        "1",
	})

	assert.Equal(t, p.GetLeader().LeaderId, "1")
}

func TestParliamentService_GetParliament(t *testing.T) {
	p := SetParliamentService()

	parliament := p.GetParliament()

	assert.Equal(t, parliament.RepresentativeTable["1"].ID, "1")
}

func TestParliamentService_GetRepresentativeById(t *testing.T) {
	p := SetParliamentService()

	rep := p.GetRepresentativeById("1")

	assert.Equal(t, rep.ID, "1")
}

func TestParliamentService_GetRepresentativeTable(t *testing.T) {
	p := SetParliamentService()

	table := p.GetRepresentativeTable()

	assert.Equal(t, table["1"].ID, "1")
}

func TestParliamentService_SetLeader(t *testing.T) {
	p := SetParliamentService()

	p.SetLeader(&pbft.Representative{
		ID:        "1",
		IpAddress: "1",
	})

	assert.Equal(t, p.GetLeader().LeaderId, "1")
}

func TestNewParliamentService(t *testing.T) {
	p := SetParliamentService()

	assert.Equal(t, p.GetRepresentativeById("1").IpAddress, "1")
}

func SetParliamentService() *adapter.ParliamentService {
	parliament := pbft.NewParliament()
	repository := mem.NewPeerReopository()
	repository.Save(p2p.Peer{
		IpAddress: "1",
		PeerId: p2p.PeerId{
			Id: "1",
		},
	})
	repository.Save(p2p.Peer{
		IpAddress: "2",
		PeerId: p2p.PeerId{
			Id: "2",
		},
	})
	repository.Save(p2p.Peer{
		IpAddress: "3",
		PeerId: p2p.PeerId{
			Id: "3",
		},
	})
	api := api_gateway.NewPeerQueryApi(repository)
	p := adapter.NewParliamentService(parliament, api)
	p.Build()

	return p
}
