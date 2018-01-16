package auth

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"crypto/elliptic"
)

func TestRSAKeyGenerator_KeyGenerate(t *testing.T) {

	keygen := &rsaKeyGenerator{1024}
	key, err := keygen.KeyGenerate(nil)
	assert.NoError(t, err)
	assert.NotNil(t, key)

	rsaKey, valid := key.(*rsaPrivateKey)
	assert.True(t, valid)
	assert.NotNil(t, rsaKey)
	assert.Equal(t, rsaKey.priv.N.BitLen(), 1024)

}

func TestECDSAKeyGenerator_KeyGenerate(t *testing.T) {

	keygen := &ecdsaKeyGenerator{elliptic.P256()}
	key, err := keygen.KeyGenerate(nil)
	assert.NoError(t, err)
	assert.NotNil(t, key)

	ecdsaKey, valid := key.(*ecdsaPrivateKey)
	assert.True(t, valid)
	assert.NotNil(t, ecdsaKey)
	assert.Equal(t, ecdsaKey.priv.Curve, elliptic.P256())

}

func TestRSAKeyGenerator_InvalidInput(t *testing.T) {

	keygen := &rsaKeyGenerator{-1}

	_, err := keygen.KeyGenerate(nil)
	assert.Error(t, err)

}

func TestECDSAKeyGenerator_NilInput(t *testing.T) {

	keygen := &ecdsaKeyGenerator{nil}

	_, err := keygen.KeyGenerate(nil)
	assert.Error(t, err)

}