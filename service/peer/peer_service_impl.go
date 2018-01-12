package peer

import (
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	pb "it-chain/sample/gossip/grpc"
)

type P2PServiceImpl struct{
	gossipTable *GossipTable
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

