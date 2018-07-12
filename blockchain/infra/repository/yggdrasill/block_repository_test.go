package yggdrasill
//
//import (
//	"testing"
//
//	"github.com/stretchr/testify/assert"
//)
//
//func TestNewBlockRepository(t *testing.T) {
//	dbPath := "./.db"
//	opts := map[string]interface{}{
//		"db_path": dbPath,
//	}
//	TestUser := "TestUser"
//
//	blockRepository, err := NewBlockRepository(dbPath, opts, TestUser)
//
//	assert.NoError(t, err)
//	assert.NotNil(t, blockRepository.GetValidator())
//	assert.Equal(t, blockRepository.Creator, TestUser)
//}
