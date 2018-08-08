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

package api

import (
	"github.com/it-chain/engine/common"
	"github.com/it-chain/engine/common/event"
	"github.com/it-chain/engine/p2p"
)

type PeerApi interface {
	Save(peer p2p.Peer) error
	Remove(peerId p2p.PeerId) error
}

type PeerApiImpl struct {
	peerRepository p2p.PeerRepository
	eventService   common.EventService
}

func (ps *PeerApiImpl) Save(peer p2p.Peer) error {

	err := ps.peerRepository.Save(peer)

	if err != nil {
		return err
	}

	event := event.PeerCreated{
		PeerId:    peer.PeerId.Id,
		IpAddress: peer.IpAddress,
	}

	return ps.eventService.Publish("peer.created", event)
}

func (ps *PeerApiImpl) Remove(peerId p2p.PeerId) error {

	err := ps.peerRepository.Remove(peerId.Id)

	if err != nil {
		return err
	}

	event := event.PeerDeleted{
		PeerId: peerId.Id,
	}

	return ps.eventService.Publish("peer.deleted", event)

}
