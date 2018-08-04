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

type PeerService struct {
	peerRepository PeerRepository
	publishService PublishService
}

func NewPeerService(peerRepository PeerRepository, publishService PublishService) PeerService {
	return PeerService{
		peerRepository: peerRepository,
		publishService: publishService,
	}
}

func (ps *PeerService) Save(peer Peer) error {

	if peer.IpAddress == "" {
		return ErrEmptyAddress
	}

	ps.peerRepository.Save(peer)

	ps.publishService.PeerCreated(peer)

	return nil
}

func (ps *PeerService) Remove(peerId PeerId) error {

	ps.peerRepository.Delete(peerId.Id)

	ps.publishService.PeerDeleted(peerId)

	return nil
}
