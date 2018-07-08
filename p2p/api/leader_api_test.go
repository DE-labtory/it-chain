package api_test
//
//import (
//	"testing"
//
//	"github.com/it-chain/it-chain-Engine/p2p"
//	"github.com/it-chain/it-chain-Engine/p2p/api"
//	"github.com/magiconair/properties/assert"
//)
//
//type MockReadOnlyLeaderRepository struct{}
//
//func (mrolr MockReadOnlyLeaderRepository) GetLeader() p2p.Leader { return p2p.Leader{} }
//
//func TestLeaderApi_UpdateLeader(t *testing.T) {
//
//	tests := map[string]struct {
//		input p2p.Leader
//		err   error
//	}{
//		"empty leader id test": {
//			input: p2p.Leader{
//				LeaderId: p2p.LeaderId{
//					Id: "",
//				},
//			},
//			err: api.ErrEmptyLeaderId,
//		},
//		"first leader update test": {
//			input: p2p.Leader{
//				LeaderId: p2p.LeaderId{
//					Id: "1",
//				},
//			},
//			err: nil,
//		},
//	}
//
//	leaderApi := SetupLeaderApi()
//
//	for testName, test := range tests {
//
//		t.Logf("Running test case %s", testName)
//		err := leaderApi.UpdateLeader(test.input)
//		assert.Equal(t, err, test.err)
//	}
//}
//
//func SetupLeaderApi() *api.LeaderApi {
//
//	leaderRepository := MockReadOnlyLeaderRepository{}
//	leaderApi := api.NewLeaderApi(leaderRepository, &p2p.Peer{PeerId: p2p.PeerId{Id: "123"}})
//
//	return leaderApi
//}
