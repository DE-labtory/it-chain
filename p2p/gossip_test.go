package p2p

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"time"
)

func getMyGossipTableAndPeerInfo() (*GossipTable,*PeerInfo){
	peerInfo := &PeerInfo{}
	peerInfo.peerID = "1"
	peerInfo.counter = 1
	peerInfo.ipAddress = "127.0.0.1"
	gossipTable := &GossipTable{}
	gossipTable.peerList = make([]*PeerInfo,0)
	gossipTable.peerList = append(gossipTable.peerList, peerInfo)
	gossipTable.updatedTimeStamp = time.Now()
	gossipTable.myID = peerInfo.peerID

	return gossipTable, peerInfo
}

func TestCreateNewGossipTable(t *testing.T) {

	//when
	peerInfo := &PeerInfo{}
	peerInfo.peerID = "1"
	peerInfo.counter = 1
	peerInfo.ipAddress = "127.0.0.1"

	//then
	gossipTable := CreateNewGossipTable(peerInfo)

	//result
	assert.Equal(t,1,len(gossipTable.peerList))
	assert.Equal(t,peerInfo,gossipTable.peerList[0])
	assert.Equal(t,peerInfo.peerID,gossipTable.myID)

}

func TestGossipTable_FindPeerByPeerID(t *testing.T) {

	//when
	gossipTable, peerInfo := getMyGossipTableAndPeerInfo()

	//then
	pi := gossipTable.FindPeerByPeerID(peerInfo.peerID)
	pi2 := gossipTable.FindPeerByPeerID("2")

	//result
	assert.Equal(t,peerInfo,pi)
	assert.Nil(t,pi2)
}

func TestGossipTable_IncrementMyCounter(t *testing.T) {
	//when
	gossipTable, _ := getMyGossipTableAndPeerInfo()

	//then
	gossipTable.IncrementMyCounter()

	//result
	assert.Equal(t,2,gossipTable.peerList[0].counter)
}

func TestGossipTable_addPeer(t *testing.T){

	//when
	gossipTable, _ := getMyGossipTableAndPeerInfo()
	peerInfo2 := &PeerInfo{}
	peerInfo2.peerID = "3"
	peerInfo2.counter = 2
	peerInfo2.ipAddress = "127.0.0.1"

	gossipTable.addPeerInfo(peerInfo2)

	assert.Equal(t,2,len(gossipTable.peerList))
	assert.Equal(t,peerInfo2,gossipTable.peerList[1])
}


//todo for문 돌면서 addpeer에 대한 검증이 필요함
func TestGossipTable_UpdateGossipTable(t *testing.T) {

	//when
	gossipTable, _ := getMyGossipTableAndPeerInfo()
	peerInfo2 := &PeerInfo{}
	peerInfo2.peerID = "3"
	peerInfo2.counter = 2
	peerInfo2.ipAddress = "127.0.0.2"
	gossipTable.peerList = append(gossipTable.peerList, peerInfo2)

	gossipTable2, _ := getMyGossipTableAndPeerInfo()
	peerInfo3 := &PeerInfo{}
	peerInfo3.peerID = "3"
	peerInfo3.counter = 3
	peerInfo3.ipAddress = "127.0.1.3"
	gossipTable2.peerList = append(gossipTable2.peerList,peerInfo3)

	//then
	gossipTable.UpdateGossipTable(*gossipTable2)

	assert.Equal(t,3,gossipTable.peerList[1].counter)
}

func TestGossipTable_SelectRandomPeerInfo(t *testing.T) {
	gossipTable, _ := getMyGossipTableAndPeerInfo()
	peerInfo2 := &PeerInfo{}
	peerInfo2.peerID = "2"
	peerInfo2.counter = 2
	peerInfo2.ipAddress = "127.0.0.2"
	gossipTable.peerList = append(gossipTable.peerList, peerInfo2)

	peerInfo3 := &PeerInfo{}
	peerInfo3.peerID = "3"
	peerInfo3.counter = 3
	peerInfo3.ipAddress = "127.0.1.3"
	gossipTable.peerList = append(gossipTable.peerList,peerInfo3)

	var ipAddresslist []string
	var err error

	ipAddresslist, err = gossipTable.SelectRandomPeerInfo(1)

	if err != nil{
		assert.Fail(t,err.Error())
	}

	assert.Equal(t,3,len(ipAddresslist))

	ipAddresslist, err = gossipTable.SelectRandomPeerInfo(0.7)

	if err != nil{
		assert.Fail(t,err.Error())
	}

	assert.Equal(t,2,len(ipAddresslist))
}