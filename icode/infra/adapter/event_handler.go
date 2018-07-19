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
