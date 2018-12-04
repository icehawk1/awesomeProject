package networking

import (
	"awesomeProject/blockchain"
	"github.com/emirpasic/gods/sets/treeset"
	"net/url"
)

type Peer struct {
	Address url.URL
	score   int
}
var PeerList = make([]Peer, 0, 5)

func BroadcastBlock(block blockchain.Block) {
	// TODO: Implement
}

func BroadcastTransaction(tx blockchain.Transaction) {
	// TODO: Implement
}

func SelectTransactionsForNextBlock(pendingTx *treeset.Set) []blockchain.Transaction {
	// TODO: Implement
	return nil
}


func CreatePeer(address string) Peer {
	parsed, _ := url.Parse(address)
	result := Peer{*parsed, 100}
	return result
}

func (self Peer) Validate() bool {
	if self.Address.Scheme != "http" && self.Address.Scheme != "https" {
		return false
	}
	return true
}
