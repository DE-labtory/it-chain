package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateRepository(t *testing.T) {
	repoName := "test-chain"
	_, _, err := CreateRepository(repoName)

	assert.NoError(t, err)
	//assert.Equal(t, "https://github.com/"+, r.GetCloneURL())
}
