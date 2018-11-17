package main

import (
	"testing"
)

func TestInitializeChain(t *testing.T) {
	actual := CreateChain("my message")
	if actual.Genesis.Payload != "my message" {
		t.Errorf("Message didn't get through: %s",actual.Genesis.Payload)
	}

	if len(actual.Genesis.Hash) <=0 {
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
		t.Errorf("Block 2 is missing")
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

