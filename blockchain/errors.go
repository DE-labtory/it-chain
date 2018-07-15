package blockchain

import "errors"

var ErrTxListMarshal = errors.New("tx list marshal failed")
var ErrTxListUnmarshal = errors.New("tx list unmarshal failed")
var ErrSetConfig = errors.New("error when get Config")
var ErrDecodingEmptyBlock = errors.New("Empty Block decoding failed")
var ErrBuildingTxSeal = errors.New("Error in building tx seal")
var ErrBuildingSeal = errors.New("Error in building seal")
var ErrCreatingEvent = errors.New("Error in creating event")
var ErrOnEvent = errors.New("Error on event")
var ErrSyncedCheck = errors.New("Error when checking if synchronizing is needed")
var ErrConstruct = errors.New("Error when constructing")
