package peer

import (
	"testing"
	"fmt"
)

type MockHandler struct {

}

func (mh *MockHandler) handle(ms interface{}) (interface{},error){

	fmt.Print(ms)

	return "success", nil
}


func TestPeerServiceImpl_GetPeerInfoByPeerID(t *testing.T) {
	//var peerServiceImpl := PeerServiceImpl{}
}
