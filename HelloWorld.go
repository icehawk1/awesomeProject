package main

import (
	"encoding/json"
	"fmt"
)

func main() {
	chain := Blockchain{}
	chain.Mine()
	chain.Mine()
	chain.Mine()

	bytearr, _ := json.Marshal(chain)
	encoded := fmt.Sprintf("%s",bytearr)
	fmt.Println("chain: ",encoded)

	var secondChain Blockchain
	json.Unmarshal(bytearr, secondChain)
	fmt.Println(secondChain)
}
