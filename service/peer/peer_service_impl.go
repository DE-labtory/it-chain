package peer

type PeerServiceImpl struct {
	peerTable *PeerTable
}

func NewPeerServiceImpl(peerTable *PeerTable) *PeerServiceImpl{

	return &PeerServiceImpl{
		peerTable: peerTable,
	}
}

func (ps *PeerServiceImpl) GetPeerTable() PeerTable{
	return *ps.peerTable
}

//peer info 찾기
func (ps *PeerServiceImpl) GetPeerInfoByPeerID(peerID string) (*PeerInfo){

	peerInfo := ps.peerTable.FindPeerByPeerID(peerID)
	return peerInfo
}

//peer info
func (ps *PeerServiceImpl) PushPeerTable(peerIDs []string){

}

//주기적으로 handle 함수가 콜 된다.
//주기적으로 peerTable의 peerlist에게 peerTable을 전송한다.
func (ps *PeerServiceImpl) handle(interface{}){

}

////peer message를 peersIP에 전파
//func (gs GossipServiceImpl) Push(gossipTable GossipTable, peersIP []string){
//
//	for _, ip := range peersIP{
//		go func(){
//			conn, err := grpc.Dial(ip, grpc.WithInsecure())
//			if err != nil {
//				logger.Fatalf("did not connect: %v", err)
//			}
//
//			defer conn.Close()
//			c := pb.NewGossipClient(conn)
//
//			r, err := c.PushGossip(context.Background(), gossipTable.toProto())
//
//			if err != nil {
//				logger.Error("fail to push peer ",r.String())
//			}
//
//			logger.Info("success to push peer",r.String())
//		}()
//	}
//}
//
//func (gs *GossipServiceImpl) Update(gossipTable GossipTable){
//
//	gs.gossipTable.UpdateGossipTable(gossipTable)
//	logger.Info("update peer", gs.gossipTable)
//}
//
//func (gs GossipServiceImpl) Listen(gossipTable GossipTable){
//
//}
//
//func (gs GossipServiceImpl) Pull() GossipTable{
//
//	return *gs.gossipTable
//}
//
//func (gs GossipServiceImpl) Stop(){
//
//}
//
//func (gs GossipServiceImpl) GetMyGossipTable() *GossipTable{
//	return gs.gossipTable
//}

