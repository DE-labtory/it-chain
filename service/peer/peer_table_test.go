package peer

import (
	"testing"
	"time"
	"github.com/stretchr/testify/assert"
	"fmt"
)

func MockCreateNewPeerInfo(peerID string) *PeerInfo{

	return  &PeerInfo{
		PeerID: peerID,
		Port: "8080",
		IpAddress: "127.0.0.1",
		HeartBeat: 1,
		TimeStamp: time.Now(),
	}
}

func MockCreateNewPeerTable(peerID string) *PeerTable{

	peerInfo := MockCreateNewPeerInfo(peerID)
	peerTable,err := NewPeerTable(peerInfo)

	if err != nil{

	}

	return peerTable
}

func TestPeerInfo_Update(t *testing.T) {
	peer1 := MockCreateNewPeerInfo("test1")
	peer2 := MockCreateNewPeerInfo("test1")
	peer2.HeartBeat = 5
	peer2.Port = "7777"
	peer2.IpAddress = "127.0.0.2"

	peer1.Update(peer2)

	assert.Equal(t,peer1.HeartBeat ,5)
	assert.Equal(t,peer1.Port ,"7777")
	assert.Equal(t,peer1.IpAddress ,"127.0.0.2")
}

func TestNewPeerTable(t *testing.T) {

	peerInfo := MockCreateNewPeerInfo("test1")
	peerTable,err := NewPeerTable(peerInfo)

	if err != nil{
		assert.Fail(t,"fail to create new peertable")
	}

	assert.Equal(t,len(peerTable.PeerMap),1)
	assert.Equal(t,peerTable.PeerMap[peerInfo.PeerID],peerInfo)
}

func TestPeerTable_FindPeerByPeerID(t *testing.T) {

	peerInfo := MockCreateNewPeerInfo("test1")
	peerTable,err := NewPeerTable(peerInfo)
	if err != nil{
		assert.Fail(t,"fail to create new peertable")
	}

	//when peer exist
	assert.Equal(t,peerTable.FindPeerByPeerID(peerInfo.PeerID),peerInfo)

	//when peer does not exist
	assert.Nil(t,peerTable.FindPeerByPeerID("test2"))
}

func TestPeerTable_AddPeerInfo(t *testing.T) {
	//when
	peerInfo := MockCreateNewPeerInfo("test1")
	peerTable,err := NewPeerTable(peerInfo)

	if err != nil{
		assert.Fail(t,"fail to create new peertable")
	}

	peerInfo2 := MockCreateNewPeerInfo("test2")
	peerTable.AddPeerInfo(peerInfo2)

	assert.Equal(t,len(peerTable.PeerMap),2)
	assert.Equal(t,peerTable.PeerMap[peerInfo2.PeerID],peerInfo2)
}

func TestPeerTable_UpdatePeerTable(t *testing.T) {

	peerInfo := MockCreateNewPeerInfo("test1")
	peerTable,err := NewPeerTable(peerInfo)

	if err != nil{
		assert.Fail(t,"fail to create new peertable")
	}

	peerInfo2 := MockCreateNewPeerInfo("test1")
	peerInfo2.HeartBeat = 3
	peerInfo2.IpAddress = "127.0.0.2"
	peerInfo2.Port = "7070"

	peerInfo3 := MockCreateNewPeerInfo("test3")

	peerTable2,err := NewPeerTable(peerInfo2)

	if err != nil{
		assert.Fail(t,"fail to create new peertable")
	}

	peerTable2.PeerMap[peerInfo3.PeerID] = peerInfo3


	//////
	peerTable.UpdatePeerTable(*peerTable2)


	assert.Equal(t,len(peerTable.PeerMap),2)
	assert.Equal(t,peerTable.PeerMap[peerInfo2.PeerID].HeartBeat,peerInfo2.HeartBeat)
	assert.Equal(t,peerTable.PeerMap[peerInfo2.PeerID].IpAddress,peerInfo2.IpAddress)
	assert.Equal(t,peerTable.PeerMap[peerInfo2.PeerID].Port,peerInfo2.Port)

	assert.Equal(t,peerTable.PeerMap[peerInfo3.PeerID].HeartBeat,peerInfo3.HeartBeat)
	assert.Equal(t,peerTable.PeerMap[peerInfo3.PeerID].IpAddress,peerInfo3.IpAddress)
	assert.Equal(t,peerTable.PeerMap[peerInfo3.PeerID].Port,peerInfo3.Port)
}

func TestPeerTable_IncrementHeartBeat(t *testing.T) {

	peerTable := MockCreateNewPeerTable("test1")

	peerTable.IncrementHeartBeat()

	assert.Equal(t,peerTable.PeerMap["test1"].HeartBeat,2)
}

func TestPeerTable_string(t *testing.T){
	peerTable := MockCreateNewPeerTable("test1")
	fmt.Println(peerTable)
}

////todo for문 돌면서 addpeer에 대한 검증이 필요함
//func TestGossipTable_UpdateGossipTable(t *testing.T) {
//
//	//when
//	gossipTable, _ := getMyGossipTableAndPeerInfo()
//	peerInfo2 := &PeerInfo{}
//	peerInfo2.peerID = "3"
//	peerInfo2.counter = 2
//	peerInfo2.ipAddress = "127.0.0.2"
//	gossipTable.peerList = append(gossipTable.peerList, peerInfo2)
//
//	gossipTable2, _ := getMyGossipTableAndPeerInfo()
//	peerInfo3 := &PeerInfo{}
//	peerInfo3.peerID = "3"
//	peerInfo3.counter = 3
//	peerInfo3.ipAddress = "127.0.1.3"
//	gossipTable2.peerList = append(gossipTable2.peerList,peerInfo3)
//
//	//then
//	gossipTable.UpdateGossipTable(*gossipTable2)
//
//	assert.Equal(t,3,gossipTable.peerList[1].counter)
//}
//
//func TestGossipTable_SelectRandomPeerInfo(t *testing.T) {
//	gossipTable, _ := getMyGossipTableAndPeerInfo()
//	peerInfo2 := &PeerInfo{}
//	peerInfo2.peerID = "2"
//	peerInfo2.counter = 2
//	peerInfo2.ipAddress = "127.0.0.2"
//	gossipTable.peerList = append(gossipTable.peerList, peerInfo2)
//
//	peerInfo3 := &PeerInfo{}
//	peerInfo3.peerID = "3"
//	peerInfo3.counter = 3
//	peerInfo3.ipAddress = "127.0.1.3"
//	gossipTable.peerList = append(gossipTable.peerList,peerInfo3)
//
//	var ipAddresslist []string
//	var err error
//
//	ipAddresslist, err = gossipTable.SelectRandomPeerInfo(1)
//
//	if err != nil{
//		assert.Fail(t,err.Error())
//	}
//
//	assert.Equal(t,3,len(ipAddresslist))
//
//	ipAddresslist, err = gossipTable.SelectRandomPeerInfo(0.7)
//
//	if err != nil{
//		assert.Fail(t,err.Error())
//	}
//
//	assert.Equal(t,2,len(ipAddresslist))
//}
//
////todo refactor와 go routine exit action추가 해야함
//func TestGossip_Process(t *testing.T){
//
//	peerInfo := &PeerInfo{}
//	peerInfo.peerID = "1"
//	peerInfo.counter = 1
//	peerInfo.ipAddress = "127.0.0.1"
//	peerInfo.timeStamp = time.Now()
//	gossiptable1,_ := CreateNewGossipTable(peerInfo)
//
//	peerInfo2 := &PeerInfo{}
//	peerInfo2.peerID = "2"
//	peerInfo2.counter = 1
//	peerInfo2.ipAddress = "127.0.0.2"
//	peerInfo2.timeStamp = time.Now()
//	gossiptable2,_ := CreateNewGossipTable(peerInfo2)
//	gossiptable1.addPeerInfo(peerInfo2)
//
//	peerInfo3 := &PeerInfo{}
//	peerInfo3.peerID = "3"
//	peerInfo3.counter = 1
//	peerInfo3.ipAddress = "127.0.0.3"
//	peerInfo3.timeStamp = time.Now()
//	gossiptable3,_ := CreateNewGossipTable(peerInfo3)
//	gossiptable2.addPeerInfo(peerInfo3)
//
//	gossipTableList := []*GossipTable{gossiptable1,gossiptable2,gossiptable3}
//	fmt.Print(gossipTableList)
//	ticker := time.NewTicker(time.Millisecond * 500)
//
//	var wg sync.WaitGroup
//	wg.Add(6)
//
//	messages1 := make(chan GossipTable)
//	messages2 := make(chan GossipTable)
//	messages3 := make(chan GossipTable)
//
//	channelList := make(map[string](chan GossipTable))
//	channelList["127.0.0.1"] = messages1
//	channelList["127.0.0.2"] = messages2
//	channelList["127.0.0.3"] = messages3
//
//	go func() {
//		defer wg.Done()
//		go func() {
//			for{
//				fmt.Println("Received1")
//				newGossipTable := <- messages1
//				gossiptable1.Lock()
//				gossiptable1.UpdateGossipTable(newGossipTable)
//				gossiptable1.Unlock()
//				fmt.Println(gossiptable1.peerList[1].counter)
//			}
//		}()
//		for t := range ticker.C {
//			gossiptable1.Lock()
//			gossiptable1.IncrementMyCounter()
//			peerList,err := gossiptable1.SelectRandomPeerInfo(0.6)
//
//			if err != nil {
//				fmt.Println(err)
//			}
//			fmt.Println(len(peerList))
//			for _, address := range peerList{
//				fmt.Println("channel",channelList[address])
//				channelList[address] <- *gossiptable1
//			}
//			gossiptable1.Unlock()
//			fmt.Println(t)
//		}
//	}()
//
//	go func() {
//		defer wg.Done()
//		go func() {
//			for {
//				fmt.Println("Received2")
//				newGossipTable := <-messages2
//				gossiptable2.Lock()
//				fmt.Println("time",newGossipTable.peerList[0].timeStamp)
//				gossiptable2.UpdateGossipTable(newGossipTable)
//				gossiptable2.Unlock()
//			}
//		}()
//		for t := range ticker.C {
//			gossiptable2.Lock()
//			gossiptable2.IncrementMyCounter()
//
//			peerList,err := gossiptable2.SelectRandomPeerInfo(0.6)
//			if err != nil {
//			}
//			for _, address := range peerList{
//				channelList[address] <- *gossiptable2
//			}
//
//			gossiptable2.Unlock()
//			fmt.Println("Tick at", t)
//		}
//	}()
//
//	go func() {
//		defer wg.Done()
//		go func() {
//			for {
//				fmt.Println("Received3")
//				newGossipTable := <- messages3
//				gossiptable3.Lock()
//				gossiptable3.UpdateGossipTable(newGossipTable)
//				gossiptable3.Unlock()
//				fmt.Println(gossiptable3)
//			}
//		}()
//		for t := range ticker.C {
//			gossiptable3.Lock()
//			gossiptable3.IncrementMyCounter()
//
//			peerList,err := gossiptable3.SelectRandomPeerInfo(0.6)
//			if err != nil {
//			}
//			for _, address := range peerList{
//				channelList[address] <- *gossiptable3
//			}
//
//			gossiptable3.Unlock()
//			fmt.Println("Tick at", t)
//		}
//	}()
//
//	time.Sleep(time.Millisecond * 20000)
//	ticker.Stop()
//	fmt.Println("Ticker stopped")
//
//	wg.Wait()
//}
