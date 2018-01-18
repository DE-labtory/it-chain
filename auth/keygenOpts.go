package auth

type RSAKeyGenOpts struct {
	Temporary bool
}

func (opts *RSAKeyGenOpts) Ephemeral() bool {
	return opts.Temporary
}

type ECDSAKeyGenOpts struct {
	Temporary bool
}

func (opts *ECDSAKeyGenOpts) Ephemeral() bool {
	return opts.Temporary
}