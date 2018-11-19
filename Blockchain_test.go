package main

import (
	"crypto/ecdsa"
	"fmt"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

var utxoList = make([]txoutput,0,20)
var keylist = make([]ecdsa.PrivateKey,0,20)

func TestMain(m *testing.M) {
	for i:=0; i<cap(utxoList); i++ {
		keylist = append(keylist, CreateKeypair())
		utxoList = append(utxoList, CreateTxOutput(2*i,keylist[i].PublicKey))
	}
	fmt.Println(utxoList)

	os.Exit(m.Run())
}

func TestInitializeChain(t *testing.T) {
	actual := CreateChain()

	if len(actual.Blocklist[0].Hash) <=0 {
		t.Errorf("Wrong Hash: %s",actual.Blocklist[0].Hash)
	}
}

func TestMineTwoBlocks(t *testing.T) {
	actual := CreateChain()
	actual.Mine()
	actual.Mine()

	assert.Equal(t,len(actual.Blocklist),3,"Blockchain has incorrect length")
}

func TestValidateTransactionValid(t *testing.T) {
	inputlist := make([]txinput,0,10)
	for i:=0; i<cap(inputlist); i++ {
		inputlist = append(inputlist, CreateTxInput(&utxoList[i],keylist[i]))
	}

	// Demonstrate that there can be more outputs than inputs
	outputlist := make([]txoutput,0,11)
	for value :=0; value<cap(outputlist); value++ {
		outputlist = append(outputlist, CreateTxOutput(value, CreateKeypair().PublicKey))
	}

	tx := transaction{Outputs: outputlist,Inputs:inputlist}
	assert.True(t,tx.Validate(),fmt.Sprintf("Transaction %s should be valid",tx))
}

func TestValidateTransactionInvalidValue(t *testing.T) {
	inputlist := make([]txinput,0,10)
	for i:=0; i<cap(inputlist); i++ {
		inputlist = append(inputlist, CreateTxInput(&utxoList[i],keylist[i]))
	}

	// Demonstrate that there can be more outputs than inputs
	outputlist := make([]txoutput,0,11)
	for value :=0; value<cap(outputlist); value++ {
		outputlist = append(outputlist, CreateTxOutput(value*2, CreateKeypair().PublicKey))
	}

	tx := transaction{Outputs: outputlist,Inputs:inputlist}
	assert.False(t,tx.Validate(),fmt.Sprintf("Transaction %s should NOT be valid",tx))
}

func TestValidateTransactionInvalidInput(t *testing.T) {
	inputlist := make([]txinput,0,10)
	for i:=0; i<cap(inputlist)-1; i++ {
		inputlist = append(inputlist, CreateTxInput(&utxoList[i],keylist[i]))
	}
	inputlist= append(inputlist, CreateTxInput(&txoutput{0,keylist[0].PublicKey},keylist[1]))

	// Demonstrate that there can be more outputs than inputs
	outputlist := make([]txoutput,0,11)
	for value :=0; value<cap(outputlist); value++ {
		outputlist = append(outputlist, CreateTxOutput(value, CreateKeypair().PublicKey))
	}

	tx := transaction{Outputs: outputlist,Inputs:inputlist}
	assert.False(t,tx.Validate(),fmt.Sprintf("Transaction %s should NOT be valid",tx))
}

func TestValidateChain(t *testing.T) {
	actual := CreateChain()
	actual.Mine()
	actual.Mine()
	assert.True(t,actual.Validate(),"Chain should be valid")
}