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
	"github.com/it-chain/engine/p2p"
	"github.com/it-chain/engine/p2p/infra/mem"
)

type PeerApi struct {
	peerRepository mem.PeerRepository
	publishService p2p.PublishService
}

func (ps *PeerApi) Save(peer p2p.Peer) {

	ps.peerRepository.Save(peer)

	ps.publishService.PeerCreated(peer)
}

func (ps *PeerApi) Remove(peerId p2p.PeerId) {

	ps.peerRepository.Delete(peerId.Id)

	ps.publishService.PeerDeleted(peerId)
}
