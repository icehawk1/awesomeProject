package blockchain

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
	ComputeHashByte() string
}

type MerkleTree struct {
	left  *MerkleTree
	right *MerkleTree
	value *Transaction
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

	// Prevent second preimage attack
	if self.IsLeaf() {
		return computeSha256Hex(fmt.Sprintf("00%s%s", leftHash, rightHash))
	} else {
		return computeSha256Hex(fmt.Sprintf("01%s%s", leftHash, rightHash))
	}
}

func (self MerkleTree) IsLeaf() bool {
	return self.left == nil && self.right == nil
}

func CreateMerkleTree(txlist []Transaction) *MerkleTree {
	var result []*MerkleTree
	for _,tx := range txlist {
		result = append(result, &MerkleTree{value:&tx})
	}

	for len(result)>1 {
		var nextresult = make([]*MerkleTree, 0, len(result)+1)
		for i:=0; i<len(result); i+=2 {
			if (i+1)<len(result) {
				nextresult = append(nextresult, &MerkleTree{left: result[i], right: result[i+1]})
			} else {
				nextresult = append(nextresult, &MerkleTree{left: result[i], right: nil})
			}
		}
		result = nextresult
	}
	return result[0]
}

func Max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func Min(a, b int) int {
	if a < b {
		return a
	}
	return b
}