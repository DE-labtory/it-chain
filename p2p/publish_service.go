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

import "github.com/it-chain/engine/common/event"

type Publish func(topic string, data interface{}) (err error) // 나중에 의존성 주입을 해준다.

type PublishService interface {
	PeerCreated(peer Peer) error
	PeerDeleted(peerId PeerId) error
	LeaderUpdated(leader Leader) error
}

type PublishServiceImpl struct {
	publish Publish
}

func (ps *PublishServiceImpl) PeerCreated(peer Peer) error {

	event := event.PeerCreated{
		PeerId:    peer.PeerId.Id,
		IpAddress: peer.IpAddress,
	}

	return ps.publish("peer.created", event)
}

func (ps *PublishServiceImpl) PeerDeleted(peerId PeerId) error {

	event := event.PeerDeleted{
		PeerId: peerId.Id,
	}

	return ps.publish("peer.deleted", event)
}

func (ps *PublishServiceImpl) LeaderUpdated(leader Leader) error {

	event := event.LeaderUpdated{
		LeaderId: leader.LeaderId.Id,
	}

	return ps.publish("leader.updated", event)
}
