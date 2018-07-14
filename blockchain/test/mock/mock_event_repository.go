package mock

import "github.com/it-chain/midgard"

type EventRepository struct {
	LoadFunc func(aggregate midgard.Aggregate, aggregateID string) error
	SaveFunc func(aggregateID string, events ...midgard.Event) error
	CloseFunc func()
}

func (er EventRepository) Load(aggregate midgard.Aggregate, aggregateID string) error {
	return er.LoadFunc(aggregate, aggregateID)
}
func (er EventRepository) Save(aggregateID string, events ...midgard.Event) error {
	return er.SaveFunc(aggregateID, events...)
}
func (er EventRepository) Close() {
	er.CloseFunc()
}
