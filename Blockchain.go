package main

import (
	"crypto/ecdsa"
	"fmt"
	"math"
)

type Blockchain struct {
	Genesis Block
}

type Block struct {
	Hash         string
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

const OUTPUT_MINVALUE = 0
const OUTPUT_MAXVALUE = math.MaxInt32

type txoutput struct {
	value  int
	pubkey ecdsa.PublicKey
}

func CreateChain() Blockchain {
	result := Blockchain{Genesis: CreateBlock(make([]transaction, 0))}
	return result
}

func CreateBlock(txlist []transaction) Block {
	result := Block{transactions: txlist}
	result.Hash = result.ComputeHash()
	return result
}

func CreateTxInput(from *txoutput, key ecdsa.PrivateKey) txinput {
	result := txinput{from: from}
	SignInput(&result, key)
	return result
}

func CreateTxOutput(value int, key ecdsa.PublicKey) txoutput {
	return txoutput{value, key}
}

func (self *Blockchain) Mine() {
	_, oldhead := self.ComputeBlockHeight()
	newhead := Block{Prev: oldhead}
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

type Validatable interface {
	Validate() bool
}

func (self Blockchain) Validate() bool {
	for current := &self.Genesis; current != nil; current = current.Next {
		if current.Prev != nil && current.Prev.Next != current {
			return false
		}

		if !current.Validate() {
			return false
		}
	}

	return true
}

const MAX_TRANSACTIONS_PER_BLOCK = 4096

func (self Block) Validate() bool {
	if len(self.transactions) > MAX_TRANSACTIONS_PER_BLOCK {
		return false
	}

	for _, tx := range self.transactions {
		if !tx.Validate() {
			return false
		}
	}
	return true
}

const MAX_INPUTS_PER_TX = 1024
const MAX_OUTPUTS_PER_TX = 1024

func (self transaction) Validate() bool {
	if len(self.Inputs) > MAX_INPUTS_PER_TX || len(self.Outputs) > MAX_OUTPUTS_PER_TX {
		return false
	}

	sum_inputs := 0
	for _, input := range self.Inputs {
		if !input.Validate() {
			return false
		}
		sum_inputs += input.from.value
	}

	sum_outputs := 0
	for _, output := range self.Outputs {
		if !output.Validate() {
			return false
		}
		sum_outputs += output.value
	}

	if sum_outputs > sum_inputs {
		return false
	}

	// TODO: Sobald ich einen full node implementiert habe, hier prüfen ob Inputs auf UTXOs verweisen

	return true
}
func (self txinput) Validate() bool {
	return self.from != nil && CheckInputSignature(self)
}
func (self txoutput) Validate() bool {
	return self.value >= OUTPUT_MINVALUE && self.value <= OUTPUT_MAXVALUE
}

func (self *Block) ComputeHash() string {
	return fmt.Sprintf("%X", self.ComputeHashByte())
}
func (self *Block) ComputeHashByte() []byte {
	if self != nil {
		input := fmt.Sprintf("block%d", self.nonce)

		for i, tx := range self.transactions {
			input += fmt.Sprintf("%d%X", i, tx.ComputeHashByte())
		}

		if self.Prev != nil {
			return computeSha256(input + self.Prev.Hash)
		} else {
			return computeSha256(input)
		}
	} else {
		return computeSha256("")
	}
}

func (self *transaction) ComputeHash() string { return fmt.Sprintf("%X", self.ComputeHashByte()) }
func (self *transaction) ComputeHashByte() []byte {
	if self != nil {
		hashinput := "tx"
		for _, output := range self.Outputs {
			hashinput += output.ComputeHash()
		}
		for _, input := range self.Inputs {
			hashinput += input.ComputeHash()
		}

		return computeSha256(hashinput)
	} else {
		return computeSha256("")
	}
}

func (self *txinput) ComputeHash() string { return fmt.Sprintf("%X", self.ComputeHashByte()) }
func (self *txinput) ComputeHashByte() []byte {
	if self != nil {
		return computeSha256(fmt.Sprintf("input%X", self.sig.hash))
	} else {
		return computeSha256("")
	}
}

func (self *txoutput) ComputeHash() string { return fmt.Sprintf("%X", self.ComputeHashByte()) }
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
	return fmt.Sprintf("Block(Hash='%s', Genesis=%t, head=%t)",
		self.Hash, self.Prev == nil, self.Next == nil)
}
func (self transaction) String() string {
	return fmt.Sprintf("Transaction[num_outputs=%d,num_inputs=%d]", len(self.Outputs), len(self.Inputs))
}
func (self txinput) String() string {
	return fmt.Sprintf("Input[from=%s]", self.from)
}
func (self txoutput) String() string {
	return fmt.Sprintf("Output[value=%d]", self.value)
}
