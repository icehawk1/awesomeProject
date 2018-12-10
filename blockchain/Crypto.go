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

func PubkeyEqual(key1 ecdsa.PublicKey, key2 ecdsa.PublicKey) bool {
	x := key1.X.Cmp(key2.X)==0
	y := key1.Y.Cmp(key2.Y)==0
	return x && y
}
func PrivkeyEqual(key1 ecdsa.PrivateKey, key2 ecdsa.PrivateKey) bool {
	return key1.D==key2.D
}

func MarshalPubkey(pubkey ecdsa.PublicKey) []byte {
	return elliptic.Marshal(DefaultCurve, pubkey.X,pubkey.Y)
}
func UnmarshalPubkey(pubkey []byte) ecdsa.PublicKey {
	x,y := elliptic.Unmarshal(DefaultCurve, pubkey)
	return ecdsa.PublicKey{DefaultCurve,x,y}
}

func SignInput(input Txinput, key ecdsa.PrivateKey) Signature {
	//hash := sha256.Sum256([]byte(strconv.Itoa(input.From.Value)))
	hash := input.From.ComputeHashByte()
	r, s, err := ecdsa.Sign(rand.Reader, &key, hash[:])
	if err != nil {
		panic(err)
	}
	return Signature{R: *r, S: *s}
}

func CheckInputSignature(input Txinput) bool {
	pubkey := UnmarshalPubkey(input.From.Pubkey)
	valid := ecdsa.Verify(&pubkey, input.From.ComputeHashByte(), &input.Sig.R, &input.Sig.S)
	return valid
}
