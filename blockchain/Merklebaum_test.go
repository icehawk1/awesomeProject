package blockchain

import (
	"awesomeProject/util"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"testing"
)

var contentlist = make([]Transaction, 0, 10)
var baum Merklebaum

func TestIsValid(t *testing.T) {
	assert.True(t, baum.IsValid())
}

func TestContains(t *testing.T) {
	assert.True(t, baum.Contains(createTestTx("Transaction 7")))
}

func TestHasNode_valid(t *testing.T) {
	proof, ok := baum.CreateSpvProof(createTestTx("Transaction 3"))
	assert.True(t, ok)
	assert.True(t, baum.HasNode(proof))
}

func TestHasNode_invalid(t *testing.T) {
	proof, ok := baum.CreateSpvProof(createTestTx("Transaction 666"))
	assert.True(t, ok)
	assert.False(t, baum.HasNode(proof))
}

func TestCreateLevel(t *testing.T) {
	bäume := []*Merklebaum{&Merklebaum{Hash: "abcd"}, &Merklebaum{Hash: "efgh"}, &Merklebaum{Hash: "ijkl"},
		&Merklebaum{Hash: "mnop"}, &Merklebaum{Hash: "qrst"}}
	assert.Equal(t, len(createMerkleLevel(bäume[:3])), 2)
	assert.Equal(t, len(createMerkleLevel(bäume[:4])), 2)
	assert.Equal(t, len(createMerkleLevel(bäume[:5])), 3)

	assert.Equal(t, createMerkleLevel(bäume)[1].Hash, util.ComputeSha256Hex("0ijklmnop"))
}

func TestCreateMerklebaum_oddlen(t *testing.T) {
	actual := CreateMerklebaum(contentlist[:3])
	assert.Equal(t, 3, len(actual.GetElements()))
}

func TestGetLeafes(t *testing.T) {
	actual := baum.GetElements()
	assert.Equal(t, len(actual), len(contentlist))
	hashprefix := "tx"
	assert.Equal(t, util.ComputeSha256Hex(hashprefix+"Transaction 3"), actual[3].ComputeHash())
}

func TestConvertToJson(t *testing.T) {
	encoded, error := baum.toJson()
	assert.NoError(t,error)

	var neueleafs []Transaction
	error = json.Unmarshal(encoded,&neueleafs)
	assert.NoError(t,error)

	neuerbaum := CreateMerklebaum(neueleafs)
	
	assert.Equal(t, baum.Hash, neuerbaum.Hash)
	assert.Equal(t, baum.GetElements(),neuerbaum.GetElements())
}

func createTestTx(msg string) Transaction {
	return Transaction{Message:msg, Outputs:[]Txoutput{},Inputs:[]Txinput{}}
}