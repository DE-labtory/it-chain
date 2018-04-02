package repository

import "github.com/it-chain/it-chain-Engine/blockchain/domain/model/block"

type BlockRepository interface {
	Close()
	AddBlock(block block.Block)
	GetLastBlock(block block.Block)
	GetTransactionsById(id int) // block과 관련된 정보 조회 예시
}
