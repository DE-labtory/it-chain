package consensus

import (
	"github.com/it-chain/midgard"
	"errors"
	"fmt"
)

type MemberId struct {
	Id string
}

type Member struct {
	MemberId MemberId
}

func (mid MemberId) ToString() string {
	return string(mid.Id)
}

func (m *Member) StringMemberId() string {
	return m.MemberId.ToString()
}

func (m Member) GetId() string {
	return m.StringMemberId()
}

func (m *Member) On(event midgard.Event) error {
	switch v := event.(type) {

	case *MemberJoinedEvent:
		m.MemberId = MemberId{v.ID}

	default:
		return errors.New(fmt.Sprintf("unhandled event [%s]", v))
	}

	return nil
}