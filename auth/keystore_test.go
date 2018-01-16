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
	err = ks.StoreKey(ecdsaKey)
	assert.NoError(t, err)

	err = ks.StoreKey(nil)
	assert.Error(t, err)

}