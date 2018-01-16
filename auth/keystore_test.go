package auth

import (
	"testing"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"github.com/stretchr/testify/assert"
)

func TestKeyStore_StoreKey(t *testing.T) {
	ks := &keyStore{}

	generatedKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	assert.NoError(t, err)
	assert.NotNil(t, generatedKey)

	ecdsaKey := &ecdsaPrivateKey{generatedKey}
	err = ks.StoreKey(ecdsaKey, PEM)
	assert.NoError(t, err)

	err = ks.StoreKey(nil, PEM)
	assert.Error(t, err)

	err = ks.StoreKey(ecdsaKey, -1)
	assert.Error(t, err)

}