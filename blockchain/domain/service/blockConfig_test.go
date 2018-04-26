package service

import (
	"testing"
	"github.com/stretchr/testify/assert"

)

func TestConfigFromJson(t *testing.T) {
	_,err1 := ConfigFromJson("GenesisBlockConfig.json")
	assert.NoError(t,err1)
	_,err2 := ConfigFromJson("WrongFileName.json")
	assert.Error(t,err2)
}