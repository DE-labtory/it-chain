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

package adapter

import (
	"errors"

	"github.com/it-chain/engine/icode"
)

type EventHandler struct {
	iCodeMetaRepository icode.MetaRepository
}

func NewEventHandler(ICodeMetaRepository icode.MetaRepository) *EventHandler {
	return &EventHandler{
		iCodeMetaRepository: ICodeMetaRepository,
	}
}

func (handler EventHandler) HandleMetaCreatedEvent(event icode.MetaCreatedEvent) error {
	if event.ID == "" {
		return errors.New("Empty event id err")
	}
	m := icode.Meta{}
	m.On(event)
	return handler.iCodeMetaRepository.Save(m)
}

func (handler EventHandler) HandleMetaDeletedEvent(event icode.MetaDeletedEvent) error {
	if event.ID == "" {
		return errors.New("Empty event id err")
	}
	return handler.iCodeMetaRepository.Remove(event.ID)
}
