/*
 * Copyright 2018 It-chain
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * https://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package blockchain

import (
	"bytes"
	"crypto/sha256"
	"errors"
	"time"

	"github.com/it-chain/yggdrasill/common"
)

// ErrHashCalculationFailed 변수는 Hash 계산 중 발생한 에러를 정의한다.
var ErrHashCalculationFailed = errors.New("Hash Calculation Failed Error")
var ErrInsufficientFields = errors.New("Previous seal or transaction list seal is not set")
var ErrEmptyTxList = errors.New("Empty TxList")

type Validator = common.Validator

// DefaultValidator 객체는 Validator interface를 구현한 객체.
type DefaultValidator struct{}

func (t *DefaultValidator) ValidateTree(tree Tree) (bool, error) {

	root := tree.GetRoot()
	if root.IsLeafNode() {
		return validateGenesis(root)
	}

	calculatedMerkleRoot, err := validateLeaf(tree.GetRoot())
	if err != nil {
		return false, err
	}

	if bytes.Compare(tree.GetTxSealRoot(), calculatedMerkleRoot) == 0 {
		return true, nil
	}
	return false, nil
}

func validateGenesis(node Node) (bool, error) {
	return bytes.Equal(node.GetTxSeal(), []byte("genesis")), nil
}

func validateLeaf(node Node) ([]byte, error) {

	if node.IsLeafNode() {
		return node.GetTransaction().CalculateSeal()
	}

	rightHash, err := validateLeaf(node.GetRight())
	if err != nil {
		return nil, err
	}

	leftHash, err := validateLeaf(node.GetLeft())
	if err != nil {
		return nil, err
	}

	return calculateIntermediateLeafHash(leftHash, rightHash), nil

}

// ValidateSeal 함수는 원래 Seal 값과 주어진 Seal 값(comparisonSeal)을 비교하여, 올바른지 검증한다.
func (t *DefaultValidator) ValidateSeal(seal []byte, comparisonBlock Block) (bool, error) {

	txSealRoot := comparisonBlock.GetTxSealRoot()

	comparisonSeal, err := t.BuildSeal(comparisonBlock.GetTimestamp(), comparisonBlock.GetPrevSeal(), txSealRoot, comparisonBlock.GetCreator())
	if err != nil {
		return false, err
	}

	return bytes.Compare(seal, comparisonSeal) == 0, nil
}

//TODO: Tree 구조가 안정되면 지우기
// ValidateTxSeal 함수는 주어진 Transaction 리스트에 따라 주어진 transaction Seal을 검증함.
func (t *DefaultValidator) ValidateTxSeal(txSeal [][]byte, txList []Transaction) (bool, error) {
	leafNodeIndex := 0
	if len(txList)%2 != 0 {
		txList = append(txList, txList[len(txList)-1])
	}
	for i, n := range txSeal {
		leftIndex, rightIndex := (i+1)*2-1, (i+1)*2
		if rightIndex >= len(txSeal) {
			// Check Leaf Node

			calculatedHash, error := txList[leafNodeIndex].CalculateSeal()
			if error != nil {
				return false, ErrHashCalculationFailed
			}

			if bytes.Compare(n, calculatedHash) != 0 {
				return false, nil
			}
			leafNodeIndex++
		} else {
			// Check Intermediate Node
			leftNode, rightNode := txSeal[leftIndex], txSeal[rightIndex]
			calculatedHash := calculateIntermediateLeafHash(leftNode, rightNode)
			if bytes.Compare(n, calculatedHash) != 0 {
				return false, nil
			}
		}
	}

	return true, nil
}

// ValidateTransaction 함수는 주어진 Transaction이 이 txSeal에 올바로 있는지를 확인한다.
func (t *DefaultValidator) ValidateTransaction(txSeal [][]byte, transaction Transaction) (bool, error) {
	hash, error := transaction.CalculateSeal()
	if error != nil {
		return false, error
	}

	index := -1
	for i, h := range txSeal {
		if bytes.Compare(h, hash) == 0 {
			index = i
		}
	}

	if index == -1 {
		return false, nil
	}

	var siblingIndex, parentIndex int
	for index > 0 {
		var isLeft bool
		if index%2 == 0 {
			siblingIndex = index - 1
			parentIndex = (index - 1) / 2
			isLeft = false
		} else {
			siblingIndex = index + 1
			parentIndex = index / 2
			isLeft = true
		}

		var parentHash []byte
		if isLeft {
			parentHash = calculateIntermediateLeafHash(txSeal[index], txSeal[siblingIndex])
		} else {
			parentHash = calculateIntermediateLeafHash(txSeal[siblingIndex], txSeal[index])
		}

		if bytes.Compare(parentHash, txSeal[parentIndex]) != 0 {
			return false, nil
		}

		index = parentIndex
	}

	return true, nil
}

// BuildSeal 함수는 block 객체를 받아서 Seal 값을 만들고, Seal 값을 반환한다.
// 인풋 파라미터의 block에 자동으로 할당해주지는 않는다.
func (t *DefaultValidator) BuildSeal(timeStamp time.Time, prevSeal []byte, txSealRoot []byte, creator string) ([]byte, error) {
	timestamp, err := timeStamp.MarshalText()
	if err != nil {
		return nil, err
	}

	//ToDo: tree == nil 수정.
	if prevSeal == nil || txSealRoot == nil || creator == "" {
		return nil, ErrInsufficientFields
	}

	//ToDo: txSealRoot가 없을 때 에러 처리.
	combined := append(prevSeal, txSealRoot...)
	combined = append(combined, timestamp...)

	seal := calculateHash(combined)
	return seal, nil
}

func (t *DefaultValidator) BuildTree(txList []Transaction) (Tree, error) {

	root, primeLeafs, err := buildMerkleTree(GetBackTxType(txList))
	if err != nil {
		return nil, err
	}

	return &DefaultTree{
		Root:       root,
		TxSealRoot: root.TxSeal,
		PrimeLeafs: primeLeafs,
	}, nil
}

//ToDo: legacy code 중 buildTree를 지우면, buildTree로 네이밍 변환
func buildMerkleTree(txList []*DefaultTransaction) (*DefaultNode, []*DefaultNode, error) {

	if isEmpty(txList) {
		return &DefaultNode{
			IsLeaf: true,
			TxSeal: []byte("genesis"),
		}, nil, nil
	}

	Leafs := []*DefaultNode{}
	for _, tx := range txList {
		txSeal, err := tx.CalculateSeal()
		if err != nil {
			return nil, nil, err
		}

		Leafs = append(Leafs, &DefaultNode{
			TxSeal:      txSeal,
			Transaction: tx,
			IsLeaf:      true,
		})
	}

	if len(Leafs)%2 == 1 {
		duplicatedNode := &DefaultNode{
			TxSeal:       Leafs[len(Leafs)-1].TxSeal,
			Transaction:  Leafs[len(Leafs)-1].Transaction,
			IsLeaf:       true,
			IsDuplicated: true,
		}
		Leafs = append(Leafs, duplicatedNode)
	}

	root, err := buildRoot(Leafs)
	if err != nil {
		return nil, nil, err
	}

	return root, Leafs, nil

}

func isEmpty(txList []*DefaultTransaction) bool {
	if len(txList) == 0 {
		return true
	}
	return false
}

func buildRoot(nodeList []*DefaultNode) (*DefaultNode, error) {
	return buildIntermediate(nodeList)
}

func buildIntermediate(nodeList []*DefaultNode) (*DefaultNode, error) {

	intermediateLeafList := []*DefaultNode{}
	for i := 0; i < len(nodeList); i += 2 {
		h := sha256.New()
		left, right := i, i+1

		if i+1 == len(nodeList) {
			right = i
		}

		joinedSeal := append(nodeList[left].TxSeal, nodeList[right].TxSeal...)
		if _, err := h.Write(joinedSeal); err != nil {
			return nil, err
		}

		newLeaf := &DefaultNode{
			Left:   nodeList[left],
			Right:  nodeList[right],
			TxSeal: h.Sum(nil),
		}
		intermediateLeafList = append(intermediateLeafList, newLeaf)

		if len(nodeList) == 2 {
			root := newLeaf
			return root, nil
		}
	}
	return buildIntermediate(intermediateLeafList)
}

//ToDo: Tree 구조가 안정되면 지우기
// BuildTxSeal 함수는 Transaction 배열을 받아서 TxSeal을 생성하여 반환한다.
func (t *DefaultValidator) BuildTxSeal(txList []Transaction) ([][]byte, error) {
	if len(txList) == 0 {
		return nil, ErrEmptyTxList
	}

	leafNodeList := make([][]byte, 0)

	for _, tx := range txList {
		leafNode, err := tx.CalculateSeal()
		if err != nil {
			return nil, err
		}

		leafNodeList = append(leafNodeList, leafNode)
	}

	// leafNodeList의 개수는 짝수개로 맞춤. (홀수 일 경우 마지막 Tx를 중복 저장.)
	if len(leafNodeList)%2 != 0 {
		leafNodeList = append(leafNodeList, leafNodeList[len(leafNodeList)-1])
	}

	tree, err := buildTree(leafNodeList, leafNodeList)
	if err != nil {
		return nil, err
	}

	// DefaultValidator 는 Merkle Tree의 루트노드(tree[0])를 Proof로 간주함
	return tree, nil
}

func buildTree(nodeList [][]byte, fullNodeList [][]byte) ([][]byte, error) {
	intermediateNodeList := make([][]byte, 0)
	for i := 0; i < len(nodeList); i += 2 {
		leftIndex, rightIndex := i, i+1
		leftNode, rightNode := nodeList[leftIndex], nodeList[rightIndex]

		intermediateNode := calculateIntermediateLeafHash(leftNode, rightNode)

		intermediateNodeList = append(intermediateNodeList, intermediateNode)

		if len(nodeList) == 2 {
			return append(intermediateNodeList, fullNodeList...), nil
		}
	}

	newFullNodeList := append(intermediateNodeList, fullNodeList...)

	return buildTree(intermediateNodeList, newFullNodeList)
}

func calculateIntermediateLeafHash(leftHash []byte, rightHash []byte) []byte {
	combinedHash := append(leftHash, rightHash...)

	return calculateHash(combinedHash)
}
