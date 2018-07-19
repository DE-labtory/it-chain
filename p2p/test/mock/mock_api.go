package mock

import "github.com/it-chain/engine/p2p"

type MockPLTableApi struct {
	getPLTableFunc func() p2p.PLTable
}

func (mplta *MockPLTableApi) GetPLTable() p2p.PLTable {

	return mplta.getPLTableFunc()
}

type MockLeaderApi struct {
}

func (mla *MockLeaderApi) UpdateLeaderWithAddress(ipAddress string) error {
	return mla.UpdateLeaderWithAddress(ipAddress)
}

func (mla *MockLeaderApi) UpdateLeaderWithLargePeerTable(oppositePLTable p2p.PLTable) error {
	return mla.UpdateLeaderWithLargePeerTable(oppositePLTable)
}

type MockCommunicationApi struct {
	DeliverPLTableFunc       func(connectionId string) error
	DialToUnConnectedNodeFuc func(peerTable map[string]p2p.Peer) error
}

func (mca *MockCommunicationApi) DeliverPLTable(connectionId string) error {

	return mca.DeliverPLTableFunc(connectionId)
}

func (mca *MockCommunicationApi) DialToUnConnectedNode(peerTable map[string]p2p.Peer) error {

	return mca.DialToUnConnectedNodeFuc(peerTable)
}
