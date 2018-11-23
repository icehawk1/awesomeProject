package networking

import "net/url"

type Peer struct {
	Address url.URL
	score   int
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
