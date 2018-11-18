package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)


func TestCanVerifyValid(t *testing.T)  {
	key := CreateKeypair()
	output := txoutput{44, key.PublicKey}
	input := txinput{from: &output}

	SignInput(&input,key)
	assert.NotNil(t,input.sig.hash)
	if !CheckInputSignature(input) {
		t.Errorf("Illegal txoutput")
	}
}

func TestCanDetectChangeInOutput(t *testing.T)  {
	key := CreateKeypair()
	output := txoutput{44, key.PublicKey}
	input := txinput{from: &output}

	SignInput(&input,key)

	input.from.value = 44000

	if CheckInputSignature(input) {
		t.Errorf("change of txoutput should have been detected")
	}

	input.from = &txoutput{44000, key.PublicKey}

	if CheckInputSignature(input) {
		t.Errorf("change of value should have been detected")
	}

}