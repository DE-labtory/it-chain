package repository

import (
	"sync"

	"github.com/it-chain/it-chain-Engine/consensus/domain/model/parliament"
)

type ParlimentRepository interface {
	Get() parliament.Parliament
	Save(parliament.Parliament) error
}

type ParlimentRepository_impl struct {
	lock       sync.Mutex
	parliament parliament.Parliament
}

func (pr *ParlimentRepository_impl) Get() parliament.Parliament {
	pr.lock.Lock()
	defer pr.lock.Unlock()

	return pr.parliament
}

func (pr *ParlimentRepository_impl) Insert(parliament parliament.Parliament) error {
	pr.lock.Lock()
	defer pr.lock.Unlock()

	pr.parliament = parliament

	return nil
}
