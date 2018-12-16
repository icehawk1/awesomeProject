package networking

import (
	"awesomeProject/blockchain"
	"awesomeProject/util"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type Peer struct {
	Address url.URL
	score   int
}

const MAX_PEERS = 5

var PeerList = make([]Peer, 0, MAX_PEERS)
var SelfAddr Peer
var httpclient = &http.Client{
	Timeout: time.Second * 2,
}

func CreatePeer(address string) *Peer {
	parsed, err := url.Parse(strings.TrimRight(address, "/"))
	if err==nil {
		result := Peer{*parsed, 50}
		return &result
	} else {
		return nil
	}
}

// Etabliert ein simples P2P-Netzwerk bei dem nicht antwortende oder überlastete Peers langsam herausgefiltert werden
// Stabil laufende Peers werden bevorzugt kontaktiert
func ContactPeer(i int) bool {
	if i < 0 || i > len(PeerList) {
		return false
	}

	peer := PeerList[i]
	addr := peer.Address.String() + "/peers?url=" + url.QueryEscape(SelfAddr.Address.String())
	response, err := http.Get(addr)
	if err != nil {
		peer.score = util.Min(0, peer.score-5)
		return false
	}
	defer response.Body.Close()

	if response.StatusCode >= 200 && response.StatusCode <= 299 {
		peer.score = util.Max(100, peer.score+1)
		return true
	} else {
		peer.score = util.Min(0, peer.score-1)
		return false
	}
}

func FillPeerList(peer string) {
	if peer=="" {return}

	addr := strings.TrimRight(peer, "/") + "/peers?url=" + url.QueryEscape(SelfAddr.Address.String())
	response, err := http.Get(addr)
	if err != nil {
		return
	}
	defer response.Body.Close()

	if response.StatusCode >= 200 && response.StatusCode <= 299 {
		created := CreatePeer(peer)
		if created != nil {AddPeer(*created)}

		var remotePeers []Peer
		json.NewDecoder(response.Body).Decode(&remotePeers)
		if remotePeers != nil {
			for _,rpeer := range remotePeers {
				AddPeer(rpeer)
			}
		}
	}
}

// Fügt einen Peer zur PeerList hinzu, wenn dadurch ein Peer mit schlechterem Score ersetzt wird
func AddPeer(newPeer Peer) bool {
	if newPeer.Validate() && !KnownPeer(newPeer) {
		if len(PeerList) < MAX_PEERS {
			PeerList = append(PeerList, newPeer)
			return true
		} else {
			for i, existingPeer := range PeerList {
				if existingPeer.score < newPeer.score {
					PeerList[i] = newPeer
					return true
				}
			}
			return false
		}
	}
	return false
}

// true falls der Peer bereits in der Peerlist ist
func KnownPeer(peer Peer) bool {
	for _, elem := range PeerList {
		if peer.Address == elem.Address {
			return true
		}
	}

	return peer.Address==SelfAddr.Address
}

func BroadcastBlock(block blockchain.Block) []int {
	result := make([]int, 0, len(PeerList))
	for _, peer := range PeerList {
		addr := peer.Address.String() + "/block/"

		buf := new(bytes.Buffer)
		json.NewEncoder(buf).Encode(block)

		response, err := httpclient.Post(addr, "application/json", buf)
		if err != nil {
			result = append(result, -1)
			continue
		}
		defer response.Body.Close()

		result = append(result, response.StatusCode)
	}

	return result
}

func BroadcastTransaction(tx blockchain.Transaction) []int {
	result := make([]int, 0, len(PeerList))
	for _, peer := range PeerList {
		addr := peer.Address.String() + "/pending_transaction"

		buf := new(bytes.Buffer)
		json.NewEncoder(buf).Encode(tx)

		response, err := httpclient.Post(addr, "application/json", buf)
		if err != nil {
			result = append(result, -1)
		} else {
			defer response.Body.Close()
			result = append(result, response.StatusCode)
		}
	}

	return result
}

func (self Peer) Validate() bool {
	if self.Address.Scheme != "http" && self.Address.Scheme != "https" {
		return false
	}

	if self.score < 0 || self.score > 100 {
		return false
	}

	return true
}

func (self Peer) String() string {
	return fmt.Sprintf("Peer{Address='%s', score=%d}", self.Address.String(), self.score)
}
