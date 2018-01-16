package auth

import (
	"crypto/rsa"
	"crypto/ecdsa"
	"crypto/elliptic"
	"github.com/stretchr/testify/assert"
	"testing"
	"crypto/rand"
)

func TestRSAPublicKeyToPem(t *testing.T) {

	generatedRSAKey, err := rsa.GenerateKey(rand.Reader, 1024)
	assert.NoError(t, err)

	rsaKey := &rsaPrivateKey{generatedRSAKey}
	pub, err := rsaKey.PublicKey()
	assert.NoError(t, err)
	assert.NotNil(t, pub)

	data, err := PublicKeyToPem(pub.(*rsaPublicKey).pub)
	assert.NoError(t, err)
	assert.NotNil(t, data)

}

func TestRSAPrivateKeyToPem(t *testing.T) {

	generatedRSAKey, err := rsa.GenerateKey(rand.Reader, 1024)
	assert.NoError(t, err)

	rsaKey := &rsaPrivateKey{generatedRSAKey}

	data, err := PrivateKeyToPem(rsaKey.priv)
	assert.NoError(t, err)
	assert.NotNil(t, data)

}

func TestECDSAPublicKeyToPem(t *testing.T) {

	generatedECDSAKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	assert.NoError(t, err)

	ecdsaKey := &ecdsaPrivateKey{generatedECDSAKey}
	pub, err := ecdsaKey.PublicKey()
	assert.NoError(t, err)
	assert.NotNil(t, pub)

	data, err := PublicKeyToPem(pub.(*ecdsaPublicKey).pub)
	assert.NoError(t, err)
	assert.NotNil(t, data)

}

func TestECDSAPrivateKeyToPem(t *testing.T) {

	generatedECDSAKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	assert.NoError(t, err)

	rsaKey := &ecdsaPrivateKey{generatedECDSAKey}

	data, err := PrivateKeyToPem(rsaKey.priv)
	assert.NoError(t, err)
	assert.NotNil(t, data)

}