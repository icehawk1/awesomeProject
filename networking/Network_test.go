package networking

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCreatePeer(t *testing.T) {
	validPeer1 := CreatePeer("http://heise.de/")
	assert.True(t,validPeer1.Validate())

	validPeer2 := CreatePeer("https://heise:7654.de")
	assert.True(t,validPeer2.Validate())

	invalidPeer := CreatePeer("ftp://heise.de/")
	assert.False(t, invalidPeer.Validate())
}
