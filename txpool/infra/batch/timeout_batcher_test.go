package batch_test

import (
	"errors"
	"testing"
	"time"

	"sync"

	"fmt"

	"github.com/it-chain/it-chain-Engine/txpool/infra/batch"
)

//todo return error 했을 경우 test할 방법이 없음
func TestTimeoutBatcher_Run(t *testing.T) {

	wg := sync.WaitGroup{}

	counter := 0
	//given
	tests := map[string]struct {
		input struct {
			taskFunc batch.TaskFunc
			duration time.Duration
		}
		err error
	}{
		"success": {
			input: struct {
				taskFunc batch.TaskFunc
				duration time.Duration
			}{
				taskFunc: func() error {
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
				taskFunc batch.TaskFunc
				duration time.Duration
			}{
				taskFunc: func() error {

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
		quit := batcher.Run(test.input.taskFunc, test.input.duration)

		fmt.Println("wait")
		wg.Wait()
		quit <- struct{}{}
	}
}
