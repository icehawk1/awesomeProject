package main

import (
	"awesomeProject/blockchain"
	"github.com/cbergoon/merkletree"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGettingStarted(t *testing.T) {
	key1 := blockchain.CreateKeypair()

	var list []merkletree.Content
	outputlist := []blockchain.Txoutput{blockchain.CreateTxOutput(0, key1.PublicKey)}
	list = append(list, blockchain.Transaction{Outputs: outputlist})
	outputlist = []blockchain.Txoutput{blockchain.CreateTxOutput(1, key1.PublicKey)}
	list = append(list, blockchain.Transaction{Outputs: outputlist})
	outputlist = []blockchain.Txoutput{blockchain.CreateTxOutput(2, key1.PublicKey)}
	list = append(list, blockchain.Transaction{Outputs: outputlist})
	outputlist = []blockchain.Txoutput{blockchain.CreateTxOutput(3, key1.PublicKey)}
	list = append(list, blockchain.Transaction{Outputs: outputlist})

	tree, err := merkletree.NewTree(list)
	assert.NoError(t, err)

	root := tree.MerkleRoot()
	assert.NotNil(t, root)

	validtree, err := tree.VerifyTree()
	assert.NoError(t, err)
	assert.True(t, validtree)

	equal, err := list[2].Equals(tree.Leafs[2].C)
	assert.NoError(t,err)
	assert.True(t, equal, "Equals function does not work")

	isInTree, err := tree.VerifyContent(list[2])
	assert.NoError(t, err)
	assert.True(t, isInTree, "Third element was not found")

	isInTree, err = tree.VerifyContent(blockchain.Transaction{})
	assert.NoError(t, err)
	assert.False(t, isInTree)
}
