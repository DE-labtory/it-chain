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
	"encoding/json"
	"net/http"

	kitlog "github.com/go-kit/kit/log"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
)

func BlockchainApiHandler(bqa BlockQueryApi, logger kitlog.Logger) http.Handler {

	opts := []kithttp.ServerOption{
		kithttp.ServerErrorLogger(logger),
	}

	findAllCommittedBlocksHandler := kithttp.NewServer(
		makeFindCommittedBlocksEndpoint(bqa),
		decodeFindAllCommittedBlocksRequest,
		encodeResponse,
		opts...,
	)

	r := mux.NewRouter()

	r.Handle("/blocks", findAllCommittedBlocksHandler).Methods("GET")

	return r
}

func ICodeApiHandler(api ICodeQueryApi, logger kitlog.Logger) http.Handler {
	opts := []kithttp.ServerOption{
		kithttp.ServerErrorLogger(logger),
	}

	findAllMetasHandler := kithttp.NewServer(
		makeFindAllMetaEndpoint(api),
		decodeFindAllMetaRequest,
		encodeResponse,
		opts...,
	)
	r := mux.NewRouter()

	r.Handle("/metas", findAllMetasHandler).Methods("GET")

	return r
}

// this return nil because this request body is empty
func decodeFindAllUncommittedTransactionsRequest(_ context.Context, r *http.Request) (interface{}, error) {
	return nil, nil
}

func decodeFindAllCommittedBlocksRequest(_ context.Context, r *http.Request) (interface{}, error) {
	return nil, nil
}

func decodeFindAllMetaRequest(_ context.Context, r *http.Request) (interface{}, error) {
	return nil, nil
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
