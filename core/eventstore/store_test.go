package eventstore_test

import (
	"testing"

	"github.com/it-chain/it-chain-Engine/core/eventstore"
)

func TestGetConfiguration(t *testing.T) {

	eventstore.InitLevelDBStore()
}
