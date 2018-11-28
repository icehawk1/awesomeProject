package blockchain

import (
	"fmt"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	blockchain()
	merklebaum()

	os.Exit(m.Run())
}

func merklebaum() {
	for i := 0; i < cap(contentlist); i++ {
		contentlist = append(contentlist, createTestTx(fmt.Sprintf("Transaction %d", i)))
	}
	baum = CreateMerklebaum(contentlist)
}

func blockchain() {
	for i := 0; i < cap(utxoList); i++ {
		keylist = append(keylist, CreateKeypair())
		utxoList = append(utxoList, CreateTxOutput(2*i, keylist[i].PublicKey))
	}
}