package main

import (
	"crypto/sha512"

	"github.com/it-chain/bifrost/pb"
	"github.com/it-chain/heimdall/auth"
	"github.com/it-chain/heimdall/key"
)

type SingingService struct {
	au auth.Auth
}

func (s SingingService) Sign(envelope *pb.Envelope, priKey key.PriKey) *pb.Envelope {

	hash := sha512.New()
	hash.Write(envelope.Payload)
	digest := hash.Sum(nil)

	sig, _ := s.au.Sign(priKey, digest, auth.EQUAL_SHA512.SignerOptsToPSSOptions())
	envelope.Signature = sig

	return envelope
}
