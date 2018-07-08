package p2p

type MockPeerService struct {

	SaveFunc func(peer Peer) error
	RemoveFunc func(peerId PeerId) error
	FindByIdFunc func(peerId PeerId) (Peer, error)
	FindByAddressFunc func(ipAddress string) (Peer, error)
	FindAllFunc func() ([]Peer, error)
}

func (mps *MockPeerService) Save (peer Peer) error{

	return mps.SaveFunc(peer)
}

func (mps *MockPeerService) Remove (peerId PeerId) error{

	return mps.RemoveFunc(peerId)
}

func (mps *MockPeerService) FindById (peerId PeerId) (Peer, error){

	return mps.FindByIdFunc(peerId)
}

func (mps *MockPeerService) FindByAddress (ipAddress string) (Peer, error){

	return mps.FindByAddressFunc(ipAddress)
}

func (mps *MockPeerService) FindAll () ([]Peer, error){

	return mps.FindAllFunc()
}

type MockCommunicationService struct {

	DialFunc func(ipAddress string) error
	DeliverPLTableFunc func(connectionId string, pLTable PLTable) error
}

func (mcs *MockCommunicationService) Dial(ipAddress string) error {

	return mcs.DialFunc(ipAddress)
}
func (mcs *MockCommunicationService) DeliverPLTable(connectionId string, pLTable PLTable) error{

	return mcs.DeliverPLTableFunc(connectionId, pLTable)
}