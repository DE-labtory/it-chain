package rabbitmq_test

import (
	"fmt"
	"log"
	"testing"

	"sync"

	"time"

	"github.com/it-chain/midgard"
	"github.com/it-chain/midgard/bus/rabbitmq"
	"github.com/stretchr/testify/assert"
)

var wg sync.WaitGroup

func TestConnect(t *testing.T) {

	wg.Add(2)
	c := rabbitmq.Connect("")
	err := c.Subscribe("asd", "asd", &Dispatcher{})
	assert.NoError(t, err)

	err = c.Publish("asd", "asd", UserNameUpdateEvent{
		Name: "JUN",
		EventModel: midgard.EventModel{
			ID:      "123",
			Time:    time.Now(),
			Type:    "123",
			Version: 1,
		}})
	assert.NoError(t, err)

	err = c.Publish("asd", "asd", UserAddCommand{
		CommandModel: midgard.CommandModel{
			ID: "123",
		}})

	assert.NoError(t, err)

	wg.Wait()
}

type UserNameUpdateEvent struct {
	midgard.EventModel
	Name string
}

type UserAddCommand struct {
	midgard.CommandModel
}

type Dispatcher struct {
}

func (d *Dispatcher) Handle(command UserAddCommand) {
	log.Print("hello world")
	fmt.Println(command)
	wg.Done()
}

func (d *Dispatcher) HandleNameUpdateCommand(event UserNameUpdateEvent) {
	fmt.Println("hello world2")
	fmt.Println(event)
	wg.Done()
}
