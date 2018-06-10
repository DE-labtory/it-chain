package leveldb

import (
	"encoding/json"
	"errors"

	"github.com/it-chain/it-chain-Engine/p2p"
	"github.com/it-chain/leveldb-wrapper"
)

type NodeRepository struct {
	leveldb *leveldbwrapper.DB
}

// 새로운 p2p repo 생성
func NewNodeRepository(path string) *NodeRepository {
	db := leveldbwrapper.CreateNewDB(path)
	db.Open()
	return &NodeRepository{
		leveldb: db,
	}
}

// 새로운 node repo 생성
// 원래 코드가 임의 코드인 것 같아서 leveldb wrapper usage 참고하여 수정하였습니다.
// leveldb에 적합한 db를 제공하여 기존의 다른 함수들은 잘 동작할 것으로 예상되며 leveldbwrapper 참고해주시면 될 것 같습니다
//func NewNodeRepository(path string) *PeerRepository {
//	// path := "./leveldb"
//	dbProvider := leveldbwrapper.CreateNewDBProvider(path)
//	peerDb := dbProvider.getDBHandle("Peer")
//	return &PeerRepository{
//		leveldb: peerDb.db,
//	}
//}

// 새로운 p2p 를 leveldb에 저장
func (pr *NodeRepository) Save(data p2p.Node) error {

	// return empty peerID error if peerID is null
	if data.NodeId.ToString() == "" {
		return errors.New("node Id empty")
	}

	// serialize p2p and allocate to b or err if err occured
	b, err := data.Serialize()

	// return err if occured
	if err != nil {
		return err
	}

	// leveldb에 peerId 저장 중 에러가 나면 에러 리턴
	if err = pr.leveldb.Put([]byte(data.NodeId), b, true); err != nil {
		return err
	}

	return nil
}

// p2p 삭제
func (pr *NodeRepository) Remove(id p2p.NodeId) error {
	return pr.leveldb.Delete([]byte(id), true)
}

// p2p 읽어옴
func (pr *NodeRepository) FindById(id p2p.NodeId) (*p2p.Node, error) {
	b, err := pr.leveldb.Get([]byte(id))

	if err != nil {
		return nil, err
	}

	if len(b) == 0 {
		return nil, nil
	}

	// model.NodeRepository 에 읽어온 p2p 를 할당
	node := &p2p.Node{}

	err = json.Unmarshal(b, node)

	if err != nil {
		return nil, err
	}

	return node, nil
}

// 모든 피어 검색
func (pr *NodeRepository) FindAll() ([]p2p.Node, error) {
	iter := pr.leveldb.GetIteratorWithPrefix([]byte(""))
	var nodes []p2p.Node
	for iter.Next() {
		val := iter.Value()
		data := p2p.Node{}
		err := p2p.Deserialize(val, &data)

		if err != nil {
			return nil, err
		}

		nodes = append(nodes, data)
	}

	return nodes, nil
}
