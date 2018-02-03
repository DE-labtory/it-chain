package sample_smartcontract

import (
	"github.com/revel/modules/db/app"
)

type SampleSmartContract struct {

}

func (sc *SampleSmartContract) Init() {

}

func (sc *SampleSmartContract) Query(method string) {
	if method == "getA" {
		getA()
	}
}

func (sc *SampleSmartContract) Invoke() {

}

func getA() {
	return db.get('a')
}