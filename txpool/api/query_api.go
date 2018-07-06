package api

import "github.com/it-chain/it-chain-Engine/txpool"

//this api is for querying data
//this api should be used only for querying
type TxpoolQueryApi struct {
	txRepository     txpool.TransactionRepository
	leaderRepository txpool.LeaderRepository
}

func (t TxpoolQueryApi) GetAllCreatedTransactions() ([]txpool.Transaction, error) {

	txs, err := t.txRepository.FindAll()

	if err != nil {
		return nil, err
	}

	return filter(txs, func(transaction txpool.Transaction) bool {
		return transaction.Stage == txpool.CREATED
	}), nil
}

func (t TxpoolQueryApi) GetAllStagedTransactions() ([]txpool.Transaction, error) {

	txs, err := t.txRepository.FindAll()

	if err != nil {
		return nil, err
	}

	return filter(txs, func(transaction txpool.Transaction) bool {
		return transaction.Stage == txpool.STAGED
	}), nil
}

func (t TxpoolQueryApi) GetLeader() txpool.Leader {
	return t.leaderRepository.GetLeader()
}

func filter(vs []txpool.Transaction, f func(txpool.Transaction) bool) []txpool.Transaction {
	vsf := make([]txpool.Transaction, 0)
	for _, v := range vs {
		if f(v) {
			vsf = append(vsf, v)
		}
	}
	return vsf
}
