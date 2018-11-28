package blockchain

import (
	"fmt"
	. "gopkg.in/check.v1"
	"testing"
)

var contentlist = make([]Transaction, 0, 10)
var baum Merklebaum

// Hook up gocheck into the "go test" runner.
func Test(t *testing.T) { TestingT(t) }
type MySuite struct{}
var _ = Suite(&MySuite{})

func (s *MySuite) SetUpTest(c *C) {
	for i := 0; i < cap(contentlist); i++ {
		contentlist = append(contentlist, createTestTx(fmt.Sprintf("Transaction %d", i)))
	}
	baum = CreateMerklebaum(contentlist)
}

func (s *MySuite) TestIsValid(c *C) {
	c.Assert(baum.IsValid(),Equals,true)
}

func (s *MySuite) TestContains(c *C) {
	c.Assert(baum.Contains(createTestTx("Transaction 7")) ,Equals,true)
}

func (s *MySuite) TestHasNode_valid(c *C) {
	proof, ok := baum.CreateSpvProof(createTestTx("Transaction 3"))
	c.Assert(ok,Equals,true)
	c.Assert(baum.HasNode(proof), Equals,true)
}

func (s *MySuite) TestHasNode_invalid(c *C) {
	proof, ok := baum.CreateSpvProof(createTestTx("Transaction 666"))
	c.Assert(ok, Equals,true)
	c.Assert(baum.HasNode(proof), Equals, false)
}

func (s *MySuite) TestCreateLevel(c *C) {
	bäume := []*Merklebaum{{Hash: "abcd"}, {Hash: "efgh"}, {Hash: "ijkl"}, {Hash: "mnop"}, {Hash: "qrst"}}
	c.Assert( len(createMerkleLevel(bäume[:3])),Equals, 2)
	c.Assert(len(createMerkleLevel(bäume[:4])),Equals, 2)
	c.Assert(len(createMerkleLevel(bäume[:5])),Equals, 3)

	c.Assert(createMerkleLevel(bäume)[1].Hash,Equals, ComputeSha256Hex("0ijklmnop"))
}

func (s *MySuite) TestCreateMerklebaum_oddlen(c *C) {
	actual := CreateMerklebaum(contentlist[:3])
	c.Assert( len(actual.GetElements()),Equals, 3)
}

func (s *MySuite) TestGetLeafes(c *C) {
	actual := baum.GetElements()
	c.Assert( len(actual), Equals, len(contentlist))
	c.Assert( ComputeSha256Hex("Transaction 3"), Equals, actual[3].ComputeHash())
}

/*
func (s *MySuite) TestConvertToJson(c *C) {
	encoded, error := baum.toJson()
	c.Assert(error,IsNil)

	var neueleafs []Transaction
	error = json.Unmarshal(encoded,&neueleafs)
	c.Assert(error,IsNil)

	var tmp = make([]Transaction,0,len(neueleafs))
	for _,l := range neueleafs {
		tmp = append(tmp, l)
	}
	neuerbaum := CreateMerklebaum(tmp)

	c.Assert( baum.Hash, DeepEquals, neuerbaum.Hash)
	c.Assert(baum.GetElements(), DeepEquals,neuerbaum.GetElements())
}
*/

func createTestTx(msg string) Transaction {
	return Transaction{Message:msg, Outputs:[]Txoutput{},Inputs:[]Txinput{}}
}