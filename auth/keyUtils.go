package auth

import (
	"crypto/x509"
	"encoding/pem"
	"errors"
)

func PublicKeyToPem(pub Key) (data []byte, err error) {

	if pub == nil {
		return nil, errors.New("Invalid Key")
	}

	switch k := pub.(type) {
	case *rsaPublicKey:

		keyData, err := x509.MarshalPKIXPublicKey(k.pub)
		if err != nil {
			return nil, err
		}

		return pem.EncodeToMemory(
			&pem.Block{
				Type: "RSA PUBLIC KEY",
				Bytes: keyData,
			},
		), nil

	case *ecdsaPublicKey:

		keyData, err := x509.MarshalPKIXPublicKey(k.pub)
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
	case *rsaPrivateKey:
		keyData := x509.MarshalPKCS1PrivateKey(k.priv)

		return pem.EncodeToMemory(
			&pem.Block{
				Type: "RSA PRIVATE KEY",
				Bytes: keyData,
			},
		), nil

	case *ecdsaPrivateKey:
		keyData, err := x509.MarshalECPrivateKey(k.priv)
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