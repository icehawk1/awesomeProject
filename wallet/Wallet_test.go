package wallet

import (
	"awesomeProject/blockchain"
	"crypto/ecdsa"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

var keys = make([]ecdsa.PrivateKey, 0, 10)

func TestMain(m *testing.M) {
	for i:= 0; i<15; i++ {
		var key = blockchain.CreateKeypair()
		keySet[key.PublicKey] = key
		keys = append(keys, key)
	}

	for i:=0; i<10; i++ {
		utxo := blockchain.CreateTxOutput(i, keys[i].PublicKey)
		availableUtxoSet[keys[i].PublicKey] = &utxo
	}

	os.Exit(m.Run())
}

func TestComputeBalance(t *testing.T) {
	actual := ComputeBalance()
	assert.Equal(t, 45, actual)
}

func TestCreateTransaction(t *testing.T) {
	recKey := blockchain.CreateKeypair()
	changeKey := blockchain.CreateKeypair()
	receiver := blockchain.MarshalPubkey(recKey.PublicKey)
	changeAddress := blockchain.MarshalPubkey(changeKey.PublicKey)
	value := 11
	fee := 12
	actual := CreateTransaction(recKey.PublicKey, value, fee, changeKey.PublicKey)

	assert.True(t, actual.SumInputs() >= (value + fee), "Zu wenig Geld gesendet")
	possibleFee := actual.ComputePossibleFee()
	assert.Equal(t, fee, possibleFee)
	assert.Equal(t, actual.SumOutputsForAddr(receiver), value)
	assert.Equal(t, actual.SumOutputsForAddr(receiver)+actual.SumOutputsForAddr(changeAddress), actual.SumOutputs())
}

func TestTxHasBeenMined(t *testing.T) {
	tmp := blockchain.CreateTxOutput(1, keys[1].PublicKey)
	inputs := []blockchain.Txinput{blockchain.CreateTxInput(&tmp, keys[1])}
	outputs := []blockchain.Txoutput{blockchain.CreateTxOutput(2, keys[11].PublicKey)}

	tx := blockchain.Transaction{Message: "TestTxHasBeenMined", Outputs: outputs, Inputs:inputs}
	TxHasBeenMined([]blockchain.Transaction{tx})

	_, hasKey := availableUtxoSet[keys[1].PublicKey]
	assert.False(t, hasKey)
	_, hasKey = availableUtxoSet[keys[11].PublicKey]
	assert.True(t, hasKey)
}

func TestTxHasBeenPublished(t *testing.T) {
	tmp := blockchain.CreateTxOutput(1, keys[4].PublicKey)
	inputs := []blockchain.Txinput{blockchain.CreateTxInput(&tmp, keys[4])}
	tx := blockchain.Transaction{Message: "TestTxHasBeenMined", Inputs:inputs}

	TxHasBeenPublished([]blockchain.Transaction{tx})

	_, hasKey := availableUtxoSet[keys[4].PublicKey]
	assert.False(t, hasKey)
	_, hasKey = pendingUtxoSet[keys[4].PublicKey]
	assert.True(t, hasKey)
}
