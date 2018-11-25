package merklebaum

import "awesomeProject/blockchain"

type Merklebaum struct {
	Hash string
	left *Merklebaum
	right *Merklebaum
}

func CreateMerklebaum(content []blockchain.Hashable) Merklebaum {
	return Merklebaum{}
}

func (self Merklebaum) IsValid() bool {
	return false
}

func (self Merklebaum) HasNode(path []string) bool {
	return false
}

func (self Merklebaum) GetLeafes() []blockchain.Hashable {
	return []blockchain.Hashable{}
}

func (self Merklebaum) CreateSpvProof(leaf blockchain.Hashable) []string {
	return []string{}
}