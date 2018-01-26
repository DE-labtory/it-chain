package auth

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"crypto/elliptic"
)

func TestRSAKeyGenerator_GenerateKey(t *testing.T) {

	keygen := &rsaKeyGenerator{1024}
	pri, pub, err := keygen.GenerateKey(nil)
	assert.NoError(t, err)
	assert.NotNil(t, pri)
	assert.NotNil(t, pub)

	rsaPriKey, valid := pri.(*rsaPrivateKey)
	assert.True(t, valid)
	assert.NotNil(t, rsaPriKey)
	assert.Equal(t, rsaPriKey.priv.N.BitLen(), 1024)

	rsaPubKey, valid := pub.(*rsaPublicKey)
	assert.True(t, valid)
	assert.NotNil(t, rsaPubKey)
	assert.Equal(t, rsaPubKey.pub.N.BitLen(), 1024)

}

func TestECDSAKeyGenerator_GenerateKey(t *testing.T) {

	keygen := &ecdsaKeyGenerator{elliptic.P256()}
	pri, pub, err := keygen.GenerateKey(nil)
	assert.NoError(t, err)
	assert.NotNil(t, pri)
	assert.NotNil(t, pub)

	ecdsaPriKey, valid := pri.(*ecdsaPrivateKey)
	assert.True(t, valid)
	assert.NotNil(t, ecdsaPriKey)
	assert.Equal(t, ecdsaPriKey.priv.Curve, elliptic.P256())

	ecdsaPubKey, valid := pub.(*ecdsaPublicKey)
	assert.True(t, valid)
	assert.NotNil(t, ecdsaPubKey)
	assert.Equal(t, ecdsaPubKey.pub.Curve, elliptic.P256())

}

func TestRSAKeyGenerator_InvalidInput(t *testing.T) {

	keygen := &rsaKeyGenerator{-1}

	_, _, err := keygen.GenerateKey(nil)
	assert.Error(t, err)

}

func TestECDSAKeyGenerator_NilInput(t *testing.T) {

	keygen := &ecdsaKeyGenerator{nil}

	_, _, err := keygen.GenerateKey(nil)
	assert.Error(t, err)

}