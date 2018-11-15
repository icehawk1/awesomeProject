package main

import (
	"encoding/json"
	"fmt"
)


type Block struct {
	Hash    string
	Payload string
	Prev    *Block `json:"-"`
	Next    *Block
}

type Blockchain struct {
	Genesis Block
}

func (self *Block) ComputeHash() string {
	if self.Prev != nil {
		return computeSha256Hex(self.Payload +self.Prev.Hash)
	} else {
		return computeSha256Hex(self.Payload)
	}
}

func CreateChain(msg string) Blockchain {
	result := Blockchain{Genesis: Block{Payload: msg}}
	result.Genesis.Hash = result.Genesis.ComputeHash()
	return result
}

func (self *Blockchain) Mine() {
	blockheight, oldhead := self.computeBlockHeight()
	newhead := Block{Payload: fmt.Sprintf("block%d",blockheight+1), Prev:oldhead}
	newhead.Hash = newhead.ComputeHash()

	oldhead.Next =&newhead
}

func (self *Blockchain) computeBlockHeight() (int,*Block) {
	var current *Block = &self.Genesis
	i:=0
	for ; current.Next != nil; i++ {current=current.Next
	}
	return i,current
}

func toArray(chain Blockchain) []Block {
	var result []Block
	current := &chain.Genesis
	for ; current != nil; current=current.Next {
		result = append(result, *current)
	}
	return result
}

func toJson(chain Blockchain) string {
	array := toArray(chain)
	bytes, _ := json.Marshal(array)
	return fmt.Sprintf("%X", bytes)
}

func (self Blockchain) String() string {
	return fmt.Sprintf("Chain Genesis: %s", self.Genesis)
}
func (self Block) String() string {
	return fmt.Sprintf("Block(Hash='%s', msg='%s', Genesis=%t, head=%t)",
		self.Hash,self.Payload,self.Prev ==nil,self.Next ==nil)
}