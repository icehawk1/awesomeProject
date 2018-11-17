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
	Inputs  []txinput
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
	result.Hash = result.ComputeHash()
	return result
}

func (self *Blockchain) Mine() {
	blockheight, oldhead := self.ComputeBlockHeight()
	newhead := Block{Payload: fmt.Sprintf("Block%d", blockheight+1), Prev: oldhead}
	newhead.Hash = newhead.ComputeHash()

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

func (self *Block) ComputeHash() string {
	return fmt.Sprintf("%X",self.ComputeHashByte())
}

func (self *Block) ComputeHashByte() []byte {
	input := fmt.Sprintf("block%s%d", self.Payload, self.nonce)

	if self != nil {
		if self.Prev != nil {
			return computeSha256(input + self.Prev.Hash)
		} else {
			return computeSha256(input)
		}
	} else {
		return computeSha256("")
	}
}

func (self *transaction) ComputeHash() string {return fmt.Sprintf("%X",self.ComputeHashByte())}
func (self *transaction) ComputeHashByte() []byte {
	if self != nil {
		hashinput := "tx"
		for _,output := range self.Outputs {
			hashinput += output.ComputeHash()
		}
		for _,input := range self.Inputs {
			hashinput += input.ComputeHash()
		}

		return computeSha256(hashinput)
	} else {
		return computeSha256("")
	}
}

func (self *txinput) ComputeHash() string {return fmt.Sprintf("%X",self.ComputeHashByte())}
func (self *txinput) ComputeHashByte() []byte {
	if self != nil {
		return computeSha256(fmt.Sprintf("input%X", self.sig.hash))
	} else {
		return computeSha256("")
	}
}

func (self *txoutput) ComputeHash() string {return fmt.Sprintf("%X",self.ComputeHashByte())}
func (self *txoutput) ComputeHashByte() []byte {
	if self != nil {
		return computeSha256(fmt.Sprintf("output%d%s", self.value, self.pubkey))
	} else {
		return computeSha256("")
	}
}

func (self Blockchain) String() string {
	return fmt.Sprintf("Chain Genesis: %s", self.Genesis)
}
func (self Block) String() string {
	return fmt.Sprintf("Block(Hash='%s', msg='%s', Genesis=%t, head=%t)",
		self.Hash, self.Payload, self.Prev == nil, self.Next == nil)
}
