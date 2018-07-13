package consensus

func CheckConsensusPolicy(consensus Consensus) bool {
	state := consensus.CurrentState
	numberOfRepresentatives := float64(len(consensus.Representatives))
	prepareMsgs := consensus.PrepareMsgPool.messages
	commitMsgs := consensus.CommitMsgPool.messages
	var numberOfMsgs int

	switch state {

	case PREPARE_STATE:
		numberOfMsgs = len(prepareMsgs)

		if numberOfMsgs > int((numberOfRepresentatives)/3)+1 {
			return true
		}
		return false

	case COMMIT_STATE:
		numberOfMsgs = len(commitMsgs)

		if numberOfMsgs > int((numberOfRepresentatives)/3)+1 {
			return true
		}
		return false

	default:
		return false
	}
}
