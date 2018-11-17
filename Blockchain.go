package main

import (
	"crypto/ecdsa"
	"fmt"
)

type Blockchain struct {
	Genesis Block
}

type Block struct {
	Hash         string
	Payload      string
	Prev         *Block `json:"-"`
	Next         *Block
	nonce        int64
	transactions []transaction
}

type transaction struct {
	Outputs []txoutput
	inputs  []txinput
}

type txinput struct {
	from *txoutput
	sig  Signature
}

type txoutput struct {
	value  int
	pubkey ecdsa.PublicKey
}

func CreateChain(msg string) Blockchain {
	result := Blockchain{Genesis: CreateBlock(msg)}
	return result
}

func CreateBlock(msg string) Block {
	result := Block{Payload: msg}
	result.Hash = result.computeHash()
	return result
}

func (self *Blockchain) Mine() {
	blockheight, oldhead := self.ComputeBlockHeight()
	newhead := Block{Payload: fmt.Sprintf("Block%d", blockheight+1), Prev: oldhead}
	newhead.Hash = newhead.computeHash()

	oldhead.Next = &newhead
}

func (self *Blockchain) ComputeBlockHeight() (int, *Block) {
	var current *Block = &self.Genesis
	i := 0
	for ; current.Next != nil; i++ {
		current = current.Next
	}
	return i, current
}

// ------Private parts-----------------------------------------

func (self *Block) ComputeHash() string {
	input := fmt.Sprintf("%s%d",self.Payload,self.nonce)

	if self.Prev != nil {
		return computeSha256Hex(input + self.Prev.Hash)
	} else {
		return computeSha256Hex(input)
	}
}
func toArray(chain Blockchain) []Block {
	var result []Block
	current := &chain.Genesis
	for ; current != nil; current = current.Next {
		result = append(result, *current)
	}
	return result
}

func (self Blockchain) String() string {
	return fmt.Sprintf("Chain Genesis: %s", self.Genesis)
}
func (self Block) String() string {
	return fmt.Sprintf("Block(Hash='%s', msg='%s', Genesis=%t, head=%t)",
		self.Hash, self.Payload, self.Prev == nil, self.Next == nil)
}
