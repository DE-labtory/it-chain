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
	"encoding/hex"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	kitlog "github.com/go-kit/kit/log"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
)

var (
	// ErrBadRouting is returned when an expected path variable is missing.
	ErrBadRouting    = errors.New("inconsistent mapping between route and handler.")
	ErrBadConversion = errors.New("Conversion failed: invalid argument in url endpoint.")
)

func NewApiHandler(bqa *BlockQueryApi, iqa *ICodeQueryApi, iha *ICodeCommandApi, p *PeerQueryApi, cca *ConnectionCommandApi, logger kitlog.Logger) http.Handler {

	r := mux.NewRouter()

	be := MakeBlockchainEndpoints(bqa)
	ie := MakeIcodeEndpoints(iha, iqa)
	ce := MakePeerEndpoints(p, cca)
	te := MakeTransactionEndpoints(iha)

	opts := []kithttp.ServerOption{
		kithttp.ServerErrorLogger(logger),
	}

	// GET     /blocks/						retrieves all blocks committed
	// GET     /blocks?height=:height		retrieves a particular block committed
	// GET     /blocks/:seal				retrieves a particular block committed

	r.Methods("GET").Path("/blocks").Handler(kithttp.NewServer(
		be.FindAllCommittedBlocksEndpoint,
		decodeFindAllCommittedBlocksRequest,
		encodeResponse,
		opts...,
	))

	r.Methods("GET").Path("/blocks/{id}").Handler(kithttp.NewServer(
		be.FindCommittedBlockBySealEndpoint,
		decodeFindCommittedBlockBySealRequest,
		encodeResponse,
		opts...,
	))

	// GET		/icodes																retrieves all icodes deployed
	// GET		/icodes?amqpUrl=:amqpUrl											retrieves all icodes deployed using particular amqp url not in config
	// POST		/icodes																deploy icode. about post body information, see decodeDeployIcodeRequest
	// DELETE	/icodes/{icodeId}													unDeploy icode that match icodeId

	r.Methods("GET").Path("/icodes").Handler(kithttp.NewServer(
		ie.GetIcodeListEndpoint,
		decodeGetICodeListRequest,
		encodeResponse,
		opts...))

	r.Methods("POST").Path("/icodes").Handler(kithttp.NewServer(
		ie.DeployIcodeEndpoint,
		decodeDeployIcodeRequest,
		encodeResponse,
		opts...))

	r.Methods("DELETE").Path("/icodes/{id}").Handler(kithttp.NewServer(
		ie.UnDeployIcodeEndpoint,
		decodeUnDeployIcodeRequest,
		encodeResponse,
		opts...))

	// GET		/transactions			get all uncommitted transactions
	// POST 	/transactions			create transaction
	r.Methods("POST").Path("/transactions").Handler(kithttp.NewServer(
		te.CreateTransactionEndpoint,
		decodeCreateTransactionRequest,
		encodeResponse,
		opts...))

	// GET		/peers			retrieves all peers
	// GET		/peers/{id}		retrieves peers that match id
	// POST		/peers			dial or join network to address. about post body information, see decodeCreateConnectionRequest
	r.Methods("GET").Path("/peers").Handler(kithttp.NewServer(
		ce.FindAllPeerEndpoint,
		decodeFindAllConnectionRequest,
		encodeResponse,
		opts...,
	))

	r.Methods("GET").Path("/peers/{id}").Handler(kithttp.NewServer(
		ce.FindPeerByIdEndpoint,
		decodeFindConnectionByIdRequest,
		encodeResponse,
		opts...,
	))

	r.Methods("POST").Path("/peers").Handler(kithttp.NewServer(
		ce.CreateConnectionEndpoint,
		decodeCreateConnectionRequest,
		encodeResponse,
		opts...))

	return r
}

/*
txpool
*/
// this return nil because this request body is empty
func decodeFindAllUncommittedTransactionsRequest(_ context.Context, r *http.Request) (interface{}, error) {
	return nil, nil
}

func decodeCreateTransactionRequest(_ context.Context, r *http.Request) (interface{}, error) {
	body := CreateTransactionRequest{}
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

/*
block chain
*/
func decodeFindAllCommittedBlocksRequest(_ context.Context, r *http.Request) (interface{}, error) {
	if heightStr := r.URL.Query().Get("height"); heightStr != "" {
		if heightStr == "-1" {
			return FindLastCommittedBlockRequest{}, nil
		}

		height, err := strconv.ParseUint(heightStr, 10, 64)
		if err != nil {
			return nil, ErrBadConversion
		}
		return FindCommittedBlockByHeightRequest{Height: height}, nil
	}
	// length of query string is zero => means that there are no restful params
	return nil, nil
}

func decodeFindCommittedBlockBySealRequest(_ context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)

	sealStr, ok := vars["id"]
	if !ok {
		return nil, ErrBadRouting
	}

	seal, err := hex.DecodeString(sealStr)
	if err != nil {
		return nil, ErrBadConversion
	}

	seal = []byte(seal)

	return FindCommittedBlockBySealRequest{Seal: seal}, nil
}

/*
ivm
*/
func decodeGetICodeListRequest(_ context.Context, r *http.Request) (interface{}, error) {
	return r, nil
}

func decodeDeployIcodeRequest(_ context.Context, r *http.Request) (interface{}, error) {
	body := DeployIcodeRequest{}
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		return nil, err
	}

	if body.GitUrl == "" {
		return nil, ErrBadConversion
	}

	return body, nil
}

func decodeUnDeployIcodeRequest(_ context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	amqpUrl := r.FormValue("amqpUrl")

	icodeId, ok := vars["id"]
	if !ok {
		return nil, ErrBadConversion
	}

	return UnDeployIcodeRequest{
		IvmRequest: IvmRequest{
			AmqpUrl: amqpUrl,
		},
		ICodeId: icodeId,
	}, nil
}

/*
grpc gateway
*/
func decodeFindAllConnectionRequest(_ context.Context, r *http.Request) (interface{}, error) {
	return nil, nil
}

func decodeFindConnectionByIdRequest(_ context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)

	peerId, ok := vars["id"]
	if !ok {
		return nil, ErrBadRouting
	}

	return FindPeerByIdRequest{ID: peerId}, nil
}

func decodeCreateConnectionRequest(_ context.Context, r *http.Request) (interface{}, error) {
	body := CreateConnectionRequest{}
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		return nil, err
	}

	if body.Address == "" || (body.Type != "dial" && body.Type != "join") {
		return nil, ErrBadConversion
	}

	return body, nil

}

func encodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {

	if e, ok := response.(errorer); ok && e.error() != nil {
		encodeError(ctx, e.error(), w)
		return nil
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}

type errorer interface {
	error() error
}

// encode errors from business-logic
func encodeError(_ context.Context, err error, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	switch err {
	//case cargo.ErrUnknown:
	//	w.WriteHeader(http.StatusNotFound)
	//case ErrInvalidArgument:
	//	w.WriteHeader(http.StatusBadRequest)
	default:
		w.WriteHeader(http.StatusInternalServerError)
	}
	json.NewEncoder(w).Encode(map[string]interface{}{
		"error": err.Error(),
	})
}
