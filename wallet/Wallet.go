package wallet

import (
	"awesomeProject/blockchain"
	"crypto/ecdsa"
	"fmt"
)

// Txoutput will be nil if key was not used before
var availableUtxoSet = make(map[ecdsa.PublicKey]*blockchain.Txoutput)
var keySet = make(map[ecdsa.PublicKey]ecdsa.PrivateKey)
var pendingUtxoSet = make(map[ecdsa.PublicKey]*blockchain.Txoutput)

func CreateTransaction(receiver ecdsa.PublicKey, value int, fee int, changeAddress ecdsa.PublicKey) *blockchain.Transaction {
	if value+fee > ComputeBalance() {
		return nil
	}

	result := blockchain.Transaction{Message: fmt.Sprintf("Tx value:%d fee:%d", value, fee)}
	for pubkey, utxo := range availableUtxoSet {
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

// Wenn ein Txinput einen in diesem Wallet vorhandenen UTXO verbraucht, wird dieser auf pending gesetzt
func TxHasBeenPublished(txlist []blockchain.Transaction) {
	for _, tx := range txlist {
		for _, in := range tx.Inputs {
			toMove := make([]ecdsa.PublicKey, 0)
			for key, _ := range availableUtxoSet {
				inkey := blockchain.UnmarshalPubkey(in.From.Pubkey)
				if blockchain.PubkeyEqual(inkey, key) {
					toMove = append(toMove, key)
				}
			}
			for _, key := range toMove {
				pendingUtxoSet[key] = availableUtxoSet[key]
				delete(availableUtxoSet, key)
			}
		}
	}
}

// Wenn ein Txinput einen in diesem Wallet vorhandenen UTXO verbraucht, wird dieser gelöscht
// Wenn ein Txoutput einen in diesem Wallet vorhandenen Key verwendet, wird der Txoutput als UTXO hinzugefügt
func TxHasBeenMined(txlist []blockchain.Transaction) {
	for _, tx := range txlist {
		for _, in := range tx.Inputs {
			toDelete := make([]ecdsa.PublicKey, 0)
			for key, _ := range availableUtxoSet {
				inkey := blockchain.UnmarshalPubkey(in.From.Pubkey)
				if blockchain.PubkeyEqual(inkey, key) {
					toDelete = append(toDelete, key)
				}
			}
			for _, key := range toDelete {
				delete(availableUtxoSet, key)
			}
		}

		for _, out := range tx.Outputs {
			for pubkey,_ := range keySet {
				if blockchain.PubkeyEqual(blockchain.UnmarshalPubkey(out.Pubkey),pubkey) {
					availableUtxoSet[pubkey] = &out
				}
			}
		}
	}
}

func ComputeBalance() int {
	result := 0
	for _, utxo := range availableUtxoSet {
		result += utxo.Value
	}
	return result
}
