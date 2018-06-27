package leveldb

import (
	"encoding/json"
	"errors"

	"github.com/it-chain/it-chain-Engine/p2p"
	"github.com/it-chain/leveldb-wrapper"
)

type PeerRepository struct {
	leveldb *leveldbwrapper.DB
}

// 새로운 p2p repo 생성
func NewPeerRepository(path string) *PeerRepository {
	db := leveldbwrapper.CreateNewDB(path)
	db.Open()
	return &PeerRepository{
		leveldb: db,
	}
}


// 새로운 p2p 를 leveldb에 저장
func (pr *PeerRepository) Save(data p2p.Peer) error {

	// return empty peerID error if peerID is null
	if data.PeerId.ToString() == "" {
		return errors.New("empty peer id purposed")
	}

	// serialize p2p and allocate to b or err if err occured
	b, err := data.Serialize()

	// return err if occured
	if err != nil {
		return err
	}

	// leveldb에 peerId 저장 중 에러가 나면 에러 리턴
	if err = pr.leveldb.Put([]byte(data.PeerId.Id), b, true); err != nil {
		return err
	}

	return nil
}

// p2p 삭제
func (pr *PeerRepository) Remove(id p2p.PeerId) error {
	return pr.leveldb.Delete([]byte(id.ToString()), true)
}

// p2p 읽어옴
func (pr *PeerRepository) FindById(id p2p.PeerId) (*p2p.Peer, error) {
	b, err := pr.leveldb.Get([]byte(id.ToString()))

	if err != nil {
		return nil, err
	}

	if len(b) == 0 {
		return nil, nil
	}

	// model.PeerRepository 에 읽어온 p2p 를 할당
	peer := &p2p.Peer{}

	err = json.Unmarshal(b, peer)

	if err != nil {
		return nil, err
	}

	return peer, nil
}

// 모든 피어 검색
func (pr *PeerRepository) FindAll() ([]p2p.Peer, error) {
	iter := pr.leveldb.GetIteratorWithPrefix([]byte(""))
	var peers []p2p.Peer
	for iter.Next() {
		val := iter.Value()
		data := p2p.Peer{}
		err := p2p.Deserialize(val, &data)

		if err != nil {
			return nil, err
		}

		peers = append(peers, data)
	}

	return peers, nil
}

func (nr *PeerRepository) Close(){
	nr.leveldb.Close()
}