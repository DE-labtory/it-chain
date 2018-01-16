package comm

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestNewConnectionWithAddress(t *testing.T) {

	conn,err := NewConnectionWithAddress("127.0.0.1:8080",false,nil)
	defer conn.Close()

	if err != nil{
		assert.Fail(t,"fail to connect")
	}

	assert.NotNil(t,conn)
}