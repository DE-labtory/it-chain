package auth

import (
	"errors"
	"path/filepath"
	"io/ioutil"
	"encoding/hex"
	"os"
)

type keyManager struct {
	path string
}

func (ks *keyManager) Store(key Key) (err error) {

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

func (ks *keyManager) storePublicKey(key Key) (err error) {

	data, err := PublicKeyToPEM(key)
	if err != nil {
		return
	}

	path, err := ks.getFullPath(hex.EncodeToString(key.SKI()), "pub")
	if err != nil {
		return
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		err = ioutil.WriteFile(path, data, 0700)
		if err != nil {
			return err
		}
	}

	return nil
}

func (ks *keyManager) storePrivateKey(key Key) (err error) {

	data, err := PrivateKeyToPEM(key)
	if err != nil {
		return
	}

	path, err := ks.getFullPath(hex.EncodeToString(key.SKI()), "pri")
	if err != nil {
		return
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		err = ioutil.WriteFile(path, data, 0700)
		if err != nil {
			return err
		}
	}

	return nil
}

func (ks *keyManager) Load(alias string) (key interface{}, err error) {
	return nil, nil
}

func (ks *keyManager) getFullPath(alias, suffix string) (path string, err error) {
	if _, err := os.Stat(ks.path); os.IsNotExist(err) {
		err = os.MkdirAll(ks.path, 0755)
		if err != nil {
			return "", err
		}
	}

	return filepath.Join(ks.path, alias + "_" + suffix), nil
}


