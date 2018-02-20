package domain

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestGetRepositoryList(t *testing.T) {

	list, err := GetRepositoryList("yojkim")
	assert.NoError(t, err)
	assert.NotNil(t, list)

	t.Log(list)

	_, err = GetRepositoryList("")
	assert.Error(t, err)

}
