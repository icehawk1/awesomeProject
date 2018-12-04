package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCanVerifyValid(t *testing.T)  {
	key := CreateKeypair()
	output := output{44, key.PublicKey}
	input := input{from: &output}

	SignInput(&input,key)
	assert.NotNil(t,input.signature.hash)
	if !CheckOutput(output,input.signature) {
		t.Errorf("Illegal output")
	}
}
