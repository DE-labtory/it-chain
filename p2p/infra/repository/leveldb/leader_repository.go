package leveldb

import (
	"github.com/it-chain/it-chain-Engine/common"
	"github.com/it-chain/it-chain-Engine/p2p"
	"github.com/it-chain/leveldb-wrapper"
)

//todo node가 죽었다가 살아나거나, 껏다 켰을때를 고려해서 level DB로 leader repo를 구현했는데
//todo 굳이 그럴필요없이 nil이면 주변 peer에게 다시 요청해서 받아오는식으로 in memory repo 구현해도될듯
//todo 고민하고 바꿀예정
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

// get leader method
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

// set leader method
func (lr *LeaderRepository) SetLeader(leader p2p.Leader) {
	bytes, err := common.Serialize(leader)
	if err != nil {
		return
	}
	lr.leveldb.Put([]byte("leader"), bytes, true)

}
