package mock

import (
	"github.com/DE-labtory/engine/common"
	"github.com/DE-labtory/engine/consensus/pbft"
)

type ParliamentApi struct {
	AddRepresentativeFunc    func(representativeId string)
	RemoveRepresentativeFunc func(representativeId string)
	UpdateLeaderFunc         func(nodeId common.NodeID) error
	GetLeaderFunc            func() pbft.Leader
	RequestLeaderFunc        func(connectionId string)
	DeliverLeaderFunc        func(connectionId string)
}

func (p *ParliamentApi) AddRepresentative(representativeId string) {
	p.AddRepresentativeFunc(representativeId)
}
func (p *ParliamentApi) RemoveRepresentative(representativeId string) {
	p.RemoveRepresentativeFunc(representativeId)
}
func (p *ParliamentApi) UpdateLeader(nodeId common.NodeID) error {
	return p.UpdateLeaderFunc(nodeId)
}
func (p *ParliamentApi) GetLeader() pbft.Leader {
	return p.GetLeaderFunc()
}
func (p *ParliamentApi) RequestLeader(connectionId string) {
	p.RequestLeaderFunc(connectionId)
}
func (p *ParliamentApi) DeliverLeader(connectionId string) {
	p.DeliverLeaderFunc(connectionId)
}
