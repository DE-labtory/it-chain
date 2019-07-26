/*
 * Copyright 2018 DE-labtory
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
package mock

import "github.com/DE-labtory/it-chain/blockchain"

type BlockQueryService struct {
	GetStagedBlockByHeightFunc   func(height blockchain.BlockHeight) (blockchain.DefaultBlock, error)
	GetStagedBlockByIdFunc       func(blockId string) (blockchain.DefaultBlock, error)
	GetLastCommitedBlockFunc     func() (blockchain.DefaultBlock, error)
	GetCommitedBlockByHeightFunc func(height blockchain.BlockHeight) (blockchain.DefaultBlock, error)
}

func (s BlockQueryService) GetStagedBlockByHeight(height blockchain.BlockHeight) (blockchain.DefaultBlock, error) {
	return s.GetStagedBlockByHeightFunc(height)
}
func (s BlockQueryService) GetStagedBlockById(blockId string) (blockchain.DefaultBlock, error) {
	return s.GetStagedBlockByIdFunc(blockId)
}
func (s BlockQueryService) GetLastCommitedBlock() (blockchain.DefaultBlock, error) {
	return s.GetLastCommitedBlockFunc()
}
func (s BlockQueryService) GetCommitedBlockByHeight(height blockchain.BlockHeight) (blockchain.DefaultBlock, error) {
	return s.GetCommitedBlockByHeightFunc(height)
}

type BlockRepository struct {
	SaveFunc         func(block blockchain.DefaultBlock) error
	FindLastFunc     func() (blockchain.DefaultBlock, error)
	FindByHeightFunc func(height blockchain.BlockHeight) (blockchain.DefaultBlock, error)
	FindBySealFunc   func(seal []byte) (blockchain.DefaultBlock, error)
	FindAllFunc      func() ([]blockchain.DefaultBlock, error)
}

func (r BlockRepository) Save(block blockchain.DefaultBlock) error {
	return r.SaveFunc(block)
}

func (r BlockRepository) FindLast() (blockchain.DefaultBlock, error) {
	return r.FindLastFunc()
}

func (r BlockRepository) FindByHeight(height blockchain.BlockHeight) (blockchain.DefaultBlock, error) {
	return r.FindByHeightFunc(height)
}

func (r BlockRepository) FindBySeal(seal []byte) (blockchain.DefaultBlock, error) {
	return r.FindBySealFunc(seal)
}

func (r BlockRepository) FindAll() ([]blockchain.DefaultBlock, error) {
	return r.FindAllFunc()
}

type EventService struct {
	PublishFunc func(topic string, event interface{}) error
}

func (s EventService) Publish(topic string, event interface{}) error {
	return s.PublishFunc(topic, event)
}

type ConsensusService struct {
	ConsensusBlockFunc func(block blockchain.DefaultBlock) error
}

func (s ConsensusService) ConsentBlock(block blockchain.DefaultBlock) error {
	return s.ConsensusBlockFunc(block)
}

type QueryService struct {
	GetRandomPeerFunc            func() (blockchain.Peer, error)
	GetLastBlockFromPeerFunc     func(peer blockchain.Peer) (blockchain.DefaultBlock, error)
	GetBlockByHeightFromPeerFunc func(height blockchain.BlockHeight, peer blockchain.Peer) (blockchain.DefaultBlock, error)
}

func (s QueryService) GetRandomPeer() (blockchain.Peer, error) {
	return s.GetRandomPeerFunc()
}

func (s QueryService) GetLastBlockFromPeer(peer blockchain.Peer) (blockchain.DefaultBlock, error) {
	return s.GetLastBlockFromPeerFunc(peer)
}

func (s QueryService) GetBlockByHeightFromPeer(height blockchain.BlockHeight, peer blockchain.Peer) (blockchain.DefaultBlock, error) {
	return s.GetBlockByHeightFromPeerFunc(height, peer)
}
