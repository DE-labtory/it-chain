package api_gateway

import (
	"context"

	"github.com/go-kit/kit/endpoint"
)

//This file is based on the following sample.
// https://github.com/marcusolsson/goddd/blob/master/booking/endpoint.go
func makeFindUncommittedTransactionsEndpoint(t TransactionQueryApi) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		txs, err := t.FindUncommittedTransactions()

		if err != nil {
			return nil, err
		}

		return txs, nil
	}
}
