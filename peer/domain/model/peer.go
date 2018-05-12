package model

import (
	"encoding/json"

	"github.com/it-chain/it-chain-Engine/common"
)

type Peer struct {
	IpAddress string
	Id        PeerId
}

func NewPeer(ipAddress string, id PeerId) *Peer {
	return &Peer{
		IpAddress: ipAddress,
		Id:        id,
	}
}

func (p Peer) Serialize() ([]byte, error) {
	return common.Serialize(p)
}

func Deserialize(b []byte, peer *Peer) error {
	err := json.Unmarshal(b, peer)

	if err != nil {
		return err
	}

	return nil
}

type PeerId string

func (peerId PeerId) ToString() string {
	return string(peerId)
}
