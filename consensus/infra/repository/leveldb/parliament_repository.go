package leveldb

import (
	"sync"

	"github.com/it-chain/leveldb-wrapper"
)

type ParliamentRepository struct {
	leveldb *leveldbwrapper.DB
	lock *sync.RWMutex
	// todo : parliament struct 추가
	// parliament
}

func NewParliamentRepository(path string) *ParliamentRepository {
	db := leveldbwrapper.CreateNewDB(path)
	db.Open()
	return &ParliamentRepository{
		leveldb: db,
		lock: &sync.RWMutex{},
		// todo : parliament struct 추가
		// parliament: parliament.NewParliament(),
	}
}

