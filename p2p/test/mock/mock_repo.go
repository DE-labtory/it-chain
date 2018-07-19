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

package mock

import "github.com/it-chain/engine/p2p"

func MakeFakePeerTable() map[string]p2p.Peer {

	peerTable := make(map[string]p2p.Peer)

	peerTable["1"] = p2p.Peer{
		PeerId: p2p.PeerId{
			Id: "1",
		},
		IpAddress: "1",
	}
	peerTable["2"] = p2p.Peer{
		PeerId: p2p.PeerId{
			Id: "2",
		},
		IpAddress: "2",
	}
	peerTable["3"] = p2p.Peer{
		PeerId: p2p.PeerId{
			Id: "3",
		},
		IpAddress: "3",
	}

	return peerTable
}

func MakeFakePLTable() p2p.PLTable {

	peerTable := MakeFakePeerTable()
	leader := p2p.Leader{
		LeaderId: p2p.LeaderId{
			Id: "1",
		},
	}

	return p2p.PLTable{
		Leader:    leader,
		PeerTable: peerTable,
	}
}
