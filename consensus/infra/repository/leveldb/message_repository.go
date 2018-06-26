package leveldb

import (
	"github.com/it-chain/leveldb-wrapper"
)

type MessageRepository struct {
	leveldb *leveldbwrapper.DB
}

func NewMessageRepository(path string) *MessageRepository {
	db := leveldbwrapper.CreateNewDB(path)
	db.Open()
	return &MessageRepository{
		leveldb: db,
	}
}

