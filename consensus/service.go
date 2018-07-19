package consensus

import "errors"

type PropagateService interface {
	BroadcastPrePrepareMsg(msg PrePrepareMsg) error
	BroadcastPrepareMsg(msg PrepareMsg) error
	BroadcastCommitMsg(msg CommitMsg) error
}

type ConfirmService interface {
	ConfirmBlock(block ProposedBlock) error
}

func Elect(parliament Parliament) ([]*Representative, error) {
	representatives := make([]*Representative, 0)

	if !parliament.HasLeader() {
		return nil, errors.New("No Leader")
	}

	representatives = append(representatives, NewRepresentative(parliament.Leader.GetID()))

	for _, member := range parliament.Members {
		representatives = append(representatives, NewRepresentative(member.GetID()))
	}

	return representatives, nil
}
