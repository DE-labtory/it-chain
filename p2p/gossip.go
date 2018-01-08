package p2p

import (
	"time"
	"errors"
	"math/rand"
)

type PeerInfo struct{
	ipAddress string
	peerID string
	counter int
	timeStamp time.Time
}

type GossipTable struct{
	peerList []*PeerInfo
	updatedTimeStamp time.Time
	myID string
}

func CreateNewGossipTable(peerInfo *PeerInfo) *GossipTable{

	// 생성할때 넣어주는 peerInfo의 ID가 myID가 된다.
	gossipTable := &GossipTable{}
	gossipTable.peerList = make([]*PeerInfo,0)
	gossipTable.peerList = append(gossipTable.peerList, peerInfo)
	gossipTable.updatedTimeStamp = time.Now()
	gossipTable.myID = peerInfo.peerID

	return gossipTable
}

func (gt *GossipTable) FindPeerByPeerID(peerID string) (*PeerInfo){

	for _, peer := range gt.peerList {
		if peer.peerID == peerID{
			return peer
		}
	}

	return nil
}

func (gt *GossipTable) addPeerInfo(peerInfo *PeerInfo){
	// 존재하지 않을때 peerInfo를 add하는 함수
	// 존재하지 않은것이 확인된 후에 사용해야함
	gt.peerList = append(gt.peerList, peerInfo)
}

func (gt *GossipTable) UpdateGossipTable(gossipTable GossipTable){
	//fmt.Print(gossipTable)
	for _, updatePeer := range gossipTable.peerList {
		peer := gt.FindPeerByPeerID(updatePeer.peerID)
		//fmt.Print(peer)
		if peer != nil{
			// peer가 존재
			// counter를 비교하여 update
			if updatePeer.counter > peer.counter {
				//다 update하는게 좋을까?
				peer.counter = updatePeer.counter
				peer.timeStamp = time.Now()
				peer.ipAddress = updatePeer.ipAddress
			}
		}else{
			// peer를 새로 추가한다.
			peer.timeStamp = time.Now()
			gt.addPeerInfo(peer)
		}
	}

	gt.updatedTimeStamp = time.Now()
}

func (gt *GossipTable) IncrementMyCounter() error{

	myPeer := gt.FindPeerByPeerID(gt.myID)

	if myPeer == nil{
		return errors.New("myID peer does not exist error")
	}

	myPeer.counter += 1

	return nil
}

func (gt *GossipTable) SelectRandomPeerInfo(percent float64) ([]string,error) {

	if len(gt.peerList) < 1{
		return nil, errors.New("no peer in gossiptable")
	}

	num := int(percent*float64(len(gt.peerList)))

	if num < 1{
		return nil, nil
	}

	tmp := make([]*PeerInfo, len(gt.peerList))
	copy(tmp, gt.peerList)

	ipAddressList := make([]string,0)

	for i := 0; i < num; i++ {
		rand.Seed(time.Now().UTC().UnixNano())
		randNum := rand.Int() % len(tmp)
		ipAddressList = append(ipAddressList, tmp[randNum].ipAddress)

		//delete
		tmp = append(tmp[:randNum], tmp[randNum+1:]...)
	}

	return ipAddressList,nil
}

