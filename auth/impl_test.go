package auth

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"crypto/rsa"
	"crypto/rand"
	"crypto/sha256"
	"crypto"
)

func TestNew(t *testing.T) {

	_, err := New()
	assert.NoError(t, err)

}

func TestImpl_RSASign(t *testing.T) {

	cryp, err := New()
	assert.NoError(t, err)

	generatedKey, err := rsa.GenerateKey(rand.Reader, 1024)
	assert.NoError(t, err)

	privateKey := &rsaPrivateKey{generatedKey}
	publicKey, err := privateKey.PublicKey()

	rawData := []byte("RSASign Test Data")

	opts := &rsa.PSSOptions{SaltLength:rsa.PSSSaltLengthEqualsHash, Hash:crypto.SHA256}

	hash := sha256.New()
	hash.Write(rawData)
	digest := hash.Sum(nil)

	sig, err := cryp.Sign(privateKey, digest, opts)
	assert.NoError(t, err)

	// Test RSA Signer
	_, err = cryp.Sign(nil, digest, opts)
	assert.Error(t, err)

	_, err = cryp.Sign(publicKey, digest, opts)
	assert.Error(t, err)

	_, err = cryp.Sign(nil, digest, opts)
	assert.Error(t, err)

	_, err = cryp.Sign(privateKey, rawData, opts)
	assert.Error(t, err)

	_, err = cryp.Sign(privateKey, nil, opts)
	assert.Error(t, err)

	// Test RSA Verifier
	valid, err := cryp.Verify(publicKey, sig, digest, opts)
	assert.True(t, valid)

	_, err = cryp.Verify(nil, sig, digest, opts)
	assert.Error(t, err)

	_, err = cryp.Verify(privateKey, sig, digest, opts)
	assert.Error(t, err)

	_, err = cryp.Verify(nil, sig, digest, opts)
	assert.Error(t, err)

	_, err = cryp.Verify(privateKey, nil, digest, opts)
	assert.Error(t, err)

	_, err = cryp.Verify(privateKey, sig, rawData, opts)
	assert.Error(t, err)

}