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

package grpc_gateway

import "github.com/it-chain/midgard"

type ConnectionCreatedEvent struct {
	midgard.EventModel
	Address string
}

type ConnectionClosedEvent struct {
	midgard.EventModel
}

type ErrorCreatedEvent struct {
	midgard.EventModel
	Event string
	Err   string
}

type EventRepository interface {
	Load(aggregate midgard.Aggregate, aggregateID string) error
	Save(aggregateID string, events ...midgard.Event) error
}
