package batch

import (
	"testing"
	"time"
	"github.com/magiconair/properties/assert"
	"sync"
)

type MockHandler struct {
	counter int
	done    chan int
}

var mock_iter = 3

func (mh *MockHandler) handle(ms interface{}) (interface{},error){

	mh.counter = mh.counter +1

	if mh.counter > mock_iter{
		mh.done <- 3
	}
	return "success", nil
}

func MockNewGRPCMessageBatcher(done chan int) (*GRPCBatcher, *MockHandler){

	seconds := 2 * time.Second
	mockHanlder := &MockHandler{counter:0,done:done}

	return &GRPCBatcher{
		buff:     make([]*batchedMessage, 0),
		lock:     &sync.Mutex{},
		Period:   seconds,
		stopFlag: int32(0),
		handler:  mockHanlder,
		deleting: false,
	}, mockHanlder
}

func TestNewGRPCMessageBatcher(t *testing.T) {

	mockHanlder := &MockHandler{}
	seconds := 2 * time.Second
	batcher := NewGRPCMessageBatcher(seconds, mockHanlder,false)

	assert.Equal(t,batcher.stopFlag,int32(0))
	assert.Equal(t,batcher.deleting,false)
	assert.Equal(t,len(batcher.buff),0)

}

func TestGRPCBatcher_Add(t *testing.T) {

	done := make(chan int)
	batcher,_ := MockNewGRPCMessageBatcher(done)
	batcher.Add("hello world!")

	assert.Equal(t,len(batcher.buff),1)
	assert.Equal(t,batcher.buff[0].iterationsLeft,0)
}

func TestGRPCBatcher_periodic_emit(t *testing.T){

	done := make(chan int)
	batcher,handler := MockNewGRPCMessageBatcher(done)
	batcher.buff = append(batcher.buff, &batchedMessage{data: "hello", iterationsLeft: len(batcher.buff)})

	go batcher.periodicEmit()

	<-done
	assert.Equal(t,handler.counter,4)
}

func TestGRPCBatcher_Stop(t *testing.T) {

	done := make(chan int)
	batcher,handler := MockNewGRPCMessageBatcher(done)
	batcher.buff = append(batcher.buff, &batchedMessage{data: "hello", iterationsLeft: len(batcher.buff)})

	go batcher.periodicEmit()

	time.Sleep(3*time.Second)

	batcher.Stop()
	assert.Equal(t,handler.counter,1)
}

func TestGRPCBatcher_Size(t *testing.T) {
	done := make(chan int)
	batcher,_ := MockNewGRPCMessageBatcher(done)
	batcher.buff = append(batcher.buff, &batchedMessage{data: "hello", iterationsLeft: len(batcher.buff)})
	assert.Equal(t,batcher.Size(),1)
}

func TestGRPCBatcher_FuntionalTestWithDeletingTrue(t *testing.T) {

	done := make(chan int)
	mockHanlder := &MockHandler{counter:0,done:done}
	seconds := 2 * time.Second
	batcher := NewGRPCMessageBatcher(seconds, mockHanlder,true)

	batcher.Add("hello1")
	batcher.Add("hello2")
	batcher.Add("hello3")
	batcher.Add("hello4")

	assert.Equal(t,batcher.Size(),4)

	<-done
	batcher.Stop()

	assert.Equal(t,batcher.Size(),0)
}

func TestGRPCBatcher_FuntionalTest_with_deleting_false(t *testing.T) {

	done := make(chan int)
	mockHanlder := &MockHandler{counter:0,done:done}
	seconds := 2 * time.Second
	batcher := NewGRPCMessageBatcher(seconds, mockHanlder,false)

	batcher.Add("hello1")
	batcher.Add("hello2")
	batcher.Add("hello3")
	batcher.Add("hello4")

	assert.Equal(t,batcher.Size(),4)

	<-done
	batcher.Stop()

	assert.Equal(t,batcher.Size(),4)
}