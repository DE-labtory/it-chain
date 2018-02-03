package domain

import (
	"time"
	"errors"
	"it-chain/common"
	"strconv"
)

type Status int

const (
	Status_BLOCK_UNCONFIRMED Status = 0
	Status_BLOCK_CONFIRMED   Status = 1
)

type Block struct {
	Header          	*BlockHeader
	MerkleTree      	[][]string
	MerkleTreeHeight	int
	TransactionCount	int
	Transactions	 	[]*Transaction
}

func (s *Block) Reset() { *s = Block{} }

type BlockHeader struct {
	Number             uint64
	PreviousHash       string
	Version            string
	MerkleTreeRootHash string
	TimeStamp          time.Time
	BlockHeight        int
	BlockStatus        Status
	CreatedPeerID      string
	Signature          []byte
	PeerId             string
	BlockHash          string
}

func CreateNewBlock(prevBlock *Block, createPeerId string) *Block{
	var header BlockHeader
	if prevBlock == nil{
		header.Number = 0
		header.PreviousHash = ""
		header.Version = ""
		header.BlockHeight = 0
	} else {
		header.Number = prevBlock.Header.Number + 1
		header.PreviousHash = prevBlock.Header.BlockHash
		header.Version = prevBlock.Header.Version
		header.BlockHeight = prevBlock.Header.BlockHeight
	}
	header.CreatedPeerID = createPeerId
	header.TimeStamp = time.Now().Round(0)
	header.BlockStatus = Status_BLOCK_UNCONFIRMED

	return &Block{Header:&header, MerkleTree:make([][]string, 0), MerkleTreeHeight:0, TransactionCount:0, Transactions:make([]*Transaction, 0)}
}

func (s *BlockHeader) Reset() { *s = BlockHeader{} }

func (s *Block) PutTranscation(tx *Transaction) (valid bool, err error){
	if tx.Validate() == false{
		return false, errors.New("invalid tx")
	}
	if tx.TransactionStatus == status_TRANSACTION_UNKNOWN {
		if true { // Docker에서 실행하고 return이 true면 Confirmed 나중에 수정할 것.
			tx.TransactionStatus = status_TRANSACTION_CONFIRMED
		} else {
			tx.TransactionStatus = status_TRANSACTION_UNCONFIRMED
		}
	}
	for _, confirmedTx := range s.Transactions{
		if confirmedTx.TransactionHash == tx.TransactionHash{
			return false, errors.New("tx already exists")
		}
	}
	s.Transactions = append(s.Transactions, tx)
	s.TransactionCount++
	return true, nil
}

func (s Block) FindTransactionIndex(hash string) (idx int, err error){
	for idx = 0; idx < s.TransactionCount; idx++{
		if hash == s.Transactions[idx].TransactionHash{
			return idx, nil
		}
	}
	return -1, errors.New("txHash is not here")
}

func (s *Block) MakeMerkleTree(){
	var mtList []string
	for _, h := range s.Transactions{
		mtList = append(mtList, h.TransactionHash)
	}
	for {
		treeLength := len(mtList)
		s.MerkleTreeHeight++
		if treeLength <= 1 {
			s.MerkleTree = append(s.MerkleTree, mtList)
			break
		} else if treeLength % 2 == 1 {
			mtList = append(mtList, mtList[treeLength - 1])
			treeLength++
		}
		s.MerkleTree = append(s.MerkleTree, mtList)
		var tmpMtList []string
		for x := 0; x < treeLength/2; x++{
			idx := x * 2
			hashArg := []string{mtList[idx], mtList[idx+1]}
			mkHash := common.ComputeSHA256(hashArg)
			tmpMtList = append(tmpMtList, mkHash)
		}
		mtList = tmpMtList
	}
	if len(mtList) == 1 {
		s.Header.MerkleTreeRootHash = mtList[0]
	}
}

func (s Block) MakeMerklePath(idx int) (path []string){
	for i := 0; i < s.MerkleTreeHeight-1; i++{
		path = append(path, s.MerkleTree[i][(idx >> uint(i)) ^ 1])
	}
	return path
}

func (s *Block) GenerateBlockHash() error{
	if s.Header.MerkleTreeRootHash == "" {
		return errors.New("no merkle tree root hash")
	}
	str := []string{s.Header.MerkleTreeRootHash, s.Header.TimeStamp.String(), s.Header.PreviousHash}
	s.Header.BlockHash = common.ComputeSHA256(str)
	return nil
}

func (s Block) BlockSerialize() ([]byte, error){
	return common.Serialize(s)
}

func BlockDeserialize(by []byte) (Block, error) {
	block := Block{}
	err := common.Deserialize(by, &block)
	return block, err
}

// 해당 트랜잭션이 정당한지 머클패스로 검사함
func (s Block) VarifyTx(tx Transaction) (bool, error) {
	hash := tx.TransactionHash
	idx, err := s.FindTransactionIndex(hash)

	if err != nil {
		return false, err
	}

	merklePath := s.MakeMerklePath(idx)

	for _, sibling_hash := range merklePath{
		str := []string{hash, sibling_hash}
		hash = common.ComputeSHA256(str)
	}

	if hash == s.Header.MerkleTreeRootHash{
		return true, nil
	} else {
		return false, errors.New("tx is invalid")
	}
}

// 블럭내의 모든 트랜잭션들이 정당한지 머클패스로 검사함
func (s Block) VarifyBlock() (bool, error) {
	for idx := 0; idx < s.TransactionCount; idx++{
		txVarification, txErr := s.VarifyTx(*s.Transactions[idx])
		if txVarification == false  {
			err := errors.New("block is invalid --- " + strconv.Itoa(idx) + "'s " + txErr.Error())
			return false, err
		}
	}
	return true, nil
}


