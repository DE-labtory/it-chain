package domain

import (
	"time"
	"sync"
	"errors"
	"fmt"
	"encoding/json"
	"math/rand"
)

type Peer struct {
	IpAddress string
	Port      string
	PeerID    string
	HeartBeat int
	TimeStamp time.Time
	PubKey    []byte
}

func (pi *Peer) UpdateTimeStamp(){
	pi.TimeStamp = time.Now()
}

func (pi *Peer) Update(Peer *Peer) (error){

	if pi.PeerID != Peer.PeerID{
		return errors.New("different peer id")
	}

	if pi.HeartBeat < Peer.HeartBeat{
		pi.HeartBeat = Peer.HeartBeat
		pi.IpAddress = Peer.IpAddress
		pi.Port = Peer.Port
		pi.PubKey = Peer.PubKey
		pi.UpdateTimeStamp()
	}

	return nil
}

func (pi *Peer)GetEndPoint() string{
	return pi.IpAddress+":"+pi.Port
}

type PeerTable struct {
	PeerMap   map[string]*Peer
	Leader    *Peer
	TimeStamp time.Time
	MyID   string
	sync.RWMutex
}

func NewPeerTable(myInfo *Peer) (*PeerTable, error){

	// 생성할때 넣어주는 Peer의 ID가 myID가 된다.
	if myInfo.PeerID == ""{
		return nil, errors.New("Peer must have peerID")
	}

	if myInfo.IpAddress == "" || myInfo.Port == ""{
		return nil, errors.New("Peer must have ipAddress")
	}

	peerMap := make(map[string]*Peer)
	peerMap[myInfo.PeerID] = myInfo

	return &PeerTable{
		PeerMap: peerMap,
		TimeStamp: time.Now(),
		MyID: myInfo.PeerID,
	}, nil
}

//tested
func (pt *PeerTable) FindPeerByPeerID(peerID string) (*Peer){

	Peer, ok := pt.PeerMap[peerID]

	if ok {
		return Peer
	}

	return nil
}

//tested
//if does not exist insert
//if exist update all
func (pt *PeerTable) AddPeer(Peer *Peer){

	pt.PeerMap[Peer.PeerID] = Peer
	pt.PeerMap[Peer.PeerID].TimeStamp = time.Now()
}

//tested
func (pt *PeerTable) UpdatePeerTable(peerTable PeerTable){

	pt.Lock()
	defer pt.Unlock()

	for id, Peer := range peerTable.PeerMap{
		peer,ok := pt.PeerMap[id]

		if ok{
			peer.Update(Peer)
		}else{
			pt.AddPeer(Peer)
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

	myPeer := pt.FindPeerByPeerID(pt.MyID)

	if myPeer == nil{
		return errors.New("myID peer does not exist error")
	}

	myPeer.HeartBeat += 1

	return nil
}

////tested
func (pt *PeerTable) SelectRandomPeers(percent float64) ([]Peer,error){

	if len(pt.PeerMap) <= 1{
		return nil, errors.New("no peer in gossiptable")
	}

	num := int(percent*float64(len(pt.PeerMap)))

	if num < 1{
		return nil, errors.New("no peer in gossiptable")
	}

	tmp := make([]*Peer, 0)
	for _, peer := range pt.PeerMap{
		if peer.PeerID != pt.MyID{
			//내 ID는 삭제
			tmp = append(tmp, peer)
		}
	}

	PeerList := make([]Peer, 0)

	for i := 0; i < num; i++ {
		rand.Seed(time.Now().UTC().UnixNano())
		randNum := rand.Int() % len(tmp)
		PeerList = append(PeerList, *tmp[randNum])

		//delete
		tmp = append(tmp[:randNum], tmp[randNum+1:]...)
	}

	return PeerList,nil
}

//tested
//나자신을 제외한 Peer의 list를 return
func (pt *PeerTable) GetPeerList() []Peer{

	tmp := make([]Peer, 0)
	for _, peer := range pt.PeerMap{
		if peer.PeerID != pt.MyID{
			//내 ID는 삭제
			tmp = append(tmp, *peer)
		}
	}

	return tmp
}

//tested
func (pt *PeerTable) GetMyInfo() Peer{
	Peer, ok := pt.PeerMap[pt.MyID]

	if ok{
		return *Peer
	}

	return *Peer
}

func (pi Peer) String() string {

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