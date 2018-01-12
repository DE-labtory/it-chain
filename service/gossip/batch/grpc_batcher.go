package batch

import (
	"sync"
	"time"
	"sync/atomic"
)

type MessageHandler interface{
	handle([]byte) (interface{},error)
}

type GRPCBatcher struct {
	Period   time.Duration
	lock     *sync.Mutex
	stopFlag int32
	buff     []*batchedMessage
	handler  MessageHandler
	deleting bool
}

type batchedMessage struct {
	data           interface{}
	iterationsLeft int
}

func NewGRPCMessageBatcher(period time.Duration, handler MessageHandler, deleting bool) *GRPCBatcher{

	gb := &GRPCBatcher{
		buff:     make([]*batchedMessage, 0),
		lock:     &sync.Mutex{},
		Period:   period,
		stopFlag: int32(0),
		handler:  handler,
		deleting: deleting,
	}
	return gb
}


func (gb *GRPCBatcher)Add(message interface{}){

	gb.lock.Lock()
	defer gb.lock.Unlock()

	interation := len(gb.buff)
	gb.buff = append(gb.buff, &batchedMessage{data: message, iterationsLeft: interation})

}

// Stop stops the component
func (gb *GRPCBatcher)Stop(){
	atomic.StoreInt32(&(gb.stopFlag), int32(1))
}

// Size returns the amount of pending messages to be emitted
func (gb *GRPCBatcher)Size() int{

	gb.lock.Lock()
	defer gb.lock.Unlock()
	return len(gb.buff)
}

func (gb *GRPCBatcher) toDie() bool {
	return atomic.LoadInt32(&(gb.stopFlag)) == int32(1)
}

func (gb *GRPCBatcher) decrementCounters() {
	n := len(gb.buff)
	for i := 0; i < n; i++ {
		msg := gb.buff[i]
		msg.iterationsLeft--
		if msg.iterationsLeft == 0 {
			gb.buff = append(gb.buff[:i], gb.buff[i+1:]...)
			n--
			i--
		}
	}
}

func (gb *GRPCBatcher) periodicEmit() {
	for !gb.toDie() {
		time.Sleep(gb.Period)
		gb.lock.Lock()
		gb.emit()
		gb.lock.Unlock()
	}
}

func (gb *GRPCBatcher) emit() {
	//if gb.toDie() {
	//	return
	//}
	//if len(gb.buff) == 0 {
	//	return
	//}
	//msgs2beEmitted := make([]interface{}, len(gb.buff))
	//for i, v := range gb.buff {
	//	msgs2beEmitted[i] = v.data
	//}
	//
	//gb.handler.handle()
	//gb.decrementCounters()
}