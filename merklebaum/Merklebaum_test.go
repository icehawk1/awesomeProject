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

var contentlist = make([]blockchain.Hashable,0,10)
var baum Merklebaum
func TestMain(m *testing.M) {
	for i:=0; i<len(contentlist); i++ {
		contentlist = append(contentlist, TestContent{fmt.Sprintf("Transaction %d",i)})
	}
	baum = CreateMerklebaum(contentlist)

	os.Exit(m.Run())
}

func TestIsValid(t *testing.T) {
	assert.True(t, baum.IsValid())
}

func TestHasNode_valid(t *testing.T) {
	proof, ok := baum.CreateSpvProof(TestContent{"Transaction 3"})
	assert.True(t,ok)
	assert.True(t, baum.HasNode(proof))
}

func TestHasNode_invalid(t *testing.T) {
	proof,ok := baum.CreateSpvProof(TestContent{"Transaction 666"})
	assert.True(t,ok)
	assert.True(t, baum.HasNode(proof))
}

func TestGetLeafes(t *testing.T) {
	actual := baum.GetLeafes()
	assert.Equal(t, len(actual), len(contentlist))
	assert.Equal(t, actual[3].ComputeHash(), blockchain.ComputeSha256("Transaction 3"))
}