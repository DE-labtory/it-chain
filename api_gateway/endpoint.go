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

package api_gateway

import (
	"context"

	"github.com/go-kit/kit/endpoint"
	"github.com/it-chain/iLogger"
)

type Endpoints struct {
	FindAllCommittedBlocksEndpoint   endpoint.Endpoint
	FindCommittedBlockBySealEndpoint endpoint.Endpoint
	FindAllMetaEndpoint              endpoint.Endpoint
	FindAllPeerEndpoint              endpoint.Endpoint
	FindPeerByIdEndpoint             endpoint.Endpoint
}

/*
 * returns endpoints
 */

func MakeBlockchainEndpoints(b *BlockQueryApi) Endpoints {
	return Endpoints{
		FindAllCommittedBlocksEndpoint:   makeFindAllCommittedBlocksEndpoint(b),
		FindCommittedBlockBySealEndpoint: makeFindCommittedBlockBySealEndpoint(b),
	}
}

func MakeIvmEndpoints(i *ICodeQueryApi) Endpoints {
	return Endpoints{
		FindAllMetaEndpoint: makeFindAllMetaEndpoint(i),
	}
}

func MakePeerEndpoints(p *PeerQueryApi) Endpoints {
	return Endpoints{
		FindAllPeerEndpoint:  makeFindAllPeerEndpoint(p),
		FindPeerByIdEndpoint: makeFindPeerByIdEndpoint(p),
	}
}

/*
 * blockchain
 */
func makeFindAllCommittedBlocksEndpoint(b *BlockQueryApi) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {

		switch v := request.(type) {

		case FindCommittedBlockByHeightRequest:
			return b.blockRepository.FindBlockByHeight(v.Height)

		case FindLastCommittedBlockRequest:
			return b.blockRepository.FindLastBlock()

		default:
			return b.blockRepository.FindAllBlock()
		}
	}
}

func makeFindCommittedBlockBySealEndpoint(b *BlockQueryApi) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(FindCommittedBlockBySealRequest)

		block, err := b.blockRepository.FindBlockBySeal(req.Seal)
		if err != nil {
			return nil, err
		}

		return block, nil
	}
}

//ivm
func makeFindAllMetaEndpoint(i *ICodeQueryApi) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		metas, err := i.iCodeRepository.FindAllMeta()
		if err != nil {
			iLogger.Error(&iLogger.Fields{"err_message": err.Error()}, "error while find all meta endpoint")
			return nil, err
		}

		return metas, nil
	}
}

type FindCommittedBlockByHeightRequest struct {
	Height uint64
}

type FindLastCommittedBlockRequest struct {
}

type FindCommittedBlockBySealRequest struct {
	Seal []byte
}

func makeFindAllPeerEndpoint(p *PeerQueryApi) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		peers := p.GetAllPeerList()
		return peers, nil
	}
}

func makeFindPeerByIdEndpoint(p *PeerQueryApi) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(FindPeerByIdRequest)
		peer, err := p.GetPeerByID(req.ID)
		if err != nil {
			iLogger.Error(nil, "[Api-gateway] Error while finding all connections")
			return nil, err
		}
		return peer, nil
	}
}

type FindPeerByIdRequest struct {
	ID string
}
