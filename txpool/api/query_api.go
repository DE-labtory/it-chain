package api

import "github.com/it-chain/it-chain-Engine/txpool"

//this api is for querying data
//this api should be used only for querying
type TxpoolQueryApi struct {
	txRepository     txpool.TransactionRepository
	leaderRepository txpool.LeaderRepository
}

func (t TxpoolQueryApi) GetAllTransactions() ([]txpool.Transaction, error) {

	txs, err := t.txRepository.FindAll()

	if err != nil {
		return nil, err
	}

	return txs, nil
}

func (t TxpoolQueryApi) GetLeader() txpool.Leader {
	return t.leaderRepository.GetLeader()
}
