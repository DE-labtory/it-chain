package adapter

import (
	"github.com/it-chain/engine/consensus"
	"github.com/it-chain/engine/consensus/api"
)

type MemberJoinedEventHandler struct {
	parliamentApi api.ParliamentApi
}

func (handler MemberJoinedEventHandler) HandleMemberJoinedEvent(event consensus.MemberJoinedEvent) error {

	mid := consensus.MemberId{event.MemberId}
	err := handler.parliamentApi.AddMember(consensus.Member{mid})

	if err != nil {
		return err
	}
	return nil
}
