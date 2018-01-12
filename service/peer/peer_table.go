package peer

import (
	"time"
	"sync"
)

type PeerInfo struct {
	IpAddress string
	PeerID    string
	HeartBeat int
	TimeStamp time.Time
	PubKey    []byte
}

type PeerTable struct {
	PeerMap  map[string]*PeerInfo
	TimeStamp time.Time
	ID      string
	sync.RWMutex
}

////tested
//func CreateNewGossipTable(peerInfo *PeerInfo) (*GossipTable,error) {
//
//	// 생성할때 넣어주는 peerInfo의 ID가 myID가 된다.
//	if peerInfo.peerID == ""{
//		logger.Error("peerInfo must have peerID")
//		return nil, errors.New("peerInfo must have peerID")
//	}
//
//	if peerInfo.ipAddress == ""{
//		logger.Error("peerInfo must have ipAddress")
//		return nil, errors.New("peerInfo must have ipAddress")
//	}
//
//	gossipTable := &GossipTable{}
//	gossipTable.peerList = make([]*PeerInfo,0)
//	gossipTable.peerList = append(gossipTable.peerList, peerInfo)
//	gossipTable.timeStamp = time.Now()
//	gossipTable.myID = peerInfo.peerID
//
//	logger.Info("new peer table created",gossipTable)
//
//	return gossipTable, nil
//}
//
////tested
//func (gt *GossipTable) FindPeerByPeerID(peerID string) (*PeerInfo){
//
//	for _, peer := range gt.peerList {
//		if peer.peerID == peerID{
//			return peer
//		}
//	}
//
//	logger.Info("can not find peer by ",peerID)
//
//	return nil
//}
//
////tested
//func (gt *GossipTable) addPeerInfo(peerInfo *PeerInfo){
//	// 존재하지 않을때 peerInfo를 add하는 함수
//	// 존재하지 않은것이 확인된 후에 사용해야함
//	gt.peerList = append(gt.peerList, peerInfo)
//	logger.Info("new peerInfo added",peerInfo)
//}
//
////tested
//func (gt *GossipTable) UpdateGossipTable(gossipTable GossipTable){
//	//fmt.Print(gossipTable)
//	for _, updatePeer := range gossipTable.peerList {
//		peer := gt.FindPeerByPeerID(updatePeer.peerID)
//		//fmt.Print(peer)
//		if peer != nil{
//			// peer가 존재
//			// counter를 비교하여 update
//			if updatePeer.counter > peer.counter {
//				logger.Info("peerInfo updated")
//				//다 update하는게 좋을까?
//				peer.counter = updatePeer.counter
//				peer.timeStamp = time.Now()
//				peer.ipAddress = updatePeer.ipAddress
//			}
//		}else{
//			// peer를 새로 추가한다.
//			updatePeer.timeStamp = time.Now()
//			gt.addPeerInfo(updatePeer)
//		}
//	}
//
//	gt.timeStamp = time.Now()
//}
//
////tested
//func (gt *GossipTable) IncrementMyCounter() error{
//	logger.Info("increment my counter")
//
//	myPeer := gt.FindPeerByPeerID(gt.myID)
//	if myPeer == nil{
//		return errors.New("myID peer does not exist error")
//	}
//
//	myPeer.counter += 1
//
//	return nil
//}
//
////tested
//func (gt *GossipTable) SelectRandomPeerInfo(percent float64) ([]string,error) {
//
//	if len(gt.peerList) < 1{
//		return nil, errors.New("no peer in gossiptable")
//	}
//
//	num := int(percent*float64(len(gt.peerList)))
//
//	if num < 1{
//		return nil, nil
//	}
//
//	tmp := make([]*PeerInfo, len(gt.peerList))
//	copy(tmp, gt.peerList)
//
//	ipAddressList := make([]string,0)
//
//	for i := 0; i < num; i++ {
//		rand.Seed(time.Now().UTC().UnixNano())
//		randNum := rand.Int() % len(tmp)
//		ipAddressList = append(ipAddressList, tmp[randNum].ipAddress)
//
//		//delete
//		tmp = append(tmp[:randNum], tmp[randNum+1:]...)
//	}
//
//	return ipAddressList,nil
//}
//
//
//func (pi PeerInfo) String() string {
//	return fmt.Sprintf("{myID:%s, counter:%d, ipAddress:%s timeStamp:%s}", pi.peerID, pi.counter, pi.ipAddress,pi.timeStamp.Format(time.RFC3339))
//}
//
//func (gt GossipTable) String() string {
//	return fmt.Sprintf("{myID:%s, timeStamp:%s, peerList:%s}", gt.myID, gt.timeStamp.Format(time.RFC3339), gt.peerList)
//}
//
//
//func (gt GossipTable) toProto() *pb.GossipTable{
//
//	pb_gossipTable := &pb.GossipTable{
//		MyID : gt.myID,
//		PeerInfo: make([]*pb.PeerInfo,0),
//	}
//
//	for _, peerInfo := range gt.peerList{
//		pb_peerInfo := &pb.PeerInfo{
//			Counter: int64(peerInfo.counter),
//			IpAddress: peerInfo.ipAddress,
//			PeerID: peerInfo.peerID,
//		}
//		pb_gossipTable.PeerInfo = append(pb_gossipTable.PeerInfo, pb_peerInfo)
//	}
//
//	return pb_gossipTable
//}