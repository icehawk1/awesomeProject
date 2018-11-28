package blockchain

import (
	"crypto/elliptic"
	"testing"
)


func TestCanVerifyValid(t *testing.T)  {
	key := CreateKeypair()
	output := Txoutput{44, elliptic.Marshal(DefaultCurve, key.PublicKey.X, key.PublicKey.Y)}
	input := Txinput{From: &output}

	SignInput(&input,key)
	if !CheckInputSignature(input) {
		t.Errorf("Illegal Txoutput")
	}
}

func TestCanDetectChangeInOutput(t *testing.T)  {
	key := CreateKeypair()
	output := Txoutput{44, elliptic.Marshal(DefaultCurve, key.PublicKey.X, key.PublicKey.Y)}
	input := Txinput{From: &output}

	SignInput(&input,key)

	input.From.Value = 44000

	if CheckInputSignature(input) {
		t.Errorf("change of Txoutput should have been detected")
	}

	input.From = &Txoutput{44000, elliptic.Marshal(DefaultCurve, key.PublicKey.X, key.PublicKey.Y)}

	if CheckInputSignature(input) {
		t.Errorf("change of Value should have been detected")
	}

}