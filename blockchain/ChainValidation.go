package blockchain

type Validatable interface {
	Validate() bool
}

func (self Blockchain) Validate() bool {
	for _,current := range self.Blocklist {
		if !current.Validate() {
			return false
		}
	}

	return true
}

const MAX_TRANSACTIONS_PER_BLOCK = 4096

func (self Block) Validate() bool {
	if len(self.Transactions) > MAX_TRANSACTIONS_PER_BLOCK {
		return false
	}

	for _, tx := range self.Transactions {
		if !tx.Validate() {
			return false
		}
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

	// TODO: Sobald ich einen full node implementiert habe, hier prüfen ob Inputs auf UTXOs verweisen

	return true
}
func (self txinput) Validate() bool {
	return self.From != nil && CheckInputSignature(self)
}
func (self txoutput) Validate() bool {
	return self.Value >= OUTPUT_MINVALUE && self.Value <= OUTPUT_MAXVALUE
}
