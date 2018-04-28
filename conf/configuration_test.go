package conf

import (
	"fmt"
	"testing"
)

//todo : need testing code
func TestGetConfiguration(t *testing.T) {
	conf := GetConfiguration()
	fmt.Println(conf.Common.BootNodeIp)
}
