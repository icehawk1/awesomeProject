package main

import (
	"awesomeProject/blockchain"
	"github.com/cbergoon/merkletree"
	"github.com/stretchr/testify/assert"
	"testing"
)

type TestTransaction struct {
	bla     string
	Outputs []blockchain.Txoutput
	Inputs  []blockchain.Txinput
}

//CalculateHash hashes the values of a TestContent
func (self TestTransaction) CalculateHash() ([]byte, error) {
	return self.ComputeHashByte(), nil
}
func (self *TestTransaction) ComputeHashByte() []byte {
	if self != nil {
		hashinput := "tx"
		for _, output := range self.Outputs {
			hashinput += output.ComputeHash()
		}
		for _, input := range self.Inputs {
			hashinput += input.ComputeHash()
		}

		return blockchain.ComputeSha256(hashinput)
	} else {
		return blockchain.ComputeSha256("")
	}
}

//Equals tests for equality of two Contents
func (t TestTransaction) Equals(other merkletree.Content) (bool, error) {
	othertx, ok := other.(TestTransaction)
	if ok {
		return t.bla == othertx.bla, nil
	} else {
		return false, nil
	}
}

func TestGettingStarted(t *testing.T) {
	var list []merkletree.Content
	list = append(list, TestTransaction{bla: "Hello"})
	list = append(list, TestTransaction{bla: "Hi"})
	list = append(list, TestTransaction{bla: "Hey"})
	list = append(list, TestTransaction{bla: "Hola"})

	tree, err := merkletree.NewTree(list)
	assert.NoError(t, err)

	root := tree.MerkleRoot()
	assert.NotNil(t, root)

	validtree, err := tree.VerifyTree()
	assert.NoError(t, err)
	assert.True(t, validtree)

	isInTree, err := tree.VerifyContent(list[2])
	assert.NoError(t, err)
	assert.True(t, isInTree)

	isInTree, err = tree.VerifyContent(TestTransaction{bla: "Nihao"})
	assert.NoError(t, err)
	assert.False(t, isInTree)
}
