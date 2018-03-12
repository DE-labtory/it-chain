package domain

import (
	"testing"
	"time"
	"github.com/stretchr/testify/assert"
	"fmt"
)

func MockCreateNewPeer(peerID string) *Peer{

	return  &Peer{
		PeerID: peerID,
		Port: "8080",
		IpAddress: "127.0.0.1",
		HeartBeat: 1,
		TimeStamp: time.Now(),
	}
}

func MockCreateNewPeerTable(peerID string) *PeerTable{

	Peer := MockCreateNewPeer(peerID)
	peerTable,err := NewPeerTable(Peer)

	if err != nil{

	}

	return peerTable
}

func TestPeer_Update(t *testing.T) {
	peer1 := MockCreateNewPeer("test1")
	peer2 := MockCreateNewPeer("test1")
	peer2.HeartBeat = 5
	peer2.Port = "7777"
	peer2.IpAddress = "127.0.0.2"

	peer1.Update(peer2)

	assert.Equal(t,peer1.HeartBeat ,5)
	assert.Equal(t,peer1.Port ,"7777")
	assert.Equal(t,peer1.IpAddress ,"127.0.0.2")
}

func TestNewPeerTable(t *testing.T) {

	Peer := MockCreateNewPeer("test1")
	peerTable,err := NewPeerTable(Peer)

	if err != nil{
		assert.Fail(t,"fail to create new peertable")
	}

	assert.Equal(t,len(peerTable.PeerMap),1)
	assert.Equal(t,peerTable.PeerMap[Peer.PeerID],Peer)
}

func TestPeerTable_FindPeerByPeerID(t *testing.T) {

	Peer := MockCreateNewPeer("test1")
	peerTable,err := NewPeerTable(Peer)
	if err != nil{
		assert.Fail(t,"fail to create new peertable")
	}

	//when peer exist
	assert.Equal(t,peerTable.FindPeerByPeerID(Peer.PeerID),Peer)

	//when peer does not exist
	assert.Nil(t,peerTable.FindPeerByPeerID("test2"))
}

func TestPeerTable_AddPeer(t *testing.T) {
	//when
	Peer := MockCreateNewPeer("test1")
	peerTable,err := NewPeerTable(Peer)

	if err != nil{
		assert.Fail(t,"fail to create new peertable")
	}

	Peer2 := MockCreateNewPeer("test2")
	peerTable.AddPeer(Peer2)

	assert.Equal(t,len(peerTable.PeerMap),2)
	assert.Equal(t,peerTable.PeerMap[Peer2.PeerID],Peer2)
}

func TestPeerTable_UpdatePeerTable(t *testing.T) {

	Peer := MockCreateNewPeer("test1")
	peerTable,err := NewPeerTable(Peer)

	if err != nil{
		assert.Fail(t,"fail to create new peertable")
	}

	Peer2 := MockCreateNewPeer("test1")
	Peer2.HeartBeat = 3
	Peer2.IpAddress = "127.0.0.2"
	Peer2.Port = "7070"

	Peer3 := MockCreateNewPeer("test3")

	peerTable2,err := NewPeerTable(Peer2)

	if err != nil{
		assert.Fail(t,"fail to create new peertable")
	}

	peerTable2.PeerMap[Peer3.PeerID] = Peer3


	//////
	peerTable.UpdatePeerTable(*peerTable2)


	assert.Equal(t,len(peerTable.PeerMap),2)
	assert.Equal(t,peerTable.PeerMap[Peer2.PeerID].HeartBeat,Peer2.HeartBeat)
	assert.Equal(t,peerTable.PeerMap[Peer2.PeerID].IpAddress,Peer2.IpAddress)
	assert.Equal(t,peerTable.PeerMap[Peer2.PeerID].Port,Peer2.Port)

	assert.Equal(t,peerTable.PeerMap[Peer3.PeerID].HeartBeat,Peer3.HeartBeat)
	assert.Equal(t,peerTable.PeerMap[Peer3.PeerID].IpAddress,Peer3.IpAddress)
	assert.Equal(t,peerTable.PeerMap[Peer3.PeerID].Port,Peer3.Port)
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

func TestPeerTable_GetPeerList(t *testing.T) {
	peerTable := MockCreateNewPeerTable("test1")

	Peer2 := MockCreateNewPeer("test2")
	Peer2.HeartBeat = 3
	Peer2.IpAddress = "127.0.0.2"
	Peer2.Port = "7070"
	peerTable.AddPeer(Peer2)

	Peer3 := MockCreateNewPeer("test3")
	Peer3.HeartBeat = 3
	Peer3.IpAddress = "127.0.0.2"
	Peer3.Port = "7070"
	peerTable.AddPeer(Peer3)

	Peer4 := MockCreateNewPeer("test4")
	Peer4.HeartBeat = 3
	Peer4.IpAddress = "127.0.0.2"
	Peer4.Port = "7070"
	peerTable.AddPeer(Peer4)


	peerList := peerTable.GetPeerList()

	assert.Equal(t,peerList[0].PeerID,"test2")
	assert.Equal(t,peerList[1].PeerID,"test3")
	assert.Equal(t,peerList[2].PeerID,"test4")
}

func TestPeerTable_GetMyInfo(t *testing.T) {
	peerTable := MockCreateNewPeerTable("test1")

	Peer2 := MockCreateNewPeer("test2")
	Peer2.HeartBeat = 3
	Peer2.IpAddress = "127.0.0.2"
	Peer2.Port = "7070"
	peerTable.AddPeer(Peer2)

	myInfo := peerTable.GetMyInfo()

	assert.Equal(t,myInfo.PeerID,"test1")
}

////todo for문 돌면서 addpeer에 대한 검증이 필요함
//func TestGossipTable_UpdateGossipTable(t *testing.T) {
//
//	//when
//	gossipTable, _ := getMyGossipTableAndPeer()
//	Peer2 := &Peer{}
//	Peer2.peerID = "3"
//	Peer2.counter = 2
//	Peer2.ipAddress = "127.0.0.2"
//	gossipTable.peerList = append(gossipTable.peerList, Peer2)
//
//	gossipTable2, _ := getMyGossipTableAndPeer()
//	Peer3 := &Peer{}
//	Peer3.peerID = "3"
//	Peer3.counter = 3
//	Peer3.ipAddress = "127.0.1.3"
//	gossipTable2.peerList = append(gossipTable2.peerList,Peer3)
//
//	//then
//	gossipTable.UpdateGossipTable(*gossipTable2)
//
//	assert.Equal(t,3,gossipTable.peerList[1].counter)
//}
//
//func TestGossipTable_SelectRandomPeer(t *testing.T) {
//	gossipTable, _ := getMyGossipTableAndPeer()
//	Peer2 := &Peer{}
//	Peer2.peerID = "2"
//	Peer2.counter = 2
//	Peer2.ipAddress = "127.0.0.2"
//	gossipTable.peerList = append(gossipTable.peerList, Peer2)
//
//	Peer3 := &Peer{}
//	Peer3.peerID = "3"
//	Peer3.counter = 3
//	Peer3.ipAddress = "127.0.1.3"
//	gossipTable.peerList = append(gossipTable.peerList,Peer3)
//
//	var ipAddresslist []string
//	var err error
//
//	ipAddresslist, err = gossipTable.SelectRandomPeer(1)
//
//	if err != nil{
//		assert.Fail(t,err.Error())
//	}
//
//	assert.Equal(t,3,len(ipAddresslist))
//
//	ipAddresslist, err = gossipTable.SelectRandomPeer(0.7)
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
//	Peer := &Peer{}
//	Peer.peerID = "1"
//	Peer.counter = 1
//	Peer.ipAddress = "127.0.0.1"
//	Peer.timeStamp = time.Now()
//	gossiptable1,_ := CreateNewGossipTable(Peer)
//
//	Peer2 := &Peer{}
//	Peer2.peerID = "2"
//	Peer2.counter = 1
//	Peer2.ipAddress = "127.0.0.2"
//	Peer2.timeStamp = time.Now()
//	gossiptable2,_ := CreateNewGossipTable(Peer2)
//	gossiptable1.addPeer(Peer2)
//
//	Peer3 := &Peer{}
//	Peer3.peerID = "3"
//	Peer3.counter = 1
//	Peer3.ipAddress = "127.0.0.3"
//	Peer3.timeStamp = time.Now()
//	gossiptable3,_ := CreateNewGossipTable(Peer3)
//	gossiptable2.addPeer(Peer3)
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
//			peerList,err := gossiptable1.SelectRandomPeer(0.6)
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
//			peerList,err := gossiptable2.SelectRandomPeer(0.6)
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
//			peerList,err := gossiptable3.SelectRandomPeer(0.6)
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
