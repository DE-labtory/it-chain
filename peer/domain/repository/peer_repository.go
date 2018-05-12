package repository

import (
	"errors"

	"github.com/it-chain/it-chain-Engine/peer/domain/model"
)

var PeerNotFoundErr = errors.New("PeerNotFound")
var DuplicatePeerErr = errors.New("PeerAlreadyExist")
var PeerIdEmptyErr = errors.New("PeerIdEmpty")

type Peer interface {
	Save(peer model.Peer) error
	Remove(id model.PeerId) error
	FindById(id model.PeerId) (*model.Peer, error)
	FindAll() ([]*model.Peer, error)
}
