package blockchain

import "errors"

var ErrSetConfig = errors.New("error when get Config")
var ErrDecodingEmptyBlock = errors.New("Empty Block decoding failed")
var ErrBuildingTxSeal = errors.New("Error in building tx seal")
var ErrBuildingSeal = errors.New("Error in building seal")
var ErrCreatingEvent = errors.New("Error in creating event")
var ErrOnEvent = errors.New("Error on event")
var ErrGetLastCommitedBlock = errors.New("Error getting last commited block")
