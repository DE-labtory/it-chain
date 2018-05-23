package repository

// Peer Repository는 peer를 leveldb에 관한 CRUD를 담당한다.

import (
	"errors"

	"github.com/it-chain/it-chain-Engine/peer/domain/model"
)

var PeerNotFoundErr = errors.New("PeerNotFound")
var DuplicatePeerErr = errors.New("PeerAlreadyExist")
var PeerIdEmptyErr = errors.New("PeerIdEmpty")

// peer repository 인터페이스를 정의한다.
// Peer 가 아니라 PeerRepository 로 정의하는 것이 맞을것 같습니다. - 남훈

type Peer interface {
	Save(peer model.Peer) error
	Remove(id model.PeerId) error
	FindById(id model.PeerId) (*model.Peer, error)
	FindAll() ([]*model.Peer, error)
}
