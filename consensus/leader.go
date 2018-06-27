package consensus

import (
	"github.com/it-chain/midgard"
	"errors"
	"fmt"
)

type LeaderId struct {
	Id string
}

type Leader struct {
	LeaderId LeaderId
}

func (lid LeaderId) ToString() string {
	return string(lid.Id)
}

func (l *Leader) StringLeaderId() string {
	return l.LeaderId.ToString()
}

func (l Leader) GetId() string {
	return l.StringLeaderId()
}

func (l *Leader) On(event midgard.Event) error {
	switch v := event.(type) {

	case *LeaderChangedEvent:
		l.LeaderId = LeaderId{v.GetID()}

	default:
		return errors.New(fmt.Sprintf("unhandled event [%s]", v))
	}

	return nil
}
