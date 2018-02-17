package auth

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"crypto/rsa"
	"crypto/sha256"
	"crypto"
	"crypto/elliptic"
	"os"
)

func TestNew(t *testing.T) {

	// Generate Collector
	_, err := NewCryptoImpl("", &RSAKeyGenOpts{})
	assert.NoError(t, err)

	defer os.RemoveAll("./KeyRepository")

}

//todo fix testcase or hash part
func TestCollector_RSASign(t *testing.T) {

	cryp, err := NewCryptoImpl("", &RSAKeyGenOpts{})
	assert.NoError(t, err)

	privateKey, publicKey, err := cryp.GetKey()
	assert.NoError(t, err)

	defer os.RemoveAll("./KeyRepository")

	rawData := []byte("RSASign Test Data")

	opts := &rsa.PSSOptions{SaltLength:rsa.PSSSaltLengthEqualsHash, Hash:crypto.SHA256}

	hash := sha256.New()
	hash.Write(rawData)
	digest := hash.Sum(nil)

	sig, err := cryp.Sign(rawData, opts)
	assert.NoError(t, err)
	assert.NotNil(t, sig)

	// Test RSA Signer
	_, err = cryp.Sign(rawData, opts)
	assert.NoError(t, err)

	_, err = cryp.Sign(nil, opts)
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

	cryp, err := NewCryptoImpl("", &ECDSAKeyGenOpts{})
	assert.NoError(t, err)

	privateKey, publicKey, err := cryp.GetKey()
	assert.NoError(t, err)

	defer os.RemoveAll("./KeyRepository")

	rawData := []byte("ECDSA Sign Test")

	hash := sha256.New()
	hash.Write(rawData)
	digest := hash.Sum(nil)

	sig, err := cryp.Sign(rawData, nil)
	assert.NoError(t, err)
	assert.NotNil(t, sig)

	// Test RSA Signer
	_, err = cryp.Sign(rawData, nil)
	assert.NoError(t, err)

	_, err = cryp.Sign(nil, nil)
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

func TestCollector_RSAGenerateKey(t *testing.T) {

	cryp, err := NewCryptoImpl("./RSAKeyGen_Test", &RSAKeyGenOpts{})
	assert.NoError(t, err)

	defer os.RemoveAll("./RSAKeyGen_Test")

	pri, pub, err := cryp.GetKey()
	assert.NoError(t, err)
	assert.NotNil(t, pri)
	assert.NotNil(t, pub)

	rsaPriKey, valid := pri.(*rsaPrivateKey)
	assert.True(t, valid)
	assert.NotNil(t, rsaPriKey)

	rsaPubKey, valid := pub.(*rsaPublicKey)
	assert.True(t, valid)
	assert.NotNil(t, rsaPubKey)

	assert.Equal(t, 2048, rsaPriKey.priv.N.BitLen())

}

func TestCollector_ECDSAGenerateKey(t *testing.T) {

	cryp, err := NewCryptoImpl("./ECDSAKeyGen_Test", &ECDSAKeyGenOpts{})
	assert.NoError(t, err)

	defer os.RemoveAll("./ECDSAKeyGen_Test")

	pri, pub, err := cryp.GetKey()
	assert.NoError(t, err)
	assert.NotNil(t, pri)
	assert.NotNil(t, pub)

	ecdsaPriKey, valid := pri.(*ecdsaPrivateKey)
	assert.True(t, valid)
	assert.NotNil(t, ecdsaPriKey)

	ecdsaPubKey, valid := pub.(*ecdsaPublicKey)
	assert.True(t, valid)
	assert.NotNil(t, ecdsaPubKey)

	assert.Equal(t, elliptic.P521(), ecdsaPriKey.priv.Curve)

}

func TestCryptoHelper_LoadKey(t *testing.T) {

	cryp, err := NewCryptoImpl("", &RSAKeyGenOpts{})
	assert.NoError(t, err)

	defer os.RemoveAll("./KeyRepository")

	pri, pub, err := cryp.GetKey()
	assert.NoError(t, err)
	assert.NotNil(t, pri)
	assert.NotNil(t, pub)

}