package blockchain

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
	Transactions []Transaction
}

type Transaction struct {
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
	result := Blockchain{[]Block{CreateBlock(make([]Transaction, 0))}}
	return result
}

func CreateGenesisBlock() Block {
	return CreateBlock(make([]Transaction,0))
}

func CreateBlock(txlist []Transaction) Block {
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
	self.Blocklist = append(self.Blocklist, CreateBlock(make([]Transaction, 0)))
}

func (self *Blockchain) ComputeBlockHeight() int {
	return len(self.Blocklist)
}

func (self *Transaction) ComputePossibleFee() int {
	return Max(0, self.SumOutputs()-self.SumInputs())
}
func (self *Transaction) SumInputs() int {
	result := 0
	for _, input := range self.Inputs {
		result += input.From.Value
	}
	return result
}
func (self *Transaction) SumOutputs() int {
	result := 0
	for _, output := range self.Outputs {
		result += output.Value
	}
	return result
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

func (self *Transaction) ComputeHash() string { return fmt.Sprintf("%X", self.ComputeHashByte()) }
func (self *Transaction) ComputeHashByte() []byte {
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
	return fmt.Sprintf("Block(Hash='%s',Nonce=%d)", self.Hash, self.Nonce)
}
func (self Transaction) String() string {
	return fmt.Sprintf("Transaction[num_outputs=%d,num_inputs=%d]", len(self.Outputs), len(self.Inputs))
}
func (self txinput) String() string {
	return fmt.Sprintf("Input[From=%s]", self.From)
}
func (self txoutput) String() string {
	return fmt.Sprintf("Output[Value=%d]", self.Value)
}
