package leveldb

import (
	"encoding/json"

	"github.com/it-chain/it-chain-Engine/peer"
	"github.com/it-chain/leveldb-wrapper"
)

type NodeRepository struct {
	leveldb *leveldbwrapper.DB
}

// 새로운 peer repo 생성
func NewNodeRepository(path string) *NodeRepository {
	db := leveldbwrapper.CreateNewDB(path)
	db.Open()
	return &NodeRepository{
		leveldb: db,
	}
}

// 새로운 peer 를 leveldb에 저장
func (pr *NodeRepository) Save(data peer.Node) error {

	// return empty peerID error if peerID is null
	if data.Id.ToString() == "" {
		return peer.NodeIdEmptyErr
	}

	// serialize peer and allocate to b or err if err occured
	b, err := data.Serialize()

	// return err if occured
	if err != nil {
		return err
	}

	// leveldb에 peerId 저장 중 에러가 나면 에러 리턴
	if err = pr.leveldb.Put([]byte(data.Id), b, true); err != nil {
		return err
	}

	return nil
}

// peer 삭제
func (pr *NodeRepository) Remove(id peer.NodeId) error {
	return pr.leveldb.Delete([]byte(id), true)
}

// peer 읽어옴
func (pr *NodeRepository) FindById(id peer.NodeId) (*peer.Node, error) {
	b, err := pr.leveldb.Get([]byte(id))

	if err != nil {
		return nil, err
	}

	if len(b) == 0 {
		return nil, nil
	}

	// model.NodeRepository 에 읽어온 peer 를 할당
	node := &peer.Node{}

	err = json.Unmarshal(b, node)

	if err != nil {
		return nil, err
	}

	return node, nil
}

// 모든 피어 검색
func (pr *NodeRepository) FindAll() ([]*peer.Node, error) {
	iter := pr.leveldb.GetIteratorWithPrefix([]byte(""))
	var nodes []*peer.Node
	for iter.Next() {
		val := iter.Value()
		data := &peer.Node{}
		err := peer.Deserialize(val, data)

		if err != nil {
			return nil, err
		}

		nodes = append(nodes, data)
	}

	return nodes, nil
}
