package api

import "errors"

var ErrNilBlock = errors.New("block is nil")
var ErrSyncProcessing = errors.New("SyncWithPeer is in progress")
var ErrGetLastBlock = errors.New("failed get last block")
