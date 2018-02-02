package smartcontract

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestGetRepositoryList(t *testing.T) {

	err := GetRepositoryList("yojkim")
	assert.NoError(t, err)

}
