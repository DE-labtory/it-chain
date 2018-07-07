package p2p

import (
	"encoding/json"
	"github.com/it-chain/it-chain-Engine/common"
	"fmt"
	"github.com/it-chain/midgard"
	"errors"
	"github.com/it-chain/it-chain-Engine/core/eventstore"
)

var ErrEmptyPeerId = errors.New("empty peer id requested")
var ErrEmptyAddress = errors.New("empty ip address proposed")

type PeerTable map[string]Peer

// PeerId 선언
type PeerId struct {
	Id string
}

// 노드 구조체 선언.
type Peer struct {
	IpAddress string
	PeerId    PeerId
}

func (p *Peer) On(event midgard.Event) error {

	switch v := event.(type) {
	case *PeerCreatedEvent:
		p.PeerId.Id = v.ID
		p.IpAddress = v.IpAddress

	case *PeerDeletedEvent:
		p.PeerId.Id = ""
		p.IpAddress = ""

	default:
		return errors.New(fmt.Sprintf("unhandled event [%s]", v))
	}

	return nil
}

func (n Peer) GetID() string {
	return n.PeerId.ToString()
}

// 해당 노드의 ip와 Id로 새로운 피어를 생성한다.
func NewPeer(ipAddress string, id PeerId) error {

	if id.Id == ""{
		return ErrEmptyPeerId
	}

	if ipAddress == ""{
		return ErrEmptyAddress
	}

	peer := Peer{}

	event := PeerCreatedEvent{
		EventModel: midgard.EventModel{
			ID:id.Id,
			Type: "peer.created",
		},
		IpAddress:ipAddress,

	}

	peer.On(event)

	return eventstore.Save(id.Id, event)
}

func DeletePeer(peerId PeerId) error{

	event := PeerDeletedEvent{
		EventModel: midgard.EventModel{
			ID:   peerId.Id,
			Type: "peer.deleted",
		},
	}

	return eventstore.Save(peerId.Id, event)
}

// p2p 구조체를 json 으로 인코딩한다.
func (n Peer) Serialize() ([]byte, error) {
	return common.Serialize(n)
}


// 입력받은 p2p 구조체에 해당 json 인코딩 바이트 배열을 deserialize 해서 저장한다.
func Deserialize(b []byte, peer *Peer) error {
	err := json.Unmarshal(b, peer)

	if err != nil {
		return err
	}

	return nil
}

// conver peerId to String
func (peerId PeerId) ToString() string {
	return string(peerId.Id)
}

// peer repository 인터페이스를 정의한다.
type PeerRepository interface {
	Save(peer Peer) error
	Remove(id PeerId) error
	FindById(id PeerId) (Peer, error)
	FindAll() ([]Peer, error)
	FindByAddress(ipAddress string) (Peer, error)
}

func PeerFilter(vs []Peer, f func(Peer) bool) []Peer {
	vsf := make([]Peer, 0)
	for _, v := range vs {
		if f(v) {
			vsf = append(vsf, v)
		}
	}
	return vsf
}

func GetMutuallyExclusivePeers(peers1 []Peer, peers2 []Peer) ([]Peer, []Peer) {

	exclusivePeers1 := difference(peers1, peers2)

	exclusivePeers2 := difference(peers2, peers1)

	return exclusivePeers1, exclusivePeers2
}

func difference(a, b []Peer) []Peer {
	mb := map[PeerId]bool{}

	for _, x := range b {
		mb[x.PeerId] = true
	}

	ab := []Peer{}
	for _, x := range a {
		if _, ok := mb[x.PeerId]; !ok {
			ab = append(ab, x)
		}
	}

	return ab
}
