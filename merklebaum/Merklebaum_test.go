package merklebaum

import (
	"awesomeProject/blockchain"
	"fmt"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

type TestContent struct {
	text string
}

func (self TestContent) ComputeHash() string {
	return blockchain.ComputeSha256Hex(self.text)
}
func (self TestContent) ComputeHashByte() []byte {
	return blockchain.ComputeSha256(self.text)
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
	assert.Equal(t, len(actual.GetElements()), 3)
}

func TestGetLeafes(t *testing.T) {
	actual := baum.GetElements()
	assert.Equal(t, len(actual), len(contentlist))
	assert.Equal(t, actual[3].ComputeHash(), blockchain.ComputeSha256("Transaction 3"))
}
