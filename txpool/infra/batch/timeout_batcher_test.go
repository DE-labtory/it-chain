package batch_test

import (
	"testing"

	"time"

	"log"

	"sync"

	"github.com/it-chain/it-chain-Engine/txpool/infra/batch"
)

func TestGetTimeOutBatcherInstance(t *testing.T) {

	wg := sync.WaitGroup{}
	wg.Add(1)

	batcher := batch.GetTimeOutBatcherInstance()

	testValue := 0

	mockFunc := func() error {
		testValue = 1
		log.Println(testValue)
		return nil
	}

	quit := batcher.Register(mockFunc, time.Second*2)

	time.Sleep(4 * time.Second)

	quit <- struct{}{}

	//wg.Wait()
}
