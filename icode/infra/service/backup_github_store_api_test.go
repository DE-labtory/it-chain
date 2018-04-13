package service

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateRepository(t *testing.T) {

	repoName := "test-chain"
	b, err := NewBackupGithubStoreApi("", "")
	assert.NoError(t, err)

	defer func() {
		ctx := context.Background()
		_, err := client.Repositories.Delete(ctx, "steve-buzzni", repoName)
		assert.NoError(t, err)
	}()

	r, err := b.CreateRepository(repoName)

	assert.NoError(t, err)
	assert.Equal(t, "https://github.com/steve-buzzni/test-chain.git", r.GetCloneURL())
}

func TestNewBackupGithubStoreApi(t *testing.T) {
	b, err := NewBackupGithubStoreApi("", "")
	assert.NoError(t, err)

	assert.Equal(t, "https://github.com/steve-buzzni", b.homepageUrl)
	assert.Equal(t, "steve-buzzni", b.storename)
}
