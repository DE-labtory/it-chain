package smartcontract

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"fmt"
)

func TestDeploy_Deploy(t *testing.T) {
	scs := SmartContractService{"hackurity01", map[string]SmartContract{}}
	ContractPath := "junbeomlee/bloom"

	deploy_result, err := scs.Deploy(ContractPath)

	fmt.Println(deploy_result)
	if err != nil {
		assert.Fail(t, err.Error())
	}

	fmt.Println(deploy_result)
	//assert.Equal(t,nil,deploy_result)
}