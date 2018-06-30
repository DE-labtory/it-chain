package adapter

import (
	"github.com/it-chain/it-chain-Engine/icode"
	"github.com/pkg/errors"
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
	return handler.iCodeMetaRepository.Save(*event.GetMeta())
}

func (handler EventHandler) HandleMetaDeletedEvent(event icode.MetaDeletedEvent) error {
	if event.ID == "" {
		return errors.New("Empty event id err")
	}
	return handler.iCodeMetaRepository.Remove(event.ID)
}
