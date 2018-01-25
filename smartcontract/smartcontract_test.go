package smartcontract

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestDeploy_Deploy(t *testing.T) {

	smartcontract := SmartContract{map[string]string{}}
	ContractPath := "junbeomlee/bloom"

	deploy_result := smartcontract.Deploy(ContractPath)
	assert.Equal(t,nil,deploy_result)
}