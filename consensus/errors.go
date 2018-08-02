package consensus

import "errors"

var CreateConsensusError = errors.New("Create Consensus Error")
var InvalidLeaderIdError = errors.New("invalid Leader Id")
var SavePrepareMsgError = errors.New("Save PrepareMsg Error")
var SaveCommitMsgError = errors.New("Save CommitMsg Error")