package blockchain

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

var utxoList = make([]Txoutput, 0, 20)
var keylist = make([]ecdsa.PrivateKey, 0, 20)

func TestValidateTransactionValid(t *testing.T) {
	inputlist := make([]Txinput, 0, 10)
	for i := 0; i < cap(inputlist); i++ {
		inputlist = append(inputlist, CreateTxInput(&utxoList[i], keylist[i]))
	}

	// Demonstrate that there can be more outputs than inputs
	outputlist := make([]Txoutput, 0, 11)
	for value := 0; value < cap(outputlist); value++ {
		outputlist = append(outputlist, CreateTxOutput(value, CreateKeypair().PublicKey))
	}

	tx := Transaction{Outputs: outputlist, Inputs: inputlist}
	assert.True(t, tx.Validate(), fmt.Sprintf("Transaction %s should be valid", tx))
}

func TestValidateTransactionInvalidValue(t *testing.T) {
	inputlist := make([]Txinput, 0, 10)
	for i := 0; i < cap(inputlist); i++ {
		inputlist = append(inputlist, CreateTxInput(&utxoList[i], keylist[i]))
	}

	// Demonstrate that there can be more outputs than inputs
	outputlist := make([]Txoutput, 0, 11)
	for value := 0; value < cap(outputlist); value++ {
		outputlist = append(outputlist, CreateTxOutput(value*2, CreateKeypair().PublicKey))
	}

	tx := Transaction{Outputs: outputlist, Inputs: inputlist}
	assert.False(t, tx.Validate(), fmt.Sprintf("Transaction %s should NOT be valid", tx))
}

func TestValidateTransactionInvalidInput(t *testing.T) {
	inputlist := make([]Txinput, 0, 10)
	for i := 0; i < cap(inputlist)-1; i++ {
		inputlist = append(inputlist, CreateTxInput(&utxoList[i], keylist[i]))
	}
	pubkey := elliptic.Marshal(DefaultCurve, keylist[0].PublicKey.X, keylist[0].PublicKey.Y)
	inputlist = append(inputlist, CreateTxInput(&Txoutput{0, pubkey}, keylist[1]))

	// Demonstrate that there can be more outputs than inputs
	outputlist := make([]Txoutput, 0, 11)
	for value := 0; value < cap(outputlist); value++ {
		outputlist = append(outputlist, CreateTxOutput(value, CreateKeypair().PublicKey))
	}

	tx := Transaction{Outputs: outputlist, Inputs: inputlist}
	assert.False(t, tx.Validate(), fmt.Sprintf("Transaction %s should NOT be valid", tx))
}

func TestMine(t *testing.T) {
	txlist := []Transaction{createTx(1, keylist[1]), createTx(3, keylist[1])}

	genesis := CreateGenesisBlock()
	assert.True(t, genesis.Validate())

	actual := Mine(txlist, genesis.ComputeHash())
	assert.True(t, actual.Validate())
}

func TestParseJsonBlock(t *testing.T) {
	txlist := []Transaction{createTx(5, keylist[1]), createTx(6, keylist[2])}
	genesis := CreateGenesisBlock()
	input := Mine(txlist, genesis.ComputeHash())

	encodedBlock, error := json.Marshal(input)
	assert.NoError(t, error)

	var decodedBlock Block
	error = json.Unmarshal(encodedBlock, &decodedBlock)
	assert.NoError(t, error)
	assert.Equal(t, input,decodedBlock)
	
	assert.True(t,decodedBlock.Validate())
}

func TestParseJsonTx(t *testing.T) {
	inputtx := createTx(321, keylist[2])
	encodedtx, error := json.Marshal(inputtx)
	assert.NoError(t, error)

	var decodedtx Transaction
	error = json.Unmarshal(encodedtx, &decodedtx)
	assert.NoError(t, error)
	assert.Equal(t, inputtx, decodedtx)

	assert.True(t, decodedtx.Validate())
}

func createTx(value int, key ecdsa.PrivateKey) Transaction {
	outputlist := []Txoutput{CreateTxOutput(value+0, key.PublicKey), CreateTxOutput(value+1, key.PublicKey)}
	inputlist := []Txinput{CreateTxInput(&outputlist[0], key), CreateTxInput(&outputlist[1], key)}
	return Transaction{Message: fmt.Sprintf("Tx%d", value), Outputs: outputlist, Inputs: inputlist}
}
