package model

import (
	"encoding/json"

	"github.com/it-chain/it-chain-Engine/common"
)

// 피어 구조체를 선언한다.
type Peer struct {
	IpAddress string
	Id        PeerId
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
