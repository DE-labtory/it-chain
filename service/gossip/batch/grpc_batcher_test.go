package batch

import (
	"testing"
	"fmt"
	"time"
	"github.com/magiconair/properties/assert"
)

type MockHandler struct{}

func (mh *MockHandler) handle([]byte) (interface{},error){

	fmt.Print("success")
	return "success", nil
}

func TestNewGRPCMessageBatcher(t *testing.T) {

	mockHanlder := &MockHandler{}
	seconds := time.Duration(10)
	batcher := NewGRPCMessageBatcher(seconds, mockHanlder,false)

	assert.Equal(t,batcher.stopFlag,int32(0))
	assert.Equal(t,batcher.deleting,false)
	assert.Equal(t,len(batcher.buff),0)
}