package eventstore_test

import (
	"errors"
	"fmt"

	"github.com/it-chain/midgard"
)

// aggregate
type User struct {
	Name string
	midgard.AggregateModel
}

func (u *User) On(event midgard.Event) error {

	switch v := event.(type) {

	case *UserCreatedEvent:
		u.ID = v.ID

	case *UserNameUpdatedEvent:
		u.Name = v.Name

	default:
		return errors.New(fmt.Sprintf("unhandled event [%s]", v))
	}

	return nil
}

// Command
type UserCreateCommand struct {
	midgard.CommandModel
}

type UserNameUpdateCommand struct {
	midgard.CommandModel
	Name string
}

// Event
type UserCreatedEvent struct {
	midgard.EventModel
}

type UserNameUpdatedEvent struct {
	midgard.EventModel
	Name string
}
