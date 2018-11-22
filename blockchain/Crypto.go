package blockchain

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
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

func SignInput(input *Txinput, key ecdsa.PrivateKey) {
	//hash := sha256.Sum256([]byte(strconv.Itoa(input.From.Value)))
	hash := input.From.ComputeHashByte()
	r, s, err := ecdsa.Sign(rand.Reader, &key, hash[:])
	if err != nil {
		panic(err)
	}
	input.Sig = Signature{r: *r, s: *s, hash:hash[:]}
}

func CheckInputSignature(input Txinput) bool {
	valid := ecdsa.Verify(&input.From.Pubkey, input.From.ComputeHashByte(), &input.Sig.r, &input.Sig.s)
	return valid
}
