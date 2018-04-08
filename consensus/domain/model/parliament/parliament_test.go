package parliament

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParliament_NewParliament(t *testing.T) {
	p := NewParliament()

	assert.Nil(t, nil, p.Leader)
	assert.Equal(t, 0, len(p.Members))
}
