package blockchain

import "errors"

var ErrTxListMarshal = errors.New("tx list marshal failed")
var ErrTxListUnmarshal = errors.New("tx list unmarshal failed")
var ErrGetConfig = errors.New("error when get Config")
var ErrDecodingEmptyBlock = errors.New("Empty Block decoding failed")
var ErrTransactionType = errors.New("Wrong transaction type")
