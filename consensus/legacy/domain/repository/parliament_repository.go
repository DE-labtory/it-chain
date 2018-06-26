package repository

import (
	"sync"

	"github.com/it-chain/it-chain-Engine/consensus/legacy/domain/model/parliament"
)

type ParlimentRepository interface {
	Get() parliament.Parliament
	Save(parliament.Parliament) error
}

type ParlimentRepository_impl struct {
	lock       *sync.RWMutex
	parliament parliament.Parliament
}

func NewPaliamentRepository() ParlimentRepository {
	return &ParlimentRepository_impl{
		lock:       &sync.RWMutex{},
		parliament: parliament.NewParliament(),
	}
}

func (pr *ParlimentRepository_impl) Get() parliament.Parliament {
	pr.lock.Lock()
	defer pr.lock.Unlock()

	return pr.parliament
}

func (pr *ParlimentRepository_impl) Save(parliament parliament.Parliament) error {
	pr.lock.Lock()
	defer pr.lock.Unlock()

	pr.parliament = parliament

	return nil
}
