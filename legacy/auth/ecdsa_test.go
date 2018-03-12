package auth

import (
	"testing"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"github.com/stretchr/testify/assert"
)

func TestMarshalECDSASignature(t *testing.T) {

	generatedKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	assert.NoError(t, err)

	privateKey := &ecdsaPrivateKey{generatedKey}

	rawData := []byte("Now I'm ECDSA Signature Testing!")

	r, s, err := ecdsa.Sign(rand.Reader, privateKey.priv, rawData)
	assert.NoError(t, err)

	// Marshal
	sig, err := marshalECDSASignature(r, s)
	assert.NoError(t, err)
	assert.NotNil(t, sig)

	_, err = marshalECDSASignature(nil, nil)
	assert.Error(t, err)

	// UnMarshal
	r, s, err = unmarshalECDSASignature(sig)
	assert.NoError(t, err)
	assert.NotNil(t, r)
	assert.NotNil(t, s)

	_, _, err = unmarshalECDSASignature(nil)
	assert.Error(t, err)

}

func TestECDSASigner_Sign(t *testing.T) {

	signer := &ecdsaSigner{}
	verifier := &ecdsaVerifier{}

	// Generate Keys
	generatedKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	assert.NoError(t, err)

	privateKey := &ecdsaPrivateKey{generatedKey}
	publicKey, err := privateKey.PublicKey()
	assert.NoError(t, err)

	rawData := []byte("ECDSASigner Test Data")

	// Sign
	sig, err := signer.Sign(privateKey, rawData, nil)
	assert.NoError(t, err)
	assert.NotNil(t, sig)

	// Verify
	valid, err := verifier.Verify(publicKey, sig, rawData, nil)
	assert.NoError(t, err)
	assert.True(t, valid)

}