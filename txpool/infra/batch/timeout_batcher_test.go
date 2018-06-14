package batch_test

import (
	"errors"
	"testing"
	"time"

	"sync"

	"fmt"

	"github.com/it-chain/it-chain-Engine/txpool/infra/batch"
)

func TestTimeoutBatcher_Register(t *testing.T) {

	wg := sync.WaitGroup{}

	counter := 0
	//given
	tests := map[string]struct {
		input struct {
			timerFunc batch.TimerFunc
			duration  time.Duration
		}
		err error
	}{
		"success": {
			input: struct {
				timerFunc batch.TimerFunc
				duration  time.Duration
			}{
				timerFunc: func() error {
					if counter == 0 {
						wg.Done()
						fmt.Println("success done")
					}

					counter++

					return nil
				},
				duration: time.Second * 1,
			},
			err: nil,
		},
		"timer return error": {
			input: struct {
				timerFunc batch.TimerFunc
				duration  time.Duration
			}{
				timerFunc: func() error {

					if counter == 0 {
						wg.Done()
						fmt.Println("done")
					}

					counter++

					return errors.New("asd")
				},
				duration: time.Second * 1,
			},
			err: nil,
		},
	}

	batcher := batch.GetTimeOutBatcherInstance()

	for testName, test := range tests {
		t.Logf("Running test case %s", testName)

		//when
		wg.Add(1)
		counter = 0
		quit := batcher.Register(test.input.timerFunc, test.input.duration)

		fmt.Println("wait")
		wg.Wait()
		quit <- struct{}{}
	}
}
