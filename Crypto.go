package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"math/big"
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

func SignInput(input *input, key ecdsa.PrivateKey) {
	hash := sha256.Sum256([]byte(fmt.Sprintf("%X", input.from.value)))
	bla := hash[:]
	fmt.Println(bla)

	r, s, err := ecdsa.Sign(rand.Reader, &key, hash[:])
	if err != nil {
		panic(err)
	}
	input.signature = Signature{r: *r, s: *s, hash:hash[:]}
}

func CheckOutput(output output, sig Signature) bool {
	valid := ecdsa.Verify(&output.pubkey, sig.hash, &sig.r, &sig.s)
	return valid
}
