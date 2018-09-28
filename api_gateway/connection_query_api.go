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

package api_gateway

import (
	"errors"
	"sync"

	"github.com/it-chain/engine/common/event"
)

var ErrConnectionExists = errors.New("Connection already exists")

type ConnectionQueryApi struct {
	mux                  *sync.Mutex
	connectionRepository *ConnectionRepository
}

func NewConnectionQueryApi(connRepo *ConnectionRepository) *ConnectionQueryApi {
	return &ConnectionQueryApi{
		mux:                  &sync.Mutex{},
		connectionRepository: connRepo,
	}
}

func (q ConnectionQueryApi) GetAllConnectionList() ([]Connection, error) {
	return q.connectionRepository.FindAll()
}

func (q ConnectionQueryApi) GetConnectionByID(connID string) (Connection, error) {
	return q.connectionRepository.FindByID(connID)
}

type EventService interface {
	Publish(topic string, event interface{}) error
}

type ConnectionRepository struct {
	mux             *sync.RWMutex
	ConnectionTable map[string]Connection
	eventService    EventService
}

func NewConnectionRepository(eventService EventService) *ConnectionRepository {
	return &ConnectionRepository{
		mux:             &sync.RWMutex{},
		ConnectionTable: make(map[string]Connection),
		eventService:    eventService,
	}
}

func (cr *ConnectionRepository) Save(conn Connection) error {
	cr.mux.Lock()
	defer cr.mux.Unlock()

	_, exist := cr.ConnectionTable[conn.ConnectionID]
	if exist {
		return ErrConnectionExists
	}

	cr.ConnectionTable[conn.ConnectionID] = conn

	cr.eventService.Publish("connection.saved", createConnectionSavedEvent(conn))

	return nil
}

func createConnectionSavedEvent(conn Connection) event.ConnectionSaved {
	return event.ConnectionSaved{
		Address:      conn.ApiGatewayAddress,
		ConnectionID: conn.ConnectionID,
	}
}

func (cr *ConnectionRepository) Remove(connID string) error {
	cr.mux.Lock()
	defer cr.mux.Unlock()

	delete(cr.ConnectionTable, connID)

	return nil
}

func (cr *ConnectionRepository) FindAll() ([]Connection, error) {
	cr.mux.Lock()
	defer cr.mux.Unlock()

	connectionList := []Connection{}

	for _, conn := range cr.ConnectionTable {
		connectionList = append(connectionList, conn)
	}

	return connectionList, nil
}

func (cr *ConnectionRepository) FindByID(connID string) (Connection, error) {
	cr.mux.Lock()
	defer cr.mux.Unlock()

	for _, conn := range cr.ConnectionTable {
		if connID == conn.ConnectionID {
			return conn, nil
		}
	}

	return Connection{}, nil
}

type ConnectionEventListener struct {
	connectionRepository *ConnectionRepository
}

func NewConnectionEventListener(connRepo *ConnectionRepository) *ConnectionEventListener {
	return &ConnectionEventListener{
		connectionRepository: connRepo,
	}
}

func (cel *ConnectionEventListener) HandleConnectionCreatedEvent(event event.ConnectionCreated) error {

	connection := Connection{
		ConnectionID:       event.ConnectionID,
		GrpcGatewayAddress: event.GrpcGatewayAddress,
		ApiGatewayAddress:  event.ApiGatewayAddress,
	}

	err := cel.connectionRepository.Save(connection)
	if err != nil {
		return err
	}
	return nil
}

func (cel *ConnectionEventListener) HandleConnectionClosedEvent(event event.ConnectionClosed) error {
	cel.connectionRepository.Remove(event.ConnectionID)
	return nil
}

type Connection struct {
	ConnectionID       string
	GrpcGatewayAddress string
	ApiGatewayAddress  string
}
