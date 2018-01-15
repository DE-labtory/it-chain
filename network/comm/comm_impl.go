package comm

import (
	"it-chain/service/peer"
	pb "it-chain/network/protos"
)
type CommImpl struct{
	connectionMap map[string]*Connection

}

func NewCommImpl() *CommImpl{
	return &CommImpl{}
}

func (commImpl *CommImpl)CreateConn(peerInfo peer.PeerInfo) error{
	//endpoint := peerInfo.GetEndPoint()
	//NewClientConnectionWithAddress(endpoint,false,nil)

	return nil
}

func (commImpl *CommImpl) Send(envelop pb.Envelope, peerInfos []peer.PeerInfo){

}

func (commImpl *CommImpl) Stop(){

}

func (commImpl *CommImpl) Close(peerInfo peer.PeerInfo){

}

func (commImpl *CommImpl) CloseAll(){

}
