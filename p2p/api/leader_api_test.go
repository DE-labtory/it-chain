package api_test

import (
	"testing"

	"github.com/it-chain/it-chain-Engine/p2p"
	"github.com/it-chain/it-chain-Engine/p2p/api"
)

func TestLeaderApi_UpdateLeaderWithAddress(t *testing.T) {

	tests := map[string]struct {
		input p2p.Leader
		err   error
	}{
		"empty leader id test": {
			input: p2p.Leader{
				LeaderId: p2p.LeaderId{
					Id: "",
				},
			},
			err: api.ErrEmptyLeaderId,
		},
		"first leader update test": {
			input: p2p.Leader{
				LeaderId: p2p.LeaderId{
					Id: "1",
				},
			},
			err: nil,
		},
	}
}

func TestLeaderApi_UpdateLeaderWithLongerPeerList(t *testing.T) {

}

//
//func SetupLeaderApi() *api.LeaderApi {
//
//	leaderApi := api.NewLeaderApi(leaderRepository, &p2p.Peer{PeerId: p2p.PeerId{Id: "123"}})
//
//	return leaderApi
//}
