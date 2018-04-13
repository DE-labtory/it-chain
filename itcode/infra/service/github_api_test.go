package service

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateRepository(t *testing.T) {

	repoName := "test-chain"
	defer func() {
		ctx := context.Background()
		_, err := client.Repositories.Delete(ctx, "steve-buzzni", repoName)
		assert.NoError(t, err)
	}()

	r, _, err := CreateRepository(repoName)

	assert.NoError(t, err)
	assert.Equal(t, "https://github.com/steve-buzzni/test-chain.git", r.GetCloneURL())
}
