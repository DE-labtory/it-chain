package service

import (
	"testing"
	"it-chain/network/comm"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/mock"
	"it-chain/auth"
	"it-chain/domain"
	"github.com/stretchr/testify/assert"
	pb "it-chain/network/protos"
	"time"
	"fmt"
	"sync"
)

type MockConnectionManager struct{
	mock.Mock
}

func (mcm MockConnectionManager) SendStream(data *pb.StreamMessage, errorCallBack comm.OnError, connectionID string){
	mcm.MethodCalled("SendStream",data,nil,connectionID)
}

func (mcm MockConnectionManager) Stop(){

}

func (mcm MockConnectionManager) Close(connectionID string){

}

func (mcm MockConnectionManager) CreateStreamClientConn(connectionID string, ip string, handle comm.ReceiveMessageHandle) error{
	return errors.New("error")
}

func (mcm MockConnectionManager) SetOnConnectHandler(onConnectionHandler comm.OnConnectionHandler){

}

func (mcm MockConnectionManager) Size() int{
	return 0
}


func (mcm MockConnectionManager) Stream(stream pb.StreamService_StreamServer) error{
	return nil
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
	myPeer := &domain.Peer{}
	myPeer.PeerID = "peer1"

	peer1 := &domain.Peer{}
	peer1.PeerID = "peer2"

	peerTable := domain.PeerTable{}
	peerTable.MyID = myPeer.PeerID
	peerTable.PeerMap = make(map[string]*domain.Peer)
	peerTable.PeerMap[myPeer.PeerID] = myPeer
	peerTable.PeerMap[peer1.PeerID] = peer1


	return peerTable
}

//peer info 찾기
func (mps MockPeerService) GetPeerByPeerID(peerID string) *domain.Peer{
	return &domain.Peer{}
}

//peer info
func (mps MockPeerService) PushPeerTable(peerIDs []string){

}

//update peerTable
func (mps MockPeerService) UpdatePeerTable(peerTable domain.PeerTable){

}

//Add peer
func (mps MockPeerService) AddPeer(Peer *domain.Peer){

}

//Request Peer Info
func (mps MockPeerService) RequestPeer(host string, port string) *domain.Peer{
	return &domain.Peer{}
}

func (mps MockPeerService) BroadCastPeerTable(interface{}){

}


func TestNewPBFTConsensusService(t *testing.T) {

	comm:= new(MockConnectionManager)
	view := &domain.View{}
	view.ID = "123"
	view.LeaderID = "1"
	view.PeerID = []string{"1","2","3"}

	pbftService := NewPBFTConsensusService(comm,nil,&domain.Peer{PeerID:"1"},nil)

	consensusStates := pbftService.GetCurrentConsensusState()
	assert.NotNil(t,consensusStates)
}

//todo assertnumberofcall 테스트 추가해야함
func TestPBFTConsensusService_StartConsensus(t *testing.T) {

	connctionManager:= new(MockConnectionManager)
	view := &domain.View{}
	view.ID = "123"
	view.LeaderID = "1"
	view.PeerID = []string{"1","2"}

	pbftService := NewPBFTConsensusService(connctionManager,nil,&domain.Peer{PeerID:"1"},nil)
	block := &domain.Block{}

	connctionManager.On("SendStream", mock.AnythingOfType("*message.StreamMessage"), nil, "1")

	wg := sync.WaitGroup{}
	wg.Add(2)

	pbftService.StartConsensus(view,block)

	consensusStates := pbftService.GetCurrentConsensusState()
	assert.Equal(t,len(consensusStates),1)
	//
	for _, state := range consensusStates{
		assert.Equal(t,state.Block,block)
		assert.Equal(t,state.CurrentStage,domain.Prepared)
	}
	//

	wg.Wait()
	//connctionManager.AssertExpectations(t)

}

func GetMockConsensusMessage(consensusID string, msgType domain.MsgType) *pb.StreamMessage{

	consensusMessage := &pb.ConsensusMessage{}
	consensusMessage.SequenceID = time.Now().UnixNano()
	consensusMessage.ConsensusID = consensusID
	consensusMessage.MsgType = int32(msgType)
	//consensusMessage.Block = &pb.Block{}

	view := pb.View{}
	view.ViewID = "123"
	view.LeaderID = "1"
	view.PeerID = []string{"1","2","3"}

	consensusMessage.View = &view

	message := &pb.StreamMessage{}
	cm := &pb.StreamMessage_ConsensusMessage{}
	cm.ConsensusMessage = consensusMessage
	message.Content = cm

	return message
}

func TestPBFTConsensusService_ReceiveConsensusMessage(t *testing.T) {

	//when
	connctionManager:= new(MockConnectionManager)
	view := &domain.View{}
	view.ID = "123"
	view.LeaderID = "1"
	view.PeerID = []string{"1","2","3"}

	pbftService := NewPBFTConsensusService(connctionManager,nil,&domain.Peer{PeerID:"1"},nil)

	Message := GetMockConsensusMessage("1",domain.PreprepareMsg)

	outMessage := comm.OutterMessage{}
	outMessage.Message = Message

	//then
	pbftService.ReceiveConsensusMessage(outMessage)


	//result WhenNoStateExist
	addedState := pbftService.GetCurrentConsensusState()["1"]
	fmt.Println(pbftService.GetCurrentConsensusState())
	fmt.Println(Message.GetConsensusMessage())
	assert.Equal(t,addedState.ID,Message.GetConsensusMessage().ConsensusID)
	assert.Equal(t,addedState.ID,Message.GetConsensusMessage().ConsensusID)

	//2 whenStateEx
	Message2 := GetMockConsensusMessage("1",domain.CommitMsg)
	outMessage2 := comm.OutterMessage{}
	outMessage2.Message = Message2
	fmt.Println(pbftService.GetCurrentConsensusState())
	fmt.Println(Message.GetConsensusMessage())

	pbftService.ReceiveConsensusMessage(outMessage2)


	assert.Equal(t,len(pbftService.GetCurrentConsensusState()),1)

	//3 multiState
	Message3 := GetMockConsensusMessage("2",domain.CommitMsg)
	outMessage3 := comm.OutterMessage{}
	outMessage3.Message = Message3

	pbftService.ReceiveConsensusMessage(outMessage3)
	fmt.Println(pbftService.GetCurrentConsensusState())

	assert.Equal(t,len(pbftService.GetCurrentConsensusState()),2)

	//4 timeouted Message
	Message4 := GetMockConsensusMessage("3",domain.PreprepareMsg)
	Message4.GetConsensusMessage().SequenceID = time.Now().Local().Add(10*time.Minute).UnixNano()
	outMessage4 := comm.OutterMessage{}
	outMessage4.Message = Message4

	pbftService.ReceiveConsensusMessage(outMessage4)
	assert.Equal(t,len(pbftService.GetCurrentConsensusState()),2)
}