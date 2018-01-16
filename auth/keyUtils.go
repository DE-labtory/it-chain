package auth

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"crypto/ecdsa"
	"errors"
)

func PublicKeyToPem(pub Key) (data []byte, err error) {

	if pub == nil {
		return nil, errors.New("Invalid Key")
	}

	switch k := pub.(type) {
	case *rsa.PublicKey:

		keyData, err := x509.MarshalPKIXPublicKey(k)
		if err != nil {
			return nil, err
		}

		return pem.EncodeToMemory(
			&pem.Block{
				Type: "RSA PUBLIC KEY",
				Bytes: keyData,
			},
		), nil

	case *ecdsa.PublicKey:

		keyData, err := x509.MarshalPKIXPublicKey(k)
		if err != nil {
			return nil, err
		}

		return pem.EncodeToMemory(
			&pem.Block{
				Type: "RSA PUBLIC KEY",
				Bytes: keyData,
			},
		), nil

	default:
		return nil, errors.New("Unspported Public Key Type")
	}

}

func PrivateKeyToPem(pri Key) (data []byte, err error) {

	if pri == nil {
		return nil, errors.New("Invalid Private Key")
	}

	switch k := pri.(type) {
	case *rsa.PrivateKey:
		keyData := x509.MarshalPKCS1PrivateKey(k)

		return pem.EncodeToMemory(
			&pem.Block{
				Type: "RSA PRIVATE KEY",
				Bytes: keyData,
			},
		), nil

	case *ecdsa.PrivateKey:
		keyData, err := x509.MarshalECPrivateKey(k)
		if err != nil {
			return nil, err
		}

		return pem.EncodeToMemory(
			&pem.Block{
				Type: "ECDSA PRIVATE KEY",
				Bytes: keyData,
			},
		), nil

	default:
		return nil, errors.New("Unspported Private Key Type")
	}
}