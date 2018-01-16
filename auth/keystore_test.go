package auth

import (
	"testing"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"github.com/stretchr/testify/assert"
	"crypto/rsa"
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

func TestRSAPublicKeyToPem(t *testing.T) {

	generatedRSAKey, err := rsa.GenerateKey(rand.Reader, 1024)
	assert.NoError(t, err)

	rsaKey := &rsaPrivateKey{generatedRSAKey}
	pub, err := rsaKey.PublicKey()
	assert.NoError(t, err)
	assert.NotNil(t, pub)

	data, err := publicKeyToPem(pub.(*rsaPublicKey).pub)
	assert.NoError(t, err)
	assert.NotNil(t, data)

}

func TestRSAPrivateKeyToPem(t *testing.T) {

	generatedRSAKey, err := rsa.GenerateKey(rand.Reader, 1024)
	assert.NoError(t, err)

	rsaKey := &rsaPrivateKey{generatedRSAKey}

	data, err := privateKeyToPem(rsaKey.priv)
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

	data, err := publicKeyToPem(pub.(*ecdsaPublicKey).pub)
	assert.NoError(t, err)
	assert.NotNil(t, data)

}

func TestECDSAPrivateKeyToPem(t *testing.T) {

	generatedECDSAKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	assert.NoError(t, err)

	rsaKey := &ecdsaPrivateKey{generatedECDSAKey}

	data, err := privateKeyToPem(rsaKey.priv)
	assert.NoError(t, err)
	assert.NotNil(t, data)

}