package wallet

import (
	"awesomeProject/blockchain"
	"crypto/elliptic"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	for i:=0; i<10; i++ {
		var key = blockchain.CreateKeypair()
		utxo := blockchain.CreateTxOutput(i, key.PublicKey)
		utxoSet[key.PublicKey] = &utxo
		keySet[key.PublicKey] = key
	}

	os.Exit(m.Run())
}

func TestComputeBalance(t *testing.T) {
	actual := ComputeBalance()
	assert.Equal(t,45, actual)
}

func TestCreateTransaction(t *testing.T) {
	recKey := blockchain.CreateKeypair()
	changeKey := blockchain.CreateKeypair()
	receiver := elliptic.Marshal(blockchain.DefaultCurve,recKey.X,recKey.Y)
	changeAddress := elliptic.Marshal(blockchain.DefaultCurve,changeKey.X,changeKey.Y)
	fee := 12
	value := 11
	actual := CreateTransaction(recKey.PublicKey, value, fee, changeKey.PublicKey)

	assert.True(t, actual.SumInputs()>=(value+fee), "Zu wenig Geld gesendet")
	possibleFee := actual.ComputePossibleFee()
	assert.Equal(t, fee, possibleFee)
	assert.Equal(t, actual.SumOutputsForAddr(receiver),value)
	assert.Equal(t, actual.SumOutputsForAddr(receiver)+actual.SumOutputsForAddr(changeAddress), actual.SumOutputs())
}

func TestTxHasBeenPublished(t *testing.T) {

}

func TestTxHasBeenMined(t *testing.T) {

}
