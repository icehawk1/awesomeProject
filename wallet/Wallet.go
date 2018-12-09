package wallet

import (
	"awesomeProject/blockchain"
	"crypto/ecdsa"
	"fmt"
)

// Txoutput will be nil if key was not used before
var utxoSet = make(map[ecdsa.PublicKey]*blockchain.Txoutput)
var keySet = make(map[ecdsa.PublicKey]ecdsa.PrivateKey)
var pendingUtxoSet = make(map[ecdsa.PublicKey]*blockchain.Txoutput)

func CreateTransaction(receiver ecdsa.PublicKey, value int, fee int, changeAddress ecdsa.PublicKey) *blockchain.Transaction {
	if value+fee > ComputeBalance() {
		return nil
	}

	result := blockchain.Transaction{Message: fmt.Sprintf("Tx value:%d fee:%d", value, fee)}
	for pubkey, utxo := range utxoSet {
		in := blockchain.CreateTxInput(utxo, keySet[pubkey])
		result.Inputs = append(result.Inputs, in)

		if result.SumInputs() >= value+fee {
			break
		}
	}

	result.Outputs = append(result.Outputs, blockchain.CreateTxOutput(value, receiver))
	if (result.SumInputs() > value+fee) {
		change := blockchain.CreateTxOutput(result.SumInputs()-value-fee, changeAddress)
		result.Outputs = append(result.Outputs, change)
	}

	return &result
}

func TxHasBeenPublished(txlist []blockchain.Transaction) {

}

func TxHasBeenMined(txlist []blockchain.Transaction) {

}

func ComputeBalance() int {
	result := 0
	for _, utxo := range utxoSet {
		result += utxo.Value
	}
	return result
}
