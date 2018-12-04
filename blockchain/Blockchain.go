package blockchain

import (
	"awesomeProject/util"
	"crypto/ecdsa"
	"crypto/elliptic"
	"fmt"
	"github.com/cbergoon/merkletree"
	"math"
	"math/rand"
	"strings"
	"time"
)

type Block struct {
	Hash         string
	Prev         string
	Nonce        uint64
	Transactions Merklebaum
}

type Transaction struct {
	Message string
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
	Pubkey []byte
}

var prng = rand.New(rand.NewSource(time.Now().UnixNano()))

func CreateGenesisBlock() Block {
	return Mine(make([]Transaction, 0), "")
}

func CreateBlock(txlist []Transaction, prevhash string) Block {
	var result Block
	if len(txlist) > 0 {
		contentlist := make([]Transaction, len(txlist))
		for i, elem := range txlist {
			contentlist[i] = elem
		}
		if len(txlist) != len(contentlist) {
			panic("Something went horribly wrong")
		}

		tree := CreateMerklebaum(contentlist)
		result = Block{Transactions: tree, Prev: prevhash, Nonce: prng.Uint64()}
	} else {
		result = Block{Prev: prevhash, Nonce: prng.Uint64()}
	}

	result.Hash = result.ComputeHash()
	return result
}

func CreateTxInput(from *Txoutput, key ecdsa.PrivateKey) Txinput {
	result := Txinput{From: from}
	SignInput(&result, key)
	return result
}

func CreateTxOutput(value int, key ecdsa.PublicKey) Txoutput {
	return Txoutput{value, elliptic.Marshal(DefaultCurve,key.X,key.Y)}
}

const Difficulty = 1

func Mine(txlist []Transaction, prevhash string) Block {
	for {
		newblock, valid := MineAttempt(txlist,prevhash)
		if valid {
			return newblock
		}
	}
}

func MineAttempt(txlist []Transaction, prevhash string) (Block, bool) {
	requiredPrefix := strings.Repeat("0", Difficulty)
	newblock := CreateBlock(txlist, prevhash)
	return newblock, strings.HasPrefix(newblock.Hash, requiredPrefix)
}

func ComputeBlockHeight(head Block, knownBlocks *map[string]Block) int {
	i := 0
	var ok bool
	for ; head.Prev != ""; i++ {
		head, ok = (*knownBlocks)[head.Prev]
		if (!ok) {
			return -1
		}
	}
	return i
}
func (self *Block) ComputeHash() string {
	return fmt.Sprintf("%X", self.ComputeHashByte())
}
func (self *Block) ComputeHashByte() []byte {
	if self != nil {
		var roothash string
		if self.Transactions.Hash != "" {
			roothash = self.Transactions.Hash
		} else {
			roothash = util.ComputeSha256Hex("")
		}

		input := fmt.Sprintf("block%d%s%s", self.Nonce, roothash, self.Prev)
		return util.ComputeSha256(input)
	} else {
		return util.ComputeSha256("")
	}
}

func AddFees(transactions []Transaction, keypair ecdsa.PrivateKey) []Transaction {
	// TODO: Implement
	return transactions
}

func CreateRandomTransaction(utxo map[string]Txoutput, keypair ecdsa.PrivateKey) *Transaction {
	// TODO: Implement
	return nil
}

var blockreward int = 12

func CreateCoinbaseTransaction(pubkey ecdsa.PublicKey) Transaction {
	return Transaction{Outputs:[]Txoutput{CreateTxOutput(blockreward,pubkey)}}
}

func (self *Block) GetTransactions() []Transaction {
	if self != nil {
		return self.Transactions.GetElements()
	} else {
		return []Transaction{}
	}
}

func (self *Transaction) ComputePossibleFee() int {
	return util.Max(0, self.SumOutputs()-self.SumInputs())
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
func (self Transaction) ComputeHash() string { return fmt.Sprintf("%X", self.ComputeHashByte()) }
func (self Transaction) ComputeHashByte() []byte {
	hashinput := "tx"+self.Message

	for _, output := range self.Outputs {
		hashinput += output.ComputeHash()
	}
	for _, input := range self.Inputs {
		hashinput += input.ComputeHash()
	}

	return util.ComputeSha256(hashinput)
}
func (self Transaction) CalculateHash() ([]byte, error) { return self.ComputeHashByte(), nil }
func (self Transaction) Equals(other merkletree.Content) (bool, error) {
	othertx, ok := other.(Transaction)
	if ok {
		return self.ComputeHash() == othertx.ComputeHash(), nil
	} else {
		return false, nil
	}
}

func (self Txinput) ComputeHash() string { return fmt.Sprintf("%X", self.ComputeHashByte()) }
func (self Txinput) ComputeHashByte() []byte {
	return util.ComputeSha256(fmt.Sprintf("input%X%X", self.Sig.R.Bytes(),self.Sig.S.Bytes()))
}

func (self Txoutput) ComputeHash() string { return fmt.Sprintf("%X", self.ComputeHashByte()) }
func (self Txoutput) ComputeHashByte() []byte {
	return util.ComputeSha256(fmt.Sprintf("output%d%s", self.Value, self.Pubkey))
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
