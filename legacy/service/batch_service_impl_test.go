package service

import (
	"testing"
	"time"
	"github.com/stretchr/testify/assert"
	"sync"
)

type MockHandler struct {
	counter int
	done    chan int
}

var mock_iter = 3

func (mh *MockHandler) Handle(ms interface{}){

	mh.counter = mh.counter +1

	if mh.counter > mock_iter{
		mh.done <- 3
	}
}

func MockNewGRPCMessageBatcher(done chan int) (*EventBatcherServiceImpl, *MockHandler){

	seconds := 2 * time.Second
	mockHanlder := &MockHandler{counter:0,done:done}

	return &EventBatcherServiceImpl{
		buff:     make([]*batchedMessage, 0),
		lock:     &sync.Mutex{},
		Period:   seconds,
		stopFlag: int32(0),
		handle:  mockHanlder.Handle,
		deleting: false,
	}, mockHanlder
}

func TestNewGRPCMessageBatcher(t *testing.T) {

	mockHanlder := &MockHandler{}
	seconds := 2 * time.Second
	batcher := &EventBatcherServiceImpl{
		buff:     make([]*batchedMessage, 0),
		lock:     &sync.Mutex{},
		Period:   seconds,
		stopFlag: int32(0),
		handle:  mockHanlder.Handle,
		deleting: false,
	}

	assert.Equal(t,batcher.stopFlag,int32(0))
	assert.Equal(t,batcher.deleting,false)
	assert.Equal(t,len(batcher.buff),0)

}

func TestEventBatcher_Add(t *testing.T) {

	done := make(chan int)
	batcher,_ := MockNewGRPCMessageBatcher(done)
	batcher.Add("hello world!")

	assert.Equal(t,len(batcher.buff),1)
	assert.Equal(t,batcher.buff[0].iterationsLeft,0)
}

func TestEventBatcher_periodic_emit(t *testing.T){

	done := make(chan int)
	batcher,handler := MockNewGRPCMessageBatcher(done)
	batcher.buff = append(batcher.buff, &batchedMessage{data: "hello", iterationsLeft: len(batcher.buff)})

	go batcher.periodicEmit()

	<-done
	assert.Equal(t,handler.counter,4)
}

func TestEventBatcher_Stop(t *testing.T) {

	done := make(chan int)
	batcher,handler := MockNewGRPCMessageBatcher(done)
	batcher.buff = append(batcher.buff, &batchedMessage{data: "hello", iterationsLeft: len(batcher.buff)})

	go batcher.periodicEmit()

	time.Sleep(3*time.Second)

	batcher.Stop()
	assert.Equal(t,handler.counter,1)
}

func TestEventBatcher_Size(t *testing.T) {
	done := make(chan int)
	batcher,_ := MockNewGRPCMessageBatcher(done)
	batcher.buff = append(batcher.buff, &batchedMessage{data: "hello", iterationsLeft: len(batcher.buff)})
	assert.Equal(t,batcher.Size(),1)
}

func TestEventBatcher_FuntionalTestWithDeletingTrue(t *testing.T) {

	done := make(chan int)
	mockHanlder := &MockHandler{counter:0,done:done}
	seconds := 2 * time.Second
	batcher := NewBatchService(seconds, mockHanlder.Handle,true)

	batcher.Add("hello1")
	batcher.Add("hello2")
	batcher.Add("hello3")
	batcher.Add("hello4")

	assert.Equal(t,batcher.Size(),4)

	<-done
	batcher.Stop()

	assert.Equal(t,batcher.Size(),0)
}

func TestEventBatcher_FuntionalTest_with_deleting_false(t *testing.T) {

	done := make(chan int)
	mockHanlder := &MockHandler{counter:0,done:done}
	seconds := 2 * time.Second
	batcher := NewBatchService(seconds, mockHanlder.Handle,false)

	batcher.Add("hello1")
	batcher.Add("hello2")
	batcher.Add("hello3")
	batcher.Add("hello4")

	assert.Equal(t,batcher.Size(),4)

	<-done
	batcher.Stop()

	assert.Equal(t,batcher.Size(),4)
}