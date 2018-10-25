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

import "github.com/it-chain/yggdrasill/common"

type Tree = common.Tree

type DefaultTree struct {
	Root       *DefaultNode
	TxSealRoot []byte
	PrimeLeafs []*DefaultNode
}

func (t *DefaultTree) SetTxSealRoot(txSealRoot []byte) {
	t.TxSealRoot = txSealRoot
}

func (t *DefaultTree) GetRoot() Node {
	return t.Root
}

func (t *DefaultTree) GetTxSealRoot() []byte {
	return t.TxSealRoot
}

type Node = common.Node

type DefaultNode struct {
	Left         *DefaultNode
	Right        *DefaultNode
	IsLeaf       bool
	IsDuplicated bool
	TxSeal       []byte
	Transaction  *DefaultTransaction
}

func (n *DefaultNode) IsLeafNode() bool {
	return n.IsLeaf
}

func (n *DefaultNode) GetTransaction() Transaction {
	return n.Transaction
}

func (n *DefaultNode) GetRight() Node {
	return n.Right
}

func (n *DefaultNode) GetLeft() Node {
	return n.Left
}

func ConvTreeType(tree Tree) *DefaultTree {
	convTree := tree.(*DefaultTree)
	return convTree
}
