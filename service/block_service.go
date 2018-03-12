package service

import "github.com/it-chain/it-chain-Engine/domain"

type BlockService interface{
	// Confirmed 된 블록 추가
	AddBlock(blk *domain.Block) (bool, error)

	// DB close
	Close()

	// Block Chain의 마지막 블록을 반환
	GetLastBlock() (*domain.Block, error)

	// 블록을 검증
	VerifyBlock(blk *domain.Block) (bool, error)

	// 블록 조회
	LookUpBlock(arg interface{}) (*domain.Block, error)

	// 블록 생성
	CreateBlock(txList []*domain.Transaction, createPeerId string) (*domain.Block, error)
}
