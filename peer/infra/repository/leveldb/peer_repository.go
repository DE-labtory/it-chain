package leveldb

import (
	"encoding/json"

	"github.com/it-chain/leveldb-wrapper"

)

type PeerRepository struct {
	leveldb *leveldbwrapper.DB
}

// 새로운 peer repo 생성
// 원래 코드가 임의 코드인 것 같아서 leveldb wrapper usage 참고하여 수정하였습니다.
// leveldb에 적합한 db를 제공하여 기존의 다른 함수들은 잘 동작할 것으로 예상되며 leveldbwrapper 참고해주시면 될 것 같습니다
func NewPeerRepository(path string) *PeerRepository {
	// path := "./leveldb"
	dbProvider := leveldbwrapper.CreateNewDBProvider(path)
	dbHandle := dbProvider.getDBHandle("Peer")
	return &PeerRepository{
		leveldb: dbHandle.db,
	}
}

// 새로운 peer 를 leveldb에 저장
func (pr *PeerRepository) Save(peer Peer) error {

	// return empty peerID error if peerID is null
	if peer.Id.ToString() == "" {
		return PeerIdEmptyErr
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
func (pr *PeerRepository) Remove(id PeerId) error {
	return pr.leveldb.Delete([]byte(id), true)
}

// peer 읽어옴
func (pr *PeerRepository) FindById(id PeerId) (*Peer, error) {
	b, err := pr.leveldb.Get([]byte(id))

	if err != nil {
		return nil, err
	}

	if len(b) == 0 {
		return nil, nil
	}

	// model.Peer 에 읽어온 peer 를 할당
	peer := &Peer{}

	err = json.Unmarshal(b, peer)

	if err != nil {
		return nil, err
	}

	return peer, nil
}

// 모든 피어 검색
func (pr *PeerRepository) FindAll() ([]*Peer, error) {
	iter := pr.leveldb.GetIteratorWithPrefix([]byte(""))
	peers := []*Peer{}
	for iter.Next() {
		val := iter.Value()
		peer := &Peer{}
		err := Deserialize(val, peer)

		if err != nil {
			return nil, err
		}

		peers = append(peers, peer)
	}

	return peers, nil
}
