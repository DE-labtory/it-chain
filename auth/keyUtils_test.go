package auth

import (
	"crypto/rsa"
	"crypto/ecdsa"
	"crypto/elliptic"
	"github.com/stretchr/testify/assert"
	"testing"
	"crypto/rand"
	"os"
	"io/ioutil"
	"encoding/hex"
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

func TestPEMToPrivatePublicKey(t *testing.T) {

	cryp, err := NewCrypto("")
	assert.NoError(t, err)

	pri, pub, err := cryp.GenerateKey(&RSAKeyGenOpts{})
	assert.NoError(t, err)
	assert.NotNil(t, pri)
	assert.NotNil(t, pub)

	defer os.RemoveAll("./KeyRepository")

	path := "./KeyRepository/" + hex.EncodeToString(pri.SKI()) + "_pri"

	data, err := ioutil.ReadFile(path)
	assert.NoError(t, err)
	assert.NotNil(t, data)

	key, err := PEMToPrivateKey(data)
	assert.NoError(t, err)
	assert.NotNil(t, key)

	path = "./KeyRepository/" + hex.EncodeToString(pub.SKI()) + "_pub"

	data, err = ioutil.ReadFile(path)
	assert.NoError(t, err)
	assert.NotNil(t, data)

	key, err = PEMToPublicKey(data)
	assert.NoError(t, err)
	assert.NotNil(t, key)

}