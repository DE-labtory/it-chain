package auth

import (
	"testing"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"github.com/stretchr/testify/assert"
)

func TestKeyStore_StoreKey(t *testing.T) {
	ks := &keyStorer{}

	generatedKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	assert.NoError(t, err)
	assert.NotNil(t, generatedKey)

	ecdsaKey := &ecdsaPrivateKey{generatedKey}
	err = ks.Store(ecdsaKey)
	assert.NoError(t, err)

	err = ks.Store(nil)
	assert.Error(t, err)

}