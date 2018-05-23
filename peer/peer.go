package model

import (
	"encoding/json"
	"errors"

	"github.com/it-chain/it-chain-Engine/common"
	"github.com/it-chain/it-chain-Engine/peer/domain/model"
)

var PeerNotFoundErr = errors.New("PeerNotFound")
var DuplicatePeerErr = errors.New("PeerAlreadyExist")
var PeerIdEmptyErr = errors.New("PeerIdEmpty")

// 피어 구조체를 선언한다.
type Peer struct {
	IpAddress string
	Id        PeerId
}

// peer repository 인터페이스를 정의한다.
// Peer 가 아니라 PeerRepository 로 정의했습니다. - 남훈

type PeerRepository interface {
	Save(peer model.Peer) error
	Remove(id model.PeerId) error
	FindById(id model.PeerId) (*model.Peer, error)
	FindAll() ([]*model.Peer, error)
}

// 해당 피어의 ip와 peerId로 새로운 피어를 생성한다.
func NewPeer(ipAddress string, id PeerId) *Peer {
	return &Peer{
		IpAddress: ipAddress,
		Id:        id,
	}
}

// peer 구조체를 json 으로 인코딩한다.
func (p Peer) Serialize() ([]byte, error) {
	return common.Serialize(p)
}

// 입력받은 peer 구조체에 해당 json 인코딩 바이트 배열을 deserialize 해서 저장한다.
func Deserialize(b []byte, peer *Peer) error {
	err := json.Unmarshal(b, peer)

	if err != nil {
		return err
	}

	return nil
}

// peerId 선언
type PeerId string

// conver peerId to String
func (peerId PeerId) ToString() string {
	return string(peerId)
}
