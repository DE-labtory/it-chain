package bus_test

import (
	"fmt"
	"testing"

	"encoding/json"

	"github.com/it-chain/midgard"
	"github.com/it-chain/midgard/bus"
	"github.com/stretchr/testify/assert"
)

func TestNewParamBasedRouter(t *testing.T) {
	d, err := bus.NewParamBasedRouter()
	assert.NoError(t, err)

	err = d.SetHandler(&Dispatcher{})
	assert.NoError(t, err)

	cmd := UserNameUpdateCommand{
		Name: "jun",
	}

	b, err := json.Marshal(cmd)
	assert.NoError(t, err)

	fmt.Println(b)

	err = d.Route(b, "UserNameUpdateCommand")
	assert.NoError(t, err)
}

type UserNameUpdateCommand struct {
	midgard.EventModel
	Name string
}

type Dispatcher struct {
}

func (d *Dispatcher) HandleNameUpdateCommand(command UserNameUpdateCommand) {
	fmt.Println("hello world2")
}
