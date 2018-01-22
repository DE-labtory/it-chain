package auth

import (
	"crypto/x509"
	"encoding/pem"
	"errors"
	"crypto/rsa"
	"crypto/ecdsa"
)

func PublicKeyToPEM(pub Key) (data []byte, err error) {

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

func PrivateKeyToPEM(pri Key) (data []byte, err error) {

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

func PEMToPublicKey(data []byte) (key interface{}, err error) {
	if len(data) == 0 {
		return nil, errors.New("Input data should not be NIL")
	}

	block, _ := pem.Decode(data)
	if block == nil {
		return nil, errors.New("Failed to decode data")
	}

	key, err = DERToPublicKey(block.Bytes)
	if err != nil {
		return nil, errors.New("Failed to convert PEM data to public key")
	}

	return
}

func PEMToPrivateKey(data []byte) (key interface{}, err error) {
	if len(data) == 0 {
		return nil, errors.New("Input data should not be NIL")
	}

	block, _ := pem.Decode(data)
	if block == nil {
		return nil, errors.New("Failed to decode data")
	}

	key, err = DERToPrivateKey(block.Bytes)
	if err != nil {
		return nil, errors.New("Failed to convert PEM data to private key")
	}

	return

}

func DERToPublicKey(data []byte) (key interface{}, err error) {
	if len(data) == 0 {
		return nil, errors.New("Input data should not be NIL")
	}

	key, err = x509.ParsePKIXPublicKey(data)
	if err != nil {
		return nil, errors.New("Failed to Parse data")
	}

	return
}

func DERToPrivateKey(data []byte) (key interface{}, err error) {
	if len(data) == 0 {
		return nil, errors.New("Input data should not be NIL")
	}

	if key, err = x509.ParsePKCS1PrivateKey(data); err == nil {
		return key, err
	}

	if key, err = x509.ParseECPrivateKey(data); err == nil {
		return key, err
	}

	return nil, errors.New("Unspported Private Key Type")
}



