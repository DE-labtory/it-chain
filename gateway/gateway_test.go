package gateway

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStart(t *testing.T) {
	err := Start()

	assert.NoError(t, err)
}
