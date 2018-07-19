package conf

import (
	"testing"

	"github.com/magiconair/properties/assert"
)

//todo : need testing code
func TestGetConfiguration(t *testing.T) {
	conf := GetConfiguration()
	assert.Equal(t, conf.GrpcGateway.Address, "127.0.0.1")
	assert.Equal(t, conf.Engine.Mode, "solo")
}
