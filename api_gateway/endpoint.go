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
	"github.com/pkg/errors"
)

type Endpoints struct {
	FindAllCommittedBlocksEndpoint   endpoint.Endpoint
	FindCommittedBlockBySealEndpoint endpoint.Endpoint

	FindAllPeerEndpoint  endpoint.Endpoint
	FindPeerByIdEndpoint endpoint.Endpoint

	GetIcodeListEndpoint  endpoint.Endpoint
	DeployIcodeEndpoint   endpoint.Endpoint
	UnDeployIcodeEndpoint endpoint.Endpoint

	FindAllUncommittedTransactionEndpoint endpoint.Endpoint
	CreateTransactionEndpoint             endpoint.Endpoint
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

func MakeIvmEndpoints(i *ICodeCommandApi, iqa *ICodeQueryApi) Endpoints {
	return Endpoints{
		GetIcodeListEndpoint:  makeFindAllICodeEndpoint(iqa),
		DeployIcodeEndpoint:   makeDeployIcodeEndpoint(i),
		UnDeployIcodeEndpoint: makeUnDeployIcodeEndpoint(i),
	}
}

func MakePeerEndpoints(p *PeerQueryApi) Endpoints {
	return Endpoints{
		FindAllPeerEndpoint:  makeFindAllPeerEndpoint(p),
		FindPeerByIdEndpoint: makeFindPeerByIdEndpoint(p),
	}
}

func MakeTransactionEndpoints(i *ICodeCommandApi) Endpoints {
	return Endpoints{
		FindAllUncommittedTransactionEndpoint: makeFindAllUncommittedTransactionEndpoint(),
		CreateTransactionEndpoint:             makeCreateTransactionEndpoint(i),
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

//icode
func makeFindAllICodeEndpoint(i *ICodeQueryApi) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		icodes, err := i.iCodeRepository.FindAllMeta()
		if err != nil {
			iLogger.Error(&iLogger.Fields{"err_message": err.Error()}, "error while find all icode endpoint")
			return nil, err
		}
		return icodes, nil
	}
}

func makeDeployIcodeEndpoint(i *ICodeCommandApi) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(DeployIcodeRequest)
		icodeId, err := i.deploy(req.AmqpUrl, req.GitUrl, req.SshPath, req.SshPassWord)
		if err != nil {
			iLogger.Error(&iLogger.Fields{"err_message": err.Error()}, "error while deploy icode endpoint")
			return nil, err
		}
		return icodeId, nil
	}
}

func makeUnDeployIcodeEndpoint(i *ICodeCommandApi) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(UnDeployIcodeRequest)
		err = i.unDeploy(req.AmqpUrl, req.ICodeId)
		if err != nil {
			iLogger.Error(&iLogger.Fields{"err_message": err.Error()}, "error while deploy icode endpoint")
			return nil, err
		}
		return nil, nil
	}
}

//transaction

//todo impl
func makeFindAllUncommittedTransactionEndpoint() endpoint.Endpoint {
	panic("impl plz")
}

func makeCreateTransactionEndpoint(i *ICodeCommandApi) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(CreateTransactionRequest)
		switch req.Type {
		case "invoke":
			txId, err := i.invoke(req.AmqpUrl, req.ICodeId, req.FuncName, req.Args)
			if err != nil {
				iLogger.Error(&iLogger.Fields{"err_message": err.Error()}, "error while invoke icode endpoint")
				return nil, err
			}
			return txId, err
		case "query":
			results, err := i.query(req.AmqpUrl, req.ICodeId, req.FuncName, req.Args)
			if err != nil {
				iLogger.Error(&iLogger.Fields{"err_message": err.Error()}, "error while query icode endpoint")
				return nil, err
			}
			return results, nil
		default:
			iLogger.Error(nil, "error while create transaction endpoint. unknown type err")
			return nil, errors.New("unknown type err")
		}
	}
}

//grpc gateway
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

//request struct
type IvmRequest struct {
	AmqpUrl string
}

// block chain request struct
type FindCommittedBlockByHeightRequest struct {
	Height uint64
}

type FindLastCommittedBlockRequest struct {
}

type FindCommittedBlockBySealRequest struct {
	Seal []byte
}

// ivm request struct

type DeployIcodeRequest struct {
	IvmRequest
	GitUrl      string
	SshPath     string
	SshPassWord string
}

type UnDeployIcodeRequest struct {
	IvmRequest
	ICodeId string
}

type CreateTransactionRequest struct {
	IvmRequest
	Type     string
	ICodeId  string
	FuncName string
	Args     []string
}

// grpc request struct
type FindConnectionByIdRequest struct {
	ConnectionId string
}
