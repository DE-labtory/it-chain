package auth

import (
	"errors"
	"path/filepath"
	"io/ioutil"
	"encoding/hex"
	"os"
	"strings"
	"crypto/rsa"
	"crypto/ecdsa"
)

type keyManager struct {
	path string
}

type keyType string

const (
	PRIVATE_KEY keyType = "pri"
	PUBLIC_KEY	keyType = "pub"
)

func (km *keyManager) Init(path string) {

	if len(path) == 0 {
		km.path = "./KeyRepository"
	} else {
		if !strings.HasPrefix(path, "./") {
			km.path = "./" + path
		} else {
			km.path = path
		}
	}

}

func (km *keyManager) Store(keys... Key) (err error) {

	if len(keys) == 0 {
		return errors.New("Input values should not be NIL")
	}

	for _, key := range keys {
		switch k := key.(type) {
		case *rsaPrivateKey:
			err = km.storeKey(k, PRIVATE_KEY)
		case *rsaPublicKey:
			err = km.storeKey(k, PUBLIC_KEY)
		case *ecdsaPrivateKey:
			err = km.storeKey(k, PRIVATE_KEY)
		case *ecdsaPublicKey:
			err = km.storeKey(k, PUBLIC_KEY)
		default:
			return errors.New("Unspported Key Type.")
		}
	}

	return nil
}

func (km *keyManager) storeKey(key Key, keyType keyType) (err error) {

	var data []byte

	switch keyType {
	case PRIVATE_KEY:
		data, err = PrivateKeyToPEM(key)
	case PUBLIC_KEY:
		data, err = PublicKeyToPEM(key)
	default:
		return errors.New("Unsupported Key Type")
	}

	if err != nil {
		return
	}

	path, err := km.getFullPath(hex.EncodeToString(key.SKI()), string(keyType))
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

func (km *keyManager) Load() (pri, pub Key, err error) {

	if _, err := os.Stat(km.path); os.IsNotExist(err) {
		return nil, nil, errors.New("Keys are not exist")
	}

	files, err := ioutil.ReadDir(km.path)
	if err != nil {
		return nil, nil, errors.New("Failed to read key repository directory")
	}

	for _, file := range files {

		suffix, valid := km.getSuffix(file.Name())
		if valid == true {
			alias := strings.Split(file.Name(), "_")[0]
			switch suffix {
			case "pri":
				key, err := km.loadKey(alias, PRIVATE_KEY)
				if err != nil {
					return nil, nil, err
				}

				switch key.(type) {
				case *rsa.PrivateKey:
					pri = &rsaPrivateKey{key.(*rsa.PrivateKey)}
				case *ecdsa.PrivateKey:
					pri = &ecdsaPrivateKey{key.(*ecdsa.PrivateKey)}
				default:
					return nil, nil, errors.New("Failed to load Key")
				}

			case "pub":
				key, err := km.loadKey(alias, PUBLIC_KEY)
				if err != nil {
					return nil, nil, err
				}

				switch key.(type) {
				case *rsa.PublicKey:
					pub = &rsaPublicKey{key.(*rsa.PublicKey)}
				case *ecdsa.PublicKey:
					pub = &ecdsaPublicKey{key.(*ecdsa.PublicKey)}
				default:
					return nil, nil, errors.New("Failed to load Key")
				}
			}
		}
	}

	if pri == nil || pub == nil {
		return nil, nil, errors.New("Failed to load Key")
	}

	return

}

func (km *keyManager) loadKey(alias string, keyType keyType) (key interface{}, err error) {

	if len(alias) == 0 {
		return nil, errors.New("Input value should not be blank")
	}

	path, err := km.getFullPath(alias, string(keyType))
	if err != nil {
		return nil, err
	}

	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	switch keyType {
	case PRIVATE_KEY:
		key, err = PEMToPrivateKey(data)
	case PUBLIC_KEY:
		key, err = PEMToPublicKey(data)
	}

	if err != nil {
		return nil, err
	}

	return

}

func (km *keyManager) getSuffix(name string) (suffix string, valid bool) {

	if strings.HasSuffix(name, "pri") {
		return "pri", true
	} else if strings.HasSuffix(name, "pub") {
		return "pub", true
	}

	return "", false

}

func (km *keyManager) getFullPath(alias, suffix string) (path string, err error) {
	if _, err := os.Stat(km.path); os.IsNotExist(err) {
		err = os.MkdirAll(km.path, 0755)
		if err != nil {
			return "", err
		}
	}

	return filepath.Join(km.path, alias + "_" + suffix), nil
}


