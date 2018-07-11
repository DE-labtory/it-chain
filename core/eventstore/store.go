package eventstore

import (
	"errors"
	"sync"

	"github.com/it-chain/midgard"
	"github.com/it-chain/midgard/bus/rabbitmq"
	"github.com/it-chain/midgard/store"
	"github.com/it-chain/midgard/store/leveldb"
	"github.com/it-chain/midgard/store/mongodb"
)

var ErrNilStore = errors.New("event store is nil")

var Instance *Store

type Store struct {
	repo midgard.EventRepository

	//serializer는 event struct의 list를가지고 있고 byte[]로 저장된 event를 deserialize할때 등록된 event로 deserialize한다.
	//serializer는 store의 내부에서 작동하므로, 동적으로 event를 등록하기 위해 따로 밖에서 instance를 가지고 있다.
	//db에서 event를 복구하기 위해서는 꼭 event가 등록 되어있어야 하므로, RegisterEvents함수를 이용해 저장하는 event들을 등록해야한다.
	serializer store.EventSerializer
}

var once sync.Once

//** init function should be cavar serializer ll once **//

//Default setting
func InitDefault() {

	store, serializer := initDefaultStore()
	publisher := initDefaultPublisher()
	repo := midgard.NewRepo(store, publisher)

	Instance = &Store{
		repo:       repo,
		serializer: serializer,
	}
}

//todo path, dbname from viper
func initDefaultStore() (midgard.EventStore, store.EventSerializer) {
	if Instance != nil {
		panic("eventstore is already initialized")
	}

	path := "mongodb://localhost:27017"
	dbname := "test"

	serializer := store.NewSerializer()

	mongoStore, err := mongodb.NewEventStore(path, dbname, serializer)

	if err != nil {
		panic(err.Error())
	}

	return mongoStore, serializer
}

//todo rabbitmq url
func initDefaultPublisher() midgard.Publisher {

	client := rabbitmq.Connect("")

	return client
}

//this function is for testing
func InitForMock(repository midgard.EventRepository) {
	Instance = &Store{
		repo: repository,
	}
}

//todo CustomMongoStore init part
func InitMongoStore(path string, dbname string, publisher midgard.Publisher, events ...midgard.Event) {
	if Instance != nil {
		panic("eventstore is already initialized")
	}
}

//
func InitLevelDBStore(path string, publisher midgard.Publisher, events ...midgard.Event) {
	if Instance != nil {
		panic("eventstore is already initialized")
	}

	serializer := store.NewSerializer(events...)
	store := leveldb.NewEventStore(path, serializer)

	Instance = &Store{
		repo:       midgard.NewRepo(store, publisher),
		serializer: serializer,
	}
}

func RegisterEvents(events ...midgard.Event) error {

	if Instance.serializer == nil {
		return errors.New("nil event register")
	}

	for _, event := range events {
		Instance.serializer.Register(event)
	}

	return nil
}

func Save(aggregateID string, events ...midgard.Event) error {

	if Instance == nil {
		return ErrNilStore
	}

	return Instance.repo.Save(aggregateID, events...)
}

func Load(aggregate midgard.Aggregate, aggregateID string) error {

	if Instance == nil {
		return ErrNilStore
	}

	return Instance.repo.Load(aggregate, aggregateID)
}
func Close() {

	if Instance == nil {
		return
	}

	Instance.repo.Close()
	Instance = nil
}
