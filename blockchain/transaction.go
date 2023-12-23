package blockchain

type Transaction struct {
	ID      []byte
	Inputs  []TxInput
	Outputs []TxOutput
}

type TxInput struct {
	Value  int
	PubKey string
}

type TxOutput struct {
	ID  []byte
	Out int
	Sig string
}
