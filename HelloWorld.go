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
	fmt.Println(chain)

	bytearr, _ := json.Marshal(chain)
	encoded := fmt.Sprintf("%s",bytearr)
	fmt.Println("chain: ",encoded)
}

