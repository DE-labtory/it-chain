package peer

import (
	"time"
	"sync"
	"errors"
	"fmt"
	"encoding/json"
	"math/rand"
)

type PeerInfo struct {
	IpAddress string
	Port      string
	PeerID    string
	HeartBeat int
	TimeStamp time.Time
	PubKey    []byte
}

func (pi *PeerInfo) UpdateTimeStamp(){
	pi.TimeStamp = time.Now()
}

func (pi *PeerInfo) Update(peerInfo *PeerInfo) (error){

	if pi.PeerID != peerInfo.PeerID{
		return errors.New("different peer id")
	}

	if pi.HeartBeat < peerInfo.HeartBeat{
		pi.HeartBeat = peerInfo.HeartBeat
		pi.IpAddress = peerInfo.IpAddress
		pi.Port = peerInfo.Port
		pi.PubKey = peerInfo.PubKey
		pi.UpdateTimeStamp()
	}

	return nil
}

func (pi *PeerInfo)GetEndPoint() string{
	return pi.IpAddress+":"+pi.Port
}

type PeerTable struct {
	PeerMap   map[string]*PeerInfo
	TimeStamp time.Time
	OwnerID   string
	sync.RWMutex
}

func NewPeerTable(myInfo *PeerInfo) (*PeerTable, error){

	// 생성할때 넣어주는 peerInfo의 ID가 myID가 된다.
	if myInfo.PeerID == ""{
		return nil, errors.New("peerInfo must have peerID")
	}

	if myInfo.IpAddress == "" || myInfo.Port == ""{
		return nil, errors.New("peerInfo must have ipAddress")
	}

	peerMap := make(map[string]*PeerInfo)
	peerMap[myInfo.PeerID] = myInfo

	return &PeerTable{
		PeerMap: peerMap,
		TimeStamp: time.Now(),
		OwnerID: myInfo.PeerID,
	}, nil
}

//tested
func (pt *PeerTable) FindPeerByPeerID(peerID string) (*PeerInfo){

	peerInfo, ok := pt.PeerMap[peerID]

	if ok {
		return peerInfo
	}

	return nil
}

//tested
//if does not exist insert
//if exist update all
func (pt *PeerTable) AddPeerInfo(peerInfo *PeerInfo){

	pt.PeerMap[peerInfo.PeerID] = peerInfo
	pt.PeerMap[peerInfo.PeerID].TimeStamp = time.Now()
}

//tested
func (pt *PeerTable) UpdatePeerTable(peerTable PeerTable){

	pt.Lock()
	defer pt.Unlock()

	for id, peerInfo := range peerTable.PeerMap{
		peer,ok := pt.PeerMap[id]

		if ok{
			peer.Update(peerInfo)
		}else{
			pt.AddPeerInfo(peerInfo)
		}
	}

	pt.UpdateTimeStamp()
}

//tested
func (pt *PeerTable) UpdateTimeStamp(){

	pt.TimeStamp = time.Now()
}

////tested
func (pt *PeerTable) IncrementHeartBeat() error{

	pt.Lock()
	defer pt.Unlock()

	myPeer := pt.FindPeerByPeerID(pt.OwnerID)

	if myPeer == nil{
		return errors.New("myID peer does not exist error")
	}

	myPeer.HeartBeat += 1

	return nil
}

////tested
func (pt *PeerTable) SelectRandomPeerInfos(percent float64) ([]PeerInfo,error){

	if len(pt.PeerMap) <= 1{
		return nil, errors.New("no peer in gossiptable")
	}

	num := int(percent*float64(len(pt.PeerMap)))

	if num < 1{
		return nil, errors.New("no peer in gossiptable")
	}

	tmp := make([]*PeerInfo, 0)
	for _, peer := range pt.PeerMap{
		if peer.PeerID != pt.OwnerID{
			//내 ID는 삭제
			tmp = append(tmp, peer)
		}
	}

	peerInfoList := make([]PeerInfo, 0)

	for i := 0; i < num; i++ {
		rand.Seed(time.Now().UTC().UnixNano())
		randNum := rand.Int() % len(tmp)
		peerInfoList = append(peerInfoList, *tmp[randNum])

		//delete
		tmp = append(tmp[:randNum], tmp[randNum+1:]...)
	}

	return peerInfoList,nil
}

//tested
//나자신을 제외한 peerInfo의 list를 return
func (pt *PeerTable) GetPeerList() []PeerInfo{

	tmp := make([]PeerInfo, 0)
	for _, peer := range pt.PeerMap{
		if peer.PeerID != pt.OwnerID{
			//내 ID는 삭제
			tmp = append(tmp, *peer)
		}
	}

	return tmp
}

//tested
func (pt *PeerTable) GetMyInfo() PeerInfo{
	peerInfo, ok := pt.PeerMap[pt.OwnerID]

	if ok{
		return *peerInfo
	}

	logger.Println("No my peer info")

	return *peerInfo
}

func (pi PeerInfo) String() string {

	b, err := json.Marshal(pi)
	if err != nil {
		fmt.Println("error:", err)
	}

	return string(b)
}

func (pt PeerTable) String() string {

	b, err := json.Marshal(pt)
	if err != nil {
		fmt.Println("error:", err)
	}

	return string(b)
}