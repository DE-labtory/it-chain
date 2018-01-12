package auth

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"crypto/rsa"
	"crypto/rand"
	"crypto/sha256"
	"crypto"
	"crypto/ecdsa"
	"crypto/elliptic"
)

func TestNew(t *testing.T) {

	// Generate Collector
	_, err := NewCollector()
	assert.NoError(t, err)

}

func TestCollector_RSASign(t *testing.T) {

	cryp, err := NewCollector()
	assert.NoError(t, err)

	// Generate an RSA Key
	generatedKey, err := rsa.GenerateKey(rand.Reader, 1024)
	assert.NoError(t, err)

	// Set private key
	privateKey := &rsaPrivateKey{generatedKey}

	// Get public key
	publicKey, err := privateKey.PublicKey()

	rawData := []byte("RSASign Test Data")

	opts := &rsa.PSSOptions{SaltLength:rsa.PSSSaltLengthEqualsHash, Hash:crypto.SHA256}

	hash := sha256.New()
	hash.Write(rawData)
	digest := hash.Sum(nil)

	sig, err := cryp.Sign(privateKey, digest, opts)
	assert.NoError(t, err)
	assert.NotNil(t, sig)

	// Test RSA Signer
	_, err = cryp.Sign(nil, digest, opts)
	assert.Error(t, err)

	_, err = cryp.Sign(publicKey, digest, opts)
	assert.Error(t, err)

	_, err = cryp.Sign(privateKey, rawData, opts)
	assert.Error(t, err)

	_, err = cryp.Sign(privateKey, nil, opts)
	assert.Error(t, err)

	// Test RSA Verifier
	valid, err := cryp.Verify(publicKey, sig, digest, opts)
	assert.NoError(t, err)
	assert.True(t, valid)

	_, err = cryp.Verify(nil, sig, digest, opts)
	assert.Error(t, err)

	_, err = cryp.Verify(privateKey, sig, digest, opts)
	assert.Error(t, err)

	_, err = cryp.Verify(publicKey, nil, digest, opts)
	assert.Error(t, err)

	_, err = cryp.Verify(publicKey, sig, rawData, opts)
	assert.Error(t, err)

}

func TestCollector_ECDSASign(t *testing.T) {

	cryp, err := NewCollector()
	assert.NoError(t, err)

	generatedKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	assert.NoError(t, err)

	privateKey := &ecdsaPrivateKey{generatedKey}
	publicKey, err := privateKey.PublicKey()
	assert.NoError(t, err)

	rawData := []byte("ECDSA Sign Test")

	hash := sha256.New()
	hash.Write(rawData)
	digest := hash.Sum(nil)

	sig, err := cryp.Sign(privateKey, digest, nil)
	assert.NoError(t, err)
	assert.NotNil(t, sig)

	// Test RSA Signer
	_, err = cryp.Sign(nil, digest, nil)
	assert.Error(t, err)

	_, err = cryp.Sign(publicKey, digest, nil)
	assert.Error(t, err)

	_, err = cryp.Sign(privateKey, nil, nil)
	assert.Error(t, err)

	// Test RSA Verifier
	valid, err := cryp.Verify(publicKey, sig, digest, nil)
	assert.NoError(t, err)
	assert.True(t, valid)

	_, err = cryp.Verify(nil, sig, digest, nil)
	assert.Error(t, err)

	_, err = cryp.Verify(privateKey, sig, digest, nil)
	assert.Error(t, err)

	_, err = cryp.Verify(publicKey, nil, digest, nil)
	assert.Error(t, err)


}

func TestCollector_KeyGenerate(t *testing.T) {

	

}