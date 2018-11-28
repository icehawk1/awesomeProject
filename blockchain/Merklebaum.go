package blockchain

import (
	"awesomeProject/util"
	"encoding/json"
	"fmt"
	"reflect"
)

type Merklebaum struct {
	Hash  string
	Left  *Merklebaum
	Right *Merklebaum
	Elem  *Transaction
}

func (self Merklebaum) Less(other Merklebaum) bool {
	return self.Hash < other.Hash
}
func (self Merklebaum) ComputeHash() string {
	return fmt.Sprintf("%X", self.ComputeHashByte())
}
func (self Merklebaum) ComputeHashByte() []byte {
	if self.Elem != nil {
		return (*self.Elem).ComputeHashByte()
	} else {

		// Prepend 0 to prevent second preimage attack
		input := "0"
		if self.Left != nil {
			input += self.Left.Hash
		}

		if self.Right != nil {
			input += self.Right.Hash
		}

		return util.ComputeSha256(input)
	}
}

func CreateMerklebaum(content []Transaction) Merklebaum {
	if len(content) == 0 {
		return Merklebaum{}
	}

	var leafs = make([]*Merklebaum, 0, len(content))
	for i := 0; i < len(content); i++ {
		leafs = append(leafs, &Merklebaum{Hash: content[i].ComputeHash(), Elem: &content[i]})
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
		neuerbaum := Merklebaum{Left: bäume[i], Right: bäume[i+1]}
		neuerbaum.Hash = neuerbaum.ComputeHash()
		result = append(result, &neuerbaum)
	}

	if len(bäume)%2 == 1 {
		neuerbaum := &Merklebaum{Left: bäume[len(bäume)-1]}
		neuerbaum.Hash = neuerbaum.ComputeHash()
		result = append(result, neuerbaum)
	}

	return result
}

func (self Merklebaum) IsValid() bool {
	// Nur Leafs dürfen Transaktionen halten
	if self.Elem != nil && !self.IsLeaf() {
		return false
	}
	// Alle Leafs müssen eine Transaktion haben
	if self.IsLeaf() && self.Elem == nil {
		return false
	}

	leftValid := self.Left == nil || self.Left.IsValid()
	rightValid := self.Right == nil || self.Right.IsValid()
	return leftValid && rightValid && self.Hash != ""
}

func (self Merklebaum) IsLeaf() bool {
	return self.Left == nil && self.Right == nil
}

func (self Merklebaum) GetElements() []Transaction {
	return self.collectElements([]Transaction{})
}
func (self Merklebaum) collectElements(collectedSoFar []Transaction) []Transaction {
	if self.IsLeaf() {
		if self.Elem != nil {
			return append(collectedSoFar, *self.Elem)
		} else {
			return collectedSoFar
		}
	} else {
		if self.Left != nil {
			collectedSoFar = self.Left.collectElements(collectedSoFar)
		}

		if self.Right != nil {
			collectedSoFar = self.Right.collectElements(collectedSoFar)
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
		if current.Left != nil && current.Left.Hash == path[i+1] {
			current = current.Left
		} else if current.Right != nil && current.Right.Hash == path[i+1] {
			current = current.Right
		} else {
			return false
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
		return append(path, self.Hash), (*self.Elem).ComputeHash() == elemhash
	} else {
		path = append(path, self.Hash)
		if self.Left != nil {
			newpath, found := self.Left.findPath(elemhash, path)
			if found {
				return newpath, true
			}
		}

		if self.Right != nil {
			newpath, found := self.Right.findPath(elemhash, path)
			if found {
				return newpath, true
			}
		}

		return path, false
	}
}

func (self Merklebaum) Contains(leaf Transaction) bool {
	if self.IsLeaf() {
		return reflect.DeepEqual(*self.Elem, leaf)
	} else {
		return self.Left.Contains(leaf) || self.Right.Contains(leaf)
	}
}

func (self Merklebaum) toJson() ([]byte, error) {
	elements := self.GetElements()
	return json.Marshal(elements)
}
