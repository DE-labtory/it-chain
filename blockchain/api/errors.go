package api

import "errors"

var ErrNilBlock = errors.New("block is nil")
var ErrSyncProcessing = errors.New("SyncWithPeer is in progress")
var ErrGetLastBlock = errors.New("failed get last block")
var ErrGetRandomPeer = errors.New("error while getting random peer")
var ErrSyncWithPeer = errors.New("error while synchronizing with a given peer")
