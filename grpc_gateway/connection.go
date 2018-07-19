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

import (
	"errors"
	"fmt"

	"github.com/it-chain/engine/core/eventstore"
	"github.com/it-chain/midgard"
)

type Connection struct {
	midgard.AggregateModel
	Address string
}

func (c *Connection) On(event midgard.Event) error {
	switch v := event.(type) {

	case *ConnectionCreatedEvent:
		c.ID = v.ID
		c.Address = v.Address

	case *ConnectionClosedEvent:
		c.ID = ""
		c.Address = ""

	default:
		return errors.New(fmt.Sprintf("unhandled event [%s]", v))
	}

	return nil
}

func NewConnection(connectionID string, address string) (Connection, error) {

	c := Connection{}

	event := &ConnectionCreatedEvent{
		EventModel: midgard.EventModel{
			ID:   connectionID,
			Type: "connection.created",
		},
		Address: address,
	}

	c.On(event)

	return c, eventstore.Save(connectionID, event)
}

func CloseConnection(connectionID string) error {

	return eventstore.Save(connectionID, ConnectionClosedEvent{
		EventModel: midgard.EventModel{
			ID:   connectionID,
			Type: "connection.closed",
		},
	})
}
