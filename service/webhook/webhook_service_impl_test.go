package webhook

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"it-chain/domain"
)

func TestNewWebhookService(t *testing.T) {

	_, err := NewWebhookService()
	assert.NoError(t, err)

}

func TestWebhookServiceImpl_SendConfirmedBlock(t *testing.T) {

	ws, err := NewWebhookService()
	assert.NoError(t, err)

	err = ws.SendConfirmedBlock(&domain.Block{})
	assert.NoError(t, err)

	err = ws.SendConfirmedBlock(nil)
	assert.Error(t, err)

}