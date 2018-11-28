package blockchain

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"math/big"
)

type Signature struct {
	R big.Int
	S big.Int
}

var DefaultCurve = elliptic.P256()

func CreateKeypair() ecdsa.PrivateKey {
	result, err := ecdsa.GenerateKey(DefaultCurve, rand.Reader)
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
	input.Sig = Signature{R: *r, S: *s}
}

func CheckInputSignature(input Txinput) bool {
	x,y := elliptic.Unmarshal(DefaultCurve,input.From.Pubkey)
	valid := ecdsa.Verify(&ecdsa.PublicKey{DefaultCurve,x,y}, input.From.ComputeHashByte(), &input.Sig.R, &input.Sig.S)
	return valid
}
