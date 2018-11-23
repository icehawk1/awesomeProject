package blockchain

import (
	"awesomeProject/merkletree"
	"crypto/ecdsa"
	"fmt"
	"math"
	"strings"
)

type Block struct {
	Hash         string
	prev         string
	Nonce        int64
	Transactions merkletree.MTree
}

type Transaction struct {
	Outputs []Txoutput
	Inputs  []Txinput
}

type Txinput struct {
	From *Txoutput
	Sig  Signature
}

const OUTPUT_MINVALUE = 0
const OUTPUT_MAXVALUE = math.MaxInt32

type Txoutput struct {
	Value  int
	Pubkey ecdsa.PublicKey
}

func CreateGenesisBlock() Block {
	return CreateBlock(make([]Transaction, 0), "")
}

func CreateBlock(txlist []Transaction, prevhash string) Block {
	var tree merkletree.MTree = *merkletree.CreateMTree(ConvertTxToMerkable(txlist))
	result := Block{Transactions: tree, prev: prevhash}
	result.Hash = result.ComputeHash()
	return result
}

func CreateTxInput(from *Txoutput, key ecdsa.PrivateKey) Txinput {
	result := Txinput{From: from}
	SignInput(&result, key)
	return result
}

func CreateTxOutput(value int, key ecdsa.PublicKey) Txoutput {
	return Txoutput{value, key}
}

const Difficulty = 1

func Mine(txlist []Transaction, prevhash string) Block {
	requiredPrefix := strings.Repeat("0", Difficulty)

	for {
		newblock := CreateBlock(txlist, prevhash)
		if strings.HasPrefix(newblock.Hash, requiredPrefix) {
			return newblock
		}
	}
}

func ComputeBlockHeight(head Block, knownBlocks *map[string]Block) int {
	i := 0
	var ok bool
	for ; head.prev != ""; i++ {
		head, ok = (*knownBlocks)[head.prev]
		if (!ok) {
			return -1
		}
	}
	return i
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

func (self *Block) ComputeHash() string {
	return fmt.Sprintf("%X", self.ComputeHashByte())
}
func (self *Block) ComputeHashByte() []byte {
	if self != nil {
		input := fmt.Sprintf("block%d%X%s", self.Nonce, self.Transactions.RootHash, self.prev)
		return computeSha256(input)
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
func (self Transaction) Hash() []byte { return self.ComputeHashByte() }

func (self *Txinput) ComputeHash() string { return fmt.Sprintf("%X", self.ComputeHashByte()) }
func (self *Txinput) ComputeHashByte() []byte {
	if self != nil {
		return computeSha256(fmt.Sprintf("input%X", self.Sig.hash))
	} else {
		return computeSha256("")
	}
}

func (self *Txoutput) ComputeHash() string { return fmt.Sprintf("%X", self.ComputeHashByte()) }
func (self *Txoutput) ComputeHashByte() []byte {
	if self != nil {
		return computeSha256(fmt.Sprintf("output%d%s", self.Value, self.Pubkey))
	} else {
		return computeSha256("")
	}
}

/* Since []Transaction does not implement []Merkables using
 * Go, you have to do this converion manually
 */
func ConvertTxToMerkable(data []Transaction) []merkletree.Merkable{
	merkables := make([]merkletree.Merkable, len(data))
	for i, v := range data {
		merkables[i] = merkletree.Merkable(v)
	}
	return merkables
}

func (self Block) String() string {
	return fmt.Sprintf("Block(ComputeHash='%s',Nonce=%d)", self.Hash, self.Nonce)
}
func (self Transaction) String() string {
	return fmt.Sprintf("Transaction[num_outputs=%d,num_inputs=%d]", len(self.Outputs), len(self.Inputs))
}
func (self Txinput) String() string {
	return fmt.Sprintf("Input[From=%s]", self.From)
}
func (self Txoutput) String() string {
	return fmt.Sprintf("Output[Value=%d]", self.Value)
}
