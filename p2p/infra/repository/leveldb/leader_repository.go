package leveldb

import (
	"github.com/it-chain/it-chain-Engine/common"
	"github.com/it-chain/it-chain-Engine/p2p"
	"github.com/it-chain/leveldb-wrapper"
)

type LeaderRepository struct {
	leveldb *leveldbwrapper.DB
}

func NewLeaderRepository(path string) *LeaderRepository {
	db := leveldbwrapper.CreateNewDB(path)
	db.Open()
	return &LeaderRepository{
		leveldb: db,
	}
}
func (lr *LeaderRepository) GetLeader() *p2p.Leader {
	b, err := lr.leveldb.Get([]byte("leader"))
	if err != nil {
		return nil
	}

	if len(b) == 0 {
		return nil
	}
	leader := &p2p.Leader{}
	err = common.Deserialize(b, leader)
	if err != nil {
		return nil
	}
	return leader
}

func (lr *LeaderRepository) SetLeader(leader p2p.Leader) {
	bytes, err := common.Serialize(leader)
	if err != nil {
		return
	}
	lr.leveldb.Put([]byte("leader"), bytes, true)

}
