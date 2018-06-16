package yggdrasill

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewBlockRepository(t *testing.T) {
	dbPath := "./.db"
	opts := map[string]interface{}{
		"db_path": dbPath,
	}
	blockRepository, err := NewBlockRepository(dbPath, opts)
	assert.NoError(t, err)
	assert.NotNil(t, blockRepository.GetValidator())
}
