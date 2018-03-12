package service

import (
	"testing"
	"time"
	"github.com/stretchr/testify/assert"
	"os"
	"github.com/it-chain/it-chain-Engine/domain"
	"fmt"
)

func TestLedger_CreateBlock(t *testing.T) {
	path := "./test_db"
	ledger := NewLedger(path)
	defer func(){
		ledger.Close()
		os.RemoveAll(path)
	}()
	txList := make([]*domain.Transaction, 0)
	for i := 0; i < 999; i++{
		tx := domain.CreateNewTransaction("1", fmt.Sprintf("%d", i), domain.General, time.Now(), domain.SetTxData("", domain.Invoke, domain.SetTxMethodParameters(0, "", []string{""}), ""))
		tx.GenerateHash()
		txList = append(txList, tx)
	}

	_, err := ledger.CreateBlock(txList, "123")

	assert.NoError(t, err)
}

func TestLedger_AddBlock(t *testing.T) {
	path := "./test_db"
	ledger := NewLedger(path)
	defer func(){
		ledger.Close()
		os.RemoveAll(path)
	}()
	txList := make([]*domain.Transaction, 0)
	for i := 0; i < 999; i++{
		tx := domain.CreateNewTransaction("1", fmt.Sprintf("%d", i), domain.General, time.Now(), domain.SetTxData("", domain.Invoke, domain.SetTxMethodParameters(0, "", []string{""}), ""))
		tx.GenerateHash()
		txList = append(txList, tx)
	}
	blk, err := ledger.CreateBlock(txList, "123")
	assert.NoError(t, err)

	_, err = ledger.VerifyBlock(blk)
	assert.NoError(t, err)

	_, err = ledger.AddBlock(blk)
	assert.NoError(t, err)
}

func TestLedger_GetLastBlock(t *testing.T) {
	path := "./test_db"
	ledger := NewLedger(path)
	defer func(){
		ledger.Close()
		os.RemoveAll(path)
	}()
	txList := make([]*domain.Transaction, 0)
	for i := 0; i < 999; i++{
		tx := domain.CreateNewTransaction("1", fmt.Sprintf("%d", i), domain.General, time.Now(), domain.SetTxData("", domain.Invoke, domain.SetTxMethodParameters(0, "", []string{""}), ""))
		tx.GenerateHash()
		txList = append(txList, tx)
	}
	blk, err := ledger.CreateBlock(txList, "123")
	assert.NoError(t, err)

	_, err = ledger.VerifyBlock(blk)
	assert.NoError(t, err)

	_, err = ledger.AddBlock(blk)
	assert.NoError(t, err)

	lastblk, err := ledger.GetLastBlock()

	assert.Equal(t, blk.Header.BlockHash, lastblk.Header.BlockHash)
}

func TestLedger_LookUpBlock(t *testing.T) {
	path := "./test_db"
	ledger := NewLedger(path)
	defer func(){
		ledger.Close()
		os.RemoveAll(path)
	}()
	txList := make([]*domain.Transaction, 0)
	for i := 0; i < 999; i++{
		tx := domain.CreateNewTransaction("1", fmt.Sprintf("%d", i), domain.General, time.Now(), domain.SetTxData("", domain.Invoke, domain.SetTxMethodParameters(0, "", []string{""}), ""))
		tx.GenerateHash()
		txList = append(txList, tx)
	}
	blk, err := ledger.CreateBlock(txList, "123")
	assert.NoError(t, err)

	_, err = ledger.VerifyBlock(blk)
	assert.NoError(t, err)

	_, err = ledger.AddBlock(blk)
	assert.NoError(t, err)

	// 렛저에 저장되어 있지 않은 블록 인덱스를 검색할 때에 렛저에 없는 인덱스라고 에러가 리턴되지 않고 실행됨을 확인, 디비에서 해당 인덱스가 있는지 없는지 에러 리턴할 필요 있어보임
	retBlk, err := ledger.LookUpBlock(0)
	assert.NoError(t, err)
	assert.Equal(t, blk.Header.BlockHash, retBlk.Header.BlockHash)

	retBlk, err = ledger.LookUpBlock(blk.Header.BlockHash)
	assert.NoError(t, err)
	assert.Equal(t, blk.Header.BlockHash, retBlk.Header.BlockHash)
}

func TestLedger_VerifyBlock(t *testing.T) {
	path := "./test_db"
	ledger := NewLedger(path)
	defer func(){
		ledger.Close()
		os.RemoveAll(path)
	}()
	txList := make([]*domain.Transaction, 0)
	for i := 0; i < 999; i++{
		tx := domain.CreateNewTransaction("1", fmt.Sprintf("%d", i), domain.General, time.Now(), domain.SetTxData("", domain.Invoke, domain.SetTxMethodParameters(0, "", []string{""}), ""))
		tx.GenerateHash()
		txList = append(txList, tx)
	}
	blk, err := ledger.CreateBlock(txList, "123")
	assert.NoError(t, err)

	_, err = ledger.VerifyBlock(blk)
	assert.NoError(t, err)

	_, err = ledger.AddBlock(blk)
	assert.NoError(t, err)
}