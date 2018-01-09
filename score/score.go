package score

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/junbeomlee/it-chain/service/blockchain"
	"sync"
)


// = blockchain.BlockChain
// *임시 블록체인 구조체
// (init,deploy 등을 메소드로 구현했기 때문에 임시 구조체 사용 : 이 후, 블록체인 구조체와 독립적으로 돌아 갈 수 있도록 인자형태로 받는게 좋을 듯)
type SimpleChaincode struct {
	sync.RWMutex
	Header *blockchain.ChainHeader
	Blocks []*blockchain.Block
}

const CONTRACTSTATEKEY string = "ContractStateKey"
const SCOREVERSION string = "1.0"

/* asset and contract state
----------------------------------------*/
type ContractState struct {
	Version string `json:"version"`
}

type AssetState struct {
	AssetID     *string      `json:"assetID,omitempty"`
	Temperature *float64     `json:"temperature,omitempty"`
	Carrier     *string      `json:"carrier,omitempty"`
}

var contractState = ContractState{SCOREVERSION}

/* deploy callback mode (블록체인 구조체를 인자로 받으면서 독립시킬 필요있음)
----------------------------------------*/
func (t *SimpleChaincode) Init(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {
	var stateArg ContractState
	var err error
	if len(args) != 1 {
		return nil, errors.New("init expects one argument, a JSON string with tagged version string")
	}
	err = json.Unmarshal([]byte(args[0]), &stateArg)
	if err != nil {
		return nil, errors.New("Version argument unmarshal failed: " + fmt.Sprint(err))
	}
	if stateArg.Version != SCOREVERSION {
		return nil, errors.New("Contract version " + SCOREVERSION + " must match version argument: " + stateArg.Version)
	}
	contractStateJSON, err := json.Marshal(stateArg)
	if err != nil {
		return nil, errors.New("Marshal failed for contract state" + fmt.Sprint(err))
	}
	err = stub.PutState(CONTRACTSTATEKEY, contractStateJSON)
	if err != nil {
		return nil, errors.New("Contract state failed PUT to ledger: " + fmt.Sprint(err))
	}
	return nil, nil
}


/* deploy and invoke callback mode
----------------------------------------*/
func (t *SimpleChaincode) Invoke() {
}

/* query callback mode
----------------------------------------*/
func (t *SimpleChaincode) Query() {
}