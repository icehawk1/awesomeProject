package blockchain

import "strings"

type Validatable interface {
	Validate() bool
}

const MAX_TRANSACTIONS_PER_BLOCK = 4096

func (self Block) Validate() bool {
	if len(self.Transactions.GetElements()) > MAX_TRANSACTIONS_PER_BLOCK {
		return false
	}

	for _, tx := range self.Transactions.GetElements() {
		if !tx.Validate() {
			return false
		}
	}

	if !strings.HasPrefix(self.ComputeHash(), strings.Repeat("0", Difficulty)) {
		return false
	}

	return true
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
