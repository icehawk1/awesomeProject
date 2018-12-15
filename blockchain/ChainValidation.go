package blockchain

import (
	"reflect"
)

type Validatable interface {
	Validate() bool
}

const MAX_TRANSACTIONS_PER_BLOCK = 128

func (self Block) Validate() bool {
	hashIsCorrect := self.Hash == self.ComputeHash()
	difficultsIsMet := BlockhashSatifiesDifficulty(self.Hash)
	numTxIsWithinBounds := len(self.Transactions.GetElements()) <= MAX_TRANSACTIONS_PER_BLOCK
	txAreValid := true
	for _, tx := range self.Transactions.GetElements() {
		if !tx.Validate() {
			txAreValid = false
		}
	}

	return hashIsCorrect && difficultsIsMet && numTxIsWithinBounds && txAreValid
}

const MAX_INPUTS_PER_TX = 1024
const MAX_OUTPUTS_PER_TX = 1024

func (self Transaction) Validate() bool {
	if len(self.Inputs) > MAX_INPUTS_PER_TX || len(self.Outputs) > MAX_OUTPUTS_PER_TX {
		return false
	}

	for _, input := range self.Inputs {
		if !input.Validate() {
			return false
		}

		// Do not accept self referencing
		for _, output := range self.Inputs {
			if reflect.DeepEqual(input.From, output) {
				return false
			}
		}
	}

	for _, output := range self.Outputs {
		if !output.Validate() {
			return false
		}
	}

	if self.SumOutputs() > self.SumInputs() {
		return false
	}

	return true
}
func (self Txinput) Validate() bool {
	return self.From != nil && CheckInputSignature(self)
}
func (self Txoutput) Validate() bool {
	return self.Value >= OUTPUT_MINVALUE && self.Value <= OUTPUT_MAXVALUE
}
