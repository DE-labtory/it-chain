package api

import (
	"errors"

	"github.com/it-chain/it-chain-Engine/consensus/domain/factory"
	cs "github.com/it-chain/it-chain-Engine/consensus/domain/model/consensus"
	"github.com/it-chain/it-chain-Engine/consensus/domain/model/msg"
	"github.com/it-chain/it-chain-Engine/consensus/domain/model/parliament"
	"github.com/it-chain/it-chain-Engine/consensus/domain/repository"
	"github.com/it-chain/it-chain-Engine/consensus/domain/service"
)

type ConsensusApi struct {
	consensusRepository repository.ConsensusRepository
	parlimentRepository repository.ParlimentRepository
	msgPool             msg.MsgPool
	messageApi          MessageApi
}

func (cApi ConsensusApi) StartConsensus(userId parliament.PeerID, block cs.Block) error {

	//Paliament의 Validate Check
	parliament := cApi.parlimentRepository.Get()

	if parliament == nil {
		return errors.New("No parliament")
	}

	if parliament.IsNeedConsensus() {
		consensus, err := factory.CreateConsensus(*parliament, block)

		if err != nil {
			return err
		}

		consensus.Start()
		cApi.consensusRepository.Save(*consensus)

		PreprepareMessage := factory.CreatePreprepareMsg(*consensus)
		cApi.messageApi.BroadCastMsg(PreprepareMessage, consensus.Representatives)
	} else {
		cApi.messageApi.ConfirmedBlock(block)
	}

	return nil
}

func (cApi ConsensusApi) ReceivePrepareMsg(msg msg.PrepareMsg) {

	cApi.msgPool.InsertPrepareMsg(msg)
	consensus := cApi.consensusRepository.FindById(msg.ConsensusID)

	if service.CheckPreparePolicy(*consensus, cApi.msgPool) {
		CommitMsg := factory.CreateCommitMsg(*consensus)
		cApi.messageApi.BroadCastMsg(CommitMsg, consensus.Representatives)
	} else {
		return
	}
}

func (cApi ConsensusApi) ReceiveCommitMsg(msg msg.CommitMsg) {
	cApi.msgPool.InsertCommitMsg(msg)
	consensus := cApi.consensusRepository.FindById(msg.ConsensusID)

	if service.CheckCommitPolicy(*consensus, cApi.msgPool) {
		cApi.messageApi.ConfirmedBlock(consensus.Block)
	} else {
		return
	}
}

func (cApi ConsensusApi) ReceivePreprepareMsg(msg msg.PreprepareMsg) {

	consensus := msg.Consensus
	parliament := cApi.parlimentRepository.Get()

	flag := parliament.ValidateRepresentative(consensus.Representatives)

	if !flag {
		return
	}

	consensus.Start()
	cApi.consensusRepository.Save(consensus)
	PrepareMsg := factory.CreatePrepareMsg(consensus)
	cApi.messageApi.BroadCastMsg(PrepareMsg, consensus.Representatives)
}
