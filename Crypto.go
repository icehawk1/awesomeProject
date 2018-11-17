package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"math/big"
	"strconv"
)

type Signature struct {
	r    big.Int
	s    big.Int
	hash []byte
}

func CreateKeypair() ecdsa.PrivateKey {
	result, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		panic(err)
	}
	return *result
}

func SignInput(input *txinput, key ecdsa.PrivateKey) {
	hash := sha256.Sum256([]byte(strconv.Itoa(input.from.value)))
	r, s, err := ecdsa.Sign(rand.Reader, &key, hash[:])
	if err != nil {
		panic(err)
	}
	input.sig = Signature{r: *r, s: *s, hash:hash[:]}
}

func CheckInput(input txinput) bool {
	valid := ecdsa.Verify(&input.from.pubkey, input.sig.hash, &input.sig.r, &input.sig.s)
	return valid
}
