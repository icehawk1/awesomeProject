package blockchain

import "testing"

func TestGenerateHash(t *testing.T) {
	actual := computeSha256Hex("abcdef")
	expected := "BEF57EC7F53A6D40BEB640A780A639C83BC29AC8A9816F1FC6C5C6DCD93C4721"
	if actual != expected {
		t.Errorf("Got wrong ComputeHash: %s",actual)
	}

	actual = computeSha256Hex("a")
	expected = "CA978112CA1BBDCAFAC231B39A23DC4DA786EFF8147C4E72B9807785AFEE48BB"
	if actual != expected {
		t.Errorf("Got wrong ComputeHash: %s",actual)
	}
}

func TestCreateMerkleTree(t *testing.T) {

}