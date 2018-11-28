package blockchain

import (
	"encoding/json"
	"fmt"
	"reflect"
)

type Merklebaum struct {
	Hash  string
	left  *Merklebaum
	right *Merklebaum
	elem  *Transaction
}

func (self Merklebaum) Less(other Merklebaum) bool {
	return self.Hash < other.Hash
}
func (self Merklebaum) ComputeHash() string {
	return fmt.Sprintf("%X", self.ComputeHashByte())
}
func (self Merklebaum) ComputeHashByte() []byte {
	if self.elem != nil {
		return (*self.elem).ComputeHashByte()
	} else {

		// Prepend 0 to prevent second preimage attack
		input := "0"
		if self.left != nil {
			input += self.left.Hash
		}

		if self.right != nil {
			input += self.right.Hash
		}

		return ComputeSha256(input)
	}
}

func CreateMerklebaum(content []Transaction) Merklebaum {
	if len(content) == 0 {
		return Merklebaum{}
	}

	var leafs = make([]*Merklebaum, 0, len(content))
	for i := 0; i < len(content); i++ {
		leafs = append(leafs, &Merklebaum{Hash: content[i].ComputeHash(), elem: &content[i]})
	}

	var bäume = leafs
	for len(bäume) > 1 {
		bäume = createMerkleLevel(bäume)
	}

	return *bäume[0]
}
func createMerkleLevel(bäume []*Merklebaum) []*Merklebaum {
	var result = make([]*Merklebaum, 0, len(bäume)/2+1)

	// In case of odd number of trees, skip last tree for later
	for i := 0; i+1 < len(bäume); i += 2 {
		neuerbaum := Merklebaum{left: bäume[i], right: bäume[i+1]}
		neuerbaum.Hash = neuerbaum.ComputeHash()
		result = append(result, &neuerbaum)
	}

	if len(bäume)%2 == 1 {
		neuerbaum := &Merklebaum{left: bäume[len(bäume)-1]}
		neuerbaum.Hash = neuerbaum.ComputeHash()
		result = append(result, neuerbaum)
	}

	return result
}

func (self Merklebaum) IsValid() bool {
	// Nur Leafs dürfen Transaktionen halten
	if self.elem != nil && !self.IsLeaf() {
		return false
	}
	// Alle Leafs müssen eine Transaktion haben
	if self.IsLeaf() && self.elem == nil {
		return false
	}

	leftValid := self.left == nil || self.left.IsValid()
	rightValid := self.right == nil || self.right.IsValid()
	return leftValid && rightValid && self.Hash != ""
}

func (self Merklebaum) IsLeaf() bool {
	return self.left == nil && self.right == nil
}

func (self Merklebaum) GetElements() []Transaction {
	return self.collectElements([]Transaction{})
}
func (self Merklebaum) collectElements(collectedSoFar []Transaction) []Transaction {
	if self.IsLeaf() {
		if self.elem != nil {
			return append(collectedSoFar, *self.elem)
		} else {
			return collectedSoFar
		}
	} else {
		if self.left != nil {
			collectedSoFar = self.left.collectElements(collectedSoFar)
		}

		if self.right != nil {
			collectedSoFar = self.right.collectElements(collectedSoFar)
		}

		return collectedSoFar
	}
}

func (self Merklebaum) HasNode(path []string) bool {
	if len(path) == 0 {
		return false
	}

	var current = &self
	for i := 0; i < len(path)-1; i++ {
		if current.Hash != path[i] {
			return false
		}
		if current.left.Hash == path[i+1] {
			current = current.left
		}
		if current.right.Hash == path[i+1] {
			current = current.right
		}
	}

	return current.Hash == path[len(path)-1]
}

func (self Merklebaum) CreateSpvProof(elem Transaction) ([]string, bool) {
	proof, found := self.findPath(elem.ComputeHash(), []string{})
	if found {
		return proof, true
	} else {
		return []string{}, true
	}
}
func (self Merklebaum) findPath(elemhash string, path []string) ([]string, bool) {
	if self.IsLeaf() {
		return append(path, self.Hash), (*self.elem).ComputeHash() == elemhash
	} else {
		path = append(path, self.Hash)
		if self.left != nil {
			newpath, found := self.left.findPath(elemhash, path)
			if found {
				return newpath, true
			}
		}

		if self.right != nil {
			newpath, found := self.right.findPath(elemhash, path)
			if found {
				return newpath, true
			}
		}

		return path, false
	}
}

func (self Merklebaum) Contains(leaf Transaction) bool {
	if self.IsLeaf() {
		return reflect.DeepEqual(*self.elem, leaf)
	} else {
		return self.left.Contains(leaf) || self.right.Contains(leaf)
	}
}

func (self Merklebaum) toJson() ([]byte, error) {
	elements := self.GetElements()
	return json.Marshal(elements)
}
