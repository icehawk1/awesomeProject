package main

import (
	"strings"
	"testing"
)

func TestInitializeChain(t *testing.T) {
	actual := CreateChain("my message")
	if actual.Genesis.Payload != "my message" {
		t.Errorf("Message didn't get through: %s",actual.Genesis.Payload)
	}

	if actual.Genesis.Hash != "EA38E30F75767D7E6C21EBA85B14016646A3B60ADE426CA966DAC940A5DB1BAB" {
		t.Errorf("Wrong Hash: %s",actual.Genesis.Hash)
	}

	if actual.Genesis.Prev != nil {
		t.Errorf("Genesis hat Vorgänger")
	}
}

func TestMineTwoBlocks(t *testing.T) {
	actual := CreateChain("andere nachricht")
	actual.Mine()
	actual.Mine()

	if actual.Genesis.Next == nil {
		t.Errorf("block 2 is missing")
	}

	if actual.Genesis.Next.Prev != &actual.Genesis {
		t.Errorf("Man zeigt mit einem angezogenen Pointer auf einen nackten Block!")
	}

	if actual.Genesis.Next.Next == nil {
		t.Errorf("Block 3 is missing")
	}

	if actual.Genesis.Next.Next.Next != nil {
		t.Errorf("Thats too many blocks")
	}
}

func TestJsonEncode(t *testing.T) {
	chain := CreateChain("andere nachricht")
	chain.Mine()
	chain.Mine()
	chain.Mine()

	actual := toJson(chain)
	if !strings.Contains(actual, "block3") {
		t.Errorf("block3 is missing")
	}
}