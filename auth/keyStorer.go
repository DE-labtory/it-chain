package auth

import (
	"errors"
	"path/filepath"
	"io/ioutil"
	"encoding/hex"
)

type keyStorer struct {
	path string
}

func (ks *keyStorer) Store(key Key) (err error) {

	if key == nil {
		return errors.New("Failed to get Key Data")
	}

	switch k := key.(type) {
	case *rsaPrivateKey:
		err = ks.storePrivateKey(k)
	case *rsaPublicKey:
		err = ks.storePublicKey(k)
	case *ecdsaPrivateKey:
		err = ks.storePrivateKey(k)
	case *ecdsaPublicKey:
		err = ks.storePublicKey(k)
	default:
		return errors.New("Unspported Key Type.")
	}
	
	return
}

func (ks *keyStorer) storePublicKey(key Key) (err error) {

	data, err := PublicKeyToPem(key)
	if err != nil {
		return
	}

	err = ioutil.WriteFile(ks.getFullPath(hex.EncodeToString(key.SKI()), "pub"), data, 0700)
	if err != nil {
		return
	}

	return nil
}

func (ks *keyStorer) storePrivateKey(key Key) (err error) {

	data, err := PrivateKeyToPem(key)
	if err != nil {
		return
	}

	err = ioutil.WriteFile(ks.getFullPath(hex.EncodeToString(key.SKI()), "pri"), data, 0700)
	if err != nil {
		return
	}

	return nil
}

func (ks *keyStorer) getFullPath(alias, suffix string) string {
	return filepath.Join(ks.path, alias + "_" + suffix)
}


