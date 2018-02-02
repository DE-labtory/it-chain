package consensus

import (
	"testing"
	"it-chain/network/protos"
	"it-chain/network/comm"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/mock"
)

type MockConnectionManager struct{
	mock.Mock
}

func (mcm MockConnectionManager) SendStream(envelope message.Envelope, errorCallBack comm.OnError, connectionID string){

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

}




func TestNewPBFTConsensusService(t *testing.T) {
	comm:=MockConnectionManager{}
	//NewPBFTConsensusService(comm,)
}

func TestPBFTConsensusService_broadcastMessage(t *testing.T) {

}

func TestPBFTConsensusService_StartConsensus(t *testing.T) {

}
