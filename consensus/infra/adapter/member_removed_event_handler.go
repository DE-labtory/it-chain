package adapter

import (
	"github.com/it-chain/engine/consensus"
	"github.com/it-chain/engine/consensus/api"
)

type MemberRemovedEventHandler struct {
	parliamentApi api.ParliamentApi
}

func (handler MemberRemovedEventHandler) HandleMemberRemovedEvent(event consensus.MemberRemovedEvent) error {

	mid := consensus.MemberId{event.MemberId}

	err := handler.parliamentApi.RemoveMember(mid)

	if err != nil {
		return err
	}
	return nil
}
