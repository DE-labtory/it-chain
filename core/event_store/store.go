package event_store

import (
	"sync"

	"errors"

	"github.com/it-chain/midgard"
)

var ErrNilStore = errors.New("event store is nil")

var instance midgard.EventRepository

var once sync.Once

//** init function should be call once **//

//Default setting
func init() {

}

//this function is for testing
func InitForMock(repository midgard.EventRepository) {
	instance = repository
}

//
func InitMongoStore() {

}

//
func InitLevelDBStore() {

}

func RegisterEvents() {

}

func Save(aggregateID string, events ...midgard.Event) error {

	if instance == nil {
		return ErrNilStore
	}

	return instance.Save(aggregateID, events...)
}

func Load(aggregate midgard.Aggregate, aggregateID string) error {

	if instance == nil {
		return ErrNilStore
	}

	return instance.Load(aggregate, aggregateID)
}
func Close() {

	if instance == nil {
		return
	}

	instance.Close()
}
