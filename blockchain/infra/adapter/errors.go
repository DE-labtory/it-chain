package adapter

import "errors"

var ErrBlockInfoDeliver = errors.New("block info deliver failed")
var ErrGetBlock = errors.New("error when Getting block")
var ErrResponseBlock = errors.New("error when response block")
var ErrGetLastBlock = errors.New("error when get last block")
var ErrSyncCheckResponse = errors.New("error when sync check response")
var ErrEmptyNodeId = errors.New("empty nodeid proposed")
var ErrEmptyBlockSeal = errors.New("empty block seal")
var ErrBlockMissingProperties = errors.New("error when block miss some properties")
