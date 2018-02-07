package auth

type Key interface {

	SKI() (ski []byte)

}

type Crypto interface {

	Sign(digest []byte, opts SignerOpts) ([]byte, error)

	Verify(key Key, signature, digest []byte, opts SignerOpts) (bool, error)

	GenerateKey(opts KeyGenOpts) (pri, pub Key, err error)

	LoadKey() (pri, pub Key, err error)
}