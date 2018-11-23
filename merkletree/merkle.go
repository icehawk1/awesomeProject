package merkletree

import (
	"crypto/sha256"
	// "fmt"
)

/* You can add any type of data to the merkle tree
 * as long as it has a method that returns the
 * hash of the object
 */
type Merkable interface {
	Hash() []byte
}

// StringMerkle implements Merkable, with just a string
type StringMerkle struct {
	Data string
}

func (s StringMerkle) Hash() []byte{
	h := sha256.New()
	h.Write([]byte(s.Data))
	return h.Sum(nil)
}

/* Since []StringMerkables does not implement []Merkables using
 * Go, you have to do this converion manually
 */
func ConvertMerkleDataToMerkable(data []StringMerkle) []Merkable{
	merkables := make([]Merkable, len(data))
	for i, v := range data {
		merkables[i] = Merkable(v)
	}
	return merkables
}

type MNode struct {
	Hash       []byte
	LeftChild  *MNode
	RightChild *MNode
}

type MTree struct {
	RootHash	 []byte
	Data         []Merkable
}

func EmptyMTree() *MTree {
	emptyDataList  := []Merkable{}
	emptyBlockData := CreateMTree(emptyDataList)
	return emptyBlockData
}

func (mt *MTree) HasData(m Merkable) bool {
	seenData := mt.Data
	for _, seenDatum := range seenData{
		if string(seenDatum.Hash()) == string(m.Hash()){
			return true
		}
	}
	return false // if there are no matching transactions
}

func CreateMTree(data []Merkable)(*MTree){
	if len(data) == 0 {
		return &MTree{RootHash: []byte{}, Data: data}
	}

	baseTree   := createBaseTree(data)
	resultTree := buildTree(baseTree)
	root       := resultTree[len(resultTree) - 1][0]
	mt         := MTree{RootHash: root.Hash, Data: data}

	return &mt
}

// given a slice of transactions, returns a baseTree for use in buildTree()
func createBaseTree(data []Merkable) (baseTree [][]*MNode){
	var baseLevel []*MNode

	for _ , d := range data {
		h := sha256.New()
		h.Write(d.Hash())
		// h.Write(d.Sig)
		m := &MNode{Hash: h.Sum(nil), LeftChild: nil, RightChild: nil}

		baseLevel = append(baseLevel, m)
	}

	// baseLevel = sortLevel(baseLevel)
	baseTree = append(baseTree, baseLevel)
	return baseTree
}

func addDataToTree(mt MTree, datum Merkable) MTree{
	newData        := append(mt.Data, datum)

	baseTree       := createBaseTree(newData)
	resultTree     := buildTree(baseTree)
	resultTreeRoot := resultTree[len(resultTree) - 1][0]

	finalTree := MTree{RootHash: resultTreeRoot.Hash, Data: newData}

	return finalTree
}

// recursive implementation
func buildTree(inputTree [][]*MNode) (resultTree [][]*MNode) {
	numLevels := len(inputTree)
	highestLevel := inputTree[numLevels - 1]

	if len(highestLevel) > 1 {
		nextLevel := buildNextLevel(highestLevel) //create the net level
		inputTree = append(inputTree, nextLevel) // add to tree

		return buildTree(inputTree) // check if its got a singlular output
	} else { // we have the top level
		return inputTree
	}
}

func buildNextLevel(level []*MNode) (nextLevel []*MNode) {

	if len(level) % 2 != 0{  // uneven number of levels
		level = append(level, level[len(level) - 1])
	}

	for i:=0; i < len(level); i = i + 2 { // interate every two nodes
		nextLevel = append(nextLevel, hashNodes(level[i], level[i+1]))
	}
	return nextLevel
}

func hashNodes(leftChild, rightChild *MNode) (*MNode) {
	h := sha256.New()
	h.Write(leftChild.Hash)
	h.Write(rightChild.Hash)

	parent := MNode{Hash: h.Sum(nil), LeftChild: leftChild, RightChild: rightChild}

	return &parent
}

func (mt MTree) VerifyMTree() bool {
	tree := CreateMTree(mt.Data)
	if string(mt.RootHash) != string(tree.RootHash) {
		return false
	}
	return true
}
