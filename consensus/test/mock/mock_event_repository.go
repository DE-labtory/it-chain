package mock

import "github.com/it-chain/midgard"

type MockEventRepository struct {
	LoadFunc  func(aggregate midgard.Aggregate, aggregateID string) error
	SaveFunc  func(aggregateID string, events ...midgard.Event) error
	CloseFunc func()
}

func (er MockEventRepository) Load(aggregate midgard.Aggregate, aggregateID string) error {
	return er.LoadFunc(aggregate, aggregateID)
}
func (er MockEventRepository) Save(aggregateID string, events ...midgard.Event) error {
	return er.SaveFunc(aggregateID, events...)
}
func (er MockEventRepository) Close() {
	er.CloseFunc()
}
