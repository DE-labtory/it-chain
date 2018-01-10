package auth

import (
	"testing"
	"crypto/rsa"
	"crypto/rand"
	"github.com/stretchr/testify/assert"
	"crypto"
	"crypto/sha256"
)

func TestRsaSigner_Sign(t *testing.T) {

	signer := &rsaSigner{}
	verifier := &rsaVerifier{}

	// Generate keys
	generatedKey, err := rsa.GenerateKey(rand.Reader, 1024)
	assert.NoError(t, err)
	privateKey := &rsaPrivateKey{generatedKey}
	publicKey, err := privateKey.PublicKey()


	rawData := []byte("RSASigner Test Data")

	opts := &rsa.PSSOptions{SaltLength:rsa.PSSSaltLengthEqualsHash, Hash:crypto.SHA256}

	// Sign
	_, err = signer.Sign(privateKey, rawData, opts)
	assert.Error(t, err)

	_, err = signer.Sign(privateKey, rawData, nil)
	assert.Error(t, err)

	hash := sha256.New()
	hash.Write(rawData)
	digest := hash.Sum(nil)
	sig, err := signer.Sign(privateKey, digest, opts)
	assert.NoError(t, err)

	err = rsa.VerifyPSS(&privateKey.priv.PublicKey, crypto.SHA256, digest, sig, opts)
	assert.NoError(t, err)

	// Verify
	valid, err := verifier.Verify(publicKey, sig, digest, opts)
	assert.True(t, valid)

	_, err = verifier.Verify(publicKey, sig, digest, nil)
	assert.Error(t, err)

}