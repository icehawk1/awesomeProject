package blockchain

import (
	"github.com/stretchr/testify/assert"
	"testing"
)


func TestCanVerifyValid(t *testing.T)  {
	key := CreateKeypair()
	output := Txoutput{44, key.PublicKey}
	input := Txinput{From: &output}

	SignInput(&input,key)
	assert.NotNil(t,input.Sig.hash)
	if !CheckInputSignature(input) {
		t.Errorf("Illegal Txoutput")
	}
}

func TestCanDetectChangeInOutput(t *testing.T)  {
	key := CreateKeypair()
	output := Txoutput{44, key.PublicKey}
	input := Txinput{From: &output}

	SignInput(&input,key)

	input.From.Value = 44000

	if CheckInputSignature(input) {
		t.Errorf("change of Txoutput should have been detected")
	}

	input.From = &Txoutput{44000, key.PublicKey}

	if CheckInputSignature(input) {
		t.Errorf("change of Value should have been detected")
	}

}