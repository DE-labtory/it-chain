package auth

import (
	"crypto/rsa"
	"crypto/ecdsa"
	"crypto/elliptic"
	"github.com/stretchr/testify/assert"
	"testing"
	"crypto/rand"
)

func TestRSAPublicKeyToPEM(t *testing.T) {

	generatedRSAKey, err := rsa.GenerateKey(rand.Reader, 1024)
	assert.NoError(t, err)

	rsaKey := &rsaPrivateKey{generatedRSAKey}
	pub, err := rsaKey.PublicKey()
	assert.NoError(t, err)
	assert.NotNil(t, pub)

	data, err := PublicKeyToPEM(pub.(*rsaPublicKey))
	assert.NoError(t, err)
	assert.NotNil(t, data)

}

func TestRSAPrivateKeyToPEM(t *testing.T) {

	generatedRSAKey, err := rsa.GenerateKey(rand.Reader, 1024)
	assert.NoError(t, err)

	rsaKey := &rsaPrivateKey{generatedRSAKey}

	data, err := PrivateKeyToPEM(rsaKey)
	assert.NoError(t, err)
	assert.NotNil(t, data)

}

func TestECDSAPublicKeyToPEM(t *testing.T) {

	generatedECDSAKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	assert.NoError(t, err)

	ecdsaKey := &ecdsaPrivateKey{generatedECDSAKey}
	pub, err := ecdsaKey.PublicKey()
	assert.NoError(t, err)
	assert.NotNil(t, pub)

	data, err := PublicKeyToPEM(pub.(*ecdsaPublicKey))
	assert.NoError(t, err)
	assert.NotNil(t, data)

}

func TestECDSAPrivateKeyToPEM(t *testing.T) {

	generatedECDSAKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	assert.NoError(t, err)

	rsaKey := &ecdsaPrivateKey{generatedECDSAKey}

	data, err := PrivateKeyToPEM(rsaKey)
	assert.NoError(t, err)
	assert.NotNil(t, data)

}