package auth

import (
	"testing"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"github.com/stretchr/testify/assert"
	"os"
	"path/filepath"
	"encoding/hex"
	"crypto/rsa"
)

func TestKeyStore_StoreKey(t *testing.T) {
	ks := &keyStorer{"./testStorer"}

	defer os.RemoveAll(ks.path)

	rsaRawKey, err := rsa.GenerateKey(rand.Reader, 1024)

	rsaPriKey := &rsaPrivateKey{rsaRawKey}
	err = ks.Store(rsaPriKey)
	assert.NoError(t, err)

	rsaPubKey := &rsaPublicKey{&rsaPriKey.priv.PublicKey}
	err = ks.Store(rsaPubKey)
	assert.NoError(t, err)

	ecdsaRawKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)

	ecdsaPriKey := &ecdsaPrivateKey{ecdsaRawKey}
	err = ks.Store(ecdsaPriKey)
	assert.NoError(t, err)

	ecdsaPubKey := &ecdsaPublicKey{&ecdsaPriKey.priv.PublicKey}
	err = ks.Store(ecdsaPubKey)
	assert.NoError(t, err)

	// check whether file is exist
	path := filepath.Join(ks.path, hex.EncodeToString(rsaPriKey.SKI()) + "_pri")
	assert.FileExists(t, path)

	path = filepath.Join(ks.path, hex.EncodeToString(rsaPubKey.SKI()) + "_pub")
	assert.FileExists(t, path)

	path = filepath.Join(ks.path, hex.EncodeToString(ecdsaPriKey.SKI()) + "_pri")
	assert.FileExists(t, path)

	path = filepath.Join(ks.path, hex.EncodeToString(ecdsaPubKey.SKI()) + "_pub")
	assert.FileExists(t, path)

	// Invalid input for KeyStore
	err = ks.Store(nil)
	assert.Error(t, err)
}

func TestKeyStorer_InvalidInput(t *testing.T) {

	ks := &keyStorer{"./testStorer"}

	defer os.RemoveAll(ks.path)

	rsaRawKey, err := rsa.GenerateKey(rand.Reader, 1024)

	rsaPriKey := &rsaPrivateKey{rsaRawKey}
	err = ks.storePublicKey(rsaPriKey)
	assert.Error(t, err)

	rsaPubKey := &rsaPublicKey{&rsaPriKey.priv.PublicKey}
	err = ks.storePrivateKey(rsaPubKey)
	assert.Error(t, err)

	ecdsaRawKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)

	ecdsaPriKey := &ecdsaPrivateKey{ecdsaRawKey}
	err = ks.storePublicKey(ecdsaPriKey)
	assert.Error(t, err)

	ecdsaPubKey := &ecdsaPublicKey{&ecdsaPriKey.priv.PublicKey}
	err = ks.storePrivateKey(ecdsaPubKey)
	assert.Error(t, err)

}