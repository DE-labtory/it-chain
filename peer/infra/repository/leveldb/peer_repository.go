package leveldb

import (
	"encoding/json"

	"github.com/it-chain/it-chain-Engine/peer/domain/model"
	"github.com/it-chain/it-chain-Engine/peer/domain/repository"
	"github.com/it-chain/leveldb-wrapper"
)

type PeerRepository struct {
	leveldb *leveldbwrapper.DB
}

func NewPeerRepository(path string) *PeerRepository {
	db := leveldbwrapper.CreateNewDB(path)
	db.Open()
	return &PeerRepository{
		leveldb: db,
	}
}

func (pr *PeerRepository) Save(peer model.Peer) error {
	if peer.Id.ToString() == "" {
		return repository.PeerIdEmptyErr
	}

	b, err := peer.Serialize()

	if err != nil {
		return err
	}

	if err = pr.leveldb.Put([]byte(peer.Id), b, true); err != nil {
		return err
	}

	return nil
}

func (pr *PeerRepository) Remove(id model.PeerId) error {
	return pr.leveldb.Delete([]byte(id), true)
}

func (pr *PeerRepository) FindById(id model.PeerId) (*model.Peer, error) {
	b, err := pr.leveldb.Get([]byte(id))

	if err != nil {
		return nil, err
	}

	if len(b) == 0 {
		return nil, nil
	}

	peer := &model.Peer{}

	err = json.Unmarshal(b, peer)

	if err != nil {
		return nil, err
	}

	return peer, nil
}

func (pr *PeerRepository) FindAll() ([]*model.Peer, error) {
	iter := pr.leveldb.GetIteratorWithPrefix([]byte(""))
	peers := []*model.Peer{}
	for iter.Next() {
		val := iter.Value()
		peer := &model.Peer{}
		err := model.Deserialize(val, peer)

		if err != nil {
			return nil, err
		}

		peers = append(peers, peer)
	}

	return peers, nil
}
