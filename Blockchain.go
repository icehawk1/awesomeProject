package main

import (
	"crypto/ecdsa"
	"fmt"
	"math"
)

type Blockchain struct {
	Blocklist []Block
}

type Block struct {
	Hash         string
	Nonce        int64
	Transactions []transaction
}

type transaction struct {
	Outputs []txoutput
	Inputs  []txinput
}

type txinput struct {
	From *txoutput
	Sig  Signature
}

const OUTPUT_MINVALUE = 0
const OUTPUT_MAXVALUE = math.MaxInt32

type txoutput struct {
	Value  int
	Pubkey ecdsa.PublicKey
}

func CreateChain() Blockchain {
	result := Blockchain{[]Block{CreateBlock(make([]transaction, 0))}}
	return result
}

func CreateBlock(txlist []transaction) Block {
	result := Block{Transactions: txlist}
	result.Hash = result.ComputeHash(nil)
	return result
}

func CreateTxInput(from *txoutput, key ecdsa.PrivateKey) txinput {
	result := txinput{From: from}
	SignInput(&result, key)
	return result
}

func CreateTxOutput(value int, key ecdsa.PublicKey) txoutput {
	return txoutput{value, key}
}

func (self *Blockchain) Mine() {
	self.Blocklist = append(self.Blocklist, CreateBlock(make([]transaction,0)))
}

func (self *Blockchain) ComputeBlockHeight() int {
	return len(self.Blocklist)
}

type Validatable interface {
	Validate() bool
}

func (self Blockchain) Validate() bool {
	for _,current := range self.Blocklist {
		if !current.Validate() {
			return false
		}
	}

	return true
}

const MAX_TRANSACTIONS_PER_BLOCK = 4096

func (self Block) Validate() bool {
	if len(self.Transactions) > MAX_TRANSACTIONS_PER_BLOCK {
		return false
	}

	for _, tx := range self.Transactions {
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
		sum_inputs += input.From.Value
	}

	sum_outputs := 0
	for _, output := range self.Outputs {
		if !output.Validate() {
			return false
		}
		sum_outputs += output.Value
	}

	if sum_outputs > sum_inputs {
		return false
	}

	// TODO: Sobald ich einen full node implementiert habe, hier prüfen ob Inputs auf UTXOs verweisen

	return true
}
func (self txinput) Validate() bool {
	return self.From != nil && CheckInputSignature(self)
}
func (self txoutput) Validate() bool {
	return self.Value >= OUTPUT_MINVALUE && self.Value <= OUTPUT_MAXVALUE
}

func (self *Block) ComputeHash(prev *Block) string {
	return fmt.Sprintf("%X", self.ComputeHashByte(prev))
}
func (self *Block) ComputeHashByte(prev *Block) []byte {
	if self != nil {
		input := fmt.Sprintf("block%d", self.Nonce)

		for i, tx := range self.Transactions {
			input += fmt.Sprintf("%d%X", i, tx.ComputeHashByte())
		}

		if prev != nil {
			return computeSha256(input + prev.Hash)
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
		return computeSha256(fmt.Sprintf("input%X", self.Sig.hash))
	} else {
		return computeSha256("")
	}
}

func (self *txoutput) ComputeHash() string { return fmt.Sprintf("%X", self.ComputeHashByte()) }
func (self *txoutput) ComputeHashByte() []byte {
	if self != nil {
		return computeSha256(fmt.Sprintf("output%d%s", self.Value, self.Pubkey))
	} else {
		return computeSha256("")
	}
}

func (self Blockchain) String() string {
	return fmt.Sprintf("Chain Genesis: %s", self.Blocklist[0])
}
func (self Block) String() string {
	return fmt.Sprintf("Block(Hash='%s',Nonce=%d)", self.Hash,self.Nonce)
}
func (self transaction) String() string {
	return fmt.Sprintf("Transaction[num_outputs=%d,num_inputs=%d]", len(self.Outputs), len(self.Inputs))
}
func (self txinput) String() string {
	return fmt.Sprintf("Input[From=%s]", self.From)
}
func (self txoutput) String() string {
	return fmt.Sprintf("Output[Value=%d]", self.Value)
}
