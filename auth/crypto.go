package auth

type Key interface {

	SKI() (ski []byte)

}

type KeyManager interface {

	Store(key Key) (err error)

	Load(alias string) (key interface{}, err error)

}

type Crypto interface {

	Sign(key Key, digest []byte, opts SignerOpts) (signature []byte, err error)

	Verify(key Key, signature, digest []byte, opts SignerOpts) (valid bool, err error)

	KeyGenerate(opts KeyGenOpts) (key Key, err error)

	KeyImport(data interface{}, opts KeyImporterOpts) (key Key, err error)

}