package blockchain

type BlockService interface{
	// Confirmed 된 블록 추가
	AddBlock(blk *Block) (bool, error)

	// Block Chain의 마지막 블록을 반환
	GetLastBlock() (*Block, error)

	// 블록을 검증
	VerifyBlock(blk *Block) (bool, error)

	// 블록 조회
	LookUpBlock(arg interface{}) (*Block, error)

	// 블록 생성
	CreateBlock(txList []*Transaction, createPeerId string) (*Block, error)
}
