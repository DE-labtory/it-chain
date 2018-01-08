package p2p

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"time"
	"fmt"
	"sync"
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

func TestGossip_Process(t *testing.T){
	peerInfo := &PeerInfo{}
	peerInfo.peerID = "1"
	peerInfo.counter = 1
	peerInfo.ipAddress = "127.0.0.1"
	peerInfo.timeStamp = time.Now()
	gossiptable1 := CreateNewGossipTable(peerInfo)

	peerInfo2 := &PeerInfo{}
	peerInfo2.peerID = "2"
	peerInfo2.counter = 1
	peerInfo2.ipAddress = "127.0.0.2"
	peerInfo2.timeStamp = time.Now()
	gossiptable2 := CreateNewGossipTable(peerInfo2)
	gossiptable1.addPeerInfo(peerInfo2)

	peerInfo3 := &PeerInfo{}
	peerInfo3.peerID = "3"
	peerInfo3.counter = 1
	peerInfo3.ipAddress = "127.0.0.3"
	peerInfo3.timeStamp = time.Now()
	gossiptable3 := CreateNewGossipTable(peerInfo3)
	gossiptable2.addPeerInfo(peerInfo3)

	gossipTableList := []*GossipTable{gossiptable1,gossiptable2,gossiptable3}
	fmt.Print(gossipTableList)
	ticker := time.NewTicker(time.Millisecond * 500)

	var wg sync.WaitGroup
	wg.Add(6)

	messages1 := make(chan GossipTable)
	messages2 := make(chan GossipTable)
	messages3 := make(chan GossipTable)

	channelList := make(map[string](chan GossipTable))
	channelList["127.0.0.1"] = messages1
	channelList["127.0.0.2"] = messages2
	channelList["127.0.0.3"] = messages3

	go func() {
		defer wg.Done()
		go func() {
			for{
				fmt.Println("Received1")
				newGossipTable := <- messages1
				gossiptable1.Lock()
				gossiptable1.UpdateGossipTable(newGossipTable)
				gossiptable1.Unlock()
				fmt.Println(gossiptable1.peerList[1].counter)
			}
		}()
		for t := range ticker.C {
			gossiptable1.Lock()
			gossiptable1.IncrementMyCounter()
			peerList,err := gossiptable1.SelectRandomPeerInfo(0.6)

			if err != nil {
				fmt.Println(err)
			}
			fmt.Println(len(peerList))
			for _, address := range peerList{
				fmt.Println("channel",channelList[address])
				channelList[address] <- *gossiptable1
			}
			gossiptable1.Unlock()
			fmt.Println(t)
		}
	}()

	go func() {
		defer wg.Done()
		go func() {
			for {
				fmt.Println("Received2")
				newGossipTable := <-messages2
				gossiptable2.Lock()
				fmt.Println("time",newGossipTable.peerList[0].timeStamp)
				gossiptable2.UpdateGossipTable(newGossipTable)
				gossiptable2.Unlock()
			}
		}()
		for t := range ticker.C {
			gossiptable2.Lock()
			gossiptable2.IncrementMyCounter()

			peerList,err := gossiptable2.SelectRandomPeerInfo(0.6)
			if err != nil {
			}
			for _, address := range peerList{
				channelList[address] <- *gossiptable2
			}

			gossiptable2.Unlock()
			fmt.Println("Tick at", t)
		}
	}()

	go func() {
		defer wg.Done()
		go func() {
			for {
				fmt.Println("Received3")
				newGossipTable := <- messages3
				gossiptable3.Lock()
				gossiptable3.UpdateGossipTable(newGossipTable)
				gossiptable3.Unlock()
				fmt.Println(gossiptable3)
			}
		}()
		for t := range ticker.C {
			gossiptable3.Lock()
			gossiptable3.IncrementMyCounter()

			peerList,err := gossiptable3.SelectRandomPeerInfo(0.6)
			if err != nil {
			}
			for _, address := range peerList{
				channelList[address] <- *gossiptable3
			}

			gossiptable3.Unlock()
			fmt.Println("Tick at", t)
		}
	}()

	time.Sleep(time.Millisecond * 20000)
	ticker.Stop()
	fmt.Println("Ticker stopped")

	wg.Wait()
}
