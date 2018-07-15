package mock

import "github.com/it-chain/it-chain-Engine/p2p"

type MockPLTableApi struct {

	getPLTableFunc func() p2p.PLTable
}


func (mplta *MockPLTableApi) GetPLTable() p2p.PLTable {

	return mplta.getPLTableFunc()
}

