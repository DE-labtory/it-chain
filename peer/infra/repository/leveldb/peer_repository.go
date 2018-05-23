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

// 새로운 peer repo 생성
func NewPeerRepository(path string) *PeerRepository {
	db := leveldbwrapper.CreateNewDB(path)
	db.Open()
	return &PeerRepository{
		leveldb: db,
	}
}


// 새로운 peer 를 leveldb에 저장
func (pr *PeerRepository) Save(peer model.Peer) error {

	// return empty peerID error if peerID is null
	if peer.Id.ToString() == "" {
		return repository.PeerIdEmptyErr
	}

	// serialize peer and allocate to b or err if err occured
	b, err := peer.Serialize()

	// return err if occured
	if err != nil {
		return err
	}

	// leveldb에 peerId 저장 중 에러가 나면 에러 리턴
	if err = pr.leveldb.Put([]byte(peer.Id), b, true); err != nil {
		return err
	}

	return nil
}


// peer 삭제
func (pr *PeerRepository) Remove(id model.PeerId) error {
	return pr.leveldb.Delete([]byte(id), true)
}


// peer 읽어옴
func (pr *PeerRepository) FindById(id model.PeerId) (*model.Peer, error) {
	b, err := pr.leveldb.Get([]byte(id))

	if err != nil {
		return nil, err
	}

	if len(b) == 0 {
		return nil, nil
	}

	// model.Peer 에 읽어온 peer 를 할당
	peer := &model.Peer{}

	err = json.Unmarshal(b, peer)

	if err != nil {
		return nil, err
	}

	return peer, nil
}


// 모든 피어 검색
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
