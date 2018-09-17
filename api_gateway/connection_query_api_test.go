package api_gateway_test

import (
	"strconv"
	"testing"

	"github.com/it-chain/engine/api_gateway"
	"github.com/it-chain/engine/common/event"
	"github.com/stretchr/testify/assert"
)

func TestConnectionQueryApi_GetAllConnectionList(t *testing.T) {
	repo := api_gateway.NewConnectionRepository()
	api := api_gateway.NewConnectionQueryApi(*repo)

	for i := 0; i < 10; i++ {
		id := strconv.Itoa(i)
		err := repo.Save(api_gateway.Connection{
			ConnectionID: id,
			Address:      "address" + id,
		})
		assert.NoError(t, err)
	}

	connectionList, err := api.GetAllConnectionList()
	assert.NoError(t, err)

	for i := 0; i < 11; i++ {
		var id string
		if i != 10 {
			id = connectionList[i].ConnectionID
		} else {
			id = strconv.Itoa(i)
		}

		connection, err := api.GetConnectionByID(id)
		assert.NoError(t, err)

		if i != 10 {
			assert.NotEqual(t, api_gateway.Connection{}, connection)
		} else {
			assert.Equal(t, api_gateway.Connection{}, connection)
		}
	}
}

func TestConnectionQueryApi_GetConnectionByID(t *testing.T) {
	repo := api_gateway.NewConnectionRepository()
	api := api_gateway.NewConnectionQueryApi(*repo)
	ids := [4]string{"it-chain", "engine", "api_gateway", "makehoney"}

	for i := 0; i < len(ids); i++ {
		err := repo.Save(api_gateway.Connection{
			ConnectionID: ids[i],
			Address:      "address",
		})
		assert.NoError(t, err)
	}

	for i := 0; i < len(ids); i++ {
		connection, err := api.GetConnectionByID(ids[i])
		assert.NoError(t, err)

		assert.NotEqual(t, api_gateway.Connection{}, connection)
	}
	connection, err := api.GetConnectionByID("wrongID")
	assert.NoError(t, err)
	assert.Equal(t, api_gateway.Connection{}, connection)
}

func TestConnectionRepository_Save(t *testing.T) {
	repo := api_gateway.NewConnectionRepository()
	err := repo.Save(api_gateway.Connection{
		ConnectionID: "1",
		Address:      "address",
	})
	assert.NoError(t, err)

	tests := map[string]struct {
		Input  api_gateway.Connection
		Output error
	}{
		"success": {
			Input: api_gateway.Connection{
				ConnectionID: "0",
				Address:      "address",
			},
			Output: nil,
		},
		"fail": {
			Input: api_gateway.Connection{
				ConnectionID: "1",
				Address:      "address",
			},
			Output: api_gateway.ErrConnectionExists,
		},
	}

	for testName, test := range tests {
		t.Logf("Running '%s' test, caseName: %s", t.Name(), testName)
		//given

		err := repo.Save(test.Input)
		if err != nil {
			assert.Equal(t, test.Output, err)
			continue
		}

		c, exist := repo.ConnectionTable[test.Input.ConnectionID]
		assert.True(t, exist)
		assert.Equal(t, test.Output, err)
		assert.Equal(t, test.Input.ConnectionID, c.ConnectionID)
		assert.Equal(t, test.Input.ConnectionID, c.ConnectionID)
	}

}

func TestConnectionRepository_Remove(t *testing.T) {
	repo := api_gateway.NewConnectionRepository()
	err := repo.Save(api_gateway.Connection{
		ConnectionID: "0",
		Address:      "address",
	})
	assert.NoError(t, err)

	c, exist := repo.ConnectionTable["0"]
	assert.True(t, exist)

	repo.Remove(c.ConnectionID)
	_, exist = repo.ConnectionTable["0"]
	assert.False(t, exist)
}

func TestConnectionEventListener_HandleConnectionCreatedEvent(t *testing.T) {
	repo := api_gateway.NewConnectionRepository()
	listener := api_gateway.NewConnectionEventListener(*repo)

	err := listener.HandleConnectionCreatedEvent(event.ConnectionCreated{
		ConnectionID: "0",
		Address:      "address",
	})
	assert.NoError(t, err)
}

func TestConnectionEventListener_HandleConnectionClosedEvent(t *testing.T) {
	repo := api_gateway.NewConnectionRepository()
	listener := api_gateway.NewConnectionEventListener(*repo)

	err := listener.HandleConnectionCreatedEvent(event.ConnectionCreated{
		ConnectionID: "0",
		Address:      "address",
	})
	assert.NoError(t, err)

	c, exist := repo.ConnectionTable["0"]
	assert.True(t, exist)

	listener.HandleConnectionClosedEvent(event.ConnectionClosed{
		ConnectionID: c.ConnectionID,
	})
	_, exist = repo.ConnectionTable["0"]
	assert.False(t, exist)
}
