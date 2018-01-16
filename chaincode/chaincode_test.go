package chaincode

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestDeploy_Deploy(t *testing.T) {

	chaincode := Chaincode{map[string]string{}}
	ContractPath := "junbeomlee/bloom"

	deploy_result := chaincode.Deploy(ContractPath)
	assert.Equal(t,nil,deploy_result)
}