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

import "github.com/it-chain/midgard"

//publish
type LeaderChangedEvent struct {
	midgard.EventModel
}

//handle
type ConnectionCreatedEvent struct {
	midgard.EventModel
	Address string
}

//handle
type ConnectionDisconnectedEvent struct {
	midgard.EventModel
}

// node created event
type PeerCreatedEvent struct {
	midgard.EventModel
	IpAddress string
}

type PeerDeletedEvent struct {
	midgard.EventModel
}

// handle leader received event
type LeaderUpdatedEvent struct {
	midgard.EventModel
}

type LeaderDeliveredEvent struct {
	midgard.EventModel
}

//todo add to event doc
type LeaderDeletedEvent struct {
	midgard.EventModel
}
