package main

import (
	"awesomeProject/blockchain"
	"crypto/sha256"
	"github.com/cbergoon/merkletree"
	"github.com/stretchr/testify/assert"
	"testing"
)

type TestContent struct {
	x string
}

//CalculateHash hashes the values of a TestContent
func (t TestContent) CalculateHash() ([]byte, error) {
	h := sha256.New()
	if _, err := h.Write([]byte(t.x)); err != nil {
		return nil, err
	}

	return h.Sum(nil), nil
}

//Equals tests for equality of two Contents
func (t TestContent) Equals(other merkletree.Content) (bool, error) {
	return t.x == other.(TestContent).x, nil
}

type TestTransaction struct {
	bla     string
	Outputs []blockchain.Txoutput
	Inputs  []blockchain.Txinput
}

//CalculateHash hashes the values of a TestContent
func (t TestTransaction) CalculateHash() ([]byte, error) {
	h := sha256.New()
	if _, err := h.Write([]byte(t.bla)); err != nil {
		return nil, err
	}

	return h.Sum(nil), nil
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

	isInTree, err = tree.VerifyContent(TestContent{x: "Nihao"})
	assert.NoError(t, err)
	assert.False(t, isInTree)
}
