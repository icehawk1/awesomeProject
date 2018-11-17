package main

import (
	"crypto/sha256"
	"fmt"
)

func computeSha256Hex(input string) string {
	return fmt.Sprintf("%X", computeSha256(input))
}

func computeSha256(input string) []byte {
	hash := sha256.New()
	hash.Write([]byte(input))
	return hash.Sum(nil)
}

type Hashable interface {
	ComputeHash() string
}

type MerkleTree struct {
	left  *MerkleTree
	right *MerkleTree
	value *transaction
}

func (self MerkleTree) ComputeHash() string {
	var leftHash, rightHash string
	if self.left != nil {
		leftHash = self.left.ComputeHash()
	} else {
		leftHash = computeSha256Hex("")
	}

	if self.right != nil {
		rightHash = self.right.ComputeHash()
	} else {
		rightHash = computeSha256Hex("")
	}

	return computeSha256Hex(fmt.Sprintf("%s%s", leftHash, rightHash))
}

func (self MerkleTree) IsLeaf() bool {
	return self.left == nil && self.right == nil
}

func (self MerkleTree) PutIntoTree(tx *transaction) {
	if self.left == nil {
		self.left = &MerkleTree{value: tx}
	} else if self.right == nil {
		self.right = &MerkleTree{value: tx}
	} else {
		// TODO: Baum ordentlich balancieren
		self.left.PutIntoTree(tx)
	}
}
