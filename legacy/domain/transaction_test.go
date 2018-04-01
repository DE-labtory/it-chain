package domain

import (
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCreateNewTransactionTest(t *testing.T) {
	tx := CreateNewTransaction(strconv.Itoa(1), strconv.Itoa(1), General, time.Now(), SetTxData("", Invoke, SetTxMethodParameters(0, "", []string{""}), ""))
	assert.Equal(t, Status_TRANSACTION_UNKNOWN, tx.TransactionStatus)
}

func TestSetTxMethodParameters(t *testing.T) {
	var par = SetTxMethodParameters(0, "void", []string{""})
	assert.Equal(t, 0, par.ParamsType)
}

func TestSetTxData(t *testing.T) {
	var txdata = SetTxData("temp", Query, SetTxMethodParameters(0, "void", []string{""}), "111")
	assert.Equal(t, Query, txdata.Method)
}

func TestTransaction_Validate(t *testing.T) {
	tx := CreateNewTransaction(strconv.Itoa(1), strconv.Itoa(1), General, time.Now(), SetTxData("", Invoke, SetTxMethodParameters(0, "", []string{""}), ""))
	tx.GenerateHash()
	assert.Equal(t, true, tx.Validate())
}

//func TestTransaction_SignHash(t *testing.T) {
//	tx := CreateNewTransaction(strconv.Itoa(1), strconv.Itoa(1), General, time.Now(), SetTxData("", Invoke, SetTxMethodParameters(0, "", []string{""}), ""))
//	_, err := tx.SignHash()
//	assert.NoError(t, err)
//}
