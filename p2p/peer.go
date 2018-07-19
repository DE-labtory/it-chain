/*
 * Copyright 2018 It-chain
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * https://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package p2p

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/it-chain/engine/common"
	"github.com/it-chain/engine/core/eventstore"
	"github.com/it-chain/midgard"
)

var ErrEmptyPeerId = errors.New("empty peer id requested")
var ErrEmptyAddress = errors.New("empty ip address proposed")
var ErrNoMatchingPeerId = errors.New("no matching peer id")

// 노드 구조체 선언.
type Peer struct {
	IpAddress string
	PeerId    PeerId
}

// PeerId 선언
type PeerId struct {
	Id string
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

	if id.Id == "" {
		return ErrEmptyPeerId
	}

	if ipAddress == "" {
		return ErrEmptyAddress
	}

	peer := Peer{}

	event := PeerCreatedEvent{
		EventModel: midgard.EventModel{
			ID:   id.Id,
			Type: "peer.created",
		},
		IpAddress: ipAddress,
	}

	peer.On(event)

	return eventstore.Save(id.Id, event)
}

func DeletePeer(peerId PeerId) error {

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

func PeerFilter(vs []Peer, f func(Peer) bool) []Peer {
	vsf := make([]Peer, 0)
	for _, v := range vs {
		if f(v) {
			vsf = append(vsf, v)
		}
	}
	return vsf
}
