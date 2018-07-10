package mock

import (
	"github.com/it-chain/it-chain-Engine/p2p"
)

type MockPeerService struct {

	SaveFunc func(peer p2p.Peer) error
	RemoveFunc func(peerId p2p.PeerId) error
	FindByIdFunc func(peerId p2p.PeerId) (p2p.Peer, error)
	FindByAddressFunc func(ipAddress string) (p2p.Peer, error)
	FindAllFunc func() ([]p2p.Peer, error)
}

func (mps *MockPeerService) Save (peer p2p.Peer) error{

	return mps.SaveFunc(peer)
}

func (mps *MockPeerService) Remove (peerId p2p.PeerId) error{

	return mps.RemoveFunc(peerId)
}

func (mps *MockPeerService) FindById (peerId p2p.PeerId) (p2p.Peer, error){

	return mps.FindByIdFunc(peerId)
}

func (mps *MockPeerService) FindByAddress (ipAddress string) (p2p.Peer, error){

	return mps.FindByAddressFunc(ipAddress)
}

func (mps *MockPeerService) FindAll () ([]p2p.Peer, error){

	return mps.FindAllFunc()
}

type MockLeaderService struct {
	SetFunc func(leader p2p.Leader) error
}

func (ls *MockLeaderService) Set(leader p2p.Leader) error {

	return ls.SetFunc(leader)
}

type MockCommunicationService struct {

	DialFunc func(ipAddress string) error
	DeliverPLTableFunc func(connectionId string, pLTable p2p.PLTable) error
}

func (mcs *MockCommunicationService) Dial(ipAddress string) error {

	return mcs.DialFunc(ipAddress)
}

func (mcs *MockCommunicationService) DeliverPLTable(connectionId string, pLTable p2p.PLTable) error{

	return mcs.DeliverPLTableFunc(connectionId, pLTable)
}

type MockPeerQueryService struct{

	FindByIdFunc func(peerId p2p.PeerId) (p2p.Peer, error)
	FindAllFunc func() ([]p2p.Peer, error)
	FindByAddressFunc func(ipAddress string) (p2p.Peer, error)
}

func (mpqs *MockPeerQueryService) FindById(peerId p2p.PeerId) (p2p.Peer, error){

	return mpqs.FindByIdFunc(peerId)
}

func (mpqs *MockPeerQueryService) FindAll() ([]p2p.Peer, error){

	return mpqs.FindAll()
}

func (mpqs *MockPeerQueryService) FindByAddress(ipAddress string) (p2p.Peer, error){

	return mpqs.FindByAddress(ipAddress)
}

type MockPLTableQueryService struct{
	GetPLTableFunc func() (p2p.PLTable, error)
}

func (mpltqs *MockPLTableQueryService) GetPLTable() (p2p.PLTable, error){
	return mpltqs.GetPLTableFunc()
}
