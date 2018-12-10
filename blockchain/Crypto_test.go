package blockchain

import (
	"github.com/stretchr/testify/assert"
	"testing"
)


func TestPubkeyEqual_marshalled(t *testing.T) {
	key1 := CreateKeypair().PublicKey
	key1Marshalled := UnmarshalPubkey(MarshalPubkey(key1))
	assert.True(t, PubkeyEqual(key1, key1Marshalled))
}

func TestPubkeyEqual_nonequal(t *testing.T) {
	key1 := CreateKeypair().PublicKey
	key2 := CreateKeypair().PublicKey
	assert.False(t, PubkeyEqual(key1, key2))
}

func TestCanVerifyValid(t *testing.T)  {
	key := CreateKeypair()
	output := Txoutput{44, MarshalPubkey(key.PublicKey)}
	input := Txinput{From: &output}

	input.Sig = SignInput(input,key)
	if !CheckInputSignature(input) {
		t.Errorf("Illegal Txoutput")
	}
}

func TestCanDetectChangeInOutput(t *testing.T)  {
	key := CreateKeypair()
	output := Txoutput{44, MarshalPubkey(key.PublicKey)}
	input := Txinput{From: &output}

	input.Sig = SignInput(input,key)

	input.From.Value = 44000

	if CheckInputSignature(input) {
		t.Errorf("change of Txoutput should have been detected")
	}

	input.From = &Txoutput{44000, MarshalPubkey(key.PublicKey)}

	if CheckInputSignature(input) {
		t.Errorf("change of Value should have been detected")
	}
}

