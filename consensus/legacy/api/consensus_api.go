package api

import (
	"github.com/it-chain/it-chain-Engine/consensus/legacy/domain/factory"
	cs "github.com/it-chain/it-chain-Engine/consensus/legacy/domain/model/consensus"
	"github.com/it-chain/it-chain-Engine/consensus/legacy/domain/model/msg"
	"github.com/it-chain/it-chain-Engine/consensus/legacy/domain/model/parliament"
	"github.com/it-chain/it-chain-Engine/consensus/legacy/domain/repository"
	"github.com/it-chain/it-chain-Engine/consensus/legacy/domain/service"
)

type ConsensusApi struct {
	consensusRepository repository.ConsensusRepository
	parlimentRepository repository.ParlimentRepository
	msgPool             msg.MsgPool
	messageService      service.MessageService
}

func (cApi ConsensusApi) StartConsensus(userId parliament.PeerID, block cs.Block) error {

	//todo start consensus timeout batcher
	//Paliament의 Validate Check
	parliament := cApi.parlimentRepository.Get()

	if parliament.IsNeedConsensus() {
		consensus, err := factory.CreateConsensus(parliament, block)

		if err != nil {
			return err
		}

		consensus.Start()
		cApi.consensusRepository.Save(*consensus)

		PreprepareMessage := factory.CreatePreprepareMsg(*consensus)
		cApi.messageService.BroadCastMsg(PreprepareMessage, consensus.Representatives)
	} else {
		cApi.messageService.ConfirmedBlock(block)
	}

	return nil
}

func (cApi ConsensusApi) ReceivePrepareMsg(msg msg.PrepareMsg) {

	cApi.msgPool.InsertPrepareMsg(msg)
	consensus := cApi.consensusRepository.FindById(msg.ConsensusID)

	if service.CheckPreparePolicy(*consensus, cApi.msgPool) {
		CommitMsg := factory.CreateCommitMsg(*consensus)
		consensus.ToCommitState()
		cApi.messageService.BroadCastMsg(CommitMsg, consensus.Representatives)
	} else {
		return
	}
}

func (cApi ConsensusApi) ReceiveCommitMsg(msg msg.CommitMsg) {
	cApi.msgPool.InsertCommitMsg(msg)
	consensus := cApi.consensusRepository.FindById(msg.ConsensusID)

	if service.CheckCommitPolicy(*consensus, cApi.msgPool) {
		cApi.messageService.ConfirmedBlock(consensus.Block)
		//todo delete consensus and remove all message
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
	cApi.messageService.BroadCastMsg(PrepareMsg, consensus.Representatives)
}
