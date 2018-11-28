package blockchain

import (
	"awesomeProject/blockchain"
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

type TestContent struct {
	Text string
}

func (self TestContent) ComputeHash() string {
	return blockchain.ComputeSha256Hex(self.Text)
}
func (self TestContent) ComputeHashByte() []byte {
	return blockchain.ComputeSha256(self.Text)
}

var contentlist = make([]blockchain.Hashable, 0, 10)
var baum Merklebaum

func TestMain(m *testing.M) {
	for i := 0; i < cap(contentlist); i++ {
		contentlist = append(contentlist, TestContent{fmt.Sprintf("Transaction %d", i)})
	}
	baum = CreateMerklebaum(contentlist)

	os.Exit(m.Run())
}

func TestIsValid(t *testing.T) {
	assert.True(t, baum.IsValid())
}

func TestContains(t *testing.T) {
	assert.True(t, baum.Contains(TestContent{"Transaction 7"}))
}

func TestHasNode_valid(t *testing.T) {
	proof, ok := baum.CreateSpvProof(TestContent{"Transaction 3"})
	assert.True(t, ok)
	assert.True(t, baum.HasNode(proof))
}

func TestHasNode_invalid(t *testing.T) {
	proof, ok := baum.CreateSpvProof(TestContent{"Transaction 666"})
	assert.True(t, ok)
	assert.False(t, baum.HasNode(proof))
}

func TestCreateLevel(t *testing.T) {
	bäume := []*Merklebaum{&Merklebaum{Hash: "abcd"}, &Merklebaum{Hash: "efgh"}, &Merklebaum{Hash: "ijkl"},
		&Merklebaum{Hash: "mnop"}, &Merklebaum{Hash: "qrst"}}
	assert.Equal(t, len(createMerkleLevel(bäume[:3])), 2)
	assert.Equal(t, len(createMerkleLevel(bäume[:4])), 2)
	assert.Equal(t, len(createMerkleLevel(bäume[:5])), 3)

	assert.Equal(t, createMerkleLevel(bäume)[1].Hash, blockchain.ComputeSha256Hex("0ijklmnop"))
}

func TestCreateMerklebaum_oddlen(t *testing.T) {
	actual := CreateMerklebaum(contentlist[:3])
	assert.Equal(t, 3, len(actual.GetElements()))
}

func TestGetLeafes(t *testing.T) {
	actual := baum.GetElements()
	assert.Equal(t, len(actual), len(contentlist))
	assert.Equal(t, blockchain.ComputeSha256Hex("Transaction 3"), actual[3].ComputeHash())
}

func TestConvertToJson(t *testing.T) {
	encoded, error := baum.toJson()
	assert.NoError(t,error)

	var neueleafs []TestContent
	error = json.Unmarshal(encoded,&neueleafs)
	assert.NoError(t,error)

	var tmp = make([]blockchain.Hashable,0,len(neueleafs))
	for _,l := range neueleafs {
		tmp = append(tmp, l)
	}
	neuerbaum := CreateMerklebaum(tmp)
	
	assert.Equal(t, baum.Hash, neuerbaum.Hash)
	assert.Equal(t, baum.GetElements(),neuerbaum.GetElements())
}
