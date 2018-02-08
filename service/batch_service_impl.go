package service

import (
	"sync"
	"time"
	"sync/atomic"
)

type Handle func(interface{})


type EventBatcherServiceImpl struct {
	Period   time.Duration
	lock     *sync.Mutex
	stopFlag int32
	buff     []*batchedMessage
	handle  Handle
	deleting bool
}

type batchedMessage struct {
	data           interface{}
	iterationsLeft int
}

//batcher T시간 간격으로 handler에게 메세지를 전달해준다
//deleting option에 따라서 전달한 message를 지울껀지 아니면 계속 남겨둘지를 설정한다.
//buff: protos queue
//lock: sync
//period: T time
//stopflag: to stop batcher
//handler: messaging handler
func NewBatchService(period time.Duration, handle Handle, deleting bool) BatchService{

	gb := &EventBatcherServiceImpl{
		buff:     make([]*batchedMessage, 0),
		lock:     &sync.Mutex{},
		Period:   period,
		stopFlag: int32(0),
		handle:  handle,
		deleting: deleting,
	}

	go gb.periodicEmit()

	return gb
}

//tested
func (gb *EventBatcherServiceImpl)Add(message interface{}){

	gb.lock.Lock()
	defer gb.lock.Unlock()

	iteration := len(gb.buff)
	gb.buff = append(gb.buff, &batchedMessage{data: message, iterationsLeft: iteration})
}

//tested
func (gb *EventBatcherServiceImpl)Stop(){
	atomic.StoreInt32(&(gb.stopFlag), int32(1))
}

//tested
func (gb *EventBatcherServiceImpl)Size() int{

	gb.lock.Lock()
	defer gb.lock.Unlock()
	return len(gb.buff)
}

//tested
func (gb *EventBatcherServiceImpl) toDie() bool {
	return atomic.LoadInt32(&(gb.stopFlag)) == int32(1)
}

//tested
func (gb *EventBatcherServiceImpl) periodicEmit() {
	for !gb.toDie() {
		time.Sleep(gb.Period)
		gb.lock.Lock()
		gb.emit()
		gb.lock.Unlock()
	}
}

//tested
func (gb *EventBatcherServiceImpl) emit() {

	if gb.toDie(){
		return
	}

	if len(gb.buff) == 0{
		return
	}

	for _, message := range gb.buff{
		gb.handle(message.data)
	}

	gb.vacate()
}

//test
func (gb *EventBatcherServiceImpl) vacate() {

	if gb.deleting{
		gb.buff = gb.buff[0:0]
	}
}