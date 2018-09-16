package adapter_test

import (
	"net/http"
	"testing"

	"time"

	"log"

	"encoding/json"

	"github.com/it-chain/engine/blockchain"
	"github.com/it-chain/engine/blockchain/infra/adapter"
	"github.com/it-chain/engine/blockchain/test/mock"
	"github.com/stretchr/testify/assert"
)

func TestHttpBlockAdapter_GetLastBlock(t *testing.T) {

	lastBlock := mock.GetNewBlock([]byte("prevseal"), 0)

	go createServer(*lastBlock)

	time.Sleep(3)

	a := adapter.HttpBlockAdapter{}

	_, err := a.GetLastBlock("127.0.0.1:8080/")
	assert.NoError(t, err)

	//wg := sync.WaitGroup{}
	//wg.Add(1)
	//wg.Wait()

}

func createServer(block blockchain.DefaultBlock) {

	b, _ := json.Marshal(block)

	handler := func(w http.ResponseWriter, req *http.Request) {
		if lastStr := req.URL.Query().Get("last"); lastStr == "true" {
			w.Write(b)
		}
	}

	mux := http.NewServeMux()

	mux.HandleFunc("/blocks", handler)

	log.Fatal(http.ListenAndServe(":8080", mux))

}
