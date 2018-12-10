package networking

import (
	"awesomeProject/blockchain"
	"bytes"
	"encoding/json"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type Peer struct {
	Address url.URL
	score   int
}
var PeerList = make([]Peer, 0, 5)
var httpclient = &http.Client{
	Timeout: time.Second * 2,
}

func BroadcastBlock(block blockchain.Block) []int {
	// TODO: Implement
	return nil
}

func BroadcastTransaction(tx blockchain.Transaction) []int {
	result := make([]int,0,len(PeerList))
	for _,peer := range PeerList {
		url := strings.TrimRight(peer.Address.String(),"/") + "/pending_transaction"
		statusCode := funcName(tx, url)
		result = append(result, statusCode)
	}

	return result
}

func funcName(tx blockchain.Transaction, url string) int {
	buf := new(bytes.Buffer)
	json.NewEncoder(buf).Encode(tx)

	response, err := httpclient.Post(url, "application/json", buf)
	if err != nil {
		return -1
	}
	defer response.Body.Close()
	return response.StatusCode
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
