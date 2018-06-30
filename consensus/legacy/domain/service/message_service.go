package service

import "github.com/it-chain/it-chain-Engine/consensus/legacy/domain/model/consensus"

type Serializable interface {
	ToByte() ([]byte, error)
}

type MessageService interface {
	ConfirmedBlock(block consensus.Block)
	BroadCastMsg(Msg Serializable, representatives []*consensus.Representative)
}
