package leveldb

import (
	"testing"

	"github.com/it-chain/it-chain-Engine/p2p"
	"github.com/magiconair/properties/assert"
)

func TestPeerRepository_Save(t *testing.T) {
	tests := map[string] struct{
		input struct{
			firstInput p2p.Peer
			secondInput p2p.Peer
		}
		err error
	}{
		"success":{
			input: struct {
				firstInput  p2p.Peer
				secondInput p2p.Peer
			}{firstInput: p2p.Peer{PeerId:p2p.PeerId{Id:"1"}}, secondInput: p2p.Peer{PeerId:p2p.PeerId{Id:"2"}}},
			err:nil,
		},
		"empty peer id test":{
			input: struct {
				firstInput  p2p.Peer
				secondInput p2p.Peer
			}{firstInput: p2p.Peer{PeerId:p2p.PeerId{Id:"1"}}, secondInput: p2p.Peer{PeerId:p2p.PeerId{Id:""}}},
			err: ErrEmptyPeerId,
		},
		"already exist peer test":{
			input: struct {
				firstInput  p2p.Peer
				secondInput p2p.Peer
			}{firstInput: p2p.Peer{PeerId:p2p.PeerId{Id:"1"}}, secondInput: p2p.Peer{PeerId:p2p.PeerId{Id:"1"}}},
			err: ErrExistPeer,
		},
	}
	peerRepository := NewPeerRepository()
	ClearPeerTable()

	for testName, test := range tests{
		t.Logf("running test case %s", testName)
		peerRepository.Save(test.input.firstInput)
		err := peerRepository.Save(test.input.secondInput)
		assert.Equal(t, err, test.err)
	}

}

func TestPeerRepository_Remove(t *testing.T) {
	tests := map[string] struct{
		input p2p.PeerId
		err error
	}{
		"success":{
			input: p2p.PeerId{Id:"1"},
			err:nil,
		},
		"no matching peer test":{
			input: p2p.PeerId{Id:"1"},
			err: ErrNoMatchingPeer,
		},
		"empty peer id test":{
			input: p2p.PeerId{Id:""},
			err: ErrEmptyPeerId,
		},
	}
	peerRepository := NewPeerRepository()
	peer := p2p.Peer{PeerId:p2p.PeerId{Id:"1"}}
	peerRepository.Save(peer)

	for testName, test := range tests{
		t.Logf("running test case %s", testName)
		err := peerRepository.Remove(test.input)
		assert.Equal(t, err, test.err)
	}
}

func TestPeerRepository_FindAll(t *testing.T) {
	tests := map[string] struct{
		input p2p.Peer
		err error
	}{
		"empty peer list test":{
			input:p2p.Peer{PeerId:p2p.PeerId{Id:""}},
			err: ErrEmptyPeerTable,
		},
		"success":{
			input:p2p.Peer{PeerId:p2p.PeerId{Id:"1"}},
			err:nil,
		},
	}
	peerRepository := NewPeerRepository()
	ClearPeerTable()

	for testName, test := range tests{
		t.Logf("running test case %s", testName)
		ClearPeerTable()
		peerRepository.Save(test.input)

		_, err := peerRepository.FindAll()
		assert.Equal(t, err, test.err)
	}

}

func TestPeerRepository_FindById(t *testing.T) {
	tests := map[string] struct{
		input struct{
			insertPeer p2p.Peer
			searchPeer p2p.Peer
		}
		err error
	}{
		"success":{
			input: struct {
				insertPeer p2p.Peer
				searchPeer p2p.Peer
			}{insertPeer: p2p.Peer{PeerId:p2p.PeerId{Id:"1"}}, searchPeer: p2p.Peer{PeerId:p2p.PeerId{Id:"1"}}},
			err: nil,
		},
		"empty peer id test":{
			input: struct {
				insertPeer p2p.Peer
				searchPeer p2p.Peer
			}{insertPeer: p2p.Peer{PeerId:p2p.PeerId{Id:"1"}}, searchPeer: p2p.Peer{PeerId:p2p.PeerId{Id:""}}},
			err:ErrEmptyPeerId,
		},
		"no matching peer test":{
			input: struct {
				insertPeer p2p.Peer
				searchPeer p2p.Peer
			}{insertPeer: p2p.Peer{PeerId:p2p.PeerId{Id:"1"}}, searchPeer: p2p.Peer{PeerId:p2p.PeerId{Id:"2"}}},
			err:ErrNoMatchingPeer,
		},
	}
	peerRepository := NewPeerRepository()

	for testName, test := range tests{
		t.Logf("running test case %s", testName)
		ClearPeerTable()
		peerRepository.Save(test.input.insertPeer)
		_, err := peerRepository.FindById(test.input.searchPeer.PeerId)

		assert.Equal(t, err, test.err)
	}
}

func ClearPeerTable(){
	peerRepository := NewPeerRepository()
	peer := p2p.Peer{PeerId:p2p.PeerId{Id:"1"}}
	_, exist := peerTable[peer.PeerId.Id]
	if exist{
		peerRepository.Remove(peer.PeerId)
	}
}