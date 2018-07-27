package consensus

type PropagateService interface {
	BroadcastPrePrepareMsg(msg PrePrepareMsg) error
	BroadcastPrepareMsg(msg PrepareMsg) error
	BroadcastCommitMsg(msg CommitMsg) error
}

type ConfirmService interface {
	ConfirmBlock(block ProposedBlock) error
}

// 연결된 peer 중에서 consensus 에 참여할 representative 들을 선출
func Elect(parliament []MemberId) ([]*Representative, error) {
	representatives := make([]*Representative, 0)

	for _, peerId := range parliament {
		representatives = append(representatives, NewRepresentative(peerId.ToString()))
	}

	return representatives, nil
}
