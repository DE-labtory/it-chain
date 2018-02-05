package service

import (
	"testing"
	"it-chain/network/comm"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/mock"
	"it-chain/auth"
	"it-chain/domain"
)

type MockConnectionManager struct{
	mock.Mock
}

func (mcm MockConnectionManager) SendStream(data interface{}, errorCallBack comm.OnError, connectionID string){
	//mcm.Called(data)
	//mcm.Called(data,errorCallBack,connectionID)
	//mcm.Called(data,nil,connectionID)
	mcm.MethodCalled("SendStream",data,nil,connectionID)
}

func (mcm MockConnectionManager) Stop(){

}

func (mcm MockConnectionManager) Close(connectionID string){

}

func (mcm MockConnectionManager) CreateStreamConn(connectionID string, ip string, handle comm.ReceiveMessageHandle) error{
	return errors.New("error")
}

func (mcm MockConnectionManager) Size() int{
	return 0
}

type MockCrypto struct{
	mock.Mock
}

func (mc MockCrypto) Sign(digest []byte, opts auth.SignerOpts) (signature []byte, err error){
	return []byte("ASD"), nil
}

func (mc MockCrypto) Verify(key auth.Key, signature, digest []byte, opts auth.SignerOpts) (valid bool, err error){
	return true, nil
}

func (mc MockCrypto) GenerateKey(opts auth.KeyGenOpts) (pri, pub auth.Key, err error){
	return nil, nil, errors.New("asd")
}

func (mc MockCrypto) LoadKey() (pri, pub auth.Key, err error){
	return nil,nil,errors.New("asd")
}

type MockPeerService struct{

}

//peer table 조회
func (mps MockPeerService) GetPeerTable() domain.PeerTable{
	myPeer := &domain.PeerInfo{}
	myPeer.PeerID = "peer1"

	peer1 := &domain.PeerInfo{}
	peer1.PeerID = "peer2"

	peerTable := domain.PeerTable{}
	peerTable.OwnerID = myPeer.PeerID
	peerTable.PeerMap = make(map[string]*domain.PeerInfo)
	peerTable.PeerMap[myPeer.PeerID] = myPeer
	peerTable.PeerMap[peer1.PeerID] = peer1


	return peerTable
}

//peer info 찾기
func (mps MockPeerService) GetPeerInfoByPeerID(peerID string) *domain.PeerInfo{
	return &domain.PeerInfo{}
}

//peer info
func (mps MockPeerService) PushPeerTable(peerIDs []string){

}

//update peerTable
func (mps MockPeerService) UpdatePeerTable(peerTable domain.PeerTable){

}

//Add peer
func (mps MockPeerService) AddPeerInfo(peerInfo *domain.PeerInfo){

}

//Request Peer Info
func (mps MockPeerService) RequestPeerInfo(host string, port string) *domain.PeerInfo{
	return &domain.PeerInfo{}
}

func (mps MockPeerService) BroadCastPeerTable(interface{}){

}

func TestNewPBFTConsensusService(t *testing.T) {
	comm:= new(MockConnectionManager)
	peerService := new(MockPeerService)

	pbftService := NewPBFTConsensusService(comm,peerService)

	block := &domain.Block{}

	comm.On("SendStream", mock.AnythingOfType("ConsensusMessage"), nil, "peer2")

	pbftService.StartConsensus(block)

	comm.AssertExpectations(t)
	//comm.AssertNumberOfCalls(t,"SendStream",1)
	//comm.AssertNumberOfCalls(t,"SendStream",2)
}

func TestPBFTConsensusService_broadcastMessage(t *testing.T) {

}

func TestPBFTConsensusService_StartConsensus(t *testing.T) {

}
